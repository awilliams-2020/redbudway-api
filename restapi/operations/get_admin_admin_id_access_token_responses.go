// Code generated by go-swagger; DO NOT EDIT.

package operations

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"
)

// GetAdminAdminIDAccessTokenOKCode is the HTTP code returned for type GetAdminAdminIDAccessTokenOK
const GetAdminAdminIDAccessTokenOKCode int = 200

/*GetAdminAdminIDAccessTokenOK A valid admin accessToken

swagger:response getAdminAdminIdAccessTokenOK
*/
type GetAdminAdminIDAccessTokenOK struct {

	/*
	  In: Body
	*/
	Payload *GetAdminAdminIDAccessTokenOKBody `json:"body,omitempty"`
}

// NewGetAdminAdminIDAccessTokenOK creates GetAdminAdminIDAccessTokenOK with default headers values
func NewGetAdminAdminIDAccessTokenOK() *GetAdminAdminIDAccessTokenOK {

	return &GetAdminAdminIDAccessTokenOK{}
}

// WithPayload adds the payload to the get admin admin Id access token o k response
func (o *GetAdminAdminIDAccessTokenOK) WithPayload(payload *GetAdminAdminIDAccessTokenOKBody) *GetAdminAdminIDAccessTokenOK {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get admin admin Id access token o k response
func (o *GetAdminAdminIDAccessTokenOK) SetPayload(payload *GetAdminAdminIDAccessTokenOKBody) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetAdminAdminIDAccessTokenOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(200)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}
