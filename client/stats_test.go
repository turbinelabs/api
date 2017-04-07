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
	"bytes"
	"compress/gzip"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"math"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/golang/mock/gomock"

	"github.com/turbinelabs/api"
	apihttp "github.com/turbinelabs/api/http"
	"github.com/turbinelabs/api/http/envelope"
	httperr "github.com/turbinelabs/api/http/error"
	apiheader "github.com/turbinelabs/api/http/header"
	statsapi "github.com/turbinelabs/api/service/stats"
	"github.com/turbinelabs/nonstdlib/executor"
	"github.com/turbinelabs/nonstdlib/ptr"
	tbntime "github.com/turbinelabs/nonstdlib/time"
	"github.com/turbinelabs/test/assert"
)

const metricName1 = "group.metric"

var (
	when1Micros = tbntime.ToUnixMicro(time.Now())

	payload = &statsapi.Payload{
		Source: sourceString1,
		Stats: []statsapi.Stat{
			{
				Name:      metricName1,
				Value:     ptr.Float64(1.41421),
				Timestamp: when1Micros,
				Tags:      map[string]string{"tag": "tag-value"},
			},
		},
	}

	badPayload = &statsapi.Payload{
		Source: sourceString1,
		Stats: []statsapi.Stat{
			{
				Name:      metricName1,
				Value:     ptr.Float64(math.Inf(1)),
				Timestamp: when1Micros,
				Tags:      map[string]string{},
			},
		},
	}

	endpoint, _ = apihttp.NewEndpoint(apihttp.HTTP, "example.com", 8080)
)

func TestEncodePayload(t *testing.T) {
	expectedJson :=
		fmt.Sprintf(
			`{"source":"%s","stats":[{"name":"%s","value":%g,"timestamp":%d,"tags":{"%s":"%s"}}]}`+"\n",
			sourceString1,
			metricName1,
			1.41421,
			when1Micros,
			"tag",
			"tag-value",
		)
	var expectedBytes bytes.Buffer
	gw := gzip.NewWriter(&expectedBytes)
	gw.Write([]byte(expectedJson))
	gw.Close()

	json, err := encodePayload(payload)
	assert.Nil(t, err)
	assert.DeepEqual(t, json, expectedBytes.Bytes())
}

func TestEncodePayloadError(t *testing.T) {
	json, err := encodePayload(badPayload)
	assert.Nil(t, json)
	assert.NonNil(t, err)
}

type forwardResult struct {
	result *statsapi.ForwardResult
	err    error
}

type resultFunc func() (*statsapi.ForwardResult, error)
type requestFunc func(statsapi.StatsService) (*statsapi.ForwardResult, error)
type newStatsFunc func(
	apihttp.Endpoint,
	string,
	executor.Executor,
) (statsapi.StatsService, error)

func prepareStatsClientTest(
	t *testing.T,
	e apihttp.Endpoint,
	reqFunc requestFunc,
) (executor.Func, executor.CallbackFunc, resultFunc) {
	ctrl := gomock.NewController(assert.Tracing(t))

	funcChan := make(chan executor.Func, 1)
	callbackFuncChan := make(chan executor.CallbackFunc, 1)

	mockExec := executor.NewMockExecutor(ctrl)
	mockExec.EXPECT().
		Exec(gomock.Any(), gomock.Any()).
		Do(
			func(f executor.Func, cb executor.CallbackFunc) {
				funcChan <- f
				callbackFuncChan <- cb
			},
		)

	client, err := NewStatsClient(e, clientTestAPIKey, clientTestApp, mockExec)
	assert.Nil(t, err)

	rvChan := make(chan forwardResult, 1)

	go func() {
		r, err := reqFunc(client)
		rvChan <- forwardResult{r, err}
	}()

	f := <-funcChan
	cb := <-callbackFuncChan

	return f, cb, func() (*statsapi.ForwardResult, error) {
		defer ctrl.Finish()
		rv := <-rvChan
		return rv.result, rv.err
	}
}

func payloadForward(p *statsapi.Payload) func(client statsapi.StatsService) (*statsapi.ForwardResult, error) {
	return func(client statsapi.StatsService) (*statsapi.ForwardResult, error) {
		return client.Forward(p)
	}
}

var simpleForward = payloadForward(payload)

func TestStatsClientForward(t *testing.T) {
	_, cb, getResult := prepareStatsClientTest(t, endpoint, simpleForward)

	expectedResult := &statsapi.ForwardResult{NumAccepted: 1}
	cb(executor.NewReturn(expectedResult))

	result, err := getResult()
	assert.SameInstance(t, result, expectedResult)
	assert.Nil(t, err)
}

func TestStatsClientForwardFailure(t *testing.T) {
	_, cb, getResult := prepareStatsClientTest(t, endpoint, simpleForward)

	expectedErr := errors.New("failure")
	cb(executor.NewError(expectedErr))

	result, err := getResult()
	assert.Nil(t, result)
	assert.SameInstance(t, err, expectedErr)
}

type testHandler struct {
	t               *testing.T
	requestPayload  *statsapi.Payload
	responsePayload *statsapi.ForwardResult
	responseError   *httperr.Error
}

func (h *testHandler) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	handler := verifyingHandler{
		fn: func(rr apihttp.RichRequest) {
			body := rr.Underlying().Body
			assert.NonNil(h.t, body)

			gzipReader, err := gzip.NewReader(body)
			assert.Nil(h.t, err)

			bytes, err := ioutil.ReadAll(gzipReader)
			defer gzipReader.Close()
			assert.Nil(h.t, err)

			stats := &statsapi.Payload{}
			err = json.Unmarshal(bytes, stats)
			assert.Nil(h.t, err)
			h.requestPayload = stats
		},
		status:   200,
		response: &envelope.Response{Error: h.responseError, Payload: h.responsePayload},
	}

	handler.ServeHTTP(resp, req)
}

func runStatsClientFuncTest(
	t *testing.T,
	requestPayload *statsapi.Payload,
	responsePayload *statsapi.ForwardResult,
	httpErr *httperr.Error,
) (*statsapi.Payload, *statsapi.ForwardResult, error) {
	handler := &testHandler{responsePayload: responsePayload, responseError: httpErr}
	server := httptest.NewServer(handler)
	defer server.Close()

	host, portStr, _ := net.SplitHostPort(server.Listener.Addr().String())
	port, _ := net.LookupPort(server.Listener.Addr().Network(), portStr)

	endpoint, err := apihttp.NewEndpoint(apihttp.HTTP, host, port)
	assert.Nil(t, err)

	f, cb, _ := prepareStatsClientTest(t, endpoint, payloadForward(requestPayload))
	cb(executor.NewError(errors.New("don't care about this")))

	ctxt := context.Background()

	response, err := f(ctxt)
	if response == nil {
		return handler.requestPayload, nil, err
	}
	return handler.requestPayload, response.(*statsapi.ForwardResult), err
}

func TestStatsClientForwardExecFunc(t *testing.T) {
	expectedResult := &statsapi.ForwardResult{NumAccepted: 12}

	gotPayload, result, err := runStatsClientFuncTest(t, payload, expectedResult, nil)
	assert.DeepEqual(t, gotPayload, payload)
	assert.DeepEqual(t, result, expectedResult)
	assert.Nil(t, err)
}

func TestStatsClientForwardExecFuncFailure(t *testing.T) {
	expectedErr := httperr.AuthorizationError()

	gotPayload, result, err := runStatsClientFuncTest(t, payload, nil, expectedErr)
	assert.DeepEqual(t, gotPayload, payload)
	assert.Nil(t, result)
	assert.Equal(t, err.Error(), expectedErr.Error())
}

func TestStatsClientQuerySuccess(t *testing.T) {
	wantQueryStr := `{"zone_name":"","time_range":{"granularity":"seconds"},"timeseries":null}`

	want := &statsapi.QueryResult{}

	verifier := verifyingHandler{
		fn: func(rr apihttp.RichRequest) {
			assert.Equal(t, rr.Underlying().URL.Path, queryPath)
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
			assert.Equal(t, rr.Underlying().URL.Path, queryPath)
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
	endpoint := newTestEndpoint("example.com", 80)

	r, err := endpoint.NewRequest("GET", "/index.html", apihttp.Params{}, nil)
	assert.Nil(t, err)
	assert.Equal(t, len(r.Header), 0)

	client, err := newInternalStatsClient(endpoint, clientTestAPIKey, clientTestApp, nil)
	assert.Nil(t, err)
	assert.NonNil(t, client)

	statsEndpoint := client.(*httpStatsV1).dest

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
