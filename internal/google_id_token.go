package internal

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
)

// googleTokeninfoURL is the tokeninfo endpoint base (tests may point this at httptest).
var googleTokeninfoURL = "https://oauth2.googleapis.com/tokeninfo"

// VerifyGoogleIDToken validates a Google Sign-In ID token (GIS credential JWT).
// It uses Google's tokeninfo endpoint; CLIENT_ID must match the OAuth web client used in Angular (gClientId).
func VerifyGoogleIDToken(idToken string) (email string, err error) {
	aud := strings.TrimSpace(os.Getenv("CLIENT_ID"))
	if aud == "" {
		return "", fmt.Errorf("CLIENT_ID is not set (Google OAuth web client ID)")
	}
	return verifyGoogleIDTokenWithAud(idToken, aud)
}

// tokeninfoEmailVerified parses email_verified from Google's tokeninfo JSON.
// Google may return it as the string "true" or as a boolean true.
func tokeninfoEmailVerified(raw json.RawMessage) bool {
	if len(raw) == 0 {
		return false
	}
	var s string
	if err := json.Unmarshal(raw, &s); err == nil {
		return s == "true"
	}
	var b bool
	if err := json.Unmarshal(raw, &b); err == nil {
		return b
	}
	return false
}

func verifyGoogleIDTokenWithAud(idToken, expectedAud string) (string, error) {
	if idToken == "" || expectedAud == "" {
		return "", fmt.Errorf("missing token or audience")
	}
	u := googleTokeninfoURL + "?id_token=" + url.QueryEscape(idToken)
	resp, err := http.Get(u)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("tokeninfo: status %d: %s", resp.StatusCode, strings.TrimSpace(string(body)))
	}

	var m map[string]json.RawMessage
	if err := json.Unmarshal(body, &m); err != nil {
		return "", err
	}

	rawAud, ok := m["aud"]
	if !ok {
		return "", fmt.Errorf("no aud in token")
	}
	var aud string
	if err := json.Unmarshal(rawAud, &aud); err != nil {
		return "", fmt.Errorf("invalid aud in token: %w", err)
	}
	if aud != expectedAud {
		return "", fmt.Errorf("invalid token audience")
	}

	rawEmail, ok := m["email"]
	if !ok {
		return "", fmt.Errorf("no email in token")
	}
	var email string
	if err := json.Unmarshal(rawEmail, &email); err != nil || email == "" {
		return "", fmt.Errorf("no email in token")
	}

	if !tokeninfoEmailVerified(m["email_verified"]) {
		return "", fmt.Errorf("email not verified by Google")
	}

	return email, nil
}
