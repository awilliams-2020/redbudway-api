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

// PostTradespersonTradespersonIDBillingManualInvoiceInvoiceIDRefundHandlerFunc turns a function with the right signature into a post tradesperson tradesperson ID billing manual invoice invoice ID refund handler
type PostTradespersonTradespersonIDBillingManualInvoiceInvoiceIDRefundHandlerFunc func(PostTradespersonTradespersonIDBillingManualInvoiceInvoiceIDRefundParams, interface{}) middleware.Responder

// Handle executing the request and returning a response
func (fn PostTradespersonTradespersonIDBillingManualInvoiceInvoiceIDRefundHandlerFunc) Handle(params PostTradespersonTradespersonIDBillingManualInvoiceInvoiceIDRefundParams, principal interface{}) middleware.Responder {
	return fn(params, principal)
}

// PostTradespersonTradespersonIDBillingManualInvoiceInvoiceIDRefundHandler interface for that can handle valid post tradesperson tradesperson ID billing manual invoice invoice ID refund params
type PostTradespersonTradespersonIDBillingManualInvoiceInvoiceIDRefundHandler interface {
	Handle(PostTradespersonTradespersonIDBillingManualInvoiceInvoiceIDRefundParams, interface{}) middleware.Responder
}

// NewPostTradespersonTradespersonIDBillingManualInvoiceInvoiceIDRefund creates a new http.Handler for the post tradesperson tradesperson ID billing manual invoice invoice ID refund operation
func NewPostTradespersonTradespersonIDBillingManualInvoiceInvoiceIDRefund(ctx *middleware.Context, handler PostTradespersonTradespersonIDBillingManualInvoiceInvoiceIDRefundHandler) *PostTradespersonTradespersonIDBillingManualInvoiceInvoiceIDRefund {
	return &PostTradespersonTradespersonIDBillingManualInvoiceInvoiceIDRefund{Context: ctx, Handler: handler}
}

/* PostTradespersonTradespersonIDBillingManualInvoiceInvoiceIDRefund swagger:route POST /tradesperson/{tradespersonId}/billing/manual-invoice/{invoiceId}/refund postTradespersonTradespersonIdBillingManualInvoiceInvoiceIdRefund

PostTradespersonTradespersonIDBillingManualInvoiceInvoiceIDRefund post tradesperson tradesperson ID billing manual invoice invoice ID refund API

*/
type PostTradespersonTradespersonIDBillingManualInvoiceInvoiceIDRefund struct {
	Context *middleware.Context
	Handler PostTradespersonTradespersonIDBillingManualInvoiceInvoiceIDRefundHandler
}

func (o *PostTradespersonTradespersonIDBillingManualInvoiceInvoiceIDRefund) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, rCtx, _ := o.Context.RouteInfo(r)
	if rCtx != nil {
		*r = *rCtx
	}
	var Params = NewPostTradespersonTradespersonIDBillingManualInvoiceInvoiceIDRefundParams()
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

// PostTradespersonTradespersonIDBillingManualInvoiceInvoiceIDRefundOKBody post tradesperson tradesperson ID billing manual invoice invoice ID refund o k body
//
// swagger:model PostTradespersonTradespersonIDBillingManualInvoiceInvoiceIDRefundOKBody
type PostTradespersonTradespersonIDBillingManualInvoiceInvoiceIDRefundOKBody struct {

	// refunded
	Refunded bool `json:"refunded"`
}

// Validate validates this post tradesperson tradesperson ID billing manual invoice invoice ID refund o k body
func (o *PostTradespersonTradespersonIDBillingManualInvoiceInvoiceIDRefundOKBody) Validate(formats strfmt.Registry) error {
	return nil
}

// ContextValidate validates this post tradesperson tradesperson ID billing manual invoice invoice ID refund o k body based on context it is used
func (o *PostTradespersonTradespersonIDBillingManualInvoiceInvoiceIDRefundOKBody) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (o *PostTradespersonTradespersonIDBillingManualInvoiceInvoiceIDRefundOKBody) MarshalBinary() ([]byte, error) {
	if o == nil {
		return nil, nil
	}
	return swag.WriteJSON(o)
}

// UnmarshalBinary interface implementation
func (o *PostTradespersonTradespersonIDBillingManualInvoiceInvoiceIDRefundOKBody) UnmarshalBinary(b []byte) error {
	var res PostTradespersonTradespersonIDBillingManualInvoiceInvoiceIDRefundOKBody
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*o = res
	return nil
}
