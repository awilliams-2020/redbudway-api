// Code generated by go-swagger; DO NOT EDIT.

package operations

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"
)

// PostTradespersonTradespersonIDCouponCouponIDPromoOKCode is the HTTP code returned for type PostTradespersonTradespersonIDCouponCouponIDPromoOK
const PostTradespersonTradespersonIDCouponCouponIDPromoOKCode int = 200

/*PostTradespersonTradespersonIDCouponCouponIDPromoOK If promo was created

swagger:response postTradespersonTradespersonIdCouponCouponIdPromoOK
*/
type PostTradespersonTradespersonIDCouponCouponIDPromoOK struct {

	/*
	  In: Body
	*/
	Payload *PostTradespersonTradespersonIDCouponCouponIDPromoOKBody `json:"body,omitempty"`
}

// NewPostTradespersonTradespersonIDCouponCouponIDPromoOK creates PostTradespersonTradespersonIDCouponCouponIDPromoOK with default headers values
func NewPostTradespersonTradespersonIDCouponCouponIDPromoOK() *PostTradespersonTradespersonIDCouponCouponIDPromoOK {

	return &PostTradespersonTradespersonIDCouponCouponIDPromoOK{}
}

// WithPayload adds the payload to the post tradesperson tradesperson Id coupon coupon Id promo o k response
func (o *PostTradespersonTradespersonIDCouponCouponIDPromoOK) WithPayload(payload *PostTradespersonTradespersonIDCouponCouponIDPromoOKBody) *PostTradespersonTradespersonIDCouponCouponIDPromoOK {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the post tradesperson tradesperson Id coupon coupon Id promo o k response
func (o *PostTradespersonTradespersonIDCouponCouponIDPromoOK) SetPayload(payload *PostTradespersonTradespersonIDCouponCouponIDPromoOKBody) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *PostTradespersonTradespersonIDCouponCouponIDPromoOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(200)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}
