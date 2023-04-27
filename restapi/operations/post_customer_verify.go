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

// PostCustomerVerifyHandlerFunc turns a function with the right signature into a post customer verify handler
type PostCustomerVerifyHandlerFunc func(PostCustomerVerifyParams, interface{}) middleware.Responder

// Handle executing the request and returning a response
func (fn PostCustomerVerifyHandlerFunc) Handle(params PostCustomerVerifyParams, principal interface{}) middleware.Responder {
	return fn(params, principal)
}

// PostCustomerVerifyHandler interface for that can handle valid post customer verify params
type PostCustomerVerifyHandler interface {
	Handle(PostCustomerVerifyParams, interface{}) middleware.Responder
}

// NewPostCustomerVerify creates a new http.Handler for the post customer verify operation
func NewPostCustomerVerify(ctx *middleware.Context, handler PostCustomerVerifyHandler) *PostCustomerVerify {
	return &PostCustomerVerify{Context: ctx, Handler: handler}
}

/* PostCustomerVerify swagger:route POST /customer/verify postCustomerVerify

PostCustomerVerify post customer verify API

*/
type PostCustomerVerify struct {
	Context *middleware.Context
	Handler PostCustomerVerifyHandler
}

func (o *PostCustomerVerify) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, rCtx, _ := o.Context.RouteInfo(r)
	if rCtx != nil {
		*r = *rCtx
	}
	var Params = NewPostCustomerVerifyParams()
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

// PostCustomerVerifyOKBody post customer verify o k body
//
// swagger:model PostCustomerVerifyOKBody
type PostCustomerVerifyOKBody struct {

	// verified
	Verified bool `json:"verified"`
}

// Validate validates this post customer verify o k body
func (o *PostCustomerVerifyOKBody) Validate(formats strfmt.Registry) error {
	return nil
}

// ContextValidate validates this post customer verify o k body based on context it is used
func (o *PostCustomerVerifyOKBody) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (o *PostCustomerVerifyOKBody) MarshalBinary() ([]byte, error) {
	if o == nil {
		return nil, nil
	}
	return swag.WriteJSON(o)
}

// UnmarshalBinary interface implementation
func (o *PostCustomerVerifyOKBody) UnmarshalBinary(b []byte) error {
	var res PostCustomerVerifyOKBody
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*o = res
	return nil
}
