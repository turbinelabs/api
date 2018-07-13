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
	"testing"

	"github.com/turbinelabs/test/assert"
)

func getListeners() (Listener, Listener) {
	l := Listener{
		ListenerKey: "lkey1",
		ZoneKey:     "zkey1",
		Name:        "name1",
		IP:          "127.0.0.1",
		Port:        80,
		Protocol:    "http",
		DomainKeys:  []DomainKey{"dkey1", "dkey2"},
		TracingConfig: &TracingConfig{
			Ingress:               true,
			RequestHeadersForTags: []string{"x-foo"},
		},
		OrgKey:   "okey1",
		Checksum: Checksum{"csum1"},
	}

	return l, l
}

func TestListenerEquals(t *testing.T) {
	l1, l2 := getListeners()

	assert.True(t, l1.Equals(l2))
	assert.True(t, l2.Equals(l1))
}

func TestListenerEqualsDiffListenerKey(t *testing.T) {
	l1, l2 := getListeners()
	l2.ListenerKey = "lkey2"
	assert.False(t, l1.Equals(l2))
	assert.False(t, l2.Equals(l1))
}

func TestListenerEqualsDiffZoneKey(t *testing.T) {
	l1, l2 := getListeners()
	l2.ZoneKey = "zkey2"
	assert.False(t, l1.Equals(l2))
	assert.False(t, l2.Equals(l1))
}

func TestListenerEqualsDiffName(t *testing.T) {
	l1, l2 := getListeners()
	l2.Name = "name2"
	assert.False(t, l1.Equals(l2))
	assert.False(t, l2.Equals(l1))
}

func TestListenerEqualsDiffIP(t *testing.T) {
	l1, l2 := getListeners()
	l2.IP = "127.0.0.2"
	assert.False(t, l1.Equals(l2))
	assert.False(t, l2.Equals(l1))
}

func TestListenerEqualsDiffPort(t *testing.T) {
	l1, l2 := getListeners()
	l2.Port = 81
	assert.False(t, l1.Equals(l2))
	assert.False(t, l2.Equals(l1))
}

func TestListenerEqualsDiffProtocol(t *testing.T) {
	l1, l2 := getListeners()
	l2.Protocol = "http2"
	assert.False(t, l1.Equals(l2))
	assert.False(t, l2.Equals(l1))
}

func TestListenerEqualsDiffTracingConfig(t *testing.T) {
	l1, l2 := getListeners()
	l2.TracingConfig = &TracingConfig{
		Ingress:               false,
		RequestHeadersForTags: []string{"x-bar"},
	}
	assert.False(t, l1.Equals(l2))
	assert.False(t, l2.Equals(l1))
}

func TestListenerEqualsDiffOrgKey(t *testing.T) {
	l1, l2 := getListeners()
	l2.OrgKey = "okey2"
	assert.False(t, l1.Equals(l2))
	assert.False(t, l2.Equals(l1))
}

func TestListenerEqualsDiffChecksum(t *testing.T) {
	l1, l2 := getListeners()
	l2.Checksum = Checksum{"csum2"}
	assert.False(t, l1.Equals(l2))
	assert.False(t, l2.Equals(l1))
}

func TestListenerEqualsDiffDomains(t *testing.T) {
	l1, l2 := getListeners()
	l2.DomainKeys = []DomainKey{"dkey1"}
	assert.False(t, l1.Equals(l2))
	assert.False(t, l2.Equals(l1))
}

func TestListenerEqualsDiffDomainOrder(t *testing.T) {
	p1, p2 := getProxies()
	p2.DomainKeys = []DomainKey{"dkey2", "dkey1"}
	assert.True(t, p1.Equals(p2))
	assert.True(t, p2.Equals(p1))
}

func mkTestL() Listener {
	return Listener{
		ListenerKey: "lkey1",
		ZoneKey:     "zkey1",
		Name:        "name1",
		IP:          "127.0.0.1",
		Port:        80,
		Protocol:    "http",
		DomainKeys:  []DomainKey{"dkey1", "dkey2"},
		TracingConfig: &TracingConfig{
			Ingress:               true,
			RequestHeadersForTags: []string{"x-foo"},
		},
		OrgKey:   "okey1",
		Checksum: Checksum{"csum1"},
	}
}

func TestListenerIsValid(t *testing.T) {
	l := mkTestL()
	assert.Nil(t, l.IsValid())
}

func TestListenerIsValidNoListenerKey(t *testing.T) {
	l := mkTestL()
	l.ListenerKey = ""
	assert.NonNil(t, l.IsValid())
}

func TestListenerIsValidBadListenerKey(t *testing.T) {
	l := mkTestL()
	l.ListenerKey = "aoeunthi-!!"
	assert.NonNil(t, l.IsValid())
}

func TestListenerIsValidNoName(t *testing.T) {
	l := mkTestL()
	l.Name = ""
	assert.NonNil(t, l.IsValid())
}

func TestListenerIsValidNoZoneKey(t *testing.T) {
	l := mkTestL()
	l.ZoneKey = ""
	assert.NonNil(t, l.IsValid())
}

func TestListenerIsValidBadZoneKey(t *testing.T) {
	l := mkTestL()
	l.ZoneKey = "aoeunthi-!!"
	assert.NonNil(t, l.IsValid())
}

func TestListenerIsValidBadName(t *testing.T) {
	l := mkTestL()
	l.Name = "!!! name["
	assert.NonNil(t, l.IsValid())
}

func TestListenerIsValidBadIP(t *testing.T) {
	l := mkTestL()
	l.IP = "!!! name["
	assert.NonNil(t, l.IsValid())
}

func TestListenerIsValidBadPort(t *testing.T) {
	l := mkTestL()
	l.Port = -1
	assert.NonNil(t, l.IsValid())
}

func TestListenerIsValidBadProtocol(t *testing.T) {
	l := mkTestL()
	l.Protocol = "foo"
	assert.NonNil(t, l.IsValid())
}

func TestListenerIsValidBadDomainKeys(t *testing.T) {
	l := mkTestL()
	badKey := "anethircgoenith]]"
	l.DomainKeys = []DomainKey{DomainKey(badKey)}
	gotErr := l.IsValid()
	assert.DeepEqual(t, gotErr, &ValidationError{[]ErrorCase{
		{
			fmt.Sprintf("listener.domain_keys[%v]", badKey),
			"must match pattern: ^[0-9a-zA-Z]+(-[0-9a-zA-Z]+)*$",
		},
	}})
}

func TestListenerIsValidDupeDomainKeys(t *testing.T) {
	l := mkTestL()
	l.DomainKeys = append(l.DomainKeys, l.DomainKeys[0])
	gotErr := l.IsValid()
	assert.DeepEqual(t, gotErr, &ValidationError{[]ErrorCase{
		{"listener.domain_keys", fmt.Sprintf("duplicate domain key '%v'", l.DomainKeys[0])},
	}})
}

func TestListenerIsValidBadOrgKey(t *testing.T) {
	l := mkTestL()
	l.OrgKey = "---"
	assert.NonNil(t, l.IsValid())
}
