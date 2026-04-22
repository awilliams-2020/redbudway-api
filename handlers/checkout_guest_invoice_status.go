package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"redbudway-api/database"
	"redbudway-api/internal"

	"github.com/stripe/stripe-go/v82"
)

type postCheckoutGuestInvoiceStatusBody struct {
	TradespersonID  string `json:"tradespersonId"`
	StripeInvoiceID string `json:"stripeInvoiceId"`
	Token           string `json:"token"`
}

type postCheckoutGuestInvoiceStatusResponse struct {
	TradespersonID  string `json:"tradespersonId,omitempty"`
	StripeInvoiceID string `json:"stripeInvoiceId,omitempty"`
	StripeQuoteID   string `json:"stripeQuoteId,omitempty"`
	InvoiceStatus   string `json:"invoiceStatus"`
	// DepositPct is 0–100 from the quote; -1 if unknown (same semantics as guest-accept-quote).
	DepositPct      int64  `json:"depositPct"`
	ProviderEmail   string `json:"providerEmail,omitempty"`
	ServiceName     string `json:"serviceName,omitempty"`
	Request         string `json:"request,omitempty"`
	AmountPaid      int64  `json:"amountPaid"`
	AmountRemaining int64  `json:"amountRemaining"`
	AmountTotal     *int64 `json:"amountTotal,omitempty"`
	PaidInFull      bool   `json:"paidInFull"`
}

// PostCheckoutGuestInvoiceStatusHTTP returns current Stripe invoice totals for a public pay link (reCAPTCHA).
func PostCheckoutGuestInvoiceStatusHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "method not allowed"})
		return
	}

	var body postCheckoutGuestInvoiceStatusBody
	if err := json.NewDecoder(http.MaxBytesReader(w, r.Body, 65536)).Decode(&body); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "invalid JSON"})
		return
	}

	body.TradespersonID = strings.TrimSpace(body.TradespersonID)
	body.StripeInvoiceID = strings.TrimSpace(body.StripeInvoiceID)
	body.Token = strings.TrimSpace(body.Token)

	if body.TradespersonID == "" || body.StripeInvoiceID == "" || body.Token == "" {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "tradespersonId, stripeInvoiceId, and token are required"})
		return
	}

	ok, err := internal.VerifyReCaptcha(body.Token)
	if err != nil || !ok {
		if err != nil {
			log.Printf("guest invoice status reCAPTCHA error: %v", err)
		}
		w.WriteHeader(http.StatusUnauthorized)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "reCAPTCHA verification failed"})
		return
	}

	env, loadCode, loadErr := guestCheckoutLoadInvoiceForGuestPay(body.TradespersonID, body.StripeInvoiceID)
	if loadCode != 0 {
		w.WriteHeader(loadCode)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": loadErr})
		return
	}

	inv := env.Inv
	paidInFull := inv.Status == stripe.InvoiceStatusPaid || inv.AmountRemaining <= 0

	pageDetails, pdErr := database.GetGuestAcceptQuotePageDetails(body.TradespersonID, env.QuoteStripe)
	if pdErr != nil {
		log.Printf("guest invoice status GetGuestAcceptQuotePageDetails: %v", pdErr)
	}
	total := inv.Total

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(postCheckoutGuestInvoiceStatusResponse{
		TradespersonID:  body.TradespersonID,
		StripeInvoiceID: body.StripeInvoiceID,
		StripeQuoteID:   env.QuoteStripe,
		InvoiceStatus:   string(inv.Status),
		DepositPct:      env.DepositPct,
		ProviderEmail:   pageDetails.ProviderEmail,
		ServiceName:     pageDetails.ServiceName,
		Request:         pageDetails.Request,
		AmountPaid:      inv.AmountPaid,
		AmountRemaining: inv.AmountRemaining,
		AmountTotal:     &total,
		PaidInFull:      paidInFull,
	})
}
