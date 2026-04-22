package email

import (
	"os"
	"testing"
)

func TestResendFromAddress_defaults(t *testing.T) {
	_ = os.Unsetenv("RESEND_FROM")
	if got := resendFromAddress(); got != "Redbud Way <onboarding@resend.dev>" {
		t.Fatalf("got %q", got)
	}
}

func TestResendFromAddress_resendFromEnv(t *testing.T) {
	t.Setenv("RESEND_FROM", "App <app@example.com>")
	if got := resendFromAddress(); got != "App <app@example.com>" {
		t.Fatalf("got %q", got)
	}
}
