package stripe

import (
	"os"
	"redbudway-api/restapi/operations"

	"github.com/stripe/stripe-go/v72"
	"github.com/stripe/stripe-go/v72/billingportal/session"
	"github.com/stripe/stripe-go/v72/customer"
)

func CreateCustomerStripeAccount(_customer operations.PostCustomerBody) (*stripe.Customer, error) {
	params := &stripe.CustomerParams{
		Name:  stripe.String(*_customer.Name),
		Email: stripe.String(_customer.Email.String()),
		Address: &stripe.AddressParams{
			City:       stripe.String(_customer.Address.City),
			Line1:      stripe.String(_customer.Address.LineOne),
			Line2:      stripe.String(_customer.Address.LineTwo),
			PostalCode: stripe.String(_customer.Address.ZipCode),
			State:      stripe.String(_customer.Address.State),
		},
		Phone: stripe.String(*_customer.Number),
	}
	return customer.New(params)
}

func GetCustomerBillingLink(stripeID string) (*stripe.BillingPortalSession, error) {
	params := &stripe.BillingPortalSessionParams{
		Customer:  stripe.String(stripeID),
		ReturnURL: stripe.String("https://" + os.Getenv("SUBDOMAIN") + "redbudway.com"),
	}
	return session.New(params)
}
