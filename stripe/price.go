package stripe

import (
	"redbudway-api/models"

	"github.com/stripe/stripe-go/v72"
	"github.com/stripe/stripe-go/v72/price"
)

func CreatePrice(fixedPrice *models.ServiceDetails) (*stripe.Price, error) {

	product, err := CreateProduct(fixedPrice)
	if err != nil {
		return &stripe.Price{}, err
	}

	priceInteger := int64(fixedPrice.Price * float64(100.00))
	priceParams := &stripe.PriceParams{
		Currency:   stripe.String(string(stripe.CurrencyUSD)),
		Product:    stripe.String(product.ID),
		UnitAmount: stripe.Int64(priceInteger),
	}
	if fixedPrice.Subscription && (fixedPrice.Interval == "week" || fixedPrice.Interval == "month" || fixedPrice.Interval == "year") {
		priceParams.Recurring = &stripe.PriceRecurringParams{
			Interval: stripe.String(fixedPrice.Interval),
		}
	}
	sPrice, err := price.New(priceParams)
	if err != nil {
		return sPrice, err
	}

	return sPrice, nil
}
