package handlers

import (
	"database/sql"
	"fmt"
	"log"
	"math"
	"os"
	"redbudway-api/database"
	"redbudway-api/email"
	"redbudway-api/models"
	"redbudway-api/restapi/operations"

	"github.com/go-openapi/runtime/middleware"
	"github.com/stripe/stripe-go/v72"
	"github.com/stripe/stripe-go/v72/customer"
	"github.com/stripe/stripe-go/v72/invoice"
	"github.com/stripe/stripe-go/v72/product"
	"github.com/stripe/stripe-go/v72/quote"
	"github.com/stripe/stripe-go/v72/refund"
)

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
		stripeQuote, err := quote.Get(quoteID, nil)
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

		stripeCustomer, err := customer.Get(stripeQuote.Customer.ID, nil)
		if err != nil {
			log.Printf("Failed to get stripe customer, %v", err)
			return response
		}

		if stripeCustomer.Deleted {
			if stripeQuote.Status == "accepted" {
				stripeInvoice, err := invoice.Get(stripeQuote.Invoice.ID, nil)
				if err != nil {
					log.Printf("Failed to get stripe invoice with ID %s, %s", &stripeQuote.Invoice.ID, err)
					return response
				}
				_quote.Customer = *stripeInvoice.CustomerName
			}
		} else {
			_quote.Customer = stripeCustomer.Name
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

	db := database.GetConnection()

	stmt, err := db.Prepare("SELECT tq.request, q.title, q.description FROM tradesperson_quotes tq INNER JOIN quotes q ON tq.quoteId=q.id WHERE tq.tradespersonId=? AND tq.quote=?")
	if err != nil {
		log.Printf("Failed to create prepared statement, %v", err)
		return response
	}
	defer stmt.Close()

	var message, title, description string
	row := stmt.QueryRow(tradespersonID, quoteID)
	switch err = row.Scan(&message, &title, &description); err {
	case sql.ErrNoRows:
		log.Printf("Tradesperson with ID %s has no quote %s", tradespersonID, quoteID)
	case nil:

		stripeQuote, err := quote.Get(quoteID, nil)
		if err != nil {
			log.Printf("Failed to get stripe quote, %v", err)
		}

		_quote.Request = message
		_quote.Created = stripeQuote.Created
		_quote.Status = string(stripeQuote.Status)
		_quote.Number = stripeQuote.Number
		_quote.Description = stripeQuote.Description
		_quote.Expires = stripeQuote.ExpiresAt
		if stripeQuote.Status == "accepted" {
			_quote.InvoiceID = stripeQuote.Invoice.ID
		}

		service := &models.QuoteDetailsService{}
		service.Title = title
		service.Description = description
		products := []*models.Product{}

		params := &stripe.QuoteListLineItemsParams{Quote: stripe.String(quoteID)}
		i := quote.ListLineItems(params)
		for i.Next() {
			lineItem := i.LineItem()
			stripeProduct, err := product.Get(lineItem.Price.Product.ID, nil)
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

		stripeCustomer, err := customer.Get(stripeQuote.Customer.ID, nil)
		if err != nil {
			log.Printf("Failed to get stripe customer, %v", err)
		}

		_customer := &models.Customer{}
		if stripeCustomer.Deleted {
			if stripeQuote.Status == "accepted" {
				stripeInvoice, err := invoice.Get(stripeQuote.Invoice.ID, nil)
				if err != nil {
					log.Printf("Failed to get stripe invoice with ID %s, %s", &stripeQuote.Invoice.ID, err)
					return response
				}
				_customer.Name = *stripeInvoice.CustomerName
				_customer.Email = stripeInvoice.CustomerEmail
				_customer.Phone = *stripeInvoice.CustomerPhone
				_customer.Address = &models.Address{
					LineOne: stripeInvoice.CustomerAddress.Line1,
					LineTwo: stripeInvoice.CustomerAddress.Line2,
					City:    stripeInvoice.CustomerAddress.City,
					State:   stripeInvoice.CustomerAddress.State,
					ZipCode: stripeInvoice.CustomerAddress.PostalCode,
				}
			}
		} else {
			_customer.Name = stripeCustomer.Name
			_customer.Email = stripeCustomer.Email
			_customer.Phone = stripeCustomer.Phone
			_customer.Address = &models.Address{
				LineOne: stripeCustomer.Address.Line1,
				LineTwo: stripeCustomer.Address.Line2,
				City:    stripeCustomer.Address.City,
				State:   stripeCustomer.Address.State,
				ZipCode: stripeCustomer.Address.PostalCode,
			}
		}

		_quote.Customer = _customer
	default:
		log.Printf("Unknown, %v", err)
	}
	response.SetPayload(&_quote)
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

	lineItems := []*stripe.QuoteLineItemParams{}
	total := float64(0.0)
	for _, _product := range products {
		params := &stripe.ProductParams{
			Name: stripe.String(_product.Title),
		}
		stripeProduct, err := product.New(params)
		if err != nil {
			log.Printf("Failed to get create stripe product, %v", err)
			return response
		}
		lineItem := &stripe.QuoteLineItemParams{
			PriceData: &stripe.QuoteLineItemPriceDataParams{
				Currency:   stripe.String("USD"),
				Product:    stripe.String(stripeProduct.ID),
				UnitAmount: &_product.Price,
			},
			Quantity: &_product.Quantity,
		}
		total += float64(_product.Price)
		lineItems = append(lineItems, lineItem)
	}

	sellingFee, err := database.GetTradespersonSellingFee(tradespersonID)
	if err != nil {
		log.Printf("Failed to get tradesperson %s selling fee, %v", tradespersonID, err)
		sellingFee = float64(0.06)
	}
	appFee := total * sellingFee
	fee := int64(appFee * 100)

	quoteParams := &stripe.QuoteParams{
		Description: stripe.String(description),
	}

	if len(products) != 0 {
		quoteParams.LineItems = lineItems
		quoteParams.ApplicationFeeAmount = &fee
	}
	_, err = quote.Update(
		quoteID,
		quoteParams,
	)
	if err != nil {
		log.Printf("Failed ot update quote %s, %v", quoteID, err)
		return response
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
		stripeQuote, err := quote.Cancel(quoteID, nil)
		if err != nil {
			log.Printf("Failed to cancel quote %s, %v", quoteID, err)
		}
		if stripeQuote.Status == "canceled" {

			if stripeQuote.Invoice == nil {
				_, err := database.DeleteQuote(tradespersonID, quoteID)
				if err != nil {
					log.Printf("Failed to delete tradesperson quote, %v", err)
					return response
				}
			}
			payload.Canceled = true
			response.SetPayload(&payload)

			tradesperson, err := database.GetTradespersonAccount(tradespersonID)
			if err != nil {
				log.Printf("Failed to get tradesperson account, %v", err)
				return response
			}
			stripeCustomer, err := customer.Get(stripeQuote.Customer.ID, nil)
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
		stripeQuote, err := quote.FinalizeQuote(quoteID, nil)
		if err != nil {
			log.Printf("Failed to finalize quote %s, %v", quoteID, err)
		}
		if stripeQuote.Status == "open" {
			payload.Finalized = true
			response.SetPayload(&payload)

			stripeCustomer, err := customer.Get(stripeQuote.Customer.ID, nil)
			if err != nil {
				log.Printf("Failed to get stripe customer, %v", err)
				return response
			}
			customerID, err := database.GetCustomerID(stripeCustomer.ID)
			if err != nil {
				log.Printf("Failed to get customer ID, %v", err)
			}
			quoteURL := fmt.Sprintf("https://%sredbudway.com/#/session/customer-login?customerId=%s&quoteId=%s", os.Getenv("SUBDOMAIN"), customerID, quoteID)
			if err := email.SentQuote(stripeCustomer, quoteURL); err != nil {
				log.Printf("Failed to send customer email, %v", err)
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
		stripeInvoice, err := invoice.Get(invoiceID, nil)
		if err != nil {
			log.Printf("Failed to get stripe invoice with ID %s, %s", invoiceID, err)
			return response
		}

		_invoice.Created = stripeInvoice.Created
		_invoice.Description = stripeInvoice.Description
		_invoice.Total = stripeInvoice.Total
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
		params := &stripe.QuoteListLineItemsParams{Quote: stripe.String(quoteID)}
		i := quote.ListLineItems(params)
		for i.Next() {
			lineItem := i.LineItem()
			stripeProduct, err := product.Get(lineItem.Price.Product.ID, nil)
			if err != nil {
				log.Printf("Failed to get stripe product, %v", err)
			}
			_product := &models.Product{}
			_product.Title = stripeProduct.Name
			_product.Price = lineItem.Price.UnitAmount
			_product.Quantity = lineItem.Quantity
			products = append(products, _product)
		}
		_invoice.Products = products

		customer := models.Customer{}
		customer.Name = *stripeInvoice.CustomerName
		customer.Email = stripeInvoice.CustomerEmail
		customer.Phone = *stripeInvoice.CustomerPhone
		address := models.Address{}
		address.LineOne = stripeInvoice.CustomerAddress.Line1
		address.LineTwo = stripeInvoice.CustomerAddress.Line2
		address.City = stripeInvoice.CustomerAddress.City
		address.State = stripeInvoice.CustomerAddress.State
		address.ZipCode = stripeInvoice.CustomerAddress.PostalCode
		customer.Address = &address
		_invoice.Customer = &customer

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
		stripeInvoice, err := invoice.FinalizeInvoice(
			invoiceID,
			nil,
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
		params := &stripe.InvoiceParams{
			Description: stripe.String(description),
		}
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
		stripeInvoice, err := invoice.VoidInvoice(invoiceID, nil)
		if err != nil {
			log.Printf("Failed to void invoice %s, %v", invoiceID, err)
			return response
		}

		if stripeInvoice.Status == "void" {
			payload.Voided = true
			response.SetPayload(&payload)

			tradesperson, err := database.GetTradespersonAccount(tradespersonID)
			if err != nil {
				log.Printf("Failed to get tradesperson account, %v", err)
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
		stripeInvoice, err := invoice.Get(invoiceID, nil)
		if err != nil {
			log.Printf("Failed to get invoice %s, %v", invoiceID, err)
			return response
		}

		params := &stripe.RefundParams{
			Charge:          stripe.String(stripeInvoice.Charge.ID),
			ReverseTransfer: stripe.Bool(true),
		}
		stripeRefund, err := refund.New(params)
		if err != nil {
			log.Printf("Failed to refund charge for invoice, %s", err)
			return response
		}

		if stripeRefund.Status == "succeeded" || stripeRefund.Status == "pending" {
			err := database.CreateInvoiceRefund(invoiceID, stripeRefund.ID)
			if err != nil {
				log.Printf("Failed to create refund in database, %v", err)
			}

			payload.Refunded = true
			response.SetPayload(&payload)

			stripeProduct, err := product.Get(stripeInvoice.Lines.Data[0].Price.ID, nil)
			if err != nil {
				return response
			}

			decimalPrice := float64(stripeInvoice.Lines.Data[0].Price.UnitAmount / 100.00)

			err = email.SendCustomerRefund(stripeInvoice, stripeProduct, decimalPrice)
			if err != nil {
				log.Printf("Failed to send customer refund email, %v", err)
			}
		}
	default:
		log.Printf("Unkown %v", err)
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
		stripeInvoice, err := invoice.MarkUncollectible(
			invoiceID,
			nil,
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
