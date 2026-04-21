package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"redbudway-api/internal"

	"github.com/stripe/stripe-go/v82"
	"github.com/stripe/stripe-go/v82/checkout/session"
	"github.com/stripe/stripe-go/v82/invoice"
)

// postCheckoutGuestInvoiceCompleteBody links a successful guest Checkout PaymentIntent to the billing invoice (partial payment).
type postCheckoutGuestInvoiceCompleteBody struct {
	TradespersonID  string `json:"tradespersonId"`
	StripeInvoiceID string `json:"stripeInvoiceId"`
	CheckoutSessionID string `json:"checkoutSessionId"`
	Token             string `json:"token"`
}

type postCheckoutGuestInvoiceCompleteResponse struct {
	Attached bool   `json:"attached"`
	Message  string `json:"message,omitempty"`
}

// PostCheckoutGuestInvoiceCompleteHTTP calls invoice.AttachPayment for guest Checkout sessions created with guest_invoice_attach metadata.
func PostCheckoutGuestInvoiceCompleteHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "method not allowed"})
		return
	}

	var body postCheckoutGuestInvoiceCompleteBody
	if err := json.NewDecoder(http.MaxBytesReader(w, r.Body, 65536)).Decode(&body); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "invalid JSON"})
		return
	}

	body.TradespersonID = strings.TrimSpace(body.TradespersonID)
	body.StripeInvoiceID = strings.TrimSpace(body.StripeInvoiceID)
	body.CheckoutSessionID = strings.TrimSpace(body.CheckoutSessionID)
	body.Token = strings.TrimSpace(body.Token)

	if body.TradespersonID == "" || body.StripeInvoiceID == "" || body.CheckoutSessionID == "" || body.Token == "" {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "tradespersonId, stripeInvoiceId, checkoutSessionId, and token are required"})
		return
	}

	ok, err := internal.VerifyReCaptcha(body.Token)
	if err != nil || !ok {
		if err != nil {
			log.Printf("guest invoice complete reCAPTCHA error: %v", err)
		}
		w.WriteHeader(http.StatusUnauthorized)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "reCAPTCHA verification failed"})
		return
	}

	sess, err := session.Get(body.CheckoutSessionID, &stripe.CheckoutSessionParams{
		Params: stripe.Params{
			Expand: []*string{
				stripe.String("payment_intent"),
			},
		},
	})
	if err != nil {
		log.Printf("guest invoice complete session.Get: %v", err)
		w.WriteHeader(http.StatusNotFound)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "checkout session not found"})
		return
	}

	if sess.Metadata == nil || sess.Metadata["guest_pay"] != "true" {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "not a guest checkout session"})
		return
	}
	if sess.Metadata["tradesperson_id"] != body.TradespersonID || sess.Metadata["stripe_invoice_id"] != body.StripeInvoiceID {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "session does not match request"})
		return
	}
	if sess.Metadata[guestCheckoutMetadataInvoiceAttach] != "true" {
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(postCheckoutGuestInvoiceCompleteResponse{Attached: false, Message: "no invoice attach needed"})
		return
	}

	if sess.PaymentIntent == nil || sess.PaymentIntent.ID == "" {
		w.WriteHeader(http.StatusConflict)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "payment is not ready yet"})
		return
	}
	if sess.PaymentIntent.Status != stripe.PaymentIntentStatusSucceeded {
		w.WriteHeader(http.StatusConflict)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "payment has not succeeded"})
		return
	}

	piID := sess.PaymentIntent.ID

	invLoaded, err := invoice.Get(body.StripeInvoiceID, &stripe.InvoiceParams{
		Params: stripe.Params{
			Expand: []*string{
				stripe.String("payments.data.payment.payment_intent"),
			},
		},
	})
	if err != nil {
		log.Printf("guest invoice complete invoice.Get %s: %v", body.StripeInvoiceID, err)
		w.WriteHeader(http.StatusBadGateway)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "could not load invoice"})
		return
	}
	if guestInvoicePaymentIntentAttached(invLoaded, piID) {
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(postCheckoutGuestInvoiceCompleteResponse{Attached: true, Message: "payment already linked to invoice"})
		return
	}

	if _, err := invoice.AttachPayment(body.StripeInvoiceID, &stripe.InvoiceAttachPaymentParams{
		PaymentIntent: stripe.String(piID),
	}); err != nil {
		log.Printf("guest invoice complete AttachPayment %s invoice %s: %v", piID, body.StripeInvoiceID, err)
		// Race: another request may have attached between Get and Attach.
		if inv2, err2 := invoice.Get(body.StripeInvoiceID, &stripe.InvoiceParams{
			Params: stripe.Params{
				Expand: []*string{stripe.String("payments.data.payment.payment_intent")},
			},
		}); err2 == nil && guestInvoicePaymentIntentAttached(inv2, piID) {
			w.WriteHeader(http.StatusOK)
			_ = json.NewEncoder(w).Encode(postCheckoutGuestInvoiceCompleteResponse{Attached: true, Message: "payment linked to invoice"})
			return
		}
		w.WriteHeader(http.StatusBadGateway)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "could not link payment to invoice"})
		return
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(postCheckoutGuestInvoiceCompleteResponse{Attached: true, Message: "payment linked to invoice"})
}

func guestInvoicePaymentIntentAttached(inv *stripe.Invoice, paymentIntentID string) bool {
	if inv == nil || paymentIntentID == "" || inv.Payments == nil {
		return false
	}
	for _, row := range inv.Payments.Data {
		if row == nil || row.Payment == nil || row.Payment.PaymentIntent == nil {
			continue
		}
		if row.Payment.PaymentIntent.ID == paymentIntentID {
			return true
		}
	}
	return false
}
