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

	"github.com/go-openapi/runtime/middleware"
)

func GetProfileVanityOrIDHandler(params operations.GetProfileVanityOrIDParams) middleware.Responder {
	vanityOrID := params.VanityOrID
	db := database.GetConnection()

	tradesperson := &models.Tradesperson{}
	response := operations.NewGetProfileVanityOrIDOK()
	response.SetPayload(tradesperson)

	stmt, err := db.Prepare("SELECT ta.tradespersonId, ts.number, ts.email, ts.address  FROM tradesperson_account ta INNER JOIN tradesperson_settings ts ON ts.tradespersonId=ta.tradespersonId WHERE ta.tradespersonId=? OR ts.vanityURL=?")
	if err != nil {
		log.Printf("Failed to create select statement %s", err)
		return response
	}
	defer stmt.Close()

	var tradespersonID string
	var number, email, address bool
	row := stmt.QueryRow(vanityOrID, vanityOrID)
	switch err = row.Scan(&tradespersonID, &number, &email, &address); err {
	case sql.ErrNoRows:
		//
	case nil:
		_tradesperson, err := database.GetTradespersonProfile(tradespersonID)
		if err != nil {
			log.Printf("Failed to get tradesperson profile %s", err)
		}
		tradesperson.Name = _tradesperson.Name
		tradesperson.Image = _tradesperson.Image
		tradesperson.Description = _tradesperson.Description
		tradesperson.Address = &models.Address{}
		if address {
			tradesperson.Address.City = _tradesperson.Address.City
			tradesperson.Address.State = _tradesperson.Address.State
			tradesperson.Address.LineOne = _tradesperson.Address.LineOne
			tradesperson.Address.LineTwo = _tradesperson.Address.LineTwo
			tradesperson.Address.ZipCode = _tradesperson.Address.ZipCode
		}
		if number {
			tradesperson.Number = _tradesperson.Number
		}
		if email {
			tradesperson.Email = _tradesperson.Email
		}

		jobs, err := database.GetTradespersonJobs(tradespersonID)
		if err != nil {
			log.Printf("Failed to get tradesperson job count %s", err)
			return response
		}
		tradesperson.Jobs = jobs

		rating, reviews, err := database.GetTradespersonRatingReviews(tradespersonID)
		if err != nil {
			log.Printf("Failed to get tradesperson rating & reviews %s", err)
			return response
		}
		tradesperson.Rating = rating
		tradesperson.Reviews = reviews

		response.SetPayload(tradesperson)
	default:
		log.Printf("Unknown %v", err)
	}
	return response
}

func GetProfileVanityOrIDFixedPricesHandler(params operations.GetProfileVanityOrIDFixedPricesParams) middleware.Responder {
	vanityOrID := params.VanityOrID

	db := database.GetConnection()

	fixedPrices := []*models.Service{}
	response := operations.NewGetProfileVanityOrIDFixedPricesOK().WithPayload(fixedPrices)

	stmt, err := db.Prepare("SELECT fp.tradespersonId, fp.id, fp.priceId, fp.title, fp.price, fp.description, fp.subscription, fp.subInterval FROM fixed_prices fp INNER JOIN tradesperson_settings ts ON ts.tradespersonId=fp.tradespersonId WHERE fp.archived=false AND (fp.tradespersonId=? OR ts.vanityURL=?)")
	if err != nil {
		log.Printf("Failed to create select statement %s", err)
		return response
	}
	defer stmt.Close()

	rows, err := stmt.Query(vanityOrID, vanityOrID)
	if err != nil {
		log.Printf("Failed to execute select statement %s", err)
		return response
	}

	var tradespersonID string
	var id, price int64
	var interval sql.NullString
	for rows.Next() {
		fixedPrice := &models.Service{}
		if err := rows.Scan(&tradespersonID, &id, &fixedPrice.PriceID, &fixedPrice.Title, &price, &fixedPrice.Description, &fixedPrice.Subscription, &interval); err != nil {
			log.Printf("Failed to scan for profile fixed prices, %s", err)
			return response
		}
		if interval.Valid {
			fixedPrice.Interval = interval.String
		}
		strPrice := fmt.Sprintf("%.2f", float64(price)/float64(100.00))
		floatPrice, err := strconv.ParseFloat(strPrice, 64)
		if err != nil {
			log.Printf("Failed to parse float, %v", err)
			return response
		}
		fixedPrice.Price = floatPrice
		fixedPrice.Image, err = internal.GetImage(fixedPrice.PriceID, tradespersonID)
		if err != nil {
			log.Printf("Failed to get fixedPrice image %s", err)
		}
		fixedPrice.AvailableTimeSlots, err = database.GetAvailableTimeSlots(id, fixedPrice.Subscription)
		if err != nil {
			log.Printf("Failed to get timeslots %s", err)
		}

		fixedPrice.Reviews, fixedPrice.Rating, err = database.GetFixedPriceReviewsRating(id)
		if err != nil {
			log.Printf("Failed to get reviews and rating %s", err)
		}

		fixedPrices = append(fixedPrices, fixedPrice)
	}
	response.SetPayload(fixedPrices)

	return response
}

func GetProfileVanityOrIDQuotesHandler(params operations.GetProfileVanityOrIDQuotesParams) middleware.Responder {
	vanityOrID := params.VanityOrID

	db := database.GetConnection()

	quotes := []*models.Service{}
	response := operations.NewGetProfileVanityOrIDQuotesOK().WithPayload(quotes)

	stmt, err := db.Prepare("SELECT q.tradespersonId, q.id, q.quote, q.title FROM quotes q INNER JOIN tradesperson_settings ts ON ts.tradespersonId=q.tradespersonId WHERE q.archived=false AND (q.tradespersonId=? OR ts.vanityURL=?)")
	if err != nil {
		log.Printf("Failed to create select statement %s", err)
		return response
	}
	defer stmt.Close()

	rows, err := stmt.Query(vanityOrID, vanityOrID)
	if err != nil {
		log.Printf("Failed to execute select statement %s", err)
		return response
	}

	var ID int64
	var tradespersonID, quoteID, title string
	for rows.Next() {
		if err := rows.Scan(&tradespersonID, &ID, &quoteID, &title); err != nil {
			log.Printf("Failed to scan for profile quotes, %s", err)
			return response
		}
		quote := &models.Service{}
		quote.Title = title
		quote.QuoteID = quoteID
		quote.Reviews, quote.Rating, err = database.GetQuoteRating(ID)
		if err != nil {
			log.Printf("Failed to get quote reviews and rating %s", err)
		}

		quote.Image, err = internal.GetImage(quoteID, tradespersonID)
		if err != nil {
			log.Printf("Failed to get quote image %s", err)
		}

		quotes = append(quotes, quote)
	}
	response.SetPayload(quotes)

	return response
}
