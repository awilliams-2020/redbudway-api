package database

import (
	"database/sql"
	"fmt"
	"log"
	"redbudway-api/models"
	"redbudway-api/restapi/operations"
	"strconv"

	"github.com/stripe/stripe-go/v72/price"
	"github.com/stripe/stripe-go/v72/product"
)

func GetFixedPriceServiceDetails(priceID, state, city string) (*models.ServiceDetails, *operations.GetFixedPricePriceIDOKBodyBusiness, error) {
	fixedPrice := &models.ServiceDetails{}
	business := &operations.GetFixedPricePriceIDOKBodyBusiness{}

	stmt, err := db.Prepare("SELECT ta.name, ta.tradespersonId, ts.vanityURL, fp.id, fp.category, fp.subCategory, fp.subscription, fp.subInterval FROM fixed_prices fp INNER JOIN tradesperson_account ta ON ta.tradespersonId=fp.tradespersonId INNER JOIN tradesperson_settings ts ON ts.tradespersonId=fp.tradespersonId LEFT JOIN fixed_price_state_cities fpsc ON fpsc.fixedPriceId=fp.id WHERE (fp.selectPlaces=false OR fpsc.state=?) AND (fp.selectPlaces=false OR JSON_CONTAINS(fpsc.cities, JSON_OBJECT('name', ?))) AND fp.archived=false AND fp.priceId=?")
	if err != nil {
		return fixedPrice, business, err
	}
	defer stmt.Close()

	row := stmt.QueryRow(state, city, priceID)
	var fixedPriceID int64
	var vanityURL, interval sql.NullString
	var name, tradespersonID, category, subCategory string
	var subscription, archived bool
	switch err = row.Scan(&name, &tradespersonID, &vanityURL, &fixedPriceID, &category, &subCategory, &subscription, &interval); err {
	case sql.ErrNoRows:
		return fixedPrice, business, err
	case nil:

		business.Name = name
		business.VanityURL = vanityURL.String
		business.TradespersonID = tradespersonID

		fixedPrice.Category = &category
		fixedPrice.SubCategory = subCategory
		fixedPrice.Subscription = subscription
		if interval.Valid {
			fixedPrice.Interval = interval.String
		}
		fixedPrice.Archived = archived
		p, err := price.Get(priceID, nil)
		if err != nil {
			return fixedPrice, business, err
		}
		pr, err := product.Get(p.Product.ID, nil)
		if err != nil {
			return fixedPrice, business, err
		}
		strPrice := fmt.Sprintf("%.2f", p.UnitAmountDecimal/float64(100.00))
		floatPrice, err := strconv.ParseFloat(strPrice, 64)
		if err != nil {
			return fixedPrice, business, err
		}
		fixedPrice.Price = floatPrice
		fixedPrice.Images = pr.Images
		fixedPrice.Title = &pr.Name
		fixedPrice.Description = &pr.Description
		fixedPrice.TimeSlots, err = GetPublicTimeSlots(fixedPriceID, subscription)
		if err != nil {
			return fixedPrice, business, err
		}
		fixedPrice.Filters, err = GetFilters(fixedPriceID)
		if err != nil {
			return fixedPrice, business, err
		}

		fixedPrice.Reviews, fixedPrice.Rating, err = GetFixedPriceReviewsRating(fixedPriceID)
		if err != nil {
			log.Printf("Failed to get reviews and rating %s", err)
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
