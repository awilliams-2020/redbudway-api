package database

import "time"

func SaveSubscription(stripeSubscriptionID, cuStripeID, tradpespersonID string, fixedPriceID, created int64) (int64, error) {
	subscriptionCreated := time.Unix(created, 0)
	var subscriptionID int64
	stmt, err := db.Prepare("INSERT INTO tradesperson_subscriptions (subscriptionId, fixedPriceId, cuStripeId, tradespersonId, created) VALUES (?, ?, ?, ?, ?)")
	if err != nil {
		return subscriptionID, err
	}
	defer stmt.Close()

	results, err := stmt.Exec(stripeSubscriptionID, fixedPriceID, cuStripeID, tradpespersonID, subscriptionCreated)
	if err != nil {
		return subscriptionID, err
	}

	subscriptionID, err = results.LastInsertId()
	if err != nil {
		return subscriptionID, err
	}

	return subscriptionID, nil
}

func DeleteSubscription(subscriptionID, tradespersonID string) (bool, error) {
	stmt, err := db.Prepare("DELETE FROM tradesperson_subscriptions WHERE tradespersonId=? AND subscriptionId=?")
	if err != nil {
		return false, err
	}
	defer stmt.Close()

	results, err := stmt.Exec(tradespersonID, subscriptionID)
	if err != nil {
		return false, err
	}

	rowsAffected, err := results.RowsAffected()
	if err != nil {
		return false, err
	}

	return rowsAffected == 1, nil
}
