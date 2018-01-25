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

// Package fixture provides a set of fixtures for use in testing that requires
// a fairly sane universe of api objects.
package fixture

import (
	"time"

	"github.com/turbinelabs/api"
)

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
	ClusterRequireTLS1 bool           // should cluster 1 require TLS communications
	ClusterChecksum1   api.Checksum   // the checksum for cluster 1
	ClusterOrgKey1     api.OrgKey
	ClusterKey2        api.ClusterKey // UUId of cluster 2
	ClusterZone2       api.ZoneKey    // zone key for cluster 2
	ClusterName2       string         // name of cluster 2
	ClusterRequireTLS2 bool           // should cluster 2 require TLS communications
	ClusterChecksum2   api.Checksum   // the checksum for cluster 2
	ClusterOrgKey2     api.OrgKey
	Cluster1           api.Cluster  // instance of cluster 1
	Cluster2           api.Cluster  // instance of cluster 1
	Instance21         api.Instance // first instance on cluster 2
	Instance22         api.Instance // first instance on cluster 2
	ClusterSlice       api.Clusters // slice of the two clusters
	PublicClusterSlice api.Clusters

	DomainKey1         api.DomainKey // UUID of domain 1
	DomainZone1        api.ZoneKey   // zone of domain 1
	DomainName1        string        // name of domain 1
	DomainPort1        int           // port of domain 1
	DomainSSLConfig1   *api.SSLConfig
	DomainRedirects1   api.Redirects // part of domain 1
	DomainGzipEnabled1 bool          // part of domain 1
	DomainCorsConfig1  *api.CorsConfig
	DomainAliases1     api.DomainAliases
	DomainChecksum1    api.Checksum // checks for domain 1
	DomainOrgKey1      api.OrgKey
	DomainKey2         api.DomainKey // UUID of domain 2
	DomainName2        string        // name of domain 2
	DomainZone2        api.ZoneKey   // zone of domain 2
	DomainPort2        int           // port of domain 2
	DomainSSLConfig2   *api.SSLConfig
	DomainRedirects2   api.Redirects // part of domain 2
	DomainGzipEnabled2 bool          // part of domain 2
	DomainCorsConfig2  *api.CorsConfig
	DomainAliases2     api.DomainAliases
	DomainOrgKey2      api.OrgKey
	DomainChecksum2    api.Checksum // checks for domain 2
	Domain1            api.Domain   // domain 1
	Domain2            api.Domain   // domain 2
	DomainSlice        api.Domains  // slice of the two domains
	PublicDomainSlice  api.Domains

	ProxyKey1        api.ProxyKey
	ProxyZone1       api.ZoneKey
	ProxyMetadata1   api.Metadata
	ProxyName1       string
	ProxyDomainKeys1 []api.DomainKey
	ProxyChecksum1   api.Checksum
	ProxyOrgKey1     api.OrgKey
	ProxyKey2        api.ProxyKey
	ProxyZone2       api.ZoneKey
	ProxyMetadata2   api.Metadata
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
	RouteResponseData1   api.ResponseData
	RouteCohortSeed1     *api.CohortSeed
	RouteRetryPolicy1    *api.RetryPolicy
	RouteChecksum1       api.Checksum
	RouteOrgKey1         api.OrgKey
	RouteKey2            api.RouteKey
	RouteDomain2         api.DomainKey
	RouteZone2           api.ZoneKey
	RoutePath2           string
	RouteSharedRulesKey2 api.SharedRulesKey
	RouteRules2          api.Rules
	RouteResponseData2   api.ResponseData
	RouteCohortSeed2     *api.CohortSeed
	RouteRetryPolicy2    *api.RetryPolicy
	RouteChecksum2       api.Checksum
	RouteOrgKey2         api.OrgKey
	Route1               api.Route
	Route2               api.Route
	RouteSlice           api.Routes
	PublicRouteSlice     api.Routes

	SharedRulesKey1          api.SharedRulesKey
	SharedRulesName1         string
	SharedRulesZone1         api.ZoneKey
	SharedRulesDefault1      api.AllConstraints
	SharedRulesRules1        api.Rules
	SharedRulesResponseData1 api.ResponseData
	SharedRulesCohortSeed1   *api.CohortSeed
	SharedRulesProperties1   api.Metadata
	SharedRulesRetryPolicy1  *api.RetryPolicy
	SharedRulesChecksum1     api.Checksum
	SharedRulesOrgKey1       api.OrgKey
	SharedRulesKey2          api.SharedRulesKey
	SharedRulesName2         string
	SharedRulesZone2         api.ZoneKey
	SharedRulesDefault2      api.AllConstraints
	SharedRulesRules2        api.Rules
	SharedRulesResponseData2 api.ResponseData
	SharedRulesCohortSeed2   *api.CohortSeed
	SharedRulesProperties2   api.Metadata
	SharedRulesRetryPolicy2  *api.RetryPolicy
	SharedRulesChecksum2     api.Checksum
	SharedRulesOrgKey2       api.OrgKey
	SharedRules1             api.SharedRules
	SharedRules2             api.SharedRules
	SharedRulesSlice         api.SharedRulesSlice
	PublicSharedRulesSlice   api.SharedRulesSlice

	AccessToken1            api.AccessToken
	AccessTokenKey1         api.AccessTokenKey
	AccessTokenDescription1 string
	AccessTokenUserKey1     api.UserKey
	AccessTokenOrgKey1      api.OrgKey
	AccessTokenCreatedAt1   *time.Time
	AccessTokenChecksum1    api.Checksum

	AccessToken2            api.AccessToken
	AccessTokenKey2         api.AccessTokenKey
	AccessTokenDescription2 string
	AccessTokenUserKey2     api.UserKey
	AccessTokenOrgKey2      api.OrgKey
	AccessTokenCreatedAt2   *time.Time
	AccessTokenChecksum2    api.Checksum

	AccessTokenSlice api.AccessTokens
	AccessTokenZone2 api.ZoneKey
}

// Provides access to key data within the store; simple values are set here
// while complex values are constructed in init()
// TODO: convert this to a function producing a new DataFixturesT
func New() DataFixturesT {
	df := DataFixturesT{
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

		ClusterKey1:        "98a13568-a599-4c8d-4ae8-657f3917e2cf",
		ClusterZone1:       "zk1",
		ClusterName1:       "cluster1",
		ClusterRequireTLS1: false,
		ClusterChecksum1:   api.Checksum{"cluster-checksum-1"},
		ClusterOrgKey1:     "1",
		ClusterKey2:        "2794c958-d44c-418c-5cac-4d1af020df99",
		ClusterZone2:       "zk2",
		ClusterName2:       "cluster2",
		ClusterRequireTLS2: true,
		ClusterOrgKey2:     "1",
		ClusterChecksum2:   api.Checksum{"cluster-checksum-2"},

		DomainKey1:       "asonetuhasonetuh",
		DomainZone1:      "zk1",
		DomainName1:      "domain-1",
		DomainPort1:      8080,
		DomainSSLConfig1: nil,
		DomainRedirects1: api.Redirects{{
			"redirect1",
			".*",
			"http://www.example.com",
			api.PermanentRedirect,
			api.HeaderConstraints{{"x-random-header", "", false, true}},
		}},
		DomainGzipEnabled1: true,
		DomainCorsConfig1: &api.CorsConfig{
			AllowedOrigins:   []string{"*"},
			AllowCredentials: true,
			ExposedHeaders:   []string{"x-expose-1", "x-expose-2"},
			MaxAge:           600,
			AllowedMethods:   []string{"GET", "POST"},
			AllowedHeaders:   []string{"x-allowed-1", "x-allowed-2"},
		},
		DomainAliases1:     api.DomainAliases{"example.com", "*.example.com"},
		DomainChecksum1:    api.Checksum{"ck1"},
		DomainOrgKey1:      "1",
		DomainKey2:         "sntaohesntahoesuntaohe",
		DomainZone2:        "zk2",
		DomainName2:        "domain-2",
		DomainPort2:        5050,
		DomainSSLConfig2:   nil,
		DomainRedirects2:   nil,
		DomainGzipEnabled2: false,
		DomainCorsConfig2:  nil,
		DomainAliases2:     nil,
		DomainOrgKey2:      "1",
		DomainChecksum2:    api.Checksum{"ck2"},

		ProxyKey1:      "proxy-1",
		ProxyZone1:     "proxy-zone-1",
		ProxyName1:     "proxy-name-1",
		ProxyChecksum1: api.Checksum{"proxy-cs-1"},
		ProxyOrgKey1:   "1",
		ProxyKey2:      "proxy-2",
		ProxyZone2:     "proxy-zone-2",
		ProxyName2:     "proxy-name-2",
		ProxyOrgKey2:   "1",
		ProxyChecksum2: api.Checksum{"proxy-cs-2"},

		RouteKey1:         "route-key-1",
		RouteDomain1:      "route-dom-1",
		RouteZone1:        "route-zone-1",
		RoutePath1:        "for/bar/path",
		RouteCohortSeed1:  &api.CohortSeed{api.CohortSeedCookie, "cookie-cohort-data", false},
		RouteRetryPolicy1: nil,
		RouteOrgKey1:      "1",
		RouteChecksum1:    api.Checksum{"route-cs-1"},
		RouteKey2:         "route-key-2",
		RouteDomain2:      "route-dom-2",
		RouteZone2:        "route-zone-2",
		RoutePath2:        "quix/qux/quuuuux",
		RouteCohortSeed2:  nil,
		RouteRetryPolicy2: nil,
		RouteOrgKey2:      "1",
		RouteChecksum2:    api.Checksum{"route-cs-2"},

		SharedRulesKey1:         "shared-rules-key-1",
		SharedRulesName1:        "shared-rules-name-1",
		SharedRulesZone1:        "shared-rules-zone-1",
		SharedRulesCohortSeed1:  &api.CohortSeed{api.CohortSeedHeader, "x-cohort-data", true},
		SharedRulesProperties1:  api.Metadata{{"pk1", "pv1"}, {"pk12", "pv12"}},
		SharedRulesRetryPolicy1: nil,
		SharedRulesOrgKey1:      "1",
		SharedRulesChecksum1:    api.Checksum{"shared-rules-cs-1"},
		SharedRulesKey2:         "shared-rules-key-2",
		SharedRulesName2:        "shared-rules-name-2",
		SharedRulesZone2:        "shared-rules-zone-2",
		SharedRulesCohortSeed2:  nil,
		SharedRulesProperties2:  api.Metadata{{"pk2", "pv2"}, {"pk22", "pv22"}},
		SharedRulesRetryPolicy2: nil,
		SharedRulesOrgKey2:      "1",
		SharedRulesChecksum2:    api.Checksum{"shared-rules-cs-2"},

		AccessTokenKey1:         "access-token-key-1",
		AccessTokenDescription1: "access-token-descirption-1",
		AccessTokenUserKey1:     "access-token-user-key-1",
		AccessTokenOrgKey1:      "access-token-org-key-1",
		AccessTokenCreatedAt1:   nil,
		AccessTokenChecksum1:    api.Checksum{"access-token-cs-1"},

		AccessTokenKey2:         "access-token-key-2",
		AccessTokenDescription2: "access-token-descirption-2",
		AccessTokenUserKey2:     "access-token-user-key-2",
		AccessTokenOrgKey2:      "access-token-org-key-2",
		AccessTokenCreatedAt2:   nil,
		AccessTokenChecksum2:    api.Checksum{"access-token-cs-2"},
	}

	df.Org1 = api.Org{
		df.OrgKey1,
		df.OrgName1,
		df.OrgContactEmail1,
		df.OrgChecksum1,
	}
	df.Org2 = api.Org{
		df.OrgKey2,
		df.OrgName2,
		df.OrgContactEmail2,
		df.OrgChecksum2,
	}
	df.OrgSlice = api.Orgs{df.Org1, df.Org2}

	ts := time.Date(2015, 2, 28, 12, 30, 0, 0, time.UTC)
	df.UserDeletedAt2 = &ts
	df.Zone1 = api.Zone{
		df.ZoneKey1,
		df.ZoneName1,
		df.ZoneOrgKey1,
		df.ZoneChecksum1,
	}
	df.Zone2 = api.Zone{
		df.ZoneKey2,
		df.ZoneName2,
		df.ZoneOrgKey2,
		df.ZoneChecksum2,
	}
	df.ZoneSlice = api.Zones{df.Zone1, df.Zone2}
	df.PublicZoneSlice = make(api.Zones, len(df.ZoneSlice))
	for i, z := range df.ZoneSlice {
		z.OrgKey = ""
		df.PublicZoneSlice[i] = z
	}

	df.User1 = api.User{
		UserKey:    df.UserKey1,
		LoginEmail: df.UserLoginEmail1,
		APIAuthKey: df.UserAPIAuthKey1,
		OrgKey:     df.UserOrgKey1,
		DeletedAt:  df.UserDeletedAt1,
		Checksum:   df.UserChecksum1,
	}
	df.User2 = api.User{
		UserKey:    df.UserKey2,
		LoginEmail: df.UserLoginEmail2,
		APIAuthKey: df.UserAPIAuthKey2,
		OrgKey:     df.UserOrgKey2,
		DeletedAt:  df.UserDeletedAt2,
		Checksum:   df.UserChecksum2,
	}
	df.UserSlice = api.Users{df.User1, df.User2}
	df.PublicUserSlice = make(api.Users, len(df.UserSlice))
	for i, u := range df.UserSlice {
		df.PublicUserSlice[i] = u
	}

	df.Cluster1 = api.Cluster{
		ClusterKey: df.ClusterKey1,
		ZoneKey:    df.ClusterZone1,
		Name:       df.ClusterName1,
		RequireTLS: df.ClusterRequireTLS1,
		OrgKey:     df.ClusterOrgKey1,
		Checksum:   df.ClusterChecksum1,
	}

	df.Instance21 = api.Instance{
		Host: "int-host", Port: 1234, Metadata: api.Metadata{{"key1", "value1"}, {"key2", "value2"}}}

	df.Instance22 = api.Instance{Host: "int-host-2", Port: 1234}

	df.Cluster2 = api.Cluster{
		ClusterKey: df.ClusterKey2,
		ZoneKey:    df.ClusterZone2,
		Name:       df.ClusterName2,
		RequireTLS: df.ClusterRequireTLS2,
		Instances:  api.Instances{df.Instance21, df.Instance22},
		OrgKey:     df.ClusterOrgKey2,
		Checksum:   df.ClusterChecksum2,
	}

	df.ClusterSlice = []api.Cluster{df.Cluster1, df.Cluster2}
	df.PublicClusterSlice = make(api.Clusters, len(df.ClusterSlice))
	for i, c := range df.ClusterSlice {
		c.OrgKey = ""
		df.PublicClusterSlice[i] = c
	}

	// domain setup
	df.Domain1 = api.Domain{
		df.DomainKey1,
		df.DomainZone1,
		df.DomainName1,
		df.DomainPort1,
		df.DomainSSLConfig1,
		df.DomainRedirects1,
		df.DomainGzipEnabled1,
		df.DomainCorsConfig1,
		df.DomainAliases1,
		df.DomainOrgKey1,
		df.DomainChecksum1,
	}

	df.Domain2 = api.Domain{
		df.DomainKey2,
		df.DomainZone2,
		df.DomainName2,
		df.DomainPort2,
		df.DomainSSLConfig2,
		df.DomainRedirects2,
		df.DomainGzipEnabled2,
		df.DomainCorsConfig2,
		df.DomainAliases2,
		df.DomainOrgKey2,
		df.DomainChecksum2,
	}

	df.DomainSlice = api.Domains{df.Domain1, df.Domain2}
	df.PublicDomainSlice = make(api.Domains, len(df.DomainSlice))
	for i, d := range df.DomainSlice {
		d.OrgKey = ""
		df.PublicDomainSlice[i] = d
	}

	// proxy setup
	df.ProxyDomainKeys1 = []api.DomainKey{
		df.Domain1.DomainKey,
		df.Domain2.DomainKey,
	}
	df.Proxy1 = api.Proxy{
		df.ProxyKey1,
		df.ProxyZone1,
		df.ProxyName1,
		df.ProxyDomainKeys1,
		df.ProxyOrgKey1,
		df.ProxyChecksum1,
	}

	df.ProxyDomainKeys2 = df.ProxyDomainKeys1
	df.Proxy2 = api.Proxy{
		df.ProxyKey2,
		df.ProxyZone2,
		df.ProxyName2,
		df.ProxyDomainKeys2,
		df.ProxyOrgKey2,
		df.ProxyChecksum2,
	}

	df.ProxySlice = []api.Proxy{df.Proxy1, df.Proxy2}
	df.PublicProxySlice = make(api.Proxies, len(df.ProxySlice))
	for i, p := range df.ProxySlice {
		p.OrgKey = ""
		df.PublicProxySlice[i] = p
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
				{
					"cc-0",
					"ckey2",
					api.Metadata{{"key-2", "value-2"}},
					api.Metadata{{"state", "test"}},
					api.ResponseData{},
					1234,
				}}},
		nil,
	}

	routeRule2 := api.Rule{
		"rk-0-1",
		[]string{"PUT", "DELETE"},
		api.Matches{
			api.Match{api.CookieMatchKind, api.Metadatum{"x-2", "value"}, api.Metadatum{"other", "true"}}},
		api.AllConstraints{
			Tap: api.ClusterConstraints{
				{"cc-1", "ckey3", api.Metadata{{"key-2", "value-2"}}, api.Metadata{}, api.ResponseData{}, 1234}},
			Light: api.ClusterConstraints{
				{"cc-2", "ckey2", api.Metadata{{"key-2", "value-2"}}, api.Metadata{}, api.ResponseData{}, 1234}}},
		nil,
	}

	df.RouteRules1 = api.Rules{routeRule1}
	df.RouteResponseData1 = api.ResponseData{
		Headers: []api.HeaderDatum{
			{
				api.ResponseDatum{
					Name:  "X-Route1-Tbn-Server-Addr",
					Value: "server-addr",
				},
			},
			{
				ResponseDatum: api.ResponseDatum{
					Name:           "X-Route1-Tbn-Server-Literal",
					Value:          "some literal value",
					ValueIsLiteral: true,
				},
			},
		},
		Cookies: []api.CookieDatum{
			{
				ResponseDatum: api.ResponseDatum{
					Name:  "route1-server-version",
					Value: "server-version",
				},
				Secure:   true,
				SameSite: api.SameSiteStrict,
			},
			{
				ResponseDatum: api.ResponseDatum{
					Name:  "route1-server-addr",
					Value: "server-addr",
				},
				Secure:   true,
				SameSite: api.SameSiteStrict,
			},
		},
	}
	df.RouteRules2 = api.Rules{routeRule1, routeRule2}
	df.RouteResponseData2 = api.ResponseData{}
	df.Route1 = api.Route{
		df.RouteKey1,
		df.RouteDomain1,
		df.RouteZone1,
		df.RoutePath1,
		df.SharedRulesKey1,
		df.RouteRules1,
		df.RouteResponseData1,
		df.RouteCohortSeed1,
		df.RouteRetryPolicy1,
		df.RouteOrgKey1,
		df.RouteChecksum1,
	}

	df.Route2 = api.Route{
		df.RouteKey2,
		df.RouteDomain2,
		df.RouteZone2,
		df.RoutePath2,
		df.SharedRulesKey2,
		df.RouteRules2,
		df.RouteResponseData2,
		df.RouteCohortSeed2,
		df.RouteRetryPolicy2,
		df.RouteOrgKey2,
		df.RouteChecksum2,
	}

	df.RouteSlice = api.Routes{df.Route1, df.Route2}
	df.PublicRouteSlice = make(api.Routes, len(df.RouteSlice))
	for i, r := range df.RouteSlice {
		r.OrgKey = ""
		df.PublicRouteSlice[i] = r
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
				{"cc-0", "ckey2", api.Metadata{{"key-2", "value-2"}}, api.Metadata{{"state", "test"}}, api.ResponseData{}, 1234}}},
		nil,
	}

	sharedRulesRule2 := api.Rule{
		"srk-0-1",
		[]string{"PUT", "DELETE"},
		api.Matches{
			api.Match{api.CookieMatchKind, api.Metadatum{"x-2", "value"}, api.Metadatum{"other", "true"}}},
		api.AllConstraints{
			Tap: api.ClusterConstraints{
				{"cc-1", "ckey3", api.Metadata{{"key-2", "value-2"}}, api.Metadata{}, api.ResponseData{}, 1234}},
			Light: api.ClusterConstraints{
				{"cc-2", "ckey2", api.Metadata{{"key-2", "value-2"}}, api.Metadata{}, api.ResponseData{}, 1234}}},
		nil,
	}

	sharedRulesDefault1 := api.AllConstraints{
		Light: api.ClusterConstraints{
			{"cc-4", api.HeaderMatchKind, api.Metadata{{"k", "v"}, {"k2", "v2"}}, api.Metadata{{"state", "released"}}, api.ResponseData{}, 23}}}
	sharedRulesDefault2 := sharedRulesDefault1

	df.SharedRulesDefault1 = sharedRulesDefault1
	df.SharedRulesRules1 = api.Rules{sharedRulesRule1}
	df.SharedRulesResponseData1 = api.ResponseData{
		Headers: []api.HeaderDatum{
			{
				api.ResponseDatum{
					Name:  "X-Tbn-Server-Addr",
					Value: "server-addr",
				},
			},
			{
				api.ResponseDatum{
					Name:           "X-Tbn-Server-Literal",
					Value:          "some literal value",
					ValueIsLiteral: true,
				},
			},
		},
		Cookies: []api.CookieDatum{
			{
				ResponseDatum: api.ResponseDatum{
					Name:  "server-version",
					Value: "server-version",
				},
				Secure:   true,
				SameSite: api.SameSiteStrict,
			},
			{
				ResponseDatum: api.ResponseDatum{
					Name:  "server-addr",
					Value: "server-addr",
				},
				Secure:   true,
				SameSite: api.SameSiteStrict,
			},
		},
	}
	df.SharedRulesDefault2 = sharedRulesDefault2
	df.SharedRulesRules2 = api.Rules{sharedRulesRule1, sharedRulesRule2}
	df.SharedRulesResponseData2 = api.ResponseData{}
	df.SharedRules1 = api.SharedRules{
		df.SharedRulesKey1,
		df.SharedRulesName1,
		df.SharedRulesZone1,
		df.SharedRulesDefault1,
		df.SharedRulesRules1,
		df.SharedRulesResponseData1,
		df.SharedRulesCohortSeed1,
		df.SharedRulesProperties1,
		df.SharedRulesRetryPolicy1,
		df.SharedRulesOrgKey1,
		df.SharedRulesChecksum1,
	}

	df.SharedRules2 = api.SharedRules{
		df.SharedRulesKey2,
		df.SharedRulesName2,
		df.SharedRulesZone2,
		df.SharedRulesDefault2,
		df.SharedRulesRules2,
		df.SharedRulesResponseData2,
		df.SharedRulesCohortSeed2,
		df.SharedRulesProperties2,
		df.SharedRulesRetryPolicy2,
		df.SharedRulesOrgKey2,
		df.SharedRulesChecksum2,
	}

	df.SharedRulesSlice = api.SharedRulesSlice{df.SharedRules1, df.SharedRules2}
	df.PublicSharedRulesSlice = make(api.SharedRulesSlice, len(df.SharedRulesSlice))
	for i, r := range df.SharedRulesSlice {
		r.OrgKey = ""
		df.PublicSharedRulesSlice[i] = r
	}

	now := time.Now()
	df.AccessTokenCreatedAt1 = &now
	df.AccessTokenCreatedAt2 = &now

	df.AccessToken1 = api.AccessToken{
		df.AccessTokenKey1,
		df.AccessTokenDescription1,
		"",
		df.AccessTokenUserKey1,
		df.AccessTokenOrgKey1,
		df.AccessTokenCreatedAt1,
		df.AccessTokenChecksum1,
	}

	df.AccessToken2 = api.AccessToken{
		df.AccessTokenKey2,
		df.AccessTokenDescription2,
		"",
		df.AccessTokenUserKey2,
		df.AccessTokenOrgKey2,
		df.AccessTokenCreatedAt2,
		df.AccessTokenChecksum2,
	}

	df.AccessTokenSlice = api.AccessTokens{
		df.AccessToken1,
		df.AccessToken2,
	}

	return df
}
