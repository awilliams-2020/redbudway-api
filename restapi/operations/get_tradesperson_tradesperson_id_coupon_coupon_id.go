// Code generated by go-swagger; DO NOT EDIT.

package operations

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"net/http"

	"github.com/go-openapi/runtime/middleware"
)

// GetTradespersonTradespersonIDCouponCouponIDHandlerFunc turns a function with the right signature into a get tradesperson tradesperson ID coupon coupon ID handler
type GetTradespersonTradespersonIDCouponCouponIDHandlerFunc func(GetTradespersonTradespersonIDCouponCouponIDParams, interface{}) middleware.Responder

// Handle executing the request and returning a response
func (fn GetTradespersonTradespersonIDCouponCouponIDHandlerFunc) Handle(params GetTradespersonTradespersonIDCouponCouponIDParams, principal interface{}) middleware.Responder {
	return fn(params, principal)
}

// GetTradespersonTradespersonIDCouponCouponIDHandler interface for that can handle valid get tradesperson tradesperson ID coupon coupon ID params
type GetTradespersonTradespersonIDCouponCouponIDHandler interface {
	Handle(GetTradespersonTradespersonIDCouponCouponIDParams, interface{}) middleware.Responder
}

// NewGetTradespersonTradespersonIDCouponCouponID creates a new http.Handler for the get tradesperson tradesperson ID coupon coupon ID operation
func NewGetTradespersonTradespersonIDCouponCouponID(ctx *middleware.Context, handler GetTradespersonTradespersonIDCouponCouponIDHandler) *GetTradespersonTradespersonIDCouponCouponID {
	return &GetTradespersonTradespersonIDCouponCouponID{Context: ctx, Handler: handler}
}

/* GetTradespersonTradespersonIDCouponCouponID swagger:route GET /tradesperson/{tradespersonId}/coupon/{couponId} getTradespersonTradespersonIdCouponCouponId

GetTradespersonTradespersonIDCouponCouponID get tradesperson tradesperson ID coupon coupon ID API

*/
type GetTradespersonTradespersonIDCouponCouponID struct {
	Context *middleware.Context
	Handler GetTradespersonTradespersonIDCouponCouponIDHandler
}

func (o *GetTradespersonTradespersonIDCouponCouponID) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, rCtx, _ := o.Context.RouteInfo(r)
	if rCtx != nil {
		*r = *rCtx
	}
	var Params = NewGetTradespersonTradespersonIDCouponCouponIDParams()
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
