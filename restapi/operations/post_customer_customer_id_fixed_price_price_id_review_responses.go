// Code generated by go-swagger; DO NOT EDIT.

package operations

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"
)

// PostCustomerCustomerIDFixedPricePriceIDReviewOKCode is the HTTP code returned for type PostCustomerCustomerIDFixedPricePriceIDReviewOK
const PostCustomerCustomerIDFixedPricePriceIDReviewOKCode int = 200

/*PostCustomerCustomerIDFixedPricePriceIDReviewOK Reviewed service

swagger:response postCustomerCustomerIdFixedPricePriceIdReviewOK
*/
type PostCustomerCustomerIDFixedPricePriceIDReviewOK struct {

	/*
	  In: Body
	*/
	Payload *PostCustomerCustomerIDFixedPricePriceIDReviewOKBody `json:"body,omitempty"`
}

// NewPostCustomerCustomerIDFixedPricePriceIDReviewOK creates PostCustomerCustomerIDFixedPricePriceIDReviewOK with default headers values
func NewPostCustomerCustomerIDFixedPricePriceIDReviewOK() *PostCustomerCustomerIDFixedPricePriceIDReviewOK {

	return &PostCustomerCustomerIDFixedPricePriceIDReviewOK{}
}

// WithPayload adds the payload to the post customer customer Id fixed price price Id review o k response
func (o *PostCustomerCustomerIDFixedPricePriceIDReviewOK) WithPayload(payload *PostCustomerCustomerIDFixedPricePriceIDReviewOKBody) *PostCustomerCustomerIDFixedPricePriceIDReviewOK {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the post customer customer Id fixed price price Id review o k response
func (o *PostCustomerCustomerIDFixedPricePriceIDReviewOK) SetPayload(payload *PostCustomerCustomerIDFixedPricePriceIDReviewOKBody) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *PostCustomerCustomerIDFixedPricePriceIDReviewOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(200)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}
