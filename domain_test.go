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
	"testing"

	"github.com/turbinelabs/test/assert"
)

func getDomains() (Domain, Domain) {
	d := Domain{
		"dkey",
		"zkey",
		"name",
		1234,
		Redirects{{
			"redir",
			".*",
			"http://www.example.com",
			PermanentRedirect,
			HeaderConstraints{{"x-tbn-api-key", "", false, false}},
		}},
		true,
		mkCC(),
		DomainAliases{},
		"okey",
		Checksum{"aoeusnth"},
	}
	d2 := d
	d2.CorsConfig = mkCC()
	return d, d2
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

func TestDomainEqualsRedirctChanged(t *testing.T) {
	d1, d2 := getDomains()
	d2.Redirects = make(Redirects, len(d1.Redirects))
	copy(d2.Redirects, d1.Redirects)
	d2.Redirects[0].From = "aoeu"
	assert.False(t, d2.Equals(d1))
	assert.False(t, d1.Equals(d2))
}

func TestDomainEqualsGzipChanged(t *testing.T) {
	d1, d2 := getDomains()
	d2.GzipEnabled = false
	assert.False(t, d2.Equals(d1))
	assert.False(t, d1.Equals(d2))
}

func TestDomainEqualsCorsConfigNilNil(t *testing.T) {
	d1, d2 := getDomains()
	d1.CorsConfig = nil
	d2.CorsConfig = nil

	assert.True(t, d2.Equals(d1))
	assert.True(t, d1.Equals(d2))
}

func TestDomainEqualsCorsConfigSomeNil(t *testing.T) {
	d1, d2 := getDomains()
	d2.CorsConfig = nil

	assert.False(t, d2.Equals(d1))
	assert.False(t, d1.Equals(d2))
}

func TestDomainEqualsCorsConfigChanges(t *testing.T) {
	d1, d2 := getDomains()
	d2.CorsConfig.AllowCredentials = !d1.CorsConfig.AllowCredentials

	assert.False(t, d2.Equals(d1))
	assert.False(t, d1.Equals(d2))
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
	cc := mkCC()
	cc.AllowedOrigins = []string{"*"}
	return Domain{
		"dkey",
		"zk",
		"name",
		1234,
		Redirects{{
			"sample-redirect",
			".*",
			"http://www.example.com",
			PermanentRedirect,
			HeaderConstraints{{"x-tbn-api-key", "", false, false}},
		}},
		true,
		cc,
		DomainAliases{},
		"okey",
		Checksum{},
	}
}

func TestDomainIsValidSuccess(t *testing.T) {
	d1 := getDomain()

	assert.Nil(t, d1.IsValid())
}

func TestDomainIsValidNoKey(t *testing.T) {
	d := getDomain()
	d.DomainKey = ""
	assert.NonNil(t, d.IsValid())
}

func TestDomainIsValidBadKey(t *testing.T) {
	d := getDomain()
	d.DomainKey = "-aoeu"
	assert.NonNil(t, d.IsValid())
}

func TestDomainIsValidNoName(t *testing.T) {
	d := getDomain()
	d.Name = ""
	assert.NonNil(t, d.IsValid())
}

func TestDomainIsValidBadName(t *testing.T) {
	d := getDomain()
	d.Name = "bad[name]"
	assert.NonNil(t, d.IsValid())
}

func TestDomainIsValidBadPort(t *testing.T) {
	d := getDomain()
	d.Port = 0
	assert.NonNil(t, d.IsValid())
}

func TestDomainIsValidNoOrg(t *testing.T) {
	d := getDomain()
	d.OrgKey = ""
	assert.NonNil(t, d.IsValid())
}

func TestDomainIsValidBadOrg(t *testing.T) {
	d := getDomain()
	d.OrgKey = "aoeu*snth"
	assert.NonNil(t, d.IsValid())
}

func TestDomainIsValidNoCorsConfig(t *testing.T) {
	d1 := getDomain()
	d1.CorsConfig = nil

	assert.Nil(t, d1.IsValid())
}

func TestDomainIsValidFailsOnCorsConfig(t *testing.T) {
	d1 := getDomain()
	d1.CorsConfig.AllowedOrigins = nil

	assert.NonNil(t, d1.IsValid())
}

func TestDomainIsValidFailedDkey(t *testing.T) {
	d1 := getDomain()
	d1.DomainKey = ""

	assert.NonNil(t, d1.IsValid())
}

func TestDomainIsValidFailedName(t *testing.T) {
	d1 := getDomain()
	d1.Name = ""

	assert.NonNil(t, d1.IsValid())
}

func TestDomainIsValidFailedPort(t *testing.T) {
	d1 := getDomain()
	d1.Port = 0

	assert.NonNil(t, d1.IsValid())
}

func TestDomainIsValidFailedDuplicateRedirect(t *testing.T) {
	d1 := getDomain()
	d1.Redirects = append(d1.Redirects, d1.Redirects[0])
	err := d1.IsValid()
	assert.NonNil(t, err)
	assert.HasSameElements(t, err.Errors, []ErrorCase{
		{
			"domain.redirects",
			fmt.Sprintf(
				"name must be unique, multiple redirects found called '%v'",
				d1.Redirects[0].Name,
			),
		},
	})
}

func TestDomainIsValidFailedRedirect(t *testing.T) {
	d1 := getDomain()
	d1.Redirects[0].To = ""
	err := d1.IsValid()
	assert.NonNil(t, err)
	assert.HasSameElements(t, err.Errors, []ErrorCase{
		{"domain.redirects[sample-redirect].to", "must not be empty"},
	})
}

func getThreeDomains() (Domain, Domain, Domain) {
	d1 := Domain{"dkey-1", "zk", "name", 10, nil, true, nil, DomainAliases{}, "okey", Checksum{}}
	d2 := Domain{"dkey-2", "zk", "name", 20, nil, true, nil, DomainAliases{}, "okey", Checksum{}}
	d3 := Domain{"dkey-3", "zk", "name", 30, nil, true, nil, DomainAliases{}, "okey", Checksum{}}

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
	d1 := Domain{"dkey-1", "zk", "name", 10, nil, true, nil, DomainAliases{}, "okey", Checksum{}}
	d2 := Domain{"dkey-2", "zk", "name", 20, nil, true, nil, DomainAliases{}, "okey", Checksum{}}
	d3 := Domain{"dkey-3", "zk", "name", 30, nil, true, nil, DomainAliases{}, "okey", Checksum{}}
	ds := Domains{d3, d2, d1}

	assert.Nil(t, ds.IsValid())
}

func TestDomainsIsValidFailureDupe(t *testing.T) {
	d1 := Domain{"dkey-1", "zk", "name", 10, nil, true, nil, DomainAliases{}, "okey", Checksum{}}
	d2 := Domain{"dkey-2", "zk", "name", 20, nil, true, nil, DomainAliases{}, "okey", Checksum{}}
	d3 := Domain{"dkey-3", "zk", "name", 30, nil, true, nil, DomainAliases{}, "okey", Checksum{}}
	ds := Domains{d3, d2, d1, d3}

	assert.NonNil(t, ds.IsValid())
}

func TestDomainsIsValidFailureBadDomain(t *testing.T) {
	d1 := Domain{"dkey-1", "zk", "name", 10, nil, true, nil, DomainAliases{}, "okey", Checksum{}}
	d2 := Domain{"dkey-2", "", "name", 20, nil, true, nil, DomainAliases{}, "okey", Checksum{}}
	d3 := Domain{"dkey-3", "zk", "name", 30, nil, true, nil, DomainAliases{}, "okey", Checksum{}}
	ds := Domains{d3, d2, d1}

	assert.NonNil(t, ds.IsValid())
}

func mkCC() *CorsConfig {
	return &CorsConfig{
		AllowedOrigins:   []string{"a", "b"},
		AllowCredentials: true,
		ExposedHeaders:   []string{"b", "c"},
		MaxAge:           500,
		AllowedMethods:   []string{"GET", "PUT"},
		AllowedHeaders:   []string{"h1", "h2"},
	}

}

func TestCorsConfigEqualsTrue(t *testing.T) {
	cc := *mkCC()
	cc2 := *mkCC()
	assert.True(t, cc.Equals(cc2))
	assert.True(t, cc2.Equals(cc))
}

func TestCorsConfigEqualsTrueOrderChanged(t *testing.T) {
	cc := *mkCC()
	cc2 := *mkCC()

	swap := func(s1 []string) {
		s1[0], s1[1] = s1[1], s1[0]
	}
	swap(cc2.AllowedOrigins)
	swap(cc2.ExposedHeaders)
	swap(cc2.AllowedMethods)
	swap(cc2.AllowedHeaders)

	assert.True(t, cc.Equals(cc2))
	assert.True(t, cc2.Equals(cc))
}

func TestCorsConfigEqualsFalseAge(t *testing.T) {
	cc := *mkCC()
	cc2 := *mkCC()
	cc2.MaxAge = cc.MaxAge + 1

	assert.False(t, cc.Equals(cc2))
	assert.False(t, cc2.Equals(cc))
}

func TestCorsConfigEqualsFalseAllowedOrigins(t *testing.T) {
	cc := *mkCC()
	cc2 := *mkCC()
	cc2.AllowedOrigins = append(cc2.AllowedOrigins, "new-element")

	assert.False(t, cc.Equals(cc2))
	assert.False(t, cc2.Equals(cc))
}

func TestCorsConfigEqualsFalseAllowCredentials(t *testing.T) {
	cc := *mkCC()
	cc2 := *mkCC()
	cc2.AllowCredentials = !cc2.AllowCredentials

	assert.False(t, cc.Equals(cc2))
	assert.False(t, cc2.Equals(cc))
}

func TestCorsConfigEqualsFalseExposedHeaders(t *testing.T) {
	cc := *mkCC()
	cc2 := *mkCC()
	cc2.ExposedHeaders = append(cc2.ExposedHeaders, "new-element")

	assert.False(t, cc.Equals(cc2))
	assert.False(t, cc2.Equals(cc))
}

func TestCorsConfigEqualsFalseAllowedMethods(t *testing.T) {
	cc := *mkCC()
	cc2 := *mkCC()
	cc2.AllowedMethods = append(cc2.AllowedMethods, "new-element")

	assert.False(t, cc.Equals(cc2))
	assert.False(t, cc2.Equals(cc))
}

func TestCorsConfigEqualsFalseAllowedHeaders(t *testing.T) {
	cc := *mkCC()
	cc2 := *mkCC()
	cc2.AllowedHeaders = append(cc2.AllowedHeaders, "new-element")

	assert.False(t, cc.Equals(cc2))
	assert.False(t, cc2.Equals(cc))
}

func TestDomainAliasEquals(t *testing.T) {
	da1 := DomainAlias("example.com")
	da2 := DomainAlias("example.com")
	assert.True(t, da1.Equals(da2))
	assert.True(t, da2.Equals(da1))
}

func TestDomainAliasEqualsFails(t *testing.T) {
	da1 := DomainAlias("www.google.com")
	da2 := DomainAlias("www.bing.com")
	assert.False(t, da1.Equals(da2))
	assert.False(t, da2.Equals(da1))
}

func TestDomainAliasIsValidSuccess(t *testing.T) {
	da := DomainAlias("example.com")
	assert.Nil(t, da.IsValid())
}

func TestDomainAliasIsValidFails(t *testing.T) {
	test := func(in string, fail bool) {
		got := DomainAlias(in).IsValid()
		var want *ValidationError

		if fail {
			want = &ValidationError{[]ErrorCase{{"", AliasPatternFailure}}}
		}

		if !assert.DeepEqual(t, got, want) {
			t.Logf("failed validation test on domain alias: %v\n--------------------", in)
		}
	}

	fail := func(in string) { test(in, true) }
	pass := func(in string) { test(in, false) }

	fail("*example.com")
	fail("example.com*")
	fail("*.test.*")
	fail("test.*.com")
	fail("test..com")
	fail(".example.com")
	fail("example.com.")
	fail(".example.com")
	fail("*.")
	fail(".*")
	fail("")
	pass("*.example.com")
	pass("test.*")
	pass("test.test.example.com")
}

func TestDomainAliasesIsValid(t *testing.T) {
	da := DomainAliases{
		"example.com",
		"*.example.com",
		"test.*",
		"bar.example.com",
	}

	assert.Nil(t, da.IsValid())
}

func TestDomainAliasesIsValidDupes(t *testing.T) {
	da := DomainAliases{
		"test.com",
		"*.test.com",
		"test.com",
	}

	assert.DeepEqual(t, da.IsValid(), &ValidationError{[]ErrorCase{
		{"domain_aliases", "duplicate alias found test.com"},
	}})
}

func TestDomainAliasesIsValidFailure(t *testing.T) {
	da := DomainAlias("*.*.*")
	daerr := da.IsValid()

	das := DomainAliases{da}
	assert.DeepEqual(t, das.IsValid(), &ValidationError{[]ErrorCase{
		{fmt.Sprintf("domain_aliases[%v]", da), daerr.Errors[0].Msg},
	}})
}

func TestDomainAliasesEquals(t *testing.T) {
	das1 := DomainAliases{"test.com", "example.com"}
	das2 := DomainAliases{"example.com", "test.com"}

	assert.True(t, das1.Equals(das2))
	assert.True(t, das2.Equals(das1))
}

func TestDomainAliasesEqualsFailure(t *testing.T) {
	das1 := DomainAliases{"test.com", "example.com", "foo.com"}
	das2 := DomainAliases{"example.com", "test.com"}

	assert.False(t, das1.Equals(das2))
	assert.False(t, das2.Equals(das1))
}
