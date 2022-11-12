package email

import (
	_ "embed"
	"fmt"
	"redbudway-api/internal"
	"redbudway-api/models"
	"strings"

	"github.com/go-gomail/gomail"
	"github.com/stripe/stripe-go/v72"
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

func SendTradespersonMessage(businessName, businessEmail, service, message string, stripeCustomer *stripe.Customer, images []string) ([]string, error) {
	m := gomail.NewMessage()
	m.SetAddressHeader("To", businessEmail, businessName)
	m.SetAddressHeader("From", stripeCustomer.Email, stripeCustomer.Name)
	m.SetHeader("Subject", service)
	m.SetBody("text/plain", message)

	images, _ = internal.ProcessEmailImages(stripeCustomer.Email, images)
	for _, image := range images {
		m.Attach(image)
	}

	d := gomail.NewDialer("mail.redbudway.com", 25, "help@redbudway.com", "MerCedEsAmgGt22$")

	return images, d.DialAndSend(m)
}

func SendTradespersonBooking(tradesperson models.Tradesperson, stripeCustomer *stripe.Customer, stripeProduct *stripe.Product, timeAndPrice string) error {
	body := tradespersonBooking

	customerInfo := fmt.Sprintf("%s<br>%s<br>%s", stripeCustomer.Name, stripeCustomer.Email, stripeCustomer.Phone)
	body = strings.Replace(body, "{CUSTOMER_INFO}", customerInfo, -1)
	body = strings.Replace(body, "{SERVICE_NAME}", stripeProduct.Name, -1)
	body = strings.Replace(body, "{TIME_AND_PRICE}", timeAndPrice, -1)

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

func SendTradespersonQuoteRequest(tradesperson models.Tradesperson, stripeCustomer *stripe.Customer, message string, quote *models.ServiceDetails, images []string) ([]string, error) {
	m := gomail.NewMessage()
	m.SetHeader("From", "help@redbudway.com")
	m.SetAddressHeader("To", tradesperson.Email, tradesperson.Name)
	m.SetHeader("Subject", "Quote Request")

	images, _ = internal.ProcessEmailImages(stripeCustomer.Email, images)
	for _, image := range images {
		m.Attach(image)
	}

	body := tradespersonQuote

	customerInfo := fmt.Sprintf("%s<br>%s<br>%s", stripeCustomer.Name, stripeCustomer.Email, stripeCustomer.Phone)
	body = strings.Replace(body, "{CUSTOMER_MESSAGE}", message, -1)
	body = strings.Replace(body, "{SERVICE_NAME}", *quote.Title, -1)
	body = strings.Replace(body, "{CUSTOMER_INFO}", customerInfo, -1)

	m.SetBody("text/html", body)

	d := gomail.NewDialer("mail.redbudway.com", 25, "help@redbudway.com", "MerCedEsAmgGt22$")

	return images, d.DialAndSend(m)
}

func SendTradespersonQuoteAccepted(tradesperson models.Tradesperson, stripeCustomer *stripe.Customer, message string, quote *models.ServiceDetails) error {
	body := tradespersonQuoteAccepted

	customerInfo := fmt.Sprintf("%s<br>%s<br>%s", stripeCustomer.Name, stripeCustomer.Email, stripeCustomer.Phone)
	body = strings.Replace(body, "{CUSTOMER_MESSAGE}", message, -1)
	body = strings.Replace(body, "{SERVICE_NAME}", *quote.Title, -1)
	body = strings.Replace(body, "{CUSTOMER_INFO}", customerInfo, -1)

	return email(tradesperson.Email, tradesperson.Name, "Quote Accepted", body)
}
