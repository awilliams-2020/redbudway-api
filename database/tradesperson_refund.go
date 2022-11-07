package database

import (
	"database/sql"
	"log"

	"github.com/stripe/stripe-go/v72/refund"
)

func CreateInvoiceRefund(invoiceID, refundID string) error {
	stmt, err := db.Prepare("INSERT INTO tradesperson_refunds (invoiceId, refundId) VALUES (?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(invoiceID, refundID)
	if err != nil {
		return err
	}

	return nil
}

func GetInvoiceRefund(invoiceID string) (string, int64, error) {
	var status string
	var refunded int64
	stmt, err := db.Prepare("SELECT refundId FROM tradesperson_refunds WHERE invoiceId=?")
	if err != nil {
		return status, refunded, err
	}
	defer stmt.Close()

	var refundID string
	row := stmt.QueryRow(invoiceID)
	switch err = row.Scan(&refundID); err {
	case sql.ErrNoRows:
		break
	case nil:
		stripeRefund, err := refund.Get(
			refundID,
			nil,
		)
		if err != nil {
			log.Printf("Failed to get stripe refund %s", refundID)
			return status, refunded, err
		}
		status = stripeRefund.Object
		refunded = stripeRefund.Created
	default:
		log.Printf("Unknown %v", err)
	}
	return status, refunded, nil
}
