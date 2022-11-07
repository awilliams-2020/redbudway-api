// Code generated by go-swagger; DO NOT EDIT.

package operations

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"
	"net/http"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/validate"
)

// NewPostCustomerCustomerIDQuoteQuoteIDReviewParams creates a new PostCustomerCustomerIDQuoteQuoteIDReviewParams object
//
// There are no default values defined in the spec.
func NewPostCustomerCustomerIDQuoteQuoteIDReviewParams() PostCustomerCustomerIDQuoteQuoteIDReviewParams {

	return PostCustomerCustomerIDQuoteQuoteIDReviewParams{}
}

// PostCustomerCustomerIDQuoteQuoteIDReviewParams contains all the bound params for the post customer customer ID quote quote ID review operation
// typically these are obtained from a http.Request
//
// swagger:parameters PostCustomerCustomerIDQuoteQuoteIDReview
type PostCustomerCustomerIDQuoteQuoteIDReviewParams struct {

	// HTTP Request Object
	HTTPRequest *http.Request `json:"-"`

	/*Customer ID
	  Required: true
	  In: path
	*/
	CustomerID string
	/*Quote ID with review
	  Required: true
	  In: path
	*/
	QuoteID string
	/*The quote review to create.
	  In: body
	*/
	Review PostCustomerCustomerIDQuoteQuoteIDReviewBody
}

// BindRequest both binds and validates a request, it assumes that complex things implement a Validatable(strfmt.Registry) error interface
// for simple values it will use straight method calls.
//
// To ensure default values, the struct must have been initialized with NewPostCustomerCustomerIDQuoteQuoteIDReviewParams() beforehand.
func (o *PostCustomerCustomerIDQuoteQuoteIDReviewParams) BindRequest(r *http.Request, route *middleware.MatchedRoute) error {
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

	if runtime.HasBody(r) {
		defer r.Body.Close()
		var body PostCustomerCustomerIDQuoteQuoteIDReviewBody
		if err := route.Consumer.Consume(r.Body, &body); err != nil {
			res = append(res, errors.NewParseError("review", "body", "", err))
		} else {
			// validate body object
			if err := body.Validate(route.Formats); err != nil {
				res = append(res, err)
			}

			ctx := validate.WithOperationRequest(context.Background())
			if err := body.ContextValidate(ctx, route.Formats); err != nil {
				res = append(res, err)
			}

			if len(res) == 0 {
				o.Review = body
			}
		}
	}
	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

// bindCustomerID binds and validates parameter CustomerID from path.
func (o *PostCustomerCustomerIDQuoteQuoteIDReviewParams) bindCustomerID(rawData []string, hasKey bool, formats strfmt.Registry) error {
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
func (o *PostCustomerCustomerIDQuoteQuoteIDReviewParams) bindQuoteID(rawData []string, hasKey bool, formats strfmt.Registry) error {
	var raw string
	if len(rawData) > 0 {
		raw = rawData[len(rawData)-1]
	}

	// Required: true
	// Parameter is provided by construction from the route
	o.QuoteID = raw

	return nil
}
