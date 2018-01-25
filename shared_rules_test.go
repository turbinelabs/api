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

func getSharedRulesDefaults() (SharedRules, SharedRules) {
	defaultCC := AllConstraints{
		Light: ClusterConstraints{
			ClusterConstraint{
				"cckey1",
				HeaderMatchKind,
				Metadata{{"k", "v"}, {"k2", "v2"}},
				Metadata{{"state", "released"}},
				getRD(),
				23}}}

	rule1, rule2 := getRulesDefaults()
	rules := Rules{rule1, rule2}

	r1 := SharedRules{
		"routekey",
		"rule1-name",
		"zkey",
		defaultCC,
		rules,
		getRD(),
		&CohortSeed{CohortSeedHeader, "x-cohort-data", false},
		Metadata{{"pk", "pv"}, {"pk2", "pv2"}},
		nil,
		"1",
		Checksum{"cs-1"},
	}

	r2 := SharedRules{
		"routekey",
		"rule1-name",
		"zkey",
		defaultCC,
		rules,
		getRD(),
		&CohortSeed{CohortSeedHeader, "x-cohort-data", false},
		Metadata{{"pk", "pv"}, {"pk2", "pv2"}},
		nil,
		"1",
		Checksum{"cs-1"},
	}

	return r1, r2
}

// SharedRules.Equals
func TestSharedRulesEqualsSuccess(t *testing.T) {
	r1, r2 := getSharedRulesDefaults()

	assert.True(t, r1.Equals(r2))
	assert.True(t, r2.Equals(r1))
}

func TestSharedRulesEqualsCohortSeedNilNil(t *testing.T) {
	r1, r2 := getSharedRulesDefaults()
	r1.CohortSeed = nil
	r2.CohortSeed = nil

	assert.True(t, r1.Equals(r2))
	assert.True(t, r2.Equals(r1))
}

func TestSharedRulesEqualsCohortSeedNotNilNil(t *testing.T) {
	r1, r2 := getSharedRulesDefaults()
	r2.CohortSeed = nil

	assert.False(t, r1.Equals(r2))
	assert.False(t, r2.Equals(r1))
}

func TestSharedRulesEqualsCohortSeedNilNotNil(t *testing.T) {
	r1, r2 := getSharedRulesDefaults()
	r1.CohortSeed = nil

	assert.False(t, r1.Equals(r2))
	assert.False(t, r2.Equals(r1))
}

func TestSharedRulesEqualsCohortSeedVaries(t *testing.T) {
	r1, r2 := getSharedRulesDefaults()
	r1.CohortSeed.UseZeroValueSeed = !r2.CohortSeed.UseZeroValueSeed

	assert.False(t, r1.Equals(r2))
	assert.False(t, r2.Equals(r1))
}

func TestSharedRulesEqualsResponseDataVaries(t *testing.T) {
	r1, r2 := getSharedRulesDefaults()
	r1.ResponseData.Headers[0].Value += "-new"

	assert.False(t, r1.Equals(r2))
	assert.False(t, r2.Equals(r1))
}

func TestSharedRulesEqualsPropertiesVaries(t *testing.T) {
	r1, r2 := getSharedRulesDefaults()
	r1.Properties[1].Value = "asdasd"

	assert.False(t, r1.Equals(r2))
	assert.False(t, r2.Equals(r1))
}

func TestSharedRulesEqualsOrgVaries(t *testing.T) {
	r1, r2 := getSharedRulesDefaults()
	r1.OrgKey = "snth"

	assert.False(t, r1.Equals(r2))
	assert.False(t, r2.Equals(r1))
}

func TestSharedRulesEqualsSharedRulesKeyVaries(t *testing.T) {
	r1, r2 := getSharedRulesDefaults()
	r1.SharedRulesKey = "asonteuh"

	assert.False(t, r1.Equals(r2))
	assert.False(t, r2.Equals(r1))
}

func TestSharedRulesEqualsZoneKeyVaries(t *testing.T) {
	r1, r2 := getSharedRulesDefaults()
	r1.ZoneKey = "saotehuasontehu"

	assert.False(t, r1.Equals(r2))
	assert.False(t, r2.Equals(r1))
}

func TestSharedRulesEqualsDefaultVaries(t *testing.T) {
	defaultCC := AllConstraints{
		Light: ClusterConstraints{
			ClusterConstraint{"cckey1", HeaderMatchKind, Metadata{{"k1", "v1"}, {"k2", "v2"}}, nil, ResponseData{}, 23},
			ClusterConstraint{"cckey2", HeaderMatchKind, Metadata{{"k2", "v2"}, {"k2", "v2"}}, nil, ResponseData{}, 23}}}

	r1, r2 := getSharedRulesDefaults()
	r1.Default = defaultCC

	assert.False(t, r1.Equals(r2))
	assert.False(t, r2.Equals(r1))
}

func TestSharedRulesEqualsRulesVaryOrder(t *testing.T) {
	r1, r2 := getSharedRulesDefaults()
	rule2, rule1 := getRulesDefaults()
	r1.Rules = Rules{rule1, rule2}

	assert.False(t, r1.Equals(r2))
	assert.False(t, r2.Equals(r1))
}

func TestSharedRulesEqualsRulesVaryContents(t *testing.T) {
	r1, r2 := getSharedRulesDefaults()
	rule1, rule2 := getRulesDefaults()
	rule1.Methods = []string{"DELETE", "POST"}
	r1.Rules = Rules{rule1, rule2}

	assert.False(t, r1.Equals(r2))
	assert.False(t, r2.Equals(r1))
}

func TestSharedRulesEqualsChecksumVaries(t *testing.T) {
	r1, r2 := getSharedRulesDefaults()
	r1.Checksum = Checksum{"asontehuasoneht"}

	assert.False(t, r1.Equals(r2))
	assert.False(t, r2.Equals(r1))
}

func TestSharedRulesIsValidSuccess(t *testing.T) {
	r, _ := getSharedRulesDefaults()

	assert.Nil(t, r.IsValid())
}

func TestSharedRulesIsValidBadResponseData(t *testing.T) {
	r, _ := getSharedRulesDefaults()
	r.ResponseData.Headers[0].Value = ""
	n := r.ResponseData.Headers[0].Name

	assert.DeepEqual(t, r.IsValid(), &ValidationError{[]ErrorCase{
		{"shared_rules.response_data.headers[" + n + "].value", "may not be empty"},
	}})
}

func TestSharedRulesIsValidBadKey(t *testing.T) {
	r, _ := getSharedRulesDefaults()
	r.SharedRulesKey = "aoeu&snth"

	assert.DeepEqual(t, r.IsValid(), &ValidationError{[]ErrorCase{
		{"shared_rules.shared_rules_key", "must match pattern: ^[0-9a-zA-Z]+(-[0-9a-zA-Z]+)*$"},
	}})
}

func TestSharedRulesIsValidNoKey(t *testing.T) {
	r, _ := getSharedRulesDefaults()
	r.SharedRulesKey = ""

	assert.DeepEqual(t, r.IsValid(), &ValidationError{[]ErrorCase{
		{"shared_rules.shared_rules_key", "may not be empty"},
	}})
}

func TestSharedRulesIsValidBadName(t *testing.T) {
	r, _ := getSharedRulesDefaults()
	r.Name = "name[name]"

	assert.DeepEqual(t, r.IsValid(), &ValidationError{[]ErrorCase{
		{"shared_rules.name", "may not contain [ or ] characters"},
	}})
}

func TestSharedRulesIsValidNoName(t *testing.T) {
	r, _ := getSharedRulesDefaults()
	r.Name = ""

	assert.NonNil(t, r.IsValid())
}

func TestSharedRulesIsValidBadZoneKey(t *testing.T) {
	r, _ := getSharedRulesDefaults()
	r.ZoneKey = "1234(5678"

	assert.NonNil(t, r.IsValid())
}

func TestSharedRulesIsValidNoZoneKey(t *testing.T) {
	r, _ := getSharedRulesDefaults()
	r.ZoneKey = ""

	assert.NonNil(t, r.IsValid())
}

func TestSharedRulesIsValidBadDefault(t *testing.T) {
	r, _ := getSharedRulesDefaults()
	r.Default = AllConstraints{}

	assert.NonNil(t, r.IsValid())
}

func TestSharedRulesIsValidNoCohort(t *testing.T) {
	r, _ := getSharedRulesDefaults()
	r.CohortSeed = nil

	assert.Nil(t, r.IsValid())
}

func TestSharedRulesIsValidBadCohort(t *testing.T) {
	r, _ := getSharedRulesDefaults()
	r.CohortSeed.Name = ""

	assert.DeepEqual(t, r.IsValid(), &ValidationError{[]ErrorCase{
		{"shared_rules.cohort_seed.name", "may not be empty"},
	}})
}

func TestSharedRulesIsValidBadProperty(t *testing.T) {
	r, _ := getSharedRulesDefaults()
	r.Properties[0].Key = ""
	assert.DeepEqual(t, r.IsValid(), &ValidationError{[]ErrorCase{
		{"shared_rules.properties[].key", "must not be empty"},
	}})
}

func TestSharedRulesIsValidNoMetadataNoProperty(t *testing.T) {
	r, _ := getSharedRulesDefaults()
	r.Properties = nil
	assert.Nil(t, r.IsValid())
}

func TestSharedRulesIsValidBadRules(t *testing.T) {
	r, _ := getSharedRulesDefaults()
	rule1, _ := getRulesDefaults()
	rule1.Methods = []string{"WHEE"}
	r.Rules[0] = rule1
	r.Default.Light[0].Weight = 0

	assert.DeepEqual(t, r.IsValid(), &ValidationError{[]ErrorCase{
		{"shared_rules.default.light[cckey1].weight", "must be greater than 0"},
		{"shared_rules.rules[rk0].methods", "WHEE is not a valid method"},
	}})
}
