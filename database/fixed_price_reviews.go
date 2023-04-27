package database

import (
	"database/sql"
	"log"
	"redbudway-api/restapi/operations"

	"github.com/stripe/stripe-go/v72/customer"
)

func GetFixedPriceReviews(priceID string, page int64) ([]*operations.GetFixedPricePriceIDReviewsOKBodyReviewsItems0, error) {
	reviews := []*operations.GetFixedPricePriceIDReviewsOKBodyReviewsItems0{}

	stmt, err := db.Prepare("SELECT fp.tradespersonId, fpr.customerId, fpr.rating, fpr.message, DATE_FORMAT(fpr.date, '%M %D %Y') date, fpr.responded, fpr.respMsg, DATE_FORMAT(fpr.respDate, '%M %D %Y') respDate FROM fixed_price_reviews fpr INNER JOIN fixed_prices fp ON fpr.fixedPriceId=fp.id WHERE fp.priceId=? ORDER BY fpr.date DESC LIMIT ?, 10")
	if err != nil {
		log.Printf("Failed to create select statement %s", err)
		return reviews, err
	}
	defer stmt.Close()

	offSet := (page - 1) * 10
	rows, err := stmt.Query(priceID, offSet)
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

		review := &operations.GetFixedPricePriceIDReviewsOKBodyReviewsItems0{}
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
			tradesperson, err := GetTradespersonProfile(tradespersonID)
			if err != nil {
				log.Printf("Failed to get tradesperson profile %s", err)
				return reviews, err
			}
			review.Tradesperson = tradesperson.Name
		}
		review.Responded = responded
		reviews = append(reviews, review)
	}
	return reviews, nil
}

func GetFixedPriceRatings(priceID string) (operations.GetFixedPricePriceIDReviewsOKBody, error) {

	reviews := operations.GetFixedPricePriceIDReviewsOKBody{}

	stmt, err := db.Prepare("SELECT sum(case when fr.rating = '1' then 1 else 0 end) AS oneStars, sum(case when fr.rating = '2' then 1 else 0 end) AS twoStars, sum(case when fr.rating = '3' then 1 else 0 end) AS threeStars, sum(case when fr.rating = '4' then 1 else 0 end) AS fourStars, sum(case when fr.rating = '5' then 1 else 0 end) AS fiveStars  FROM fixed_price_reviews fr INNER JOIN fixed_prices f ON f.id=fr.fixedPriceId WHERE f.priceId=?")
	if err != nil {
		return reviews, err
	}
	defer stmt.Close()

	var oneStars, twoStars, threeStars, fourStars, fiveStars int64
	row := stmt.QueryRow(priceID)
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

func GetFixedPriceReviewsRating(fixedPriceId int64) (int64, float64, error) {
	reviews := int64(0)
	businessRating := float64(0.0)

	stmt, err := db.Prepare("SELECT fpr.rating FROM fixed_price_reviews fpr INNER JOIN fixed_prices fp ON fp.id=fpr.fixedPriceId WHERE fpr.fixedPriceId=?")
	if err != nil {
		return reviews, businessRating, err
	}
	defer stmt.Close()

	rows, err := stmt.Query(fixedPriceId)
	if err != nil {
		log.Printf("Failed to execute select statement %s", err)
		return reviews, businessRating, err
	}

	var customerRating float64
	total := float64(0)
	for rows.Next() {
		if err := rows.Scan(&customerRating); err != nil {
			return reviews, businessRating, err
		}
		total += customerRating
		reviews += 1
	}
	if reviews != 0 {
		businessRating = total / float64(reviews)
	}
	return reviews, businessRating, nil
}

func CustomerReviewedFixedPrice(customerID, priceID string) (bool, error) {
	reviewed := true
	stmt, err := db.Prepare("SELECT sr.rating FROM fixed_price_reviews sr INNER JOIN fixed_prices s ON s.id=sr.fixedPriceId WHERE sr.customerId=? AND s.priceId=?")
	if err != nil {
		log.Printf("Failed to create select statement %s", err)
		return reviewed, err
	}
	defer stmt.Close()

	var rating int64
	row := stmt.QueryRow(customerID, priceID)
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

func CustomerReviewedSubscription(customerID, priceID string) (bool, error) {
	reviewed := true
	stmt, err := db.Prepare("SELECT sr.rating FROM fixed_price_reviews sr INNER JOIN fixed_prices s ON s.id=sr.fixedPriceId WHERE sr.customerId=? AND s.priceId=?")
	if err != nil {
		log.Printf("Failed to create select statement %s", err)
		return reviewed, err
	}
	defer stmt.Close()

	var rating int64
	row := stmt.QueryRow(customerID, priceID)
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

func CreateFixedPriceReview(customerID, message string, fixedPriceID, rating int64) (bool, error) {
	stmt, err := db.Prepare("INSERT INTO fixed_price_reviews (fixedPriceId, customerId, rating, message, date) VALUES (?, ?, ?, ?, NOW())")
	if err != nil {
		return false, err
	}
	defer stmt.Close()

	results, err := stmt.Exec(fixedPriceID, customerID, rating, message)
	if err != nil {
		return false, err
	}

	rowsAffected, err := results.RowsAffected()
	if err != nil {
		return false, err
	}
	return rowsAffected == 1, nil
}
