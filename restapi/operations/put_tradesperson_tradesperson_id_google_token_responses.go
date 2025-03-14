// Code generated by go-swagger; DO NOT EDIT.

package operations

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"
)

// PutTradespersonTradespersonIDGoogleTokenOKCode is the HTTP code returned for type PutTradespersonTradespersonIDGoogleTokenOK
const PutTradespersonTradespersonIDGoogleTokenOKCode int = 200

/*PutTradespersonTradespersonIDGoogleTokenOK Valid google token

swagger:response putTradespersonTradespersonIdGoogleTokenOK
*/
type PutTradespersonTradespersonIDGoogleTokenOK struct {

	/*
	  In: Body
	*/
	Payload *PutTradespersonTradespersonIDGoogleTokenOKBody `json:"body,omitempty"`
}

// NewPutTradespersonTradespersonIDGoogleTokenOK creates PutTradespersonTradespersonIDGoogleTokenOK with default headers values
func NewPutTradespersonTradespersonIDGoogleTokenOK() *PutTradespersonTradespersonIDGoogleTokenOK {

	return &PutTradespersonTradespersonIDGoogleTokenOK{}
}

// WithPayload adds the payload to the put tradesperson tradesperson Id google token o k response
func (o *PutTradespersonTradespersonIDGoogleTokenOK) WithPayload(payload *PutTradespersonTradespersonIDGoogleTokenOKBody) *PutTradespersonTradespersonIDGoogleTokenOK {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the put tradesperson tradesperson Id google token o k response
func (o *PutTradespersonTradespersonIDGoogleTokenOK) SetPayload(payload *PutTradespersonTradespersonIDGoogleTokenOKBody) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *PutTradespersonTradespersonIDGoogleTokenOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(200)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}
