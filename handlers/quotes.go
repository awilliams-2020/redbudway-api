package handlers

import (
	"database/sql"
	"log"
	"math"
	"redbudway-api/database"
	"redbudway-api/internal"
	"redbudway-api/models"
	"redbudway-api/restapi/operations"
	"strings"

	"github.com/go-openapi/runtime/middleware"
)

func processQuoteRows(db *sql.DB, rows *sql.Rows, quotes []*models.Service) ([]*models.Service, error) {
	var id int64
	var vanityURL, icon sql.NullString
	var stripeID, quoteID, title, description, tradespersonID string
	for rows.Next() {
		if err := rows.Scan(&stripeID, &id, &tradespersonID, &quoteID, &title, &description, &vanityURL, &icon); err != nil {
			return quotes, err
		}
		tradesperson, err := database.GetTradespersonProfile(tradespersonID)
		if err != nil {
			log.Printf("Failed to get tradesperson profile %s", err)
		}
		quote := &models.Service{}
		quote.Title = title
		if len(description) > 84 {
			description = description[:84] + "..."
		}
		quote.Description = description
		quote.QuoteID = quoteID
		quote.TradespersonID = tradespersonID
		quote.Business = &models.Business{
			Name:      tradesperson.Name,
			Icon:      icon.String,
			VanityURL: vanityURL.String,
		}
		quote.Reviews, quote.Rating, err = database.GetQuoteRating(id)
		if err != nil {
			log.Printf("Failed to get quote reviews and rating %s", err)
		}

		quote.Image, err = internal.GetImage(quoteID, tradespersonID)
		if err != nil {
			log.Printf("Failed to get quote image %s", err)
		}

		jobs, err := database.GetQuoteJobs(id, tradespersonID)
		if err != nil {
			log.Printf("Failed to get quote jobs %s", err)
		}
		quote.Jobs = jobs

		quotes = append(quotes, quote)
	}
	return quotes, nil
}

func getQuotesWithSpecialties(state, city, category, subCategory, specialties string, page int64) ([]*models.Service, error) {
	specialtyArry := strings.Split(specialties, ",")
	query := ""
	for _, specialty := range specialtyArry {
		query += "'" + specialty + "',"
	}
	query = query[:len(query)-1]
	quotes := []*models.Service{}

	db := database.GetConnection()

	stmt, err := db.Prepare("SELECT a.stripeId, q.id, q.tradespersonId, q.quote, q.title, q.description, s.vanityURL, b.icon FROM tradesperson_account a INNER JOIN tradesperson_settings s ON a.tradespersonId=s.tradespersonId INNER JOIN quotes q ON a.tradespersonId=q.tradespersonId LEFT JOIN quote_state_cities c ON c.quoteId=q.id LEFT JOIN tradesperson_branding b ON b.tradespersonId=q.tradespersonId LEFT JOIN quote_specialties qf ON qf.quoteId=q.id WHERE q.category=? AND q.subcategory=? AND ((? = '' AND ? = '') OR q.selectPlaces=false OR c.state=? OR JSON_CONTAINS(c.cities, JSON_OBJECT('name', ?))) AND qf.specialty IN (" + query + ") AND q.archived=false GROUP BY q.id ORDER BY q.id DESC LIMIT ?, ?")
	if err != nil {
		log.Printf("Failed to create select statement %s", err)
		return quotes, err
	}
	defer stmt.Close()

	offset := (page - 1) * int64(PAGE_SIZE)
	rows, err := stmt.Query(category, subCategory, state, city, state, city, offset, PAGE_SIZE)
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

	stmt, err := db.Prepare("SELECT a.stripeId, q.id, q.tradespersonId, q.quote, q.title, q.description, s.vanityURL, b.icon FROM tradesperson_account a INNER JOIN tradesperson_settings s ON a.tradespersonId=s.tradespersonId INNER JOIN quotes q ON a.tradespersonId=q.tradespersonId LEFT JOIN quote_state_cities c ON c.quoteId=q.id LEFT JOIN tradesperson_branding b ON b.tradespersonId=q.tradespersonId WHERE q.category=? AND q.subcategory=? AND ((? = '' AND ? = '') OR q.selectPlaces=false OR c.state=? OR JSON_CONTAINS(c.cities, JSON_OBJECT('name', ?))) AND q.archived=false GROUP BY q.id ORDER BY q.id DESC LIMIT ?, ?")
	if err != nil {
		log.Printf("Failed to create select statement %s", err)
		return quotes, err
	}
	defer stmt.Close()

	offset := (page - 1) * int64(PAGE_SIZE)
	rows, err := stmt.Query(category, subCategory, state, city, state, city, offset, PAGE_SIZE)
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

	stmt, err := db.Prepare("SELECT a.stripeId, q.id, q.tradespersonId, q.quote, q.title, q.description, s.vanityURL, b.icon FROM tradesperson_account a INNER JOIN tradesperson_settings s ON a.tradespersonId=s.tradespersonId INNER JOIN quotes q ON a.tradespersonId=q.tradespersonId LEFT JOIN quote_state_cities c ON c.quoteId=q.id LEFT JOIN tradesperson_branding b ON b.tradespersonId=q.tradespersonId WHERE q.category=? AND ((? = '' AND ? = '') OR q.selectPlaces=false OR c.state=? OR JSON_CONTAINS(c.cities, JSON_OBJECT('name', ?))) AND q.archived=false GROUP BY q.id ORDER BY q.id DESC LIMIT ?, ?")
	if err != nil {
		log.Printf("Failed to create select statement %s", err)
		return quotes, err
	}
	defer stmt.Close()

	offset := (page - 1) * int64(PAGE_SIZE)
	rows, err := stmt.Query(category, state, city, state, city, offset, PAGE_SIZE)
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

	stmt, err := db.Prepare("SELECT a.stripeId, q.id, q.tradespersonId, q.quote, q.title, q.description, s.vanityURL, b.icon FROM tradesperson_account a INNER JOIN tradesperson_settings s ON a.tradespersonId=s.tradespersonId INNER JOIN quotes q ON a.tradespersonId=q.tradespersonId LEFT JOIN quote_state_cities c ON c.quoteId=q.id LEFT JOIN tradesperson_branding b ON b.tradespersonId=q.tradespersonId WHERE ((? = '' AND ? = '') OR q.selectPlaces=false OR c.state=? OR JSON_CONTAINS(c.cities, JSON_OBJECT('name', ?))) AND q.archived=false GROUP BY q.id ORDER BY q.id DESC LIMIT ?, ?")
	if err != nil {
		return quotes, err
	}
	defer stmt.Close()

	offset := (page - 1) * int64(PAGE_SIZE)
	rows, err := stmt.Query(state, city, state, city, offset, PAGE_SIZE)
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
		if params.Specialties == nil {
			quotes, err = getSubCategoryQuotes(state, city, *params.Category, *params.SubCategory, page)
			if err != nil {
				log.Printf("%s", err)
				return response
			}
		} else {
			quotes, err = getQuotesWithSpecialties(state, city, *params.Category, *params.SubCategory, *params.Specialties, page)
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
func getQuotesWithSpecialtiesPages(pages float64, state, city, category, subCategory, specialties string) (float64, error) {
	specialtyArry := strings.Split(specialties, ",")
	query := ""
	for _, specialty := range specialtyArry {
		query += "'" + specialty + "',"
	}
	query = query[:len(query)-1]

	db := database.GetConnection()

	stmt, err := db.Prepare("SELECT COUNT(*) FROM quotes q LEFT JOIN quote_state_cities c ON c.quoteId=q.id LEFT JOIN quote_specialties qf ON qf.quoteId=q.id WHERE q.category=? AND q.subcategory=? AND ((? = '' AND ? = '') OR q.selectPlaces=false OR c.state=? OR JSON_CONTAINS(c.cities, JSON_OBJECT('name', ?))) AND  qf.specialty IN (" + query + ") AND q.archived=false")
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

	return pages, nil
}

func getSubCategoryQuotePages(pages float64, state, city, category, subCategory string) (float64, error) {
	db := database.GetConnection()

	stmt, err := db.Prepare("SELECT COUNT(*) FROM quotes q LEFT JOIN quote_state_cities c ON c.quoteId=q.id WHERE q.category=? AND q.subcategory=? AND ((? = '' AND ? = '') OR q.selectPlaces=false OR c.state=? OR JSON_CONTAINS(c.cities, JSON_OBJECT('name', ?))) AND q.archived=false")
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

func getCategoryQuotePages(pages float64, state, city, category string) (float64, error) {
	db := database.GetConnection()

	stmt, err := db.Prepare("SELECT COUNT(*) FROM quotes q LEFT JOIN quote_state_cities c ON c.quoteId=q.id WHERE q.category=? AND ((? = '' AND ? = '') OR q.selectPlaces=false OR c.state=? OR JSON_CONTAINS(c.cities, JSON_OBJECT('name', ?))) AND q.archived=false")
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

func getAllQuotePages(pages float64, state, city string) (float64, error) {

	db := database.GetConnection()

	stmt, err := db.Prepare("SELECT COUNT(*) FROM quotes q LEFT JOIN quote_state_cities c ON c.quoteId=q.id WHERE ((? = '' AND ? = '') OR q.selectPlaces=false OR c.state=? OR JSON_CONTAINS(c.cities, JSON_OBJECT('name', ?))) AND q.archived=false")
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
		if params.Specialties == nil {
			pages, err = getSubCategoryQuotePages(pages, state, city, *params.Category, *params.SubCategory)
			if err != nil {
				log.Printf("%s", err)
				return response
			}
		} else {
			pages, err = getQuotesWithSpecialtiesPages(pages, state, city, *params.Category, *params.SubCategory, *params.Specialties)
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
