package handlers

import (
	"log"

	"github.com/go-openapi/runtime/middleware"

	"redbudway-api/database"
	"redbudway-api/restapi/operations"
)

func GetTradespersonTradespersonIDScheduleHandler(params operations.GetTradespersonTradespersonIDScheduleParams, principal interface{}) middleware.Responder {
	tradespersonID := params.TradespersonID

	response, err := database.GetTradespersonSchedule(tradespersonID)
	if err != nil {
		log.Printf("Failed to retrieve tradesperson schedule %s", err)
		return response
	}

	return response
}
