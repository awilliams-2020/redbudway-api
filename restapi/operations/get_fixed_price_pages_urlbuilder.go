// Code generated by go-swagger; DO NOT EDIT.

package operations

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"errors"
	"net/url"
	golangswaggerpaths "path"
)

// GetFixedPricePagesURL generates an URL for the get fixed price pages operation
type GetFixedPricePagesURL struct {
	Category    *string
	City        string
	Filters     *string
	State       string
	SubCategory *string

	_basePath string
	// avoid unkeyed usage
	_ struct{}
}

// WithBasePath sets the base path for this url builder, only required when it's different from the
// base path specified in the swagger spec.
// When the value of the base path is an empty string
func (o *GetFixedPricePagesURL) WithBasePath(bp string) *GetFixedPricePagesURL {
	o.SetBasePath(bp)
	return o
}

// SetBasePath sets the base path for this url builder, only required when it's different from the
// base path specified in the swagger spec.
// When the value of the base path is an empty string
func (o *GetFixedPricePagesURL) SetBasePath(bp string) {
	o._basePath = bp
}

// Build a url path and query string
func (o *GetFixedPricePagesURL) Build() (*url.URL, error) {
	var _result url.URL

	var _path = "/fixed-price/pages"

	_basePath := o._basePath
	if _basePath == "" {
		_basePath = "/v1"
	}
	_result.Path = golangswaggerpaths.Join(_basePath, _path)

	qs := make(url.Values)

	var categoryQ string
	if o.Category != nil {
		categoryQ = *o.Category
	}
	if categoryQ != "" {
		qs.Set("category", categoryQ)
	}

	cityQ := o.City
	if cityQ != "" {
		qs.Set("city", cityQ)
	}

	var filtersQ string
	if o.Filters != nil {
		filtersQ = *o.Filters
	}
	if filtersQ != "" {
		qs.Set("filters", filtersQ)
	}

	stateQ := o.State
	if stateQ != "" {
		qs.Set("state", stateQ)
	}

	var subCategoryQ string
	if o.SubCategory != nil {
		subCategoryQ = *o.SubCategory
	}
	if subCategoryQ != "" {
		qs.Set("subCategory", subCategoryQ)
	}

	_result.RawQuery = qs.Encode()

	return &_result, nil
}

// Must is a helper function to panic when the url builder returns an error
func (o *GetFixedPricePagesURL) Must(u *url.URL, err error) *url.URL {
	if err != nil {
		panic(err)
	}
	if u == nil {
		panic("url can't be nil")
	}
	return u
}

// String returns the string representation of the path with query string
func (o *GetFixedPricePagesURL) String() string {
	return o.Must(o.Build()).String()
}

// BuildFull builds a full url with scheme, host, path and query string
func (o *GetFixedPricePagesURL) BuildFull(scheme, host string) (*url.URL, error) {
	if scheme == "" {
		return nil, errors.New("scheme is required for a full url on GetFixedPricePagesURL")
	}
	if host == "" {
		return nil, errors.New("host is required for a full url on GetFixedPricePagesURL")
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
func (o *GetFixedPricePagesURL) StringFull(scheme, host string) string {
	return o.Must(o.BuildFull(scheme, host)).String()
}
