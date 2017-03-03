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

func TestMatchEqualsSuccess(t *testing.T) {
	m1 := Match{HeaderMatchKind, Metadatum{"x-random", "value"}, Metadatum{"randomflag", "true"}}
	m2 := Match{HeaderMatchKind, Metadatum{"x-random", "value"}, Metadatum{"randomflag", "true"}}

	assert.True(t, m1.Equals(m2))
	assert.True(t, m2.Equals(m1))
}

func TestMatchEqualsFromVaries(t *testing.T) {
	m1 := Match{HeaderMatchKind, Metadatum{"x-random", "value"}, Metadatum{"randomflag", "true"}}
	m2 := Match{HeaderMatchKind, Metadatum{"x-other", "value"}, Metadatum{"randomflag", "true"}}

	assert.False(t, m1.Equals(m2))
	assert.False(t, m2.Equals(m1))
}

func TestMatchEqualsToVaries(t *testing.T) {
	m1 := Match{HeaderMatchKind, Metadatum{"x-random", "value"}, Metadatum{"randomflag", "true"}}
	m2 := Match{HeaderMatchKind, Metadatum{"x-random", "value"}, Metadatum{"specificflag", "true"}}

	assert.False(t, m1.Equals(m2))
	assert.False(t, m2.Equals(m1))
}

func TestMatchEqualsMatchKindVaries(t *testing.T) {
	m1 := Match{HeaderMatchKind, Metadatum{"x-random", "value"}, Metadatum{"randomflag", "true"}}
	m2 := Match{CookieMatchKind, Metadatum{"x-random", "value"}, Metadatum{"randomflag", "true"}}

	assert.False(t, m1.Equals(m2))
	assert.False(t, m2.Equals(m1))
}

func TestMatchesEqualsSuccess(t *testing.T) {
	m1 := Match{HeaderMatchKind, Metadatum{"x-random", "value"}, Metadatum{"randomflag", "true"}}
	m2 := Match{CookieMatchKind, Metadatum{"x-random", "value"}, Metadatum{"randomflag", "true"}}

	slice1 := Matches{m1, m2}
	slice2 := Matches{m2, m1}

	assert.True(t, slice1.Equals(slice2))
	assert.True(t, slice2.Equals(slice1))
}

func TestMatchesEqualsLengthMismatch(t *testing.T) {
	m1 := Match{HeaderMatchKind, Metadatum{"x-random", "value"}, Metadatum{"randomflag", "true"}}
	m2 := Match{CookieMatchKind, Metadatum{"x-random", "value"}, Metadatum{"randomflag", "true"}}

	slice1 := Matches{m1, m2}
	slice2 := Matches{m2}

	assert.False(t, slice1.Equals(slice2))
	assert.False(t, slice2.Equals(slice1))
}

func TestMatchesEqualsContentDiffers(t *testing.T) {
	m1 := Match{HeaderMatchKind, Metadatum{"x-random", "value"}, Metadatum{"randomflag", "true"}}
	m2 := Match{CookieMatchKind, Metadatum{"x-random", "value"}, Metadatum{"randomflag", "true"}}
	m3 := Match{CookieMatchKind, Metadatum{"x-random", "value"}, Metadatum{"randomflag", "false"}}

	slice1 := Matches{m1, m2}
	slice2 := Matches{m2, m3}

	assert.False(t, slice1.Equals(slice2))
	assert.False(t, slice2.Equals(slice1))
}

func TestMatchIsValidSucces(t *testing.T) {
	m := Match{CookieMatchKind, Metadatum{"aoeu", ""}, Metadatum{"snth", "1234"}}

	assert.Nil(t, m.IsValid())
}

func TestMatchIsValidFailedFromBadKey(t *testing.T) {
	m := Match{CookieMatchKind, Metadatum{"from]", ""}, Metadatum{"snth", "1234"}}

	assert.NonNil(t, m.IsValid())
}

func TestMatchIsValidFailedFromEmpty(t *testing.T) {
	m := Match{CookieMatchKind, Metadatum{}, Metadatum{"snth", "1234"}}

	assert.NonNil(t, m.IsValid())
}

func TestMatchIsValidSuccessEmptyTo(t *testing.T) {
	m := Match{CookieMatchKind, Metadatum{"aoeu", ""}, Metadatum{"", ""}}

	assert.Nil(t, m.IsValid())
}

func TestMatchIsValidFailedKind(t *testing.T) {
	m := Match{"snth", Metadatum{"aoeu", ""}, Metadatum{"snth", "1234"}}

	assert.NonNil(t, m.IsValid())
}

func TestMatchIsValidFailedToNoKey(t *testing.T) {
	m := Match{"aoeu", Metadatum{"snth", ""}, Metadatum{"", "snth"}}
	assert.NonNil(t, m.IsValid())
}

func TestMatchesIsValidSuccess(t *testing.T) {
	m1 := Match{CookieMatchKind, Metadatum{"aoeu", ""}, Metadatum{"snth", "1234"}}
	m2 := Match{HeaderMatchKind, Metadatum{"aoeu", ""}, Metadatum{"snth", "1234"}}
	m3 := Match{QueryMatchKind, Metadatum{"aoeu", ""}, Metadatum{"snth", "1234"}}
	m := Matches{m1, m2, m3}

	assert.Nil(t, m.IsValid())
}

func TestMatchesIsValidEmpty(t *testing.T) {
	m := Matches{}

	assert.Nil(t, m.IsValid())
}

func TestMatchesIsValidNil(t *testing.T) {
	var m Matches
	assert.Nil(t, m.IsValid())
}

func TestMatchesIsValidFailure(t *testing.T) {
	m1 := Match{CookieMatchKind, Metadatum{"aoeu", ""}, Metadatum{"snth", "1234"}}
	m2 := Match{"badmatch", Metadatum{"aoeu", ""}, Metadatum{"snth", "1234"}}
	m3 := Match{QueryMatchKind, Metadatum{"aoeu", ""}, Metadatum{"snth", "1234"}}
	m := Matches{m1, m2, m3}

	assert.NonNil(t, m.IsValid())
}

func TestMatchesIsValidDupeMatch(t *testing.T) {
	m1 := Match{CookieMatchKind, Metadatum{"type", "chocolate chip"}, Metadatum{"texture", "dense"}}
	m2 := Match{CookieMatchKind, Metadatum{"type", "amaretti"}, Metadatum{"texture", "chewy"}}

	m := Matches{m1, m2}
	assert.DeepEqual(t, m.IsValid(), &ValidationError{[]ErrorCase{
		{"", "duplicate match found cookie:type"},
	}})
}
