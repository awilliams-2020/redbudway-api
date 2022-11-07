// Code generated by go-swagger; DO NOT EDIT.

package operations

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	"redbudway-api/models"
)

// GetTradespersonTradespersonIDOKCode is the HTTP code returned for type GetTradespersonTradespersonIDOK
const GetTradespersonTradespersonIDOKCode int = 200

/*GetTradespersonTradespersonIDOK A tradesperson

swagger:response getTradespersonTradespersonIdOK
*/
type GetTradespersonTradespersonIDOK struct {

	/*
	  In: Body
	*/
	Payload *models.Tradesperson `json:"body,omitempty"`
}

// NewGetTradespersonTradespersonIDOK creates GetTradespersonTradespersonIDOK with default headers values
func NewGetTradespersonTradespersonIDOK() *GetTradespersonTradespersonIDOK {

	return &GetTradespersonTradespersonIDOK{}
}

// WithPayload adds the payload to the get tradesperson tradesperson Id o k response
func (o *GetTradespersonTradespersonIDOK) WithPayload(payload *models.Tradesperson) *GetTradespersonTradespersonIDOK {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get tradesperson tradesperson Id o k response
func (o *GetTradespersonTradespersonIDOK) SetPayload(payload *models.Tradesperson) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetTradespersonTradespersonIDOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(200)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}
