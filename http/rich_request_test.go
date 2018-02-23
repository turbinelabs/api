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
	"bytes"
	"io/ioutil"
	"net/http"
	"net/url"
	"testing"

	httperr "github.com/turbinelabs/api/http/error"
	"github.com/turbinelabs/test/assert"
)

func richRequestVerify(
	t *testing.T,
	gotValue string,
	gotBool bool,
	wantValue string,
	wantBool bool,
) {
	assert.Equal(t, gotValue, wantValue)
	assert.Equal(t, gotBool, wantBool)
}

func richRequestVerifyGotBody(
	t *testing.T,
	gotValue []byte,
	gotErr error,
	wantValue []byte,
	wantErr error,
) {
	assert.DeepEqual(t, gotValue, wantValue)
	assert.DeepEqual(t, gotErr, wantErr)
}

func mkRR(t *testing.T, urlstr string) richRequest {
	req := &http.Request{}
	newurl, err := url.Parse("http://foo.com" + urlstr)
	if err != nil {
		t.Fatalf("Failure to construct test object: %v", err)
	}

	req.URL = newurl
	return richRequest{req}
}

func TestQueryArgOkNotFound(t *testing.T) {
	rr := mkRR(t, "?arg=asonetuh")
	v, b := rr.QueryArgOk("nope")
	richRequestVerify(t, v, b, "", false)
}

func TestQueryArgOkHit(t *testing.T) {
	rr := mkRR(t, "?arg=test")
	v, b := rr.QueryArgOk("arg")
	richRequestVerify(t, v, b, "test", true)
}

func TestQueryArgOkHitMultiple(t *testing.T) {
	rr := mkRR(t, "?arg=test&arg=test2")
	v, b := rr.QueryArgOk("arg")
	richRequestVerify(t, v, b, "test", true)
}

func TestQueryArgHit(t *testing.T) {
	rr := mkRR(t, "?arg=test")
	v := rr.QueryArg("arg")
	richRequestVerify(t, v, true, "test", true)
}

func TestQueryArgNotFound(t *testing.T) {
	rr := mkRR(t, "?argaoeu=test")
	v := rr.QueryArg("arg")
	richRequestVerify(t, v, true, "", true)
}

func TestQueryArgOr(t *testing.T) {
	rr := mkRR(t, "")
	v := rr.QueryArgOr("arg", "default")
	richRequestVerify(t, v, true, "default", true)
}

func TestQueryArgOrHit(t *testing.T) {
	rr := mkRR(t, "?arg=foo&arg=bar")
	v := rr.QueryArgOr("arg", "default")
	richRequestVerify(t, v, true, "foo", true)
}

func TestQueryArgEscaped(t *testing.T) {
	input := "{\"'&/?<#"
	rr := mkRR(t, "?arg="+url.QueryEscape(input))
	v := rr.QueryArg("arg")
	richRequestVerify(t, v, true, input, true)
}

func TestGetBody(t *testing.T) {
	rr := mkRR(t, "")
	want := []byte("this is some input\n\nover multiple lines")
	rr.Body = ioutil.NopCloser(bytes.NewBuffer(want))
	got, e := rr.GetBody()
	richRequestVerifyGotBody(t, got, e, want, nil)
}

func TestGetBodyNil(t *testing.T) {
	rr := mkRR(t, "")
	wantErr := httperr.New400("no body available", httperr.UnknownNoBodyCode)
	gotB, gotErr := rr.GetBody()
	richRequestVerifyGotBody(t, gotB, gotErr, nil, wantErr)
}

type obj struct {
	Foo string `json:"foo"`
	Bar int    `json:"bar"`
}

func TestGetBodyObject(t *testing.T) {
	jsonSrc := `{"foo":"whee", "bar":1234}`
	rr := mkRR(t, "")
	rr.Body = ioutil.NopCloser(bytes.NewBuffer([]byte(jsonSrc)))
	want := obj{"whee", 1234}
	got := obj{}
	err := rr.GetBodyObject(&got)
	assert.Nil(t, err)
	assert.DeepEqual(t, got, want)
}

func TestGetBodyObjectBrokenJSON(t *testing.T) {
	jsonSrc := "nope"
	rr := mkRR(t, "")
	rr.Body = ioutil.NopCloser(bytes.NewBuffer([]byte(jsonSrc)))
	want := obj{}
	got := obj{}
	err := rr.GetBodyObject(&got)

	assert.DeepEqual(t, got, want)
	wantErr := httperr.NewDetailed400(
		"error handling JSON content",
		httperr.UnknownDecodingCode,
		map[string]string{
			"error":   "invalid character 'o' in literal null (expecting 'u')",
			"content": jsonSrc,
		},
	)
	assert.DeepEqual(t, err, wantErr)
}
