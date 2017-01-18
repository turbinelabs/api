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
	"testing"

	"github.com/turbinelabs/test/assert"
)

func getDomains() (Domain, Domain) {
	d := Domain{
		"dkey",
		"zkey",
		"name",
		1234,
		Redirects{{"redir", ".*", "http://www.example.com", PermanentRedirect}},
		"okey",
		Checksum{"aoeusnth"},
	}
	return d, d
}

func TestDomainEqualsSuccess(t *testing.T) {
	d1, d2 := getDomains()

	assert.True(t, d2.Equals(d1))
	assert.True(t, d1.Equals(d2))
	assert.True(t, d2.Equivalent(d1))
	assert.True(t, d1.Equivalent(d2))
}

func TestDomainEqualsOrgKeyVaries(t *testing.T) {
	d1, d2 := getDomains()
	d2.OrgKey = "okey2"

	assert.False(t, d2.Equals(d1))
	assert.False(t, d1.Equals(d2))
	assert.False(t, d2.Equivalent(d1))
	assert.False(t, d1.Equivalent(d2))
}

func TestDomainEquivalentVsEquals(t *testing.T) {
	d1, d2 := getDomains()
	d2.DomainKey = "dkey2"
	d2.Checksum = Checksum{"aoeu"}

	assert.False(t, d2.Equals(d1))
	assert.False(t, d1.Equals(d2))
	assert.True(t, d2.Equivalent(d1))
	assert.True(t, d1.Equivalent(d2))
}

func TestDomainNotEqualsKeyVaries(t *testing.T) {
	d1, d2 := getDomains()
	d2.DomainKey = "dkey2"

	assert.False(t, d1.Equals(d2))
	assert.False(t, d2.Equals(d1))
}

func TestDomainNotEqualsZoneKeyVaries(t *testing.T) {
	d1, d2 := getDomains()
	d2.ZoneKey = "zkey2"

	assert.False(t, d1.Equals(d2))
	assert.False(t, d2.Equals(d1))
}

func TestDomainNotEqualsNameVaries(t *testing.T) {
	d1, d2 := getDomains()
	d2.Name = "name2"

	assert.False(t, d1.Equals(d2))
	assert.False(t, d2.Equals(d1))
}

func TestDomainNotEqualsPortVaries(t *testing.T) {
	d1, d2 := getDomains()
	d2.Port = 12342

	assert.False(t, d1.Equals(d2))
	assert.False(t, d2.Equals(d1))
}

func getDomain() Domain {
	return Domain{
		"dkey",
		"zk",
		"name",
		1234,
		Redirects{{"sample-redirect", ".*", "http://www.example.com", PermanentRedirect}},
		"okey",
		Checksum{},
	}
}

func TestDomainIsValidSuccessPreCreation(t *testing.T) {
	d1 := getDomain()
	d1.DomainKey = ""

	assert.Nil(t, d1.IsValid(true))
	assert.NonNil(t, d1.IsValid(false))
}

func TestDomainIsValidSuccess(t *testing.T) {
	d1 := getDomain()

	assert.Nil(t, d1.IsValid(true))
	assert.Nil(t, d1.IsValid(false))
}

func TestDomainIsValidFailedDkey(t *testing.T) {
	d1 := getDomain()
	d1.DomainKey = ""

	assert.Nil(t, d1.IsValid(true))
	assert.NonNil(t, d1.IsValid(false))
}

func TestDomainIsValidFailedName(t *testing.T) {
	d1 := getDomain()
	d1.Name = ""

	assert.NonNil(t, d1.IsValid(true))
	assert.NonNil(t, d1.IsValid(false))
}

func TestDomainIsValidFailedPort(t *testing.T) {
	d1 := getDomain()
	d1.Port = 0

	assert.NonNil(t, d1.IsValid(true))
	assert.NonNil(t, d1.IsValid(false))
}

func TestDomainIsValidFailedRedirect(t *testing.T) {
	d1 := getDomain()
	d1.Redirects[0].To = ""
	err := d1.IsValid(true)
	assert.NonNil(t, err)
	assert.HasSameElements(t, err.Errors, []ErrorCase{
		{"domain[dkey].redirects[sample-redirect].to", "must not be empty"},
	})
}

func getThreeDomains() (Domain, Domain, Domain) {
	d1 := Domain{"dkey-1", "zk", "name", 10, nil, "okey", Checksum{}}
	d2 := Domain{"dkey-2", "zk", "name", 20, nil, "okey", Checksum{}}
	d3 := Domain{"dkey-3", "zk", "name", 30, nil, "okey", Checksum{}}

	return d1, d2, d3
}

func TestDomainsEqualsSuccess(t *testing.T) {
	d1, d2, d3 := getThreeDomains()
	ds1 := Domains{d1, d2, d3}
	ds2 := Domains{d1, d2, d3}

	assert.True(t, ds1.Equals(ds2))
	assert.True(t, ds2.Equals(ds1))
}

func TestDomainsEqualsOrderSuccess(t *testing.T) {
	d1, d2, d3 := getThreeDomains()
	ds1 := Domains{d3, d2, d1}
	ds2 := Domains{d3, d1, d2}

	assert.True(t, ds1.Equals(ds2))
	assert.True(t, ds2.Equals(ds1))
}

func TestDomainsEqualsFailure(t *testing.T) {
	d1, d2, d3 := getThreeDomains()
	ds1 := Domains{d3, d2, d1}
	ds2 := Domains{d3, d1}

	assert.False(t, ds1.Equals(ds2))
	assert.False(t, ds2.Equals(ds1))
}

func TestDomainsIsValidSuccess(t *testing.T) {
	d1 := Domain{"dkey-1", "zk", "name", 10, nil, "okey", Checksum{}}
	d2 := Domain{"dkey-2", "zk", "name", 20, nil, "okey", Checksum{}}
	d3 := Domain{"dkey-3", "zk", "name", 30, nil, "okey", Checksum{}}
	ds := Domains{d3, d2, d1}

	assert.Nil(t, ds.IsValid(true))
	assert.Nil(t, ds.IsValid(false))
}

func TestDomainsIsValidFailureDupe(t *testing.T) {
	d1 := Domain{"dkey-1", "zk", "name", 10, nil, "okey", Checksum{}}
	d2 := Domain{"dkey-2", "zk", "name", 20, nil, "okey", Checksum{}}
	d3 := Domain{"dkey-3", "zk", "name", 30, nil, "okey", Checksum{}}
	ds := Domains{d3, d2, d1, d3}

	assert.NonNil(t, ds.IsValid(true))
	assert.NonNil(t, ds.IsValid(false))
}

func TestDomainsIsValidFailureBadDomain(t *testing.T) {
	d1 := Domain{"dkey-1", "zk", "name", 10, nil, "okey", Checksum{}}
	d2 := Domain{"dkey-2", "", "name", 20, nil, "okey", Checksum{}}
	d3 := Domain{"dkey-3", "zk", "name", 30, nil, "okey", Checksum{}}
	ds := Domains{d3, d2, d1}

	assert.NonNil(t, ds.IsValid(true))
	assert.NonNil(t, ds.IsValid(false))
}
