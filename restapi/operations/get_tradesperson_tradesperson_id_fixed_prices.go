// Code generated by go-swagger; DO NOT EDIT.

package operations

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"net/http"

	"github.com/go-openapi/runtime/middleware"
)

// GetTradespersonTradespersonIDFixedPricesHandlerFunc turns a function with the right signature into a get tradesperson tradesperson ID fixed prices handler
type GetTradespersonTradespersonIDFixedPricesHandlerFunc func(GetTradespersonTradespersonIDFixedPricesParams, interface{}) middleware.Responder

// Handle executing the request and returning a response
func (fn GetTradespersonTradespersonIDFixedPricesHandlerFunc) Handle(params GetTradespersonTradespersonIDFixedPricesParams, principal interface{}) middleware.Responder {
	return fn(params, principal)
}

// GetTradespersonTradespersonIDFixedPricesHandler interface for that can handle valid get tradesperson tradesperson ID fixed prices params
type GetTradespersonTradespersonIDFixedPricesHandler interface {
	Handle(GetTradespersonTradespersonIDFixedPricesParams, interface{}) middleware.Responder
}

// NewGetTradespersonTradespersonIDFixedPrices creates a new http.Handler for the get tradesperson tradesperson ID fixed prices operation
func NewGetTradespersonTradespersonIDFixedPrices(ctx *middleware.Context, handler GetTradespersonTradespersonIDFixedPricesHandler) *GetTradespersonTradespersonIDFixedPrices {
	return &GetTradespersonTradespersonIDFixedPrices{Context: ctx, Handler: handler}
}

/* GetTradespersonTradespersonIDFixedPrices swagger:route GET /tradesperson/{tradespersonId}/fixed-prices getTradespersonTradespersonIdFixedPrices

GetTradespersonTradespersonIDFixedPrices get tradesperson tradesperson ID fixed prices API

*/
type GetTradespersonTradespersonIDFixedPrices struct {
	Context *middleware.Context
	Handler GetTradespersonTradespersonIDFixedPricesHandler
}

func (o *GetTradespersonTradespersonIDFixedPrices) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, rCtx, _ := o.Context.RouteInfo(r)
	if rCtx != nil {
		*r = *rCtx
	}
	var Params = NewGetTradespersonTradespersonIDFixedPricesParams()
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
