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
