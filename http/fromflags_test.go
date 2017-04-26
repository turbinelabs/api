/*
Copyright 2017 Turbine Labs, Inc.

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
	"fmt"
	"net/http"
	"testing"

	tbnflag "github.com/turbinelabs/nonstdlib/flag"
	"github.com/turbinelabs/nonstdlib/ptr"
	"github.com/turbinelabs/test/assert"
)

func TestNewFromFlags(t *testing.T) {
	flagset := tbnflag.NewTestFlagSet()

	ff := NewFromFlags("api.turbinelabs.io", flagset.Scope("api", "API"))
	ffImpl := ff.(*fromFlags)

	assert.Equal(t, ffImpl.hostPort, "api.turbinelabs.io")

	flagset.Parse([]string{
		"-api.host=example.com:999",
		"-api.ssl=false",
		"-api.insecure=true",
		"-api.header=fred: flintstone",
		"-api.header=barney: rubble",
	})

	assert.Equal(t, ffImpl.hostPort, "example.com:999")
	assert.False(t, ffImpl.ssl)
	assert.True(t, ffImpl.insecure)
	assert.ArrayEqual(t, ffImpl.headers.Strings, []string{"fred: flintstone", "barney: rubble"})
}

func TestFromFlagsValidate(t *testing.T) {
	ff := &fromFlags{
		hostPort: "example.com:443",
		headers:  tbnflag.NewStrings(),
	}

	assert.Nil(t, ff.Validate())

	ff.hostPort = "example.com:::443"
	assert.ErrorContains(t, ff.Validate(), "too many colons")

	ff.hostPort = "example.com:99999"
	assert.ErrorContains(t, ff.Validate(), "invalid port")

	ff.hostPort = "example.com:"
	assert.Nil(t, ff.Validate())
	assert.Equal(t, ff.hostPort, "example.com:80")

	ff.ssl = true
	ff.hostPort = "example.com"
	assert.Nil(t, ff.Validate())
	assert.Equal(t, ff.hostPort, "example.com:443")

	ff.hostPort = "example.com:443"
	ff.headers.Set("not a header")
	assert.ErrorContains(t, ff.Validate(), "invalid header")

	ff.headers.ResetDefault()
	ff.headers.Set("X-Header: Value")
	assert.Nil(t, ff.Validate())
}

func TestFromFlagsMakeClient(t *testing.T) {
	ff := &fromFlags{}
	client := ff.makeClient()
	assert.Nil(t, client.Transport)
}

func TestFromFlagsMakeClientInsecure(t *testing.T) {
	ff := &fromFlags{insecure: true}
	client := ff.makeClient()

	switch transport := client.Transport.(type) {
	case *http.Transport:
		assert.True(t, transport.TLSClientConfig.InsecureSkipVerify)

	default:
		t.Fatal("bad type")
	}
}

func TestFromFlagsMakeEndpoint(t *testing.T) {
	ff := &fromFlags{
		hostPort: "example.com:80",
		headers:  tbnflag.NewStrings(),
	}

	ff.headers.Set("x-fred: flintstone,x-barney: rubble")

	e, err := ff.MakeEndpoint()
	assert.Nil(t, err)
	assert.Equal(t, e.urlBase.String(), "http://example.com:80")
	assert.NonNil(t, e.client)
	assert.Equal(t, e.header.Get("X-Fred"), "flintstone")
	assert.Equal(t, e.header.Get("X-Barney"), "rubble")

	ff.hostPort = "example.com:443"
	ff.ssl = true

	e, err = ff.MakeEndpoint()
	assert.Nil(t, err)
	assert.Equal(t, e.urlBase.String(), "https://example.com:443")
	assert.NonNil(t, e.client)

	ff.hostPort = "not a domain"

	_, err = ff.MakeEndpoint()
	assert.NonNil(t, err)
}

func TestCheckHostPort(t *testing.T) {
	testCases := []struct {
		hostPort              string
		ssl                   bool
		expectedHostPort      string
		expectedErrorContains *string
	}{
		{":80", false, "", ptr.String("missing hostname or address")},
		{":http", false, "", ptr.String("missing hostname or address")},
		{":::999", false, "", ptr.String("too many colons")},
		{"x:99999", false, "", ptr.String("invalid port")},
		{"x:neverheardofit", false, "", ptr.String("invalid port")},
		{"[::1", false, "", ptr.String("missing ']' in address")},
		{"[::1]", false, "[::1]:80", nil},
		{"[::1]:", false, "[::1]:80", nil},
		{"[::1]:999", false, "[::1]:999", nil},
		{"[::1]:http", false, "[::1]:http", nil},
		{"[::1]", true, "[::1]:443", nil},
		{"[::1]:", true, "[::1]:443", nil},
		{"[::1]:999", true, "[::1]:999", nil},
		{"[::1]:http", true, "[::1]:http", nil},
		{"10.0.0.1", false, "10.0.0.1:80", nil},
		{"10.0.0.1:", false, "10.0.0.1:80", nil},
		{"10.0.0.1:999", false, "10.0.0.1:999", nil},
		{"10.0.0.1:http", false, "10.0.0.1:http", nil},
		{"10.0.0.1", true, "10.0.0.1:443", nil},
		{"10.0.0.1:", true, "10.0.0.1:443", nil},
		{"10.0.0.1:999", true, "10.0.0.1:999", nil},
		{"10.0.0.1:http", true, "10.0.0.1:http", nil},
		{"example.com", false, "example.com:80", nil},
		{"example.com:", false, "example.com:80", nil},
		{"example.com", true, "example.com:443", nil},
		{"example.com:", true, "example.com:443", nil},
		{"localhost:1234", false, "localhost:1234", nil},
		{"localhost:1234", true, "localhost:1234", nil},
		{"example.com:http", false, "example.com:http", nil},
		{"example.com:http", true, "example.com:http", nil},
	}

	for _, tc := range testCases {
		assert.Group(
			fmt.Sprintf("%q (ssl: %t)", tc.hostPort, tc.ssl),
			t,
			func(g *assert.G) {
				gotHostPort, err := checkHostPort(tc.hostPort, tc.ssl)
				assert.Equal(g, gotHostPort, tc.expectedHostPort)
				if tc.expectedErrorContains != nil {
					assert.ErrorContains(
						g,
						err,
						*tc.expectedErrorContains,
					)
				} else {
					assert.Nil(g, err)
				}
			},
		)
	}
}
