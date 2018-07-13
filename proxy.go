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

// A Proxy is a named Instance, responsible for serving one or more Listeners.
//
// Current behavior:
// For backwards compatibility, when the DomainKeys field is populated it indicates
// that the specified Domains should be attached to Listeners with matching ports.
// If no such Listener exists the consumer should create a default one to support the
// Domain. Note that because Listeners allow specification of an address to bind to,
// and Domains do not, is possible to create a confusing configuration where multiple
// Listeners are configured for a Proxy port, and we can't determine which of the
// Listeners a Domain should be attached to. In this case consumers should attach
// the Domain to _every_ Listener configured for the given port. Note that Domains
// cannot be mapped to Listeners that are configured for non-http protocols.
//
// Future behavior:
// In the future we will remove the DomainKeys field from the Proxy. Proxies will create
// the Listeners indicated in ListenerKeys, which in turn will have mapped Domains.
type Proxy struct {
	ProxyKey     ProxyKey      `json:"proxy_key"` // overwritten on create
	ZoneKey      ZoneKey       `json:"zone_key"`
	Name         string        `json:"name"`
	DomainKeys   []DomainKey   `json:"domain_keys"`
	ListenerKeys []ListenerKey `json:"listener_keys"`
	OrgKey       OrgKey        `json:"-"`
	Checksum
}

func (p Proxy) GetZoneKey() ZoneKey   { return p.ZoneKey }
func (p Proxy) GetOrgKey() OrgKey     { return p.OrgKey }
func (p Proxy) Key() string           { return string(p.ProxyKey) }
func (p Proxy) GetChecksum() Checksum { return p.Checksum }

// Check validity of a new or existing proxy. A Valid proxy requires a
// ProxyKey (unless new), a ZoneKey, and valid sub objects (Instance, Listeners & Domains).
func (p Proxy) IsValid() *ValidationError {
	scope := func(n string) string { return "proxy." + n }
	errs := &ValidationError{}

	errCheckKey(string(p.ProxyKey), errs, scope("proxy_key"))
	errCheckKey(string(p.ZoneKey), errs, scope("zone_key"))
	errCheckIndex(string(p.Name), errs, scope("name"))

	seenListener := map[string]bool{}
	for _, lk := range p.ListenerKeys {
		ldk := string(lk)
		if seenListener[ldk] {
			errs.AddNew(ErrorCase{scope("listener_keys"), fmt.Sprintf("duplicate listener key '%v'", ldk)})
		}
		seenListener[ldk] = true
		errCheckKey(ldk, errs, fmt.Sprintf("proxy.listener_keys[%v]", ldk))
	}

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
