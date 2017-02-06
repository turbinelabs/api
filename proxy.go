package api

import (
	"fmt"
)

type Proxies []Proxy

type ProxyKey string

// A Proxy is a named Instance, responsible for serving one or more Domains.
type Proxy struct {
	Instance               // TODO: we should replace Instance with HostSpecifier or something
	ProxyKey   ProxyKey    `json:"proxy_key"` // overwritten on create
	ZoneKey    ZoneKey     `json:"zone_key"`
	Name       string      `json:"name"`
	DomainKeys []DomainKey `json:"domain_keys"`
	OrgKey     OrgKey      `json:"-"`
	Checksum
}

// Check validity of a new or existing proxy. A Valid proxy requires a
// ProxyKey (unless new), a ZoneKey, and valid sub objects (Instance & Domains).
func (p Proxy) IsValid(precreation bool) *ValidationError {
	ecase := func(f, m string) ErrorCase {
		return ErrorCase{fmt.Sprintf("proxy.%s", f), m}
	}

	errs := &ValidationError{}

	keyValid := precreation || p.ProxyKey != ""
	zoneValid := p.ZoneKey != ""
	nameValid := p.Name != ""

	if !keyValid {
		errs.AddNew(ecase("proxy_key", "must not be empty"))
	}

	if !zoneValid {
		errs.AddNew(ecase("zone_key", "must not be empty"))
	}

	if !nameValid {
		errs.AddNew(ecase("name", "must not be empty"))
	}

	errs.MergePrefixed(p.Instance.IsValid(precreation), "")

	return errs.OrNil()
}

func (p Proxy) IsNil() bool {
	return p.Equals(Proxy{})
}

// Check if one Proxy exactly equals another. This checks all fields (including
// derived fields).
func (p Proxy) Equals(o Proxy) bool {
	var (
		instEq = p.Instance.Equals(o.Instance)
		keyEq  = p.ProxyKey == o.ProxyKey
		zoneEq = p.ZoneKey == o.ZoneKey
		nameEq = p.Name == o.Name
		orgEq  = p.OrgKey == o.OrgKey
		csEq   = p.Checksum.Equals(o.Checksum)
	)

	if !(instEq && keyEq && zoneEq && nameEq && csEq && orgEq) {
		return false
	}

	if len(p.DomainKeys) != len(o.DomainKeys) {
		return false
	}

	hasDomain := make(map[DomainKey]bool)

	for _, dk := range p.DomainKeys {
		hasDomain[dk] = true
	}

	for _, dk := range o.DomainKeys {
		if !hasDomain[dk] {
			return false
		}
	}

	return true
}
