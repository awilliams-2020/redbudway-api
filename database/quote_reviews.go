package database

import (
	"database/sql"
	"log"
)

func CustomerReviewedQuote(customerID string, ID int64) (bool, error) {
	reviewed := true
	stmt, err := db.Prepare("SELECT qr.rating FROM quote_reviews qr INNER JOIN quotes q ON q.id=qr.quoteId WHERE qr.customerId=? AND qr.quoteId=?")
	if err != nil {
		log.Printf("Failed to create select statement %s", err)
		return reviewed, err
	}
	defer stmt.Close()

	var rating int64
	row := stmt.QueryRow(customerID, ID)
	switch err = row.Scan(&rating); err {
	case sql.ErrNoRows:
		reviewed = false
	case nil:
		//
	default:
		log.Printf("Unknown %v", err)
	}

	return reviewed, nil
}

func CreateQuoteReview(customerID, message string, quoteID, rating int64) (bool, error) {
	stmt, err := db.Prepare("INSERT INTO quote_reviews (quoteId, customerId, rating, message, date) VALUES (?, ?, ?, ?, NOW())")
	if err != nil {
		return false, err
	}
	defer stmt.Close()

	results, err := stmt.Exec(quoteID, customerID, rating, message)
	if err != nil {
		return false, err
	}

	rowsAffected, err := results.RowsAffected()
	if err != nil {
		return false, err
	}
	return rowsAffected == 1, nil
}
