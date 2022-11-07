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

// PostTradespersonTradespersonIDBillingQuoteQuoteIDInvoiceInvoiceIDUncollectibleHandlerFunc turns a function with the right signature into a post tradesperson tradesperson ID billing quote quote ID invoice invoice ID uncollectible handler
type PostTradespersonTradespersonIDBillingQuoteQuoteIDInvoiceInvoiceIDUncollectibleHandlerFunc func(PostTradespersonTradespersonIDBillingQuoteQuoteIDInvoiceInvoiceIDUncollectibleParams, interface{}) middleware.Responder

// Handle executing the request and returning a response
func (fn PostTradespersonTradespersonIDBillingQuoteQuoteIDInvoiceInvoiceIDUncollectibleHandlerFunc) Handle(params PostTradespersonTradespersonIDBillingQuoteQuoteIDInvoiceInvoiceIDUncollectibleParams, principal interface{}) middleware.Responder {
	return fn(params, principal)
}

// PostTradespersonTradespersonIDBillingQuoteQuoteIDInvoiceInvoiceIDUncollectibleHandler interface for that can handle valid post tradesperson tradesperson ID billing quote quote ID invoice invoice ID uncollectible params
type PostTradespersonTradespersonIDBillingQuoteQuoteIDInvoiceInvoiceIDUncollectibleHandler interface {
	Handle(PostTradespersonTradespersonIDBillingQuoteQuoteIDInvoiceInvoiceIDUncollectibleParams, interface{}) middleware.Responder
}

// NewPostTradespersonTradespersonIDBillingQuoteQuoteIDInvoiceInvoiceIDUncollectible creates a new http.Handler for the post tradesperson tradesperson ID billing quote quote ID invoice invoice ID uncollectible operation
func NewPostTradespersonTradespersonIDBillingQuoteQuoteIDInvoiceInvoiceIDUncollectible(ctx *middleware.Context, handler PostTradespersonTradespersonIDBillingQuoteQuoteIDInvoiceInvoiceIDUncollectibleHandler) *PostTradespersonTradespersonIDBillingQuoteQuoteIDInvoiceInvoiceIDUncollectible {
	return &PostTradespersonTradespersonIDBillingQuoteQuoteIDInvoiceInvoiceIDUncollectible{Context: ctx, Handler: handler}
}

/* PostTradespersonTradespersonIDBillingQuoteQuoteIDInvoiceInvoiceIDUncollectible swagger:route POST /tradesperson/{tradespersonId}/billing/quote/{quoteId}/invoice/{invoiceId}/uncollectible postTradespersonTradespersonIdBillingQuoteQuoteIdInvoiceInvoiceIdUncollectible

PostTradespersonTradespersonIDBillingQuoteQuoteIDInvoiceInvoiceIDUncollectible post tradesperson tradesperson ID billing quote quote ID invoice invoice ID uncollectible API

*/
type PostTradespersonTradespersonIDBillingQuoteQuoteIDInvoiceInvoiceIDUncollectible struct {
	Context *middleware.Context
	Handler PostTradespersonTradespersonIDBillingQuoteQuoteIDInvoiceInvoiceIDUncollectibleHandler
}

func (o *PostTradespersonTradespersonIDBillingQuoteQuoteIDInvoiceInvoiceIDUncollectible) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, rCtx, _ := o.Context.RouteInfo(r)
	if rCtx != nil {
		*r = *rCtx
	}
	var Params = NewPostTradespersonTradespersonIDBillingQuoteQuoteIDInvoiceInvoiceIDUncollectibleParams()
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

// PostTradespersonTradespersonIDBillingQuoteQuoteIDInvoiceInvoiceIDUncollectibleOKBody post tradesperson tradesperson ID billing quote quote ID invoice invoice ID uncollectible o k body
//
// swagger:model PostTradespersonTradespersonIDBillingQuoteQuoteIDInvoiceInvoiceIDUncollectibleOKBody
type PostTradespersonTradespersonIDBillingQuoteQuoteIDInvoiceInvoiceIDUncollectibleOKBody struct {

	// uncollectible
	Uncollectible bool `json:"uncollectible"`
}

// Validate validates this post tradesperson tradesperson ID billing quote quote ID invoice invoice ID uncollectible o k body
func (o *PostTradespersonTradespersonIDBillingQuoteQuoteIDInvoiceInvoiceIDUncollectibleOKBody) Validate(formats strfmt.Registry) error {
	return nil
}

// ContextValidate validates this post tradesperson tradesperson ID billing quote quote ID invoice invoice ID uncollectible o k body based on context it is used
func (o *PostTradespersonTradespersonIDBillingQuoteQuoteIDInvoiceInvoiceIDUncollectibleOKBody) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (o *PostTradespersonTradespersonIDBillingQuoteQuoteIDInvoiceInvoiceIDUncollectibleOKBody) MarshalBinary() ([]byte, error) {
	if o == nil {
		return nil, nil
	}
	return swag.WriteJSON(o)
}

// UnmarshalBinary interface implementation
func (o *PostTradespersonTradespersonIDBillingQuoteQuoteIDInvoiceInvoiceIDUncollectibleOKBody) UnmarshalBinary(b []byte) error {
	var res PostTradespersonTradespersonIDBillingQuoteQuoteIDInvoiceInvoiceIDUncollectibleOKBody
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*o = res
	return nil
}
