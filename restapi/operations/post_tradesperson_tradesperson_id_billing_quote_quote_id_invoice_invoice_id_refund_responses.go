// Code generated by go-swagger; DO NOT EDIT.

package operations

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"
)

// PostTradespersonTradespersonIDBillingQuoteQuoteIDInvoiceInvoiceIDRefundOKCode is the HTTP code returned for type PostTradespersonTradespersonIDBillingQuoteQuoteIDInvoiceInvoiceIDRefundOK
const PostTradespersonTradespersonIDBillingQuoteQuoteIDInvoiceInvoiceIDRefundOKCode int = 200

/*PostTradespersonTradespersonIDBillingQuoteQuoteIDInvoiceInvoiceIDRefundOK If invoice refunded

swagger:response postTradespersonTradespersonIdBillingQuoteQuoteIdInvoiceInvoiceIdRefundOK
*/
type PostTradespersonTradespersonIDBillingQuoteQuoteIDInvoiceInvoiceIDRefundOK struct {

	/*
	  In: Body
	*/
	Payload *PostTradespersonTradespersonIDBillingQuoteQuoteIDInvoiceInvoiceIDRefundOKBody `json:"body,omitempty"`
}

// NewPostTradespersonTradespersonIDBillingQuoteQuoteIDInvoiceInvoiceIDRefundOK creates PostTradespersonTradespersonIDBillingQuoteQuoteIDInvoiceInvoiceIDRefundOK with default headers values
func NewPostTradespersonTradespersonIDBillingQuoteQuoteIDInvoiceInvoiceIDRefundOK() *PostTradespersonTradespersonIDBillingQuoteQuoteIDInvoiceInvoiceIDRefundOK {

	return &PostTradespersonTradespersonIDBillingQuoteQuoteIDInvoiceInvoiceIDRefundOK{}
}

// WithPayload adds the payload to the post tradesperson tradesperson Id billing quote quote Id invoice invoice Id refund o k response
func (o *PostTradespersonTradespersonIDBillingQuoteQuoteIDInvoiceInvoiceIDRefundOK) WithPayload(payload *PostTradespersonTradespersonIDBillingQuoteQuoteIDInvoiceInvoiceIDRefundOKBody) *PostTradespersonTradespersonIDBillingQuoteQuoteIDInvoiceInvoiceIDRefundOK {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the post tradesperson tradesperson Id billing quote quote Id invoice invoice Id refund o k response
func (o *PostTradespersonTradespersonIDBillingQuoteQuoteIDInvoiceInvoiceIDRefundOK) SetPayload(payload *PostTradespersonTradespersonIDBillingQuoteQuoteIDInvoiceInvoiceIDRefundOKBody) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *PostTradespersonTradespersonIDBillingQuoteQuoteIDInvoiceInvoiceIDRefundOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(200)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}
