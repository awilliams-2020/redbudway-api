package database

import (
	"database/sql"
	"log"

	"redbudway-api/restapi/operations"

	"github.com/stripe/stripe-go/v72/customer"
)

func GetQuoteReviews(quoteID string, page int64) ([]*operations.GetQuoteQuoteIDReviewsOKBodyReviewsItems0, error) {
	reviews := []*operations.GetQuoteQuoteIDReviewsOKBodyReviewsItems0{}

	stmt, err := db.Prepare("SELECT q.tradespersonId, qr.customerId, qr.rating, qr.message, DATE_FORMAT(qr.date, '%M %D %Y') date, qr.responded, qr.respMsg, DATE_FORMAT(qr.respDate, '%M %D %Y') respDate FROM quote_reviews qr INNER JOIN quotes q ON qr.quoteId=q.id WHERE q.quote=? ORDER BY qr.date DESC LIMIT ?, 10")
	if err != nil {
		log.Printf("Failed to create select statement %s", err)
		return reviews, err
	}
	defer stmt.Close()

	offSet := (page - 1) * 10
	rows, err := stmt.Query(quoteID, offSet)
	if err != nil {
		log.Printf("Failed to execute select statement %s", err)
		return reviews, err
	}

	var customerID, tradespersonID, date string
	var message, respMsg, respDate sql.NullString
	var rating int64
	var responded bool
	for rows.Next() {
		if err := rows.Scan(&tradespersonID, &customerID, &rating, &message, &date, &responded, &respMsg, &respDate); err != nil {
			log.Printf("Failed to scan for profile fixed prices, %s", err)
			return reviews, err
		}

		cuStripeID, err := GetCustomerStripeID(customerID)
		if err != nil {
			log.Printf("Failed to get customer %s account, %v", customerID, err)
			return reviews, err
		}

		stripeCustomer, err := customer.Get(cuStripeID, nil)
		if err != nil {
			log.Printf("Failed to get stripe customer %s, %v", customerID, err)
			return reviews, err
		}

		review := &operations.GetQuoteQuoteIDReviewsOKBodyReviewsItems0{}
		review.Customer = stripeCustomer.Name
		review.Rating = rating
		if message.Valid {
			review.Message = message.String
		}
		if respMsg.Valid {
			review.RespMsg = respMsg.String
		}
		if respDate.Valid {
			review.RespDate = respDate.String
		}
		review.Date = date
		if responded {
			tradesperson, err := GetTradespersonAccount(tradespersonID)
			if err != nil {
				log.Printf("Failed to get tradesperson account %s, %v", tradespersonID, err)
				return reviews, err
			}
			review.Tradesperson = tradesperson.Name
		}
		review.Responded = responded
		reviews = append(reviews, review)
	}
	return reviews, nil
}

func GetQuoteRatings(quoteID string) (operations.GetQuoteQuoteIDReviewsOKBody, error) {

	reviews := operations.GetQuoteQuoteIDReviewsOKBody{}

	stmt, err := db.Prepare("SELECT sum(case when qr.rating = '1' then 1 else 0 end) AS oneStars, sum(case when qr.rating = '2' then 1 else 0 end) AS twoStars, sum(case when qr.rating = '3' then 1 else 0 end) AS threeStars, sum(case when qr.rating = '4' then 1 else 0 end) AS fourStars, sum(case when qr.rating = '5' then 1 else 0 end) AS fiveStars  FROM quote_reviews qr INNER JOIN quotes q ON q.id=qr.quoteId WHERE q.quote=?")
	if err != nil {
		return reviews, err
	}
	defer stmt.Close()

	var oneStars, twoStars, threeStars, fourStars, fiveStars int64
	row := stmt.QueryRow(quoteID)
	switch err = row.Scan(&oneStars, &twoStars, &threeStars, &fourStars, &fiveStars); err {
	case sql.ErrNoRows:
		return reviews, err
	case nil:
		reviews.OneStars = oneStars
		reviews.TwoStars = twoStars
		reviews.ThreeStars = threeStars
		reviews.FourStars = fourStars
		reviews.FiveStars = fiveStars
	default:
		log.Printf("Unknown %v", err)
	}

	return reviews, err
}

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

func CreateQuoteReview(customerID, message string, ID, rating int64) (bool, error) {
	stmt, err := db.Prepare("INSERT INTO quote_reviews (quoteId, customerId, rating, message, date) VALUES (?, ?, ?, ?, NOW())")
	if err != nil {
		return false, err
	}
	defer stmt.Close()

	log.Printf("ID: %s", ID)
	results, err := stmt.Exec(ID, customerID, rating, message)
	if err != nil {
		return false, err
	}

	rowsAffected, err := results.RowsAffected()
	if err != nil {
		return false, err
	}
	return rowsAffected == 1, nil
}
