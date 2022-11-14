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

// NewGetFixedPricePagesParams creates a new GetFixedPricePagesParams object
//
// There are no default values defined in the spec.
func NewGetFixedPricePagesParams() GetFixedPricePagesParams {

	return GetFixedPricePagesParams{}
}

// GetFixedPricePagesParams contains all the bound params for the get fixed price pages operation
// typically these are obtained from a http.Request
//
// swagger:parameters GetFixedPricePages
type GetFixedPricePagesParams struct {

	// HTTP Request Object
	HTTPRequest *http.Request `json:"-"`

	/*
	  In: query
	*/
	Category *string
	/*
	  Required: true
	  In: query
	*/
	City string
	/*
	  In: query
	*/
	Filters *string
	/*
	  Required: true
	  In: query
	*/
	State string
	/*
	  In: query
	*/
	SubCategory *string
}

// BindRequest both binds and validates a request, it assumes that complex things implement a Validatable(strfmt.Registry) error interface
// for simple values it will use straight method calls.
//
// To ensure default values, the struct must have been initialized with NewGetFixedPricePagesParams() beforehand.
func (o *GetFixedPricePagesParams) BindRequest(r *http.Request, route *middleware.MatchedRoute) error {
	var res []error

	o.HTTPRequest = r

	qs := runtime.Values(r.URL.Query())

	qCategory, qhkCategory, _ := qs.GetOK("category")
	if err := o.bindCategory(qCategory, qhkCategory, route.Formats); err != nil {
		res = append(res, err)
	}

	qCity, qhkCity, _ := qs.GetOK("city")
	if err := o.bindCity(qCity, qhkCity, route.Formats); err != nil {
		res = append(res, err)
	}

	qFilters, qhkFilters, _ := qs.GetOK("filters")
	if err := o.bindFilters(qFilters, qhkFilters, route.Formats); err != nil {
		res = append(res, err)
	}

	qState, qhkState, _ := qs.GetOK("state")
	if err := o.bindState(qState, qhkState, route.Formats); err != nil {
		res = append(res, err)
	}

	qSubCategory, qhkSubCategory, _ := qs.GetOK("subCategory")
	if err := o.bindSubCategory(qSubCategory, qhkSubCategory, route.Formats); err != nil {
		res = append(res, err)
	}
	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

// bindCategory binds and validates parameter Category from query.
func (o *GetFixedPricePagesParams) bindCategory(rawData []string, hasKey bool, formats strfmt.Registry) error {
	var raw string
	if len(rawData) > 0 {
		raw = rawData[len(rawData)-1]
	}

	// Required: false
	// AllowEmptyValue: false

	if raw == "" { // empty values pass all other validations
		return nil
	}
	o.Category = &raw

	return nil
}

// bindCity binds and validates parameter City from query.
func (o *GetFixedPricePagesParams) bindCity(rawData []string, hasKey bool, formats strfmt.Registry) error {
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

// bindFilters binds and validates parameter Filters from query.
func (o *GetFixedPricePagesParams) bindFilters(rawData []string, hasKey bool, formats strfmt.Registry) error {
	var raw string
	if len(rawData) > 0 {
		raw = rawData[len(rawData)-1]
	}

	// Required: false
	// AllowEmptyValue: false

	if raw == "" { // empty values pass all other validations
		return nil
	}
	o.Filters = &raw

	return nil
}

// bindState binds and validates parameter State from query.
func (o *GetFixedPricePagesParams) bindState(rawData []string, hasKey bool, formats strfmt.Registry) error {
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

// bindSubCategory binds and validates parameter SubCategory from query.
func (o *GetFixedPricePagesParams) bindSubCategory(rawData []string, hasKey bool, formats strfmt.Registry) error {
	var raw string
	if len(rawData) > 0 {
		raw = rawData[len(rawData)-1]
	}

	// Required: false
	// AllowEmptyValue: false

	if raw == "" { // empty values pass all other validations
		return nil
	}
	o.SubCategory = &raw

	return nil
}
