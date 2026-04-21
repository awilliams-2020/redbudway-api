package handlers

import (
	"log"
	"net/http"
	"os"
	"strings"

	"redbudway-api/restapi/operations"

	"github.com/go-openapi/runtime/middleware"
	"github.com/resend/resend-go/v2"
)

// GetTradespersonTradespersonIDBillingQuoteQuoteIDNotifyStatusHandler returns the Resend delivery
// status for a previously sent quote-update notification email.
func GetTradespersonTradespersonIDBillingQuoteQuoteIDNotifyStatusHandler(params operations.GetTradespersonTradespersonIDBillingQuoteQuoteIDNotifyStatusParams, principal interface{}) middleware.Responder {
	tradespersonID := params.TradespersonID
	token := params.HTTPRequest.Header.Get("Authorization")

	valid, err := ValidateTradespersonAccessToken(tradespersonID, token)
	if err != nil {
		log.Printf("notify status: validate tradesperson %s: %v", tradespersonID, err)
		return middleware.Error(http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
	}
	if !valid {
		return middleware.Error(http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
	}

	key := strings.TrimSpace(os.Getenv("RESEND_API_KEY"))
	if key == "" {
		return middleware.Error(http.StatusServiceUnavailable, map[string]string{"error": "email status not available"})
	}

	client := resend.NewClient(key)
	emailRecord, err := client.Emails.Get(params.EmailID)
	if err != nil {
		log.Printf("notify status Emails.Get %s: %v", params.EmailID, err)
		return middleware.Error(http.StatusBadGateway, map[string]string{"error": "unable to retrieve email status"})
	}

	payload := &operations.GetTradespersonTradespersonIDBillingQuoteQuoteIDNotifyStatusOKBody{
		LastEvent: emailRecord.LastEvent,
	}
	return operations.NewGetTradespersonTradespersonIDBillingQuoteQuoteIDNotifyStatusOK().WithPayload(payload)
}
