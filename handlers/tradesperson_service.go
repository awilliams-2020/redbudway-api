package handlers

import (
	"database/sql"
	"log"
	"redbudway-api/database"
	"redbudway-api/restapi/operations"

	"github.com/go-openapi/runtime/middleware"
	"github.com/stripe/stripe-go/v72/customer"
)

func PostTradespersonTradespersonIDFixedPriceHandler(params operations.PostTradespersonTradespersonIDFixedPriceParams, principal interface{}) middleware.Responder {
	tradespersonID := params.TradespersonID
	fixedPrice := params.FixedPrice

	db := database.GetConnection()

	response := operations.NewPostTradespersonTradespersonIDFixedPriceCreated()
	payload := &operations.PostTradespersonTradespersonIDFixedPriceCreatedBody{}
	created := false
	payload.Created = created
	response.SetPayload(payload)

	stmt, err := db.Prepare("SELECT stripeId FROM tradesperson_account WHERE tradespersonId=?")
	if err != nil {
		log.Printf("Failed to create select statement %s", err)
		return response
	}
	defer stmt.Close()

	row := stmt.QueryRow(tradespersonID)
	var stripeId string
	switch err = row.Scan(&stripeId); err {
	case sql.ErrNoRows:
		log.Printf("Tradesperson with id %s doesn't exist", tradespersonID)
	case nil:
		created, err = database.CreateFixedPrice(tradespersonID, fixedPrice)
		if err != nil {
			log.Printf("Failed to create fixed price, %s", err)
		}
		payload.Created = created
		response.SetPayload(payload)
	default:
		log.Printf("Unknown %v", err)
	}

	return response
}

func GetTradespersonTradespersonIDFixedPricePriceIDHandler(params operations.GetTradespersonTradespersonIDFixedPricePriceIDParams, principal interface{}) middleware.Responder {
	tradespersonID := params.TradespersonID
	priceID := params.PriceID

	response := operations.NewGetTradespersonTradespersonIDFixedPriceOK()
	payload := operations.GetTradespersonTradespersonIDFixedPriceOKBody{}
	fixedPrice, fixedPriceID, err := database.GetTradespersonFixedPrice(tradespersonID, priceID)
	if err != nil {
		log.Printf("Failed to get fixed price, %s", err)
		return response
	}
	otherServices, err := database.GetOtherServices(tradespersonID, fixedPriceID)
	if err != nil {
		log.Printf("Failed to get other fixed prices, %s", err)
		return response
	}
	payload.FixedPrice = fixedPrice
	payload.OtherServices = otherServices
	response.SetPayload(&payload)

	return response
}

func PutTradespersonTradespersonIDFixedPricePriceIDHandler(params operations.PutTradespersonTradespersonIDFixedPricePriceIDParams, principal interface{}) middleware.Responder {
	tradespersonID := params.TradespersonID
	priceID := params.PriceID
	fixedPrice := params.FixedPrice

	response := operations.NewPutTradespersonTradespersonIDFixedPriceOK()
	payload := operations.PutTradespersonTradespersonIDFixedPriceOKBody{}
	updated := false
	payload.Updated = updated

	var err error
	payload.Updated, err = database.UpdateFixedPrice(tradespersonID, priceID, fixedPrice)
	if err != nil {
		log.Printf("Failed to update fixed price, %s", err)
		return response
	}
	response.SetPayload(&payload)

	return response
}

func GetTradespersonTradespersonIDFixedPricesHandler(params operations.GetTradespersonTradespersonIDFixedPricesParams, principal interface{}) middleware.Responder {
	tradespersonID := params.TradespersonID

	response := operations.NewGetTradespersonTradespersonIDFixedPricesOK()
	payload := database.GetTradespersonFixedPrices(tradespersonID)
	response.SetPayload(payload)

	return response
}

func GetTradespersonTradespersonIDFixedPriceReviewsHandler(params operations.GetTradespersonTradespersonIDFixedPriceReviewsParams, principal interface{}) middleware.Responder {
	tradespersonID := params.TradespersonID

	db := database.GetConnection()

	response := operations.NewGetTradespersonTradespersonIDFixedPriceReviewsOK()
	reviews := []*operations.GetTradespersonTradespersonIDFixedPriceReviewsOKBodyItems0{}
	response.SetPayload(reviews)

	stmt, err := db.Prepare("SELECT fpr.id, fpr.customerId, fpr.rating, fpr.message, DATE_FORMAT(fpr.date, '%M %D %Y') date FROM fixed_price_reviews fpr INNER JOIN fixed_prices fp ON fpr.fixedPriceId=fp.id WHERE fp.tradespersonId=? AND fpr.responded=0")
	if err != nil {
		log.Printf("Failed to create select statement %s", err)
		return response
	}
	defer stmt.Close()

	rows, err := stmt.Query(tradespersonID)
	if err != nil {
		log.Printf("Failed to execute select statement %s", err)
		return response
	}

	var id, rating int64
	var customerID, message, date string
	for rows.Next() {
		if err := rows.Scan(&id, &customerID, &rating, &message, &date); err != nil {
			return response
		}
		review := &operations.GetTradespersonTradespersonIDFixedPriceReviewsOKBodyItems0{}
		review.ID = id
		review.Rating = rating
		cuStripeID, err := database.GetCustomerStripeID(customerID)
		if err != nil {
			log.Printf("Failed to get customer %s account %v", customerID, err)
		}
		stripeCustomer, err := customer.Get(cuStripeID, nil)
		if err != nil {
			log.Printf("%s", err)
			continue
		}
		review.Customer = stripeCustomer.Name
		review.Message = message
		review.Date = date
		reviews = append(reviews, review)
	}
	response.SetPayload(reviews)

	return response
}

func PostTradespersonTradespersonIDFixedPriceReviewHandler(params operations.PostTradespersonTradespersonIDFixedPriceReviewParams, principal interface{}) middleware.Responder {
	tradespersonID := params.TradespersonID
	respMsg := *params.Review.Response
	reviewID := *params.Review.ReviewID

	db := database.GetConnection()

	response := operations.NewPostTradespersonTradespersonIDFixedPriceReviewOK()
	payload := operations.PostTradespersonTradespersonIDFixedPriceReviewOKBody{Responded: false}
	response.SetPayload(&payload)

	stmt, err := db.Prepare("SELECT fpr.id FROM fixed_price_reviews fpr INNER JOIN fixed_prices fp ON fpr.fixedPriceId=fp.id WHERE fp.tradespersonId=? AND fpr.responded=0 AND fpr.id=?")
	if err != nil {
		log.Printf("Failed to create select statement %s", err)
		return response
	}
	defer stmt.Close()

	row := stmt.QueryRow(tradespersonID, reviewID)

	switch err = row.Scan(&reviewID); err {
	case sql.ErrNoRows:
		log.Printf("Review %d does not exit to respond, %v", reviewID, err)
		return response
	case nil:
		stmt, err := db.Prepare("UPDATE fixed_price_reviews SET responded = 1, respMsg = ?, respDate = NOW() WHERE id = ?")
		if err != nil {
			log.Printf("Failed to create prepare statement, %v", err)
			return response
		}
		defer stmt.Close()

		results, err := stmt.Exec(respMsg, reviewID)
		if err != nil {
			log.Printf("Failed to exec statement, %v", err)
			return response
		}

		rowsAffected, err := results.RowsAffected()
		if err != nil {
			log.Printf("Failed to retrieve affected rows, %v", err)
			return response
		}

		if rowsAffected == 1 {
			payload.Responded = true
			response.SetPayload(&payload)
		}
	default:
		log.Printf("Unknown %v", err)
	}

	return response
}

func PostTradespersonTradespersonIDQuoteHandler(params operations.PostTradespersonTradespersonIDQuoteParams, principal interface{}) middleware.Responder {
	tradespersonID := params.TradespersonID
	quote := params.Quote

	db := database.GetConnection()

	response := operations.NewPostTradespersonTradespersonIDQuoteCreated()
	payload := &operations.PostTradespersonTradespersonIDQuoteCreatedBody{}
	created := false
	payload.Created = created
	response.SetPayload(payload)

	stmt, err := db.Prepare("SELECT stripeId FROM tradesperson_account WHERE tradespersonId=?")
	if err != nil {
		log.Printf("Failed to create select statement %s", err)
		return response
	}
	defer stmt.Close()

	row := stmt.QueryRow(tradespersonID)
	var stripeId string
	switch err = row.Scan(&stripeId); err {
	case sql.ErrNoRows:
		log.Printf("Tradesperson with id %s doesn't exist", tradespersonID)
	case nil:
		created, err = database.CreateQuote(tradespersonID, quote)
		if err != nil {
			log.Printf("Failed to create quote, %s", err)
		}
		payload.Created = created
		response.SetPayload(payload)
	default:
		log.Printf("Unknown %v", err)
	}

	return response
}

func GetTradespersonTradespersonIDQuoteQuoteIDHandler(params operations.GetTradespersonTradespersonIDQuoteQuoteIDParams, principal interface{}) middleware.Responder {
	tradespersonID := params.TradespersonID
	quoteID := params.QuoteID

	response := operations.NewGetTradespersonTradespersonIDQuoteOK()
	quote, err := database.GetTradespersonQuote(tradespersonID, quoteID)
	if err != nil {
		log.Printf("Failed to get quote, %s", err)
		return response
	}

	response.SetPayload(quote)

	return response
}

func PutTradespersonTradespersonIDQuoteQuoteIDHandler(params operations.PutTradespersonTradespersonIDQuoteQuoteIDParams, principal interface{}) middleware.Responder {
	tradespersonID := params.TradespersonID
	quoteID := params.QuoteID
	quote := params.Quote

	response := operations.NewPutTradespersonTradespersonIDQuoteOK()
	payload := operations.PutTradespersonTradespersonIDQuoteOKBody{}

	payload.Updated = false

	var err error
	payload.Updated, err = database.UpdateTradespersonQuote(tradespersonID, quoteID, quote)
	if err != nil {
		log.Printf("Failed to update quote, %s", err)
		return response
	}
	response.SetPayload(&payload)

	return response
}

func GetTradespersonTradespersonIDQuoteReviewsHandler(params operations.GetTradespersonTradespersonIDQuoteReviewsParams, principal interface{}) middleware.Responder {
	tradespersonID := params.TradespersonID

	db := database.GetConnection()

	response := operations.NewGetTradespersonTradespersonIDQuoteReviewsOK()
	reviews := []*operations.GetTradespersonTradespersonIDQuoteReviewsOKBodyItems0{}
	response.SetPayload(reviews)

	stmt, err := db.Prepare("SELECT qr.id, qr.customerId, qr.rating, qr.message, DATE_FORMAT(qr.date, '%M %D %Y') date, q.tradespersonId FROM quote_reviews qr INNER JOIN quotes q ON qr.quoteId=q.id WHERE q.tradespersonId=? AND qr.responded=0")
	if err != nil {
		log.Printf("Failed to create select statement %s", err)
		return response
	}
	defer stmt.Close()

	rows, err := stmt.Query(tradespersonID)
	if err != nil {
		log.Printf("Failed to execute select statement %s", err)
		return response
	}

	var id, rating int64
	var customerID, message, date string
	for rows.Next() {
		if err := rows.Scan(&id, &customerID, &rating, &message, &date); err != nil {
			return response
		}
		review := &operations.GetTradespersonTradespersonIDQuoteReviewsOKBodyItems0{}
		review.ID = id
		review.Rating = rating
		cuStripeID, err := database.GetCustomerStripeID(customerID)
		if err != nil {
			log.Printf("Failed to get customer %s account %v", customerID, err)
		}
		stripeCustomer, err := customer.Get(cuStripeID, nil)
		if err != nil {
			log.Printf("%s", err)
			continue
		}
		review.Customer = stripeCustomer.Name
		review.Message = message
		review.Date = date
		reviews = append(reviews, review)
	}
	response.SetPayload(reviews)

	return response

	return response
}

func GetTradespersonTradespersonIDQuotesHandler(params operations.GetTradespersonTradespersonIDQuotesParams, principal interface{}) middleware.Responder {
	tradespersonID := params.TradespersonID

	response := operations.NewGetTradespersonTradespersonIDQuotesOK()
	services := database.GetTradespersonQuotes(tradespersonID)
	response.SetPayload(services)

	return response
}
