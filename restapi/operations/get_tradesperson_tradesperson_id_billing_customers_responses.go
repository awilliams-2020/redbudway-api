// Code generated by go-swagger; DO NOT EDIT.

package operations

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	"redbudway-api/models"
)

// GetTradespersonTradespersonIDBillingCustomersOKCode is the HTTP code returned for type GetTradespersonTradespersonIDBillingCustomersOK
const GetTradespersonTradespersonIDBillingCustomersOKCode int = 200

/*GetTradespersonTradespersonIDBillingCustomersOK Customers that've done business

swagger:response getTradespersonTradespersonIdBillingCustomersOK
*/
type GetTradespersonTradespersonIDBillingCustomersOK struct {

	/*
	  In: Body
	*/
	Payload []*models.Customer `json:"body,omitempty"`
}

// NewGetTradespersonTradespersonIDBillingCustomersOK creates GetTradespersonTradespersonIDBillingCustomersOK with default headers values
func NewGetTradespersonTradespersonIDBillingCustomersOK() *GetTradespersonTradespersonIDBillingCustomersOK {

	return &GetTradespersonTradespersonIDBillingCustomersOK{}
}

// WithPayload adds the payload to the get tradesperson tradesperson Id billing customers o k response
func (o *GetTradespersonTradespersonIDBillingCustomersOK) WithPayload(payload []*models.Customer) *GetTradespersonTradespersonIDBillingCustomersOK {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get tradesperson tradesperson Id billing customers o k response
func (o *GetTradespersonTradespersonIDBillingCustomersOK) SetPayload(payload []*models.Customer) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetTradespersonTradespersonIDBillingCustomersOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(200)
	payload := o.Payload
	if payload == nil {
		// return empty array
		payload = make([]*models.Customer, 0, 50)
	}

	if err := producer.Produce(rw, payload); err != nil {
		panic(err) // let the recovery middleware deal with this
	}
}
