package api

import (
	"testing"

	"github.com/turbinelabs/test/assert"
)

func getRules() (Rule, Rule) {
	r1 := Rule{
		"rkey1",
		[]string{"GET", "POST"},
		Matches{
			Match{HeaderMatchKind, Metadatum{"x-random", "value"}, Metadatum{"randomflag", "true"}},
			Match{CookieMatchKind, Metadatum{"x-other", "value"}, Metadatum{"otherflag", "true"}}},
		AllConstraints{
			Light: ClusterConstraints{
				ClusterConstraint{
					"cckey1",
					"ckey1",
					Metadata{{"key", "value"}, {"key2", "value2"}},
					nil,
					1234},
				ClusterConstraint{
					"cckey2",
					"ckey2",
					Metadata{{"key-2", "value-2"}},
					Metadata{{"state", "testing"}},
					1234}}},
	}

	r2 := Rule{
		"rkey1",
		[]string{"POST", "GET"},
		Matches{
			Match{CookieMatchKind, Metadatum{"x-other", "value"}, Metadatum{"otherflag", "true"}},
			Match{HeaderMatchKind, Metadatum{"x-random", "value"}, Metadatum{"randomflag", "true"}}},
		AllConstraints{
			Light: ClusterConstraints{
				ClusterConstraint{
					"cckey1",
					"ckey1",
					Metadata{{"key", "value"}, {"key2", "value2"}},
					nil,
					1234},
				ClusterConstraint{
					"cckey2",
					"ckey2",
					Metadata{{"key-2", "value-2"}},
					Metadata{{"state", "testing"}},
					1234}}},
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
			ClusterConstraint{"cckey1", "ckey2", Metadata{{"key-2", "value-2"}}, nil, 1234}},
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

	assert.Nil(t, r.IsValid(true))
	assert.Nil(t, r.IsValid(false))
}

func TestRuleIsValidBadRuleKey(t *testing.T) {
	r := getRuleValid()
	r.RuleKey = ""

	assert.NonNil(t, r.IsValid(true))
	assert.NonNil(t, r.IsValid(false))
}

func TestRuleIsValidNoMethodOrMatches(t *testing.T) {
	r := getRuleValid()
	r.Matches = Matches{}
	r.Methods = []string{}
	assert.NonNil(t, r.IsValid(true))
	assert.NonNil(t, r.IsValid(false))
}

func TestRuleIsValidBadMethod(t *testing.T) {
	r := getRuleValid()
	r.Methods = []string{"POST", "PUT", "GRAB_THAT_RESOURCE"}

	assert.NonNil(t, r.IsValid(true))
	assert.NonNil(t, r.IsValid(false))
}

func TestRuleIsValidBadMatches(t *testing.T) {
	r := getRuleValid()
	r.Matches = Matches{
		Match{CookieMatchKind, Metadatum{"x-other", "value"}, Metadatum{"otherflag", "true"}},
		Match{HeaderMatchKind, Metadatum{"x-random", "value"}, Metadatum{"", "aoeu"}},
	}

	assert.NonNil(t, r.IsValid(true))
	assert.NonNil(t, r.IsValid(false))
}

func TestRuleIsValidBadConstraints(t *testing.T) {
	r := getRuleValid()
	r.Constraints = AllConstraints{
		Dark: ClusterConstraints{{"cckey0", "ckey2", Metadata{{"key-2", "value-2"}}, Metadata{{"aoeu", "snth"}}, 1234}}}

	assert.NonNil(t, r.IsValid(true))
	assert.NonNil(t, r.IsValid(false))
}

func getRulesValidTestRules() (Rule, Rule) {
	r1 := Rule{
		"rkey0",
		[]string{"POST", "PUT"},
		Matches{
			Match{CookieMatchKind, Metadatum{"x-other", "value"}, Metadatum{"otherflag", "true"}},
			Match{HeaderMatchKind, Metadatum{"x-random", "value"}, Metadatum{"randomflag", "true"}}},
		AllConstraints{
			Light: ClusterConstraints{{"ck0", "ckey2", Metadata{{"key-2", "value-2"}}, nil, 1234}},
		},
	}

	r2 := Rule{
		"rkey1",
		[]string{"GET"},
		Matches{
			Match{CookieMatchKind, Metadatum{"other", "v"}, Metadatum{"flag", "true"}},
			Match{HeaderMatchKind, Metadatum{"random", "v"}, Metadatum{"random", "true"}}},
		AllConstraints{
			Light: ClusterConstraints{{"ck1", "ckey2", Metadata{{"key-2", "value-2"}}, nil, 1234}},
		},
	}

	return r1, r2
}

func TestRulesIsValidSucces(t *testing.T) {
	r1, r2 := getRulesValidTestRules()
	r := Rules{r1, r2}

	assert.Nil(t, r.IsValid(true))
	assert.Nil(t, r.IsValid(false))
}

func TestRulesIsValidFailureOnDupeKey(t *testing.T) {
	r1, r2 := getRulesValidTestRules()
	r2.RuleKey = r1.RuleKey
	r := Rules{r1, r2}

	assert.NonNil(t, r.IsValid(true))
	assert.NonNil(t, r.IsValid(false))
}

func TestRulesIsValidEmptySuccess(t *testing.T) {
	r := Rules{}

	assert.Nil(t, r.IsValid(true))
	assert.Nil(t, r.IsValid(false))
}

func TestRulesIsValidFailureBadMatches(t *testing.T) {
	r1, r2 := getRulesValidTestRules()
	badMatch := Match{"whee", Metadatum{"other", "v"}, Metadatum{"flag", "true"}}
	r2.Matches[1] = badMatch
	r := Rules{r1, r2}

	assert.NonNil(t, r.IsValid(true))
	assert.NonNil(t, r.IsValid(false))
}

func TestRulesIsValidFailureBadConstraints(t *testing.T) {
	r1, r2 := getRulesValidTestRules()
	badc := r2.Constraints
	badc.Dark = badc.Light
	badc.Light = ClusterConstraints{}
	r2.Constraints = badc

	r := Rules{r1, r2}

	assert.NonNil(t, r.IsValid(true))
	assert.NonNil(t, r.IsValid(false))
}
