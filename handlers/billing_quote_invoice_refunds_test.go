package handlers

import (
	"testing"

	"github.com/stripe/stripe-go/v82"
)

func TestGuestInvoicePaymentSegment(t *testing.T) {
	t.Parallel()
	inv := &stripe.Invoice{AmountPaid: 0}
	if g := guestInvoicePaymentSegment(inv, true, 30); g != "deposit" {
		t.Fatalf("deposit partial: got %q", g)
	}
	inv2 := &stripe.Invoice{AmountPaid: 5000}
	if g := guestInvoicePaymentSegment(inv2, false, 30); g != "balance" {
		t.Fatalf("balance: got %q", g)
	}
	if g := guestInvoicePaymentSegment(inv2, false, 0); g != "full" {
		t.Fatalf("no deposit pct: got %q", g)
	}
}
