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

func processQuoteRows(db *sql.DB, rows *sql.Rows, quotes []*models.Service) ([]*models.Service, error) {
	var id int64
	var vanityURL sql.NullString
	var stripeID, quoteID, title, description, tradespersonID string
	for rows.Next() {
		if err := rows.Scan(&stripeID, &id, &tradespersonID, &quoteID, &title, &description, &vanityURL); err != nil {
			return quotes, err
		}
		tradesperson, err := database.GetTradespersonProfile(tradespersonID)
		if err != nil {
			log.Printf("Failed to get tradesperson profile %s", err)
		}
		quote := &models.Service{}
		quote.Business = tradesperson.Name
		quote.Title = title
		if len(description) > 84 {
			description = description[:84] + "..."
		}
		quote.Description = description
		quote.QuoteID = quoteID
		quote.TradespersonID = tradespersonID
		if vanityURL.Valid {
			quote.VanityURL = vanityURL.String
		}
		quote.Reviews, quote.Rating, err = database.GetQuoteRating(id)
		if err != nil {
			log.Printf("Failed to get quote reviews and rating %s", err)
		}
		quote.Image, err = database.GetImage(id, "quote")
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

	stmt, err := db.Prepare("SELECT a.stripeId, q.id, q.tradespersonId, q.quote, q.title, q.description, s.vanityURL FROM tradesperson_account a INNER JOIN tradesperson_settings s ON a.tradespersonId=s.tradespersonId INNER JOIN quotes q ON a.tradespersonId=q.tradespersonId LEFT JOIN quote_state_cities c ON c.quoteId=q.id LEFT JOIN quote_filters qf ON qf.quoteId=q.id WHERE q.category=? AND q.subcategory=? AND (q.selectPlaces=false OR c.state=? OR JSON_CONTAINS(c.cities, JSON_OBJECT('name', ?))) AND qf.filter IN (" + query + ") AND q.archived=false GROUP BY q.id ORDER BY q.id DESC LIMIT ?, ?")
	if err != nil {
		log.Printf("Failed to create select statement %s", err)
		return quotes, err
	}
	defer stmt.Close()

	offset := (page - 1) * int64(PAGE_SIZE)
	rows, err := stmt.Query(category, subCategory, state, city, offset, PAGE_SIZE)
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

	stmt, err := db.Prepare("SELECT a.stripeId, q.id, q.tradespersonId, q.quote, q.title, q.description, s.vanityURL FROM tradesperson_account a INNER JOIN tradesperson_settings s ON a.tradespersonId=s.tradespersonId INNER JOIN quotes q ON a.tradespersonId=q.tradespersonId LEFT JOIN quote_state_cities c ON c.quoteId=q.id WHERE q.category=? AND q.subcategory=? AND (q.selectPlaces=false OR c.state=? OR JSON_CONTAINS(c.cities, JSON_OBJECT('name', ?))) AND q.archived=false GROUP BY q.id ORDER BY q.id DESC LIMIT ?, ?")
	if err != nil {
		log.Printf("Failed to create select statement %s", err)
		return quotes, err
	}
	defer stmt.Close()

	offset := (page - 1) * int64(PAGE_SIZE)
	rows, err := stmt.Query(category, subCategory, state, city, offset, PAGE_SIZE)
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

	stmt, err := db.Prepare("SELECT a.stripeId, q.id, q.tradespersonId, q.quote, q.title, q.description, s.vanityURL FROM tradesperson_account a INNER JOIN tradesperson_settings s ON a.tradespersonId=s.tradespersonId INNER JOIN quotes q ON a.tradespersonId=q.tradespersonId LEFT JOIN quote_state_cities c ON c.quoteId=q.id WHERE q.category=? AND (q.selectPlaces=false OR c.state=? OR JSON_CONTAINS(c.cities, JSON_OBJECT('name', ?))) AND q.archived=false GROUP BY q.id ORDER BY q.id DESC LIMIT ?, ?")
	if err != nil {
		log.Printf("Failed to create select statement %s", err)
		return quotes, err
	}
	defer stmt.Close()

	offset := (page - 1) * int64(PAGE_SIZE)
	rows, err := stmt.Query(category, state, city, offset, PAGE_SIZE)
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

	stmt, err := db.Prepare("SELECT a.stripeId, q.id, q.tradespersonId, q.quote, q.title, q.description, s.vanityURL FROM tradesperson_account a INNER JOIN tradesperson_settings s ON a.tradespersonId=s.tradespersonId INNER JOIN quotes q ON a.tradespersonId=q.tradespersonId LEFT JOIN quote_state_cities c ON c.quoteId=q.id WHERE (q.selectPlaces=false OR c.state=? OR JSON_CONTAINS(c.cities, JSON_OBJECT('name', ?))) AND q.archived=false GROUP BY q.id ORDER BY q.id DESC LIMIT ?, ?")
	if err != nil {
		return quotes, err
	}
	defer stmt.Close()

	offset := (page - 1) * int64(PAGE_SIZE)
	rows, err := stmt.Query(state, city, offset, PAGE_SIZE)
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

func getQuotesWithFiltersWOL(category, subCategory, filters string, page int64) ([]*models.Service, error) {
	filterArry := strings.Split(filters, ",")
	query := ""
	for _, filter := range filterArry {
		query += "'" + filter + "',"
	}
	query = query[:len(query)-1]
	log.Println(query)
	quotes := []*models.Service{}

	db := database.GetConnection()

	stmt, err := db.Prepare("SELECT a.stripeId, q.id, q.tradespersonId, q.quote, q.title, q.description, s.vanityURL FROM tradesperson_account a INNER JOIN tradesperson_settings s ON a.tradespersonId=s.tradespersonId INNER JOIN quotes q ON a.tradespersonId=q.tradespersonId WHERE q.category=? AND q.subcategory=? AND qf.filter IN (" + query + ") AND q.archived=false GROUP BY q.id ORDER BY q.id DESC LIMIT ?, ?")
	if err != nil {
		log.Printf("Failed to create select statement %s", err)
		return quotes, err
	}
	defer stmt.Close()

	offset := (page - 1) * int64(PAGE_SIZE)
	rows, err := stmt.Query(category, subCategory, offset, PAGE_SIZE)
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

func getSubCategoryQuotesWOL(category, subCategory string, page int64) ([]*models.Service, error) {
	quotes := []*models.Service{}

	db := database.GetConnection()

	stmt, err := db.Prepare("SELECT a.stripeId, q.id, q.tradespersonId, q.quote, q.title, q.description, s.vanityURL FROM tradesperson_account a INNER JOIN tradesperson_settings s ON a.tradespersonId=s.tradespersonId INNER JOIN quotes q ON a.tradespersonId=q.tradespersonId WHERE q.category=? AND q.subcategory=? AND q.archived=false GROUP BY q.id ORDER BY q.id DESC LIMIT ?, ?")
	if err != nil {
		log.Printf("Failed to create select statement %s", err)
		return quotes, err
	}
	defer stmt.Close()

	offset := (page - 1) * int64(PAGE_SIZE)
	rows, err := stmt.Query(category, subCategory, offset, PAGE_SIZE)
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

func getCategoryQuotesWOL(category string, page int64) ([]*models.Service, error) {
	quotes := []*models.Service{}

	db := database.GetConnection()

	stmt, err := db.Prepare("SELECT a.stripeId, q.id, q.tradespersonId, q.quote, q.title, q.description, s.vanityURL FROM tradesperson_account a INNER JOIN tradesperson_settings s ON a.tradespersonId=s.tradespersonId INNER JOIN quotes q ON a.tradespersonId=q.tradespersonId WHERE q.category=? AND q.archived=false GROUP BY q.id ORDER BY q.id DESC LIMIT ?, ?")
	if err != nil {
		log.Printf("Failed to create select statement %s", err)
		return quotes, err
	}
	defer stmt.Close()

	offset := (page - 1) * int64(PAGE_SIZE)
	rows, err := stmt.Query(category, offset, PAGE_SIZE)
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

func getAllQuotesWOL(page int64) ([]*models.Service, error) {
	quotes := []*models.Service{}

	db := database.GetConnection()

	stmt, err := db.Prepare("SELECT a.stripeId, q.id, q.tradespersonId, q.quote, q.title, q.description, s.vanityURL FROM tradesperson_account a INNER JOIN tradesperson_settings s ON a.tradespersonId=s.tradespersonId INNER JOIN quotes q ON a.tradespersonId=q.tradespersonId WHERE q.archived=false GROUP BY q.id ORDER BY q.id DESC LIMIT ?, ?")
	if err != nil {
		return quotes, err
	}
	defer stmt.Close()

	offset := (page - 1) * int64(PAGE_SIZE)
	rows, err := stmt.Query(offset, PAGE_SIZE)
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
	state := ""
	city := ""
	if params.State != nil {
		state = *params.State
	}
	if params.City != nil {
		city = *params.City
	}
	page := *params.Page
	response := operations.NewGetQuotesOK()
	quotes := []*models.Service{}
	response.SetPayload(quotes)

	var err error
	if params.Category == nil && params.SubCategory == nil {
		if city == "" && state == "" {
			quotes, err = getAllQuotesWOL(page)
			if err != nil {
				log.Printf("%s", err)
				return response
			}
		} else {
			quotes, err = getAllQuotes(state, city, page)
			if err != nil {
				log.Printf("%s", err)
				return response
			}
		}
	} else if params.Category != nil && params.SubCategory == nil {
		if city == "" && state == "" {
			quotes, err = getCategoryQuotesWOL(*params.Category, page)
			if err != nil {
				log.Printf("%s", err)
				return response
			}
		} else {
			quotes, err = getCategoryQuotes(state, city, *params.Category, page)
			if err != nil {
				log.Printf("%s", err)
				return response
			}
		}
	} else if params.Category != nil && params.SubCategory != nil {
		if city == "" && state == "" {
			if params.Filters == nil {
				quotes, err = getSubCategoryQuotesWOL(*params.Category, *params.SubCategory, page)
				if err != nil {
					log.Printf("%s", err)
					return response
				}
			} else {
				quotes, err = getQuotesWithFiltersWOL(*params.Category, *params.SubCategory, *params.Filters, page)
				if err != nil {
					log.Printf("%s", err)
					return response
				}
			}
		} else {
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

	stmt, err := db.Prepare("SELECT COUNT(*) FROM quotes q LEFT JOIN quote_state_cities c ON c.quoteId=q.id LEFT JOIN quote_filters qf ON qf.quoteId=q.id WHERE q.category=? AND q.subcategory=? AND (q.selectPlaces=false OR c.state=? OR JSON_CONTAINS(c.cities, JSON_OBJECT('name', ?))) AND  qf.filter IN (" + query + ") AND q.archived=false")
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

	stmt, err := db.Prepare("SELECT COUNT(*) FROM quotes q LEFT JOIN quote_state_cities c ON c.quoteId=q.id WHERE q.category=? AND q.subcategory=? AND (q.selectPlaces=false OR c.state=? OR JSON_CONTAINS(c.cities, JSON_OBJECT('name', ?))) AND q.archived=false")
	if err != nil {
		log.Printf("Failed to create select statement %s", err)
		return pages, err
	}
	defer stmt.Close()

	err = stmt.QueryRow(category, subCategory, state, city).Scan(&pages)
	if err != nil {
		return pages, err
	}

	return pages, nil
}

func getCategoryQuotePages(pages float64, state, city, category string) (float64, error) {
	db := database.GetConnection()

	stmt, err := db.Prepare("SELECT COUNT(*) FROM quotes q LEFT JOIN quote_state_cities c ON c.quoteId=q.id WHERE q.category=? AND (q.selectPlaces=false OR c.state=? OR JSON_CONTAINS(c.cities, JSON_OBJECT('name', ?))) AND q.archived=false")
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

	stmt, err := db.Prepare("SELECT COUNT(*) FROM quotes q LEFT JOIN quote_state_cities c ON c.quoteId=q.id WHERE (q.selectPlaces=false OR c.state=? OR JSON_CONTAINS(c.cities, JSON_OBJECT('name', ?))) AND q.archived=false")
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

func getQuotesWithFiltersPagesWOL(pages float64, category, subCategory, filters string) (float64, error) {
	filterArry := strings.Split(filters, ",")
	query := ""
	for _, filter := range filterArry {
		query += "'" + filter + "',"
	}
	query = query[:len(query)-1]

	db := database.GetConnection()

	stmt, err := db.Prepare("SELECT COUNT(*) FROM quotes q LEFT JOIN quote_filters qf ON qf.quoteId=q.id WHERE q.category=? AND q.subcategory=? AND qf.filter IN (" + query + ") AND q.archived=false")
	if err != nil {
		log.Printf("Failed to create select statement %s", err)
		return pages, err
	}
	defer stmt.Close()

	err = stmt.QueryRow(category, subCategory).Scan(&pages)
	if err != nil {
		log.Printf("Failed to execute select statement %s", err)
		return pages, err
	}

	return pages, nil
}

func getSubCategoryQuotePagesWOL(pages float64, category, subCategory string) (float64, error) {
	db := database.GetConnection()

	stmt, err := db.Prepare("SELECT COUNT(*) FROM quotes WHERE category=? AND subcategory=? AND q.archived=false")
	if err != nil {
		log.Printf("Failed to create select statement %s", err)
		return pages, err
	}
	defer stmt.Close()

	err = stmt.QueryRow(category, subCategory).Scan(&pages)
	if err != nil {
		return pages, err
	}

	return pages, nil
}

func getCategoryQuotePagesWOL(pages float64, category string) (float64, error) {
	db := database.GetConnection()

	stmt, err := db.Prepare("SELECT COUNT(*) FROM quotes WHERE category=? AND archived=false")
	if err != nil {
		return pages, err
	}
	defer stmt.Close()

	err = stmt.QueryRow(category).Scan(&pages)
	if err != nil {
		return pages, err
	}

	return pages, nil
}

func getAllQuotePagesWOL(pages float64) (float64, error) {

	db := database.GetConnection()

	stmt, err := db.Prepare("SELECT COUNT(*) FROM quotes WHERE archived=false")
	if err != nil {
		return pages, err
	}
	defer stmt.Close()

	err = stmt.QueryRow().Scan(&pages)
	if err != nil {
		return pages, err
	}

	return pages, nil
}

func GetQuotePagesHandler(params operations.GetQuotePagesParams) middleware.Responder {
	state := ""
	city := ""
	if params.State != nil {
		state = *params.State
	}
	if params.City != nil {
		city = *params.City
	}

	pages := float64(1)
	response := operations.NewGetQuotePagesOK().WithPayload(int64(pages))

	var err error
	if params.Category == nil && params.SubCategory == nil {
		if city == "" && state == "" {
			pages, err = getAllQuotePagesWOL(pages)
			if err != nil {
				log.Printf("%s", err)
				return response
			}
		} else {
			pages, err = getAllQuotePages(pages, state, city)
			if err != nil {
				log.Printf("%s", err)
				return response
			}
		}
	} else if params.Category != nil && params.SubCategory == nil {
		if city == "" && state == "" {
			pages, err = getCategoryQuotePagesWOL(pages, *params.Category)
			if err != nil {
				log.Printf("%s", err)
				return response
			}
		} else {
			pages, err = getCategoryQuotePages(pages, state, city, *params.Category)
			if err != nil {
				log.Printf("%s", err)
				return response
			}
		}
	} else if params.Category != nil && params.SubCategory != nil {
		if city == "" && state == "" {
			if params.Filters == nil {
				pages, err = getSubCategoryQuotePagesWOL(pages, *params.Category, *params.SubCategory)
				if err != nil {
					log.Printf("%s", err)
					return response
				}
			} else {
				pages, err = getQuotesWithFiltersPagesWOL(pages, *params.Category, *params.SubCategory, *params.Filters)
				if err != nil {
					log.Printf("%s", err)
					return response
				}
			}
		} else {
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
	}

	if pages == float64(0) {
		pages = float64(1)
	}

	pages = math.Ceil(pages / PAGE_SIZE)

	response.SetPayload(int64(pages))

	return response
}
