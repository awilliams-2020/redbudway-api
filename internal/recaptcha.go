package internal

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

// recaptchaSiteVerifyURL is the Google siteverify endpoint; tests may replace it with httptest.
var recaptchaSiteVerifyURL = "https://www.google.com/recaptcha/api/siteverify"

func recaptchaSecret() (string, error) {
	if s := os.Getenv("RECAPTCHA_SECRET"); s != "" {
		return s, nil
	}
	return "", fmt.Errorf("RECAPTCHA_SECRET is not set (must pair with reCaptchaKey in Angular environments)")
}

func VerifyReCaptcha(token string) (bool, error) {
	valid := false
	secret, err := recaptchaSecret()
	if err != nil {
		log.Printf("%v", err)
		return valid, err
	}
	URL := fmt.Sprintf("%s?secret=%s&response=%s", recaptchaSiteVerifyURL, secret, token)
	resp, err := http.Post(URL, "application/json", bytes.NewBuffer([]byte{}))
	if err != nil {
		log.Printf("Failed to recaptcha token, %v", err)
		return valid, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Failed to read body in %v", err)
	}
	var r map[string]interface{}
	err = json.Unmarshal(body, &r)
	if err != nil {
		log.Printf("Failed to unmarshal response, %v", err)
		return valid, err
	}
	valid = r["success"].(bool)

	return valid, nil
}
