package handlers

import (
	"github.com/stripe/stripe-go/v82"
)

// defaultSellingFeeFraction is used when the selling_fee table has no row for a tradesperson.
// 3% is a typical marketplace take for service trades (Stripe processing is paid by the connected account on direct charges).
const defaultSellingFeeFraction = 0.03

// billingQuoteParams scopes Quote API calls to a Connect account (direct charges).
func billingQuoteParams(connectAccountID string, expand []*string) *stripe.QuoteParams {
	p := &stripe.QuoteParams{
		Params: stripe.Params{
			Expand: expand,
		},
	}
	if connectAccountID != "" {
		p.SetStripeAccount(connectAccountID)
	}
	return p
}
