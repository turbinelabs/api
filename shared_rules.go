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
func (r SharedRules) IsValid(precreation bool) *ValidationError {
	ecase := func(f, m string) ErrorCase {
		return ErrorCase{"shared_rules." + f, m}
	}

	errs := &ValidationError{}
	var (
		validName    = r.Name != ""
		validKey     = precreation || r.SharedRulesKey != ""
		validZoneKey = r.ZoneKey != ""
	)

	if !validKey {
		errs.AddNew(ecase("shared_rules_key", "must not be empty"))
	}

	if !validName {
		errs.AddNew(ecase("name", "must not be empty"))
	}

	if !validZoneKey {
		errs.AddNew(ecase("zone_key", "must not be empty"))
	}

	errs.MergePrefixed(r.Default.IsValid(precreation), "shared_rules.default")
	errs.MergePrefixed(r.Rules.IsValid(precreation), "shared_rules.rules")

	return errs.OrNil()
}
