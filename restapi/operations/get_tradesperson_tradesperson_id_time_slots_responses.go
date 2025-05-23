// Code generated by go-swagger; DO NOT EDIT.

package operations

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"
)

// GetTradespersonTradespersonIDTimeSlotsOKCode is the HTTP code returned for type GetTradespersonTradespersonIDTimeSlotsOK
const GetTradespersonTradespersonIDTimeSlotsOKCode int = 200

/*GetTradespersonTradespersonIDTimeSlotsOK Service time slots

swagger:response getTradespersonTradespersonIdTimeSlotsOK
*/
type GetTradespersonTradespersonIDTimeSlotsOK struct {

	/*
	  In: Body
	*/
	Payload *GetTradespersonTradespersonIDTimeSlotsOKBody `json:"body,omitempty"`
}

// NewGetTradespersonTradespersonIDTimeSlotsOK creates GetTradespersonTradespersonIDTimeSlotsOK with default headers values
func NewGetTradespersonTradespersonIDTimeSlotsOK() *GetTradespersonTradespersonIDTimeSlotsOK {

	return &GetTradespersonTradespersonIDTimeSlotsOK{}
}

// WithPayload adds the payload to the get tradesperson tradesperson Id time slots o k response
func (o *GetTradespersonTradespersonIDTimeSlotsOK) WithPayload(payload *GetTradespersonTradespersonIDTimeSlotsOKBody) *GetTradespersonTradespersonIDTimeSlotsOK {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get tradesperson tradesperson Id time slots o k response
func (o *GetTradespersonTradespersonIDTimeSlotsOK) SetPayload(payload *GetTradespersonTradespersonIDTimeSlotsOKBody) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetTradespersonTradespersonIDTimeSlotsOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(200)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}
