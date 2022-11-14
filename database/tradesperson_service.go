package database

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"redbudway-api/internal"
	"redbudway-api/models"
	"redbudway-api/stripe"
	"strconv"
	"time"

	"github.com/stripe/stripe-go/v72/price"
	"github.com/stripe/stripe-go/v72/product"
)

func insertTimeSlots(fixedPriceID int64, fixedPrice *models.ServiceDetails) error {
	for _, timeSlot := range fixedPrice.TimeSlots {
		t, err := time.Parse("1/2/2006, 3:04:00 PM", timeSlot.StartTime)
		if err != nil {
			return err
		}
		formattedDate := t.Format("2006-01-02 15:04:00")
		stmt, err := db.Prepare("INSERT INTO fixed_price_time_slots (fixedPriceId, startTime, segmentSize, taken, takenBy, cuStripeId) VALUES (?, ?, ?, DEFAULT, DEFAULT, DEFAULT)")
		if err != nil {
			return err
		}
		defer stmt.Close()
		_, err = stmt.Exec(fixedPriceID, formattedDate, timeSlot.SegmentSize)
		if err != nil {
			return err
		}
	}
	return nil
}

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
	stmt, err := db.Prepare("INSERT INTO fixed_prices (tradespersonId, category, subcategory, priceId, subscription, subInterval, selectPlaces, archived) VALUES (?, ?, ?, ?, ?, ?, ?, false)")
	if err != nil {
		return false, err
	}
	defer stmt.Close()

	results, err := stmt.Exec(tradespersonID, fixedPrice.Category, fixedPrice.SubCategory, price.ID, fixedPrice.Subscription, fixedPrice.Interval, fixedPrice.SelectPlaces)
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

		if err = insertStatesAndCities(fixedPriceID, fixedPrice); err != nil {
			return false, err
		}

		if err = insertTimeSlots(fixedPriceID, fixedPrice); err != nil {
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

func GetStatesAndCities(fixedPriceID int64) ([]*models.ServiceDetailsStatesAndCitiesItems0, error) {
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
		stmt, err := db.Prepare("SELECT startTime, segmentSize, taken, takenBy FROM fixed_price_time_slots WHERE fixedPriceId=?")
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

		var startTime, segmentSize, takenBy string
		var taken bool
		for rows.Next() {
			if err := rows.Scan(&startTime, &segmentSize, &taken, &takenBy); err != nil {
				log.Printf("Failed to scan statement, %s", err)
				continue
			}
			timeSlot := &models.TimeSlot{}
			timeSlot.StartTime = startTime
			timeSlot.SegmentSize = segmentSize
			timeSlot.Taken = taken
			timeSlot.TakenBy = takenBy
			timeSlots = append(timeSlots, timeSlot)
		}
		service.TimeSlots = timeSlots
		otherServices = append(otherServices, service)
	}
	return otherServices, nil
}

func GetTradespersonFixedPrice(tradespersonID string, priceID string) (*models.ServiceDetails, int64, error) {
	fixedPrice := &models.ServiceDetails{}

	stmt, err := db.Prepare("SELECT f.id, f.category, f.subCategory, f.subscription, f.subInterval, f.selectPlaces, f.archived FROM tradesperson_account a INNER JOIN fixed_prices f ON a.tradespersonId=f.tradespersonId WHERE a.tradespersonId=? AND f.priceId=?")
	if err != nil {
		return fixedPrice, 0, err
	}
	defer stmt.Close()

	row := stmt.QueryRow(tradespersonID, priceID)
	var fixedPriceID int64
	var interval sql.NullString
	var category, subCategory string
	var subscription, selectPlaces, archived bool
	switch err = row.Scan(&fixedPriceID, &category, &subCategory, &subscription, &interval, &selectPlaces, &archived); err {
	case sql.ErrNoRows:
		return fixedPrice, fixedPriceID, err
	case nil:
		fixedPrice.Category = &category
		fixedPrice.SubCategory = subCategory
		fixedPrice.Subscription = subscription
		if interval.Valid {
			fixedPrice.Interval = interval.String
		}
		fixedPrice.SelectPlaces = &selectPlaces
		fixedPrice.Archived = archived
		stripePrice, err := price.Get(priceID, nil)
		if err != nil {
			return fixedPrice, fixedPriceID, err
		}
		stripeProduct, err := product.Get(stripePrice.Product.ID, nil)
		if err != nil {
			return fixedPrice, fixedPriceID, err
		}
		strPrice := fmt.Sprintf("%.2f", stripePrice.UnitAmountDecimal/float64(100.00))
		floatPrice, err := strconv.ParseFloat(strPrice, 64)
		if err != nil {
			return fixedPrice, fixedPriceID, err
		}
		fixedPrice.Price = floatPrice
		fixedPrice.Images = stripeProduct.Images
		fixedPrice.Title = &stripeProduct.Name
		fixedPrice.Description = &stripeProduct.Description
		fixedPrice.TimeSlots, err = GetTimeSlots(fixedPriceID)
		if err != nil {
			return fixedPrice, fixedPriceID, err
		}
		fixedPrice.StatesAndCities, err = GetStatesAndCities(fixedPriceID)
		if err != nil {
			return fixedPrice, fixedPriceID, err
		}
		fixedPrice.Filters, err = GetFilters(fixedPriceID)
		if err != nil {
			return fixedPrice, fixedPriceID, err
		}
	default:
		return fixedPrice, fixedPriceID, err
	}

	return fixedPrice, fixedPriceID, nil
}

func updateTimeSlots(fixedPriceID int64, fixedPrice *models.ServiceDetails) error {
	stmt, err := db.Prepare("SELECT startTime, segmentSize FROM fixed_price_time_slots WHERE fixedPriceId=?")
	if err != nil {
		return err
	}
	defer stmt.Close()

	rows, err := stmt.Query(fixedPriceID)
	if err != nil {
		return err
	}

	existingTimeSlots := []*models.TimeSlotsItems0{}
	var startTime, segmentSize string
	for rows.Next() {
		if err := rows.Scan(&startTime, &segmentSize); err != nil {
			return err
		}
		existingTimeSlot := &models.TimeSlotsItems0{}
		existingTimeSlot.StartTime = startTime
		existingTimeSlot.SegmentSize = segmentSize
		existingTimeSlots = append(existingTimeSlots, existingTimeSlot)
	}

	for _, existingTimeSlot := range existingTimeSlots {
		found := false
		for _, timeSlot := range fixedPrice.TimeSlots {
			t, err := time.Parse("1/2/2006, 3:04:00 PM", timeSlot.StartTime)
			if err != nil {
				return err
			}
			formattedDate := t.Format("2006-01-02 15:04:00")
			if formattedDate == existingTimeSlot.StartTime {
				found = true
				if timeSlot.SegmentSize != existingTimeSlot.SegmentSize {
					stmt, err := db.Prepare("UPDATE fixed_price_time_slots SET segmentSize=? WHERE fixedPriceId=? AND startTime=?")
					if err != nil {
						return err
					}
					defer stmt.Close()
					_, err = stmt.Exec(existingTimeSlot.SegmentSize, fixedPriceID, formattedDate)
					if err != nil {
						return err
					}
				}
			}
		}
		if !found {
			stmt, err := db.Prepare("DELETE FROM fixed_price_time_slots WHERE takenBy IS NULL AND fixedPriceId=? AND startTime=?")
			if err != nil {
				return err
			}
			defer stmt.Close()
			_, err = stmt.Exec(fixedPriceID, existingTimeSlot.StartTime)
			if err != nil {
				return err
			}
		}
	}

	for _, timeSlot := range fixedPrice.TimeSlots {
		t, err := time.Parse("1/2/2006, 3:04:00 PM", timeSlot.StartTime)
		if err != nil {
			return err
		}
		formattedDate := t.Format("2006-01-02 15:04:00")
		found := false
		for _, existingTimeSlot := range existingTimeSlots {
			if formattedDate == existingTimeSlot.StartTime {
				found = true
			}
		}
		if !found {
			stmt, err := db.Prepare("INSERT INTO fixed_price_time_slots (fixedPriceId, startTime, segmentSize, taken) VALUES (?, ?, ?, DEFAULT)")
			if err != nil {
				return err
			}
			defer stmt.Close()
			_, err = stmt.Exec(fixedPriceID, formattedDate, timeSlot.SegmentSize)
			if err != nil {
				return err
			}
		}
	}

	return nil
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
		return updated, errors.New(fmt.Sprintf("FixedPriced with priceId %s does not exist", priceID))
	case nil:
		stmt, err := db.Prepare("UPDATE fixed_prices SET category=?, subcategory=?, selectPlaces=?, archived=? WHERE priceId=?")
		if err != nil {
			return updated, err
		}
		defer stmt.Close()

		_, err = stmt.Exec(fixedPrice.Category, fixedPrice.SubCategory, fixedPrice.SelectPlaces, fixedPrice.Archived, priceID)
		if err != nil {
			return updated, err
		}

		if err := updateStatesAndCities(fixedPriceID, fixedPrice); err != nil {
			return updated, err
		}

		if err := updateFilters(fixedPriceID, fixedPrice); err != nil {
			return updated, err
		}

		if err := updateTimeSlots(fixedPriceID, fixedPrice); err != nil {
			return updated, err
		}
		updated = true
	default:
		log.Println(err)
	}

	return updated, nil
}

func GetTradespersonFixedPrices(tradespersonID string) []*models.Service {
	fixedPrices := []*models.Service{}

	stmt, err := db.Prepare("SELECT id, priceId, subscription, subInterval, selectPlaces, archived FROM fixed_prices WHERE tradespersonId=?")
	if err != nil {
		log.Printf("Failed to create select statement %s", err)
		return fixedPrices
	}
	defer stmt.Close()

	rows, err := stmt.Query(tradespersonID)
	if err != nil {
		log.Printf("Failed to execute select statement %s", err)
		return fixedPrices
	}

	var id int64
	var interval sql.NullString
	var subscription, selectPlaces, archived bool
	var priceID string
	for rows.Next() {
		if err := rows.Scan(&id, &priceID, &subscription, &interval, &selectPlaces, &archived); err != nil {
			log.Printf("Failed to scan for fixed_price, %v", err)
			return fixedPrices
		}
		fixedPrice := &models.Service{}
		fixedPrice.Subscription = subscription
		if interval.Valid {
			fixedPrice.Interval = interval.String
		}
		stripePrice, err := price.Get(priceID, nil)
		if err != nil {
			log.Printf("Failed to get stripe price, %v", err)
			return fixedPrices
		}
		stripeProduct, err := product.Get(stripePrice.Product.ID, nil)
		if err != nil {
			log.Printf("Failed to get stripe product, %v", err)
			return fixedPrices
		}
		fixedPrice.PriceID = priceID
		strPrice := fmt.Sprintf("%.2f", stripePrice.UnitAmountDecimal/float64(100.00))
		floatPrice, err := strconv.ParseFloat(strPrice, 64)
		if err != nil {
			log.Printf("Failed to parse float, %v", err)
			return fixedPrices
		}
		fixedPrice.Price = floatPrice
		fixedPrice.Title = stripeProduct.Name
		fixedPrice.Image = stripeProduct.Images[0]

		fixedPrice.AvailableTimeSlots, err = GetAvailableTimeSlots(id, subscription)
		if err != nil {
			log.Printf("Failed to get timeslots %s", err)
			return fixedPrices
		}

		fixedPrice.Reviews, fixedPrice.Rating, err = GetFixedPriceReviewsRating(id)
		if err != nil {
			log.Printf("Failed to get reviews and rating %s", err)
			return fixedPrices

		}
		fixedPrices = append(fixedPrices, fixedPrice)
	}

	return fixedPrices
}

//Find better way
func updateImages(ID int64, images []*string) error {
	stmt, err := db.Prepare("DELETE FROM quote_images WHERE quoteId=?")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(ID)
	if err != nil {
		return err
	}

	for _, url := range images {
		stmt, err := db.Prepare("INSERT INTO quote_images (quoteId, url) VALUES (?, ?)")
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

		if err = updateImages(ID, images); err != nil {
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

func GetQuoteImage(ID int64) (string, error) {
	url := ""

	stmt, err := db.Prepare("SELECT url FROM quote_images WHERE quoteId=?")
	if err != nil {
		return url, err
	}
	defer stmt.Close()

	row := stmt.QueryRow(ID)
	if err != nil {
		return url, err
	}
	row.Scan(&url)

	return url, nil
}

func GetQuoteImages(ID int64) ([]string, error) {
	images := []string{}

	stmt, err := db.Prepare("SELECT url FROM quote_images WHERE quoteId=?")
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
	var category, subCategory, title, description string
	var selectPlaces, archived bool
	switch err = row.Scan(&ID, &category, &subCategory, &title, &description, &selectPlaces, &archived); err {
	case sql.ErrNoRows:
		return quote, err
	case nil:
		quote.ID = ID
		quote.Category = &category
		quote.SubCategory = subCategory
		quote.SelectPlaces = &selectPlaces
		quote.Archived = archived
		quote.Title = &title
		quote.Description = &description

		quote.Images, err = GetQuoteImages(ID)
		if err != nil {
			return quote, err
		}
		quote.StatesAndCities, err = GetStatesAndCities(ID)
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

		if err := updateImages(ID, images); err != nil {
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

func GetTradespersonQuotes(tradespersonID string) []*models.Service {

	quotes := []*models.Service{}

	stmt, err := db.Prepare("SELECT id, quote, title FROM quotes WHERE tradespersonId=?")
	if err != nil {
		log.Printf("Failed to create select statement %s", err)
		return quotes
	}
	defer stmt.Close()

	rows, err := stmt.Query(tradespersonID)
	if err != nil {
		log.Printf("Failed to execute select statement %s", err)
		return quotes
	}

	var ID int64
	var quoteID, title string
	for rows.Next() {
		if err := rows.Scan(&ID, &quoteID, &title); err != nil {
			log.Printf("Failed to scan row %s", err)
			return quotes
		}
		quote := models.Service{}
		quote.Title = title
		quote.QuoteID = quoteID
		quote.Reviews, quote.Rating, err = GetQuoteRating(ID)
		if err != nil {
			log.Printf("Failed to get quote reviews and rating %s", err)
		}
		quote.Image, err = GetQuoteImage(ID)
		if err != nil {
			log.Printf("Failed to get quote image %s", err)
		}
		quotes = append(quotes, &quote)
	}
	return quotes
}
