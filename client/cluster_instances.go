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
	"strconv"

	"github.com/turbinelabs/api"
	apihttp "github.com/turbinelabs/api/http"
	"github.com/turbinelabs/api/queryargs"
)

func (hc *httpClusterV1) AddInstance(
	clusterKey api.ClusterKey,
	checksum api.Checksum,
	instance api.Instance,
) (api.Cluster, error) {
	encoded := ""

	if b, err := json.Marshal(instance); err != nil {
		return api.Cluster{}, err
	} else {
		encoded = string(b)
	}

	reqFn := func() (*http.Request, error) {
		return hc.post(
			fmt.Sprintf("/%s/instance", url.QueryEscape(string(clusterKey))),
			apihttp.Params{queryargs.Checksum: checksum.Checksum},
			encoded)
	}
	response := api.Cluster{}

	if err := hc.requestHandler.Do(reqFn, &response); err != nil {
		return api.Cluster{}, err
	}

	return response, nil
}

func (hc *httpClusterV1) RemoveInstance(
	clusterKey api.ClusterKey,
	checksum api.Checksum,
	instance api.Instance,
) (api.Cluster, error) {
	ckey := url.QueryEscape(string(clusterKey))
	host := url.QueryEscape(instance.Host)
	port := url.QueryEscape(strconv.Itoa(instance.Port))
	instPath := fmt.Sprintf("/%s/instance/%s:%s", ckey, host, port)

	reqFn := func() (*http.Request, error) {
		return hc.delete(instPath, apihttp.Params{queryargs.Checksum: checksum.Checksum})
	}
	response := api.Cluster{}

	if err := hc.requestHandler.Do(reqFn, &response); err != nil {
		return api.Cluster{}, err
	}

	return response, nil
}
