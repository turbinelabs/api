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

// This file was automatically generated by genny.
// Any changes will be lost if this file is regenerated.
// see https://github.com/falun/genny

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

type httpUserV1 struct {
	dest apihttp.Endpoint

	requestHandler apihttp.RequestHandler
}

// Construct a new HTTP backed User API implementation.
//
// Parameters:
//	dest - service handling our HTTP requests; cf. NewService
func NewUserV1(
	dest apihttp.Endpoint,
) (*httpUserV1, error) {
	return &httpUserV1{dest, apihttp.NewRequestHandler(dest.Client())}, nil
}

// creates a user-scoped version of the specified path
func (hc *httpUserV1) path(p string) string {
	return "/v1.0/admin/user" + p
}

// Construct a request to the associated user Endpoint with a specified
// method, path, query params, and body.
func (hc *httpUserV1) request(
	method httpMethod,
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

func (hc *httpUserV1) get(path string, params apihttp.Params) (*http.Request, error) {
	return hc.request(mGET, path, params, "")
}

func (hc *httpUserV1) post(
	path string,
	params apihttp.Params,
	body string,
) (*http.Request, error) {
	return hc.request(mPOST, path, params, body)
}

func (hc *httpUserV1) put(
	path string,
	params apihttp.Params,
	body string,
) (*http.Request, error) {
	return hc.request(mPUT, path, params, body)
}

func (hc *httpUserV1) delete(path string, params apihttp.Params) (*http.Request, error) {
	return hc.request(mDELETE, path, params, "")
}

func (hc *httpUserV1) Index(filters ...service.UserFilter) (api.Users, error) {
	params := apihttp.Params{}

	if filters != nil && len(filters) != 0 {
		filterBytes, err := json.Marshal(filters)
		if err != nil {
			return nil, httperr.New400(
				fmt.Sprintf("unable to encode user filters: %v: %s", filters, err),
				httperr.UnknownUnclassifiedCode,
			)
		}

		params[queryargs.IndexFilters] = string(filterBytes)
	}

	response := make(api.Users, 0, 10)
	reqFn := func() (*http.Request, error) { return hc.get("", params) }

	if err := hc.requestHandler.Do(reqFn, &response); err != nil {
		return nil, err
	}

	return response, nil
}

func (hc *httpUserV1) Get(key api.UserKey) (api.User, error) {
	if key == "" {
		return api.User{}, httperr.New400(
			"UserKey is a required parameter", httperr.ObjectKeyRequiredErrorCode)
	}

	reqFn := func() (*http.Request, error) {
		return hc.get("/"+url.QueryEscape(string(key)), nil)
	}

	response := api.User{}
	if err := hc.requestHandler.Do(reqFn, &response); err != nil {
		return api.User{}, err
	}

	return response, nil
}

func mkEncodeUserError(user api.User) *httperr.Error {
	msg := fmt.Sprintf("could not encode provided user: %+v", user)
	return httperr.New400(msg, httperr.UnknownEncodingCode)
}

func (hc *httpUserV1) Create(newUser api.User) (api.User, error) {
	encoded := ""

	if b, err := json.Marshal(newUser); err == nil {
		encoded = string(b)
	} else {
		return api.User{}, mkEncodeUserError(newUser)
	}

	reqFn := func() (*http.Request, error) { return hc.post("", nil, encoded) }
	response := api.User{}
	if err := hc.requestHandler.Do(reqFn, &response); err != nil {
		return api.User{}, err
	}

	return response, nil
}

func (hc *httpUserV1) Modify(user api.User) (api.User, error) {
	encoded := ""

	if b, err := json.Marshal(user); err == nil {
		encoded = string(b)
	} else {
		return api.User{}, mkEncodeUserError(user)
	}

	response := api.User{}
	reqFn := func() (*http.Request, error) {
		return hc.put("/"+url.QueryEscape(string(user.UserKey)), nil, encoded)
	}

	if err := hc.requestHandler.Do(reqFn, &response); err != nil {
		return api.User{}, err
	}

	return response, nil
}

func (hc *httpUserV1) Delete(
	userKey api.UserKey,
	checksum api.Checksum,
) error {
	reqFn := func() (*http.Request, error) {
		return hc.delete(
			"/"+url.QueryEscape(string(userKey)),
			apihttp.Params{queryargs.Checksum: checksum.Checksum},
		)
	}

	if err := hc.requestHandler.Do(reqFn, nil); err != nil {
		return err
	}

	return nil
}

func (hc *httpUserV1) Purge(_ api.UserKey, _ api.Checksum) error {
	return httperr.New501("Purge not implemented", httperr.MiscErrorCode)
}
