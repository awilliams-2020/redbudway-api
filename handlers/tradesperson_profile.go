package handlers

import (
	"database/sql"
	"log"
	"redbudway-api/database"
	"redbudway-api/models"
	"redbudway-api/restapi/operations"
	"redbudway-api/stripe"

	"github.com/go-openapi/runtime/middleware"
)

func GetTradespersonTradespersonIDProfileHandler(params operations.GetTradespersonTradespersonIDProfileParams, principal interface{}) middleware.Responder {
	tradespersonID := params.TradespersonID
	token := params.HTTPRequest.Header.Get("Authorization")

	response := operations.NewGetTradespersonTradespersonIDProfileOK()

	valid, err := ValidateTradespersonAccessToken(tradespersonID, token)
	if err != nil {
		log.Printf("Failed to validate tradesperson %s, accessToken %s", tradespersonID, token)
		return response
	} else if !valid {
		log.Printf("Bad actor tradesperson %s, accessToken %s", tradespersonID, token)
		return response
	}

	tradesperson := models.Tradesperson{}

	tradesperson, err = database.GetTradespersonProfile(tradespersonID)
	if err != nil {
		log.Printf("Failed to get tradesperson profile %s", err)
	}

	jobs, err := database.GetTradespersonJobs(tradespersonID)
	if err != nil {
		log.Printf("Failed to get tradesperson job count %s", err)
	}
	tradesperson.Jobs = jobs

	rating, reviews, err := database.GetTradespersonRatingReviews(tradespersonID)
	if err != nil {
		log.Printf("Failed to get tradesperson rating & reviews %s", err)
	}
	tradesperson.Rating = rating
	tradesperson.Reviews = reviews

	response.SetPayload(&tradesperson)

	return response
}

func PutTradespersonTradespersonIDProfileHandler(params operations.PutTradespersonTradespersonIDProfileParams, principal interface{}) middleware.Responder {
	tradespersonID := params.TradespersonID
	tradesperson := params.Tradesperson
	token := params.HTTPRequest.Header.Get("Authorization")

	payload := operations.PutTradespersonTradespersonIDProfileOKBody{}
	payload.Updated = false
	response := operations.NewPutTradespersonTradespersonIDProfileOK()
	response.SetPayload(&payload)

	valid, err := ValidateTradespersonAccessToken(tradespersonID, token)
	if err != nil {
		log.Printf("Failed to validate tradesperson %s, accessToken %s", tradespersonID, token)
		return response
	} else if !valid {
		log.Printf("Bad actor tradesperson %s, accessToken %s", tradespersonID, token)
		return response
	}

	stripeID, err := database.GetTradespersonStripeID(tradespersonID)
	if err != nil {
		log.Printf("Failed to get tradesperson %s stripeID, %s", tradesperson, err)
		return response
	}

	if tradesperson.Name != "" {
		if err := database.UpdateTradespersonProfileName(tradespersonID, tradesperson.Name); err != nil {
			log.Printf("Failed to update tradesperson profile name, %v", err)
		}
		if err := stripe.UpdateBusinessProfileName(stripeID, tradesperson.Name); err != nil {
			log.Printf("Failed to update tradesperson connect business profile name, %v", err)
		}
	}

	if tradesperson.Number != "" {
		if err := database.UpdateTradespersonProfileNumber(tradespersonID, tradesperson.Number); err != nil {
			log.Printf("Failed to update tradesperson profile number, %v", err)
		}
		if err := stripe.UpdateBusinessProfileNumber(stripeID, tradesperson.Number); err != nil {
			log.Printf("Failed to update tradesperson connect business profile number, %v", err)
		}
	}

	if tradesperson.Email != "" {
		if err := database.UpdateTradespersonProfileEmail(tradespersonID, tradesperson.Email); err != nil {
			log.Printf("Failed to update tradesperson profile email, %v", err)
		}
		if err := stripe.UpdateBusinessProfileEmail(stripeID, tradesperson.Email); err != nil {
			log.Printf("Failed to update tradesperson connect business profile email, %v", err)
		}
	}

	if tradesperson.Image != "" {
		if err := database.UpdateTradespersonProfileImage(tradespersonID, tradesperson.Image); err != nil {
			log.Printf("Failed to update tradesperson profile image, %v", err)
		}
	}

	if tradesperson.Description != "" {
		if err := database.UpdateTradespersonProfileDescription(tradespersonID, tradesperson.Description); err != nil {
			log.Printf("Failed to update tradesperson profile description, %v", err)
		}
	}

	if tradesperson.Address != nil {
		if err := database.UpdateTradespersonProfileAddress(tradespersonID, tradesperson.Address); err != nil {
			log.Printf("Failed to update tradesperson profile address, %v", err)
		}
		if err := stripe.UpdateBusinessProfileAddress(stripeID, tradesperson.Address); err != nil {
			log.Printf("Failed to update tradesperson connect business profile address, %v", err)
		}
	}

	payload.Updated = true
	response = operations.NewPutTradespersonTradespersonIDProfileOK()
	response.SetPayload(&payload)
	return response
}

func GetTradespersonTradespersonIDSettingsHandler(params operations.GetTradespersonTradespersonIDSettingsParams, principal interface{}) middleware.Responder {
	tradespersonID := params.TradespersonID
	token := params.HTTPRequest.Header.Get("Authorization")

	payload := operations.GetTradespersonTradespersonIDSettingsOKBody{}
	response := operations.NewGetTradespersonTradespersonIDSettingsOK()
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

	stmt, err := db.Prepare("SELECT vanityURL, number, email, address, timeZone FROM tradesperson_settings WHERE tradespersonId=?")
	if err != nil {
		log.Printf("Failed to create select statement %s", err)
		return response
	}
	defer stmt.Close()

	var vanityURL sql.NullString
	var timeZone string
	var displayNumber, displayEmail, displayAddress bool
	row := stmt.QueryRow(tradespersonID)
	switch err = row.Scan(&vanityURL, &displayNumber, &displayEmail, &displayAddress, &timeZone); err {
	case sql.ErrNoRows:
		log.Printf("Tradesperson with ID %s doesn't exist", tradespersonID)
	case nil:
		if vanityURL.Valid {
			payload.VanityURL = vanityURL.String
		}
		payload.DisplayNumber = displayNumber
		payload.DisplayEmail = displayEmail
		payload.DisplayAddress = displayAddress
		payload.TimeZone = timeZone
		response.SetPayload(&payload)
	default:
		log.Printf("Unknown %v", err)
	}
	return response
}

func PutTradespersonTradespersonIDSettingsHandler(params operations.PutTradespersonTradespersonIDSettingsParams, principal interface{}) middleware.Responder {
	tradespersonID := params.TradespersonID
	settings := params.Settings
	token := params.HTTPRequest.Header.Get("Authorization")

	response := operations.NewPutTradespersonTradespersonIDSettingsOK()

	valid, err := ValidateTradespersonAccessToken(tradespersonID, token)
	if err != nil {
		log.Printf("Failed to validate tradesperson %s, accessToken %s", tradespersonID, token)
		return response
	} else if !valid {
		log.Printf("Bad actor tradesperson %s, accessToken %s", tradespersonID, token)
		return response
	}

	updated := false

	updated, err = database.UpdateTradespersonDisplaySettings(tradespersonID, settings)
	if err != nil {
		log.Printf("Failed to update tradesperson display settings %s", err)
	}
	updated, err = database.UpdateTradespersonVanitySettings(tradespersonID, settings)
	if err != nil {
		log.Printf("Failed to update tradesperson vanity settings %s", err)
	}

	payload := operations.PutTradespersonTradespersonIDSettingsOKBody{Updated: updated}

	response.SetPayload(&payload)

	return response
}

func PutTradespersonTradespersonIDTimeZoneHandler(params operations.PutTradespersonTradespersonIDTimeZoneParams, principal interface{}) middleware.Responder {
	tradespersonID := params.TradespersonID
	settings := params.Settings
	token := params.HTTPRequest.Header.Get("Authorization")

	response := operations.NewPutTradespersonTradespersonIDTimeZoneOK()

	valid, err := ValidateTradespersonAccessToken(tradespersonID, token)
	if err != nil {
		log.Printf("Failed to validate tradesperson %s, accessToken %s", tradespersonID, token)
		return response
	} else if !valid {
		log.Printf("Bad actor tradesperson %s, accessToken %s", tradespersonID, token)
		return response
	}

	updated, err := database.UpdateTradespersonTimeZoneSettings(tradespersonID, settings)
	if err != nil {
		log.Printf("Failed to update tradesperson timezone settings %s", err)
	}

	payload := operations.PutTradespersonTradespersonIDTimeZoneOKBody{Updated: updated}

	response.SetPayload(&payload)

	return response
}
