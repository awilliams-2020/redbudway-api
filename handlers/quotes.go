package handlers

import (
	"database/sql"
	"log"
	"redbudway-api/database"
	"redbudway-api/internal"
	"redbudway-api/models"
	"redbudway-api/restapi/operations"
	"strings"

	"github.com/go-openapi/runtime/middleware"
)

func processQuoteRows(db *sql.DB, rows *sql.Rows, quotes []*models.Service, city string) ([]*models.Service, error) {
	var id int64
	var vanityURL, citiesJson sql.NullString
	var selectPlaces bool
	var name, quoteID, title, tradespersonID string
	for rows.Next() {
		if err := rows.Scan(&name, &id, &tradespersonID, &quoteID, &title, &selectPlaces, &vanityURL, &citiesJson); err != nil {
			return quotes, err
		}
		quote := &models.Service{}
		quote.Business = name
		quote.Title = title
		quote.QuoteID = quoteID
		quote.TradespersonID = tradespersonID
		if vanityURL.Valid {
			quote.VanityURL = vanityURL.String
		}
		var err error
		quote.Reviews, quote.Rating, err = database.GetQuoteRating(id)
		if err != nil {
			log.Printf("Failed to get quote reviews and rating %s", err)
		}
		quote.Image, err = database.GetQuoteImage(id)
		if err != nil {
			log.Printf("Failed to get quote image %s", err)
		}

		if selectPlaces {
			if citiesJson.Valid {
				cityExist, err := internal.SelectedCities(citiesJson.String, city)
				if err != nil {
					return quotes, err
				}
				if cityExist {
					quotes = append(quotes, quote)
				}
			}
		} else {
			quotes = append(quotes, quote)
		}
	}
	return quotes, nil
}

func getQuotesWithFilters(state, city, category, subCategory, filters string) ([]*models.Service, error) {
	filterArry := strings.Split(filters, ",")
	query := ""
	for _, filter := range filterArry {
		query += "'" + filter + "',"
	}
	query = query[:len(query)-1]
	log.Println(query)
	quotes := []*models.Service{}

	db := database.GetConnection()

	stmt, err := db.Prepare("SELECT a.name, q.id, q.tradespersonId, q.quote, q.title, q.selectPlaces, s.vanityURL, c.cities FROM tradesperson_account a INNER JOIN tradesperson_settings s ON a.tradespersonId=s.tradespersonId INNER JOIN quotes q ON a.tradespersonId=q.tradespersonId LEFT JOIN quote_state_cities c ON c.quoteId=q.id LEFT JOIN quote_filters fi ON fi.quoteId=q.id WHERE q.category=? AND q.subcategory=? AND (q.selectPlaces=false OR c.state=?) AND fi.filter IN (" + query + ") AND q.archived=false GROUP BY q.id ORDER BY rand()")
	if err != nil {
		log.Printf("Failed to create select statement %s", err)
		return quotes, err
	}
	defer stmt.Close()

	rows, err := stmt.Query(category, subCategory, state)
	if err != nil {
		log.Printf("Failed to execute select statement %s", err)
		return quotes, err
	}

	quotes, err = processQuoteRows(db, rows, quotes, city)
	if err != nil {
		log.Println("Failed to process rows")
		return quotes, err
	}

	return quotes, nil
}

func getSubCategoryQuotes(state, city, category, subCategory string) ([]*models.Service, error) {
	quotes := []*models.Service{}

	db := database.GetConnection()

	stmt, err := db.Prepare("SELECT a.name, q.id, q.tradespersonId, q.quote, q.title, q.selectPlaces, s.vanityURL, c.cities FROM tradesperson_account a INNER JOIN tradesperson_settings s ON a.tradespersonId=s.tradespersonId INNER JOIN quotes q ON a.tradespersonId=q.tradespersonId LEFT JOIN quote_state_cities c ON c.quoteId=q.id WHERE q.category=? AND q.subcategory=? AND (q.selectPlaces=false OR c.state=?) AND q.archived=false GROUP BY q.id ORDER BY rand()")
	if err != nil {
		log.Printf("Failed to create select statement %s", err)
		return quotes, err
	}
	defer stmt.Close()

	rows, err := stmt.Query(category, subCategory, state)
	if err != nil {
		log.Printf("Failed to execute select statement %s", err)
		return quotes, err
	}

	quotes, err = processQuoteRows(db, rows, quotes, city)
	if err != nil {
		log.Println("Failed to process rows")
		return quotes, err
	}

	return quotes, nil
}

func getCategoryQuotes(state, city, category string) ([]*models.Service, error) {
	quotes := []*models.Service{}

	db := database.GetConnection()

	stmt, err := db.Prepare("SELECT a.name, q.id, q.tradespersonId, q.quote, q.title, q.selectPlaces, s.vanityURL, c.cities FROM tradesperson_account a INNER JOIN tradesperson_settings s ON a.tradespersonId=s.tradespersonId INNER JOIN quotes q ON a.tradespersonId=q.tradespersonId LEFT JOIN quote_state_cities c ON c.quoteId=q.id WHERE q.category=? AND (q.selectPlaces=false OR c.state=?) AND q.archived=false GROUP BY q.id ORDER BY rand()")
	if err != nil {
		log.Printf("Failed to create select statement %s", err)
		return quotes, err
	}
	defer stmt.Close()

	rows, err := stmt.Query(category, state)
	if err != nil {
		log.Printf("Failed to execute select statement %s", err)
		return quotes, err
	}

	quotes, err = processQuoteRows(db, rows, quotes, city)
	if err != nil {
		log.Println("Failed to process rows")
		return quotes, err
	}

	return quotes, nil
}

func getAllQuotes(state, city string) ([]*models.Service, error) {
	quotes := []*models.Service{}

	db := database.GetConnection()

	stmt, err := db.Prepare("SELECT a.name, q.id, q.tradespersonId, q.quote, q.title, q.selectPlaces, s.vanityURL, c.cities FROM tradesperson_account a INNER JOIN tradesperson_settings s ON a.tradespersonId=s.tradespersonId INNER JOIN quotes q ON a.tradespersonId=q.tradespersonId LEFT JOIN quote_state_cities c ON c.quoteId=q.id WHERE (q.selectPlaces=false OR c.state=?) AND q.archived=false GROUP BY q.id ORDER BY rand() LIMIT 50")
	if err != nil {
		return quotes, err
	}
	defer stmt.Close()

	rows, err := stmt.Query(state)
	if err != nil {
		log.Println("Failed to execute select statement")
		return quotes, err
	}

	quotes, err = processQuoteRows(db, rows, quotes, city)
	if err != nil {
		log.Println("Failed to process rows")
		return quotes, err
	}

	return quotes, nil
}

func GetQuotesHandler(params operations.GetQuotesParams) middleware.Responder {
	state := params.State
	city := params.City
	response := operations.NewGetQuotesOK()
	quotes := []*models.Service{}
	response.SetPayload(quotes)

	var err error
	if params.Category == nil && params.SubCategory == nil {
		quotes, err = getAllQuotes(state, city)
		if err != nil {
			log.Printf("%s", err)
			return response
		}
	} else if params.Category != nil && params.SubCategory == nil {
		quotes, err = getCategoryQuotes(state, city, *params.Category)
		if err != nil {
			log.Printf("%s", err)
			return response
		}
	} else if params.Category != nil && params.SubCategory != nil {

		if params.Filters == nil {
			quotes, err = getSubCategoryQuotes(state, city, *params.Category, *params.SubCategory)
			if err != nil {
				log.Printf("%s", err)
				return response
			}
		} else {
			quotes, err = getQuotesWithFilters(state, city, *params.Category, *params.SubCategory, *params.Filters)
			if err != nil {
				log.Printf("%s", err)
				return response
			}
		}
	}

	response.SetPayload(quotes)

	return response
}
