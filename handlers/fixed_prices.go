package handlers

import (
	"database/sql"
	"fmt"
	"log"
	"redbudway-api/database"
	"redbudway-api/internal"
	"redbudway-api/models"
	"redbudway-api/restapi/operations"
	"strconv"
	"strings"

	"github.com/go-openapi/runtime/middleware"
	"github.com/stripe/stripe-go/v72/price"
	"github.com/stripe/stripe-go/v72/product"
)

func processFixedPriceRows(db *sql.DB, rows *sql.Rows, fixedPrices []*models.Service, city string) ([]*models.Service, error) {
	var id int64
	var vanityURL, citiesJson, interval sql.NullString
	var subscription, selectPlaces bool
	var name, priceID, tradespersonID string
	for rows.Next() {
		if err := rows.Scan(&name, &id, &tradespersonID, &priceID, &subscription, &interval, &selectPlaces, &vanityURL, &citiesJson); err != nil {
			return fixedPrices, err
		}
		fixedPrice := &models.Service{}
		fixedPrice.Subscription = subscription
		if interval.Valid {
			fixedPrice.Interval = interval.String
		}
		stripePrice, err := price.Get(priceID, nil)
		if err != nil {
			log.Printf("Failed to get stripe price, %v", err)
			return fixedPrices, err
		}
		stripeProduct, err := product.Get(stripePrice.Product.ID, nil)
		if err != nil {
			log.Printf("Failed to get stripe product, %v", err)
			return fixedPrices, err
		}
		fixedPrice.PriceID = priceID
		strPrice := fmt.Sprintf("%.2f", stripePrice.UnitAmountDecimal/float64(100.00))
		floatPrice, err := strconv.ParseFloat(strPrice, 64)
		if err != nil {
			log.Printf("Failed to parse float, %v", err)
			return fixedPrices, err
		}
		fixedPrice.Price = floatPrice
		fixedPrice.Title = &stripeProduct.Name
		fixedPrice.Image = &stripeProduct.Images[0]
		fixedPrice.VanityURL = vanityURL.String
		fixedPrice.Business = name
		fixedPrice.TradespersonID = tradespersonID
		fixedPrice.AvailableTimeSlots, err = database.GetAvailableTimeSlots(id, subscription)
		if err != nil {
			log.Printf("Failed to get timeslots %s", err)
		}

		fixedPrice.Reviews, fixedPrice.Rating, err = database.GetFixedPriceReviewsRating(id)
		if err != nil {
			log.Printf("Failed to get reviews and rating %s", err)
		}

		if selectPlaces {
			if citiesJson.Valid {
				cityExist, err := internal.SelectedCities(citiesJson.String, city)
				if err != nil {
					return fixedPrices, err
				}
				if cityExist {
					fixedPrices = append(fixedPrices, fixedPrice)
				}
			}
		} else {
			fixedPrices = append(fixedPrices, fixedPrice)
		}
	}
	return fixedPrices, nil
}

func getFixedPricesWithFilters(state, city, category, subCategory, filters string) ([]*models.Service, error) {
	filterArry := strings.Split(filters, ",")
	query := ""
	for _, filter := range filterArry {
		query += "'" + filter + "',"
	}
	query = query[:len(query)-1]
	log.Println(query)
	fixedPrices := []*models.Service{}

	db := database.GetConnection()

	stmt, err := db.Prepare("SELECT a.name, f.id, f.tradespersonId, f.priceId, f.subscription, f.subInterval, f.selectPlaces, s.vanityURL, c.cities FROM tradesperson_account a INNER JOIN tradesperson_settings s ON a.tradespersonId=s.tradespersonId INNER JOIN fixed_prices f ON a.tradespersonId=f.tradespersonId LEFT JOIN fixed_price_state_cities c ON c.fixedPriceId=f.id LEFT JOIN fixed_price_filters fi ON fi.fixedPriceId=f.id WHERE f.category=? AND f.subcategory=? AND (f.selectPlaces=false OR c.state=?) AND fi.filter IN (" + query + ") AND f.archived=false GROUP BY f.id ORDER BY rand()")
	if err != nil {
		log.Printf("Failed to create select statement %s", err)
		return fixedPrices, err
	}
	defer stmt.Close()

	rows, err := stmt.Query(category, subCategory, state)
	if err != nil {
		log.Printf("Failed to execute select statement %s", err)
		return fixedPrices, err
	}

	fixedPrices, err = processFixedPriceRows(db, rows, fixedPrices, city)
	if err != nil {
		log.Println("Failed to process rows")
		return fixedPrices, err
	}

	return fixedPrices, nil
}

func getSubCategoryFixedPrices(state, city, category, subCategory string) ([]*models.Service, error) {
	fixedPrices := []*models.Service{}

	db := database.GetConnection()

	stmt, err := db.Prepare("SELECT a.name, f.id, f.tradespersonId, f.priceId, f.subscription, f.subInterval, f.selectPlaces, s.vanityURL, c.cities FROM tradesperson_account a INNER JOIN tradesperson_settings s ON a.tradespersonId=s.tradespersonId INNER JOIN fixed_prices f ON a.tradespersonId=f.tradespersonId LEFT JOIN fixed_price_state_cities c ON c.fixedPriceId=f.id WHERE f.category=? AND f.subcategory=? AND (f.selectPlaces=false OR c.state=?) AND f.archived=false GROUP BY f.id ORDER BY rand()")
	if err != nil {
		log.Printf("Failed to create select statement %s", err)
		return fixedPrices, err
	}
	defer stmt.Close()

	rows, err := stmt.Query(category, subCategory, state)
	if err != nil {
		log.Printf("Failed to execute select statement %s", err)
		return fixedPrices, err
	}

	fixedPrices, err = processFixedPriceRows(db, rows, fixedPrices, city)
	if err != nil {
		log.Println("Failed to process rows")
		return fixedPrices, err
	}

	return fixedPrices, nil
}

func getCategoryFixedPrices(state, city, category string) ([]*models.Service, error) {
	fixedPrices := []*models.Service{}

	db := database.GetConnection()

	stmt, err := db.Prepare("SELECT a.name, f.id, f.tradespersonId, f.priceId, f.subscription, f.subInterval, f.selectPlaces, s.vanityURL, c.cities FROM tradesperson_account a INNER JOIN tradesperson_settings s ON a.tradespersonId=s.tradespersonId INNER JOIN fixed_prices f ON a.tradespersonId=f.tradespersonId LEFT JOIN fixed_price_state_cities c ON c.fixedPriceId=f.id WHERE f.category=? AND (f.selectPlaces=false OR c.state=?) AND f.archived=false GROUP BY f.id ORDER BY rand()")
	if err != nil {
		log.Printf("Failed to create select statement %s", err)
		return fixedPrices, err
	}
	defer stmt.Close()

	rows, err := stmt.Query(category, state)
	if err != nil {
		log.Printf("Failed to execute select statement %s", err)
		return fixedPrices, err
	}

	fixedPrices, err = processFixedPriceRows(db, rows, fixedPrices, city)
	if err != nil {
		log.Println("Failed to process rows")
		return fixedPrices, err
	}

	return fixedPrices, nil
}

func getAllFixedPrices(state, city string) ([]*models.Service, error) {
	fixedPrices := []*models.Service{}

	db := database.GetConnection()

	stmt, err := db.Prepare("SELECT a.name, f.id, f.tradespersonId, f.priceId, f.subscription, f.subInterval, f.selectPlaces, s.vanityURL, c.cities FROM tradesperson_account a INNER JOIN tradesperson_settings s ON a.tradespersonId=s.tradespersonId INNER JOIN fixed_prices f ON a.tradespersonId=f.tradespersonId LEFT JOIN fixed_price_state_cities c ON c.fixedPriceId=f.id WHERE (f.selectPlaces=false OR c.state=?) AND f.archived=false GROUP BY f.id ORDER BY rand() LIMIT 100")
	if err != nil {
		return fixedPrices, err
	}
	defer stmt.Close()

	rows, err := stmt.Query(state)
	if err != nil {
		log.Println("Failed to execute select statement")
		return fixedPrices, err
	}

	fixedPrices, err = processFixedPriceRows(db, rows, fixedPrices, city)
	if err != nil {
		log.Println("Failed to process rows")
		return fixedPrices, err
	}

	return fixedPrices, nil
}

func GetFixedPricesHandler(params operations.GetFixedPricesParams) middleware.Responder {
	state := params.State
	city := params.City
	response := operations.NewGetFixedPricesOK()
	fixedPrices := []*models.Service{}
	response.SetPayload(fixedPrices)

	var err error
	if params.Category == nil && params.SubCategory == nil {
		fixedPrices, err = getAllFixedPrices(state, city)
		if err != nil {
			log.Printf("%s", err)
			return response
		}
	} else if params.Category != nil && params.SubCategory == nil {
		fixedPrices, err = getCategoryFixedPrices(state, city, *params.Category)
		if err != nil {
			log.Printf("%s", err)
			return response
		}
	} else if params.Category != nil && params.SubCategory != nil {

		if params.Filters == nil {
			fixedPrices, err = getSubCategoryFixedPrices(state, city, *params.Category, *params.SubCategory)
			if err != nil {
				log.Printf("%s", err)
				return response
			}
		} else {
			fixedPrices, err = getFixedPricesWithFilters(state, city, *params.Category, *params.SubCategory, *params.Filters)
			if err != nil {
				log.Printf("%s", err)
				return response
			}
		}
	}

	response.SetPayload(fixedPrices)

	return response
}
