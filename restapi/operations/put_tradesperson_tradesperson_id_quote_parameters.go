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

	"redbudway-api/models"
)

// NewPutTradespersonTradespersonIDQuoteParams creates a new PutTradespersonTradespersonIDQuoteParams object
//
// There are no default values defined in the spec.
func NewPutTradespersonTradespersonIDQuoteParams() PutTradespersonTradespersonIDQuoteParams {

	return PutTradespersonTradespersonIDQuoteParams{}
}

// PutTradespersonTradespersonIDQuoteParams contains all the bound params for the put tradesperson tradesperson ID quote operation
// typically these are obtained from a http.Request
//
// swagger:parameters PutTradespersonTradespersonIDQuote
type PutTradespersonTradespersonIDQuoteParams struct {

	// HTTP Request Object
	HTTPRequest *http.Request `json:"-"`

	/*The quote to create.
	  In: body
	*/
	Quote *models.ServiceDetails
	/*
	  Required: true
	  In: path
	*/
	TradespersonID string
}

// BindRequest both binds and validates a request, it assumes that complex things implement a Validatable(strfmt.Registry) error interface
// for simple values it will use straight method calls.
//
// To ensure default values, the struct must have been initialized with NewPutTradespersonTradespersonIDQuoteParams() beforehand.
func (o *PutTradespersonTradespersonIDQuoteParams) BindRequest(r *http.Request, route *middleware.MatchedRoute) error {
	var res []error

	o.HTTPRequest = r

	if runtime.HasBody(r) {
		defer r.Body.Close()
		var body models.ServiceDetails
		if err := route.Consumer.Consume(r.Body, &body); err != nil {
			res = append(res, errors.NewParseError("quote", "body", "", err))
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
				o.Quote = &body
			}
		}
	}

	rTradespersonID, rhkTradespersonID, _ := route.Params.GetOK("tradespersonId")
	if err := o.bindTradespersonID(rTradespersonID, rhkTradespersonID, route.Formats); err != nil {
		res = append(res, err)
	}
	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

// bindTradespersonID binds and validates parameter TradespersonID from path.
func (o *PutTradespersonTradespersonIDQuoteParams) bindTradespersonID(rawData []string, hasKey bool, formats strfmt.Registry) error {
	var raw string
	if len(rawData) > 0 {
		raw = rawData[len(rawData)-1]
	}

	// Required: true
	// Parameter is provided by construction from the route
	o.TradespersonID = raw

	return nil
}
