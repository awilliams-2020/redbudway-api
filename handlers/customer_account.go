package handlers

import (
	"bytes"
	"database/sql"
	"log"
	"os"
	"redbudway-api/database"
	"redbudway-api/email"
	"redbudway-api/internal"
	"redbudway-api/models"
	"redbudway-api/restapi/operations"
	_stripe "redbudway-api/stripe"
	"time"

	"github.com/go-openapi/runtime/middleware"
	"github.com/stripe/stripe-go/v72"
	"github.com/stripe/stripe-go/v72/customer"
	"github.com/stripe/stripe-go/v72/invoice"
	"github.com/stripe/stripe-go/v72/invoiceitem"
	"github.com/stripe/stripe-go/v72/price"
	"github.com/stripe/stripe-go/v72/product"
	"github.com/stripe/stripe-go/v72/quote"
	"github.com/stripe/stripe-go/v72/sub"
)

func PostCustomerHandler(params operations.PostCustomerParams) middleware.Responder {
	customer := params.Customer

	db := database.GetConnection()

	payload := operations.PostCustomerCreatedBody{Created: false}
	response := operations.NewPostCustomerCreated().WithPayload(&payload)

	stmt, err := db.Prepare("SELECT email FROM customer_account WHERE email=?")
	if err != nil {
		log.Printf("Failed to create select statement %s", err)
		return response
	}
	defer stmt.Close()

	row := stmt.QueryRow(customer.Email)
	var email string
	switch err = row.Scan(&email); err {
	case sql.ErrNoRows:
		stripeAccount, err := _stripe.CreateCustomerStripeAccount(customer)
		if err != nil {
			log.Printf("Failed creating customer stripe connect account %s", err)
			return response
		}
		customerID, err := database.CreateCustomerAccount(customer, stripeAccount)
		if err != nil {
			log.Printf("Failed creating customer account %s", err)
			return response
		}
		payload.Created = true
		payload.CustomerID = customerID.String()
		accessToken, err := internal.GenerateToken(customerID.String(), "customer", "access", time.Minute*15)
		if err != nil {
			log.Printf("Failed to generate JWT, %s", err)
			return response
		}
		payload.AccessToken = accessToken

		refreshToken, err := internal.GenerateToken(customerID.String(), "customer", "refresh", time.Minute*20)
		if err != nil {
			log.Printf("Failed to generate JWT, %s", err)
			return response
		}
		payload.RefreshToken = refreshToken

		saved, err := database.SaveCustomerTokens(customerID.String(), refreshToken, accessToken)
		if err != nil {
			log.Printf("Failed to save customer tokens, %s", err)
			return response
		}
		if !saved {
			log.Printf("No issues, but failed to save customer")
		}

		response.SetPayload(&payload)
	case nil:
		log.Printf("Customer with email %s already exist", email)
	default:
		log.Printf("Unknown %v", err)
	}

	return response
}

func GetCustomerCustomerIDBillingLinkHandler(params operations.GetCustomerCustomerIDBillingLinkParams, principal interface{}) middleware.Responder {
	customerID := params.CustomerID
	token := params.HTTPRequest.Header.Get("Authorization")

	payload := operations.GetCustomerCustomerIDBillingLinkOKBody{}
	response := operations.NewGetCustomerCustomerIDBillingLinkOK()

	valid, err := ValidateCustomerAccessToken(customerID, token)
	if err != nil {
		log.Printf("Failed to validate customer %s, accessToken %s", customerID, token)
		return response
	} else if !valid {
		log.Printf("Bad actor customer %s, accessToken %s", customerID, token)
		return response
	}

	db := database.GetConnection()

	stmt, err := db.Prepare("SELECT stripeId FROM customer_account WHERE customerId=?")
	if err != nil {
		log.Printf("Failed to create select statement %s", err)
		return response
	}
	defer stmt.Close()

	row := stmt.QueryRow(customerID)
	var stripeID string
	switch err = row.Scan(&stripeID); err {
	case sql.ErrNoRows:
		log.Printf("Customer with ID %s does not exist", customerID)
	case nil:
		session, err := _stripe.GetCustomerBillingLink(stripeID)
		if err != nil {
			log.Printf("Failed to create billing session, %s", err)
		}
		payload.URL = session.URL
		response.SetPayload(&payload)
	default:
		log.Printf("Unknown %v", err)
	}

	return response
}

func emailHelper(tradesperson models.Tradesperson, stripePrice *stripe.Price, stripeProduct *stripe.Product, cuStripeID, customerID, body string) {

	stripeCustomer, err := customer.Get(cuStripeID, nil)
	if err != nil {
		log.Printf("Failed to retrieve customer, %s", customerID)
		return
	}

	if err := email.SendCustomerConfirmation(tradesperson, stripeCustomer, stripeProduct, body); err != nil {
		log.Printf("Failed to send customer receipt email, %v", err)
	}

	if err := email.SendTradespersonBooking(tradesperson, stripeCustomer, stripeProduct, body); err != nil {
		log.Printf("Failed to send tradesperson receipt email, %v", err)
	}
}

func PostCustomerCustomerIDFixedPricePriceIDBookHandler(params operations.PostCustomerCustomerIDFixedPricePriceIDBookParams, principal interface{}) middleware.Responder {
	customerID := params.CustomerID
	priceID := params.PriceID
	timeSlots := params.Booking.TimeSlots
	token := params.HTTPRequest.Header.Get("Authorization")

	payload := operations.PostCustomerCustomerIDFixedPricePriceIDBookCreatedBody{Booked: false}
	response := operations.NewPostCustomerCustomerIDFixedPricePriceIDBookCreated().WithPayload(&payload)

	valid, err := ValidateCustomerAccessToken(customerID, token)
	if err != nil {
		log.Printf("Failed to validate customer %s, accessToken %s", customerID, token)
		return response
	} else if !valid {
		log.Printf("Bad actor customer %s, accessToken %s", customerID, token)
		return response
	}

	cuStripeID, err := database.GetCustomerStripeID(customerID)
	if err != nil {
		log.Printf("Failed to retrieve customer account, %v", priceID)
		return response
	}

	tradesperson, tpStripeID, tradespersonID, err := database.GetTradespersonAccountByPriceID(priceID)
	if err != nil {
		log.Printf("Failed to retrieve tradesperson from price id, %s", priceID)
		return response
	}

	stripePrice, err := price.Get(priceID, nil)
	if err != nil {
		log.Printf("Failed to retrieve stripe price, %s", priceID)
		return response
	}

	stripeProduct, err := product.Get(stripePrice.Product.ID, nil)
	if err != nil {
		log.Printf("Failed to retrieve stripe product, %s", &stripePrice.Product.ID)
	}

	sellingFee, err := database.GetTradespersonSellingFee(tradespersonID)
	if err != nil {
		log.Printf("Failed to retrieve tradesperson selling fee, %v", sellingFee)
		return response
	}

	fixedPriceID, err := database.GetFixedPriceID(priceID)
	if err != nil {
		log.Printf("Failed to get fixed price ID, %v", err)
		return response
	}

	var body bytes.Buffer
	for _, timeSlot := range timeSlots {
		decimalPrice := (stripePrice.UnitAmountDecimal / float64(100.00)) * float64(timeSlot.Quantity)
		appFee := decimalPrice * sellingFee
		fee := int64(appFee * 100)

		endDate, err := internal.GetEndDate(timeSlot.EndTime)
		if err != nil {
			log.Printf("Failed to get endDate from endTime, %v", err)
		}
		timeStamp := endDate.Unix()
		invoiceParams := &stripe.InvoiceParams{
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
		invoiceParams.AddMetadata("tradesperson_id", tradespersonID)
		stripeInvoice, err := invoice.New(invoiceParams)
		if err != nil {
			log.Printf("Failed to create new invoice %v", err)
			return response
		}

		invoiceItemParams := &stripe.InvoiceItemParams{
			Customer:    stripe.String(cuStripeID),
			Price:       stripe.String(stripePrice.ID),
			Quantity:    &timeSlot.Quantity,
			Description: stripe.String(stripeProduct.Description),
			Invoice:     stripe.String(stripeInvoice.ID),
		}
		_, err = invoiceitem.New(invoiceItemParams)
		if err != nil {
			log.Printf("Failed to create new invoice item, %v", err)
			return response
		}

		_, err = database.SaveInvoice(stripeInvoice.ID, customerID, tradespersonID, fixedPriceID, stripeInvoice.Created)
		if err != nil {
			log.Printf("Failed to save invoice %v", err)
			return response
		}

		_, err = database.UpdateTakenTimeSlot(stripeInvoice.ID, cuStripeID, timeSlot.StartTime, fixedPriceID, timeSlot.Quantity, timeSlot.ID)
		if err != nil {
			log.Printf("Failed to update time slots %v", err)
			return response
		}

		timeAndPrice, err := internal.CreateTimeAndPrice(timeSlot.StartTime, timeSlot.EndTime, decimalPrice)
		if err != nil {
			log.Printf("Failed to create time and price, %v", err)
			return response
		}
		body.WriteString(timeAndPrice)
	}

	go emailHelper(tradesperson, stripePrice, stripeProduct, cuStripeID, customerID, body.String())

	payload.Booked = true
	response.SetPayload(&payload)

	return response
}

func emailSubscriptionHelper(tradesperson models.Tradesperson, stripePrice *stripe.Price, cuStripeID, body string) {
	stripeProduct, err := product.Get(stripePrice.Product.ID, nil)
	if err != nil {
		log.Printf("Failed to retrieve stripe product, %s", &stripePrice.Product.ID)
		return
	}

	stripeCustomer, err := customer.Get(cuStripeID, nil)
	if err != nil {
		log.Printf("Failed to retrieve customer, %s", &cuStripeID)
		return
	}

	if err := email.SendCustomerSubscriptionConfirmation(tradesperson, stripeCustomer, stripeProduct, body); err != nil {
		log.Printf("Failed to send customer confirmation email, %v", err)
	}

	if err := email.SendTradespersonSubscriptionBooking(tradesperson, stripeCustomer, stripeProduct, body); err != nil {
		log.Printf("Failed to send tradesperson confirmation email, %v", err)
	}
}

func PostCustomerCustomerIDSubscriptionPriceIDBookHandler(params operations.PostCustomerCustomerIDSubscriptionPriceIDBookParams, principal interface{}) middleware.Responder {
	customerID := params.CustomerID
	priceID := params.PriceID
	timeSlots := params.Booking.TimeSlots
	token := params.HTTPRequest.Header.Get("Authorization")

	payload := operations.PostCustomerCustomerIDSubscriptionPriceIDBookCreatedBody{}
	payload.Booked = false
	response := operations.NewPostCustomerCustomerIDSubscriptionPriceIDBookCreated().WithPayload(&payload)

	valid, err := ValidateCustomerAccessToken(customerID, token)
	if err != nil {
		log.Printf("Failed to validate customer %s, accessToken %s", customerID, token)
		return response
	} else if !valid {
		log.Printf("Bad actor customer %s, accessToken %s", customerID, token)
		return response
	}

	cuStripeID, err := database.GetCustomerStripeID(customerID)
	if err != nil {
		log.Printf("Failed to retrieve customer account, %v", customerID)
		return response
	}

	stripePrice, err := price.Get(priceID, nil)
	if err != nil {
		log.Printf("Failed to retrieve stripe price, %s", priceID)
		return response
	}

	tradesperson, tpStripeID, tradespersonID, err := database.GetTradespersonAccountByPriceID(priceID)
	if err != nil {
		log.Printf("Failed to retrieve tradesperson from price id, %s", priceID)
	}

	sellingFee, err := database.GetTradespersonSellingFee(tradespersonID)
	if err != nil {
		log.Printf("Failed to retrieve tradesperson selling fee, %v", sellingFee)
		return response
	}

	fee := sellingFee * float64(100)

	fixedPriceID, err := database.GetFixedPriceID(priceID)
	if err != nil {
		log.Printf("Failed to get fixed price ID, %v", err)
		return response
	}

	interval, err := database.GetFixedPriceInterval(priceID)
	if err != nil {
		log.Printf("Failed to get fixed price interval, %v", err)
		return response
	}

	var body bytes.Buffer
	for _, timeSlot := range timeSlots {
		decimalPrice := (stripePrice.UnitAmountDecimal / float64(100.00)) * float64(timeSlot.Quantity)

		startDate, err := internal.GetStartDate(timeSlot.FutureTime)
		if err != nil {
			log.Printf("Failed to get startDate from futureTime, %v", err)
		}
		timeStamp := startDate.Unix()

		subscriptionParams := &stripe.SubscriptionParams{
			Customer: stripe.String(cuStripeID),
			Items: []*stripe.SubscriptionItemsParams{
				{
					Price:    stripe.String(stripePrice.ID),
					Quantity: &timeSlot.Quantity,
				},
			},
			BillingCycleAnchor:    &timeStamp,
			ApplicationFeePercent: &fee,
			TransferData: &stripe.SubscriptionTransferDataParams{
				Destination: stripe.String(tpStripeID),
			},
			ProrationBehavior: stripe.String("none"),
			PaymentBehavior:   stripe.String("default_incomplete"),
			OnBehalfOf:        stripe.String(tpStripeID),
		}
		subscriptionParams.AddMetadata("tradesperson_id", tradespersonID)
		stripeSubscription, err := sub.New(subscriptionParams)
		if err != nil {
			log.Printf("Failed to create subscription, %v", err)
			return response
		}

		_, err = database.SaveSubscription(stripeSubscription.ID, cuStripeID, tradespersonID, fixedPriceID, stripeSubscription.Created)
		if err != nil {
			log.Printf("Failed to save subscription %v", err)
			return response
		}

		if interval == "week" {
			_, err = database.UpdateWeeklyTimeSlot(stripeSubscription.ID, cuStripeID, timeSlot.StartTime, fixedPriceID, timeSlot.Quantity, timeSlot.ID)
			if err != nil {
				log.Printf("Failed to update time slots %v", err)
				return response
			}
		} else if interval == "month" {
			_, err = database.UpdateMonthlyTimeSlot(stripeSubscription.ID, cuStripeID, timeSlot.StartTime, fixedPriceID, timeSlot.Quantity, timeSlot.ID)
			if err != nil {
				log.Printf("Failed to update time slots %v", err)
				return response
			}
		} else if interval == "year" {
			_, err = database.UpdateYearlyTimeSlot(stripeSubscription.ID, cuStripeID, timeSlot.StartTime, fixedPriceID, timeSlot.Quantity, timeSlot.ID)
			if err != nil {
				log.Printf("Failed to update time slots %v", err)
				return response
			}
		}

		timeAndPrice, err := internal.CreateSubscriptionTimeAndPrice(interval, timeSlot.StartTime, timeSlot.EndTime, decimalPrice)
		if err != nil {
			log.Printf("Failed to create time and price, %v", err)
			return response
		}
		body.WriteString(timeAndPrice)

	}

	go emailSubscriptionHelper(tradesperson, stripePrice, cuStripeID, body.String())

	payload.Booked = true
	response.SetPayload(&payload)

	return response
}

func GetCustomerCustomerIDPaymentDefaultHandler(params operations.GetCustomerCustomerIDPaymentDefaultParams, principal interface{}) middleware.Responder {
	customerID := params.CustomerID
	token := params.HTTPRequest.Header.Get("Authorization")

	payload := operations.GetCustomerCustomerIDPaymentDefaultOKBody{DefaultPayment: false}
	response := operations.NewGetCustomerCustomerIDPaymentDefaultOK().WithPayload(&payload)

	valid, err := ValidateCustomerAccessToken(customerID, token)
	if err != nil {
		log.Printf("Failed to validate customer %s, accessToken %s", customerID, token)
		return response
	} else if !valid {
		log.Printf("Bad actor customer %s, accessToken %s", customerID, token)
		return response
	}

	cuStripeID, err := database.GetCustomerStripeID(customerID)
	if err != nil {
		log.Printf("Failed to retrieve customer account, %v", customerID)
		return response
	}

	stripeCustomer, err := customer.Get(cuStripeID, nil)
	if err != nil {
		log.Printf("Failed to get customer %s stripe account, %v", cuStripeID, err)
		return response
	}

	if stripeCustomer.DefaultSource != nil {
		payload.DefaultPayment = true
		response.SetPayload(&payload)
	}

	return response
}

func emailQuoteHelper(tradesperson models.Tradesperson, quote *models.ServiceDetails, images []string, cuStripeID, message string) {
	stripeCustomer, err := customer.Get(cuStripeID, nil)
	if err != nil {
		log.Printf("Failed to get stripe customer, %v", err)
		return
	}
	if err := email.SendCustomerQuoteConfirmation(tradesperson, stripeCustomer, message, quote); err != nil {
		log.Printf("Failed to send customer email, %v", err)
	}
	images, err = email.SendTradespersonQuoteRequest(tradesperson, stripeCustomer, message, quote, images)
	if err != nil {
		log.Printf("Failed to send tradesperson email, %v", err)
	}
	for _, imagePath := range images {
		err := os.Remove(imagePath)
		if err != nil {
			log.Printf("Failed to delete image, %s", imagePath)
		}
	}
}

func PostCustomerCustomerIDQuoteQuoteIDRequestHandler(params operations.PostCustomerCustomerIDQuoteQuoteIDRequestParams, principal interface{}) middleware.Responder {
	customerID := params.CustomerID
	quoteID := params.QuoteID
	images := params.Request.Images
	message := params.Request.Message
	token := params.HTTPRequest.Header.Get("Authorization")

	response := operations.NewPostCustomerCustomerIDQuoteQuoteIDRequestCreated()
	payload := operations.PostCustomerCustomerIDQuoteQuoteIDRequestCreatedBody{Requested: false}
	response.SetPayload(&payload)

	valid, err := ValidateCustomerAccessToken(customerID, token)
	if err != nil {
		log.Printf("Failed to validate customer %s, accessToken %s", customerID, token)
		return response
	} else if !valid {
		log.Printf("Bad actor customer %s, accessToken %s", customerID, token)
		return response
	}

	db := database.GetConnection()

	stmt, err := db.Prepare("SELECT stripeId FROM customer_account WHERE customerId=?")
	if err != nil {
		log.Printf("Failed to create select statement %s", err)
		return response
	}
	defer stmt.Close()

	row := stmt.QueryRow(customerID)
	var cuStripeID string
	switch err = row.Scan(&cuStripeID); err {
	case sql.ErrNoRows:
	case nil:
		tradesperson, tpStripeID, tradespersonID, err := database.GetTradespersonAccountByQuoteID(quoteID)
		if err != nil {
			log.Printf("Failed to get tradesperson info, %v", err)
			return response
		}
		_quote, err := database.GetTradespersonQuote(tradespersonID, quoteID)
		if err != nil {
			log.Printf("Failed to get tradesperson quote, %v", err)
			return response
		}

		daysDue := int64(7)
		params := &stripe.QuoteParams{
			Customer:         stripe.String(cuStripeID),
			CollectionMethod: stripe.String("send_invoice"),
			TransferData: &stripe.QuoteTransferDataParams{
				Destination: stripe.String(tpStripeID),
			},
			InvoiceSettings: &stripe.QuoteInvoiceSettingsParams{
				DaysUntilDue: &daysDue,
			},
			OnBehalfOf: stripe.String(tpStripeID),
		}
		stripeQuote, err := quote.New(params)
		if err != nil {
			log.Printf("Failed to create stripe quote, %v", err)
			return response
		}

		if stripeQuote.Status == "draft" {
			created, err := database.SaveQuote(stripeQuote.ID, customerID, tradespersonID, message, _quote.ID, stripeQuote.Created)
			if err != nil {
				log.Printf("Failed to save tradesperson quote, %v", err)
				return response
			}
			if created {
				go emailQuoteHelper(tradesperson, _quote, images, cuStripeID, message)

				payload.Requested = true
				response.SetPayload(&payload)
			}
		}
	default:
		log.Printf("Unknown: %v", err)
	}

	return response
}

func DeleteCustomerCustomerIDHandler(params operations.DeleteCustomerCustomerIDParams, principal interface{}) middleware.Responder {
	customerID := params.CustomerID
	token := params.HTTPRequest.Header.Get("Authorization")

	payload := operations.DeleteCustomerCustomerIDOKBody{Deleted: false, Tradespeople: []string{}}
	response := operations.NewDeleteCustomerCustomerIDOK().WithPayload(&payload)

	valid, err := ValidateCustomerAccessToken(customerID, token)
	if err != nil {
		log.Printf("Failed to validate customer %s, accessToken %s", customerID, token)
		return response
	} else if !valid {
		log.Printf("Bad actor customer %s, accessToken %s", customerID, token)
		return response
	}

	_, err = database.GetCustomerStripeID(customerID)
	if err != nil {
		log.Printf("Failed to get customer %s stripe ID, %v", customerID, err)
		return response
	}

	deleted, err := database.DeleteCustomerAccount(customerID)
	if err != nil {
		log.Printf("Failed to delete customer %s account, %v", customerID, err)
		return response
	}
	payload.Deleted = deleted
	response.SetPayload(&payload)

	return response
}

func GetCustomerCustomerIDQuotesHandler(params operations.GetCustomerCustomerIDQuotesParams, principal interface{}) middleware.Responder {
	customerID := params.CustomerID
	token := params.HTTPRequest.Header.Get("Authorization")

	response := operations.NewGetCustomerCustomerIDQuotesOK()
	quotes := []*operations.GetCustomerCustomerIDQuotesOKBodyItems0{}

	valid, err := ValidateCustomerAccessToken(customerID, token)
	if err != nil {
		log.Printf("Failed to validate customer %s, accessToken %s", customerID, token)
		return response
	} else if !valid {
		log.Printf("Bad actor customer %s, accessToken %s", customerID, token)
		return response
	}

	db := database.GetConnection()

	stmt, err := db.Prepare("SELECT quote FROM tradesperson_quotes WHERE customerId=?")
	if err != nil {
		log.Printf("Failed to create prepare statement, %v", err)
		return response
	}
	defer stmt.Close()

	rows, err := stmt.Query(customerID)
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
		if stripeQuote.Status == "open" || stripeQuote.Status == "draft" {
			_quote := &operations.GetCustomerCustomerIDQuotesOKBodyItems0{}
			_quote.Status = string(stripeQuote.Status)
			_quote.Number = stripeQuote.Number
			_quote.QuoteID = quoteID

			quotes = append(quotes, _quote)
		}
	}
	response.SetPayload(quotes)

	return response
}

func GetCustomerCustomerIDQuoteQuoteIDHandler(params operations.GetCustomerCustomerIDQuoteQuoteIDParams, principal interface{}) middleware.Responder {
	customerID := params.CustomerID
	quoteID := params.QuoteID
	token := params.HTTPRequest.Header.Get("Authorization")

	response := operations.NewGetCustomerCustomerIDQuoteQuoteIDOK()
	_quote := operations.GetCustomerCustomerIDQuoteQuoteIDOKBody{}

	valid, err := ValidateCustomerAccessToken(customerID, token)
	if err != nil {
		log.Printf("Failed to validate customer %s, accessToken %s", customerID, token)
		return response
	} else if !valid {
		log.Printf("Bad actor customer %s, accessToken %s", customerID, token)
		return response
	}

	db := database.GetConnection()

	stmt, err := db.Prepare("SELECT tq.tradespersonID, tq.request, q.title, q.description FROM tradesperson_quotes tq INNER JOIN quotes q ON tq.quoteId=q.id WHERE tq.customerId=? AND tq.quote=?")
	if err != nil {
		log.Printf("Failed to create prepared statement, %v", err)
		return response
	}
	defer stmt.Close()

	var tradespersonID, message, title, description string
	row := stmt.QueryRow(customerID, quoteID)
	switch err = row.Scan(&tradespersonID, &message, &title, &description); err {
	case sql.ErrNoRows:
		log.Printf("Customer %s has no quote %s", customerID, quoteID)
	case nil:

		stripeQuote, err := quote.Get(quoteID, nil)
		if err != nil {
			log.Printf("Failed to get stripe quote, %v", err)
		}

		if stripeQuote.Status == "draft" || stripeQuote.Status == "open" {

			_quote.Request = message
			_quote.Created = stripeQuote.Created
			_quote.Status = string(stripeQuote.Status)
			_quote.Number = stripeQuote.Number
			_quote.Description = stripeQuote.Description
			_quote.Expires = stripeQuote.ExpiresAt

			service := &operations.GetCustomerCustomerIDQuoteQuoteIDOKBodyService{}
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

			tradesperson, err := database.GetTradespersonProfile(tradespersonID)
			if err != nil {
				log.Printf("Failed to get tradesperson profile %s", err)
			}
			_tradesperson := &models.Tradesperson{}
			_tradesperson.Email = tradesperson.Email
			_tradesperson.Name = tradesperson.Name
			//check if wanted phone displayed
			_quote.Tradesperson = &tradesperson
		}
	default:
		log.Printf("Unknown, %v", err)
	}
	response.SetPayload(&_quote)
	return response
}

func PostCustomerCustomerIDQuoteQuoteIDAcceptHandler(params operations.PostCustomerCustomerIDQuoteQuoteIDAcceptParams, principal interface{}) middleware.Responder {
	customerID := params.CustomerID
	quoteID := params.QuoteID
	token := params.HTTPRequest.Header.Get("Authorization")

	payload := operations.PostCustomerCustomerIDQuoteQuoteIDAcceptOKBody{Accepted: false}
	response := operations.NewPostCustomerCustomerIDQuoteQuoteIDAcceptOK().WithPayload(&payload)

	valid, err := ValidateCustomerAccessToken(customerID, token)
	if err != nil {
		log.Printf("Failed to validate customer %s, accessToken %s", customerID, token)
		return response
	} else if !valid {
		log.Printf("Bad actor customer %s, accessToken %s", customerID, token)
		return response
	}

	db := database.GetConnection()

	stmt, err := db.Prepare("SELECT q.quote, tq.tradespersonID, tq.request FROM tradesperson_quotes tq INNER JOIN quotes q ON q.id=tq.quoteId WHERE tq.customerId=? AND tq.quote=?")
	if err != nil {
		log.Printf("Failed to create prepared statement, %v", err)
		return response
	}
	defer stmt.Close()

	var _quoteID, tradespersonID, message string
	row := stmt.QueryRow(customerID, quoteID)
	switch err = row.Scan(&_quoteID, &tradespersonID, &message); err {
	case sql.ErrNoRows:
		log.Printf("Customer %s has no quote %s", customerID, quoteID)
	case nil:
		stripeQuote, err := quote.Accept(quoteID, nil)
		if err != nil {
			log.Printf("Failed to accept stripe quote, %v", err)
			return response
		}
		if stripeQuote.Status == "accepted" {
			payload.Accepted = true
			response.SetPayload(&payload)

			tradesperson, err := database.GetTradespersonProfile(tradespersonID)
			if err != nil {
				log.Printf("Failed to get tradesperson profile %s", err)
				return response
			}

			_quote, err := database.GetTradespersonQuote(tradespersonID, _quoteID)
			if err != nil {
				log.Printf("Failed to get tradesperson quote, %v", err)
				return response
			}

			stripeCustomer, err := customer.Get(stripeQuote.Customer.ID, nil)
			if err != nil {
				log.Printf("Failed to get stripe customer, %v", err)
				return response
			}

			if err := email.SendTradespersonQuoteAccepted(tradesperson, stripeCustomer, message, _quote); err != nil {
				log.Printf("Failed to send customer email, %v", err)
			}
		}
	default:
		log.Printf("Unknown, %v", err)
	}
	return response
}

func PutCustomerCustomerIDHandler(params operations.PutCustomerCustomerIDParams, principal interface{}) middleware.Responder {
	customerID := params.CustomerID
	account := params.Account
	token := params.HTTPRequest.Header.Get("Authorization")

	response := operations.NewPutCustomerCustomerIDOK()
	payload := &operations.PutCustomerCustomerIDOKBody{}
	payload.Updated = false
	response.SetPayload(payload)

	valid, err := ValidateCustomerAccessToken(customerID, token)
	if err != nil {
		log.Printf("Failed to validate customer %s, accessToken %s", customerID, token)
		return response
	} else if !valid {
		log.Printf("Bad actor customer %s, accessToken %s", customerID, token)
		return response
	}

	db := database.GetConnection()

	stmt, err := db.Prepare("SELECT stripeId, email, password FROM customer_account WHERE customerId=?")
	if err != nil {
		log.Printf("Failed to create select statement %s", err)
		return response
	}
	defer stmt.Close()

	row := stmt.QueryRow(customerID)
	var stripeID, hashPassword, accountEmail string
	switch err = row.Scan(&stripeID, &accountEmail, &hashPassword); err {
	case sql.ErrNoRows:
		log.Printf("Tradesperson with ID %s doesn't exist", customerID)
	case nil:
		if internal.CheckPasswordHash(account.CurPassword, hashPassword) {
			stmt, err := db.Prepare("UPDATE customer_account SET password=? WHERE customerId = ?")
			if err != nil {
				return response
			}
			defer stmt.Close()

			newHashPassword, err := internal.HashPassword(account.NewPassword)
			if err != nil {
				log.Printf("%s", err)
				return response
			}
			results, err := stmt.Exec(newHashPassword, customerID)
			if err != nil {
				return response
			}

			rowsAffected, err := results.RowsAffected()
			if err != nil {
				return response
			}

			payload.Updated = rowsAffected == 1
			if payload.Updated {
				stripeCustomer, err := customer.Get(stripeID, nil)
				if err != nil {
					log.Printf("Failed to retrieve customer, %s", customerID)
					return response
				}
				if err := email.PasswordUpdated(accountEmail, stripeCustomer.Name); err != nil {
					log.Printf("Failed to send customer email, %v", err)
					return response
				}
			}
		}
	default:
		log.Printf("Unknown %v", err)
	}

	response = operations.NewPutCustomerCustomerIDOK()
	response.SetPayload(payload)
	return response
}
