package handlers

import (
	"log"

	"github.com/go-openapi/runtime/middleware"

	"redbudway-api/database"
	"redbudway-api/restapi/operations"
)

func GetTradespersonTradespersonIDScheduleHandler(params operations.GetTradespersonTradespersonIDScheduleParams, principal interface{}) middleware.Responder {
	tradespersonID := params.TradespersonID
	accessToken := params.AccessToken
	token := params.HTTPRequest.Header.Get("Authorization")

	response := operations.NewGetTradespersonTradespersonIDScheduleOK()

	valid, err := ValidateTradespersonAccessToken(tradespersonID, token)
	if err != nil {
		log.Printf("Failed to validate tradesperson %s, accessToken %s", tradespersonID, token)
		return response
	} else if !valid {
		log.Printf("Bad actor tradesperson %s, accessToken %s", tradespersonID, token)
		return response
	}

	response, err = database.GetTradespersonSchedule(tradespersonID, accessToken)
	if err != nil {
		log.Printf("Failed to retrieve tradesperson schedule %s", err)
		return response
	}

	return response
}
