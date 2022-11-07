// Code generated by go-swagger; DO NOT EDIT.

package operations

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"
)

// GetCustomerCustomerIDPaymentDefaultOKCode is the HTTP code returned for type GetCustomerCustomerIDPaymentDefaultOK
const GetCustomerCustomerIDPaymentDefaultOKCode int = 200

/*GetCustomerCustomerIDPaymentDefaultOK Customer default payment

swagger:response getCustomerCustomerIdPaymentDefaultOK
*/
type GetCustomerCustomerIDPaymentDefaultOK struct {

	/*
	  In: Body
	*/
	Payload *GetCustomerCustomerIDPaymentDefaultOKBody `json:"body,omitempty"`
}

// NewGetCustomerCustomerIDPaymentDefaultOK creates GetCustomerCustomerIDPaymentDefaultOK with default headers values
func NewGetCustomerCustomerIDPaymentDefaultOK() *GetCustomerCustomerIDPaymentDefaultOK {

	return &GetCustomerCustomerIDPaymentDefaultOK{}
}

// WithPayload adds the payload to the get customer customer Id payment default o k response
func (o *GetCustomerCustomerIDPaymentDefaultOK) WithPayload(payload *GetCustomerCustomerIDPaymentDefaultOKBody) *GetCustomerCustomerIDPaymentDefaultOK {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get customer customer Id payment default o k response
func (o *GetCustomerCustomerIDPaymentDefaultOK) SetPayload(payload *GetCustomerCustomerIDPaymentDefaultOKBody) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetCustomerCustomerIDPaymentDefaultOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(200)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}
