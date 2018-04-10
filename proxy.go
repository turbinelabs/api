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
	"fmt"
)

type Proxies []Proxy

type ProxyKey string

// A Proxy is a named Instance, responsible for serving one or more Domains.
type Proxy struct {
	ProxyKey   ProxyKey    `json:"proxy_key"` // overwritten on create
	ZoneKey    ZoneKey     `json:"zone_key"`
	Name       string      `json:"name"`
	DomainKeys []DomainKey `json:"domain_keys"`
	OrgKey     OrgKey      `json:"-"`
	Checksum
}

func (p Proxy) GetZoneKey() ZoneKey   { return p.ZoneKey }
func (p Proxy) GetOrgKey() OrgKey     { return p.OrgKey }
func (p Proxy) Key() string           { return string(p.ProxyKey) }
func (p Proxy) GetChecksum() Checksum { return p.Checksum }

// Check validity of a new or existing proxy. A Valid proxy requires a
// ProxyKey (unless new), a ZoneKey, and valid sub objects (Instance & Domains).
func (p Proxy) IsValid() *ValidationError {
	scope := func(n string) string { return "proxy." + n }
	errs := &ValidationError{}

	errCheckKey(string(p.ProxyKey), errs, scope("proxy_key"))
	errCheckKey(string(p.ZoneKey), errs, scope("zone_key"))
	errCheckIndex(string(p.Name), errs, scope("name"))

	seenDomain := map[string]bool{}
	for _, dk := range p.DomainKeys {
		sdk := string(dk)
		if seenDomain[sdk] {
			errs.AddNew(ErrorCase{scope("domain_keys"), fmt.Sprintf("duplicate domain key '%v'", sdk)})
		}
		seenDomain[sdk] = true
		errCheckKey(sdk, errs, fmt.Sprintf("proxy.domain_keys[%v]", sdk))
	}

	errCheckKey(string(p.OrgKey), errs, scope("org_key"))

	return errs.OrNil()
}

func (p Proxy) IsNil() bool {
	return p.Equals(Proxy{})
}

// Check if one Proxy exactly equals another. This checks all fields (including
// derived fields).
func (p Proxy) Equals(o Proxy) bool {
	var (
		keyEq  = p.ProxyKey == o.ProxyKey
		zoneEq = p.ZoneKey == o.ZoneKey
		nameEq = p.Name == o.Name
		orgEq  = p.OrgKey == o.OrgKey
		csEq   = p.Checksum.Equals(o.Checksum)
	)

	if !(keyEq && zoneEq && nameEq && csEq && orgEq) {
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
