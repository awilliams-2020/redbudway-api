// Code generated by go-swagger; DO NOT EDIT.

package operations

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"
)

// PostCustomerCustomerIDVerifyOKCode is the HTTP code returned for type PostCustomerCustomerIDVerifyOK
const PostCustomerCustomerIDVerifyOKCode int = 200

/*PostCustomerCustomerIDVerifyOK Verified customer account

swagger:response postCustomerCustomerIdVerifyOK
*/
type PostCustomerCustomerIDVerifyOK struct {

	/*
	  In: Body
	*/
	Payload *PostCustomerCustomerIDVerifyOKBody `json:"body,omitempty"`
}

// NewPostCustomerCustomerIDVerifyOK creates PostCustomerCustomerIDVerifyOK with default headers values
func NewPostCustomerCustomerIDVerifyOK() *PostCustomerCustomerIDVerifyOK {

	return &PostCustomerCustomerIDVerifyOK{}
}

// WithPayload adds the payload to the post customer customer Id verify o k response
func (o *PostCustomerCustomerIDVerifyOK) WithPayload(payload *PostCustomerCustomerIDVerifyOKBody) *PostCustomerCustomerIDVerifyOK {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the post customer customer Id verify o k response
func (o *PostCustomerCustomerIDVerifyOK) SetPayload(payload *PostCustomerCustomerIDVerifyOKBody) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *PostCustomerCustomerIDVerifyOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(200)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}
