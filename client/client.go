/*
	Implements the Turbine service interface by constructing a HTTP client to a
	specified backend api server.

	In order to configure the desired Turbine service and Endpoint must be
	specified which then may be passed into the constructor for a HTTP service:
		endpoint, err := NewEndpoint(HTTP, "dev.turbinelabs.io", 8080)
		if err != nil {
			return err
		}

		service := NewService(Endpoint, apiKey, nil)

	For a detailed discussion about what each of these values mean see method
	docs.
*/
package client

import (
	"net/http"

	"github.com/turbinelabs/api/service"
	tbnhttp "github.com/turbinelabs/client/http"
)

type httpMethod string

const (
	mGET    httpMethod = "GET"
	mPUT               = "PUT"
	mPOST              = "POST"
	mDELETE            = "DELETE"
)

const apiClientID string = "tbn-api-client (v0.1)"

// Create a new Service backed by a Turbine api server at dest. Communication
// with this server will happen via HTTP (or HTTPS as specified in the
// Endpoint) and will sign your requests with the provided apiKey.
//
// Service creation can not fail but it does not guarantee that the target
// Endpoint is a valid, or live, Turbine service.
//
// Parameters:
//	dest - a server we want to communicate with, construct via tbnhttp.NewEndpoint
//	apiKey - an API key assigned to your organization during setup process
// 	httpClient
//	  nil (should be the most common value) will default to a standard
//	  http.Client that has been modified to carry forward headers if it
//	  sees a 3xx redirect (cf. tbnhttp.HeaderPreserving()). Other values
//	  may be used if custom behavior is necessary.
func NewAll(
	dest tbnhttp.Endpoint,
	apiKey string,
	httpClient *http.Client,
) (service.All, error) {
	if httpClient == nil {
		httpClient = tbnhttp.HeaderPreserving()
	}

	c, err := NewClusterV1(dest, apiKey, httpClient)
	if err != nil {
		return nil, err
	}
	d, err := NewDomainV1(dest, apiKey, httpClient)
	if err != nil {
		return nil, err
	}
	r, err := NewRouteV1(dest, apiKey, httpClient)
	if err != nil {
		return nil, err
	}
	s, err := NewSharedRulesV1(dest, apiKey, httpClient)
	if err != nil {
		return nil, err
	}
	p, err := NewProxyV1(dest, apiKey, httpClient)
	if err != nil {
		return nil, err
	}
	z, err := NewZoneV1(dest, apiKey, httpClient)
	if err != nil {
		return nil, err
	}
	h, err := NewHistoryV1(dest, apiKey, httpClient)
	if err != nil {
		return nil, err
	}

	httpService := httpServiceV1{c, d, r, s, p, z, h}

	return &httpService, nil
}

// Create a new Admin backed by a Turbine api server at dest. Communication
// with this server will happen via HTTP (or HTTPS as specified in the
// Endpoint) and will sign your requests with the provided apiKey.
//
// Service creation can not fail but it does not guarantee that the target
// Endpoint is a valid, or live, Turbine service.
//
// Parameters: See NewService.
func NewAdmin(
	dest tbnhttp.Endpoint,
	apiKey string,
	httpClient *http.Client,
) (service.Admin, error) {
	if httpClient == nil {
		httpClient = tbnhttp.HeaderPreserving()
	}

	u, err := NewUserV1(dest, apiKey, httpClient)
	if err != nil {
		return nil, err
	}

	httpAdmin := httpAdminV1{u}

	return &httpAdmin, nil
}

// v1 http-backed service that implements service.All via HTTP calls to some
// backend. For the implementation of each sub interface see in gen_XYZ.go
type httpServiceV1 struct {
	clusterV1     *httpClusterV1
	domainV1      *httpDomainV1
	routeV1       *httpRouteV1
	sharedRulesV1 *httpSharedRulesV1
	proxyV1       *httpProxyV1
	zoneV1        *httpZoneV1
	historyV1     *httpHistoryV1
}

// Returns an implementation of service.Cluster.
func (hs *httpServiceV1) Cluster() service.Cluster {
	return hs.clusterV1
}

// Returns an implementation of service.Domain.
func (hs *httpServiceV1) Domain() service.Domain {
	return hs.domainV1
}

// Returns an implementation of service.SharedRules.
func (hs *httpServiceV1) SharedRules() service.SharedRules {
	return hs.sharedRulesV1
}

// Returns an implementation of service.Route.
func (hs *httpServiceV1) Route() service.Route {
	return hs.routeV1
}

// Returns an implementation of service.Proxy.
func (hs *httpServiceV1) Proxy() service.Proxy {
	return hs.proxyV1
}

// Returns an implementation of service.Zone.
func (hs *httpServiceV1) Zone() service.Zone {
	return hs.zoneV1
}

// Returns an implementation of service.History.
func (hs *httpServiceV1) History() service.History {
	return hs.historyV1
}

// v1 http-backed service that implements Admin via HTTP calls to some
// backend. For the implementation of each sub interface see in gen_XYZ.go
type httpAdminV1 struct {
	userV1 *httpUserV1
}

// Returns an implementation of service.User.
func (as *httpAdminV1) User() service.User {
	return as.userV1
}
