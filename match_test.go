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

	assert.Nil(t, m.IsValid(true))
	assert.Nil(t, m.IsValid(false))
}

func TestMatchIsValidFailedFrom(t *testing.T) {
	m := Match{CookieMatchKind, Metadatum{}, Metadatum{"snth", "1234"}}

	assert.NonNil(t, m.IsValid(true))
	assert.NonNil(t, m.IsValid(false))
}

func TestMatchIsValidSuccessEmptyTo(t *testing.T) {
	m := Match{CookieMatchKind, Metadatum{"aoeu", ""}, Metadatum{"", ""}}

	assert.Nil(t, m.IsValid(true))
	assert.Nil(t, m.IsValid(false))
}

func TestMatchIsValidFailedKind(t *testing.T) {
	m := Match{"snth", Metadatum{"aoeu", ""}, Metadatum{"snth", "1234"}}

	assert.NonNil(t, m.IsValid(true))
	assert.NonNil(t, m.IsValid(false))
}

func TestMatchIsValidFailedToNoKey(t *testing.T) {
	m := Match{"aoeu", Metadatum{"snth", ""}, Metadatum{"", "snth"}}
	assert.NonNil(t, m.IsValid(true))
	assert.NonNil(t, m.IsValid(false))
}

func TestMatchesIsValidSuccess(t *testing.T) {
	m1 := Match{CookieMatchKind, Metadatum{"aoeu", ""}, Metadatum{"snth", "1234"}}
	m2 := Match{HeaderMatchKind, Metadatum{"aoeu", ""}, Metadatum{"snth", "1234"}}
	m3 := Match{QueryMatchKind, Metadatum{"aoeu", ""}, Metadatum{"snth", "1234"}}
	m := Matches{m1, m2, m3}

	assert.Nil(t, m.IsValid(true))
	assert.Nil(t, m.IsValid(false))
}

func TestMatchesIsValidEmpty(t *testing.T) {
	m := Matches{}

	assert.Nil(t, m.IsValid(true))
	assert.Nil(t, m.IsValid(false))
}

func TestMatchesIsValidFailure(t *testing.T) {
	m1 := Match{CookieMatchKind, Metadatum{"aoeu", ""}, Metadatum{"snth", "1234"}}
	m2 := Match{"badmatch", Metadatum{"aoeu", ""}, Metadatum{"snth", "1234"}}
	m3 := Match{QueryMatchKind, Metadatum{"aoeu", ""}, Metadatum{"snth", "1234"}}
	m := Matches{m1, m2, m3}

	assert.NonNil(t, m.IsValid(true))
	assert.NonNil(t, m.IsValid(false))
}
