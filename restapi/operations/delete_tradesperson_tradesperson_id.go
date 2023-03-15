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

// DeleteTradespersonTradespersonIDHandlerFunc turns a function with the right signature into a delete tradesperson tradesperson ID handler
type DeleteTradespersonTradespersonIDHandlerFunc func(DeleteTradespersonTradespersonIDParams, interface{}) middleware.Responder

// Handle executing the request and returning a response
func (fn DeleteTradespersonTradespersonIDHandlerFunc) Handle(params DeleteTradespersonTradespersonIDParams, principal interface{}) middleware.Responder {
	return fn(params, principal)
}

// DeleteTradespersonTradespersonIDHandler interface for that can handle valid delete tradesperson tradesperson ID params
type DeleteTradespersonTradespersonIDHandler interface {
	Handle(DeleteTradespersonTradespersonIDParams, interface{}) middleware.Responder
}

// NewDeleteTradespersonTradespersonID creates a new http.Handler for the delete tradesperson tradesperson ID operation
func NewDeleteTradespersonTradespersonID(ctx *middleware.Context, handler DeleteTradespersonTradespersonIDHandler) *DeleteTradespersonTradespersonID {
	return &DeleteTradespersonTradespersonID{Context: ctx, Handler: handler}
}

/* DeleteTradespersonTradespersonID swagger:route DELETE /tradesperson/{tradespersonId} deleteTradespersonTradespersonId

DeleteTradespersonTradespersonID delete tradesperson tradesperson ID API

*/
type DeleteTradespersonTradespersonID struct {
	Context *middleware.Context
	Handler DeleteTradespersonTradespersonIDHandler
}

func (o *DeleteTradespersonTradespersonID) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, rCtx, _ := o.Context.RouteInfo(r)
	if rCtx != nil {
		*r = *rCtx
	}
	var Params = NewDeleteTradespersonTradespersonIDParams()
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

// DeleteTradespersonTradespersonIDOKBody delete tradesperson tradesperson ID o k body
//
// swagger:model DeleteTradespersonTradespersonIDOKBody
type DeleteTradespersonTradespersonIDOKBody struct {

	// deleted
	Deleted bool `json:"deleted"`
}

// Validate validates this delete tradesperson tradesperson ID o k body
func (o *DeleteTradespersonTradespersonIDOKBody) Validate(formats strfmt.Registry) error {
	return nil
}

// ContextValidate validates this delete tradesperson tradesperson ID o k body based on context it is used
func (o *DeleteTradespersonTradespersonIDOKBody) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (o *DeleteTradespersonTradespersonIDOKBody) MarshalBinary() ([]byte, error) {
	if o == nil {
		return nil, nil
	}
	return swag.WriteJSON(o)
}

// UnmarshalBinary interface implementation
func (o *DeleteTradespersonTradespersonIDOKBody) UnmarshalBinary(b []byte) error {
	var res DeleteTradespersonTradespersonIDOKBody
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*o = res
	return nil
}
