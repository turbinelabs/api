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

package service

import (
	"testing"

	"github.com/turbinelabs/api"
	"github.com/turbinelabs/test/assert"
)

func getSharedRulesFilterTestSharedRules() api.SharedRules {
	defaultCC := api.AllConstraints{
		Light: api.ClusterConstraints{
			{"cc1", api.HeaderMatchKind, api.Metadata{{"k", "v"}, {"k2", "v2"}}, nil, api.ResponseData{}, 23}}}

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
				{"cckey2", "ckey2", api.Metadata{{"key-2", "value-2"}}, api.Metadata{{"state", "releasing"}}, api.ResponseData{}, 1234}}},
		nil,
	}

	rules := api.Rules{rule1}

	r1 := api.SharedRules{
		"shared-rules-key",
		"shared-rules-name",
		"zkey",
		defaultCC,
		rules,
		api.ResponseData{},
		nil,
		api.Metadata{{"pk", "pv"}, {"pk2", "pv2"}},
		nil,
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
				{"cckey2", "ckey2", api.Metadata{{"key-2", "value-2"}}, api.Metadata{{"state", "releasing"}}, api.ResponseData{}, 1234}}},
		nil,
	}

	rules := api.Rules{rule1}

	r1 := api.Route{
		"routekey",
		"dkey",
		"zkey",
		"host/slug/some/other",
		api.SharedRulesKey("shared-rules-key"),
		rules,
		api.ResponseData{},
		nil,
		nil,
		"123",
		api.Checksum{"cs-1"},
	}

	return r1
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
