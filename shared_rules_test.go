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
				23}}}

	rule1, rule2 := getRulesDefaults()
	rules := Rules{rule1, rule2}

	r1 := SharedRules{
		"routekey",
		"rule1-name",
		"zkey",
		defaultCC,
		rules,
		"1",
		Checksum{"cs-1"},
	}

	r2 := SharedRules{
		"routekey",
		"rule1-name",
		"zkey",
		defaultCC,
		rules,
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
			ClusterConstraint{"cckey1", HeaderMatchKind, Metadata{{"k1", "v1"}, {"k2", "v2"}}, nil, 23},
			ClusterConstraint{"cckey2", HeaderMatchKind, Metadata{{"k2", "v2"}, {"k2", "v2"}}, nil, 23}}}

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

	assert.Nil(t, r.IsValid(true))
	assert.Nil(t, r.IsValid(false))
}

func TestSharedRulesIsValidPrecreationOnlySuccess(t *testing.T) {
	r, _ := getSharedRulesDefaults()
	r.SharedRulesKey = ""

	assert.Nil(t, r.IsValid(true))
	assert.NonNil(t, r.IsValid(false))
}

func TestSharedRulesIsValidBadName(t *testing.T) {
	r, _ := getSharedRulesDefaults()
	r.Name = ""

	assert.NonNil(t, r.IsValid(true))
	assert.NonNil(t, r.IsValid(false))
}

func TestSharedRulesIsValidBadZoneKey(t *testing.T) {
	r, _ := getSharedRulesDefaults()
	r.ZoneKey = ""

	assert.NonNil(t, r.IsValid(true))
	assert.NonNil(t, r.IsValid(false))
}

func TestSharedRulesIsValidBadDefault(t *testing.T) {
	r, _ := getSharedRulesDefaults()
	r.Default = AllConstraints{}

	assert.NonNil(t, r.IsValid(true))
	assert.NonNil(t, r.IsValid(false))
}

func TestSharedRulesIsValidBadRules(t *testing.T) {
	r, _ := getSharedRulesDefaults()
	rule1, _ := getRulesDefaults()
	rule1.Matches = Matches{}
	rule1.Methods = []string{}
	r.Rules[0] = rule1

	assert.NonNil(t, r.IsValid(true))
	assert.NonNil(t, r.IsValid(false))
}
