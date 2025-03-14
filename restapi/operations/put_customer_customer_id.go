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

// PutCustomerCustomerIDHandlerFunc turns a function with the right signature into a put customer customer ID handler
type PutCustomerCustomerIDHandlerFunc func(PutCustomerCustomerIDParams, interface{}) middleware.Responder

// Handle executing the request and returning a response
func (fn PutCustomerCustomerIDHandlerFunc) Handle(params PutCustomerCustomerIDParams, principal interface{}) middleware.Responder {
	return fn(params, principal)
}

// PutCustomerCustomerIDHandler interface for that can handle valid put customer customer ID params
type PutCustomerCustomerIDHandler interface {
	Handle(PutCustomerCustomerIDParams, interface{}) middleware.Responder
}

// NewPutCustomerCustomerID creates a new http.Handler for the put customer customer ID operation
func NewPutCustomerCustomerID(ctx *middleware.Context, handler PutCustomerCustomerIDHandler) *PutCustomerCustomerID {
	return &PutCustomerCustomerID{Context: ctx, Handler: handler}
}

/* PutCustomerCustomerID swagger:route PUT /customer/{customerId} putCustomerCustomerId

PutCustomerCustomerID put customer customer ID API

*/
type PutCustomerCustomerID struct {
	Context *middleware.Context
	Handler PutCustomerCustomerIDHandler
}

func (o *PutCustomerCustomerID) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, rCtx, _ := o.Context.RouteInfo(r)
	if rCtx != nil {
		*r = *rCtx
	}
	var Params = NewPutCustomerCustomerIDParams()
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

// PutCustomerCustomerIDBody put customer customer ID body
//
// swagger:model PutCustomerCustomerIDBody
type PutCustomerCustomerIDBody struct {

	// cur password
	CurPassword string `json:"curPassword,omitempty"`

	// new password
	NewPassword string `json:"newPassword,omitempty"`
}

// Validate validates this put customer customer ID body
func (o *PutCustomerCustomerIDBody) Validate(formats strfmt.Registry) error {
	return nil
}

// ContextValidate validates this put customer customer ID body based on context it is used
func (o *PutCustomerCustomerIDBody) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (o *PutCustomerCustomerIDBody) MarshalBinary() ([]byte, error) {
	if o == nil {
		return nil, nil
	}
	return swag.WriteJSON(o)
}

// UnmarshalBinary interface implementation
func (o *PutCustomerCustomerIDBody) UnmarshalBinary(b []byte) error {
	var res PutCustomerCustomerIDBody
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*o = res
	return nil
}

// PutCustomerCustomerIDOKBody put customer customer ID o k body
//
// swagger:model PutCustomerCustomerIDOKBody
type PutCustomerCustomerIDOKBody struct {

	// updated
	Updated bool `json:"updated"`
}

// Validate validates this put customer customer ID o k body
func (o *PutCustomerCustomerIDOKBody) Validate(formats strfmt.Registry) error {
	return nil
}

// ContextValidate validates this put customer customer ID o k body based on context it is used
func (o *PutCustomerCustomerIDOKBody) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (o *PutCustomerCustomerIDOKBody) MarshalBinary() ([]byte, error) {
	if o == nil {
		return nil, nil
	}
	return swag.WriteJSON(o)
}

// UnmarshalBinary interface implementation
func (o *PutCustomerCustomerIDOKBody) UnmarshalBinary(b []byte) error {
	var res PutCustomerCustomerIDOKBody
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*o = res
	return nil
}
