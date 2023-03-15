package handlers

import (
	"database/sql"
	"log"
	"math"
	"redbudway-api/database"
	"redbudway-api/internal"
	"redbudway-api/models"
	"redbudway-api/restapi/operations"

	"github.com/go-openapi/runtime/middleware"
	"github.com/stripe/stripe-go/v72/customer"
)

func PostTradespersonTradespersonIDFixedPriceHandler(params operations.PostTradespersonTradespersonIDFixedPriceParams, principal interface{}) middleware.Responder {
	tradespersonID := params.TradespersonID
	fixedPrice := params.FixedPrice
	token := params.HTTPRequest.Header.Get("Authorization")

	response := operations.NewPostTradespersonTradespersonIDFixedPriceCreated()
	payload := &operations.PostTradespersonTradespersonIDFixedPriceCreatedBody{}
	created := false
	payload.Created = created
	response.SetPayload(payload)

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
		log.Printf("Failed to create select statement %s", err)
		return response
	}
	defer stmt.Close()

	row := stmt.QueryRow(tradespersonID)
	var stripeID string
	switch err = row.Scan(&stripeID); err {
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
	accessToken := params.AccessToken
	token := params.HTTPRequest.Header.Get("Authorization")

	response := operations.NewGetTradespersonTradespersonIDFixedPricePriceIDOK()
	payload := operations.GetTradespersonTradespersonIDFixedPricePriceIDOKBody{}

	valid, err := ValidateTradespersonAccessToken(tradespersonID, token)
	if err != nil {
		log.Printf("Failed to validate tradesperson %s, accessToken %s", tradespersonID, token)
		return response
	} else if !valid {
		log.Printf("Bad actor tradesperson %s, accessToken %s", tradespersonID, token)
		return response
	}

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
	googleTimeSlots := models.GoogleTimeSlots{}
	if accessToken != nil {
		googleTimeSlots = internal.GetGoogleTimeSlots(*accessToken)
		if err != nil {
			log.Printf("Failed to get google time slots, %s", err)
			return response
		}
	}
	payload.FixedPrice = fixedPrice
	payload.OtherServices = otherServices
	payload.GoogleTimeSlots = googleTimeSlots
	response.SetPayload(&payload)

	return response
}

func PutTradespersonTradespersonIDFixedPricePriceIDHandler(params operations.PutTradespersonTradespersonIDFixedPricePriceIDParams, principal interface{}) middleware.Responder {
	tradespersonID := params.TradespersonID
	priceID := params.PriceID
	fixedPrice := params.FixedPrice
	token := params.HTTPRequest.Header.Get("Authorization")

	response := operations.NewPutTradespersonTradespersonIDFixedPricePriceIDOK()
	payload := operations.PutTradespersonTradespersonIDFixedPricePriceIDOKBody{}
	updated := false
	payload.Updated = updated

	valid, err := ValidateTradespersonAccessToken(tradespersonID, token)
	if err != nil {
		log.Printf("Failed to validate tradesperson %s, accessToken %s", tradespersonID, token)
		return response
	} else if !valid {
		log.Printf("Bad actor tradesperson %s, accessToken %s", tradespersonID, token)
		return response
	}

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
	page := params.Page
	token := params.HTTPRequest.Header.Get("Authorization")

	response := operations.NewGetTradespersonTradespersonIDFixedPricesOK()

	valid, err := ValidateTradespersonAccessToken(tradespersonID, token)
	if err != nil {
		log.Printf("Failed to validate tradesperson %s, accessToken %s", tradespersonID, token)
		return response
	} else if !valid {
		log.Printf("Bad actor tradesperson %s, accessToken %s", tradespersonID, token)
		return response
	}

	payload := database.GetTradespersonFixedPrices(tradespersonID, page)
	response.SetPayload(payload)

	return response
}

func GetTradespersonTradespersonIDFixedPricePagesHandler(params operations.GetTradespersonTradespersonIDFixedPricePagesParams, principal interface{}) middleware.Responder {
	tradespersonID := params.TradespersonID
	token := params.HTTPRequest.Header.Get("Authorization")

	pages := float64(1)
	response := operations.NewGetTradespersonTradespersonIDFixedPricePagesOK().WithPayload(int64(pages))

	valid, err := ValidateTradespersonAccessToken(tradespersonID, token)
	if err != nil {
		log.Printf("Failed to validate tradesperson %s, accessToken %s", tradespersonID, token)
		return response
	} else if !valid {
		log.Printf("Bad actor tradesperson %s, accessToken %s", tradespersonID, token)
		return response
	}

	db := database.GetConnection()

	stmt, err := db.Prepare("SELECT COUNT(*) FROM fixed_prices WHERE tradespersonId=?")
	if err != nil {
		log.Printf("Failed to create select statement %s", err)
		return response
	}
	defer stmt.Close()

	err = stmt.QueryRow(tradespersonID).Scan(&pages)
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

func GetTradespersonTradespersonIDFixedPriceReviewsHandler(params operations.GetTradespersonTradespersonIDFixedPriceReviewsParams, principal interface{}) middleware.Responder {
	tradespersonID := params.TradespersonID
	token := params.HTTPRequest.Header.Get("Authorization")

	response := operations.NewGetTradespersonTradespersonIDFixedPriceReviewsOK()
	reviews := []*operations.GetTradespersonTradespersonIDFixedPriceReviewsOKBodyItems0{}
	response.SetPayload(reviews)

	valid, err := ValidateTradespersonAccessToken(tradespersonID, token)
	if err != nil {
		log.Printf("Failed to validate tradesperson %s, accessToken %s", tradespersonID, token)
		return response
	} else if !valid {
		log.Printf("Bad actor tradesperson %s, accessToken %s", tradespersonID, token)
		return response
	}

	db := database.GetConnection()

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
	token := params.HTTPRequest.Header.Get("Authorization")

	response := operations.NewPostTradespersonTradespersonIDFixedPriceReviewOK()
	payload := operations.PostTradespersonTradespersonIDFixedPriceReviewOKBody{Responded: false}
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
	token := params.HTTPRequest.Header.Get("Authorization")

	response := operations.NewPostTradespersonTradespersonIDQuoteCreated()
	payload := &operations.PostTradespersonTradespersonIDQuoteCreatedBody{}
	created := false
	payload.Created = created
	response.SetPayload(payload)

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
	token := params.HTTPRequest.Header.Get("Authorization")

	response := operations.NewGetTradespersonTradespersonIDQuoteQuoteIDOK()

	valid, err := ValidateTradespersonAccessToken(tradespersonID, token)
	if err != nil {
		log.Printf("Failed to validate tradesperson %s, accessToken %s", tradespersonID, token)
		return response
	} else if !valid {
		log.Printf("Bad actor tradesperson %s, accessToken %s", tradespersonID, token)
		return response
	}

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
	token := params.HTTPRequest.Header.Get("Authorization")

	response := operations.NewPutTradespersonTradespersonIDQuoteQuoteIDOK()
	payload := operations.PutTradespersonTradespersonIDQuoteQuoteIDOKBody{}
	payload.Updated = false

	valid, err := ValidateTradespersonAccessToken(tradespersonID, token)
	if err != nil {
		log.Printf("Failed to validate tradesperson %s, accessToken %s", tradespersonID, token)
		return response
	} else if !valid {
		log.Printf("Bad actor tradesperson %s, accessToken %s", tradespersonID, token)
		return response
	}

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
	token := params.HTTPRequest.Header.Get("Authorization")

	response := operations.NewGetTradespersonTradespersonIDQuoteReviewsOK()
	reviews := []*operations.GetTradespersonTradespersonIDQuoteReviewsOKBodyItems0{}
	response.SetPayload(reviews)

	valid, err := ValidateTradespersonAccessToken(tradespersonID, token)
	if err != nil {
		log.Printf("Failed to validate tradesperson %s, accessToken %s", tradespersonID, token)
		return response
	} else if !valid {
		log.Printf("Bad actor tradesperson %s, accessToken %s", tradespersonID, token)
		return response
	}

	db := database.GetConnection()

	stmt, err := db.Prepare("SELECT qr.id, qr.customerId, qr.rating, qr.message, DATE_FORMAT(qr.date, '%M %D %Y') date FROM quote_reviews qr INNER JOIN quotes q ON qr.quoteId=q.id WHERE q.tradespersonId=? AND qr.responded=0")
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
			log.Printf("Failed to scan row, %s", err)
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
}

func PostTradespersonTradespersonIDQuoteReviewHandler(params operations.PostTradespersonTradespersonIDQuoteReviewParams, principal interface{}) middleware.Responder {
	tradespersonID := params.TradespersonID
	respMsg := *params.Review.Response
	reviewID := *params.Review.ReviewID
	token := params.HTTPRequest.Header.Get("Authorization")

	response := operations.NewPostTradespersonTradespersonIDQuoteReviewOK()
	payload := operations.PostTradespersonTradespersonIDQuoteReviewOKBody{Responded: false}
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

	stmt, err := db.Prepare("SELECT qr.id FROM quote_reviews qr INNER JOIN quotes q ON qr.quoteId=q.id WHERE q.tradespersonId=? AND qr.responded=0 AND qr.id=?")
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
		stmt, err := db.Prepare("UPDATE quote_reviews SET responded = 1, respMsg = ?, respDate = NOW() WHERE id = ?")
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

func GetTradespersonTradespersonIDQuotesHandler(params operations.GetTradespersonTradespersonIDQuotesParams, principal interface{}) middleware.Responder {
	tradespersonID := params.TradespersonID
	page := params.Page
	response := operations.NewGetTradespersonTradespersonIDQuotesOK()
	token := params.HTTPRequest.Header.Get("Authorization")

	valid, err := ValidateTradespersonAccessToken(tradespersonID, token)
	if err != nil {
		log.Printf("Failed to validate tradesperson %s, accessToken %s", tradespersonID, token)
		return response
	} else if !valid {
		log.Printf("Bad actor tradesperson %s, accessToken %s", tradespersonID, token)
		return response
	}

	services := database.GetTradespersonQuotes(tradespersonID, page)
	response.SetPayload(services)

	return response
}

func GetTradespersonTradespersonIDQuotePagesHandler(params operations.GetTradespersonTradespersonIDQuotePagesParams, principal interface{}) middleware.Responder {
	tradespersonID := params.TradespersonID
	token := params.HTTPRequest.Header.Get("Authorization")

	pages := float64(1)
	response := operations.NewGetTradespersonTradespersonIDQuotePagesOK().WithPayload(int64(pages))

	valid, err := ValidateTradespersonAccessToken(tradespersonID, token)
	if err != nil {
		log.Printf("Failed to validate tradesperson %s, accessToken %s", tradespersonID, token)
		return response
	} else if !valid {
		log.Printf("Bad actor tradesperson %s, accessToken %s", tradespersonID, token)
		return response
	}

	db := database.GetConnection()

	stmt, err := db.Prepare("SELECT COUNT(*) FROM quotes WHERE tradespersonId=?")
	if err != nil {
		log.Printf("Failed to create select statement %s", err)
		return response
	}
	defer stmt.Close()

	err = stmt.QueryRow(tradespersonID).Scan(&pages)
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
