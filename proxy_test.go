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
	"testing"

	"github.com/turbinelabs/test/assert"
)

func getProxies() (Proxy, Proxy) {
	p := Proxy{
		ProxyKey:   "pkey1",
		Name:       "name1",
		ZoneKey:    "zkey1",
		OrgKey:     "okey1",
		Checksum:   Checksum{"csum1"},
		DomainKeys: []DomainKey{"dkey1", "dkey2"},
	}

	return p, p
}
func TestProxyEquals(t *testing.T) {
	p1, p2 := getProxies()

	assert.True(t, p1.Equals(p2))
	assert.True(t, p2.Equals(p1))
}

func TestProxyEqualsDiffProxyKey(t *testing.T) {
	p1, p2 := getProxies()
	p2.ProxyKey = "pkey2"
	assert.False(t, p1.Equals(p2))
	assert.False(t, p2.Equals(p1))
}

func TestProxyEqualsDiffZoneKey(t *testing.T) {
	p1, p2 := getProxies()
	p2.ZoneKey = "zkey2"
	assert.False(t, p1.Equals(p2))
	assert.False(t, p2.Equals(p1))
}

func TestProxyEqualsDiffName(t *testing.T) {
	p1, p2 := getProxies()
	p2.Name = "name2"
	assert.False(t, p1.Equals(p2))
	assert.False(t, p2.Equals(p1))
}

func TestProxyEquasDiffOrgKey(t *testing.T) {
	p1, p2 := getProxies()
	p2.OrgKey = "okey2"
	assert.False(t, p1.Equals(p2))
	assert.False(t, p2.Equals(p1))
}

func TestProxyEquasDiffChecksum(t *testing.T) {
	p1, p2 := getProxies()
	p2.Checksum = Checksum{"csum2"}
	assert.False(t, p1.Equals(p2))
	assert.False(t, p2.Equals(p1))
}

func TestProxyEqualsDiffDomains(t *testing.T) {
	p1, p2 := getProxies()
	p2.DomainKeys = []DomainKey{"dkey1"}
	assert.False(t, p1.Equals(p2))
	assert.False(t, p2.Equals(p1))
}

func TestProxyEqualsDiffDomainOrder(t *testing.T) {
	p1, p2 := getProxies()
	p2.DomainKeys = []DomainKey{"dkey2", "dkey1"}
	assert.True(t, p1.Equals(p2))
	assert.True(t, p2.Equals(p1))
}

func mkTestP() Proxy {
	return Proxy{
		ProxyKey:   "pk-1",
		ZoneKey:    "zk-1",
		Name:       "my neat proxy!",
		DomainKeys: []DomainKey{"dk-1", "dk-2"},
		OrgKey:     "ok-1",
	}
}

func TestProxyIsValid(t *testing.T) {
	p := mkTestP()
	assert.Nil(t, p.IsValid())
}

func TestProxyIsValidNoProxyKey(t *testing.T) {
	p := mkTestP()
	p.ProxyKey = ""
	assert.NonNil(t, p.IsValid())
}

func TestProxyIsValidNoName(t *testing.T) {
	p := mkTestP()
	p.Name = ""
	assert.NonNil(t, p.IsValid())
}

func TestProxyIsValidNoZoneKey(t *testing.T) {
	p := mkTestP()
	p.ZoneKey = ""
	assert.NonNil(t, p.IsValid())
}

func TestProxyIsValidBadKey(t *testing.T) {
	p := mkTestP()
	p.ProxyKey = "aosnetuh-!!!"
	assert.NonNil(t, p.IsValid())
}
func TestProxyIsValidBadName(t *testing.T) {
	p := mkTestP()
	p.Name = "some weird name["
	assert.NonNil(t, p.IsValid())
}

func TestProxyIsValidBadZoneKey(t *testing.T) {
	p := mkTestP()
	p.ZoneKey = "111-222-##"
	assert.NonNil(t, p.IsValid())
}

func TestProxyIsValidDupeDomainKeys(t *testing.T) {
	p := mkTestP()
	p.DomainKeys = append(p.DomainKeys, p.DomainKeys[0])
	assert.NonNil(t, p.IsValid())
}

func TestProxyIsValidBadOrgKey(t *testing.T) {
	p := mkTestP()
	p.OrgKey = "---"
	assert.NonNil(t, p.IsValid())
}
