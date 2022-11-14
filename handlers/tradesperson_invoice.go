package handlers

import (
	"database/sql"
	"log"
	"math"
	"redbudway-api/database"
	"redbudway-api/email"
	"redbudway-api/internal"
	"redbudway-api/models"
	"redbudway-api/restapi/operations"

	"github.com/go-openapi/runtime/middleware"
	"github.com/stripe/stripe-go/v72"
	"github.com/stripe/stripe-go/v72/invoice"
	"github.com/stripe/stripe-go/v72/product"
	"github.com/stripe/stripe-go/v72/refund"
)

func GetTradespersonTradespersonIDBillingInvoiceInvoiceIDHandler(params operations.GetTradespersonTradespersonIDBillingInvoiceInvoiceIDParams, principal interface{}) middleware.Responder {
	tradespersonID := params.TradespersonID
	invoiceID := params.InvoiceID

	response := operations.NewGetTradespersonTradespersonIDBillingInvoiceInvoiceIDOK()
	_invoice := operations.GetTradespersonTradespersonIDBillingInvoiceInvoiceIDOKBody{}

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
		log.Printf("Tradesperson with ID %s has no invoice %s", tradespersonID, invoiceID)
		return response
	case nil:
		stripeInvoice, err := invoice.Get(
			invoiceID,
			nil,
		)
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

		stripeProduct, err := product.Get(stripeInvoice.Lines.Data[0].Price.Product.ID, nil)
		if err != nil {
			log.Printf("Failed to get stripe product with ID %s, %s", stripeInvoice.Lines.Data[0].Price.Product.ID, err)
			return response
		}
		service := operations.GetTradespersonTradespersonIDBillingInvoiceInvoiceIDOKBodyService{}
		service.Title = stripeProduct.Name
		service.Description = stripeProduct.Description
		_invoice.Service = &service

		startTime, segmentSize, err := database.GetInvoiceStartTimeSegmentSize(invoiceID)
		if err != nil {
			log.Printf("Failed to get invoice %s startTime, segmentSize, %s", invoiceID, err)
		}
		if startTime != "" && segmentSize != "" {
			timeSlot := &operations.GetTradespersonTradespersonIDBillingInvoiceInvoiceIDOKBodyTimeSlot{}
			timeSlot.StartTime = startTime
			endTime, err := internal.CreateEndTime(startTime, segmentSize)
			if err != nil {
				log.Printf("Failed to create endTime, %v", err)
				return response
			}
			timeSlot.EndTime = endTime
			_invoice.TimeSlot = timeSlot
		}

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

func PutTradespersonTradespersonIDBillingInvoiceInvoiceIDHandler(params operations.PutTradespersonTradespersonIDBillingInvoiceInvoiceIDParams, principal interface{}) middleware.Responder {
	tradespersonID := params.TradespersonID
	invoiceID := params.InvoiceID
	description := params.Body.Description

	response := operations.NewPutTradespersonTradespersonIDBillingInvoiceInvoiceIDOK()
	payload := operations.PutTradespersonTradespersonIDBillingInvoiceInvoiceIDOKBody{Updated: false}
	response.SetPayload(&payload)

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
		log.Printf("Tradesperson with ID %s has no invoice %s", tradespersonID, invoiceID)
		return response
	case nil:
		params := &stripe.InvoiceParams{
			Description: stripe.String(description),
		}
		_, err := invoice.Update(
			invoiceID,
			params,
		)
		if err != nil {
			log.Printf("Failed to updated invoice with ID %s description, %s", invoiceID, err)
			return response
		}
		payload.Updated = true
		response.SetPayload(&payload)
	default:
		log.Printf("Unknown default switch case, %v", err)
	}
	return response
}

func DeleteTradespersonTradespersonIDBillingInvoiceInvoiceIDHandler(params operations.DeleteTradespersonTradespersonIDBillingInvoiceInvoiceIDParams, principal interface{}) middleware.Responder {
	tradespersonID := params.TradespersonID
	invoiceID := params.InvoiceID

	response := operations.NewDeleteTradespersonTradespersonIDBillingInvoiceInvoiceIDOK()
	payload := operations.DeleteTradespersonTradespersonIDBillingInvoiceInvoiceIDOKBody{Deleted: false}
	response.SetPayload(&payload)

	db := database.GetConnection()

	stmt, err := db.Prepare("SELECT a.name, a.email, a.number FROM tradesperson_invoices i INNER JOIN tradesperson_account a ON i.tradespersonId=a.tradespersonId WHERE i.tradespersonId=? AND i.invoiceId=?")
	if err != nil {
		return response
	}
	defer stmt.Close()

	var tradespersonName, tradespersonEmail, tradespersonNumber string
	row := stmt.QueryRow(tradespersonID, invoiceID)
	switch err = row.Scan(&tradespersonName, &tradespersonEmail, &tradespersonNumber); err {
	case sql.ErrNoRows:
		log.Printf("Tradesperson with ID %s has no invoice %s", tradespersonID, invoiceID)
		return response
	case nil:
		stripeInvoice, err := invoice.Get(
			invoiceID,
			nil,
		)
		if err != nil {
			log.Printf("Failed to get invoice %s, %v", invoiceID, err)
		}

		stripeProduct, err := product.Get(stripeInvoice.Lines.Data[0].Price.Product.ID, nil)
		if err != nil {
			return response
		}

		decimalPrice := float64(stripeInvoice.Lines.Data[0].Price.UnitAmount / 100.00)

		startTime, segmentSize, err := database.GetInvoiceStartTimeSegmentSize(invoiceID)
		if err != nil {
			log.Printf("Failed to get invoice %s startTime and segmentSize, %v", invoiceID, err)
			return response
		}

		endTime, err := internal.CreateEndTime(startTime, segmentSize)
		if err != nil {
			log.Printf("Failed to create endTime for invoice %s, %v", invoiceID, err)
			return response
		}

		timeAndPrice, err := internal.CreateTimeAndPriceFrmDB(startTime, endTime, decimalPrice)
		if err != nil {
			log.Printf("Failed to create time and price, %v", err)
			return response
		}

		if err := database.ResetTakenTimeSlotByInvoice(invoiceID); err != nil {
			log.Printf("Failed to reset taken time slot, %v", err)
			return response
		}

		deleted, err := database.DeleteInvoice(tradespersonID, invoiceID)
		if err != nil {
			log.Printf("Failed to delete database invoice %s, %v", invoiceID, err)
			return response
		}
		if deleted {
			_, err = invoice.Del(invoiceID, nil)
			if err != nil {
				log.Printf("Failed to delete stripe invoice %s description, %s", invoiceID, err)
				return response
			}
			if err := email.SendCustomerCancellation(tradespersonName, tradespersonEmail, tradespersonNumber, timeAndPrice, stripeInvoice, stripeProduct); err != nil {
				log.Printf("Failed to send customer email, %v", err)
			}
			payload.Deleted = true
			response.SetPayload(&payload)
		}
	default:
		log.Printf("Unkown %v", err)
	}

	return response
}

func PostTradespersonTradespersonIDBillingInvoiceInvoiceIDFinalizeHandler(params operations.PostTradespersonTradespersonIDBillingInvoiceInvoiceIDFinalizeParams, principal interface{}) middleware.Responder {
	tradespersonID := params.TradespersonID
	invoiceID := params.InvoiceID

	response := operations.NewPostTradespersonTradespersonIDBillingInvoiceInvoiceIDFinalizeOK()
	payload := operations.PostTradespersonTradespersonIDBillingInvoiceInvoiceIDFinalizeOKBody{Finalized: false}
	response.SetPayload(&payload)

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
		log.Printf("Tradesperson with ID %s has no invoice %s", tradespersonID, invoiceID)
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

func PostTradespersonTradespersonIDBillingInvoiceInvoiceIDVoidHandler(params operations.PostTradespersonTradespersonIDBillingInvoiceInvoiceIDVoidParams, principal interface{}) middleware.Responder {
	tradespersonID := params.TradespersonID
	invoiceID := params.InvoiceID

	response := operations.NewPostTradespersonTradespersonIDBillingInvoiceInvoiceIDVoidOK()
	payload := operations.PostTradespersonTradespersonIDBillingInvoiceInvoiceIDVoidOKBody{Voided: false}
	response.SetPayload(&payload)

	db := database.GetConnection()

	stmt, err := db.Prepare("SELECT a.name, a.email, a.number FROM tradesperson_invoices i INNER JOIN tradesperson_account a ON i.tradespersonId=a.tradespersonId WHERE i.tradespersonId=? AND i.invoiceId=?")
	if err != nil {
		return response
	}
	defer stmt.Close()

	var tradespersonName, tradespersonEmail, tradespersonNumber string
	row := stmt.QueryRow(tradespersonID, invoiceID)
	switch err = row.Scan(&tradespersonName, &tradespersonEmail, &tradespersonNumber); err {
	case sql.ErrNoRows:
		log.Printf("Tradesperson with ID %s has no invoice %s", tradespersonID, invoiceID)
		return response
	case nil:
		stripeInvoice, err := invoice.Get(
			invoiceID,
			nil,
		)
		if err != nil {
			log.Printf("Failed to get invoice %s, %v", invoiceID, err)
		}

		stripeProduct, err := product.Get(stripeInvoice.Lines.Data[0].Price.ID, nil)
		if err != nil {
			return response
		}

		decimalPrice := float64(stripeInvoice.Lines.Data[0].Price.UnitAmount / 100.00)

		startTime, segmentSize, err := database.GetInvoiceStartTimeSegmentSize(invoiceID)
		if err != nil {
			log.Printf("Failed to get invoice %s startTime and segmentSize, %v", invoiceID, err)
			return response
		}

		endTime, err := internal.CreateEndTime(startTime, segmentSize)
		if err != nil {
			log.Printf("Failed to create endTime for invoice %s, %v", invoiceID, err)
			return response
		}

		timeAndPrice, err := internal.CreateTimeAndPriceFrmDB(startTime, endTime, decimalPrice)
		if err != nil {
			log.Printf("Failed to create time and price, %v", err)
			return response
		}

		if err := database.ResetTakenTimeSlotByInvoice(invoiceID); err != nil {
			log.Printf("Failed to reset taken time slot, %v", err)
			return response
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
			if err := email.SendCustomerCancellation(tradespersonName, tradespersonEmail, tradespersonNumber, timeAndPrice, stripeInvoice, stripeProduct); err != nil {
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

func PostTradespersonTradespersonIDBillingInvoiceInvoiceIDRefundHandler(params operations.PostTradespersonTradespersonIDBillingInvoiceInvoiceIDRefundParams, principal interface{}) middleware.Responder {
	tradespersonID := params.TradespersonID
	invoiceID := params.InvoiceID

	response := operations.NewPostTradespersonTradespersonIDBillingInvoiceInvoiceIDRefundOK()
	payload := operations.PostTradespersonTradespersonIDBillingInvoiceInvoiceIDRefundOKBody{Refunded: false}
	response.SetPayload(&payload)

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
		log.Printf("Tradesperson with ID %s has no invoice %s", tradespersonID, invoiceID)
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
			err = database.ResetTakenTimeSlotByInvoice(invoiceID)
			if err != nil {
				log.Printf("Failed to reset taken time slot for invoice %s, %v", &invoiceID, err)
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

func PostTradespersonTradespersonIDBillingInvoiceInvoiceIDUncollectibleHandler(params operations.PostTradespersonTradespersonIDBillingInvoiceInvoiceIDUncollectibleParams, principal interface{}) middleware.Responder {
	tradespersonID := params.TradespersonID
	invoiceID := params.InvoiceID

	response := operations.NewPostTradespersonTradespersonIDBillingInvoiceInvoiceIDUncollectibleOK()
	payload := operations.PostTradespersonTradespersonIDBillingInvoiceInvoiceIDUncollectibleOKBody{Uncollectible: false}
	response.SetPayload(&payload)

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
		log.Printf("Tradesperson with ID %s has no invoice %s", tradespersonID, invoiceID)
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

func GetTradespersonTradespersonIDBillingInvoicesHandler(params operations.GetTradespersonTradespersonIDBillingInvoicesParams, principal interface{}) middleware.Responder {
	tradespersonID := params.TradespersonID
	quarter := params.Quarter
	year := params.Year
	page := *params.Page

	response := operations.NewGetTradespersonTradespersonIDBillingInvoicesOK()
	invoices := []*operations.GetTradespersonTradespersonIDBillingInvoicesOKBodyItems0{}

	db := database.GetConnection()

	stmt, err := db.Prepare("SELECT invoiceId FROM tradesperson_invoices WHERE tradespersonId=? AND QUARTER(created) = ? AND YEAR(created) = ? ORDER BY created DESC LIMIT ?, ?")
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

		invoice := &operations.GetTradespersonTradespersonIDBillingInvoicesOKBodyItems0{}
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

func GetTradespersonTradespersonIDBillingInvoicePagesHandler(params operations.GetTradespersonTradespersonIDBillingInvoicePagesParams, principal interface{}) middleware.Responder {
	tradespersonID := params.TradespersonID
	quarter := params.Quarter
	year := params.Year

	pages := float64(1)
	response := operations.NewGetTradespersonTradespersonIDBillingInvoicePagesOK().WithPayload(int64(pages))

	db := database.GetConnection()

	stmt, err := db.Prepare("SELECT COUNT(*) FROM tradesperson_invoices WHERE tradespersonId=? AND QUARTER(created) = ? AND YEAR(created) = ?")
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
