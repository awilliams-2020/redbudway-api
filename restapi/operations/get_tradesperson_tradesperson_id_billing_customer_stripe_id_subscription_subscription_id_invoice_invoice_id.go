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

	"redbudway-api/models"
)

// GetTradespersonTradespersonIDBillingCustomerStripeIDSubscriptionSubscriptionIDInvoiceInvoiceIDHandlerFunc turns a function with the right signature into a get tradesperson tradesperson ID billing customer stripe ID subscription subscription ID invoice invoice ID handler
type GetTradespersonTradespersonIDBillingCustomerStripeIDSubscriptionSubscriptionIDInvoiceInvoiceIDHandlerFunc func(GetTradespersonTradespersonIDBillingCustomerStripeIDSubscriptionSubscriptionIDInvoiceInvoiceIDParams, interface{}) middleware.Responder

// Handle executing the request and returning a response
func (fn GetTradespersonTradespersonIDBillingCustomerStripeIDSubscriptionSubscriptionIDInvoiceInvoiceIDHandlerFunc) Handle(params GetTradespersonTradespersonIDBillingCustomerStripeIDSubscriptionSubscriptionIDInvoiceInvoiceIDParams, principal interface{}) middleware.Responder {
	return fn(params, principal)
}

// GetTradespersonTradespersonIDBillingCustomerStripeIDSubscriptionSubscriptionIDInvoiceInvoiceIDHandler interface for that can handle valid get tradesperson tradesperson ID billing customer stripe ID subscription subscription ID invoice invoice ID params
type GetTradespersonTradespersonIDBillingCustomerStripeIDSubscriptionSubscriptionIDInvoiceInvoiceIDHandler interface {
	Handle(GetTradespersonTradespersonIDBillingCustomerStripeIDSubscriptionSubscriptionIDInvoiceInvoiceIDParams, interface{}) middleware.Responder
}

// NewGetTradespersonTradespersonIDBillingCustomerStripeIDSubscriptionSubscriptionIDInvoiceInvoiceID creates a new http.Handler for the get tradesperson tradesperson ID billing customer stripe ID subscription subscription ID invoice invoice ID operation
func NewGetTradespersonTradespersonIDBillingCustomerStripeIDSubscriptionSubscriptionIDInvoiceInvoiceID(ctx *middleware.Context, handler GetTradespersonTradespersonIDBillingCustomerStripeIDSubscriptionSubscriptionIDInvoiceInvoiceIDHandler) *GetTradespersonTradespersonIDBillingCustomerStripeIDSubscriptionSubscriptionIDInvoiceInvoiceID {
	return &GetTradespersonTradespersonIDBillingCustomerStripeIDSubscriptionSubscriptionIDInvoiceInvoiceID{Context: ctx, Handler: handler}
}

/* GetTradespersonTradespersonIDBillingCustomerStripeIDSubscriptionSubscriptionIDInvoiceInvoiceID swagger:route GET /tradesperson/{tradespersonId}/billing/customer/{stripeId}/subscription/{subscriptionId}/invoice/{invoiceId} getTradespersonTradespersonIdBillingCustomerStripeIdSubscriptionSubscriptionIdInvoiceInvoiceId

GetTradespersonTradespersonIDBillingCustomerStripeIDSubscriptionSubscriptionIDInvoiceInvoiceID get tradesperson tradesperson ID billing customer stripe ID subscription subscription ID invoice invoice ID API

*/
type GetTradespersonTradespersonIDBillingCustomerStripeIDSubscriptionSubscriptionIDInvoiceInvoiceID struct {
	Context *middleware.Context
	Handler GetTradespersonTradespersonIDBillingCustomerStripeIDSubscriptionSubscriptionIDInvoiceInvoiceIDHandler
}

func (o *GetTradespersonTradespersonIDBillingCustomerStripeIDSubscriptionSubscriptionIDInvoiceInvoiceID) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, rCtx, _ := o.Context.RouteInfo(r)
	if rCtx != nil {
		*r = *rCtx
	}
	var Params = NewGetTradespersonTradespersonIDBillingCustomerStripeIDSubscriptionSubscriptionIDInvoiceInvoiceIDParams()
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

// GetTradespersonTradespersonIDBillingCustomerStripeIDSubscriptionSubscriptionIDInvoiceInvoiceIDOKBody get tradesperson tradesperson ID billing customer stripe ID subscription subscription ID invoice invoice ID o k body
//
// swagger:model GetTradespersonTradespersonIDBillingCustomerStripeIDSubscriptionSubscriptionIDInvoiceInvoiceIDOKBody
type GetTradespersonTradespersonIDBillingCustomerStripeIDSubscriptionSubscriptionIDInvoiceInvoiceIDOKBody struct {

	// created
	Created int64 `json:"created,omitempty"`

	// customer
	Customer *models.Customer `json:"customer,omitempty"`

	// description
	Description string `json:"description,omitempty"`

	// interval
	Interval string `json:"interval,omitempty"`

	// number
	Number string `json:"number,omitempty"`

	// pdf
	Pdf string `json:"pdf,omitempty"`

	// refunded
	Refunded int64 `json:"refunded,omitempty"`

	// service
	Service *GetTradespersonTradespersonIDBillingCustomerStripeIDSubscriptionSubscriptionIDInvoiceInvoiceIDOKBodyService `json:"service,omitempty"`

	// status
	Status string `json:"status,omitempty"`

	// time slot
	TimeSlot *models.TimeSlot `json:"timeSlot,omitempty"`

	// total
	Total int64 `json:"total,omitempty"`

	// url
	URL string `json:"url,omitempty"`
}

// Validate validates this get tradesperson tradesperson ID billing customer stripe ID subscription subscription ID invoice invoice ID o k body
func (o *GetTradespersonTradespersonIDBillingCustomerStripeIDSubscriptionSubscriptionIDInvoiceInvoiceIDOKBody) Validate(formats strfmt.Registry) error {
	var res []error

	if err := o.validateCustomer(formats); err != nil {
		res = append(res, err)
	}

	if err := o.validateService(formats); err != nil {
		res = append(res, err)
	}

	if err := o.validateTimeSlot(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (o *GetTradespersonTradespersonIDBillingCustomerStripeIDSubscriptionSubscriptionIDInvoiceInvoiceIDOKBody) validateCustomer(formats strfmt.Registry) error {
	if swag.IsZero(o.Customer) { // not required
		return nil
	}

	if o.Customer != nil {
		if err := o.Customer.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("getTradespersonTradespersonIdBillingCustomerStripeIdSubscriptionSubscriptionIdInvoiceInvoiceIdOK" + "." + "customer")
			} else if ce, ok := err.(*errors.CompositeError); ok {
				return ce.ValidateName("getTradespersonTradespersonIdBillingCustomerStripeIdSubscriptionSubscriptionIdInvoiceInvoiceIdOK" + "." + "customer")
			}
			return err
		}
	}

	return nil
}

func (o *GetTradespersonTradespersonIDBillingCustomerStripeIDSubscriptionSubscriptionIDInvoiceInvoiceIDOKBody) validateService(formats strfmt.Registry) error {
	if swag.IsZero(o.Service) { // not required
		return nil
	}

	if o.Service != nil {
		if err := o.Service.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("getTradespersonTradespersonIdBillingCustomerStripeIdSubscriptionSubscriptionIdInvoiceInvoiceIdOK" + "." + "service")
			} else if ce, ok := err.(*errors.CompositeError); ok {
				return ce.ValidateName("getTradespersonTradespersonIdBillingCustomerStripeIdSubscriptionSubscriptionIdInvoiceInvoiceIdOK" + "." + "service")
			}
			return err
		}
	}

	return nil
}

func (o *GetTradespersonTradespersonIDBillingCustomerStripeIDSubscriptionSubscriptionIDInvoiceInvoiceIDOKBody) validateTimeSlot(formats strfmt.Registry) error {
	if swag.IsZero(o.TimeSlot) { // not required
		return nil
	}

	if o.TimeSlot != nil {
		if err := o.TimeSlot.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("getTradespersonTradespersonIdBillingCustomerStripeIdSubscriptionSubscriptionIdInvoiceInvoiceIdOK" + "." + "timeSlot")
			} else if ce, ok := err.(*errors.CompositeError); ok {
				return ce.ValidateName("getTradespersonTradespersonIdBillingCustomerStripeIdSubscriptionSubscriptionIdInvoiceInvoiceIdOK" + "." + "timeSlot")
			}
			return err
		}
	}

	return nil
}

// ContextValidate validate this get tradesperson tradesperson ID billing customer stripe ID subscription subscription ID invoice invoice ID o k body based on the context it is used
func (o *GetTradespersonTradespersonIDBillingCustomerStripeIDSubscriptionSubscriptionIDInvoiceInvoiceIDOKBody) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	var res []error

	if err := o.contextValidateCustomer(ctx, formats); err != nil {
		res = append(res, err)
	}

	if err := o.contextValidateService(ctx, formats); err != nil {
		res = append(res, err)
	}

	if err := o.contextValidateTimeSlot(ctx, formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (o *GetTradespersonTradespersonIDBillingCustomerStripeIDSubscriptionSubscriptionIDInvoiceInvoiceIDOKBody) contextValidateCustomer(ctx context.Context, formats strfmt.Registry) error {

	if o.Customer != nil {
		if err := o.Customer.ContextValidate(ctx, formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("getTradespersonTradespersonIdBillingCustomerStripeIdSubscriptionSubscriptionIdInvoiceInvoiceIdOK" + "." + "customer")
			} else if ce, ok := err.(*errors.CompositeError); ok {
				return ce.ValidateName("getTradespersonTradespersonIdBillingCustomerStripeIdSubscriptionSubscriptionIdInvoiceInvoiceIdOK" + "." + "customer")
			}
			return err
		}
	}

	return nil
}

func (o *GetTradespersonTradespersonIDBillingCustomerStripeIDSubscriptionSubscriptionIDInvoiceInvoiceIDOKBody) contextValidateService(ctx context.Context, formats strfmt.Registry) error {

	if o.Service != nil {
		if err := o.Service.ContextValidate(ctx, formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("getTradespersonTradespersonIdBillingCustomerStripeIdSubscriptionSubscriptionIdInvoiceInvoiceIdOK" + "." + "service")
			} else if ce, ok := err.(*errors.CompositeError); ok {
				return ce.ValidateName("getTradespersonTradespersonIdBillingCustomerStripeIdSubscriptionSubscriptionIdInvoiceInvoiceIdOK" + "." + "service")
			}
			return err
		}
	}

	return nil
}

func (o *GetTradespersonTradespersonIDBillingCustomerStripeIDSubscriptionSubscriptionIDInvoiceInvoiceIDOKBody) contextValidateTimeSlot(ctx context.Context, formats strfmt.Registry) error {

	if o.TimeSlot != nil {
		if err := o.TimeSlot.ContextValidate(ctx, formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("getTradespersonTradespersonIdBillingCustomerStripeIdSubscriptionSubscriptionIdInvoiceInvoiceIdOK" + "." + "timeSlot")
			} else if ce, ok := err.(*errors.CompositeError); ok {
				return ce.ValidateName("getTradespersonTradespersonIdBillingCustomerStripeIdSubscriptionSubscriptionIdInvoiceInvoiceIdOK" + "." + "timeSlot")
			}
			return err
		}
	}

	return nil
}

// MarshalBinary interface implementation
func (o *GetTradespersonTradespersonIDBillingCustomerStripeIDSubscriptionSubscriptionIDInvoiceInvoiceIDOKBody) MarshalBinary() ([]byte, error) {
	if o == nil {
		return nil, nil
	}
	return swag.WriteJSON(o)
}

// UnmarshalBinary interface implementation
func (o *GetTradespersonTradespersonIDBillingCustomerStripeIDSubscriptionSubscriptionIDInvoiceInvoiceIDOKBody) UnmarshalBinary(b []byte) error {
	var res GetTradespersonTradespersonIDBillingCustomerStripeIDSubscriptionSubscriptionIDInvoiceInvoiceIDOKBody
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*o = res
	return nil
}

// GetTradespersonTradespersonIDBillingCustomerStripeIDSubscriptionSubscriptionIDInvoiceInvoiceIDOKBodyService get tradesperson tradesperson ID billing customer stripe ID subscription subscription ID invoice invoice ID o k body service
//
// swagger:model GetTradespersonTradespersonIDBillingCustomerStripeIDSubscriptionSubscriptionIDInvoiceInvoiceIDOKBodyService
type GetTradespersonTradespersonIDBillingCustomerStripeIDSubscriptionSubscriptionIDInvoiceInvoiceIDOKBodyService struct {

	// description
	Description string `json:"description,omitempty"`

	// title
	Title string `json:"title,omitempty"`
}

// Validate validates this get tradesperson tradesperson ID billing customer stripe ID subscription subscription ID invoice invoice ID o k body service
func (o *GetTradespersonTradespersonIDBillingCustomerStripeIDSubscriptionSubscriptionIDInvoiceInvoiceIDOKBodyService) Validate(formats strfmt.Registry) error {
	return nil
}

// ContextValidate validates this get tradesperson tradesperson ID billing customer stripe ID subscription subscription ID invoice invoice ID o k body service based on context it is used
func (o *GetTradespersonTradespersonIDBillingCustomerStripeIDSubscriptionSubscriptionIDInvoiceInvoiceIDOKBodyService) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (o *GetTradespersonTradespersonIDBillingCustomerStripeIDSubscriptionSubscriptionIDInvoiceInvoiceIDOKBodyService) MarshalBinary() ([]byte, error) {
	if o == nil {
		return nil, nil
	}
	return swag.WriteJSON(o)
}

// UnmarshalBinary interface implementation
func (o *GetTradespersonTradespersonIDBillingCustomerStripeIDSubscriptionSubscriptionIDInvoiceInvoiceIDOKBodyService) UnmarshalBinary(b []byte) error {
	var res GetTradespersonTradespersonIDBillingCustomerStripeIDSubscriptionSubscriptionIDInvoiceInvoiceIDOKBodyService
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*o = res
	return nil
}
