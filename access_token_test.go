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
	"strings"
	"testing"
	"time"

	"github.com/turbinelabs/test/assert"
)

func getat() (AccessToken, AccessToken) {
	n := time.Now()
	mkat := func() AccessToken {
		return AccessToken{
			"access-token-key",
			"description",
			"signed-token",
			"user-key",
			"org-key",
			&n,
			Checksum{"checksum"},
		}
	}
	return mkat(), mkat()
}

func TestAccessTokenEquals(t *testing.T) {
	t1, t2 := getat()

	assert.True(t, t1.Equals(t2))
	assert.True(t, t2.Equals(t1))
}

func TestAccessTokenEqualsKeyChange(t *testing.T) {
	t1, t2 := getat()
	t2.AccessTokenKey = t2.AccessTokenKey + "aoeu"

	assert.False(t, t1.Equals(t2))
	assert.False(t, t2.Equals(t1))
}

func TestAccessTokenEqualsDescriptionChange(t *testing.T) {
	t1, t2 := getat()
	t2.Description = t2.Description + "aosentuh"

	assert.False(t, t1.Equals(t2))
	assert.False(t, t2.Equals(t1))
}

func TestAccessTokenEqualsSignedTokenChange(t *testing.T) {
	t1, t2 := getat()
	t2.SignedToken = t2.SignedToken + "aoeu"
	assert.False(t, t1.Equals(t2))
	assert.False(t, t2.Equals(t1))
}

func TestAccessTokenEqualsUserKey(t *testing.T) {
	t1, t2 := getat()
	t2.UserKey = "asonetuh"

	assert.False(t, t1.Equals(t2))
	assert.False(t, t2.Equals(t1))
}

func TestAccessTokenEqualsCreatedAtSomeNil(t *testing.T) {
	t1, t2 := getat()
	t2.CreatedAt = nil

	assert.False(t, t1.Equals(t2))
	assert.False(t, t2.Equals(t1))
}

func TestAccessTokenEqualsCreatedAtZeroNil(t *testing.T) {
	t1, t2 := getat()
	t1.CreatedAt = &time.Time{}
	t2.CreatedAt = nil

	assert.False(t, t1.Equals(t2))
	assert.False(t, t2.Equals(t1))
}

func TestAccessTokenEqualsCreatedAtNilNil(t *testing.T) {
	t1, t2 := getat()
	t1.CreatedAt = nil
	t2.CreatedAt = nil

	assert.True(t, t1.Equals(t2))
	assert.True(t, t2.Equals(t1))
}

func TestAccessTokenEqualsChecksumChanges(t *testing.T) {
	t1, t2 := getat()
	t2.Checksum = Checksum{"asoneteaoseouah"}

	assert.False(t, t1.Equals(t2))
	assert.False(t, t2.Equals(t1))
}

func TestAccessTokenIsValid(t *testing.T) {
	t1, _ := getat()

	assert.Nil(t, t1.IsValid())
}

func TestAccessTokenIsValidNoKey(t *testing.T) {
	t1, _ := getat()
	t1.AccessTokenKey = ""
	assert.DeepEqual(t, t1.IsValid(), &ValidationError{[]ErrorCase{
		{"access_token.access_token_key", "may not be empty"},
	}})

	t1.AccessTokenKey = "aonethu anotehu"
	assert.DeepEqual(t, t1.IsValid(), &ValidationError{[]ErrorCase{
		{"access_token.access_token_key", KeyPatternMatchFailure},
	}})
}

func TestAccessTokenIsValidNoUser(t *testing.T) {
	t1, _ := getat()
	t1.UserKey = ""
	assert.DeepEqual(t, t1.IsValid(), &ValidationError{[]ErrorCase{
		{"access_token.user_key", "may not be empty"},
	}})

	t1.UserKey = "aonethu anotehu"
	assert.DeepEqual(t, t1.IsValid(), &ValidationError{[]ErrorCase{
		{"access_token.user_key", KeyPatternMatchFailure},
	}})
}

func TestAccessTokenIsValidNoOrg(t *testing.T) {
	t1, _ := getat()
	t1.OrgKey = ""
	assert.DeepEqual(t, t1.IsValid(), &ValidationError{[]ErrorCase{
		{"access_token.org_key", "may not be empty"},
	}})

	t1.OrgKey = "aonethu anotehu"
	assert.DeepEqual(t, t1.IsValid(), &ValidationError{[]ErrorCase{
		{"access_token.org_key", KeyPatternMatchFailure},
	}})
}

func TestAccessTokenIsValidDescriptionBad(t *testing.T) {
	t1, _ := getat()
	t1.Description = ""
	assert.DeepEqual(t, t1.IsValid(), &ValidationError{[]ErrorCase{
		{"access_token.description", tooShortErr},
	}})

	t1.Description = strings.Repeat("a", DescriptionLen+1)
	assert.DeepEqual(t, t1.IsValid(), &ValidationError{[]ErrorCase{
		{"access_token.description", tooLongErr},
	}})
}
