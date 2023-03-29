package handlers

import (
	"log"
	"redbudway-api/database"
	"redbudway-api/restapi/operations"

	"github.com/go-openapi/runtime/middleware"
)

func GetFixedPricePriceIDHandler(params operations.GetFixedPricePriceIDParams) middleware.Responder {
	priceID := params.PriceID

	payload := operations.GetFixedPricePriceIDOKBody{}
	response := operations.NewGetFixedPricePriceIDOK().WithPayload(&payload)
	fixedPrice, business, err := database.GetFixedPriceServiceDetails(priceID)
	if err != nil {
		log.Printf("Failed to get public fixed price, %s", err)
		return response
	}
	payload.Service = fixedPrice
	payload.Business = business
	response.SetPayload(&payload)
	return response

}

func GetFixedPricePriceIDReviewsHandler(params operations.GetFixedPricePriceIDReviewsParams) middleware.Responder {
	priceID := params.PriceID
	page := params.Page

	reviews := operations.GetFixedPricePriceIDReviewsOKBody{}
	response := operations.NewGetFixedPricePriceIDReviewsOK()

	var err error
	reviews, err = database.GetFixedPriceRatings(priceID)
	if err != nil {
		log.Printf("Failed to get fixed price %s ratings, %v", priceID, err)
	}

	reviews.Reviews, err = database.GetFixedPriceReviews(priceID, page)
	if err != nil {
		log.Printf("Failed to get fixed price %s reviews, %v", priceID, err)
	}

	response.SetPayload(&reviews)

	return response
}
