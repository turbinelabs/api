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

type SharedRulesKey string

/*
  SharedRules define mappings from a request to a pool of Instances, shared by
  a number of Routes. The left side of the mappings are defined by a vector of
  Rules.

  If none of the Rules applies to a given request, the Default
  AllConstraints are used; these define a default weighted set of
  Constraints. The weights determine the likelihood that one Constraint
  will be used over another.

  If one or more Rules applies, the order of the rules informs which is
  tried first. If a Rule fails to produce an Instance, the next applicable
  Rule is tried.

  It is possible to set a cohort seed on a SharedRules, Route, or Rule object.
  Only one of these will apply to any given request. SharedRules is the most
  generic of these objects and a seed set on either a Route or Rule will take
  precedence.

  See CohortSeed docs for additional details of what a cohort seed does.
*/
type SharedRules struct {
	SharedRulesKey SharedRulesKey `json:"shared_rules_key"` // overwritten for create
	Name           string         `json:"name"`
	ZoneKey        ZoneKey        `json:"zone_key"`
	Default        AllConstraints `json:"default"`
	Rules          Rules          `json:"rules"`
	ResponseData   ResponseData   `json:"response_data"`
	CohortSeed     *CohortSeed    `json:"cohort_seed"`
	Properties     Metadata       `json:"properties"`
	RetryPolicy    *RetryPolicy   `json:"retry_policy"`
	OrgKey         OrgKey         `json:"-"`
	Checksum
}

func (o SharedRules) GetZoneKey() ZoneKey   { return o.ZoneKey }
func (o SharedRules) GetOrgKey() OrgKey     { return o.OrgKey }
func (o SharedRules) Key() string           { return string(o.SharedRulesKey) }
func (o SharedRules) GetChecksum() Checksum { return o.Checksum }

type SharedRulesSlice []SharedRules

func (r SharedRules) IsNil() bool {
	return r.Equals(SharedRules{})
}

// Checks for exact equality between this SharedRules and another. Exact
// equality means each field must be equal (== or Equal, as appropriate) to the
// corresponding field in the parameter.
func (r SharedRules) Equals(o SharedRules) bool {
	var (
		eqKey  = r.SharedRulesKey == o.SharedRulesKey
		eqZone = r.ZoneKey == o.ZoneKey
		eqName = r.Name == o.Name
		eqCS   = r.Checksum.Equals(o.Checksum)
		eqOrg  = r.OrgKey == o.OrgKey
		eqRd   = r.ResponseData.Equals(o.ResponseData)
		eqCs   = CohortSeedPtrEquals(r.CohortSeed, o.CohortSeed)
		eqPr   = r.Properties.Equals(o.Properties)
		eqRp   = RetryPolicyEquals(r.RetryPolicy, o.RetryPolicy)
	)

	if !(eqKey && eqName && eqZone && eqCS && eqOrg && eqRd && eqCs && eqPr && eqRp) {
		return false
	}

	return r.Rules.Equals(o.Rules) && r.Default.Equals(o.Default)
}

// Checks validity of a SharedRules. For a route to be valid it must have a non-empty
// SharedRulesKey (or be precreation), have a ZoneKey, a Path, and valid Default +
// Rules.
func (r SharedRules) IsValid() *ValidationError {
	scope := func(s string) string { return "shared_rules." + s }
	errs := &ValidationError{}

	errCheckKey(string(r.SharedRulesKey), errs, scope("shared_rules_key"))
	errCheckIndex(r.Name, errs, scope("name"))
	errCheckKey(string(r.ZoneKey), errs, scope("zone_key"))

	errs.MergePrefixed(r.Default.IsValid("default"), "shared_rules")
	errs.MergePrefixed(r.Rules.IsValid(), "shared_rules")
	errs.MergePrefixed(r.ResponseData.IsValid(), scope("response_data"))
	if r.CohortSeed != nil {
		errs.MergePrefixed(r.CohortSeed.IsValid(), "shared_rules")
	}
	errs.MergePrefixed(SharedRulesPropertiesValid(r.Properties), "shared_rules")
	if r.RetryPolicy != nil {
		errs.MergePrefixed(r.RetryPolicy.IsValid(), "shared_rules")
	}

	return errs.OrNil()
}

// SharedRulesPropertiesValid ensures that the metadata has no duplicate
// or empty keys.
func SharedRulesPropertiesValid(m Metadata) *ValidationError {
	return MetadataValid("properties", m, MetadataCheckNonEmptyKeys)
}
