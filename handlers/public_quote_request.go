package handlers

import (
	"crypto/hmac"
	"crypto/sha256"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"redbudway-api/database"
	"redbudway-api/internal"
	"redbudway-api/internal/httputil"
	"redbudway-api/internal/quoteratelimit"

	"github.com/stripe/stripe-go/v82"
	"github.com/stripe/stripe-go/v82/customer"
	"github.com/stripe/stripe-go/v82/quote"
)

const quoteRequestReturningSecretEnv = "QUOTE_REQUEST_RETURNING_SECRET"

var (
	quoteReqIPLimiter    *quoteratelimit.WindowLimiter
	quoteReqEmailLimiter *quoteratelimit.WindowLimiter
)

func init() {
	// Anonymous quote-request abuse limits (US-2). Tune via env.
	quoteReqIPLimiter = quoteratelimit.NewWindowLimiter(
		quoteratelimit.EnvInt("QUOTE_REQ_PER_IP_MAX", 60),
		quoteratelimit.EnvDurationMinutes("QUOTE_REQ_PER_IP_WINDOW_MIN", 60*time.Minute),
	)
	quoteReqEmailLimiter = quoteratelimit.NewWindowLimiter(
		quoteratelimit.EnvInt("QUOTE_REQ_PER_EMAIL_MAX", 30),
		quoteratelimit.EnvDurationMinutes("QUOTE_REQ_PER_EMAIL_WINDOW_MIN", 24*60*time.Minute),
	)
}

type publicQuoteRequestBody struct {
	Email          string   `json:"email"`
	Message        string   `json:"message"`
	Name           string   `json:"name"`
	Images         []string `json:"images"`
	Token          string   `json:"token"`
	ReturningToken string   `json:"returningToken"`
}

type publicQuoteRequestOK struct {
	Requested      bool   `json:"requested"`
	StripeQuoteID  string `json:"stripeQuoteId"`
	ReturningToken string `json:"returningToken,omitempty"`
}

type publicQuoteRequestErr struct {
	Error string `json:"error"`
}

type quoteReturningPayload struct {
	Email   string `json:"email"`
	QuoteID string `json:"quoteId"`
	Exp     int64  `json:"exp"`
}

func quoteReturningSecret() string {
	return strings.TrimSpace(os.Getenv(quoteRequestReturningSecretEnv))
}

func signQuoteReturningToken(email, catalogQuoteID string) (string, error) {
	secret := quoteReturningSecret()
	if secret == "" {
		return "", nil
	}
	payload := quoteReturningPayload{
		Email:   strings.TrimSpace(strings.ToLower(email)),
		QuoteID: catalogQuoteID,
		Exp:     time.Now().Add(7 * 24 * time.Hour).Unix(),
	}
	raw, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}
	b64 := base64.RawURLEncoding.EncodeToString(raw)
	mac := hmac.New(sha256.New, []byte(secret))
	if _, err := mac.Write([]byte(b64)); err != nil {
		return "", err
	}
	sig := base64.RawURLEncoding.EncodeToString(mac.Sum(nil))
	return b64 + "." + sig, nil
}

func validQuoteReturningToken(token, email, catalogQuoteID string) bool {
	secret := quoteReturningSecret()
	if secret == "" || strings.TrimSpace(token) == "" {
		return false
	}
	parts := strings.Split(token, ".")
	if len(parts) != 2 {
		return false
	}
	b64, sigB64 := parts[0], parts[1]
	mac := hmac.New(sha256.New, []byte(secret))
	if _, err := mac.Write([]byte(b64)); err != nil {
		return false
	}
	expectedSig := mac.Sum(nil)
	sig, err := base64.RawURLEncoding.DecodeString(sigB64)
	if err != nil || !hmac.Equal(sig, expectedSig) {
		return false
	}
	raw, err := base64.RawURLEncoding.DecodeString(b64)
	if err != nil {
		return false
	}
	var p quoteReturningPayload
	if err := json.Unmarshal(raw, &p); err != nil {
		return false
	}
	if time.Now().Unix() > p.Exp {
		return false
	}
	if !strings.EqualFold(strings.TrimSpace(p.Email), strings.TrimSpace(email)) {
		return false
	}
	if strings.TrimSpace(p.QuoteID) != strings.TrimSpace(catalogQuoteID) {
		return false
	}
	return true
}

// PostPublicQuoteRequestHTTP handles POST /v1/quote/{catalogQuoteId}/request (anonymous quote request with reCAPTCHA).
func PostPublicQuoteRequestHTTP(w http.ResponseWriter, r *http.Request, catalogQuoteID string) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		_ = json.NewEncoder(w).Encode(publicQuoteRequestErr{Error: "method not allowed"})
		return
	}

	catalogQuoteID = strings.TrimSpace(catalogQuoteID)
	if catalogQuoteID == "" {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(publicQuoteRequestErr{Error: "quote id is required"})
		return
	}

	var body publicQuoteRequestBody
	if err := json.NewDecoder(http.MaxBytesReader(w, r.Body, 65536)).Decode(&body); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(publicQuoteRequestErr{Error: "invalid JSON"})
		return
	}

	body.Email = strings.TrimSpace(body.Email)
	body.Message = strings.TrimSpace(body.Message)
	body.Name = strings.TrimSpace(body.Name)
	body.Token = strings.TrimSpace(body.Token)
	body.ReturningToken = strings.TrimSpace(body.ReturningToken)

	if body.Email == "" || body.Message == "" {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(publicQuoteRequestErr{Error: "email and message are required"})
		return
	}

	ipKey := httputil.ClientIP(r)
	emailKey := strings.ToLower(body.Email)
	if !quoteReqIPLimiter.Allow(ipKey) {
		w.WriteHeader(http.StatusTooManyRequests)
		_ = json.NewEncoder(w).Encode(publicQuoteRequestErr{Error: "too many requests; try again later"})
		return
	}
	if !quoteReqEmailLimiter.Allow(emailKey) {
		w.WriteHeader(http.StatusTooManyRequests)
		_ = json.NewEncoder(w).Encode(publicQuoteRequestErr{Error: "too many requests for this email; try again later"})
		return
	}

	skipCaptcha := validQuoteReturningToken(body.ReturningToken, body.Email, catalogQuoteID)
	if !skipCaptcha {
		if body.Token == "" {
			w.WriteHeader(http.StatusBadRequest)
			_ = json.NewEncoder(w).Encode(publicQuoteRequestErr{Error: "token is required"})
			return
		}
		ok, err := internal.VerifyReCaptcha(body.Token)
		if err != nil {
			log.Printf("public quote request reCAPTCHA error: %v", err)
			w.WriteHeader(http.StatusUnauthorized)
			_ = json.NewEncoder(w).Encode(publicQuoteRequestErr{Error: "reCAPTCHA verification failed"})
			return
		}
		if !ok {
			w.WriteHeader(http.StatusUnauthorized)
			_ = json.NewEncoder(w).Encode(publicQuoteRequestErr{Error: "reCAPTCHA verification failed"})
			return
		}
	}

	tradesperson, tpStripeID, tradespersonID, err := database.GetTradespersonAccountByQuoteID(catalogQuoteID)
	if err != nil {
		if err == sql.ErrNoRows {
			w.WriteHeader(http.StatusNotFound)
			_ = json.NewEncoder(w).Encode(publicQuoteRequestErr{Error: "quote not found"})
			return
		}
		log.Printf("public quote request GetTradespersonAccountByQuoteID: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(publicQuoteRequestErr{Error: "internal error"})
		return
	}

	customerID, cuStripeID, err := database.GetCustomerAccountByEmail(body.Email)
	if err != nil && err != sql.ErrNoRows {
		log.Printf("public quote request GetCustomerAccountByEmail: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(publicQuoteRequestErr{Error: "internal error"})
		return
	}
	if err == sql.ErrNoRows {
		customerID, cuStripeID, err = database.CreateGuestCustomerAccount(body.Email, body.Name)
		if err != nil {
			log.Printf("public quote request CreateGuestCustomerAccount: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			_ = json.NewEncoder(w).Encode(publicQuoteRequestErr{Error: "internal error"})
			return
		}
	} else if body.Name != "" {
		// Best-effort: sync display name to Stripe when client supplied one.
		if _, err := customer.Update(cuStripeID, &stripe.CustomerParams{Name: stripe.String(body.Name)}); err != nil {
			log.Printf("public quote request customer.Update name: %v", err)
		}
	}

	_quote, err := database.GetTradespersonQuote(tradespersonID, catalogQuoteID)
	if err != nil {
		if err == sql.ErrNoRows {
			w.WriteHeader(http.StatusNotFound)
			_ = json.NewEncoder(w).Encode(publicQuoteRequestErr{Error: "quote not found"})
			return
		}
		log.Printf("public quote request GetTradespersonQuote: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(publicQuoteRequestErr{Error: "internal error"})
		return
	}

	connectCu, err := database.GetOrCreateStripeCustomerOnConnect(tradespersonID, customerID, tpStripeID)
	if err != nil {
		log.Printf("public quote request GetOrCreateStripeCustomerOnConnect: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(publicQuoteRequestErr{Error: "internal error"})
		return
	}

	daysDue := int64(7)
	qparams := &stripe.QuoteParams{
		Customer:         stripe.String(connectCu),
		CollectionMethod: stripe.String("send_invoice"),
		InvoiceSettings: &stripe.QuoteInvoiceSettingsParams{
			DaysUntilDue: &daysDue,
		},
	}
	qparams.SetStripeAccount(tpStripeID)
	stripeQuote, err := quote.New(qparams)
	if err != nil {
		log.Printf("public quote request quote.New: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(publicQuoteRequestErr{Error: "internal error"})
		return
	}

	if stripeQuote.Status != "draft" {
		log.Printf("public quote request: unexpected stripe quote status %s", stripeQuote.Status)
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(publicQuoteRequestErr{Error: "internal error"})
		return
	}

	created, err := database.SaveQuote(stripeQuote.ID, customerID, tradespersonID, body.Message, _quote.ID, stripeQuote.Created)
	if err != nil {
		log.Printf("public quote request SaveQuote: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(publicQuoteRequestErr{Error: "internal error"})
		return
	}
	if !created {
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(publicQuoteRequestErr{Error: "internal error"})
		return
	}

	images := body.Images
	if images == nil {
		images = []string{}
	}
	go EmailQuoteHelper(tradesperson, _quote, images, cuStripeID, body.Message, stripeQuote.ID)

	w.WriteHeader(http.StatusCreated)
	resp := publicQuoteRequestOK{Requested: true, StripeQuoteID: stripeQuote.ID}
	if rt, err := signQuoteReturningToken(body.Email, catalogQuoteID); err == nil && rt != "" {
		resp.ReturningToken = rt
	}
	_ = json.NewEncoder(w).Encode(resp)
}
