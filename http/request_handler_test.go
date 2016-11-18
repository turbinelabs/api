package http

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/turbinelabs/api/http/envelope"
	httperr "github.com/turbinelabs/api/http/error"
	"github.com/turbinelabs/api/http/header"
	"github.com/turbinelabs/test/assert"
	testio "github.com/turbinelabs/test/io"
)

const (
	testAPIKey   = "some-api-key"
	testClientID = "some-client-id"
)

type testPayload struct {
	N int `json:"n"`
}

type testMalformedPayload struct {
	N testPayload `json:"n"`
}

func TestNewRequestHandler(t *testing.T) {
	rh := NewRequestHandler(http.DefaultClient, testAPIKey, testClientID)
	assert.SameInstance(t, rh.client, http.DefaultClient)
	assert.Equal(t, rh.apiKey, testAPIKey)
	assert.Equal(t, rh.clientID, testClientID)
}

func TestRequestHandlerGetBody(t *testing.T) {
	const theBody = "some body content"
	resp := &http.Response{
		Body: ioutil.NopCloser(strings.NewReader(theBody)),
	}

	body, err := getBody(resp)
	assert.Nil(t, err)
	assert.DeepEqual(t, body, []byte(theBody))
}

func TestRequestHandlerGetBodyError(t *testing.T) {
	resp := &http.Response{
		Body: testio.NewFailingReader(),
	}

	body, err := getBody(resp)
	assert.NonNil(t, err)
	assert.Nil(t, body)
}

func mkBody(t *testing.T, err *httperr.Error, payload interface{}) string {
	env := envelope.Response{
		Error:   err,
		Payload: payload,
	}

	encoded, encodingErr := json.Marshal(&env)
	assert.Nil(t, encodingErr)
	return string(encoded)
}

func mkBodyReader(t *testing.T, err *httperr.Error, payload interface{}) io.ReadCloser {
	encoded := mkBody(t, err, payload)
	return ioutil.NopCloser(strings.NewReader(string(encoded)))
}

func TestExpectsNoPayload(t *testing.T) {
	resp := &http.Response{
		StatusCode: http.StatusOK,
	}
	assert.Nil(t, expectsNoPayload(resp))

	resp.StatusCode = http.StatusInternalServerError
	assert.ErrorContains(
		t,
		expectsNoPayload(resp),
		"error response with no additional information",
	)

	resp.Body = testio.NewFailingReader()
	assert.ErrorContains(t, expectsNoPayload(resp), testio.FailingReaderMessage)

	resp.Body = ioutil.NopCloser(strings.NewReader("not json"))
	assert.ErrorContains(t, expectsNoPayload(resp), "malformed response")

	httpErr := httperr.New400("reasons", httperr.MiscErrorCode)
	resp.Body = mkBodyReader(t, httpErr, nil)

	// overwrites envelope error status with HTTP response code
	assert.DeepEqual(
		t,
		expectsNoPayload(resp),
		httperr.New500("reasons", httperr.MiscErrorCode),
	)
}

func TestExpectsPayload(t *testing.T) {
	payloadDest := &testPayload{}

	resp := &http.Response{
		StatusCode: http.StatusOK,
	}
	assert.ErrorContains(
		t,
		expectsPayload(resp, payloadDest),
		"expected payload but response (200) included no content",
	)
	assert.DeepEqual(t, payloadDest, &testPayload{})

	resp.Body = testio.NewFailingReader()
	assert.ErrorContains(t, expectsPayload(resp, payloadDest), testio.FailingReaderMessage)
	assert.DeepEqual(t, payloadDest, &testPayload{})

	resp.Body = ioutil.NopCloser(strings.NewReader("not json"))
	assert.ErrorContains(t, expectsPayload(resp, payloadDest), "malformed response")
	assert.DeepEqual(t, payloadDest, &testPayload{})

	resp.Body = mkBodyReader(t, nil, &testMalformedPayload{N: testPayload{10}})
	assert.ErrorContains(t, expectsPayload(resp, payloadDest), "malformed response")
	assert.DeepEqual(t, payloadDest, &testPayload{})

	httpErr := httperr.New400("reasons", httperr.MiscErrorCode)
	resp.StatusCode = http.StatusInternalServerError
	resp.Body = mkBodyReader(t, httpErr, nil)
	// overwrites envelope error status with HTTP response code
	assert.DeepEqual(
		t,
		expectsPayload(resp, payloadDest),
		httperr.New500("reasons", httperr.MiscErrorCode),
	)
	assert.DeepEqual(t, payloadDest, &testPayload{})

	expectedPayload := &testPayload{N: 100}
	resp.Body = mkBodyReader(t, httpErr, expectedPayload)
	resp.StatusCode = http.StatusInternalServerError
	// overwrites envelope error status with HTTP response code, does not unmarshal payload
	assert.DeepEqual(
		t,
		expectsPayload(resp, payloadDest),
		httperr.New500("reasons", httperr.MiscErrorCode),
	)
	assert.DeepEqual(t, payloadDest, &testPayload{})

	resp.Body = mkBodyReader(t, nil, expectedPayload)
	resp.StatusCode = http.StatusOK
	assert.Nil(t, expectsPayload(resp, payloadDest))
	assert.DeepEqual(t, payloadDest, expectedPayload)
}

func TestDoMkReqError(t *testing.T) {
	rh := NewRequestHandler(http.DefaultClient, testAPIKey, testClientID)
	assert.ErrorContains(
		t,
		rh.Do(
			func() (*http.Request, error) {
				return nil, errors.New("boom")
			},
			nil,
		),
		"could not create request: boom",
	)
}

func TestDoRequestFailure(t *testing.T) {
	rh := NewRequestHandler(http.DefaultClient, testAPIKey, testClientID)

	mkReq := func() (*http.Request, error) {
		return http.NewRequest("GET", "http://127.0.0.1:1/nope", nil)
	}

	err := rh.Do(mkReq, nil)
	assert.Equal(t, err.Status, 400)
	assert.ErrorContains(t, err, "could not successfully make request")
}

func TestDoWithoutPayload(t *testing.T) {
	server := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, r.Header.Get(header.APIKey), testAPIKey)
			assert.Equal(t, r.Header.Get(header.ClientID), testClientID)
			w.WriteHeader(http.StatusOK)
			fmt.Fprintln(w, "OK")
		}),
	)
	defer server.Close()

	rh := NewRequestHandler(http.DefaultClient, testAPIKey, testClientID)

	mkReq := func() (*http.Request, error) {
		return http.NewRequest("GET", server.URL, nil)
	}

	assert.Nil(t, rh.Do(mkReq, nil))
}

func TestDoWithPayload(t *testing.T) {
	expectedPayload := &testPayload{N: 99}

	server := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, r.Header.Get(header.APIKey), testAPIKey)
			assert.Equal(t, r.Header.Get(header.ClientID), testClientID)
			w.WriteHeader(http.StatusOK)
			fmt.Fprintln(w, mkBody(t, nil, expectedPayload))
		}),
	)
	defer server.Close()

	rh := NewRequestHandler(http.DefaultClient, testAPIKey, testClientID)

	mkReq := func() (*http.Request, error) {
		return http.NewRequest("GET", server.URL, nil)
	}

	payload := &testPayload{}

	err := rh.Do(mkReq, payload)
	assert.Nil(t, err)
	assert.DeepEqual(t, payload, expectedPayload)
}