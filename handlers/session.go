package handlers

import (
	"database/sql"
	"log"
	"redbudway-api/database"
	"redbudway-api/email"
	"redbudway-api/internal"
	"redbudway-api/restapi/operations"
	"strings"
	"time"

	"github.com/go-openapi/runtime/middleware"
	"github.com/stripe/stripe-go/v72/account"
	"github.com/stripe/stripe-go/v72/customer"
)

const (
	ACCESS_TIME  = 15
	REFRESH_TIME = 15
)

func ValidateCustomerRefreshToken(customerID, bearerHeader string) (bool, error) {
	accessToken := strings.Split(bearerHeader, " ")[1]
	valid := false
	claims, err := internal.GetRegisteredClaims(bearerHeader)
	if err != nil {
		log.Printf("Failed to get registered claims from token, %v", err)
		return valid, err
	}
	if claims.Audience[0] == "customer" && claims.Subject == customerID && claims.ID == "refresh" {
		valid, err = database.CheckCustomerRefreshToken(customerID, accessToken)
		if err != nil {
			log.Printf("Failed to check customer refresh token, %s", err)
			return valid, err
		}
	}
	return valid, nil
}

func PostCustomerCustomerIDAccessTokenHandler(params operations.PostCustomerCustomerIDAccessTokenParams, principal interface{}) middleware.Responder {
	bearerHeader := params.HTTPRequest.Header.Get("Authorization")
	customerID := params.CustomerID

	payload := operations.PostCustomerCustomerIDAccessTokenOKBody{}
	response := operations.NewPostCustomerCustomerIDAccessTokenOK().WithPayload(&payload)

	valid, err := ValidateCustomerRefreshToken(customerID, bearerHeader)
	if err != nil || !valid {
		if !valid {
			cleared, err := database.ClearCustomerTokens(customerID)
			if !cleared || err != nil {
				log.Printf("Failed to clear customer tokens, %v", err)
				return response
			}
		}
		log.Printf("Failed to validate customer (%s) refresh token, %v", customerID, err)
		return response
	}

	//CHECK ACCESS TOKEN IS EXPIRED, IF NOT RESET TOKENS; BAD ACTOR

	accessToken, err := internal.GenerateToken(customerID, "customer", "access", time.Minute*ACCESS_TIME)
	if err != nil {
		log.Printf("Failed to generate access token, %s", err)
		return response
	}
	refreshToken, err := internal.GenerateToken(customerID, "customer", "refresh", time.Minute*REFRESH_TIME)
	if err != nil {
		log.Printf("Failed to generate refresh token, %s", err)
		return response
	}

	updated, err := database.UpdateCustomerTokens(customerID, refreshToken, accessToken)
	if err != nil {
		log.Printf("Failed to update customer tokens, %v", err)
		return response
	}

	if updated {
		payload.RefreshToken = refreshToken
		payload.AccessToken = accessToken
		response.SetPayload(&payload)
	}

	return response
}

func ValidateTradespersonRefreshToken(tradespersonID, bearerHeader string) (bool, error) {
	accessToken := strings.Split(bearerHeader, " ")[1]
	valid := false
	claims, err := internal.GetRegisteredClaims(bearerHeader)
	if err != nil {
		log.Printf("Failed to get registered claims from token, %v", err)
		return valid, err
	}
	if claims.Audience[0] == "tradesperson" && claims.Subject == tradespersonID && claims.ID == "refresh" {
		valid, err = database.CheckTradespersonRefreshToken(tradespersonID, accessToken)
		if err != nil {
			log.Printf("Failed to check tradesperson refresh token, %s", err)
			return valid, err
		}
	}
	return valid, nil
}

func PostTradespersonTradespersonIDAccessTokenHandler(params operations.PostTradespersonTradespersonIDAccessTokenParams, principal interface{}) middleware.Responder {
	bearerHeader := params.HTTPRequest.Header.Get("Authorization")
	tradespersonID := params.TradespersonID

	payload := operations.PostTradespersonTradespersonIDAccessTokenOKBody{}
	response := operations.NewPostTradespersonTradespersonIDAccessTokenOK()

	valid, err := ValidateTradespersonRefreshToken(tradespersonID, bearerHeader)
	if err != nil || !valid {
		if !valid {
			cleared, err := database.ClearTradespersonTokens(tradespersonID)
			if !cleared || err != nil {
				log.Printf("Failed to clear tradesperson tokens, %v", err)
				return response
			}
		}
		log.Printf("Failed to validate tradesperson (%s) refresh token, %v", &tradespersonID, err)
		return response
	}

	//CHECK ACCESS TOKEN IS EXPIRED?

	accessToken, err := internal.GenerateToken(tradespersonID, "tradesperson", "access", time.Minute*ACCESS_TIME)
	if err != nil {
		log.Printf("Failed to generate access token, %s", err)
		return response
	}
	refreshToken, err := internal.GenerateToken(tradespersonID, "tradesperson", "refresh", time.Minute*REFRESH_TIME)
	if err != nil {
		log.Printf("Failed to generate refresh token, %s", err)
		return response
	}

	updated, err := database.UpdateTradespersonTokens(tradespersonID, refreshToken, accessToken)
	if err != nil {
		log.Printf("Failed to save tradesperson refresh token, %v", err)
	}

	if updated {
		payload.AccessToken = accessToken
		payload.RefreshToken = refreshToken
		response.SetPayload(&payload)
	}

	return response
}

func ValidateTradespersonAccessToken(tradespersonID, bearerHeader string) (bool, error) {
	accessToken := strings.Split(bearerHeader, " ")[1]
	valid := false
	claims, err := internal.GetRegisteredClaims(bearerHeader)
	if err != nil {
		log.Printf("Failed to get registered claims from token, %v", err)
		return valid, err
	}
	if claims.Audience[0] == "tradesperson" && claims.Subject == tradespersonID && claims.ID == "access" {
		valid, err = database.CheckTradespersonAccessToken(tradespersonID, accessToken)
		if err != nil {
			log.Printf("Failed to check tradesperson access token, %s", err)
			return valid, err
		}
	}
	return valid, nil
}

func GetTradespersonTradespersonIDAccessTokenHandler(params operations.GetTradespersonTradespersonIDAccessTokenParams, principal interface{}) middleware.Responder {
	tradespersonID := params.TradespersonID
	bearerHeader := params.HTTPRequest.Header.Get("Authorization")

	payload := operations.GetTradespersonTradespersonIDAccessTokenOKBody{Valid: false}
	response := operations.NewGetTradespersonTradespersonIDAccessTokenOK().WithPayload(&payload)

	var err error
	payload.Valid, err = ValidateTradespersonAccessToken(tradespersonID, bearerHeader)
	if err != nil {
		log.Printf("Failed to validate tradesperson (%s) access token, %v", &tradespersonID, err)
	}

	response.SetPayload(&payload)

	return response
}

func ValidateCustomerAccessToken(customerID, bearerHeader string) (bool, error) {
	accessToken := strings.Split(bearerHeader, " ")[1]
	valid := false
	claims, err := internal.GetRegisteredClaims(bearerHeader)
	if err != nil {
		log.Printf("Failed to get registered claims from token, %v", err)
		return valid, err
	}
	if claims.Audience[0] == "customer" && claims.Subject == customerID && claims.ID == "access" {
		valid, err = database.CheckCustomerAccessToken(customerID, accessToken)
		if err != nil {
			log.Printf("Failed to check customer access token, %s", err)
			return valid, err
		}
	}
	return valid, nil
}

func GetCustomerCustomerIDAccessTokenHandler(params operations.GetCustomerCustomerIDAccessTokenParams, principal interface{}) middleware.Responder {
	customerID := params.CustomerID
	bearerHeader := params.HTTPRequest.Header.Get("Authorization")

	payload := operations.GetCustomerCustomerIDAccessTokenOKBody{Valid: false}
	response := operations.NewGetCustomerCustomerIDAccessTokenOK().WithPayload(&payload)

	var err error
	payload.Valid, err = ValidateCustomerAccessToken(customerID, bearerHeader)
	if err != nil {
		log.Printf("Failed to validate customer (%s) access token, %v", customerID, err)
	}

	response.SetPayload(&payload)
	return response
}

func PostTradespersonLoginHandler(params operations.PostTradespersonLoginParams) middleware.Responder {
	tradesperson := params.Tradesperson

	payload := operations.PostTradespersonLoginOKBody{Valid: false}
	response := operations.NewPostTradespersonLoginOK().WithPayload(&payload)

	db := database.GetConnection()
	stmt, err := db.Prepare("SELECT tradespersonId, password FROM tradesperson_account WHERE email=?")
	if err != nil {
		log.Printf("Failed to create select statement %s", err)
		return response
	}
	defer stmt.Close()

	var tradespersonID string
	var hashPassword string
	row := stmt.QueryRow(*tradesperson.Email)
	switch err = row.Scan(&tradespersonID, &hashPassword); err {
	case sql.ErrNoRows:
		log.Println("Tradesperson doesn't exist")
	case nil:
		if internal.CheckPasswordHash(*tradesperson.Password, hashPassword) {
			accessToken, err := internal.GenerateToken(tradespersonID, "tradesperson", "access", time.Minute*ACCESS_TIME)
			if err != nil {
				log.Printf("Failed to generate access token, %s", err)
				return response
			}
			refreshToken, err := internal.GenerateToken(tradespersonID, "tradesperson", "refresh", time.Minute*REFRESH_TIME)
			if err != nil {
				log.Printf("Failed to generate refresh token, %s", err)
				return response
			}

			isSaved, err := database.SaveTradespersonTokens(tradespersonID, refreshToken, accessToken)
			if err != nil {
				log.Printf("Failed to save tradesperson tokens, %v", err)
			}

			payload.Valid = isSaved
			payload.TradespersonID = tradespersonID
			payload.RefreshToken = refreshToken
			payload.AccessToken = accessToken
			response.SetPayload(&payload)
		}
	default:
		log.Printf("Unkown default case %s", err)
	}

	return response
}

func PostTradespersonTradespersonIDLogoutHandler(params operations.PostTradespersonTradespersonIDLogoutParams, principal interface{}) middleware.Responder {
	tradespersonID := params.TradespersonID

	payload := operations.PostTradespersonTradespersonIDLogoutOKBody{LoggedOut: false}
	response := operations.NewPostTradespersonTradespersonIDLogoutOK().WithPayload(&payload)

	db := database.GetConnection()
	stmt, err := db.Prepare("SELECT tradespersonId FROM tradesperson_account WHERE tradespersonId=?")
	if err != nil {
		log.Printf("Failed to create select statement %s", err)
		return response
	}
	defer stmt.Close()

	row := stmt.QueryRow(tradespersonID)
	switch err = row.Scan(&tradespersonID); err {
	case sql.ErrNoRows:
		log.Println("Tradesperson doesn't exist")
	case nil:
		cleared, err := database.ClearTradespersonTokens(tradespersonID)
		if err != nil {
			log.Printf("Failed to log out tradesperson %s, %v", tradespersonID, err)
		}
		payload.LoggedOut = cleared
		response.SetPayload(&payload)
	default:
		log.Printf("Failed scanning for tradesperson and password %s", err)
	}

	return response
}

func PostCustomerLoginHandler(params operations.PostCustomerLoginParams) middleware.Responder {
	customer := params.Customer

	payload := operations.PostCustomerLoginOKBody{Valid: false}
	response := operations.NewPostCustomerLoginOK().WithPayload(&payload)

	db := database.GetConnection()
	stmt, err := db.Prepare("SELECT customerId, password FROM customer_account WHERE email=?")
	if err != nil {
		log.Printf("Failed to create select statement %s", err)
		return response
	}
	defer stmt.Close()
	var customerID string
	var hashPassword string
	row := stmt.QueryRow(*customer.Email)
	switch err = row.Scan(&customerID, &hashPassword); err {
	case sql.ErrNoRows:
		log.Println("Customer doesn't exist")
	case nil:
		if internal.CheckPasswordHash(*customer.Password, hashPassword) {
			accessToken, err := internal.GenerateToken(customerID, "customer", "access", time.Minute*ACCESS_TIME)
			if err != nil {
				log.Printf("Failed to generate access token, %s", err)
				return response
			}
			refreshToken, err := internal.GenerateToken(customerID, "customer", "refresh", time.Minute*REFRESH_TIME)
			if err != nil {
				log.Printf("Failed to generate refresh token, %s", err)
				return response
			}

			isSaved, err := database.SaveCustomerTokens(customerID, refreshToken, accessToken)
			if err != nil {
				log.Printf("Failed to save customer tokens, %v", err)
				return response
			}

			payload.Valid = isSaved
			payload.CustomerID = customerID
			payload.RefreshToken = refreshToken
			payload.AccessToken = accessToken
			response.SetPayload(&payload)
		}
	default:
		log.Printf("%v", err)
	}

	return response
}

func PostCustomerCustomerIDLogoutHandler(params operations.PostCustomerCustomerIDLogoutParams, principal interface{}) middleware.Responder {
	customerID := params.CustomerID

	payload := operations.PostCustomerCustomerIDLogoutOKBody{LoggedOut: false}
	response := operations.NewPostCustomerCustomerIDLogoutOK().WithPayload(&payload)

	db := database.GetConnection()
	stmt, err := db.Prepare("SELECT customerId FROM customer_account WHERE customerId=?")
	if err != nil {
		log.Printf("Failed to create select statement %s", err)
		return response
	}
	defer stmt.Close()

	row := stmt.QueryRow(customerID)
	switch err = row.Scan(&customerID); err {
	case sql.ErrNoRows:
		log.Println("Tradesperson doesn't exist")
	case nil:
		cleared, err := database.ClearCustomerTokens(customerID)
		if err != nil {
			log.Printf("Failed to log out tradesperson %s, %v", customerID, err)
		}
		payload.LoggedOut = cleared
		response.SetPayload(&payload)
	default:
		log.Printf("Failed scanning for tradesperson and password %s", err)
	}

	return response
}

func PostCustomerCustomerIDVerifyHandler(params operations.PostCustomerCustomerIDVerifyParams) middleware.Responder {
	customerID := params.CustomerID
	token := params.Body.AccessToken

	db := database.GetConnection()

	verified := false
	payload := operations.GetCustomerCustomerIDVerifyOKBody{Verified: verified}
	response := operations.NewGetCustomerCustomerIDVerifyOK().WithPayload(&payload)

	stmt, err := db.Prepare("SELECT verified FROM customer_account WHERE customerId=?")
	if err != nil {
		log.Printf("Failed to create select statement %s", err)
		return response
	}
	defer stmt.Close()

	row := stmt.QueryRow(customerID)
	switch err = row.Scan(&verified); err {
	case sql.ErrNoRows:
		log.Printf("Customer with ID %s does not exist", customerID)
	case nil:
		if verified {
			payload.Verified = verified
			response.SetPayload(&payload)
		} else {
			if token != "" {
				claims, err := internal.GetRegisteredClaims("bearer " + token)
				if err != nil {
					log.Printf("Failed to get registred claims from token, %s", err)
					return response
				}
				if claims.Audience[0] == "customer" {
					if claims.ID == "verification" {
						stmt, err := db.Prepare("UPDATE customer_account SET verified=True WHERE customerId = ?")
						if err != nil {
							return response
						}
						defer stmt.Close()

						results, err := stmt.Exec(customerID)
						if err != nil {
							return response
						}

						rowsAffected, err := results.RowsAffected()
						if err != nil {
							return response
						}

						if rowsAffected == 1 {
							verified = true
							payload.Verified = verified
							response.SetPayload(&payload)
						}
					}
				}
			}
		}
	default:
		log.Printf("Unknown %v", err)
	}

	return response
}

func GetCustomerCustomerIDVerifyHandler(params operations.GetCustomerCustomerIDVerifyParams) middleware.Responder {
	customerID := params.CustomerID

	db := database.GetConnection()

	verified := false
	payload := operations.GetCustomerCustomerIDVerifyOKBody{Verified: verified}
	response := operations.NewGetCustomerCustomerIDVerifyOK().WithPayload(&payload)

	stmt, err := db.Prepare("SELECT verified FROM customer_account WHERE customerId=?")
	if err != nil {
		log.Printf("Failed to create select statement %s", err)
		return response
	}
	defer stmt.Close()

	row := stmt.QueryRow(customerID)
	switch err = row.Scan(&verified); err {
	case sql.ErrNoRows:
		log.Printf("Customer with ID %s does not exist", customerID)
	case nil:
		payload.Verified = verified
		response.SetPayload(&payload)
	default:
		log.Printf("Unknown %v", err)
	}

	return response
}

func GetCustomerCustomerIDReverifyHandler(params operations.GetCustomerCustomerIDReverifyParams) middleware.Responder {
	customerID := params.CustomerID

	db := database.GetConnection()

	response := operations.NewGetCustomerCustomerIDReverifyOK()

	stmt, err := db.Prepare("SELECT stripeId, email FROM customer_account WHERE customerId=?")
	if err != nil {
		log.Printf("Failed to create select statement %s", err)
		return response
	}
	defer stmt.Close()

	var stripeId, customerEmail string
	row := stmt.QueryRow(customerID)
	switch err = row.Scan(&stripeId, &customerEmail); err {
	case sql.ErrNoRows:
		log.Printf("Customer with ID %s does not exist", customerID)
	case nil:
		err := email.SendCustomerVerification("", customerEmail, customerID)
		if err != nil {
			log.Printf("Faield to send email verification, %s", err)
		}
	default:
		log.Printf("Unknown %v", err)
	}

	return response

}

func GetForgotPasswordHandler(params operations.GetForgotPasswordParams) middleware.Responder {
	userEmail := params.Email
	accountType := params.AccountType

	response := operations.NewGetForgotPasswordOK()

	db := database.GetConnection()

	var sqlStmt string
	if accountType == "customer" {
		sqlStmt = "SELECT stripeId, customerId FROM customer_account WHERE email=?"
	} else {
		sqlStmt = "SELECT stripeId, tradespersonId FROM tradesperson_account WHERE email=?"
	}

	stmt, err := db.Prepare(sqlStmt)
	if err != nil {
		log.Printf("Failed to create select statement %s", err)
		return response
	}
	defer stmt.Close()

	var stripeID, userID string
	row := stmt.QueryRow(userEmail)
	switch err = row.Scan(&stripeID, &userID); err {
	case sql.ErrNoRows:
		log.Printf("Account with email %s does not exist", userEmail)
	case nil:
		token, err := internal.GenerateToken(userID, accountType, "password", time.Hour*24)
		if err != nil {
			log.Printf("Failed to generate JWT, %s", err)
			return response
		}
		if accountType == "customer" {
			stripeCustomer, err := customer.Get(stripeID, nil)
			if err != nil {
				log.Printf("Failed to get stripe customer, %v", err)
				return response
			}
			if err := email.ForgotPassword(stripeCustomer.Email, stripeCustomer.Name, token, accountType, userID); err != nil {
				log.Printf("Failed to send customer email, %v", err)
				return response
			}
		} else {
			tradesperson, err := database.GetTradespersonAccount(userID)
			if err != nil {
				log.Printf("Failed to get tradesperson account , %v", err)
				return response
			}
			stripeAccount, err := account.GetByID(stripeID, nil)
			if err != nil {
				log.Printf("Failed to get stripe customer, %v", err)
				return response
			}
			if err := email.ForgotPassword(stripeAccount.Email, tradesperson.Name, token, accountType, userID); err != nil {
				log.Printf("Failed to send tradesperson email, %v", err)
				return response
			}
		}
	default:
		log.Printf("Unknown %v", err)
	}

	return response
}

func PostResetPasswordHandler(params operations.PostResetPasswordParams) middleware.Responder {
	userID := params.User.UserID
	accountType := params.User.AccountType
	password := params.User.Password

	payload := operations.PostResetPasswordOKBody{Updated: false}
	response := operations.NewPostResetPasswordOK().WithPayload(&payload)

	db := database.GetConnection()

	var sqlStmt string
	if *accountType == "customer" {
		sqlStmt = "UPDATE customer_account set password=? WHERE customerId=?"
	} else {
		sqlStmt = "UPDATE tradesperson_account set password=? WHERE tradespersonId=?"
	}

	stmt, err := db.Prepare(sqlStmt)
	if err != nil {
		return response
	}
	defer stmt.Close()

	passwordHash, err := internal.HashPassword(*password)
	if err != nil {
		return response
	}

	results, err := stmt.Exec(passwordHash, userID)
	if err != nil {
		return response
	}

	rowsAffected, err := results.RowsAffected()
	if err != nil {
		return response
	}

	payload.Updated = rowsAffected == 1
	response.SetPayload(&payload)

	return response
}
