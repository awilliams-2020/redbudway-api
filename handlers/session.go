package handlers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"redbudway-api/database"
	"redbudway-api/email"
	"redbudway-api/internal"
	"redbudway-api/restapi/operations"
	"strings"
	"time"

	"github.com/go-openapi/runtime/middleware"
	"github.com/stripe/stripe-go/v72/customer"
)

const (
	ACCESS_TIME  = 2
	REFRESH_TIME = 2
)

func ValidateCustomerRefreshToken(customerID, bearerHeader string) (bool, error) {
	accessToken := strings.Split(bearerHeader, " ")[1]
	valid := false
	claims, valid, err := internal.GetRegisteredClaims(bearerHeader)
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

func ValidateCustomerAccessToken(customerID, bearerHeader string) (bool, error) {
	accessToken := strings.Split(bearerHeader, " ")[1]
	valid := false
	claims, valid, err := internal.GetRegisteredClaims(bearerHeader)
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

func ValidateTradespersonRefreshToken(tradespersonID, bearerHeader string) (bool, error) {
	accessToken := strings.Split(bearerHeader, " ")[1]
	valid := false
	claims, valid, err := internal.GetRegisteredClaims(bearerHeader)
	if err != nil {
		log.Printf("Failed to get registered claims from token, %v", err)
		return valid, err
	}
	if (claims.Audience[0] == "tradesperson" || claims.Audience[0] == "admin") && claims.Subject == tradespersonID && claims.ID == "refresh" {
		valid, err = database.CheckTradespersonRefreshToken(tradespersonID, accessToken)
		if err != nil {
			log.Printf("Failed to check tradesperson refresh token, %s", err)
			return valid, err
		}
	}
	return valid, nil
}

func ValidateTradespersonAccessToken(tradespersonID, bearerHeader string) (bool, error) {
	accessToken := strings.Split(bearerHeader, " ")[1]
	valid := false
	claims, valid, err := internal.GetRegisteredClaims(bearerHeader)
	if err != nil {
		log.Printf("Failed to get registered claims from token, %v", err)
		return valid, err
	}

	if claims.Audience[0] == "admin" {
		valid, err = database.CheckTradespersonAccessToken(claims.Subject, accessToken)
		if err != nil {
			log.Printf("Failed to check tradesperson access token, %s", err)
			return valid, err
		}
	} else if claims.Audience[0] == "tradesperson" && claims.Subject == tradespersonID {
		valid, err = database.CheckTradespersonAccessToken(tradespersonID, accessToken)
		if err != nil {
			log.Printf("Failed to check tradesperson access token, %s", err)
			return valid, err
		}
	}
	return valid, nil
}

func ValidateAdminAccessToken(adminID, bearerHeader string) (bool, error) {
	accessToken := strings.Split(bearerHeader, " ")[1]
	valid := false
	claims, valid, err := internal.GetRegisteredClaims(bearerHeader)
	if err != nil {
		log.Printf("Failed to get registered claims from token, %v", err)
		return valid, err
	}
	if claims.Audience[0] == "admin" && claims.Subject == adminID && claims.ID == "access" {
		valid, err = database.CheckTradespersonAccessToken(adminID, accessToken)
		if err != nil {
			log.Printf("Failed to check tradesperson access token, %s", err)
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

func PostTradespersonTradespersonIDAccessTokenHandler(params operations.PostTradespersonTradespersonIDAccessTokenParams, principal interface{}) middleware.Responder {
	bearerHeader := params.HTTPRequest.Header.Get("Authorization")
	tradespersonID := params.TradespersonID

	payload := operations.PostTradespersonTradespersonIDAccessTokenOKBody{}
	response := operations.NewPostTradespersonTradespersonIDAccessTokenOK()

	valid, err := ValidateTradespersonRefreshToken(tradespersonID, bearerHeader)
	if err != nil || !valid {
		if !valid {
			log.Printf("Invalid tradesperson (%s) refresh token (%s)\n", tradespersonID, bearerHeader)
			cleared, err := database.ClearTradespersonTokens(tradespersonID)
			if !cleared || err != nil {
				log.Printf("Failed to clear tradesperson tokens, %v\n", err)
				return response
			}
		} else if err != nil {
			log.Printf("Error validating tradesperson refresh token, %s", err)
		}
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

func PostTradespersonTradespersonIDGoogleTokenHandler(params operations.PostTradespersonTradespersonIDGoogleTokenParams, principal interface{}) middleware.Responder {
	tradespersonID := params.TradespersonID
	code := params.Google.Code

	payload := operations.PostTradespersonTradespersonIDGoogleTokenOKBody{}
	response := operations.NewPostTradespersonTradespersonIDGoogleTokenOK().WithPayload(&payload)
	data := url.Values{
		"client_id":     {os.Getenv("CLIENT_ID")},
		"client_secret": {os.Getenv("CLIENT_SECRET")},
		"code":          {code},
		"grant_type":    {"authorization_code"},
		"redirect_uri":  {"https://" + os.Getenv("SUBDOMAIN") + "redbudway.com"},
	}

	resp, err := http.PostForm("https://oauth2.googleapis.com/token", data)
	if err != nil {
		log.Printf("Failed to get user google token, %v", err)
		return response
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Failed to read body in %v", err)
	}
	var res map[string]interface{}

	err = json.Unmarshal(body, &res)
	if err != nil {
		log.Printf("Failed to unmarshal response, %v", err)
		return response
	}

	accessToken := res["access_token"].(string)
	expiresIn := res["expires_in"].(float64)
	refreshToken := res["refresh_token"].(string)
	idToken := res["id_token"].(string)
	userInfo, err := internal.DecodeJWT(idToken)
	if err != nil {
		log.Printf("Failed to get user google info, %v", err)
		return response
	}

	updated, err := database.SaveTradespersonGoogleTokens(tradespersonID, refreshToken, accessToken)
	if err != nil {
		log.Printf("Failed to save tradesperson google tokens, %v", err)
	}
	if updated {
		payload.AccessToken = accessToken
		payload.ExpiresIn = expiresIn
		payload.Email = userInfo["email"].(string)
		payload.Picture = userInfo["picture"].(string)
		response.SetPayload(&payload)
	}

	return response
}

func GetTradespersonTradespersonIDGoogleTokenHandler(tradespersonID, accessToken string) (map[string]interface{}, error) {
	var res map[string]interface{}

	refreshToken, err := database.GetGoogleRefreshToken(tradespersonID, accessToken)
	if err != nil {
		log.Printf("Failed to get tradesperson %s google refresh token from access token %s, %v", tradespersonID, accessToken, err)
		return res, err
	}
	data := url.Values{
		"client_id":     {os.Getenv("CLIENT_ID")},
		"client_secret": {os.Getenv("CLIENT_SECRET")},
		"refresh_token": {refreshToken},
		"grant_type":    {"refresh_token"},
		"access_type":   {"offline"},
	}

	resp, err := http.PostForm("https://oauth2.googleapis.com/token", data)
	if err != nil {
		log.Printf("Failed to refresh user google token, %v", err)
		return res, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Failed to read body in %v", err)
	}
	err = json.Unmarshal(body, &res)
	if err != nil {
		log.Printf("Failed to unmarshal response, %v", err)
		return res, err
	}
	if res["error"] != nil {
		return res, errors.New(string(body))
	}

	return res, nil
}

func PutTradespersonTradespersonIDGoogleTokenHandler(params operations.PutTradespersonTradespersonIDGoogleTokenParams, principal interface{}) middleware.Responder {
	tradespersonID := params.TradespersonID
	accessToken := *params.Google.AccessToken

	payload := operations.PutTradespersonTradespersonIDGoogleTokenOKBody{}
	response := operations.NewPutTradespersonTradespersonIDGoogleTokenOK().WithPayload(&payload)

	res, err := GetTradespersonTradespersonIDGoogleTokenHandler(tradespersonID, accessToken)
	if err != nil {
		log.Printf("Failed to get tradesperson google token, %v", err)
		return response
	}

	accessToken = res["access_token"].(string)
	updated, err := database.UpdateTradespersonGoogleAccessToken(tradespersonID, accessToken)
	if err != nil || !updated {
		log.Printf("Failed to update tradesperson google access token, %s", err)
		return response
	}
	expiresIn := res["expires_in"].(float64)

	payload.Email, payload.Picture, err = GetTradespersonGoogleInfo(tradespersonID, accessToken)
	if err != nil {
		log.Printf("Failed to get tradesperson google info , %v", err)
		return response
	}

	payload.AccessToken = accessToken
	payload.ExpiresIn = expiresIn
	response.SetPayload(&payload)

	return response
}

func DeleteTradespersonTradespersonIDGoogleTokenHandler(params operations.DeleteTradespersonTradespersonIDGoogleTokenParams, principal interface{}) middleware.Responder {
	tradespersonID := params.TradespersonID
	accessToken := params.AccessToken

	payload := operations.DeleteTradespersonTradespersonIDGoogleTokenOKBody{Revoked: false}
	response := operations.NewDeleteTradespersonTradespersonIDGoogleTokenOK().WithPayload(&payload)

	cleared, err := database.ClearGoogleTokens(tradespersonID, accessToken)
	if err != nil {
		log.Printf("Failed to delete tradesperson google token, %v", err)
		return response
	}
	if cleared {
		resp, err := http.PostForm("https://oauth2.googleapis.com/revoke?token="+accessToken+"", nil)
		if err != nil {
			log.Printf("Failed to delete user google token, %v", err)
			return response
		}

		if resp.StatusCode == http.StatusOK {
			payload.Revoked = true
			response.SetPayload(&payload)
		}
	}

	return response
}

func GetAdminAdminIDAccessTokenHandler(params operations.GetAdminAdminIDAccessTokenParams, principal interface{}) middleware.Responder {
	adminID := params.AdminID
	bearerHeader := params.HTTPRequest.Header.Get("Authorization")

	payload := operations.GetAdminAdminIDAccessTokenOKBody{Valid: false}
	response := operations.NewGetAdminAdminIDAccessTokenOK().WithPayload(&payload)

	var err error
	payload.Valid, err = ValidateAdminAccessToken(adminID, bearerHeader)
	if err != nil {
		log.Printf("Failed to validate admin (%s) access token, %v\n", adminID, err)
	}

	response.SetPayload(&payload)

	return response
}

func GetTradespersonTradespersonIDAccessTokenHandler(params operations.GetTradespersonTradespersonIDAccessTokenParams, principal interface{}) middleware.Responder {
	tradespersonID := params.TradespersonID
	bearerHeader := params.HTTPRequest.Header.Get("Authorization")

	payload := operations.GetTradespersonTradespersonIDAccessTokenOKBody{Valid: false}
	response := operations.NewGetTradespersonTradespersonIDAccessTokenOK().WithPayload(&payload)

	var err error
	payload.Valid, err = ValidateTradespersonAccessToken(tradespersonID, bearerHeader)
	if err != nil {
		log.Printf("Failed to validate tradesperson (%s) access token, %v\n", tradespersonID, err)
	}

	response.SetPayload(&payload)

	return response
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

func GetTradespersonGoogleInfo(tradespersonID, accessToken string) (string, string, error) {
	var res map[string]interface{}

	var email, picture string
	client := &http.Client{}
	req, err := http.NewRequest(http.MethodGet, "https://www.googleapis.com/oauth2/v1/userinfo", nil)
	if err != nil {
		log.Printf("Failed to create new request, %v", err)
		return email, picture, err
	}
	req.Header.Add("content-type", "application/json")
	req.Header.Add("Authorization", "Bearer "+accessToken)

	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Failed to get tradesperson google info, %v", err)
		return email, picture, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Failed to read body in %v", err)
	}
	err = json.Unmarshal(body, &res)
	if err != nil {
		log.Printf("Failed to unmarshal response, %v", err)
		return email, picture, err
	}
	email = res["email"].(string)
	picture = res["picture"].(string)

	return email, picture, nil
}

func PostTradespersonLoginHandler(params operations.PostTradespersonLoginParams) middleware.Responder {
	tradesperson := params.Tradesperson

	payload := operations.PostTradespersonLoginOKBody{Valid: false}
	response := operations.NewPostTradespersonLoginOK().WithPayload(&payload)

	db := database.GetConnection()
	stmt, err := db.Prepare("SELECT tradespersonId, password, admin FROM tradesperson_account WHERE email=?")
	if err != nil {
		log.Printf("Failed to create select statement %s", err)
		return response
	}
	defer stmt.Close()

	var tradespersonID string
	var hashPassword string
	var admin bool
	row := stmt.QueryRow(*tradesperson.Email)
	switch err = row.Scan(&tradespersonID, &hashPassword, &admin); err {
	case sql.ErrNoRows:
		log.Println("Tradesperson doesn't exist")
	case nil:
		if internal.CheckPasswordHash(*tradesperson.Password, hashPassword) {
			accountType := "tradesperson"
			if admin {
				accountType = "admin"
			}
			accessToken, err := internal.GenerateToken(tradespersonID, accountType, "access", time.Minute*ACCESS_TIME)
			if err != nil {
				log.Printf("Failed to generate access token, %s", err)
				return response
			}
			refreshToken, err := internal.GenerateToken(tradespersonID, accountType, "refresh", time.Minute*REFRESH_TIME)
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
			payload.Admin = admin
			response.SetPayload(&payload)

			gAccessToken, err := database.GetTradespersonGoogleAccessToken(tradespersonID)
			if err != nil {
				log.Printf("Failed to get tradesperson google access token, %v", err)
				return response
			}

			if gAccessToken == "" {
				return response
			}

			res, err := GetTradespersonTradespersonIDGoogleTokenHandler(tradespersonID, gAccessToken)
			if err != nil {
				log.Printf("Failed to get tradesperson google token , %v", err)
				return response
			}
			gAccessToken = res["access_token"].(string)
			expiresIn := res["expires_in"].(float64)
			payload.GoogleAccessToken = gAccessToken
			payload.ExpiresIn = expiresIn

			updated, err := database.UpdateTradespersonGoogleAccessToken(tradespersonID, gAccessToken)
			if err != nil || !updated {
				log.Printf("Failed to update tradesperson google access token, %s", err)
				return response
			}

			payload.Email, payload.Picture, err = GetTradespersonGoogleInfo(tradespersonID, gAccessToken)
			if err != nil {
				log.Printf("Failed to get tradesperson google info , %v", err)
				return response
			}
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

func PostCustomerCustomerIDVerifyHandler(params operations.PostCustomerCustomerIDVerifyParams, principal interface{}) middleware.Responder {
	customerID := params.CustomerID
	bearerHeader := params.HTTPRequest.Header.Get("Authorization")

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
			claims, valid, err := internal.GetRegisteredClaims(bearerHeader)
			if err != nil {
				log.Printf("Failed to get registred claims from token, %s", err)
				return response
			}
			if claims.Audience[0] == "customer" && valid {
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
	default:
		log.Printf("Unknown %v", err)
	}

	return response
}

func GetCustomerCustomerIDVerifyHandler(params operations.GetCustomerCustomerIDVerifyParams, principal interface{}) middleware.Responder {
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

	var stripeID, customerEmail string
	row := stmt.QueryRow(customerID)
	switch err = row.Scan(&stripeID, &customerEmail); err {
	case sql.ErrNoRows:
		log.Printf("Customer with ID %s does not exist", customerID)
	case nil:
		token, err := internal.GenerateToken(customerID, "customer", "verification", time.Minute*15)
		if err != nil {
			log.Printf("Failed to generate JWT, %s", err)
			return response
		}
		stripeCustomer, err := customer.Get(stripeID, nil)
		if err != nil {
			log.Printf("Failed to get stripe customer, %v", err)
			return response
		}
		err = email.SendCustomerVerification(stripeCustomer.Name, customerEmail, customerID, token)
		if err != nil {
			log.Printf("Faield to send email verification, %s", err)
		}
	default:
		log.Printf("Unknown %v", err)
	}

	return response
}

func GetForgotPasswordHandler(params operations.GetForgotPasswordParams) middleware.Responder {
	accountEmail := params.Email
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
	row := stmt.QueryRow(accountEmail)
	switch err = row.Scan(&stripeID, &userID); err {
	case sql.ErrNoRows:
		log.Printf("Account with email %s does not exist", accountEmail)
	case nil:
		token, err := internal.GenerateToken(userID, accountType, "password", time.Minute*15)
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
			if err := email.ForgotPassword(accountEmail, stripeCustomer.Name, token, accountType); err != nil {
				log.Printf("Failed to send customer email, %v", err)
				return response
			}
		} else {
			tradesperson, err := database.GetTradespersonProfile(userID)
			if err != nil {
				log.Printf("Failed to get tradesperson profile %s", err)
				return response
			}
			if err := email.ForgotPassword(accountEmail, tradesperson.Name, token, accountType); err != nil {
				log.Printf("Failed to send tradesperson email, %v", err)
				return response
			}
		}
	default:
		log.Printf("Unknown %v", err)
	}

	return response
}

func PostResetPasswordHandler(params operations.PostResetPasswordParams, principal interface{}) middleware.Responder {
	password := *params.User.Password
	bearerHeader := params.HTTPRequest.Header.Get("Authorization")

	payload := operations.PostResetPasswordOKBody{Updated: false}
	response := operations.NewPostResetPasswordOK().WithPayload(&payload)

	db := database.GetConnection()

	claims, _, err := internal.GetRegisteredClaims(bearerHeader)
	if err != nil {
		log.Printf("Failed to get registered claims from token, %v", err)
		return response
	}

	var sqlStmt string
	if claims.Subject == "customer" {
		sqlStmt = "SELECT stripeId, email FROM customer_account WHERE customerId=?"
	} else {
		sqlStmt = "SELECT stripeId, email FROM tradesperson_account WHERE tradespersonId=?"
	}

	stmt, err := db.Prepare(sqlStmt)
	if err != nil {
		log.Printf("Failed to create select statement %s", err)
		return response
	}
	defer stmt.Close()

	var stripeID, accountEmail string
	row := stmt.QueryRow(claims.Subject)
	switch err = row.Scan(&stripeID, &accountEmail); err {
	case sql.ErrNoRows:
		log.Printf("Account with userID %s does not exist", claims.Subject)
	case nil:
		if claims.Audience[0] == "customer" {
			sqlStmt = "UPDATE customer_account set password=? WHERE customerId=?"
		} else {
			sqlStmt = "UPDATE tradesperson_account set password=? WHERE tradespersonId=?"
		}

		stmt, err := db.Prepare(sqlStmt)
		if err != nil {
			return response
		}
		defer stmt.Close()

		passwordHash, err := internal.HashPassword(password)
		if err != nil {
			return response
		}

		results, err := stmt.Exec(passwordHash, claims.Subject)
		if err != nil {
			return response
		}

		rowsAffected, err := results.RowsAffected()
		if err != nil {
			return response
		}

		payload.Updated = rowsAffected == 1
		response.SetPayload(&payload)
		if payload.Updated {
			if claims.Audience[0] == "customer" {
				stripeCustomer, err := customer.Get(stripeID, nil)
				if err != nil {
					log.Printf("Failed to get stripe customer, %v", err)
					return response
				}
				if err := email.PasswordUpdated(accountEmail, stripeCustomer.Name); err != nil {
					log.Printf("Failed to send customer email, %v", err)
					return response
				}
			} else {
				tradesperson, err := database.GetTradespersonProfile(claims.Subject)
				if err != nil {
					log.Printf("Failed to get tradesperson profile %s", err)
					return response
				}
				if err := email.PasswordUpdated(accountEmail, tradesperson.Name); err != nil {
					log.Printf("Failed to send tradesperson email, %v", err)
					return response
				}
			}
		}
	default:
		log.Printf("Unknown %v", err)
	}

	return response
}
