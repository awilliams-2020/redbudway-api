// Code generated by go-swagger; DO NOT EDIT.

package operations

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"errors"
	"net/url"
	golangswaggerpaths "path"
	"strings"
)

// GetTradespersonTradespersonIDBillingQuoteQuoteIDInvoiceInvoiceIDURL generates an URL for the get tradesperson tradesperson ID billing quote quote ID invoice invoice ID operation
type GetTradespersonTradespersonIDBillingQuoteQuoteIDInvoiceInvoiceIDURL struct {
	InvoiceID      string
	QuoteID        string
	TradespersonID string

	_basePath string
	// avoid unkeyed usage
	_ struct{}
}

// WithBasePath sets the base path for this url builder, only required when it's different from the
// base path specified in the swagger spec.
// When the value of the base path is an empty string
func (o *GetTradespersonTradespersonIDBillingQuoteQuoteIDInvoiceInvoiceIDURL) WithBasePath(bp string) *GetTradespersonTradespersonIDBillingQuoteQuoteIDInvoiceInvoiceIDURL {
	o.SetBasePath(bp)
	return o
}

// SetBasePath sets the base path for this url builder, only required when it's different from the
// base path specified in the swagger spec.
// When the value of the base path is an empty string
func (o *GetTradespersonTradespersonIDBillingQuoteQuoteIDInvoiceInvoiceIDURL) SetBasePath(bp string) {
	o._basePath = bp
}

// Build a url path and query string
func (o *GetTradespersonTradespersonIDBillingQuoteQuoteIDInvoiceInvoiceIDURL) Build() (*url.URL, error) {
	var _result url.URL

	var _path = "/tradesperson/{tradespersonId}/billing/quote/{quoteId}/invoice/{invoiceId}"

	invoiceID := o.InvoiceID
	if invoiceID != "" {
		_path = strings.Replace(_path, "{invoiceId}", invoiceID, -1)
	} else {
		return nil, errors.New("invoiceId is required on GetTradespersonTradespersonIDBillingQuoteQuoteIDInvoiceInvoiceIDURL")
	}

	quoteID := o.QuoteID
	if quoteID != "" {
		_path = strings.Replace(_path, "{quoteId}", quoteID, -1)
	} else {
		return nil, errors.New("quoteId is required on GetTradespersonTradespersonIDBillingQuoteQuoteIDInvoiceInvoiceIDURL")
	}

	tradespersonID := o.TradespersonID
	if tradespersonID != "" {
		_path = strings.Replace(_path, "{tradespersonId}", tradespersonID, -1)
	} else {
		return nil, errors.New("tradespersonId is required on GetTradespersonTradespersonIDBillingQuoteQuoteIDInvoiceInvoiceIDURL")
	}

	_basePath := o._basePath
	if _basePath == "" {
		_basePath = "/v1"
	}
	_result.Path = golangswaggerpaths.Join(_basePath, _path)

	return &_result, nil
}

// Must is a helper function to panic when the url builder returns an error
func (o *GetTradespersonTradespersonIDBillingQuoteQuoteIDInvoiceInvoiceIDURL) Must(u *url.URL, err error) *url.URL {
	if err != nil {
		panic(err)
	}
	if u == nil {
		panic("url can't be nil")
	}
	return u
}

// String returns the string representation of the path with query string
func (o *GetTradespersonTradespersonIDBillingQuoteQuoteIDInvoiceInvoiceIDURL) String() string {
	return o.Must(o.Build()).String()
}

// BuildFull builds a full url with scheme, host, path and query string
func (o *GetTradespersonTradespersonIDBillingQuoteQuoteIDInvoiceInvoiceIDURL) BuildFull(scheme, host string) (*url.URL, error) {
	if scheme == "" {
		return nil, errors.New("scheme is required for a full url on GetTradespersonTradespersonIDBillingQuoteQuoteIDInvoiceInvoiceIDURL")
	}
	if host == "" {
		return nil, errors.New("host is required for a full url on GetTradespersonTradespersonIDBillingQuoteQuoteIDInvoiceInvoiceIDURL")
	}

	base, err := o.Build()
	if err != nil {
		return nil, err
	}

	base.Scheme = scheme
	base.Host = host
	return base, nil
}

// StringFull returns the string representation of a complete url
func (o *GetTradespersonTradespersonIDBillingQuoteQuoteIDInvoiceInvoiceIDURL) StringFull(scheme, host string) string {
	return o.Must(o.BuildFull(scheme, host)).String()
}
