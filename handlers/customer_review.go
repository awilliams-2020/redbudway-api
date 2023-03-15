package handlers

import (
	"database/sql"
	"log"
	"redbudway-api/database"
	"redbudway-api/email"
	"redbudway-api/restapi/operations"

	"github.com/go-openapi/runtime/middleware"
	"github.com/stripe/stripe-go/v72/price"
	"github.com/stripe/stripe-go/v72/product"
	"github.com/stripe/stripe-go/v72/quote"
	"github.com/stripe/stripe-go/v72/sub"
)

func GetCustomerCustomerIDQuoteQuoteIDReviewHandler(params operations.GetCustomerCustomerIDQuoteQuoteIDReviewParams, prinicpal interface{}) middleware.Responder {
	customerID := params.CustomerID
	quoteID := params.QuoteID
	token := params.HTTPRequest.Header.Get("Authorization")

	payload := operations.GetCustomerCustomerIDQuoteQuoteIDReviewOKBody{Reviewed: true}
	response := operations.NewGetCustomerCustomerIDQuoteQuoteIDReviewOK().WithPayload(&payload)

	valid, err := ValidateCustomerAccessToken(customerID, token)
	if err != nil {
		log.Printf("Failed to validate customer %s, accessToken %s", customerID, token)
		return response
	} else if !valid {
		log.Printf("Bad actor customer %s, accessToken %s", customerID, token)
		return response
	}

	db := database.GetConnection()

	stmt, err := db.Prepare("SELECT q.id, tq.quote FROM quotes q INNER JOIN tradesperson_quotes tq ON q.id=tq.quoteId WHERE q.quote=?")
	if err != nil {
		log.Printf("Failed to create select statement %s", err)
		return response
	}
	defer stmt.Close()

	rows, err := stmt.Query(quoteID)
	if err != nil {
		log.Printf("Failed to execute select statement %s", err)
		return response
	}

	var ID int64
	var _quote string
	for rows.Next() {
		if err := rows.Scan(&ID, &_quote); err != nil {
			return response
		}
		stripeQuote, err := quote.Get(_quote, nil)
		if err != nil {
			return response
		}

		if stripeQuote.Invoice != nil {
			payload.Reviewed, err = database.CustomerReviewedQuote(customerID, ID)
			if err != nil {
				return response
			}
			if payload.Reviewed {
				response.SetPayload(&payload)
				break
			}
		}
	}
	return response
}

func GetCustomerCustomerIDFixedPricePriceIDReviewHandler(params operations.GetCustomerCustomerIDFixedPricePriceIDReviewParams, prinicpal interface{}) middleware.Responder {
	customerID := params.CustomerID
	priceID := params.PriceID
	token := params.HTTPRequest.Header.Get("Authorization")

	response := operations.NewGetCustomerCustomerIDFixedPricePriceIDReviewOK()
	payload := operations.GetCustomerCustomerIDFixedPricePriceIDReviewOKBody{Reviewed: true}

	valid, err := ValidateCustomerAccessToken(customerID, token)
	if err != nil {
		log.Printf("Failed to validate customer %s, accessToken %s", customerID, token)
		return response
	} else if !valid {
		log.Printf("Bad actor customer %s, accessToken %s", customerID, token)
		return response
	}

	db := database.GetConnection()

	stmt, err := db.Prepare("SELECT ti.fixedPriceId FROM tradesperson_invoices ti INNER JOIN fixed_prices fp ON ti.fixedPriceId=fp.id WHERE ti.customerId=? AND fp.priceId=?")
	if err != nil {
		log.Printf("Failed to create select statement %s", err)
		return response
	}
	defer stmt.Close()

	row := stmt.QueryRow(customerID, priceID)
	var ID int64
	switch err = row.Scan(&ID); err {
	case sql.ErrNoRows:
		//
	case nil:
		var err error
		payload.Reviewed, err = database.CustomerReviewedFixedPrice(customerID, priceID)
		if err != nil {
			log.Printf("Failed to get customer review, %v", err)
		}
	default:
		log.Printf("Unknown %v", err)
	}

	response.SetPayload(&payload)
	return response
}

func GetCustomerCustomerIDSubscriptionPriceIDReviewHandler(params operations.GetCustomerCustomerIDSubscriptionPriceIDReviewParams, prinicpal interface{}) middleware.Responder {
	customerID := params.CustomerID
	priceID := params.PriceID
	token := params.HTTPRequest.Header.Get("Authorization")

	response := operations.NewGetCustomerCustomerIDSubscriptionPriceIDReviewOK()
	payload := operations.GetCustomerCustomerIDSubscriptionPriceIDReviewOKBody{Reviewed: true}

	valid, err := ValidateCustomerAccessToken(customerID, token)
	if err != nil {
		log.Printf("Failed to validate customer %s, accessToken %s", customerID, token)
		return response
	} else if !valid {
		log.Printf("Bad actor customer %s, accessToken %s", customerID, token)
		return response
	}

	db := database.GetConnection()

	stmt, err := db.Prepare("SELECT ts.subscriptionId FROM tradesperson_subscriptions ts INNER JOIN fixed_prices fp ON ts.fixedPriceId=fp.id WHERE ts.cuStripeId=? AND fp.priceId=?")
	if err != nil {
		log.Printf("Failed to create select statement %s", err)
		return response
	}
	defer stmt.Close()

	cuStripeID, err := database.GetCustomerStripeID(customerID)
	if err != nil {
		log.Printf("Failed to get customer stripe ID, %v", err)
		return response
	}

	var subscriptionID string
	row := stmt.QueryRow(cuStripeID, priceID)
	switch err = row.Scan(&subscriptionID); err {
	case sql.ErrNoRows:
		//
	case nil:
		var err error
		stripeSubscription, err := sub.Get(subscriptionID, nil)
		if err != nil {
			log.Printf("Failed to get stripe subscription, %v", err)
			return response
		}
		if stripeSubscription.LatestInvoice != nil {
			payload.Reviewed, err = database.CustomerReviewedSubscription(customerID, priceID)
			if err != nil {
				log.Printf("Failed to get customer review, %v", err)
			}
		}
	default:
		log.Printf("Unknown %v", err)
	}

	response.SetPayload(&payload)
	return response
}

func PostCustomerCustomerIDFixedPricePriceIDReviewHandler(params operations.PostCustomerCustomerIDFixedPricePriceIDReviewParams, prinicpal interface{}) middleware.Responder {
	customerID := params.CustomerID
	priceID := params.PriceID
	message := params.Review.Message
	rating := params.Review.Rating
	token := params.HTTPRequest.Header.Get("Authorization")

	response := operations.NewPostCustomerCustomerIDFixedPricePriceIDReviewOK()
	payload := operations.PostCustomerCustomerIDFixedPricePriceIDReviewOKBody{Rated: false}

	valid, err := ValidateCustomerAccessToken(customerID, token)
	if err != nil {
		log.Printf("Failed to validate customer %s, accessToken %s", customerID, token)
		return response
	} else if !valid {
		log.Printf("Bad actor customer %s, accessToken %s", customerID, token)
		return response
	}

	db := database.GetConnection()

	stmt, err := db.Prepare("SELECT id, tradespersonId FROM fixed_prices WHERE priceId=?")
	if err != nil {
		log.Printf("Failed to create select statement %s", err)
		return response
	}
	defer stmt.Close()

	var fixedPriceID int64
	var tradespersonID string
	row := stmt.QueryRow(priceID)
	switch err = row.Scan(&fixedPriceID, &tradespersonID); err {
	case sql.ErrNoRows:
		log.Printf("Customer with ID %v doesn't have invoice to review", customerID)
	case nil:
		created, err := database.CreateFixedPriceReview(customerID, message, fixedPriceID, rating)
		if err != nil {
			log.Printf("Failed to create fixed price review, %v", err)
			return response
		}

		if created {
			payload.Rated = true
			response.SetPayload(&payload)

			stripePrice, err := price.Get(priceID, nil)
			if err != nil {
				log.Printf("Failed to retrieve stripe price, %s", priceID)
				return response
			}

			stripeProduct, err := product.Get(stripePrice.Product.ID, nil)
			if err != nil {
				log.Printf("Failed to get stripe product %s, %v", stripePrice.Product.ID, err)
				return response
			}

			tradesperson, err := database.GetTradespersonAccount(tradespersonID)
			if err != nil {
				log.Printf("Failed to get tradesperson account %s", tradespersonID)
				return response
			}

			if err := email.SendTradespersonReview(tradesperson, stripeProduct.Name); err != nil {
				log.Printf("Failed to send email to tradesperon %s, %v", tradespersonID, err)
			}
		}
	default:
		log.Printf("Unknown %v", err)
	}

	return response
}

func PostCustomerCustomerIDQuoteQuoteIDReviewHandler(params operations.PostCustomerCustomerIDQuoteQuoteIDReviewParams, prinicpal interface{}) middleware.Responder {
	customerID := params.CustomerID
	quoteID := params.QuoteID
	message := params.Review.Message
	rating := params.Review.Rating
	token := params.HTTPRequest.Header.Get("Authorization")

	response := operations.NewPostCustomerCustomerIDQuoteQuoteIDReviewOK()
	payload := operations.PostCustomerCustomerIDQuoteQuoteIDReviewOKBody{Rated: false}

	valid, err := ValidateCustomerAccessToken(customerID, token)
	if err != nil {
		log.Printf("Failed to validate customer %s, accessToken %s", customerID, token)
		return response
	} else if !valid {
		log.Printf("Bad actor customer %s, accessToken %s", customerID, token)
		return response
	}

	db := database.GetConnection()

	stmt, err := db.Prepare("SELECT id, tradespersonId FROM quotes WHERE quote=?")
	if err != nil {
		log.Printf("Failed to create select statement %s", err)
		return response
	}
	defer stmt.Close()

	var ID int64
	var tradespersonID string
	row := stmt.QueryRow(quoteID)
	switch err = row.Scan(&ID, &tradespersonID); err {
	case sql.ErrNoRows:
		log.Printf("Customer %v doesn't have invoice to review", customerID)
	case nil:
		created, err := database.CreateQuoteReview(customerID, message, ID, rating)
		if err != nil {
			log.Printf("Failed to create fixed price review, %v", err)
			return response
		}

		if created {
			payload.Rated = true
			response.SetPayload(&payload)

			quote, err := database.GetTradespersonQuote(tradespersonID, quoteID)
			if err != nil {
				log.Printf("Failed to get tradesperson quote, %v", err)
				return response
			}

			tradesperson, err := database.GetTradespersonAccount(tradespersonID)
			if err != nil {
				log.Printf("Failed to get tradesperson account %s", tradespersonID)
				return response
			}

			if err := email.SendTradespersonReview(tradesperson, *quote.Title); err != nil {
				log.Printf("Failed to send email to tradesperon %s, %v", tradespersonID, err)
			}
		}
	default:
		log.Printf("Unknown %v", err)
	}

	return response
}
