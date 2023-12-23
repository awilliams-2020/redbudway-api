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

// DeleteTradespersonTradespersonIDCouponCouponIDHandlerFunc turns a function with the right signature into a delete tradesperson tradesperson ID coupon coupon ID handler
type DeleteTradespersonTradespersonIDCouponCouponIDHandlerFunc func(DeleteTradespersonTradespersonIDCouponCouponIDParams, interface{}) middleware.Responder

// Handle executing the request and returning a response
func (fn DeleteTradespersonTradespersonIDCouponCouponIDHandlerFunc) Handle(params DeleteTradespersonTradespersonIDCouponCouponIDParams, principal interface{}) middleware.Responder {
	return fn(params, principal)
}

// DeleteTradespersonTradespersonIDCouponCouponIDHandler interface for that can handle valid delete tradesperson tradesperson ID coupon coupon ID params
type DeleteTradespersonTradespersonIDCouponCouponIDHandler interface {
	Handle(DeleteTradespersonTradespersonIDCouponCouponIDParams, interface{}) middleware.Responder
}

// NewDeleteTradespersonTradespersonIDCouponCouponID creates a new http.Handler for the delete tradesperson tradesperson ID coupon coupon ID operation
func NewDeleteTradespersonTradespersonIDCouponCouponID(ctx *middleware.Context, handler DeleteTradespersonTradespersonIDCouponCouponIDHandler) *DeleteTradespersonTradespersonIDCouponCouponID {
	return &DeleteTradespersonTradespersonIDCouponCouponID{Context: ctx, Handler: handler}
}

/* DeleteTradespersonTradespersonIDCouponCouponID swagger:route DELETE /tradesperson/{tradespersonId}/coupon/{couponId} deleteTradespersonTradespersonIdCouponCouponId

DeleteTradespersonTradespersonIDCouponCouponID delete tradesperson tradesperson ID coupon coupon ID API

*/
type DeleteTradespersonTradespersonIDCouponCouponID struct {
	Context *middleware.Context
	Handler DeleteTradespersonTradespersonIDCouponCouponIDHandler
}

func (o *DeleteTradespersonTradespersonIDCouponCouponID) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, rCtx, _ := o.Context.RouteInfo(r)
	if rCtx != nil {
		*r = *rCtx
	}
	var Params = NewDeleteTradespersonTradespersonIDCouponCouponIDParams()
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

// DeleteTradespersonTradespersonIDCouponCouponIDOKBody delete tradesperson tradesperson ID coupon coupon ID o k body
//
// swagger:model DeleteTradespersonTradespersonIDCouponCouponIDOKBody
type DeleteTradespersonTradespersonIDCouponCouponIDOKBody struct {

	// deleted
	Deleted bool `json:"deleted"`
}

// Validate validates this delete tradesperson tradesperson ID coupon coupon ID o k body
func (o *DeleteTradespersonTradespersonIDCouponCouponIDOKBody) Validate(formats strfmt.Registry) error {
	return nil
}

// ContextValidate validates this delete tradesperson tradesperson ID coupon coupon ID o k body based on context it is used
func (o *DeleteTradespersonTradespersonIDCouponCouponIDOKBody) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (o *DeleteTradespersonTradespersonIDCouponCouponIDOKBody) MarshalBinary() ([]byte, error) {
	if o == nil {
		return nil, nil
	}
	return swag.WriteJSON(o)
}

// UnmarshalBinary interface implementation
func (o *DeleteTradespersonTradespersonIDCouponCouponIDOKBody) UnmarshalBinary(b []byte) error {
	var res DeleteTradespersonTradespersonIDCouponCouponIDOKBody
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*o = res
	return nil
}
