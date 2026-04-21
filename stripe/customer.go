package stripe

import (
	"os"
	"redbudway-api/restapi/operations"

	"github.com/stripe/stripe-go/v82"
	"github.com/stripe/stripe-go/v82/billingportal/session"
	"github.com/stripe/stripe-go/v82/customer"
)

func CreateCustomerStripeAccount(_customer operations.PostCustomerBody) (*stripe.Customer, error) {
	params := &stripe.CustomerParams{
		Name:  stripe.String(*_customer.Name),
		Email: stripe.String(_customer.Email.String()),
		Phone: stripe.String(*_customer.Number),
	}
	if _customer.Address != nil {
		a := _customer.Address
		if a.LineOne != "" || a.LineTwo != "" || a.City != "" || a.State != "" || a.ZipCode != "" {
			params.Address = &stripe.AddressParams{
				City:       stripe.String(a.City),
				Line1:      stripe.String(a.LineOne),
				Line2:      stripe.String(a.LineTwo),
				PostalCode: stripe.String(a.ZipCode),
				State:      stripe.String(a.State),
			}
		}
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
