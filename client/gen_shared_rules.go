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

// This file was automatically generated by
//   github.com/turbinelabs/api/client/gen.go
// from
//   object.template.
// Any changes will be lost if this file is regenerated.

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

type httpSharedRulesV1 struct {
	dest apihttp.Endpoint

	requestHandler apihttp.RequestHandler
}

// Construct a new HTTP backed api.SharedRules API implementation.
//
// Parameters:
//	dest - service handling our HTTP requests; cf. NewService
func NewSharedRulesV1(
	dest apihttp.Endpoint,
) (*httpSharedRulesV1, error) {
	return &httpSharedRulesV1{dest, apihttp.NewRequestHandler(dest.Client())}, nil
}

// creates a sharedRules-scoped version of the specified path
func (hc *httpSharedRulesV1) path(p string) string {
	return "/v1.0/shared_rules" + p
}

// Construct a request to the associated sharedRules Endpoint with a specified
// method, path, query params, and body.
func (hc *httpSharedRulesV1) request(
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

func (hc *httpSharedRulesV1) get(path string, params apihttp.Params) (*http.Request, error) {
	return hc.request(http.MethodGet, path, params, "")
}

func (hc *httpSharedRulesV1) post(
	path string,
	params apihttp.Params,
	body string,
) (*http.Request, error) {
	return hc.request(http.MethodPost, path, params, body)
}

func (hc *httpSharedRulesV1) put(
	path string,
	params apihttp.Params,
	body string,
) (*http.Request, error) {
	return hc.request(http.MethodPut, path, params, body)
}

func (hc *httpSharedRulesV1) delete(path string, params apihttp.Params) (*http.Request, error) {
	return hc.request(http.MethodDelete, path, params, "")
}

func (hc *httpSharedRulesV1) Index(filters ...service.SharedRulesFilter) (api.SharedRulesSlice, error) {
	params := apihttp.Params{}

	if filters != nil && len(filters) != 0 {
		filterBytes, err := json.Marshal(filters)
		if err != nil {
			return nil, httperr.New400(
				fmt.Sprintf("unable to encode sharedRules filters: %v: %s", filters, err),
				httperr.UnknownUnclassifiedCode,
			)
		}

		params[queryargs.IndexFilters] = string(filterBytes)
	}

	response := make(api.SharedRulesSlice, 0, 10)
	reqFn := func() (*http.Request, error) { return hc.get("", params) }

	if err := hc.requestHandler.Do(reqFn, &response); err != nil {
		return nil, err
	}

	return response, nil
}

func (hc *httpSharedRulesV1) Get(key api.SharedRulesKey) (api.SharedRules, error) {
	if key == "" {
		return api.SharedRules{}, httperr.New400(
			"SharedRulesKey is a required parameter", httperr.ObjectKeyRequiredErrorCode)
	}

	reqFn := func() (*http.Request, error) {
		return hc.get("/"+url.QueryEscape(string(key)), nil)
	}

	response := api.SharedRules{}
	if err := hc.requestHandler.Do(reqFn, &response); err != nil {
		return api.SharedRules{}, err
	}

	return response, nil
}

func mkEncodeSharedRulesError(sharedRules api.SharedRules) *httperr.Error {
	msg := fmt.Sprintf("could not encode provided sharedRules: %+v", sharedRules)
	return httperr.New400(msg, httperr.UnknownEncodingCode)
}

func (hc *httpSharedRulesV1) Create(newSharedRules api.SharedRules) (api.SharedRules, error) {
	encoded := ""

	if b, err := json.Marshal(newSharedRules); err == nil {
		encoded = string(b)
	} else {
		return api.SharedRules{}, mkEncodeSharedRulesError(newSharedRules)
	}

	reqFn := func() (*http.Request, error) { return hc.post("", nil, encoded) }
	response := api.SharedRules{}
	if err := hc.requestHandler.Do(reqFn, &response); err != nil {
		return api.SharedRules{}, err
	}

	return response, nil
}

func (hc *httpSharedRulesV1) Modify(sharedRules api.SharedRules) (api.SharedRules, error) {
	encoded := ""

	if b, err := json.Marshal(sharedRules); err == nil {
		encoded = string(b)
	} else {
		return api.SharedRules{}, mkEncodeSharedRulesError(sharedRules)
	}

	response := api.SharedRules{}
	reqFn := func() (*http.Request, error) {
		return hc.put("/"+url.QueryEscape(string(sharedRules.SharedRulesKey)), nil, encoded)
	}

	if err := hc.requestHandler.Do(reqFn, &response); err != nil {
		return api.SharedRules{}, err
	}

	return response, nil
}

func (hc *httpSharedRulesV1) Delete(
	sharedRulesKey api.SharedRulesKey,
	checksum api.Checksum,
) error {
	reqFn := func() (*http.Request, error) {
		return hc.delete(
			"/"+url.QueryEscape(string(sharedRulesKey)),
			apihttp.Params{queryargs.Checksum: checksum.Checksum},
		)
	}

	if err := hc.requestHandler.Do(reqFn, nil); err != nil {
		return err
	}

	return nil
}

func (hc *httpSharedRulesV1) Purge(_ api.SharedRulesKey, _ api.Checksum) error {
	return httperr.New501("Purge not implemented", httperr.MiscErrorCode)
}
