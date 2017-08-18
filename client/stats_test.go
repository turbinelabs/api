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

package client

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/turbinelabs/api"
	apihttp "github.com/turbinelabs/api/http"
	"github.com/turbinelabs/api/http/envelope"
	httperr "github.com/turbinelabs/api/http/error"
	apiheader "github.com/turbinelabs/api/http/header"
	statsapi "github.com/turbinelabs/api/service/stats"
	"github.com/turbinelabs/test/assert"
)

func TestStatsClientQuerySuccess(t *testing.T) {
	wantQueryStr := `{"zone_name":"","time_range":{"granularity":"seconds"},"timeseries":null}`

	want := &statsapi.QueryResult{}

	verifier := verifyingHandler{
		fn: func(rr apihttp.RichRequest) {
			assert.Equal(t, rr.Underlying().URL.Path, v1QueryPath)
			assert.NonNil(t, rr.Underlying().URL.Query()["query"])
			assert.Equal(t, len(rr.Underlying().URL.Query()["query"]), 1)
			assert.Equal(t, rr.Underlying().URL.Query()["query"][0], wantQueryStr)
		},
		status:   http.StatusOK,
		response: want,
	}

	server := httptest.NewServer(verifier)
	defer server.Close()

	endpoint := newTestEndpointFromServer(server)
	client, _ := NewStatsClient(endpoint, clientTestAPIKey, clientTestApp, nil)

	got, gotErr := client.Query(&statsapi.Query{})
	assert.Nil(t, gotErr)

	assert.DeepEqual(t, got, want)
}

func TestStatsClientQueryError(t *testing.T) {
	wantQueryStr := `{"zone_name":"","time_range":{"granularity":"seconds"},"timeseries":null}`

	wantErr := httperr.New500("Gah!", httperr.UnknownUnclassifiedCode)

	verifier := verifyingHandler{
		fn: func(rr apihttp.RichRequest) {
			assert.Equal(t, rr.Underlying().URL.Path, v1QueryPath)
			assert.NonNil(t, rr.Underlying().URL.Query()["query"])
			assert.Equal(t, len(rr.Underlying().URL.Query()["query"]), 1)
			assert.Equal(t, rr.Underlying().URL.Query()["query"][0], wantQueryStr)
		},
		status:   http.StatusInternalServerError,
		response: envelope.Response{wantErr, nil},
	}

	server := httptest.NewServer(verifier)
	defer server.Close()

	endpoint := newTestEndpointFromServer(server)
	client, _ := NewStatsClient(endpoint, clientTestAPIKey, clientTestApp, nil)

	got, gotErr := client.Query(&statsapi.Query{})
	assert.DeepEqual(t, gotErr, wantErr)
	assert.Nil(t, got)
}

func TestNewInternalStatsClientCopiesEndpoint(t *testing.T) {
	endpoint := newTestEndpoint("example.com:80")

	r, err := endpoint.NewRequest("GET", "/index.html", apihttp.Params{}, nil)
	assert.Nil(t, err)
	assert.Equal(t, len(r.Header), 0)

	client, err := newInternalStatsClient(
		endpoint,
		v1ForwardPath,
		clientTestAPIKey,
		clientTestApp,
		nil,
	)
	assert.Nil(t, err)
	assert.NonNil(t, client)

	statsEndpoint := client.dest

	r, err = endpoint.NewRequest("GET", "/index.html", apihttp.Params{}, nil)
	assert.Nil(t, err)
	assert.Equal(t, len(r.Header), 0)

	r, err = statsEndpoint.NewRequest("GET", "/index.html", apihttp.Params{}, nil)
	assert.Nil(t, err)
	assert.Equal(t, len(r.Header), 6)
	assert.ArrayEqual(t, r.Header[apiheader.Authorization], []string{clientTestAPIKey})
	assert.ArrayEqual(t, r.Header[apiheader.ClientType], []string{clientType})
	assert.ArrayEqual(t, r.Header[apiheader.ClientVersion], []string{api.TbnPublicVersion})
	assert.ArrayEqual(t, r.Header[apiheader.ClientApp], []string{string(clientTestApp)})
	assert.ArrayEqual(t, r.Header["Content-Type"], []string{"application/json"})
	assert.ArrayEqual(t, r.Header["Content-Encoding"], []string{"gzip"})
}
