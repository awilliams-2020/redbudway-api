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

// PutTradespersonTradespersonIDFixedPriceHandlerFunc turns a function with the right signature into a put tradesperson tradesperson ID fixed price handler
type PutTradespersonTradespersonIDFixedPriceHandlerFunc func(PutTradespersonTradespersonIDFixedPriceParams, interface{}) middleware.Responder

// Handle executing the request and returning a response
func (fn PutTradespersonTradespersonIDFixedPriceHandlerFunc) Handle(params PutTradespersonTradespersonIDFixedPriceParams, principal interface{}) middleware.Responder {
	return fn(params, principal)
}

// PutTradespersonTradespersonIDFixedPriceHandler interface for that can handle valid put tradesperson tradesperson ID fixed price params
type PutTradespersonTradespersonIDFixedPriceHandler interface {
	Handle(PutTradespersonTradespersonIDFixedPriceParams, interface{}) middleware.Responder
}

// NewPutTradespersonTradespersonIDFixedPrice creates a new http.Handler for the put tradesperson tradesperson ID fixed price operation
func NewPutTradespersonTradespersonIDFixedPrice(ctx *middleware.Context, handler PutTradespersonTradespersonIDFixedPriceHandler) *PutTradespersonTradespersonIDFixedPrice {
	return &PutTradespersonTradespersonIDFixedPrice{Context: ctx, Handler: handler}
}

/* PutTradespersonTradespersonIDFixedPrice swagger:route PUT /tradesperson/{tradespersonId}/fixed-price putTradespersonTradespersonIdFixedPrice

PutTradespersonTradespersonIDFixedPrice put tradesperson tradesperson ID fixed price API

*/
type PutTradespersonTradespersonIDFixedPrice struct {
	Context *middleware.Context
	Handler PutTradespersonTradespersonIDFixedPriceHandler
}

func (o *PutTradespersonTradespersonIDFixedPrice) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, rCtx, _ := o.Context.RouteInfo(r)
	if rCtx != nil {
		*r = *rCtx
	}
	var Params = NewPutTradespersonTradespersonIDFixedPriceParams()
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

// PutTradespersonTradespersonIDFixedPriceOKBody put tradesperson tradesperson ID fixed price o k body
//
// swagger:model PutTradespersonTradespersonIDFixedPriceOKBody
type PutTradespersonTradespersonIDFixedPriceOKBody struct {

	// updated
	Updated bool `json:"updated"`
}

// Validate validates this put tradesperson tradesperson ID fixed price o k body
func (o *PutTradespersonTradespersonIDFixedPriceOKBody) Validate(formats strfmt.Registry) error {
	return nil
}

// ContextValidate validates this put tradesperson tradesperson ID fixed price o k body based on context it is used
func (o *PutTradespersonTradespersonIDFixedPriceOKBody) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (o *PutTradespersonTradespersonIDFixedPriceOKBody) MarshalBinary() ([]byte, error) {
	if o == nil {
		return nil, nil
	}
	return swag.WriteJSON(o)
}

// UnmarshalBinary interface implementation
func (o *PutTradespersonTradespersonIDFixedPriceOKBody) UnmarshalBinary(b []byte) error {
	var res PutTradespersonTradespersonIDFixedPriceOKBody
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*o = res
	return nil
}
