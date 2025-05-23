// Code generated by go-swagger; DO NOT EDIT.

package operations

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"
)

// PostTradespersonTradespersonIDEmailOKCode is the HTTP code returned for type PostTradespersonTradespersonIDEmailOK
const PostTradespersonTradespersonIDEmailOKCode int = 200

/*PostTradespersonTradespersonIDEmailOK Email sent to tradesperson

swagger:response postTradespersonTradespersonIdEmailOK
*/
type PostTradespersonTradespersonIDEmailOK struct {

	/*
	  In: Body
	*/
	Payload *PostTradespersonTradespersonIDEmailOKBody `json:"body,omitempty"`
}

// NewPostTradespersonTradespersonIDEmailOK creates PostTradespersonTradespersonIDEmailOK with default headers values
func NewPostTradespersonTradespersonIDEmailOK() *PostTradespersonTradespersonIDEmailOK {

	return &PostTradespersonTradespersonIDEmailOK{}
}

// WithPayload adds the payload to the post tradesperson tradesperson Id email o k response
func (o *PostTradespersonTradespersonIDEmailOK) WithPayload(payload *PostTradespersonTradespersonIDEmailOKBody) *PostTradespersonTradespersonIDEmailOK {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the post tradesperson tradesperson Id email o k response
func (o *PostTradespersonTradespersonIDEmailOK) SetPayload(payload *PostTradespersonTradespersonIDEmailOKBody) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *PostTradespersonTradespersonIDEmailOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(200)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}
