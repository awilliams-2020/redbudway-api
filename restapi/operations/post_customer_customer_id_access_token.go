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

// PostCustomerCustomerIDAccessTokenHandlerFunc turns a function with the right signature into a post customer customer ID access token handler
type PostCustomerCustomerIDAccessTokenHandlerFunc func(PostCustomerCustomerIDAccessTokenParams, interface{}) middleware.Responder

// Handle executing the request and returning a response
func (fn PostCustomerCustomerIDAccessTokenHandlerFunc) Handle(params PostCustomerCustomerIDAccessTokenParams, principal interface{}) middleware.Responder {
	return fn(params, principal)
}

// PostCustomerCustomerIDAccessTokenHandler interface for that can handle valid post customer customer ID access token params
type PostCustomerCustomerIDAccessTokenHandler interface {
	Handle(PostCustomerCustomerIDAccessTokenParams, interface{}) middleware.Responder
}

// NewPostCustomerCustomerIDAccessToken creates a new http.Handler for the post customer customer ID access token operation
func NewPostCustomerCustomerIDAccessToken(ctx *middleware.Context, handler PostCustomerCustomerIDAccessTokenHandler) *PostCustomerCustomerIDAccessToken {
	return &PostCustomerCustomerIDAccessToken{Context: ctx, Handler: handler}
}

/* PostCustomerCustomerIDAccessToken swagger:route POST /customer/{customerId}/access-token postCustomerCustomerIdAccessToken

PostCustomerCustomerIDAccessToken post customer customer ID access token API

*/
type PostCustomerCustomerIDAccessToken struct {
	Context *middleware.Context
	Handler PostCustomerCustomerIDAccessTokenHandler
}

func (o *PostCustomerCustomerIDAccessToken) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, rCtx, _ := o.Context.RouteInfo(r)
	if rCtx != nil {
		*r = *rCtx
	}
	var Params = NewPostCustomerCustomerIDAccessTokenParams()
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

// PostCustomerCustomerIDAccessTokenOKBody post customer customer ID access token o k body
//
// swagger:model PostCustomerCustomerIDAccessTokenOKBody
type PostCustomerCustomerIDAccessTokenOKBody struct {

	// access token
	AccessToken string `json:"accessToken,omitempty"`

	// refresh token
	RefreshToken string `json:"refreshToken,omitempty"`
}

// Validate validates this post customer customer ID access token o k body
func (o *PostCustomerCustomerIDAccessTokenOKBody) Validate(formats strfmt.Registry) error {
	return nil
}

// ContextValidate validates this post customer customer ID access token o k body based on context it is used
func (o *PostCustomerCustomerIDAccessTokenOKBody) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (o *PostCustomerCustomerIDAccessTokenOKBody) MarshalBinary() ([]byte, error) {
	if o == nil {
		return nil, nil
	}
	return swag.WriteJSON(o)
}

// UnmarshalBinary interface implementation
func (o *PostCustomerCustomerIDAccessTokenOKBody) UnmarshalBinary(b []byte) error {
	var res PostCustomerCustomerIDAccessTokenOKBody
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*o = res
	return nil
}
