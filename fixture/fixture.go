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

package fixture

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/turbinelabs/api"
)

// Map() exposes a copy of the fixture test data in a map fromat. As this
// map is a copy it's fine to mutate and it won't update the original fixtures.
func Map() map[string]string {
	data := make(map[string]string)
	for k, v := range initialTestData {
		data[k] = v
	}

	return data
}

// struct containing fixture data
type DataFixturesT struct {
	APIKey        string     // an api key that has been configured to point to ValidOrgID
	InvalidAPIKey string     // an api key that isn't in the store
	ValidOrgID    api.OrgKey // a valid org ID
	InvalidOrgID  api.OrgKey // an invalid org ID

	ZoneKey1        api.ZoneKey
	ZoneName1       string
	ZoneOrgKey1     api.OrgKey
	ZoneChecksum1   api.Checksum
	ZoneKey2        api.ZoneKey
	ZoneName2       string
	ZoneOrgKey2     api.OrgKey
	ZoneChecksum2   api.Checksum
	Zone1           api.Zone
	Zone2           api.Zone
	ZoneSlice       api.Zones
	PublicZoneSlice api.Zones

	OrgKey1          api.OrgKey
	OrgName1         string
	OrgContactEmail1 string
	OrgChecksum1     api.Checksum
	OrgKey2          api.OrgKey
	OrgName2         string
	OrgContactEmail2 string
	OrgChecksum2     api.Checksum
	Org1             api.Org
	Org2             api.Org
	OrgSlice         api.Orgs

	UserKey1        api.UserKey
	UserLoginEmail1 string
	UserAPIAuthKey1 api.APIAuthKey
	UserDeletedAt1  *time.Time
	UserChecksum1   api.Checksum
	UserOrgKey1     api.OrgKey
	UserKey2        api.UserKey
	UserLoginEmail2 string
	UserAPIAuthKey2 api.APIAuthKey
	UserOrgKey2     api.OrgKey
	UserDeletedAt2  *time.Time
	UserChecksum2   api.Checksum
	User1           api.User
	User2           api.User
	UserSlice       api.Users
	PublicUserSlice api.Users

	ClusterKey1        api.ClusterKey // UUID of cluster 1
	ClusterZone1       api.ZoneKey    // zone key for cluster 1
	ClusterName1       string         // name of cluster 1
	ClusterChecksum1   api.Checksum   // the checksum for cluster 1
	ClusterOrgKey1     api.OrgKey
	ClusterKey2        api.ClusterKey // UUId of cluster 2
	ClusterZone2       api.ZoneKey    // zone key for cluster 2
	ClusterName2       string         // name of cluster 2
	ClusterChecksum2   api.Checksum   // the checksum for cluster 2
	ClusterOrgKey2     api.OrgKey
	Cluster1           api.Cluster  // instance of cluster 1
	Cluster2           api.Cluster  // instance of cluster 1
	Instance21         api.Instance // first instance on cluster 2
	Instance22         api.Instance // first instance on cluster 2
	ClusterSlice       api.Clusters // slice of the two clusters
	PublicClusterSlice api.Clusters

	DomainKey1        api.DomainKey // UUID of domain 1
	DomainZone1       api.ZoneKey   // zone of domain 1
	DomainName1       string        // name of domain 1
	DomainPort1       int           // port of domain 1
	DomainChecksum1   api.Checksum  // checks for domain 1
	DomainOrgKey1     api.OrgKey
	DomainKey2        api.DomainKey // UUID of domain 2
	DomainName2       string        // name of domain 2
	DomainZone2       api.ZoneKey   // zone of domain 2
	DomainPort2       int           // port of domain 2
	DomainOrgKey2     api.OrgKey
	DomainChecksum2   api.Checksum // checks for domain 2
	Domain1           api.Domain   // domain 1
	Domain2           api.Domain   // domain 2
	DomainSlice       api.Domains  // slice of the two domains
	PublicDomainSlice api.Domains

	ProxyKey1        api.ProxyKey
	ProxyZone1       api.ZoneKey
	ProxyMetadata1   api.Metadata
	ProxyInstance1   api.Instance
	ProxyName1       string
	ProxyDomainKeys1 []api.DomainKey
	ProxyChecksum1   api.Checksum
	ProxyOrgKey1     api.OrgKey
	ProxyKey2        api.ProxyKey
	ProxyZone2       api.ZoneKey
	ProxyMetadata2   api.Metadata
	ProxyInstance2   api.Instance
	ProxyName2       string
	ProxyDomainKeys2 []api.DomainKey
	ProxyChecksum2   api.Checksum
	ProxyOrgKey2     api.OrgKey
	ProxyDomain21    api.Domain
	ProxyDomain22    api.Domain
	Proxy1           api.Proxy
	Proxy2           api.Proxy
	ProxySlice       api.Proxies
	PublicProxySlice api.Proxies

	RouteKey1            api.RouteKey
	RouteDomain1         api.DomainKey
	RouteZone1           api.ZoneKey
	RoutePath1           string
	RouteSharedRulesKey1 api.SharedRulesKey
	RouteRules1          api.Rules
	RouteChecksum1       api.Checksum
	RouteOrgKey1         api.OrgKey
	RouteKey2            api.RouteKey
	RouteDomain2         api.DomainKey
	RouteZone2           api.ZoneKey
	RoutePath2           string
	RouteSharedRulesKey2 api.SharedRulesKey
	RouteRules2          api.Rules
	RouteChecksum2       api.Checksum
	RouteOrgKey2         api.OrgKey
	Route1               api.Route
	Route2               api.Route
	RouteSlice           api.Routes
	PublicRouteSlice     api.Routes

	SharedRulesKey1        api.SharedRulesKey
	SharedRulesName1       string
	SharedRulesZone1       api.ZoneKey
	SharedRulesDefault1    api.AllConstraints
	SharedRulesRules1      api.Rules
	SharedRulesChecksum1   api.Checksum
	SharedRulesOrgKey1     api.OrgKey
	SharedRulesKey2        api.SharedRulesKey
	SharedRulesName2       string
	SharedRulesZone2       api.ZoneKey
	SharedRulesDefault2    api.AllConstraints
	SharedRulesRules2      api.Rules
	SharedRulesChecksum2   api.Checksum
	SharedRulesOrgKey2     api.OrgKey
	SharedRules1           api.SharedRules
	SharedRules2           api.SharedRules
	SharedRulesSlice       api.SharedRulesSlice
	PublicSharedRulesSlice api.SharedRulesSlice
}

// Provides access to key data within the store; simple values are set here
// while complex values are constructed in init()
var DataFixtures DataFixturesT = DataFixturesT{
	APIKey:        "key-present",
	InvalidAPIKey: "key-missing",
	ValidOrgID:    "1",
	InvalidOrgID:  "nope",

	ZoneKey1:      api.ZoneKey("zone-1"),
	ZoneName1:     "us-west",
	ZoneOrgKey1:   api.OrgKey("1"),
	ZoneChecksum1: api.Checksum{"Z-CS-1"},
	ZoneKey2:      api.ZoneKey("zone-2"),
	ZoneName2:     "us-east",
	ZoneOrgKey2:   api.OrgKey("2"),
	ZoneChecksum2: api.Checksum{"Z-CS-2"},

	OrgKey1:          api.OrgKey("1"),
	OrgName1:         "ExampleCo",
	OrgContactEmail1: "adminco1@example.com",
	OrgChecksum1:     api.Checksum{"Org-CS-1"},
	OrgKey2:          api.OrgKey("2"),
	OrgName2:         "ExampleCo2",
	OrgContactEmail2: "adminco2@alt.example.com",
	OrgChecksum2:     api.Checksum{"Org-CS-2"},

	UserKey1:        api.UserKey("1"),
	UserLoginEmail1: "someuser@example.com",
	UserAPIAuthKey1: api.APIAuthKey("user-api-key-1"),
	UserOrgKey1:     "1",
	UserDeletedAt1:  nil,
	UserChecksum1:   api.Checksum{"user-cs-1"},

	UserKey2:        api.UserKey("2"),
	UserLoginEmail2: "otheruser@example.com",
	UserAPIAuthKey2: api.APIAuthKey("user-api-key-2"),
	UserOrgKey2:     "1",
	UserDeletedAt2:  nil,
	UserChecksum2:   api.Checksum{"user-cs-2"},

	ClusterKey1:      "98a13568-a599-4c8d-4ae8-657f3917e2cf",
	ClusterZone1:     "zk1",
	ClusterName1:     "cluster1",
	ClusterChecksum1: api.Checksum{"cluster-checksum-1"},
	ClusterOrgKey1:   "1",
	ClusterKey2:      "2794c958-d44c-418c-5cac-4d1af020df99",
	ClusterZone2:     "zk2",
	ClusterName2:     "cluster2",
	ClusterOrgKey2:   "1",
	ClusterChecksum2: api.Checksum{"cluster-checksum-2"},

	DomainKey1:      "asonetuhasonetuh",
	DomainZone1:     "zk1",
	DomainName1:     "domain-1",
	DomainPort1:     8080,
	DomainChecksum1: api.Checksum{"ck1"},
	DomainOrgKey1:   "1",
	DomainKey2:      "sntaohesntahoesuntaohe",
	DomainZone2:     "zk2",
	DomainName2:     "domain-2",
	DomainPort2:     5050,
	DomainOrgKey2:   "1",
	DomainChecksum2: api.Checksum{"ck2"},

	ProxyKey1:      "proxy-1",
	ProxyZone1:     "proxy-zone-1",
	ProxyName1:     "proxy-name-1",
	ProxyChecksum1: api.Checksum{"proxy-cs-1"},
	ProxyInstance1: api.Instance{"proxy-host-1", 8085, api.Metadata{{"key", "value"}}},
	ProxyOrgKey1:   "1",
	ProxyKey2:      "proxy-2",
	ProxyZone2:     "proxy-zone-2",
	ProxyName2:     "proxy-name-2",
	ProxyOrgKey2:   "1",
	ProxyInstance2: api.Instance{"proxy-host-2", 8085, api.Metadata{}},
	ProxyChecksum2: api.Checksum{"proxy-cs-2"},

	RouteKey1:      "route-key-1",
	RouteDomain1:   "route-dom-1",
	RouteZone1:     "route-zone-1",
	RoutePath1:     "for/bar/path",
	RouteOrgKey1:   "1",
	RouteChecksum1: api.Checksum{"route-cs-1"},
	RouteKey2:      "route-key-2",
	RouteDomain2:   "route-dom-2",
	RouteZone2:     "route-zone-2",
	RoutePath2:     "quix/qux/quuuuux",
	RouteOrgKey2:   "1",
	RouteChecksum2: api.Checksum{"route-cs-2"},

	SharedRulesKey1:      "shared-rules-key-1",
	SharedRulesName1:     "shared-rules-name-1",
	SharedRulesZone1:     "shared-rules-zone-1",
	SharedRulesOrgKey1:   "1",
	SharedRulesChecksum1: api.Checksum{"shared-rules-cs-1"},
	SharedRulesKey2:      "shared-rules-key-2",
	SharedRulesName2:     "shared-rules-name-2",
	SharedRulesZone2:     "shared-rules-zone-2",
	SharedRulesOrgKey2:   "1",
	SharedRulesChecksum2: api.Checksum{"shared-rules-cs-2"},
}

var initialTestData map[string]string = make(map[string]string)

func init() {
	DataFixtures.Org1 = api.Org{
		DataFixtures.OrgKey1,
		DataFixtures.OrgName1,
		DataFixtures.OrgContactEmail1,
		DataFixtures.OrgChecksum1,
	}
	DataFixtures.Org2 = api.Org{
		DataFixtures.OrgKey2,
		DataFixtures.OrgName2,
		DataFixtures.OrgContactEmail2,
		DataFixtures.OrgChecksum2,
	}
	DataFixtures.OrgSlice = api.Orgs{DataFixtures.Org1, DataFixtures.Org2}

	ts := time.Date(2015, 2, 28, 12, 30, 0, 0, time.UTC)
	DataFixtures.UserDeletedAt2 = &ts
	DataFixtures.Zone1 = api.Zone{
		DataFixtures.ZoneKey1,
		DataFixtures.ZoneName1,
		DataFixtures.ZoneOrgKey1,
		DataFixtures.ZoneChecksum1,
	}
	DataFixtures.Zone2 = api.Zone{
		DataFixtures.ZoneKey2,
		DataFixtures.ZoneName2,
		DataFixtures.ZoneOrgKey2,
		DataFixtures.ZoneChecksum2,
	}
	DataFixtures.ZoneSlice = api.Zones{DataFixtures.Zone1, DataFixtures.Zone2}
	DataFixtures.PublicZoneSlice = make(api.Zones, len(DataFixtures.ZoneSlice))
	for i, z := range DataFixtures.ZoneSlice {
		z.OrgKey = ""
		DataFixtures.PublicZoneSlice[i] = z
	}

	DataFixtures.User1 = api.User{
		UserKey:    DataFixtures.UserKey1,
		LoginEmail: DataFixtures.UserLoginEmail1,
		APIAuthKey: DataFixtures.UserAPIAuthKey1,
		OrgKey:     DataFixtures.UserOrgKey1,
		DeletedAt:  DataFixtures.UserDeletedAt1,
		Checksum:   DataFixtures.UserChecksum1,
	}
	DataFixtures.User2 = api.User{
		UserKey:    DataFixtures.UserKey2,
		LoginEmail: DataFixtures.UserLoginEmail2,
		APIAuthKey: DataFixtures.UserAPIAuthKey2,
		OrgKey:     DataFixtures.UserOrgKey2,
		DeletedAt:  DataFixtures.UserDeletedAt2,
		Checksum:   DataFixtures.UserChecksum2,
	}
	DataFixtures.UserSlice = api.Users{DataFixtures.User1, DataFixtures.User2}
	DataFixtures.PublicUserSlice = make(api.Users, len(DataFixtures.UserSlice))
	for i, u := range DataFixtures.UserSlice {
		DataFixtures.PublicUserSlice[i] = u
	}

	DataFixtures.Cluster1 = api.Cluster{
		ClusterKey: DataFixtures.ClusterKey1,
		ZoneKey:    DataFixtures.ClusterZone1,
		Name:       DataFixtures.ClusterName1,
		OrgKey:     DataFixtures.ClusterOrgKey1,
		Checksum:   DataFixtures.ClusterChecksum1,
	}

	DataFixtures.Instance21 = api.Instance{
		Host: "int-host", Port: 1234, Metadata: api.Metadata{{"key1", "value1"}, {"key2", "value2"}}}

	DataFixtures.Instance22 = api.Instance{Host: "int-host-2", Port: 1234}

	DataFixtures.Cluster2 = api.Cluster{
		ClusterKey: DataFixtures.ClusterKey2,
		ZoneKey:    DataFixtures.ClusterZone2,
		Name:       DataFixtures.ClusterName2,
		Instances:  api.Instances{DataFixtures.Instance21, DataFixtures.Instance22},
		OrgKey:     DataFixtures.ClusterOrgKey2,
		Checksum:   DataFixtures.ClusterChecksum2,
	}

	DataFixtures.ClusterSlice = []api.Cluster{DataFixtures.Cluster1, DataFixtures.Cluster2}
	DataFixtures.PublicClusterSlice = make(api.Clusters, len(DataFixtures.ClusterSlice))
	for i, c := range DataFixtures.ClusterSlice {
		c.OrgKey = ""
		DataFixtures.PublicClusterSlice[i] = c
	}

	// domain setup
	DataFixtures.Domain1 = api.Domain{
		DataFixtures.DomainKey1,
		DataFixtures.DomainZone1,
		DataFixtures.DomainName1,
		DataFixtures.DomainPort1,
		DataFixtures.DomainOrgKey1,
		DataFixtures.DomainChecksum1,
	}

	DataFixtures.Domain2 = api.Domain{
		DataFixtures.DomainKey2,
		DataFixtures.DomainZone2,
		DataFixtures.DomainName2,
		DataFixtures.DomainPort2,
		DataFixtures.DomainOrgKey2,
		DataFixtures.DomainChecksum2,
	}

	DataFixtures.DomainSlice = api.Domains{DataFixtures.Domain1, DataFixtures.Domain2}
	DataFixtures.PublicDomainSlice = make(api.Domains, len(DataFixtures.DomainSlice))
	for i, d := range DataFixtures.DomainSlice {
		d.OrgKey = ""
		DataFixtures.PublicDomainSlice[i] = d
	}

	// proxy setup
	DataFixtures.ProxyDomainKeys1 = []api.DomainKey{
		DataFixtures.Domain1.DomainKey,
		DataFixtures.Domain2.DomainKey,
	}
	DataFixtures.Proxy1 = api.Proxy{
		DataFixtures.ProxyInstance1,
		DataFixtures.ProxyKey1,
		DataFixtures.ProxyZone1,
		DataFixtures.ProxyName1,
		DataFixtures.ProxyDomainKeys1,
		DataFixtures.ProxyOrgKey1,
		DataFixtures.ProxyChecksum1,
	}

	DataFixtures.ProxyDomainKeys2 = DataFixtures.ProxyDomainKeys1
	DataFixtures.Proxy2 = api.Proxy{
		DataFixtures.ProxyInstance2,
		DataFixtures.ProxyKey2,
		DataFixtures.ProxyZone2,
		DataFixtures.ProxyName2,
		DataFixtures.ProxyDomainKeys2,
		DataFixtures.ProxyOrgKey2,
		DataFixtures.ProxyChecksum2,
	}

	DataFixtures.ProxySlice = []api.Proxy{DataFixtures.Proxy1, DataFixtures.Proxy2}
	DataFixtures.PublicProxySlice = make(api.Proxies, len(DataFixtures.ProxySlice))
	for i, p := range DataFixtures.ProxySlice {
		p.OrgKey = ""
		DataFixtures.PublicProxySlice[i] = p
	}

	// route setup
	routeRule1 := api.Rule{
		"rk-1-0",
		[]string{"GET", "POST"},
		api.Matches{
			api.Match{api.HeaderMatchKind, api.Metadatum{"x-1", "value"}, api.Metadatum{"flag", "true"}},
			api.Match{api.CookieMatchKind, api.Metadatum{"x-2", ""}, api.Metadatum{"other", "true"}}},
		api.AllConstraints{
			Light: api.ClusterConstraints{
				{"cc-0", "ckey2", api.Metadata{{"key-2", "value-2"}}, api.Metadata{{"state", "test"}}, 1234}}},
	}

	routeRule2 := api.Rule{
		"rk-0-1",
		[]string{"PUT", "DELETE"},
		api.Matches{
			api.Match{api.CookieMatchKind, api.Metadatum{"x-2", "value"}, api.Metadatum{"other", "true"}}},
		api.AllConstraints{
			Tap: api.ClusterConstraints{
				{"cc-1", "ckey3", api.Metadata{{"key-2", "value-2"}}, api.Metadata{}, 1234}},
			Light: api.ClusterConstraints{
				{"cc-2", "ckey2", api.Metadata{{"key-2", "value-2"}}, api.Metadata{}, 1234}}},
	}

	DataFixtures.RouteRules1 = api.Rules{routeRule1}
	DataFixtures.RouteRules2 = api.Rules{routeRule1, routeRule2}
	DataFixtures.Route1 = api.Route{
		DataFixtures.RouteKey1,
		DataFixtures.RouteDomain1,
		DataFixtures.RouteZone1,
		DataFixtures.RoutePath1,
		DataFixtures.SharedRulesKey1,
		DataFixtures.RouteRules1,
		DataFixtures.RouteOrgKey1,
		DataFixtures.RouteChecksum1,
	}

	DataFixtures.Route2 = api.Route{
		DataFixtures.RouteKey2,
		DataFixtures.RouteDomain2,
		DataFixtures.RouteZone2,
		DataFixtures.RoutePath2,
		DataFixtures.SharedRulesKey2,
		DataFixtures.RouteRules2,
		DataFixtures.RouteOrgKey2,
		DataFixtures.RouteChecksum2,
	}

	DataFixtures.RouteSlice = api.Routes{DataFixtures.Route1, DataFixtures.Route2}
	DataFixtures.PublicRouteSlice = make(api.Routes, len(DataFixtures.RouteSlice))
	for i, r := range DataFixtures.RouteSlice {
		r.OrgKey = ""
		DataFixtures.PublicRouteSlice[i] = r
	}

	// sharedRules setup
	sharedRulesRule1 := api.Rule{
		"srk-1-0",
		[]string{"GET", "POST"},
		api.Matches{
			api.Match{api.HeaderMatchKind, api.Metadatum{"x-1", "value"}, api.Metadatum{"flag", "true"}},
			api.Match{api.CookieMatchKind, api.Metadatum{"x-2", ""}, api.Metadatum{"other", "true"}}},
		api.AllConstraints{
			Light: api.ClusterConstraints{
				{"cc-0", "ckey2", api.Metadata{{"key-2", "value-2"}}, api.Metadata{{"state", "test"}}, 1234}}},
	}

	sharedRulesRule2 := api.Rule{
		"srk-0-1",
		[]string{"PUT", "DELETE"},
		api.Matches{
			api.Match{api.CookieMatchKind, api.Metadatum{"x-2", "value"}, api.Metadatum{"other", "true"}}},
		api.AllConstraints{
			Tap: api.ClusterConstraints{
				{"cc-1", "ckey3", api.Metadata{{"key-2", "value-2"}}, api.Metadata{}, 1234}},
			Light: api.ClusterConstraints{
				{"cc-2", "ckey2", api.Metadata{{"key-2", "value-2"}}, api.Metadata{}, 1234}}},
	}

	sharedRulesDefault1 := api.AllConstraints{
		Light: api.ClusterConstraints{
			{"cc-4", api.HeaderMatchKind, api.Metadata{{"k", "v"}, {"k2", "v2"}}, api.Metadata{{"state", "released"}}, 23}}}
	sharedRulesDefault2 := sharedRulesDefault1

	DataFixtures.SharedRulesDefault1 = sharedRulesDefault1
	DataFixtures.SharedRulesRules1 = api.Rules{sharedRulesRule1}
	DataFixtures.SharedRulesDefault2 = sharedRulesDefault2
	DataFixtures.SharedRulesRules2 = api.Rules{sharedRulesRule1, sharedRulesRule2}
	DataFixtures.SharedRules1 = api.SharedRules{
		DataFixtures.SharedRulesKey1,
		DataFixtures.SharedRulesName1,
		DataFixtures.SharedRulesZone1,
		DataFixtures.SharedRulesDefault1,
		DataFixtures.SharedRulesRules1,
		DataFixtures.SharedRulesOrgKey1,
		DataFixtures.SharedRulesChecksum1,
	}

	DataFixtures.SharedRules2 = api.SharedRules{
		DataFixtures.SharedRulesKey2,
		DataFixtures.SharedRulesName2,
		DataFixtures.SharedRulesZone2,
		DataFixtures.SharedRulesDefault2,
		DataFixtures.SharedRulesRules2,
		DataFixtures.SharedRulesOrgKey2,
		DataFixtures.SharedRulesChecksum2,
	}

	DataFixtures.SharedRulesSlice = api.SharedRulesSlice{DataFixtures.SharedRules1, DataFixtures.SharedRules2}
	DataFixtures.PublicSharedRulesSlice = make(api.SharedRulesSlice, len(DataFixtures.SharedRulesSlice))
	for i, r := range DataFixtures.SharedRulesSlice {
		r.OrgKey = ""
		DataFixtures.PublicSharedRulesSlice[i] = r
	}

	// install api key
	initialTestData["/tbn/api-keys/"+DataFixtures.APIKey] = string(DataFixtures.ValidOrgID)

	// install clusters
	c1Encoded, err := json.Marshal(DataFixtures.Cluster1)
	if err != nil {
		log.Fatal(err)
	}
	key1 := fmt.Sprintf("/tbn/api/%s/cluster/%s", DataFixtures.ValidOrgID, string(DataFixtures.ClusterKey1))
	initialTestData[key1] = string(c1Encoded)

	c2Encoded, err := json.Marshal(DataFixtures.Cluster2)
	if err != nil {
		log.Fatal(err)
	}
	key2 := fmt.Sprintf("/tbn/api/%s/cluster/%s", DataFixtures.ValidOrgID, string(DataFixtures.ClusterKey2))
	initialTestData[key2] = string(c2Encoded)

	// install domains
	d1Encoded, err := json.Marshal(DataFixtures.Domain1)
	if err != nil {
		log.Fatal(err)
	}
	dkey1 := fmt.Sprintf("/tbn/api/%s/domain/%s", DataFixtures.ValidOrgID, string(DataFixtures.DomainKey1))
	initialTestData[dkey1] = string(d1Encoded)

	d2Encoded, err := json.Marshal(DataFixtures.Domain2)
	if err != nil {
		log.Fatal(err)
	}
	dkey2 := fmt.Sprintf("/tbn/api/%s/domain/%s", DataFixtures.ValidOrgID, string(DataFixtures.DomainKey2))
	initialTestData[dkey2] = string(d2Encoded)

	// install proxies
	p1Encoded, err := json.Marshal(DataFixtures.Proxy1)
	if err != nil {
		log.Fatal(err)
	}
	pkey1 := fmt.Sprintf("/tbn/api/%s/proxy/%s", DataFixtures.ValidOrgID, string(DataFixtures.ProxyKey1))
	initialTestData[pkey1] = string(p1Encoded)

	p2Encoded, err := json.Marshal(DataFixtures.Proxy2)
	if err != nil {
		log.Fatal(err)
	}
	pkey2 := fmt.Sprintf("/tbn/api/%s/proxy/%s", DataFixtures.ValidOrgID, string(DataFixtures.ProxyKey2))
	initialTestData[pkey2] = string(p2Encoded)

	// install routes
	r1Encoded, err := json.Marshal(DataFixtures.Route1)
	if err != nil {
		log.Fatal(err)
	}
	rkey1 := fmt.Sprintf("/tbn/api/%s/route/%s", DataFixtures.ValidOrgID, string(DataFixtures.RouteKey1))
	initialTestData[rkey1] = string(r1Encoded)

	r2Encoded, err := json.Marshal(DataFixtures.Route2)
	if err != nil {
		log.Fatal(err)
	}
	rkey2 := fmt.Sprintf("/tbn/api/%s/route/%s", DataFixtures.ValidOrgID, string(DataFixtures.RouteKey2))
	initialTestData[rkey2] = string(r2Encoded)
}
