package email

import (
	_ "embed"
	"fmt"
	"log"
	"redbudway-api/internal"
	"redbudway-api/models"
	"strings"

	"github.com/stripe/stripe-go/v82"
)

//go:embed html/tradesperson-booking.html
var tradespersonBooking string

//go:embed html/tradesperson-quote-request.html
var tradespersonQuote string

//go:embed html/tradesperson-quote-accepted.html
var tradespersonQuoteAccepted string

//go:embed html/tradesperson-subscription-booking.html
var tradespersonSubscriptionBooking string

//go:embed html/fixed-price-review.html
var fixedPriceReview string

//go:embed html/welcome.html
var welcomeMessage string

func SendProviderWelcome(accountEmail string) error {
	return email(accountEmail, accountEmail, "Welcome to Redbud Way", welcomeMessage)
}

func SendTradespersonMessage(businessName, businessEmail, service, message string, stripeCustomer *stripe.Customer, images []string) ([]string, error) {
	images, err := internal.ProcessEmailImages(stripeCustomer.Email, "", images)
	if err != nil {
		log.Printf("Failed to process email images, %s", err)
		return images, nil
	}
	replyTo := formatMailbox(stripeCustomer.Email, stripeCustomer.Name)
	return images, sendTextResendWithAttachments(businessEmail, businessName, service, message, replyTo, images)
}

func SendTradespersonBooking(tradesperson models.Tradesperson, stripeCustomer *stripe.Customer, stripeProduct *stripe.Product, timeAndPrice, formRowsCols string) error {
	body := tradespersonBooking

	customerInfo := fmt.Sprintf("%s<br>%s<br>%s", stripeCustomer.Name, stripeCustomer.Email, stripeCustomer.Phone)
	body = strings.Replace(body, "{CUSTOMER_INFO}", customerInfo, -1)
	body = strings.Replace(body, "{SERVICE_NAME}", stripeProduct.Name, -1)
	body = strings.Replace(body, "{TIME_AND_PRICE}", timeAndPrice, -1)
	body = strings.Replace(body, "{FORM_ROWS_COLS}", formRowsCols, -1)

	return email(tradesperson.Email, tradesperson.Name, "Booking", body)
}

func SendTradespersonReview(tradesperson models.Tradesperson, serviceName string) error {
	body := fixedPriceReview

	body = strings.Replace(body, "{SERVICE_NAME}", serviceName, -1)

	return email(tradesperson.Email, tradesperson.Name, "Review", body)
}

func SendTradespersonSubscriptionBooking(tradesperson models.Tradesperson, stripeCustomer *stripe.Customer, stripeProduct *stripe.Product, timeAndPrice string) error {
	body := tradespersonSubscriptionBooking

	customerInfo := fmt.Sprintf("%s<br>%s<br>%s", stripeCustomer.Name, stripeCustomer.Email, stripeCustomer.Phone)
	body = strings.Replace(body, "{CUSTOMER_INFO}", customerInfo, -1)
	body = strings.Replace(body, "{SERVICE_NAME}", stripeProduct.Name, -1)
	body = strings.Replace(body, "{TIME_AND_PRICE}", timeAndPrice, -1)

	return email(tradesperson.Email, tradesperson.Name, "Booking", body)
}

func SendTradespersonQuoteRequest(tradesperson models.Tradesperson, stripeCustomer *stripe.Customer, message, quoteID string, quote *models.ServiceDetails, images []string) ([]string, error) {
	images, err := internal.ProcessEmailImages(stripeCustomer.Email, quoteID, images)
	if err != nil {
		log.Printf("Failed to process email images, %s", err)
		return images, err
	}

	body := tradespersonQuote

	customerInfo := fmt.Sprintf("%s<br>%s<br>%s", stripeCustomer.Name, stripeCustomer.Email, stripeCustomer.Phone)
	body = strings.Replace(body, "{CUSTOMER_MESSAGE}", message, -1)
	body = strings.Replace(body, "{SERVICE_NAME}", *quote.Title, -1)
	body = strings.Replace(body, "{CUSTOMER_INFO}", customerInfo, -1)

	return images, SendProviderQuoteRequest(tradesperson.Email, tradesperson.Name, "Quote Request", body, images)
}

func SendTradespersonQuoteAccepted(tradesperson models.Tradesperson, stripeCustomer *stripe.Customer, message string, quote *models.ServiceDetails) error {
	body := tradespersonQuoteAccepted

	customerInfo := fmt.Sprintf("%s<br>%s<br>%s", stripeCustomer.Name, stripeCustomer.Email, stripeCustomer.Phone)
	body = strings.Replace(body, "{CUSTOMER_MESSAGE}", message, -1)
	body = strings.Replace(body, "{SERVICE_NAME}", *quote.Title, -1)
	body = strings.Replace(body, "{CUSTOMER_INFO}", customerInfo, -1)

	return email(tradesperson.Email, tradesperson.Name, "Quote Accepted", body)
}
