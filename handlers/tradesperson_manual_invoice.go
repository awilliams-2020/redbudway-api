package handlers

import (
	"database/sql"
	"fmt"
	"log"
	"math"
	"redbudway-api/database"
	"redbudway-api/email"
	"redbudway-api/internal"
	"redbudway-api/models"
	"redbudway-api/restapi/operations"

	"github.com/go-openapi/runtime/middleware"
	"github.com/stripe/stripe-go/v72"
	"github.com/stripe/stripe-go/v72/customer"
	"github.com/stripe/stripe-go/v72/invoice"
	"github.com/stripe/stripe-go/v72/invoiceitem"
	"github.com/stripe/stripe-go/v72/price"
	"github.com/stripe/stripe-go/v72/product"
	"github.com/stripe/stripe-go/v72/refund"
)

func GetTradespersonTradespersonIDBillingCustomersHandler(params operations.GetTradespersonTradespersonIDBillingCustomersParams, principal interface{}) middleware.Responder {
	tradespersonID := params.TradespersonID
	token := params.HTTPRequest.Header.Get("Authorization")

	response := operations.NewGetTradespersonTradespersonIDBillingCustomersOK()
	customers := []*models.Customer{}
	response.SetPayload(customers)

	valid, err := ValidateTradespersonAccessToken(tradespersonID, token)
	if err != nil {
		log.Printf("Failed to validate tradesperson %s, accessToken %s", tradespersonID, token)
		return response
	} else if !valid {
		log.Printf("Bad actor tradesperson %s, accessToken %s", tradespersonID, token)
		return response
	}

	db := database.GetConnection()
	stmt, err := db.Prepare("SELECT a.stripeId FROM customer_account a LEFT JOIN tradesperson_invoices i ON a.customerId=i.customerId LEFT JOIN tradesperson_subscriptions s ON a.stripeId=s.cuStripeId LEFT JOIN tradesperson_quotes q ON a.customerId=q.customerId WHERE i.tradespersonId=? OR s.tradespersonId=? OR q.tradespersonId=? GROUP BY a.stripeId")
	if err != nil {
		log.Printf("Failed to create select statement %s", err)
		return response
	}
	defer stmt.Close()

	rows, err := stmt.Query(tradespersonID, tradespersonID, tradespersonID)
	if err != nil {
		log.Printf("Failed to exec select statement %s", err)
		return response
	}

	var cuStripeID string
	for rows.Next() {
		if err := rows.Scan(&cuStripeID); err != nil {
			log.Printf("Failed to scan row, %v", err)
			return response
		}
		_customer := models.Customer{}
		stripeCustomer, err := customer.Get(cuStripeID, nil)
		if err != nil {
			log.Printf("Failed to retrieve customer, %s", &cuStripeID)
		}
		_customer.ID = stripeCustomer.ID
		_customer.Name = stripeCustomer.Name
		_customer.Address = &models.Address{
			City:    stripeCustomer.Address.City,
			LineOne: stripeCustomer.Address.Line1,
			LineTwo: stripeCustomer.Address.Line2,
			ZipCode: stripeCustomer.Address.PostalCode,
			State:   stripeCustomer.Address.State,
		}
		_customer.Email = stripeCustomer.Email
		_customer.Phone = stripeCustomer.Phone
		customers = append(customers, &_customer)
	}
	response.SetPayload(customers)

	return response
}

func GetTradespersonTradespersonIDBillingManualInvoicesHandler(params operations.GetTradespersonTradespersonIDBillingManualInvoicesParams, principal interface{}) middleware.Responder {
	tradespersonID := params.TradespersonID
	quarter := params.Quarter
	year := params.Year
	page := *params.Page
	token := params.HTTPRequest.Header.Get("Authorization")

	response := operations.NewGetTradespersonTradespersonIDBillingManualInvoicesOK()
	invoices := []*operations.GetTradespersonTradespersonIDBillingManualInvoicesOKBodyItems0{}

	valid, err := ValidateTradespersonAccessToken(tradespersonID, token)
	if err != nil {
		log.Printf("Failed to validate tradesperson %s, accessToken %s", tradespersonID, token)
		return response
	} else if !valid {
		log.Printf("Bad actor tradesperson %s, accessToken %s", tradespersonID, token)
		return response
	}

	db := database.GetConnection()

	stmt, err := db.Prepare("SELECT invoiceId FROM tradesperson_manual_invoices WHERE tradespersonId=? AND QUARTER(created) = ? AND YEAR(created) = ? ORDER BY created DESC LIMIT ?, ?")
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

	var invoiceID string
	for rows.Next() {
		if err := rows.Scan(&invoiceID); err != nil {
			log.Printf("Failed to scan row %v", err)
			return response
		}

		stripeInvoice, err := invoice.Get(invoiceID, nil)
		if err != nil {
			log.Printf("Failed to get stripe invoice, %v", err)
			return response
		}

		invoice := &operations.GetTradespersonTradespersonIDBillingManualInvoicesOKBodyItems0{}
		invoice.Status = string(stripeInvoice.Status)
		invoice.Number = stripeInvoice.Number
		invoice.InvoiceID = stripeInvoice.ID
		invoice.Customer = *stripeInvoice.CustomerName

		status, _, err := database.GetInvoiceRefund(stripeInvoice.ID)
		if err != nil {
			log.Printf("Failed to get invoice refund, %v", err)
		}
		if status != "" {
			invoice.Status = status
		}

		invoices = append(invoices, invoice)
	}
	response.SetPayload(invoices)

	return response
}

func PostTradespersonTradespersonIDBillingManualInvoiceHandler(params operations.PostTradespersonTradespersonIDBillingManualInvoiceParams, principal interface{}) middleware.Responder {
	tradespersonID := params.TradespersonID
	manualInvoice := params.Invoice
	cuStripeID := manualInvoice.CuStripeID
	products := manualInvoice.Products
	token := params.HTTPRequest.Header.Get("Authorization")

	response := operations.NewPostTradespersonTradespersonIDBillingManualInvoiceOK()
	payload := operations.PostTradespersonTradespersonIDBillingManualInvoiceOKBody{Created: false}
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

	stmt, err := db.Prepare("SELECT stripeId FROM tradesperson_account WHERE tradespersonId=?")
	if err != nil {
		log.Printf("Failed to create select statement, %v", err)
		return response
	}
	defer stmt.Close()

	var tpStripeID string
	row := stmt.QueryRow(tradespersonID)
	switch err = row.Scan(&tpStripeID); err {
	case sql.ErrNoRows:
		log.Printf("Tradesperson %s doesn't exist, %s", tradespersonID, err)
		return response
	case nil:
		customerID, err := database.GetCustomerID(cuStripeID)
		if err != nil {
			log.Printf("Failed to get customer account ID, %v", err)
		}

		var decimalPrice float64
		for _, product := range products {
			decimalPrice += float64(product.Price / 100)
		}

		sellingFee, err := database.GetTradespersonSellingFee(tradespersonID)
		if err != nil {
			log.Printf("Failed to retrieve tradesperson selling fee, %v", sellingFee)
			return response
		}

		appFee := decimalPrice * sellingFee
		fee := int64(appFee * 100)

		dueDate, err := internal.GetDueDate(manualInvoice.DueDate)
		if err != nil {
			log.Printf("Failed to get dueDate from date, %v", err)
			return response
		}
		timeStamp := dueDate.Unix()
		invoiceParams := &stripe.InvoiceParams{
			AutomaticTax: &stripe.InvoiceAutomaticTaxParams{
				Enabled: stripe.Bool(true),
			},
			Customer:             stripe.String(cuStripeID),
			CollectionMethod:     stripe.String("send_invoice"),
			ApplicationFeeAmount: &fee,
			TransferData: &stripe.InvoiceTransferDataParams{
				Destination: stripe.String(tpStripeID),
			},
			DueDate:                     &timeStamp,
			OnBehalfOf:                  stripe.String(tpStripeID),
			PendingInvoiceItemsBehavior: stripe.String("exclude"),
		}
		stripeInvoice, err := invoice.New(invoiceParams)
		if err != nil {
			log.Printf("Failed to create new manual invoice %v", err)
			return response
		}

		for _, product := range products {
			invoiceItemParams := &stripe.InvoiceItemParams{
				Customer:    stripe.String(cuStripeID),
				Currency:    stripe.String("USD"),
				UnitAmount:  &product.Price,
				Quantity:    &product.Quantity,
				Description: stripe.String(product.Title),
				Invoice:     stripe.String(stripeInvoice.ID),
			}
			_, err = invoiceitem.New(invoiceItemParams)
			if err != nil {
				log.Printf("Failed to create new invoice item, %v", err)
				return response
			}
		}

		if stripeInvoice.Status == "draft" {
			err := database.SaveManualInvoice(stripeInvoice.ID, customerID, tradespersonID, stripeInvoice.Created)
			if err != nil {
				log.Printf("Failed to save manual invoice to database %v", err)
				return response
			}
			stripeInvoice, err := invoice.FinalizeInvoice(
				stripeInvoice.ID,
				nil,
			)
			if err != nil {
				log.Printf("Failed to finalize invoice %s, %v", stripeInvoice.ID, err)
				return response
			}

			if stripeInvoice.Status == "open" {
				if err := email.SentInvoice(stripeInvoice); err != nil {
					log.Printf("Failed to email customer sent invoice, %v", err)
				}
				payload.Created = true
				response.SetPayload(&payload)
			}
		}
	default:
		log.Printf("Unknown %v", err)
	}

	return response
}

func GetTradespersonTradespersonIDBillingManualInvoiceInvoiceIDHandler(params operations.GetTradespersonTradespersonIDBillingManualInvoiceInvoiceIDParams, principal interface{}) middleware.Responder {
	tradespersonID := params.TradespersonID
	invoiceID := params.InvoiceID
	token := params.HTTPRequest.Header.Get("Authorization")

	response := operations.NewGetTradespersonTradespersonIDBillingManualInvoiceInvoiceIDOK()
	manualInvoice := operations.GetTradespersonTradespersonIDBillingManualInvoiceInvoiceIDOKBody{}

	valid, err := ValidateTradespersonAccessToken(tradespersonID, token)
	if err != nil {
		log.Printf("Failed to validate tradesperson %s, accessToken %s", tradespersonID, token)
		return response
	} else if !valid {
		log.Printf("Bad actor tradesperson %s, accessToken %s", tradespersonID, token)
		return response
	}

	db := database.GetConnection()

	stmt, err := db.Prepare("SELECT id FROM tradesperson_manual_invoices WHERE tradespersonId=? AND invoiceId=?")
	if err != nil {
		return response
	}
	defer stmt.Close()

	var id int64
	row := stmt.QueryRow(tradespersonID, invoiceID)
	switch err = row.Scan(&id); err {
	case sql.ErrNoRows:
		log.Printf("Tradesperson with ID %s has no manual invoice %s", tradespersonID, invoiceID)
		return response
	case nil:
		stripeInvoice, err := invoice.Get(
			invoiceID,
			nil,
		)
		if err != nil {
			log.Printf("Failed to get stripe manual invoice with ID %s, %s", invoiceID, err)
			return response
		}

		manualInvoice.Created = stripeInvoice.Created
		manualInvoice.DueDate = stripeInvoice.DueDate
		manualInvoice.Description = stripeInvoice.Description
		manualInvoice.Total = stripeInvoice.Total
		manualInvoice.Paid = stripeInvoice.Paid
		manualInvoice.Status = string(stripeInvoice.Status)
		manualInvoice.Number = stripeInvoice.Number
		manualInvoice.Pdf = stripeInvoice.InvoicePDF
		manualInvoice.URL = stripeInvoice.HostedInvoiceURL

		status, refunded, err := database.GetInvoiceRefund(invoiceID)
		if err != nil {
			log.Printf("Failed to get manual invoice refund, %v", err)
		}
		if status != "" && refunded != 0 {
			manualInvoice.Status = status
			manualInvoice.Refunded = refunded
		}

		cost := []*models.Product{}
		for _, data := range stripeInvoice.Lines.Data {
			item := models.Product{}
			stripeProduct, err := product.Get(data.Price.Product.ID, nil)
			if err != nil {
				log.Printf("Failed to get stripe product with ID %s, %s", data.Price.Product.ID, err)
				return response
			}
			item.Title = stripeProduct.Name
			item.Price = data.Price.UnitAmount
			item.Quantity = data.Quantity
			cost = append(cost, &item)

		}
		manualInvoice.Cost = cost

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
		manualInvoice.Customer = &customer

		response.SetPayload(&manualInvoice)
	default:
		log.Printf("Unknown %v", err)
	}

	return response
}

func PostTradespersonTradespersonIDBillingManualInvoiceInvoiceIDVoidHandler(params operations.PostTradespersonTradespersonIDBillingManualInvoiceInvoiceIDVoidParams, principal interface{}) middleware.Responder {
	tradespersonID := params.TradespersonID
	invoiceID := params.InvoiceID
	token := params.HTTPRequest.Header.Get("Authorization")

	response := operations.NewPostTradespersonTradespersonIDBillingManualInvoiceInvoiceIDVoidOK()
	payload := operations.PostTradespersonTradespersonIDBillingManualInvoiceInvoiceIDVoidOKBody{Voided: false}
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

	stmt, err := db.Prepare("SELECT a.name, a.email, a.number FROM tradesperson_manual_invoices i INNER JOIN tradesperson_account a ON i.tradespersonId=a.tradespersonId WHERE i.tradespersonId=? AND i.invoiceId=?")
	if err != nil {
		return response
	}
	defer stmt.Close()

	var tradespersonName, tradespersonEmail, tradespersonNumber string
	row := stmt.QueryRow(tradespersonID, invoiceID)
	switch err = row.Scan(&tradespersonName, &tradespersonEmail, &tradespersonNumber); err {
	case sql.ErrNoRows:
		log.Printf("Tradesperson with ID %s has no invoice %s", &tradespersonID, &invoiceID)
		return response
	case nil:
		stripeInvoice, err := invoice.Get(
			invoiceID,
			nil,
		)
		if err != nil {
			log.Printf("Failed to get invoice %s, %v", invoiceID, err)
		}

		var itemsAndPrice string
		for _, data := range stripeInvoice.Lines.Data {
			stripeProduct, err := product.Get(data.Price.Product.ID, nil)
			if err != nil {
				log.Printf("Failed to get stripe product with ID %s, %s", data.Price.Product.ID, err)
				return response
			}
			itemsAndPrice += fmt.Sprintf("<b>%v</b><br>", stripeProduct.Name)
			itemsAndPrice += fmt.Sprintf("Price: $%v<br>", float64(data.Price.UnitAmount)/float64(100))
			itemsAndPrice += fmt.Sprintf("Quantity: %v<br><br>", data.Quantity)
		}

		stripeInvoice, err = invoice.VoidInvoice(
			invoiceID,
			nil,
		)
		if err != nil {
			log.Printf("Failed to void invoice %s, %v", invoiceID, err)
			return response
		}

		if stripeInvoice.Status == "void" {
			if err := email.SendCustomerVoid(tradespersonName, tradespersonEmail, tradespersonNumber, itemsAndPrice, stripeInvoice); err != nil {
				log.Printf("Failed to send customer email, %v", err)
			}
			payload.Voided = true
			response.SetPayload(&payload)
		}

	default:
		log.Printf("Unkown %v", err)
	}

	return response
}

func PostTradespersonTradespersonIDBillingManualInvoiceInvoiceIDRefundHandler(params operations.PostTradespersonTradespersonIDBillingManualInvoiceInvoiceIDRefundParams, principal interface{}) middleware.Responder {
	tradespersonID := params.TradespersonID
	invoiceID := params.InvoiceID
	token := params.HTTPRequest.Header.Get("Authorization")

	response := operations.NewPostTradespersonTradespersonIDBillingManualInvoiceInvoiceIDRefundOK()
	payload := operations.PostTradespersonTradespersonIDBillingManualInvoiceInvoiceIDRefundOKBody{Refunded: false}
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

	stmt, err := db.Prepare("SELECT id FROM tradesperson_invoices WHERE tradespersonId=? AND invoiceId=?")
	if err != nil {
		return response
	}
	defer stmt.Close()

	var id int64
	row := stmt.QueryRow(tradespersonID, invoiceID)
	switch err = row.Scan(&id); err {
	case sql.ErrNoRows:
		log.Printf("Tradesperson with ID %s has no invoice %s", &tradespersonID, &invoiceID)
		return response
	case nil:
		stripeInvoice, err := invoice.Get(
			invoiceID,
			nil,
		)
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
			if err := database.UpdateTimeSlotByInvoice(invoiceID, stripeInvoice.Lines.Data[0].Quantity); err != nil {
				log.Printf("Failed to update time slot current people, %s, %v", invoiceID, err)
			}

			payload.Refunded = true
			response.SetPayload(&payload)

			stripePrice, err := price.Get(stripeInvoice.Lines.Data[0].Price.ID, nil)
			if err != nil {
				log.Printf("Failed to retrieve stripe price, %s", stripeInvoice.Lines.Data[0].Price.ID)
				return response
			}

			stripeProduct, err := product.Get(stripePrice.Product.ID, nil)
			if err != nil {
				log.Printf("Failed to get stripe product %s, %v", stripePrice.Product.ID, err)
				return response
			}

			decimalPrice := stripePrice.UnitAmountDecimal / float64(100.00)

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

func PostTradespersonTradespersonIDBillingManualInvoiceInvoiceIDUncollectibleHandler(params operations.PostTradespersonTradespersonIDBillingManualInvoiceInvoiceIDUncollectibleParams, principal interface{}) middleware.Responder {
	tradespersonID := params.TradespersonID
	invoiceID := params.InvoiceID
	token := params.HTTPRequest.Header.Get("Authorization")

	response := operations.NewPostTradespersonTradespersonIDBillingManualInvoiceInvoiceIDUncollectibleOK()
	payload := operations.PostTradespersonTradespersonIDBillingManualInvoiceInvoiceIDUncollectibleOKBody{Uncollectible: false}
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

	stmt, err := db.Prepare("SELECT id FROM tradesperson_manual_invoices WHERE tradespersonId=? AND invoiceId=?")
	if err != nil {
		return response
	}
	defer stmt.Close()

	var id int64
	row := stmt.QueryRow(tradespersonID, invoiceID)
	switch err = row.Scan(&id); err {
	case sql.ErrNoRows:
		log.Printf("Tradesperson %s has no manual invoice %s", tradespersonID, invoiceID)
		return response
	case nil:
		stripeInvoice, err := invoice.MarkUncollectible(
			invoiceID,
			nil,
		)
		if err != nil {
			log.Printf("Failed to mark manual invoice %s uncollectible, %v", invoiceID, err)
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

func GetTradespersonTradespersonIDBillingManualInvoicePagesHandler(params operations.GetTradespersonTradespersonIDBillingManualInvoicePagesParams, principal interface{}) middleware.Responder {
	tradespersonID := params.TradespersonID
	quarter := params.Quarter
	year := params.Year
	token := params.HTTPRequest.Header.Get("Authorization")

	pages := float64(1)
	response := operations.NewGetTradespersonTradespersonIDBillingManualInvoicePagesOK().WithPayload(int64(pages))

	valid, err := ValidateTradespersonAccessToken(tradespersonID, token)
	if err != nil {
		log.Printf("Failed to validate tradesperson %s, accessToken %s", tradespersonID, token)
		return response
	} else if !valid {
		log.Printf("Bad actor tradesperson %s, accessToken %s", tradespersonID, token)
		return response
	}

	db := database.GetConnection()

	stmt, err := db.Prepare("SELECT COUNT(*) FROM tradesperson_manual_invoices WHERE tradespersonId=? AND QUARTER(created) = ? AND YEAR(created) = ?")
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
	if pages == float64(0) {
		pages = float64(1)
	}
	response.SetPayload(int64(pages))
	return response
}
