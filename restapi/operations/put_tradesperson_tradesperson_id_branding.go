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

// PutTradespersonTradespersonIDBrandingHandlerFunc turns a function with the right signature into a put tradesperson tradesperson ID branding handler
type PutTradespersonTradespersonIDBrandingHandlerFunc func(PutTradespersonTradespersonIDBrandingParams, interface{}) middleware.Responder

// Handle executing the request and returning a response
func (fn PutTradespersonTradespersonIDBrandingHandlerFunc) Handle(params PutTradespersonTradespersonIDBrandingParams, principal interface{}) middleware.Responder {
	return fn(params, principal)
}

// PutTradespersonTradespersonIDBrandingHandler interface for that can handle valid put tradesperson tradesperson ID branding params
type PutTradespersonTradespersonIDBrandingHandler interface {
	Handle(PutTradespersonTradespersonIDBrandingParams, interface{}) middleware.Responder
}

// NewPutTradespersonTradespersonIDBranding creates a new http.Handler for the put tradesperson tradesperson ID branding operation
func NewPutTradespersonTradespersonIDBranding(ctx *middleware.Context, handler PutTradespersonTradespersonIDBrandingHandler) *PutTradespersonTradespersonIDBranding {
	return &PutTradespersonTradespersonIDBranding{Context: ctx, Handler: handler}
}

/* PutTradespersonTradespersonIDBranding swagger:route PUT /tradesperson/{tradespersonId}/branding putTradespersonTradespersonIdBranding

PutTradespersonTradespersonIDBranding put tradesperson tradesperson ID branding API

*/
type PutTradespersonTradespersonIDBranding struct {
	Context *middleware.Context
	Handler PutTradespersonTradespersonIDBrandingHandler
}

func (o *PutTradespersonTradespersonIDBranding) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, rCtx, _ := o.Context.RouteInfo(r)
	if rCtx != nil {
		*r = *rCtx
	}
	var Params = NewPutTradespersonTradespersonIDBrandingParams()
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

// PutTradespersonTradespersonIDBrandingOKBody put tradesperson tradesperson ID branding o k body
//
// swagger:model PutTradespersonTradespersonIDBrandingOKBody
type PutTradespersonTradespersonIDBrandingOKBody struct {

	// updated
	Updated bool `json:"updated"`
}

// Validate validates this put tradesperson tradesperson ID branding o k body
func (o *PutTradespersonTradespersonIDBrandingOKBody) Validate(formats strfmt.Registry) error {
	return nil
}

// ContextValidate validates this put tradesperson tradesperson ID branding o k body based on context it is used
func (o *PutTradespersonTradespersonIDBrandingOKBody) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (o *PutTradespersonTradespersonIDBrandingOKBody) MarshalBinary() ([]byte, error) {
	if o == nil {
		return nil, nil
	}
	return swag.WriteJSON(o)
}

// UnmarshalBinary interface implementation
func (o *PutTradespersonTradespersonIDBrandingOKBody) UnmarshalBinary(b []byte) error {
	var res PutTradespersonTradespersonIDBrandingOKBody
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*o = res
	return nil
}
