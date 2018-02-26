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

package client

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/turbinelabs/api"
	"github.com/turbinelabs/api/fixture"
	apihttp "github.com/turbinelabs/api/http"
	"github.com/turbinelabs/api/http/envelope"
	apiheader "github.com/turbinelabs/api/http/header"
	"github.com/turbinelabs/api/service"
	"github.com/turbinelabs/test/assert"
)

const (
	clientTestAPIKey     = "whee-whee-whee"
	clientTestApp        = App("app")
	clusterCommonURL     = "/v1.0/cluster"
	domainCommonURL      = "/v1.0/domain"
	proxyCommonURL       = "/v1.0/proxy"
	routeCommonURL       = "/v1.0/route"
	sharedRulesCommonURL = "/v1.0/shared_rules"
	zoneCommonURL        = "/v1.0/zone"
	userCommonURL        = "/v1.0/admin/user"
)

var fixtures = fixture.New()

// Used for verifying http client tests. It does clever things to decide how
// to write out a response:
//
//   If response is a X the verifier handler writes Y:
//    string ------------- exactly those bytes
//    envelope.Response -- the marshaled version of that object
//    *envelope.Response - the marshaled version of that object
//    Something else ----- an envelope.Response with the response parameter as the "response" field of the envelope
type verifyingHandler struct {
	fn       func(apihttp.RichRequest)
	status   int
	response interface{}
}

func (w verifyingHandler) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	rr := apihttp.NewRichRequest(r)
	rrw := apihttp.RichResponseWriter{ResponseWriter: rw}

	apiKey := rr.Underlying().Header.Get(apiheader.Authorization)
	if apiKey != clientTestAPIKey {
		rw.WriteHeader(400)
		rw.Write([]byte(
			fmt.Sprintf(
				"wrong %s header, got %s, want %s",
				apiheader.Authorization,
				apiKey,
				clientTestAPIKey,
			),
		))
		return
	}

	cType := rr.Underlying().Header.Get(apiheader.ClientType)
	if cType != clientType {
		rw.WriteHeader(400)
		rw.Write([]byte(
			fmt.Sprintf(
				"wrong %s header: got %s, want %s",
				apiheader.ClientType,
				cType,
				clientType,
			),
		))
		return
	}

	clientVersion := rr.Underlying().Header.Get(apiheader.ClientVersion)
	if clientVersion != api.TbnPublicVersion {
		rw.WriteHeader(400)
		rw.Write([]byte(
			fmt.Sprintf(
				"wrong %s header: got %s, want %s",
				apiheader.ClientVersion,
				clientVersion,
				api.TbnPublicVersion,
			),
		))
		return
	}

	cApp := App(rr.Underlying().Header.Get(apiheader.ClientApp))
	if cApp != clientTestApp {
		rw.WriteHeader(400)
		rw.Write([]byte(
			fmt.Sprintf(
				"wrong %s header: got %s, want %s",
				apiheader.ClientApp,
				cApp,
				clientTestApp,
			),
		))
		return
	}

	w.fn(rr)

	if w.response != nil {
		switch t := w.response.(type) {
		case string:
			rw.WriteHeader(w.status)
			rw.Write([]byte(t))
		case envelope.Response:
			rrw.WriteEnvelope(t.Error, t.Payload)
		case *envelope.Response:
			rrw.WriteEnvelope(t.Error, t.Payload)
		default:
			rrw.WriteEnvelope(nil, w.response)
		}
	}
}

func stripURLPrefix(url, prefix string) string {
	return url[len(prefix):]
}

func newTestEndpoint(hostPort string) apihttp.Endpoint {
	e, err := apihttp.NewEndpoint(apihttp.HTTP, hostPort)
	if err != nil {
		log.Fatal(err)
	}
	return e
}

func newTestEndpointFromServer(server *httptest.Server) apihttp.Endpoint {
	u, e := url.Parse(server.URL)
	if e != nil {
		log.Fatal(e)
	}

	return newTestEndpoint(u.Host)
}

func getAllInterface(server *httptest.Server) service.All {
	endpoint := newTestEndpointFromServer(server)
	serviceall, err := NewAll(endpoint, clientTestAPIKey, clientTestApp)
	if err != nil {
		log.Fatal(err)
	}
	return serviceall
}

func getAdminInterface(server *httptest.Server) service.Admin {
	endpoint := newTestEndpointFromServer(server)
	admin, err := NewAdmin(endpoint, clientTestAPIKey, clientTestApp)
	if err != nil {
		log.Fatal(err)
	}
	return admin
}

func TestNewAllCopiesEndpoint(t *testing.T) {
	e := newTestEndpoint("example.com:80")

	r, err := e.NewRequest("GET", "/index.html", apihttp.Params{}, nil)
	assert.Nil(t, err)
	assert.Equal(t, len(r.Header), 0)

	all, err := NewAll(endpoint, clientTestAPIKey, clientTestApp)
	assert.Nil(t, err)
	assert.NonNil(t, all)

	allEndpoint := all.(*httpServiceV1).clusterV1.dest

	r, err = e.NewRequest("GET", "/index.html", apihttp.Params{}, nil)
	assert.Nil(t, err)
	assert.Equal(t, len(r.Header), 0)

	r, err = allEndpoint.NewRequest("GET", "/index.html", apihttp.Params{}, nil)
	assert.Nil(t, err)
	assert.Equal(t, len(r.Header), 4)
	assert.ArrayEqual(t, r.Header[apiheader.Authorization], []string{clientTestAPIKey})
	assert.ArrayEqual(t, r.Header[apiheader.ClientType], []string{clientType})
	assert.ArrayEqual(t, r.Header[apiheader.ClientVersion], []string{api.TbnPublicVersion})
	assert.ArrayEqual(t, r.Header[apiheader.ClientApp], []string{string(clientTestApp)})
}

func TestNewAdminCopiesEndpoint(t *testing.T) {
	e := newTestEndpoint("example.com:80")

	r, err := e.NewRequest("GET", "/index.html", apihttp.Params{}, nil)
	assert.Nil(t, err)
	assert.Equal(t, len(r.Header), 0)

	all, err := NewAdmin(endpoint, clientTestAPIKey, clientTestApp)
	assert.Nil(t, err)
	assert.NonNil(t, all)

	adminEndpoint := all.(*httpAdminV1).userV1.dest

	r, err = e.NewRequest("GET", "/index.html", apihttp.Params{}, nil)
	assert.Nil(t, err)
	assert.Equal(t, len(r.Header), 0)

	r, err = adminEndpoint.NewRequest("GET", "/index.html", apihttp.Params{}, nil)
	assert.Nil(t, err)
	assert.Equal(t, len(r.Header), 4)
	assert.ArrayEqual(t, r.Header[apiheader.Authorization], []string{clientTestAPIKey})
	assert.ArrayEqual(t, r.Header[apiheader.ClientType], []string{clientType})
	assert.ArrayEqual(t, r.Header[apiheader.ClientVersion], []string{api.TbnPublicVersion})
	assert.ArrayEqual(t, r.Header[apiheader.ClientApp], []string{string(clientTestApp)})
}
