package handlers

import (
	"log"
	"os"

	"github.com/go-openapi/runtime/middleware"

	"github.com/stripe/stripe-go/v72/customer"
	"github.com/stripe/stripe-go/v72/price"
	"github.com/stripe/stripe-go/v72/product"

	"redbudway-api/database"
	_email "redbudway-api/email"
	"redbudway-api/restapi/operations"
)

func PostTradespersonTradespersonIDEmailHandler(params operations.PostTradespersonTradespersonIDEmailParams, principal interface{}) middleware.Responder {
	message := params.Email.Message
	images := params.Email.Images
	priceID := params.Email.PriceID
	quoteID := params.Email.QuoteID
	customerID := params.Email.CustomerID
	tradespersonID := params.TradespersonID
	token := params.HTTPRequest.Header.Get("Authorization")

	payload := operations.PostTradespersonTradespersonIDEmailOKBody{}
	sent := false
	payload.Sent = sent
	response := operations.NewPostTradespersonTradespersonIDEmailOK()
	response.SetPayload(&payload)

	valid, err := ValidateCustomerAccessToken(*customerID, token)
	if err != nil {
		log.Printf("Failed to validate tradesperson %s, accessToken %s", tradespersonID, token)
		return response
	} else if !valid {
		log.Printf("Bad actor tradesperson %s, accessToken %s", tradespersonID, token)
		return response
	}

	tradesperson, err := database.GetTradespersonProfile(tradespersonID)
	if err != nil {
		log.Printf("Failed to get tradesperson profile, %v\n", err)
		return response
	}

	cuStripeID, err := database.GetCustomerStripeID(*customerID)
	if err != nil {
		log.Printf("Failed to get customer stripe ID, %v", err)
	}

	stripeCustomer, err := customer.Get(cuStripeID, nil)
	if err != nil {
		return response
	}

	service := ""
	if priceID != "" {
		p, err := price.Get(priceID, nil)
		if err != nil {
			log.Printf("Failed to get stripe price with ID %s, %v", priceID, err)
			return response
		}
		pr, err := product.Get(p.Product.ID, nil)
		if err != nil {
			log.Printf("Failed to get stripe product with ID %s, %v", p.Product.ID, err)
			return response
		}
		service = pr.Name
	} else if quoteID != "" {

	}

	images, err = _email.SendTradespersonMessage(tradesperson.Name, tradesperson.Email, service, *message, stripeCustomer, images)
	if err != nil {
		log.Printf("Failed to send email, %v", err)
	} else {
		sent = true
		payload.Sent = sent
	}

	for _, imagePath := range images {
		err := os.Remove(imagePath)
		if err != nil {
			log.Printf("Failed to delete image, %s", imagePath)
		}
	}
	response.SetPayload(&payload)

	return response
}
