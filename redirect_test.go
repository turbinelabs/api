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
	"fmt"
	"testing"

	"github.com/turbinelabs/test/assert"
)

func getRedir() Redirect {
	return Redirect{
		"0-Name",
		"(.*)",
		"http://www.example.com?original=$1",
		PermanentRedirect,
	}
}

func TestRedirectEqualsTrue(t *testing.T) {
	r1 := getRedir()
	r2 := getRedir()
	assert.True(t, r1.Equals(r2))
	assert.True(t, r2.Equals(r1))
}

func TestRedirectEqualsFalseName(t *testing.T) {
	r1 := getRedir()
	r1.Name = "asoentuh"
	r2 := getRedir()
	assert.False(t, r1.Equals(r2))
	assert.False(t, r2.Equals(r1))
}

func TestRedirectEqualsFalseFrom(t *testing.T) {
	r1 := getRedir()
	r1.From = "asoentuh"
	r2 := getRedir()
	assert.False(t, r1.Equals(r2))
	assert.False(t, r2.Equals(r1))
}

func TestRedirectEqualsFalseTo(t *testing.T) {
	r1 := getRedir()
	r1.To = "asoentuh"
	r2 := getRedir()
	assert.False(t, r1.Equals(r2))
	assert.False(t, r2.Equals(r1))
}

func TestRedirectEqualsFalsePermanent(t *testing.T) {
	r1 := getRedir()
	r1.RedirectType = TemporaryRedirect
	r2 := getRedir()
	assert.False(t, r1.Equals(r2))
	assert.False(t, r2.Equals(r1))
}

func TestRedirectsEquals(t *testing.T) {
	rs1 := Redirects{getRedir(), getRedir(), getRedir()}
	rs2 := Redirects{getRedir(), getRedir(), getRedir()}
	assert.True(t, rs1.Equals(rs2))
	assert.True(t, rs2.Equals(rs1))
}

func TestRedirectsEqualsUnordered(t *testing.T) {
	r2 := getRedir()
	r2.Name = "asonethu"
	rs1 := Redirects{getRedir(), r2, getRedir()}
	rs2 := Redirects{getRedir(), getRedir(), r2}
	assert.False(t, rs1.Equals(rs2))
	assert.False(t, rs2.Equals(rs1))
}

func TestRedirectsEqualsFalseLength(t *testing.T) {
	r2 := getRedir()
	r2.Name = "asonethu"
	rs1 := Redirects{getRedir(), r2}
	rs2 := Redirects{getRedir(), getRedir(), r2}
	assert.False(t, rs1.Equals(rs2))
	assert.False(t, rs2.Equals(rs1))
}

func TestRedirectsEqualsFalse(t *testing.T) {
	r1 := getRedir()
	r1.Name = "r1"
	r2 := getRedir()
	r2.Name = "asonethu"
	rs1 := Redirects{getRedir(), r2}
	rs2 := Redirects{getRedir(), r1}
	assert.False(t, rs1.Equals(rs2))
	assert.False(t, rs2.Equals(rs1))
}

func TestRedirectIsValid(t *testing.T) {
	r := getRedir()
	assert.Nil(t, r.IsValid())
}

func TestRedirectIsValidFailsNameBadChar(t *testing.T) {
	r := getRedir()
	r.Name = "aoeu snth"
	assert.DeepEqual(t,
		r.IsValid(),
		&ValidationError{[]ErrorCase{{
			"redirects[aoeu snth].name",
			fmt.Sprintf("must match %s", NamePattern.String()),
		}}},
	)
}

func TestRedirectIsValidFailsNoName(t *testing.T) {
	r := getRedir()
	r.Name = ""
	assert.DeepEqual(t,
		r.IsValid(),
		&ValidationError{[]ErrorCase{{"redirects[].name", "must not be empty"}}},
	)
}

func TestRedirectIsValidFailsNoFrom(t *testing.T) {
	r := getRedir()
	r.From = ""
	assert.DeepEqual(t,
		r.IsValid(),
		&ValidationError{[]ErrorCase{{"redirects[0-Name].from", "must not be empty"}}},
	)
}

func TestRedirectIsValidFailsNoTo(t *testing.T) {
	r := getRedir()
	r.To = ""
	assert.DeepEqual(t,
		r.IsValid(),
		&ValidationError{[]ErrorCase{{"redirects[0-Name].to", "must not be empty"}}},
	)
}

func TestRedirectIsValidFailsBadRegex(t *testing.T) {
	r := getRedir()
	r.From = `asonetu(asontehu`
	err := r.IsValid()
	assert.DeepEqual(
		t,
		err,
		&ValidationError{[]ErrorCase{{
			"redirects[0-Name].from",
			"invalid url match expression 'error parsing regexp: missing closing ): `asonetu(asontehu`'",
		}}},
	)
}

func TestRedirectsIsValid(t *testing.T) {
	r1 := getRedir()
	r1.Name = "r1"
	r2 := getRedir()
	r2.Name = "r2"

	rs := Redirects{r1, r2}
	err := rs.IsValid()

	assert.Nil(t, err)
}

func TestRedirectsAsMap(t *testing.T) {
	r1 := getRedir()
	r1.Name = "r1"
	r2 := getRedir()
	r2.Name = "r2"
	rs := Redirects{r1, r2}

	assert.DeepEqual(t, rs.AsMap(), map[string]Redirect{"r1": r1, "r2": r2})
}

func TestRedirectsAsMapNil(t *testing.T) {
	var rs Redirects
	assert.DeepEqual(t, rs.AsMap(), map[string]Redirect{})
}

func TestRedirectsKeys(t *testing.T) {
	r1 := getRedir()
	r1.Name = "r1"
	r2 := getRedir()
	r2.Name = "r2"
	rs := Redirects{r1, r2}

	assert.DeepEqual(t, rs.Keys(), []string{"r1", "r2"})
}

func TestRedirectsKeysNil(t *testing.T) {
	var rs Redirects
	assert.HasSameElements(t, rs.Keys(), []string{})
}
