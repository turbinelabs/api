package service

//go:generate mockgen -source $GOFILE -destination mock_$GOFILE -package $GOPACKAGE

import (
	"strings"

	"github.com/turbinelabs/api"
)

/*
	All defines the interface for the public JSON/REST API.

	We define it in Go because it's convenient to allow our
	Go clients and servers to share a common language-level
	interface.

	Where necessary, we add commentary to describe things about
	the REST api that aren't documented by the language itself
	(eg paths, methods, etc)

	Each of the sub-interfaces presents an organization-scoped view.
	The method of scoping is not specified here, but in HTTP implementations
	will be the use of an authorization header.
*/
type All interface {
	Cluster() Cluster
	Domain() Domain
	SharedRules() SharedRules
	Route() Route
	Proxy() Proxy
	Zone() Zone
	History() History
}

type ClusterFilter struct {
	ClusterKey api.ClusterKey `json:"cluster_key"`
	Name       string         `json:"name"`
	ZoneKey    api.ZoneKey    `json:"zone_key"`
	OrgKey     api.OrgKey     `json:"org_key"`
}

func (cf ClusterFilter) Matches(c api.Cluster) bool {
	var (
		keyMatch  = cf.ClusterKey == "" || cf.ClusterKey == c.ClusterKey
		nameMatch = cf.Name == "" || cf.Name == c.Name
		zoneMatch = cf.ZoneKey == "" || cf.ZoneKey == c.ZoneKey
		orgMatch  = cf.OrgKey == "" || cf.OrgKey == c.OrgKey
	)

	return keyMatch && nameMatch && zoneMatch && orgMatch
}

func (cf ClusterFilter) IsNil() bool {
	return cf.Equals(ClusterFilter{})
}

func (cf ClusterFilter) Equals(o ClusterFilter) bool {
	return cf == o
}

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

type DomainFilter struct {
	DomainKey api.DomainKey `json:"domain_key"`
	Name      string        `json:"name"`
	ZoneKey   api.ZoneKey   `json:"zone_key"`
	OrgKey    api.OrgKey    `json:"org_key"`
}

func (df DomainFilter) Matches(d api.Domain) bool {
	var (
		keyMatches  = df.DomainKey == "" || df.DomainKey == d.DomainKey
		nameMatches = df.Name == "" || df.Name == d.Name
		zoneMatches = df.ZoneKey == "" || df.ZoneKey == d.ZoneKey
		orgMatches  = df.OrgKey == "" || df.OrgKey == d.OrgKey
	)

	return keyMatches && nameMatches && zoneMatches && orgMatches
}

func (df DomainFilter) IsNil() bool {
	return df.Equals(DomainFilter{})
}

func (df DomainFilter) Equals(o DomainFilter) bool {
	return df == o
}

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

type ProxyFilter struct {
	ProxyKey   api.ProxyKey    `json:"proxy_key"`
	Instance   api.Instance    `json:"instance"` // Matches Host/Port, ignores Metadata
	Name       string          `json:"name"`
	DomainKeys []api.DomainKey `json:"domain_keys"` // Matches Proxies with a superset of the specified DomainKeys
	ZoneKey    api.ZoneKey     `json:"zone_key"`
	OrgKey     api.OrgKey      `json:"org_key"`
}

func (pf ProxyFilter) Matches(p api.Proxy) bool {
	var (
		keyMatch  = pf.ProxyKey == "" || pf.ProxyKey == p.ProxyKey
		instMatch = pf.Instance.IsNil() || pf.Instance.Equivalent(p.Instance)
		nameMatch = pf.Name == "" || pf.Name == p.Name
		zoneMatch = pf.ZoneKey == "" || pf.ZoneKey == p.ZoneKey
		orgMatch  = pf.OrgKey == "" || pf.OrgKey == p.OrgKey
	)

	if !(keyMatch && instMatch && nameMatch && zoneMatch && orgMatch) {
		return false
	}

	if len(pf.DomainKeys) == 0 {
		return true
	}

	hasDomain := make(map[api.DomainKey]bool)
	for _, dk := range p.DomainKeys {
		hasDomain[dk] = true
	}

	for _, dk := range pf.DomainKeys {
		if !hasDomain[dk] {
			return false
		}
	}

	return true
}

func (pf ProxyFilter) IsNil() bool {
	return pf.Equals(ProxyFilter{})
}

func (pf ProxyFilter) Equals(o ProxyFilter) bool {
	if !(pf.ProxyKey == o.ProxyKey &&
		pf.Instance.Equals(o.Instance) &&
		pf.Name == o.Name &&
		pf.ZoneKey == o.ZoneKey &&
		pf.OrgKey == o.OrgKey &&
		len(pf.DomainKeys) == len(o.DomainKeys)) {
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

	return true
}

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

type SharedRulesFilter struct {
	SharedRulesKey api.SharedRulesKey `json:"shared_rules_key"`
	Name           string             `json:"name"`
	ZoneKey        api.ZoneKey        `json:"zone_key"`
	OrgKey         api.OrgKey         `json:"org_key"`
}

func (rf SharedRulesFilter) Matches(r api.SharedRules) bool {
	var (
		eqName = rf.Name == "" || rf.Name == r.Name
		eqKey  = rf.SharedRulesKey == "" || rf.SharedRulesKey == r.SharedRulesKey
		eqZone = rf.ZoneKey == "" || r.ZoneKey == rf.ZoneKey
		eqOrg  = rf.OrgKey == "" || r.OrgKey == rf.OrgKey
	)

	return eqKey && eqName && eqZone && eqOrg
}

func (rf SharedRulesFilter) IsNil() bool {
	return rf.Equals(SharedRulesFilter{})
}

func (rf SharedRulesFilter) Equals(o SharedRulesFilter) bool {
	return rf == o
}

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

type RouteFilter struct {
	RouteKey       api.RouteKey       `json:"route_key"`
	DomainKey      api.DomainKey      `json:"domain_key"`
	SharedRulesKey api.SharedRulesKey `json:"shared_rules_key"`
	Path           string             `json:"path"`
	PathPrefix     string             `json:"path_prefix"`
	ZoneKey        api.ZoneKey        `json:"zone_key"`
	OrgKey         api.OrgKey         `json:"org_key"`
}

func (rf RouteFilter) Matches(r api.Route) bool {
	var (
		eqKey    = rf.RouteKey == "" || rf.RouteKey == r.RouteKey
		eqDomain = rf.DomainKey == "" || r.DomainKey == rf.DomainKey
		eqSRK    = rf.SharedRulesKey == "" || r.SharedRulesKey == rf.SharedRulesKey
		eqPath   = rf.Path == "" || r.Path == rf.Path
		eqPrefix = rf.PathPrefix == "" || strings.HasPrefix(r.Path, rf.PathPrefix)
		eqZone   = rf.ZoneKey == "" || r.ZoneKey == rf.ZoneKey
		eqOrg    = rf.OrgKey == "" || r.OrgKey == rf.OrgKey
	)

	return eqKey && eqDomain && eqSRK && eqPath && eqPrefix && eqZone && eqOrg
}

func (rf RouteFilter) IsNil() bool {
	return rf.Equals(RouteFilter{})
}

func (rf RouteFilter) Equals(o RouteFilter) bool {
	return rf == o
}

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

type ZoneFilter struct {
	ZoneKey api.ZoneKey `json:"zone_key"`
	Name    string      `json:"name"`
	OrgKey  api.OrgKey  `json:"org_key"`
}

func (z ZoneFilter) IsNil() bool {
	return z.Equals(ZoneFilter{})
}

func (z ZoneFilter) Equals(o ZoneFilter) bool {
	return z == o
}
