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

type Routes []Route

type RouteKey string

/*
	A Route defines a mapping from a request to a pool of Instances.
	The left side of the mapping is defined by a Zone, a Domain, a Path,
	and a vector of Rules.

	If none of the Rules applies to a given request, the Default
	AllConstraints are used; these define a default weighted set of
	Constraints. The weights determine the likelihood that one Constraint
	will be used over another.

	If one or more Rules applies, the order of the rules informs which is
	tried first. If a Rule fails to produce an Instance, the next applicable
	Rule is tried.

	It is possible to set a cohort seed on a SharedRules, Route, or Rule object.
	Only one of these will apply to any given request. Route is a subset of
	SharedRules and not as specific as a Rule so a seed specified by the Route
	will override one inherited from a SharedRules but not one specified in a
	Rule (regardless of the Rule source).

	See CohortSeed docs for additional details of what a cohort seed does.
*/
type Route struct {
	RouteKey       RouteKey       `json:"route_key"` // overwritten for create
	DomainKey      DomainKey      `json:"domain_key"`
	ZoneKey        ZoneKey        `json:"zone_key"`
	Path           string         `json:"path"`
	SharedRulesKey SharedRulesKey `json:"shared_rules_key"`
	Rules          Rules          `json:"rules"`
	ResponseData   ResponseData   `json:"response_data"`
	CohortSeed     *CohortSeed    `json:"cohort_seed"`
	RetryPolicy    *RetryPolicy   `json:"retry_policy"`
	OrgKey         OrgKey         `json:"-"`
	Checksum
}

func (o Route) GetZoneKey() ZoneKey   { return o.ZoneKey }
func (o Route) GetOrgKey() OrgKey     { return o.OrgKey }
func (o Route) Key() string           { return string(o.RouteKey) }
func (o Route) GetChecksum() Checksum { return o.Checksum }

func (r Route) IsNil() bool {
	return r.Equals(Route{})
}

// Checks for exact equality between this route and another. Exactly equality
// means each field must be equal (== or Equal, as appropriate) to the
// corresponding field in the parameter.
func (r Route) Equals(o Route) bool {
	var (
		eqKey    = r.RouteKey == o.RouteKey
		eqDom    = r.DomainKey == o.DomainKey
		eqZone   = r.ZoneKey == o.ZoneKey
		eqPath   = r.Path == o.Path
		eqCS     = r.Checksum.Equals(o.Checksum)
		eqOrg    = r.OrgKey == o.OrgKey
		eqSRKey  = r.SharedRulesKey == o.SharedRulesKey
		eqRd     = r.ResponseData.Equals(o.ResponseData)
		eqCohort = CohortSeedPtrEquals(r.CohortSeed, o.CohortSeed)
		eqRp     = RetryPolicyEquals(r.RetryPolicy, o.RetryPolicy)
	)

	if !(eqKey && eqDom && eqZone && eqPath && eqCS &&
		eqOrg && eqSRKey && eqRd && eqCohort && eqRp) {
		return false
	}

	return r.Rules.Equals(o.Rules)
}

// Checks validity of a Route. For a route to be valid it must have a non-empty
// RouteKey (or be precreation), have a DomainKey, a ZoneKey, a Path, and valid
// Default + Rules.
func (r Route) IsValid() *ValidationError {
	scope := func(s string) string { return "route." + s }

	errs := &ValidationError{}
	errCheckKey(string(r.RouteKey), errs, scope("route_key"))
	errCheckKey(string(r.SharedRulesKey), errs, scope("shared_rules_key"))
	errCheckKey(string(r.DomainKey), errs, scope("domain_key"))
	errCheckKey(string(r.ZoneKey), errs, scope("zone_key"))

	if r.Path == "" {
		errs.AddNew(ErrorCase{scope("path"), "must not be empty"})
	}

	errs.MergePrefixed(r.Rules.IsValid(), "route")
	errs.MergePrefixed(r.ResponseData.IsValid(), scope("response_data"))
	if r.CohortSeed != nil {
		errs.MergePrefixed(r.CohortSeed.IsValid(), "route")
	}
	if r.RetryPolicy != nil {
		errs.MergePrefixed(r.RetryPolicy.IsValid(), "route")
	}

	return errs.OrNil()
}
