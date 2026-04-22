package database

import (
	"database/sql"
	"log"
	"strings"

	"github.com/stripe/stripe-go/v82"
	"github.com/stripe/stripe-go/v82/customer"
)

// GetOrCreateStripeCustomerOnConnect returns the Stripe Customer ID on the given Connect account
// for this app customer + tradesperson pair. Creates the Connect customer on first use.
func GetOrCreateStripeCustomerOnConnect(tradespersonID, customerID, connectAccountID string) (string, error) {
	tradespersonID = strings.TrimSpace(tradespersonID)
	customerID = strings.TrimSpace(customerID)
	connectAccountID = strings.TrimSpace(connectAccountID)
	if tradespersonID == "" || customerID == "" || connectAccountID == "" {
		return "", sql.ErrNoRows
	}

	var existing string
	err := GetConnection().QueryRow(
		`SELECT stripeCustomerId FROM customer_tradesperson_stripe WHERE customerId=? AND tradespersonId=?`,
		customerID, tradespersonID,
	).Scan(&existing)
	if err == nil && strings.TrimSpace(existing) != "" {
		return existing, nil
	}
	if err != nil && err != sql.ErrNoRows {
		return "", err
	}

	var email string
	err = GetConnection().QueryRow(`SELECT email FROM customer_account WHERE customerId=?`, customerID).Scan(&email)
	if err != nil {
		return "", err
	}
	email = strings.TrimSpace(email)
	if email == "" {
		return "", sql.ErrNoRows
	}

	cp := &stripe.CustomerParams{
		Email: stripe.String(email),
	}
	cp.SetStripeAccount(connectAccountID)
	cu, err := customer.New(cp)
	if err != nil {
		return "", err
	}

	_, err = GetConnection().Exec(
		`INSERT INTO customer_tradesperson_stripe (customerId, tradespersonId, stripeCustomerId) VALUES (?, ?, ?)`,
		customerID, tradespersonID, cu.ID,
	)
	if err != nil {
		log.Printf("customer_tradesperson_stripe insert %s %s: %v", customerID, tradespersonID, err)
		return "", err
	}
	return cu.ID, nil
}
