package handlers

import (
	"database/sql"
	"log"
	"math"
	"redbudway-api/database"
	"redbudway-api/models"
	"redbudway-api/restapi/operations"
	"strings"

	"github.com/go-openapi/runtime/middleware"
)

const QUOTE_PAGE_SIZE = 3

func processQuoteRows(db *sql.DB, rows *sql.Rows, quotes []*models.Service) ([]*models.Service, error) {
	var id int64
	var vanityURL sql.NullString
	var name, quoteID, title, tradespersonID string
	for rows.Next() {
		if err := rows.Scan(&name, &id, &tradespersonID, &quoteID, &title, &vanityURL); err != nil {
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

		quotes = append(quotes, quote)
	}
	return quotes, nil
}

func getQuotesWithFilters(state, city, category, subCategory, filters string, page int64) ([]*models.Service, error) {
	filterArry := strings.Split(filters, ",")
	query := ""
	for _, filter := range filterArry {
		query += "'" + filter + "',"
	}
	query = query[:len(query)-1]
	log.Println(query)
	quotes := []*models.Service{}

	db := database.GetConnection()

	stmt, err := db.Prepare("SELECT a.name, q.id, q.tradespersonId, q.quote, q.title, s.vanityURL FROM tradesperson_account a INNER JOIN tradesperson_settings s ON a.tradespersonId=s.tradespersonId INNER JOIN quotes q ON a.tradespersonId=q.tradespersonId LEFT JOIN quote_state_cities c ON c.quoteId=q.id LEFT JOIN quote_filters qf ON qf.quoteId=q.id WHERE q.category=? AND q.subcategory=? AND (q.selectPlaces=false OR c.state=?) AND (q.selectPlaces=false OR JSON_CONTAINS(c.cities, JSON_OBJECT('name', ?))) AND qf.filter IN (" + query + ") AND q.archived=false GROUP BY q.id ORDER BY q.id DESC LIMIT ?, ?")
	if err != nil {
		log.Printf("Failed to create select statement %s", err)
		return quotes, err
	}
	defer stmt.Close()

	rows, err := stmt.Query(category, subCategory, state, city, page, QUOTE_PAGE_SIZE)
	if err != nil {
		log.Printf("Failed to execute select statement %s", err)
		return quotes, err
	}

	quotes, err = processQuoteRows(db, rows, quotes)
	if err != nil {
		log.Println("Failed to process rows")
		return quotes, err
	}

	return quotes, nil
}

func getSubCategoryQuotes(state, city, category, subCategory string, page int64) ([]*models.Service, error) {
	quotes := []*models.Service{}

	db := database.GetConnection()

	stmt, err := db.Prepare("SELECT a.name, q.id, q.tradespersonId, q.quote, q.title, s.vanityURL FROM tradesperson_account a INNER JOIN tradesperson_settings s ON a.tradespersonId=s.tradespersonId INNER JOIN quotes q ON a.tradespersonId=q.tradespersonId LEFT JOIN quote_state_cities c ON c.quoteId=q.id WHERE q.category=? AND q.subcategory=? AND (q.selectPlaces=false OR c.state=?) AND (q.selectPlaces=false OR JSON_CONTAINS(c.cities, JSON_OBJECT('name', ?))) AND q.archived=false GROUP BY q.id ORDER BY q.id DESC LIMIT ?, ?")
	if err != nil {
		log.Printf("Failed to create select statement %s", err)
		return quotes, err
	}
	defer stmt.Close()

	rows, err := stmt.Query(category, subCategory, state, city, page, QUOTE_PAGE_SIZE)
	if err != nil {
		log.Printf("Failed to execute select statement %s", err)
		return quotes, err
	}

	quotes, err = processQuoteRows(db, rows, quotes)
	if err != nil {
		log.Println("Failed to process rows")
		return quotes, err
	}

	return quotes, nil
}

func getCategoryQuotes(state, city, category string, page int64) ([]*models.Service, error) {
	quotes := []*models.Service{}

	db := database.GetConnection()

	stmt, err := db.Prepare("SELECT a.name, q.id, q.tradespersonId, q.quote, q.title, s.vanityURL FROM tradesperson_account a INNER JOIN tradesperson_settings s ON a.tradespersonId=s.tradespersonId INNER JOIN quotes q ON a.tradespersonId=q.tradespersonId LEFT JOIN quote_state_cities c ON c.quoteId=q.id WHERE q.category=? AND (q.selectPlaces=false OR c.state=?) AND (q.selectPlaces=false OR JSON_CONTAINS(c.cities, JSON_OBJECT('name', ?))) AND q.archived=false GROUP BY q.id ORDER BY q.id DESC LIMIT ?, ?")
	if err != nil {
		log.Printf("Failed to create select statement %s", err)
		return quotes, err
	}
	defer stmt.Close()

	rows, err := stmt.Query(category, state, city, page, QUOTE_PAGE_SIZE)
	if err != nil {
		log.Printf("Failed to execute select statement %s", err)
		return quotes, err
	}

	quotes, err = processQuoteRows(db, rows, quotes)
	if err != nil {
		log.Println("Failed to process rows")
		return quotes, err
	}

	return quotes, nil
}

func getAllQuotes(state, city string, page int64) ([]*models.Service, error) {
	quotes := []*models.Service{}

	db := database.GetConnection()

	stmt, err := db.Prepare("SELECT a.name, q.id, q.tradespersonId, q.quote, q.title, s.vanityURL FROM tradesperson_account a INNER JOIN tradesperson_settings s ON a.tradespersonId=s.tradespersonId INNER JOIN quotes q ON a.tradespersonId=q.tradespersonId LEFT JOIN quote_state_cities c ON c.quoteId=q.id WHERE (q.selectPlaces=false OR c.state=?) AND (q.selectPlaces=false OR JSON_CONTAINS(c.cities, JSON_OBJECT('name', ?))) AND q.archived=false GROUP BY q.id ORDER BY q.id DESC LIMIT ?, ?")
	if err != nil {
		return quotes, err
	}
	defer stmt.Close()

	rows, err := stmt.Query(state, city, page, QUOTE_PAGE_SIZE)
	if err != nil {
		log.Println("Failed to execute select statement")
		return quotes, err
	}

	quotes, err = processQuoteRows(db, rows, quotes)
	if err != nil {
		log.Println("Failed to process rows")
		return quotes, err
	}

	return quotes, nil
}

func GetQuotesHandler(params operations.GetQuotesParams) middleware.Responder {
	state := params.State
	city := params.City
	page := *params.Page
	response := operations.NewGetQuotesOK()
	quotes := []*models.Service{}
	response.SetPayload(quotes)

	var err error
	if params.Category == nil && params.SubCategory == nil {
		quotes, err = getAllQuotes(state, city, page)
		if err != nil {
			log.Printf("%s", err)
			return response
		}
	} else if params.Category != nil && params.SubCategory == nil {
		quotes, err = getCategoryQuotes(state, city, *params.Category, page)
		if err != nil {
			log.Printf("%s", err)
			return response
		}
	} else if params.Category != nil && params.SubCategory != nil {

		if params.Filters == nil {
			quotes, err = getSubCategoryQuotes(state, city, *params.Category, *params.SubCategory, page)
			if err != nil {
				log.Printf("%s", err)
				return response
			}
		} else {
			quotes, err = getQuotesWithFilters(state, city, *params.Category, *params.SubCategory, *params.Filters, page)
			if err != nil {
				log.Printf("%s", err)
				return response
			}
		}
	}

	response.SetPayload(quotes)

	return response
}

//PAGES
func getQuotesWithFiltersPages(pages float64, state, city, category, subCategory, filters string) (float64, error) {
	filterArry := strings.Split(filters, ",")
	query := ""
	for _, filter := range filterArry {
		query += "'" + filter + "',"
	}
	query = query[:len(query)-1]

	db := database.GetConnection()

	stmt, err := db.Prepare("SELECT COUNT(*) FROM quotes q LEFT JOIN quote_state_cities c ON c.quoteId=q.id LEFT JOIN quote_filters qf ON qf.quoteId=q.id WHERE q.category=? AND q.subcategory=? AND (q.selectPlaces=false OR c.state=?) AND (q.selectPlaces=false OR JSON_CONTAINS(c.cities, JSON_OBJECT('name', ?))) AND  qf.filter IN (" + query + ") AND q.archived=false")
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

func getSubCategoryQuotePages(pages float64, state, city, category, subCategory string) (float64, error) {
	db := database.GetConnection()

	stmt, err := db.Prepare("SELECT COUNT(*) FROM quotes q LEFT JOIN quote_state_cities c ON c.quoteId=q.id WHERE q.category=? AND q.subcategory=? AND (q.selectPlaces=false OR c.state=?) AND (q.selectPlaces=false OR JSON_CONTAINS(c.cities, JSON_OBJECT('name', ?))) AND q.archived=false")
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

func getCategoryQuotePages(pages float64, state, city, category string) (float64, error) {
	db := database.GetConnection()

	stmt, err := db.Prepare("SELECT COUNT(*) FROM quotes q LEFT JOIN quote_state_cities c ON c.quoteId=q.id WHERE q.category=? AND (q.selectPlaces=false OR c.state=?) AND (q.selectPlaces=false OR JSON_CONTAINS(c.cities, JSON_OBJECT('name', ?))) AND q.archived=false")
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

func getAllQuotePages(pages float64, state, city string) (float64, error) {

	db := database.GetConnection()

	stmt, err := db.Prepare("SELECT COUNT(*) FROM quotes q LEFT JOIN quote_state_cities c ON c.quoteId=q.id WHERE (q.selectPlaces=false OR c.state=?) AND (q.selectPlaces=false OR JSON_CONTAINS(c.cities, JSON_OBJECT('name', ?))) AND q.archived=false")
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

func GetQuotePagesHandler(params operations.GetQuotePagesParams) middleware.Responder {
	state := params.State
	city := params.City

	pages := float64(1)
	response := operations.NewGetQuotePagesOK().WithPayload(int64(pages))

	var err error
	if params.Category == nil && params.SubCategory == nil {
		pages, err = getAllQuotePages(pages, state, city)
		if err != nil {
			log.Printf("%s", err)
			return response
		}
	} else if params.Category != nil && params.SubCategory == nil {
		pages, err = getCategoryQuotePages(pages, state, city, *params.Category)
		if err != nil {
			log.Printf("%s", err)
			return response
		}
	} else if params.Category != nil && params.SubCategory != nil {
		if params.Filters == nil {
			pages, err = getSubCategoryQuotePages(pages, state, city, *params.Category, *params.SubCategory)
			if err != nil {
				log.Printf("%s", err)
				return response
			}
		} else {
			pages, err = getQuotesWithFiltersPages(pages, state, city, *params.Category, *params.SubCategory, *params.Filters)
			if err != nil {
				log.Printf("%s", err)
				return response
			}
		}
	}

	if pages == float64(0) {
		pages = float64(1)
	}

	pages = math.Ceil(pages / QUOTE_PAGE_SIZE)

	response.SetPayload(int64(pages))

	return response
}
