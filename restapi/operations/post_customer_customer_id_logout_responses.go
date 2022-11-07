// Code generated by go-swagger; DO NOT EDIT.

package operations

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"
)

// PostCustomerCustomerIDLogoutOKCode is the HTTP code returned for type PostCustomerCustomerIDLogoutOK
const PostCustomerCustomerIDLogoutOKCode int = 200

/*PostCustomerCustomerIDLogoutOK Log customer out

swagger:response postCustomerCustomerIdLogoutOK
*/
type PostCustomerCustomerIDLogoutOK struct {

	/*
	  In: Body
	*/
	Payload *PostCustomerCustomerIDLogoutOKBody `json:"body,omitempty"`
}

// NewPostCustomerCustomerIDLogoutOK creates PostCustomerCustomerIDLogoutOK with default headers values
func NewPostCustomerCustomerIDLogoutOK() *PostCustomerCustomerIDLogoutOK {

	return &PostCustomerCustomerIDLogoutOK{}
}

// WithPayload adds the payload to the post customer customer Id logout o k response
func (o *PostCustomerCustomerIDLogoutOK) WithPayload(payload *PostCustomerCustomerIDLogoutOKBody) *PostCustomerCustomerIDLogoutOK {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the post customer customer Id logout o k response
func (o *PostCustomerCustomerIDLogoutOK) SetPayload(payload *PostCustomerCustomerIDLogoutOKBody) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *PostCustomerCustomerIDLogoutOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(200)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}
