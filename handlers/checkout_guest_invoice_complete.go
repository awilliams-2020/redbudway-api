package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"redbudway-api/database"
	"redbudway-api/internal"

	"github.com/stripe/stripe-go/v82"
	"github.com/stripe/stripe-go/v82/checkout/session"
	"github.com/stripe/stripe-go/v82/invoice"
)

// postCheckoutGuestInvoiceCompleteBody links a successful guest Checkout PaymentIntent to the billing invoice (partial payment).
type postCheckoutGuestInvoiceCompleteBody struct {
	TradespersonID    string `json:"tradespersonId"`
	StripeInvoiceID   string `json:"stripeInvoiceId"`
	CheckoutSessionID string `json:"checkoutSessionId"`
	Token             string `json:"token"`
}

type postCheckoutGuestInvoiceCompleteResponse struct {
	Attached        bool   `json:"attached"`
	Message         string `json:"message,omitempty"`
	TradespersonID  string `json:"tradespersonId,omitempty"`
	StripeInvoiceID string `json:"stripeInvoiceId,omitempty"`
	StripeQuoteID   string `json:"stripeQuoteId,omitempty"`
	// DepositPct is 0–100 from the quote; -1 if unknown (aligned with guest-accept-quote).
	DepositPct      int64  `json:"depositPct"`
	ProviderEmail   string `json:"providerEmail,omitempty"`
	ServiceName     string `json:"serviceName,omitempty"`
	Request         string `json:"request,omitempty"`
	AmountPaid      *int64 `json:"amountPaid,omitempty"`
	AmountRemaining *int64 `json:"amountRemaining,omitempty"`
	AmountTotal     *int64 `json:"amountTotal,omitempty"`
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

	tpStripeID, _ := database.GetTradespersonStripeID(body.TradespersonID)
	sessParams := &stripe.CheckoutSessionParams{
		Params: stripe.Params{
			Expand: []*string{
				stripe.String("payment_intent"),
			},
		},
	}
	if tpStripeID != "" {
		sessParams.SetStripeAccount(tpStripeID)
	}
	sess, err := session.Get(body.CheckoutSessionID, sessParams)
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
		var inv *stripe.Invoice
		if invTry, invErr := guestConnectInvoiceGetExpand(tpStripeID, body.StripeInvoiceID, []*string{stripe.String("parent.quote_details")}); invErr == nil {
			inv = invTry
		}
		resp := guestInvoiceCompleteResponseWithAmounts(false, "no invoice attach needed", inv, body.TradespersonID, body.StripeInvoiceID)
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(resp)
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

	expandPI := []*string{
		stripe.String("payments.data.payment.payment_intent"),
		stripe.String("parent.quote_details"),
	}
	invLoaded, err := guestConnectInvoiceGetExpand(tpStripeID, body.StripeInvoiceID, expandPI)
	if err != nil {
		log.Printf("guest invoice complete invoice.Get %s: %v", body.StripeInvoiceID, err)
		w.WriteHeader(http.StatusBadGateway)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "could not load invoice"})
		return
	}
	if guestInvoicePaymentIntentAttached(invLoaded, piID) {
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(guestInvoiceCompleteResponseWithAmounts(true, "payment already linked to invoice", invLoaded, body.TradespersonID, body.StripeInvoiceID))
		return
	}

	attachParams := &stripe.InvoiceAttachPaymentParams{
		PaymentIntent: stripe.String(piID),
	}
	if tpStripeID != "" {
		attachParams.SetStripeAccount(tpStripeID)
	}
	if _, err := invoice.AttachPayment(body.StripeInvoiceID, attachParams); err != nil {
		log.Printf("guest invoice complete AttachPayment %s invoice %s: %v", piID, body.StripeInvoiceID, err)
		// Race: another request may have attached between Get and Attach.
		if inv2, err2 := guestConnectInvoiceGetExpand(tpStripeID, body.StripeInvoiceID, expandPI); err2 == nil && guestInvoicePaymentIntentAttached(inv2, piID) {
			w.WriteHeader(http.StatusOK)
			_ = json.NewEncoder(w).Encode(guestInvoiceCompleteResponseWithAmounts(true, "payment linked to invoice", inv2, body.TradespersonID, body.StripeInvoiceID))
			return
		}
		w.WriteHeader(http.StatusBadGateway)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "could not link payment to invoice"})
		return
	}

	invFinal, errFinal := guestConnectInvoiceGetExpand(tpStripeID, body.StripeInvoiceID, expandPI)
	if errFinal != nil {
		log.Printf("guest invoice complete invoice refresh after attach %s: %v", body.StripeInvoiceID, errFinal)
		// Expand+payments can fail transiently; plain Get still returns paid/remaining for the client UI.
		if invPlain, plainErr := guestConnectInvoiceGetExpand(tpStripeID, body.StripeInvoiceID, []*string{stripe.String("parent.quote_details")}); plainErr == nil {
			w.WriteHeader(http.StatusOK)
			_ = json.NewEncoder(w).Encode(guestInvoiceCompleteResponseWithAmounts(true, "payment linked to invoice", invPlain, body.TradespersonID, body.StripeInvoiceID))
			return
		}
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(guestInvoiceCompleteResponseWithAmounts(true, "payment linked to invoice", nil, body.TradespersonID, body.StripeInvoiceID))
		return
	}
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(guestInvoiceCompleteResponseWithAmounts(true, "payment linked to invoice", invFinal, body.TradespersonID, body.StripeInvoiceID))
}

func guestInvoiceCompleteResponseWithAmounts(attached bool, msg string, inv *stripe.Invoice, tradespersonID, invoiceID string) postCheckoutGuestInvoiceCompleteResponse {
	resp := postCheckoutGuestInvoiceCompleteResponse{Attached: attached, Message: msg, DepositPct: -1}
	if tradespersonID != "" {
		resp.TradespersonID = tradespersonID
	}
	if invoiceID != "" {
		resp.StripeInvoiceID = invoiceID
	}
	if inv != nil {
		ap, ar := inv.AmountPaid, inv.AmountRemaining
		resp.AmountPaid = &ap
		resp.AmountRemaining = &ar
		at := inv.Total
		resp.AmountTotal = &at
		dep, pe, sn, rq := guestQuotePublicFieldsFromInvoice(tradespersonID, inv)
		resp.DepositPct = dep
		resp.ProviderEmail = pe
		resp.ServiceName = sn
		resp.Request = rq
		if inv.Parent != nil && inv.Parent.QuoteDetails != nil {
			if q := strings.TrimSpace(inv.Parent.QuoteDetails.Quote); q != "" {
				resp.StripeQuoteID = q
			}
		}
	}
	return resp
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
