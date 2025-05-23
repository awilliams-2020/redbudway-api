// Code generated by go-swagger; DO NOT EDIT.

package operations

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	"redbudway-api/models"
)

// GetProfileVanityOrIDFixedPricesOKCode is the HTTP code returned for type GetProfileVanityOrIDFixedPricesOK
const GetProfileVanityOrIDFixedPricesOKCode int = 200

/*GetProfileVanityOrIDFixedPricesOK Tradesperson fixed-price services

swagger:response getProfileVanityOrIdFixedPricesOK
*/
type GetProfileVanityOrIDFixedPricesOK struct {

	/*
	  In: Body
	*/
	Payload []*models.Service `json:"body,omitempty"`
}

// NewGetProfileVanityOrIDFixedPricesOK creates GetProfileVanityOrIDFixedPricesOK with default headers values
func NewGetProfileVanityOrIDFixedPricesOK() *GetProfileVanityOrIDFixedPricesOK {

	return &GetProfileVanityOrIDFixedPricesOK{}
}

// WithPayload adds the payload to the get profile vanity or Id fixed prices o k response
func (o *GetProfileVanityOrIDFixedPricesOK) WithPayload(payload []*models.Service) *GetProfileVanityOrIDFixedPricesOK {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get profile vanity or Id fixed prices o k response
func (o *GetProfileVanityOrIDFixedPricesOK) SetPayload(payload []*models.Service) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetProfileVanityOrIDFixedPricesOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(200)
	payload := o.Payload
	if payload == nil {
		// return empty array
		payload = make([]*models.Service, 0, 50)
	}

	if err := producer.Produce(rw, payload); err != nil {
		panic(err) // let the recovery middleware deal with this
	}
}
