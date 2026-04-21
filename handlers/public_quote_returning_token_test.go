package handlers

import (
	"os"
	"testing"
)

func TestQuoteReturningSecret_empty(t *testing.T) {
	t.Setenv(quoteRequestReturningSecretEnv, "")
	if quoteReturningSecret() != "" {
		t.Fatal("expected empty")
	}
}

func TestSignAndValidQuoteReturningToken(t *testing.T) {
	const secret = "unit-test-secret-at-least-16-bytes"
	t.Setenv(quoteRequestReturningSecretEnv, secret)
	t.Cleanup(func() { _ = os.Unsetenv(quoteRequestReturningSecretEnv) })

	tok, err := signQuoteReturningToken("User@Example.com", "quote-abc")
	if err != nil {
		t.Fatalf("sign: %v", err)
	}
	if tok == "" {
		t.Fatal("expected token")
	}
	if !validQuoteReturningToken(tok, "user@example.com", "quote-abc") {
		t.Fatal("should validate same email case-insensitive")
	}
	if validQuoteReturningToken(tok, "other@example.com", "quote-abc") {
		t.Fatal("wrong email")
	}
	if validQuoteReturningToken(tok, "user@example.com", "other-quote") {
		t.Fatal("wrong quote id")
	}
}

func TestSignQuoteReturningToken_noSecret(t *testing.T) {
	t.Setenv(quoteRequestReturningSecretEnv, "")
	tok, err := signQuoteReturningToken("a@b.co", "q1")
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if tok != "" {
		t.Fatalf("expected empty token when secret unset, got %q", tok)
	}
}

func TestValidQuoteReturningToken_malformed(t *testing.T) {
	t.Setenv(quoteRequestReturningSecretEnv, "somesecret")
	if validQuoteReturningToken("not-a-jwt", "a@b.co", "q") {
		t.Fatal("malformed")
	}
	if validQuoteReturningToken("", "a@b.co", "q") {
		t.Fatal("empty")
	}
}
