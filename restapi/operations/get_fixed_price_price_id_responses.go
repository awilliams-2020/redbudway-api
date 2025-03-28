// Code generated by go-swagger; DO NOT EDIT.

package operations

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"
)

// GetFixedPricePriceIDOKCode is the HTTP code returned for type GetFixedPricePriceIDOK
const GetFixedPricePriceIDOKCode int = 200

/*GetFixedPricePriceIDOK Fixed Price Details

swagger:response getFixedPricePriceIdOK
*/
type GetFixedPricePriceIDOK struct {

	/*
	  In: Body
	*/
	Payload *GetFixedPricePriceIDOKBody `json:"body,omitempty"`
}

// NewGetFixedPricePriceIDOK creates GetFixedPricePriceIDOK with default headers values
func NewGetFixedPricePriceIDOK() *GetFixedPricePriceIDOK {

	return &GetFixedPricePriceIDOK{}
}

// WithPayload adds the payload to the get fixed price price Id o k response
func (o *GetFixedPricePriceIDOK) WithPayload(payload *GetFixedPricePriceIDOKBody) *GetFixedPricePriceIDOK {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get fixed price price Id o k response
func (o *GetFixedPricePriceIDOK) SetPayload(payload *GetFixedPricePriceIDOKBody) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetFixedPricePriceIDOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(200)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}
