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

// GoogleTimeSlots google time slots
//
// swagger:model GoogleTimeSlots
type GoogleTimeSlots []*GoogleTimeSlotsItems0

// Validate validates this google time slots
func (m GoogleTimeSlots) Validate(formats strfmt.Registry) error {
	var res []error

	for i := 0; i < len(m); i++ {
		if swag.IsZero(m[i]) { // not required
			continue
		}

		if m[i] != nil {
			if err := m[i].Validate(formats); err != nil {
				if ve, ok := err.(*errors.Validation); ok {
					return ve.ValidateName(strconv.Itoa(i))
				} else if ce, ok := err.(*errors.CompositeError); ok {
					return ce.ValidateName(strconv.Itoa(i))
				}
				return err
			}
		}

	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

// ContextValidate validate this google time slots based on the context it is used
func (m GoogleTimeSlots) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	var res []error

	for i := 0; i < len(m); i++ {

		if m[i] != nil {
			if err := m[i].ContextValidate(ctx, formats); err != nil {
				if ve, ok := err.(*errors.Validation); ok {
					return ve.ValidateName(strconv.Itoa(i))
				} else if ce, ok := err.(*errors.CompositeError); ok {
					return ce.ValidateName(strconv.Itoa(i))
				}
				return err
			}
		}

	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

// GoogleTimeSlotsItems0 google time slots items0
//
// swagger:model GoogleTimeSlotsItems0
type GoogleTimeSlotsItems0 struct {

	// end time
	EndTime string `json:"endTime,omitempty"`

	// recurrence
	Recurrence string `json:"recurrence,omitempty"`

	// start time
	StartTime string `json:"startTime,omitempty"`

	// time zone
	TimeZone string `json:"timeZone,omitempty"`

	// title
	Title string `json:"title,omitempty"`
}

// Validate validates this google time slots items0
func (m *GoogleTimeSlotsItems0) Validate(formats strfmt.Registry) error {
	return nil
}

// ContextValidate validates this google time slots items0 based on context it is used
func (m *GoogleTimeSlotsItems0) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (m *GoogleTimeSlotsItems0) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *GoogleTimeSlotsItems0) UnmarshalBinary(b []byte) error {
	var res GoogleTimeSlotsItems0
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
