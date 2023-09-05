package handlers

import (
	"database/sql"
	"fmt"
	"log"
	"math"
	"redbudway-api/database"
	"redbudway-api/internal"
	"redbudway-api/models"
	"redbudway-api/restapi/operations"
	"strconv"
	"strings"

	"github.com/go-openapi/runtime/middleware"
)

const PAGE_SIZE = float64(9)

//  AND ((SELECT COUNT(*) FROM fixed_price_time_slots WHERE fixedPriceId=f.id AND startTime > CURDATE()) > 0  OR f.subscription=True)

func processFixedPriceRows(db *sql.DB, rows *sql.Rows, fixedPrices []*models.Service) ([]*models.Service, error) {
	var id, price int64
	var vanityURL, interval sql.NullString
	var stripeID, tradespersonID string
	for rows.Next() {
		fixedPrice := &models.Service{}
		if err := rows.Scan(&stripeID, &id, &tradespersonID, &fixedPrice.PriceID, &fixedPrice.Title, &price, &fixedPrice.Subscription, &interval, &vanityURL); err != nil {
			return fixedPrices, err
		}

		tradesperson, err := database.GetTradespersonProfile(tradespersonID)
		if err != nil {
			log.Printf("Failed to get tradesperson profile %s", err)
		}
		fixedPrice.Business = tradesperson.Name
		fixedPrice.Interval = interval.String
		strPrice := fmt.Sprintf("%.2f", float64(price)/float64(100.00))
		floatPrice, err := strconv.ParseFloat(strPrice, 64)
		if err != nil {
			log.Printf("Failed to parse float, %v", err)
			return fixedPrices, err
		}
		fixedPrice.Price = floatPrice
		fixedPrice.Image, err = internal.GetImage(fixedPrice.PriceID, tradespersonID)
		if err != nil {
			log.Printf("Failed to get fixedPrice image %s", err)
		}
		fixedPrice.VanityURL = vanityURL.String
		fixedPrice.TradespersonID = tradespersonID
		fixedPrice.AvailableTimeSlots, err = database.GetAvailableTimeSlots(id, fixedPrice.Subscription)
		if err != nil {
			log.Printf("Failed to get timeslots %s", err)
		}

		fixedPrice.Reviews, fixedPrice.Rating, err = database.GetFixedPriceReviewsRating(id)
		if err != nil {
			log.Printf("Failed to get reviews and rating %s", err)
		}

		repeat, err := database.GetFixedPriceRepeatCustomers(id, tradespersonID)
		if err != nil {
			log.Printf("Failed to get fixed price repeat customers %s", err)
		}
		fixedPrice.Repeat = repeat

		jobs, err := database.GetFixedPriceJobs(id, tradespersonID)
		if err != nil {
			log.Printf("Failed to get fixed price jobs %s", err)
		}
		fixedPrice.Jobs = jobs

		fixedPrices = append(fixedPrices, fixedPrice)
	}
	return fixedPrices, nil
}

func getFixedPricesWithFilters(state, city, category, subCategory, sort string, specialties, fromDate, toDate *string, page int64, min, max *int64) ([]*models.Service, error) {
	query := ""
	fixedPrices := []*models.Service{}
	db := database.GetConnection()

	if specialties != nil {
		specialtyArry := strings.Split(*specialties, ",")
		for _, specialty := range specialtyArry {
			query += "'" + specialty + "',"
		}
		query = query[:len(query)-1]
		query = "AND fi.specialty IN (" + query + ")"
	}

	rows := &sql.Rows{}
	if specialties != nil && min == nil && max == nil && fromDate == nil && toDate == nil {
		stmt, err := db.Prepare("SELECT a.stripeId, f.id, f.tradespersonId, f.priceId, f.title, f.price, f.subscription, f.subInterval, s.vanityURL FROM tradesperson_account a INNER JOIN tradesperson_settings s ON a.tradespersonId=s.tradespersonId INNER JOIN fixed_prices f ON a.tradespersonId=f.tradespersonId LEFT JOIN fixed_price_state_cities c ON c.fixedPriceId=f.id LEFT JOIN fixed_price_specialties fi ON fi.fixedPriceId=f.id WHERE f.category=? AND f.subcategory=? AND ((? = '' AND ? = '') OR f.selectPlaces=false OR c.state=? OR JSON_CONTAINS(c.cities, JSON_OBJECT('name', ?))) " + query + " AND f.archived=false GROUP BY f.id ORDER BY " + sort + " LIMIT ?, ?")
		if err != nil {
			log.Printf("Failed to create select statement %s", err)
			return fixedPrices, err
		}
		defer stmt.Close()

		offset := (page - 1) * int64(PAGE_SIZE)
		rows, err = stmt.Query(category, subCategory, state, city, state, city, offset, PAGE_SIZE)
		if err != nil {
			log.Printf("Failed to execute select statement %s", err)
			return fixedPrices, err
		}
	} else if specialties != nil && min != nil && max == nil && fromDate == nil && toDate == nil {
		stmt, err := db.Prepare("SELECT a.stripeId, f.id, f.tradespersonId, f.priceId, f.title, f.price, f.subscription, f.subInterval, s.vanityURL FROM tradesperson_account a INNER JOIN tradesperson_settings s ON a.tradespersonId=s.tradespersonId INNER JOIN fixed_prices f ON a.tradespersonId=f.tradespersonId LEFT JOIN fixed_price_state_cities c ON c.fixedPriceId=f.id LEFT JOIN fixed_price_specialties fi ON fi.fixedPriceId=f.id WHERE f.category=? AND f.subcategory=? AND ((? = '' AND ? = '') OR f.selectPlaces=false OR c.state=? OR JSON_CONTAINS(c.cities, JSON_OBJECT('name', ?))) " + query + " AND f.archived=false AND f.price >= ? GROUP BY f.id ORDER BY " + sort + " LIMIT ?, ?")
		if err != nil {
			log.Printf("Failed to create select statement %s", err)
			return fixedPrices, err
		}
		defer stmt.Close()

		offset := (page - 1) * int64(PAGE_SIZE)
		rows, err = stmt.Query(category, subCategory, state, city, state, city, *min, offset, PAGE_SIZE)
		if err != nil {
			log.Printf("Failed to execute select statement %s", err)
			return fixedPrices, err
		}
	} else if specialties != nil && min != nil && max != nil && fromDate == nil && toDate == nil {
		stmt, err := db.Prepare("SELECT a.stripeId, f.id, f.tradespersonId, f.priceId, f.title, f.price, f.subscription, f.subInterval, s.vanityURL FROM tradesperson_account a INNER JOIN tradesperson_settings s ON a.tradespersonId=s.tradespersonId INNER JOIN fixed_prices f ON a.tradespersonId=f.tradespersonId LEFT JOIN fixed_price_state_cities c ON c.fixedPriceId=f.id LEFT JOIN fixed_price_specialties fi ON fi.fixedPriceId=f.id WHERE f.category=? AND f.subcategory=? AND ((? = '' AND ? = '') OR f.selectPlaces=false OR c.state=? OR JSON_CONTAINS(c.cities, JSON_OBJECT('name', ?))) " + query + " AND f.archived=false AND f.price >= ? AND f.price <= ? GROUP BY f.id ORDER BY " + sort + " LIMIT ?, ?")
		if err != nil {
			log.Printf("Failed to create select statement %s", err)
			return fixedPrices, err
		}
		defer stmt.Close()

		offset := (page - 1) * int64(PAGE_SIZE)
		rows, err = stmt.Query(category, subCategory, state, city, state, city, *min, *max, offset, PAGE_SIZE)
		if err != nil {
			log.Printf("Failed to execute select statement %s", err)
			return fixedPrices, err
		}
	} else if specialties != nil && min != nil && max != nil && fromDate != nil && toDate == nil {
		stmt, err := db.Prepare("SELECT a.stripeId, f.id, f.tradespersonId, f.priceId, f.title, f.price, f.subscription, f.subInterval, s.vanityURL FROM tradesperson_account a INNER JOIN tradesperson_settings s ON a.tradespersonId=s.tradespersonId INNER JOIN fixed_prices f ON a.tradespersonId=f.tradespersonId LEFT JOIN fixed_price_state_cities c ON c.fixedPriceId=f.id LEFT JOIN fixed_price_specialties fi ON fi.fixedPriceId=f.id WHERE f.category=? AND f.subcategory=? AND ((? = '' AND ? = '') OR f.selectPlaces=false OR c.state=? OR JSON_CONTAINS(c.cities, JSON_OBJECT('name', ?))) " + query + " AND f.archived=false AND f.price >= ? AND f.price <= ? AND f.id=(SELECT fixedPriceId from fixed_price_time_slots WHERE startTime >= DATE(?) AND fixedPriceId=f.id LIMIT 1) GROUP BY f.id ORDER BY " + sort + " LIMIT ?, ?")
		if err != nil {
			log.Printf("Failed to create select statement %s", err)
			return fixedPrices, err
		}
		defer stmt.Close()

		offset := (page - 1) * int64(PAGE_SIZE)
		rows, err = stmt.Query(category, subCategory, state, city, state, city, *min, *max, *fromDate, offset, PAGE_SIZE)
		if err != nil {
			log.Printf("Failed to execute select statement %s", err)
			return fixedPrices, err
		}
	} else if specialties != nil && min != nil && max != nil && fromDate != nil && toDate != nil {
		stmt, err := db.Prepare("SELECT a.stripeId, f.id, f.tradespersonId, f.priceId, f.title, f.price, f.subscription, f.subInterval, s.vanityURL FROM tradesperson_account a INNER JOIN tradesperson_settings s ON a.tradespersonId=s.tradespersonId INNER JOIN fixed_prices f ON a.tradespersonId=f.tradespersonId LEFT JOIN fixed_price_state_cities c ON c.fixedPriceId=f.id LEFT JOIN fixed_price_specialties fi ON fi.fixedPriceId=f.id WHERE f.category=? AND f.subcategory=? AND ((? = '' AND ? = '') OR f.selectPlaces=false OR c.state=? OR JSON_CONTAINS(c.cities, JSON_OBJECT('name', ?))) " + query + " AND f.archived=false AND f.price >= ? AND f.price <= ? AND f.id=(SELECT fixedPriceId from fixed_price_time_slots WHERE startTime >= DATE(?) AND startTime <= DATE(?) AND fixedPriceId=f.id LIMIT 1) GROUP BY f.id ORDER BY " + sort + " LIMIT ?, ?")
		if err != nil {
			log.Printf("Failed to create select statement %s", err)
			return fixedPrices, err
		}
		defer stmt.Close()

		offset := (page - 1) * int64(PAGE_SIZE)
		rows, err = stmt.Query(category, subCategory, state, city, state, city, *min, *max, *fromDate, *toDate, offset, PAGE_SIZE)
		if err != nil {
			log.Printf("Failed to execute select statement %s", err)
			return fixedPrices, err
		}
	} else if specialties != nil && min == nil && max != nil && fromDate == nil && toDate == nil {
		stmt, err := db.Prepare("SELECT a.stripeId, f.id, f.tradespersonId, f.priceId, f.title, f.price, f.subscription, f.subInterval, s.vanityURL FROM tradesperson_account a INNER JOIN tradesperson_settings s ON a.tradespersonId=s.tradespersonId INNER JOIN fixed_prices f ON a.tradespersonId=f.tradespersonId LEFT JOIN fixed_price_state_cities c ON c.fixedPriceId=f.id LEFT JOIN fixed_price_specialties fi ON fi.fixedPriceId=f.id WHERE f.category=? AND f.subcategory=? AND ((? = '' AND ? = '') OR f.selectPlaces=false OR c.state=? OR JSON_CONTAINS(c.cities, JSON_OBJECT('name', ?))) " + query + " AND f.archived=false AND f.price <= ? GROUP BY f.id ORDER BY " + sort + " LIMIT ?, ?")
		if err != nil {
			log.Printf("Failed to create select statement %s", err)
			return fixedPrices, err
		}
		defer stmt.Close()

		offset := (page - 1) * int64(PAGE_SIZE)
		rows, err = stmt.Query(category, subCategory, state, city, state, city, *max, offset, PAGE_SIZE)
		if err != nil {
			log.Printf("Failed to execute select statement %s", err)
			return fixedPrices, err
		}
	} else if specialties != nil && min == nil && max != nil && fromDate != nil && toDate == nil {
		stmt, err := db.Prepare("SELECT a.stripeId, f.id, f.tradespersonId, f.priceId, f.title, f.price, f.subscription, f.subInterval, s.vanityURL FROM tradesperson_account a INNER JOIN tradesperson_settings s ON a.tradespersonId=s.tradespersonId INNER JOIN fixed_prices f ON a.tradespersonId=f.tradespersonId LEFT JOIN fixed_price_state_cities c ON c.fixedPriceId=f.id LEFT JOIN fixed_price_specialties fi ON fi.fixedPriceId=f.id WHERE f.category=? AND f.subcategory=? AND ((? = '' AND ? = '') OR f.selectPlaces=false OR c.state=? OR JSON_CONTAINS(c.cities, JSON_OBJECT('name', ?))) " + query + " AND f.archived=false AND f.price <= ? AND f.id=(SELECT fixedPriceId from fixed_price_time_slots WHERE startTime >= DATE(?) AND fixedPriceId=f.id LIMIT 1) GROUP BY f.id ORDER BY " + sort + " LIMIT ?, ?")
		if err != nil {
			log.Printf("Failed to create select statement %s", err)
			return fixedPrices, err
		}
		defer stmt.Close()

		offset := (page - 1) * int64(PAGE_SIZE)
		rows, err = stmt.Query(category, subCategory, state, city, state, city, *max, *fromDate, offset, PAGE_SIZE)
		if err != nil {
			log.Printf("Failed to execute select statement %s", err)
			return fixedPrices, err
		}
	} else if specialties != nil && min == nil && max != nil && fromDate != nil && toDate != nil {
		stmt, err := db.Prepare("SELECT a.stripeId, f.id, f.tradespersonId, f.priceId, f.title, f.price, f.subscription, f.subInterval, s.vanityURL FROM tradesperson_account a INNER JOIN tradesperson_settings s ON a.tradespersonId=s.tradespersonId INNER JOIN fixed_prices f ON a.tradespersonId=f.tradespersonId LEFT JOIN fixed_price_state_cities c ON c.fixedPriceId=f.id LEFT JOIN fixed_price_specialties fi ON fi.fixedPriceId=f.id WHERE f.category=? AND f.subcategory=? AND ((? = '' AND ? = '') OR f.selectPlaces=false OR c.state=? OR JSON_CONTAINS(c.cities, JSON_OBJECT('name', ?))) " + query + " AND f.archived=false AND f.price <= ? AND f.id=(SELECT fixedPriceId from fixed_price_time_slots WHERE startTime >= DATE(?) AND startTime <= DATE(?) AND fixedPriceId=f.id LIMIT 1) GROUP BY f.id ORDER BY " + sort + " LIMIT ?, ?")
		if err != nil {
			log.Printf("Failed to create select statement %s", err)
			return fixedPrices, err
		}
		defer stmt.Close()

		offset := (page - 1) * int64(PAGE_SIZE)
		rows, err = stmt.Query(category, subCategory, state, city, state, city, *max, *fromDate, *toDate, offset, PAGE_SIZE)
		if err != nil {
			log.Printf("Failed to execute select statement %s", err)
			return fixedPrices, err
		}
	} else if specialties != nil && min == nil && max == nil && fromDate != nil && toDate == nil {
		stmt, err := db.Prepare("SELECT a.stripeId, f.id, f.tradespersonId, f.priceId, f.title, f.price, f.subscription, f.subInterval, s.vanityURL FROM tradesperson_account a INNER JOIN tradesperson_settings s ON a.tradespersonId=s.tradespersonId INNER JOIN fixed_prices f ON a.tradespersonId=f.tradespersonId LEFT JOIN fixed_price_state_cities c ON c.fixedPriceId=f.id LEFT JOIN fixed_price_specialties fi ON fi.fixedPriceId=f.id WHERE f.category=? AND f.subcategory=? AND ((? = '' AND ? = '') OR f.selectPlaces=false OR c.state=? OR JSON_CONTAINS(c.cities, JSON_OBJECT('name', ?))) " + query + " AND f.archived=false AND f.id=(SELECT fixedPriceId from fixed_price_time_slots WHERE startTime >= DATE(?) AND fixedPriceId=f.id LIMIT 1) GROUP BY f.id ORDER BY " + sort + " LIMIT ?, ?")
		if err != nil {
			log.Printf("Failed to create select statement %s", err)
			return fixedPrices, err
		}
		defer stmt.Close()

		offset := (page - 1) * int64(PAGE_SIZE)
		rows, err = stmt.Query(category, subCategory, state, city, state, city, *fromDate, offset, PAGE_SIZE)
		if err != nil {
			log.Printf("Failed to execute select statement %s", err)
			return fixedPrices, err
		}
	} else if specialties != nil && min == nil && max == nil && fromDate != nil && toDate != nil {
		stmt, err := db.Prepare("SELECT a.stripeId, f.id, f.tradespersonId, f.priceId, f.title, f.price, f.subscription, f.subInterval, s.vanityURL FROM tradesperson_account a INNER JOIN tradesperson_settings s ON a.tradespersonId=s.tradespersonId INNER JOIN fixed_prices f ON a.tradespersonId=f.tradespersonId LEFT JOIN fixed_price_state_cities c ON c.fixedPriceId=f.id LEFT JOIN fixed_price_specialties fi ON fi.fixedPriceId=f.id WHERE f.category=? AND f.subcategory=? AND ((? = '' AND ? = '') OR f.selectPlaces=false OR c.state=? OR JSON_CONTAINS(c.cities, JSON_OBJECT('name', ?))) " + query + " AND f.archived=false AND f.id=(SELECT fixedPriceId from fixed_price_time_slots WHERE startTime >= DATE(?) AND startTime <= DATE(?) AND fixedPriceId=f.id LIMIT 1) GROUP BY f.id ORDER BY " + sort + " LIMIT ?, ?")
		if err != nil {
			log.Printf("Failed to create select statement %s", err)
			return fixedPrices, err
		}
		defer stmt.Close()

		offset := (page - 1) * int64(PAGE_SIZE)
		rows, err = stmt.Query(category, subCategory, state, city, state, city, *fromDate, *toDate, offset, PAGE_SIZE)
		if err != nil {
			log.Printf("Failed to execute select statement %s", err)
			return fixedPrices, err
		}
	} else if specialties == nil && min != nil && max == nil && fromDate == nil && toDate == nil {
		stmt, err := db.Prepare("SELECT a.stripeId, f.id, f.tradespersonId, f.priceId, f.title, f.price, f.subscription, f.subInterval, s.vanityURL FROM tradesperson_account a INNER JOIN tradesperson_settings s ON a.tradespersonId=s.tradespersonId INNER JOIN fixed_prices f ON a.tradespersonId=f.tradespersonId LEFT JOIN fixed_price_state_cities c ON c.fixedPriceId=f.id WHERE f.category=? AND f.subcategory=? AND ((? = '' AND ? = '') OR f.selectPlaces=false OR c.state=? OR JSON_CONTAINS(c.cities, JSON_OBJECT('name', ?))) AND f.archived=false AND f.price >= ? GROUP BY f.id ORDER BY " + sort + " LIMIT ?, ?")
		if err != nil {
			log.Printf("Failed to create select statement %s", err)
			return fixedPrices, err
		}
		defer stmt.Close()

		offset := (page - 1) * int64(PAGE_SIZE)
		rows, err = stmt.Query(category, subCategory, state, city, state, city, *min, offset, PAGE_SIZE)
		if err != nil {
			log.Printf("Failed to execute select statement %s", err)
			return fixedPrices, err
		}
	} else if specialties == nil && min != nil && max != nil && fromDate == nil && toDate == nil {
		stmt, err := db.Prepare("SELECT a.stripeId, f.id, f.tradespersonId, f.priceId, f.title, f.price, f.subscription, f.subInterval, s.vanityURL FROM tradesperson_account a INNER JOIN tradesperson_settings s ON a.tradespersonId=s.tradespersonId INNER JOIN fixed_prices f ON a.tradespersonId=f.tradespersonId LEFT JOIN fixed_price_state_cities c ON c.fixedPriceId=f.id WHERE f.category=? AND f.subcategory=? AND ((? = '' AND ? = '') OR f.selectPlaces=false OR c.state=? OR JSON_CONTAINS(c.cities, JSON_OBJECT('name', ?))) AND f.archived=false AND f.price >= ? AND f.price <= ? GROUP BY f.id ORDER BY " + sort + " LIMIT ?, ?")
		if err != nil {
			log.Printf("Failed to create select statement %s", err)
			return fixedPrices, err
		}
		defer stmt.Close()

		offset := (page - 1) * int64(PAGE_SIZE)
		rows, err = stmt.Query(category, subCategory, state, city, state, city, *min, *max, offset, PAGE_SIZE)
		if err != nil {
			log.Printf("Failed to execute select statement %s", err)
			return fixedPrices, err
		}
	} else if specialties == nil && min != nil && max != nil && fromDate != nil && toDate == nil {
		stmt, err := db.Prepare("SELECT a.stripeId, f.id, f.tradespersonId, f.priceId, f.title, f.price, f.subscription, f.subInterval, s.vanityURL FROM tradesperson_account a INNER JOIN tradesperson_settings s ON a.tradespersonId=s.tradespersonId INNER JOIN fixed_prices f ON a.tradespersonId=f.tradespersonId LEFT JOIN fixed_price_state_cities c ON c.fixedPriceId=f.id WHERE f.category=? AND f.subcategory=? AND ((? = '' AND ? = '') OR f.selectPlaces=false OR c.state=? OR JSON_CONTAINS(c.cities, JSON_OBJECT('name', ?))) AND f.archived=false AND f.price >= ? AND f.price <= ? AND f.id=(SELECT fixedPriceId from fixed_price_time_slots WHERE startTime >= DATE(?) AND fixedPriceId=f.id LIMIT 1) GROUP BY f.id ORDER BY " + sort + " LIMIT ?, ?")
		if err != nil {
			log.Printf("Failed to create select statement %s", err)
			return fixedPrices, err
		}
		defer stmt.Close()

		offset := (page - 1) * int64(PAGE_SIZE)
		rows, err = stmt.Query(category, subCategory, state, city, state, city, *min, *max, *fromDate, offset, PAGE_SIZE)
		if err != nil {
			log.Printf("Failed to execute select statement %s", err)
			return fixedPrices, err
		}
	} else if specialties == nil && min != nil && max != nil && fromDate != nil && toDate != nil {
		stmt, err := db.Prepare("SELECT a.stripeId, f.id, f.tradespersonId, f.priceId, f.title, f.price, f.subscription, f.subInterval, s.vanityURL FROM tradesperson_account a INNER JOIN tradesperson_settings s ON a.tradespersonId=s.tradespersonId INNER JOIN fixed_prices f ON a.tradespersonId=f.tradespersonId LEFT JOIN fixed_price_state_cities c ON c.fixedPriceId=f.id WHERE f.category=? AND f.subcategory=? AND ((? = '' AND ? = '') OR f.selectPlaces=false OR c.state=? OR JSON_CONTAINS(c.cities, JSON_OBJECT('name', ?))) AND f.archived=false AND f.price >= ? AND f.price <= ? AND f.id=(SELECT fixedPriceId from fixed_price_time_slots WHERE startTime >= DATE(?) AND startTime <= DATE(?) AND fixedPriceId=f.id LIMIT 1) GROUP BY f.id ORDER BY " + sort + " LIMIT ?, ?")
		if err != nil {
			log.Printf("Failed to create select statement %s", err)
			return fixedPrices, err
		}
		defer stmt.Close()

		offset := (page - 1) * int64(PAGE_SIZE)
		rows, err = stmt.Query(category, subCategory, state, city, state, city, *min, *max, *fromDate, *toDate, offset, PAGE_SIZE)
		if err != nil {
			log.Printf("Failed to execute select statement %s", err)
			return fixedPrices, err
		}
	} else if specialties == nil && min == nil && max != nil && fromDate == nil && toDate == nil {
		stmt, err := db.Prepare("SELECT a.stripeId, f.id, f.tradespersonId, f.priceId, f.title, f.price, f.subscription, f.subInterval, s.vanityURL FROM tradesperson_account a INNER JOIN tradesperson_settings s ON a.tradespersonId=s.tradespersonId INNER JOIN fixed_prices f ON a.tradespersonId=f.tradespersonId LEFT JOIN fixed_price_state_cities c ON c.fixedPriceId=f.id WHERE f.category=? AND f.subcategory=? AND ((? = '' AND ? = '') OR f.selectPlaces=false OR c.state=? OR JSON_CONTAINS(c.cities, JSON_OBJECT('name', ?))) AND f.archived=false AND f.price <= ? GROUP BY f.id ORDER BY " + sort + " LIMIT ?, ?")
		if err != nil {
			log.Printf("Failed to create select statement %s", err)
			return fixedPrices, err
		}
		defer stmt.Close()

		offset := (page - 1) * int64(PAGE_SIZE)
		rows, err = stmt.Query(category, subCategory, state, city, state, city, *max, offset, PAGE_SIZE)
		if err != nil {
			log.Printf("Failed to execute select statement %s", err)
			return fixedPrices, err
		}
	} else if specialties == nil && min == nil && max != nil && fromDate != nil && toDate == nil {
		stmt, err := db.Prepare("SELECT a.stripeId, f.id, f.tradespersonId, f.priceId, f.title, f.price, f.subscription, f.subInterval, s.vanityURL FROM tradesperson_account a INNER JOIN tradesperson_settings s ON a.tradespersonId=s.tradespersonId INNER JOIN fixed_prices f ON a.tradespersonId=f.tradespersonId LEFT JOIN fixed_price_state_cities c ON c.fixedPriceId=f.id WHERE f.category=? AND f.subcategory=? AND ((? = '' AND ? = '') OR f.selectPlaces=false OR c.state=? OR JSON_CONTAINS(c.cities, JSON_OBJECT('name', ?))) AND f.archived=false AND f.price <= ? AND f.id=(SELECT fixedPriceId from fixed_price_time_slots WHERE startTime >= DATE(?) AND fixedPriceId=f.id LIMIT 1) GROUP BY f.id ORDER BY " + sort + " LIMIT ?, ?")
		if err != nil {
			log.Printf("Failed to create select statement %s", err)
			return fixedPrices, err
		}
		defer stmt.Close()

		offset := (page - 1) * int64(PAGE_SIZE)
		rows, err = stmt.Query(category, subCategory, state, city, state, city, *max, *fromDate, offset, PAGE_SIZE)
		if err != nil {
			log.Printf("Failed to execute select statement %s", err)
			return fixedPrices, err
		}
	} else if specialties == nil && min == nil && max != nil && fromDate != nil && toDate != nil {
		stmt, err := db.Prepare("SELECT a.stripeId, f.id, f.tradespersonId, f.priceId, f.title, f.price, f.subscription, f.subInterval, s.vanityURL FROM tradesperson_account a INNER JOIN tradesperson_settings s ON a.tradespersonId=s.tradespersonId INNER JOIN fixed_prices f ON a.tradespersonId=f.tradespersonId LEFT JOIN fixed_price_state_cities c ON c.fixedPriceId=f.id WHERE f.category=? AND f.subcategory=? AND ((? = '' AND ? = '') OR f.selectPlaces=false OR c.state=? OR JSON_CONTAINS(c.cities, JSON_OBJECT('name', ?))) AND f.archived=false AND f.price <= ? AND f.id=(SELECT fixedPriceId from fixed_price_time_slots WHERE startTime >= DATE(?) AND startTime <= DATE(?) AND fixedPriceId=f.id LIMIT 1) GROUP BY f.id ORDER BY " + sort + " LIMIT ?, ?")
		if err != nil {
			log.Printf("Failed to create select statement %s", err)
			return fixedPrices, err
		}
		defer stmt.Close()

		offset := (page - 1) * int64(PAGE_SIZE)
		rows, err = stmt.Query(category, subCategory, state, city, state, city, *max, *fromDate, *toDate, offset, PAGE_SIZE)
		if err != nil {
			log.Printf("Failed to execute select statement %s", err)
			return fixedPrices, err
		}
	} else if specialties == nil && min == nil && max == nil && fromDate != nil && toDate == nil {
		stmt, err := db.Prepare("SELECT a.stripeId, f.id, f.tradespersonId, f.priceId, f.title, f.price, f.subscription, f.subInterval, s.vanityURL FROM tradesperson_account a INNER JOIN tradesperson_settings s ON a.tradespersonId=s.tradespersonId INNER JOIN fixed_prices f ON a.tradespersonId=f.tradespersonId LEFT JOIN fixed_price_state_cities c ON c.fixedPriceId=f.id WHERE f.category=? AND f.subcategory=? AND ((? = '' AND ? = '') OR f.selectPlaces=false OR c.state=? OR JSON_CONTAINS(c.cities, JSON_OBJECT('name', ?))) AND f.archived=false AND f.id=(SELECT fixedPriceId from fixed_price_time_slots WHERE startTime >= DATE(?) AND fixedPriceId=f.id LIMIT 1) GROUP BY f.id ORDER BY " + sort + " LIMIT ?, ?")
		if err != nil {
			log.Printf("Failed to create select statement %s", err)
			return fixedPrices, err
		}
		defer stmt.Close()

		offset := (page - 1) * int64(PAGE_SIZE)
		rows, err = stmt.Query(category, subCategory, state, city, state, city, *fromDate, offset, PAGE_SIZE)
		if err != nil {
			log.Printf("Failed to execute select statement %s", err)
			return fixedPrices, err
		}
	} else if specialties == nil && min == nil && max == nil && fromDate != nil && toDate != nil {
		stmt, err := db.Prepare("SELECT a.stripeId, f.id, f.tradespersonId, f.priceId, f.title, f.price, f.subscription, f.subInterval, s.vanityURL FROM tradesperson_account a INNER JOIN tradesperson_settings s ON a.tradespersonId=s.tradespersonId INNER JOIN fixed_prices f ON a.tradespersonId=f.tradespersonId LEFT JOIN fixed_price_state_cities c ON c.fixedPriceId=f.id WHERE f.category=? AND f.subcategory=? AND ((? = '' AND ? = '') OR f.selectPlaces=false OR c.state=? OR JSON_CONTAINS(c.cities, JSON_OBJECT('name', ?))) AND f.archived=false AND f.id=(SELECT fixedPriceId from fixed_price_time_slots WHERE startTime >= DATE(?) AND startTime <= DATE(?) AND fixedPriceId=f.id LIMIT 1) GROUP BY f.id ORDER BY " + sort + " LIMIT ?, ?")
		if err != nil {
			log.Printf("Failed to create select statement %s", err)
			return fixedPrices, err
		}
		defer stmt.Close()

		offset := (page - 1) * int64(PAGE_SIZE)
		rows, err = stmt.Query(category, subCategory, state, city, state, city, *fromDate, *toDate, offset, PAGE_SIZE)
		if err != nil {
			log.Printf("Failed to execute select statement %s", err)
			return fixedPrices, err
		}
	}

	fixedPrices, err := processFixedPriceRows(db, rows, fixedPrices)
	if err != nil {
		log.Println("Failed to process rows")
		return fixedPrices, err
	}

	return fixedPrices, nil
}

func getSubCategoryFixedPrices(state, city, category, subCategory, sort string, page int64) ([]*models.Service, error) {
	fixedPrices := []*models.Service{}

	db := database.GetConnection()

	stmt, err := db.Prepare("SELECT a.stripeId, f.id, f.tradespersonId, f.priceId, f.title, f.price, f.subscription, f.subInterval, s.vanityURL FROM tradesperson_account a INNER JOIN tradesperson_settings s ON a.tradespersonId=s.tradespersonId INNER JOIN fixed_prices f ON a.tradespersonId=f.tradespersonId LEFT JOIN fixed_price_state_cities c ON c.fixedPriceId=f.id WHERE f.category=? AND f.subcategory=? AND ((? = '' AND ? = '') OR f.selectPlaces=false OR c.state=? OR JSON_CONTAINS(c.cities, JSON_OBJECT('name', ?))) AND f.archived=false GROUP BY f.id ORDER BY " + sort + " LIMIT ?, ?")
	if err != nil {
		log.Printf("Failed to create select statement %s", err)
		return fixedPrices, err
	}
	defer stmt.Close()

	offset := (page - 1) * int64(PAGE_SIZE)
	rows, err := stmt.Query(category, subCategory, state, city, state, city, offset, PAGE_SIZE)
	if err != nil {
		log.Printf("Failed to execute select statement %s", err)
		return fixedPrices, err
	}

	fixedPrices, err = processFixedPriceRows(db, rows, fixedPrices)
	if err != nil {
		log.Println("Failed to process rows")
		return fixedPrices, err
	}

	return fixedPrices, nil
}

func getCategoryFixedPrices(state, city, category, sort string, page int64) ([]*models.Service, error) {
	fixedPrices := []*models.Service{}

	db := database.GetConnection()

	stmt, err := db.Prepare("SELECT a.stripeId, f.id, f.tradespersonId, f.priceId, f.title, f.price, f.subscription, f.subInterval, s.vanityURL FROM tradesperson_account a INNER JOIN tradesperson_settings s ON a.tradespersonId=s.tradespersonId INNER JOIN fixed_prices f ON a.tradespersonId=f.tradespersonId LEFT JOIN fixed_price_state_cities c ON c.fixedPriceId=f.id WHERE f.category=? AND ((? = '' AND ? = '') OR f.selectPlaces=false OR c.state=? OR JSON_CONTAINS(c.cities, JSON_OBJECT('name', ?))) AND f.archived=false GROUP BY f.id ORDER BY " + sort + " LIMIT ?, ?")
	if err != nil {
		log.Printf("Failed to create select statement %s", err)
		return fixedPrices, err
	}
	defer stmt.Close()

	offset := (page - 1) * int64(PAGE_SIZE)
	rows, err := stmt.Query(category, state, city, state, city, offset, PAGE_SIZE)
	if err != nil {
		log.Printf("Failed to execute select statement %s", err)
		return fixedPrices, err
	}

	fixedPrices, err = processFixedPriceRows(db, rows, fixedPrices)
	if err != nil {
		log.Println("Failed to process rows")
		return fixedPrices, err
	}

	return fixedPrices, nil
}

func getAllFixedPrices(state, city, sort string, page int64) ([]*models.Service, error) {
	fixedPrices := []*models.Service{}

	db := database.GetConnection()

	stmt, err := db.Prepare("SELECT a.stripeId, f.id, f.tradespersonId, f.priceId, f.title, f.price, f.subscription, f.subInterval, s.vanityURL FROM tradesperson_account a INNER JOIN tradesperson_settings s ON a.tradespersonId=s.tradespersonId INNER JOIN fixed_prices f ON a.tradespersonId=f.tradespersonId LEFT JOIN fixed_price_state_cities c ON c.fixedPriceId=f.id WHERE (((SELECT COUNT(*) FROM fixed_price_time_slots WHERE fixedPriceId=f.id AND startTime > CURDATE()) > 0  OR f.subscription=True) OR f.subscription=True) AND ((? = '' AND ? = '') OR f.selectPlaces=false OR c.state=? OR JSON_CONTAINS(c.cities, JSON_OBJECT('name', ?))) AND f.archived=false GROUP BY f.id ORDER BY " + sort + " LIMIT ?, ?")
	if err != nil {
		return fixedPrices, err
	}
	defer stmt.Close()

	offSet := (page - 1) * int64(PAGE_SIZE)
	rows, err := stmt.Query(state, city, state, city, offSet, PAGE_SIZE)
	if err != nil {
		log.Println("Failed to execute select statement")
		return fixedPrices, err
	}

	fixedPrices, err = processFixedPriceRows(db, rows, fixedPrices)
	if err != nil {
		log.Println("Failed to process rows")
		return fixedPrices, err
	}

	return fixedPrices, nil
}

func GetFixedPricesHandler(params operations.GetFixedPricesParams) middleware.Responder {
	state := ""
	city := ""
	sort := "f.id DESC"
	if params.State != nil {
		state = *params.State
	}
	if params.City != nil {
		city = *params.City
	}
	if params.Sort != nil {
		if *params.Sort == "costasc" {
			sort = "f.price ASC"
		} else if *params.Sort == "costdesc" {
			sort = "f.price DESC"
		}
	}
	page := *params.Page
	response := operations.NewGetFixedPricesOK()
	fixedPrices := []*models.Service{}
	response.SetPayload(fixedPrices)

	var err error
	if params.Category == nil && params.SubCategory == nil {
		fixedPrices, err = getAllFixedPrices(state, city, sort, page)
		if err != nil {
			log.Printf("%s", err)
			return response
		}
	} else if params.Category != nil && params.SubCategory == nil {
		fixedPrices, err = getCategoryFixedPrices(state, city, *params.Category, sort, page)
		if err != nil {
			log.Printf("%s", err)
			return response
		}
	} else if params.Category != nil && params.SubCategory != nil {
		if params.Specialties == nil && params.Min == nil && params.Max == nil && params.FromDate == nil && params.ToDate == nil {
			fixedPrices, err = getSubCategoryFixedPrices(state, city, *params.Category, *params.SubCategory, sort, page)
			if err != nil {
				log.Printf("%s", err)
				return response
			}
		} else {
			fixedPrices, err = getFixedPricesWithFilters(state, city, *params.Category, *params.SubCategory, sort, params.Specialties, params.FromDate, params.ToDate, page, params.Min, params.Max)
			if err != nil {
				log.Printf("%s", err)
				return response
			}
		}
	}

	response.SetPayload(fixedPrices)

	return response
}

//PAGES
func getFixedPricesWithSpecialtiesPages(state, city, category, subCategory string, specialties, fromDate, toDate *string, min, max *int64) (float64, error) {
	query := ""
	db := database.GetConnection()
	var pages float64

	if specialties != nil {
		specialtyArry := strings.Split(*specialties, ",")
		for _, specialty := range specialtyArry {
			query += "'" + specialty + "',"
		}
		query = query[:len(query)-1]
		query = "AND fi.specialty IN (" + query + ")"
	}

	if specialties != nil && min == nil && max == nil && fromDate == nil && toDate == nil {
		stmt, err := db.Prepare("SELECT COUNT(*) FROM fixed_prices f LEFT JOIN fixed_price_state_cities c ON c.fixedPriceId=f.id LEFT JOIN fixed_price_specialties fi ON fi.fixedPriceId=f.id WHERE f.category=? AND f.subcategory=? AND ((? = '' AND ? = '') OR f.selectPlaces=false OR c.state=? OR JSON_CONTAINS(c.cities, JSON_OBJECT('name', ?))) " + query + " AND f.archived=false")
		if err != nil {
			log.Printf("Failed to create select statement %s", err)
			return pages, err
		}
		defer stmt.Close()

		err = stmt.QueryRow(category, subCategory, state, city, state, city).Scan(&pages)
		if err != nil {
			log.Printf("Failed to execute select statement %s", err)
			return pages, err
		}
	} else if specialties != nil && min != nil && max == nil && fromDate == nil && toDate == nil {
		stmt, err := db.Prepare("SELECT COUNT(*) FROM fixed_prices f LEFT JOIN fixed_price_state_cities c ON c.fixedPriceId=f.id LEFT JOIN fixed_price_specialties fi ON fi.fixedPriceId=f.id WHERE f.category=? AND f.subcategory=? AND ((? = '' AND ? = '') OR f.selectPlaces=false OR c.state=? OR JSON_CONTAINS(c.cities, JSON_OBJECT('name', ?))) " + query + " AND f.archived=false AND f.price >= ?")
		if err != nil {
			log.Printf("Failed to create select statement %s", err)
			return pages, err
		}
		defer stmt.Close()

		err = stmt.QueryRow(category, subCategory, state, city, state, city, *min).Scan(&pages)
		if err != nil {
			log.Printf("Failed to execute select statement %s", err)
			return pages, err
		}
	} else if specialties != nil && min != nil && max != nil && fromDate == nil && toDate == nil {
		stmt, err := db.Prepare("SELECT COUNT(*) FROM fixed_prices f LEFT JOIN fixed_price_state_cities c ON c.fixedPriceId=f.id LEFT JOIN fixed_price_specialties fi ON fi.fixedPriceId=f.id WHERE f.category=? AND f.subcategory=? AND ((? = '' AND ? = '') OR f.selectPlaces=false OR c.state=? OR JSON_CONTAINS(c.cities, JSON_OBJECT('name', ?))) " + query + " AND f.archived=false AND f.price >= ? AND f.price <= ?")
		if err != nil {
			log.Printf("Failed to create select statement %s", err)
			return pages, err
		}
		defer stmt.Close()

		err = stmt.QueryRow(category, subCategory, state, city, state, city, *min, *max).Scan(&pages)
		if err != nil {
			log.Printf("Failed to execute select statement %s", err)
			return pages, err
		}
	} else if specialties != nil && min != nil && max != nil && fromDate != nil && toDate == nil {
		stmt, err := db.Prepare("SELECT COUNT(*) FROM fixed_prices f LEFT JOIN fixed_price_state_cities c ON c.fixedPriceId=f.id LEFT JOIN fixed_price_specialties fi ON fi.fixedPriceId=f.id WHERE f.category=? AND f.subcategory=? AND ((? = '' AND ? = '') OR f.selectPlaces=false OR c.state=? OR JSON_CONTAINS(c.cities, JSON_OBJECT('name', ?))) " + query + " AND f.archived=false AND f.price >= ? AND f.price <= ? AND f.id=(SELECT fixedPriceId from fixed_price_time_slots WHERE startTime >= DATE(?) AND fixedPriceId=f.id LIMIT 1)")
		if err != nil {
			log.Printf("Failed to create select statement %s", err)
			return pages, err
		}
		defer stmt.Close()

		err = stmt.QueryRow(category, subCategory, state, city, state, city, *min, *max, *fromDate).Scan(&pages)
		if err != nil {
			log.Printf("Failed to execute select statement %s", err)
			return pages, err
		}
	} else if specialties != nil && min != nil && max != nil && fromDate != nil && toDate != nil {
		stmt, err := db.Prepare("SELECT COUNT(*) FROM fixed_prices f LEFT JOIN fixed_price_state_cities c ON c.fixedPriceId=f.id LEFT JOIN fixed_price_specialties fi ON fi.fixedPriceId=f.id WHERE f.category=? AND f.subcategory=? AND ((? = '' AND ? = '') OR f.selectPlaces=false OR c.state=? OR JSON_CONTAINS(c.cities, JSON_OBJECT('name', ?))) " + query + " AND f.archived=false AND f.price >= ? AND f.price <= ? AND f.id=(SELECT fixedPriceId from fixed_price_time_slots WHERE startTime >= DATE(?) AND startTime <= DATE(?) AND fixedPriceId=f.id LIMIT 1)")
		if err != nil {
			log.Printf("Failed to create select statement %s", err)
			return pages, err
		}
		defer stmt.Close()

		err = stmt.QueryRow(category, subCategory, state, city, state, city, *min, *max, *fromDate, *toDate).Scan(&pages)
		if err != nil {
			log.Printf("Failed to execute select statement %s", err)
			return pages, err
		}
	} else if specialties != nil && min == nil && max != nil && fromDate == nil && toDate == nil {
		stmt, err := db.Prepare("SELECT COUNT(*) FROM fixed_prices f LEFT JOIN fixed_price_state_cities c ON c.fixedPriceId=f.id LEFT JOIN fixed_price_specialties fi ON fi.fixedPriceId=f.id WHERE f.category=? AND f.subcategory=? AND ((? = '' AND ? = '') OR f.selectPlaces=false OR c.state=? OR JSON_CONTAINS(c.cities, JSON_OBJECT('name', ?))) " + query + " AND f.archived=false AND f.price <= ?")
		if err != nil {
			log.Printf("Failed to create select statement %s", err)
			return pages, err
		}
		defer stmt.Close()

		err = stmt.QueryRow(category, subCategory, state, city, state, city, *max).Scan(&pages)
		if err != nil {
			log.Printf("Failed to execute select statement %s", err)
			return pages, err
		}
	} else if specialties != nil && min == nil && max != nil && fromDate != nil && toDate == nil {
		stmt, err := db.Prepare("SELECT COUNT(*) FROM fixed_prices f LEFT JOIN fixed_price_state_cities c ON c.fixedPriceId=f.id LEFT JOIN fixed_price_specialties fi ON fi.fixedPriceId=f.id WHERE f.category=? AND f.subcategory=? AND ((? = '' AND ? = '') OR f.selectPlaces=false OR c.state=? OR JSON_CONTAINS(c.cities, JSON_OBJECT('name', ?))) " + query + " AND f.archived=false AND f.price <= ? AND f.id=(SELECT fixedPriceId from fixed_price_time_slots WHERE startTime >= DATE(?) AND fixedPriceId=f.id LIMIT 1)")
		if err != nil {
			log.Printf("Failed to create select statement %s", err)
			return pages, err
		}
		defer stmt.Close()

		err = stmt.QueryRow(category, subCategory, state, city, state, city, *max, *fromDate).Scan(&pages)
		if err != nil {
			log.Printf("Failed to execute select statement %s", err)
			return pages, err
		}
	} else if specialties != nil && min == nil && max != nil && fromDate != nil && toDate != nil {
		stmt, err := db.Prepare("SELECT COUNT(*) FROM fixed_prices f LEFT JOIN fixed_price_state_cities c ON c.fixedPriceId=f.id LEFT JOIN fixed_price_specialties fi ON fi.fixedPriceId=f.id WHERE f.category=? AND f.subcategory=? AND ((? = '' AND ? = '') OR f.selectPlaces=false OR c.state=? OR JSON_CONTAINS(c.cities, JSON_OBJECT('name', ?))) " + query + " AND f.archived=false AND f.price <= ? AND f.id=(SELECT fixedPriceId from fixed_price_time_slots WHERE startTime >= DATE(?) AND startTime <= DATE(?) AND fixedPriceId=f.id LIMIT 1)")
		if err != nil {
			log.Printf("Failed to create select statement %s", err)
			return pages, err
		}
		defer stmt.Close()

		err = stmt.QueryRow(category, subCategory, state, city, state, city, *max, *fromDate, *toDate).Scan(&pages)
		if err != nil {
			log.Printf("Failed to execute select statement %s", err)
			return pages, err
		}
	} else if specialties != nil && min == nil && max == nil && fromDate != nil && toDate == nil {
		stmt, err := db.Prepare("SELECT COUNT(*) FROM fixed_prices f LEFT JOIN fixed_price_state_cities c ON c.fixedPriceId=f.id LEFT JOIN fixed_price_specialties fi ON fi.fixedPriceId=f.id WHERE f.category=? AND f.subcategory=? AND ((? = '' AND ? = '') OR f.selectPlaces=false OR c.state=? OR JSON_CONTAINS(c.cities, JSON_OBJECT('name', ?))) " + query + " AND f.archived=false AND f.id=(SELECT fixedPriceId from fixed_price_time_slots WHERE startTime >= DATE(?) AND fixedPriceId=f.id LIMIT 1)")
		if err != nil {
			log.Printf("Failed to create select statement %s", err)
			return pages, err
		}
		defer stmt.Close()

		err = stmt.QueryRow(category, subCategory, state, city, state, city, *fromDate).Scan(&pages)
		if err != nil {
			log.Printf("Failed to execute select statement %s", err)
			return pages, err
		}
	} else if specialties != nil && min == nil && max == nil && fromDate != nil && toDate != nil {
		stmt, err := db.Prepare("SELECT COUNT(*) FROM fixed_prices f LEFT JOIN fixed_price_state_cities c ON c.fixedPriceId=f.id LEFT JOIN fixed_price_specialties fi ON fi.fixedPriceId=f.id WHERE f.category=? AND f.subcategory=? AND ((? = '' AND ? = '') OR f.selectPlaces=false OR c.state=? OR JSON_CONTAINS(c.cities, JSON_OBJECT('name', ?))) " + query + " AND f.archived=false AND f.id=(SELECT fixedPriceId from fixed_price_time_slots WHERE startTime >= DATE(?) AND startTime <= DATE(?) AND fixedPriceId=f.id LIMIT 1)")
		if err != nil {
			log.Printf("Failed to create select statement %s", err)
			return pages, err
		}
		defer stmt.Close()

		err = stmt.QueryRow(category, subCategory, state, city, state, city, *fromDate, *toDate).Scan(&pages)
		if err != nil {
			log.Printf("Failed to execute select statement %s", err)
			return pages, err
		}
	} else if specialties == nil && min != nil && max == nil && fromDate == nil && toDate == nil {
		stmt, err := db.Prepare("SELECT COUNT(*) FROM fixed_prices f LEFT JOIN fixed_price_state_cities c ON c.fixedPriceId=f.id WHERE f.category=? AND f.subcategory=? AND ((? = '' AND ? = '') OR f.selectPlaces=false OR c.state=? OR JSON_CONTAINS(c.cities, JSON_OBJECT('name', ?))) AND f.archived=false AND f.price >= ?")
		if err != nil {
			log.Printf("Failed to create select statement %s", err)
			return pages, err
		}
		defer stmt.Close()

		err = stmt.QueryRow(category, subCategory, state, city, state, city, *min).Scan(&pages)
		if err != nil {
			log.Printf("Failed to execute select statement %s", err)
			return pages, err
		}
	} else if specialties == nil && min != nil && max != nil && fromDate == nil && toDate == nil {
		stmt, err := db.Prepare("SELECT COUNT(*) FROM fixed_prices f LEFT JOIN fixed_price_state_cities c ON c.fixedPriceId=f.id WHERE f.category=? AND f.subcategory=? AND ((? = '' AND ? = '') OR f.selectPlaces=false OR c.state=? OR JSON_CONTAINS(c.cities, JSON_OBJECT('name', ?))) AND f.archived=false AND f.price >= ? AND f.price <= ?")
		if err != nil {
			log.Printf("Failed to create select statement %s", err)
			return pages, err
		}
		defer stmt.Close()

		err = stmt.QueryRow(category, subCategory, state, city, state, city, *min, *max).Scan(&pages)
		if err != nil {
			log.Printf("Failed to execute select statement %s", err)
			return pages, err
		}
	} else if specialties == nil && min != nil && max != nil && fromDate != nil && toDate == nil {
		stmt, err := db.Prepare("SELECT COUNT(*) FROM fixed_prices f LEFT JOIN fixed_price_state_cities c ON c.fixedPriceId=f.id WHERE f.category=? AND f.subcategory=? AND ((? = '' AND ? = '') OR f.selectPlaces=false OR c.state=? OR JSON_CONTAINS(c.cities, JSON_OBJECT('name', ?))) AND f.archived=false AND f.price >= ? AND f.price <= ? AND f.id=(SELECT fixedPriceId from fixed_price_time_slots WHERE startTime >= DATE(?) AND fixedPriceId=f.id LIMIT 1)")
		if err != nil {
			log.Printf("Failed to create select statement %s", err)
			return pages, err
		}
		defer stmt.Close()

		err = stmt.QueryRow(category, subCategory, state, city, state, city, *min, *max, *fromDate).Scan(&pages)
		if err != nil {
			log.Printf("Failed to execute select statement %s", err)
			return pages, err
		}
	} else if specialties == nil && min != nil && max != nil && fromDate != nil && toDate != nil {
		stmt, err := db.Prepare("SELECT COUNT(*) FROM fixed_prices f LEFT JOIN fixed_price_state_cities c ON c.fixedPriceId=f.id WHERE f.category=? AND f.subcategory=? AND ((? = '' AND ? = '') OR f.selectPlaces=false OR c.state=? OR JSON_CONTAINS(c.cities, JSON_OBJECT('name', ?))) AND f.archived=false AND f.price >= ? AND f.price <= ? AND f.id=(SELECT fixedPriceId from fixed_price_time_slots WHERE startTime >= DATE(?) AND startTime <= DATE(?) AND fixedPriceId=f.id LIMIT 1)")
		if err != nil {
			log.Printf("Failed to create select statement %s", err)
			return pages, err
		}
		defer stmt.Close()

		err = stmt.QueryRow(category, subCategory, state, city, state, city, *min, *max, *fromDate, *toDate).Scan(&pages)
		if err != nil {
			log.Printf("Failed to execute select statement %s", err)
			return pages, err
		}
	} else if specialties == nil && min == nil && max != nil && fromDate == nil && toDate == nil {
		stmt, err := db.Prepare("SELECT COUNT(*) FROM fixed_prices f LEFT JOIN fixed_price_state_cities c ON c.fixedPriceId=f.id WHERE f.category=? AND f.subcategory=? AND ((? = '' AND ? = '') OR f.selectPlaces=false OR c.state=? OR JSON_CONTAINS(c.cities, JSON_OBJECT('name', ?))) AND f.archived=false AND f.price <= ?")
		if err != nil {
			log.Printf("Failed to create select statement %s", err)
			return pages, err
		}
		defer stmt.Close()

		err = stmt.QueryRow(category, subCategory, state, city, state, city, *max).Scan(&pages)
		if err != nil {
			log.Printf("Failed to execute select statement %s", err)
			return pages, err
		}
	} else if specialties == nil && min == nil && max != nil && fromDate != nil && toDate == nil {
		stmt, err := db.Prepare("SELECT COUNT(*) FROM fixed_prices f LEFT JOIN fixed_price_state_cities c ON c.fixedPriceId=f.id WHERE f.category=? AND f.subcategory=? AND ((? = '' AND ? = '') OR f.selectPlaces=false OR c.state=? OR JSON_CONTAINS(c.cities, JSON_OBJECT('name', ?))) AND f.archived=false AND f.price <= ? AND f.id=(SELECT fixedPriceId from fixed_price_time_slots WHERE startTime >= DATE(?) AND fixedPriceId=f.id LIMIT 1)")
		if err != nil {
			log.Printf("Failed to create select statement %s", err)
			return pages, err
		}
		defer stmt.Close()

		err = stmt.QueryRow(category, subCategory, state, city, state, city, *max, *fromDate).Scan(&pages)
		if err != nil {
			log.Printf("Failed to execute select statement %s", err)
			return pages, err
		}
	} else if specialties == nil && min == nil && max != nil && fromDate != nil && toDate != nil {
		stmt, err := db.Prepare("SELECT COUNT(*) FROM fixed_prices f LEFT JOIN fixed_price_state_cities c ON c.fixedPriceId=f.id WHERE f.category=? AND f.subcategory=? AND ((? = '' AND ? = '') OR f.selectPlaces=false OR c.state=? OR JSON_CONTAINS(c.cities, JSON_OBJECT('name', ?))) AND f.archived=false AND f.price <= ? AND f.id=(SELECT fixedPriceId from fixed_price_time_slots WHERE startTime >= DATE(?) AND startTime <= DATE(?) AND fixedPriceId=f.id LIMIT 1)")
		if err != nil {
			log.Printf("Failed to create select statement %s", err)
			return pages, err
		}
		defer stmt.Close()

		err = stmt.QueryRow(category, subCategory, state, city, state, city, *max, *fromDate, *toDate).Scan(&pages)
		if err != nil {
			log.Printf("Failed to execute select statement %s", err)
			return pages, err
		}
	} else if specialties == nil && min == nil && max == nil && fromDate != nil && toDate == nil {
		stmt, err := db.Prepare("SELECT COUNT(*) FROM fixed_prices f LEFT JOIN fixed_price_state_cities c ON c.fixedPriceId=f.id INNER JOIN fixed_price_time_slots s ON s.fixedPriceId=f.id WHERE f.category=? AND f.subcategory=? AND ((? = '' AND ? = '') OR f.selectPlaces=false OR c.state=? OR JSON_CONTAINS(c.cities, JSON_OBJECT('name', ?))) AND f.archived=false AND f.id=(SELECT fixedPriceId from fixed_price_time_slots WHERE startTime >= DATE(?) AND fixedPriceId=f.id LIMIT 1)")
		if err != nil {
			log.Printf("Failed to create select statement %s", err)
			return pages, err
		}
		defer stmt.Close()

		err = stmt.QueryRow(category, subCategory, state, city, state, city, *fromDate).Scan(&pages)
		if err != nil {
			log.Printf("Failed to execute select statement %s", err)
			return pages, err
		}
	} else if specialties == nil && min == nil && max == nil && fromDate != nil && toDate != nil {
		stmt, err := db.Prepare("SELECT COUNT(*) FROM fixed_prices f LEFT JOIN fixed_price_state_cities c ON c.fixedPriceId=f.id INNER JOIN fixed_price_time_slots s ON s.fixedPriceId=f.id WHERE f.category=? AND f.subcategory=? AND ((? = '' AND ? = '') OR f.selectPlaces=false OR c.state=? OR JSON_CONTAINS(c.cities, JSON_OBJECT('name', ?))) AND f.archived=false AND f.id=(SELECT fixedPriceId from fixed_price_time_slots WHERE startTime >= DATE(?) AND startTime <= DATE(?) AND fixedPriceId=f.id LIMIT 1)")
		if err != nil {
			log.Printf("Failed to create select statement %s", err)
			return pages, err
		}
		defer stmt.Close()

		err = stmt.QueryRow(category, subCategory, state, city, state, city, *fromDate, *toDate).Scan(&pages)
		if err != nil {
			log.Printf("Failed to execute select statement %s", err)
			return pages, err
		}
	}
	return pages, nil
}

func getSubCategoryFixedPricePages(state, city, category, subCategory string) (float64, error) {
	db := database.GetConnection()
	var pages float64

	stmt, err := db.Prepare("SELECT COUNT(*) FROM fixed_prices f LEFT JOIN fixed_price_state_cities c ON c.fixedPriceId=f.id WHERE f.category=? AND f.subcategory=? AND ((? = '' AND ? = '') OR f.selectPlaces=false OR c.state=? OR JSON_CONTAINS(c.cities, JSON_OBJECT('name', ?))) AND f.archived=false")
	if err != nil {
		log.Printf("Failed to create select statement %s", err)
		return pages, err
	}
	defer stmt.Close()

	err = stmt.QueryRow(category, subCategory, state, city, state, city).Scan(&pages)
	if err != nil {
		return pages, err
	}

	return pages, nil
}

func getCategoryFixedPricePages(state, city, category string) (float64, error) {
	db := database.GetConnection()
	var pages float64

	stmt, err := db.Prepare("SELECT COUNT(*) FROM fixed_prices f LEFT JOIN fixed_price_state_cities c ON c.fixedPriceId=f.id WHERE f.category=? AND ((? = '' AND ? = '') OR f.selectPlaces=false OR c.state=? OR JSON_CONTAINS(c.cities, JSON_OBJECT('name', ?))) AND f.archived=false")
	if err != nil {
		return pages, err
	}
	defer stmt.Close()

	err = stmt.QueryRow(category, state, city, state, city).Scan(&pages)
	if err != nil {
		return pages, err
	}

	return pages, nil
}

func getAllFixedPricePages(state, city string) (float64, error) {
	db := database.GetConnection()
	var pages float64

	stmt, err := db.Prepare("SELECT COUNT(*) FROM fixed_prices f LEFT JOIN fixed_price_state_cities c ON c.fixedPriceId=f.id WHERE ((SELECT COUNT(*) FROM fixed_price_time_slots WHERE fixedPriceId=f.id AND startTime > CURDATE()) > 0  OR f.subscription=True) AND ((? = '' AND ? = '') OR f.selectPlaces=false OR c.state=? OR JSON_CONTAINS(c.cities, JSON_OBJECT('name', ?))) AND f.archived=false")
	if err != nil {
		return pages, err
	}
	defer stmt.Close()

	err = stmt.QueryRow(state, city, state, city).Scan(&pages)
	if err != nil {
		return pages, err
	}

	return pages, nil
}

func GetFixedPricePagesHandler(params operations.GetFixedPricePagesParams) middleware.Responder {
	state := ""
	city := ""
	if params.State != nil {
		state = *params.State
	}
	if params.City != nil {
		city = *params.City
	}

	pages := float64(1)
	response := operations.NewGetFixedPricePagesOK().WithPayload(int64(pages))

	var err error
	if params.Category == nil && params.SubCategory == nil {
		pages, err = getAllFixedPricePages(state, city)
		if err != nil {
			log.Printf("%s", err)
			return response
		}
	} else if params.Category != nil && params.SubCategory == nil {
		pages, err = getCategoryFixedPricePages(state, city, *params.Category)
		if err != nil {
			log.Printf("%s", err)
			return response
		}
	} else if params.Category != nil && params.SubCategory != nil {
		if params.Specialties == nil && params.Min == nil && params.Max == nil && params.FromDate == nil && params.ToDate == nil {
			pages, err = getSubCategoryFixedPricePages(state, city, *params.Category, *params.SubCategory)
			if err != nil {
				log.Printf("%s", err)
				return response
			}
		} else {
			pages, err = getFixedPricesWithSpecialtiesPages(state, city, *params.Category, *params.SubCategory, params.Specialties, params.FromDate, params.ToDate, params.Min, params.Max)
			if err != nil {
				log.Printf("%s", err)
				return response
			}
		}
	}

	if pages == float64(0) {
		pages = float64(1)
	}

	pages = math.Ceil(pages / PAGE_SIZE)

	response.SetPayload(int64(pages))

	return response
}
