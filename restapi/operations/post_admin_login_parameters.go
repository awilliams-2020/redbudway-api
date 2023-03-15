// Code generated by go-swagger; DO NOT EDIT.

package operations

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"
	"net/http"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/validate"
)

// NewPostAdminLoginParams creates a new PostAdminLoginParams object
//
// There are no default values defined in the spec.
func NewPostAdminLoginParams() PostAdminLoginParams {

	return PostAdminLoginParams{}
}

// PostAdminLoginParams contains all the bound params for the post admin login operation
// typically these are obtained from a http.Request
//
// swagger:parameters PostAdminLogin
type PostAdminLoginParams struct {

	// HTTP Request Object
	HTTPRequest *http.Request `json:"-"`

	/*Log in to admin account.
	  In: body
	*/
	Admin PostAdminLoginBody
}

// BindRequest both binds and validates a request, it assumes that complex things implement a Validatable(strfmt.Registry) error interface
// for simple values it will use straight method calls.
//
// To ensure default values, the struct must have been initialized with NewPostAdminLoginParams() beforehand.
func (o *PostAdminLoginParams) BindRequest(r *http.Request, route *middleware.MatchedRoute) error {
	var res []error

	o.HTTPRequest = r

	if runtime.HasBody(r) {
		defer r.Body.Close()
		var body PostAdminLoginBody
		if err := route.Consumer.Consume(r.Body, &body); err != nil {
			res = append(res, errors.NewParseError("admin", "body", "", err))
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
				o.Admin = body
			}
		}
	}
	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}
