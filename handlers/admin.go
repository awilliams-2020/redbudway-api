package handlers

import (
	"log"
	"redbudway-api/database"
	"redbudway-api/restapi/operations"

	"github.com/go-openapi/runtime/middleware"
)

func GetAdminAdminIDTradespeopleHandler(params operations.GetAdminAdminIDTradespeopleParams, principal interface{}) middleware.Responder {
	adminID := params.AdminID
	token := params.HTTPRequest.Header.Get("Authorization")

	tradespeople := []*operations.GetAdminAdminIDTradespeopleOKBodyItems0{}
	response := operations.NewGetAdminAdminIDTradespeopleOK().WithPayload(tradespeople)

	valid, err := ValidateAdminAccessToken(adminID, token)
	if err != nil {
		log.Printf("Failed to validate admin %s, accessToken %s", adminID, token)
		return response
	} else if !valid {
		log.Printf("Bad actor tradesperson %s, accessToken %s", adminID, token)
		return response
	}

	db := database.GetConnection()

	stmt, err := db.Prepare("SELECT tradespersonId, name, email, image FROM tradesperson_profile")
	if err != nil {
		return response
	}
	defer stmt.Close()

	var tradespersonID, name, email, image string
	rows, err := stmt.Query()
	if err != nil {
		log.Printf("Failed to execute select statement %s", err)
		return response
	}

	for rows.Next() {
		if err := rows.Scan(&tradespersonID, &name, &email, &image); err != nil {
			return response
		}
		tradesperson := operations.GetAdminAdminIDTradespeopleOKBodyItems0{}
		tradesperson.TradespersonID = tradespersonID
		tradesperson.Name = name
		tradesperson.Email = email
		tradesperson.Image = image
		tradespeople = append(tradespeople, &tradesperson)
	}
	response.SetPayload(tradespeople)

	return response
}
