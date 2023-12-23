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

// GetTradespersonTradespersonIDBillingInvoiceInvoiceIDHandlerFunc turns a function with the right signature into a get tradesperson tradesperson ID billing invoice invoice ID handler
type GetTradespersonTradespersonIDBillingInvoiceInvoiceIDHandlerFunc func(GetTradespersonTradespersonIDBillingInvoiceInvoiceIDParams, interface{}) middleware.Responder

// Handle executing the request and returning a response
func (fn GetTradespersonTradespersonIDBillingInvoiceInvoiceIDHandlerFunc) Handle(params GetTradespersonTradespersonIDBillingInvoiceInvoiceIDParams, principal interface{}) middleware.Responder {
	return fn(params, principal)
}

// GetTradespersonTradespersonIDBillingInvoiceInvoiceIDHandler interface for that can handle valid get tradesperson tradesperson ID billing invoice invoice ID params
type GetTradespersonTradespersonIDBillingInvoiceInvoiceIDHandler interface {
	Handle(GetTradespersonTradespersonIDBillingInvoiceInvoiceIDParams, interface{}) middleware.Responder
}

// NewGetTradespersonTradespersonIDBillingInvoiceInvoiceID creates a new http.Handler for the get tradesperson tradesperson ID billing invoice invoice ID operation
func NewGetTradespersonTradespersonIDBillingInvoiceInvoiceID(ctx *middleware.Context, handler GetTradespersonTradespersonIDBillingInvoiceInvoiceIDHandler) *GetTradespersonTradespersonIDBillingInvoiceInvoiceID {
	return &GetTradespersonTradespersonIDBillingInvoiceInvoiceID{Context: ctx, Handler: handler}
}

/* GetTradespersonTradespersonIDBillingInvoiceInvoiceID swagger:route GET /tradesperson/{tradespersonId}/billing/invoice/{invoiceId} getTradespersonTradespersonIdBillingInvoiceInvoiceId

GetTradespersonTradespersonIDBillingInvoiceInvoiceID get tradesperson tradesperson ID billing invoice invoice ID API

*/
type GetTradespersonTradespersonIDBillingInvoiceInvoiceID struct {
	Context *middleware.Context
	Handler GetTradespersonTradespersonIDBillingInvoiceInvoiceIDHandler
}

func (o *GetTradespersonTradespersonIDBillingInvoiceInvoiceID) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, rCtx, _ := o.Context.RouteInfo(r)
	if rCtx != nil {
		*r = *rCtx
	}
	var Params = NewGetTradespersonTradespersonIDBillingInvoiceInvoiceIDParams()
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

// GetTradespersonTradespersonIDBillingInvoiceInvoiceIDOKBody get tradesperson tradesperson ID billing invoice invoice ID o k body
//
// swagger:model GetTradespersonTradespersonIDBillingInvoiceInvoiceIDOKBody
type GetTradespersonTradespersonIDBillingInvoiceInvoiceIDOKBody struct {

	// created
	Created int64 `json:"created,omitempty"`

	// customer
	Customer *models.Customer `json:"customer,omitempty"`

	// description
	Description string `json:"description"`

	// number
	Number string `json:"number,omitempty"`

	// pdf
	Pdf string `json:"pdf,omitempty"`

	// refunded
	Refunded int64 `json:"refunded,omitempty"`

	// service
	Service *GetTradespersonTradespersonIDBillingInvoiceInvoiceIDOKBodyService `json:"service,omitempty"`

	// status
	Status string `json:"status,omitempty"`

	// time slot
	TimeSlot *GetTradespersonTradespersonIDBillingInvoiceInvoiceIDOKBodyTimeSlot `json:"timeSlot,omitempty"`

	// time zone
	TimeZone string `json:"timeZone,omitempty"`

	// total
	Total int64 `json:"total"`

	// url
	URL string `json:"url,omitempty"`
}

// Validate validates this get tradesperson tradesperson ID billing invoice invoice ID o k body
func (o *GetTradespersonTradespersonIDBillingInvoiceInvoiceIDOKBody) Validate(formats strfmt.Registry) error {
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

func (o *GetTradespersonTradespersonIDBillingInvoiceInvoiceIDOKBody) validateCustomer(formats strfmt.Registry) error {
	if swag.IsZero(o.Customer) { // not required
		return nil
	}

	if o.Customer != nil {
		if err := o.Customer.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("getTradespersonTradespersonIdBillingInvoiceInvoiceIdOK" + "." + "customer")
			} else if ce, ok := err.(*errors.CompositeError); ok {
				return ce.ValidateName("getTradespersonTradespersonIdBillingInvoiceInvoiceIdOK" + "." + "customer")
			}
			return err
		}
	}

	return nil
}

func (o *GetTradespersonTradespersonIDBillingInvoiceInvoiceIDOKBody) validateService(formats strfmt.Registry) error {
	if swag.IsZero(o.Service) { // not required
		return nil
	}

	if o.Service != nil {
		if err := o.Service.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("getTradespersonTradespersonIdBillingInvoiceInvoiceIdOK" + "." + "service")
			} else if ce, ok := err.(*errors.CompositeError); ok {
				return ce.ValidateName("getTradespersonTradespersonIdBillingInvoiceInvoiceIdOK" + "." + "service")
			}
			return err
		}
	}

	return nil
}

func (o *GetTradespersonTradespersonIDBillingInvoiceInvoiceIDOKBody) validateTimeSlot(formats strfmt.Registry) error {
	if swag.IsZero(o.TimeSlot) { // not required
		return nil
	}

	if o.TimeSlot != nil {
		if err := o.TimeSlot.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("getTradespersonTradespersonIdBillingInvoiceInvoiceIdOK" + "." + "timeSlot")
			} else if ce, ok := err.(*errors.CompositeError); ok {
				return ce.ValidateName("getTradespersonTradespersonIdBillingInvoiceInvoiceIdOK" + "." + "timeSlot")
			}
			return err
		}
	}

	return nil
}

// ContextValidate validate this get tradesperson tradesperson ID billing invoice invoice ID o k body based on the context it is used
func (o *GetTradespersonTradespersonIDBillingInvoiceInvoiceIDOKBody) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
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

func (o *GetTradespersonTradespersonIDBillingInvoiceInvoiceIDOKBody) contextValidateCustomer(ctx context.Context, formats strfmt.Registry) error {

	if o.Customer != nil {
		if err := o.Customer.ContextValidate(ctx, formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("getTradespersonTradespersonIdBillingInvoiceInvoiceIdOK" + "." + "customer")
			} else if ce, ok := err.(*errors.CompositeError); ok {
				return ce.ValidateName("getTradespersonTradespersonIdBillingInvoiceInvoiceIdOK" + "." + "customer")
			}
			return err
		}
	}

	return nil
}

func (o *GetTradespersonTradespersonIDBillingInvoiceInvoiceIDOKBody) contextValidateService(ctx context.Context, formats strfmt.Registry) error {

	if o.Service != nil {
		if err := o.Service.ContextValidate(ctx, formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("getTradespersonTradespersonIdBillingInvoiceInvoiceIdOK" + "." + "service")
			} else if ce, ok := err.(*errors.CompositeError); ok {
				return ce.ValidateName("getTradespersonTradespersonIdBillingInvoiceInvoiceIdOK" + "." + "service")
			}
			return err
		}
	}

	return nil
}

func (o *GetTradespersonTradespersonIDBillingInvoiceInvoiceIDOKBody) contextValidateTimeSlot(ctx context.Context, formats strfmt.Registry) error {

	if o.TimeSlot != nil {
		if err := o.TimeSlot.ContextValidate(ctx, formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("getTradespersonTradespersonIdBillingInvoiceInvoiceIdOK" + "." + "timeSlot")
			} else if ce, ok := err.(*errors.CompositeError); ok {
				return ce.ValidateName("getTradespersonTradespersonIdBillingInvoiceInvoiceIdOK" + "." + "timeSlot")
			}
			return err
		}
	}

	return nil
}

// MarshalBinary interface implementation
func (o *GetTradespersonTradespersonIDBillingInvoiceInvoiceIDOKBody) MarshalBinary() ([]byte, error) {
	if o == nil {
		return nil, nil
	}
	return swag.WriteJSON(o)
}

// UnmarshalBinary interface implementation
func (o *GetTradespersonTradespersonIDBillingInvoiceInvoiceIDOKBody) UnmarshalBinary(b []byte) error {
	var res GetTradespersonTradespersonIDBillingInvoiceInvoiceIDOKBody
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*o = res
	return nil
}

// GetTradespersonTradespersonIDBillingInvoiceInvoiceIDOKBodyService get tradesperson tradesperson ID billing invoice invoice ID o k body service
//
// swagger:model GetTradespersonTradespersonIDBillingInvoiceInvoiceIDOKBodyService
type GetTradespersonTradespersonIDBillingInvoiceInvoiceIDOKBodyService struct {

	// description
	Description string `json:"description,omitempty"`

	// title
	Title string `json:"title,omitempty"`
}

// Validate validates this get tradesperson tradesperson ID billing invoice invoice ID o k body service
func (o *GetTradespersonTradespersonIDBillingInvoiceInvoiceIDOKBodyService) Validate(formats strfmt.Registry) error {
	return nil
}

// ContextValidate validates this get tradesperson tradesperson ID billing invoice invoice ID o k body service based on context it is used
func (o *GetTradespersonTradespersonIDBillingInvoiceInvoiceIDOKBodyService) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (o *GetTradespersonTradespersonIDBillingInvoiceInvoiceIDOKBodyService) MarshalBinary() ([]byte, error) {
	if o == nil {
		return nil, nil
	}
	return swag.WriteJSON(o)
}

// UnmarshalBinary interface implementation
func (o *GetTradespersonTradespersonIDBillingInvoiceInvoiceIDOKBodyService) UnmarshalBinary(b []byte) error {
	var res GetTradespersonTradespersonIDBillingInvoiceInvoiceIDOKBodyService
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*o = res
	return nil
}

// GetTradespersonTradespersonIDBillingInvoiceInvoiceIDOKBodyTimeSlot get tradesperson tradesperson ID billing invoice invoice ID o k body time slot
//
// swagger:model GetTradespersonTradespersonIDBillingInvoiceInvoiceIDOKBodyTimeSlot
type GetTradespersonTradespersonIDBillingInvoiceInvoiceIDOKBodyTimeSlot struct {

	// end time
	EndTime string `json:"endTime,omitempty"`

	// start time
	StartTime string `json:"startTime,omitempty"`
}

// Validate validates this get tradesperson tradesperson ID billing invoice invoice ID o k body time slot
func (o *GetTradespersonTradespersonIDBillingInvoiceInvoiceIDOKBodyTimeSlot) Validate(formats strfmt.Registry) error {
	return nil
}

// ContextValidate validates this get tradesperson tradesperson ID billing invoice invoice ID o k body time slot based on context it is used
func (o *GetTradespersonTradespersonIDBillingInvoiceInvoiceIDOKBodyTimeSlot) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (o *GetTradespersonTradespersonIDBillingInvoiceInvoiceIDOKBodyTimeSlot) MarshalBinary() ([]byte, error) {
	if o == nil {
		return nil, nil
	}
	return swag.WriteJSON(o)
}

// UnmarshalBinary interface implementation
func (o *GetTradespersonTradespersonIDBillingInvoiceInvoiceIDOKBodyTimeSlot) UnmarshalBinary(b []byte) error {
	var res GetTradespersonTradespersonIDBillingInvoiceInvoiceIDOKBodyTimeSlot
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*o = res
	return nil
}
