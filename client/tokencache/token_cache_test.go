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

package tokencache

import (
	"encoding/json"
	"io/ioutil"
	"math"
	"os"
	"path/filepath"
	"testing"
	"time"

	"golang.org/x/oauth2"

	tbnflag "github.com/turbinelabs/nonstdlib/flag"
	"github.com/turbinelabs/nonstdlib/ptr"
	tbntime "github.com/turbinelabs/nonstdlib/time"
	"github.com/turbinelabs/test/assert"
)

func home(p string) string {
	h := os.Getenv("HOME")
	if p == "" {
		return h
	}
	return filepath.Join(h, p)
}

func TestNewPathFromFlags(t *testing.T) {
	tfs := tbnflag.NewTestFlagSet()
	pff := NewPathFromFlags("pfix", tfs)

	tfs.Parse([]string{"-token-cache", "whee"})
	assert.Equal(t, pff.Default(), home(".pfix-auth-cache"))
	assert.Equal(t, pff.CachePath(), "whee")
}

func TestNewPathFromFlagsUnset(t *testing.T) {
	tfs := tbnflag.NewTestFlagSet()
	pff := NewPathFromFlags("", tfs)

	assert.Equal(t, pff.Default(), home(".tbn-auth-cache"))
	assert.Equal(t, pff.CachePath(), pff.Default())
}

func TestNewStaticPath(t *testing.T) {
	p := "/some/path"
	sp := NewStaticPath(p)
	assert.Equal(t, sp.CachePath(), p)
	assert.Equal(t, sp.Default(), home(".tbn-auth-cache"))
}

func getToken() (*oauth2.Token, TokenCache, OAuth2Token) {
	now := time.Now().UTC()
	then := now.Add(30 * time.Minute)

	tkn := oauth2.Token{
		AccessToken:  "at",
		TokenType:    "bearer",
		RefreshToken: "rt",
		Expiry:       then,
	}

	cache := TokenCache{
		ClientID:         "id",
		ClientKey:        "key",
		ProviderURL:      "https://login.turbinelabs.io/auth/realms/turbine-labs",
		Username:         "user",
		ExpiresAt:        &now,
		RefreshExpiresAt: &then,
		Token:            &tkn,
	}

	return &tkn, cache, WrapOAuth2Token(&tkn)
}

func TestNewFromFile(t *testing.T) {
	file, err := ioutil.TempFile("", "token")
	defer func() { os.Remove(file.Name()) }()
	assert.Nil(t, err)

	_, want, _ := getToken()

	b, err := json.MarshalIndent(want, "", "  ")
	assert.Nil(t, err)
	assert.Nil(t, ioutil.WriteFile(file.Name(), b, 0600))

	tc, err := NewFromFile(file.Name())
	assert.Nil(t, err)

	assert.NonNil(t, tc.timeSource)

	// *handwave*
	want.timeSource = tc.timeSource
	assert.DeepEqual(t, tc, want)
}

func TestSave(t *testing.T) {
	file, err := ioutil.TempFile("", "token")
	defer func() { os.Remove(file.Name()) }()
	assert.Nil(t, err)

	_, want, _ := getToken()

	assert.Nil(t, want.Save(file.Name()))

	got, err := NewFromFile(file.Name())
	assert.Nil(t, err)

	// *handwave*
	got.timeSource = nil
	assert.DeepEqual(t, got, want)
}

func TestWrapOAuthToken(t *testing.T) {
	tkn, _, _ := getToken()
	ti := WrapOAuth2Token(tkn)
	assert.DeepEqual(t, ti, otw{tkn})
}

func doNilTest(t *testing.T, obj OAuth2Token) {
	_, tc, _ := getToken()

	want := tc

	tc.SetToken(obj)
	assert.Nil(t, tc.Token)
	assert.Nil(t, tc.ExpiresAt)
	assert.Nil(t, tc.RefreshExpiresAt)
	assert.Equal(t, tc.ClientID, want.ClientID)
	assert.Equal(t, tc.ClientKey, want.ClientKey)
	assert.Equal(t, tc.ProviderURL, want.ProviderURL)
	assert.Equal(t, tc.Username, want.Username)
	assert.Nil(t, tc.Token)
}

func TestSetTokenNil(t *testing.T) {
	doNilTest(t, nil)
}

func TestSetTokenWrappedNil(t *testing.T) {
	doNilTest(t, WrapOAuth2Token(nil))
}

type tot struct {
	expIn        int
	refreshExpIn int
	tkn          *oauth2.Token
}

func (t tot) Token() *oauth2.Token          { return t.tkn }
func (t tot) ExpiresIn() (int, bool)        { return t.value(t.expIn) }
func (t tot) RefreshExpiresIn() (int, bool) { return t.value(t.refreshExpIn) }
func (t tot) value(i int) (int, bool) {
	if i == -1 {
		return -1, false
	} else {
		return i, true
	}
}

func TestSetToken(t *testing.T) {
	now := time.Now().Add(-10 * time.Minute)

	tbntime.WithTimeAt(now, func(cts tbntime.ControlledSource) {
		tc := TokenCache{timeSource: cts}
		tkn, _, _ := getToken()
		tc.SetToken(tot{30, 1800, tkn})

		assert.DeepEqual(
			t,
			tc.ExpiresAt,
			ptr.Time(
				now.Add(nowOffset).Add(30*time.Second).UTC(),
			),
		)
		assert.DeepEqual(
			t,
			tc.RefreshExpiresAt,
			ptr.Time(
				now.Add(nowOffset).Add(1800*time.Second).UTC(),
			),
		)
		assert.SameInstance(t, tc.Token, tkn)
	})
}

func TestTokenCacheExpired(t *testing.T) {
	now := time.Now().Add(-10 * time.Minute)

	type tc struct {
		setTokenNil bool
		ea          *time.Time
		rea         *time.Time
		want        bool
	}

	m := time.Minute

	testcases := []tc{
		{true, nil, nil, true},
		{true, &now, ptr.Time(now.Add(5 * m)), true},
		{false, ptr.Time(now.Add(1 * m)), ptr.Time(now.Add(4 * m)), false},
		{false, ptr.Time(now.Add(-1 * m)), ptr.Time(now.Add(4 * m)), false},
		{false, ptr.Time(now.Add(-4 * m)), ptr.Time(now.Add(-1 * m)), true},
		{false, ptr.Time(now.Add(-4 * m)), nil, false},
		{false, nil, ptr.Time(now.Add(-1 * m)), false},
	}

	// when checking expiration "now" + nowOffset will be the time that it's checked at
	tbntime.WithTimeAt(now, func(cts tbntime.ControlledSource) {
		for idx, c := range testcases {
			_, cache, _ := getToken()
			cache.timeSource = cts
			if c.setTokenNil {
				cache.Token = nil
			}

			cache.ExpiresAt = c.ea
			cache.RefreshExpiresAt = c.rea
			if !assert.Equal(t, cache.Expired(), c.want) {
				t.Errorf("Expired test case %v failed", idx)
			}
		}

	})
}

func TestToOAuthConfig(t *testing.T) {
	_, cache, _ := getToken()
	c, err := ToOAuthConfig(cache)
	assert.Nil(t, err)
	assert.Equal(t, c.ClientID, cache.ClientID)
	assert.Equal(t, c.ClientSecret, cache.ClientKey)
}

func TestToOAuthConfigBadURL(t *testing.T) {
	_, cache, _ := getToken()
	cache.ProviderURL = "bad"
	_, err := ToOAuthConfig(cache)
	assert.NonNil(t, err)
}

type stubExtra struct {
	t *testing.T
	k string
	v interface{}
}

func (s stubExtra) Extra(k string) interface{} {
	assert.Equal(s.t, k, s.k)
	return s.v
}

func TestIntExtra(t *testing.T) {
	type testcase struct {
		ret    interface{}
		want   int
		wantOk bool
	}

	cases := []testcase{
		{"str", 0, false},
		{true, 0, false},
		{nil, 0, false},
		{int(10), 10, true},
		{int8(10), 10, true},
		{int16(10), 10, true},
		{int32(10), 10, true},
		{int64(10), 10, true},
		{float32(10), 10, true},
		{float64(10), 10, true},
		{uint(10), 10, true},
		{uint8(10), 10, true},
		{uint16(10), 10, true},
		{uint32(10), 10, true},
		{uint64(10), 10, true},
		{uint8(math.MaxUint8), math.MaxUint8, true},
		{uint16(math.MaxUint16), math.MaxUint16, true},
		{uint64(math.MaxUint64), 0, false},
		// skip uint32 because I don't want to deal with making my test conditional
		// on sizeof(int)
	}

	key := "saonetuh"

	v, ok := intExtra(nil, key)
	if !assert.Equal(t, v, 0) {
		t.Errorf("nil test case failed value: %v", v)
	}
	if !assert.Equal(t, ok, false) {
		t.Errorf("nil test case failed ok: %v", ok)
	}

	for i, c := range cases {
		v, ok := intExtra(stubExtra{t, key, c.ret}, key)
		if !assert.Equal(t, v, c.want) {
			t.Errorf("Test case %v failed: wanted value %v got %v", i, c.want, v)
		}
		if !assert.Equal(t, ok, c.wantOk) {
			t.Errorf("Test case %v failed: wanted OK %v got %v", i, c.wantOk, ok)
		}
	}
}
