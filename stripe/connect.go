package stripe

import (
	"log"
	"os"
	"redbudway-api/restapi/operations"

	"github.com/stripe/stripe-go/v72"
	"github.com/stripe/stripe-go/v72/account"
	"github.com/stripe/stripe-go/v72/accountlink"
)

func CreateTradespersonStripeAccount(tradesperson operations.PostTradespersonBody) (*stripe.Account, error) {
	log.Println("Creating tradesperson stripe connect account")
	params := &stripe.AccountParams{
		Type:    stripe.String("express"),
		Country: stripe.String("US"),
		Email:   stripe.String(tradesperson.Email.String()),
		BusinessProfile: &stripe.AccountBusinessProfileParams{
			Name: stripe.String(*tradesperson.Name),
			SupportAddress: &stripe.AddressParams{
				City:       stripe.String(tradesperson.Address.City),
				Line1:      stripe.String(tradesperson.Address.LineOne),
				Line2:      stripe.String(tradesperson.Address.LineTwo),
				PostalCode: stripe.String(tradesperson.Address.ZipCode),
				State:      stripe.String(tradesperson.Address.State),
			},
		},
		Settings: &stripe.AccountSettingsParams{
			Payouts: &stripe.AccountSettingsPayoutsParams{
				DebitNegativeBalances: stripe.Bool(true),
				Schedule: &stripe.PayoutScheduleParams{
					DelayDays: stripe.Int64(7),
				},
			},
		},
	}
	return account.New(params)
}

func GetConnectAccount(stripeID string) (*stripe.Account, error) {
	return account.GetByID(stripeID, nil)
}

func GetOnBoardingLink(stripeID string) (*stripe.AccountLink, error) {
	log.Print("Creating stripe connect account onboarding link")
	params := &stripe.AccountLinkParams{
		Account:    stripe.String(stripeID),
		RefreshURL: stripe.String("https://" + os.Getenv("SUBDOMAIN") + "redbudway.com"),
		ReturnURL:  stripe.String("https://" + os.Getenv("SUBDOMAIN") + "redbudway.com"),
		Type:       stripe.String("account_onboarding"),
	}
	return accountlink.New(params)
}
