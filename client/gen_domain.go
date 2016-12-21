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

type httpDomainV1 struct {
	dest apihttp.Endpoint

	requestHandler apihttp.RequestHandler
}

// Construct a new HTTP backed Domain API implementation.
//
// Parameters:
//	dest - service handling our HTTP requests; cf. NewService
func NewDomainV1(
	dest apihttp.Endpoint,
) (*httpDomainV1, error) {
	return &httpDomainV1{dest, apihttp.NewRequestHandler(dest.Client())}, nil
}

// creates a domain-scoped version of the specified path
func (hc *httpDomainV1) path(p string) string {
	return "/v1.0/domain" + p
}

// Construct a request to the associated domain Endpoint with a specified
// method, path, query params, and body.
func (hc *httpDomainV1) request(
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

func (hc *httpDomainV1) get(path string, params apihttp.Params) (*http.Request, error) {
	return hc.request(mGET, path, params, "")
}

func (hc *httpDomainV1) post(
	path string,
	params apihttp.Params,
	body string,
) (*http.Request, error) {
	return hc.request(mPOST, path, params, body)
}

func (hc *httpDomainV1) put(
	path string,
	params apihttp.Params,
	body string,
) (*http.Request, error) {
	return hc.request(mPUT, path, params, body)
}

func (hc *httpDomainV1) delete(path string, params apihttp.Params) (*http.Request, error) {
	return hc.request(mDELETE, path, params, "")
}

func (hc *httpDomainV1) Index(filters ...service.DomainFilter) (api.Domains, error) {
	params := apihttp.Params{}

	if filters != nil && len(filters) != 0 {
		filterBytes, err := json.Marshal(filters)
		if err != nil {
			return nil, httperr.New400(
				fmt.Sprintf("unable to encode domain filters: %v: %s", filters, err),
				httperr.UnknownUnclassifiedCode,
			)
		}

		params[queryargs.IndexFilters] = string(filterBytes)
	}

	response := make(api.Domains, 0, 10)
	reqFn := func() (*http.Request, error) { return hc.get("", params) }

	if err := hc.requestHandler.Do(reqFn, &response); err != nil {
		return nil, err
	}

	return response, nil
}

func (hc *httpDomainV1) Get(key api.DomainKey) (api.Domain, error) {
	if key == "" {
		return api.Domain{}, httperr.New400(
			"DomainKey is a required parameter", httperr.ObjectKeyRequiredErrorCode)
	}

	reqFn := func() (*http.Request, error) {
		return hc.get("/"+url.QueryEscape(string(key)), nil)
	}

	response := api.Domain{}
	if err := hc.requestHandler.Do(reqFn, &response); err != nil {
		return api.Domain{}, err
	}

	return response, nil
}

func mkEncodeDomainError(domain api.Domain) *httperr.Error {
	msg := fmt.Sprintf("could not encode provided domain: %+v", domain)
	return httperr.New400(msg, httperr.UnknownEncodingCode)
}

func (hc *httpDomainV1) Create(newDomain api.Domain) (api.Domain, error) {
	encoded := ""

	if b, err := json.Marshal(newDomain); err == nil {
		encoded = string(b)
	} else {
		return api.Domain{}, mkEncodeDomainError(newDomain)
	}

	reqFn := func() (*http.Request, error) { return hc.post("", nil, encoded) }
	response := api.Domain{}
	if err := hc.requestHandler.Do(reqFn, &response); err != nil {
		return api.Domain{}, err
	}

	return response, nil
}

func (hc *httpDomainV1) Modify(domain api.Domain) (api.Domain, error) {
	encoded := ""

	if b, err := json.Marshal(domain); err == nil {
		encoded = string(b)
	} else {
		return api.Domain{}, mkEncodeDomainError(domain)
	}

	response := api.Domain{}
	reqFn := func() (*http.Request, error) {
		return hc.put("/"+url.QueryEscape(string(domain.DomainKey)), nil, encoded)
	}

	if err := hc.requestHandler.Do(reqFn, &response); err != nil {
		return api.Domain{}, err
	}

	return response, nil
}

func (hc *httpDomainV1) Delete(
	domainKey api.DomainKey,
	checksum api.Checksum,
) error {
	reqFn := func() (*http.Request, error) {
		return hc.delete(
			"/"+url.QueryEscape(string(domainKey)),
			apihttp.Params{queryargs.Checksum: checksum.Checksum},
		)
	}

	if err := hc.requestHandler.Do(reqFn, nil); err != nil {
		return err
	}

	return nil
}

func (hc *httpDomainV1) Purge(_ api.DomainKey, _ api.Checksum) error {
	return httperr.New501("Purge not implemented", httperr.MiscErrorCode)
}
