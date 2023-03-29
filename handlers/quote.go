package handlers

import (
	"log"
	"redbudway-api/database"
	"redbudway-api/restapi/operations"

	"github.com/go-openapi/runtime/middleware"
)

func GetQuoteQuoteIDHandler(params operations.GetQuoteQuoteIDParams) middleware.Responder {
	quoteID := params.QuoteID

	payload := operations.GetQuoteQuoteIDOKBody{}
	response := operations.NewGetQuoteQuoteIDOK().WithPayload(&payload)
	quote, business, err := database.GetQuoteServiceDetails(quoteID)
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

	reviews := operations.GetQuoteQuoteIDReviewsOKBody{}
	response := operations.NewGetQuoteQuoteIDReviewsOK()

	var err error
	reviews, err = database.GetQuoteRatings(quoteID)
	if err != nil {
		log.Printf("Failed to get quote %s ratings, %v", &quoteID, err)
	}

	reviews.Reviews, err = database.GetQuoteReviews(quoteID, page)
	if err != nil {
		log.Printf("Failed to get quote %s reviews, %v", &quoteID, err)
	}

	response.SetPayload(&reviews)

	return response
}
