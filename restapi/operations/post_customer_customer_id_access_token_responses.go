// Code generated by go-swagger; DO NOT EDIT.

package operations

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"
)

// PostCustomerCustomerIDAccessTokenOKCode is the HTTP code returned for type PostCustomerCustomerIDAccessTokenOK
const PostCustomerCustomerIDAccessTokenOKCode int = 200

/*PostCustomerCustomerIDAccessTokenOK Valid email and password

swagger:response postCustomerCustomerIdAccessTokenOK
*/
type PostCustomerCustomerIDAccessTokenOK struct {

	/*
	  In: Body
	*/
	Payload *PostCustomerCustomerIDAccessTokenOKBody `json:"body,omitempty"`
}

// NewPostCustomerCustomerIDAccessTokenOK creates PostCustomerCustomerIDAccessTokenOK with default headers values
func NewPostCustomerCustomerIDAccessTokenOK() *PostCustomerCustomerIDAccessTokenOK {

	return &PostCustomerCustomerIDAccessTokenOK{}
}

// WithPayload adds the payload to the post customer customer Id access token o k response
func (o *PostCustomerCustomerIDAccessTokenOK) WithPayload(payload *PostCustomerCustomerIDAccessTokenOKBody) *PostCustomerCustomerIDAccessTokenOK {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the post customer customer Id access token o k response
func (o *PostCustomerCustomerIDAccessTokenOK) SetPayload(payload *PostCustomerCustomerIDAccessTokenOKBody) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *PostCustomerCustomerIDAccessTokenOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(200)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}
