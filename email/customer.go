package email

import (
	_ "embed"
	"fmt"
	"html"
	"os"
	"strings"

	"redbudway-api/models"

	"github.com/stripe/stripe-go/v82"
)

//go:embed html/verification.html
var verification string

//go:embed html/customer-confirmation.html
var customerConfirmation string

//go:embed html/customer-quote-confirmation.html
var customerQuoteConfirmation string

//go:embed html/customer-subscription-confirmation.html
var customerSubscriptionConfirmation string

//go:embed html/cancellation.html
var cancellation string

//go:embed html/quote-cancellation.html
var quoteCancellation string

//go:embed html/subscription-cancellation.html
var subscriptionCancellation string

//go:embed html/quote-invoice-void.html
var quoteInvoiceVoid string

//go:embed html/void.html
var void string

//go:embed html/invoice.html
var invoice string

//go:embed html/refund.html
var refund string

//go:embed html/billing-quote-customer-update.html
var billingQuoteCustomerUpdate string

func SendCustomerVerification(customerName, customerEmail, customerID, token string) error {
	body := verification

	body = strings.Replace(body, "{SUBDOMAIN}", os.Getenv("SUBDOMAIN"), -1)
	body = strings.Replace(body, "{CUSTOMERID}", customerID, -1)
	body = strings.Replace(body, "{TOKEN}", token, -1)

	return email(customerEmail, customerName, "Email Verification", body)
}

func SendCustomerConfirmation(tradesperson models.Tradesperson, stripeCustomer *stripe.Customer, stripeProduct *stripe.Product, timeAndPrice, formRowsCols string) error {
	body := customerConfirmation

	tradespersonInfo := fmt.Sprintf("%s<br>%s<br>%s", tradesperson.Name, tradesperson.Email, tradesperson.Number)
	body = strings.Replace(body, "{TRADESPERSON_INFO}", tradespersonInfo, -1)
	body = strings.Replace(body, "{SERVICE_NAME}", stripeProduct.Name, -1)
	body = strings.Replace(body, "{TIME_AND_PRICE}", timeAndPrice, -1)
	body = strings.Replace(body, "{FORM_ROWS_COLS}", formRowsCols, -1)

	return email(stripeCustomer.Email, stripeCustomer.Name, "Confirmation", body)
}

func SendCustomerQuoteConfirmation(tradesperson models.Tradesperson, stripeCustomer *stripe.Customer, message string, quote *models.ServiceDetails) error {
	body := customerQuoteConfirmation

	tradespersonInfo := fmt.Sprintf("%s<br>%s<br>%s", tradesperson.Name, tradesperson.Email, tradesperson.Number)
	body = strings.Replace(body, "{CUSTOMER_MESSAGE}", message, -1)
	body = strings.Replace(body, "{SERVICE_NAME}", *quote.Title, -1)
	body = strings.Replace(body, "{TRADESPERSON_INFO}", tradespersonInfo, -1)

	_, err := sendHTMLResend(stripeCustomer.Email, stripeCustomer.Name, "Confirmation", body, "")
	return err
}

func SendCustomerCancellation(tradespersonName, tradespersonEmail, tradespersonNumber, timeAndPrice string, stripeInvoice *stripe.Invoice, stripeProduct *stripe.Product) error {
	body := cancellation
	tradespersonInfo := fmt.Sprintf("%s<br>%s<br>%s", tradespersonName, tradespersonEmail, tradespersonNumber)
	body = strings.Replace(body, "{TRADESPERSON_INFO}", tradespersonInfo, -1)
	body = strings.Replace(body, "{SERVICE_NAME}", stripeProduct.Name, -1)
	body = strings.Replace(body, "{TIME_AND_PRICE}", timeAndPrice, -1)

	return email(stripeInvoice.CustomerEmail, stripeInvoice.CustomerName, "Cancellation", body)
}

func SendCustomerQuoteCancellation(tradesperson models.Tradesperson, stripeCustomer *stripe.Customer, message, title string) error {
	body := quoteCancellation

	tradespersonInfo := fmt.Sprintf("%s<br>%s<br>%s", tradesperson.Name, tradesperson.Email, tradesperson.Number)
	body = strings.Replace(body, "{CUSTOMER_MESSAGE}", message, -1)
	body = strings.Replace(body, "{SERVICE_NAME}", title, -1)
	body = strings.Replace(body, "{TRADESPERSON_INFO}", tradespersonInfo, -1)

	return email(stripeCustomer.Email, stripeCustomer.Name, "Cancellation", body)
}

func SendCustomerSubscriptionCancellation(tradesperson models.Tradesperson, stripeCustomer *stripe.Customer, stripeProduct *stripe.Product, timeAndPrice string) error {
	body := subscriptionCancellation

	tradespersonInfo := fmt.Sprintf("%s<br>%s<br>%s", tradesperson.Name, tradesperson.Email, tradesperson.Number)
	body = strings.Replace(body, "{TRADESPERSON_INFO}", tradespersonInfo, -1)
	body = strings.Replace(body, "{SERVICE_NAME}", stripeProduct.Name, -1)
	body = strings.Replace(body, "{TIME_AND_PRICE}", timeAndPrice, -1)

	return email(stripeCustomer.Email, stripeCustomer.Name, "Cancellation", body)
}

func SendCustomerQuoteInvoiceVoid(tradesperson models.Tradesperson, stripeInvoice *stripe.Invoice, message, title string) error {
	body := quoteInvoiceVoid

	tradespersonInfo := fmt.Sprintf("%s<br>%s<br>%s", tradesperson.Name, tradesperson.Email, tradesperson.Number)
	body = strings.Replace(body, "{CUSTOMER_MESSAGE}", message, -1)
	body = strings.Replace(body, "{SERVICE_NAME}", title, -1)
	body = strings.Replace(body, "{TRADESPERSON_INFO}", tradespersonInfo, -1)

	return email(stripeInvoice.CustomerEmail, stripeInvoice.CustomerName, "Invoice voided", body)
}

func SentInvoice(stripeInvoice *stripe.Invoice) error {
	body := invoice

	body = strings.Replace(body, "{DOWNLOAD_INVOICE}", stripeInvoice.InvoicePDF, -1)
	body = strings.Replace(body, "{PAY_INVOICE}", stripeInvoice.HostedInvoiceURL, -1)

	return email(stripeInvoice.CustomerEmail, stripeInvoice.CustomerName, "Invoice", body)
}

func SendCustomerRefund(stripeInvoice *stripe.Invoice, stripeProduct *stripe.Product, decimalPrice float64) error {
	body := refund

	body = strings.Replace(body, "{SERVICE_NAME}", stripeProduct.Name, -1)
	price := fmt.Sprintf("$%.2f", decimalPrice)
	body = strings.Replace(body, "{PRICE}", price, -1)

	return email(stripeInvoice.CustomerEmail, stripeInvoice.CustomerName, "Refund", body)
}

func SendCustomerSubscriptionConfirmation(tradesperson models.Tradesperson, stripeCustomer *stripe.Customer, stripeProduct *stripe.Product, timeAndPrice, formRowsCols string) error {
	body := customerSubscriptionConfirmation

	tradespersonInfo := fmt.Sprintf("%s<br>%s<br>%s", tradesperson.Name, tradesperson.Email, tradesperson.Number)
	body = strings.Replace(body, "{TRADESPERSON_INFO}", tradespersonInfo, -1)
	body = strings.Replace(body, "{SERVICE_NAME}", stripeProduct.Name, -1)
	body = strings.Replace(body, "{TIME_AND_PRICE}", timeAndPrice, -1)
	body = strings.Replace(body, "{FORM_ROWS_COLS}", formRowsCols, -1)

	return email(stripeCustomer.Email, stripeCustomer.Name, "Confirmation", body)
}

func SendCustomerVoid(tradespersonName, tradespersonEmail, tradespersonNumber, itemsAndPrice string, stripeInvoice *stripe.Invoice) error {
	body := void

	tradespersonInfo := fmt.Sprintf("%s<br>%s<br>%s", tradespersonName, tradespersonEmail, tradespersonNumber)
	body = strings.Replace(body, "{TRADESPERSON_INFO}", tradespersonInfo, -1)
	body = strings.Replace(body, "{ITEMS_AND_PRICE}", itemsAndPrice, -1)

	return email(stripeInvoice.CustomerEmail, stripeInvoice.CustomerName, "Voided Invoice", body)
}

// BillingQuoteCustomerEmailParams configures the shared billing quote customer HTML template
// (line items, description, and action block). LeadParagraphHTML and PayBlockHTML are trusted HTML from the handler.
type BillingQuoteCustomerEmailParams struct {
	CustomerEmail string
	CustomerName  string
	// ProviderEmail is the tradesperson profile email: sets Reply-To and the mailto line when non-empty.
	ProviderEmail     string
	Description       string
	LineItemsHTML     string
	PayBlockHTML      string
	Subject           string
	Preheader         string
	PageTitle         string
	HeroTitle         string
	LeadParagraphHTML string
	FooterPermission  string
}

func billingQuoteReplyHintHTML(providerEmail string) string {
	e := strings.TrimSpace(providerEmail)
	if e != "" {
		esc := html.EscapeString(e)
		return fmt.Sprintf(
			`<p style="margin: 0;">If you have questions, email <a href="mailto:%s">%s</a>.</p>`+
				`<p style="margin: 8px 0 0; font-size: 14px; color: #555;">You can also use <strong>Reply</strong> on this message — your response goes to your service provider.</p>`,
			esc, esc,
		)
	}
	return `<p style="margin: 0;">If you have questions, contact your service provider using the phone number or other contact they shared with you.</p>`
}

// SendBillingQuoteCustomerEmail sends the shared billing quote template (updates and finalize).
func SendBillingQuoteCustomerEmail(p BillingQuoteCustomerEmailParams) (string, error) {
	displayName := strings.TrimSpace(p.CustomerName)
	if displayName == "" {
		displayName = "there"
	}
	pageTitle := strings.TrimSpace(p.PageTitle)
	if pageTitle == "" {
		pageTitle = p.HeroTitle
	}
	if strings.TrimSpace(pageTitle) == "" {
		pageTitle = "Quote"
	}
	hero := strings.TrimSpace(p.HeroTitle)
	if hero == "" {
		hero = "Quote"
	}
	subject := strings.TrimSpace(p.Subject)
	if subject == "" {
		subject = "Quote from Redbud Way"
	}
	preheader := p.Preheader
	if strings.TrimSpace(preheader) == "" {
		preheader = "Your provider sent a quote through Redbud Way."
	}
	footer := p.FooterPermission
	if strings.TrimSpace(footer) == "" {
		footer = "You received this email because your provider sent a quote through Redbud Way."
	}
	body := billingQuoteCustomerUpdate
	body = strings.Replace(body, "{PAGE_TITLE}", html.EscapeString(pageTitle), -1)
	body = strings.Replace(body, "{PREHEADER}", html.EscapeString(preheader), -1)
	body = strings.Replace(body, "{HERO_TITLE}", html.EscapeString(hero), -1)
	body = strings.Replace(body, "{CUSTOMER_NAME}", html.EscapeString(displayName), -1)
	body = strings.Replace(body, "{LEAD_PARAGRAPH}", p.LeadParagraphHTML, -1)
	body = strings.Replace(body, "{DESCRIPTION}", html.EscapeString(p.Description), -1)
	body = strings.Replace(body, "{LINE_ITEMS_TABLE}", p.LineItemsHTML, -1)
	body = strings.Replace(body, "{PAY_BLOCK}", p.PayBlockHTML, -1)
	body = strings.Replace(body, "{REPLY_HINT}", billingQuoteReplyHintHTML(p.ProviderEmail), -1)
	body = strings.Replace(body, "{FOOTER_PERMISSION}", html.EscapeString(footer), -1)
	return sendHTMLResend(p.CustomerEmail, displayName, subject, body, p.ProviderEmail)
}

// SendBillingQuoteCustomerUpdate emails the customer after a provider updates a billing quote.
// lineItemsHTML and payHTMLBlock are trusted HTML from the handler.
func SendBillingQuoteCustomerUpdate(customerEmail, customerName, providerName, providerEmail, description, lineItemsHTML, payHTMLBlock string) (string, error) {
	pn := strings.TrimSpace(providerName)
	if pn == "" {
		pn = "Your provider"
	}
	lead := fmt.Sprintf(`<p style="margin: 0; color:black;"><strong>%s</strong> has updated your quote.</p>`, html.EscapeString(pn))
	return SendBillingQuoteCustomerEmail(BillingQuoteCustomerEmailParams{
		CustomerEmail:     customerEmail,
		CustomerName:      customerName,
		ProviderEmail:     providerEmail,
		Description:       description,
		LineItemsHTML:     lineItemsHTML,
		PayBlockHTML:      payHTMLBlock,
		Subject:           fmt.Sprintf("Quote update from %s", pn),
		Preheader:         "Your provider updated a quote — review details and pay when you are ready.",
		PageTitle:         "Quote update",
		HeroTitle:         "Quote updated",
		LeadParagraphHTML: lead,
		FooterPermission:  "You received this email because your provider sent a quote update through Redbud Way.",
	})
}
