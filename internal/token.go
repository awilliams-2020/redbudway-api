package internal

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/go-openapi/errors"
	"github.com/golang-jwt/jwt/v4"
)

var jwtKey string

func init() {
	log.Println("Initializing jwtKey")
	h := sha256.New()
	h.Write([]byte(os.Getenv("STRIPE_KEY")))
	jwtKey = hex.EncodeToString(h.Sum(nil))
}

func GenerateToken(userId, accountType, tokenType string, expDate time.Duration) (string, error) {
	claims := &jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(expDate)),
		Subject:   userId,
		Audience:  []string{accountType},
		ID:        tokenType,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(jwtKey))
}

func ValidateToken(bearerHeader string) (interface{}, error) {
	bearerToken := strings.Split(bearerHeader, " ")[1]
	token, err := jwt.Parse(bearerToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(jwtKey), nil
	})
	if err != nil {
		return nil, errors.New(200, err.Error())
	}

	return token.Valid, err
}

func GetRegisteredClaims(bearerHeader string) (jwt.RegisteredClaims, error) {
	bearerToken := strings.Split(bearerHeader, " ")[1]
	var claims jwt.RegisteredClaims
	_, err := jwt.ParseWithClaims(bearerToken, &claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(jwtKey), nil
	})
	if err != nil {
		return claims, err
	}
	return claims, nil
}
