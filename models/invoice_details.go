// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"
	"strconv"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
)

// InvoiceDetails invoice details
//
// swagger:model InvoiceDetails
type InvoiceDetails struct {

	// created
	Created int64 `json:"created,omitempty"`

	// customer
	Customer *Customer `json:"customer,omitempty"`

	// description
	Description string `json:"description"`

	// discount
	Discount *Discount `json:"discount,omitempty"`

	// due date
	DueDate int64 `json:"dueDate,omitempty"`

	// images
	Images []string `json:"images"`

	// interval
	Interval string `json:"interval,omitempty"`

	// number
	Number string `json:"number,omitempty"`

	// pdf
	Pdf string `json:"pdf,omitempty"`

	// products
	Products []*Product `json:"products"`

	// refunded
	Refunded int64 `json:"refunded,omitempty"`

	// service
	Service *InvoiceDetailsService `json:"service,omitempty"`

	// status
	Status string `json:"status,omitempty"`

	// time slot
	TimeSlot *TimeSlot `json:"timeSlot,omitempty"`

	// time zone
	TimeZone string `json:"timeZone,omitempty"`

	// total
	Total int64 `json:"total"`

	// url
	URL string `json:"url,omitempty"`
}

// Validate validates this invoice details
func (m *InvoiceDetails) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateCustomer(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateDiscount(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateProducts(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateService(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateTimeSlot(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *InvoiceDetails) validateCustomer(formats strfmt.Registry) error {
	if swag.IsZero(m.Customer) { // not required
		return nil
	}

	if m.Customer != nil {
		if err := m.Customer.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("customer")
			} else if ce, ok := err.(*errors.CompositeError); ok {
				return ce.ValidateName("customer")
			}
			return err
		}
	}

	return nil
}

func (m *InvoiceDetails) validateDiscount(formats strfmt.Registry) error {
	if swag.IsZero(m.Discount) { // not required
		return nil
	}

	if m.Discount != nil {
		if err := m.Discount.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("discount")
			} else if ce, ok := err.(*errors.CompositeError); ok {
				return ce.ValidateName("discount")
			}
			return err
		}
	}

	return nil
}

func (m *InvoiceDetails) validateProducts(formats strfmt.Registry) error {
	if swag.IsZero(m.Products) { // not required
		return nil
	}

	for i := 0; i < len(m.Products); i++ {
		if swag.IsZero(m.Products[i]) { // not required
			continue
		}

		if m.Products[i] != nil {
			if err := m.Products[i].Validate(formats); err != nil {
				if ve, ok := err.(*errors.Validation); ok {
					return ve.ValidateName("products" + "." + strconv.Itoa(i))
				} else if ce, ok := err.(*errors.CompositeError); ok {
					return ce.ValidateName("products" + "." + strconv.Itoa(i))
				}
				return err
			}
		}

	}

	return nil
}

func (m *InvoiceDetails) validateService(formats strfmt.Registry) error {
	if swag.IsZero(m.Service) { // not required
		return nil
	}

	if m.Service != nil {
		if err := m.Service.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("service")
			} else if ce, ok := err.(*errors.CompositeError); ok {
				return ce.ValidateName("service")
			}
			return err
		}
	}

	return nil
}

func (m *InvoiceDetails) validateTimeSlot(formats strfmt.Registry) error {
	if swag.IsZero(m.TimeSlot) { // not required
		return nil
	}

	if m.TimeSlot != nil {
		if err := m.TimeSlot.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("timeSlot")
			} else if ce, ok := err.(*errors.CompositeError); ok {
				return ce.ValidateName("timeSlot")
			}
			return err
		}
	}

	return nil
}

// ContextValidate validate this invoice details based on the context it is used
func (m *InvoiceDetails) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	var res []error

	if err := m.contextValidateCustomer(ctx, formats); err != nil {
		res = append(res, err)
	}

	if err := m.contextValidateDiscount(ctx, formats); err != nil {
		res = append(res, err)
	}

	if err := m.contextValidateProducts(ctx, formats); err != nil {
		res = append(res, err)
	}

	if err := m.contextValidateService(ctx, formats); err != nil {
		res = append(res, err)
	}

	if err := m.contextValidateTimeSlot(ctx, formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *InvoiceDetails) contextValidateCustomer(ctx context.Context, formats strfmt.Registry) error {

	if m.Customer != nil {
		if err := m.Customer.ContextValidate(ctx, formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("customer")
			} else if ce, ok := err.(*errors.CompositeError); ok {
				return ce.ValidateName("customer")
			}
			return err
		}
	}

	return nil
}

func (m *InvoiceDetails) contextValidateDiscount(ctx context.Context, formats strfmt.Registry) error {

	if m.Discount != nil {
		if err := m.Discount.ContextValidate(ctx, formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("discount")
			} else if ce, ok := err.(*errors.CompositeError); ok {
				return ce.ValidateName("discount")
			}
			return err
		}
	}

	return nil
}

func (m *InvoiceDetails) contextValidateProducts(ctx context.Context, formats strfmt.Registry) error {

	for i := 0; i < len(m.Products); i++ {

		if m.Products[i] != nil {
			if err := m.Products[i].ContextValidate(ctx, formats); err != nil {
				if ve, ok := err.(*errors.Validation); ok {
					return ve.ValidateName("products" + "." + strconv.Itoa(i))
				} else if ce, ok := err.(*errors.CompositeError); ok {
					return ce.ValidateName("products" + "." + strconv.Itoa(i))
				}
				return err
			}
		}

	}

	return nil
}

func (m *InvoiceDetails) contextValidateService(ctx context.Context, formats strfmt.Registry) error {

	if m.Service != nil {
		if err := m.Service.ContextValidate(ctx, formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("service")
			} else if ce, ok := err.(*errors.CompositeError); ok {
				return ce.ValidateName("service")
			}
			return err
		}
	}

	return nil
}

func (m *InvoiceDetails) contextValidateTimeSlot(ctx context.Context, formats strfmt.Registry) error {

	if m.TimeSlot != nil {
		if err := m.TimeSlot.ContextValidate(ctx, formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("timeSlot")
			} else if ce, ok := err.(*errors.CompositeError); ok {
				return ce.ValidateName("timeSlot")
			}
			return err
		}
	}

	return nil
}

// MarshalBinary interface implementation
func (m *InvoiceDetails) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *InvoiceDetails) UnmarshalBinary(b []byte) error {
	var res InvoiceDetails
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}

// InvoiceDetailsService invoice details service
//
// swagger:model InvoiceDetailsService
type InvoiceDetailsService struct {

	// description
	Description string `json:"description,omitempty"`

	// title
	Title string `json:"title,omitempty"`
}

// Validate validates this invoice details service
func (m *InvoiceDetailsService) Validate(formats strfmt.Registry) error {
	return nil
}

// ContextValidate validates this invoice details service based on context it is used
func (m *InvoiceDetailsService) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (m *InvoiceDetailsService) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *InvoiceDetailsService) UnmarshalBinary(b []byte) error {
	var res InvoiceDetailsService
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
