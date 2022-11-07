package handlers

import (
	"log"
	"github.com/go-openapi/runtime/middleware"
	"redbudway-api/database"
	"redbudway-api/restapi/operations"
)

func GetTradespersonTradespersonIDTimeSlotsHandler(params operations.GetTradespersonTradespersonIDTimeSlotsParams, principal interface{}) middleware.Responder {
	tradespersonID := params.TradespersonID

	response, err := database.GetTradespersonTimeslots(tradespersonID)
	if err != nil {
		log.Printf("Failed to retrieve customer schedule %s", err)
		return response
	}

	return response
}