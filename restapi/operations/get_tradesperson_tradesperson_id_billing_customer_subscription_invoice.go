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

	"redbudway-api/models"
)

// GetTradespersonTradespersonIDBillingCustomerSubscriptionInvoiceHandlerFunc turns a function with the right signature into a get tradesperson tradesperson ID billing customer subscription invoice handler
type GetTradespersonTradespersonIDBillingCustomerSubscriptionInvoiceHandlerFunc func(GetTradespersonTradespersonIDBillingCustomerSubscriptionInvoiceParams) middleware.Responder

// Handle executing the request and returning a response
func (fn GetTradespersonTradespersonIDBillingCustomerSubscriptionInvoiceHandlerFunc) Handle(params GetTradespersonTradespersonIDBillingCustomerSubscriptionInvoiceParams) middleware.Responder {
	return fn(params)
}

// GetTradespersonTradespersonIDBillingCustomerSubscriptionInvoiceHandler interface for that can handle valid get tradesperson tradesperson ID billing customer subscription invoice params
type GetTradespersonTradespersonIDBillingCustomerSubscriptionInvoiceHandler interface {
	Handle(GetTradespersonTradespersonIDBillingCustomerSubscriptionInvoiceParams) middleware.Responder
}

// NewGetTradespersonTradespersonIDBillingCustomerSubscriptionInvoice creates a new http.Handler for the get tradesperson tradesperson ID billing customer subscription invoice operation
func NewGetTradespersonTradespersonIDBillingCustomerSubscriptionInvoice(ctx *middleware.Context, handler GetTradespersonTradespersonIDBillingCustomerSubscriptionInvoiceHandler) *GetTradespersonTradespersonIDBillingCustomerSubscriptionInvoice {
	return &GetTradespersonTradespersonIDBillingCustomerSubscriptionInvoice{Context: ctx, Handler: handler}
}

/* GetTradespersonTradespersonIDBillingCustomerSubscriptionInvoice swagger:route GET /tradesperson/{tradespersonId}/billing/customer/subscription/invoice getTradespersonTradespersonIdBillingCustomerSubscriptionInvoice

GetTradespersonTradespersonIDBillingCustomerSubscriptionInvoice get tradesperson tradesperson ID billing customer subscription invoice API

*/
type GetTradespersonTradespersonIDBillingCustomerSubscriptionInvoice struct {
	Context *middleware.Context
	Handler GetTradespersonTradespersonIDBillingCustomerSubscriptionInvoiceHandler
}

func (o *GetTradespersonTradespersonIDBillingCustomerSubscriptionInvoice) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, rCtx, _ := o.Context.RouteInfo(r)
	if rCtx != nil {
		*r = *rCtx
	}
	var Params = NewGetTradespersonTradespersonIDBillingCustomerSubscriptionInvoiceParams()
	if err := o.Context.BindValidRequest(r, route, &Params); err != nil { // bind params
		o.Context.Respond(rw, r, route.Produces, route, err)
		return
	}

	res := o.Handler.Handle(Params) // actually handle the request
	o.Context.Respond(rw, r, route.Produces, route, res)

}

// GetTradespersonTradespersonIDBillingCustomerSubscriptionInvoiceBody get tradesperson tradesperson ID billing customer subscription invoice body
//
// swagger:model GetTradespersonTradespersonIDBillingCustomerSubscriptionInvoiceBody
type GetTradespersonTradespersonIDBillingCustomerSubscriptionInvoiceBody struct {

	// subscription Id
	// Required: true
	SubscriptionID *string `json:"subscriptionId"`
}

// Validate validates this get tradesperson tradesperson ID billing customer subscription invoice body
func (o *GetTradespersonTradespersonIDBillingCustomerSubscriptionInvoiceBody) Validate(formats strfmt.Registry) error {
	var res []error

	if err := o.validateSubscriptionID(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (o *GetTradespersonTradespersonIDBillingCustomerSubscriptionInvoiceBody) validateSubscriptionID(formats strfmt.Registry) error {

	if err := validate.Required("subscription"+"."+"subscriptionId", "body", o.SubscriptionID); err != nil {
		return err
	}

	return nil
}

// ContextValidate validates this get tradesperson tradesperson ID billing customer subscription invoice body based on context it is used
func (o *GetTradespersonTradespersonIDBillingCustomerSubscriptionInvoiceBody) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (o *GetTradespersonTradespersonIDBillingCustomerSubscriptionInvoiceBody) MarshalBinary() ([]byte, error) {
	if o == nil {
		return nil, nil
	}
	return swag.WriteJSON(o)
}

// UnmarshalBinary interface implementation
func (o *GetTradespersonTradespersonIDBillingCustomerSubscriptionInvoiceBody) UnmarshalBinary(b []byte) error {
	var res GetTradespersonTradespersonIDBillingCustomerSubscriptionInvoiceBody
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*o = res
	return nil
}

// GetTradespersonTradespersonIDBillingCustomerSubscriptionInvoiceOKBody get tradesperson tradesperson ID billing customer subscription invoice o k body
//
// swagger:model GetTradespersonTradespersonIDBillingCustomerSubscriptionInvoiceOKBody
type GetTradespersonTradespersonIDBillingCustomerSubscriptionInvoiceOKBody struct {

	// created
	Created string `json:"created,omitempty"`

	// customer
	Customer *GetTradespersonTradespersonIDBillingCustomerSubscriptionInvoiceOKBodyCustomer `json:"customer,omitempty"`

	// description
	Description string `json:"description,omitempty"`

	// interval
	Interval string `json:"interval,omitempty"`

	// number
	Number string `json:"number,omitempty"`

	// paid
	Paid string `json:"paid,omitempty"`

	// refund
	Refund string `json:"refund,omitempty"`

	// service
	Service *GetTradespersonTradespersonIDBillingCustomerSubscriptionInvoiceOKBodyService `json:"service,omitempty"`

	// status
	Status string `json:"status,omitempty"`

	// total
	Total string `json:"total,omitempty"`

	// url
	URL string `json:"url,omitempty"`
}

// Validate validates this get tradesperson tradesperson ID billing customer subscription invoice o k body
func (o *GetTradespersonTradespersonIDBillingCustomerSubscriptionInvoiceOKBody) Validate(formats strfmt.Registry) error {
	var res []error

	if err := o.validateCustomer(formats); err != nil {
		res = append(res, err)
	}

	if err := o.validateService(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (o *GetTradespersonTradespersonIDBillingCustomerSubscriptionInvoiceOKBody) validateCustomer(formats strfmt.Registry) error {
	if swag.IsZero(o.Customer) { // not required
		return nil
	}

	if o.Customer != nil {
		if err := o.Customer.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("getTradespersonTradespersonIdBillingCustomerSubscriptionInvoiceOK" + "." + "customer")
			} else if ce, ok := err.(*errors.CompositeError); ok {
				return ce.ValidateName("getTradespersonTradespersonIdBillingCustomerSubscriptionInvoiceOK" + "." + "customer")
			}
			return err
		}
	}

	return nil
}

func (o *GetTradespersonTradespersonIDBillingCustomerSubscriptionInvoiceOKBody) validateService(formats strfmt.Registry) error {
	if swag.IsZero(o.Service) { // not required
		return nil
	}

	if o.Service != nil {
		if err := o.Service.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("getTradespersonTradespersonIdBillingCustomerSubscriptionInvoiceOK" + "." + "service")
			} else if ce, ok := err.(*errors.CompositeError); ok {
				return ce.ValidateName("getTradespersonTradespersonIdBillingCustomerSubscriptionInvoiceOK" + "." + "service")
			}
			return err
		}
	}

	return nil
}

// ContextValidate validate this get tradesperson tradesperson ID billing customer subscription invoice o k body based on the context it is used
func (o *GetTradespersonTradespersonIDBillingCustomerSubscriptionInvoiceOKBody) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	var res []error

	if err := o.contextValidateCustomer(ctx, formats); err != nil {
		res = append(res, err)
	}

	if err := o.contextValidateService(ctx, formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (o *GetTradespersonTradespersonIDBillingCustomerSubscriptionInvoiceOKBody) contextValidateCustomer(ctx context.Context, formats strfmt.Registry) error {

	if o.Customer != nil {
		if err := o.Customer.ContextValidate(ctx, formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("getTradespersonTradespersonIdBillingCustomerSubscriptionInvoiceOK" + "." + "customer")
			} else if ce, ok := err.(*errors.CompositeError); ok {
				return ce.ValidateName("getTradespersonTradespersonIdBillingCustomerSubscriptionInvoiceOK" + "." + "customer")
			}
			return err
		}
	}

	return nil
}

func (o *GetTradespersonTradespersonIDBillingCustomerSubscriptionInvoiceOKBody) contextValidateService(ctx context.Context, formats strfmt.Registry) error {

	if o.Service != nil {
		if err := o.Service.ContextValidate(ctx, formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("getTradespersonTradespersonIdBillingCustomerSubscriptionInvoiceOK" + "." + "service")
			} else if ce, ok := err.(*errors.CompositeError); ok {
				return ce.ValidateName("getTradespersonTradespersonIdBillingCustomerSubscriptionInvoiceOK" + "." + "service")
			}
			return err
		}
	}

	return nil
}

// MarshalBinary interface implementation
func (o *GetTradespersonTradespersonIDBillingCustomerSubscriptionInvoiceOKBody) MarshalBinary() ([]byte, error) {
	if o == nil {
		return nil, nil
	}
	return swag.WriteJSON(o)
}

// UnmarshalBinary interface implementation
func (o *GetTradespersonTradespersonIDBillingCustomerSubscriptionInvoiceOKBody) UnmarshalBinary(b []byte) error {
	var res GetTradespersonTradespersonIDBillingCustomerSubscriptionInvoiceOKBody
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*o = res
	return nil
}

// GetTradespersonTradespersonIDBillingCustomerSubscriptionInvoiceOKBodyCustomer get tradesperson tradesperson ID billing customer subscription invoice o k body customer
//
// swagger:model GetTradespersonTradespersonIDBillingCustomerSubscriptionInvoiceOKBodyCustomer
type GetTradespersonTradespersonIDBillingCustomerSubscriptionInvoiceOKBodyCustomer struct {

	// address
	Address *models.Address `json:"address,omitempty"`

	// name
	Name string `json:"name,omitempty"`
}

// Validate validates this get tradesperson tradesperson ID billing customer subscription invoice o k body customer
func (o *GetTradespersonTradespersonIDBillingCustomerSubscriptionInvoiceOKBodyCustomer) Validate(formats strfmt.Registry) error {
	var res []error

	if err := o.validateAddress(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (o *GetTradespersonTradespersonIDBillingCustomerSubscriptionInvoiceOKBodyCustomer) validateAddress(formats strfmt.Registry) error {
	if swag.IsZero(o.Address) { // not required
		return nil
	}

	if o.Address != nil {
		if err := o.Address.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("getTradespersonTradespersonIdBillingCustomerSubscriptionInvoiceOK" + "." + "customer" + "." + "address")
			} else if ce, ok := err.(*errors.CompositeError); ok {
				return ce.ValidateName("getTradespersonTradespersonIdBillingCustomerSubscriptionInvoiceOK" + "." + "customer" + "." + "address")
			}
			return err
		}
	}

	return nil
}

// ContextValidate validate this get tradesperson tradesperson ID billing customer subscription invoice o k body customer based on the context it is used
func (o *GetTradespersonTradespersonIDBillingCustomerSubscriptionInvoiceOKBodyCustomer) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	var res []error

	if err := o.contextValidateAddress(ctx, formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (o *GetTradespersonTradespersonIDBillingCustomerSubscriptionInvoiceOKBodyCustomer) contextValidateAddress(ctx context.Context, formats strfmt.Registry) error {

	if o.Address != nil {
		if err := o.Address.ContextValidate(ctx, formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("getTradespersonTradespersonIdBillingCustomerSubscriptionInvoiceOK" + "." + "customer" + "." + "address")
			} else if ce, ok := err.(*errors.CompositeError); ok {
				return ce.ValidateName("getTradespersonTradespersonIdBillingCustomerSubscriptionInvoiceOK" + "." + "customer" + "." + "address")
			}
			return err
		}
	}

	return nil
}

// MarshalBinary interface implementation
func (o *GetTradespersonTradespersonIDBillingCustomerSubscriptionInvoiceOKBodyCustomer) MarshalBinary() ([]byte, error) {
	if o == nil {
		return nil, nil
	}
	return swag.WriteJSON(o)
}

// UnmarshalBinary interface implementation
func (o *GetTradespersonTradespersonIDBillingCustomerSubscriptionInvoiceOKBodyCustomer) UnmarshalBinary(b []byte) error {
	var res GetTradespersonTradespersonIDBillingCustomerSubscriptionInvoiceOKBodyCustomer
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*o = res
	return nil
}

// GetTradespersonTradespersonIDBillingCustomerSubscriptionInvoiceOKBodyService get tradesperson tradesperson ID billing customer subscription invoice o k body service
//
// swagger:model GetTradespersonTradespersonIDBillingCustomerSubscriptionInvoiceOKBodyService
type GetTradespersonTradespersonIDBillingCustomerSubscriptionInvoiceOKBodyService struct {

	// description
	Description string `json:"description,omitempty"`

	// title
	Title string `json:"title,omitempty"`
}

// Validate validates this get tradesperson tradesperson ID billing customer subscription invoice o k body service
func (o *GetTradespersonTradespersonIDBillingCustomerSubscriptionInvoiceOKBodyService) Validate(formats strfmt.Registry) error {
	return nil
}

// ContextValidate validates this get tradesperson tradesperson ID billing customer subscription invoice o k body service based on context it is used
func (o *GetTradespersonTradespersonIDBillingCustomerSubscriptionInvoiceOKBodyService) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (o *GetTradespersonTradespersonIDBillingCustomerSubscriptionInvoiceOKBodyService) MarshalBinary() ([]byte, error) {
	if o == nil {
		return nil, nil
	}
	return swag.WriteJSON(o)
}

// UnmarshalBinary interface implementation
func (o *GetTradespersonTradespersonIDBillingCustomerSubscriptionInvoiceOKBodyService) UnmarshalBinary(b []byte) error {
	var res GetTradespersonTradespersonIDBillingCustomerSubscriptionInvoiceOKBodyService
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*o = res
	return nil
}
