// Code generated by go-swagger; DO NOT EDIT.

package operations

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"
)

// PutTradespersonTradespersonIDPasswordOKCode is the HTTP code returned for type PutTradespersonTradespersonIDPasswordOK
const PutTradespersonTradespersonIDPasswordOKCode int = 200

/*PutTradespersonTradespersonIDPasswordOK Tradeperson password updated

swagger:response putTradespersonTradespersonIdPasswordOK
*/
type PutTradespersonTradespersonIDPasswordOK struct {

	/*
	  In: Body
	*/
	Payload *PutTradespersonTradespersonIDPasswordOKBody `json:"body,omitempty"`
}

// NewPutTradespersonTradespersonIDPasswordOK creates PutTradespersonTradespersonIDPasswordOK with default headers values
func NewPutTradespersonTradespersonIDPasswordOK() *PutTradespersonTradespersonIDPasswordOK {

	return &PutTradespersonTradespersonIDPasswordOK{}
}

// WithPayload adds the payload to the put tradesperson tradesperson Id password o k response
func (o *PutTradespersonTradespersonIDPasswordOK) WithPayload(payload *PutTradespersonTradespersonIDPasswordOKBody) *PutTradespersonTradespersonIDPasswordOK {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the put tradesperson tradesperson Id password o k response
func (o *PutTradespersonTradespersonIDPasswordOK) SetPayload(payload *PutTradespersonTradespersonIDPasswordOKBody) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *PutTradespersonTradespersonIDPasswordOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(200)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}
