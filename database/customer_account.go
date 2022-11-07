package database

import (
	"database/sql"
	"log"
	"redbudway-api/email"
	"redbudway-api/internal"
	"redbudway-api/restapi/operations"

	"github.com/gofrs/uuid"
	"github.com/stripe/stripe-go/v72"
)

func CreateCustomerAccount(_customer operations.PostCustomerBody, stripeAccount *stripe.Customer) (uuid.UUID, error) {
	log.Printf("Creating %s customer account", *_customer.Name)

	customerID, err := internal.GenerateUUID()
	if err != nil {
		return customerID, err
	}

	stmt, err := db.Prepare("INSERT INTO customer_account (stripeId, customerId, email, password) VALUES (?, ?, ?, ?)")
	if err != nil {
		return customerID, err
	}
	defer stmt.Close()

	passwordHash, err := internal.HashPassword(_customer.Password.String())
	if err != nil {
		return customerID, err
	}

	results, err := stmt.Exec(stripeAccount.ID, customerID, _customer.Email, passwordHash)
	if err != nil {
		return customerID, err
	}

	rowsAffected, err := results.RowsAffected()
	if err != nil {
		return customerID, err
	}

	if rowsAffected == 1 {
		if err := email.SendCustomerVerification(*_customer.Name, _customer.Email.String(), customerID.String()); err != nil {
			log.Printf("Failed to send customer verification email, %s", err)
		}
	}

	return customerID, nil
}

func GetCustomerStripeID(customerID string) (string, error) {
	var stripeID string

	stmt, err := db.Prepare("SELECT stripeId FROM customer_account WHERE customerId=?")
	if err != nil {
		return stripeID, err
	}
	defer stmt.Close()

	row := stmt.QueryRow(customerID)
	switch err = row.Scan(&stripeID); err {
	case sql.ErrNoRows:
		log.Printf("Customer with ID %s, doesn't exist", customerID)
		return stripeID, err
	case nil:
		//
	default:
		log.Printf("Unknown %v", err)
	}

	return stripeID, nil
}

func GetCustomerID(stripeID string) (string, error) {
	var customerID string

	stmt, err := db.Prepare("SELECT customerId FROM customer_account WHERE stripeId=?")
	if err != nil {
		return customerID, err
	}
	defer stmt.Close()

	row := stmt.QueryRow(stripeID)
	switch err = row.Scan(&customerID); err {
	case sql.ErrNoRows:
		log.Printf("Customer with stripe ID %s, doesn't exist", stripeID)
		return customerID, err
	case nil:
		//
	default:
		log.Printf("Unknown %v", err)
	}

	return customerID, nil
}

func DeleteCustomerAccount(customerID string) (bool, error) {
	stmt, err := db.Prepare("DELETE FROM customer_account WHERE customerId=?")
	if err != nil {
		return false, err
	}
	defer stmt.Close()

	results, err := stmt.Exec(customerID)
	if err != nil {
		return false, err
	}

	rowsAffected, err := results.RowsAffected()
	if err != nil {
		return false, err
	}
	return rowsAffected == 1, nil
}
