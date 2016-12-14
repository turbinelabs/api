package http

import (
	"fmt"
	"io"
	"net/http"
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
//
// The Endpoint object is configured with no custom headers (see
// Endpoint.AddHeader), and the net/http.Client created by
// HeaderPreservingClient. You may specify an alternate client via
// Endpoint.SetClient.
func NewEndpoint(protocol Protocol, host string, port int) (Endpoint, error) {
	url, err := url.Parse(fmt.Sprintf("%s://%s:%d", string(protocol), host, port))
	if err != nil {
		return Endpoint{}, err
	}

	return Endpoint{
		host:     host,
		port:     port,
		protocol: protocol,
		header:   http.Header{},
		client:   HeaderPreservingClient(),
		urlBase:  url,
	}, nil
}

type Endpoint struct {
	host     string
	port     int
	protocol Protocol
	header   http.Header
	client   *http.Client

	urlBase *url.URL // computed at construction
}

// Makes a copy of the Endpoint. Insures that modifications to custom
// headers of the new Endpoint are not made to the original Endpoint
// and vice versa.
func (e *Endpoint) Copy() Endpoint {
	headerCopy := make(http.Header, len(e.header))
	for header, values := range e.header {
		for _, value := range values {
			headerCopy.Add(header, value)
		}
	}

	newE := *e
	newE.header = headerCopy
	return newE
}

// Returns the net/http.Client for this Endpoint.
func (e *Endpoint) Client() *http.Client {
	return e.client
}

// Sets an alternative net/http.Client for this Endpoint.
func (e *Endpoint) SetClient(c *http.Client) {
	e.client = c
}

// Adds a header to be added to all requests created via NewRequest.
// These headers are meant to be constant across all requests (e.g. a
// client identifier). Headers specific to a particular request should
// be added directly to the net/http.Request.
func (e *Endpoint) AddHeader(header, value string) {
	e.header.Add(header, value)
}

// Construct a URL to this turbine Endpoint.
func (e *Endpoint) Url(path string, queryParams Params) string {
	newUrl := *e.urlBase
	newUrl.Path = path

	if len(queryParams) != 0 {
		q := newUrl.Query()
		for k, v := range queryParams {
			q.Set(k, v)
		}
		newUrl.RawQuery = q.Encode()
	}

	return newUrl.String()
}

// Construct a net/http.Request for this turbine Endpoint with the
// given method, path, (optional) query parameters and (optional)
// body. Headers previously configured via AddHeader are added
// automatically.
func (e *Endpoint) NewRequest(
	method string,
	path string,
	queryParams Params,
	body io.Reader,
) (*http.Request, error) {
	url := e.Url(path, queryParams)
	request, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}

	for header, values := range e.header {
		for _, value := range values {
			request.Header.Add(header, value)
		}
	}

	return request, nil
}
