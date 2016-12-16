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

type httpProxyV1 struct {
	dest apihttp.Endpoint

	requestHandler apihttp.RequestHandler
}

// Construct a new HTTP backed Proxy API implementation.
//
// Parameters:
//	dest - service handling our HTTP requests; cf. NewService
func NewProxyV1(
	dest apihttp.Endpoint,
) (*httpProxyV1, error) {
	return &httpProxyV1{dest, apihttp.NewRequestHandler(dest.Client())}, nil
}

// creates a proxy-scoped version of the specified path
func (hc *httpProxyV1) path(p string) string {
	return "/v1.0/proxy" + p
}

// Construct a request to the associated proxy Endpoint with a specified
// method, path, query params, and body.
func (hc *httpProxyV1) request(
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

func (hc *httpProxyV1) get(path string, params apihttp.Params) (*http.Request, error) {
	return hc.request(mGET, path, params, "")
}

func (hc *httpProxyV1) post(
	path string,
	params apihttp.Params,
	body string,
) (*http.Request, error) {
	return hc.request(mPOST, path, params, body)
}

func (hc *httpProxyV1) put(
	path string,
	params apihttp.Params,
	body string,
) (*http.Request, error) {
	return hc.request(mPUT, path, params, body)
}

func (hc *httpProxyV1) delete(path string, params apihttp.Params) (*http.Request, error) {
	return hc.request(mDELETE, path, params, "")
}

func (hc *httpProxyV1) Index(filters ...service.ProxyFilter) (api.Proxies, error) {
	params := apihttp.Params{}

	if filters != nil && len(filters) != 0 {
		filterBytes, err := json.Marshal(filters)
		if err != nil {
			return nil, httperr.New400(
				fmt.Sprintf("unable to encode proxy filters: %v: %s", filters, err),
				httperr.UnknownUnclassifiedCode,
			)
		}

		params[queryargs.IndexFilters] = string(filterBytes)
	}

	response := make(api.Proxies, 0, 10)
	reqFn := func() (*http.Request, error) { return hc.get("", params) }

	if err := hc.requestHandler.Do(reqFn, &response); err != nil {
		return nil, err
	}

	return response, nil
}

func (hc *httpProxyV1) Get(key api.ProxyKey) (api.Proxy, error) {
	if key == "" {
		return api.Proxy{}, httperr.New400(
			"ProxyKey is a required parameter", httperr.ObjectKeyRequiredErrorCode)
	}

	reqFn := func() (*http.Request, error) {
		return hc.get("/"+url.QueryEscape(string(key)), nil)
	}

	response := api.Proxy{}
	if err := hc.requestHandler.Do(reqFn, &response); err != nil {
		return api.Proxy{}, err
	}

	return response, nil
}

func mkEncodeProxyError(proxy api.Proxy) *httperr.Error {
	msg := fmt.Sprintf("could not encode provided proxy: %+v", proxy)
	return httperr.New400(msg, httperr.UnknownEncodingCode)
}

func (hc *httpProxyV1) Create(newProxy api.Proxy) (api.Proxy, error) {
	encoded := ""

	if b, err := json.Marshal(newProxy); err == nil {
		encoded = string(b)
	} else {
		return api.Proxy{}, mkEncodeProxyError(newProxy)
	}

	reqFn := func() (*http.Request, error) { return hc.post("", nil, encoded) }
	response := api.Proxy{}
	if err := hc.requestHandler.Do(reqFn, &response); err != nil {
		return api.Proxy{}, err
	}

	return response, nil
}

func (hc *httpProxyV1) Modify(proxy api.Proxy) (api.Proxy, error) {
	encoded := ""

	if b, err := json.Marshal(proxy); err == nil {
		encoded = string(b)
	} else {
		return api.Proxy{}, mkEncodeProxyError(proxy)
	}

	response := api.Proxy{}
	reqFn := func() (*http.Request, error) {
		return hc.put("/"+url.QueryEscape(string(proxy.ProxyKey)), nil, encoded)
	}

	if err := hc.requestHandler.Do(reqFn, &response); err != nil {
		return api.Proxy{}, err
	}

	return response, nil
}

func (hc *httpProxyV1) Delete(
	proxyKey api.ProxyKey,
	checksum api.Checksum,
) error {
	reqFn := func() (*http.Request, error) {
		return hc.delete(
			"/"+url.QueryEscape(string(proxyKey)),
			apihttp.Params{queryargs.Checksum: checksum.Checksum},
		)
	}

	if err := hc.requestHandler.Do(reqFn, nil); err != nil {
		return err
	}

	return nil
}

func (hc *httpProxyV1) Purge(_ api.ProxyKey, _ api.Checksum) error {
	return httperr.New501("Purge not implemented", httperr.MiscErrorCode)
}
