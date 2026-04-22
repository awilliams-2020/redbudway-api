package handlers

import (
	"strings"
	"testing"
)

func TestGuestCheckoutLineProductCopy(t *testing.T) {
	t.Parallel()
	t.Run("no deposit pct", func(t *testing.T) {
		t.Parallel()
		name, desc, pi := guestCheckoutLineProductCopy("Paint job", 0, false, 10000, 10000, 10000, 0)
		if !strings.Contains(name, "Full amount due now") {
			t.Fatalf("name: %q", name)
		}
		if !strings.Contains(desc, "no separate deposit") {
			t.Fatalf("desc: %q", desc)
		}
		if pi == "" {
			t.Fatal("empty pi line")
		}
	})
	t.Run("partial deposit first payment", func(t *testing.T) {
		t.Parallel()
		name, desc, pi := guestCheckoutLineProductCopy("Invoice payment", 20, true, 77500, 387500, 387500, 0)
		if !strings.Contains(name, "20% deposit") || !strings.Contains(name, "775.00") || !strings.Contains(name, "3875.00") {
			t.Fatalf("name: %q", name)
		}
		if !strings.Contains(desc, "$775.00 of $3875.00") {
			t.Fatalf("desc: %q", desc)
		}
		if !strings.Contains(desc, "Remaining balance after this payment") || !strings.Contains(desc, "$3100.00") {
			t.Fatalf("desc missing balance: %q", desc)
		}
		if strings.Contains(desc, "invoice") {
			t.Fatalf("desc should not mention invoice in checkout: %q", desc)
		}
		if strings.Contains(pi, "\n") {
			t.Fatalf("pi line should be single line: %q", pi)
		}
	})
	t.Run("remaining balance after deposit paid", func(t *testing.T) {
		t.Parallel()
		name, desc, pi := guestCheckoutLineProductCopy("Invoice payment", 20, false, 310000, 387500, 310000, 77500)
		if !strings.Contains(name, "Remaining balance due") {
			t.Fatalf("name: %q", name)
		}
		if desc != pi {
			t.Fatalf("desc and pi should match: %q vs %q", desc, pi)
		}
		if !strings.Contains(desc, "remaining") {
			t.Fatalf("desc: %q", desc)
		}
	})
	t.Run("full payment one shot 100pct", func(t *testing.T) {
		t.Parallel()
		name, desc, _ := guestCheckoutLineProductCopy("Invoice payment", 100, false, 50000, 50000, 50000, 0)
		if !strings.Contains(name, "Full payment (100% of total)") {
			t.Fatalf("name: %q", name)
		}
		if !strings.Contains(desc, "$500.00 of $500.00") {
			t.Fatalf("desc: %q", desc)
		}
	})
}
