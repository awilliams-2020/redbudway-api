package handlers

import (
	"log"
	"redbudway-api/database"
	"redbudway-api/internal"
	"redbudway-api/models"
	"redbudway-api/restapi/operations"

	"github.com/go-openapi/runtime/middleware"
)

func GetTradespersonTradespersonIDTimeSlotsHandler(params operations.GetTradespersonTradespersonIDTimeSlotsParams, principal interface{}) middleware.Responder {
	tradespersonID := params.TradespersonID
	accessToken := params.AccessToken
	token := params.HTTPRequest.Header.Get("Authorization")

	payload := operations.GetTradespersonTradespersonIDTimeSlotsOKBody{}
	response := operations.NewGetTradespersonTradespersonIDTimeSlotsOK().WithPayload(&payload)

	valid, err := ValidateTradespersonAccessToken(tradespersonID, token)
	if err != nil {
		log.Printf("Failed to validate tradesperson %s, accessToken %s", tradespersonID, token)
		return response
	} else if !valid {
		log.Printf("Bad actor tradesperson %s, accessToken %s", tradespersonID, token)
		return response
	}

	otherServices, err := database.GetOtherServices(tradespersonID, int64(0))
	if err != nil {
		log.Printf("Failed to get other fixed prices, %s", err)
		return response
	}
	payload.OtherServices = otherServices
	response.SetPayload(&payload)

	googleTimeSlots := models.GoogleTimeSlots{}
	if accessToken != nil {
		internal.GetGoogleTimeSlots(*accessToken)
		googleTimeSlots = internal.GetGoogleTimeSlots(*accessToken)
		if err != nil {
			log.Printf("Failed to get google time slots, %s", err)
			return response
		}
	}
	payload.GoogleTimeSlots = googleTimeSlots
	response.SetPayload(&payload)

	return response
}
