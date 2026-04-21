package handlers

import (
	"github.com/stripe/stripe-go/v82"
)

// guestDepositChargeCents decides how much to collect in this Checkout session.
// When depositPct is 1–100, we treat that percentage of inv.Total as the deposit target; the customer pays
// toward that target until AmountPaid reaches it, then pays the full remaining balance (second payment can use hosted invoice).
// Returns ok=false if the slice is invalid (e.g. below Stripe minimum).
func guestDepositChargeCents(inv *stripe.Invoice, depositPct int64) (charge int64, partialDeposit bool, ok bool) {
	if inv == nil {
		return 0, false, false
	}
	rem := inv.AmountRemaining
	if rem < 50 {
		return 0, false, false
	}
	if depositPct <= 0 || depositPct > 100 {
		return rem, false, true
	}

	total := inv.Total
	if total <= 0 {
		total = inv.AmountPaid + rem
	}
	targetDeposit := (total * depositPct) / 100
	if inv.AmountPaid >= targetDeposit {
		// Deposit obligation satisfied; collect full remainder (often one shot, hosted invoice OK).
		return rem, false, true
	}
	still := targetDeposit - inv.AmountPaid
	if still < 0 {
		still = 0
	}
	charge = still
	if charge > rem {
		charge = rem
	}
	if charge < 50 {
		return 0, false, false
	}
	partialDeposit = charge < rem
	return charge, partialDeposit, true
}

func stripeInvoiceCustomerID(inv *stripe.Invoice) string {
	if inv == nil || inv.Customer == nil {
		return ""
	}
	return inv.Customer.ID
}
