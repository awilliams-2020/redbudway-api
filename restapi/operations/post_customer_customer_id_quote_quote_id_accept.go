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

// PostCustomerCustomerIDQuoteQuoteIDAcceptHandlerFunc turns a function with the right signature into a post customer customer ID quote quote ID accept handler
type PostCustomerCustomerIDQuoteQuoteIDAcceptHandlerFunc func(PostCustomerCustomerIDQuoteQuoteIDAcceptParams, interface{}) middleware.Responder

// Handle executing the request and returning a response
func (fn PostCustomerCustomerIDQuoteQuoteIDAcceptHandlerFunc) Handle(params PostCustomerCustomerIDQuoteQuoteIDAcceptParams, principal interface{}) middleware.Responder {
	return fn(params, principal)
}

// PostCustomerCustomerIDQuoteQuoteIDAcceptHandler interface for that can handle valid post customer customer ID quote quote ID accept params
type PostCustomerCustomerIDQuoteQuoteIDAcceptHandler interface {
	Handle(PostCustomerCustomerIDQuoteQuoteIDAcceptParams, interface{}) middleware.Responder
}

// NewPostCustomerCustomerIDQuoteQuoteIDAccept creates a new http.Handler for the post customer customer ID quote quote ID accept operation
func NewPostCustomerCustomerIDQuoteQuoteIDAccept(ctx *middleware.Context, handler PostCustomerCustomerIDQuoteQuoteIDAcceptHandler) *PostCustomerCustomerIDQuoteQuoteIDAccept {
	return &PostCustomerCustomerIDQuoteQuoteIDAccept{Context: ctx, Handler: handler}
}

/* PostCustomerCustomerIDQuoteQuoteIDAccept swagger:route POST /customer/{customerId}/quote/{quoteId}/accept postCustomerCustomerIdQuoteQuoteIdAccept

PostCustomerCustomerIDQuoteQuoteIDAccept post customer customer ID quote quote ID accept API

*/
type PostCustomerCustomerIDQuoteQuoteIDAccept struct {
	Context *middleware.Context
	Handler PostCustomerCustomerIDQuoteQuoteIDAcceptHandler
}

func (o *PostCustomerCustomerIDQuoteQuoteIDAccept) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, rCtx, _ := o.Context.RouteInfo(r)
	if rCtx != nil {
		*r = *rCtx
	}
	var Params = NewPostCustomerCustomerIDQuoteQuoteIDAcceptParams()
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

// PostCustomerCustomerIDQuoteQuoteIDAcceptOKBody post customer customer ID quote quote ID accept o k body
//
// swagger:model PostCustomerCustomerIDQuoteQuoteIDAcceptOKBody
type PostCustomerCustomerIDQuoteQuoteIDAcceptOKBody struct {

	// accepted
	Accepted bool `json:"accepted"`
}

// Validate validates this post customer customer ID quote quote ID accept o k body
func (o *PostCustomerCustomerIDQuoteQuoteIDAcceptOKBody) Validate(formats strfmt.Registry) error {
	return nil
}

// ContextValidate validates this post customer customer ID quote quote ID accept o k body based on context it is used
func (o *PostCustomerCustomerIDQuoteQuoteIDAcceptOKBody) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (o *PostCustomerCustomerIDQuoteQuoteIDAcceptOKBody) MarshalBinary() ([]byte, error) {
	if o == nil {
		return nil, nil
	}
	return swag.WriteJSON(o)
}

// UnmarshalBinary interface implementation
func (o *PostCustomerCustomerIDQuoteQuoteIDAcceptOKBody) UnmarshalBinary(b []byte) error {
	var res PostCustomerCustomerIDQuoteQuoteIDAcceptOKBody
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*o = res
	return nil
}
