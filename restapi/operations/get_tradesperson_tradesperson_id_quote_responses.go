// Code generated by go-swagger; DO NOT EDIT.

package operations

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	"redbudway-api/models"
)

// GetTradespersonTradespersonIDQuoteOKCode is the HTTP code returned for type GetTradespersonTradespersonIDQuoteOK
const GetTradespersonTradespersonIDQuoteOKCode int = 200

/*GetTradespersonTradespersonIDQuoteOK Tradesperson fixed-price services

swagger:response getTradespersonTradespersonIdQuoteOK
*/
type GetTradespersonTradespersonIDQuoteOK struct {

	/*
	  In: Body
	*/
	Payload *models.ServiceDetails `json:"body,omitempty"`
}

// NewGetTradespersonTradespersonIDQuoteOK creates GetTradespersonTradespersonIDQuoteOK with default headers values
func NewGetTradespersonTradespersonIDQuoteOK() *GetTradespersonTradespersonIDQuoteOK {

	return &GetTradespersonTradespersonIDQuoteOK{}
}

// WithPayload adds the payload to the get tradesperson tradesperson Id quote o k response
func (o *GetTradespersonTradespersonIDQuoteOK) WithPayload(payload *models.ServiceDetails) *GetTradespersonTradespersonIDQuoteOK {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get tradesperson tradesperson Id quote o k response
func (o *GetTradespersonTradespersonIDQuoteOK) SetPayload(payload *models.ServiceDetails) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetTradespersonTradespersonIDQuoteOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(200)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}
