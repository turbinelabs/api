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
