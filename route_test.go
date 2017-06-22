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

package api

import (
	"testing"

	"github.com/turbinelabs/test/assert"
)

func getRulesDefaults() (Rule, Rule) {
	rule1 := Rule{
		"rk0",
		[]string{"GET", "POST"},
		Matches{
			Match{HeaderMatchKind, Metadatum{"x-1", "value"}, Metadatum{"flag", "true"}},
			Match{CookieMatchKind, Metadatum{"x-2", "value"}, Metadatum{"other", "true"}}},
		AllConstraints{
			Light: ClusterConstraints{
				ClusterConstraint{"cckey1", "ckey2", Metadata{{"key-2", "value-2"}}, nil, ResponseData{}, 1234}}},
		nil,
	}

	rule2 := Rule{
		"rk1",
		[]string{"PUT", "DELETE"},
		Matches{
			Match{CookieMatchKind, Metadatum{"x-2", "value"}, Metadatum{"other", "true"}}},
		AllConstraints{
			Tap: ClusterConstraints{
				ClusterConstraint{"cckey1", "ckey3", Metadata{{"key-2", "value-2"}}, nil, ResponseData{}, 1234}},
			Light: ClusterConstraints{
				ClusterConstraint{"cckey2", "ckey2", Metadata{{"key-2", "value-2"}}, nil, ResponseData{}, 1234}}},
		nil,
	}

	return rule1, rule2
}

func getRouteDefaults() (Route, Route) {
	srk := SharedRulesKey("srk-1")
	rule1, rule2 := getRulesDefaults()
	rules := Rules{rule1, rule2}

	r1 := Route{
		"routekey",
		"dkey",
		"zkey",
		"host/slug/some/other",
		srk,
		rules,
		getRD(),
		&CohortSeed{CohortSeedHeader, "x-cohort-seed", true},
		"1",
		Checksum{"cs-1"},
	}

	r2 := Route{
		"routekey",
		"dkey",
		"zkey",
		"host/slug/some/other",
		srk,
		rules,
		getRD(),
		&CohortSeed{CohortSeedHeader, "x-cohort-seed", true},
		"1",
		Checksum{"cs-1"},
	}

	return r1, r2
}

// Route.Equals
func TestRouteEqualsSuccess(t *testing.T) {
	r1, r2 := getRouteDefaults()

	assert.True(t, r1.Equals(r2))
	assert.True(t, r2.Equals(r1))
}

func TestRouteEqualsCohortSeedVaries(t *testing.T) {
	r1, r2 := getRouteDefaults()
	r2.CohortSeed.Name = r1.CohortSeed.Name + "aosentuh"

	assert.False(t, r1.Equals(r2))
	assert.False(t, r2.Equals(r1))
}

func TestRouteEqualsCohortSeedNotNilNil(t *testing.T) {
	r1, r2 := getRouteDefaults()
	r2.CohortSeed = nil

	assert.False(t, r1.Equals(r2))
	assert.False(t, r2.Equals(r1))
}

func TestRouteEqualsCohortSeedNilNil(t *testing.T) {
	r1, r2 := getRouteDefaults()
	r1.CohortSeed = nil
	r2.CohortSeed = nil

	assert.True(t, r1.Equals(r2))
	assert.True(t, r2.Equals(r1))
}

func TestRouteEqualsResponseDataVaries(t *testing.T) {
	r1, r2 := getRouteDefaults()
	r1.ResponseData.Headers[0].Value += "aosenuth"

	assert.False(t, r1.Equals(r2))
	assert.False(t, r2.Equals(r1))
}

func TestRouteEqualsOrgVaries(t *testing.T) {
	r1, r2 := getRouteDefaults()
	r1.OrgKey = "snth"

	assert.False(t, r1.Equals(r2))
	assert.False(t, r2.Equals(r1))
}

func TestRouteEqualsRouteKeyVaries(t *testing.T) {
	r1, r2 := getRouteDefaults()
	r1.RouteKey = "asonteuh"

	assert.False(t, r1.Equals(r2))
	assert.False(t, r2.Equals(r1))
}

func TestRouteEqualsDomainKeyVaries(t *testing.T) {
	r1, r2 := getRouteDefaults()
	r1.DomainKey = "saeouashetnoas"

	assert.False(t, r1.Equals(r2))
	assert.False(t, r2.Equals(r1))
}

func TestRouteEqualsZoneKeyVaries(t *testing.T) {
	r1, r2 := getRouteDefaults()
	r1.ZoneKey = "saotehuasontehu"

	assert.False(t, r1.Equals(r2))
	assert.False(t, r2.Equals(r1))
}

func TestRouteEqualsPathVaries(t *testing.T) {
	r1, r2 := getRouteDefaults()
	r1.Path = "saonteuh"

	assert.False(t, r1.Equals(r2))
	assert.False(t, r2.Equals(r1))
}

func TestRouteEqualsSharedRulesKeyVaries(t *testing.T) {
	r1, r2 := getRouteDefaults()
	r1.SharedRulesKey = SharedRulesKey("newaoesutnhao")

	assert.False(t, r1.Equals(r2))
	assert.False(t, r2.Equals(r1))
}

func TestRouteEqualsRulesVaryOrder(t *testing.T) {
	r1, r2 := getRouteDefaults()
	rule2, rule1 := getRulesDefaults()
	r1.Rules = Rules{rule1, rule2}

	assert.False(t, r1.Equals(r2))
	assert.False(t, r2.Equals(r1))
}

func TestRouteEqualsRulesVaryContents(t *testing.T) {
	r1, r2 := getRouteDefaults()
	rule1, rule2 := getRulesDefaults()
	rule1.Methods = []string{"DELETE", "POST"}
	r1.Rules = Rules{rule1, rule2}

	assert.False(t, r1.Equals(r2))
	assert.False(t, r2.Equals(r1))
}

func TestRouteEqualsChecksumVaries(t *testing.T) {
	r1, r2 := getRouteDefaults()
	r1.Checksum = Checksum{"asontehuasoneht"}

	assert.False(t, r1.Equals(r2))
	assert.False(t, r2.Equals(r1))
}

func TestRouteIsValidSuccess(t *testing.T) {
	r, _ := getRouteDefaults()

	assert.Nil(t, r.IsValid())
}

func TestRouteIsValidBadCohortseed(t *testing.T) {
	r, _ := getRouteDefaults()
	r.CohortSeed.Name = ""

	assert.DeepEqual(t, r.IsValid(), &ValidationError{[]ErrorCase{
		{"route.cohort_seed.name", "may not be empty"},
	}})
}

func TestRouteIsValidBadResponseData(t *testing.T) {
	r, _ := getRouteDefaults()
	r.ResponseData.Headers[0].Value = ""
	n := r.ResponseData.Headers[0].Name

	assert.DeepEqual(t, r.IsValid(), &ValidationError{[]ErrorCase{
		{"route.response_data.headers[" + n + "].value", "may not be empty"},
	}})
}

func TestRouteIsValidNoDomainKey(t *testing.T) {
	r, _ := getRouteDefaults()
	r.DomainKey = ""

	assert.NonNil(t, r.IsValid())
}

func TestRouteIsValidBadDomainKey(t *testing.T) {
	r, _ := getRouteDefaults()
	r.DomainKey = "key $"

	assert.NonNil(t, r.IsValid())
}

func TestRouteIsValidBadZoneKey(t *testing.T) {
	r, _ := getRouteDefaults()
	r.ZoneKey = ""

	assert.NonNil(t, r.IsValid())
}

func TestRouteIsValidBadPath(t *testing.T) {
	r, _ := getRouteDefaults()
	r.Path = ""

	assert.NonNil(t, r.IsValid())
}

func TestRouteIsValidBadSharedRulesKeyRef(t *testing.T) {
	r, _ := getRouteDefaults()
	r.SharedRulesKey = ""

	assert.NonNil(t, r.IsValid())
}

func TestRouteIsValidBadRules(t *testing.T) {
	r, _ := getRouteDefaults()
	rule1, _ := getRulesDefaults()
	rule1.Matches = Matches{}
	rule1.Methods = []string{}
	r.Rules[0] = rule1

	assert.NonNil(t, r.IsValid())
}

func TestRouteIsValidDupeRules(t *testing.T) {
	r, _ := getRouteDefaults()
	rule1, _ := getRulesDefaults()
	r.Rules = Rules{rule1, rule1}

	errs := r.IsValid()
	assert.DeepEqual(t, errs, &ValidationError{[]ErrorCase{
		{"route.rules", "multiple instances of key " + string(rule1.RuleKey)},
	}})
}
