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

// Package service defines interfaces representing the Turbine Labs public API
package service

//go:generate mockgen -source $GOFILE -destination mock_$GOFILE -package $GOPACKAGE --write_package_comment=false

import (
	"time"

	"github.com/turbinelabs/api"
	tbntime "github.com/turbinelabs/nonstdlib/time"
)

// None is a monitor value used to incidate an empty array
const None = "-"

// All defines the interface for the public JSON/REST API.
//
// We define it in Go because it's convenient to allow our
// Go clients and servers to share a common language-level
// interface.
//
// Where necessary, we add commentary to describe things about
// the REST api that aren't documented by the language itself
// (eg paths, methods, etc)
//
// Each of the sub-interfaces presents an organization-scoped view.
// The method of scoping is not specified here, but in HTTP implementations
// will be the use of an authorization header.
type All interface {
	Cluster() Cluster
	Domain() Domain
	SharedRules() SharedRules
	Route() Route
	Proxy() Proxy
	Listener() Listener
	Zone() Zone
	History() History
}

// ClusterFilter describes a filter on the full list of Clusters
type ClusterFilter struct {
	ClusterKey api.ClusterKey `json:"cluster_key"`
	Name       string         `json:"name"`
	ZoneKey    api.ZoneKey    `json:"zone_key"`
	OrgKey     api.OrgKey     `json:"org_key"`
}

// IsNil returns true if the receiver is the zero value
func (cf ClusterFilter) IsNil() bool {
	return cf.Equals(ClusterFilter{})
}

// Equals returns true if the target is equal to the receiver
func (cf ClusterFilter) Equals(o ClusterFilter) bool {
	return cf == o
}

// Cluster describes the CRUD interface for api.Clusters
type Cluster interface {
	// GET /v1.0/cluster
	//
	// Index returns all Clusters to which the given filters apply. All non-empty
	// fields in a filter must apply for the filter to apply. Any Cluster to which
	// any filter applies is included in the result.
	//
	// If no filters are supplied, all Clusters are returned.
	Index(filters ...ClusterFilter) (api.Clusters, error)

	// GET /v1.0/cluster/<string:clusterKey>[?include_deleted]
	//
	// Get returns a Cluster for the given ClusterKey. If the Cluster does not
	// exist, an error is returned.
	Get(clusterKey api.ClusterKey) (api.Cluster, error)

	// POST /v1.0/cluster
	//
	// Create creates the given Cluster. Cluster Names must be unique for a given
	// ZoneKey. If a ClusterKey is specified in the Cluster, it is ignored and
	// replaced in the result with the authoritative ClusterKey.
	Create(cluster api.Cluster) (api.Cluster, error)

	// PUT /v1.0/cluster/<string:clusterKey>
	//
	// Modify modifies the given Cluster. Cluster Names must be unique for a given
	// ZoneKey. The given Cluster Checksum must match the existing Checksum.
	Modify(cluster api.Cluster) (api.Cluster, error)

	// DELETE /v1.0/cluster/<string:clusterKey>?checksum=<string:checksum>
	//
	// Delete completely removes the Cluster from the database.
	// If the checksum does not match no action is taken and an error
	// is returned.
	Delete(clusterKey api.ClusterKey, checksum api.Checksum) error

	// POST /v1.0/cluster/<string:clusterKey>/instances
	//
	// AddInstance adds an Instance to the Cluster corresponding to the given
	// ClusterKey. The given Cluster Checksum must match the existing Checksum.
	// If the Instance already exists, it will be updated.
	AddInstance(
		clusterKey api.ClusterKey,
		checksum api.Checksum,
		instance api.Instance,
	) (api.Cluster, error)

	// DELETE /v1.0/cluster/<string:clusterKey>/instances/<string:host>:<int:port>
	//
	// RemoveInstance removes an Instance from the Cluster corresponding to the
	// given ClusterKey. The given Cluster Checksum must match the existing
	// Checksum.
	RemoveInstance(
		clusterKey api.ClusterKey,
		checksum api.Checksum,
		instance api.Instance,
	) (api.Cluster, error)
}

// DomainFilter describes a filter on the full list of Domains
// per #5628 we should add the ability to find Listeners a Domain is attached to
type DomainFilter struct {
	DomainKey api.DomainKey `json:"domain_key"`
	Name      string        `json:"name"`
	ZoneKey   api.ZoneKey   `json:"zone_key"`
	OrgKey    api.OrgKey    `json:"org_key"`
	// ProxyKeys matches Domains with a superset of the specified ProxyKeys. A
	// slice with a single value of "-" will produce Domains with no linked
	// Proxies.
	ProxyKeys []api.ProxyKey `json:"proxy_keys"`
}

// HasNoProxies returns true if ProxyKeys has been set to the monitor value
// indicating a filter for Domains with no linked Proxies.
func (df DomainFilter) HasNoProxies() bool {
	return len(df.ProxyKeys) == 1 && df.ProxyKeys[0] == None
}

// IsNil returns true if the receiver is the zero value
func (df DomainFilter) IsNil() bool {
	return df.Equals(DomainFilter{})
}

// Equals returns true if the target is equal to the receiver
func (df DomainFilter) Equals(o DomainFilter) bool {
	if !(df.DomainKey == o.DomainKey &&
		df.Name == o.Name &&
		df.ZoneKey == o.ZoneKey &&
		df.OrgKey == o.OrgKey &&
		len(df.ProxyKeys) == len(o.ProxyKeys)) {
		return false
	}

	m := make(map[string]bool)
	for _, e := range df.ProxyKeys {
		m[string(e)] = true
	}
	for _, e := range o.ProxyKeys {
		if !m[string(e)] {
			return false
		}
	}

	return true
}

// Domain describes the CRUD interface for api.Domains
type Domain interface {
	// GET /v1.0/domains
	//
	// Index returns all Domains to which the given filters apply. All non-empty
	// fields in a filter must apply for the filter to apply. Any Domain to which
	// any filter applies is included in the result.
	//
	// If no filters are supplied, all Domains are returned.
	Index(filters ...DomainFilter) (api.Domains, error)

	// GET /v1.0/domains/<string:domainKey>[?include_deleted]
	//
	// Get returns a Domain for the given DomainKey. If the Domain does not
	// exist, an error is returned.
	Get(domainKey api.DomainKey) (api.Domain, error)

	// POST /v1.0/domains
	//
	// Create creates the given Domain. The tuple of (Host, Port, ZoneKey) must be
	// unique. If a DomainKey is specified in the Domain, it is ignored and
	// replaced in the result with the authoritative DomainKey.
	Create(domain api.Domain) (api.Domain, error)

	// PUT /v1.0/domains/<string:domainKey>
	//
	// Modify modifies the given Domain. The tuple of (Host, Port, ZoneKey) must
	// be unique. The given Domain Checksum must match the existing Checksum.
	Modify(domain api.Domain) (api.Domain, error)

	// DELETE /v1.0/domains/<string:domainKey>?checksum=<string:checksum>
	//
	// Delete completely removes the Domain from the database.
	// If the checksum does not match no action is taken and an error
	// is returned.
	Delete(domainKey api.DomainKey, checksum api.Checksum) error
}

// ProxyFilter describes a filter on the full list of Proxies
type ProxyFilter struct {
	ProxyKey api.ProxyKey `json:"proxy_key"`
	Name     string       `json:"name"`
	// DomainKeys matches Proxies with a superset of the specified DomainKeys. A
	// slice with a single value of "-" will produce Proxies with no linked
	// Domains.
	DomainKeys []api.DomainKey `json:"domain_keys"`
	// ListenerKeys matches Proxies with a superset of the specified ListenerKeys. A
	// slice with a single value of "-" will produce Proxies with no linked
	// Listeners.
	ListenerKeys []api.ListenerKey `json:"listener_keys"`
	ZoneKey      api.ZoneKey       `json:"zone_key"`
	OrgKey       api.OrgKey        `json:"org_key"`
}

// HasNoDomains returns true if DomainKeys has been set to the monitor
// value indicating a filter for Proxies with no linked Domains.
func (pf ProxyFilter) HasNoDomains() bool {
	return len(pf.DomainKeys) == 1 && pf.DomainKeys[0] == None
}

// HasNoListeners returns true if ListenerKeys has been set to the monitor
// value indicating a filter for Proxies with no linked Listeners.
func (pf ProxyFilter) HasNoListeners() bool {
	return len(pf.ListenerKeys) == 1 && pf.ListenerKeys[0] == None
}

// IsNil returns true if the receiver is the zero value
func (pf ProxyFilter) IsNil() bool {
	return pf.Equals(ProxyFilter{})
}

// Equals returns true if the target is equal to the receiver
func (pf ProxyFilter) Equals(o ProxyFilter) bool {
	if !(pf.ProxyKey == o.ProxyKey &&
		pf.Name == o.Name &&
		pf.ZoneKey == o.ZoneKey &&
		pf.OrgKey == o.OrgKey &&
		len(pf.DomainKeys) == len(o.DomainKeys) &&
		len(pf.ListenerKeys) == len(o.ListenerKeys)) {
		return false
	}

	m := make(map[string]bool)
	for _, e := range pf.DomainKeys {
		m[string(e)] = true
	}
	for _, e := range o.DomainKeys {
		if !m[string(e)] {
			return false
		}
	}

	ml := make(map[string]bool)
	for _, e := range pf.ListenerKeys {
		ml[string(e)] = true
	}
	for _, e := range o.ListenerKeys {
		if !ml[string(e)] {
			return false
		}
	}

	return true
}

// Proxy describes the CRUD interface for api.Proxy
type Proxy interface {
	// GET /v1.0/proxies
	//
	// Index returns all Proxies to which the given filters apply. All non-empty
	// fields in a filter must apply for the filter to apply. Any Proxy to which
	// any filter applies is included in the result.
	//
	// If no filters are supplied, all Proxies are returned.
	Index(filters ...ProxyFilter) (api.Proxies, error)

	// GET /v1.0/proxies/<string:proxyKey>[?include_deleted]
	//
	// Get returns a Proxy for the given ProxyKey. If the Proxy does not
	// exist, an error is returned.
	Get(proxyKey api.ProxyKey) (api.Proxy, error)

	// POST /v1.0/proxies
	//
	// Create creates the given Proxy. The tuple of (Host, Port, ZoneKey) must be
	// unique. If a ProxyKey is specified in the Proxy, it is ignored and
	// replaced in the result with the authoritative ProxyKey.
	Create(proxy api.Proxy) (api.Proxy, error)

	// PUT /v1.0/proxies/<string:proxyKey>
	//
	// Modify Modifies the given Proxy. The tuple of (Host, Port, ZoneKey) must be
	// unique. The given Proxy Checksum must match the existing Checksum.
	Modify(proxy api.Proxy) (api.Proxy, error)

	// DELETE /v1.0/proxies/<string:proxyKey>?checksum=<string:checksum>
	//
	// Delete completely removes the Proxy from the database.
	// If the checksum does not match no action is taken and an error
	// is returned.
	Delete(proxyKey api.ProxyKey, checksum api.Checksum) error
}

// ListenerFilter describes a filter on the full list of Listeners
// per #5629 we should add the ability to find Proxies a Listener is attached to
type ListenerFilter struct {
	ListenerKey api.ListenerKey `json:"listener_key"`
	Name        string          `json:"name"`
	// DomainKeys matches Listeners with a superset of the specified DomainKeys. A
	// slice with a single value of "-" will produce Listeners with no linked
	// Domains.
	DomainKeys []api.DomainKey `json:"domain_keys"`
	ZoneKey    api.ZoneKey     `json:"zone_key"`
	OrgKey     api.OrgKey      `json:"org_key"`
}

// HasNoDomains returns true if DomainKeys has been set to the monitor
// value indicating a filter for Listeners with no linked Domains.
func (lf ListenerFilter) HasNoDomains() bool {
	return len(lf.DomainKeys) == 1 && lf.DomainKeys[0] == None
}

// IsNil returns true if the receiver is the zero value
func (lf ListenerFilter) IsNil() bool {
	return lf.Equals(ListenerFilter{})
}

// Equals returns true if the target is equal to the receiver
func (lf ListenerFilter) Equals(o ListenerFilter) bool {
	if !(lf.ListenerKey == o.ListenerKey &&
		lf.Name == o.Name &&
		lf.ZoneKey == o.ZoneKey &&
		lf.OrgKey == o.OrgKey &&
		len(lf.DomainKeys) == len(o.DomainKeys)) {
		return false
	}

	m := make(map[string]bool)
	for _, e := range lf.DomainKeys {
		m[string(e)] = true
	}
	for _, e := range o.DomainKeys {
		if !m[string(e)] {
			return false
		}
	}

	return true
}

// Listener describes the CRUD interface for api.Listener
type Listener interface {
	// GET /v1.0/listeners
	//
	// Index returns all listeners to which the given filters apply. All non-empty
	// fields in a filter must apply for the filter to apply. Any Listener to which
	// any filter applies is included in the result.
	//
	// If no filters are supplied, all Listeners are returned.
	Index(filters ...ListenerFilter) (api.Listeners, error)

	// GET /v1.0/listeners/<string:listenerKey>[?include_deleted]
	//
	// Get returns a Listener for the given ListenerKey. If the Listener does not
	// exist, an error is returned.
	Get(listenerKey api.ListenerKey) (api.Listener, error)

	// POST /v1.0/listeners
	//
	// Create creates the given Listener. The tuple of (Host, Port, ZoneKey) must be
	// unique. If a ListenerKey is specified in the Listener, it is ignored and
	// replaced in the result with the authoritative ListenerKey.
	Create(listener api.Listener) (api.Listener, error)

	// PUT /v1.0/listeners/<string:listenerKey>
	//
	// Modify Modifies the given Listener. The tuple of (Host, Port, ZoneKey) must be
	// unique. The given Listener Checksum must match the existing Checksum.
	Modify(listener api.Listener) (api.Listener, error)

	// DELETE /v1.0/listeners/<string:listenerKey>?checksum=<string:checksum>
	//
	// Delete completely removes the Listener from the database.
	// If the checksum does not match no action is taken and an error
	// is returned.
	Delete(listenerKey api.ListenerKey, checksum api.Checksum) error
}

// SharedRulesFilter describes a filter on the full list of SharedRules
type SharedRulesFilter struct {
	SharedRulesKey api.SharedRulesKey `json:"shared_rules_key"`
	Name           string             `json:"name"`
	ZoneKey        api.ZoneKey        `json:"zone_key"`
	OrgKey         api.OrgKey         `json:"org_key"`
}

// IsNil returns true if the receiver is the zero value
func (rf SharedRulesFilter) IsNil() bool {
	return rf.Equals(SharedRulesFilter{})
}

// Equals returns true if the target is equal to the receiver
func (rf SharedRulesFilter) Equals(o SharedRulesFilter) bool {
	return rf == o
}

// SharedRules describes the CRUD interface for api.SharedRules
type SharedRules interface {
	// GET /v1.0/shared_rules/<string:sharedRulesKey>
	//
	// Index returns all SharedRules to which the given filters apply. All non-empty
	// fields in a filter must apply for the filter to apply. Any SharedRules to which
	// any filter applies is included in the result.
	//
	// If no filters are supplied, all SharedRules are returned.
	Index(filters ...SharedRulesFilter) (api.SharedRulesSlice, error)

	// GET /v1.0/shared_rules/<string:sharedRulesKey>[?include_deleted]
	//
	// Get returns a SharedRules for the given SharedRulesKey. If the SharedRules does not
	// exist, an error is returned.
	Get(sharedRulesKey api.SharedRulesKey) (api.SharedRules, error)

	// POST /v1.0/shared_rules
	//
	// Create creates the given SharedRules. The Path must be unique for a given
	// ZoneKey. If a SharedRulesKey is specified in the SharedRules, it is ignored
	// and replaced in the result with the authoritative SharedRulesKey.
	Create(route api.SharedRules) (api.SharedRules, error)

	// PUT /v1.0/shared_rules/string:sharedRulesKey>
	//
	// Modify modifies the given SharedRules. The Path must be unique for a given
	// ZoneKey. The given SharedRules Checksum must match the existing Checksum.
	Modify(route api.SharedRules) (api.SharedRules, error)

	// DELETE /v1.0/shared_rules/<string:sharedRulesKey>?checksum=<string:checksum>
	//
	// Delete completely removes the SharedRules from the database.
	// If the checksum does not match no action is taken and an error
	// is returned.
	Delete(sharedRulesKey api.SharedRulesKey, checksum api.Checksum) error
}

// RouteFilter describes a filter on the full list of Routes
type RouteFilter struct {
	RouteKey       api.RouteKey       `json:"route_key"`
	DomainKey      api.DomainKey      `json:"domain_key"`
	SharedRulesKey api.SharedRulesKey `json:"shared_rules_key"`
	Path           string             `json:"path"`
	PathPrefix     string             `json:"path_prefix"`
	ZoneKey        api.ZoneKey        `json:"zone_key"`
	OrgKey         api.OrgKey         `json:"org_key"`
}

// IsNil returns true if the receiver is the zero value
func (rf RouteFilter) IsNil() bool {
	return rf.Equals(RouteFilter{})
}

// Equals returns true if the target is equal to the receiver
func (rf RouteFilter) Equals(o RouteFilter) bool {
	return rf == o
}

// Route describes the CRUD interface for api.Routes
type Route interface {
	// GET /v1.0/routes/<string:routeKey>
	//
	// Index returns all Routes to which the given filters apply. All non-empty
	// fields in a filter must apply for the filter to apply. Any Route to which
	// any filter applies is included in the result.
	//
	// If no filters are supplied, all Routes are returned.
	Index(filters ...RouteFilter) (api.Routes, error)

	// GET /v1.0/routes/<string:routeKey>[?include_deleted]
	//
	// Get returns a Route for the given RouteKey. If the Route does not
	// exist, an error is returned.
	Get(routeKey api.RouteKey) (api.Route, error)

	// POST /v1.0/routes
	//
	// Create creates the given Route. The Path must be unique for a given
	// ZoneKey.  If a RouteKey is specified in the Route, it is ignored
	// and replaced in the result with the authoritative RouteKey.
	Create(route api.Route) (api.Route, error)

	// PUT /v1.0/routes/string:routeKey>
	//
	// Modify modifies the given Route. The Path must be unique for a given
	// ZoneKey. The given Route Checksum must match the existing Checksum.
	Modify(route api.Route) (api.Route, error)

	// DELETE /v1.0/routes/<string:routeKey>?checksum=<string:checksum>
	//
	// Delete completely removes the Route from the database.
	// If the checksum does not match no action is taken and an error
	// is returned.
	Delete(routeKey api.RouteKey, checksum api.Checksum) error
}

// Zone describes the CRUD interface for api.Zones
type Zone interface {
	// GET /v1.0/admin/zone
	//
	// Index returns all Zones to which the given filters apply. All non-empty
	// fields in a filter must apply for the filter to apply. Any Zone to which
	// any filter applies is included in the result.
	//
	// If no filters are supplied, all Zones are returned.
	Index(filters ...ZoneFilter) (api.Zones, error)

	// GET /v1.0/admin/zone/<string:zoneKey>[?include_deleted]
	//
	// Get returns a Zone for the given ZoneKey. If the Zone does not
	// exist, an error is returned.
	Get(zoneKey api.ZoneKey) (api.Zone, error)

	// POST /v1.0/admin/zone
	//
	// Create creates the given Zone. Zone Names must be unique for a given
	// ZoneKey. If a ZoneKey is specified in the Zone, it is ignored and
	// replaced in the result with the authoritative ZoneKey.
	Create(zone api.Zone) (api.Zone, error)

	// PUT /v1.0/admin/zone/<string:zoneKey>
	//
	// Modify modifies the given Zone. Zone Names must be unique for a given
	// ZoneKey. The given Zone Checksum must match the existing Checksum.
	Modify(zone api.Zone) (api.Zone, error)

	// DELETE /v1.0/admin/zone/<string:zoneKey>?checksum=<string:checksum>
	//
	// Delete completely removes the Zone data from the database.
	// If the checksum does not match no action is taken and an error
	// is returned.
	Delete(zoneKey api.ZoneKey, checksum api.Checksum) error
}

// ZoneFilter describes a filter on the full list of Zones
type ZoneFilter struct {
	ZoneKey api.ZoneKey `json:"zone_key"`
	Name    string      `json:"name"`
	OrgKey  api.OrgKey  `json:"org_key"`
}

// IsNil returns true if the receiver is the zero value
func (z ZoneFilter) IsNil() bool {
	return z.Equals(ZoneFilter{})
}

// Equals returns true if the target is equal to the receiver
func (z ZoneFilter) Equals(o ZoneFilter) bool {
	return z == o
}

// AccessTokenFilter describes a filter on the full list of AccessTokens
type AccessTokenFilter struct {
	Description    string             `json:"description"`
	AccessTokenKey api.AccessTokenKey `json:"access_token_key"`
	UserKey        api.UserKey        `json:"user_key"`
	OrgKey         api.OrgKey         `json:"org_key"`
	CreatedAfter   *time.Time         `json:"created_after"`
	CreatedBefore  *time.Time         `json:"created_before"`
}

// IsNil returns true if the receiver is the zero value
func (of AccessTokenFilter) IsNil() bool {
	return of.Equals(AccessTokenFilter{})
}

// Equals returns true if the target is equal to the receiver
func (of AccessTokenFilter) Equals(o AccessTokenFilter) bool {
	return of.Description == o.Description &&
		of.AccessTokenKey == o.AccessTokenKey &&
		of.UserKey == o.UserKey &&
		of.OrgKey == o.OrgKey &&
		tbntime.Equal(of.CreatedAfter, o.CreatedAfter) &&
		tbntime.Equal(of.CreatedBefore, o.CreatedBefore)
}
