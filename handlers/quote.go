package handlers

import (
	"log"
	"redbudway-api/database"
	"redbudway-api/restapi/operations"

	"github.com/go-openapi/runtime/middleware"
)

func GetQuoteQuoteIDHandler(params operations.GetQuoteQuoteIDParams) middleware.Responder {
	quoteID := params.QuoteID
	state := params.State
	city := params.City

	payload := operations.GetQuoteQuoteIDOKBody{}
	response := operations.NewGetQuoteQuoteIDOK()
	quote, business, err := database.GetQuoteServiceDetails(quoteID, state, city)
	if err != nil {
		log.Printf("Failed to get quote, %s", err)
		return response
	}
	payload.Service = quote
	payload.Business = business
	response.SetPayload(&payload)

	return response
}

func GetQuoteQuoteIDReviewsHandler(params operations.GetQuoteQuoteIDReviewsParams) middleware.Responder {
	quoteID := params.QuoteID
	page := params.Page

	reviews := []*operations.GetQuoteQuoteIDReviewsOKBodyItems0{}
	response := operations.NewGetQuoteQuoteIDReviewsOK()
	db := database.GetConnection()

	stmt, err := db.Prepare("SELECT qr.customerId, qr.rating, qr.message, DATE_FORMAT(qr.date, '%M %D %Y') date, qr.responded, qr.response, DATE_FORMAT(qr.respDate, '%M %D %Y') respDate, q.tradespersonId FROM quote_reviews qr INNER JOIN quotes q ON qr.quoteId=q.id WHERE q.quote=? ORDER BY qr.date DESC LIMIT ?, 10")
	if err != nil {
		log.Printf("Failed to create select statement %s", err)
		return response
	}
	defer stmt.Close()

	rows, err := stmt.Query(quoteID, page)
	if err != nil {
		log.Printf("Failed to execute select statement %s", err)
		return response
	}

	var customerID, message, date, respMsg, respDate string
	var rating int64
	var responded bool
	for rows.Next() {
		if err := rows.Scan(&customerID, &rating, &message, &date, &responded, &respMsg, &respDate); err != nil {
			log.Printf("Failed to scan for quote reviews, %s", err)
			return response
		}

		review := &operations.GetQuoteQuoteIDReviewsOKBodyItems0{}
		review.Rating = rating
		review.Message = message
		review.Date = date
		review.Responded = responded
		review.RespMsg = respMsg
		review.RespDate = respDate
		reviews = append(reviews, review)
	}

	response.SetPayload(reviews)

	return response
}
