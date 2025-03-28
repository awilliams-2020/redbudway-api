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

// PostCustomerCustomerIDVerifyHandlerFunc turns a function with the right signature into a post customer customer ID verify handler
type PostCustomerCustomerIDVerifyHandlerFunc func(PostCustomerCustomerIDVerifyParams, interface{}) middleware.Responder

// Handle executing the request and returning a response
func (fn PostCustomerCustomerIDVerifyHandlerFunc) Handle(params PostCustomerCustomerIDVerifyParams, principal interface{}) middleware.Responder {
	return fn(params, principal)
}

// PostCustomerCustomerIDVerifyHandler interface for that can handle valid post customer customer ID verify params
type PostCustomerCustomerIDVerifyHandler interface {
	Handle(PostCustomerCustomerIDVerifyParams, interface{}) middleware.Responder
}

// NewPostCustomerCustomerIDVerify creates a new http.Handler for the post customer customer ID verify operation
func NewPostCustomerCustomerIDVerify(ctx *middleware.Context, handler PostCustomerCustomerIDVerifyHandler) *PostCustomerCustomerIDVerify {
	return &PostCustomerCustomerIDVerify{Context: ctx, Handler: handler}
}

/* PostCustomerCustomerIDVerify swagger:route POST /customer/{customerId}/verify postCustomerCustomerIdVerify

PostCustomerCustomerIDVerify post customer customer ID verify API

*/
type PostCustomerCustomerIDVerify struct {
	Context *middleware.Context
	Handler PostCustomerCustomerIDVerifyHandler
}

func (o *PostCustomerCustomerIDVerify) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, rCtx, _ := o.Context.RouteInfo(r)
	if rCtx != nil {
		*r = *rCtx
	}
	var Params = NewPostCustomerCustomerIDVerifyParams()
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

// PostCustomerCustomerIDVerifyOKBody post customer customer ID verify o k body
//
// swagger:model PostCustomerCustomerIDVerifyOKBody
type PostCustomerCustomerIDVerifyOKBody struct {

	// verified
	Verified bool `json:"verified"`
}

// Validate validates this post customer customer ID verify o k body
func (o *PostCustomerCustomerIDVerifyOKBody) Validate(formats strfmt.Registry) error {
	return nil
}

// ContextValidate validates this post customer customer ID verify o k body based on context it is used
func (o *PostCustomerCustomerIDVerifyOKBody) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (o *PostCustomerCustomerIDVerifyOKBody) MarshalBinary() ([]byte, error) {
	if o == nil {
		return nil, nil
	}
	return swag.WriteJSON(o)
}

// UnmarshalBinary interface implementation
func (o *PostCustomerCustomerIDVerifyOKBody) UnmarshalBinary(b []byte) error {
	var res PostCustomerCustomerIDVerifyOKBody
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*o = res
	return nil
}
