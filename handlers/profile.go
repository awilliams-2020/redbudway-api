package handlers

import (
	"database/sql"
	"fmt"
	"log"
	"redbudway-api/database"
	"redbudway-api/internal"
	"redbudway-api/models"
	"redbudway-api/restapi/operations"
	"redbudway-api/stripe"
	"strconv"

	"github.com/go-openapi/runtime/middleware"
	"github.com/stripe/stripe-go/v72/price"
	"github.com/stripe/stripe-go/v72/product"
)

func GetProfileVanityOrIDHandler(params operations.GetProfileVanityOrIDParams) middleware.Responder {
	vanityOrID := params.VanityOrID
	db := database.GetConnection()

	tradesperson := &models.Tradesperson{}
	response := operations.NewGetProfileVanityOrIDOK()
	response.SetPayload(tradesperson)

	stmt, err := db.Prepare("SELECT ta.tradespersonId, ta.stripeId, ta.name, ta.description, ta.image, IF(ts.number=true, ta.number, '') as number, IF(ts.email=true, ta.email, '') as email, ts.address  FROM tradesperson_account ta INNER JOIN tradesperson_settings ts ON ts.tradespersonId=ta.tradespersonId WHERE ta.tradespersonId=? OR ts.vanityURL=?")
	if err != nil {
		log.Printf("Failed to create select statement %s", err)
		return response
	}
	defer stmt.Close()

	var name, number, email, stripeID, tradespersonID string
	var description, image sql.NullString
	var address bool
	row := stmt.QueryRow(vanityOrID, vanityOrID)
	switch err = row.Scan(&tradespersonID, &stripeID, &name, &description, &image, &number, &email, &address); err {
	case sql.ErrNoRows:
		//
	case nil:
		if description.Valid {
			tradesperson.Description = description.String
		}
		if image.Valid {
			tradesperson.Image = image.String
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

		tradesperson.Name = name
		tradesperson.Number = number
		tradesperson.Email = email
		if address {
			stripe, err := stripe.GetConnectAccount(stripeID)
			if err != nil {
				log.Print("Failed to get stripe account for tradesperson with ID %s", tradespersonID)
				return response
			}

			tradesperson.Address = &models.Address{}
			tradesperson.Address.City = stripe.BusinessProfile.SupportAddress.City
			tradesperson.Address.State = stripe.BusinessProfile.SupportAddress.State
			tradesperson.Address.LineOne = stripe.BusinessProfile.SupportAddress.Line1
			tradesperson.Address.LineTwo = stripe.BusinessProfile.SupportAddress.Line2
			tradesperson.Address.ZipCode = stripe.BusinessProfile.SupportAddress.PostalCode
		}
		response.SetPayload(tradesperson)
	default:
		log.Printf("Unknown %v", err)
	}
	return response
}

func GetProfileVanityOrIDFixedPricesHandler(params operations.GetProfileVanityOrIDFixedPricesParams) middleware.Responder {
	vanityOrID := params.VanityOrID
	state := params.State
	city := params.City

	db := database.GetConnection()

	fixedPrices := []*models.Service{}
	response := operations.NewGetProfileVanityOrIDFixedPricesOK().WithPayload(fixedPrices)

	stmt, err := db.Prepare("SELECT fp.id, fp.priceId, fp.subscription, fp.subInterval, fp.selectPlaces, fpsc.cities FROM fixed_prices fp INNER JOIN tradesperson_settings ts ON ts.tradespersonId=fp.tradespersonId LEFT JOIN fixed_price_state_cities fpsc ON fpsc.fixedPriceId=fp.id WHERE (fp.selectPlaces=false OR fpsc.state=?) AND fp.archived=false AND (fp.tradespersonId=? OR ts.vanityURL=?)")
	if err != nil {
		log.Printf("Failed to create select statement %s", err)
		return response
	}
	defer stmt.Close()

	rows, err := stmt.Query(state, vanityOrID, vanityOrID)
	if err != nil {
		log.Printf("Failed to execute select statement %s", err)
		return response
	}

	var id int64
	var interval, citiesJson sql.NullString
	var subscription, selectPlaces bool
	var priceID string
	for rows.Next() {
		if err := rows.Scan(&id, &priceID, &subscription, &interval, &selectPlaces, &citiesJson); err != nil {
			log.Printf("Failed to scan for profile fixed prices, %s", err)
			return response
		}

		fixedPrice := &models.Service{}
		fixedPrice.Subscription = subscription
		if interval.Valid {
			fixedPrice.Interval = interval.String
		}
		stripePrice, err := price.Get(priceID, nil)
		if err != nil {
			log.Printf("Failed to get stripe price, %v", err)
			return response
		}
		stripeProduct, err := product.Get(stripePrice.Product.ID, nil)
		if err != nil {
			log.Printf("Failed to get stripe product, %v", err)
			return response
		}
		fixedPrice.PriceID = priceID
		strPrice := fmt.Sprintf("%.2f", stripePrice.UnitAmountDecimal/float64(100.00))
		floatPrice, err := strconv.ParseFloat(strPrice, 64)
		if err != nil {
			log.Printf("Failed to parse float, %v", err)
			return response
		}
		fixedPrice.Price = floatPrice
		fixedPrice.Title = stripeProduct.Name
		fixedPrice.Image = stripeProduct.Images[0]

		if !fixedPrice.Subscription {
			fixedPrice.AvailableTimeSlots, err = database.GetAvailableTimeSlots(id, subscription)
			if err != nil {
				log.Printf("Failed to get timeslots %s", err)
			}
		}

		fixedPrice.Reviews, fixedPrice.Rating, err = database.GetFixedPriceReviewsRating(id)
		if err != nil {
			log.Printf("Failed to get reviews and rating %s", err)
		}

		if selectPlaces {
			if citiesJson.Valid {
				cityExist, _ := internal.SelectedCities(citiesJson.String, city)
				if cityExist {
					fixedPrices = append(fixedPrices, fixedPrice)
				}
			}
		} else {
			fixedPrices = append(fixedPrices, fixedPrice)
		}
		response.SetPayload(fixedPrices)
	}

	return response
}

func GetProfileVanityOrIDQuotesHandler(params operations.GetProfileVanityOrIDQuotesParams) middleware.Responder {
	vanityOrID := params.VanityOrID
	state := params.State
	city := params.City

	db := database.GetConnection()

	quotes := []*models.Service{}
	response := operations.NewGetProfileVanityOrIDQuotesOK().WithPayload(quotes)

	stmt, err := db.Prepare("SELECT q.id, q.quote, q.title, q.selectPlaces, qsc.cities FROM quotes q INNER JOIN tradesperson_settings ts ON ts.tradespersonId=q.tradespersonId LEFT JOIN quote_state_cities qsc ON qsc.quoteId=q.id WHERE (q.selectPlaces=false OR qsc.state=?) AND q.archived=false AND (q.tradespersonId=? OR ts.vanityURL=?)")
	if err != nil {
		log.Printf("Failed to create select statement %s", err)
		return response
	}
	defer stmt.Close()

	rows, err := stmt.Query(state, vanityOrID, vanityOrID)
	if err != nil {
		log.Printf("Failed to execute select statement %s", err)
		return response
	}

	var ID int64
	var citiesJson sql.NullString
	var quoteID, title string
	var selectPlaces bool
	for rows.Next() {
		if err := rows.Scan(&ID, &quoteID, &title, &selectPlaces, &citiesJson); err != nil {
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

		quote.Image, err = database.GetQuoteImage(ID)
		if err != nil {
			log.Printf("Failed to get quote image %s", err)
		}

		if selectPlaces {
			if citiesJson.Valid {
				cityExist, _ := internal.SelectedCities(citiesJson.String, city)
				if cityExist {
					quotes = append(quotes, quote)
				}
			}
		} else {
			quotes = append(quotes, quote)
		}
		response.SetPayload(quotes)
	}

	return response
}
