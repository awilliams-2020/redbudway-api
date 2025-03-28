// Code generated by go-swagger; DO NOT EDIT.

package operations

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"
)

// PutTradespersonTradespersonIDBillingQuoteQuoteIDOKCode is the HTTP code returned for type PutTradespersonTradespersonIDBillingQuoteQuoteIDOK
const PutTradespersonTradespersonIDBillingQuoteQuoteIDOKCode int = 200

/*PutTradespersonTradespersonIDBillingQuoteQuoteIDOK Updated quote

swagger:response putTradespersonTradespersonIdBillingQuoteQuoteIdOK
*/
type PutTradespersonTradespersonIDBillingQuoteQuoteIDOK struct {

	/*
	  In: Body
	*/
	Payload *PutTradespersonTradespersonIDBillingQuoteQuoteIDOKBody `json:"body,omitempty"`
}

// NewPutTradespersonTradespersonIDBillingQuoteQuoteIDOK creates PutTradespersonTradespersonIDBillingQuoteQuoteIDOK with default headers values
func NewPutTradespersonTradespersonIDBillingQuoteQuoteIDOK() *PutTradespersonTradespersonIDBillingQuoteQuoteIDOK {

	return &PutTradespersonTradespersonIDBillingQuoteQuoteIDOK{}
}

// WithPayload adds the payload to the put tradesperson tradesperson Id billing quote quote Id o k response
func (o *PutTradespersonTradespersonIDBillingQuoteQuoteIDOK) WithPayload(payload *PutTradespersonTradespersonIDBillingQuoteQuoteIDOKBody) *PutTradespersonTradespersonIDBillingQuoteQuoteIDOK {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the put tradesperson tradesperson Id billing quote quote Id o k response
func (o *PutTradespersonTradespersonIDBillingQuoteQuoteIDOK) SetPayload(payload *PutTradespersonTradespersonIDBillingQuoteQuoteIDOKBody) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *PutTradespersonTradespersonIDBillingQuoteQuoteIDOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(200)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}
