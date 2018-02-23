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
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/turbinelabs/test/assert"
)

type redirectingHandler struct {
	numRedirectsRemaining int
	numRequests           int
	numRedirected         int
	headers               []http.Header
}

func (r *redirectingHandler) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	headerCopy := http.Header{}
	for k, v := range req.Header {
		headerCopy[k] = v
	}
	r.headers = append(r.headers, headerCopy)

	r.numRequests++
	if r.numRedirectsRemaining > 0 {
		r.numRedirectsRemaining--
		r.numRedirected++
		http.Redirect(
			rw,
			req,
			fmt.Sprintf("/redirect/%d", r.numRedirected),
			http.StatusMovedPermanently,
		)
	} else {
		fmt.Fprintln(rw, "OK")
	}
}

func TestHeaderPreservingClientNoRedirects(t *testing.T) {
	handler := &redirectingHandler{}
	server := httptest.NewServer(handler)
	defer server.Close()

	resp, err := HeaderPreservingClient().Get(server.URL)
	assert.Nil(t, err)
	result, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	assert.Nil(t, err)

	assert.Equal(t, resp.StatusCode, 200)
	assert.Equal(t, string(result), "OK\n")
	assert.Equal(t, handler.numRequests, 1)
	assert.Equal(t, handler.numRedirected, 0)
}

func TestHeaderPreservingClientSomeRedirects(t *testing.T) {
	handler := &redirectingHandler{numRedirectsRemaining: 5}
	server := httptest.NewServer(handler)
	defer server.Close()

	resp, err := HeaderPreservingClient().Get(server.URL)
	assert.Nil(t, err)
	result, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	assert.Nil(t, err)

	assert.Equal(t, resp.StatusCode, 200)
	assert.Equal(t, string(result), "OK\n")
	assert.Equal(t, handler.numRequests, 6)
	assert.Equal(t, handler.numRedirected, 5)
}

func TestHeaderPreservingClientTooManyRedirects(t *testing.T) {
	handler := &redirectingHandler{numRedirectsRemaining: 10}
	server := httptest.NewServer(handler)
	defer server.Close()

	resp, err := HeaderPreservingClient().Get(server.URL)
	assert.NonNil(t, err)
	assert.NotEqual(t, resp.StatusCode, 200)
	assert.Equal(t, handler.numRequests, 6)
	assert.Equal(t, handler.numRedirected, 6)
}

func TestHeaderPreservingClientPreservesHeaders(t *testing.T) {
	handler := &redirectingHandler{numRedirectsRemaining: 5}
	server := httptest.NewServer(handler)
	defer server.Close()

	resp, err := HeaderPreservingClient().Get(server.URL)
	assert.Nil(t, err)
	ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()

	assert.Equal(t, handler.numRequests, 6)
	assert.Equal(t, handler.numRedirected, 5)

	assert.Equal(t, len(handler.headers), 6)
	for i, headers := range handler.headers {
		if i > 0 {
			assert.DeepEqual(t, headers, handler.headers[0])
		}
	}
}
