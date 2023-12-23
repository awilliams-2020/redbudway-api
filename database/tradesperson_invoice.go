package database

import (
	"log"
	"time"
)

func SaveInvoice(stripeInvoiceID, customerID, tradpespersonID, timeZone string, fixedPriceID, created int64) (int64, error) {
	invoiceCreate := time.Unix(created, 0)
	var invoiceID int64
	stmt, err := db.Prepare("INSERT INTO tradesperson_invoices (invoiceId, fixedPriceId, customerId, tradespersonId, created, timeZone) VALUES (?, ?, ?, ?, ?, ?)")
	if err != nil {
		return invoiceID, err
	}
	defer stmt.Close()

	results, err := stmt.Exec(stripeInvoiceID, fixedPriceID, customerID, tradpespersonID, invoiceCreate, timeZone)
	if err != nil {
		return invoiceID, err
	}

	invoiceID, err = results.LastInsertId()
	if err != nil {
		return invoiceID, err
	}

	return invoiceID, nil
}

func DeleteInvoice(tradespersonID, invoiceID string) (bool, error) {
	stmt, err := db.Prepare("DELETE FROM tradesperson_invoices WHERE tradespersonId=? AND invoiceId=?")
	if err != nil {
		return false, err
	}
	defer stmt.Close()

	results, err := stmt.Exec(tradespersonID, invoiceID)
	if err != nil {
		return false, err
	}

	rowsAffected, err := results.RowsAffected()
	if err != nil {
		return false, err
	}

	return rowsAffected == 1, nil
}

func SaveManualInvoice(stripeInvoiceID, customerID, tradpespersonID string, created int64) error {
	invoiceCreate := time.Unix(created, 0)
	stmt, err := db.Prepare("INSERT INTO tradesperson_manual_invoices (invoiceId, customerId, tradespersonId, created) VALUES (?, ?, ?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(stripeInvoiceID, customerID, tradpespersonID, invoiceCreate)
	if err != nil {
		return err
	}

	return nil
}

func OpenInvoices(tradespersonID string) (bool, error) {
	stmt, err := db.Prepare("SELECT invoiceId FROM tradesperson_invoices WHERE tradespersonId=? AND QUARTER(created) = ? AND YEAR(created) = ? GROUP BY id ORDER BY created DESC")
	if err != nil {
		log.Printf("Failed to create prepare statement, %v", err)
		return false, err
	}
	defer stmt.Close()

	rows, err := stmt.Query(tradespersonID)
	if err != nil {
		log.Printf("Failed to execute select statement %s", err)
		return false, err
	}

	var cuStripeID string
	for rows.Next() {
		if err := rows.Scan(&cuStripeID); err != nil {
			log.Printf("Failed to scan row %v", err)
			return false, err
		}
	}

	return false, nil
}
