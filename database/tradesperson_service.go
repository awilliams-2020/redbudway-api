package database

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"redbudway-api/internal"
	"redbudway-api/models"
	"redbudway-api/stripe"
	"strconv"

	"github.com/stripe/stripe-go/v72/price"
)

const PAGE_SIZE = float64(9)

func insertFilters(fixedPriceID int64, fixedPrice *models.ServiceDetails) error {
	for _, filter := range fixedPrice.Filters {
		stmt, err := db.Prepare("INSERT INTO fixed_price_filters (fixedPriceId, filter) VALUES (?, ?)")
		if err != nil {
			return err
		}
		defer stmt.Close()
		_, err = stmt.Exec(fixedPriceID, filter)
		if err != nil {
			return err
		}
	}
	return nil
}

func insertStatesAndCities(fixedPriceID int64, fixedPrice *models.ServiceDetails) error {
	for _, obj := range fixedPrice.StatesAndCities {
		cities, err := json.Marshal(obj.Cities)
		if err != nil {
			return err
		}
		stmt, err := db.Prepare("INSERT INTO fixed_price_state_cities (fixedPriceId, state, cities) VALUES (?, ?, ?)")
		if err != nil {
			return err
		}
		defer stmt.Close()
		_, err = stmt.Exec(fixedPriceID, obj.State, string(cities))
		if err != nil {
			return err
		}
	}
	return nil
}

func CreateFixedPrice(tradespersonID string, fixedPrice *models.ServiceDetails) (bool, error) {

	price, err := stripe.CreatePrice(fixedPrice)
	if err != nil {
		return false, err
	}
	stmt, err := db.Prepare("INSERT INTO fixed_prices (tradespersonId, priceId, category, subcategory, title, price, description, subscription, subInterval, selectPlaces, archived) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, false)")
	if err != nil {
		return false, err
	}
	defer stmt.Close()

	results, err := stmt.Exec(tradespersonID, price.ID, fixedPrice.Category, fixedPrice.SubCategory, fixedPrice.Title, price.UnitAmountDecimal, fixedPrice.Description, fixedPrice.Subscription, fixedPrice.Interval, fixedPrice.SelectPlaces)
	if err != nil {
		return false, err
	}
	rowsAffected, err := results.RowsAffected()
	if err != nil {
		return false, err
	}
	fixedPriceID, err := results.LastInsertId()
	if err != nil {
		return false, err
	}
	if rowsAffected == 1 {
		images, err := internal.ProcessImages(tradespersonID, price.ID, fixedPrice)
		if err != nil {
			return false, err
		}
		if err := stripe.UpdateProduct(images, fixedPrice, price); err != nil {
			return false, err
		}

		if err := updateImages(fixedPriceID, images, "fixed_price"); err != nil {
			return false, err
		}

		if err = insertStatesAndCities(fixedPriceID, fixedPrice); err != nil {
			return false, err
		}

		if err = InsertTimeSlots(fixedPriceID, fixedPrice); err != nil {
			return false, err
		}

		if err = insertFilters(fixedPriceID, fixedPrice); err != nil {
			return false, err
		}
	}

	return true, nil
}

func GetFilters(fixedPriceID int64) ([]string, error) {
	filters := []string{}

	stmt, err := db.Prepare("SELECT filter FROM fixed_price_filters WHERE fixedPriceId=?")
	if err != nil {
		return filters, err
	}
	defer stmt.Close()

	rows, err := stmt.Query(fixedPriceID)
	if err != nil {
		return filters, err
	}

	var filter string
	for rows.Next() {
		if err := rows.Scan(&filter); err != nil {
			return filters, err
		}
		filters = append(filters, filter)
	}
	return filters, nil
}

func GetIncludes(fixedPriceID int64) ([]string, []string, error) {
	includes := []string{}
	notIncludes := []string{}

	stmt, err := db.Prepare("SELECT included, items FROM fixed_price_includes WHERE fixedPriceId=?")
	if err != nil {
		return includes, notIncludes, err
	}
	defer stmt.Close()

	rows, err := stmt.Query(fixedPriceID)
	if err != nil {
		return includes, notIncludes, err
	}

	var included bool
	var includesJSON string
	for rows.Next() {
		if err := rows.Scan(&included, &includesJSON); err != nil {
			return includes, notIncludes, err
		}
		var tempIncludes []string
		err := json.Unmarshal([]byte(includesJSON), &tempIncludes)
		if err != nil {
			return includes, notIncludes, err
		}
		if included {
			includes = tempIncludes
		} else {
			notIncludes = tempIncludes
		}
	}
	return includes, notIncludes, nil
}

func GetFixedPriceStatesAndCities(fixedPriceID int64) ([]*models.ServiceDetailsStatesAndCitiesItems0, error) {
	StatesAndCities := []*models.ServiceDetailsStatesAndCitiesItems0{}

	stmt, err := db.Prepare("SELECT state, cities FROM fixed_price_state_cities WHERE fixedPriceId=?")
	if err != nil {
		return StatesAndCities, err
	}
	defer stmt.Close()

	rows, err := stmt.Query(fixedPriceID)
	if err != nil {
		return StatesAndCities, err
	}

	var state, citiesJSON string
	for rows.Next() {
		if err := rows.Scan(&state, &citiesJSON); err != nil {
			return StatesAndCities, err
		}
		statesAndCitiesItem := &models.ServiceDetailsStatesAndCitiesItems0{}
		statesAndCitiesItem.State = state
		var cities []*models.ServiceDetailsStatesAndCitiesItems0CitiesItems0
		err := json.Unmarshal([]byte(citiesJSON), &cities)
		if err != nil {
			fmt.Println("error:", err)
		}
		statesAndCitiesItem.Cities = cities
		StatesAndCities = append(StatesAndCities, statesAndCitiesItem)
	}
	return StatesAndCities, nil
}

func GetQuoteStatesAndCities(quoteID int64) ([]*models.ServiceDetailsStatesAndCitiesItems0, error) {
	StatesAndCities := []*models.ServiceDetailsStatesAndCitiesItems0{}

	stmt, err := db.Prepare("SELECT state, cities FROM quote_state_cities WHERE quoteId=?")
	if err != nil {
		return StatesAndCities, err
	}
	defer stmt.Close()

	rows, err := stmt.Query(quoteID)
	if err != nil {
		return StatesAndCities, err
	}

	var state, citiesBlob string
	for rows.Next() {
		if err := rows.Scan(&state, &citiesBlob); err != nil {
			return StatesAndCities, err
		}
		statesAndCitiesItem := &models.ServiceDetailsStatesAndCitiesItems0{}
		statesAndCitiesItem.State = state
		var cities []*models.ServiceDetailsStatesAndCitiesItems0CitiesItems0
		err := json.Unmarshal([]byte(citiesBlob), &cities)
		if err != nil {
			fmt.Println("error:", err)
		}
		statesAndCitiesItem.Cities = cities
		StatesAndCities = append(StatesAndCities, statesAndCitiesItem)
	}
	return StatesAndCities, nil
}

func GetOtherServices(tradespersonID string, fixedPriceID int64) ([]*models.OtherServicesItems0, error) {
	otherServices := []*models.OtherServicesItems0{}

	stmt, err := db.Prepare("SELECT id, subscription, subInterval FROM fixed_prices WHERE tradespersonId=? AND id != ? AND archived = false")
	if err != nil {
		return otherServices, err
	}
	defer stmt.Close()

	rows, err := stmt.Query(tradespersonID, fixedPriceID)
	if err != nil {
		return otherServices, err
	}

	var subscription bool
	var interval sql.NullString
	for rows.Next() {
		if err := rows.Scan(&fixedPriceID, &subscription, &interval); err != nil {
			return otherServices, err
		}
		service := &models.OtherServicesItems0{}
		service.Subscription = subscription
		if interval.Valid {
			service.Interval = interval.String
		}
		stmt, err := db.Prepare("SELECT startTime, endTime FROM fixed_price_time_slots WHERE fixedPriceId=?")
		if err != nil {
			log.Printf("Failed to create prepared statement, %s", err)
			continue
		}
		defer stmt.Close()

		rows, err := stmt.Query(fixedPriceID)
		if err != nil {
			log.Printf("Failed to execute prepared statement, %s", err)
			continue
		}

		timeSlots := []*models.TimeSlot{}

		var startTime, endTime string
		for rows.Next() {
			if err := rows.Scan(&startTime, &endTime); err != nil {
				log.Printf("Failed to scan statement, %s", err)
				continue
			}
			timeSlot := &models.TimeSlot{}
			timeSlot.StartTime = startTime
			timeSlot.EndTime = endTime
			timeSlots = append(timeSlots, timeSlot)
		}
		service.TimeSlots = timeSlots
		otherServices = append(otherServices, service)
	}
	return otherServices, nil
}

func GetTradespersonFixedPrice(tradespersonID string, priceID string) (*models.ServiceDetails, int64, error) {
	fixedPrice := &models.ServiceDetails{}

	stmt, err := db.Prepare("SELECT f.id, f.category, f.subCategory, f.title, f.price, f.description, f.subscription, f.subInterval, f.selectPlaces, f.archived FROM tradesperson_account a INNER JOIN fixed_prices f ON a.tradespersonId=f.tradespersonId WHERE a.tradespersonId=? AND f.priceId=?")
	if err != nil {
		return fixedPrice, 0, err
	}
	defer stmt.Close()

	row := stmt.QueryRow(tradespersonID, priceID)
	var ID, price int64
	var interval sql.NullString
	switch err = row.Scan(&ID, &fixedPrice.Category, &fixedPrice.SubCategory, &fixedPrice.Title, &price, &fixedPrice.Description, &fixedPrice.Subscription, &interval, &fixedPrice.SelectPlaces, &fixedPrice.Archived); err {
	case sql.ErrNoRows:
		return fixedPrice, ID, err
	case nil:
		if interval.Valid {
			fixedPrice.Interval = interval.String
		}
		strPrice := fmt.Sprintf("%.2f", float64(price)/float64(100.00))
		floatPrice, err := strconv.ParseFloat(strPrice, 64)
		if err != nil {
			return fixedPrice, ID, err
		}
		fixedPrice.Price = floatPrice

		fixedPrice.Images, err = GetImages(ID, "fixed_price")
		if err != nil {
			return fixedPrice, ID, err
		}
		fixedPrice.TimeSlots, err = GetTimeSlots(ID)
		if err != nil {
			return fixedPrice, ID, err
		}
		fixedPrice.StatesAndCities, err = GetFixedPriceStatesAndCities(ID)
		if err != nil {
			return fixedPrice, ID, err
		}
		fixedPrice.Filters, err = GetFilters(ID)
		if err != nil {
			return fixedPrice, ID, err
		}
		fixedPrice.Includes, fixedPrice.NotIncludes, err = GetIncludes(ID)
		if err != nil {
			return fixedPrice, ID, err
		}
	default:
		return fixedPrice, ID, err
	}

	return fixedPrice, ID, nil
}

func updateFilters(fixedPriceID int64, fixedPrice *models.ServiceDetails) error {
	stmt, err := db.Prepare("SELECT filter FROM fixed_price_filters WHERE fixedPriceId=?")
	if err != nil {
		return err
	}
	defer stmt.Close()

	rows, err := stmt.Query(fixedPriceID)
	if err != nil {
		return err
	}

	existingFilters := []string{}
	var filter string
	for rows.Next() {
		if err := rows.Scan(&filter); err != nil {
			return err
		}
		existingFilters = append(existingFilters, filter)
	}

	for _, existingFilter := range existingFilters {
		found := false
		for _, filter := range fixedPrice.Filters {
			if existingFilter == filter {
				found = true
			}
		}
		if !found {
			stmt, err := db.Prepare("DELETE FROM fixed_price_filters WHERE fixedPriceId=? AND filter=?")
			if err != nil {
				return err
			}
			defer stmt.Close()
			_, err = stmt.Exec(fixedPriceID, existingFilter)
			if err != nil {
				return err
			}
		}
	}

	for _, filter := range fixedPrice.Filters {
		found := false
		for _, existingFilter := range existingFilters {
			if filter == existingFilter {
				found = true
			}
		}
		if !found {
			stmt, err := db.Prepare("INSERT INTO fixed_price_filters (fixedPriceId, filter) VALUES (?, ?)")
			if err != nil {
				return err
			}
			defer stmt.Close()
			_, err = stmt.Exec(fixedPriceID, filter)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func updateStatesAndCities(fixedPriceID int64, fixedPrice *models.ServiceDetails) error {
	stmt, err := db.Prepare("SELECT state FROM fixed_price_state_cities WHERE fixedPriceId=?")
	if err != nil {
		return err
	}
	defer stmt.Close()

	rows, err := stmt.Query(fixedPriceID)
	if err != nil {
		return err
	}

	existingStates := []string{}
	var state string
	for rows.Next() {
		if err := rows.Scan(&state); err != nil {
			return err
		}
		existingStates = append(existingStates, state)
	}

	for _, existingState := range existingStates {
		found := false
		for _, obj := range fixedPrice.StatesAndCities {
			if existingState == obj.State {
				found = true
			}
		}
		if !found {
			stmt, err := db.Prepare("DELETE FROM fixed_price_state_cities WHERE fixedPriceId=? AND state=?")
			if err != nil {
				return err
			}
			defer stmt.Close()
			_, err = stmt.Exec(fixedPriceID, existingState)
			if err != nil {
				return err
			}
		}
	}

	for _, obj := range fixedPrice.StatesAndCities {
		cities, err := json.Marshal(obj.Cities)
		if err != nil {
			return err
		}
		stmt, err := db.Prepare("SELECT id FROM fixed_price_state_cities WHERE state=? AND fixedPriceId=?")
		if err != nil {
			return err
		}
		defer stmt.Close()

		var id int64
		row := stmt.QueryRow(obj.State, fixedPriceID)
		switch err = row.Scan(&id); err {
		case sql.ErrNoRows:
			stmt, err := db.Prepare("INSERT INTO fixed_price_state_cities (fixedPriceId, state, cities) VALUES (?, ?, ?)")
			if err != nil {
				return err
			}
			defer stmt.Close()
			_, err = stmt.Exec(fixedPriceID, obj.State, string(cities))
			if err != nil {
				return err
			}
		case nil:
			stmt, err := db.Prepare("UPDATE fixed_price_state_cities SET cities=? WHERE fixedPriceId=? AND state=?")
			if err != nil {
				return err
			}
			defer stmt.Close()
			_, err = stmt.Exec(string(cities), fixedPriceID, obj.State)
			if err != nil {
				return err
			}
		default:
			log.Println(err)
		}
	}
	return nil
}

func updateIncludes(fixedPriceID int64, fixedPrice *models.ServiceDetails) error {
	includes, err := json.Marshal(fixedPrice.Includes)
	if err != nil {
		return err
	}
	notIncludes, err := json.Marshal(fixedPrice.NotIncludes)
	if err != nil {
		return err
	}
	stmt, err := db.Prepare("SELECT included FROM fixed_price_includes WHERE fixedPriceId=?")
	if err != nil {
		return err
	}
	defer stmt.Close()

	rows, err := stmt.Query(fixedPriceID)
	if err != nil {
		return err
	}

	exist := false
	var included bool
	for rows.Next() {
		if err := rows.Scan(&included); err != nil {
			return err
		}
		exist = true
		if included {
			stmt, err := db.Prepare("UPDATE fixed_price_includes SET items=? WHERE fixedPriceId=? AND included=?")
			if err != nil {
				return err
			}
			defer stmt.Close()
			_, err = stmt.Exec(string(includes), fixedPriceID, included)
			if err != nil {
				return err
			}
		} else {
			stmt, err := db.Prepare("UPDATE fixed_price_includes SET items=? WHERE fixedPriceId=? AND included=?")
			if err != nil {
				return err
			}
			defer stmt.Close()
			_, err = stmt.Exec(string(notIncludes), fixedPriceID, included)
			if err != nil {
				return err
			}
		}
	}
	if !exist {
		stmt, err := db.Prepare("INSERT INTO fixed_price_includes (fixedPriceId, included, items) VALUES (?, ?, ?), (?, ?, ?)")
		if err != nil {
			return err
		}
		defer stmt.Close()
		_, err = stmt.Exec(fixedPriceID, true, string(includes), fixedPriceID, false, string(notIncludes))
		if err != nil {
			return err
		}
	}

	return nil
}

func UpdateFixedPrice(tradespersonID, priceID string, fixedPrice *models.ServiceDetails) (bool, error) {
	updated := false

	images, err := internal.ProcessImages(tradespersonID, priceID, fixedPrice)
	if err != nil {
		return updated, err
	}

	p, err := price.Get(priceID, nil)
	if err != nil {
		return updated, err
	}

	if err := stripe.UpdateProduct(images, fixedPrice, p); err != nil {
		return updated, err
	}

	stmt, err := db.Prepare("SELECT id FROM fixed_prices WHERE tradespersonId=? AND priceId=?")
	if err != nil {
		return updated, err
	}
	defer stmt.Close()

	var fixedPriceID int64
	row := stmt.QueryRow(tradespersonID, priceID)
	switch err = row.Scan(&fixedPriceID); err {
	case sql.ErrNoRows:
		return updated, fmt.Errorf("FixedPriced with priceId %s does not exist", priceID)
	case nil:
		stmt, err := db.Prepare("UPDATE fixed_prices SET category=?, subcategory=?, title=?, description=?, selectPlaces=?, archived=? WHERE priceId=?")
		if err != nil {
			return updated, err
		}
		defer stmt.Close()

		_, err = stmt.Exec(fixedPrice.Category, fixedPrice.SubCategory, fixedPrice.Title, fixedPrice.Description, fixedPrice.SelectPlaces, fixedPrice.Archived, priceID)
		if err != nil {
			return updated, err
		}

		if err := updateImages(fixedPriceID, images, "fixed_price"); err != nil {
			return updated, err
		}

		if err := updateStatesAndCities(fixedPriceID, fixedPrice); err != nil {
			return updated, err
		}

		if err := updateFilters(fixedPriceID, fixedPrice); err != nil {
			return updated, err
		}

		if err := UpdateTimeSlots(fixedPriceID, fixedPrice); err != nil {
			return updated, err
		}

		if err := updateIncludes(fixedPriceID, fixedPrice); err != nil {
			return updated, err
		}
		updated = true
	default:
		log.Println(err)
	}

	return updated, nil
}

func GetTradespersonFixedPrices(tradespersonID string, page int64) []*models.Service {
	fixedPrices := []*models.Service{}

	stmt, err := db.Prepare("SELECT id, priceId, title, price, description, subscription, subInterval, archived FROM fixed_prices WHERE tradespersonId=? ORDER BY id DESC LIMIT ?, ?")
	if err != nil {
		log.Printf("Failed to create select statement %s", err)
		return fixedPrices
	}
	defer stmt.Close()

	offSet := (page - 1) * int64(PAGE_SIZE)
	rows, err := stmt.Query(tradespersonID, offSet, PAGE_SIZE)
	if err != nil {
		log.Printf("Failed to execute select statement %s", err)
		return fixedPrices
	}

	var ID, price int64
	var interval sql.NullString
	for rows.Next() {
		fixedPrice := &models.Service{}
		if err := rows.Scan(&ID, &fixedPrice.PriceID, &fixedPrice.Title, &price, &fixedPrice.Description, &fixedPrice.Subscription, &interval, &fixedPrice.Archived); err != nil {
			log.Printf("Failed to scan for fixed_price, %v", err)
			return fixedPrices
		}
		if interval.Valid {
			fixedPrice.Interval = interval.String
		}
		strPrice := fmt.Sprintf("%.2f", float64(price)/float64(100.00))
		floatPrice, err := strconv.ParseFloat(strPrice, 64)
		if err != nil {
			log.Printf("Failed to parse float, %v", err)
			return fixedPrices
		}
		fixedPrice.Price = floatPrice
		fixedPrice.Image, err = GetImage(ID, "fixed_price")
		if err != nil {
			log.Printf("Failed to get fixedPrice image %s", err)
		}

		fixedPrice.AvailableTimeSlots, err = GetAvailableTimeSlots(ID, fixedPrice.Subscription)
		if err != nil {
			log.Printf("Failed to get timeslots %s", err)
			return fixedPrices
		}

		fixedPrice.Reviews, fixedPrice.Rating, err = GetFixedPriceReviewsRating(ID)
		if err != nil {
			log.Printf("Failed to get reviews and rating %s", err)
			return fixedPrices

		}
		fixedPrices = append(fixedPrices, fixedPrice)
	}

	return fixedPrices
}

//Find better way
func updateImages(ID int64, images []*string, table string) error {
	deleteSql := "DELETE FROM quote_images WHERE quoteId=?"
	insertSql := "INSERT INTO quote_images (quoteId, url) VALUES (?, ?)"
	if table == "fixed_price" {
		deleteSql = "DELETE FROM fixed_price_images WHERE fixedPriceId=?"
		insertSql = "INSERT INTO fixed_price_images (fixedPriceId, url) VALUES (?, ?)"
	}
	stmt, err := db.Prepare(deleteSql)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(ID)
	if err != nil {
		return err
	}

	for _, url := range images {
		stmt, err := db.Prepare(insertSql)
		if err != nil {
			return err
		}
		defer stmt.Close()
		_, err = stmt.Exec(ID, url)
		if err != nil {
			return err
		}
	}

	return nil
}

func insertQuoteFilters(ID int64, quote *models.ServiceDetails) error {
	for _, filter := range quote.Filters {
		stmt, err := db.Prepare("INSERT INTO quote_filters (quoteId, filter) VALUES (?, ?)")
		if err != nil {
			return err
		}
		defer stmt.Close()
		_, err = stmt.Exec(ID, filter)
		if err != nil {
			return err
		}
	}
	return nil
}

func insertQuoteStatesAndCities(ID int64, quote *models.ServiceDetails) error {
	for _, obj := range quote.StatesAndCities {
		cities, err := json.Marshal(obj.Cities)
		if err != nil {
			return err
		}
		stmt, err := db.Prepare("INSERT INTO quote_state_cities (quoteId, state, cities) VALUES (?, ?, ?)")
		if err != nil {
			return err
		}
		defer stmt.Close()
		_, err = stmt.Exec(ID, obj.State, string(cities))
		if err != nil {
			return err
		}
	}
	return nil
}

func CreateQuote(tradespersonID string, quote *models.ServiceDetails) (bool, error) {

	quoteID := "quote_" + internal.GenerateQuoteSuffix()

	stmt, err := db.Prepare("INSERT INTO quotes (quote, tradespersonId, category, subcategory, title, description, selectPlaces, archived) VALUES (?, ?, ?, ?, ?, ?, ?, false)")
	if err != nil {
		return false, err
	}
	defer stmt.Close()

	results, err := stmt.Exec(quoteID, tradespersonID, quote.Category, quote.SubCategory, quote.Title, quote.Description, quote.SelectPlaces)
	if err != nil {
		return false, err
	}
	rowsAffected, err := results.RowsAffected()
	if err != nil {
		return false, err
	}
	ID, err := results.LastInsertId()
	if err != nil {
		return false, err
	}
	if rowsAffected == 1 {
		images, err := internal.ProcessImages(tradespersonID, quoteID, quote)
		if err != nil {
			return false, err
		}

		if err = updateImages(ID, images, "quote"); err != nil {
			return false, err
		}

		if err = insertQuoteStatesAndCities(ID, quote); err != nil {
			return false, err
		}

		if err = insertQuoteFilters(ID, quote); err != nil {
			return false, err
		}
	}

	return true, nil
}

func GetImage(ID int64, table string) (string, error) {
	url := ""
	selectSQL := "SELECT url FROM quote_images WHERE quoteId=?"

	if table == "fixed_price" {
		selectSQL = "SELECT url FROM fixed_price_images WHERE fixedPriceId=?"
	}

	stmt, err := db.Prepare(selectSQL)
	if err != nil {
		return url, err
	}
	defer stmt.Close()

	row := stmt.QueryRow(ID)
	if err != nil {
		return url, err
	}
	row.Scan(&url)
	if url == "" {
		url = "https://" + os.Getenv("SUBDOMAIN") + "redbudway.com/assets/images/deal.svg"
	}

	return url, nil
}

func GetImages(ID int64, table string) ([]string, error) {
	images := []string{}
	selectSQL := "SELECT url FROM quote_images WHERE quoteId=?"
	if table == "fixed_price" {
		selectSQL = "SELECT url FROM fixed_price_images WHERE fixedPriceId=?"
	}
	stmt, err := db.Prepare(selectSQL)
	if err != nil {
		return images, err
	}
	defer stmt.Close()

	rows, err := stmt.Query(ID)
	if err != nil {
		return images, err
	}

	var url string
	for rows.Next() {
		if err := rows.Scan(&url); err != nil {
			return images, err
		}
		images = append(images, url)
	}
	if len(images) == 0 {
		images = append(images, "https://"+os.Getenv("SUBDOMAIN")+"redbudway.com/assets/images/deal.svg")
	}

	return images, nil
}

func GetQuoteRating(ID int64) (int64, float64, error) {
	reviews := int64(0)
	businessRating := float64(0.0)
	stmt, err := db.Prepare("SELECT qr.rating FROM quote_reviews qr INNER JOIN quotes q ON q.id=qr.quoteId WHERE q.id=?")
	if err != nil {
		return reviews, businessRating, err
	}
	defer stmt.Close()

	rows, err := stmt.Query(ID)
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

func GetTradespersonQuote(tradespersonID, quoteID string) (*models.ServiceDetails, error) {
	quote := &models.ServiceDetails{}

	stmt, err := db.Prepare("SELECT q.id, q.category, q.subcategory, q.title, q.description, q.selectPlaces, q.archived FROM tradesperson_account m INNER JOIN quotes q ON m.tradespersonId=q.tradespersonId WHERE m.tradespersonId=? AND q.quote=?")
	if err != nil {
		return quote, err
	}
	defer stmt.Close()

	row := stmt.QueryRow(tradespersonID, quoteID)
	var ID int64
	switch err = row.Scan(&ID, &quote.Category, &quote.SubCategory, &quote.Title, &quote.Description, &quote.SelectPlaces, &quote.Archived); err {
	case sql.ErrNoRows:
		return quote, err
	case nil:
		quote.ID = ID
		quote.Images, err = GetImages(ID, "quote")
		if err != nil {
			return quote, err
		}
		quote.StatesAndCities, err = GetQuoteStatesAndCities(ID)
		if err != nil {
			return quote, err
		}
		quote.Filters, err = GetFilters(ID)
		if err != nil {
			return quote, err
		}
	default:
		return quote, err
	}

	return quote, err
}

func updateQuoteFilters(ID int64, quote *models.ServiceDetails) error {
	stmt, err := db.Prepare("SELECT filter FROM quote_filters WHERE quoteId=?")
	if err != nil {
		return err
	}
	defer stmt.Close()

	rows, err := stmt.Query(ID)
	if err != nil {
		return err
	}

	existingFilters := []string{}
	var filter string
	for rows.Next() {
		if err := rows.Scan(&filter); err != nil {
			return err
		}
		existingFilters = append(existingFilters, filter)
	}

	for _, existingFilter := range existingFilters {
		found := false
		for _, filter := range quote.Filters {
			if existingFilter == filter {
				found = true
			}
		}
		if !found {
			stmt, err := db.Prepare("DELETE FROM quote_filters WHERE quoteId=? AND filter=?")
			if err != nil {
				return err
			}
			defer stmt.Close()
			_, err = stmt.Exec(ID, existingFilter)
			if err != nil {
				return err
			}
		}
	}

	for _, filter := range quote.Filters {
		found := false
		for _, existingFilter := range existingFilters {
			if filter == existingFilter {
				found = true
			}
		}
		if !found {
			stmt, err := db.Prepare("INSERT INTO quote_filters (quoteId, filter) VALUES (?, ?)")
			if err != nil {
				return err
			}
			defer stmt.Close()
			_, err = stmt.Exec(ID, filter)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func updateQuoteStatesAndCities(ID int64, quote *models.ServiceDetails) error {
	stmt, err := db.Prepare("SELECT state FROM quote_state_cities WHERE quoteId=?")
	if err != nil {
		return err
	}
	defer stmt.Close()

	rows, err := stmt.Query(ID)
	if err != nil {
		return err
	}

	existingStates := []string{}
	var state string
	for rows.Next() {
		if err := rows.Scan(&state); err != nil {
			return err
		}
		existingStates = append(existingStates, state)
	}

	for _, existingState := range existingStates {
		found := false
		for _, obj := range quote.StatesAndCities {
			if existingState == obj.State {
				found = true
			}
		}
		if !found {
			stmt, err := db.Prepare("DELETE FROM quote_state_cities WHERE quoteId=? AND state=?")
			if err != nil {
				return err
			}
			defer stmt.Close()
			_, err = stmt.Exec(ID, existingState)
			if err != nil {
				return err
			}
		}
	}

	for _, obj := range quote.StatesAndCities {
		cities, err := json.Marshal(obj.Cities)
		if err != nil {
			return err
		}
		stmt, err := db.Prepare("SELECT id FROM quote_state_cities WHERE state=? AND quoteId=?")
		if err != nil {
			return err
		}
		defer stmt.Close()

		var id int64
		row := stmt.QueryRow(obj.State, ID)
		switch err = row.Scan(&id); err {
		case sql.ErrNoRows:
			stmt, err := db.Prepare("INSERT INTO quote_state_cities (quoteId, state, cities) VALUES (?, ?, ?)")
			if err != nil {
				return err
			}
			defer stmt.Close()
			_, err = stmt.Exec(ID, obj.State, string(cities))
			if err != nil {
				return err
			}
		case nil:
			stmt, err := db.Prepare("UPDATE quote_state_cities SET cities=? WHERE quoteId=? AND state=?")
			if err != nil {
				return err
			}
			defer stmt.Close()
			_, err = stmt.Exec(string(cities), ID, obj.State)
			if err != nil {
				return err
			}
		default:
			log.Println(err)
		}
	}
	return nil
}

func UpdateTradespersonQuote(tradespersonID string, quoteID string, quote *models.ServiceDetails) (bool, error) {
	updated := false

	stmt, err := db.Prepare("SELECT id FROM quotes WHERE tradespersonId=? AND quote=?")
	if err != nil {
		return updated, err
	}
	defer stmt.Close()

	var ID int64
	row := stmt.QueryRow(tradespersonID, quoteID)
	switch err = row.Scan(&ID); err {
	case sql.ErrNoRows:
		return updated, fmt.Errorf("quote with quoteId %s does not exist", quoteID)
	case nil:
		stmt, err := db.Prepare("UPDATE quotes SET category=?, subcategory=?, title=?, description=?, selectPlaces=?, archived=? WHERE quote=?")
		if err != nil {
			return updated, err
		}
		defer stmt.Close()

		_, err = stmt.Exec(quote.Category, quote.SubCategory, quote.Title, quote.Description, quote.SelectPlaces, quote.Archived, quoteID)
		if err != nil {
			return updated, err
		}

		images, err := internal.ProcessImages(tradespersonID, quoteID, quote)
		if err != nil {
			return updated, err
		}

		if err := updateImages(ID, images, "quote"); err != nil {
			return updated, err
		}

		if err := updateQuoteStatesAndCities(ID, quote); err != nil {
			return updated, err
		}

		if err := updateQuoteFilters(ID, quote); err != nil {
			return updated, err
		}

		updated = true
	default:
		log.Println(err)
	}

	return updated, nil
}

func GetTradespersonQuotes(tradespersonID string, page int64) []*models.Service {

	quotes := []*models.Service{}

	stmt, err := db.Prepare("SELECT id, quote, title, description, archived FROM quotes WHERE tradespersonId=? ORDER BY id DESC LIMIT ?, ?")
	if err != nil {
		log.Printf("Failed to create select statement %s", err)
		return quotes
	}
	defer stmt.Close()

	offSet := (page - 1) * int64(PAGE_SIZE)
	rows, err := stmt.Query(tradespersonID, offSet, PAGE_SIZE)
	if err != nil {
		log.Printf("Failed to execute select statement %s", err)
		return quotes
	}

	var ID int64
	var description string
	for rows.Next() {
		quote := models.Service{}
		if err := rows.Scan(&ID, &quote.QuoteID, &quote.Title, &description, &quote.Archived); err != nil {
			log.Printf("Failed to scan row %s\n", err)
			return quotes
		}
		if len(description) > 65 {
			description = description[:65] + "..."
		}
		quote.Description = description
		quote.Reviews, quote.Rating, err = GetQuoteRating(ID)
		if err != nil {
			log.Printf("Failed to get quote reviews and rating %s\n", err)
		}
		quote.Image, err = GetImage(ID, "quote")
		if err != nil {
			log.Printf("Failed to get quote image %s\n", err)
		}
		quotes = append(quotes, &quote)
	}
	return quotes
}
