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

type httpClusterV1 struct {
	dest apihttp.Endpoint

	requestHandler apihttp.RequestHandler
}

// Construct a new HTTP backed api.Cluster API implementation.
//
// Parameters:
//	dest - service handling our HTTP requests; cf. NewService
func NewClusterV1(
	dest apihttp.Endpoint,
) (*httpClusterV1, error) {
	return &httpClusterV1{dest, apihttp.NewRequestHandler(dest.Client())}, nil
}

// creates a cluster-scoped version of the specified path
func (hc *httpClusterV1) path(p string) string {
	return "/v1.0/cluster" + p
}

// Construct a request to the associated cluster Endpoint with a specified
// method, path, query params, and body.
func (hc *httpClusterV1) request(
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

func (hc *httpClusterV1) get(path string, params apihttp.Params) (*http.Request, error) {
	return hc.request(http.MethodGet, path, params, "")
}

func (hc *httpClusterV1) post(
	path string,
	params apihttp.Params,
	body string,
) (*http.Request, error) {
	return hc.request(http.MethodPost, path, params, body)
}

func (hc *httpClusterV1) put(
	path string,
	params apihttp.Params,
	body string,
) (*http.Request, error) {
	return hc.request(http.MethodPut, path, params, body)
}

func (hc *httpClusterV1) delete(path string, params apihttp.Params) (*http.Request, error) {
	return hc.request(http.MethodDelete, path, params, "")
}

func (hc *httpClusterV1) Index(filters ...service.ClusterFilter) (api.Clusters, error) {
	params := apihttp.Params{}

	if filters != nil && len(filters) != 0 {
		filterBytes, err := json.Marshal(filters)
		if err != nil {
			return nil, httperr.New400(
				fmt.Sprintf("unable to encode cluster filters: %v: %s", filters, err),
				httperr.UnknownUnclassifiedCode,
			)
		}

		params[queryargs.IndexFilters] = string(filterBytes)
	}

	response := make(api.Clusters, 0, 10)
	reqFn := func() (*http.Request, error) { return hc.get("", params) }

	if err := hc.requestHandler.Do(reqFn, &response); err != nil {
		return nil, err
	}

	return response, nil
}

func (hc *httpClusterV1) Get(key api.ClusterKey) (api.Cluster, error) {
	if key == "" {
		return api.Cluster{}, httperr.New400(
			"ClusterKey is a required parameter", httperr.ObjectKeyRequiredErrorCode)
	}

	reqFn := func() (*http.Request, error) {
		return hc.get("/"+url.QueryEscape(string(key)), nil)
	}

	response := api.Cluster{}
	if err := hc.requestHandler.Do(reqFn, &response); err != nil {
		return api.Cluster{}, err
	}

	return response, nil
}

func mkEncodeClusterError(cluster api.Cluster) *httperr.Error {
	msg := fmt.Sprintf("could not encode provided cluster: %+v", cluster)
	return httperr.New400(msg, httperr.UnknownEncodingCode)
}

func (hc *httpClusterV1) Create(newCluster api.Cluster) (api.Cluster, error) {
	encoded := ""

	if b, err := json.Marshal(newCluster); err == nil {
		encoded = string(b)
	} else {
		return api.Cluster{}, mkEncodeClusterError(newCluster)
	}

	reqFn := func() (*http.Request, error) { return hc.post("", nil, encoded) }
	response := api.Cluster{}
	if err := hc.requestHandler.Do(reqFn, &response); err != nil {
		return api.Cluster{}, err
	}

	return response, nil
}

func (hc *httpClusterV1) Modify(cluster api.Cluster) (api.Cluster, error) {
	encoded := ""

	if b, err := json.Marshal(cluster); err == nil {
		encoded = string(b)
	} else {
		return api.Cluster{}, mkEncodeClusterError(cluster)
	}

	response := api.Cluster{}
	reqFn := func() (*http.Request, error) {
		return hc.put("/"+url.QueryEscape(string(cluster.ClusterKey)), nil, encoded)
	}

	if err := hc.requestHandler.Do(reqFn, &response); err != nil {
		return api.Cluster{}, err
	}

	return response, nil
}

func (hc *httpClusterV1) Delete(
	clusterKey api.ClusterKey,
	checksum api.Checksum,
) error {
	reqFn := func() (*http.Request, error) {
		return hc.delete(
			"/"+url.QueryEscape(string(clusterKey)),
			apihttp.Params{queryargs.Checksum: checksum.Checksum},
		)
	}

	if err := hc.requestHandler.Do(reqFn, nil); err != nil {
		return err
	}

	return nil
}

func (hc *httpClusterV1) Purge(_ api.ClusterKey, _ api.Checksum) error {
	return httperr.New501("Purge not implemented", httperr.MiscErrorCode)
}
