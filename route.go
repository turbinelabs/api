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
*/
type Route struct {
	RouteKey       RouteKey       `json:"route_key"` // overwritten for create
	DomainKey      DomainKey      `json:"domain_key"`
	ZoneKey        ZoneKey        `json:"zone_key"`
	Path           string         `json:"path"`
	SharedRulesKey SharedRulesKey `json:"shared_rules_key"`
	Rules          Rules          `json:"rules"`
	OrgKey         OrgKey         `json:"-"`
	Checksum
}

func (r Route) IsNil() bool {
	return r.Equals(Route{})
}

// Checks for exact equality between this route and another. Exactly equality
// means each field must be equal (== or Equal, as appropriate) to the
// corresponding field in the parameter.
func (r Route) Equals(o Route) bool {
	var (
		eqKey   = r.RouteKey == o.RouteKey
		eqDom   = r.DomainKey == o.DomainKey
		eqZone  = r.ZoneKey == o.ZoneKey
		eqPath  = r.Path == o.Path
		eqCS    = r.Checksum.Equals(o.Checksum)
		eqOrg   = r.OrgKey == o.OrgKey
		eqSRKey = r.SharedRulesKey == o.SharedRulesKey
	)

	if !(eqKey && eqDom && eqZone && eqPath && eqCS && eqOrg && eqSRKey) {
		return false
	}

	return r.Rules.Equals(o.Rules)
}

// Checks validity of a Route. For a route to be valid it must have a non-empty
// RouteKey (or be precreation), have a DomainKey, a ZoneKey, a Path, and valid
// Default + Rules.
func (r Route) IsValid(precreation bool) *ValidationError {
	ecase := func(f, m string) ErrorCase {
		return ErrorCase{"route." + f, m}
	}

	errs := &ValidationError{}
	var (
		validKey            = precreation || r.RouteKey != ""
		validDomainKey      = r.DomainKey != ""
		validSharedRulesKey = r.SharedRulesKey != ""
		validZoneKey        = r.ZoneKey != ""
		validPath           = r.Path != ""
	)

	if !validSharedRulesKey {
		errs.AddNew(ecase("shared_rules_key", "must not be empty"))
	}

	if !validKey {
		errs.AddNew(ecase("route_key", "must not be empty"))
	}

	if !validDomainKey {
		errs.AddNew(ecase("domain_key", "must not be empty"))
	}

	if !validZoneKey {
		errs.AddNew(ecase("zone_key", "must not be empty"))
	}

	if !validPath {
		errs.AddNew(ecase("path", "must not be empty"))
	}

	errs.MergePrefixed(r.Rules.IsValid(precreation), "route.rules")

	return errs.OrNil()
}
