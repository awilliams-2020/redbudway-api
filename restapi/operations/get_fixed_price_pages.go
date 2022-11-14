// Code generated by go-swagger; DO NOT EDIT.

package operations

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"net/http"

	"github.com/go-openapi/runtime/middleware"
)

// GetFixedPricePagesHandlerFunc turns a function with the right signature into a get fixed price pages handler
type GetFixedPricePagesHandlerFunc func(GetFixedPricePagesParams) middleware.Responder

// Handle executing the request and returning a response
func (fn GetFixedPricePagesHandlerFunc) Handle(params GetFixedPricePagesParams) middleware.Responder {
	return fn(params)
}

// GetFixedPricePagesHandler interface for that can handle valid get fixed price pages params
type GetFixedPricePagesHandler interface {
	Handle(GetFixedPricePagesParams) middleware.Responder
}

// NewGetFixedPricePages creates a new http.Handler for the get fixed price pages operation
func NewGetFixedPricePages(ctx *middleware.Context, handler GetFixedPricePagesHandler) *GetFixedPricePages {
	return &GetFixedPricePages{Context: ctx, Handler: handler}
}

/* GetFixedPricePages swagger:route GET /fixed-price/pages getFixedPricePages

GetFixedPricePages get fixed price pages API

*/
type GetFixedPricePages struct {
	Context *middleware.Context
	Handler GetFixedPricePagesHandler
}

func (o *GetFixedPricePages) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, rCtx, _ := o.Context.RouteInfo(r)
	if rCtx != nil {
		*r = *rCtx
	}
	var Params = NewGetFixedPricePagesParams()
	if err := o.Context.BindValidRequest(r, route, &Params); err != nil { // bind params
		o.Context.Respond(rw, r, route.Produces, route, err)
		return
	}

	res := o.Handler.Handle(Params) // actually handle the request
	o.Context.Respond(rw, r, route.Produces, route, res)

}
