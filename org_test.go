package api

import (
	"testing"

	"github.com/turbinelabs/test/assert"
)

func getOrgs() (Org, Org) {
	org := Org{Name: "name1", ContactEmail: "bar", OrgKey: "okey1", Checksum: Checksum{"csum1"}}
	return org, org
}

func TestOrgEquals(t *testing.T) {
	org1, org2 := getOrgs()

	assert.True(t, org1.Equals(org2))
	assert.True(t, org2.Equals(org1))
}

func TestOrgEqualsDiffName(t *testing.T) {
	org1, org2 := getOrgs()
	org2.Name = "name2"

	assert.False(t, org1.Equals(org2))
	assert.False(t, org2.Equals(org1))
}

func TestOrgEqualsDiffEmail(t *testing.T) {
	org1, org2 := getOrgs()
	org2.ContactEmail = "email2"

	assert.False(t, org1.Equals(org2))
	assert.False(t, org2.Equals(org1))
}

func TestOrgEqualsDiffOrg(t *testing.T) {
	org1, org2 := getOrgs()
	org2.OrgKey = "okey2"

	assert.False(t, org1.Equals(org2))
	assert.False(t, org2.Equals(org1))
}

func TestOrgEqualsDiffChecksum(t *testing.T) {
	org1, org2 := getOrgs()
	org2.Checksum = Checksum{"csum2"}

	assert.False(t, org1.Equals(org2))
	assert.False(t, org2.Equals(org1))
}

func getOrg() Org {
	return Org{Name: "name1", ContactEmail: "email1", OrgKey: "okey1"}
}

func TestOrgIsValid(t *testing.T) {
	org := getOrg()

	assert.Nil(t, org.IsValid(true))
	assert.Nil(t, org.IsValid(false))
}

func TestOrgIsValidNoOrgKey(t *testing.T) {
	org := getOrg()
	org.OrgKey = ""

	assert.Nil(t, org.IsValid(true))
	assert.NonNil(t, org.IsValid(false))
}

func TestOrgIsValidNoName(t *testing.T) {
	org := getOrg()
	org.Name = ""

	assert.NonNil(t, org.IsValid(true))
	assert.NonNil(t, org.IsValid(false))
}

func TestOrgIsValidNoEmail(t *testing.T) {
	org := getOrg()
	org.ContactEmail = ""

	assert.NonNil(t, org.IsValid(true))
	assert.NonNil(t, org.IsValid(false))
}
