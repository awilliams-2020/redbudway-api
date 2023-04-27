// Code generated by go-swagger; DO NOT EDIT.

package operations

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"context"
	"net/http"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
	"github.com/go-openapi/validate"
)

// PutTradespersonTradespersonIDAccountEmailHandlerFunc turns a function with the right signature into a put tradesperson tradesperson ID account email handler
type PutTradespersonTradespersonIDAccountEmailHandlerFunc func(PutTradespersonTradespersonIDAccountEmailParams, interface{}) middleware.Responder

// Handle executing the request and returning a response
func (fn PutTradespersonTradespersonIDAccountEmailHandlerFunc) Handle(params PutTradespersonTradespersonIDAccountEmailParams, principal interface{}) middleware.Responder {
	return fn(params, principal)
}

// PutTradespersonTradespersonIDAccountEmailHandler interface for that can handle valid put tradesperson tradesperson ID account email params
type PutTradespersonTradespersonIDAccountEmailHandler interface {
	Handle(PutTradespersonTradespersonIDAccountEmailParams, interface{}) middleware.Responder
}

// NewPutTradespersonTradespersonIDAccountEmail creates a new http.Handler for the put tradesperson tradesperson ID account email operation
func NewPutTradespersonTradespersonIDAccountEmail(ctx *middleware.Context, handler PutTradespersonTradespersonIDAccountEmailHandler) *PutTradespersonTradespersonIDAccountEmail {
	return &PutTradespersonTradespersonIDAccountEmail{Context: ctx, Handler: handler}
}

/* PutTradespersonTradespersonIDAccountEmail swagger:route PUT /tradesperson/{tradespersonId}/account-email putTradespersonTradespersonIdAccountEmail

PutTradespersonTradespersonIDAccountEmail put tradesperson tradesperson ID account email API

*/
type PutTradespersonTradespersonIDAccountEmail struct {
	Context *middleware.Context
	Handler PutTradespersonTradespersonIDAccountEmailHandler
}

func (o *PutTradespersonTradespersonIDAccountEmail) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, rCtx, _ := o.Context.RouteInfo(r)
	if rCtx != nil {
		*r = *rCtx
	}
	var Params = NewPutTradespersonTradespersonIDAccountEmailParams()
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

// PutTradespersonTradespersonIDAccountEmailBody put tradesperson tradesperson ID account email body
//
// swagger:model PutTradespersonTradespersonIDAccountEmailBody
type PutTradespersonTradespersonIDAccountEmailBody struct {

	// email
	// Required: true
	Email *string `json:"email"`
}

// Validate validates this put tradesperson tradesperson ID account email body
func (o *PutTradespersonTradespersonIDAccountEmailBody) Validate(formats strfmt.Registry) error {
	var res []error

	if err := o.validateEmail(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (o *PutTradespersonTradespersonIDAccountEmailBody) validateEmail(formats strfmt.Registry) error {

	if err := validate.Required("account"+"."+"email", "body", o.Email); err != nil {
		return err
	}

	return nil
}

// ContextValidate validates this put tradesperson tradesperson ID account email body based on context it is used
func (o *PutTradespersonTradespersonIDAccountEmailBody) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (o *PutTradespersonTradespersonIDAccountEmailBody) MarshalBinary() ([]byte, error) {
	if o == nil {
		return nil, nil
	}
	return swag.WriteJSON(o)
}

// UnmarshalBinary interface implementation
func (o *PutTradespersonTradespersonIDAccountEmailBody) UnmarshalBinary(b []byte) error {
	var res PutTradespersonTradespersonIDAccountEmailBody
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*o = res
	return nil
}

// PutTradespersonTradespersonIDAccountEmailOKBody put tradesperson tradesperson ID account email o k body
//
// swagger:model PutTradespersonTradespersonIDAccountEmailOKBody
type PutTradespersonTradespersonIDAccountEmailOKBody struct {

	// updated
	Updated bool `json:"updated"`
}

// Validate validates this put tradesperson tradesperson ID account email o k body
func (o *PutTradespersonTradespersonIDAccountEmailOKBody) Validate(formats strfmt.Registry) error {
	return nil
}

// ContextValidate validates this put tradesperson tradesperson ID account email o k body based on context it is used
func (o *PutTradespersonTradespersonIDAccountEmailOKBody) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (o *PutTradespersonTradespersonIDAccountEmailOKBody) MarshalBinary() ([]byte, error) {
	if o == nil {
		return nil, nil
	}
	return swag.WriteJSON(o)
}

// UnmarshalBinary interface implementation
func (o *PutTradespersonTradespersonIDAccountEmailOKBody) UnmarshalBinary(b []byte) error {
	var res PutTradespersonTradespersonIDAccountEmailOKBody
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*o = res
	return nil
}
