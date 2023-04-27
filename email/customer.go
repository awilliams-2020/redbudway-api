package email

import (
	_ "embed"
	"fmt"
	"os"
	"strings"

	"redbudway-api/models"

	"github.com/stripe/stripe-go/v72"
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

//go:embed html/quote.html
var quote string

//go:embed html/refund.html
var refund string

func SendCustomerVerification(customerName, customerEmail, customerID, token string) error {
	body := verification

	body = strings.Replace(body, "{SUBDOMAIN}", os.Getenv("SUBDOMAIN"), -1)
	body = strings.Replace(body, "{CUSTOMERID}", customerID, -1)
	body = strings.Replace(body, "{TOKEN}", token, -1)

	return email(customerEmail, customerName, "Email Verification", body)
}

func SendCustomerConfirmation(tradesperson models.Tradesperson, stripeCustomer *stripe.Customer, stripeProduct *stripe.Product, timeAndPrice string) error {
	body := customerConfirmation

	tradespersonInfo := fmt.Sprintf("%s<br>%s<br>%s", tradesperson.Name, tradesperson.Email, tradesperson.Number)
	body = strings.Replace(body, "{TRADESPERSON_INFO}", tradespersonInfo, -1)
	body = strings.Replace(body, "{SERVICE_NAME}", stripeProduct.Name, -1)
	body = strings.Replace(body, "{TIME_AND_PRICE}", timeAndPrice, -1)

	return email(stripeCustomer.Email, stripeCustomer.Name, "Confirmation", body)
}

func SendCustomerQuoteConfirmation(tradesperson models.Tradesperson, stripeCustomer *stripe.Customer, message string, quote *models.ServiceDetails) error {
	body := customerQuoteConfirmation

	tradespersonInfo := fmt.Sprintf("%s<br>%s<br>%s", tradesperson.Name, tradesperson.Email, tradesperson.Number)
	body = strings.Replace(body, "{CUSTOMER_MESSAGE}", message, -1)
	body = strings.Replace(body, "{SERVICE_NAME}", *quote.Title, -1)
	body = strings.Replace(body, "{TRADESPERSON_INFO}", tradespersonInfo, -1)

	return email(stripeCustomer.Email, stripeCustomer.Name, "Confirmation", body)
}

func SendCustomerCancellation(tradespersonName, tradespersonEmail, tradespersonNumber, timeAndPrice string, stripeInvoice *stripe.Invoice, stripeProduct *stripe.Product) error {
	body := cancellation
	tradespersonInfo := fmt.Sprintf("%s<br>%s<br>%s", tradespersonName, tradespersonEmail, tradespersonNumber)
	body = strings.Replace(body, "{TRADESPERSON_INFO}", tradespersonInfo, -1)
	body = strings.Replace(body, "{SERVICE_NAME}", stripeProduct.Name, -1)
	body = strings.Replace(body, "{TIME_AND_PRICE}", timeAndPrice, -1)

	return email(stripeInvoice.CustomerEmail, *stripeInvoice.CustomerName, "Cancellation", body)
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

	return email(stripeInvoice.CustomerEmail, *stripeInvoice.CustomerName, "Invoice voided", body)
}

func SentInvoice(stripeInvoice *stripe.Invoice) error {
	body := invoice

	body = strings.Replace(body, "{DOWNLOAD_INVOICE}", stripeInvoice.InvoicePDF, -1)
	body = strings.Replace(body, "{PAY_INVOICE}", stripeInvoice.HostedInvoiceURL, -1)

	return email(stripeInvoice.CustomerEmail, *stripeInvoice.CustomerName, "Invoice", body)
}

func SentQuote(stripeCustomer *stripe.Customer, quoteURL string) error {
	body := quote

	body = strings.Replace(body, "{VIEW_QUOTE}", quoteURL, -1)

	return email(stripeCustomer.Email, stripeCustomer.Name, "Quote", body)
}

func SendCustomerRefund(stripeInvoice *stripe.Invoice, stripeProduct *stripe.Product, decimalPrice float64) error {
	body := refund

	body = strings.Replace(body, "{SERVICE_NAME}", stripeProduct.Name, -1)
	price := fmt.Sprintf("$%.2f", decimalPrice)
	body = strings.Replace(body, "{PRICE}", price, -1)

	return email(stripeInvoice.CustomerEmail, *stripeInvoice.CustomerName, "Refund", body)
}

func SendCustomerSubscriptionConfirmation(tradesperson models.Tradesperson, stripeCustomer *stripe.Customer, stripeProduct *stripe.Product, timeAndPrice string) error {
	body := customerSubscriptionConfirmation

	tradespersonInfo := fmt.Sprintf("%s<br>%s<br>%s", tradesperson.Name, tradesperson.Email, tradesperson.Number)
	body = strings.Replace(body, "{TRADESPERSON_INFO}", tradespersonInfo, -1)
	body = strings.Replace(body, "{SERVICE_NAME}", stripeProduct.Name, -1)
	body = strings.Replace(body, "{TIME_AND_PRICE}", timeAndPrice, -1)

	return email(stripeCustomer.Email, stripeCustomer.Name, "Confirmation", body)
}

func SendCustomerVoid(tradespersonName, tradespersonEmail, tradespersonNumber, itemsAndPrice string, stripeInvoice *stripe.Invoice) error {
	body := void

	tradespersonInfo := fmt.Sprintf("%s<br>%s<br>%s", tradespersonName, tradespersonEmail, tradespersonNumber)
	body = strings.Replace(body, "{TRADESPERSON_INFO}", tradespersonInfo, -1)
	body = strings.Replace(body, "{ITEMS_AND_PRICE}", itemsAndPrice, -1)

	return email(stripeInvoice.CustomerEmail, *stripeInvoice.CustomerName, "Voided Invoice", body)
}
