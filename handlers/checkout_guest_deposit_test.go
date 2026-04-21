package handlers

import (
	"testing"

	"github.com/stripe/stripe-go/v82"
)

func TestGuestDepositChargeCents(t *testing.T) {
	inv := &stripe.Invoice{
		Total:           10000,
		AmountPaid:      0,
		AmountRemaining: 10000,
	}
	got, partial, ok := guestDepositChargeCents(inv, 30)
	if !ok || got != 3000 || !partial {
		t.Fatalf("30%% of 100: got %d partial=%v ok=%v", got, partial, ok)
	}

	inv2 := &stripe.Invoice{
		Total:           10000,
		AmountPaid:      3000,
		AmountRemaining: 7000,
	}
	got2, partial2, ok2 := guestDepositChargeCents(inv2, 30)
	if !ok2 || got2 != 7000 || partial2 {
		t.Fatalf("after deposit: want full remainder 7000 partial=false, got %d partial=%v", got2, partial2)
	}

	inv3 := &stripe.Invoice{
		Total:           10000,
		AmountPaid:      0,
		AmountRemaining: 10000,
	}
	got3, partial3, ok3 := guestDepositChargeCents(inv3, 0)
	if !ok3 || got3 != 10000 || partial3 {
		t.Fatalf("0%% deposit: want full amount, got %d partial=%v", got3, partial3)
	}
}
