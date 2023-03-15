package handlers

import (
	"database/sql"
	"log"
	"time"

	"github.com/go-openapi/runtime/middleware"

	"redbudway-api/database"
	"redbudway-api/internal"
	"redbudway-api/models"
	"redbudway-api/restapi/operations"
	_stripe "redbudway-api/stripe"

	"github.com/stripe/stripe-go/v72"
	"github.com/stripe/stripe-go/v72/account"
	"github.com/stripe/stripe-go/v72/invoice"
	"github.com/stripe/stripe-go/v72/sub"
)

func PostTradespersonHandler(params operations.PostTradespersonParams) middleware.Responder {
	tradesperson := params.Tradesperson

	db := database.GetConnection()

	payload := operations.PostTradespersonCreatedBody{Created: false}
	response := operations.NewPostTradespersonCreated().WithPayload(&payload)

	stmt, err := db.Prepare("SELECT email FROM tradesperson_account WHERE email=?")
	if err != nil {
		log.Printf("Failed to create select statement %s", err)
		return response
	}
	defer stmt.Close()

	row := stmt.QueryRow(tradesperson.Email)
	var email string
	switch err = row.Scan(&email); err {
	case sql.ErrNoRows:
		stripeAccount, err := _stripe.CreateTradespersonStripeAccount(tradesperson)
		if err != nil {
			log.Printf("Failed creating tradesperson stripe connect account %s", err)
			return response
		}
		tradespersonID, err := database.CreateTradespersonAccount(tradesperson, stripeAccount)
		if err != nil {
			log.Printf("Failed creating tradesperson account %s", err)
			return response
		}
		onBoarding, err := _stripe.GetOnBoardingLink(stripeAccount.ID)
		if err != nil {
			log.Printf("Failed creating tradesperson onboarding link %s", err)
			return response
		}
		payload.Created = true
		payload.TradespersonID = tradespersonID.String()
		payload.URL = onBoarding.URL

		accessToken, err := internal.GenerateToken(tradespersonID.String(), "tradesperson", "access", time.Minute*15)
		if err != nil {
			log.Printf("Failed to generate JWT, %s", err)
			return response
		}
		payload.AccessToken = accessToken

		refreshToken, err := internal.GenerateToken(tradespersonID.String(), "tradesperson", "refresh", time.Minute*20)
		if err != nil {
			log.Printf("Failed to generate JWT, %s", err)
			return response
		}
		payload.RefreshToken = refreshToken

		saved, err := database.SaveTradespersonTokens(tradespersonID.String(), refreshToken, accessToken)
		if err != nil {
			log.Printf("Failed to save tradesperson tokens, %s", err)
			return response
		}
		if !saved {
			log.Printf("No issues, but failed to save tradesperson")
		}
		response.SetPayload(&payload)
	case nil:
		log.Printf("Tradesperson with email %s already exist", email)
	default:
		log.Printf("Unknown %v", err)
	}

	return response
}

func DeleteTradespersonTradespersonIDHandler(params operations.DeleteTradespersonTradespersonIDParams, principal interface{}) middleware.Responder {
	tradespersonID := params.TradespersonID
	token := params.HTTPRequest.Header.Get("Authorization")

	payload := operations.DeleteTradespersonTradespersonIDOKBody{Deleted: false}
	response := operations.NewDeleteTradespersonTradespersonIDOK().WithPayload(&payload)

	valid, err := ValidateTradespersonAccessToken(tradespersonID, token)
	if err != nil {
		log.Printf("Failed to validate tradesperson %s, accessToken %s", tradespersonID, token)
		return response
	} else if !valid {
		log.Printf("Bad actor tradesperson %s, accessToken %s", tradespersonID, token)
		return response
	}

	//Handle open invoices, subscriptions, quotes?
	subParams := &stripe.SubscriptionListParams{
		Status: "active",
	}
	s := sub.List(subParams)
	for s.Next() {
		return response
	}
	invParams := &stripe.InvoiceListParams{
		Status: stripe.String("open"),
	}
	i := invoice.List(invParams)
	for i.Next() {
		return response
	}

	stripeID, err := database.GetTradespersonStripeID(tradespersonID)
	if err != nil {
		log.Printf("Failed to get tradesperson %s stripe ID, %v", tradespersonID, err)
		return response
	}
	stripeAccount, err := account.Del(stripeID, nil)
	if err != nil {
		log.Printf("Failed to delete tradesperson %s stripe account, %v", &stripeID, err)
		return response
	}
	if stripeAccount.Deleted {
		payload.Deleted, err = database.DeleteTradespersonAccount(tradespersonID, stripeID)
		if err != nil {
			log.Printf("Failed to delete tradesperson database account, %v", tradespersonID, err)
			return response
		}
		response.SetPayload(&payload)
	}

	return response
}

func GetTradespersonTradespersonID(tradespersonID string) *models.Tradesperson {
	db := database.GetConnection()

	tradesperson := &models.Tradesperson{}

	stmt, err := db.Prepare("SELECT name, number, description, image, email, stripeId FROM tradesperson_account WHERE tradespersonId=?")
	if err != nil {
		log.Printf("Failed to create select statement %s", err)
		return tradesperson
	}
	defer stmt.Close()

	var name, number, email, stripeID string
	var description, image sql.NullString
	row := stmt.QueryRow(tradespersonID)
	switch err = row.Scan(&name, &number, &description, &image, &email, &stripeID); err {
	case sql.ErrNoRows:
		log.Printf("Tradesperson with ID %s doesn't exist", tradespersonID)
	case nil:
		if description.Valid {
			tradesperson.Description = description.String
		}
		if image.Valid {
			tradesperson.Image = image.String
		}
		tradesperson.Name = name
		tradesperson.Number = number
		tradesperson.Email = email

		jobs, err := database.GetTradespersonJobs(tradespersonID)
		if err != nil {
			log.Printf("Failed to get tradesperson job count %s", err)
			return tradesperson
		}
		tradesperson.Jobs = jobs

		rating, reviews, err := database.GetTradespersonRatingReviews(tradespersonID)
		if err != nil {
			log.Printf("Failed to get tradesperson rating & reviews %s", err)
			return tradesperson
		}
		tradesperson.Rating = rating
		tradesperson.Reviews = reviews

		connect, err := _stripe.GetConnectAccount(stripeID)
		if err != nil {
			log.Print("Failed to get stripe connect account for tradesperson with ID %s", tradespersonID)
			return tradesperson
		}

		tradesperson.Address = &models.Address{}
		tradesperson.Address.City = connect.BusinessProfile.SupportAddress.City
		tradesperson.Address.State = connect.BusinessProfile.SupportAddress.State
		tradesperson.Address.LineOne = connect.BusinessProfile.SupportAddress.Line1
		tradesperson.Address.LineTwo = connect.BusinessProfile.SupportAddress.Line2
		tradesperson.Address.ZipCode = connect.BusinessProfile.SupportAddress.PostalCode
	default:
		log.Printf("Unknown %v", err)
	}

	return tradesperson
}

func GetTradespersonTradespersonIDHandler(params operations.GetTradespersonTradespersonIDParams, principal interface{}) middleware.Responder {
	tradespersonID := params.TradespersonID
	token := params.HTTPRequest.Header.Get("Authorization")

	response := operations.NewGetTradespersonTradespersonIDOK()

	valid, err := ValidateTradespersonAccessToken(tradespersonID, token)
	if err != nil {
		log.Printf("Failed to validate tradesperson %s, accessToken %s", tradespersonID, token)
		return response
	} else if !valid {
		log.Printf("Bad actor tradesperson %s, accessToken %s", tradespersonID, token)
		return response
	}

	tradeperson := GetTradespersonTradespersonID(tradespersonID)
	response.SetPayload(tradeperson)

	return response
}

func PutTradespersonTradespersonIDHandler(params operations.PutTradespersonTradespersonIDParams, principal interface{}) middleware.Responder {
	tradespersonID := params.TradespersonID
	tradesperson := params.Tradesperson
	token := params.HTTPRequest.Header.Get("Authorization")

	payload := operations.PutTradespersonTradespersonIDOKBody{}
	updated := false
	payload.Updated = updated
	response := operations.NewPutTradespersonTradespersonIDOK()
	response.SetPayload(&payload)

	valid, err := ValidateTradespersonAccessToken(tradespersonID, token)
	if err != nil {
		log.Printf("Failed to validate tradesperson %s, accessToken %s", tradespersonID, token)
		return response
	} else if !valid {
		log.Printf("Bad actor tradesperson %s, accessToken %s", tradespersonID, token)
		return response
	}

	if err := database.UpdateTradespersonDescription(tradespersonID, tradesperson.Description); err != nil {
		log.Printf("Failed to update tradesperson description, %v", err)
	}

	if tradesperson.Image != "" {
		if err := database.UpdateTradespersonImage(tradespersonID, tradesperson.Image); err != nil {
			log.Printf("Failed to update tradesperson image, %v", err)
		}
	}

	payload.Updated = true
	response = operations.NewPutTradespersonTradespersonIDOK()
	response.SetPayload(&payload)
	return response
}

func GetTradespersonTradespersonIDSettingsHandler(params operations.GetTradespersonTradespersonIDSettingsParams, principal interface{}) middleware.Responder {
	tradespersonID := params.TradespersonID
	token := params.HTTPRequest.Header.Get("Authorization")

	payload := operations.GetTradespersonTradespersonIDSettingsOKBody{}
	response := operations.NewGetTradespersonTradespersonIDSettingsOK()
	response.SetPayload(&payload)

	valid, err := ValidateTradespersonAccessToken(tradespersonID, token)
	if err != nil {
		log.Printf("Failed to validate tradesperson %s, accessToken %s", tradespersonID, token)
		return response
	} else if !valid {
		log.Printf("Bad actor tradesperson %s, accessToken %s", tradespersonID, token)
		return response
	}

	db := database.GetConnection()

	stmt, err := db.Prepare("SELECT vanityURL, number, email, address FROM tradesperson_settings WHERE tradespersonId=?")
	if err != nil {
		log.Printf("Failed to create select statement %s", err)
		return response
	}
	defer stmt.Close()

	var vanityURL sql.NullString
	var displayNumber, displayEmail, displayAddress bool
	row := stmt.QueryRow(tradespersonID)
	switch err = row.Scan(&vanityURL, &displayNumber, &displayEmail, &displayAddress); err {
	case sql.ErrNoRows:
		log.Printf("Tradesperson with ID %s doesn't exist", tradespersonID)
	case nil:
		if vanityURL.Valid {
			payload.VanityURL = vanityURL.String
		}
		payload.DisplayNumber = displayNumber
		payload.DisplayEmail = displayEmail
		payload.DisplayAddress = displayAddress

		response.SetPayload(&payload)
	default:
		log.Printf("Unknown %v", err)
	}
	return response
}

func PutTradespersonTradespersonIDSettingsHandler(params operations.PutTradespersonTradespersonIDSettingsParams, principal interface{}) middleware.Responder {
	tradespersonID := params.TradespersonID
	settings := params.Settings
	token := params.HTTPRequest.Header.Get("Authorization")

	response := operations.NewPutTradespersonTradespersonIDSettingsOK()

	valid, err := ValidateTradespersonAccessToken(tradespersonID, token)
	if err != nil {
		log.Printf("Failed to validate tradesperson %s, accessToken %s", tradespersonID, token)
		return response
	} else if !valid {
		log.Printf("Bad actor tradesperson %s, accessToken %s", tradespersonID, token)
		return response
	}

	updated, err := database.UpdateTradespersonDisplaySettings(tradespersonID, settings)
	if err != nil {
		log.Printf("Failed to update tradesperson display settings %s", err)
	}
	if updated {
		updated, err = database.UpdateTradespersonVanitySettings(tradespersonID, settings)
		if err != nil {
			log.Printf("Failed to update tradesperson vanity settings %s", err)
		}
	}
	payload := operations.PutTradespersonTradespersonIDSettingsOKBody{Updated: updated}

	response.SetPayload(&payload)

	return response
}

func GetTradespersonTradespersonIDStatusHandler(params operations.GetTradespersonTradespersonIDStatusParams, principal interface{}) middleware.Responder {
	tradespersonID := params.TradespersonID
	token := params.HTTPRequest.Header.Get("Authorization")

	payload := operations.GetTradespersonTradespersonIDStatusOKBody{Enabled: false, Submitted: false}
	response := operations.NewGetTradespersonTradespersonIDStatusOK()
	response.SetPayload(&payload)

	valid, err := ValidateTradespersonAccessToken(tradespersonID, token)
	if err != nil {
		log.Printf("Failed to validate tradesperson %s, accessToken %s", tradespersonID, token)
		return response
	} else if !valid {
		log.Printf("Bad actor tradesperson %s, accessToken %s", tradespersonID, token)
		return response
	}

	db := database.GetConnection()

	stmt, err := db.Prepare("SELECT stripeId FROM tradesperson_account WHERE tradespersonId=?")
	if err != nil {
		log.Printf("Failed to create select statement %s", err)
		return response
	}
	defer stmt.Close()

	row := stmt.QueryRow(tradespersonID)
	var stripeID string
	switch err = row.Scan(&stripeID); err {
	case sql.ErrNoRows:
		log.Printf("Tradesperson with ID %s doesn't exist", tradespersonID)
	case nil:
		connect, err := _stripe.GetConnectAccount(stripeID)
		if err != nil {
			log.Print("Failed to get stripe account for tradesperson with ID %s", tradespersonID)
			return response
		}
		payload.Enabled = connect.ChargesEnabled
		payload.Submitted = connect.DetailsSubmitted
		response.SetPayload(&payload)
	default:
		log.Printf("Unknown default switch case, %v", err)
	}

	return response
}

func PutTradespersonTradespersonIDPasswordHandler(params operations.PutTradespersonTradespersonIDPasswordParams, principal interface{}) middleware.Responder {
	tradesperson := params.Tradesperson
	tradespersonID := params.TradespersonID
	token := params.HTTPRequest.Header.Get("Authorization")

	response := operations.NewPutTradespersonTradespersonIDPasswordOK()
	payload := &operations.PutTradespersonTradespersonIDPasswordOKBody{}
	payload.Updated = false
	response.SetPayload(payload)

	valid, err := ValidateTradespersonAccessToken(tradespersonID, token)
	if err != nil {
		log.Printf("Failed to validate tradesperson %s, accessToken %s", tradespersonID, token)
		return response
	} else if !valid {
		log.Printf("Bad actor tradesperson %s, accessToken %s", tradespersonID, token)
		return response
	}

	db := database.GetConnection()

	stmt, err := db.Prepare("SELECT name, email, password FROM tradesperson_account WHERE tradespersonId=?")
	if err != nil {
		log.Printf("Failed to create select statement %s", err)
		return response
	}
	defer stmt.Close()

	row := stmt.QueryRow(tradespersonID)
	var name, email, hashPassword string
	switch err = row.Scan(&name, &email, &hashPassword); err {
	case sql.ErrNoRows:
		log.Printf("Tradesperson with ID %s doesn't exist", tradespersonID)
	case nil:
		if internal.CheckPasswordHash(*tradesperson.CurPassword, hashPassword) {
			stmt, err := db.Prepare("UPDATE tradesperson_account SET password=? WHERE tradespersonId = ?")
			if err != nil {
				return response
			}
			defer stmt.Close()

			newHashPassword, err := internal.HashPassword(*tradesperson.NewPassword)
			if err != nil {
				log.Printf("%s", err)
				return response
			}
			results, err := stmt.Exec(newHashPassword, tradespersonID)
			if err != nil {
				return response
			}

			rowsAffected, err := results.RowsAffected()
			if err != nil {
				return response
			}

			if rowsAffected == 1 {
				payload.Updated = true
				response.SetPayload(payload)
			}
		}
	default:
		log.Printf("Unknown %v", err)
	}

	return response
}

func GetTradespersonTradespersonIDOnboardHandler(params operations.GetTradespersonTradespersonIDOnboardParams, principal interface{}) middleware.Responder {
	tradespersonID := params.TradespersonID
	token := params.HTTPRequest.Header.Get("Authorization")

	response := operations.NewGetTradespersonTradespersonIDOnboardOK()
	payload := operations.GetTradespersonTradespersonIDOnboardOKBody{}
	response.SetPayload(&payload)

	valid, err := ValidateTradespersonAccessToken(tradespersonID, token)
	if err != nil {
		log.Printf("Failed to validate tradesperson %s, accessToken %s", tradespersonID, token)
		return response
	} else if !valid {
		log.Printf("Bad actor tradesperson %s, accessToken %s", tradespersonID, token)
		return response
	}

	db := database.GetConnection()

	stmt, err := db.Prepare("SELECT stripeId FROM tradesperson_account WHERE tradespersonId=?")
	if err != nil {
		log.Printf("Failed to create select statement %s", err)
		return response
	}
	defer stmt.Close()

	var stripeID string
	row := stmt.QueryRow(tradespersonID)
	switch err = row.Scan(&stripeID); err {
	case sql.ErrNoRows:
		log.Printf("Tradesperson with ID %v doesn't exist", tradespersonID)
	case nil:
		onBoarding, err := _stripe.GetOnBoardingLink(stripeID)
		if err != nil {
			log.Printf("Failed creating tradesperson onboarding link %s", err)
			return response
		}
		payload.URL = onBoarding.URL
		response.SetPayload(&payload)
	default:
		log.Printf("Unknown %v", err)
	}

	return response
}

func GetTradespersonTradespersonIDLoginLinkHandler(params operations.GetTradespersonTradespersonIDLoginLinkParams, principal interface{}) middleware.Responder {
	tradespersonID := params.TradespersonID
	token := params.HTTPRequest.Header.Get("Authorization")
	response := operations.NewGetTradespersonTradespersonIDLoginLinkOK()

	valid, err := ValidateTradespersonAccessToken(tradespersonID, token)
	if err != nil {
		log.Printf("Failed to validate tradesperson %s, accessToken %s", tradespersonID, token)
		return response
	} else if !valid {
		log.Printf("Bad actor tradesperson %s, accessToken %s", tradespersonID, token)
		return response
	}

	stripeID, err := database.GetTradespersonStripeID(tradespersonID)
	if err != nil {
		log.Printf("Failed to get tradesperson %s stripe ID, %v", tradespersonID, err)
		return response
	}

	loginLink, err := _stripe.GetTradespersonLoginLink(stripeID)
	if err != nil {
		log.Printf("Failed to get tradesperson login link, %v", err)
		return response
	}
	payload := operations.GetTradespersonTradespersonIDLoginLinkOKBody{}
	payload.URL = loginLink.URL
	response.SetPayload(&payload)

	return response
}

func GetTradespersonTradespersonIDSellingFeeHandler(params operations.GetTradespersonTradespersonIDSellingFeeParams, principal interface{}) middleware.Responder {

	return operations.NewGetTradespersonTradespersonIDSellingFeeOK()
}
