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
	"strings"

	tbnstrings "github.com/turbinelabs/nonstdlib/strings"
)

type DomainKey string

// A Domain represents the TLD or subdomain under which which a set of Routes is served.
type Domain struct {
	DomainKey   DomainKey   `json:"domain_key"` // overwritten for create
	ZoneKey     ZoneKey     `json:"zone_key"`
	Name        string      `json:"name"`
	Port        int         `json:"port"`
	Redirects   Redirects   `json:"redirects"`
	GzipEnabled bool        `json:"gzip_enabled"`
	CorsConfig  *CorsConfig `json:"cors_config"`
	OrgKey      OrgKey      `json:"-"`
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

	parent := fmt.Sprintf("domain[%v]", d.DomainKey)
	errs.MergePrefixed(d.Redirects.IsValid(), parent)

	if d.CorsConfig != nil {
		errs.MergePrefixed(d.CorsConfig.IsValid(), parent)
	}

	return errs.OrNil()
}

// Check if all fields of this domain are exactly equal to fields of another
// domain.
func (d Domain) Equals(o Domain) bool {
	dCCNil := d.CorsConfig == nil
	oCCNil := o.CorsConfig == nil

	if dCCNil != oCCNil {
		return false
	}
	ccEq := oCCNil || d.CorsConfig.Equals(*o.CorsConfig)

	return d.DomainKey == o.DomainKey &&
		d.ZoneKey == o.ZoneKey &&
		d.Name == o.Name &&
		d.Port == o.Port &&
		d.Checksum.Equals(o.Checksum) &&
		d.OrgKey == o.OrgKey &&
		d.GzipEnabled == o.GzipEnabled &&
		d.Redirects.Equals(o.Redirects) &&
		ccEq
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

// CorsConfig describes how the domain should respond to OPTIONS requests.
// For a detailed discussion of what each attribute means see
// https://developer.mozilla.org/docs/Web/HTTP/Access_control_CORS
type CorsConfig struct {
	AllowedOrigins   []string `json:"allowed_origins"`
	AllowCredentials bool     `json:"allow_credentials"`
	ExposedHeaders   []string `json:"exposed_headers"`
	MaxAge           int      `json:"max_age"`
	AllowedMethods   []string `json:"allowed_methods"`
	AllowedHeaders   []string `json:"allowed_headers"`
}

// Equals compares two CorsConfig objects returning true if they are the same.
// AllowedOrigins, ExposedHeaders, AllowedMethods, and AllowedHeaders are
// compared without regard for ordering of their content.
func (cc CorsConfig) Equals(o CorsConfig) bool {
	cmp := func(ccs, os []string) bool {
		s1 := tbnstrings.NewSet(ccs...)
		s2 := tbnstrings.NewSet(os...)
		return s1.Equals(s2)
	}

	return cc.MaxAge == o.MaxAge &&
		cc.AllowCredentials == o.AllowCredentials &&
		cmp(cc.AllowedOrigins, o.AllowedOrigins) &&
		cmp(cc.ExposedHeaders, o.ExposedHeaders) &&
		cmp(cc.AllowedMethods, o.AllowedMethods) &&
		cmp(cc.AllowedHeaders, o.AllowedHeaders)
}

// MethodString produces a comma-separated string for the AllowedMethods slice.
func (cc CorsConfig) MethodString() string {
	m := make([]string, len(cc.AllowedMethods))
	copy(m, cc.AllowedMethods)
	for i, j := range m {
		m[i] = strings.ToUpper(j)
	}

	return strings.Join(m, ", ")
}

// ExposedHeadersString produces a comma-separated string for the ExposedHeaders
// slice.
func (cc CorsConfig) ExposedHeadersString() string {
	m := make([]string, len(cc.ExposedHeaders))
	copy(m, cc.ExposedHeaders)
	return strings.Join(m, ", ")
}

// AllowHeadersString produces a comma-separated string for the AllowedHeaders
// slice.
func (cc CorsConfig) AllowHeadersString() string {
	m := make([]string, len(cc.AllowedHeaders))
	copy(m, cc.AllowedHeaders)
	return strings.Join(m, ", ")
}

var isAllowedMethod = map[string]bool{
	"GET":    true,
	"HEAD":   true,
	"PUT":    true,
	"POST":   true,
	"DELETE": true,
}

// IsValid checks a CorsConfig object for validity.
func (cc CorsConfig) IsValid() *ValidationError {
	errs := &ValidationError{}
	ec := func(f, m string) ErrorCase {
		return ErrorCase{"cors_config." + f, m}
	}

	lao := len(cc.AllowedOrigins)
	if lao == 0 {
		errs.AddNew(ec("allowed_origins", "must have at least one element"))
	}

	if lao > 1 {
		// temporary until we build a more powerful proxy plugin instead of relying
		// on config-only solution
		errs.AddNew(ec(
			"allowed_origins",
			"currently Allowed-Origins only supports wildcard or a single target"))

		if tbnstrings.NewSet(cc.AllowedOrigins...).Contains("*") {
			errs.AddNew(ec("allowed_origins", "may not mix wildcard (*) with specific origins"))
		}
	}

	if len(cc.AllowedMethods) == 0 {
		errs.AddNew(ec("allow_methods", "must have at least one element"))
	}

	for _, m := range cc.AllowedMethods {
		if !isAllowedMethod[m] {
			errs.AddNew(ec("allow_methods", fmt.Sprintf("%s is not a valid method", m)))
		}
	}

	if cc.MaxAge < 0 {
		errs.AddNew(ec("max_age", "must be greater than or equal to 0"))
	}

	return errs.OrNil()
}
