// Code generated by go-swagger; DO NOT EDIT.

package operations

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"
)

// GetTradespersonTradespersonIDAccessTokenOKCode is the HTTP code returned for type GetTradespersonTradespersonIDAccessTokenOK
const GetTradespersonTradespersonIDAccessTokenOKCode int = 200

/*GetTradespersonTradespersonIDAccessTokenOK A valid tradesperson accessToken

swagger:response getTradespersonTradespersonIdAccessTokenOK
*/
type GetTradespersonTradespersonIDAccessTokenOK struct {

	/*
	  In: Body
	*/
	Payload *GetTradespersonTradespersonIDAccessTokenOKBody `json:"body,omitempty"`
}

// NewGetTradespersonTradespersonIDAccessTokenOK creates GetTradespersonTradespersonIDAccessTokenOK with default headers values
func NewGetTradespersonTradespersonIDAccessTokenOK() *GetTradespersonTradespersonIDAccessTokenOK {

	return &GetTradespersonTradespersonIDAccessTokenOK{}
}

// WithPayload adds the payload to the get tradesperson tradesperson Id access token o k response
func (o *GetTradespersonTradespersonIDAccessTokenOK) WithPayload(payload *GetTradespersonTradespersonIDAccessTokenOKBody) *GetTradespersonTradespersonIDAccessTokenOK {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get tradesperson tradesperson Id access token o k response
func (o *GetTradespersonTradespersonIDAccessTokenOK) SetPayload(payload *GetTradespersonTradespersonIDAccessTokenOKBody) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetTradespersonTradespersonIDAccessTokenOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(200)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}
