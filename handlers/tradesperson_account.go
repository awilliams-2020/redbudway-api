package handlers

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/go-openapi/runtime/middleware"

	"redbudway-api/database"
	"redbudway-api/email"
	"redbudway-api/internal"
	"redbudway-api/restapi/operations"
	_stripe "redbudway-api/stripe"

	"github.com/stripe/stripe-go/v72/account"
)

func PostTradespersonHandler(params operations.PostTradespersonParams) middleware.Responder {
	tradesperson := params.Tradesperson

	db := database.GetConnection()

	payload := operations.PostTradespersonCreatedBody{Created: false}
	response := operations.NewPostTradespersonCreated().WithPayload(&payload)
	valid, err := internal.VerifyReCaptcha(*tradesperson.Token)
	if err != nil {
		log.Printf("Verifying recaptcha failed, %v", err)
		return response
	}
	if !valid {
		log.Printf("Signup failed recaptcha")
		return response
	}

	stmt, err := db.Prepare("SELECT email FROM tradesperson_account WHERE email=?")
	if err != nil {
		log.Printf("Failed to create select statement %s", err)
		return response
	}
	defer stmt.Close()

	row := stmt.QueryRow(tradesperson.Email)
	var _email string
	switch err = row.Scan(&_email); err {
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
		onBoarding, err := _stripe.GetOnBoardingLink(stripeAccount.ID, tradespersonID.String())
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
		if err := email.SendProviderWelcome(tradesperson.Email.String()); err != nil {
			log.Printf("Failed to send tradesperson welcome email, %v", err)
		}
		response.SetPayload(&payload)
	case nil:
		log.Printf("Tradesperson with email %s already exist", _email)
	default:
		log.Printf("Unknown %v", err)
	}

	return response
}

func DeleteTradespersonTradespersonIDHandler(params operations.DeleteTradespersonTradespersonIDParams, principal interface{}) middleware.Responder {
	tradespersonID := params.TradespersonID
	token := params.HTTPRequest.Header.Get("Authorization")

	payload := operations.DeleteTradespersonTradespersonIDOKBody{Deleted: true}
	response := operations.NewDeleteTradespersonTradespersonIDOK().WithPayload(&payload)

	valid, err := ValidateTradespersonAccessToken(tradespersonID, token)
	if err != nil {
		log.Printf("Failed to validate tradesperson %s, accessToken %s", tradespersonID, token)
		return response
	} else if !valid {
		log.Printf("Bad actor tradesperson %s, accessToken %s", tradespersonID, token)
		return response
	}
	deleteAccount(tradespersonID)

	return response
}

func deleteAccount(tradespersonID string) {

	stripeID, err := database.GetTradespersonStripeID(tradespersonID)
	if err != nil {
		log.Printf("Failed to get tradesperson %s stripe ID, %v", tradespersonID, err)
		return
	}
	stripeAccount, err := account.Del(stripeID, nil)
	if err != nil {
		log.Printf("Failed to delete tradesperson %s stripe account, %v", &stripeID, err)
		return
	}
	if stripeAccount.Deleted {
		_, err = database.DeleteTradespersonAccount(tradespersonID, stripeID)
		if err != nil {
			log.Printf("Failed to delete tradesperson database account, %v", tradespersonID, err)
			return
		}
	}
}

func GetTradespersonTradespersonIDHandler(params operations.GetTradespersonTradespersonIDParams, principal interface{}) middleware.Responder {
	tradespersonID := params.TradespersonID
	token := params.HTTPRequest.Header.Get("Authorization")

	payload := operations.GetTradespersonTradespersonIDOKBody{Enabled: false, Submitted: false}
	response := operations.NewGetTradespersonTradespersonIDOK()
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

func GetTradespersonTradespersonIDSyncHandler(params operations.GetTradespersonTradespersonIDSyncParams, principal interface{}) middleware.Responder {
	tradespersonID := params.TradespersonID
	token := params.HTTPRequest.Header.Get("Authorization")

	response := operations.NewGetTradespersonTradespersonIDSyncOK()

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
		log.Printf("Failed to get tradesperson %s, stripeID, %v", tradespersonID, err)
		return response
	}
	connect, err := _stripe.GetConnectAccount(stripeID)
	if err != nil {
		log.Printf("Failed to get tradesperson %s, connect account, %v", tradespersonID, err)
		return response
	}

	if connect.BusinessProfile.Name != "" {
		if err := database.UpdateTradespersonProfileName(tradespersonID, connect.BusinessProfile.Name); err != nil {
			log.Printf("Failed to update tradesperson profile name, %v", err)
		}
	} else if connect.BusinessProfile.URL != "" {
		if err := database.UpdateTradespersonProfileName(tradespersonID, connect.BusinessProfile.URL); err != nil {
			log.Printf("Failed to update tradesperson profile name, %v", err)
		}
	}

	// if connect.BusinessProfile.SupportEmail != "" {
	// 	if err := database.UpdateTradespersonProfileEmail(tradespersonID, connect.BusinessProfile.SupportEmail); err != nil {
	// 		log.Printf("Failed to update tradesperson profile email, %v", err)
	// 	}
	// }

	// if connect.BusinessProfile.SupportPhone != "" {
	// 	if err := database.UpdateTradespersonProfileEmail(tradespersonID, connect.BusinessProfile.SupportPhone); err != nil {
	// 		log.Printf("Failed to update tradesperson profile number, %v", err)
	// 	}
	// }

	// if connect.BusinessProfile.SupportAddress != nil {
	// 	address := &models.Address{}
	// 	address.LineOne = connect.BusinessProfile.SupportAddress.Line1
	// 	address.LineTwo = connect.BusinessProfile.SupportAddress.Line2
	// 	address.City = connect.BusinessProfile.SupportAddress.City
	// 	address.State = connect.BusinessProfile.SupportAddress.State
	// 	address.ZipCode = connect.BusinessProfile.SupportAddress.PostalCode
	// 	if err := database.UpdateTradespersonProfileAddress(tradespersonID, address); err != nil {
	// 		log.Printf("Failed to update tradesperson profile address, %v", err)
	// 	}
	// }

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
		onBoarding, err := _stripe.GetOnBoardingLink(stripeID, tradespersonID)
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

func PutTradespersonTradespersonIDHandler(params operations.PutTradespersonTradespersonIDParams, principal interface{}) middleware.Responder {
	tradespersonID := params.TradespersonID
	account := params.Account
	token := params.HTTPRequest.Header.Get("Authorization")

	response := operations.NewPutTradespersonTradespersonIDOK()
	payload := &operations.PutTradespersonTradespersonIDOKBody{}
	payload.Updated = false
	response.SetPayload(payload)

	var err error
	if account.Email != "" && account.CurPassword != "" {
		payload.Updated, err = putTradespersonEmail(tradespersonID, account.Email, account.CurPassword, token)
		if err != nil {
			log.Printf("Failed to update tradesperson account email, %s", err)
		}
	} else if account.NewPassword != "" && account.CurPassword != "" {
		payload.Updated, err = putTradespersonPassword(tradespersonID, account.CurPassword, account.NewPassword, token)
		if err != nil {
			log.Printf("Failed to update tradesperson account password, %s", err)
		}
	} else if account.Email != "" {
		payload.Updated, err = revertTradespersonEmail(tradespersonID, account.Email, token)
		if err != nil {
			log.Printf("Failed to revert tradesperson account email, %s", err)
		}
	}

	response = operations.NewPutTradespersonTradespersonIDOK()
	response.SetPayload(payload)
	return response
}

func putTradespersonPassword(tradespersonID, curPassword, newPassword, token string) (bool, error) {
	updated := false

	valid, err := ValidateTradespersonAccessToken(tradespersonID, token)
	if err != nil {
		log.Printf("Failed to validate tradesperson %s, accessToken %s", tradespersonID, token)
		return updated, err
	} else if !valid {
		return updated, fmt.Errorf("Bad actor tradesperson %s, accessToken %s", tradespersonID, token)
	}

	db := database.GetConnection()
	stmt, err := db.Prepare("SELECT email, password FROM tradesperson_account WHERE tradespersonId=?")
	if err != nil {
		log.Printf("Failed to create select statement %s", err)
		return updated, err
	}
	defer stmt.Close()

	row := stmt.QueryRow(tradespersonID)
	var hashPassword, accountEmail string
	switch err = row.Scan(&accountEmail, &hashPassword); err {
	case sql.ErrNoRows:
		log.Printf("Tradesperson with ID %s doesn't exist", tradespersonID)
	case nil:
		if internal.CheckPasswordHash(curPassword, hashPassword) {
			stmt, err := db.Prepare("UPDATE tradesperson_account SET password=? WHERE tradespersonId = ?")
			if err != nil {
				return updated, err
			}
			defer stmt.Close()

			newHashPassword, err := internal.HashPassword(newPassword)
			if err != nil {
				log.Printf("%s", err)
				return updated, err
			}
			results, err := stmt.Exec(newHashPassword, tradespersonID)
			if err != nil {
				return updated, err
			}

			rowsAffected, err := results.RowsAffected()
			if err != nil {
				return updated, err
			}

			updated = rowsAffected == 1
			if updated {
				tradesperson, err := database.GetTradespersonProfile(tradespersonID)
				if err != nil {
					log.Printf("Failed to get tradesperson profile %s", err)
					return updated, nil
				}
				if err := email.PasswordUpdated(accountEmail, tradesperson.Name); err != nil {
					log.Printf("Failed to send tradesperson email, %v", err)
					return updated, nil
				}
			}
		}
	default:
		log.Printf("Unknown %v", err)
	}

	return updated, nil
}

func putTradespersonEmail(tradespersonID, newEmail, curPassword, token string) (bool, error) {
	updated := false

	valid, err := ValidateTradespersonAccessToken(tradespersonID, token)
	if err != nil {
		log.Printf("Failed to validate tradesperson %s, accessToken %s", tradespersonID, token)
		return updated, err
	} else if !valid {
		return updated, fmt.Errorf("Bad actor tradesperson %s, accessToken %s", tradespersonID, token)
	}

	db := database.GetConnection()

	stmt, err := db.Prepare("SELECT email, password FROM tradesperson_account WHERE tradespersonId=?")
	if err != nil {
		log.Printf("Failed to create select statement %s", err)
		return updated, err
	}
	defer stmt.Close()

	row := stmt.QueryRow(tradespersonID)
	var accountEmail, hashPassword string
	switch err = row.Scan(&accountEmail, &hashPassword); err {
	case sql.ErrNoRows:
		log.Printf("Tradesperson with ID %s doesn't exist", tradespersonID)
	case nil:
		if internal.CheckPasswordHash(curPassword, hashPassword) {
			stmt, err := db.Prepare("UPDATE tradesperson_account SET email=? WHERE tradespersonId = ?")
			if err != nil {
				return updated, err
			}
			defer stmt.Close()

			results, err := stmt.Exec(newEmail, tradespersonID)
			if err != nil {
				return updated, err
			}

			rowsAffected, err := results.RowsAffected()
			if err != nil {
				return updated, err
			}

			updated = rowsAffected == 1
			if updated {
				token, err := internal.GenerateToken(tradespersonID, "tradesperson", accountEmail, time.Hour*1)
				if err != nil {
					log.Printf("Failed to generate JWT, %s", err)
					return updated, err
				}
				tradesperson, err := database.GetTradespersonProfile(tradespersonID)
				if err != nil {
					log.Printf("Failed to get tradesperson profile %s", err)
					return updated, err
				}
				if err := email.EmailUpdated(accountEmail, tradesperson.Name, token, tradespersonID); err != nil {
					log.Printf("Failed to send tradesperson email, %v", err)
					return updated, err
				}
			}
		}
	default:
		log.Printf("Unknown %v", err)
	}

	return updated, nil
}

func revertTradespersonEmail(tradespersonID, oldEmail, token string) (bool, error) {

	claims, _, err := internal.GetRegisteredClaims(token)
	if err != nil {
		log.Printf("Failed to get registered claims from token, %v", err)
		return false, err
	}
	if claims.Subject != tradespersonID || claims.ID != oldEmail {
		return false, fmt.Errorf("tradespersonID %s doesnt match %s or email %s doesnt match %s", claims.Subject, tradespersonID, claims.ID, oldEmail)
	}
	db := database.GetConnection()

	stmt, err := db.Prepare("UPDATE tradesperson_account SET email=? WHERE tradespersonId = ?")
	if err != nil {
		return false, err
	}
	defer stmt.Close()

	results, err := stmt.Exec(oldEmail, tradespersonID)
	if err != nil {
		return false, err
	}

	rowsAffected, err := results.RowsAffected()
	if err != nil {
		return false, err
	}

	return rowsAffected == 1, nil
}
