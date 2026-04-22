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

// Metadata on guest Checkout Sessions when a partial deposit must be linked to the invoice after payment.
const guestCheckoutMetadataInvoiceAttach = "guest_invoice_attach"

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

// guestCheckoutInvoiceEnv is the loaded Stripe invoice plus ids after the same checks used for guest-session.
type guestCheckoutInvoiceEnv struct {
	Inv         *stripe.Invoice
	TpStripeID  string
	QuoteStripe string
	DepositPct  int64
}

// guestCheckoutLoadInvoiceForGuestPay runs authorization + quote validation through finalize (same path as guest-session).
// On failure returns non-zero HTTP status and an error message for JSON {"error": msg}.
func guestCheckoutLoadInvoiceForGuestPay(tradespersonID, stripeInvoiceID string) (*guestCheckoutInvoiceEnv, int, string) {
	tpStripeID, err := database.GetTradespersonStripeID(tradespersonID)
	if err != nil || tpStripeID == "" {
		return nil, http.StatusNotFound, "provider not found"
	}

	inv, err := guestConnectInvoiceGetExpand(tpStripeID, stripeInvoiceID, []*string{
		stripe.String("quote"),
		stripe.String("customer"),
	})
	if err != nil {
		log.Printf("guest checkout invoice.Get: %v", err)
		return nil, http.StatusNotFound, "invoice not found"
	}

	if inv.Parent == nil || inv.Parent.QuoteDetails == nil || inv.Parent.QuoteDetails.Quote == "" {
		return nil, http.StatusBadRequest, "invoice is not from a quote; guest pay is only available after the quote is accepted"
	}

	quoteStripeID := inv.Parent.QuoteDetails.Quote

	dbTpID, err := database.GetTradespersonIDByBillingStripeQuote(quoteStripeID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, http.StatusNotFound, "quote not found for this provider"
		}
		log.Printf("guest checkout GetTradespersonIDByBillingStripeQuote: %v", err)
		return nil, http.StatusInternalServerError, "database error"
	}
	if dbTpID != tradespersonID {
		return nil, http.StatusNotFound, "invoice does not match provider"
	}

	sq, err := quote.Get(quoteStripeID, billingQuoteParams(tpStripeID, []*string{
		stripe.String("invoice"),
	}))
	if err != nil {
		log.Printf("guest checkout quote.Get: %v", err)
		return nil, http.StatusBadRequest, "unable to load Stripe quote"
	}

	if sq.Status != stripe.QuoteStatusAccepted {
		return nil, http.StatusConflict, "quote must be accepted before the customer can pay the invoice"
	}

	if sq.Invoice == nil || sq.Invoice.ID != inv.ID {
		return nil, http.StatusBadRequest, "invoice does not match this quote"
	}

	acct, err := _stripe.GetConnectAccount(tpStripeID)
	if err != nil {
		log.Printf("guest checkout GetConnectAccount: %v", err)
		return nil, http.StatusBadRequest, "unable to load Connect account"
	}
	if !acct.ChargesEnabled {
		return nil, http.StatusConflict, "provider cannot accept charges yet"
	}

	if !quoteConnectMatchesProvider(sq, tpStripeID) {
		return nil, http.StatusBadRequest, "quote does not belong to this Connect account"
	}

	inv, err = guestConnectInvoiceGet(tpStripeID, inv.ID)
	if err != nil {
		log.Printf("guest checkout invoice refresh %s: %v", stripeInvoiceID, err)
		return nil, http.StatusBadGateway, "invoice not available"
	}

	ensureQuoteInvoiceDepositFooter(tpStripeID, tradespersonID, quoteStripeID, sq)

	inv, err = finalizeInvoiceIfDraft(inv, tpStripeID)
	if err != nil {
		log.Printf("guest checkout finalizeInvoiceIfDraft: %v", err)
		return nil, http.StatusBadGateway, "could not finalize invoice for payment"
	}

	var depositPct int64
	if err := database.GetConnection().QueryRow(
		`SELECT q.depositPct FROM tradesperson_quotes tq INNER JOIN quotes q ON tq.quoteId=q.id WHERE tq.tradespersonId=? AND tq.quote=?`,
		tradespersonID, quoteStripeID,
	).Scan(&depositPct); err != nil {
		depositPct = -1
	}

	return &guestCheckoutInvoiceEnv{
		Inv:         inv,
		TpStripeID:  tpStripeID,
		QuoteStripe: quoteStripeID,
		DepositPct:  depositPct,
	}, 0, ""
}

// guestQuotePublicFieldsFromInvoice returns DB-backed guest UI fields when the invoice was created from a billing quote
// (same source as guest-accept-quote). Requires inv.Parent.QuoteDetails (request expand parent.quote_details on invoice Get when needed).
func guestQuotePublicFieldsFromInvoice(tradespersonID string, inv *stripe.Invoice) (depositPct int64, providerEmail, serviceName, request string) {
	depositPct = -1
	if inv == nil || inv.Parent == nil || inv.Parent.QuoteDetails == nil {
		return depositPct, "", "", ""
	}
	qid := strings.TrimSpace(inv.Parent.QuoteDetails.Quote)
	if qid == "" {
		return depositPct, "", "", ""
	}
	dbTpID, err := database.GetTradespersonIDByBillingStripeQuote(qid)
	if err != nil || dbTpID != tradespersonID {
		return depositPct, "", "", ""
	}
	if err := database.GetConnection().QueryRow(
		`SELECT q.depositPct FROM tradesperson_quotes tq INNER JOIN quotes q ON tq.quoteId=q.id WHERE tq.tradespersonId=? AND tq.quote=?`,
		tradespersonID, qid,
	).Scan(&depositPct); err != nil {
		depositPct = -1
	}
	pd, err := database.GetGuestAcceptQuotePageDetails(tradespersonID, qid)
	if err != nil {
		log.Printf("guestQuotePublicFieldsFromInvoice GetGuestAcceptQuotePageDetails: %v", err)
		return depositPct, "", "", ""
	}
	return depositPct, pd.ProviderEmail, pd.ServiceName, pd.Request
}

type postCheckoutGuestAcceptQuoteBody struct {
	TradespersonID string `json:"tradespersonId"`
	StripeQuoteID  string `json:"stripeQuoteId"`
	Token          string `json:"token"`
}

type postCheckoutGuestAcceptQuoteResponse struct {
	TradespersonID  string `json:"tradespersonId"`
	StripeInvoiceID string `json:"stripeInvoiceId"`
	// StripeQuoteID echoes the billing quote id (qt_…) from the request for symmetry with guest-invoice-status.
	StripeQuoteID string `json:"stripeQuoteId,omitempty"`
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

// guestConnectInvoiceGetExpand loads an invoice for a Connect billing flow (try connected account first).
func guestConnectInvoiceGetExpand(connectAccountID, invoiceID string, expand []*string) (*stripe.Invoice, error) {
	baseParams := &stripe.InvoiceParams{
		Params: stripe.Params{
			Expand: expand,
		},
	}
	if connectAccountID != "" {
		p := &stripe.InvoiceParams{
			Params: stripe.Params{
				Expand: expand,
			},
		}
		p.SetStripeAccount(connectAccountID)
		inv, err := invoice.Get(invoiceID, p)
		if err == nil {
			return inv, nil
		}
		log.Printf("guest invoice.Get %s with Connect account: %v", invoiceID, err)
	}
	return invoice.Get(invoiceID, baseParams)
}

// guestConnectInvoiceGet loads an invoice for a Connect billing flow (try connected account first).
func guestConnectInvoiceGet(connectAccountID, invoiceID string) (*stripe.Invoice, error) {
	return guestConnectInvoiceGetExpand(connectAccountID, invoiceID, nil)
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
		"Deposit terms: %d%% of the quoted total ($%s) is due with the first payment. Remaining balance after that payment: $%s.",
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

// guestCheckoutLineProductCopy builds Stripe Checkout line item name/description and a line for PaymentIntent text.
// Checkout does not show the invoice, so copy uses dollar amounts against the invoice total instead of “same invoice” language.
func guestCheckoutLineProductCopy(
	base string,
	depositPct int64,
	partialDeposit bool,
	chargeCents, invoiceTotalCents, amountRemainingBefore, amountPaidBefore int64,
) (name string, productDescription string, paymentIntentLine string) {
	dp := depositPct
	if dp < 0 || dp > 100 {
		dp = 0
	}
	base = strings.TrimSpace(base)
	if base == "" {
		base = "Invoice payment"
	}
	format := formatCentsForFooter

	balanceAfter := amountRemainingBefore - chargeCents
	if balanceAfter < 0 {
		balanceAfter = 0
	}

	switch {
	case dp == 0:
		suf := "Full amount due now (no separate deposit)"
		return combineCheckoutLineItemDisplayName(base, suf),
			"Full quoted amount due with this payment (no separate deposit).",
			"Full quoted amount due with this payment (no separate deposit)."
	case partialDeposit && dp > 0:
		suf := fmt.Sprintf("%d%% deposit ($%s of $%s)", dp, format(chargeCents), format(invoiceTotalCents))
		desc := fmt.Sprintf(
			"$%s of $%s total (%d%% deposit).\nRemaining balance after this payment: $%s.",
			format(chargeCents), format(invoiceTotalCents), dp, format(balanceAfter),
		)
		return combineCheckoutLineItemDisplayName(base, suf), desc, strings.ReplaceAll(desc, "\n", " ")
	case !partialDeposit && dp > 0 && amountPaidBefore > 0:
		suf := "Remaining balance due"
		line := fmt.Sprintf("Pay the remaining $%s of $%s total.", format(amountRemainingBefore), format(invoiceTotalCents))
		return combineCheckoutLineItemDisplayName(base, suf), line, line
	default:
		// Full remaining in one Checkout session (e.g. 100%% deposit or no split).
		suf := fmt.Sprintf("Full payment (%d%% of total)", dp)
		line := fmt.Sprintf("$%s of $%s total due with this payment.", format(chargeCents), format(invoiceTotalCents))
		return combineCheckoutLineItemDisplayName(base, suf), line, line
	}
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

func combineCheckoutLineItemDisplayName(base, suffix string) string {
	base = strings.TrimSpace(base)
	suf := strings.TrimSpace(suffix)
	if suf == "" {
		return truncateStripeLineItemName(base, stripeCheckoutLineItemNameMaxBytes)
	}
	if base == "" {
		base = "Invoice payment"
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

// ensureQuoteInvoiceDepositFooter sets deposit/payment terms on the invoice footer while the invoice is still a draft.
// Stripe rejects footer updates after finalization (draft → open), and Connect invoices require Stripe-Account on API calls.
func ensureQuoteInvoiceDepositFooter(connectAccountID, tradespersonID, stripeQuoteID string, sq *stripe.Quote) {
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
	inv, err := guestConnectInvoiceGet(connectAccountID, sq.Invoice.ID)
	if err != nil {
		log.Printf("ensureQuoteInvoiceDepositFooter invoice.Get %s: %v", sq.Invoice.ID, err)
		return
	}
	if strings.TrimSpace(inv.Footer) != "" {
		return
	}
	if inv.Status != stripe.InvoiceStatusDraft {
		return
	}
	up := &stripe.InvoiceParams{Footer: stripe.String(footer)}
	if connectAccountID != "" {
		up.SetStripeAccount(connectAccountID)
	}
	if _, err := invoice.Update(sq.Invoice.ID, up); err != nil {
		log.Printf("ensureQuoteInvoiceDepositFooter invoice.Update footer %s: %v", sq.Invoice.ID, err)
	}
}

// AfterQuoteAcceptedEnsuresInvoiceDepositNote applies payment/deposit footer text on the draft invoice, then finalizes if needed.
// Use after quote.Accept (e.g. logged-in customer flow) so the generated invoice matches guest accept behavior.
func AfterQuoteAcceptedEnsuresInvoiceDepositNote(tradespersonID, stripeQuoteID string) {
	tpStripeID, errStripe := database.GetTradespersonStripeID(tradespersonID)
	if errStripe != nil {
		log.Printf("AfterQuoteAccepted GetTradespersonStripeID %s: %v", tradespersonID, errStripe)
	}
	sq, err := quote.Get(stripeQuoteID, billingQuoteParams(tpStripeID, []*string{stripe.String("invoice")}))
	if err != nil {
		log.Printf("AfterQuoteAccepted quote.Get %s: %v", stripeQuoteID, err)
		return
	}
	if sq.Invoice == nil || sq.Invoice.ID == "" {
		return
	}
	ensureQuoteInvoiceDepositFooter(tpStripeID, tradespersonID, stripeQuoteID, sq)
	inv, err := guestConnectInvoiceGet(tpStripeID, sq.Invoice.ID)
	if err != nil {
		log.Printf("AfterQuoteAccepted invoice.Get %s: %v", sq.Invoice.ID, err)
		return
	}
	if _, err := finalizeInvoiceIfDraft(inv, tpStripeID); err != nil {
		log.Printf("AfterQuoteAccepted finalizeInvoiceIfDraft %s: %v", sq.Invoice.ID, err)
	}
}

func notifyTradespersonQuoteAcceptedGuest(connectAccountID, tradespersonID, stripeQuoteID string, sq *stripe.Quote) {
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
	cuParams := &stripe.CustomerParams{}
	if connectAccountID != "" {
		cuParams.SetStripeAccount(connectAccountID)
	}
	stripeCustomer, err := customer.Get(sq.Customer.ID, cuParams)
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

	expandList := []*string{
		stripe.String("customer"),
		stripe.String("invoice"),
	}

	sq, err := quote.Get(body.StripeQuoteID, billingQuoteParams(tpStripeID, expandList))
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
		acceptParams := &stripe.QuoteAcceptParams{}
		acceptParams.SetStripeAccount(tpStripeID)
		sq, err = quote.Accept(body.StripeQuoteID, acceptParams)
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

	sq, err = quote.Get(body.StripeQuoteID, billingQuoteParams(tpStripeID, expandList))
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

	ensureQuoteInvoiceDepositFooter(tpStripeID, body.TradespersonID, body.StripeQuoteID, sq)

	invAccept, err := guestConnectInvoiceGet(tpStripeID, sq.Invoice.ID)
	if err != nil {
		log.Printf("guest accept invoice.Get %s: %v", sq.Invoice.ID, err)
		w.WriteHeader(http.StatusBadGateway)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "invoice not available"})
		return
	}
	if _, err = finalizeInvoiceIfDraft(invAccept, tpStripeID); err != nil {
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

	if didAccept {
		notifyTradespersonQuoteAcceptedGuest(tpStripeID, body.TradespersonID, body.StripeQuoteID, sq)
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
		StripeQuoteID:   body.StripeQuoteID,
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
	if tpStripeID == "" {
		return false
	}
	// Legacy destination-charge quotes (platform invoice + transfer_data).
	if sq.TransferData != nil && sq.TransferData.Destination != nil && sq.TransferData.Destination.ID != "" {
		return sq.TransferData.Destination.ID == tpStripeID
	}
	if sq.OnBehalfOf != nil && sq.OnBehalfOf.ID != "" {
		return sq.OnBehalfOf.ID == tpStripeID
	}
	// Direct charges: quote and invoice live on the connected account.
	return true
}

// finalizeInvoiceIfDraft turns a post–accept quote invoice from draft → open. Without this, customers see
// "invoice is not open for payment" until the provider finalizes in the dashboard or app.
// connectAccountID must be set when the invoice lives on a Stripe Connect account.
func finalizeInvoiceIfDraft(inv *stripe.Invoice, connectAccountID string) (*stripe.Invoice, error) {
	if inv.Status != stripe.InvoiceStatusDraft {
		return inv, nil
	}
	finalizeParams := &stripe.InvoiceFinalizeInvoiceParams{}
	if connectAccountID != "" {
		finalizeParams.SetStripeAccount(connectAccountID)
	}
	if _, err := invoice.FinalizeInvoice(inv.ID, finalizeParams); err != nil {
		return nil, err
	}
	getParams := &stripe.InvoiceParams{
		Params: stripe.Params{
			Expand: []*string{
				stripe.String("quote"),
				stripe.String("customer"),
			},
		},
	}
	if connectAccountID != "" {
		getParams.SetStripeAccount(connectAccountID)
	}
	return invoice.Get(inv.ID, getParams)
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

	env, loadCode, loadErr := guestCheckoutLoadInvoiceForGuestPay(body.TradespersonID, body.StripeInvoiceID)
	if loadCode != 0 {
		w.WriteHeader(loadCode)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": loadErr})
		return
	}
	inv := env.Inv
	tpStripeID := env.TpStripeID
	quoteStripeID := env.QuoteStripe
	depositPct := env.DepositPct

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

	dpForCalc := depositPct
	if dpForCalc < 0 || dpForCalc > 100 {
		dpForCalc = 0
	}
	amountCents, partialDeposit, okAmount := guestDepositChargeCents(inv, dpForCalc)
	if !okAmount {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "amount due is too small"})
		return
	}

	if partialDeposit && stripeInvoiceCustomerID(inv) == "" {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "deposit requires a Stripe customer on the invoice; use the full hosted invoice link or contact the provider."})
		return
	}

	// Hosted invoice page cannot collect a custom partial amount — use Checkout for deposit slices.
	if u := strings.TrimSpace(inv.HostedInvoiceURL); u != "" && !partialDeposit {
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

	displayDepositPct := depositPct
	if displayDepositPct < 0 || displayDepositPct > 100 {
		displayDepositPct = 0
	}

	invoiceTotalCents := inv.Total
	if invoiceTotalCents <= 0 {
		invoiceTotalCents = inv.AmountPaid + inv.AmountRemaining
	}

	sellingFee, err := database.GetTradespersonSellingFee(body.TradespersonID)
	if err != nil {
		log.Printf("guest checkout selling fee: %v", err)
		sellingFee = defaultSellingFeeFraction
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

	lineName, productDescription, piLine := guestCheckoutLineProductCopy(
		desc,
		displayDepositPct,
		partialDeposit,
		amountCents,
		invoiceTotalCents,
		inv.AmountRemaining,
		inv.AmountPaid,
	)

	paymentSegment := guestInvoicePaymentSegment(inv, partialDeposit, displayDepositPct)

	productData := &stripe.CheckoutSessionLineItemPriceDataProductDataParams{
		Name: stripe.String(lineName),
	}
	if productDescription != "" {
		productData.Description = stripe.String(productDescription)
	}

	piDesc := desc
	if piLine != "" {
		piDesc = fmt.Sprintf("%s — %s", desc, piLine)
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
			Metadata: map[string]string{
				"tradesperson_id":            body.TradespersonID,
				"stripe_invoice_id":          inv.ID,
				"stripe_quote_id":            quoteStripeID,
				"guest_pay":                  "true",
				"stripe_connect_account_id":  tpStripeID,
				invoicePaymentSegmentMetaKey: paymentSegment,
			},
		},
	}
	params.SetStripeAccount(tpStripeID)
	params.AddMetadata("tradesperson_id", body.TradespersonID)
	params.AddMetadata("stripe_invoice_id", inv.ID)
	params.AddMetadata("stripe_quote_id", quoteStripeID)
	params.AddMetadata("guest_pay", "true")
	params.AddMetadata("stripe_connect_account_id", tpStripeID)
	params.AddMetadata(invoicePaymentSegmentMetaKey, paymentSegment)
	if partialDeposit {
		params.AddMetadata(guestCheckoutMetadataInvoiceAttach, "true")
		params.Customer = stripe.String(stripeInvoiceCustomerID(inv))
		params.PaymentIntentData.AddMetadata(guestCheckoutMetadataInvoiceAttach, "true")
	}

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
