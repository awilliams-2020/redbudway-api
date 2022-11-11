package stripe

import (
	"redbudway-api/models"

	"github.com/stripe/stripe-go/v72"
	"github.com/stripe/stripe-go/v72/product"
)

func CreateProduct(fixedPrice *models.ServiceDetails) (*stripe.Product, error) {
	productParams := &stripe.ProductParams{
		Name:        stripe.String(*fixedPrice.Title),
		Description: stripe.String(*fixedPrice.Description),
	}
	product, _ := product.New(productParams)

	return product, nil
}

func UpdateProduct(images []*string, fixedPrice *models.ServiceDetails, price *stripe.Price) error {
	params := &stripe.ProductParams{}
	params.Name = fixedPrice.Title
	params.Description = fixedPrice.Description
	params.Images = images
	_, err := product.Update(
		price.Product.ID,
		params,
	)
	if err != nil {
		return err
	}
	return nil
}
