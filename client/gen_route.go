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
) /*
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

type httpRouteV1 struct {
	dest apihttp.Endpoint

	requestHandler apihttp.RequestHandler
}

// Construct a new HTTP backed Route API implementation.
//
// Parameters:
//	dest - service handling our HTTP requests; cf. NewService
func NewRouteV1(
	dest apihttp.Endpoint,
) (*httpRouteV1, error) {
	return &httpRouteV1{dest, apihttp.NewRequestHandler(dest.Client())}, nil
}

// creates a route-scoped version of the specified path
func (hc *httpRouteV1) path(p string) string {
	return "/v1.0/route" + p
}

// Construct a request to the associated route Endpoint with a specified
// method, path, query params, and body.
func (hc *httpRouteV1) request(
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

func (hc *httpRouteV1) get(path string, params apihttp.Params) (*http.Request, error) {
	return hc.request(http.MethodGet, path, params, "")
}

func (hc *httpRouteV1) post(
	path string,
	params apihttp.Params,
	body string,
) (*http.Request, error) {
	return hc.request(http.MethodPost, path, params, body)
}

func (hc *httpRouteV1) put(
	path string,
	params apihttp.Params,
	body string,
) (*http.Request, error) {
	return hc.request(http.MethodPut, path, params, body)
}

func (hc *httpRouteV1) delete(path string, params apihttp.Params) (*http.Request, error) {
	return hc.request(http.MethodDelete, path, params, "")
}

func (hc *httpRouteV1) Index(filters ...service.RouteFilter) (api.Routes, error) {
	params := apihttp.Params{}

	if filters != nil && len(filters) != 0 {
		filterBytes, err := json.Marshal(filters)
		if err != nil {
			return nil, httperr.New400(
				fmt.Sprintf("unable to encode route filters: %v: %s", filters, err),
				httperr.UnknownUnclassifiedCode,
			)
		}

		params[queryargs.IndexFilters] = string(filterBytes)
	}

	response := make(api.Routes, 0, 10)
	reqFn := func() (*http.Request, error) { return hc.get("", params) }

	if err := hc.requestHandler.Do(reqFn, &response); err != nil {
		return nil, err
	}

	return response, nil
}

func (hc *httpRouteV1) Get(key api.RouteKey) (api.Route, error) {
	if key == "" {
		return api.Route{}, httperr.New400(
			"RouteKey is a required parameter", httperr.ObjectKeyRequiredErrorCode)
	}

	reqFn := func() (*http.Request, error) {
		return hc.get("/"+url.QueryEscape(string(key)), nil)
	}

	response := api.Route{}
	if err := hc.requestHandler.Do(reqFn, &response); err != nil {
		return api.Route{}, err
	}

	return response, nil
}

func mkEncodeRouteError(route api.Route) *httperr.Error {
	msg := fmt.Sprintf("could not encode provided route: %+v", route)
	return httperr.New400(msg, httperr.UnknownEncodingCode)
}

func (hc *httpRouteV1) Create(newRoute api.Route) (api.Route, error) {
	encoded := ""

	if b, err := json.Marshal(newRoute); err == nil {
		encoded = string(b)
	} else {
		return api.Route{}, mkEncodeRouteError(newRoute)
	}

	reqFn := func() (*http.Request, error) { return hc.post("", nil, encoded) }
	response := api.Route{}
	if err := hc.requestHandler.Do(reqFn, &response); err != nil {
		return api.Route{}, err
	}

	return response, nil
}

func (hc *httpRouteV1) Modify(route api.Route) (api.Route, error) {
	encoded := ""

	if b, err := json.Marshal(route); err == nil {
		encoded = string(b)
	} else {
		return api.Route{}, mkEncodeRouteError(route)
	}

	response := api.Route{}
	reqFn := func() (*http.Request, error) {
		return hc.put("/"+url.QueryEscape(string(route.RouteKey)), nil, encoded)
	}

	if err := hc.requestHandler.Do(reqFn, &response); err != nil {
		return api.Route{}, err
	}

	return response, nil
}

func (hc *httpRouteV1) Delete(
	routeKey api.RouteKey,
	checksum api.Checksum,
) error {
	reqFn := func() (*http.Request, error) {
		return hc.delete(
			"/"+url.QueryEscape(string(routeKey)),
			apihttp.Params{queryargs.Checksum: checksum.Checksum},
		)
	}

	if err := hc.requestHandler.Do(reqFn, nil); err != nil {
		return err
	}

	return nil
}

func (hc *httpRouteV1) Purge(_ api.RouteKey, _ api.Checksum) error {
	return httperr.New501("Purge not implemented", httperr.MiscErrorCode)
}
