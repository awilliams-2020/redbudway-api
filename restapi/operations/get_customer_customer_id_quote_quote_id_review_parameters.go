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

// NewGetCustomerCustomerIDQuoteQuoteIDReviewParams creates a new GetCustomerCustomerIDQuoteQuoteIDReviewParams object
//
// There are no default values defined in the spec.
func NewGetCustomerCustomerIDQuoteQuoteIDReviewParams() GetCustomerCustomerIDQuoteQuoteIDReviewParams {

	return GetCustomerCustomerIDQuoteQuoteIDReviewParams{}
}

// GetCustomerCustomerIDQuoteQuoteIDReviewParams contains all the bound params for the get customer customer ID quote quote ID review operation
// typically these are obtained from a http.Request
//
// swagger:parameters GetCustomerCustomerIDQuoteQuoteIDReview
type GetCustomerCustomerIDQuoteQuoteIDReviewParams struct {

	// HTTP Request Object
	HTTPRequest *http.Request `json:"-"`

	/*Customer ID
	  Required: true
	  In: path
	*/
	CustomerID string
	/*Quote with review
	  Required: true
	  In: path
	*/
	QuoteID string
}

// BindRequest both binds and validates a request, it assumes that complex things implement a Validatable(strfmt.Registry) error interface
// for simple values it will use straight method calls.
//
// To ensure default values, the struct must have been initialized with NewGetCustomerCustomerIDQuoteQuoteIDReviewParams() beforehand.
func (o *GetCustomerCustomerIDQuoteQuoteIDReviewParams) BindRequest(r *http.Request, route *middleware.MatchedRoute) error {
	var res []error

	o.HTTPRequest = r

	rCustomerID, rhkCustomerID, _ := route.Params.GetOK("customerId")
	if err := o.bindCustomerID(rCustomerID, rhkCustomerID, route.Formats); err != nil {
		res = append(res, err)
	}

	rQuoteID, rhkQuoteID, _ := route.Params.GetOK("quoteId")
	if err := o.bindQuoteID(rQuoteID, rhkQuoteID, route.Formats); err != nil {
		res = append(res, err)
	}
	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

// bindCustomerID binds and validates parameter CustomerID from path.
func (o *GetCustomerCustomerIDQuoteQuoteIDReviewParams) bindCustomerID(rawData []string, hasKey bool, formats strfmt.Registry) error {
	var raw string
	if len(rawData) > 0 {
		raw = rawData[len(rawData)-1]
	}

	// Required: true
	// Parameter is provided by construction from the route
	o.CustomerID = raw

	return nil
}

// bindQuoteID binds and validates parameter QuoteID from path.
func (o *GetCustomerCustomerIDQuoteQuoteIDReviewParams) bindQuoteID(rawData []string, hasKey bool, formats strfmt.Registry) error {
	var raw string
	if len(rawData) > 0 {
		raw = rawData[len(rawData)-1]
	}

	// Required: true
	// Parameter is provided by construction from the route
	o.QuoteID = raw

	return nil
}
