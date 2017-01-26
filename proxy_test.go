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

func TestIsValid(t *testing.T) {
	p := Proxy{ProxyKey: "pkey1", Name: "name1", ZoneKey: "zkey1"}
	assert.Nil(t, p.IsValid(true))
	assert.Nil(t, p.IsValid(false))
}

func TestIsValidNoProxyKey(t *testing.T) {
	p := Proxy{Name: "name1", ZoneKey: "zkey1"}
	assert.Nil(t, p.IsValid(true))
	assert.NonNil(t, p.IsValid(false))
}

func TestIsValidNoName(t *testing.T) {
	p := Proxy{ProxyKey: "pkey1", ZoneKey: "zkey1"}
	assert.NonNil(t, p.IsValid(true))
	assert.NonNil(t, p.IsValid(false))
}

func TestIsValidNoZoneKey(t *testing.T) {
	p := Proxy{ProxyKey: "pkey1", Name: "name1"}
	assert.NonNil(t, p.IsValid(true))
	assert.NonNil(t, p.IsValid(false))
}
