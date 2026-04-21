package handlers

import "github.com/stripe/stripe-go/v82"

// chargeIDFromInvoice returns the first charge ID found among an invoice's paid payments.
// Invoice.Charge and Invoice.PaymentIntent were removed in the 2025-03-31 API; charge IDs
// are now accessed via Invoice.Payments (expand "payments" when fetching the invoice).
func chargeIDFromInvoice(inv *stripe.Invoice) string {
	if inv.Payments == nil {
		return ""
	}
	for _, p := range inv.Payments.Data {
		if p.Status != "paid" || p.Payment == nil {
			continue
		}
		if p.Payment.Charge != nil && p.Payment.Charge.ID != "" {
			return p.Payment.Charge.ID
		}
		if p.Payment.PaymentIntent != nil && p.Payment.PaymentIntent.LatestCharge != nil {
			return p.Payment.PaymentIntent.LatestCharge.ID
		}
	}
	return ""
}
