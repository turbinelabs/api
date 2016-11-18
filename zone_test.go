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
	assert.Nil(t, z.IsValid(true))
	assert.Nil(t, z.IsValid(false))
}

func TestZoneIsValidNoZoneKey(t *testing.T) {
	z := getZone()
	z.ZoneKey = ""
	assert.Nil(t, z.IsValid(true))
	assert.NonNil(t, z.IsValid(false))
}

func TestZoneIsValidNoName(t *testing.T) {
	z := getZone()
	z.Name = ""

	assert.NonNil(t, z.IsValid(true))
	assert.NonNil(t, z.IsValid(false))
}

func TestZoneIsValidNoOrg(t *testing.T) {
	z := getZone()
	z.OrgKey = ""

	assert.NonNil(t, z.IsValid(true))
	assert.NonNil(t, z.IsValid(false))
}