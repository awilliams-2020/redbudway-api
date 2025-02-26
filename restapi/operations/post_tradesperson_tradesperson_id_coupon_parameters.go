// Code generated by go-swagger; DO NOT EDIT.

package operations

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"
	"io"
	"net/http"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/validate"

	"redbudway-api/models"
)

// NewPostTradespersonTradespersonIDCouponParams creates a new PostTradespersonTradespersonIDCouponParams object
//
// There are no default values defined in the spec.
func NewPostTradespersonTradespersonIDCouponParams() PostTradespersonTradespersonIDCouponParams {

	return PostTradespersonTradespersonIDCouponParams{}
}

// PostTradespersonTradespersonIDCouponParams contains all the bound params for the post tradesperson tradesperson ID coupon operation
// typically these are obtained from a http.Request
//
// swagger:parameters PostTradespersonTradespersonIDCoupon
type PostTradespersonTradespersonIDCouponParams struct {

	// HTTP Request Object
	HTTPRequest *http.Request `json:"-"`

	/*Coupon to create
	  Required: true
	  In: body
	*/
	Coupon *models.Coupon
	/*Tradesperson ID
	  Required: true
	  In: path
	*/
	TradespersonID string
}

// BindRequest both binds and validates a request, it assumes that complex things implement a Validatable(strfmt.Registry) error interface
// for simple values it will use straight method calls.
//
// To ensure default values, the struct must have been initialized with NewPostTradespersonTradespersonIDCouponParams() beforehand.
func (o *PostTradespersonTradespersonIDCouponParams) BindRequest(r *http.Request, route *middleware.MatchedRoute) error {
	var res []error

	o.HTTPRequest = r

	if runtime.HasBody(r) {
		defer r.Body.Close()
		var body models.Coupon
		if err := route.Consumer.Consume(r.Body, &body); err != nil {
			if err == io.EOF {
				res = append(res, errors.Required("coupon", "body", ""))
			} else {
				res = append(res, errors.NewParseError("coupon", "body", "", err))
			}
		} else {
			// validate body object
			if err := body.Validate(route.Formats); err != nil {
				res = append(res, err)
			}

			ctx := validate.WithOperationRequest(context.Background())
			if err := body.ContextValidate(ctx, route.Formats); err != nil {
				res = append(res, err)
			}

			if len(res) == 0 {
				o.Coupon = &body
			}
		}
	} else {
		res = append(res, errors.Required("coupon", "body", ""))
	}

	rTradespersonID, rhkTradespersonID, _ := route.Params.GetOK("tradespersonId")
	if err := o.bindTradespersonID(rTradespersonID, rhkTradespersonID, route.Formats); err != nil {
		res = append(res, err)
	}
	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

// bindTradespersonID binds and validates parameter TradespersonID from path.
func (o *PostTradespersonTradespersonIDCouponParams) bindTradespersonID(rawData []string, hasKey bool, formats strfmt.Registry) error {
	var raw string
	if len(rawData) > 0 {
		raw = rawData[len(rawData)-1]
	}

	// Required: true
	// Parameter is provided by construction from the route
	o.TradespersonID = raw

	return nil
}
