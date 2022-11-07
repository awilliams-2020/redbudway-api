package database

import (
	"database/sql"
	"encoding/base64"
	"fmt"
	"log"
	"os"
	"redbudway-api/internal"
	"redbudway-api/models"
	"redbudway-api/restapi/operations"
	"strings"

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
	log.Printf("Creating %v tradesperson account", tradesperson.Name)

	tradespersonID, err := internal.GenerateUUID()
	if err != nil {
		return tradespersonID, err
	}

	stmt, err := db.Prepare("INSERT INTO tradesperson_account (tradespersonId, name, description, number, email, image, password, stripeId) VALUES (?, ?, ?, ?, ?, '', ?, ?)")
	if err != nil {
		return tradespersonID, err
	}
	defer stmt.Close()

	passwordHash, err := internal.HashPassword(tradesperson.Password.String())
	if err != nil {
		return tradespersonID, err
	}

	results, err := stmt.Exec(tradespersonID, tradesperson.Name, tradesperson.Description, tradesperson.Number, tradesperson.Email, passwordHash, stripeAccount.ID)
	if err != nil {
		return tradespersonID, err
	}

	rowsAffected, err := results.RowsAffected()
	if err != nil {
		return tradespersonID, err
	}
	if rowsAffected == 1 {
		stmt, err := db.Prepare("INSERT INTO tradesperson_settings (tradespersonId, vanityURL) VALUES (?, '')")
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
	return tradespersonID, nil
}

func GetTradespersonAccount(tradespersonID string) (models.Tradesperson, error) {
	tradesperson := models.Tradesperson{}

	stmt, err := db.Prepare("SELECT email, number, name FROM tradesperson_account WHERE tradespersonId=?")
	if err != nil {
		return tradesperson, err
	}
	defer stmt.Close()

	row := stmt.QueryRow(tradespersonID)
	switch err = row.Scan(&tradesperson.Email, &tradesperson.Number, &tradesperson.Name); err {
	case sql.ErrNoRows:
		return tradesperson, err
	case nil:
		//
	default:
		log.Printf("Unknown %v", err)
	}

	return tradesperson, err
}

func GetTradespersonAccountByPriceID(priceID string) (models.Tradesperson, string, string, error) {
	tradesperson := models.Tradesperson{}
	var stripeID, tradespersonID string

	stmt, err := db.Prepare("SELECT m.stripeId, m.tradespersonId, m.email, m.number, m.name FROM tradesperson_account m INNER JOIN fixed_prices s ON s.tradespersonId=m.tradespersonId WHERE s.priceId=?")
	if err != nil {
		return tradesperson, stripeID, tradespersonID, err
	}
	defer stmt.Close()

	row := stmt.QueryRow(priceID)
	switch err = row.Scan(&stripeID, &tradespersonID, &tradesperson.Email, &tradesperson.Number, &tradesperson.Name); err {
	case sql.ErrNoRows:
		return tradesperson, stripeID, tradespersonID, err
	case nil:
		//
	default:
		log.Printf("Unknown %v", err)
	}

	return tradesperson, stripeID, tradespersonID, err
}

func GetTradespersonAccountByQuoteID(quoteID string) (models.Tradesperson, string, string, error) {
	tradesperson := models.Tradesperson{}
	var stripeID, tradespersonID string

	stmt, err := db.Prepare("SELECT t.stripeId, t.tradespersonId, t.email, t.number, t.name FROM tradesperson_account t INNER JOIN quotes q ON q.tradespersonId=t.tradespersonId WHERE q.quote=?")
	if err != nil {
		return tradesperson, stripeID, tradespersonID, err
	}
	defer stmt.Close()

	row := stmt.QueryRow(quoteID)
	switch err = row.Scan(&stripeID, &tradespersonID, &tradesperson.Email, &tradesperson.Number, &tradesperson.Name); err {
	case sql.ErrNoRows:
		return tradesperson, stripeID, tradespersonID, err
	case nil:
		//
	default:
		log.Printf("Unknown %v", err)
	}

	return tradesperson, stripeID, tradespersonID, err
}

func UpdateTradespersonDescription(tradespersonID, description string) error {

	stmt, err := db.Prepare("UPDATE tradesperson_account SET description =? WHERE tradespersonId=?")
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

func UpdateTradespersonImage(tradespersonID, image string) error {
	data := strings.Split(image, ",")

	dec, err := base64.StdEncoding.DecodeString(data[1])
	if err != nil {
		log.Println("Failed to decode")
		return err
	}
	format := ""
	switch data[0] {
	case "data:image/jpeg;base64":
		format = ".jpeg"
	case "data:image/png;base64":
		format = ".png"
	case "data:image/webp;base64":
		format = ".webp"
	}

	path := fmt.Sprintf("%s/%s", "images", tradespersonID)
	//add to Util package
	if _, err := os.Stat(path); os.IsNotExist(err) {
		err = os.MkdirAll(path, 0755)
		if err != nil {
			return err
		}
	}

	fileName := fmt.Sprintf("%s/%s%s", path, tradespersonID, format)
	f, err := os.Create(fileName)
	if err != nil {
		log.Println("Failed to create file with name %s", fileName)
		return err
	}
	defer f.Close()

	if _, err := f.Write(dec); err != nil {
		return err
	}
	if err := f.Sync(); err != nil {
		return err
	}

	stmt, err := db.Prepare("UPDATE tradesperson_account SET image =? WHERE tradespersonId=?")
	if err != nil {
		return err
	}
	defer stmt.Close()

	url := fmt.Sprintf("https://"+os.Getenv("SUBDOMAIN")+"redbudway.com/%s", fileName)
	_, err = stmt.Exec(url, tradespersonID)
	if err != nil {
		return err
	}

	return nil
}

func UpdateTradespersonDisplaySettings(tradespersonID string, settings operations.PutTradespersonTradespersonIDSettingsBody) (bool, error) {
	stmt, err := db.Prepare("UPDATE tradesperson_settings SET displayEmail=?, displayNumber=?, displayAddress=? WHERE tradespersonId=?")
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
