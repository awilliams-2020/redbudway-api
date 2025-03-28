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

// PutTradespersonTradespersonIDPasswordHandlerFunc turns a function with the right signature into a put tradesperson tradesperson ID password handler
type PutTradespersonTradespersonIDPasswordHandlerFunc func(PutTradespersonTradespersonIDPasswordParams, interface{}) middleware.Responder

// Handle executing the request and returning a response
func (fn PutTradespersonTradespersonIDPasswordHandlerFunc) Handle(params PutTradespersonTradespersonIDPasswordParams, principal interface{}) middleware.Responder {
	return fn(params, principal)
}

// PutTradespersonTradespersonIDPasswordHandler interface for that can handle valid put tradesperson tradesperson ID password params
type PutTradespersonTradespersonIDPasswordHandler interface {
	Handle(PutTradespersonTradespersonIDPasswordParams, interface{}) middleware.Responder
}

// NewPutTradespersonTradespersonIDPassword creates a new http.Handler for the put tradesperson tradesperson ID password operation
func NewPutTradespersonTradespersonIDPassword(ctx *middleware.Context, handler PutTradespersonTradespersonIDPasswordHandler) *PutTradespersonTradespersonIDPassword {
	return &PutTradespersonTradespersonIDPassword{Context: ctx, Handler: handler}
}

/* PutTradespersonTradespersonIDPassword swagger:route PUT /tradesperson/{tradespersonId}/password putTradespersonTradespersonIdPassword

PutTradespersonTradespersonIDPassword put tradesperson tradesperson ID password API

*/
type PutTradespersonTradespersonIDPassword struct {
	Context *middleware.Context
	Handler PutTradespersonTradespersonIDPasswordHandler
}

func (o *PutTradespersonTradespersonIDPassword) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, rCtx, _ := o.Context.RouteInfo(r)
	if rCtx != nil {
		*r = *rCtx
	}
	var Params = NewPutTradespersonTradespersonIDPasswordParams()
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

// PutTradespersonTradespersonIDPasswordBody put tradesperson tradesperson ID password body
//
// swagger:model PutTradespersonTradespersonIDPasswordBody
type PutTradespersonTradespersonIDPasswordBody struct {

	// cur password
	// Required: true
	CurPassword *string `json:"curPassword"`

	// new password
	// Required: true
	NewPassword *string `json:"newPassword"`
}

// Validate validates this put tradesperson tradesperson ID password body
func (o *PutTradespersonTradespersonIDPasswordBody) Validate(formats strfmt.Registry) error {
	var res []error

	if err := o.validateCurPassword(formats); err != nil {
		res = append(res, err)
	}

	if err := o.validateNewPassword(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (o *PutTradespersonTradespersonIDPasswordBody) validateCurPassword(formats strfmt.Registry) error {

	if err := validate.Required("tradesperson"+"."+"curPassword", "body", o.CurPassword); err != nil {
		return err
	}

	return nil
}

func (o *PutTradespersonTradespersonIDPasswordBody) validateNewPassword(formats strfmt.Registry) error {

	if err := validate.Required("tradesperson"+"."+"newPassword", "body", o.NewPassword); err != nil {
		return err
	}

	return nil
}

// ContextValidate validates this put tradesperson tradesperson ID password body based on context it is used
func (o *PutTradespersonTradespersonIDPasswordBody) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (o *PutTradespersonTradespersonIDPasswordBody) MarshalBinary() ([]byte, error) {
	if o == nil {
		return nil, nil
	}
	return swag.WriteJSON(o)
}

// UnmarshalBinary interface implementation
func (o *PutTradespersonTradespersonIDPasswordBody) UnmarshalBinary(b []byte) error {
	var res PutTradespersonTradespersonIDPasswordBody
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*o = res
	return nil
}

// PutTradespersonTradespersonIDPasswordOKBody put tradesperson tradesperson ID password o k body
//
// swagger:model PutTradespersonTradespersonIDPasswordOKBody
type PutTradespersonTradespersonIDPasswordOKBody struct {

	// updated
	Updated bool `json:"updated"`
}

// Validate validates this put tradesperson tradesperson ID password o k body
func (o *PutTradespersonTradespersonIDPasswordOKBody) Validate(formats strfmt.Registry) error {
	return nil
}

// ContextValidate validates this put tradesperson tradesperson ID password o k body based on context it is used
func (o *PutTradespersonTradespersonIDPasswordOKBody) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (o *PutTradespersonTradespersonIDPasswordOKBody) MarshalBinary() ([]byte, error) {
	if o == nil {
		return nil, nil
	}
	return swag.WriteJSON(o)
}

// UnmarshalBinary interface implementation
func (o *PutTradespersonTradespersonIDPasswordOKBody) UnmarshalBinary(b []byte) error {
	var res PutTradespersonTradespersonIDPasswordOKBody
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*o = res
	return nil
}
