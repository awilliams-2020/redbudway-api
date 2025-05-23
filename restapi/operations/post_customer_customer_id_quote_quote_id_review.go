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

// PostCustomerCustomerIDQuoteQuoteIDReviewHandlerFunc turns a function with the right signature into a post customer customer ID quote quote ID review handler
type PostCustomerCustomerIDQuoteQuoteIDReviewHandlerFunc func(PostCustomerCustomerIDQuoteQuoteIDReviewParams, interface{}) middleware.Responder

// Handle executing the request and returning a response
func (fn PostCustomerCustomerIDQuoteQuoteIDReviewHandlerFunc) Handle(params PostCustomerCustomerIDQuoteQuoteIDReviewParams, principal interface{}) middleware.Responder {
	return fn(params, principal)
}

// PostCustomerCustomerIDQuoteQuoteIDReviewHandler interface for that can handle valid post customer customer ID quote quote ID review params
type PostCustomerCustomerIDQuoteQuoteIDReviewHandler interface {
	Handle(PostCustomerCustomerIDQuoteQuoteIDReviewParams, interface{}) middleware.Responder
}

// NewPostCustomerCustomerIDQuoteQuoteIDReview creates a new http.Handler for the post customer customer ID quote quote ID review operation
func NewPostCustomerCustomerIDQuoteQuoteIDReview(ctx *middleware.Context, handler PostCustomerCustomerIDQuoteQuoteIDReviewHandler) *PostCustomerCustomerIDQuoteQuoteIDReview {
	return &PostCustomerCustomerIDQuoteQuoteIDReview{Context: ctx, Handler: handler}
}

/* PostCustomerCustomerIDQuoteQuoteIDReview swagger:route POST /customer/{customerId}/quote/{quoteId}/review postCustomerCustomerIdQuoteQuoteIdReview

PostCustomerCustomerIDQuoteQuoteIDReview post customer customer ID quote quote ID review API

*/
type PostCustomerCustomerIDQuoteQuoteIDReview struct {
	Context *middleware.Context
	Handler PostCustomerCustomerIDQuoteQuoteIDReviewHandler
}

func (o *PostCustomerCustomerIDQuoteQuoteIDReview) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, rCtx, _ := o.Context.RouteInfo(r)
	if rCtx != nil {
		*r = *rCtx
	}
	var Params = NewPostCustomerCustomerIDQuoteQuoteIDReviewParams()
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

// PostCustomerCustomerIDQuoteQuoteIDReviewBody post customer customer ID quote quote ID review body
//
// swagger:model PostCustomerCustomerIDQuoteQuoteIDReviewBody
type PostCustomerCustomerIDQuoteQuoteIDReviewBody struct {

	// message
	Message string `json:"message,omitempty"`

	// rating
	Rating int64 `json:"rating,omitempty"`
}

// Validate validates this post customer customer ID quote quote ID review body
func (o *PostCustomerCustomerIDQuoteQuoteIDReviewBody) Validate(formats strfmt.Registry) error {
	return nil
}

// ContextValidate validates this post customer customer ID quote quote ID review body based on context it is used
func (o *PostCustomerCustomerIDQuoteQuoteIDReviewBody) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (o *PostCustomerCustomerIDQuoteQuoteIDReviewBody) MarshalBinary() ([]byte, error) {
	if o == nil {
		return nil, nil
	}
	return swag.WriteJSON(o)
}

// UnmarshalBinary interface implementation
func (o *PostCustomerCustomerIDQuoteQuoteIDReviewBody) UnmarshalBinary(b []byte) error {
	var res PostCustomerCustomerIDQuoteQuoteIDReviewBody
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*o = res
	return nil
}

// PostCustomerCustomerIDQuoteQuoteIDReviewOKBody post customer customer ID quote quote ID review o k body
//
// swagger:model PostCustomerCustomerIDQuoteQuoteIDReviewOKBody
type PostCustomerCustomerIDQuoteQuoteIDReviewOKBody struct {

	// rated
	Rated bool `json:"rated"`
}

// Validate validates this post customer customer ID quote quote ID review o k body
func (o *PostCustomerCustomerIDQuoteQuoteIDReviewOKBody) Validate(formats strfmt.Registry) error {
	return nil
}

// ContextValidate validates this post customer customer ID quote quote ID review o k body based on context it is used
func (o *PostCustomerCustomerIDQuoteQuoteIDReviewOKBody) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (o *PostCustomerCustomerIDQuoteQuoteIDReviewOKBody) MarshalBinary() ([]byte, error) {
	if o == nil {
		return nil, nil
	}
	return swag.WriteJSON(o)
}

// UnmarshalBinary interface implementation
func (o *PostCustomerCustomerIDQuoteQuoteIDReviewOKBody) UnmarshalBinary(b []byte) error {
	var res PostCustomerCustomerIDQuoteQuoteIDReviewOKBody
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*o = res
	return nil
}
