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
	"time"

	"github.com/turbinelabs/test/assert"
)

func getUsers() (User, User) {
	now := time.Now()
	user1 := User{
		UserKey:    "ukey1",
		LoginEmail: "email1",
		APIAuthKey: "akey1",
		OrgKey:     "okey1",
		DeletedAt:  &now,
		Checksum:   Checksum{"csum1"},
	}

	user2 := user1
	now2 := now
	user2.DeletedAt = &now2

	return user1, user2
}

func TestUserEqualsWithNoDeletion(t *testing.T) {
	u1, u2 := getUsers()
	u1.DeletedAt = nil
	u2.DeletedAt = nil

	assert.True(t, u1.Equals(u2))
	assert.True(t, u2.Equals(u1))
}

func TestUserEqualsWithDeletion(t *testing.T) {
	u1, u2 := getUsers()

	assert.True(t, u1.Equals(u2))
	assert.True(t, u2.Equals(u1))
}

func TestEqualsDiffUserKey(t *testing.T) {
	u1, u2 := getUsers()
	u2.UserKey = "ukey2"

	assert.False(t, u1.Equals(u2))
	assert.False(t, u2.Equals(u1))
}

func TestUserEqualsDiffEmail(t *testing.T) {
	u1, u2 := getUsers()
	u2.LoginEmail = "email2"

	assert.False(t, u1.Equals(u2))
	assert.False(t, u2.Equals(u1))
}

func TestUserEqualsDiffAuthKey(t *testing.T) {
	u1, u2 := getUsers()
	u2.APIAuthKey = "akey2"

	assert.False(t, u1.Equals(u2))
	assert.False(t, u2.Equals(u1))
}

func TestUserEqualsDiffOrgKey(t *testing.T) {
	u1, u2 := getUsers()
	u2.OrgKey = "okey2"

	assert.False(t, u1.Equals(u2))
	assert.False(t, u2.Equals(u1))
}

func TestUserEqualsDiffDeletedAt(t *testing.T) {
	u1, u2 := getUsers()
	ts := time.Now()
	ts = ts.Add(10 * time.Second)
	u2.DeletedAt = &ts

	assert.False(t, u1.Equals(u2))
	assert.False(t, u2.Equals(u1))
}

func TestUserEqualsDiffChecksum(t *testing.T) {
	u1, u2 := getUsers()
	u2.Checksum = Checksum{"csum2"}

	assert.False(t, u1.Equals(u2))
	assert.False(t, u2.Equals(u1))
}

func getUser() User {
	now := time.Now()
	return User{
		UserKey:    "ukey1",
		LoginEmail: "email1",
		APIAuthKey: "akey1",
		OrgKey:     "okey1",
		DeletedAt:  &now,
	}
}

func TestUserIsValid(t *testing.T) {
	u := getUser()

	assert.Nil(t, u.IsValid())
}

func TestUserIsValidNoDeletedAt(t *testing.T) {
	u := getUser()
	u.DeletedAt = nil

	assert.Nil(t, u.IsValid())
}

func TestUserIsValidNoUserKey(t *testing.T) {
	u := getUser()
	u.UserKey = ""

	assert.NonNil(t, u.IsValid())
}

func TestUserIsValidBadUserKey(t *testing.T) {
	u := getUser()
	u.UserKey = "bad-user-@"
	assert.NonNil(t, u.IsValid())
}

func TestUserIsValidNoEmail(t *testing.T) {
	u := getUser()
	u.LoginEmail = ""

	assert.NonNil(t, u.IsValid())
}

func TestUserIsValidNoAuthKey(t *testing.T) {
	u := getUser()
	u.APIAuthKey = ""

	assert.Nil(t, u.IsValid())
}

func TestUserIsValidDelegatedAuthKey(t *testing.T) {
	u := getUser()
	u.APIAuthKey = "Bearer " +
		"https://login.turbinelabs.io/auth/realms/turbine-labs " +
		"12341234-1234-1234-1234-123412341234"
	assert.Nil(t, u.IsValid())
}

func TestUserIsValidBadAuthKey(t *testing.T) {
	u := getUser()
	u.APIAuthKey = "Bearer [turbine-labs] snth-snth-snth"
	assert.NonNil(t, u.IsValid())
}

func TestUserIsValidNoOrgKey(t *testing.T) {
	u := getUser()
	u.OrgKey = ""

	assert.NonNil(t, u.IsValid())
}

func TestUserIsValidBadOrgKey(t *testing.T) {
	u := getUser()
	u.OrgKey = "bad-(-key"
	assert.NonNil(t, u.IsValid())
}
