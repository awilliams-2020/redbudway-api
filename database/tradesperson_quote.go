package database

import (
	"database/sql"
	"log"
	"strings"
	"time"
)

// GuestAcceptQuotePageDetails is safe copy for the public accept-quote page (email link flow).
type GuestAcceptQuotePageDetails struct {
	Request       string
	ServiceName   string
	ProviderEmail string
}

// GetGuestAcceptQuotePageDetails loads provider email, catalog service label, and customer request for a Stripe billing quote id.
func GetGuestAcceptQuotePageDetails(tradespersonID, stripeQuoteID string) (GuestAcceptQuotePageDetails, error) {
	var d GuestAcceptQuotePageDetails
	var req, title, sub, email sql.NullString
	err := db.QueryRow(
		`SELECT tq.request, q.title, q.subcategory, p.email
		FROM tradesperson_quotes tq
		INNER JOIN quotes q ON tq.quoteId = q.id
		INNER JOIN tradesperson_profile p ON p.tradespersonId = tq.tradespersonId
		WHERE tq.tradespersonId = ? AND tq.quote = ?`,
		tradespersonID, stripeQuoteID,
	).Scan(&req, &title, &sub, &email)
	if err == sql.ErrNoRows {
		return d, nil
	}
	if err != nil {
		return d, err
	}
	if req.Valid {
		d.Request = strings.TrimSpace(req.String)
	}
	if email.Valid {
		d.ProviderEmail = strings.TrimSpace(email.String)
	}
	if title.Valid {
		t := strings.TrimSpace(title.String)
		s := ""
		if sub.Valid {
			s = strings.TrimSpace(sub.String)
		}
		if t != "" && s != "" {
			d.ServiceName = t + " — " + s
		} else if t != "" {
			d.ServiceName = t
		} else if s != "" {
			d.ServiceName = s
		}
	}
	return d, nil
}

// GetTradespersonIDByBillingStripeQuote returns the provider id for a Stripe billing quote id (qt_…) stored on tradesperson_quotes.quote.
func GetTradespersonIDByBillingStripeQuote(stripeQuoteID string) (string, error) {
	var tid string
	err := db.QueryRow("SELECT tradespersonId FROM tradesperson_quotes WHERE quote = ? LIMIT 1", stripeQuoteID).Scan(&tid)
	if err != nil {
		return "", err
	}
	return tid, nil
}

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
