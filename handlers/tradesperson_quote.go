package handlers

import (
	"database/sql"
	"log"
	"math"
	"net/http"
	"sort"
	"strings"

	"redbudway-api/database"
	"redbudway-api/email"
	"redbudway-api/internal"
	"redbudway-api/models"
	"redbudway-api/restapi/operations"

	"github.com/go-openapi/runtime/middleware"
	"github.com/stripe/stripe-go/v82"
	"github.com/stripe/stripe-go/v82/customer"
	"github.com/stripe/stripe-go/v82/invoice"
	"github.com/stripe/stripe-go/v82/product"
	"github.com/stripe/stripe-go/v82/quote"
)

// productSnap is a comparable line-item snapshot (Stripe or request) for change detection.
type productSnap struct {
	title string
	price int64
	qty   int64
}

func clampDepositPct64(v int64) int64 {
	if v < 0 {
		return 0
	}
	if v > 100 {
		return 100
	}
	return v
}

func getBillingQuoteDepositPct(tradespersonID, stripeQuoteID string) (int64, error) {
	var dp sql.NullInt64
	err := database.GetConnection().QueryRow(
		`SELECT q.depositPct FROM tradesperson_quotes tq INNER JOIN quotes q ON tq.quoteId=q.id WHERE tq.tradespersonId=? AND tq.quote=?`,
		tradespersonID, stripeQuoteID,
	).Scan(&dp)
	if err == sql.ErrNoRows {
		return 0, nil
	}
	if err != nil {
		return 0, err
	}
	if !dp.Valid {
		return 0, nil
	}
	return dp.Int64, nil
}

func loadStripeQuoteProductSnaps(quoteID, connectAccountID string) ([]productSnap, error) {
	var out []productSnap
	p := &stripe.QuoteListLineItemsParams{Quote: stripe.String(quoteID)}
	if connectAccountID != "" {
		p.SetStripeAccount(connectAccountID)
	}
	it := quote.ListLineItems(p)
	for it.Next() {
		li := it.LineItem()
		if li == nil || li.Price == nil || li.Price.Product == nil || li.Price.Product.ID == "" {
			continue
		}
		pp := &stripe.ProductParams{}
		if connectAccountID != "" {
			pp.SetStripeAccount(connectAccountID)
		}
		sp, err := product.Get(li.Price.Product.ID, pp)
		if err != nil {
			return nil, err
		}
		qty := li.Quantity
		if qty < 1 {
			qty = 1
		}
		ua := li.Price.UnitAmount
		out = append(out, productSnap{title: strings.TrimSpace(sp.Name), price: ua, qty: qty})
	}
	if err := it.Err(); err != nil {
		return nil, err
	}
	return out, nil
}

func productSnapsFromModels(products []*models.Product) []productSnap {
	var out []productSnap
	for _, p := range products {
		if p == nil {
			continue
		}
		qty := p.Quantity
		if qty < 1 {
			qty = 1
		}
		out = append(out, productSnap{title: strings.TrimSpace(p.Title), price: p.Price, qty: qty})
	}
	return out
}

func equalSortedProductSnaps(a, b []productSnap) bool {
	if len(a) != len(b) {
		return false
	}
	key := func(s []productSnap) {
		sort.Slice(s, func(i, j int) bool {
			if s[i].title != s[j].title {
				return s[i].title < s[j].title
			}
			if s[i].price != s[j].price {
				return s[i].price < s[j].price
			}
			return s[i].qty < s[j].qty
		})
	}
	aa := append([]productSnap(nil), a...)
	bb := append([]productSnap(nil), b...)
	key(aa)
	key(bb)
	for i := range aa {
		if aa[i].title != bb[i].title || aa[i].price != bb[i].price || aa[i].qty != bb[i].qty {
			return false
		}
	}
	return true
}

func GetTradespersonTradespersonIDBillingQuotesHandler(params operations.GetTradespersonTradespersonIDBillingQuotesParams, principal interface{}) middleware.Responder {
	tradespersonID := params.TradespersonID
	quarter := params.Quarter
	year := params.Year
	page := *params.Page
	token := params.HTTPRequest.Header.Get("Authorization")

	response := operations.NewGetTradespersonTradespersonIDBillingQuotesOK()
	quotes := []*operations.GetTradespersonTradespersonIDBillingQuotesOKBodyItems0{}

	valid, err := ValidateTradespersonAccessToken(tradespersonID, token)
	if err != nil {
		log.Printf("Failed to validate tradesperson %s, accessToken %s", tradespersonID, token)
		return response
	} else if !valid {
		log.Printf("Bad actor tradesperson %s, accessToken %s", tradespersonID, token)
		return response
	}

	tpStripeID, err := database.GetTradespersonStripeID(tradespersonID)
	if err != nil || tpStripeID == "" {
		log.Printf("billing quotes GetTradespersonStripeID: %v", err)
		return response
	}

	db := database.GetConnection()

	stmt, err := db.Prepare("SELECT quote FROM tradesperson_quotes WHERE tradespersonId=? AND QUARTER(created) = ? AND YEAR(created) = ? GROUP BY id ORDER BY created DESC LIMIT ?, ?")
	if err != nil {
		log.Printf("Failed to create prepare statement, %v", err)
		return response
	}
	defer stmt.Close()

	offSet := (page - 1) * int64(PAGE_SIZE)
	rows, err := stmt.Query(tradespersonID, quarter, year, offSet, PAGE_SIZE)
	if err != nil {
		log.Printf("Failed to execute select statement %s", err)
		return response
	}

	var quoteID string
	for rows.Next() {
		if err := rows.Scan(&quoteID); err != nil {
			log.Printf("Failed to scan row %v", err)
			return response
		}
		stripeQuote, err := quote.Get(quoteID, billingQuoteParams(tpStripeID, nil))
		if err != nil {
			log.Printf("Failed to get stripe quote, %v", err)
			return response
		}

		_quote := &operations.GetTradespersonTradespersonIDBillingQuotesOKBodyItems0{}
		_quote.Status = string(stripeQuote.Status)
		_quote.Number = stripeQuote.Number
		_quote.QuoteID = quoteID
		if stripeQuote.Invoice != nil {
			_quote.InvoiceID = stripeQuote.Invoice.ID
		}

		cuP := &stripe.CustomerParams{}
		cuP.SetStripeAccount(tpStripeID)
		stripeCustomer, err := customer.Get(stripeQuote.Customer.ID, cuP)
		if err != nil {
			log.Printf("Failed to get stripe customer, %v", err)
			return response
		}

		if stripeCustomer.Deleted {
			if stripeQuote.Status == "accepted" {
				stripeInvoice, err := loadStripeInvoiceForBillingQuote(tpStripeID, stripeQuote.Invoice.ID)
				if err != nil {
					log.Printf("Failed to get stripe invoice with ID %v, %v", stripeQuote.Invoice.ID, err)
					return response
				}
				_quote.Customer = stripeInvoice.CustomerName
				if stripeInvoice.CustomerEmail != "" {
					_quote.CustomerEmail = stripeInvoice.CustomerEmail
				}
			}
		} else {
			_quote.Customer = stripeCustomer.Name
			if stripeCustomer.Email != "" {
				_quote.CustomerEmail = stripeCustomer.Email
			}
		}

		quotes = append(quotes, _quote)

	}
	response.SetPayload(quotes)
	return response
}

func GetTradespersonTradespersonIDBillingQuoteQuoteIDHandler(params operations.GetTradespersonTradespersonIDBillingQuoteQuoteIDParams, principal interface{}) middleware.Responder {
	tradespersonID := params.TradespersonID
	quoteID := params.QuoteID
	token := params.HTTPRequest.Header.Get("Authorization")

	response := operations.NewGetTradespersonTradespersonIDBillingQuoteQuoteIDOK()
	_quote := models.QuoteDetails{}

	valid, err := ValidateTradespersonAccessToken(tradespersonID, token)
	if err != nil {
		log.Printf("Failed to validate tradesperson %s, accessToken %s", tradespersonID, token)
		return response
	} else if !valid {
		log.Printf("Bad actor tradesperson %s, accessToken %s", tradespersonID, token)
		return response
	}

	tpStripeID, err := database.GetTradespersonStripeID(tradespersonID)
	if err != nil || tpStripeID == "" {
		log.Printf("quote detail GetTradespersonStripeID %s: %v", tradespersonID, err)
		return errMap(http.StatusBadGateway, "provider billing account unavailable")
	}

	db := database.GetConnection()

	stmt, err := db.Prepare("SELECT tq.request, q.title, q.description, tq.customerId, q.depositPct FROM tradesperson_quotes tq INNER JOIN quotes q ON tq.quoteId=q.id WHERE tq.tradespersonId=? AND tq.quote=?")
	if err != nil {
		log.Printf("Failed to create prepared statement, %v", err)
		return response
	}
	defer stmt.Close()

	var message, title, description, customerID string
	var depositPct int64
	row := stmt.QueryRow(tradespersonID, quoteID)
	switch err = row.Scan(&message, &title, &description, &customerID, &depositPct); err {
	case sql.ErrNoRows:
		log.Printf("Tradesperson with ID %s has no quote %s", tradespersonID, quoteID)
	case nil:

		stripeQuote, err := quote.Get(quoteID, billingQuoteParams(tpStripeID, nil))
		if err != nil || stripeQuote == nil {
			log.Printf("Failed to get stripe quote %s: %v", quoteID, err)
			return errMap(http.StatusBadGateway, "unable to load quote from billing")
		}

		_quote.Request = message
		_quote.Created = stripeQuote.Created
		_quote.Status = string(stripeQuote.Status)
		_quote.Number = stripeQuote.Number
		_quote.Description = stripeQuote.Description
		_quote.Expires = stripeQuote.ExpiresAt
		if stripeQuote.Status == stripe.QuoteStatusAccepted && stripeQuote.Invoice != nil {
			_quote.InvoiceID = stripeQuote.Invoice.ID
		}

		service := &models.QuoteDetailsService{}
		service.Title = title
		service.Description = description
		service.DepositPct = depositPct
		products := []*models.Product{}

		params := &stripe.QuoteListLineItemsParams{Quote: stripe.String(quoteID)}
		params.SetStripeAccount(tpStripeID)
		i := quote.ListLineItems(params)
		for i.Next() {
			lineItem := i.LineItem()
			pp := &stripe.ProductParams{}
			pp.SetStripeAccount(tpStripeID)
			stripeProduct, err := product.Get(lineItem.Price.Product.ID, pp)
			if err != nil {
				log.Printf("Failed to get stripe product, %v", err)
			}
			_product := &models.Product{}
			_product.Title = stripeProduct.Name
			_product.Price = lineItem.Price.UnitAmount
			_product.Quantity = lineItem.Quantity
			products = append(products, _product)
		}
		service.Products = products
		_quote.Service = service

		if stripeQuote.Customer == nil || stripeQuote.Customer.ID == "" {
			log.Printf("stripe quote %s has no customer id", quoteID)
			return errMap(http.StatusBadGateway, "quote has no customer in billing")
		}

		cuP := &stripe.CustomerParams{}
		cuP.SetStripeAccount(tpStripeID)
		stripeCustomer, err := customer.Get(stripeQuote.Customer.ID, cuP)
		if err != nil || stripeCustomer == nil {
			log.Printf("Failed to get stripe customer %s: %v", stripeQuote.Customer.ID, err)
			return errMap(http.StatusBadGateway, "unable to load customer from billing")
		}

		_customer := &models.Customer{}
		if stripeCustomer.Deleted {
			if stripeQuote.Status == stripe.QuoteStatusAccepted && stripeQuote.Invoice != nil {
				stripeInvoice, err := loadStripeInvoiceForBillingQuote(tpStripeID, stripeQuote.Invoice.ID)
				if err != nil {
					log.Printf("Failed to get stripe invoice with ID %v, %v", stripeQuote.Invoice.ID, err)
					return errMap(http.StatusBadGateway, "unable to load invoice from billing")
				}
				_customer.Name = stripeInvoice.CustomerName
				_customer.Email = stripeInvoice.CustomerEmail
				_customer.Phone = stripeInvoice.CustomerPhone
				// Invoice snapshot may include billing address; Customer object often has none for deleted test users.
				_customer.Address = quoteDetailAddressFromStripe(stripeInvoice.CustomerAddress)
			}
		} else {
			_customer.Name = stripeCustomer.Name
			_customer.Email = stripeCustomer.Email
			_customer.Phone = stripeCustomer.Phone
			// Do not map Customer.address: checkout-created cus_ objects usually have no address; dashboard only needs contact fields.
		}

		_quote.Images, err = internal.GetQuoteImages(_customer.Email, stripeQuote.ID)
		if err != nil {
			log.Printf("Failed to get quote email images, %v", err)
		}

		_quote.Customer = _customer
	default:
		log.Printf("Unknown, %v", err)
	}
	response.SetPayload(&_quote)
	return response
}

func GetTradespersonTradespersonIDBillingQuoteQuoteIDPdfHandler(params operations.GetTradespersonTradespersonIDBillingQuoteQuoteIDPdfParams, principal interface{}) middleware.Responder {
	tradespersonID := params.TradespersonID
	quoteID := params.QuoteID
	token := params.HTTPRequest.Header.Get("Authorization")

	response := operations.NewGetTradespersonTradespersonIDBillingQuoteQuoteIDPdfOK()

	valid, err := ValidateTradespersonAccessToken(tradespersonID, token)
	if err != nil {
		log.Printf("Failed to validate tradesperson %s, accessToken %s", tradespersonID, token)
		return response
	} else if !valid {
		log.Printf("Bad actor tradesperson %s, accessToken %s", tradespersonID, token)
		return response
	}

	tpStripeID, err := database.GetTradespersonStripeID(tradespersonID)
	if err != nil || tpStripeID == "" {
		log.Printf("quote PDF GetTradespersonStripeID %s: %v", tradespersonID, err)
		return response
	}

	db := database.GetConnection()

	stmt, err := db.Prepare("SELECT tq.quote FROM tradesperson_quotes tq INNER JOIN quotes q ON tq.quoteId=q.id WHERE tq.tradespersonId=? AND tq.quote=?")
	if err != nil {
		log.Printf("Failed to create prepared statement, %v", err)
		return response
	}
	defer stmt.Close()

	row := stmt.QueryRow(tradespersonID, quoteID)
	switch err = row.Scan(&quoteID); err {
	case sql.ErrNoRows:
		log.Printf("Tradesperson with ID %s has no quote %s", tradespersonID, quoteID)
	case nil:

		params := &stripe.QuotePDFParams{}
		params.SetStripeAccount(tpStripeID)
		resp, err := quote.PDF(quoteID, params)
		if err != nil {
			log.Printf("Failed to get stripe quote, %v", err)
			return response
		}
		response.SetPayload(resp.LastResponse.Body)
	default:
		log.Printf("Unknown, %v", err)
	}
	return response
}

func PutTradespersonTradespersonIDBillingQuoteQuoteIDHandler(params operations.PutTradespersonTradespersonIDBillingQuoteQuoteIDParams, principal interface{}) middleware.Responder {
	tradespersonID := params.TradespersonID
	quoteID := params.QuoteID
	description := params.Quote.Description
	products := params.Quote.Products
	token := params.HTTPRequest.Header.Get("Authorization")

	response := operations.NewPutTradespersonTradespersonIDBillingQuoteQuoteIDOK()
	payload := operations.PutTradespersonTradespersonIDBillingQuoteQuoteIDOKBody{}
	payload.Updated = false
	response.SetPayload(&payload)

	valid, err := ValidateTradespersonAccessToken(tradespersonID, token)
	if err != nil {
		log.Printf("Failed to validate tradesperson %s, accessToken %s", tradespersonID, token)
		return response
	} else if !valid {
		log.Printf("Bad actor tradesperson %s, accessToken %s", tradespersonID, token)
		return response
	}

	tpStripeID, err := database.GetTradespersonStripeID(tradespersonID)
	if err != nil || tpStripeID == "" {
		log.Printf("Put GetTradespersonStripeID %s: %v", tradespersonID, err)
		return response
	}

	sq, err := quote.Get(quoteID, billingQuoteParams(tpStripeID, nil))
	if err != nil {
		log.Printf("Put quote.Get %s: %v", quoteID, err)
		return response
	}
	if sq.Status == stripe.QuoteStatusAccepted {
		return errMap(http.StatusConflict, "quote is accepted and cannot be edited")
	}

	origQuoteID := quoteID

	currentDeposit, err := getBillingQuoteDepositPct(tradespersonID, origQuoteID)
	if err != nil {
		log.Printf("Put getBillingQuoteDepositPct %s: %v", origQuoteID, err)
		return response
	}
	incomingDesc := strings.TrimSpace(description)
	existingDesc := strings.TrimSpace(sq.Description)
	incomingDeposit := currentDeposit
	if params.Quote.DepositPct != nil {
		incomingDeposit = clampDepositPct64(*params.Quote.DepositPct)
	}
	descChanged := incomingDesc != existingDesc
	depositNeedsWrite := incomingDeposit != currentDeposit

	var linesUnchanged bool
	if len(products) > 0 {
		existingSnaps, snapErr := loadStripeQuoteProductSnaps(origQuoteID, tpStripeID)
		if snapErr != nil {
			log.Printf("Put loadStripeQuoteProductSnaps %s: %v", origQuoteID, snapErr)
			linesUnchanged = false
		} else {
			linesUnchanged = equalSortedProductSnaps(existingSnaps, productSnapsFromModels(products))
		}
	} else {
		linesUnchanged = true
	}

	// Draft (or non-open) with no line items in the request and nothing to change — avoid calling Stripe.
	if len(products) == 0 && !descChanged && !depositNeedsWrite && sq.Status != stripe.QuoteStatusOpen {
		payload.Updated = true
		response.SetPayload(&payload)
		return response
	}

	// Stripe does not allow line_items on quote.Update when status is `open`.
	// Description-only updates are allowed; line-item changes require creating a revision draft, then finalize.
	if sq.Status == stripe.QuoteStatusOpen && len(products) == 0 {
		if !descChanged && !depositNeedsWrite {
			payload.Updated = true
			response.SetPayload(&payload)
			return response
		}
		if descChanged {
			descParams := &stripe.QuoteParams{
				Description: stripe.String(description),
			}
			descParams.SetStripeAccount(tpStripeID)
			_, err = quote.Update(quoteID, descParams)
			if err != nil {
				log.Printf("Failed to update open quote description %s: %v", quoteID, err)
				return response
			}
		}
		if depositNeedsWrite {
			db := database.GetConnection()
			_, err := db.Exec(
				`UPDATE quotes q INNER JOIN tradesperson_quotes tq ON q.id = tq.quoteId SET q.depositPct = ? WHERE tq.tradespersonId = ? AND tq.quote = ?`,
				incomingDeposit, tradespersonID, origQuoteID,
			)
			if err != nil {
				log.Printf("Failed to update depositPct for quote %s: %v", origQuoteID, err)
				return response
			}
		}
		payload.Updated = true
		response.SetPayload(&payload)
		return response
	}

	// Nothing changed (including priced lines): idempotent success, no Stripe churn or customer email.
	if len(products) > 0 && linesUnchanged && !descChanged && !depositNeedsWrite {
		payload.Updated = true
		response.SetPayload(&payload)
		return response
	}

	// Open quote, line items unchanged: update Stripe description and/or DB deposit only — no revision, no email.
	if sq.Status == stripe.QuoteStatusOpen && len(products) > 0 && linesUnchanged {
		if descChanged {
			descParams := &stripe.QuoteParams{
				Description: stripe.String(description),
			}
			descParams.SetStripeAccount(tpStripeID)
			_, err = quote.Update(quoteID, descParams)
			if err != nil {
				log.Printf("Failed to update open quote description %s: %v", quoteID, err)
				return response
			}
		}
		if depositNeedsWrite {
			db := database.GetConnection()
			_, err := db.Exec(
				`UPDATE quotes q INNER JOIN tradesperson_quotes tq ON q.id = tq.quoteId SET q.depositPct = ? WHERE tq.tradespersonId = ? AND tq.quote = ?`,
				incomingDeposit, tradespersonID, origQuoteID,
			)
			if err != nil {
				log.Printf("Failed to update depositPct for quote %s: %v", origQuoteID, err)
				return response
			}
		}
		payload.Updated = true
		response.SetPayload(&payload)
		return response
	}

	lineItems := []*stripe.QuoteLineItemParams{}
	for _, _product := range products {
		prodParams := &stripe.ProductParams{
			Name: stripe.String(_product.Title),
		}
		prodParams.SetStripeAccount(tpStripeID)
		stripeProduct, err := product.New(prodParams)
		if err != nil {
			log.Printf("Failed to get create stripe product, %v", err)
			return response
		}
		qty := _product.Quantity
		if qty < 1 {
			qty = 1
		}
		lineItem := &stripe.QuoteLineItemParams{
			PriceData: &stripe.QuoteLineItemPriceDataParams{
				Currency:   stripe.String("USD"),
				Product:    stripe.String(stripeProduct.ID),
				UnitAmount: &_product.Price,
			},
			Quantity: stripe.Int64(qty),
		}
		lineItems = append(lineItems, lineItem)
	}

	effectiveID := origQuoteID
	revisedFromOpen := false
	// Revise only when the quote is open and priced lines actually changed vs Stripe.
	if sq.Status == stripe.QuoteStatusOpen && len(products) > 0 && !linesUnchanged {
		revParams := &stripe.QuoteParams{
			FromQuote: &stripe.QuoteFromQuoteParams{
				Quote:      stripe.String(origQuoteID),
				IsRevision: stripe.Bool(true),
			},
		}
		revParams.SetStripeAccount(tpStripeID)
		revQ, err := quote.New(revParams)
		if err != nil {
			log.Printf("Put quote revision from %s: %v", origQuoteID, err)
			return response
		}
		if revQ.Status != stripe.QuoteStatusDraft {
			log.Printf("Put quote revision: expected draft, got %s", revQ.Status)
			return response
		}
		effectiveID = revQ.ID
		revisedFromOpen = true
	}

	quoteParams := &stripe.QuoteParams{
		Description: stripe.String(description),
	}
	quoteParams.SetStripeAccount(tpStripeID)
	if len(products) != 0 {
		quoteParams.LineItems = lineItems
	}

	_, err = quote.Update(effectiveID, quoteParams)
	if err != nil {
		log.Printf("Failed to update quote %s: %v", effectiveID, err)
		return response
	}

	if revisedFromOpen {
		finalizeParams := &stripe.QuoteFinalizeQuoteParams{}
		finalizeParams.SetStripeAccount(tpStripeID)
		stripeQuote, err := quote.FinalizeQuote(effectiveID, finalizeParams)
		if err != nil {
			log.Printf("Put quote FinalizeQuote %s: %v", effectiveID, err)
			return response
		}
		if stripeQuote.Status != stripe.QuoteStatusOpen {
			log.Printf("Put quote FinalizeQuote %s: expected open, got %s", effectiveID, stripeQuote.Status)
			return response
		}
		updated, err := database.UpdateQuote(tradespersonID, effectiveID, origQuoteID)
		if err != nil || !updated {
			log.Printf("Put quote UpdateQuote DB %s -> %s: updated=%v err=%v", origQuoteID, effectiveID, updated, err)
			return response
		}
		sqMail, qerr := quote.Get(effectiveID, billingQuoteParams(tpStripeID, nil))
		if qerr != nil {
			log.Printf("Put quote quote.Get after finalize %s: %v", effectiveID, qerr)
			sqMail = stripeQuote
		}
		if _, err := sendFinalizedQuoteReadyEmail(tradespersonID, effectiveID, sqMail); err != nil {
			log.Printf("Put quote sendFinalizedQuoteReadyEmail %s: %v", effectiveID, err)
		}
		payload.QuoteID = effectiveID
	}

	depositStripeID := origQuoteID
	if revisedFromOpen {
		depositStripeID = effectiveID
	}
	if depositNeedsWrite {
		db := database.GetConnection()
		_, err := db.Exec(
			`UPDATE quotes q INNER JOIN tradesperson_quotes tq ON q.id = tq.quoteId SET q.depositPct = ? WHERE tq.tradespersonId = ? AND tq.quote = ?`,
			incomingDeposit, tradespersonID, depositStripeID,
		)
		if err != nil {
			log.Printf("Failed to update depositPct for quote %s: %v", depositStripeID, err)
			return response
		}
	}

	payload.Updated = true
	response.SetPayload(&payload)

	return response
}

func PostTradespersonTradespersonIDBillingQuoteQuoteIDCancelHandler(params operations.PostTradespersonTradespersonIDBillingQuoteQuoteIDCancelParams, principal interface{}) middleware.Responder {
	tradespersonID := params.TradespersonID
	quoteID := params.QuoteID
	token := params.HTTPRequest.Header.Get("Authorization")

	response := operations.NewPostTradespersonTradespersonIDBillingQuoteQuoteIDCancelOK()
	payload := operations.PostTradespersonTradespersonIDBillingQuoteQuoteIDCancelOKBody{}

	valid, err := ValidateTradespersonAccessToken(tradespersonID, token)
	if err != nil {
		log.Printf("Failed to validate tradesperson %s, accessToken %s", tradespersonID, token)
		return response
	} else if !valid {
		log.Printf("Bad actor tradesperson %s, accessToken %s", tradespersonID, token)
		return response
	}

	tpStripeID, err := database.GetTradespersonStripeID(tradespersonID)
	if err != nil || tpStripeID == "" {
		log.Printf("cancel quote GetTradespersonStripeID %s: %v", tradespersonID, err)
		return response
	}

	db := database.GetConnection()

	stmt, err := db.Prepare("SELECT q.title, tq.request FROM tradesperson_quotes tq INNER JOIN quotes q ON q.id=tq.quoteId WHERE tq.tradespersonId=? AND tq.quote=?")
	if err != nil {
		log.Printf("Failed to create prepared statement, %v", err)
		return response
	}
	defer stmt.Close()

	var title, message string
	row := stmt.QueryRow(tradespersonID, quoteID)
	switch err = row.Scan(&title, &message); err {
	case sql.ErrNoRows:
		log.Printf("Tradesperson with ID %s has no quote %s", tradespersonID, quoteID)
	case nil:
		cancelParams := &stripe.QuoteCancelParams{}
		cancelParams.SetStripeAccount(tpStripeID)
		stripeQuote, err := quote.Cancel(quoteID, cancelParams)
		if err != nil {
			log.Printf("Failed to cancel quote %s, %v", quoteID, err)
		}
		if stripeQuote.Status == "canceled" {

			// if stripeQuote.Invoice == nil {
			// 	_, err := database.DeleteQuote(tradespersonID, quoteID)
			// 	if err != nil {
			// 		log.Printf("Failed to delete tradesperson quote, %v", err)
			// 		return response
			// 	}
			// }
			payload.Canceled = true
			response.SetPayload(&payload)

			tradesperson, err := database.GetTradespersonProfile(tradespersonID)
			if err != nil {
				log.Printf("Failed to get tradesperson profile %s", err)
				return response
			}
			cuP := &stripe.CustomerParams{}
			cuP.SetStripeAccount(tpStripeID)
			stripeCustomer, err := customer.Get(stripeQuote.Customer.ID, cuP)
			if err != nil {
				log.Printf("Failed to get stripe customer, %v", err)
				return response
			}
			if err := email.SendCustomerQuoteCancellation(tradesperson, stripeCustomer, message, title); err != nil {
				log.Printf("Failed to send customer email, %v", err)
			}
		}
	default:
		log.Printf("Unknown, %v", err)
	}
	return response
}

func PostTradespersonTradespersonIDBillingQuoteQuoteIDFinalizeHandler(params operations.PostTradespersonTradespersonIDBillingQuoteQuoteIDFinalizeParams, principal interface{}) middleware.Responder {
	tradespersonID := params.TradespersonID
	quoteID := params.QuoteID
	token := params.HTTPRequest.Header.Get("Authorization")

	payload := operations.PostTradespersonTradespersonIDBillingQuoteQuoteIDFinalizeOKBody{Finalized: false}
	response := operations.NewPostTradespersonTradespersonIDBillingQuoteQuoteIDFinalizeOK().WithPayload(&payload)

	valid, err := ValidateTradespersonAccessToken(tradespersonID, token)
	if err != nil {
		log.Printf("Failed to validate tradesperson %s, accessToken %s", tradespersonID, token)
		return response
	} else if !valid {
		log.Printf("Bad actor tradesperson %s, accessToken %s", tradespersonID, token)
		return response
	}

	tpStripeID, err := database.GetTradespersonStripeID(tradespersonID)
	if err != nil || tpStripeID == "" {
		log.Printf("finalize quote GetTradespersonStripeID %s: %v", tradespersonID, err)
		return response
	}

	db := database.GetConnection()

	stmt, err := db.Prepare("SELECT request FROM tradesperson_quotes WHERE tradespersonId=? AND quote=?")
	if err != nil {
		log.Printf("Failed to create prepared statement, %v", err)
		return response
	}
	defer stmt.Close()

	var message string
	row := stmt.QueryRow(tradespersonID, quoteID)
	switch err = row.Scan(&message); err {
	case sql.ErrNoRows:
		log.Printf("Tradesperson with ID %s has no quote %s", tradespersonID, quoteID)
	case nil:
		finalizeParams := &stripe.QuoteFinalizeQuoteParams{}
		finalizeParams.SetStripeAccount(tpStripeID)
		stripeQuote, err := quote.FinalizeQuote(quoteID, finalizeParams)
		if err != nil {
			log.Printf("Failed to finalize quote %s, %v", quoteID, err)
		}
		if stripeQuote.Status == "open" {
			payload.Finalized = true
			response.SetPayload(&payload)

			sq, qerr := quote.Get(quoteID, billingQuoteParams(tpStripeID, nil))
			if qerr != nil {
				log.Printf("Failed to refresh quote %s after finalize: %v", quoteID, qerr)
				sq = stripeQuote
			}

			_, err = sendFinalizedQuoteReadyEmail(tradespersonID, quoteID, sq)
			if err != nil {
				log.Printf("Failed to send customer finalize email, %v", err)
			}
		}
	default:
		log.Printf("Unknown, %v", err)
	}
	return response
}

func PostTradespersonTradespersonIDBillingQuoteQuoteIDReviseHandler(params operations.PostTradespersonTradespersonIDBillingQuoteQuoteIDReviseParams, principal interface{}) middleware.Responder {
	tradespersonID := params.TradespersonID
	quoteID := params.QuoteID
	token := params.HTTPRequest.Header.Get("Authorization")

	payload := operations.PostTradespersonTradespersonIDBillingQuoteQuoteIDReviseOKBody{Revised: false}
	response := operations.NewPostTradespersonTradespersonIDBillingQuoteQuoteIDReviseOK().WithPayload(&payload)

	valid, err := ValidateTradespersonAccessToken(tradespersonID, token)
	if err != nil {
		log.Printf("Failed to validate tradesperson %s, accessToken %s", tradespersonID, token)
		return response
	} else if !valid {
		log.Printf("Bad actor tradesperson %s, accessToken %s", tradespersonID, token)
		return response
	}

	tpStripeID, err := database.GetTradespersonStripeID(tradespersonID)
	if err != nil || tpStripeID == "" {
		log.Printf("revise quote GetTradespersonStripeID %s: %v", tradespersonID, err)
		return response
	}

	db := database.GetConnection()

	stmt, err := db.Prepare("SELECT tradespersonId FROM tradesperson_quotes WHERE tradespersonId=? AND quote=?")
	if err != nil {
		log.Printf("Failed to create prepared statement, %v", err)
		return response
	}
	defer stmt.Close()

	row := stmt.QueryRow(tradespersonID, quoteID)
	switch err = row.Scan(&tradespersonID); err {
	case sql.ErrNoRows:
		log.Printf("Tradesperson with ID %s has no quote %s", tradespersonID, quoteID)
	case nil:
		params := &stripe.QuoteParams{
			FromQuote: &stripe.QuoteFromQuoteParams{
				Quote:      &quoteID,
				IsRevision: stripe.Bool(true),
			},
		}
		params.SetStripeAccount(tpStripeID)
		stripeQuote, err := quote.New(params)
		if err != nil {
			log.Printf("Failed to revise stripe quote, %v", err)
			return response
		}
		if stripeQuote.Status == "draft" {
			updated, err := database.UpdateQuote(tradespersonID, stripeQuote.ID, quoteID)
			if err != nil {
				log.Printf("Failed to update quote %v", err)
			}
			payload.Revised = updated
			payload.QuoteID = stripeQuote.ID
			response.SetPayload(&payload)
		}
	default:
		log.Printf("Unknown, %v", err)
	}
	return response
}

// loadStripeInvoiceForBillingQuote retrieves an invoice created under Stripe Connect when applicable.
func loadStripeInvoiceForBillingQuote(connectAccountID, invoiceID string) (*stripe.Invoice, error) {
	if connectAccountID != "" {
		p := &stripe.InvoiceParams{}
		p.SetStripeAccount(connectAccountID)
		inv, err := invoice.Get(invoiceID, p)
		if err == nil {
			return inv, nil
		}
		log.Printf("invoice.Get %s with Connect account: %v", invoiceID, err)
	}
	return invoice.Get(invoiceID, nil)
}

func GetTradespersonTradespersonIDBillingQuoteQuoteIDInvoiceInvoiceIDHandler(params operations.GetTradespersonTradespersonIDBillingQuoteQuoteIDInvoiceInvoiceIDParams, principal interface{}) middleware.Responder {
	tradespersonID := params.TradespersonID
	quoteID := params.QuoteID
	invoiceID := params.InvoiceID
	token := params.HTTPRequest.Header.Get("Authorization")

	response := operations.NewGetTradespersonTradespersonIDBillingQuoteQuoteIDInvoiceInvoiceIDOK()
	_invoice := operations.GetTradespersonTradespersonIDBillingQuoteQuoteIDInvoiceInvoiceIDOKBody{}

	valid, err := ValidateTradespersonAccessToken(tradespersonID, token)
	if err != nil {
		log.Printf("Failed to validate tradesperson %s, accessToken %s", tradespersonID, token)
		return response
	} else if !valid {
		log.Printf("Bad actor tradesperson %s, accessToken %s", tradespersonID, token)
		return response
	}

	db := database.GetConnection()

	stmt, err := db.Prepare("SELECT id FROM tradesperson_quotes WHERE tradespersonId=? AND quote=?")
	if err != nil {
		return response
	}
	defer stmt.Close()

	var id int64
	row := stmt.QueryRow(tradespersonID, quoteID)
	switch err = row.Scan(&id); err {
	case sql.ErrNoRows:
		log.Printf("Tradesperson %s has no quote %s", tradespersonID, quoteID)
		return response
	case nil:
		tpStripeID, stripeAccErr := database.GetTradespersonStripeID(tradespersonID)
		if stripeAccErr != nil {
			log.Printf("GetTradespersonStripeID %s: %v", tradespersonID, stripeAccErr)
		}

		stripeInvoice, err := loadStripeInvoiceForBillingQuote(tpStripeID, invoiceID)
		if err != nil {
			log.Printf("Failed to get stripe invoice with ID %s, %s", invoiceID, err)
			return errMap(http.StatusBadGateway, "unable to load invoice")
		}

		_invoice.Created = stripeInvoice.Created
		_invoice.Description = stripeInvoice.Description
		_invoice.Total = stripeInvoice.Total
		_invoice.AmountPaid = stripeInvoice.AmountPaid
		_invoice.AmountRemaining = stripeInvoice.AmountRemaining
		_invoice.Status = string(stripeInvoice.Status)
		_invoice.Number = stripeInvoice.Number
		_invoice.Pdf = stripeInvoice.InvoicePDF
		_invoice.URL = stripeInvoice.HostedInvoiceURL

		status, refunded, err := database.GetInvoiceRefund(invoiceID)
		if err != nil {
			log.Printf("Failed to get invoice refund, %v", err)
		}
		if status != "" && refunded != 0 {
			_invoice.Status = status
			_invoice.Refunded = refunded
		}

		products := []*models.Product{}
		lineParams := &stripe.QuoteListLineItemsParams{Quote: stripe.String(quoteID)}
		if tpStripeID != "" {
			lineParams.SetStripeAccount(tpStripeID)
		}
		i := quote.ListLineItems(lineParams)
		for i.Next() {
			lineItem := i.LineItem()
			if lineItem.Price == nil || lineItem.Price.Product == nil || lineItem.Price.Product.ID == "" {
				log.Printf("quote %s line item missing price or product", quoteID)
				continue
			}
			pp := &stripe.ProductParams{}
			if tpStripeID != "" {
				pp.SetStripeAccount(tpStripeID)
			}
			stripeProduct, err := product.Get(lineItem.Price.Product.ID, pp)
			if err != nil {
				log.Printf("Failed to get stripe product, %v", err)
				continue
			}
			_product := &models.Product{}
			_product.Title = stripeProduct.Name
			_product.Price = lineItem.Price.UnitAmount
			_product.Quantity = lineItem.Quantity
			products = append(products, _product)
		}
		_invoice.Products = products

		customer := models.Customer{}
		customer.Name = stripeInvoice.CustomerName
		customer.Email = stripeInvoice.CustomerEmail
		customer.Phone = stripeInvoice.CustomerPhone
		if stripeInvoice.CustomerAddress != nil {
			ca := stripeInvoice.CustomerAddress
			address := models.Address{}
			address.LineOne = ca.Line1
			address.LineTwo = ca.Line2
			address.City = ca.City
			address.State = ca.State
			address.ZipCode = ca.PostalCode
			customer.Address = &address
		}
		_invoice.Customer = &customer

		sqCard, errSqCard := quote.Get(quoteID, billingQuoteParams(tpStripeID, nil))
		if errSqCard != nil {
			log.Printf("invoice card payments quote.Get %s: %v", quoteID, errSqCard)
		} else {
			stripeInvExp, errInv := loadStripeInvoiceForBillingQuoteExpand(tpStripeID, invoiceID, []*string{
				stripe.String("payments.data.payment.payment_intent"),
			})
			if errInv != nil {
				log.Printf("invoice expand for card payments %s: %v", invoiceID, errInv)
			} else if cardRows, tr, tfr, errB := BuildQuoteInvoiceCardPayments(stripeInvExp, sqCard, tpStripeID); errB != nil {
				log.Printf("BuildQuoteInvoiceCardPayments %s: %v", invoiceID, errB)
			} else {
				_invoice.CardPayments = cardRows
				_invoice.TotalCardRefundableCents = tr
				_invoice.TotalCardRefundedCents = tfr
			}
		}

		response.SetPayload(&_invoice)
	default:
		log.Printf("Unknown %v", err)
	}

	return response
}

func PostTradespersonTradespersonIDBillingQuoteQuoteIDInvoiceInvoiceIDFinalizeHandler(params operations.PostTradespersonTradespersonIDBillingQuoteQuoteIDInvoiceInvoiceIDFinalizeParams, principal interface{}) middleware.Responder {
	tradespersonID := params.TradespersonID
	quoteID := params.QuoteID
	invoiceID := params.InvoiceID
	token := params.HTTPRequest.Header.Get("Authorization")

	response := operations.NewPostTradespersonTradespersonIDBillingQuoteQuoteIDInvoiceInvoiceIDFinalizeOK()
	payload := operations.PostTradespersonTradespersonIDBillingQuoteQuoteIDInvoiceInvoiceIDFinalizeOKBody{Finalized: false}

	valid, err := ValidateTradespersonAccessToken(tradespersonID, token)
	if err != nil {
		log.Printf("Failed to validate tradesperson %s, accessToken %s", tradespersonID, token)
		return response
	} else if !valid {
		log.Printf("Bad actor tradesperson %s, accessToken %s", tradespersonID, token)
		return response
	}

	response.SetPayload(&payload)

	db := database.GetConnection()

	stmt, err := db.Prepare("SELECT id FROM tradesperson_quotes WHERE tradespersonId=? AND quote=?")
	if err != nil {
		return response
	}
	defer stmt.Close()

	var id int64
	row := stmt.QueryRow(tradespersonID, quoteID)
	switch err = row.Scan(&id); err {
	case sql.ErrNoRows:
		log.Printf("Tradesperson %s has no quote %s", tradespersonID, quoteID)
		return response
	case nil:
		tpStripeID, stripeAccErr := database.GetTradespersonStripeID(tradespersonID)
		if stripeAccErr != nil || tpStripeID == "" {
			log.Printf("finalize invoice GetTradespersonStripeID %s: %v", tradespersonID, stripeAccErr)
			return response
		}
		finalizeInv := &stripe.InvoiceFinalizeInvoiceParams{}
		finalizeInv.SetStripeAccount(tpStripeID)
		stripeInvoice, err := invoice.FinalizeInvoice(
			invoiceID,
			finalizeInv,
		)
		if err != nil {
			log.Printf("Failed to finalize invoice %s, %v", invoiceID, err)
			return response
		}

		if stripeInvoice.Status == "open" {
			if err := email.SentInvoice(stripeInvoice); err != nil {
				log.Printf("Failed to email customer sent invoice, %v", err)
			}
			payload.Finalized = true
			response.SetPayload(&payload)
		}

	default:
		log.Printf("Unkown %v", err)
	}

	return response
}

func PutTradespersonTradespersonIDBillingQuoteQuoteIDInvoiceInvoiceIDHandler(params operations.PutTradespersonTradespersonIDBillingQuoteQuoteIDInvoiceInvoiceIDParams, principal interface{}) middleware.Responder {
	tradespersonID := params.TradespersonID
	invoiceID := params.InvoiceID
	quoteID := params.QuoteID
	description := params.Body.Description
	token := params.HTTPRequest.Header.Get("Authorization")

	response := operations.NewPutTradespersonTradespersonIDBillingQuoteQuoteIDInvoiceInvoiceIDOK()
	payload := operations.PutTradespersonTradespersonIDBillingQuoteQuoteIDInvoiceInvoiceIDOKBody{Updated: false}
	response.SetPayload(&payload)

	valid, err := ValidateTradespersonAccessToken(tradespersonID, token)
	if err != nil {
		log.Printf("Failed to validate tradesperson %s, accessToken %s", tradespersonID, token)
		return response
	} else if !valid {
		log.Printf("Bad actor tradesperson %s, accessToken %s", tradespersonID, token)
		return response
	}

	db := database.GetConnection()

	stmt, err := db.Prepare("SELECT id FROM tradesperson_quotes WHERE tradespersonId=? AND quote=?")
	if err != nil {
		return response
	}
	defer stmt.Close()

	var id int64
	row := stmt.QueryRow(tradespersonID, quoteID)
	switch err = row.Scan(&id); err {
	case sql.ErrNoRows:
		log.Printf("Tradesperson %s has no quote %s", tradespersonID, quoteID)
		return response
	case nil:
		tpStripeID, stripeAccErr := database.GetTradespersonStripeID(tradespersonID)
		if stripeAccErr != nil || tpStripeID == "" {
			log.Printf("invoice update GetTradespersonStripeID %s: %v", tradespersonID, stripeAccErr)
			return response
		}
		params := &stripe.InvoiceParams{
			Description: stripe.String(description),
		}
		params.SetStripeAccount(tpStripeID)
		_, err := invoice.Update(invoiceID, params)
		if err != nil {
			log.Printf("Failed to updated invoice %s description, %s", invoiceID, err)
			return response
		}
		payload.Updated = true
		response.SetPayload(&payload)
	default:
		log.Printf("Unknown default switch case, %v", err)
	}
	return response
}

func PostTradespersonTradespersonIDBillingQuoteQuoteIDInvoiceInvoiceIDVoidHandler(params operations.PostTradespersonTradespersonIDBillingQuoteQuoteIDInvoiceInvoiceIDVoidParams, principal interface{}) middleware.Responder {
	tradespersonID := params.TradespersonID
	quoteID := params.QuoteID
	invoiceID := params.InvoiceID
	token := params.HTTPRequest.Header.Get("Authorization")

	response := operations.NewPostTradespersonTradespersonIDBillingQuoteQuoteIDInvoiceInvoiceIDVoidOK()
	payload := operations.PostTradespersonTradespersonIDBillingQuoteQuoteIDInvoiceInvoiceIDVoidOKBody{Voided: false}
	response.SetPayload(&payload)

	valid, err := ValidateTradespersonAccessToken(tradespersonID, token)
	if err != nil {
		log.Printf("Failed to validate tradesperson %s, accessToken %s", tradespersonID, token)
		return response
	} else if !valid {
		log.Printf("Bad actor tradesperson %s, accessToken %s", tradespersonID, token)
		return response
	}

	db := database.GetConnection()

	stmt, err := db.Prepare("SELECT q.title, tq.request FROM tradesperson_quotes tq INNER JOIN quotes q ON q.id=tq.quoteId WHERE tq.tradespersonId=? AND tq.quote=?")
	if err != nil {
		return response
	}
	defer stmt.Close()

	var title, message string
	row := stmt.QueryRow(tradespersonID, quoteID)
	switch err = row.Scan(&title, &message); err {
	case sql.ErrNoRows:
		log.Printf("Tradesperson %s has no quote %s", tradespersonID, quoteID)
		return response
	case nil:
		tpStripeID, stripeAccErr := database.GetTradespersonStripeID(tradespersonID)
		if stripeAccErr != nil || tpStripeID == "" {
			log.Printf("void invoice GetTradespersonStripeID %s: %v", tradespersonID, stripeAccErr)
			return response
		}
		voidParams := &stripe.InvoiceVoidInvoiceParams{}
		voidParams.SetStripeAccount(tpStripeID)
		stripeInvoice, err := invoice.VoidInvoice(invoiceID, voidParams)
		if err != nil {
			log.Printf("Failed to void invoice %s, %v", invoiceID, err)
			return response
		}

		if stripeInvoice.Status == "void" {
			payload.Voided = true
			response.SetPayload(&payload)

			tradesperson, err := database.GetTradespersonProfile(tradespersonID)
			if err != nil {
				log.Printf("Failed to get tradesperson profile %s", err)
				return response
			}
			if err := email.SendCustomerQuoteInvoiceVoid(tradesperson, stripeInvoice, message, title); err != nil {
				log.Printf("Failed to send customer email, %v", err)
			}
		}

	default:
		log.Printf("Unkown %v", err)
	}

	return response
}

func PostTradespersonTradespersonIDBillingQuoteQuoteIDInvoiceInvoiceIDRefundHandler(params operations.PostTradespersonTradespersonIDBillingQuoteQuoteIDInvoiceInvoiceIDRefundParams, principal interface{}) middleware.Responder {
	tradespersonID := params.TradespersonID
	quoteID := params.QuoteID
	invoiceID := params.InvoiceID
	token := params.HTTPRequest.Header.Get("Authorization")

	response := operations.NewPostTradespersonTradespersonIDBillingQuoteQuoteIDInvoiceInvoiceIDRefundOK()
	payload := operations.PostTradespersonTradespersonIDBillingQuoteQuoteIDInvoiceInvoiceIDRefundOKBody{Refunded: false}
	response.SetPayload(&payload)

	valid, err := ValidateTradespersonAccessToken(tradespersonID, token)
	if err != nil {
		log.Printf("Failed to validate tradesperson %s, accessToken %s", tradespersonID, token)
		return errMap(http.StatusUnauthorized, "unauthorized")
	} else if !valid {
		log.Printf("Bad actor tradesperson %s, accessToken %s", tradespersonID, token)
		return errMap(http.StatusUnauthorized, "unauthorized")
	}

	db := database.GetConnection()

	stmt, err := db.Prepare("SELECT id FROM tradesperson_quotes WHERE tradespersonId=? AND quote=?")
	if err != nil {
		return errMap(http.StatusInternalServerError, "internal error")
	}
	defer stmt.Close()

	var id int64
	row := stmt.QueryRow(tradespersonID, quoteID)
	switch err = row.Scan(&id); err {
	case sql.ErrNoRows:
		log.Printf("Tradesperson %s has no quote %s", tradespersonID, quoteID)
		return errMap(http.StatusNotFound, "quote not found")
	case nil:
		tpStripeID, stripeAccErr := database.GetTradespersonStripeID(tradespersonID)
		if stripeAccErr != nil || tpStripeID == "" {
			log.Printf("refund GetTradespersonStripeID %s: %v", tradespersonID, stripeAccErr)
			return errMap(http.StatusBadRequest, "provider billing account unavailable")
		}
		sq, err := quote.Get(quoteID, billingQuoteParams(tpStripeID, nil))
		if err != nil {
			log.Printf("refund quote.Get %s: %v", quoteID, err)
			return errMap(http.StatusBadRequest, "unable to load quote")
		}
		if sq.Status != stripe.QuoteStatusAccepted {
			return errMap(http.StatusConflict, "quote must be accepted before refunding the invoice")
		}
		if sq.Invoice == nil || sq.Invoice.ID != invoiceID {
			return errMap(http.StatusNotFound, "invoice does not belong to this quote")
		}

		stripeInvoice, err := loadStripeInvoiceForBillingQuoteExpand(tpStripeID, invoiceID, []*string{
			stripe.String("payments.data.payment.payment_intent"),
		})
		if err != nil {
			log.Printf("Failed to get invoice %s, %v", invoiceID, err)
			return errMap(http.StatusBadRequest, "unable to load invoice")
		}

		cardRows, _, _, err := BuildQuoteInvoiceCardPayments(stripeInvoice, sq, tpStripeID)
		if err != nil {
			log.Printf("BuildQuoteInvoiceCardPayments %s: %v", invoiceID, err)
			return errMap(http.StatusBadRequest, "unable to read invoice payments")
		}
		if len(cardRows) == 0 {
			return errMap(http.StatusConflict, "no card charge found for this invoice; complete the refund in Stripe if needed")
		}

		req := params.Body
		if req != nil && req.ChargeID == "" && req.AmountCents > 0 {
			return errMap(http.StatusBadRequest, "chargeId is required when amountCents is set")
		}

		byCharge := map[string]*models.BillingQuoteInvoiceCardPayment{}
		for _, r := range cardRows {
			if r != nil && r.ChargeID != "" {
				byCharge[r.ChargeID] = r
			}
		}

		refundAll := req == nil || (req.ChargeID == "" && req.AmountCents == 0)
		var totalRefundedCents int64
		var anyRefund bool

		if refundAll {
			for _, r := range cardRows {
				if r == nil || r.AmountRefundableCents <= 0 {
					continue
				}
				sr, _, err := QuoteInvoiceStripeRefund(r.ChargeID, 0, sq, tpStripeID, invoiceID)
				if err != nil {
					log.Printf("Failed to refund charge %s: %v", r.ChargeID, err)
					return errMap(http.StatusBadRequest, err.Error())
				}
				if sr.Status != "succeeded" && sr.Status != "pending" {
					return errMap(http.StatusBadRequest, "refund did not complete")
				}
				totalRefundedCents += sr.Amount
				anyRefund = true
				if err := database.CreateInvoiceRefund(invoiceID, sr.ID); err != nil {
					log.Printf("Failed to create refund in database, %v", err)
				}
			}
			if !anyRefund {
				return errMap(http.StatusConflict, "no refundable card balance on this invoice")
			}
		} else {
			ch := strings.TrimSpace(req.ChargeID)
			if ch == "" {
				return errMap(http.StatusBadRequest, "chargeId is required unless refunding all card charges")
			}
			if _, ok := byCharge[ch]; !ok {
				return errMap(http.StatusBadRequest, "charge does not belong to this invoice")
			}
			amt := req.AmountCents
			if amt < 0 {
				return errMap(http.StatusBadRequest, "amountCents must be zero or positive")
			}
			sr, _, err := QuoteInvoiceStripeRefund(ch, amt, sq, tpStripeID, invoiceID)
			if err != nil {
				log.Printf("Failed to refund charge for invoice, %s", err)
				return errMap(http.StatusBadRequest, err.Error())
			}
			if sr.Status != "succeeded" && sr.Status != "pending" {
				return errMap(http.StatusBadRequest, "refund did not complete")
			}
			totalRefundedCents = sr.Amount
			anyRefund = true
			if err := database.CreateInvoiceRefund(invoiceID, sr.ID); err != nil {
				log.Printf("Failed to create refund in database, %v", err)
			}
		}

		payload.Refunded = anyRefund
		payload.AmountRefundedCents = totalRefundedCents
		response.SetPayload(&payload)

		stripeProduct := &stripe.Product{Name: "Your purchase"}
		if len(stripeInvoice.Lines.Data) > 0 {
			line := stripeInvoice.Lines.Data[0]
			if line.Pricing != nil && line.Pricing.PriceDetails != nil && line.Pricing.PriceDetails.Product != "" {
				pp := &stripe.ProductParams{}
				if sq.TransferData == nil || sq.TransferData.Destination == nil {
					pp.SetStripeAccount(tpStripeID)
				}
				if p, err := product.Get(line.Pricing.PriceDetails.Product, pp); err == nil && p != nil {
					stripeProduct = p
				}
			}
		}
		decimalPrice := float64(totalRefundedCents) / 100.0

		if err := email.SendCustomerRefund(stripeInvoice, stripeProduct, decimalPrice); err != nil {
			log.Printf("Failed to send customer refund email, %v", err)
		}
	default:
		log.Printf("Unkown %v", err)
		return errMap(http.StatusInternalServerError, "internal error")
	}

	return response
}

func PostTradespersonTradespersonIDBillingQuoteQuoteIDInvoiceInvoiceIDUncollectibleHandler(params operations.PostTradespersonTradespersonIDBillingQuoteQuoteIDInvoiceInvoiceIDUncollectibleParams, principal interface{}) middleware.Responder {
	tradespersonID := params.TradespersonID
	quoteID := params.QuoteID
	invoiceID := params.InvoiceID
	token := params.HTTPRequest.Header.Get("Authorization")

	response := operations.NewPostTradespersonTradespersonIDBillingQuoteQuoteIDInvoiceInvoiceIDUncollectibleOK()
	payload := operations.PostTradespersonTradespersonIDBillingQuoteQuoteIDInvoiceInvoiceIDUncollectibleOKBody{Uncollectible: false}
	response.SetPayload(&payload)

	valid, err := ValidateTradespersonAccessToken(tradespersonID, token)
	if err != nil {
		log.Printf("Failed to validate tradesperson %s, accessToken %s", tradespersonID, token)
		return response
	} else if !valid {
		log.Printf("Bad actor tradesperson %s, accessToken %s", tradespersonID, token)
		return response
	}

	db := database.GetConnection()

	stmt, err := db.Prepare("SELECT id FROM tradesperson_quotes WHERE tradespersonId=? AND quote=?")
	if err != nil {
		return response
	}
	defer stmt.Close()

	var id int64
	row := stmt.QueryRow(tradespersonID, quoteID)
	switch err = row.Scan(&id); err {
	case sql.ErrNoRows:
		log.Printf("Tradesperson %s has no quote %s", tradespersonID, quoteID)
		return response
	case nil:
		tpStripeID, stripeAccErr := database.GetTradespersonStripeID(tradespersonID)
		if stripeAccErr != nil || tpStripeID == "" {
			log.Printf("uncollectible GetTradespersonStripeID %s: %v", tradespersonID, stripeAccErr)
			return response
		}
		uncParams := &stripe.InvoiceMarkUncollectibleParams{}
		uncParams.SetStripeAccount(tpStripeID)
		stripeInvoice, err := invoice.MarkUncollectible(
			invoiceID,
			uncParams,
		)
		if err != nil {
			log.Printf("Failed to mark invoice %s uncollectible, %v", invoiceID, err)
			return response
		}

		if stripeInvoice.Status == "uncollectible" {
			payload.Uncollectible = true
			response.SetPayload(&payload)
		}
	default:
		log.Printf("Unkown %v", err)
	}

	return response
}

func GetTradespersonTradespersonIDBillingQuotePagesHandler(params operations.GetTradespersonTradespersonIDBillingQuotePagesParams, principal interface{}) middleware.Responder {
	tradespersonID := params.TradespersonID
	quarter := params.Quarter
	year := params.Year
	token := params.HTTPRequest.Header.Get("Authorization")

	pages := float64(1)
	response := operations.NewGetTradespersonTradespersonIDBillingQuotePagesOK().WithPayload(int64(pages))

	valid, err := ValidateTradespersonAccessToken(tradespersonID, token)
	if err != nil {
		log.Printf("Failed to validate tradesperson %s, accessToken %s", tradespersonID, token)
		return response
	} else if !valid {
		log.Printf("Bad actor tradesperson %s, accessToken %s", tradespersonID, token)
		return response
	}

	db := database.GetConnection()

	stmt, err := db.Prepare("SELECT COUNT(*) FROM tradesperson_quotes WHERE tradespersonId=? AND QUARTER(created) = ? AND YEAR(created) = ?")
	if err != nil {
		log.Printf("Failed to create prepare statement, %v", err)
		return response
	}
	defer stmt.Close()

	err = stmt.QueryRow(tradespersonID, quarter, year).Scan(&pages)
	if err != nil {
		log.Printf("Failed to execute select statement %s", err)
		return response
	}

	if pages == float64(0) {
		pages = float64(1)
	}

	pages = math.Ceil(pages / PAGE_SIZE)

	response.SetPayload(int64(pages))
	return response
}

func quoteDetailAddressFromStripe(a *stripe.Address) *models.Address {
	if a == nil {
		return nil
	}
	return &models.Address{
		LineOne: a.Line1,
		LineTwo: a.Line2,
		City:    a.City,
		State:   a.State,
		ZipCode: a.PostalCode,
	}
}
