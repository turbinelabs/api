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

package service

import (
	"testing"

	"github.com/turbinelabs/api"
	"github.com/turbinelabs/test/assert"
)

func getSharedRulesFilterTestSharedRules() api.SharedRules {
	defaultCC := api.AllConstraints{
		Light: api.ClusterConstraints{
			{"cc1", api.HeaderMatchKind, api.Metadata{{"k", "v"}, {"k2", "v2"}}, nil, 23}}}

	rule1 := api.Rule{
		"rk1",
		[]string{"GET", "POST"},
		api.Matches{
			api.Match{
				api.HeaderMatchKind,
				api.Metadatum{"x-1", "value"},
				api.Metadatum{"flag", "true"}},
			api.Match{
				api.CookieMatchKind,
				api.Metadatum{"x-2", "value"},
				api.Metadatum{"other", "true"}}},
		api.AllConstraints{
			Light: api.ClusterConstraints{
				{"cckey2", "ckey2", api.Metadata{{"key-2", "value-2"}}, api.Metadata{{"state", "releasing"}}, 1234}}},
	}

	rules := api.Rules{rule1}

	r1 := api.SharedRules{
		"shared-rules-key",
		"shared-rules-name",
		"zkey",
		defaultCC,
		rules,
		"123",
		api.Checksum{"cs-1"},
	}

	return r1
}

func getRouteFilterTestRoute() api.Route {
	rule1 := api.Rule{
		"rk1",
		[]string{"GET", "POST"},
		api.Matches{
			api.Match{
				api.HeaderMatchKind,
				api.Metadatum{"x-1", "value"},
				api.Metadatum{"flag", "true"}},
			api.Match{
				api.CookieMatchKind,
				api.Metadatum{"x-2", "value"},
				api.Metadatum{"other", "true"}}},
		api.AllConstraints{
			Light: api.ClusterConstraints{
				{"cckey2", "ckey2", api.Metadata{{"key-2", "value-2"}}, api.Metadata{{"state", "releasing"}}, 1234}}},
	}

	rules := api.Rules{rule1}

	r1 := api.Route{
		"routekey",
		"dkey",
		"zkey",
		"host/slug/some/other",
		api.SharedRulesKey("shared-rules-key"),
		rules,
		"123",
		api.Checksum{"cs-1"},
	}

	return r1
}

// SharedRulesFilter.Match
func TestSharedRulesFilterMatchOnEmpty(t *testing.T) {
	r1 := getSharedRulesFilterTestSharedRules()
	f := SharedRulesFilter{}
	assert.True(t, f.Matches(r1))
}

func TestSharedRulesFilterMatchOnOrgKey(t *testing.T) {
	r1 := getSharedRulesFilterTestSharedRules()
	f := SharedRulesFilter{OrgKey: "123"}
	assert.True(t, f.Matches(r1))
}

func TestSharedRulesFilterMatchOnZoneKey(t *testing.T) {
	r1 := getSharedRulesFilterTestSharedRules()
	f := SharedRulesFilter{ZoneKey: r1.ZoneKey}
	assert.True(t, f.Matches(r1))
}

func TestSharedRulesFilterMismatchOnZoneKey(t *testing.T) {
	r1 := getSharedRulesFilterTestSharedRules()
	f := SharedRulesFilter{ZoneKey: r1.ZoneKey + r1.ZoneKey}
	assert.False(t, f.Matches(r1))
}

func TestSharedRulesFilterMatchAndMismatchCorrectly(t *testing.T) {
	r1 := getSharedRulesFilterTestSharedRules()
	f := SharedRulesFilter{ZoneKey: r1.ZoneKey + r1.ZoneKey, SharedRulesKey: api.SharedRulesKey("nope")}
	assert.False(t, f.Matches(r1))
}

// RouteFilter.Match
func TestRouteFilterMatchOnEmpty(t *testing.T) {
	r1 := getRouteFilterTestRoute()
	f := RouteFilter{}
	assert.True(t, f.Matches(r1))
}

func TestRouteFilterMatchOnOrgKey(t *testing.T) {
	r1 := getRouteFilterTestRoute()
	f := RouteFilter{OrgKey: "123"}
	assert.True(t, f.Matches(r1))
}

func TestRouteFilterMatchOnPath(t *testing.T) {
	r1 := getRouteFilterTestRoute()
	f := RouteFilter{Path: r1.Path}
	assert.True(t, f.Matches(r1))
}

func TestRouteFilterMatchOnPathPrefix(t *testing.T) {
	r1 := getRouteFilterTestRoute()
	f := RouteFilter{PathPrefix: "host/slug"}
	assert.True(t, f.Matches(r1))
}

func TestRouteFilterMatchOnZoneKey(t *testing.T) {
	r1 := getRouteFilterTestRoute()
	f := RouteFilter{ZoneKey: r1.ZoneKey}
	assert.True(t, f.Matches(r1))
}

func TestRouteFilterMatchOnDomain(t *testing.T) {
	r1 := getRouteFilterTestRoute()
	f := RouteFilter{DomainKey: r1.DomainKey}
	assert.True(t, f.Matches(r1))
}

func TestRouteFilterMismatchOnPath(t *testing.T) {
	r1 := getRouteFilterTestRoute()
	f := RouteFilter{Path: r1.Path + r1.Path}
	assert.False(t, f.Matches(r1))
}

func TestRouteFilterMismatchOnPathPrefix(t *testing.T) {
	r1 := getRouteFilterTestRoute()
	f := RouteFilter{PathPrefix: "nopehost/slug"}
	assert.False(t, f.Matches(r1))
}

func TestRouteFilterMismatchOnZoneKey(t *testing.T) {
	r1 := getRouteFilterTestRoute()
	f := RouteFilter{ZoneKey: r1.ZoneKey + r1.ZoneKey}
	assert.False(t, f.Matches(r1))
}

func TestRouteFilterMismatchOnDomain(t *testing.T) {
	r1 := getRouteFilterTestRoute()
	f := RouteFilter{DomainKey: r1.DomainKey + r1.DomainKey}
	assert.False(t, f.Matches(r1))
}

func TestRouteFilterMatchAndMismatchCorrectly(t *testing.T) {
	r1 := getRouteFilterTestRoute()
	f := RouteFilter{DomainKey: r1.DomainKey, Path: "someothershit"}
	assert.False(t, f.Matches(r1))
}

func TestProxyFilterMatches(t *testing.T) {
	type pf ProxyFilter
	type dk []api.DomainKey
	type testcase struct {
		name        string
		f1          pf
		f2          pf
		shouldMatch bool
	}

	run := func(tc testcase) {
		assert.Group(tc.name, t, func(tg *assert.G) {
			assert.Equal(tg, ProxyFilter(tc.f1).Equals(ProxyFilter(tc.f2)), tc.shouldMatch)
		})
	}

	cases := []testcase{
		{"two empty filters", pf{}, pf{}, true},
		{"only domainkeys, nil", pf{DomainKeys: dk{"a", "b"}}, pf{}, false},
		{"different domainkeys, nil", pf{DomainKeys: dk{"a", "b"}}, pf{DomainKeys: dk{"a", "c"}}, false},
		{"same domain keys, different order", pf{DomainKeys: dk{"a", "b"}}, pf{DomainKeys: dk{"b", "a"}}, true},
	}

	for _, c := range cases {
		run(c)
	}
}

func getProxyFilter() ProxyFilter {
	return ProxyFilter{
		"proxy-key",
		api.Instance{"host", 8080, nil},
		"proxy-name",
		[]api.DomainKey{"key", "key2"},
		"zone-key",
		"org-key",
	}
}

func TestProxyFilterEquals(t *testing.T) {
	p1 := getProxyFilter()
	p2 := getProxyFilter()

	assert.True(t, p1.Equals(p2))
}

func TestProxyFilterEqualsMismatchProxyKey(t *testing.T) {
	p1 := getProxyFilter()
	p2 := getProxyFilter()
	p2.ProxyKey = "aoesutnh"

	assert.False(t, p1.Equals(p2))
}

func TestProxyFilterEqualsMismatchInstance(t *testing.T) {
	p1 := getProxyFilter()
	p2 := getProxyFilter()
	p2.Instance = api.Instance{"asoetunh", 8080, nil}

	assert.False(t, p1.Equals(p2))
}

func TestProxyFilterEqualsMismatchName(t *testing.T) {
	p1 := getProxyFilter()
	p2 := getProxyFilter()
	p2.Name = "aoesutnh"

	assert.False(t, p1.Equals(p2))
}

func TestProxyFilterEqualsMismatchDomainKeys(t *testing.T) {
	p1 := getProxyFilter()
	p2 := getProxyFilter()
	p2.DomainKeys = []api.DomainKey{"key2"}

	assert.False(t, p1.Equals(p2))
}

func TestProxyFilterEqualsMismatchZoneKey(t *testing.T) {
	p1 := getProxyFilter()
	p2 := getProxyFilter()
	p2.ZoneKey = "aoesutnh"

	assert.False(t, p1.Equals(p2))
}

func TestProxyFilterEqualsMismatchOrgKey(t *testing.T) {
	p1 := getProxyFilter()
	p2 := getProxyFilter()
	p2.OrgKey = "aoesutnh"

	assert.False(t, p1.Equals(p2))
}
