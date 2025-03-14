// Code generated by go-swagger; DO NOT EDIT.

package operations

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"
)

// PostCustomerCustomerIDFixedPricePriceIDBookCreatedCode is the HTTP code returned for type PostCustomerCustomerIDFixedPricePriceIDBookCreated
const PostCustomerCustomerIDFixedPricePriceIDBookCreatedCode int = 201

/*PostCustomerCustomerIDFixedPricePriceIDBookCreated If fixed price was booked

swagger:response postCustomerCustomerIdFixedPricePriceIdBookCreated
*/
type PostCustomerCustomerIDFixedPricePriceIDBookCreated struct {

	/*
	  In: Body
	*/
	Payload *PostCustomerCustomerIDFixedPricePriceIDBookCreatedBody `json:"body,omitempty"`
}

// NewPostCustomerCustomerIDFixedPricePriceIDBookCreated creates PostCustomerCustomerIDFixedPricePriceIDBookCreated with default headers values
func NewPostCustomerCustomerIDFixedPricePriceIDBookCreated() *PostCustomerCustomerIDFixedPricePriceIDBookCreated {

	return &PostCustomerCustomerIDFixedPricePriceIDBookCreated{}
}

// WithPayload adds the payload to the post customer customer Id fixed price price Id book created response
func (o *PostCustomerCustomerIDFixedPricePriceIDBookCreated) WithPayload(payload *PostCustomerCustomerIDFixedPricePriceIDBookCreatedBody) *PostCustomerCustomerIDFixedPricePriceIDBookCreated {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the post customer customer Id fixed price price Id book created response
func (o *PostCustomerCustomerIDFixedPricePriceIDBookCreated) SetPayload(payload *PostCustomerCustomerIDFixedPricePriceIDBookCreatedBody) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *PostCustomerCustomerIDFixedPricePriceIDBookCreated) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(201)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}
