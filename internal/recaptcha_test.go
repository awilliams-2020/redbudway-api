package internal

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRecaptchaSecret_requiresEnv(t *testing.T) {
	t.Setenv("RECAPTCHA_SECRET", "")

	_, err := recaptchaSecret()
	if err == nil {
		t.Fatal("expected error when RECAPTCHA_SECRET is empty")
	}
}

func TestRecaptchaSecret_fromEnv(t *testing.T) {
	t.Setenv("RECAPTCHA_SECRET", "test-secret-value")
	s, err := recaptchaSecret()
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if s != "test-secret-value" {
		t.Fatalf("got %q", s)
	}
}

func TestVerifyReCaptcha_siteverifySuccess(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"success": true}`))
	}))
	t.Cleanup(srv.Close)

	old := recaptchaSiteVerifyURL
	recaptchaSiteVerifyURL = srv.URL
	t.Cleanup(func() { recaptchaSiteVerifyURL = old })

	t.Setenv("RECAPTCHA_SECRET", "unit-test-secret")

	ok, err := VerifyReCaptcha("unit-test-token")
	if err != nil {
		t.Fatalf("VerifyReCaptcha: %v", err)
	}
	if !ok {
		t.Fatal("expected success true")
	}
}

func TestVerifyReCaptcha_siteverifyFailure(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"success": false}`))
	}))
	t.Cleanup(srv.Close)

	old := recaptchaSiteVerifyURL
	recaptchaSiteVerifyURL = srv.URL
	t.Cleanup(func() { recaptchaSiteVerifyURL = old })

	t.Setenv("RECAPTCHA_SECRET", "unit-test-secret")

	ok, err := VerifyReCaptcha("bad-token")
	if err != nil {
		t.Fatalf("VerifyReCaptcha: %v", err)
	}
	if ok {
		t.Fatal("expected success false")
	}
}

func TestVerifyReCaptcha_missingSecret(t *testing.T) {
	t.Setenv("RECAPTCHA_SECRET", "")

	_, err := VerifyReCaptcha("x")
	if err == nil {
		t.Fatal("expected error when RECAPTCHA_SECRET unset")
	}
}
