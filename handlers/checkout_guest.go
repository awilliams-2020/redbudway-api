package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"net/http"
	"os"
	"strings"
	"unicode/utf8"

	"redbudway-api/database"
	"redbudway-api/email"
	"redbudway-api/internal"
	_stripe "redbudway-api/stripe"

	"github.com/stripe/stripe-go/v82"
	"github.com/stripe/stripe-go/v82/checkout/session"
	"github.com/stripe/stripe-go/v82/customer"
	"github.com/stripe/stripe-go/v82/invoice"
	"github.com/stripe/stripe-go/v82/quote"
)

// Guest checkout: payment is tied to a Stripe Invoice created after the quote is accepted (not the quote itself).
type postCheckoutGuestSessionBody struct {
	TradespersonID  string `json:"tradespersonId"`
	StripeInvoiceID string `json:"stripeInvoiceId"`
	Token           string `json:"token"`
	SuccessPath     string `json:"successPath"`
	CancelPath      string `json:"cancelPath"`
}

type postCheckoutGuestSessionResponse struct {
	ID  string `json:"id"`
	URL string `json:"url"`
}

type postCheckoutGuestAcceptQuoteBody struct {
	TradespersonID string `json:"tradespersonId"`
	StripeQuoteID  string `json:"stripeQuoteId"`
	Token          string `json:"token"`
}

type postCheckoutGuestAcceptQuoteResponse struct {
	TradespersonID  string `json:"tradespersonId"`
	StripeInvoiceID string `json:"stripeInvoiceId"`
	// DepositPct is 0–100 from the quote; -1 if unknown (DB error).
	DepositPct int64 `json:"depositPct"`
	// AlreadyAccepted is true when the quote was accepted before this request (revisit / refresh).
	AlreadyAccepted bool   `json:"alreadyAccepted"`
	ProviderEmail   string `json:"providerEmail,omitempty"`
	ServiceName     string `json:"serviceName,omitempty"`
	Request         string `json:"request,omitempty"`
	// Stripe invoice cents (set when the invoice is loaded after finalize; omit on load failure).
	AmountPaid      *int64 `json:"amountPaid,omitempty"`
	AmountRemaining *int64 `json:"amountRemaining,omitempty"`
	AmountTotal     *int64 `json:"amountTotal,omitempty"`
}

// guestConnectInvoiceGet loads an invoice for a Connect billing flow (try connected account first).
func guestConnectInvoiceGet(connectAccountID, invoiceID string) (*stripe.Invoice, error) {
	if connectAccountID != "" {
		p := &stripe.InvoiceParams{}
		p.SetStripeAccount(connectAccountID)
		inv, err := invoice.Get(invoiceID, p)
		if err == nil {
			return inv, nil
		}
		log.Printf("guest invoice.Get %s with Connect account: %v", invoiceID, err)
	}
	return invoice.Get(invoiceID, nil)
}

// depositInvoiceFooter returns text for the Stripe invoice footer (hosted invoice, PDF, email) so customers
// see payment/deposit expectations on the generated invoice after quote acceptance.
func depositInvoiceFooter(depositPct, amountTotalCents int64) string {
	if amountTotalCents <= 0 {
		return ""
	}
	if depositPct < 0 || depositPct > 100 {
		return ""
	}
	if depositPct == 0 {
		return "Payment terms: The full amount shown on this invoice is due at payment (no separate deposit)."
	}
	depositDue := (amountTotalCents * depositPct) / 100
	balance := amountTotalCents - depositDue
	return fmt.Sprintf(
		"Deposit terms: %d%% of the quoted total is due with this invoice ($%s). Remaining balance on this same invoice after this payment: $%s.",
		depositPct,
		formatCentsForFooter(depositDue),
		formatCentsForFooter(balance),
	)
}

func formatCentsForFooter(cents int64) string {
	if cents < 0 {
		cents = 0
	}
	return fmt.Sprintf("%.2f", float64(cents)/100.0)
}

// checkoutLineItemDepositDescription is a short line for Stripe Checkout when ProductData description is shown.
func checkoutLineItemDepositDescription(depositPct int64) string {
	if depositPct < 0 || depositPct > 100 {
		return ""
	}
	if depositPct == 0 {
		return "Full quoted amount due with this payment (no separate deposit)."
	}
	return fmt.Sprintf("%d%% of the quoted total is due with this payment; the remaining balance stays on this invoice until it is paid.", depositPct)
}

// checkoutLineItemNameSuffix is a compact fragment merged into the Checkout line item name (Stripe often hides ProductData.Description).
func checkoutLineItemNameSuffix(depositPct int64) string {
	if depositPct < 0 || depositPct > 100 {
		return ""
	}
	if depositPct == 0 {
		return "Full amount due now (no separate deposit)"
	}
	return fmt.Sprintf("%d%% deposit now; balance on same invoice", depositPct)
}

const stripeCheckoutLineItemNameMaxBytes = 500

func truncateStripeLineItemName(s string, maxBytes int) string {
	if maxBytes <= 0 || len(s) <= maxBytes {
		return s
	}
	ellipsis := "…"
	if maxBytes <= len(ellipsis) {
		return ellipsis[:maxBytes]
	}
	end := maxBytes - len(ellipsis)
	for end > 0 && !utf8.RuneStart(s[end]) {
		end--
	}
	if end == 0 {
		return ellipsis
	}
	return s[:end] + ellipsis
}

func checkoutLineItemDisplayName(base string, depositPct int64) string {
	base = strings.TrimSpace(base)
	if base == "" {
		base = "Invoice payment"
	}
	suf := checkoutLineItemNameSuffix(depositPct)
	if suf == "" {
		return truncateStripeLineItemName(base, stripeCheckoutLineItemNameMaxBytes)
	}
	sep := " — "
	combined := base + sep + suf
	if len(combined) <= stripeCheckoutLineItemNameMaxBytes {
		return combined
	}
	room := stripeCheckoutLineItemNameMaxBytes - len(sep) - len(suf)
	if room < 12 {
		return truncateStripeLineItemName(combined, stripeCheckoutLineItemNameMaxBytes)
	}
	return truncateStripeLineItemName(base, room) + sep + suf
}

func ensureQuoteInvoiceDepositFooter(tradespersonID, stripeQuoteID string, sq *stripe.Quote) {
	if sq == nil || sq.Invoice == nil || sq.Invoice.ID == "" {
		return
	}
	var depositPct int64
	err := database.GetConnection().QueryRow(
		`SELECT q.depositPct FROM tradesperson_quotes tq INNER JOIN quotes q ON tq.quoteId=q.id WHERE tq.tradespersonId=? AND tq.quote=?`,
		tradespersonID, stripeQuoteID,
	).Scan(&depositPct)
	if err != nil {
		return
	}
	footer := depositInvoiceFooter(depositPct, sq.AmountTotal)
	if footer == "" {
		return
	}
	inv, err := invoice.Get(sq.Invoice.ID, nil)
	if err != nil {
		log.Printf("ensureQuoteInvoiceDepositFooter invoice.Get %s: %v", sq.Invoice.ID, err)
		return
	}
	if strings.TrimSpace(inv.Footer) != "" {
		return
	}
	if _, err := invoice.Update(sq.Invoice.ID, &stripe.InvoiceParams{Footer: stripe.String(footer)}); err != nil {
		log.Printf("ensureQuoteInvoiceDepositFooter invoice.Update footer %s: %v", sq.Invoice.ID, err)
	}
}

// AfterQuoteAcceptedEnsuresInvoiceDepositNote finalizes a draft invoice if needed and applies payment/deposit footer text.
// Use after quote.Accept (e.g. logged-in customer flow) so the generated invoice matches guest accept behavior.
func AfterQuoteAcceptedEnsuresInvoiceDepositNote(tradespersonID, stripeQuoteID string) {
	expand := &stripe.QuoteParams{
		Params: stripe.Params{
			Expand: []*string{stripe.String("invoice")},
		},
	}
	sq, err := quote.Get(stripeQuoteID, expand)
	if err != nil {
		log.Printf("AfterQuoteAccepted quote.Get %s: %v", stripeQuoteID, err)
		return
	}
	if sq.Invoice == nil || sq.Invoice.ID == "" {
		return
	}
	inv, err := invoice.Get(sq.Invoice.ID, nil)
	if err != nil {
		log.Printf("AfterQuoteAccepted invoice.Get %s: %v", sq.Invoice.ID, err)
		return
	}
	if _, err := finalizeInvoiceIfDraft(inv); err != nil {
		log.Printf("AfterQuoteAccepted finalizeInvoiceIfDraft %s: %v", sq.Invoice.ID, err)
	}
	sq, err = quote.Get(stripeQuoteID, expand)
	if err != nil {
		log.Printf("AfterQuoteAccepted quote.Get refresh %s: %v", stripeQuoteID, err)
		return
	}
	ensureQuoteInvoiceDepositFooter(tradespersonID, stripeQuoteID, sq)
}

func notifyTradespersonQuoteAcceptedGuest(tradespersonID, stripeQuoteID string, sq *stripe.Quote) {
	if sq == nil || sq.Customer == nil || sq.Customer.ID == "" {
		return
	}
	var catalogQuoteSlug, message string
	err := database.GetConnection().QueryRow(
		`SELECT q.quote, tq.request FROM tradesperson_quotes tq INNER JOIN quotes q ON q.id=tq.quoteId WHERE tq.tradespersonId=? AND tq.quote=?`,
		tradespersonID, stripeQuoteID,
	).Scan(&catalogQuoteSlug, &message)
	if err != nil {
		log.Printf("guest accept notify lookup: %v", err)
		return
	}
	tradesperson, err := database.GetTradespersonProfile(tradespersonID)
	if err != nil {
		log.Printf("guest accept GetTradespersonProfile: %v", err)
		return
	}
	_quote, err := database.GetTradespersonQuote(tradespersonID, catalogQuoteSlug)
	if err != nil {
		log.Printf("guest accept GetTradespersonQuote: %v", err)
		return
	}
	stripeCustomer, err := customer.Get(sq.Customer.ID, nil)
	if err != nil {
		log.Printf("guest accept customer.Get: %v", err)
		return
	}
	if err := email.SendTradespersonQuoteAccepted(tradesperson, stripeCustomer, message, _quote); err != nil {
		log.Printf("guest accept SendTradespersonQuoteAccepted: %v", err)
	}
}

// PostCheckoutGuestAcceptQuoteHTTP accepts an open Stripe billing quote (email link flow), creates the invoice, then the client redirects to guest invoice pay.
func PostCheckoutGuestAcceptQuoteHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "method not allowed"})
		return
	}

	var body postCheckoutGuestAcceptQuoteBody
	if err := json.NewDecoder(http.MaxBytesReader(w, r.Body, 65536)).Decode(&body); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "invalid JSON"})
		return
	}

	body.TradespersonID = strings.TrimSpace(body.TradespersonID)
	body.StripeQuoteID = strings.TrimSpace(body.StripeQuoteID)
	body.Token = strings.TrimSpace(body.Token)

	if body.TradespersonID == "" || body.StripeQuoteID == "" || body.Token == "" {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "tradespersonId, stripeQuoteId, and token are required"})
		return
	}

	ok, err := internal.VerifyReCaptcha(body.Token)
	if err != nil || !ok {
		if err != nil {
			log.Printf("guest accept reCAPTCHA error: %v", err)
		}
		w.WriteHeader(http.StatusUnauthorized)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "reCAPTCHA verification failed"})
		return
	}

	dbTpID, err := database.GetTradespersonIDByBillingStripeQuote(body.StripeQuoteID)
	if err != nil {
		if err == sql.ErrNoRows {
			w.WriteHeader(http.StatusNotFound)
			_ = json.NewEncoder(w).Encode(map[string]string{"error": "quote not found for this provider"})
			return
		}
		log.Printf("guest accept GetTradespersonIDByBillingStripeQuote: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "database error"})
		return
	}
	if dbTpID != body.TradespersonID {
		w.WriteHeader(http.StatusNotFound)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "quote does not match provider"})
		return
	}

	tpStripeID, err := database.GetTradespersonStripeID(body.TradespersonID)
	if err != nil || tpStripeID == "" {
		w.WriteHeader(http.StatusNotFound)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "provider not found"})
		return
	}

	expand := &stripe.QuoteParams{
		Params: stripe.Params{
			Expand: []*string{
				stripe.String("customer"),
				stripe.String("transfer_data.destination"),
				stripe.String("on_behalf_of"),
				stripe.String("invoice"),
			},
		},
	}

	sq, err := quote.Get(body.StripeQuoteID, expand)
	if err != nil {
		log.Printf("guest accept quote.Get: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "unable to load quote"})
		return
	}

	if !quoteConnectMatchesProvider(sq, tpStripeID) {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "quote does not belong to this Connect account"})
		return
	}

	var didAccept bool
	switch sq.Status {
	case stripe.QuoteStatusOpen:
		sq, err = quote.Accept(body.StripeQuoteID, nil)
		if err != nil {
			log.Printf("guest accept quote.Accept: %v", err)
			w.WriteHeader(http.StatusBadGateway)
			_ = json.NewEncoder(w).Encode(map[string]string{"error": "unable to accept quote"})
			return
		}
		didAccept = true
	case stripe.QuoteStatusAccepted:
		didAccept = false
	default:
		if sq.Status == stripe.QuoteStatusCanceled {
			w.WriteHeader(http.StatusConflict)
			_ = json.NewEncoder(w).Encode(map[string]string{"error": "quote is canceled"})
			return
		}
		w.WriteHeader(http.StatusConflict)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "quote is not open for acceptance"})
		return
	}

	sq, err = quote.Get(body.StripeQuoteID, expand)
	if err != nil {
		log.Printf("guest accept quote.Get after accept: %v", err)
		w.WriteHeader(http.StatusBadGateway)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "unable to load quote after accept"})
		return
	}

	if sq.Status != stripe.QuoteStatusAccepted {
		w.WriteHeader(http.StatusConflict)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "quote was not accepted"})
		return
	}

	if sq.Invoice == nil || sq.Invoice.ID == "" {
		w.WriteHeader(http.StatusBadGateway)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "invoice not available yet"})
		return
	}

	invAccept, err := guestConnectInvoiceGet(tpStripeID, sq.Invoice.ID)
	if err != nil {
		log.Printf("guest accept invoice.Get %s: %v", sq.Invoice.ID, err)
		w.WriteHeader(http.StatusBadGateway)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "invoice not available"})
		return
	}
	if _, err = finalizeInvoiceIfDraft(invAccept); err != nil {
		log.Printf("guest accept finalizeInvoiceIfDraft %s: %v", sq.Invoice.ID, err)
		w.WriteHeader(http.StatusBadGateway)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "could not finalize invoice for payment"})
		return
	}

	var amountPaid, amountRemaining, amountTotal *int64
	if invLatest, errAmt := guestConnectInvoiceGet(tpStripeID, sq.Invoice.ID); errAmt != nil {
		log.Printf("guest accept invoice amount refresh %s: %v", sq.Invoice.ID, errAmt)
	} else {
		ap := invLatest.AmountPaid
		ar := invLatest.AmountRemaining
		at := invLatest.Total
		amountPaid = &ap
		amountRemaining = &ar
		amountTotal = &at
	}

	ensureQuoteInvoiceDepositFooter(body.TradespersonID, body.StripeQuoteID, sq)

	if didAccept {
		notifyTradespersonQuoteAcceptedGuest(body.TradespersonID, body.StripeQuoteID, sq)
	}

	var depositPct int64
	if err := database.GetConnection().QueryRow(
		`SELECT q.depositPct FROM tradesperson_quotes tq INNER JOIN quotes q ON tq.quoteId=q.id WHERE tq.tradespersonId=? AND tq.quote=?`,
		body.TradespersonID, body.StripeQuoteID,
	).Scan(&depositPct); err != nil {
		depositPct = -1
	}

	pageDetails, pdErr := database.GetGuestAcceptQuotePageDetails(body.TradespersonID, body.StripeQuoteID)
	if pdErr != nil {
		log.Printf("guest accept GetGuestAcceptQuotePageDetails: %v", pdErr)
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(postCheckoutGuestAcceptQuoteResponse{
		TradespersonID:  body.TradespersonID,
		StripeInvoiceID: sq.Invoice.ID,
		DepositPct:      depositPct,
		AlreadyAccepted: !didAccept,
		ProviderEmail:   pageDetails.ProviderEmail,
		ServiceName:     pageDetails.ServiceName,
		Request:         pageDetails.Request,
		AmountPaid:      amountPaid,
		AmountRemaining: amountRemaining,
		AmountTotal:     amountTotal,
	})
}

func baseAppOrigin() string {
	sub := os.Getenv("SUBDOMAIN")
	return "https://" + sub + "redbudway.com"
}

func joinOriginAndPath(origin, path string) string {
	path = strings.TrimSpace(path)
	if path == "" {
		return origin + "/"
	}
	if strings.HasPrefix(path, "http://") || strings.HasPrefix(path, "https://") {
		return path
	}
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}
	return origin + path
}

func guestCheckoutApplicationFeeCents(amountCents int64, sellingFee float64) int64 {
	if amountCents <= 0 || sellingFee <= 0 {
		return 0
	}
	decimalAmount := float64(amountCents) / 100.0
	appFee := decimalAmount * sellingFee
	return int64(math.Floor(appFee * 100))
}

func quoteConnectMatchesProvider(sq *stripe.Quote, tpStripeID string) bool {
	if sq.TransferData != nil && sq.TransferData.Destination != nil && sq.TransferData.Destination.ID != "" {
		return sq.TransferData.Destination.ID == tpStripeID
	}
	if sq.OnBehalfOf != nil && sq.OnBehalfOf.ID != "" {
		return sq.OnBehalfOf.ID == tpStripeID
	}
	return false
}

// finalizeInvoiceIfDraft turns a post–accept quote invoice from draft → open. Without this, customers see
// "invoice is not open for payment" until the provider finalizes in the dashboard or app.
func finalizeInvoiceIfDraft(inv *stripe.Invoice) (*stripe.Invoice, error) {
	if inv.Status != stripe.InvoiceStatusDraft {
		return inv, nil
	}
	if _, err := invoice.FinalizeInvoice(inv.ID, nil); err != nil {
		return nil, err
	}
	return invoice.Get(inv.ID, &stripe.InvoiceParams{
		Params: stripe.Params{
			Expand: []*string{
				stripe.String("quote"),
			},
		},
	})
}

// PostCheckoutGuestSessionHTTP creates a payment path for an open Stripe invoice from an accepted quote (verified in DB + Stripe).
func PostCheckoutGuestSessionHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "method not allowed"})
		return
	}

	var body postCheckoutGuestSessionBody
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
			log.Printf("guest checkout reCAPTCHA error: %v", err)
		}
		w.WriteHeader(http.StatusUnauthorized)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "reCAPTCHA verification failed"})
		return
	}

	inv, err := invoice.Get(body.StripeInvoiceID, &stripe.InvoiceParams{
		Params: stripe.Params{
			Expand: []*string{
				stripe.String("quote"),
			},
		},
	})
	if err != nil {
		log.Printf("guest checkout invoice.Get: %v", err)
		w.WriteHeader(http.StatusNotFound)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "invoice not found"})
		return
	}

	inv, err = finalizeInvoiceIfDraft(inv)
	if err != nil {
		log.Printf("guest checkout finalizeInvoiceIfDraft: %v", err)
		w.WriteHeader(http.StatusBadGateway)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "could not finalize invoice for payment"})
		return
	}

	if inv.Parent == nil || inv.Parent.QuoteDetails == nil || inv.Parent.QuoteDetails.Quote == "" {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "invoice is not from a quote; guest pay is only available after the quote is accepted"})
		return
	}

	quoteStripeID := inv.Parent.QuoteDetails.Quote

	dbTpID, err := database.GetTradespersonIDByBillingStripeQuote(quoteStripeID)
	if err != nil {
		if err == sql.ErrNoRows {
			w.WriteHeader(http.StatusNotFound)
			_ = json.NewEncoder(w).Encode(map[string]string{"error": "quote not found for this provider"})
			return
		}
		log.Printf("guest checkout GetTradespersonIDByBillingStripeQuote: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "database error"})
		return
	}
	if dbTpID != body.TradespersonID {
		w.WriteHeader(http.StatusNotFound)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "invoice does not match provider"})
		return
	}

	sq, err := quote.Get(quoteStripeID, &stripe.QuoteParams{
		Params: stripe.Params{
			Expand: []*string{
				stripe.String("transfer_data.destination"),
				stripe.String("on_behalf_of"),
				stripe.String("invoice"),
			},
		},
	})
	if err != nil {
		log.Printf("guest checkout quote.Get: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "unable to load Stripe quote"})
		return
	}

	if sq.Status != stripe.QuoteStatusAccepted {
		w.WriteHeader(http.StatusConflict)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "quote must be accepted before the customer can pay the invoice"})
		return
	}

	if sq.Invoice == nil || sq.Invoice.ID != inv.ID {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "invoice does not match this quote"})
		return
	}

	tpStripeID, err := database.GetTradespersonStripeID(body.TradespersonID)
	if err != nil || tpStripeID == "" {
		w.WriteHeader(http.StatusNotFound)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "provider not found"})
		return
	}

	acct, err := _stripe.GetConnectAccount(tpStripeID)
	if err != nil {
		log.Printf("guest checkout GetConnectAccount: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "unable to load Connect account"})
		return
	}
	if !acct.ChargesEnabled {
		w.WriteHeader(http.StatusConflict)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "provider cannot accept charges yet"})
		return
	}

	if !quoteConnectMatchesProvider(sq, tpStripeID) {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "quote does not belong to this Connect account"})
		return
	}

	if inv.Status == stripe.InvoiceStatusPaid || inv.AmountRemaining <= 0 {
		w.WriteHeader(http.StatusConflict)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "invoice is already paid"})
		return
	}

	if inv.Status != stripe.InvoiceStatusOpen {
		w.WriteHeader(http.StatusConflict)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "invoice is not open for payment"})
		return
	}

	amountCents := inv.AmountRemaining
	if amountCents < 50 {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "amount due is too small"})
		return
	}

	if u := strings.TrimSpace(inv.HostedInvoiceURL); u != "" {
		ensureQuoteInvoiceDepositFooter(body.TradespersonID, quoteStripeID, sq)
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(postCheckoutGuestSessionResponse{ID: inv.ID, URL: u})
		return
	}

	currency := strings.ToLower(string(inv.Currency))
	if currency == "" {
		currency = "usd"
	}

	desc := strings.TrimSpace(inv.Description)
	if desc == "" {
		desc = "Invoice payment"
	}

	var depositPct int64
	if err := database.GetConnection().QueryRow(
		`SELECT q.depositPct FROM tradesperson_quotes tq INNER JOIN quotes q ON tq.quoteId=q.id WHERE tq.tradespersonId=? AND tq.quote=?`,
		body.TradespersonID, quoteStripeID,
	).Scan(&depositPct); err != nil {
		depositPct = -1
	}

	ensureQuoteInvoiceDepositFooter(body.TradespersonID, quoteStripeID, sq)

	sellingFee, err := database.GetTradespersonSellingFee(body.TradespersonID)
	if err != nil {
		log.Printf("guest checkout selling fee: %v", err)
		sellingFee = 0.06
	}
	appFeeCents := guestCheckoutApplicationFeeCents(amountCents, sellingFee)
	if appFeeCents >= amountCents {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "fee configuration prevents checkout"})
		return
	}

	origin := baseAppOrigin()
	successURL := joinOriginAndPath(origin, body.SuccessPath)
	cancelURL := joinOriginAndPath(origin, body.CancelPath)
	if !strings.Contains(successURL, "{CHECKOUT_SESSION_ID}") {
		if strings.Contains(successURL, "?") {
			successURL = successURL + "&session_id={CHECKOUT_SESSION_ID}"
		} else {
			successURL = successURL + "?session_id={CHECKOUT_SESSION_ID}"
		}
	}

	productData := &stripe.CheckoutSessionLineItemPriceDataProductDataParams{
		Name: stripe.String(checkoutLineItemDisplayName(desc, depositPct)),
	}
	if depLine := checkoutLineItemDepositDescription(depositPct); depLine != "" {
		productData.Description = stripe.String(depLine)
	}

	piDesc := desc
	if depLine := checkoutLineItemDepositDescription(depositPct); depLine != "" {
		piDesc = fmt.Sprintf("%s — %s", desc, depLine)
	}
	piDesc = truncateStripeLineItemName(piDesc, 1000)

	params := &stripe.CheckoutSessionParams{
		Mode:       stripe.String(string(stripe.CheckoutSessionModePayment)),
		SuccessURL: stripe.String(successURL),
		CancelURL:  stripe.String(cancelURL),
		LineItems: []*stripe.CheckoutSessionLineItemParams{
			{
				Quantity: stripe.Int64(1),
				PriceData: &stripe.CheckoutSessionLineItemPriceDataParams{
					Currency:    stripe.String(currency),
					ProductData: productData,
					UnitAmount:  stripe.Int64(amountCents),
				},
			},
		},
		PaymentIntentData: &stripe.CheckoutSessionPaymentIntentDataParams{
			Description:          stripe.String(piDesc),
			ApplicationFeeAmount: stripe.Int64(appFeeCents),
			OnBehalfOf:           stripe.String(tpStripeID),
			TransferData: &stripe.CheckoutSessionPaymentIntentDataTransferDataParams{
				Destination: stripe.String(tpStripeID),
			},
			Metadata: map[string]string{
				"tradesperson_id":   body.TradespersonID,
				"stripe_invoice_id": inv.ID,
				"stripe_quote_id":   quoteStripeID,
				"guest_pay":         "true",
			},
		},
	}
	params.AddMetadata("tradesperson_id", body.TradespersonID)
	params.AddMetadata("stripe_invoice_id", inv.ID)
	params.AddMetadata("stripe_quote_id", quoteStripeID)
	params.AddMetadata("guest_pay", "true")

	sess, err := session.New(params)
	if err != nil {
		log.Printf("guest checkout session.New: %v", err)
		w.WriteHeader(http.StatusBadGateway)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "stripe error creating checkout session"})
		return
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(postCheckoutGuestSessionResponse{ID: sess.ID, URL: sess.URL})
}
