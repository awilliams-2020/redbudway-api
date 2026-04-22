package handlers

import (
	"fmt"
	"log"
	"strings"

	"redbudway-api/models"

	"github.com/stripe/stripe-go/v82"
	"github.com/stripe/stripe-go/v82/charge"
	"github.com/stripe/stripe-go/v82/invoice"
	"github.com/stripe/stripe-go/v82/refund"
)

// invoicePaymentSegmentMetaKey is stored on guest Checkout PaymentIntents for refund UX.
const invoicePaymentSegmentMetaKey = "invoice_payment_segment"

// guestInvoicePaymentSegment classifies a guest Checkout payment (deposit vs balance vs full).
func guestInvoicePaymentSegment(inv *stripe.Invoice, partialDeposit bool, displayDepositPct int64) string {
	dp := displayDepositPct
	if dp < 0 || dp > 100 {
		dp = 0
	}
	if dp == 0 {
		return "full"
	}
	if partialDeposit {
		return "deposit"
	}
	if inv != nil && inv.AmountPaid > 0 {
		return "balance"
	}
	return "full"
}

func loadStripeInvoiceForBillingQuoteExpand(connectAccountID, invoiceID string, expand []*string) (*stripe.Invoice, error) {
	if connectAccountID != "" {
		p := &stripe.InvoiceParams{
			Params: stripe.Params{Expand: expand},
		}
		p.SetStripeAccount(connectAccountID)
		inv, err := invoice.Get(invoiceID, p)
		if err == nil {
			return inv, nil
		}
		log.Printf("invoice.Get %s with Connect account: %v", invoiceID, err)
	}
	p := &stripe.InvoiceParams{
		Params: stripe.Params{Expand: expand},
	}
	return invoice.Get(invoiceID, p)
}

func chargeIDFromStripeInvoicePayment(ip *stripe.InvoicePayment) string {
	if ip == nil || ip.Payment == nil {
		return ""
	}
	if ip.Payment.Type == stripe.InvoicePaymentPaymentTypeCharge && ip.Payment.Charge != nil && ip.Payment.Charge.ID != "" {
		return ip.Payment.Charge.ID
	}
	if ip.Payment.PaymentIntent != nil && ip.Payment.PaymentIntent.LatestCharge != nil && ip.Payment.PaymentIntent.LatestCharge.ID != "" {
		return ip.Payment.PaymentIntent.LatestCharge.ID
	}
	if ip.Payment.Charge != nil && ip.Payment.Charge.ID != "" {
		return ip.Payment.Charge.ID
	}
	return ""
}

func paymentIntentIDFromStripeInvoicePayment(ip *stripe.InvoicePayment) string {
	if ip == nil || ip.Payment == nil || ip.Payment.PaymentIntent == nil {
		return ""
	}
	return ip.Payment.PaymentIntent.ID
}

func segmentFromInvoicePaymentPI(pi *stripe.PaymentIntent) string {
	if pi == nil || pi.Metadata == nil {
		return ""
	}
	if s := strings.TrimSpace(pi.Metadata[invoicePaymentSegmentMetaKey]); s != "" {
		return strings.ToLower(s)
	}
	if pi.Metadata[guestCheckoutMetadataInvoiceAttach] == "true" {
		return "deposit"
	}
	return ""
}

func cardPaymentLabel(segment string, index int) string {
	switch segment {
	case "deposit":
		return "Deposit payment"
	case "balance":
		return "Balance payment"
	case "full":
		return "Card payment"
	default:
		if index == 0 {
			return "Card payment"
		}
		return fmt.Sprintf("Card payment %d", index+1)
	}
}

func quoteUsesReverseTransfer(sq *stripe.Quote) bool {
	return sq != nil && sq.TransferData != nil && sq.TransferData.Destination != nil && sq.TransferData.Destination.ID != ""
}

func stripeChargeGetForQuote(chargeID string, reverseTransfer bool, tpStripeID string) (*stripe.Charge, error) {
	cp := &stripe.ChargeParams{}
	if !reverseTransfer && tpStripeID != "" {
		cp.SetStripeAccount(tpStripeID)
	}
	return charge.Get(chargeID, cp)
}

// BuildQuoteInvoiceCardPayments lists paid card charges on the invoice with refundable amounts.
func BuildQuoteInvoiceCardPayments(stripeInv *stripe.Invoice, sq *stripe.Quote, tpStripeID string) ([]*models.BillingQuoteInvoiceCardPayment, int64, int64, error) {
	if stripeInv == nil || stripeInv.Payments == nil {
		return []*models.BillingQuoteInvoiceCardPayment{}, 0, 0, nil
	}
	reverseTransfer := quoteUsesReverseTransfer(sq)

	var rows []*models.BillingQuoteInvoiceCardPayment
	var totalRefundable, totalRefunded int64
	seen := map[string]struct{}{}
	idx := 0

	for _, ip := range stripeInv.Payments.Data {
		if ip == nil || !strings.EqualFold(ip.Status, "paid") {
			continue
		}
		chID := chargeIDFromStripeInvoicePayment(ip)
		if chID == "" {
			continue
		}
		if _, dup := seen[chID]; dup {
			continue
		}
		seen[chID] = struct{}{}

		ch, err := stripeChargeGetForQuote(chID, reverseTransfer, tpStripeID)
		if err != nil {
			return nil, 0, 0, fmt.Errorf("charge %s: %w", chID, err)
		}
		refunded := ch.AmountRefunded
		refundable := ch.Amount - refunded
		if refundable < 0 {
			refundable = 0
		}

		seg := ""
		if ip.Payment != nil && ip.Payment.PaymentIntent != nil {
			seg = segmentFromInvoicePaymentPI(ip.Payment.PaymentIntent)
		}
		if seg == "" {
			seg = "unknown"
		}

		piID := paymentIntentIDFromStripeInvoicePayment(ip)
		row := &models.BillingQuoteInvoiceCardPayment{
			ChargeID:              chID,
			PaymentIntentID:       piID,
			AmountCents:           ch.Amount,
			AmountRefundedCents:   refunded,
			AmountRefundableCents: refundable,
			Segment:               seg,
			Label:                 cardPaymentLabel(seg, idx),
		}
		rows = append(rows, row)
		totalRefundable += refundable
		totalRefunded += refunded
		idx++
	}
	return rows, totalRefundable, totalRefunded, nil
}

func quoteRefundIdempotencyKey(invoiceID, chargeID string, refundCents int64) string {
	s := fmt.Sprintf("qb-refund-%s-%s-%d", invoiceID, chargeID, refundCents)
	if len(s) > 255 {
		return s[:255]
	}
	return s
}

// QuoteInvoiceStripeRefund refunds one charge. refundAmountCents 0 means full remaining refundable on that charge.
func QuoteInvoiceStripeRefund(
	chargeID string,
	refundAmountCents int64,
	sq *stripe.Quote,
	tpStripeID string,
	invoiceID string,
) (*stripe.Refund, int64, error) {
	reverseTransfer := quoteUsesReverseTransfer(sq)
	ch, err := stripeChargeGetForQuote(chargeID, reverseTransfer, tpStripeID)
	if err != nil {
		return nil, 0, err
	}
	refundable := ch.Amount - ch.AmountRefunded
	if refundable < 0 {
		refundable = 0
	}
	if refundable == 0 {
		return nil, 0, fmt.Errorf("nothing refundable on this charge")
	}
	amountToRefund := refundable
	if refundAmountCents > 0 {
		if refundAmountCents > refundable {
			return nil, 0, fmt.Errorf("refund amount exceeds refundable balance on charge")
		}
		amountToRefund = refundAmountCents
	}

	rp := &stripe.RefundParams{
		Charge: stripe.String(chargeID),
	}
	if amountToRefund < refundable {
		rp.Amount = stripe.Int64(amountToRefund)
	}
	if reverseTransfer {
		rp.ReverseTransfer = stripe.Bool(true)
	} else {
		rp.SetStripeAccount(tpStripeID)
	}
	rp.SetIdempotencyKey(quoteRefundIdempotencyKey(invoiceID, chargeID, amountToRefund))

	sr, err := refund.New(rp)
	if err != nil {
		return nil, 0, err
	}
	return sr, amountToRefund, nil
}
