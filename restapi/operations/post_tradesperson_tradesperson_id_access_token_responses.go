// Code generated by go-swagger; DO NOT EDIT.

package operations

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"
)

// PostTradespersonTradespersonIDAccessTokenOKCode is the HTTP code returned for type PostTradespersonTradespersonIDAccessTokenOK
const PostTradespersonTradespersonIDAccessTokenOKCode int = 200

/*PostTradespersonTradespersonIDAccessTokenOK Valid email and password

swagger:response postTradespersonTradespersonIdAccessTokenOK
*/
type PostTradespersonTradespersonIDAccessTokenOK struct {

	/*
	  In: Body
	*/
	Payload *PostTradespersonTradespersonIDAccessTokenOKBody `json:"body,omitempty"`
}

// NewPostTradespersonTradespersonIDAccessTokenOK creates PostTradespersonTradespersonIDAccessTokenOK with default headers values
func NewPostTradespersonTradespersonIDAccessTokenOK() *PostTradespersonTradespersonIDAccessTokenOK {

	return &PostTradespersonTradespersonIDAccessTokenOK{}
}

// WithPayload adds the payload to the post tradesperson tradesperson Id access token o k response
func (o *PostTradespersonTradespersonIDAccessTokenOK) WithPayload(payload *PostTradespersonTradespersonIDAccessTokenOKBody) *PostTradespersonTradespersonIDAccessTokenOK {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the post tradesperson tradesperson Id access token o k response
func (o *PostTradespersonTradespersonIDAccessTokenOK) SetPayload(payload *PostTradespersonTradespersonIDAccessTokenOKBody) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *PostTradespersonTradespersonIDAccessTokenOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(200)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}
