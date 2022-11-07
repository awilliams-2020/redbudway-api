// Code generated by go-swagger; DO NOT EDIT.

package operations

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"net/http"

	"github.com/go-openapi/runtime/middleware"
)

// GetProfileVanityOrIDQuotesHandlerFunc turns a function with the right signature into a get profile vanity or ID quotes handler
type GetProfileVanityOrIDQuotesHandlerFunc func(GetProfileVanityOrIDQuotesParams) middleware.Responder

// Handle executing the request and returning a response
func (fn GetProfileVanityOrIDQuotesHandlerFunc) Handle(params GetProfileVanityOrIDQuotesParams) middleware.Responder {
	return fn(params)
}

// GetProfileVanityOrIDQuotesHandler interface for that can handle valid get profile vanity or ID quotes params
type GetProfileVanityOrIDQuotesHandler interface {
	Handle(GetProfileVanityOrIDQuotesParams) middleware.Responder
}

// NewGetProfileVanityOrIDQuotes creates a new http.Handler for the get profile vanity or ID quotes operation
func NewGetProfileVanityOrIDQuotes(ctx *middleware.Context, handler GetProfileVanityOrIDQuotesHandler) *GetProfileVanityOrIDQuotes {
	return &GetProfileVanityOrIDQuotes{Context: ctx, Handler: handler}
}

/* GetProfileVanityOrIDQuotes swagger:route GET /profile/{vanityOrId}/quotes getProfileVanityOrIdQuotes

GetProfileVanityOrIDQuotes get profile vanity or ID quotes API

*/
type GetProfileVanityOrIDQuotes struct {
	Context *middleware.Context
	Handler GetProfileVanityOrIDQuotesHandler
}

func (o *GetProfileVanityOrIDQuotes) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, rCtx, _ := o.Context.RouteInfo(r)
	if rCtx != nil {
		*r = *rCtx
	}
	var Params = NewGetProfileVanityOrIDQuotesParams()
	if err := o.Context.BindValidRequest(r, route, &Params); err != nil { // bind params
		o.Context.Respond(rw, r, route.Produces, route, err)
		return
	}

	res := o.Handler.Handle(Params) // actually handle the request
	o.Context.Respond(rw, r, route.Produces, route, res)

}
