// Code generated by go-swagger; DO NOT EDIT.

package operations

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	"redbudway-api/models"
)

// GetTradespersonTradespersonIDPromoPromoIDOKCode is the HTTP code returned for type GetTradespersonTradespersonIDPromoPromoIDOK
const GetTradespersonTradespersonIDPromoPromoIDOKCode int = 200

/*GetTradespersonTradespersonIDPromoPromoIDOK Promo

swagger:response getTradespersonTradespersonIdPromoPromoIdOK
*/
type GetTradespersonTradespersonIDPromoPromoIDOK struct {

	/*
	  In: Body
	*/
	Payload *models.Promo `json:"body,omitempty"`
}

// NewGetTradespersonTradespersonIDPromoPromoIDOK creates GetTradespersonTradespersonIDPromoPromoIDOK with default headers values
func NewGetTradespersonTradespersonIDPromoPromoIDOK() *GetTradespersonTradespersonIDPromoPromoIDOK {

	return &GetTradespersonTradespersonIDPromoPromoIDOK{}
}

// WithPayload adds the payload to the get tradesperson tradesperson Id promo promo Id o k response
func (o *GetTradespersonTradespersonIDPromoPromoIDOK) WithPayload(payload *models.Promo) *GetTradespersonTradespersonIDPromoPromoIDOK {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get tradesperson tradesperson Id promo promo Id o k response
func (o *GetTradespersonTradespersonIDPromoPromoIDOK) SetPayload(payload *models.Promo) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetTradespersonTradespersonIDPromoPromoIDOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(200)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}
