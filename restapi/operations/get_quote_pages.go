// Code generated by go-swagger; DO NOT EDIT.

package operations

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"net/http"

	"github.com/go-openapi/runtime/middleware"
)

// GetQuotePagesHandlerFunc turns a function with the right signature into a get quote pages handler
type GetQuotePagesHandlerFunc func(GetQuotePagesParams) middleware.Responder

// Handle executing the request and returning a response
func (fn GetQuotePagesHandlerFunc) Handle(params GetQuotePagesParams) middleware.Responder {
	return fn(params)
}

// GetQuotePagesHandler interface for that can handle valid get quote pages params
type GetQuotePagesHandler interface {
	Handle(GetQuotePagesParams) middleware.Responder
}

// NewGetQuotePages creates a new http.Handler for the get quote pages operation
func NewGetQuotePages(ctx *middleware.Context, handler GetQuotePagesHandler) *GetQuotePages {
	return &GetQuotePages{Context: ctx, Handler: handler}
}

/* GetQuotePages swagger:route GET /quote/pages getQuotePages

GetQuotePages get quote pages API

*/
type GetQuotePages struct {
	Context *middleware.Context
	Handler GetQuotePagesHandler
}

func (o *GetQuotePages) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, rCtx, _ := o.Context.RouteInfo(r)
	if rCtx != nil {
		*r = *rCtx
	}
	var Params = NewGetQuotePagesParams()
	if err := o.Context.BindValidRequest(r, route, &Params); err != nil { // bind params
		o.Context.Respond(rw, r, route.Produces, route, err)
		return
	}

	res := o.Handler.Handle(Params) // actually handle the request
	o.Context.Respond(rw, r, route.Produces, route, res)

}
