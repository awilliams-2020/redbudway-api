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

// NewPostTradespersonTradespersonIDBillingInvoiceInvoiceIDFinalizeParams creates a new PostTradespersonTradespersonIDBillingInvoiceInvoiceIDFinalizeParams object
//
// There are no default values defined in the spec.
func NewPostTradespersonTradespersonIDBillingInvoiceInvoiceIDFinalizeParams() PostTradespersonTradespersonIDBillingInvoiceInvoiceIDFinalizeParams {

	return PostTradespersonTradespersonIDBillingInvoiceInvoiceIDFinalizeParams{}
}

// PostTradespersonTradespersonIDBillingInvoiceInvoiceIDFinalizeParams contains all the bound params for the post tradesperson tradesperson ID billing invoice invoice ID finalize operation
// typically these are obtained from a http.Request
//
// swagger:parameters PostTradespersonTradespersonIDBillingInvoiceInvoiceIDFinalize
type PostTradespersonTradespersonIDBillingInvoiceInvoiceIDFinalizeParams struct {

	// HTTP Request Object
	HTTPRequest *http.Request `json:"-"`

	/*Invoice ID
	  Required: true
	  In: path
	*/
	InvoiceID string
	/*Tradesperson ID
	  Required: true
	  In: path
	*/
	TradespersonID string
}

// BindRequest both binds and validates a request, it assumes that complex things implement a Validatable(strfmt.Registry) error interface
// for simple values it will use straight method calls.
//
// To ensure default values, the struct must have been initialized with NewPostTradespersonTradespersonIDBillingInvoiceInvoiceIDFinalizeParams() beforehand.
func (o *PostTradespersonTradespersonIDBillingInvoiceInvoiceIDFinalizeParams) BindRequest(r *http.Request, route *middleware.MatchedRoute) error {
	var res []error

	o.HTTPRequest = r

	rInvoiceID, rhkInvoiceID, _ := route.Params.GetOK("invoiceId")
	if err := o.bindInvoiceID(rInvoiceID, rhkInvoiceID, route.Formats); err != nil {
		res = append(res, err)
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

// bindInvoiceID binds and validates parameter InvoiceID from path.
func (o *PostTradespersonTradespersonIDBillingInvoiceInvoiceIDFinalizeParams) bindInvoiceID(rawData []string, hasKey bool, formats strfmt.Registry) error {
	var raw string
	if len(rawData) > 0 {
		raw = rawData[len(rawData)-1]
	}

	// Required: true
	// Parameter is provided by construction from the route
	o.InvoiceID = raw

	return nil
}

// bindTradespersonID binds and validates parameter TradespersonID from path.
func (o *PostTradespersonTradespersonIDBillingInvoiceInvoiceIDFinalizeParams) bindTradespersonID(rawData []string, hasKey bool, formats strfmt.Registry) error {
	var raw string
	if len(rawData) > 0 {
		raw = rawData[len(rawData)-1]
	}

	// Required: true
	// Parameter is provided by construction from the route
	o.TradespersonID = raw

	return nil
}
