package internal

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestVerifyGoogleIDTokenWithAud_OK(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("id_token") != "tok" {
			t.Errorf("expected id_token query")
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"aud":"my-aud","email":"a@example.com","email_verified":"true"}`))
	}))
	defer ts.Close()
	old := googleTokeninfoURL
	googleTokeninfoURL = ts.URL
	t.Cleanup(func() { googleTokeninfoURL = old })

	email, err := verifyGoogleIDTokenWithAud("tok", "my-aud")
	if err != nil {
		t.Fatal(err)
	}
	if email != "a@example.com" {
		t.Fatalf("email %q", email)
	}
}

func TestVerifyGoogleIDTokenWithAud_EmailVerifiedBool(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"aud":"my-aud","email":"b@example.com","email_verified":true}`))
	}))
	defer ts.Close()
	old := googleTokeninfoURL
	googleTokeninfoURL = ts.URL
	t.Cleanup(func() { googleTokeninfoURL = old })

	email, err := verifyGoogleIDTokenWithAud("tok", "my-aud")
	if err != nil {
		t.Fatal(err)
	}
	if email != "b@example.com" {
		t.Fatalf("email %q", email)
	}
}

func TestVerifyGoogleIDTokenWithAud_WrongAud(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`{"aud":"other","email":"a@example.com","email_verified":"true"}`))
	}))
	defer ts.Close()
	old := googleTokeninfoURL
	googleTokeninfoURL = ts.URL
	t.Cleanup(func() { googleTokeninfoURL = old })

	_, err := verifyGoogleIDTokenWithAud("tok", "my-aud")
	if err == nil {
		t.Fatal("expected error")
	}
}
