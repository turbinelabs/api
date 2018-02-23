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

func getZones() (Zone, Zone) {
	zone := Zone{ZoneKey: "zkey1", Name: "name1", OrgKey: "okey1", Checksum: Checksum{"csum1"}}
	return zone, zone
}

func TestZoneEquals(t *testing.T) {
	z1, z2 := getZones()

	assert.True(t, z1.Equals(z2))
	assert.True(t, z2.Equals(z1))
}

func TestZoneEqualsDiffZoneKey(t *testing.T) {
	z1, z2 := getZones()
	z2.ZoneKey = "zkey2"

	assert.False(t, z1.Equals(z2))
	assert.False(t, z2.Equals(z1))
}

func TestZoneEqualsDiffName(t *testing.T) {
	z1, z2 := getZones()
	z2.Name = "name2"

	assert.False(t, z1.Equals(z2))
	assert.False(t, z2.Equals(z1))
}

func TestZoneEqualsDiffOrgKey(t *testing.T) {
	z1, z2 := getZones()
	z2.OrgKey = "okey2"

	assert.False(t, z1.Equals(z2))
	assert.False(t, z2.Equals(z1))
}

func TestZoneEqualsDiffChecksum(t *testing.T) {
	z1, z2 := getZones()
	z2.Checksum = Checksum{"csum2"}

	assert.False(t, z1.Equals(z2))
	assert.False(t, z2.Equals(z1))
}

func getZone() Zone {
	return Zone{ZoneKey: "zkey1", Name: "name1", OrgKey: "okey1"}
}

func TestZoneIsValid(t *testing.T) {
	z := getZone()
	assert.Nil(t, z.IsValid())
}

func TestZoneIsValidNoZoneKey(t *testing.T) {
	z := getZone()
	z.ZoneKey = ""
	assert.NonNil(t, z.IsValid())
}

func TestZoneIsValidBadZoneKey(t *testing.T) {
	z := getZone()
	z.ZoneKey = "-!"
	assert.NonNil(t, z.IsValid())
}

func TestZoneIsValidNoName(t *testing.T) {
	z := getZone()
	z.Name = ""

	assert.NonNil(t, z.IsValid())
}

func TestZoneIsValidBadName(t *testing.T) {
	z := getZone()
	z.Name = "[]"

	assert.NonNil(t, z.IsValid())
}

func TestZoneIsValidBadOrg(t *testing.T) {
	z := getZone()
	z.OrgKey = "org#"

	assert.NonNil(t, z.IsValid())
}

func TestZoneIsValidNoOrg(t *testing.T) {
	z := getZone()
	z.OrgKey = ""

	assert.NonNil(t, z.IsValid())
}
