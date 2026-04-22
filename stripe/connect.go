package stripe

import (
	"fmt"
	"log"
	"os"
	"redbudway-api/models"
	"redbudway-api/restapi/operations"
	"strings"

	"github.com/stripe/stripe-go/v82"
	"github.com/stripe/stripe-go/v82/account"
	"github.com/stripe/stripe-go/v82/accountlink"
	"github.com/stripe/stripe-go/v82/file"
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
				Schedule: &stripe.AccountSettingsPayoutsScheduleParams{
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

// ConnectAccountDisplayName picks a customer-facing business name from a Connect Account.
// Stripe stores the public name in business_profile.name, but legal-entity edits (and some
// Express Dashboard flows) update company.name or individual instead — only reading
// business_profile.name misses those updates.
func ConnectAccountDisplayName(ac *stripe.Account) string {
	if ac == nil {
		return ""
	}
	if ac.BusinessProfile != nil {
		if n := strings.TrimSpace(ac.BusinessProfile.Name); n != "" {
			return n
		}
	}
	if ac.Company != nil {
		if n := strings.TrimSpace(ac.Company.Name); n != "" {
			return n
		}
	}
	if ac.Individual != nil {
		first := strings.TrimSpace(ac.Individual.FirstName)
		last := strings.TrimSpace(ac.Individual.LastName)
		if first != "" || last != "" {
			return strings.TrimSpace(first + " " + last)
		}
	}
	return ""
}

func GetOnBoardingLink(stripeID string) (*stripe.AccountLink, error) {
	log.Print("Creating stripe connect account onboarding link")
	base := "https://" + os.Getenv("SUBDOMAIN") + "redbudway.com"
	params := &stripe.AccountLinkParams{
		Account:    stripe.String(stripeID),
		RefreshURL: stripe.String(base + "/dashboard"),
		ReturnURL:  stripe.String(base + "/dashboard"),
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

func UpdateBusinessBranding(stripeID, logoURL, iconURL, primary, secondary, tradespersonID string) error {
	params := &stripe.AccountParams{}
	params.Settings = &stripe.AccountSettingsParams{}

	branding := &stripe.AccountSettingsBrandingParams{}
	if logoURL != "" {
		path := fmt.Sprintf("images/%s/logo.png", tradespersonID)
		ID, err := createLogoFile(path)
		if err != nil {
			return err
		}
		branding.Logo = &ID
	}
	if iconURL != "" {
		path := fmt.Sprintf("images/%s/icon.png", tradespersonID)
		ID, err := createIconFile(path)
		if err != nil {
			return err
		}
		branding.Icon = &ID
	}
	if primary != "" {
		branding.PrimaryColor = &primary
	}
	if secondary != "" {
		branding.SecondaryColor = &secondary
	}
	params.Settings.Branding = branding

	_, err := account.Update(
		stripeID,
		params,
	)
	if err != nil {
		return err
	}
	return nil
}

func createLogoFile(filePath string) (string, error) {
	fp, _ := os.Open(filePath)
	params := &stripe.FileParams{
		FileReader: fp,
		Filename:   stripe.String("logo.png"),
		Purpose:    stripe.String(string(stripe.FilePurposeBusinessLogo)),
	}
	f, err := file.New(params)
	if err != nil {
		return f.ID, err
	}
	return f.ID, nil
}

func createIconFile(filePath string) (string, error) {
	fp, _ := os.Open(filePath)
	params := &stripe.FileParams{
		FileReader: fp,
		Filename:   stripe.String("icon.png"),
		Purpose:    stripe.String(string(stripe.FilePurposeBusinessIcon)),
	}
	f, err := file.New(params)
	if err != nil {
		return f.ID, err
	}
	return f.ID, nil
}
