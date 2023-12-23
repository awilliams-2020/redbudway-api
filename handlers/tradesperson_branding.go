package handlers

import (
	"log"
	"redbudway-api/database"
	"redbudway-api/internal"
	"redbudway-api/restapi/operations"
	"redbudway-api/stripe"
	"strings"

	"github.com/go-openapi/runtime/middleware"
)

func GetTradespersonTradespersonIDBrandingHandler(params operations.GetTradespersonTradespersonIDBrandingParams, principal interface{}) middleware.Responder {
	tradespersonID := params.TradespersonID
	token := params.HTTPRequest.Header.Get("Authorization")

	response := operations.NewGetTradespersonTradespersonIDBrandingOK()

	valid, err := ValidateTradespersonAccessToken(tradespersonID, token)
	if err != nil {
		log.Printf("Failed to validate tradesperson %s, accessToken %s", tradespersonID, token)
		return response
	} else if !valid {
		log.Printf("Bad actor tradesperson %s, accessToken %s", tradespersonID, token)
		return response
	}

	branding, err := database.GetTradespersonBranding(tradespersonID)
	if err != nil {
		log.Printf("Failed to get tradesperson profile %s", err)
	}

	response.SetPayload(branding)

	return response
}

func PutTradespersonTradespersonIDBrandingHandler(params operations.PutTradespersonTradespersonIDBrandingParams, principal interface{}) middleware.Responder {
	tradespersonID := params.TradespersonID
	branding := params.Branding
	token := params.HTTPRequest.Header.Get("Authorization")

	payload := operations.PutTradespersonTradespersonIDBrandingOKBody{}
	payload.Updated = false
	response := operations.NewPutTradespersonTradespersonIDBrandingOK()
	response.SetPayload(&payload)

	valid, err := ValidateTradespersonAccessToken(tradespersonID, token)
	if err != nil {
		log.Printf("Failed to validate tradesperson %s, accessToken %s", tradespersonID, token)
		return response
	} else if !valid {
		log.Printf("Bad actor tradesperson %s, accessToken %s", tradespersonID, token)
		return response
	}

	var logoURL, iconURL string
	if branding.Logo != "" && !strings.Contains(branding.Logo, "https://") {
		logoURL, err = internal.SaveImage(tradespersonID, branding.Logo, "logo")
		if err != nil {
			log.Printf("Failed to save tradesperson %s, logo", tradespersonID)
			return response
		}
	} else if strings.Contains(branding.Logo, "https://") {
		logoURL = branding.Logo
	}
	if branding.Icon != "" && !strings.Contains(branding.Icon, "https://") {
		iconURL, err = internal.SaveImage(tradespersonID, branding.Icon, "icon")
		if err != nil {
			log.Printf("Failed to save tradesperson %s, icon", tradespersonID)
			return response
		}
	} else if strings.Contains(branding.Icon, "https://") {
		iconURL = branding.Icon
	}

	stripeID, err := database.UpdateTradespersonBranding(tradespersonID, logoURL, iconURL, branding)
	if err != nil {
		log.Printf("Failed to save tradesperson %s, branding, %s", tradespersonID, err)
		return response
	}

	if stripeID == "" {
		stripeID, err = database.GetTradespersonStripeID(tradespersonID)
		if err != nil {
			log.Printf("Failed to get tradesperson %s stripeID, %s", tradespersonID, err)
			return response
		}
	}

	err = stripe.UpdateBusinessBranding(stripeID, logoURL, iconURL, branding.Primary, branding.Secondary, tradespersonID)
	if err != nil {

	}

	payload.Updated = true
	response = operations.NewPutTradespersonTradespersonIDBrandingOK()
	response.SetPayload(&payload)
	return response
}
