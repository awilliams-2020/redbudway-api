// Code generated by go-swagger; DO NOT EDIT.

package operations

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"context"
	"net/http"

	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
)

// PutTradespersonTradespersonIDBillingQuoteQuoteIDInvoiceInvoiceIDHandlerFunc turns a function with the right signature into a put tradesperson tradesperson ID billing quote quote ID invoice invoice ID handler
type PutTradespersonTradespersonIDBillingQuoteQuoteIDInvoiceInvoiceIDHandlerFunc func(PutTradespersonTradespersonIDBillingQuoteQuoteIDInvoiceInvoiceIDParams, interface{}) middleware.Responder

// Handle executing the request and returning a response
func (fn PutTradespersonTradespersonIDBillingQuoteQuoteIDInvoiceInvoiceIDHandlerFunc) Handle(params PutTradespersonTradespersonIDBillingQuoteQuoteIDInvoiceInvoiceIDParams, principal interface{}) middleware.Responder {
	return fn(params, principal)
}

// PutTradespersonTradespersonIDBillingQuoteQuoteIDInvoiceInvoiceIDHandler interface for that can handle valid put tradesperson tradesperson ID billing quote quote ID invoice invoice ID params
type PutTradespersonTradespersonIDBillingQuoteQuoteIDInvoiceInvoiceIDHandler interface {
	Handle(PutTradespersonTradespersonIDBillingQuoteQuoteIDInvoiceInvoiceIDParams, interface{}) middleware.Responder
}

// NewPutTradespersonTradespersonIDBillingQuoteQuoteIDInvoiceInvoiceID creates a new http.Handler for the put tradesperson tradesperson ID billing quote quote ID invoice invoice ID operation
func NewPutTradespersonTradespersonIDBillingQuoteQuoteIDInvoiceInvoiceID(ctx *middleware.Context, handler PutTradespersonTradespersonIDBillingQuoteQuoteIDInvoiceInvoiceIDHandler) *PutTradespersonTradespersonIDBillingQuoteQuoteIDInvoiceInvoiceID {
	return &PutTradespersonTradespersonIDBillingQuoteQuoteIDInvoiceInvoiceID{Context: ctx, Handler: handler}
}

/* PutTradespersonTradespersonIDBillingQuoteQuoteIDInvoiceInvoiceID swagger:route PUT /tradesperson/{tradespersonId}/billing/quote/{quoteId}/invoice/{invoiceId} putTradespersonTradespersonIdBillingQuoteQuoteIdInvoiceInvoiceId

PutTradespersonTradespersonIDBillingQuoteQuoteIDInvoiceInvoiceID put tradesperson tradesperson ID billing quote quote ID invoice invoice ID API

*/
type PutTradespersonTradespersonIDBillingQuoteQuoteIDInvoiceInvoiceID struct {
	Context *middleware.Context
	Handler PutTradespersonTradespersonIDBillingQuoteQuoteIDInvoiceInvoiceIDHandler
}

func (o *PutTradespersonTradespersonIDBillingQuoteQuoteIDInvoiceInvoiceID) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, rCtx, _ := o.Context.RouteInfo(r)
	if rCtx != nil {
		*r = *rCtx
	}
	var Params = NewPutTradespersonTradespersonIDBillingQuoteQuoteIDInvoiceInvoiceIDParams()
	uprinc, aCtx, err := o.Context.Authorize(r, route)
	if err != nil {
		o.Context.Respond(rw, r, route.Produces, route, err)
		return
	}
	if aCtx != nil {
		*r = *aCtx
	}
	var principal interface{}
	if uprinc != nil {
		principal = uprinc.(interface{}) // this is really a interface{}, I promise
	}

	if err := o.Context.BindValidRequest(r, route, &Params); err != nil { // bind params
		o.Context.Respond(rw, r, route.Produces, route, err)
		return
	}

	res := o.Handler.Handle(Params, principal) // actually handle the request
	o.Context.Respond(rw, r, route.Produces, route, res)

}

// PutTradespersonTradespersonIDBillingQuoteQuoteIDInvoiceInvoiceIDBody put tradesperson tradesperson ID billing quote quote ID invoice invoice ID body
//
// swagger:model PutTradespersonTradespersonIDBillingQuoteQuoteIDInvoiceInvoiceIDBody
type PutTradespersonTradespersonIDBillingQuoteQuoteIDInvoiceInvoiceIDBody struct {

	// description
	Description string `json:"description,omitempty"`
}

// Validate validates this put tradesperson tradesperson ID billing quote quote ID invoice invoice ID body
func (o *PutTradespersonTradespersonIDBillingQuoteQuoteIDInvoiceInvoiceIDBody) Validate(formats strfmt.Registry) error {
	return nil
}

// ContextValidate validates this put tradesperson tradesperson ID billing quote quote ID invoice invoice ID body based on context it is used
func (o *PutTradespersonTradespersonIDBillingQuoteQuoteIDInvoiceInvoiceIDBody) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (o *PutTradespersonTradespersonIDBillingQuoteQuoteIDInvoiceInvoiceIDBody) MarshalBinary() ([]byte, error) {
	if o == nil {
		return nil, nil
	}
	return swag.WriteJSON(o)
}

// UnmarshalBinary interface implementation
func (o *PutTradespersonTradespersonIDBillingQuoteQuoteIDInvoiceInvoiceIDBody) UnmarshalBinary(b []byte) error {
	var res PutTradespersonTradespersonIDBillingQuoteQuoteIDInvoiceInvoiceIDBody
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*o = res
	return nil
}

// PutTradespersonTradespersonIDBillingQuoteQuoteIDInvoiceInvoiceIDOKBody put tradesperson tradesperson ID billing quote quote ID invoice invoice ID o k body
//
// swagger:model PutTradespersonTradespersonIDBillingQuoteQuoteIDInvoiceInvoiceIDOKBody
type PutTradespersonTradespersonIDBillingQuoteQuoteIDInvoiceInvoiceIDOKBody struct {

	// updated
	Updated bool `json:"updated"`
}

// Validate validates this put tradesperson tradesperson ID billing quote quote ID invoice invoice ID o k body
func (o *PutTradespersonTradespersonIDBillingQuoteQuoteIDInvoiceInvoiceIDOKBody) Validate(formats strfmt.Registry) error {
	return nil
}

// ContextValidate validates this put tradesperson tradesperson ID billing quote quote ID invoice invoice ID o k body based on context it is used
func (o *PutTradespersonTradespersonIDBillingQuoteQuoteIDInvoiceInvoiceIDOKBody) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (o *PutTradespersonTradespersonIDBillingQuoteQuoteIDInvoiceInvoiceIDOKBody) MarshalBinary() ([]byte, error) {
	if o == nil {
		return nil, nil
	}
	return swag.WriteJSON(o)
}

// UnmarshalBinary interface implementation
func (o *PutTradespersonTradespersonIDBillingQuoteQuoteIDInvoiceInvoiceIDOKBody) UnmarshalBinary(b []byte) error {
	var res PutTradespersonTradespersonIDBillingQuoteQuoteIDInvoiceInvoiceIDOKBody
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*o = res
	return nil
}
