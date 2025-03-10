// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"

	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
)

// Discount discount
//
// swagger:model Discount
type Discount struct {

	// amount
	Amount float64 `json:"amount,omitempty"`

	// code
	Code string `json:"code,omitempty"`

	// coupon Id
	CouponID string `json:"couponId,omitempty"`

	// duration
	Duration string `json:"duration,omitempty"`

	// months
	Months int64 `json:"months,omitempty"`

	// percent
	Percent float64 `json:"percent,omitempty"`

	// type
	Type string `json:"type,omitempty"`

	// valid
	Valid bool `json:"valid"`
}

// Validate validates this discount
func (m *Discount) Validate(formats strfmt.Registry) error {
	return nil
}

// ContextValidate validates this discount based on context it is used
func (m *Discount) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (m *Discount) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *Discount) UnmarshalBinary(b []byte) error {
	var res Discount
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
