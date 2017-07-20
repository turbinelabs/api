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

// Package client implements the Turbine service interface by constructing an
// HTTP client to a specified backend api server.
//
// In order to configure the desired Turbine service and Endpoint must be
// specified which then may be passed into the constructor for a HTTP service:
// 	endpoint, err := NewEndpoint(HTTP, "dev.turbinelabs.io", 8080)
// 	if err != nil {
// 		return err
// 	}
//
// 	service := NewAll(Endpoint, apiKey)
//
// For a detailed discussion about what each of these values mean see method
// docs.
package client

import (
	"github.com/turbinelabs/api"
	apihttp "github.com/turbinelabs/api/http"
	apiheader "github.com/turbinelabs/api/http/header"
	"github.com/turbinelabs/api/service"
)

type httpMethod string

// App is passed in the X-Tbn-Client-App header in API calls
type App string

const (
	mGET    httpMethod = "GET"
	mPUT               = "PUT"
	mPOST              = "POST"
	mDELETE            = "DELETE"
)

// clientType is the value sent for the X-Tbn-Client-Type header
const clientType = "github.com/turbinelabs/api/client"

// Create a new Service backed by a Turbine api server at dest. Communication
// with this server will happen via HTTP (or HTTPS as specified in the
// Endpoint) and will sign your requests with the provided apiKey. The
// Endpoint is copied so changes to headers or clients must be made before
// invoking NewAdmin.
//
// Service creation can not fail but it does not guarantee that the target
// Endpoint is a valid, or live, Turbine service.
//
// Parameters:
//	dest - a server we want to communicate with, construct via apihttp.NewEndpoint
//	apiKey - an API key assigned to your organization during setup process
func NewAll(
	dest apihttp.Endpoint,
	apiKey string,
	clientApp App,
) (service.All, error) {
	dest = configureEndpoint(dest, apiKey, clientApp)
	c, err := NewClusterV1(dest)
	if err != nil {
		return nil, err
	}
	d, err := NewDomainV1(dest)
	if err != nil {
		return nil, err
	}
	r, err := NewRouteV1(dest)
	if err != nil {
		return nil, err
	}
	s, err := NewSharedRulesV1(dest)
	if err != nil {
		return nil, err
	}
	p, err := NewProxyV1(dest)
	if err != nil {
		return nil, err
	}
	z, err := NewZoneV1(dest)
	if err != nil {
		return nil, err
	}
	h, err := NewHistoryV1(dest)
	if err != nil {
		return nil, err
	}

	httpService := httpServiceV1{c, d, r, s, p, z, h}

	return &httpService, nil
}

// Create a new Admin backed by a Turbine api server at dest. Communication
// with this server will happen via HTTP (or HTTPS as specified in the
// Endpoint) and will sign your requests with the provided apiKey. The
// Endpoint is copied so changes to headers or clients must be made before
// invoking NewAdmin.
//
// Service creation can not fail but it does not guarantee that the target
// Endpoint is a valid, or live, Turbine service.
//
// Parameters: See NewAll.
func NewAdmin(
	dest apihttp.Endpoint,
	apiKey string,
	clientApp App,
) (service.Admin, error) {
	dest = configureEndpoint(dest, apiKey, clientApp)

	u, err := NewUserV1(dest)
	if err != nil {
		return nil, err
	}

	httpAdmin := httpAdminV1{u}

	return &httpAdmin, nil
}

func configureEndpoint(dest apihttp.Endpoint, apiKey string, clientApp App) apihttp.Endpoint {
	// Copy the Endpoint to avoid polluting the original with our
	// headers.
	dest = dest.Copy()
	if apiKey != "" {
		dest.AddHeader(apiheader.Authorization, apiKey)
	}
	dest.AddHeader(apiheader.ClientType, clientType)
	dest.AddHeader(apiheader.ClientVersion, api.TbnPublicVersion)
	dest.AddHeader(apiheader.ClientApp, string(clientApp))
	return dest
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
