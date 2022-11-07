package database

import (
	"log"
	"time"
)

func SaveQuote(ID, customerID, tradespersonID, message string, quoteID, created int64) (bool, error) {
	quoteCreated := time.Unix(created, 0)

	stmt, err := db.Prepare("INSERT INTO tradesperson_quotes (quote, quoteId, customerId, tradespersonId, request, created) VALUES (?, ?, ?, ?, ?, ?)")
	if err != nil {
		return false, err
	}
	defer stmt.Close()

	results, err := stmt.Exec(ID, quoteID, customerID, tradespersonID, message, quoteCreated)
	if err != nil {
		return false, err
	}

	numRows, err := results.RowsAffected()
	if err != nil {
		log.Printf("Failed to get rows affected, %v", err)
	}
	return numRows == 1, nil
}

func DeleteQuote(tradespersonID, quoteID string) (bool, error) {
	stmt, err := db.Prepare("DELETE FROM tradesperson_quotes WHERE tradespersonId=? AND quote=?")
	if err != nil {
		return false, err
	}
	defer stmt.Close()

	results, err := stmt.Exec(tradespersonID, quoteID)
	if err != nil {
		return false, err
	}

	numRows, err := results.RowsAffected()
	if err != nil {
		log.Printf("Failed to get rows affected, %v", err)
	}
	return numRows == 1, nil
}

func UpdateQuote(tradespersonID, revisedQuoteID, quoteID string) (bool, error) {
	stmt, err := db.Prepare("UPDATE tradesperson_quotes SET quote=? WHERE tradespersonId=? AND quote=?")
	if err != nil {
		return false, err
	}
	defer stmt.Close()

	results, err := stmt.Exec(revisedQuoteID, tradespersonID, quoteID)
	if err != nil {
		return false, err
	}

	numRows, err := results.RowsAffected()
	if err != nil {
		log.Printf("Failed to get rows affected, %v", err)
	}
	return numRows == 1, nil
}
