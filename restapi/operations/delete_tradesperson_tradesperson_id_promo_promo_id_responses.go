// Code generated by go-swagger; DO NOT EDIT.

package operations

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"
)

// DeleteTradespersonTradespersonIDPromoPromoIDOKCode is the HTTP code returned for type DeleteTradespersonTradespersonIDPromoPromoIDOK
const DeleteTradespersonTradespersonIDPromoPromoIDOKCode int = 200

/*DeleteTradespersonTradespersonIDPromoPromoIDOK If promo was deleted

swagger:response deleteTradespersonTradespersonIdPromoPromoIdOK
*/
type DeleteTradespersonTradespersonIDPromoPromoIDOK struct {

	/*
	  In: Body
	*/
	Payload *DeleteTradespersonTradespersonIDPromoPromoIDOKBody `json:"body,omitempty"`
}

// NewDeleteTradespersonTradespersonIDPromoPromoIDOK creates DeleteTradespersonTradespersonIDPromoPromoIDOK with default headers values
func NewDeleteTradespersonTradespersonIDPromoPromoIDOK() *DeleteTradespersonTradespersonIDPromoPromoIDOK {

	return &DeleteTradespersonTradespersonIDPromoPromoIDOK{}
}

// WithPayload adds the payload to the delete tradesperson tradesperson Id promo promo Id o k response
func (o *DeleteTradespersonTradespersonIDPromoPromoIDOK) WithPayload(payload *DeleteTradespersonTradespersonIDPromoPromoIDOKBody) *DeleteTradespersonTradespersonIDPromoPromoIDOK {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the delete tradesperson tradesperson Id promo promo Id o k response
func (o *DeleteTradespersonTradespersonIDPromoPromoIDOK) SetPayload(payload *DeleteTradespersonTradespersonIDPromoPromoIDOKBody) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *DeleteTradespersonTradespersonIDPromoPromoIDOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(200)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}
