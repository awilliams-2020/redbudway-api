package database

import (
	"database/sql"
	"log"
	"redbudway-api/email"
	"redbudway-api/internal"
	"redbudway-api/restapi/operations"
	"time"

	"github.com/gofrs/uuid"
	"github.com/stripe/stripe-go/v72"
	"github.com/stripe/stripe-go/v72/coupon"
	"github.com/stripe/stripe-go/v72/promotioncode"
)

func CreateCustomerAccount(_customer operations.PostCustomerBody, stripeAccount *stripe.Customer) (uuid.UUID, error) {
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
		token, err := internal.GenerateToken(customerID.String(), "customer", "verification", time.Minute*15)
		if err != nil {
			log.Printf("Failed to generate JWT, %s", err)
			return customerID, err
		}
		if err := email.SendCustomerVerification(*_customer.Name, _customer.Email.String(), customerID.String(), token); err != nil {
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

func GetDiscount(priceID, code string) (*operations.GetCustomerCustomerIDPromoOKBody, error) {
	discount := &operations.GetCustomerCustomerIDPromoOKBody{}
	var promoID string
	stmt, err := db.Prepare("SELECT tp.couponId, tp.promoId, tp.code, tc.discountType, tc.duration, tc.months FROM tradesperson_promos tp INNER JOIN tradesperson_coupons tc ON tp.couponId=tc.couponId WHERE JSON_CONTAINS(tc.services, ?) AND tp.code=? AND tp.active=true")
	if err != nil {
		log.Printf("Failed to create prepare statement, %v", err)
		return discount, err
	}
	defer stmt.Close()
	row := stmt.QueryRow("\""+priceID+"\"", code)
	switch err = row.Scan(&discount.CouponID, &promoID, &discount.Code, &discount.Type, &discount.Duration, &discount.Months); err {
	case sql.ErrNoRows:
	case nil:
		couponParams := &stripe.CouponParams{}
		stripeCoupon, err := coupon.Get(discount.CouponID, couponParams)
		if err != nil {
			log.Printf("Failed to retrieve coupon %s, %v", discount.CouponID, err)
			return discount, err
		}
		codeParams := &stripe.PromotionCodeParams{}
		stripePromo, err := promotioncode.Get(promoID, codeParams)
		if err != nil {
			log.Printf("Failed to retrieve promo %s, %v", promoID, err)
			return discount, err
		}
		discount.Valid = stripePromo.Active
		if stripePromo.Active {
			if discount.Type == "percent_off" {
				discount.Percent = stripeCoupon.PercentOff
			} else {
				discount.Amount = float64(stripeCoupon.AmountOff) / float64(100)
			}
		}
	default:
		log.Printf("Unknown %v", err)
	}

	return discount, nil
}
