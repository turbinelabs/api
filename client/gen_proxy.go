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
	httperr "github.com/turbinelabs/api/http/error"
	"github.com/turbinelabs/api/queryargs"
	"github.com/turbinelabs/api/service"
	tbnhttp "github.com/turbinelabs/client/http"
)

type httpProxyV1 struct {
	dest tbnhttp.Endpoint

	requestHandler tbnhttp.RequestHandler
}

// Construct a new HTTP backed Proxy API implementation.
//
// Parameters:
//	dest - service handling our HTTP requests; cf. NewService
//	apiKey - key used to sign our API requests; cf. NewService
//	client - HTTP client used to make these requests; must NOT be nil
func NewProxyV1(
	dest tbnhttp.Endpoint,
	apiKey string,
	client *http.Client,
) (*httpProxyV1, error) {
	if client == nil {
		// Future investigation note: when nil is passed in here the actual failure
		// is way upstream in a curious way; investigating could lend understanding
		// of some cool go internals
		return nil, fmt.Errorf("Attempting to configure Proxy with nil *http.Client")
	}
	return &httpProxyV1{dest, tbnhttp.NewRequestHandler(client, apiKey, apiClientID)}, nil
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
	params tbnhttp.Params,
	body string,
) (*http.Request, error) {
	rdr := strings.NewReader(body)
	req, err := http.NewRequest(string(method), hc.dest.Url(hc.path(path), params), rdr)

	if err != nil {
		return nil, err
	}

	return req, nil
}

func (hc *httpProxyV1) get(path string, params tbnhttp.Params) (*http.Request, error) {
	return hc.request(mGET, path, params, "")
}

func (hc *httpProxyV1) post(
	path string,
	params tbnhttp.Params,
	body string,
) (*http.Request, error) {
	return hc.request(mPOST, path, params, body)
}

func (hc *httpProxyV1) put(
	path string,
	params tbnhttp.Params,
	body string,
) (*http.Request, error) {
	return hc.request(mPUT, path, params, body)
}

func (hc *httpProxyV1) delete(path string, params tbnhttp.Params) (*http.Request, error) {
	return hc.request(mDELETE, path, params, "")
}

func (hc *httpProxyV1) Index(filters ...service.ProxyFilter) (api.Proxies, error) {
	params := tbnhttp.Params{}

	if filters != nil && len(filters) != 0 {
		filterBytes, e := json.Marshal(filters)
		if e != nil {
			return nil, httperr.New400(
				fmt.Sprintf("unable to encode proxy filters: %v", filters),
				httperr.UnknownUnclassifiedCode)
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
			tbnhttp.Params{queryargs.Checksum: checksum.Checksum},
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
