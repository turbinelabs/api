package http

import (
	"fmt"
	"net/url"
)

// Indicates which transport we should use for communicating with our service.
type Protocol string

const (
	HTTP  Protocol = "http"
	HTTPS Protocol = "https"
)

// Holds URL query arg name -> value mappings
type Params map[string]string

// Constructs a new HTTP Endpoint. This is used to configure the HTTP service
// implementation.
//
// Parameters:
// 	protocol
//	  specifies whether we should be using HTTP or HTTPS there is currently no
// 	  special configuration for HTTPS (certificate pinning, custom root CAs,
//	  etc.)
//	host - DNS entry, or IP, for the service we should connect to
// 	port - open port on the host above
//
// Returns a new Endpoint object and an error if there was a problem. Currently
// the only error possible is the result of a failed call to url.Parse which
// will be passed directly to the caller.
func NewEndpoint(protocol Protocol, host string, port int) (Endpoint, error) {
	url, err := url.Parse(fmt.Sprintf("%s://%s:%d", string(protocol), host, port))
	if err != nil {
		return Endpoint{}, err
	}

	return Endpoint{host, port, protocol, *url}, nil
}

type Endpoint struct {
	host     string
	port     int
	protocol Protocol

	urlBase url.URL // computed at construction
}

// construct a URL to this turbine Endpoint
func (e Endpoint) Url(path string, queryParams Params) string {
	newUrl := e.urlBase
	newUrl.Path = path

	if queryParams != nil && len(queryParams) != 0 {
		q := newUrl.Query()
		for k, v := range queryParams {
			q.Set(k, v)
		}
		newUrl.RawQuery = q.Encode()
	}

	return newUrl.String()
}
