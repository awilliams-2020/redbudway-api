package handlers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	"redbudway-api/database"
	"redbudway-api/internal"
	"redbudway-api/restapi/operations"
)

// PostTradespersonGoogleLoginHTTP handles POST /v1/tradesperson/google-login (wired outside swagger).
// Verifies the Google ID token, finds a provider by email, and returns the same body as POST /tradesperson/login.
func PostTradespersonGoogleLoginHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		_ = json.NewEncoder(w).Encode(&operations.PostTradespersonLoginOKBody{Valid: false})
		return
	}
	var body struct {
		GoogleIDToken string `json:"googleIdToken"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		log.Printf("google login decode: %v", err)
		_ = json.NewEncoder(w).Encode(&operations.PostTradespersonLoginOKBody{Valid: false})
		return
	}
	emailAddr, err := internal.VerifyGoogleIDToken(body.GoogleIDToken)
	if err != nil {
		log.Printf("google id token login: %v", err)
		_ = json.NewEncoder(w).Encode(&operations.PostTradespersonLoginOKBody{Valid: false})
		return
	}

	db := database.GetConnection()
	stmt, err := db.Prepare("SELECT tradespersonId, admin FROM tradesperson_account WHERE email=?")
	if err != nil {
		log.Printf("google login prepare: %v", err)
		_ = json.NewEncoder(w).Encode(&operations.PostTradespersonLoginOKBody{Valid: false})
		return
	}
	defer stmt.Close()

	var tradespersonID string
	var admin bool
	switch err = stmt.QueryRow(emailAddr).Scan(&tradespersonID, &admin); err {
	case sql.ErrNoRows:
		log.Printf("google login: no account for %s", emailAddr)
		_ = json.NewEncoder(w).Encode(&operations.PostTradespersonLoginOKBody{Valid: false})
		return
	case nil:
		out := completeTradespersonLoginSession(tradespersonID, admin)
		_ = json.NewEncoder(w).Encode(out)
		return
	default:
		log.Printf("google login scan: %v", err)
		_ = json.NewEncoder(w).Encode(&operations.PostTradespersonLoginOKBody{Valid: false})
	}
}
