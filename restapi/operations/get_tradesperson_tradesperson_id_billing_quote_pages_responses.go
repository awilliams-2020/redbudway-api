// Code generated by go-swagger; DO NOT EDIT.

package operations

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"
)

// GetTradespersonTradespersonIDBillingQuotePagesOKCode is the HTTP code returned for type GetTradespersonTradespersonIDBillingQuotePagesOK
const GetTradespersonTradespersonIDBillingQuotePagesOKCode int = 200

/*GetTradespersonTradespersonIDBillingQuotePagesOK Number of quote pages

swagger:response getTradespersonTradespersonIdBillingQuotePagesOK
*/
type GetTradespersonTradespersonIDBillingQuotePagesOK struct {

	/*
	  In: Body
	*/
	Payload int64 `json:"body,omitempty"`
}

// NewGetTradespersonTradespersonIDBillingQuotePagesOK creates GetTradespersonTradespersonIDBillingQuotePagesOK with default headers values
func NewGetTradespersonTradespersonIDBillingQuotePagesOK() *GetTradespersonTradespersonIDBillingQuotePagesOK {

	return &GetTradespersonTradespersonIDBillingQuotePagesOK{}
}

// WithPayload adds the payload to the get tradesperson tradesperson Id billing quote pages o k response
func (o *GetTradespersonTradespersonIDBillingQuotePagesOK) WithPayload(payload int64) *GetTradespersonTradespersonIDBillingQuotePagesOK {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get tradesperson tradesperson Id billing quote pages o k response
func (o *GetTradespersonTradespersonIDBillingQuotePagesOK) SetPayload(payload int64) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetTradespersonTradespersonIDBillingQuotePagesOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(200)
	payload := o.Payload
	if err := producer.Produce(rw, payload); err != nil {
		panic(err) // let the recovery middleware deal with this
	}
}
