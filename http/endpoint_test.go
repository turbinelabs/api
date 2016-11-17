package http

import (
	"testing"

	"github.com/turbinelabs/test/assert"
)

func TestNewEndpointHttp(t *testing.T) {
	e, err := NewEndpoint(HTTP, "example.com", 80)
	assert.Nil(t, err)
	assert.Equal(t, e.host, "example.com")
	assert.Equal(t, e.port, 80)
	assert.Equal(t, e.protocol, HTTP)

	assert.NonNil(t, e.urlBase)
	assert.Equal(t, e.urlBase.String(), "http://example.com:80")
}

func TestNewEndpointHttps(t *testing.T) {
	e, err := NewEndpoint(HTTPS, "example.com", 443)
	assert.Nil(t, err)
	assert.Equal(t, e.host, "example.com")
	assert.Equal(t, e.port, 443)
	assert.Equal(t, e.protocol, HTTPS)

	assert.NonNil(t, e.urlBase)
	assert.Equal(t, e.urlBase.String(), "https://example.com:443")
}

func TestNewEndpointParseError(t *testing.T) {
	e, err := NewEndpoint(HTTP, "not a domain", 99)
	assert.NonNil(t, err)
	assert.DeepEqual(t, e, Endpoint{})
}

func TestEndpointUrl(t *testing.T) {
	e, _ := NewEndpoint(HTTP, "example.com", 80)
	u := e.Url("/admin/user", Params{})
	assert.Equal(t, u, "http://example.com:80/admin/user")

	u2 := e.Url("/admin/user", Params{"q": "encode me!"})
	assert.Equal(t, u2, "http://example.com:80/admin/user?q=encode+me%21")
}
