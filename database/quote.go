package database

import (
	"database/sql"
	"log"
	"os"
	"redbudway-api/internal"
	"redbudway-api/models"
	"redbudway-api/restapi/operations"
)

func GetQuoteServiceDetails(quoteID string) (*models.ServiceDetails, *operations.GetQuoteQuoteIDOKBodyBusiness, error) {
	quote := &models.ServiceDetails{}
	business := &operations.GetQuoteQuoteIDOKBodyBusiness{}

	stmt, err := db.Prepare("SELECT tp.name, tp.tradespersonId, ts.vanityURL, q.id, q.title, q.description, q.category, q.subcategory, q.selectPlaces FROM quotes q INNER JOIN tradesperson_profile tp ON tp.tradespersonId=q.tradespersonId INNER JOIN tradesperson_settings ts ON ts.tradespersonId=q.tradespersonId LEFT JOIN quote_state_cities qsc ON qsc.quoteId=q.id WHERE q.archived=false AND q.quote=?")
	if err != nil {
		return quote, business, err
	}
	defer stmt.Close()

	row := stmt.QueryRow(quoteID)
	var ID int64
	var vanityURL sql.NullString
	var name, tradespersonID, title, description, category, subCategory string
	var selectPlaces bool
	switch err = row.Scan(&name, &tradespersonID, &vanityURL, &ID, &title, &description, &category, &subCategory, &selectPlaces); err {
	case sql.ErrNoRows:
		return quote, business, err
	case nil:
		business.Name = name
		business.VanityURL = vanityURL.String
		business.TradespersonID = tradespersonID

		quote.Category = &category
		quote.SubCategory = subCategory
		quote.Title = &title
		quote.Description = &description
		quote.SelectPlaces = &selectPlaces

		var err error
		quote.Reviews, quote.Rating, err = GetQuoteRating(ID)
		if err != nil {
			log.Printf("Failed to get quote reviews and rating %s", err)
		}
		quote.Images, err = internal.GetImages(quoteID, business.TradespersonID)
		if err != nil {
			log.Printf("Failed to get quote image %s", err)
		}
		if len(quote.Images) == 0 {
			url := "https://" + os.Getenv("SUBDOMAIN") + "redbudway.com/assets/images/placeholder.svg"
			quote.Images = append(quote.Images, url)
		}
		quote.Specialties, err = GetSpecialties(ID)
		if err != nil {
			return quote, business, err
		}
		quote.StatesAndCities, err = GetQuoteStatesAndCities(ID)
		if err != nil {
			return quote, business, err
		}

	default:
		return quote, business, err
	}

	return quote, business, nil
}
