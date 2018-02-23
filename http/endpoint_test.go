/*
Copyright 2018 Turbine Labs, Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package http

import (
	"net/http"
	"testing"

	"github.com/turbinelabs/test/assert"
	"github.com/turbinelabs/test/io"
)

func TestNewEndpointHttp(t *testing.T) {
	e, err := NewEndpoint(HTTP, "example.com:80")
	assert.Nil(t, err)
	assert.Equal(t, e.hostPort, "example.com:80")
	assert.Equal(t, e.protocol, HTTP)

	assert.NonNil(t, e.urlBase)
	assert.Equal(t, e.urlBase.String(), "http://example.com:80")
}

func TestNewEndpointHttps(t *testing.T) {
	e, err := NewEndpoint(HTTPS, "example.com:443")
	assert.Nil(t, err)
	assert.Equal(t, e.hostPort, "example.com:443")
	assert.Equal(t, e.protocol, HTTPS)

	assert.NonNil(t, e.urlBase)
	assert.Equal(t, e.urlBase.String(), "https://example.com:443")
}

func TestNewEndpointParseError(t *testing.T) {
	e, err := NewEndpoint(HTTP, "not a domain:99")
	assert.NonNil(t, err)
	assert.DeepEqual(t, e, Endpoint{})
}

func TestEndpointClient(t *testing.T) {
	e, _ := NewEndpoint(HTTP, "example.com:80")
	otherClient := &http.Client{}

	assert.NonNil(t, e.Client())
	assert.NotSameInstance(t, e.Client(), otherClient)

	e.SetClient(otherClient)
	assert.SameInstance(t, e.Client(), otherClient)
}

func TestEndpointAddHeader(t *testing.T) {
	e, _ := NewEndpoint(HTTP, "example.com:80")
	assert.Equal(t, len(e.header), 0)

	e.AddHeader("foo", "1")
	e.AddHeader("foo", "2")
	e.AddHeader("bar", "3")

	assert.Equal(t, len(e.header), 2)
	assert.ArrayEqual(t, e.header["Foo"], []string{"1", "2"})
	assert.ArrayEqual(t, e.header["Bar"], []string{"3"})
}

func TestEndpointURL(t *testing.T) {
	e, _ := NewEndpoint(HTTP, "example.com:80")
	u := e.URL("/admin/user", Params{})
	assert.Equal(t, u, "http://example.com:80/admin/user")

	u2 := e.URL("/admin/user", Params{"q": "encode me!"})
	assert.Equal(t, u2, "http://example.com:80/admin/user?q=encode+me%21")
}

func TestEndpointNewRequestWithParams(t *testing.T) {
	e, _ := NewEndpoint(HTTP, "example.com:80")

	params := Params{"uid": "123"}
	r, err := e.NewRequest("GET", "/admin/user", params, nil)
	assert.Nil(t, err)
	assert.Equal(t, r.Method, "GET")
	assert.Equal(t, r.URL.String(), e.URL("/admin/user", params))
	assert.Nil(t, r.Body)
	assert.Equal(t, len(r.Header), 0)
}

func TestEndpointNewRequestWithHeader(t *testing.T) {
	e, _ := NewEndpoint(HTTP, "example.com:80")
	e.AddHeader("my-header", "my-value")
	e.AddHeader("host", "other.example.com")

	r, err := e.NewRequest("GET", "/admin/user", Params{}, nil)
	assert.Nil(t, err)
	assert.Equal(t, r.Method, "GET")
	assert.Equal(t, r.Host, "other.example.com")
	assert.Equal(t, r.URL.String(), e.URL("/admin/user", Params{}))
	assert.Nil(t, r.Body)
	assert.Equal(t, r.Header.Get(http.CanonicalHeaderKey("my-header")), "my-value")
}

func TestEndpointNewRequestWithBody(t *testing.T) {
	e, _ := NewEndpoint(HTTP, "example.com:80")

	// Use a ReadCloser so net/http doesn't do any wrapping
	body := io.NewFailingReader()

	r, err := e.NewRequest("POST", "/admin/user", Params{}, body)
	assert.Nil(t, err)
	assert.Equal(t, r.Method, "POST")
	assert.Equal(t, r.URL.String(), e.URL("/admin/user", Params{}))
	assert.SameInstance(t, r.Body, body)
	assert.Equal(t, len(r.Header), 0)
}

func TestEndpointNewRequestError(t *testing.T) {
	e, _ := NewEndpoint(HTTP, "example.com:80")

	newURLBase := *e.urlBase
	newURLBase.Host = "not a domain, hoss"

	e.urlBase = &newURLBase

	r, err := e.NewRequest("GET", "/", Params{}, nil)
	assert.Nil(t, r)
	assert.NonNil(t, err)
}

func TestEndpointCopy(t *testing.T) {
	simple, _ := NewEndpoint(HTTPS, "example.com:443")
	simpleCopy := simple.Copy()
	assert.DeepEqual(t, simpleCopy, simple)

	e, _ := NewEndpoint(HTTP, "example.com:80")
	e.AddHeader("before", "preserved")

	e2 := e.Copy()
	e2.AddHeader("before", "private-to-e2")
	e2.AddHeader("after", "private-to-e2")

	e.AddHeader("before", "private-to-e")
	e.AddHeader("after", "private-to-e")

	assert.Equal(t, len(e.header), 2)
	assert.ArrayEqual(t, e.header["Before"], []string{"preserved", "private-to-e"})
	assert.ArrayEqual(t, e.header["After"], []string{"private-to-e"})

	assert.Equal(t, len(e2.header), 2)
	assert.ArrayEqual(t, e2.header["Before"], []string{"preserved", "private-to-e2"})
	assert.ArrayEqual(t, e2.header["After"], []string{"private-to-e2"})

}
