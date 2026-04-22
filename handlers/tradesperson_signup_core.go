package handlers

import (
	"database/sql"
	"log"
	"time"

	"redbudway-api/database"
	"redbudway-api/email"
	"redbudway-api/internal"
	"redbudway-api/restapi/operations"
	_stripe "redbudway-api/stripe"
)

// executeProviderSignup creates Stripe Connect, DB rows, JWTs, and welcome email when the email is new.
// Callers must verify the user first (reCAPTCHA or Google ID token).
func executeProviderSignup(tradesperson operations.PostTradespersonBody) *operations.PostTradespersonCreatedBody {
	payload := &operations.PostTradespersonCreatedBody{Created: false}
	db := database.GetConnection()
	stmt, err := db.Prepare("SELECT email FROM tradesperson_account WHERE email=?")
	if err != nil {
		log.Printf("Failed to create select statement %s", err)
		return payload
	}
	defer stmt.Close()

	row := stmt.QueryRow(tradesperson.Email)
	var _email string
	switch err = row.Scan(&_email); err {
	case sql.ErrNoRows:
		stripeAccount, err := _stripe.CreateTradespersonStripeAccount(tradesperson)
		if err != nil {
			log.Printf("Failed creating tradesperson stripe connect account %s", err)
			return payload
		}
		tradespersonID, err := database.CreateTradespersonAccount(tradesperson, stripeAccount)
		if err != nil {
			log.Printf("Failed creating tradesperson account %s", err)
			return payload
		}
		onBoarding, err := _stripe.GetOnBoardingLink(stripeAccount.ID)
		if err != nil {
			log.Printf("Failed creating tradesperson onboarding link %s", err)
			return payload
		}
		payload.Created = true
		payload.TradespersonID = tradespersonID.String()
		payload.URL = onBoarding.URL

		accessToken, err := internal.GenerateToken(tradespersonID.String(), "tradesperson", "access", time.Minute*15)
		if err != nil {
			log.Printf("Failed to generate JWT, %s", err)
			return payload
		}
		payload.AccessToken = accessToken

		refreshToken, err := internal.GenerateToken(tradespersonID.String(), "tradesperson", "refresh", time.Minute*20)
		if err != nil {
			log.Printf("Failed to generate JWT, %s", err)
			return payload
		}
		payload.RefreshToken = refreshToken

		saved, err := database.SaveTradespersonTokens(tradespersonID.String(), refreshToken, accessToken)
		if err != nil {
			log.Printf("Failed to save tradesperson tokens, %s", err)
			return payload
		}
		if !saved {
			log.Printf("No issues, but failed to save tradesperson")
		}
		if err := email.SendProviderWelcome(tradesperson.Email.String()); err != nil {
			log.Printf("Failed to send tradesperson welcome email, %v", err)
		}
		return payload
	case nil:
		log.Printf("Tradesperson with email %s already exist", _email)
		return payload
	default:
		log.Printf("Unknown %v", err)
		return payload
	}
}
