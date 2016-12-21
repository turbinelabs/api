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
	"flag"
	"net/http"
	"testing"

	tbnflag "github.com/turbinelabs/nonstdlib/flag"
	"github.com/turbinelabs/test/assert"
)

func TestNewFromFlags(t *testing.T) {
	flagset := tbnflag.NewPrefixedFlagSet(
		flag.NewFlagSet("api/http options", flag.PanicOnError),
		"api",
		"API",
	)

	ff := NewFromFlags("api.turbinelabs.io", flagset)
	ffImpl := ff.(*fromFlags)

	assert.Equal(t, ffImpl.host, "api.turbinelabs.io")

	flagset.Parse([]string{
		"-api.host=example.com",
		"-api.port=999",
		"-api.ssl=false",
		"-api.insecure=true",
		"-api.header=fred: flintstone",
		"-api.header=barney: rubble",
	})

	assert.Equal(t, ffImpl.host, "example.com")
	assert.Equal(t, ffImpl.port, 999)
	assert.False(t, ffImpl.ssl)
	assert.True(t, ffImpl.insecure)
	assert.ArrayEqual(t, ffImpl.headers.Strings, []string{"fred: flintstone", "barney: rubble"})
}

func TestFromFlagsValidate(t *testing.T) {
	ff := &fromFlags{
		port:    443,
		headers: tbnflag.NewStrings(),
	}

	assert.Nil(t, ff.Validate())

	ff.port = 0
	assert.ErrorContains(t, ff.Validate(), "invalid API port")

	ff.port = 65536
	assert.ErrorContains(t, ff.Validate(), "invalid API port")

	ff.port = 443
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
		host:    "example.com",
		port:    80,
		headers: tbnflag.NewStrings(),
	}

	ff.headers.Set("x-fred: flintstone,x-barney: rubble")

	e, err := ff.MakeEndpoint()
	assert.Nil(t, err)
	assert.Equal(t, e.urlBase.String(), "http://example.com:80")
	assert.NonNil(t, e.client)
	assert.Equal(t, e.header.Get("X-Fred"), "flintstone")
	assert.Equal(t, e.header.Get("X-Barney"), "rubble")

	ff.port = 443
	ff.ssl = true

	e, err = ff.MakeEndpoint()
	assert.Nil(t, err)
	assert.Equal(t, e.urlBase.String(), "https://example.com:443")
	assert.NonNil(t, e.client)

	ff.host = "not a domain"

	_, err = ff.MakeEndpoint()
	assert.NonNil(t, err)
}
