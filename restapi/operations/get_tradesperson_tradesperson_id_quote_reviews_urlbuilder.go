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

// GetTradespersonTradespersonIDQuoteReviewsURL generates an URL for the get tradesperson tradesperson ID quote reviews operation
type GetTradespersonTradespersonIDQuoteReviewsURL struct {
	TradespersonID string

	_basePath string
	// avoid unkeyed usage
	_ struct{}
}

// WithBasePath sets the base path for this url builder, only required when it's different from the
// base path specified in the swagger spec.
// When the value of the base path is an empty string
func (o *GetTradespersonTradespersonIDQuoteReviewsURL) WithBasePath(bp string) *GetTradespersonTradespersonIDQuoteReviewsURL {
	o.SetBasePath(bp)
	return o
}

// SetBasePath sets the base path for this url builder, only required when it's different from the
// base path specified in the swagger spec.
// When the value of the base path is an empty string
func (o *GetTradespersonTradespersonIDQuoteReviewsURL) SetBasePath(bp string) {
	o._basePath = bp
}

// Build a url path and query string
func (o *GetTradespersonTradespersonIDQuoteReviewsURL) Build() (*url.URL, error) {
	var _result url.URL

	var _path = "/tradesperson/{tradespersonId}/quote/reviews"

	tradespersonID := o.TradespersonID
	if tradespersonID != "" {
		_path = strings.Replace(_path, "{tradespersonId}", tradespersonID, -1)
	} else {
		return nil, errors.New("tradespersonId is required on GetTradespersonTradespersonIDQuoteReviewsURL")
	}

	_basePath := o._basePath
	if _basePath == "" {
		_basePath = "/v1"
	}
	_result.Path = golangswaggerpaths.Join(_basePath, _path)

	return &_result, nil
}

// Must is a helper function to panic when the url builder returns an error
func (o *GetTradespersonTradespersonIDQuoteReviewsURL) Must(u *url.URL, err error) *url.URL {
	if err != nil {
		panic(err)
	}
	if u == nil {
		panic("url can't be nil")
	}
	return u
}

// String returns the string representation of the path with query string
func (o *GetTradespersonTradespersonIDQuoteReviewsURL) String() string {
	return o.Must(o.Build()).String()
}

// BuildFull builds a full url with scheme, host, path and query string
func (o *GetTradespersonTradespersonIDQuoteReviewsURL) BuildFull(scheme, host string) (*url.URL, error) {
	if scheme == "" {
		return nil, errors.New("scheme is required for a full url on GetTradespersonTradespersonIDQuoteReviewsURL")
	}
	if host == "" {
		return nil, errors.New("host is required for a full url on GetTradespersonTradespersonIDQuoteReviewsURL")
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
func (o *GetTradespersonTradespersonIDQuoteReviewsURL) StringFull(scheme, host string) string {
	return o.Must(o.BuildFull(scheme, host)).String()
}
