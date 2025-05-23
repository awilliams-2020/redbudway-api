// Code generated by go-swagger; DO NOT EDIT.

package operations

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"
)

// PostResetPasswordOKCode is the HTTP code returned for type PostResetPasswordOK
const PostResetPasswordOKCode int = 200

/*PostResetPasswordOK Reset password

swagger:response postResetPasswordOK
*/
type PostResetPasswordOK struct {

	/*
	  In: Body
	*/
	Payload *PostResetPasswordOKBody `json:"body,omitempty"`
}

// NewPostResetPasswordOK creates PostResetPasswordOK with default headers values
func NewPostResetPasswordOK() *PostResetPasswordOK {

	return &PostResetPasswordOK{}
}

// WithPayload adds the payload to the post reset password o k response
func (o *PostResetPasswordOK) WithPayload(payload *PostResetPasswordOKBody) *PostResetPasswordOK {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the post reset password o k response
func (o *PostResetPasswordOK) SetPayload(payload *PostResetPasswordOKBody) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *PostResetPasswordOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(200)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}
