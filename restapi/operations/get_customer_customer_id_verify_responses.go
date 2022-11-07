// Code generated by go-swagger; DO NOT EDIT.

package operations

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"
)

// GetCustomerCustomerIDVerifyOKCode is the HTTP code returned for type GetCustomerCustomerIDVerifyOK
const GetCustomerCustomerIDVerifyOKCode int = 200

/*GetCustomerCustomerIDVerifyOK Verified customer account

swagger:response getCustomerCustomerIdVerifyOK
*/
type GetCustomerCustomerIDVerifyOK struct {

	/*
	  In: Body
	*/
	Payload *GetCustomerCustomerIDVerifyOKBody `json:"body,omitempty"`
}

// NewGetCustomerCustomerIDVerifyOK creates GetCustomerCustomerIDVerifyOK with default headers values
func NewGetCustomerCustomerIDVerifyOK() *GetCustomerCustomerIDVerifyOK {

	return &GetCustomerCustomerIDVerifyOK{}
}

// WithPayload adds the payload to the get customer customer Id verify o k response
func (o *GetCustomerCustomerIDVerifyOK) WithPayload(payload *GetCustomerCustomerIDVerifyOKBody) *GetCustomerCustomerIDVerifyOK {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get customer customer Id verify o k response
func (o *GetCustomerCustomerIDVerifyOK) SetPayload(payload *GetCustomerCustomerIDVerifyOKBody) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetCustomerCustomerIDVerifyOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(200)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}
