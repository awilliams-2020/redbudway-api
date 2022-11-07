// Code generated by go-swagger; DO NOT EDIT.

package operations

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/validate"
)

// NewGetProfileVanityOrIDQuotesParams creates a new GetProfileVanityOrIDQuotesParams object
//
// There are no default values defined in the spec.
func NewGetProfileVanityOrIDQuotesParams() GetProfileVanityOrIDQuotesParams {

	return GetProfileVanityOrIDQuotesParams{}
}

// GetProfileVanityOrIDQuotesParams contains all the bound params for the get profile vanity or ID quotes operation
// typically these are obtained from a http.Request
//
// swagger:parameters GetProfileVanityOrIDQuotes
type GetProfileVanityOrIDQuotesParams struct {

	// HTTP Request Object
	HTTPRequest *http.Request `json:"-"`

	/*
	  Required: true
	  In: query
	*/
	City string
	/*
	  Required: true
	  In: query
	*/
	State string
	/*Tradesperson vanity URL or ID
	  Required: true
	  In: path
	*/
	VanityOrID string
}

// BindRequest both binds and validates a request, it assumes that complex things implement a Validatable(strfmt.Registry) error interface
// for simple values it will use straight method calls.
//
// To ensure default values, the struct must have been initialized with NewGetProfileVanityOrIDQuotesParams() beforehand.
func (o *GetProfileVanityOrIDQuotesParams) BindRequest(r *http.Request, route *middleware.MatchedRoute) error {
	var res []error

	o.HTTPRequest = r

	qs := runtime.Values(r.URL.Query())

	qCity, qhkCity, _ := qs.GetOK("city")
	if err := o.bindCity(qCity, qhkCity, route.Formats); err != nil {
		res = append(res, err)
	}

	qState, qhkState, _ := qs.GetOK("state")
	if err := o.bindState(qState, qhkState, route.Formats); err != nil {
		res = append(res, err)
	}

	rVanityOrID, rhkVanityOrID, _ := route.Params.GetOK("vanityOrId")
	if err := o.bindVanityOrID(rVanityOrID, rhkVanityOrID, route.Formats); err != nil {
		res = append(res, err)
	}
	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

// bindCity binds and validates parameter City from query.
func (o *GetProfileVanityOrIDQuotesParams) bindCity(rawData []string, hasKey bool, formats strfmt.Registry) error {
	if !hasKey {
		return errors.Required("city", "query", rawData)
	}
	var raw string
	if len(rawData) > 0 {
		raw = rawData[len(rawData)-1]
	}

	// Required: true
	// AllowEmptyValue: false

	if err := validate.RequiredString("city", "query", raw); err != nil {
		return err
	}
	o.City = raw

	return nil
}

// bindState binds and validates parameter State from query.
func (o *GetProfileVanityOrIDQuotesParams) bindState(rawData []string, hasKey bool, formats strfmt.Registry) error {
	if !hasKey {
		return errors.Required("state", "query", rawData)
	}
	var raw string
	if len(rawData) > 0 {
		raw = rawData[len(rawData)-1]
	}

	// Required: true
	// AllowEmptyValue: false

	if err := validate.RequiredString("state", "query", raw); err != nil {
		return err
	}
	o.State = raw

	return nil
}

// bindVanityOrID binds and validates parameter VanityOrID from path.
func (o *GetProfileVanityOrIDQuotesParams) bindVanityOrID(rawData []string, hasKey bool, formats strfmt.Registry) error {
	var raw string
	if len(rawData) > 0 {
		raw = rawData[len(rawData)-1]
	}

	// Required: true
	// Parameter is provided by construction from the route
	o.VanityOrID = raw

	return nil
}
