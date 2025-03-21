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

// NewPostCustomerCustomerIDQuoteQuoteIDRequestParams creates a new PostCustomerCustomerIDQuoteQuoteIDRequestParams object
//
// There are no default values defined in the spec.
func NewPostCustomerCustomerIDQuoteQuoteIDRequestParams() PostCustomerCustomerIDQuoteQuoteIDRequestParams {

	return PostCustomerCustomerIDQuoteQuoteIDRequestParams{}
}

// PostCustomerCustomerIDQuoteQuoteIDRequestParams contains all the bound params for the post customer customer ID quote quote ID request operation
// typically these are obtained from a http.Request
//
// swagger:parameters PostCustomerCustomerIDQuoteQuoteIDRequest
type PostCustomerCustomerIDQuoteQuoteIDRequestParams struct {

	// HTTP Request Object
	HTTPRequest *http.Request `json:"-"`

	/*Customer ID
	  Required: true
	  In: path
	*/
	CustomerID string
	/*Quote ID
	  Required: true
	  In: path
	*/
	QuoteID string
	/*The quote to request
	  In: body
	*/
	Request PostCustomerCustomerIDQuoteQuoteIDRequestBody
}

// BindRequest both binds and validates a request, it assumes that complex things implement a Validatable(strfmt.Registry) error interface
// for simple values it will use straight method calls.
//
// To ensure default values, the struct must have been initialized with NewPostCustomerCustomerIDQuoteQuoteIDRequestParams() beforehand.
func (o *PostCustomerCustomerIDQuoteQuoteIDRequestParams) BindRequest(r *http.Request, route *middleware.MatchedRoute) error {
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
		var body PostCustomerCustomerIDQuoteQuoteIDRequestBody
		if err := route.Consumer.Consume(r.Body, &body); err != nil {
			res = append(res, errors.NewParseError("request", "body", "", err))
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
				o.Request = body
			}
		}
	}
	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

// bindCustomerID binds and validates parameter CustomerID from path.
func (o *PostCustomerCustomerIDQuoteQuoteIDRequestParams) bindCustomerID(rawData []string, hasKey bool, formats strfmt.Registry) error {
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
func (o *PostCustomerCustomerIDQuoteQuoteIDRequestParams) bindQuoteID(rawData []string, hasKey bool, formats strfmt.Registry) error {
	var raw string
	if len(rawData) > 0 {
		raw = rawData[len(rawData)-1]
	}

	// Required: true
	// Parameter is provided by construction from the route
	o.QuoteID = raw

	return nil
}
