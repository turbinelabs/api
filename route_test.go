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
				ClusterConstraint{"cckey1", "ckey2", Metadata{{"key-2", "value-2"}}, nil, 1234}}},
	}

	rule2 := Rule{
		"rk1",
		[]string{"PUT", "DELETE"},
		Matches{
			Match{CookieMatchKind, Metadatum{"x-2", "value"}, Metadatum{"other", "true"}}},
		AllConstraints{
			Tap: ClusterConstraints{
				ClusterConstraint{"cckey1", "ckey3", Metadata{{"key-2", "value-2"}}, nil, 1234}},
			Light: ClusterConstraints{
				ClusterConstraint{"cckey2", "ckey2", Metadata{{"key-2", "value-2"}}, nil, 1234}}},
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

	assert.Nil(t, r.IsValid(true))
	assert.Nil(t, r.IsValid(false))
}

func TestRouteIsValidPrecreationOnlySuccess(t *testing.T) {
	r, _ := getRouteDefaults()
	r.RouteKey = ""

	assert.Nil(t, r.IsValid(true))
	assert.NonNil(t, r.IsValid(false))
}

func TestRouteIsValidBadDomainKey(t *testing.T) {
	r, _ := getRouteDefaults()
	r.DomainKey = ""

	assert.NonNil(t, r.IsValid(true))
	assert.NonNil(t, r.IsValid(false))
}

func TestRouteIsValidBadZoneKey(t *testing.T) {
	r, _ := getRouteDefaults()
	r.ZoneKey = ""

	assert.NonNil(t, r.IsValid(true))
	assert.NonNil(t, r.IsValid(false))
}

func TestRouteIsValidBadPath(t *testing.T) {
	r, _ := getRouteDefaults()
	r.Path = ""

	assert.NonNil(t, r.IsValid(true))
	assert.NonNil(t, r.IsValid(false))
}

func TestRouteIsValidBadSharedRulesKeyRef(t *testing.T) {
	r, _ := getRouteDefaults()
	r.SharedRulesKey = ""

	assert.NonNil(t, r.IsValid(true))
	assert.NonNil(t, r.IsValid(false))
}

func TestRouteIsValidBadRules(t *testing.T) {
	r, _ := getRouteDefaults()
	rule1, _ := getRulesDefaults()
	rule1.Matches = Matches{}
	rule1.Methods = []string{}
	r.Rules[0] = rule1

	assert.NonNil(t, r.IsValid(true))
	assert.NonNil(t, r.IsValid(false))
}
