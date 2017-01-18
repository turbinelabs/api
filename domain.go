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

import (
	"fmt"
)

type DomainKey string

// A Domain represents the TLD or subdomain under which which a set of Routes is served.
type Domain struct {
	DomainKey DomainKey `json:"domain_key"` // overwritten for create
	ZoneKey   ZoneKey   `json:"zone_key"`
	Name      string    `json:"name"`
	Port      int       `json:"port"`
	Redirects Redirects `json:"redirects"`
	OrgKey    OrgKey    `json:"-"`
	Checksum
}

func (d Domain) IsNil() bool {
	return d.Equals(Domain{})
}

type Domains []Domain

// Checks for validity of a domain. A domain is considered valid if it has a:
//  1. DomainKey OR is being checked in before creation
//  2. non empty ZoneKey
//  3. non empty Name
//  4. non zero Port
func (d Domain) IsValid(precreation bool) *ValidationError {
	ecase := func(f, m string) ErrorCase {
		return ErrorCase{fmt.Sprintf("domain[%s].%s", string(d.DomainKey), f), m}
	}

	errs := &ValidationError{}

	validDomainKey := precreation || d.DomainKey != ""
	validZoneKey := d.ZoneKey != ""
	validName := d.Name != ""
	validPort := d.Port != 0

	if !validDomainKey {
		errs.AddNew(ecase("domain_key", "must not be empty"))
	}

	if !validZoneKey {
		errs.AddNew(ecase("zone_key", "must not be empty"))
	}

	if !validName {
		errs.AddNew(ecase("name", "must not be empty"))
	}

	if !validPort {
		errs.AddNew(ecase("port", "must be non-zero"))
	}

	errs.MergePrefixed(
		d.Redirects.IsValid(),
		fmt.Sprintf("domain[%v]", d.DomainKey),
	)

	return errs.OrNil()
}

// Check if all fields of this domain are exactly equal to fields of another
// domain.
func (d Domain) Equals(o Domain) bool {
	return d.DomainKey == o.DomainKey &&
		d.ZoneKey == o.ZoneKey &&
		d.Name == o.Name &&
		d.Port == o.Port &&
		d.Checksum.Equals(o.Checksum) &&
		d.OrgKey == o.OrgKey &&
		d.Redirects.Equals(o.Redirects)
}

// Check for semantic equality between this Domain an another. Domains must
// have the same Name, Zone, and Port to be considered equivalent. Key and
// Checksum are explicitly excluded from requirements for equivalence.
func (d Domain) Equivalent(o Domain) bool {
	return d.ZoneKey == o.ZoneKey &&
		d.Name == o.Name &&
		d.Port == o.Port &&
		d.OrgKey == o.OrgKey
}

// Checks for exact contents parity between two Domains. This requires
// that each Domain with the same Key be Equal to each other.
func (ds Domains) Equals(o Domains) bool {
	if len(ds) != len(o) {
		return false
	}

	hasDomain := make(map[DomainKey]bool)

	for _, d := range ds {
		hasDomain[d.DomainKey] = true
	}

	for _, d := range o {
		if !hasDomain[d.DomainKey] {
			return false
		}
	}

	return true
}

// Checks validity of a domain slice. Requise that each domain is valid and
// that there are no domains with duplicate keys.
func (ds Domains) IsValid(precreation bool) *ValidationError {
	errs := &ValidationError{}

	keySeen := make(map[DomainKey]bool)

	for _, d := range ds {
		if keySeen[d.DomainKey] {
			errs.AddNew(ErrorCase{
				"domain_key",
				fmt.Sprintf("multiple instances of key %s", string(d.DomainKey)),
			})
		}

		keySeen[d.DomainKey] = true
		errs.MergePrefixed(d.IsValid(precreation), "")
	}

	return errs.OrNil()
}
