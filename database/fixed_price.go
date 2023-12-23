package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"redbudway-api/internal"
	"redbudway-api/models"
	"strconv"
)

func GetFixedPriceServiceDetails(priceID string) (*models.ServiceDetails, *models.Business, error) {
	fixedPrice := &models.ServiceDetails{}
	business := &models.Business{}

	stmt, err := db.Prepare("SELECT tp.name, tp.tradespersonId, ts.vanityURL, fp.id, fp.category, fp.subCategory, fp.title, fp.price, fp.description, fp.subscription, fp.subInterval, fp.selectPlaces, fp.timeZone, b.icon, ts.timeZone FROM fixed_prices fp INNER JOIN tradesperson_profile tp ON tp.tradespersonId=fp.tradespersonId INNER JOIN tradesperson_settings ts ON ts.tradespersonId=fp.tradespersonId INNER JOIN tradesperson_branding b ON b.tradespersonId=fp.tradespersonId WHERE fp.archived=false AND fp.priceId=?")
	if err != nil {
		return fixedPrice, business, err
	}
	defer stmt.Close()

	row := stmt.QueryRow(priceID)
	var fixedPriceID, price int64
	var vanityURL, interval, icon sql.NullString
	switch err = row.Scan(&business.Name, &business.TradespersonID, &vanityURL, &fixedPriceID, &fixedPrice.Category, &fixedPrice.SubCategory, &fixedPrice.Title, &price, &fixedPrice.Description, &fixedPrice.Subscription, &interval, &fixedPrice.SelectPlaces, &fixedPrice.TimeZone, &icon, &business.TimeZone); err {
	case sql.ErrNoRows:
		return fixedPrice, business, err
	case nil:
		business.Icon = icon.String
		business.VanityURL = vanityURL.String
		fixedPrice.Interval = interval.String

		strPrice := fmt.Sprintf("%.2f", float64(price)/float64(100.00))
		floatPrice, err := strconv.ParseFloat(strPrice, 64)
		if err != nil {
			return fixedPrice, business, err
		}
		fixedPrice.Price = floatPrice
		fixedPrice.Images, err = internal.GetImages(priceID, business.TradespersonID)
		if err != nil {
			return fixedPrice, business, err
		}
		if len(fixedPrice.Images) == 0 {
			url := "https://" + os.Getenv("SUBDOMAIN") + "redbudway.com/assets/images/placeholder.svg"
			fixedPrice.Images = append(fixedPrice.Images, url)
		}
		fixedPrice.TimeSlots, err = GetPublicTimeSlots(fixedPriceID, fixedPrice.Subscription)
		if err != nil {
			return fixedPrice, business, err
		}
		fixedPrice.Specialties, err = GetSpecialties(fixedPriceID)
		if err != nil {
			return fixedPrice, business, err
		}
		fixedPrice.Includes, fixedPrice.Excludes, err = GetIncludes(fixedPriceID)
		if err != nil {
			return fixedPrice, business, err
		}
		fixedPrice.Reviews, fixedPrice.Rating, err = GetFixedPriceReviewsRating(fixedPriceID)
		if err != nil {
			log.Printf("Failed to get reviews and rating %s", err)
		}
		fixedPrice.StatesAndCities, err = GetFixedPriceStatesAndCities(fixedPriceID)
		if err != nil {
			return fixedPrice, business, err
		}
		fixedPrice.Repeat, err = GetFixedPriceRepeatCustomers(fixedPriceID, business.TradespersonID)
		if err != nil {
			log.Printf("Failed to get fixed price repeat customers %s", err)
		}
		fixedPrice.Jobs, err = GetFixedPriceJobs(fixedPriceID, business.TradespersonID)
		if err != nil {
			log.Printf("Failed to get fixed price jobs %s", err)
		}

		fixedPrice.Form, err = GetFixedPriceForm(fixedPriceID)
		if err != nil {
			log.Printf("Failed to get fixed price form %s", err)
		}

	default:
		return fixedPrice, business, err
	}

	return fixedPrice, business, nil
}

func GetFixedPriceID(priceID string) (int64, error) {
	var fixedPriceID int64
	stmt, err := db.Prepare("SELECT id FROM fixed_prices WHERE priceId=?")
	if err != nil {
		return fixedPriceID, err
	}
	defer stmt.Close()

	row := stmt.QueryRow(priceID)
	switch err = row.Scan(&fixedPriceID); err {
	case sql.ErrNoRows:
		return fixedPriceID, err
	case nil:
		//
	default:
		log.Printf("Unknown, %v", err)
	}

	return fixedPriceID, nil
}

func GetFixedPriceInterval(priceID string) (string, error) {
	var interval string
	stmt, err := db.Prepare("SELECT subInterval FROM fixed_prices WHERE priceId=?")
	if err != nil {
		return interval, err
	}
	defer stmt.Close()

	row := stmt.QueryRow(priceID)
	switch err = row.Scan(&interval); err {
	case sql.ErrNoRows:
		return interval, err
	case nil:
		//
	default:
		log.Printf("Unknown, %v", err)
	}

	return interval, nil
}
