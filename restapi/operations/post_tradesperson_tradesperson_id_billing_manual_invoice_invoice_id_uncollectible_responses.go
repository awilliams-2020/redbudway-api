// Code generated by go-swagger; DO NOT EDIT.

package operations

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"
)

// PostTradespersonTradespersonIDBillingManualInvoiceInvoiceIDUncollectibleOKCode is the HTTP code returned for type PostTradespersonTradespersonIDBillingManualInvoiceInvoiceIDUncollectibleOK
const PostTradespersonTradespersonIDBillingManualInvoiceInvoiceIDUncollectibleOKCode int = 200

/*PostTradespersonTradespersonIDBillingManualInvoiceInvoiceIDUncollectibleOK If manual invoice was marked uncollectible

swagger:response postTradespersonTradespersonIdBillingManualInvoiceInvoiceIdUncollectibleOK
*/
type PostTradespersonTradespersonIDBillingManualInvoiceInvoiceIDUncollectibleOK struct {

	/*
	  In: Body
	*/
	Payload *PostTradespersonTradespersonIDBillingManualInvoiceInvoiceIDUncollectibleOKBody `json:"body,omitempty"`
}

// NewPostTradespersonTradespersonIDBillingManualInvoiceInvoiceIDUncollectibleOK creates PostTradespersonTradespersonIDBillingManualInvoiceInvoiceIDUncollectibleOK with default headers values
func NewPostTradespersonTradespersonIDBillingManualInvoiceInvoiceIDUncollectibleOK() *PostTradespersonTradespersonIDBillingManualInvoiceInvoiceIDUncollectibleOK {

	return &PostTradespersonTradespersonIDBillingManualInvoiceInvoiceIDUncollectibleOK{}
}

// WithPayload adds the payload to the post tradesperson tradesperson Id billing manual invoice invoice Id uncollectible o k response
func (o *PostTradespersonTradespersonIDBillingManualInvoiceInvoiceIDUncollectibleOK) WithPayload(payload *PostTradespersonTradespersonIDBillingManualInvoiceInvoiceIDUncollectibleOKBody) *PostTradespersonTradespersonIDBillingManualInvoiceInvoiceIDUncollectibleOK {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the post tradesperson tradesperson Id billing manual invoice invoice Id uncollectible o k response
func (o *PostTradespersonTradespersonIDBillingManualInvoiceInvoiceIDUncollectibleOK) SetPayload(payload *PostTradespersonTradespersonIDBillingManualInvoiceInvoiceIDUncollectibleOKBody) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *PostTradespersonTradespersonIDBillingManualInvoiceInvoiceIDUncollectibleOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(200)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}
