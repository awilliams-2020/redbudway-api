// Code generated by go-swagger; DO NOT EDIT.

package operations

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"
)

// GetCustomerCustomerIDAccessTokenOKCode is the HTTP code returned for type GetCustomerCustomerIDAccessTokenOK
const GetCustomerCustomerIDAccessTokenOKCode int = 200

/*GetCustomerCustomerIDAccessTokenOK A valid customer accessToken

swagger:response getCustomerCustomerIdAccessTokenOK
*/
type GetCustomerCustomerIDAccessTokenOK struct {

	/*
	  In: Body
	*/
	Payload *GetCustomerCustomerIDAccessTokenOKBody `json:"body,omitempty"`
}

// NewGetCustomerCustomerIDAccessTokenOK creates GetCustomerCustomerIDAccessTokenOK with default headers values
func NewGetCustomerCustomerIDAccessTokenOK() *GetCustomerCustomerIDAccessTokenOK {

	return &GetCustomerCustomerIDAccessTokenOK{}
}

// WithPayload adds the payload to the get customer customer Id access token o k response
func (o *GetCustomerCustomerIDAccessTokenOK) WithPayload(payload *GetCustomerCustomerIDAccessTokenOKBody) *GetCustomerCustomerIDAccessTokenOK {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get customer customer Id access token o k response
func (o *GetCustomerCustomerIDAccessTokenOK) SetPayload(payload *GetCustomerCustomerIDAccessTokenOKBody) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetCustomerCustomerIDAccessTokenOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(200)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}
