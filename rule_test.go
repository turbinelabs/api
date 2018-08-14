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

package api

import (
	"strings"
	"testing"

	"github.com/turbinelabs/test/assert"
)

func getRules() (Rule, Rule) {
	r1 := Rule{
		"rkey1",
		[]string{"GET", "POST"},
		Matches{
			{
				Kind:     HeaderMatchKind,
				Behavior: ExactMatchBehavior,
				From:     Metadatum{Key: "x-random", Value: "value"},
				To:       Metadatum{Key: "randomflag", Value: "true"},
			},
			{
				Kind:     CookieMatchKind,
				Behavior: ExactMatchBehavior,
				From:     Metadatum{Key: "x-other", Value: "value"},
				To:       Metadatum{Key: "otherflag", Value: "true"},
			},
		},
		AllConstraints{
			Light: ClusterConstraints{
				ClusterConstraint{
					"cckey1",
					"ckey1",
					Metadata{{"key", "value"}, {"key2", "value2"}},
					nil,
					ResponseData{},
					1234},
				ClusterConstraint{
					"cckey2",
					"ckey2",
					Metadata{{"key-2", "value-2"}},
					Metadata{{"state", "testing"}},
					ResponseData{},
					1234}}},
		&CohortSeed{CohortSeedHeader, "x-cohort-seed", true},
	}

	r2 := Rule{
		"rkey1",
		[]string{"POST", "GET"},
		Matches{
			{
				Kind:     CookieMatchKind,
				Behavior: ExactMatchBehavior,
				From:     Metadatum{Key: "x-other", Value: "value"},
				To:       Metadatum{Key: "otherflag", Value: "true"},
			},
			{
				Kind:     HeaderMatchKind,
				Behavior: ExactMatchBehavior,
				From:     Metadatum{Key: "x-random", Value: "value"},
				To:       Metadatum{Key: "randomflag", Value: "true"},
			},
		},
		AllConstraints{
			Light: ClusterConstraints{
				ClusterConstraint{
					"cckey1",
					"ckey1",
					Metadata{{"key", "value"}, {"key2", "value2"}},
					nil,
					ResponseData{},
					1234},
				ClusterConstraint{
					"cckey2",
					"ckey2",
					Metadata{{"key-2", "value-2"}},
					Metadata{{"state", "testing"}},
					ResponseData{},
					1234}}},
		&CohortSeed{CohortSeedHeader, "x-cohort-seed", true},
	}

	return r1, r2
}

// Rules.Equals
// Rule.Equals
func TestRuleEqualsSuccess(t *testing.T) {
	r1, r2 := getRules()

	assert.True(t, r1.Equals(r2))
	assert.True(t, r2.Equals(r1))
}

func TestRuleEqualsSuccessCohortNilNil(t *testing.T) {
	r1, r2 := getRules()
	r1.CohortSeed = nil
	r2.CohortSeed = nil

	assert.True(t, r1.Equals(r2))
	assert.True(t, r2.Equals(r1))
}

func TestRuleEqualsFailureCohortVaries(t *testing.T) {
	r1, r2 := getRules()
	r1.CohortSeed.UseZeroValueSeed = !r2.CohortSeed.UseZeroValueSeed

	assert.False(t, r1.Equals(r2))
	assert.False(t, r2.Equals(r1))
}

func TestRuleEqualsFailureCohortNilNotNil(t *testing.T) {
	r1, r2 := getRules()
	r1.CohortSeed = nil

	assert.False(t, r1.Equals(r2))
	assert.False(t, r2.Equals(r1))
}

func TestRuleEqualsFailureCohortNotNilNil(t *testing.T) {
	r1, r2 := getRules()
	r2.CohortSeed = nil

	assert.False(t, r1.Equals(r2))
	assert.False(t, r2.Equals(r1))
}

func TestRuleEqualsKeyMismatchFailure(t *testing.T) {
	r1, r2 := getRules()
	r2.RuleKey = "rkey2"

	assert.False(t, r1.Equals(r2))
	assert.False(t, r2.Equals(r1))
}

func TestRuleEqualsMethodMismatch(t *testing.T) {
	r1, r2 := getRules()
	r2.Methods = []string{"POST", "PUT"}

	assert.False(t, r1.Equals(r2))
	assert.False(t, r2.Equals(r1))
}

func TestRuleEqualsMatchesMismatch(t *testing.T) {
	r1, r2 := getRules()
	r2.Matches = Matches{}

	assert.False(t, r1.Equals(r2))
	assert.False(t, r2.Equals(r1))
}

func TestRuleEqualsConstraintsMismatch(t *testing.T) {
	r1, r2 := getRules()
	r2.Constraints = AllConstraints{
		Light: ClusterConstraints{
			ClusterConstraint{"cckey1", "ckey2", Metadata{{"key-2", "value-2"}}, nil, ResponseData{}, 1234}},
	}

	assert.False(t, r1.Equals(r2))
	assert.False(t, r2.Equals(r1))
}

func getRuleValid() Rule {
	r1, _ := getRules()
	return r1
}

func TestRuleIsValidSucces(t *testing.T) {
	r := getRuleValid()

	assert.Nil(t, r.IsValid())
}

func TestRuleIsValidBadCohort(t *testing.T) {
	r := getRuleValid()
	r.CohortSeed.Name = ""

	assert.DeepEqual(t, r.IsValid(), &ValidationError{[]ErrorCase{
		// cohort_seed is the first segment because the 'rule[$rule_key]' prefix gets
		// attached at the 'Rules' level
		{"cohort_seed.name", "may not be empty"},
	}})
}

func TestRuleIsValidNoRuleKey(t *testing.T) {
	r := getRuleValid()
	r.RuleKey = ""

	assert.NonNil(t, r.IsValid())
}

func TestRuleIsValidBadRuleKey(t *testing.T) {
	r := getRuleValid()
	r.RuleKey = "rule-key-%-1234"

	assert.NonNil(t, r.IsValid())
}

func TestRuleIsValidNoMethodOrMatches(t *testing.T) {
	r := getRuleValid()
	r.Matches = Matches{}
	r.Methods = []string{}
	assert.NonNil(t, r.IsValid())
}

func TestRuleIsValidBadMethod(t *testing.T) {
	r := getRuleValid()
	r.Methods = []string{"POST", "PUT", "GET_THAT_RESOURCE"}

	assert.NonNil(t, r.IsValid())
}

func TestRuleIsValidBadMatches(t *testing.T) {
	r := getRuleValid()
	r.Matches = Matches{
		Match{
			Kind:     CookieMatchKind,
			Behavior: ExactMatchBehavior,
			From:     Metadatum{Key: "x-other", Value: "value"},
			To:       Metadatum{Key: "otherflag", Value: "true"},
		},
		Match{
			Kind:     HeaderMatchKind,
			Behavior: ExactMatchBehavior,
			From:     Metadatum{Key: "x-random", Value: "value"},
			To:       Metadatum{Key: "", Value: "aoeu"},
		},
	}

	assert.NonNil(t, r.IsValid())
}

func TestRuleIsValidBadConstraints(t *testing.T) {
	r := getRuleValid()
	r.Constraints = AllConstraints{
		Dark: ClusterConstraints{{"cckey0", "ckey2", Metadata{{"key-2", "value-2"}}, Metadata{{"aoeu", "snth"}}, ResponseData{}, 1234}}}

	assert.NonNil(t, r.IsValid())
}

func getRulesValidTestRules() (Rule, Rule) {
	r1 := Rule{
		"rkey0",
		[]string{"POST", "PUT"},
		Matches{
			{
				Kind:     CookieMatchKind,
				Behavior: ExactMatchBehavior,
				From:     Metadatum{Key: "x-other", Value: "value"},
				To:       Metadatum{Key: "otherflag", Value: "true"},
			},
			{
				Kind:     HeaderMatchKind,
				Behavior: ExactMatchBehavior,
				From:     Metadatum{Key: "x-random", Value: "value"},
				To:       Metadatum{Key: "randomflag", Value: "true"},
			},
		},
		AllConstraints{
			Light: ClusterConstraints{{"ck0", "ckey2", Metadata{{"key-2", "value-2"}}, nil, ResponseData{}, 1234}},
		},
		&CohortSeed{CohortSeedCookie, "cohort-cookie", false},
	}

	r2 := Rule{
		"rkey1",
		[]string{"GET"},
		Matches{
			{
				Kind:     CookieMatchKind,
				Behavior: ExactMatchBehavior,
				From:     Metadatum{Key: "other", Value: "v"},
				To:       Metadatum{Key: "flag", Value: "true"},
			},
			{
				Kind:     HeaderMatchKind,
				Behavior: ExactMatchBehavior,
				From:     Metadatum{Key: "random", Value: "v"},
				To:       Metadatum{Key: "random", Value: "true"},
			},
		},
		AllConstraints{
			Light: ClusterConstraints{{"ck1", "ckey2", Metadata{{"key-2", "value-2"}}, nil, ResponseData{}, 1234}},
		},
		nil,
	}

	return r1, r2
}

func TestRulesIsValidSuccess(t *testing.T) {
	r1, r2 := getRulesValidTestRules()
	r := Rules{r1, r2}

	assert.Nil(t, r.IsValid())
}

func TestRulesIsValidFailureOnDupeKey(t *testing.T) {
	r1, r2 := getRulesValidTestRules()
	r2.RuleKey = r1.RuleKey
	r := Rules{r1, r2}

	assert.NonNil(t, r.IsValid())
}

func TestRulesIsValidFailureOnDupeMatch(t *testing.T) {
	r1, r2 := getRulesValidTestRules()
	r2.Matches = append(r2.Matches, r1.Matches[0])
	assert.Equal(t, r2.Matches[2], r1.Matches[0])

	r := Rules{r1, r2}
	errors := r.IsValid()

	assert.Equal(t, errors.Errors[0].Attribute, "rules")
	assert.Equal(
		t,
		strings.HasPrefix(errors.Errors[0].Msg, "multiple instances of match kind"),
		true,
	)
}

func TestRulesIsValidEmptySuccess(t *testing.T) {
	r := Rules{}

	assert.Nil(t, r.IsValid())
}

func TestRulesIsValidFailureBadMatches(t *testing.T) {
	r1, r2 := getRulesValidTestRules()
	badMatch := Match{
		Kind:     "whee",
		Behavior: ExactMatchBehavior,
		From:     Metadatum{Key: "other", Value: "v"},
		To:       Metadatum{Key: "flag", Value: "true"},
	}
	r2.Matches[1] = badMatch
	r := Rules{r1, r2}

	assert.NonNil(t, r.IsValid())
}

func TestRulesIsValidFailureBadNesting(t *testing.T) {
	r1, r2 := getRulesValidTestRules()
	badc := r2.Constraints
	badc.Light[0].Metadata = Metadata{Metadatum{"new-key", ""}}
	r2.Constraints = badc
	r2.Matches[0].Kind = "foo"
	r := Rules{r1, r2}

	assert.DeepEqual(t, r.IsValid(), &ValidationError{[]ErrorCase{
		{"rules[rkey1].matches[foo:exact:other].kind", `"foo" is not a valid match kind`},
		{"rules[rkey1].constraints.light[ck1].metadata[new-key].value", "must not be empty"},
	}})
}
