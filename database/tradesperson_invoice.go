package database

import "time"

func SaveInvoice(stripeInvoiceID, customerID, tradpespersonID string, fixedPriceID, created int64) (int64, error) {
	invoiceCreate := time.Unix(created, 0)
	var invoiceID int64
	stmt, err := db.Prepare("INSERT INTO tradesperson_invoices (invoiceId, fixedPriceId, customerId, tradespersonId, created) VALUES (?, ?, ?, ?, ?)")
	if err != nil {
		return invoiceID, err
	}
	defer stmt.Close()

	results, err := stmt.Exec(stripeInvoiceID, fixedPriceID, customerID, tradpespersonID, invoiceCreate)
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
