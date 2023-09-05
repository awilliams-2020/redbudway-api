package database

import (
	"database/sql"
	"log"
	"redbudway-api/internal"
	"redbudway-api/models"
	"redbudway-api/restapi/operations"

	"github.com/gofrs/uuid"
	"github.com/stripe/stripe-go/v72"
)

func GetTradespersonSellingFee(tradespersonID string) (float64, error) {
	var fee float64

	stmt, err := db.Prepare("SELECT fee FROM selling_fee WHERE tradespersonId=?")
	if err != nil {
		return fee, err
	}
	defer stmt.Close()

	row := stmt.QueryRow(tradespersonID)
	switch err = row.Scan(&fee); err {
	case sql.ErrNoRows:
		return fee, err
	case nil:
		//
	default:
		log.Printf("Unknown %v", err)
	}

	return fee, err
}

func CreateTradespersonAccount(tradesperson operations.PostTradespersonBody, stripeAccount *stripe.Account) (uuid.UUID, error) {
	log.Printf("Creating %s tradesperson account", tradesperson.Email)

	tradespersonID, err := internal.GenerateUUID()
	if err != nil {
		return tradespersonID, err
	}

	stmt, err := db.Prepare("INSERT INTO tradesperson_account (tradespersonId, email, password, stripeId) VALUES (?, ?, ?, ?)")
	if err != nil {
		return tradespersonID, err
	}
	defer stmt.Close()

	passwordHash, err := internal.HashPassword(tradesperson.Password.String())
	if err != nil {
		return tradespersonID, err
	}

	results, err := stmt.Exec(tradespersonID, tradesperson.Email, passwordHash, stripeAccount.ID)
	if err != nil {
		return tradespersonID, err
	}

	rowsAffected, err := results.RowsAffected()
	if err != nil {
		return tradespersonID, err
	}
	if rowsAffected == 1 {
		stmt, err := db.Prepare("INSERT INTO tradesperson_profile (tradespersonId, name, email) VALUES (?, ?, ?)")
		if err != nil {
			return tradespersonID, err
		}
		defer stmt.Close()

		results, err := stmt.Exec(tradespersonID, tradespersonID, tradesperson.Email)
		if err != nil {
			return tradespersonID, err
		}
		rowsAffected, err = results.RowsAffected()
		if err != nil {
			return tradespersonID, err
		}
		if rowsAffected == 1 {
			stmt, err := db.Prepare("INSERT INTO tradesperson_settings (tradespersonId, vanityURL) VALUES (?, ?)")
			if err != nil {
				return tradespersonID, err
			}
			defer stmt.Close()

			results, err := stmt.Exec(tradespersonID, tradespersonID)
			if err != nil {
				return tradespersonID, err
			}
			rowsAffected, err = results.RowsAffected()
			if err != nil {
				return tradespersonID, err
			}
			if rowsAffected == 1 {
				stmt, err := db.Prepare("INSERT INTO selling_fee (tradespersonId, fee, limited, expire) VALUES (?, 0.00, TRUE, NOW() + INTERVAL 3 MONTH )")
				if err != nil {
					return tradespersonID, err
				}
				defer stmt.Close()

				results, err := stmt.Exec(tradespersonID)
				if err != nil {
					return tradespersonID, err
				}
				rowsAffected, err = results.RowsAffected()
				if err != nil {
					return tradespersonID, err
				}
			}
		}
	}
	return tradespersonID, nil
}

func GetTradespersonProfile(tradespersonID string) (models.Tradesperson, error) {
	tradesperson := models.Tradesperson{}
	tradesperson.Address = &models.Address{}

	stmt, err := db.Prepare("SELECT name, image, description, email, number, lineOne, lineTwo, city, state, zipCode FROM tradesperson_profile WHERE tradespersonId=?")
	if err != nil {
		return tradesperson, err
	}
	defer stmt.Close()

	row := stmt.QueryRow(tradespersonID)
	switch err = row.Scan(
		&tradesperson.Name,
		&tradesperson.Image,
		&tradesperson.Description,
		&tradesperson.Email,
		&tradesperson.Number,
		&tradesperson.Address.LineOne,
		&tradesperson.Address.LineTwo,
		&tradesperson.Address.City,
		&tradesperson.Address.State,
		&tradesperson.Address.ZipCode); err {
	case sql.ErrNoRows:
		//
	case nil:
		//
	default:
		log.Printf("Unknown %v", err)
	}

	return tradesperson, err
}

func UpdateTradespersonProfileImage(tradespersonID, image string) error {
	imageURL, err := internal.SaveProfileImage(tradespersonID, image)
	if err != nil {
		log.Println("Failed to save profile image, %s", err)
	}

	stmt, err := db.Prepare("UPDATE tradesperson_profile SET image =? WHERE tradespersonId=?")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(imageURL, tradespersonID)
	if err != nil {
		return err
	}

	return nil
}

func UpdateTradespersonProfileName(tradespersonID, name string) error {
	stmt, err := db.Prepare("UPDATE tradesperson_profile SET name =? WHERE tradespersonId=?")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(name, tradespersonID)
	if err != nil {
		return err
	}

	return nil
}

func UpdateTradespersonProfileEmail(tradespersonID, email string) error {
	stmt, err := db.Prepare("UPDATE tradesperson_profile SET email =? WHERE tradespersonId=?")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(email, tradespersonID)
	if err != nil {
		return err
	}

	return nil
}

func UpdateTradespersonProfileNumber(tradespersonID, number string) error {
	stmt, err := db.Prepare("UPDATE tradesperson_profile SET number =? WHERE tradespersonId=?")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(number, tradespersonID)
	if err != nil {
		return err
	}

	return nil
}

func UpdateTradespersonProfileDescription(tradespersonID, description string) error {
	stmt, err := db.Prepare("UPDATE tradesperson_profile SET description =? WHERE tradespersonId=?")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(description, tradespersonID)
	if err != nil {
		return err
	}

	return nil
}

func UpdateTradespersonProfileAddress(tradespersonID string, address *models.Address) error {
	stmt, err := db.Prepare("UPDATE tradesperson_profile SET lineOne =?, lineTwo=?, city=?, state=?, zipCode=? WHERE tradespersonId=?")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(address.LineOne, address.LineTwo, address.City, address.State, address.ZipCode, tradespersonID)
	if err != nil {
		return err
	}

	return nil
}

func GetTradespersonAccountByPriceID(priceID string) (models.Tradesperson, string, string, error) {
	tradesperson := models.Tradesperson{}
	var stripeID, tradespersonID string

	stmt, err := db.Prepare("SELECT t.stripeId, t.tradespersonId FROM tradesperson_account t INNER JOIN fixed_prices s ON s.tradespersonId=t.tradespersonId WHERE s.priceId=?")
	if err != nil {
		return tradesperson, stripeID, tradespersonID, err
	}
	defer stmt.Close()

	row := stmt.QueryRow(priceID)
	switch err = row.Scan(&stripeID, &tradespersonID); err {
	case sql.ErrNoRows:
		return tradesperson, stripeID, tradespersonID, err
	case nil:
		tradesperson, err = GetTradespersonProfile(tradespersonID)
		if err != nil {
			log.Printf("Failed to get tradesperson profile %s", err)
		}
	default:
		log.Printf("Unknown %v", err)
	}

	return tradesperson, stripeID, tradespersonID, err
}

func GetTradespersonAccountByQuoteID(quoteID string) (models.Tradesperson, string, string, error) {
	tradesperson := models.Tradesperson{}
	var stripeID, tradespersonID string

	stmt, err := db.Prepare("SELECT t.stripeId, t.tradespersonId FROM tradesperson_account t INNER JOIN quotes q ON q.tradespersonId=t.tradespersonId WHERE q.quote=?")
	if err != nil {
		return tradesperson, stripeID, tradespersonID, err
	}
	defer stmt.Close()

	row := stmt.QueryRow(quoteID)
	switch err = row.Scan(&stripeID, &tradespersonID); err {
	case sql.ErrNoRows:
		return tradesperson, stripeID, tradespersonID, err
	case nil:
		tradesperson, err = GetTradespersonProfile(tradespersonID)
		if err != nil {
			log.Printf("Failed to get tradesperson profile %s", err)
		}
	default:
		log.Printf("Unknown %v", err)
	}

	return tradesperson, stripeID, tradespersonID, err
}

func UpdateTradespersonDisplaySettings(tradespersonID string, settings operations.PutTradespersonTradespersonIDSettingsBody) (bool, error) {
	stmt, err := db.Prepare("UPDATE tradesperson_settings SET email=?, number=?, address=? WHERE tradespersonId=?")
	if err != nil {
		return false, err
	}
	defer stmt.Close()

	_, err = stmt.Exec(settings.DisplayEmail, settings.DisplayNumber, settings.DisplayAddress, tradespersonID)
	if err != nil {
		return false, err
	}
	return true, nil
}

func UpdateTradespersonVanitySettings(tradespersonID string, settings operations.PutTradespersonTradespersonIDSettingsBody) (bool, error) {
	stmt, err := db.Prepare("UPDATE tradesperson_settings SET vanityURL=? WHERE tradespersonId=?")
	if err != nil {
		return false, err
	}
	defer stmt.Close()

	_, err = stmt.Exec(settings.VanityURL, tradespersonID)
	if err != nil {
		return false, err
	}
	return true, nil
}

func GetTradespersonJobs(tradespersonID string) (int64, error) {
	jobs := int64(0)
	count := int64(0)
	stmt, err := db.Prepare("SELECT COUNT(*) FROM tradesperson_invoices WHERE tradespersonId=?")
	if err != nil {
		return jobs, err
	}
	defer stmt.Close()

	if err := stmt.QueryRow(tradespersonID).Scan(&count); err != nil {
		return jobs, err
	}
	jobs += count

	stmt, err = db.Prepare("SELECT COUNT(*) FROM tradesperson_subscriptions WHERE tradespersonId=?")
	if err != nil {
		return jobs, err
	}
	defer stmt.Close()

	if err := stmt.QueryRow(tradespersonID).Scan(&count); err != nil {
		return 0, err
	}
	jobs += count

	return jobs, nil
}

func GetTradespersonRepeatCustomers(tradespersonID string) (int64, error) {
	repeat := int64(0)
	stmt, err := db.Prepare("SELECT COUNT(*) > 1 FROM tradesperson_invoices WHERE tradespersonId=? GROUP BY customerId")
	if err != nil {
		return repeat, err
	}
	defer stmt.Close()

	rows, err := stmt.Query(tradespersonID)
	if err != nil {
		return repeat, err
	}

	for rows.Next() {
		isRepeat := false
		if err := rows.Scan(&isRepeat); err != nil {
			continue
		}
		if isRepeat {
			repeat += 1
		}
	}

	return repeat, nil
}

func GetTradespersonRatingReviews(tradespersonID string) (int64, int64, error) {

	stmt, err := db.Prepare("SELECT sr.rating FROM fixed_price_reviews sr INNER JOIN fixed_prices s ON s.id=sr.fixedPriceId WHERE s.tradespersonId=?")
	if err != nil {
		return 0, 0, err
	}
	defer stmt.Close()

	rows, err := stmt.Query(tradespersonID)
	if err != nil {
		return 0, 0, err
	}
	total := int64(0)
	reviews := int64(0)
	rating := int64(0)
	for rows.Next() {
		if err := rows.Scan(&rating); err != nil {
			return 0, 0, err
		}
		total += rating
		reviews += int64(1)
	}

	stmt, err = db.Prepare("SELECT qr.rating FROM quote_reviews qr INNER JOIN quotes q ON q.id=qr.quoteId WHERE q.tradespersonId=?")
	if err != nil {
		return 0, 0, err
	}
	defer stmt.Close()

	rows, err = stmt.Query(tradespersonID)
	if err != nil {
		return 0, 0, err
	}

	for rows.Next() {
		if err := rows.Scan(&rating); err != nil {
			return 0, 0, err
		}
		total += rating
		reviews += int64(1)
	}

	if reviews != int64(0) {
		rating = total / reviews
	}

	return rating, reviews, nil
}

func GetTradespersonStripeID(tradespersonID string) (string, error) {
	var stripeID string

	stmt, err := db.Prepare("SELECT stripeId FROM tradesperson_account WHERE tradespersonId=?")
	if err != nil {
		return stripeID, err
	}
	defer stmt.Close()

	row := stmt.QueryRow(tradespersonID)
	switch err = row.Scan(&stripeID); err {
	case sql.ErrNoRows:
		log.Printf("Tradesperson %s, doesn't exist", tradespersonID)
		return stripeID, err
	case nil:
		//
	default:
		log.Printf("Unknown %v", err)
	}

	return stripeID, nil
}

func DeleteTradespersonAccount(tradepersonID, stripeID string) (bool, error) {
	stmt, err := db.Prepare("DELETE FROM tradesperson_account WHERE tradespersonId=? AND stripeId=?")
	if err != nil {
		return false, err
	}
	defer stmt.Close()

	results, err := stmt.Exec(tradepersonID, stripeID)
	if err != nil {
		return false, err
	}

	rowsAffected, err := results.RowsAffected()
	if err != nil {
		return false, err
	}
	return rowsAffected == 1, nil
}
