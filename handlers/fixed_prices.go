package handlers

import (
	"database/sql"
	"fmt"
	"log"
	"math"
	"redbudway-api/database"
	"redbudway-api/models"
	"redbudway-api/restapi/operations"
	"strconv"
	"strings"

	"github.com/go-openapi/runtime/middleware"
)

const PAGE_SIZE = float64(9)

func processFixedPriceRows(db *sql.DB, rows *sql.Rows, fixedPrices []*models.Service) ([]*models.Service, error) {
	var id, price int64
	var vanityURL, interval sql.NullString
	var subscription bool
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
		fixedPrice.Image, err = database.GetImage(id, "fixed_price")
		if err != nil {
			log.Printf("Failed to get fixedPrice image %s", err)
		}
		fixedPrice.VanityURL = vanityURL.String
		fixedPrice.TradespersonID = tradespersonID
		fixedPrice.AvailableTimeSlots, err = database.GetAvailableTimeSlots(id, subscription)
		if err != nil {
			log.Printf("Failed to get timeslots %s", err)
		}

		fixedPrice.Reviews, fixedPrice.Rating, err = database.GetFixedPriceReviewsRating(id)
		if err != nil {
			log.Printf("Failed to get reviews and rating %s", err)
		}

		fixedPrices = append(fixedPrices, fixedPrice)
	}
	return fixedPrices, nil
}

func getFixedPricesWithFilters(state, city, category, subCategory, filters string, page int64) ([]*models.Service, error) {
	filterArry := strings.Split(filters, ",")
	query := ""
	for _, filter := range filterArry {
		query += "'" + filter + "',"
	}
	query = query[:len(query)-1]
	fixedPrices := []*models.Service{}

	db := database.GetConnection()

	stmt, err := db.Prepare("SELECT a.stripeId, f.id, f.tradespersonId, f.priceId, f.title, f.price, f.subscription, f.subInterval, s.vanityURL FROM tradesperson_account a INNER JOIN tradesperson_settings s ON a.tradespersonId=s.tradespersonId INNER JOIN fixed_prices f ON a.tradespersonId=f.tradespersonId LEFT JOIN fixed_price_state_cities c ON c.fixedPriceId=f.id LEFT JOIN fixed_price_filters fi ON fi.fixedPriceId=f.id WHERE f.category=? AND f.subcategory=? AND (f.selectPlaces=false OR c.state=? OR JSON_CONTAINS(c.cities, JSON_OBJECT('name', ?))) AND fi.filter IN (" + query + ") AND f.archived=false GROUP BY f.id ORDER BY f.id DESC LIMIT ?, ?")
	if err != nil {
		log.Printf("Failed to create select statement %s", err)
		return fixedPrices, err
	}
	defer stmt.Close()

	offset := (page - 1) * int64(PAGE_SIZE)
	rows, err := stmt.Query(category, subCategory, state, city, offset, PAGE_SIZE)
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

func getSubCategoryFixedPrices(state, city, category, subCategory string, page int64) ([]*models.Service, error) {
	fixedPrices := []*models.Service{}

	db := database.GetConnection()

	stmt, err := db.Prepare("SELECT a.stripeId, f.id, f.tradespersonId, f.priceId, f.title, f.price, f.subscription, f.subInterval, s.vanityURL FROM tradesperson_account a INNER JOIN tradesperson_settings s ON a.tradespersonId=s.tradespersonId INNER JOIN fixed_prices f ON a.tradespersonId=f.tradespersonId LEFT JOIN fixed_price_state_cities c ON c.fixedPriceId=f.id WHERE f.category=? AND f.subcategory=? AND (f.selectPlaces=false OR c.state=? OR JSON_CONTAINS(c.cities, JSON_OBJECT('name', ?))) AND f.archived=false GROUP BY f.id ORDER BY f.id DESC LIMIT ?, ?")
	if err != nil {
		log.Printf("Failed to create select statement %s", err)
		return fixedPrices, err
	}
	defer stmt.Close()

	offset := (page - 1) * int64(PAGE_SIZE)
	rows, err := stmt.Query(category, subCategory, state, city, offset, PAGE_SIZE)
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

func getCategoryFixedPrices(state, city, category string, page int64) ([]*models.Service, error) {
	fixedPrices := []*models.Service{}

	db := database.GetConnection()

	stmt, err := db.Prepare("SELECT a.stripeId, f.id, f.tradespersonId, f.priceId, f.title, f.price, f.subscription, f.subInterval, s.vanityURL FROM tradesperson_account a INNER JOIN tradesperson_settings s ON a.tradespersonId=s.tradespersonId INNER JOIN fixed_prices f ON a.tradespersonId=f.tradespersonId LEFT JOIN fixed_price_state_cities c ON c.fixedPriceId=f.id WHERE f.category=? AND (f.selectPlaces=false OR c.state=? OR JSON_CONTAINS(c.cities, JSON_OBJECT('name', ?))) AND f.archived=false GROUP BY f.id ORDER BY f.id DESC LIMIT ?, ?")
	if err != nil {
		log.Printf("Failed to create select statement %s", err)
		return fixedPrices, err
	}
	defer stmt.Close()

	offset := (page - 1) * int64(PAGE_SIZE)
	rows, err := stmt.Query(category, state, city, offset, PAGE_SIZE)
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

func getAllFixedPrices(state, city string, page int64) ([]*models.Service, error) {
	fixedPrices := []*models.Service{}

	db := database.GetConnection()

	stmt, err := db.Prepare("SELECT a.stripeId, f.id, f.tradespersonId, f.priceId, f.title, f.price, f.subscription, f.subInterval, s.vanityURL FROM tradesperson_account a INNER JOIN tradesperson_settings s ON a.tradespersonId=s.tradespersonId INNER JOIN fixed_prices f ON a.tradespersonId=f.tradespersonId LEFT JOIN fixed_price_state_cities c ON c.fixedPriceId=f.id WHERE (f.selectPlaces=false OR c.state=? OR JSON_CONTAINS(c.cities, JSON_OBJECT('name', ?))) AND f.archived=false GROUP BY f.id ORDER BY f.id DESC LIMIT ?, ?")
	if err != nil {
		return fixedPrices, err
	}
	defer stmt.Close()

	offSet := (page - 1) * int64(PAGE_SIZE)
	rows, err := stmt.Query(state, city, offSet, PAGE_SIZE)
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

func getFixedPricesWithFiltersWithOutLocation(category, subCategory, filters string, page int64) ([]*models.Service, error) {
	filterArry := strings.Split(filters, ",")
	query := ""
	for _, filter := range filterArry {
		query += "'" + filter + "',"
	}
	query = query[:len(query)-1]
	fixedPrices := []*models.Service{}

	db := database.GetConnection()

	stmt, err := db.Prepare("SELECT a.stripeId, f.id, f.tradespersonId, f.priceId, f.title, f.price, f.subscription, f.subInterval, s.vanityURL FROM tradesperson_account a INNER JOIN tradesperson_settings s ON a.tradespersonId=s.tradespersonId INNER JOIN fixed_prices f ON a.tradespersonId=f.tradespersonId LEFT JOIN fixed_price_filters fi ON fi.fixedPriceId=f.id WHERE f.category=? AND f.subcategory=? AND fi.filter IN (" + query + ") AND f.archived=false GROUP BY f.id ORDER BY f.id DESC LIMIT ?, ?")
	if err != nil {
		log.Printf("Failed to create select statement %s", err)
		return fixedPrices, err
	}
	defer stmt.Close()

	offset := (page - 1) * int64(PAGE_SIZE)
	rows, err := stmt.Query(category, subCategory, offset, PAGE_SIZE)
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

func getSubCategoryFixedPricesWithOutLocation(category, subCategory string, page int64) ([]*models.Service, error) {
	fixedPrices := []*models.Service{}

	db := database.GetConnection()

	stmt, err := db.Prepare("SELECT a.stripeId, f.id, f.tradespersonId, f.priceId, f.title, f.price, f.subscription, f.subInterval, s.vanityURL FROM tradesperson_account a INNER JOIN tradesperson_settings s ON a.tradespersonId=s.tradespersonId INNER JOIN fixed_prices f ON a.tradespersonId=f.tradespersonId WHERE f.category=? AND f.subcategory=? AND f.archived=false GROUP BY f.id ORDER BY f.id DESC LIMIT ?, ?")
	if err != nil {
		log.Printf("Failed to create select statement %s", err)
		return fixedPrices, err
	}
	defer stmt.Close()

	offset := (page - 1) * int64(PAGE_SIZE)
	rows, err := stmt.Query(category, subCategory, offset, PAGE_SIZE)
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

func getCategoryFixedPricesWithOutLocation(category string, page int64) ([]*models.Service, error) {
	fixedPrices := []*models.Service{}

	db := database.GetConnection()

	stmt, err := db.Prepare("SELECT a.stripeId, f.id, f.tradespersonId, f.priceId, f.title, f.price, f.subscription, f.subInterval, s.vanityURL FROM tradesperson_account a INNER JOIN tradesperson_settings s ON a.tradespersonId=s.tradespersonId INNER JOIN fixed_prices f ON a.tradespersonId=f.tradespersonId WHERE f.category=? AND f.archived=false GROUP BY f.id ORDER BY f.id DESC LIMIT ?, ?")
	if err != nil {
		log.Printf("Failed to create select statement %s", err)
		return fixedPrices, err
	}
	defer stmt.Close()

	offset := (page - 1) * int64(PAGE_SIZE)
	rows, err := stmt.Query(category, offset, PAGE_SIZE)
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

func getAllFixedPricesWithOutLocation(page int64) ([]*models.Service, error) {
	fixedPrices := []*models.Service{}

	db := database.GetConnection()

	stmt, err := db.Prepare("SELECT a.stripeId, f.id, f.tradespersonId, f.priceId, f.title, f.price, f.subscription, f.subInterval, s.vanityURL FROM tradesperson_account a INNER JOIN tradesperson_settings s ON a.tradespersonId=s.tradespersonId INNER JOIN fixed_prices f ON a.tradespersonId=f.tradespersonId WHERE f.archived=false GROUP BY f.id ORDER BY f.id DESC LIMIT ?, ?")
	if err != nil {
		return fixedPrices, err
	}
	defer stmt.Close()

	offSet := (page - 1) * int64(PAGE_SIZE)
	rows, err := stmt.Query(offSet, PAGE_SIZE)
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
	if params.State != nil {
		state = *params.State
	}
	if params.City != nil {
		city = *params.City
	}
	page := *params.Page
	response := operations.NewGetFixedPricesOK()
	fixedPrices := []*models.Service{}
	response.SetPayload(fixedPrices)

	var err error
	if params.Category == nil && params.SubCategory == nil {
		if city == "" && state == "" {
			fixedPrices, err = getAllFixedPricesWithOutLocation(page)
			if err != nil {
				log.Printf("%s", err)
				return response
			}
		} else {
			fixedPrices, err = getAllFixedPrices(state, city, page)
			if err != nil {
				log.Printf("%s", err)
				return response
			}
		}
	} else if params.Category != nil && params.SubCategory == nil {
		if city == "" && state == "" {
			fixedPrices, err = getCategoryFixedPricesWithOutLocation(*params.Category, page)
			if err != nil {
				log.Printf("%s", err)
				return response
			}
		} else {
			fixedPrices, err = getCategoryFixedPrices(state, city, *params.Category, page)
			if err != nil {
				log.Printf("%s", err)
				return response
			}
		}
	} else if params.Category != nil && params.SubCategory != nil {
		if city == "" && state == "" {
			if params.Filters == nil {
				fixedPrices, err = getSubCategoryFixedPricesWithOutLocation(*params.Category, *params.SubCategory, page)
				if err != nil {
					log.Printf("%s", err)
					return response
				}
			} else {
				fixedPrices, err = getFixedPricesWithFiltersWithOutLocation(*params.Category, *params.SubCategory, *params.Filters, page)
				if err != nil {
					log.Printf("%s", err)
					return response
				}
			}
		} else {
			if params.Filters == nil {
				fixedPrices, err = getSubCategoryFixedPrices(state, city, *params.Category, *params.SubCategory, page)
				if err != nil {
					log.Printf("%s", err)
					return response
				}
			} else {
				fixedPrices, err = getFixedPricesWithFilters(state, city, *params.Category, *params.SubCategory, *params.Filters, page)
				if err != nil {
					log.Printf("%s", err)
					return response
				}
			}
		}
	}

	response.SetPayload(fixedPrices)

	return response
}

//PAGES
func getFixedPricesWithFiltersPages(pages float64, state, city, category, subCategory, filters string) (float64, error) {
	filterArry := strings.Split(filters, ",")
	query := ""
	for _, filter := range filterArry {
		query += "'" + filter + "',"
	}
	query = query[:len(query)-1]

	db := database.GetConnection()

	stmt, err := db.Prepare("SELECT COUNT(*) FROM fixed_prices f LEFT JOIN fixed_price_state_cities c ON c.fixedPriceId=f.id LEFT JOIN fixed_price_filters fi ON fi.fixedPriceId=f.id WHERE f.category=? AND f.subcategory=? AND (f.selectPlaces=false OR c.state=? OR JSON_CONTAINS(c.cities, JSON_OBJECT('name', ?))) AND fi.filter IN (" + query + ") AND f.archived=false")
	if err != nil {
		log.Printf("Failed to create select statement %s", err)
		return pages, err
	}
	defer stmt.Close()

	err = stmt.QueryRow(category, subCategory, state, city).Scan(&pages)
	if err != nil {
		log.Printf("Failed to execute select statement %s", err)
		return pages, err
	}

	return pages, nil
}

func getSubCategoryFixedPricePages(pages float64, state, city, category, subCategory string) (float64, error) {
	db := database.GetConnection()

	stmt, err := db.Prepare("SELECT COUNT(*) FROM fixed_prices f LEFT JOIN fixed_price_state_cities c ON c.fixedPriceId=f.id WHERE f.category=? AND f.subcategory=? AND (f.selectPlaces=false OR c.state=? OR JSON_CONTAINS(c.cities, JSON_OBJECT('name', ?))) AND f.archived=false")
	if err != nil {
		log.Printf("Failed to create select statement %s", err)
		return pages, err
	}
	defer stmt.Close()

	err = stmt.QueryRow(category, state, city).Scan(&pages)
	if err != nil {
		return pages, err
	}

	return pages, nil
}

func getCategoryFixedPricePages(pages float64, state, city, category string) (float64, error) {
	db := database.GetConnection()

	stmt, err := db.Prepare("SELECT COUNT(*) FROM fixed_prices f LEFT JOIN fixed_price_state_cities c ON c.fixedPriceId=f.id WHERE f.category=? AND (f.selectPlaces=false OR c.state=? OR JSON_CONTAINS(c.cities, JSON_OBJECT('name', ?))) AND f.archived=false")
	if err != nil {
		return pages, err
	}
	defer stmt.Close()

	err = stmt.QueryRow(category, state, city).Scan(&pages)
	if err != nil {
		return pages, err
	}

	return pages, nil
}

func getAllFixedPricePages(pages float64, state, city string) (float64, error) {

	db := database.GetConnection()

	stmt, err := db.Prepare("SELECT COUNT(*) FROM fixed_prices f LEFT JOIN fixed_price_state_cities c ON c.fixedPriceId=f.id WHERE (f.selectPlaces=false OR c.state=? OR JSON_CONTAINS(c.cities, JSON_OBJECT('name', ?))) AND f.archived=false")
	if err != nil {
		return pages, err
	}
	defer stmt.Close()

	err = stmt.QueryRow(state, city).Scan(&pages)
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
		pages, err = getAllFixedPricePages(pages, state, city)
		if err != nil {
			log.Printf("%s", err)
			return response
		}
	} else if params.Category != nil && params.SubCategory == nil {
		pages, err = getCategoryFixedPricePages(pages, state, city, *params.Category)
		if err != nil {
			log.Printf("%s", err)
			return response
		}
	} else if params.Category != nil && params.SubCategory != nil {
		if params.Filters == nil {
			pages, err = getSubCategoryFixedPricePages(pages, state, city, *params.Category, *params.SubCategory)
			if err != nil {
				log.Printf("%s", err)
				return response
			}
		} else {
			pages, err = getFixedPricesWithFiltersPages(pages, state, city, *params.Category, *params.SubCategory, *params.Filters)
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
