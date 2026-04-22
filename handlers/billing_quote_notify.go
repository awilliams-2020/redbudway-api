package handlers

import (
	"database/sql"
	"fmt"
	"html"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"

	"redbudway-api/database"
	"redbudway-api/email"
	"redbudway-api/restapi/operations"

	"github.com/go-openapi/runtime/middleware"
	"github.com/stripe/stripe-go/v82"
	"github.com/stripe/stripe-go/v82/customer"
	"github.com/stripe/stripe-go/v82/product"
	"github.com/stripe/stripe-go/v82/quote"
)

func errMap(code int, msg string) middleware.Responder {
	return middleware.Error(code, map[string]string{"error": msg})
}

// PostTradespersonTradespersonIDBillingQuoteQuoteIDNotifyCustomerHandler emails the Stripe customer about an updated billing quote (US-5). Generated route from swagger.
func PostTradespersonTradespersonIDBillingQuoteQuoteIDNotifyCustomerHandler(params operations.PostTradespersonTradespersonIDBillingQuoteQuoteIDNotifyCustomerParams, principal interface{}) middleware.Responder {
	tradespersonID := params.TradespersonID
	stripeQuoteID := params.QuoteID
	token := params.HTTPRequest.Header.Get("Authorization")

	valid, err := ValidateTradespersonAccessToken(tradespersonID, token)
	if err != nil {
		log.Printf("notify quote: validate tradesperson %s: %v", tradespersonID, err)
		return errMap(http.StatusUnauthorized, "unauthorized")
	}
	if !valid {
		return errMap(http.StatusUnauthorized, "unauthorized")
	}

	tpStripeID, err := database.GetTradespersonStripeID(tradespersonID)
	if err != nil || tpStripeID == "" {
		return errMap(http.StatusNotFound, "provider not found")
	}

	sq, err := quote.Get(stripeQuoteID, billingQuoteParams(tpStripeID, []*string{
		stripe.String("customer"),
		stripe.String("transfer_data.destination"),
		stripe.String("on_behalf_of"),
	}))
	if err != nil {
		log.Printf("notify quote.Get %s: %v", stripeQuoteID, err)
		return errMap(http.StatusBadRequest, "unable to load quote")
	}

	if !quoteConnectMatchesProvider(sq, tpStripeID) {
		return errMap(http.StatusNotFound, "quote not found for this account")
	}

	if sq.Status == stripe.QuoteStatusCanceled {
		return errMap(http.StatusConflict, "quote is canceled")
	}
	if sq.Status == stripe.QuoteStatusAccepted {
		return errMap(http.StatusConflict, "quote is accepted; use the invoice payment flow instead of quote update emails")
	}

	custEmail, custName, err := billingQuoteStripeCustomerEmailName(tpStripeID, sq)
	if err != nil {
		log.Printf("notify customer resolution: %v", err)
		return errMap(http.StatusInternalServerError, "unable to load customer")
	}
	if strings.TrimSpace(custEmail) == "" {
		return errMap(http.StatusBadRequest, "customer email not available for this quote")
	}

	tp, profErr := database.GetTradespersonProfile(tradespersonID)
	if profErr != nil {
		log.Printf("notify GetTradespersonProfile: %v", profErr)
	}
	providerName := strings.TrimSpace(tp.Name)
	if providerName == "" {
		providerName = "Your provider"
	}

	desc := strings.TrimSpace(sq.Description)
	if desc == "" {
		desc = "—"
	}

	depositPct := int64(0)
	_ = database.GetConnection().QueryRow(
		`SELECT q.depositPct FROM tradesperson_quotes tq INNER JOIN quotes q ON tq.quoteId=q.id WHERE tq.tradespersonId=? AND tq.quote=?`,
		tradespersonID, stripeQuoteID,
	).Scan(&depositPct)

	var payNote string
	if sq.Status == stripe.QuoteStatusOpen {
		acceptURL := guestAcceptQuoteURL(tradespersonID, stripeQuoteID)
		payNote = `<p style="margin: 0; color:black;">Review and accept this quote when you’re ready. After you accept, an invoice is created and you can pay by card using the secure payment link you receive.</p><br>` +
			buildFinalizedQuotePayHTML(acceptURL)
	} else {
		payNote = "<p>This quote is not yet open for payment. You’ll be notified when it’s ready to review and pay.</p>"
	}

	lineItemsHTML := buildLineItemsHTML(tpStripeID, stripeQuoteID, sq.AmountSubtotal, quoteTaxAmount(sq), sq.AmountTotal, depositPct)

	emailID, err := email.SendBillingQuoteCustomerUpdate(custEmail, custName, providerName, strings.TrimSpace(tp.Email), desc, lineItemsHTML, payNote)
	if err != nil {
		log.Printf("notify SendBillingQuoteCustomerUpdate: %v", err)
		return errMap(http.StatusInternalServerError, "failed to send email")
	}

	payload := &operations.PostTradespersonTradespersonIDBillingQuoteQuoteIDNotifyCustomerOKBody{
		Sent:    true,
		EmailID: emailID,
	}
	return operations.NewPostTradespersonTradespersonIDBillingQuoteQuoteIDNotifyCustomerOK().WithPayload(payload)
}

// billingQuoteStripeCustomerEmailName resolves the recipient for billing quote customer emails.
func billingQuoteStripeCustomerEmailName(connectAccountID string, sq *stripe.Quote) (custEmail, custName string, err error) {
	if sq == nil {
		return "", "", nil
	}
	if sq.Customer != nil && sq.Customer.ID != "" {
		cu := &stripe.CustomerParams{}
		if connectAccountID != "" {
			cu.SetStripeAccount(connectAccountID)
		}
		sc, e := customer.Get(sq.Customer.ID, cu)
		if e != nil {
			return "", "", e
		}
		if !sc.Deleted {
			return sc.Email, sc.Name, nil
		}
		if sq.Status == stripe.QuoteStatusAccepted && sq.Invoice != nil && sq.Invoice.ID != "" {
			inv, e := loadStripeInvoiceForBillingQuote(connectAccountID, sq.Invoice.ID)
			if e != nil {
				return "", "", e
			}
			return inv.CustomerEmail, inv.CustomerName, nil
		}
	}
	return "", "", nil
}

// sendFinalizedQuoteReadyEmail sends the same HTML email as post-finalize (hero “Your quote is ready”, Accept quote CTA).
func sendFinalizedQuoteReadyEmail(tradespersonID, stripeQuoteID string, sq *stripe.Quote) (string, error) {
	depositPct := int64(0)
	_ = database.GetConnection().QueryRow(
		`SELECT q.depositPct FROM tradesperson_quotes tq INNER JOIN quotes q ON tq.quoteId=q.id WHERE tq.tradespersonId=? AND tq.quote=?`,
		tradespersonID, stripeQuoteID,
	).Scan(&depositPct)

	tpStripeID, stripeAccErr := database.GetTradespersonStripeID(tradespersonID)
	if stripeAccErr != nil {
		return "", fmt.Errorf("provider Stripe account: %w", stripeAccErr)
	}
	if tpStripeID == "" {
		return "", fmt.Errorf("provider has no Stripe Connect account")
	}

	custEmail, custName, err := billingQuoteStripeCustomerEmailName(tpStripeID, sq)
	if err != nil {
		return "", err
	}
	if strings.TrimSpace(custEmail) == "" {
		return "", fmt.Errorf("no customer email")
	}

	tp, profErr := database.GetTradespersonProfile(tradespersonID)
	if profErr != nil {
		log.Printf("sendFinalizedQuoteReadyEmail GetTradespersonProfile: %v", profErr)
	}
	providerName := strings.TrimSpace(tp.Name)
	if providerName == "" {
		providerName = "Your provider"
	}

	desc := strings.TrimSpace(sq.Description)
	if desc == "" {
		desc = "—"
	}

	lineItemsHTML := buildLineItemsHTML(tpStripeID, stripeQuoteID, sq.AmountSubtotal, quoteTaxAmount(sq), sq.AmountTotal, depositPct)
	payBlock := buildFinalizedQuotePayHTML(guestAcceptQuoteURL(tradespersonID, stripeQuoteID))
	pn := providerName
	lead := fmt.Sprintf(
		`<p style="margin: 0; color:black;"><strong>%s</strong> has finalized this quote. Review the details below, then accept to continue to secure checkout.</p>`,
		html.EscapeString(pn),
	)
	return email.SendBillingQuoteCustomerEmail(email.BillingQuoteCustomerEmailParams{
		CustomerEmail:     custEmail,
		CustomerName:      custName,
		ProviderEmail:     strings.TrimSpace(tp.Email),
		Description:       desc,
		LineItemsHTML:     lineItemsHTML,
		PayBlockHTML:      payBlock,
		Subject:           fmt.Sprintf("Your quote from %s is ready", pn),
		Preheader:         "Your quote is ready — review the details and accept to pay.",
		PageTitle:         "Your quote is ready",
		HeroTitle:         "Your quote is ready",
		LeadParagraphHTML: lead,
		FooterPermission:  "You received this email about a payment quote through Redbud Way.",
	})
}

// PostTradespersonTradespersonIDBillingQuoteQuoteIDResendFinalizedEmailHandler resends the finalized-quote email for an open Stripe quote.
func PostTradespersonTradespersonIDBillingQuoteQuoteIDResendFinalizedEmailHandler(params operations.PostTradespersonTradespersonIDBillingQuoteQuoteIDResendFinalizedEmailParams, principal interface{}) middleware.Responder {
	tradespersonID := params.TradespersonID
	stripeQuoteID := params.QuoteID
	token := params.HTTPRequest.Header.Get("Authorization")

	valid, err := ValidateTradespersonAccessToken(tradespersonID, token)
	if err != nil {
		log.Printf("resend finalized email: validate tradesperson %s: %v", tradespersonID, err)
		return errMap(http.StatusUnauthorized, "unauthorized")
	}
	if !valid {
		return errMap(http.StatusUnauthorized, "unauthorized")
	}

	db := database.GetConnection()
	var stub int
	switch err = db.QueryRow(
		`SELECT 1 FROM tradesperson_quotes tq INNER JOIN quotes q ON tq.quoteId=q.id WHERE tq.tradespersonId=? AND tq.quote=?`,
		tradespersonID, stripeQuoteID,
	).Scan(&stub); err {
	case sql.ErrNoRows:
		return errMap(http.StatusNotFound, "quote not found")
	case nil:
	default:
		log.Printf("resend finalized email lookup: %v", err)
		return errMap(http.StatusInternalServerError, "internal error")
	}

	tpStripeID, err := database.GetTradespersonStripeID(tradespersonID)
	if err != nil || tpStripeID == "" {
		return errMap(http.StatusNotFound, "provider not found")
	}

	sq, err := quote.Get(stripeQuoteID, billingQuoteParams(tpStripeID, []*string{
		stripe.String("customer"),
		stripe.String("transfer_data.destination"),
		stripe.String("on_behalf_of"),
	}))
	if err != nil {
		log.Printf("resend finalized email quote.Get %s: %v", stripeQuoteID, err)
		return errMap(http.StatusBadRequest, "unable to load quote")
	}

	if !quoteConnectMatchesProvider(sq, tpStripeID) {
		return errMap(http.StatusNotFound, "quote not found for this account")
	}

	if sq.Status == stripe.QuoteStatusCanceled {
		return errMap(http.StatusConflict, "quote is canceled")
	}
	if sq.Status != stripe.QuoteStatusOpen {
		return errMap(http.StatusConflict, "quote must be open to resend the finalized email")
	}

	emailID, err := sendFinalizedQuoteReadyEmail(tradespersonID, stripeQuoteID, sq)
	if err != nil {
		if strings.Contains(err.Error(), "no customer email") {
			return errMap(http.StatusBadRequest, "customer email not available for this quote")
		}
		log.Printf("resend finalized email send: %v", err)
		return errMap(http.StatusInternalServerError, "failed to send email")
	}

	payload := &operations.PostTradespersonTradespersonIDBillingQuoteQuoteIDResendFinalizedEmailOKBody{
		Sent:    true,
		EmailID: emailID,
	}
	return operations.NewPostTradespersonTradespersonIDBillingQuoteQuoteIDResendFinalizedEmailOK().WithPayload(payload)
}

func formatUSDFromCents(cents int64) string {
	if cents < 0 {
		cents = 0
	}
	d := float64(cents) / 100.0
	return fmt.Sprintf("$%.2f", d)
}

func guestAcceptQuoteURL(tradespersonID, stripeQuoteID string) string {
	sub := os.Getenv("SUBDOMAIN")
	origin := "https://" + sub + "redbudway.com"
	v := url.Values{}
	v.Set("stripeQuoteId", strings.TrimSpace(stripeQuoteID))
	v.Set("tradespersonId", strings.TrimSpace(tradespersonID))
	return origin + "/accept/quote?" + v.Encode()
}

// buildFinalizedQuotePayHTML is the accept CTA for guest checkout (same destination as the “Accept quote” flow).
func buildFinalizedQuotePayHTML(acceptURL string) string {
	esc := html.EscapeString(acceptURL)
	return `<table role="presentation" border="0" cellpadding="0" cellspacing="0" width="100%">` +
		`<tr><td align="center" bgcolor="#ffffff" style="padding: 12px 0 4px;">` +
		`<table role="presentation" border="0" cellpadding="0" cellspacing="0"><tr>` +
		`<td align="center" bgcolor="#4682b4" style="border-radius: 4px;">` +
		`<a href="` + esc + `" target="_blank" rel="noopener noreferrer" style="display: inline-block; padding: 16px 36px; font-family: 'Rajdhani', sans-serif; font-size: 16px; color: #ffffff; text-decoration: none; font-weight: 600;">Accept quote</a>` +
		`</td></tr></table></td></tr></table>` +
		`<p style="margin: 16px 0 0; color:black;">The next step is secure checkout where you can pay by card as a guest — no account needed.</p>` +
		`<p style="margin: 12px 0 0; font-size: 14px; line-height: 20px; color: #444; word-break: break-all;">If the button doesn’t work, open or paste this link:<br><a href="` + esc + `" target="_blank" rel="noopener noreferrer" style="color: #1a82e2;">` + esc + `</a></p>`
}

func quoteTaxAmount(sq *stripe.Quote) int64 {
	if sq == nil || sq.TotalDetails == nil {
		return 0
	}
	return sq.TotalDetails.AmountTax
}

func buildLineItemsHTML(connectAccountID, stripeQuoteID string, amountSubtotal, amountTax, amountTotal, depositPct int64) string {
	var sb strings.Builder
	sb.WriteString(`<table border="0" cellpadding="0" cellspacing="0" width="100%" style="border-collapse: collapse; font-size: 15px; color: black;">`)
	sb.WriteString(`<tr style="background-color: #f5f5f5;">` +
		`<td style="padding: 8px 4px; border-bottom: 2px solid #d4dadf; font-weight: 700;">Item</td>` +
		`<td align="center" style="padding: 8px 4px; border-bottom: 2px solid #d4dadf; font-weight: 700;">Qty</td>` +
		`<td align="right" style="padding: 8px 4px; border-bottom: 2px solid #d4dadf; font-weight: 700;">Unit price</td>` +
		`<td align="right" style="padding: 8px 4px; border-bottom: 2px solid #d4dadf; font-weight: 700;">Amount</td>` +
		`</tr>`)

	params := &stripe.QuoteListLineItemsParams{Quote: stripe.String(stripeQuoteID)}
	if connectAccountID != "" {
		params.SetStripeAccount(connectAccountID)
	}
	it := quote.ListLineItems(params)
	for it.Next() {
		li := it.LineItem()
		name := ""
		if li.Price != nil && li.Price.Product != nil {
			pp := &stripe.ProductParams{}
			if connectAccountID != "" {
				pp.SetStripeAccount(connectAccountID)
			}
			if p, err := product.Get(li.Price.Product.ID, pp); err == nil {
				name = p.Name
			}
		}
		if name == "" {
			name = "(item)"
		}
		unitAmt := int64(0)
		if li.Price != nil {
			unitAmt = li.Price.UnitAmount
		}
		lineTotal := li.Quantity * unitAmt
		sb.WriteString(fmt.Sprintf(
			`<tr>`+
				`<td style="padding: 8px 4px; border-bottom: 1px solid #eeeeee;">%s</td>`+
				`<td align="center" style="padding: 8px 4px; border-bottom: 1px solid #eeeeee;">%d</td>`+
				`<td align="right" style="padding: 8px 4px; border-bottom: 1px solid #eeeeee;">%s</td>`+
				`<td align="right" style="padding: 8px 4px; border-bottom: 1px solid #eeeeee;">%s</td>`+
				`</tr>`,
			html.EscapeString(name), li.Quantity, formatUSDFromCents(unitAmt), formatUSDFromCents(lineTotal),
		))
	}
	if err := it.Err(); err != nil {
		log.Printf("buildLineItemsHTML ListLineItems %s: %v", stripeQuoteID, err)
	}

	if amountSubtotal != amountTotal {
		sb.WriteString(fmt.Sprintf(
			`<tr>`+
				`<td colspan="3" align="right" style="padding: 6px 4px; color: #555;">Subtotal</td>`+
				`<td align="right" style="padding: 6px 4px; color: #555;">%s</td>`+
				`</tr>`,
			formatUSDFromCents(amountSubtotal),
		))
	}
	if amountTax > 0 {
		sb.WriteString(fmt.Sprintf(
			`<tr>`+
				`<td colspan="3" align="right" style="padding: 6px 4px; color: #555;">Tax</td>`+
				`<td align="right" style="padding: 6px 4px; color: #555;">%s</td>`+
				`</tr>`,
			formatUSDFromCents(amountTax),
		))
	}
	sb.WriteString(fmt.Sprintf(
		`<tr>`+
			`<td colspan="3" align="right" style="padding: 8px 4px; font-weight: 700; border-top: 2px solid #d4dadf;">Total</td>`+
			`<td align="right" style="padding: 8px 4px; font-weight: 700; border-top: 2px solid #d4dadf;">%s</td>`+
			`</tr>`,
		formatUSDFromCents(amountTotal),
	))
	if depositPct > 0 && depositPct <= 100 && amountTotal > 0 {
		depositDue := (amountTotal * depositPct) / 100
		balance := amountTotal - depositDue
		sb.WriteString(fmt.Sprintf(
			`<tr>`+
				`<td colspan="3" align="right" style="padding: 8px 4px; color: #555; border-top: 1px solid #eeeeee;">Deposit due at checkout (%d%%)</td>`+
				`<td align="right" style="padding: 8px 4px; color: #555; border-top: 1px solid #eeeeee;">%s</td>`+
				`</tr>`,
			depositPct, formatUSDFromCents(depositDue),
		))
		sb.WriteString(fmt.Sprintf(
			`<tr>`+
				`<td colspan="3" align="right" style="padding: 6px 4px; color: #555;">Remaining balance</td>`+
				`<td align="right" style="padding: 6px 4px; color: #555;">%s</td>`+
				`</tr>`,
			formatUSDFromCents(balance),
		))
	}
	sb.WriteString(`</table>`)
	return sb.String()
}
