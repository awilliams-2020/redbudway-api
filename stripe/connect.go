package stripe

import (
	"log"
	"os"
	"redbudway-api/models"
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

func GetOnBoardingLink(stripeID, tradespersonID string) (*stripe.AccountLink, error) {
	log.Print("Creating stripe connect account onboarding link")
	params := &stripe.AccountLinkParams{
		Account:    stripe.String(stripeID),
		RefreshURL: stripe.String("https://" + os.Getenv("SUBDOMAIN") + "redbudway.com/#/tradesperson/" + tradespersonID + "/profile"),
		ReturnURL:  stripe.String("https://" + os.Getenv("SUBDOMAIN") + "redbudway.com/#/tradesperson/" + tradespersonID + "/profile"),
		Type:       stripe.String("account_onboarding"),
	}
	return accountlink.New(params)
}

func UpdateBusinessProfileName(stripeID, name string) error {
	params := &stripe.AccountParams{}
	params.BusinessProfile = &stripe.AccountBusinessProfileParams{
		Name: stripe.String(name),
	}
	_, err := account.Update(
		stripeID,
		params,
	)
	if err != nil {
		return err
	}
	return nil
}

func UpdateBusinessProfileNumber(stripeID, number string) error {
	params := &stripe.AccountParams{}
	params.BusinessProfile = &stripe.AccountBusinessProfileParams{
		SupportPhone: stripe.String(number),
	}
	_, err := account.Update(
		stripeID,
		params,
	)
	if err != nil {
		return err
	}
	return nil
}

func UpdateBusinessProfileEmail(stripeID, email string) error {
	params := &stripe.AccountParams{}
	params.BusinessProfile = &stripe.AccountBusinessProfileParams{
		SupportEmail: stripe.String(email),
	}
	_, err := account.Update(
		stripeID,
		params,
	)
	if err != nil {
		return err
	}
	return nil
}

func UpdateBusinessProfileAddress(stripeID string, address *models.Address) error {
	params := &stripe.AccountParams{}
	params.BusinessProfile = &stripe.AccountBusinessProfileParams{
		SupportAddress: &stripe.AddressParams{
			City:       stripe.String(address.City),
			State:      stripe.String(address.State),
			PostalCode: stripe.String(address.ZipCode),
			Line1:      stripe.String(address.LineOne),
			Line2:      stripe.String(address.LineTwo),
		},
	}
	_, err := account.Update(
		stripeID,
		params,
	)
	if err != nil {
		return err
	}
	return nil
}
