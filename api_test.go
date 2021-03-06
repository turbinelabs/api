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

func getVe() ValidationError {
	return ValidationError{[]ErrorCase{{"a", "b"}}}
}

func pf(b bool) string {
	if b {
		return "passed"
	} else {
		return "failed"
	}
}

func TestKeyPattern(t *testing.T) {
	test := func(in string, shouldMatch bool) {
		fmt.Printf(
			"%v - '%v'\n",
			pf(assert.Equal(t, KeyPattern.MatchString(in), shouldMatch)),
			in)
	}

	pass := func(in string) { test(in, true) }
	fail := func(in string) { test(in, false) }

	pass("1234")
	fail(" 1234")
	fail("aoeu-snth-102938-...")
	fail("aoeu-[9982829")
	pass("25db52ba-6f12-4964-6a3f-40f2e9ad47bb")
}

func TestIndexPattern(t *testing.T) {
	test := func(in string, shouldMatch bool) {
		fmt.Printf(
			"%v - '%v'\n",
			pf(assert.Equal(t, AllowedIndexPattern.MatchString(in), shouldMatch)),
			in)
	}

	pass := func(in string) { test(in, true) }
	fail := func(in string) { test(in, false) }

	pass("12341234")
	pass("a/b/c/d")
	pass("a\\b\\c\\d")
	pass("a-b-3--")
	pass("a b 3  ")
	fail("aoeu-[-snth")
	fail("aoeu-]-snth")
}

func TestValidationErrorAddNew(t *testing.T) {
	ve := getVe()
	ve.AddNew(ErrorCase{"b", "c"})

	assert.DeepEqual(t, ve, ValidationError{[]ErrorCase{{"a", "b"}, {"b", "c"}}})
}

func TestValidationErrorMergeEmpty(t *testing.T) {
	ve1 := getVe()
	ve2 := ValidationError{}
	ve1.Merge(&ve2)

	assert.DeepEqual(t, getVe(), ve1)
}

func TestValidationErrorMergeIntoEmpty(t *testing.T) {
	ve1 := ValidationError{}
	ve2 := getVe()
	ve1.Merge(&ve2)

	assert.DeepEqual(t, ve1, getVe())
}

func TestValidationErrorMergePrefixed(t *testing.T) {
	ve := getVe()
	ve2a := ValidationError{[]ErrorCase{
		{"child", "msg"},
		{"child2", "msg2"},
	}}
	ve2b := ValidationError{[]ErrorCase{
		{"child", "msg"},
		{"child2", "msg2"},
	}}

	ve.MergePrefixed(&ve2a, "parent")

	assert.DeepEqual(t, ve2a, ve2b)
	assert.DeepEqual(t, ve, ValidationError{[]ErrorCase{
		{"a", "b"},
		{"parent.child", "msg"},
		{"parent.child2", "msg2"},
	}})
}

func TestValidationErrorOrNil(t *testing.T) {
	ve := getVe()
	assert.Equal(t, &ve, ve.OrNil())
}

func TestValidationErrorOrNilNoErrors(t *testing.T) {
	ve := ValidationError{}
	assert.Nil(t, ve.OrNil())
}
