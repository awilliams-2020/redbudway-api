// Code generated by go-swagger; DO NOT EDIT.

package operations

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"
)

// GetTradespersonTradespersonIDStatusOKCode is the HTTP code returned for type GetTradespersonTradespersonIDStatusOK
const GetTradespersonTradespersonIDStatusOKCode int = 200

/*GetTradespersonTradespersonIDStatusOK Status of tradesperson stripe account

swagger:response getTradespersonTradespersonIdStatusOK
*/
type GetTradespersonTradespersonIDStatusOK struct {

	/*
	  In: Body
	*/
	Payload *GetTradespersonTradespersonIDStatusOKBody `json:"body,omitempty"`
}

// NewGetTradespersonTradespersonIDStatusOK creates GetTradespersonTradespersonIDStatusOK with default headers values
func NewGetTradespersonTradespersonIDStatusOK() *GetTradespersonTradespersonIDStatusOK {

	return &GetTradespersonTradespersonIDStatusOK{}
}

// WithPayload adds the payload to the get tradesperson tradesperson Id status o k response
func (o *GetTradespersonTradespersonIDStatusOK) WithPayload(payload *GetTradespersonTradespersonIDStatusOKBody) *GetTradespersonTradespersonIDStatusOK {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get tradesperson tradesperson Id status o k response
func (o *GetTradespersonTradespersonIDStatusOK) SetPayload(payload *GetTradespersonTradespersonIDStatusOKBody) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetTradespersonTradespersonIDStatusOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(200)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}
