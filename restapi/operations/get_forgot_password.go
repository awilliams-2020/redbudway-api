// Code generated by go-swagger; DO NOT EDIT.

package operations

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"net/http"

	"github.com/go-openapi/runtime/middleware"
)

// GetForgotPasswordHandlerFunc turns a function with the right signature into a get forgot password handler
type GetForgotPasswordHandlerFunc func(GetForgotPasswordParams) middleware.Responder

// Handle executing the request and returning a response
func (fn GetForgotPasswordHandlerFunc) Handle(params GetForgotPasswordParams) middleware.Responder {
	return fn(params)
}

// GetForgotPasswordHandler interface for that can handle valid get forgot password params
type GetForgotPasswordHandler interface {
	Handle(GetForgotPasswordParams) middleware.Responder
}

// NewGetForgotPassword creates a new http.Handler for the get forgot password operation
func NewGetForgotPassword(ctx *middleware.Context, handler GetForgotPasswordHandler) *GetForgotPassword {
	return &GetForgotPassword{Context: ctx, Handler: handler}
}

/* GetForgotPassword swagger:route GET /forgot-password getForgotPassword

GetForgotPassword get forgot password API

*/
type GetForgotPassword struct {
	Context *middleware.Context
	Handler GetForgotPasswordHandler
}

func (o *GetForgotPassword) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, rCtx, _ := o.Context.RouteInfo(r)
	if rCtx != nil {
		*r = *rCtx
	}
	var Params = NewGetForgotPasswordParams()
	if err := o.Context.BindValidRequest(r, route, &Params); err != nil { // bind params
		o.Context.Respond(rw, r, route.Produces, route, err)
		return
	}

	res := o.Handler.Handle(Params) // actually handle the request
	o.Context.Respond(rw, r, route.Produces, route, res)

}
