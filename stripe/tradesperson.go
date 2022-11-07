package stripe

import (
	"github.com/stripe/stripe-go/v72"
	"github.com/stripe/stripe-go/v72/loginlink"
)

func GetTradespersonLoginLink(stripeID string) (*stripe.LoginLink, error) {
	params := &stripe.LoginLinkParams{
		Account: stripe.String(stripeID),
	}
	ll, err := loginlink.New(params)
	return ll, err
}
