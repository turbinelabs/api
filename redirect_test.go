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
		HeaderConstraints{{"x-tbn-api-key", "", false, false}},
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
			"name",
			fmt.Sprintf("must match %s", HeaderNamePatternStr),
		}}},
	)
}

func TestRedirectIsValidFailsNoName(t *testing.T) {
	r := getRedir()
	r.Name = ""
	assert.DeepEqual(t,
		r.IsValid(),
		&ValidationError{[]ErrorCase{{"name", "must not be empty"}}},
	)
}

func TestRedirectIsValidFailsNoFrom(t *testing.T) {
	r := getRedir()
	r.From = ""
	assert.DeepEqual(t,
		r.IsValid(),
		&ValidationError{[]ErrorCase{{"from", "must not be empty"}}},
	)
}

func TestRedirectIsValidFailsNoTo(t *testing.T) {
	r := getRedir()
	r.To = ""
	assert.DeepEqual(t,
		r.IsValid(),
		&ValidationError{[]ErrorCase{{"to", "must not be empty"}}},
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
			"from",
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

func TestRedirectsIsValidDupes(t *testing.T) {
	r1 := getRedir()
	r1.Name = "r-name"
	r2 := r1

	rs := Redirects{r1, r2}
	err := rs.IsValid()

	assert.DeepEqual(t, err, &ValidationError{[]ErrorCase{
		{"redirects", "name must be unique, multiple redirects found called 'r-name'"},
	}})
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

func doTestHeaderConstraintEquals(t *testing.T, hc1, hc2 HeaderConstraint, eq bool) {
	if eq {
		assert.True(t, hc1.Equals(hc2))
		assert.True(t, hc2.Equals(hc1))
	} else {
		assert.False(t, hc1.Equals(hc2))
		assert.False(t, hc2.Equals(hc1))
	}
}

func getHC() HeaderConstraint {
	return HeaderConstraint{"name", "value", false, false}
}

func TestHeaderConstraintEquals(t *testing.T) {
	hc1 := getHC()
	hc2 := hc1
	doTestHeaderConstraintEquals(t, hc1, hc2, true)
}

func TestHeaderConstraintEqualsFalseCaseSensitivity(t *testing.T) {
	hc1 := getHC()
	hc2 := hc1
	hc2.CaseSensitive = !hc1.CaseSensitive
	doTestHeaderConstraintEquals(t, hc1, hc2, false)
}

func TestHeaderConstraintEqualsFalseName(t *testing.T) {
	hc1 := getHC()
	hc2 := hc1
	hc2.Name = "aoeu"
	doTestHeaderConstraintEquals(t, hc1, hc2, false)
}

func TestHeaderConstraintEqualsFalseValue(t *testing.T) {
	hc1 := getHC()
	hc2 := hc1
	hc2.Value = "aoeu"
	doTestHeaderConstraintEquals(t, hc1, hc2, false)
}

func TestHeaderConstraintEqualsFalseInvert(t *testing.T) {
	hc1 := getHC()
	hc2 := hc1
	hc2.Invert = !hc1.Invert
	doTestHeaderConstraintEquals(t, hc1, hc2, false)
}

func TestHeaderConstraintIsValid(t *testing.T) {
	hc1 := getHC()
	hc1.Name = "x-header-name"
	err := hc1.IsValid()
	assert.Nil(t, err)
}

func doTestHeaderConstraintIsValidFails(t *testing.T, hc HeaderConstraint, f, m string) {
	err := hc.IsValid()
	assert.DeepEqual(t, err, &ValidationError{[]ErrorCase{{f, m}}})
}

func TestHeaderConstraintIsValidFalseNamePattern1(t *testing.T) {
	hc1 := getHC()
	hc1.Name = "na me"
	doTestHeaderConstraintIsValidFails(t, hc1, "name", fmt.Sprintf("must match %s", HeaderNamePatternStr))
}

func TestHeaderConstraintIsValidFalseNamePattern2(t *testing.T) {
	hc1 := getHC()
	hc1.Name = "na;me"
	doTestHeaderConstraintIsValidFails(t, hc1, "name", fmt.Sprintf("must match %s", HeaderNamePatternStr))
}

func TestHeaderConstraintIsValidFalseNamePattern3(t *testing.T) {
	hc1 := getHC()
	hc1.Name = "na_me"
	doTestHeaderConstraintIsValidFails(t, hc1, "name", fmt.Sprintf("must match %s", HeaderNamePatternStr))
}

func TestHeaderConstraintIsValidFalseNameEmpty(t *testing.T) {
	hc1 := getHC()
	hc1.Name = ""
	doTestHeaderConstraintIsValidFails(t, hc1, "name", "may not be empty")
}

func TestHeaderConstraintsIsValidBadValueRegexp(t *testing.T) {
	hc1 := getHC()
	hc1.Value = "aoeu\\onth"
	doTestHeaderConstraintIsValidFails(
		t,
		hc1,
		"value",
		"must be a valid regexp: error parsing regexp: invalid escape sequence: `\\o`")
}

func TestHeaderConstraintsIsValid(t *testing.T) {
	hc1 := getHC()
	hcs := HeaderConstraints{hc1}
	assert.Nil(t, hcs.IsValid())
}

func TestHeaderConstraintsIsValidFalseNumber(t *testing.T) {
	hc1 := getHC()
	hc2 := getHC()
	hc2.Name = "name2"
	hcs := HeaderConstraints{hc1, hc2}

	assert.DeepEqual(t, hcs.IsValid(), &ValidationError{[]ErrorCase{
		{"header_constraints", "may only specify 0 or 1 header constraints"},
	}})
}

func TestHeaderConstraintsIsValidFalseDupe(t *testing.T) {
	hc1 := getHC()
	hc2 := getHC()
	hcs := HeaderConstraints{hc1, hc2}

	assert.DeepEqual(t, hcs.IsValid(), &ValidationError{[]ErrorCase{
		{"header_constraints", "may only specify 0 or 1 header constraints"},
		{"header_constraints[name]", "a header may only have a single constraint"},
	}})
}

func TestHeaderConstraintsIsValidFalseNested(t *testing.T) {
	hc1 := getHC()
	hc1.Name = "na;me"
	hcs := HeaderConstraints{hc1}

	assert.DeepEqual(t, hcs.IsValid(), &ValidationError{[]ErrorCase{
		{"header_constraints[na;me].name", fmt.Sprintf("must match %s", HeaderNamePatternStr)},
	}})
}
