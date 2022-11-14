// Code generated by go-swagger; DO NOT EDIT.

package operations

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"net/http"

	"github.com/go-openapi/runtime/middleware"
)

// GetTradespersonTradespersonIDFixedPricePagesHandlerFunc turns a function with the right signature into a get tradesperson tradesperson ID fixed price pages handler
type GetTradespersonTradespersonIDFixedPricePagesHandlerFunc func(GetTradespersonTradespersonIDFixedPricePagesParams, interface{}) middleware.Responder

// Handle executing the request and returning a response
func (fn GetTradespersonTradespersonIDFixedPricePagesHandlerFunc) Handle(params GetTradespersonTradespersonIDFixedPricePagesParams, principal interface{}) middleware.Responder {
	return fn(params, principal)
}

// GetTradespersonTradespersonIDFixedPricePagesHandler interface for that can handle valid get tradesperson tradesperson ID fixed price pages params
type GetTradespersonTradespersonIDFixedPricePagesHandler interface {
	Handle(GetTradespersonTradespersonIDFixedPricePagesParams, interface{}) middleware.Responder
}

// NewGetTradespersonTradespersonIDFixedPricePages creates a new http.Handler for the get tradesperson tradesperson ID fixed price pages operation
func NewGetTradespersonTradespersonIDFixedPricePages(ctx *middleware.Context, handler GetTradespersonTradespersonIDFixedPricePagesHandler) *GetTradespersonTradespersonIDFixedPricePages {
	return &GetTradespersonTradespersonIDFixedPricePages{Context: ctx, Handler: handler}
}

/* GetTradespersonTradespersonIDFixedPricePages swagger:route GET /tradesperson/{tradespersonId}/fixed-price/pages getTradespersonTradespersonIdFixedPricePages

GetTradespersonTradespersonIDFixedPricePages get tradesperson tradesperson ID fixed price pages API

*/
type GetTradespersonTradespersonIDFixedPricePages struct {
	Context *middleware.Context
	Handler GetTradespersonTradespersonIDFixedPricePagesHandler
}

func (o *GetTradespersonTradespersonIDFixedPricePages) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, rCtx, _ := o.Context.RouteInfo(r)
	if rCtx != nil {
		*r = *rCtx
	}
	var Params = NewGetTradespersonTradespersonIDFixedPricePagesParams()
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
