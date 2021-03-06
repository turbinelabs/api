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
	"bytes"
	"compress/gzip"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/turbinelabs/api"
	apihttp "github.com/turbinelabs/api/http"
	"github.com/turbinelabs/api/http/envelope"
	httperr "github.com/turbinelabs/api/http/error"
	apiheader "github.com/turbinelabs/api/http/header"
	statsapi "github.com/turbinelabs/api/service/stats"
	"github.com/turbinelabs/test/assert"
)

func TestStatsClientQueryV2Success(t *testing.T) {
	wantQueryStr := `{"time_range":{"granularity":"minutes"},"timeseries":null}`

	want := &statsapi.QueryResult{}

	verifier := verifyingHandler{
		fn: func(rr apihttp.RichRequest) {
			assert.Equal(t, rr.Underlying().URL.Path, v2QueryPath)
			assert.Nil(t, rr.Underlying().URL.Query()["query"])
			body, err := rr.GetBody()
			assert.Nil(t, err)

			gzipReader, err := gzip.NewReader(bytes.NewReader(body))
			assert.Nil(t, err)

			bodyBytes, err := ioutil.ReadAll(gzipReader)
			defer gzipReader.Close()
			assert.Nil(t, err)

			assert.Equal(t, strings.TrimSpace(string(bodyBytes)), wantQueryStr)
		},
		status:   http.StatusOK,
		response: want,
	}

	server := httptest.NewServer(verifier)
	defer server.Close()

	endpoint := newTestEndpointFromServer(server)
	client, _ := NewStatsV2Client(endpoint, clientTestAPIKey, clientTestApp, nil)

	got, gotErr := client.QueryV2(&statsapi.Query{})
	assert.Nil(t, gotErr)
	assert.DeepEqual(t, got, want)
}

func TestStatsClientQueryV2Error(t *testing.T) {
	wantQueryStr := `{"time_range":{"granularity":"minutes"},"timeseries":null}`

	wantErr := httperr.New500("Gah!", httperr.UnknownUnclassifiedCode)

	verifier := verifyingHandler{
		fn: func(rr apihttp.RichRequest) {
			assert.Equal(t, rr.Underlying().URL.Path, v2QueryPath)
			assert.Nil(t, rr.Underlying().URL.Query()["query"])
			body, err := rr.GetBody()
			assert.Nil(t, err)

			gzipReader, err := gzip.NewReader(bytes.NewReader(body))
			assert.Nil(t, err)

			bodyBytes, err := ioutil.ReadAll(gzipReader)
			defer gzipReader.Close()
			assert.Nil(t, err)

			assert.Equal(t, strings.TrimSpace(string(bodyBytes)), wantQueryStr)
		},
		status:   http.StatusInternalServerError,
		response: envelope.NewErrorResponse(wantErr, nil),
	}

	server := httptest.NewServer(verifier)
	defer server.Close()

	endpoint := newTestEndpointFromServer(server)
	client, _ := NewStatsV2Client(endpoint, clientTestAPIKey, clientTestApp, nil)

	got, gotErr := client.QueryV2(&statsapi.Query{})
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
		v2ForwardPath,
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
