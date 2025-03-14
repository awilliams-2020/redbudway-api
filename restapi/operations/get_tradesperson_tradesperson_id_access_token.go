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

// GetTradespersonTradespersonIDAccessTokenHandlerFunc turns a function with the right signature into a get tradesperson tradesperson ID access token handler
type GetTradespersonTradespersonIDAccessTokenHandlerFunc func(GetTradespersonTradespersonIDAccessTokenParams, interface{}) middleware.Responder

// Handle executing the request and returning a response
func (fn GetTradespersonTradespersonIDAccessTokenHandlerFunc) Handle(params GetTradespersonTradespersonIDAccessTokenParams, principal interface{}) middleware.Responder {
	return fn(params, principal)
}

// GetTradespersonTradespersonIDAccessTokenHandler interface for that can handle valid get tradesperson tradesperson ID access token params
type GetTradespersonTradespersonIDAccessTokenHandler interface {
	Handle(GetTradespersonTradespersonIDAccessTokenParams, interface{}) middleware.Responder
}

// NewGetTradespersonTradespersonIDAccessToken creates a new http.Handler for the get tradesperson tradesperson ID access token operation
func NewGetTradespersonTradespersonIDAccessToken(ctx *middleware.Context, handler GetTradespersonTradespersonIDAccessTokenHandler) *GetTradespersonTradespersonIDAccessToken {
	return &GetTradespersonTradespersonIDAccessToken{Context: ctx, Handler: handler}
}

/* GetTradespersonTradespersonIDAccessToken swagger:route GET /tradesperson/{tradespersonId}/access-token getTradespersonTradespersonIdAccessToken

GetTradespersonTradespersonIDAccessToken get tradesperson tradesperson ID access token API

*/
type GetTradespersonTradespersonIDAccessToken struct {
	Context *middleware.Context
	Handler GetTradespersonTradespersonIDAccessTokenHandler
}

func (o *GetTradespersonTradespersonIDAccessToken) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, rCtx, _ := o.Context.RouteInfo(r)
	if rCtx != nil {
		*r = *rCtx
	}
	var Params = NewGetTradespersonTradespersonIDAccessTokenParams()
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

// GetTradespersonTradespersonIDAccessTokenOKBody get tradesperson tradesperson ID access token o k body
//
// swagger:model GetTradespersonTradespersonIDAccessTokenOKBody
type GetTradespersonTradespersonIDAccessTokenOKBody struct {

	// valid
	Valid bool `json:"valid"`
}

// Validate validates this get tradesperson tradesperson ID access token o k body
func (o *GetTradespersonTradespersonIDAccessTokenOKBody) Validate(formats strfmt.Registry) error {
	return nil
}

// ContextValidate validates this get tradesperson tradesperson ID access token o k body based on context it is used
func (o *GetTradespersonTradespersonIDAccessTokenOKBody) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (o *GetTradespersonTradespersonIDAccessTokenOKBody) MarshalBinary() ([]byte, error) {
	if o == nil {
		return nil, nil
	}
	return swag.WriteJSON(o)
}

// UnmarshalBinary interface implementation
func (o *GetTradespersonTradespersonIDAccessTokenOKBody) UnmarshalBinary(b []byte) error {
	var res GetTradespersonTradespersonIDAccessTokenOKBody
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*o = res
	return nil
}
