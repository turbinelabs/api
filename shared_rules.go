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
*/
type SharedRules struct {
	SharedRulesKey SharedRulesKey `json:"shared_rules_key"` // overwritten for create
	Name           string         `json:"name"`
	ZoneKey        ZoneKey        `json:"zone_key"`
	Default        AllConstraints `json:"default"`
	Rules          Rules          `json:"rules"`
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
	)

	if !(eqKey && eqName && eqZone && eqCS && eqOrg) {
		return false
	}

	return r.Rules.Equals(o.Rules) && r.Default.Equals(o.Default)
}

// Checks validity of a SharedRules. For a route to be valid it must have a non-empty
// SharedRulesKey (or be precreation), have a ZoneKey, a Path, and valid Default +
// Rules.
func (r SharedRules) IsValid() *ValidationError {
	errs := &ValidationError{}

	errCheckKey(string(r.SharedRulesKey), errs, "shared_rules_key")
	errCheckIndex(r.Name, errs, "name")
	errCheckKey(string(r.ZoneKey), errs, "zone_key")

	errs.MergePrefixed(r.Default.IsValid("default"), "shared_rules")
	errs.MergePrefixed(r.Rules.IsValid(), "shared_rules")

	return errs.OrNil()
}
