// Code generated by go-swagger; DO NOT EDIT.

package operations

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/strfmt"
)

// NewGetFixedPricePriceIDParams creates a new GetFixedPricePriceIDParams object
//
// There are no default values defined in the spec.
func NewGetFixedPricePriceIDParams() GetFixedPricePriceIDParams {

	return GetFixedPricePriceIDParams{}
}

// GetFixedPricePriceIDParams contains all the bound params for the get fixed price price ID operation
// typically these are obtained from a http.Request
//
// swagger:parameters GetFixedPricePriceID
type GetFixedPricePriceIDParams struct {

	// HTTP Request Object
	HTTPRequest *http.Request `json:"-"`

	/*Price ID
	  Required: true
	  In: path
	*/
	PriceID string
}

// BindRequest both binds and validates a request, it assumes that complex things implement a Validatable(strfmt.Registry) error interface
// for simple values it will use straight method calls.
//
// To ensure default values, the struct must have been initialized with NewGetFixedPricePriceIDParams() beforehand.
func (o *GetFixedPricePriceIDParams) BindRequest(r *http.Request, route *middleware.MatchedRoute) error {
	var res []error

	o.HTTPRequest = r

	rPriceID, rhkPriceID, _ := route.Params.GetOK("priceId")
	if err := o.bindPriceID(rPriceID, rhkPriceID, route.Formats); err != nil {
		res = append(res, err)
	}
	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

// bindPriceID binds and validates parameter PriceID from path.
func (o *GetFixedPricePriceIDParams) bindPriceID(rawData []string, hasKey bool, formats strfmt.Registry) error {
	var raw string
	if len(rawData) > 0 {
		raw = rawData[len(rawData)-1]
	}

	// Required: true
	// Parameter is provided by construction from the route
	o.PriceID = raw

	return nil
}
