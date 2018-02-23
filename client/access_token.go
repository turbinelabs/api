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
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/turbinelabs/api"
	apihttp "github.com/turbinelabs/api/http"
	httperr "github.com/turbinelabs/api/http/error"
	"github.com/turbinelabs/api/queryargs"
	"github.com/turbinelabs/api/service"
)

type httpAccessTokenV1 struct {
	dest apihttp.Endpoint

	requestHandler apihttp.RequestHandler
}

// NewAccessTokenV1 constructs a new HTTP backed AccessToken API implementation.
//
// Parameters:
//   dest - service handling our HTTP requests; cf. NewService
func NewAccessTokenV1(
	dest apihttp.Endpoint,
) (*httpAccessTokenV1, error) {
	return &httpAccessTokenV1{dest, apihttp.NewRequestHandler(dest.Client())}, nil
}

// creates a accessToken-scoped version of the specified path
func (hc *httpAccessTokenV1) path(p string) string {
	return "/v1.0/admin/user/self/access_tokens" + p
}

// Construct a request to the associated accessToken Endpoint with a specified
// method, path, query params, and body.
func (hc *httpAccessTokenV1) request(
	method string,
	path string,
	params apihttp.Params,
	body string,
) (*http.Request, error) {
	rdr := strings.NewReader(body)
	req, err := hc.dest.NewRequest(string(method), hc.path(path), params, rdr)

	if err != nil {
		return nil, err
	}

	return req, nil
}

func (hc *httpAccessTokenV1) get(path string, params apihttp.Params) (*http.Request, error) {
	return hc.request(http.MethodGet, path, params, "")
}

func (hc *httpAccessTokenV1) post(
	path string,
	params apihttp.Params,
	body string,
) (*http.Request, error) {
	return hc.request(http.MethodPost, path, params, body)
}

func (hc *httpAccessTokenV1) put(
	path string,
	params apihttp.Params,
	body string,
) (*http.Request, error) {
	return hc.request(http.MethodPut, path, params, body)
}

func (hc *httpAccessTokenV1) delete(path string, params apihttp.Params) (*http.Request, error) {
	return hc.request(http.MethodDelete, path, params, "")
}

func (hc *httpAccessTokenV1) Index(filters ...service.AccessTokenFilter) (api.AccessTokens, error) {
	params := apihttp.Params{}

	if filters != nil && len(filters) != 0 {
		filterBytes, err := json.Marshal(filters)
		if err != nil {
			return nil, httperr.New400(
				fmt.Sprintf("unable to encode accessToken filters: %v: %s", filters, err),
				httperr.UnknownUnclassifiedCode,
			)
		}

		params[queryargs.IndexFilters] = string(filterBytes)
	}

	response := make(api.AccessTokens, 0, 10)
	reqFn := func() (*http.Request, error) { return hc.get("", params) }

	if err := hc.requestHandler.Do(reqFn, &response); err != nil {
		return nil, err
	}

	return response, nil
}

func (hc *httpAccessTokenV1) Get(key api.AccessTokenKey) (api.AccessToken, error) {
	if key == "" {
		return api.AccessToken{}, httperr.New400(
			"AccessTokenKey is a required parameter", httperr.ObjectKeyRequiredErrorCode)
	}

	reqFn := func() (*http.Request, error) {
		return hc.get("/"+url.QueryEscape(string(key)), nil)
	}

	response := api.AccessToken{}
	if err := hc.requestHandler.Do(reqFn, &response); err != nil {
		return api.AccessToken{}, err
	}

	return response, nil
}

func mkEncodeAccessTokenError(accessToken api.AccessToken) *httperr.Error {
	msg := fmt.Sprintf("could not encode provided accessToken: %+v", accessToken)
	return httperr.New400(msg, httperr.UnknownEncodingCode)
}

func (hc *httpAccessTokenV1) CreateAccessToken(desc string) (api.AccessToken, error) {
	return hc.Create(api.AccessToken{Description: desc})
}

func (hc *httpAccessTokenV1) Create(newAccessToken api.AccessToken) (api.AccessToken, error) {
	encoded := ""

	if b, err := json.Marshal(newAccessToken); err == nil {
		encoded = string(b)
	} else {
		return api.AccessToken{}, mkEncodeAccessTokenError(newAccessToken)
	}

	reqFn := func() (*http.Request, error) { return hc.post("", nil, encoded) }
	response := api.AccessToken{}
	if err := hc.requestHandler.Do(reqFn, &response); err != nil {
		return api.AccessToken{}, err
	}

	return response, nil
}

func (hc *httpAccessTokenV1) Delete(
	accessTokenKey api.AccessTokenKey,
	checksum api.Checksum,
) error {
	reqFn := func() (*http.Request, error) {
		return hc.delete(
			"/"+url.QueryEscape(string(accessTokenKey)),
			apihttp.Params{queryargs.Checksum: checksum.Checksum},
		)
	}

	if err := hc.requestHandler.Do(reqFn, nil); err != nil {
		return err
	}

	return nil
}
