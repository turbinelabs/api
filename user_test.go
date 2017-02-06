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

	assert.Nil(t, u.IsValid(true))
	assert.Nil(t, u.IsValid(false))
}

func TestUserIsValidNoDeletedAt(t *testing.T) {
	u := getUser()
	u.DeletedAt = nil

	assert.Nil(t, u.IsValid(true))
	assert.Nil(t, u.IsValid(false))
}

func TestUserIsValidNoUserKey(t *testing.T) {
	u := getUser()
	u.UserKey = ""

	assert.Nil(t, u.IsValid(true))
	assert.NonNil(t, u.IsValid(false))
}

func TestUserIsValidNoEmail(t *testing.T) {
	u := getUser()
	u.LoginEmail = ""

	assert.NonNil(t, u.IsValid(true))
	assert.NonNil(t, u.IsValid(false))
}

func TestUserIsValidNoAuthKey(t *testing.T) {
	u := getUser()
	u.APIAuthKey = ""

	assert.NonNil(t, u.IsValid(true))
	assert.NonNil(t, u.IsValid(false))
}

func TestUserIsValidNoOrgKey(t *testing.T) {
	u := getUser()
	u.OrgKey = ""

	assert.NonNil(t, u.IsValid(true))
	assert.NonNil(t, u.IsValid(false))
}
