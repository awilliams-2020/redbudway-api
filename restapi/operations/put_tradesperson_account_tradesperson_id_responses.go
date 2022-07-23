// Code generated by go-swagger; DO NOT EDIT.

package operations

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"
)

// PutTradespersonAccountTradespersonIDOKCode is the HTTP code returned for type PutTradespersonAccountTradespersonIDOK
const PutTradespersonAccountTradespersonIDOKCode int = 200

/*PutTradespersonAccountTradespersonIDOK Updates tradesperson profile

swagger:response putTradespersonAccountTradespersonIdOK
*/
type PutTradespersonAccountTradespersonIDOK struct {

	/*
	  In: Body
	*/
	Payload *PutTradespersonAccountTradespersonIDOKBody `json:"body,omitempty"`
}

// NewPutTradespersonAccountTradespersonIDOK creates PutTradespersonAccountTradespersonIDOK with default headers values
func NewPutTradespersonAccountTradespersonIDOK() *PutTradespersonAccountTradespersonIDOK {

	return &PutTradespersonAccountTradespersonIDOK{}
}

// WithPayload adds the payload to the put tradesperson account tradesperson Id o k response
func (o *PutTradespersonAccountTradespersonIDOK) WithPayload(payload *PutTradespersonAccountTradespersonIDOKBody) *PutTradespersonAccountTradespersonIDOK {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the put tradesperson account tradesperson Id o k response
func (o *PutTradespersonAccountTradespersonIDOK) SetPayload(payload *PutTradespersonAccountTradespersonIDOKBody) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *PutTradespersonAccountTradespersonIDOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(200)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}
