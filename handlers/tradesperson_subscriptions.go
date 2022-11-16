package handlers

import (
	"database/sql"
	"log"
	"math"
	"redbudway-api/database"
	"redbudway-api/email"
	"redbudway-api/models"
	"redbudway-api/restapi/operations"

	"github.com/stripe/stripe-go/v72"
	"github.com/stripe/stripe-go/v72/customer"
	"github.com/stripe/stripe-go/v72/invoice"
	"github.com/stripe/stripe-go/v72/price"
	"github.com/stripe/stripe-go/v72/product"
	"github.com/stripe/stripe-go/v72/refund"
	"github.com/stripe/stripe-go/v72/sub"

	"github.com/go-openapi/runtime/middleware"
)

func GetTradespersonTradespersonIDBillingSubscriptionsHandler(params operations.GetTradespersonTradespersonIDBillingSubscriptionsParams, principal interface{}) middleware.Responder {
	tradespersonID := params.TradespersonID
	page := *params.Page

	response := operations.NewGetTradespersonTradespersonIDBillingSubscriptionsOK()
	customers := []*operations.GetTradespersonTradespersonIDBillingSubscriptionsOKBodyItems0{}

	db := database.GetConnection()

	stmt, err := db.Prepare("SELECT subscriptionId, cuStripeId FROM tradesperson_subscriptions WHERE tradespersonId=? GROUP BY id ORDER BY created DESC LIMIT ?, ?")
	if err != nil {
		return response
	}
	defer stmt.Close()

	offSet := (page - 1) * int64(PAGE_SIZE)
	rows, err := stmt.Query(tradespersonID, offSet, PAGE_SIZE)
	if err != nil {
		log.Printf("Failed to execute select statement %s", err)
		return response
	}
	customerSubs := make(map[string]interface{})
	var subscriptionID, cuStripeID string
	for rows.Next() {
		if err := rows.Scan(&subscriptionID, &cuStripeID); err != nil {
			return response
		}

		stripeSubscription, err := sub.Get(subscriptionID, nil)
		if err != nil {
			log.Printf("Failed to get subscription %s, %v", subscriptionID, err)
		}

		if customerSubs[cuStripeID] != nil {
			c := customerSubs[cuStripeID]
			subscriptions := c.(map[string]interface{})
			if stripeSubscription.LatestInvoice != nil {
				//check customer isn't already set
				stripeInvoice, err := invoice.Get(stripeSubscription.LatestInvoice.ID, nil)
				if err != nil {
					log.Printf("Failed to get stripe invoice with ID %s, %s", stripeSubscription.LatestInvoice.ID, err)
					return response
				}
				subscriptions["name"] = *stripeInvoice.CustomerName

			} else {
				stripeCustomer, err := customer.Get(stripeSubscription.Customer.ID, nil)
				if err != nil {
					log.Printf("Failed to get stripe customer, %v", err)
				}
				subscriptions["name"] = stripeCustomer.Name
			}
			if stripeSubscription.Status == "active" {
				subscriptions["active"] = subscriptions["active"].(int64) + int64(1)
			} else if stripeSubscription.Status == "canceled" {
				subscriptions["canceled"] = subscriptions["canceled"].(int64) + int64(1)
			} else if stripeSubscription.Status == "incomplete" {
				subscriptions["incomplete"] = subscriptions["incomplete"].(int64) + int64(1)
			}
			customerSubs[cuStripeID] = subscriptions
		} else {
			subscriptions := map[string]interface{}{}
			subscriptions["stripeId"] = cuStripeID
			subscriptions["active"] = int64(0)
			subscriptions["canceled"] = int64(0)
			subscriptions["incomplete"] = int64(0)
			if stripeSubscription.LatestInvoice != nil {
				//check customer isn't already set
				stripeInvoice, err := invoice.Get(stripeSubscription.LatestInvoice.ID, nil)
				if err != nil {
					log.Printf("Failed to get stripe invoice with ID %s, %s", stripeSubscription.LatestInvoice.ID, err)
					return response
				}
				subscriptions["name"] = *stripeInvoice.CustomerName

			} else {
				stripeCustomer, err := customer.Get(stripeSubscription.Customer.ID, nil)
				if err != nil {
					log.Printf("Failed to get stripe customer, %v", err)
				}
				subscriptions["name"] = stripeCustomer.Name
			}
			if stripeSubscription.Status == "active" {
				subscriptions["active"] = subscriptions["active"].(int64) + int64(1)
			} else if stripeSubscription.Status == "canceled" {
				subscriptions["canceled"] = subscriptions["canceled"].(int64) + int64(1)
			} else if stripeSubscription.Status == "incomplete" {
				subscriptions["incomplete"] = subscriptions["incomplete"].(int64) + int64(1)
			}

			customerSubs[cuStripeID] = subscriptions
		}
	}

	for _, c := range customerSubs {
		subscriptions := c.(map[string]interface{})
		info := operations.GetTradespersonTradespersonIDBillingSubscriptionsOKBodyItems0{}
		info.Name = subscriptions["name"].(string)
		info.StripeID = subscriptions["stripeId"].(string)
		info.Active = subscriptions["active"].(int64)
		info.Canceled = subscriptions["canceled"].(int64)
		info.Incomplete = subscriptions["incomplete"].(int64)
		customers = append(customers, &info)
	}
	response.SetPayload(customers)

	return response
}

func GetTradespersonTradespersonIDBillingSubscriptionPagesHandler(params operations.GetTradespersonTradespersonIDBillingSubscriptionPagesParams, principal interface{}) middleware.Responder {
	tradespersonID := params.TradespersonID

	pages := float64(1)
	response := operations.NewGetTradespersonTradespersonIDBillingSubscriptionPagesOK().WithPayload(int64(pages))

	db := database.GetConnection()

	stmt, err := db.Prepare("SELECT cuStripeId FROM tradesperson_subscriptions WHERE tradespersonId=?")
	if err != nil {
		return response
	}
	defer stmt.Close()

	rows, err := stmt.Query(tradespersonID)
	if err != nil {
		log.Printf("Failed to execute select statement %s", err)
		return response
	}
	customerSubs := make(map[string]interface{})
	var cuStripeID string
	for rows.Next() {
		if err := rows.Scan(&cuStripeID); err != nil {
			return response
		}

		if customerSubs[cuStripeID] == nil {
			customerSubs[cuStripeID] = true
		}
	}

	pages = math.Ceil(float64(len(customerSubs)) / PAGE_SIZE)
	if pages == float64(0) {
		pages = float64(1)
	}
	response.SetPayload(int64(pages))

	return response
}

func GetTradespersonTradespersonIDBillingCustomerStripeIDSubscriptionsHandler(params operations.GetTradespersonTradespersonIDBillingCustomerStripeIDSubscriptionsParams, principal interface{}) middleware.Responder {
	tradespersonID := params.TradespersonID
	stripeID := params.StripeID

	response := operations.NewGetTradespersonTradespersonIDBillingCustomerStripeIDSubscriptionsOK()
	_customer := operations.GetTradespersonTradespersonIDBillingCustomerStripeIDSubscriptionsOKBody{}

	db := database.GetConnection()

	stmt, err := db.Prepare("SELECT ts.fixedPriceId, ts.subscriptionId, fp.subInterval FROM tradesperson_subscriptions ts INNER JOIN fixed_prices fp ON ts.fixedPriceId=fp.id WHERE ts.tradespersonId=? AND ts.cuStripeId=? GROUP BY ts.subscriptionId ORDER BY ts.created DESC")
	if err != nil {
		return response
	}
	defer stmt.Close()

	rows, err := stmt.Query(tradespersonID, stripeID)
	if err != nil {
		log.Printf("Failed to execute select statement %s", err)
		return response
	}

	subscriptions := _customer.Subscriptions

	var fixedPriceID int64
	var subscriptionID, interval string
	for rows.Next() {
		if err := rows.Scan(&fixedPriceID, &subscriptionID, &interval); err != nil {
			return response
		}

		stripeSubscription, err := sub.Get(subscriptionID, nil)
		if err != nil {
			log.Printf("Failed to get subscription, %v", err)
		}

		stripeProduct, err := product.Get(stripeSubscription.Items.Data[0].Price.Product.ID, nil)
		if err != nil {
			log.Printf("Failed to get product %s, %v", stripeSubscription.Items.Data[0].Price.Product.ID, err)
		}

		sameService := false
		for i, subscription := range subscriptions {
			if subscription.Title == stripeProduct.Name {
				sameService = true
				detail := operations.GetTradespersonTradespersonIDBillingCustomerStripeIDSubscriptionsOKBodySubscriptionsItems0DetailsItems0{}
				detail.SubscriptionID = subscriptionID
				detail.Status = string(stripeSubscription.Status)
				if stripeSubscription.LatestInvoice != nil {
					detail.InvoiceID = stripeSubscription.LatestInvoice.ID
					//check customer isn't already set
					stripeInvoice, err := invoice.Get(stripeSubscription.LatestInvoice.ID, nil)
					if err != nil {
						log.Printf("Failed to get stripe invoice with ID %s, %s", stripeSubscription.LatestInvoice.ID, err)
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
				} else {
					stripeCustomer, err := customer.Get(stripeID, nil)
					if err != nil {
						log.Printf("Failed to get stripe customer, %v", err)
					}
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
				timeSlot, err := database.GetSubscriptionTimeSlot(subscriptionID, fixedPriceID)
				if err != nil {
					log.Printf("Failed to get subscription %s time slot, %v", subscriptionID, err)
				}
				detail.TimeSlots = append(detail.TimeSlots, timeSlot)
				subscription.Details = append(subscription.Details, &detail)
			}
			subscriptions[i] = subscription
		}
		if !sameService {
			subscription := operations.GetTradespersonTradespersonIDBillingCustomerStripeIDSubscriptionsOKBodySubscriptionsItems0{}
			subscription.Title = stripeProduct.Name
			subscription.Description = stripeProduct.Description
			subscription.Total = stripeSubscription.Items.Data[0].Price.UnitAmount
			subscription.Interval = interval
			detail := operations.GetTradespersonTradespersonIDBillingCustomerStripeIDSubscriptionsOKBodySubscriptionsItems0DetailsItems0{}
			detail.SubscriptionID = subscriptionID
			detail.Status = string(stripeSubscription.Status)
			if stripeSubscription.LatestInvoice != nil {
				detail.InvoiceID = stripeSubscription.LatestInvoice.ID
				//check customer isn't already set
				stripeInvoice, err := invoice.Get(stripeSubscription.LatestInvoice.ID, nil)
				if err != nil {
					log.Printf("Failed to get stripe invoice with ID %s, %s", stripeSubscription.LatestInvoice.ID, err)
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
			} else {
				stripeCustomer, err := customer.Get(stripeID, nil)
				if err != nil {
					log.Printf("Failed to get stripe customer, %v", err)
				}
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
			timeSlot, err := database.GetSubscriptionTimeSlot(subscriptionID, fixedPriceID)
			if err != nil {
				log.Printf("Failed to get subscription %s time slot, %v", subscriptionID, err)
			}
			detail.TimeSlots = append(detail.TimeSlots, timeSlot)
			subscription.Details = append(subscription.Details, &detail)

			subscriptions = append(subscriptions, &subscription)
		}
	}
	_customer.Subscriptions = subscriptions
	response.SetPayload(&_customer)
	return response
}

func GetTradespersonTradespersonIDBillingCustomerStripeIDSubscriptionSubscriptionIDInvoiceInvoiceIDHandler(params operations.GetTradespersonTradespersonIDBillingCustomerStripeIDSubscriptionSubscriptionIDInvoiceInvoiceIDParams, principal interface{}) middleware.Responder {
	tradespersonID := params.TradespersonID
	subscriptionID := params.SubscriptionID
	invoiceID := params.InvoiceID

	response := operations.NewGetTradespersonTradespersonIDBillingCustomerStripeIDSubscriptionSubscriptionIDInvoiceInvoiceIDOK()
	subscriptionInvoice := operations.GetTradespersonTradespersonIDBillingCustomerStripeIDSubscriptionSubscriptionIDInvoiceInvoiceIDOKBody{}

	db := database.GetConnection()

	stmt, err := db.Prepare("SELECT fp.id, fp.subInterval FROM tradesperson_subscriptions ts LEFT JOIN fixed_prices fp ON ts.fixedPriceId=fp.id WHERE ts.tradespersonId=? AND ts.subscriptionId=?")
	if err != nil {
		log.Printf("Failed to create select statement %s", err)
		return response
	}
	defer stmt.Close()

	rows, err := stmt.Query(tradespersonID, subscriptionID)
	if err != nil {
		log.Printf("Failed to execute select statement %s", err)
		return response
	}

	var fixedPriceID int64
	var interval string
	for rows.Next() {
		if err := rows.Scan(&fixedPriceID, &interval); err != nil {
			return response
		}

		stripeInvoice, err := invoice.Get(invoiceID, nil)
		if err != nil {
			log.Printf("Failed to get stripe invoice, %v", err)
		}

		subscriptionInvoice.Interval = interval
		subscriptionInvoice.Created = stripeInvoice.Created
		subscriptionInvoice.Description = stripeInvoice.Description
		subscriptionInvoice.Status = string(stripeInvoice.Status)
		subscriptionInvoice.Number = stripeInvoice.Number
		subscriptionInvoice.Total = stripeInvoice.Total
		subscriptionInvoice.Pdf = stripeInvoice.InvoicePDF
		subscriptionInvoice.URL = stripeInvoice.HostedInvoiceURL
		subscriptionInvoice.TimeSlot, err = database.GetSubscriptionTimeSlot(subscriptionID, fixedPriceID)
		if err != nil {
			log.Printf("Failed to get subscription %s time slot, %v", subscriptionID, err)
		}

		status, refunded, err := database.GetInvoiceRefund(invoiceID)
		if err != nil {
			log.Printf("Failed to get invoice refund, %v", err)
		}
		if status != "" && refunded != 0 {
			subscriptionInvoice.Status = status
			subscriptionInvoice.Refunded = refunded
		}

		stripeProduct, err := product.Get(stripeInvoice.Lines.Data[0].Price.Product.ID, nil)
		if err != nil {
			log.Printf("Failed to get stripe product %s, %v", stripeInvoice.Lines.Data[0].Price.Product.ID, err)
		}

		subscriptionInvoice.Service = &operations.GetTradespersonTradespersonIDBillingCustomerStripeIDSubscriptionSubscriptionIDInvoiceInvoiceIDOKBodyService{}
		subscriptionInvoice.Service.Title = stripeProduct.Name
		subscriptionInvoice.Service.Description = stripeProduct.Description

		subscriptionInvoice.Customer = &models.Customer{}
		subscriptionInvoice.Customer.Name = *stripeInvoice.CustomerName
		subscriptionInvoice.Customer.Email = stripeInvoice.CustomerEmail
		subscriptionInvoice.Customer.Phone = *stripeInvoice.CustomerPhone
		subscriptionInvoice.Customer.Address = &models.Address{}
		subscriptionInvoice.Customer.Address.LineOne = stripeInvoice.CustomerAddress.Line1
		subscriptionInvoice.Customer.Address.LineTwo = stripeInvoice.CustomerAddress.Line2
		subscriptionInvoice.Customer.Address.City = stripeInvoice.CustomerAddress.City
		subscriptionInvoice.Customer.Address.State = stripeInvoice.CustomerAddress.State
		subscriptionInvoice.Customer.Address.ZipCode = stripeInvoice.CustomerAddress.PostalCode

	}
	response.SetPayload(&subscriptionInvoice)
	return response
}

func PostTradespersonTradespersonIDBillingCustomerStripeIDSubscriptionSubscriptionIDInvoiceInvoiceIDRefundHandler(params operations.PostTradespersonTradespersonIDBillingCustomerStripeIDSubscriptionSubscriptionIDInvoiceInvoiceIDRefundParams, principal interface{}) middleware.Responder {
	tradespersonID := params.TradespersonID
	stripeID := params.StripeID
	subscriptionID := params.SubscriptionID
	invoiceID := params.InvoiceID

	response := operations.NewPostTradespersonTradespersonIDBillingCustomerStripeIDSubscriptionSubscriptionIDInvoiceInvoiceIDRefundOK()
	payload := operations.PostTradespersonTradespersonIDBillingCustomerStripeIDSubscriptionSubscriptionIDInvoiceInvoiceIDRefundOKBody{Refunded: false}
	response.SetPayload(&payload)

	db := database.GetConnection()

	stmt, err := db.Prepare("SELECT ta.stripeId FROM tradesperson_subscriptions ts INNER JOIN tradesperson_account ta ON ts.tradespersonId=ta.tradespersonId WHERE ts.tradespersonId=? AND ts.subscriptionId=? AND ts.cuStripeId=?")
	if err != nil {
		return response
	}
	defer stmt.Close()

	row := stmt.QueryRow(tradespersonID, subscriptionID, stripeID)
	if err != nil {
		log.Printf("Failed to execute select statement %s", err)
		return response
	}

	var tpStripeID string
	switch err = row.Scan(&tpStripeID); err {
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
		log.Printf("Unknown, %v", err)
	}

	return response
}

func PostTradespersonTradespersonIDBillingCustomerStripeIDSubscriptionsCancelHandler(params operations.PostTradespersonTradespersonIDBillingCustomerStripeIDSubscriptionsCancelParams, principal interface{}) middleware.Responder {
	tradespersonID := params.TradespersonID
	cuStripeID := params.StripeID
	subscriptions := params.Subscriptions

	response := operations.NewPostTradespersonTradespersonIDBillingCustomerStripeIDSubscriptionsCancelOK()
	payload := operations.PostTradespersonTradespersonIDBillingCustomerStripeIDSubscriptionsCancelOKBody{Canceled: false}
	response.SetPayload(&payload)

	db := database.GetConnection()

	for _, subscriptionID := range subscriptions {
		stmt, err := db.Prepare("SELECT fpts.startTime, fpts.segmentSize FROM tradesperson_subscriptions ts INNER JOIN customer_time_slots cts ON ts.cuStripeId=cts.cuStripeId INNER JOIN fixed_price_time_slots fpts ON cts.timeSlotId=fpts.id WHERE ts.tradespersonId=? AND ts.cuStripeId=? AND ts.subscriptionId=?")
		if err != nil {
			log.Printf("Failed to create prepared statement, %v", err)
			return response
		}
		defer stmt.Close()

		var startTime, segmentSize string
		row := stmt.QueryRow(tradespersonID, cuStripeID, subscriptionID)
		switch err = row.Scan(&startTime, &segmentSize); err {
		case sql.ErrNoRows:
			log.Printf("Tradesperson %s has no subscription %s", tradespersonID, subscriptionID)
		case nil:
			_, err := sub.Cancel(subscriptionID, nil)
			if err != nil {
				log.Printf("Failed to cancel subscription %s, %v", subscriptionID, err)
			}

		default:
			log.Printf("Unknown, %v", err)
		}
	}

	payload.Canceled = true
	response.SetPayload(&payload)

	return response
}
