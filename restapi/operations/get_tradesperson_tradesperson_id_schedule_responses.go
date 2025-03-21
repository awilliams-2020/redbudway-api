// Code generated by go-swagger; DO NOT EDIT.

package operations

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"
)

// GetTradespersonTradespersonIDScheduleOKCode is the HTTP code returned for type GetTradespersonTradespersonIDScheduleOK
const GetTradespersonTradespersonIDScheduleOKCode int = 200

/*GetTradespersonTradespersonIDScheduleOK Tradesperson schedule

swagger:response getTradespersonTradespersonIdScheduleOK
*/
type GetTradespersonTradespersonIDScheduleOK struct {

	/*
	  In: Body
	*/
	Payload *GetTradespersonTradespersonIDScheduleOKBody `json:"body,omitempty"`
}

// NewGetTradespersonTradespersonIDScheduleOK creates GetTradespersonTradespersonIDScheduleOK with default headers values
func NewGetTradespersonTradespersonIDScheduleOK() *GetTradespersonTradespersonIDScheduleOK {

	return &GetTradespersonTradespersonIDScheduleOK{}
}

// WithPayload adds the payload to the get tradesperson tradesperson Id schedule o k response
func (o *GetTradespersonTradespersonIDScheduleOK) WithPayload(payload *GetTradespersonTradespersonIDScheduleOKBody) *GetTradespersonTradespersonIDScheduleOK {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get tradesperson tradesperson Id schedule o k response
func (o *GetTradespersonTradespersonIDScheduleOK) SetPayload(payload *GetTradespersonTradespersonIDScheduleOKBody) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetTradespersonTradespersonIDScheduleOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(200)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}
