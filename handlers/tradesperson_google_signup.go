package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/go-openapi/strfmt"

	"redbudway-api/internal"
	"redbudway-api/restapi/operations"
)

type googleSignupRequest struct {
	GoogleIDToken string `json:"googleIdToken"`
}

// PostTradespersonGoogleSignupHTTP handles POST /v1/tradesperson/google-signup (wired outside swagger).
func PostTradespersonGoogleSignupHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		_ = json.NewEncoder(w).Encode(&operations.PostTradespersonCreatedBody{Created: false})
		return
	}
	var body googleSignupRequest
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		log.Printf("google signup decode: %v", err)
		_ = json.NewEncoder(w).Encode(&operations.PostTradespersonCreatedBody{Created: false})
		return
	}
	emailAddr, err := internal.VerifyGoogleIDToken(body.GoogleIDToken)
	if err != nil {
		log.Printf("google id token: %v", err)
		_ = json.NewEncoder(w).Encode(&operations.PostTradespersonCreatedBody{Created: false})
		return
	}
	pw, err := internal.RandomPassword()
	if err != nil {
		log.Printf("random password: %v", err)
		_ = json.NewEncoder(w).Encode(&operations.PostTradespersonCreatedBody{Created: false})
		return
	}
	emailFmt := strfmt.Email(emailAddr)
	passFmt := strfmt.Password(pw)
	tp := operations.PostTradespersonBody{
		Email:    &emailFmt,
		Password: &passFmt,
	}
	out := executeProviderSignup(tp)
	_ = json.NewEncoder(w).Encode(out)
}
