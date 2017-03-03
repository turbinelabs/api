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

	assert.Nil(t, org.IsValid())
}

func TestOrgIsValidBadOrgKey(t *testing.T) {
	org := getOrg()
	org.OrgKey = "aoeu-%-1234"
	assert.NonNil(t, org.IsValid())
}

func TestOrgIsValidNoOrgKey(t *testing.T) {
	org := getOrg()
	org.OrgKey = ""

	assert.NonNil(t, org.IsValid())
}

func TestOrgIsValidBadName(t *testing.T) {
	org := getOrg()
	org.Name = "bad [name]"
	assert.NonNil(t, org.IsValid())
}

func TestOrgIsValidNoName(t *testing.T) {
	org := getOrg()
	org.Name = ""

	assert.NonNil(t, org.IsValid())
}

func TestOrgIsValidNoEmail(t *testing.T) {
	org := getOrg()
	org.ContactEmail = ""

	assert.NonNil(t, org.IsValid())
}
