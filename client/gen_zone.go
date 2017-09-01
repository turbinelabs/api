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

type httpZoneV1 struct {
	dest apihttp.Endpoint

	requestHandler apihttp.RequestHandler
}

// Construct a new HTTP backed Zone API implementation.
//
// Parameters:
//	dest - service handling our HTTP requests; cf. NewService
func NewZoneV1(
	dest apihttp.Endpoint,
) (*httpZoneV1, error) {
	return &httpZoneV1{dest, apihttp.NewRequestHandler(dest.Client())}, nil
}

// creates a zone-scoped version of the specified path
func (hc *httpZoneV1) path(p string) string {
	return "/v1.0/zone" + p
}

// Construct a request to the associated zone Endpoint with a specified
// method, path, query params, and body.
func (hc *httpZoneV1) request(
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

func (hc *httpZoneV1) get(path string, params apihttp.Params) (*http.Request, error) {
	return hc.request(http.MethodGet, path, params, "")
}

func (hc *httpZoneV1) post(
	path string,
	params apihttp.Params,
	body string,
) (*http.Request, error) {
	return hc.request(http.MethodPost, path, params, body)
}

func (hc *httpZoneV1) put(
	path string,
	params apihttp.Params,
	body string,
) (*http.Request, error) {
	return hc.request(http.MethodPut, path, params, body)
}

func (hc *httpZoneV1) delete(path string, params apihttp.Params) (*http.Request, error) {
	return hc.request(http.MethodDelete, path, params, "")
}

func (hc *httpZoneV1) Index(filters ...service.ZoneFilter) (api.Zones, error) {
	params := apihttp.Params{}

	if filters != nil && len(filters) != 0 {
		filterBytes, err := json.Marshal(filters)
		if err != nil {
			return nil, httperr.New400(
				fmt.Sprintf("unable to encode zone filters: %v: %s", filters, err),
				httperr.UnknownUnclassifiedCode,
			)
		}

		params[queryargs.IndexFilters] = string(filterBytes)
	}

	response := make(api.Zones, 0, 10)
	reqFn := func() (*http.Request, error) { return hc.get("", params) }

	if err := hc.requestHandler.Do(reqFn, &response); err != nil {
		return nil, err
	}

	return response, nil
}

func (hc *httpZoneV1) Get(key api.ZoneKey) (api.Zone, error) {
	if key == "" {
		return api.Zone{}, httperr.New400(
			"ZoneKey is a required parameter", httperr.ObjectKeyRequiredErrorCode)
	}

	reqFn := func() (*http.Request, error) {
		return hc.get("/"+url.QueryEscape(string(key)), nil)
	}

	response := api.Zone{}
	if err := hc.requestHandler.Do(reqFn, &response); err != nil {
		return api.Zone{}, err
	}

	return response, nil
}

func mkEncodeZoneError(zone api.Zone) *httperr.Error {
	msg := fmt.Sprintf("could not encode provided zone: %+v", zone)
	return httperr.New400(msg, httperr.UnknownEncodingCode)
}

func (hc *httpZoneV1) Create(newZone api.Zone) (api.Zone, error) {
	encoded := ""

	if b, err := json.Marshal(newZone); err == nil {
		encoded = string(b)
	} else {
		return api.Zone{}, mkEncodeZoneError(newZone)
	}

	reqFn := func() (*http.Request, error) { return hc.post("", nil, encoded) }
	response := api.Zone{}
	if err := hc.requestHandler.Do(reqFn, &response); err != nil {
		return api.Zone{}, err
	}

	return response, nil
}

func (hc *httpZoneV1) Modify(zone api.Zone) (api.Zone, error) {
	encoded := ""

	if b, err := json.Marshal(zone); err == nil {
		encoded = string(b)
	} else {
		return api.Zone{}, mkEncodeZoneError(zone)
	}

	response := api.Zone{}
	reqFn := func() (*http.Request, error) {
		return hc.put("/"+url.QueryEscape(string(zone.ZoneKey)), nil, encoded)
	}

	if err := hc.requestHandler.Do(reqFn, &response); err != nil {
		return api.Zone{}, err
	}

	return response, nil
}

func (hc *httpZoneV1) Delete(
	zoneKey api.ZoneKey,
	checksum api.Checksum,
) error {
	reqFn := func() (*http.Request, error) {
		return hc.delete(
			"/"+url.QueryEscape(string(zoneKey)),
			apihttp.Params{queryargs.Checksum: checksum.Checksum},
		)
	}

	if err := hc.requestHandler.Do(reqFn, nil); err != nil {
		return err
	}

	return nil
}

func (hc *httpZoneV1) Purge(_ api.ZoneKey, _ api.Checksum) error {
	return httperr.New501("Purge not implemented", httperr.MiscErrorCode)
}
