package email

import (
	"os"
	"testing"
)

func TestResendFromAddress_defaults(t *testing.T) {
	_ = os.Unsetenv("RESEND_FROM")
	_ = os.Unsetenv("SMTP_USER")
	if got := resendFromAddress(); got != "Redbud Way <onboarding@resend.dev>" {
		t.Fatalf("got %q", got)
	}
}

func TestResendFromAddress_resendFromEnv(t *testing.T) {
	t.Setenv("RESEND_FROM", "App <app@example.com>")
	t.Setenv("SMTP_USER", "ignored@example.com")
	if got := resendFromAddress(); got != "App <app@example.com>" {
		t.Fatalf("got %q", got)
	}
}

func TestResendFromAddress_smtpFallback(t *testing.T) {
	t.Setenv("RESEND_FROM", "")
	t.Setenv("SMTP_USER", "svc@mail.example.com")
	if got := resendFromAddress(); got != "Redbud Way <svc@mail.example.com>" {
		t.Fatalf("got %q", got)
	}
}
