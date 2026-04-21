package handlers

import "testing"

func TestGuestCheckoutApplicationFeeCents(t *testing.T) {
	tests := []struct {
		name       string
		amount     int64
		sellingFee float64
		want       int64
	}{
		{"zero amount", 0, 0.06, 0},
		{"100 USD at 6%", 10000, 0.06, 600},
		{"50 cents at 6%", 50, 0.06, 3},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := guestCheckoutApplicationFeeCents(tt.amount, tt.sellingFee)
			if got != tt.want {
				t.Fatalf("guestCheckoutApplicationFeeCents(%d, %v) = %d, want %d", tt.amount, tt.sellingFee, got, tt.want)
			}
		})
	}
}
