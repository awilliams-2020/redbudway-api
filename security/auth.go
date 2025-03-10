package security

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log"
	"os"
	"redbudway-api/database"
	"strings"

	"github.com/go-openapi/errors"
	"github.com/golang-jwt/jwt/v4"
)

var jwtKey string

func init() {
	h := sha256.New()
	h.Write([]byte(os.Getenv("STRIPE_KEY")))
	jwtKey = hex.EncodeToString(h.Sum(nil))
}

func ValidateToken(bearerHeader string) (interface{}, error) {
	accessToken := strings.Split(bearerHeader, " ")[1]
	var claims jwt.RegisteredClaims
	token, err := jwt.ParseWithClaims(accessToken, &claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(jwtKey), nil
	})
	if err != nil {
		return nil, errors.New(200, err.Error())
	}
	if !token.Valid {
		return token.Valid, nil
	}
	valid := false
	if (claims.Audience[0] == "tradesperson" || claims.Audience[0] == "admin") && claims.ID == "refresh" {
		valid, err := database.CheckTradespersonRefreshToken(claims.Subject, accessToken)
		if err != nil {
			log.Printf("Failed to check tradesperson refresh token, %s", err)
			return valid, err
		}
	} else if claims.Audience[0] == "tradesperson" && claims.ID == "access" {
		valid, err := database.CheckTradespersonAccessToken(claims.Subject, accessToken)
		if err != nil {
			log.Printf("Failed to check tradesperson access token, %s", err)
			return valid, err
		}
	} else if claims.Audience[0] == "customer" && claims.ID == "refresh" {
		valid, err = database.CheckCustomerRefreshToken(claims.Subject, accessToken)
		if err != nil {
			log.Printf("Failed to check customer refresh token, %s", err)
			return valid, err
		}
	} else if claims.Audience[0] == "customer" && claims.ID == "access" {
		valid, err = database.CheckCustomerAccessToken(claims.Subject, accessToken)
		if err != nil {
			log.Printf("Failed to check customer access token, %s", err)
			return valid, err
		}
	}
	return valid, nil
}
