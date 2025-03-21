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

// PostTradespersonTradespersonIDAccessTokenHandlerFunc turns a function with the right signature into a post tradesperson tradesperson ID access token handler
type PostTradespersonTradespersonIDAccessTokenHandlerFunc func(PostTradespersonTradespersonIDAccessTokenParams, interface{}) middleware.Responder

// Handle executing the request and returning a response
func (fn PostTradespersonTradespersonIDAccessTokenHandlerFunc) Handle(params PostTradespersonTradespersonIDAccessTokenParams, principal interface{}) middleware.Responder {
	return fn(params, principal)
}

// PostTradespersonTradespersonIDAccessTokenHandler interface for that can handle valid post tradesperson tradesperson ID access token params
type PostTradespersonTradespersonIDAccessTokenHandler interface {
	Handle(PostTradespersonTradespersonIDAccessTokenParams, interface{}) middleware.Responder
}

// NewPostTradespersonTradespersonIDAccessToken creates a new http.Handler for the post tradesperson tradesperson ID access token operation
func NewPostTradespersonTradespersonIDAccessToken(ctx *middleware.Context, handler PostTradespersonTradespersonIDAccessTokenHandler) *PostTradespersonTradespersonIDAccessToken {
	return &PostTradespersonTradespersonIDAccessToken{Context: ctx, Handler: handler}
}

/* PostTradespersonTradespersonIDAccessToken swagger:route POST /tradesperson/{tradespersonId}/access-token postTradespersonTradespersonIdAccessToken

PostTradespersonTradespersonIDAccessToken post tradesperson tradesperson ID access token API

*/
type PostTradespersonTradespersonIDAccessToken struct {
	Context *middleware.Context
	Handler PostTradespersonTradespersonIDAccessTokenHandler
}

func (o *PostTradespersonTradespersonIDAccessToken) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, rCtx, _ := o.Context.RouteInfo(r)
	if rCtx != nil {
		*r = *rCtx
	}
	var Params = NewPostTradespersonTradespersonIDAccessTokenParams()
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

// PostTradespersonTradespersonIDAccessTokenOKBody post tradesperson tradesperson ID access token o k body
//
// swagger:model PostTradespersonTradespersonIDAccessTokenOKBody
type PostTradespersonTradespersonIDAccessTokenOKBody struct {

	// access token
	AccessToken string `json:"accessToken,omitempty"`

	// refresh token
	RefreshToken string `json:"refreshToken,omitempty"`
}

// Validate validates this post tradesperson tradesperson ID access token o k body
func (o *PostTradespersonTradespersonIDAccessTokenOKBody) Validate(formats strfmt.Registry) error {
	return nil
}

// ContextValidate validates this post tradesperson tradesperson ID access token o k body based on context it is used
func (o *PostTradespersonTradespersonIDAccessTokenOKBody) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (o *PostTradespersonTradespersonIDAccessTokenOKBody) MarshalBinary() ([]byte, error) {
	if o == nil {
		return nil, nil
	}
	return swag.WriteJSON(o)
}

// UnmarshalBinary interface implementation
func (o *PostTradespersonTradespersonIDAccessTokenOKBody) UnmarshalBinary(b []byte) error {
	var res PostTradespersonTradespersonIDAccessTokenOKBody
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*o = res
	return nil
}
