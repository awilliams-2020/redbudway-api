// Code generated by go-swagger; DO NOT EDIT.

package operations

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"
)

// GetTradespersonTradespersonIDBillingQuotesOKCode is the HTTP code returned for type GetTradespersonTradespersonIDBillingQuotesOK
const GetTradespersonTradespersonIDBillingQuotesOKCode int = 200

/*GetTradespersonTradespersonIDBillingQuotesOK quotes in a quarter of some year

swagger:response getTradespersonTradespersonIdBillingQuotesOK
*/
type GetTradespersonTradespersonIDBillingQuotesOK struct {

	/*
	  In: Body
	*/
	Payload []*GetTradespersonTradespersonIDBillingQuotesOKBodyItems0 `json:"body,omitempty"`
}

// NewGetTradespersonTradespersonIDBillingQuotesOK creates GetTradespersonTradespersonIDBillingQuotesOK with default headers values
func NewGetTradespersonTradespersonIDBillingQuotesOK() *GetTradespersonTradespersonIDBillingQuotesOK {

	return &GetTradespersonTradespersonIDBillingQuotesOK{}
}

// WithPayload adds the payload to the get tradesperson tradesperson Id billing quotes o k response
func (o *GetTradespersonTradespersonIDBillingQuotesOK) WithPayload(payload []*GetTradespersonTradespersonIDBillingQuotesOKBodyItems0) *GetTradespersonTradespersonIDBillingQuotesOK {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get tradesperson tradesperson Id billing quotes o k response
func (o *GetTradespersonTradespersonIDBillingQuotesOK) SetPayload(payload []*GetTradespersonTradespersonIDBillingQuotesOKBodyItems0) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetTradespersonTradespersonIDBillingQuotesOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(200)
	payload := o.Payload
	if payload == nil {
		// return empty array
		payload = make([]*GetTradespersonTradespersonIDBillingQuotesOKBodyItems0, 0, 50)
	}

	if err := producer.Produce(rw, payload); err != nil {
		panic(err) // let the recovery middleware deal with this
	}
}
