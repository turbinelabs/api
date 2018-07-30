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
	"github.com/turbinelabs/test/check"
	"github.com/turbinelabs/test/matcher"
)

type validationErrorMatcher struct {
	attr, msg string
}

func ValidationMatcher(a, m string) matcher.Matcher {
	return &validationErrorMatcher{attr: a, msg: m}
}

func (vem *validationErrorMatcher) Matches(x interface{}) bool {
	got, ok := x.(*ValidationError)
	if !ok {
		fmt.Printf(
			"wrong got type: %+v (%T), expected %T\n",
			x,
			x,
			&ValidationError{},
		)
	}

	eq, _ := check.DeepEqual(
		got,
		&ValidationError{
			[]ErrorCase{
				{
					vem.attr,
					vem.msg,
				},
			},
		},
	)

	return eq
}

func (vem *validationErrorMatcher) String() string {
	return fmt.Sprintf(
		"validationErrorMatcher(attribute=%s, msg=%s)",
		vem.attr,
		vem.msg,
	)
}

func TestMatchEqualsSuccess(t *testing.T) {
	m1 := Match{
		Kind:     HeaderMatchKind,
		Behavior: ExactMatchBehavior,
		From:     Metadatum{Key: "x-random", Value: "value"},
		To:       Metadatum{Key: "randomflag", Value: "true"},
	}
	m2 := Match{
		Kind:     HeaderMatchKind,
		Behavior: ExactMatchBehavior,
		From:     Metadatum{Key: "x-random", Value: "value"},
		To:       Metadatum{Key: "randomflag", Value: "true"},
	}

	assert.True(t, m1.Equals(m2))
	assert.True(t, m2.Equals(m1))
}

func TestMatchEqualsFromVaries(t *testing.T) {
	m1 := Match{
		Kind:     HeaderMatchKind,
		Behavior: ExactMatchBehavior,
		From:     Metadatum{Key: "x-random", Value: "value"},
		To:       Metadatum{Key: "randomflag", Value: "true"},
	}
	m2 := Match{
		Kind:     HeaderMatchKind,
		Behavior: ExactMatchBehavior,
		From:     Metadatum{Key: "x-other", Value: "value"},
		To:       Metadatum{Key: "randomflag", Value: "true"},
	}

	assert.False(t, m1.Equals(m2))
	assert.False(t, m2.Equals(m1))
}

func TestMatchEqualsToVaries(t *testing.T) {
	m1 := Match{
		Kind:     HeaderMatchKind,
		Behavior: ExactMatchBehavior,
		From:     Metadatum{Key: "x-random", Value: "value"},
		To:       Metadatum{Key: "randomflag", Value: "true"},
	}
	m2 := Match{
		Kind:     HeaderMatchKind,
		Behavior: ExactMatchBehavior,
		From:     Metadatum{Key: "x-random", Value: "value"},
		To:       Metadatum{Key: "specificflag", Value: "true"},
	}

	assert.False(t, m1.Equals(m2))
	assert.False(t, m2.Equals(m1))
}

func TestMatchEqualsMatchKindVaries(t *testing.T) {
	m1 := Match{
		Kind:     HeaderMatchKind,
		Behavior: ExactMatchBehavior,
		From:     Metadatum{Key: "x-random", Value: "value"},
		To:       Metadatum{Key: "randomflag", Value: "true"},
	}
	m2 := Match{
		Kind:     CookieMatchKind,
		Behavior: ExactMatchBehavior,
		From:     Metadatum{Key: "x-random", Value: "value"},
		To:       Metadatum{Key: "randomflag", Value: "true"},
	}

	assert.False(t, m1.Equals(m2))
	assert.False(t, m2.Equals(m1))
}

func TestMatchEqualsMatchBehaviorVaries(t *testing.T) {
	m1 := Match{
		Kind:     HeaderMatchKind,
		Behavior: ExactMatchBehavior,
		From:     Metadatum{Key: "x-random", Value: "value"},
		To:       Metadatum{Key: "randomflag", Value: "true"},
	}
	m2 := Match{
		Kind:     HeaderMatchKind,
		Behavior: RegexMatchBehavior,
		From:     Metadatum{Key: "x-random", Value: "value"},
		To:       Metadatum{Key: "randomflag", Value: "true"},
	}

	assert.False(t, m1.Equals(m2))
	assert.False(t, m2.Equals(m1))
}

func TestMatchesEqualsSuccess(t *testing.T) {
	m1 := Match{
		Kind:     HeaderMatchKind,
		Behavior: ExactMatchBehavior,
		From:     Metadatum{Key: "x-random", Value: "value"},
		To:       Metadatum{Key: "randomflag", Value: "true"},
	}
	m2 := Match{
		Kind:     CookieMatchKind,
		Behavior: ExactMatchBehavior,
		From:     Metadatum{Key: "x-random", Value: "value"},
		To:       Metadatum{Key: "randomflag", Value: "true"},
	}

	slice1 := Matches{m1, m2}
	slice2 := Matches{m2, m1}

	assert.True(t, slice1.Equals(slice2))
	assert.True(t, slice2.Equals(slice1))
}

func TestMatchesEqualsLengthMismatch(t *testing.T) {
	m1 := Match{
		Kind:     HeaderMatchKind,
		Behavior: ExactMatchBehavior,
		From:     Metadatum{Key: "x-random", Value: "value"},
		To:       Metadatum{Key: "randomflag", Value: "true"},
	}
	m2 := Match{
		Kind:     CookieMatchKind,
		Behavior: ExactMatchBehavior,
		From:     Metadatum{Key: "x-random", Value: "value"},
		To:       Metadatum{Key: "randomflag", Value: "true"},
	}

	slice1 := Matches{m1, m2}
	slice2 := Matches{m2}

	assert.False(t, slice1.Equals(slice2))
	assert.False(t, slice2.Equals(slice1))
}

func TestMatchesEqualsContentDiffers(t *testing.T) {
	m1 := Match{
		Kind:     HeaderMatchKind,
		Behavior: ExactMatchBehavior,
		From:     Metadatum{Key: "x-random", Value: "value"},
		To:       Metadatum{Key: "randomflag", Value: "true"},
	}
	m2 := Match{
		Kind:     CookieMatchKind,
		Behavior: ExactMatchBehavior,
		From:     Metadatum{Key: "x-random", Value: "value"},
		To:       Metadatum{Key: "randomflag", Value: "true"},
	}
	m3 := Match{
		Kind:     CookieMatchKind,
		Behavior: ExactMatchBehavior,
		From:     Metadatum{Key: "x-random", Value: "value"},
		To:       Metadatum{Key: "randomflag", Value: "false"},
	}

	slice1 := Matches{m1, m2}
	slice2 := Matches{m2, m3}

	assert.False(t, slice1.Equals(slice2))
	assert.False(t, slice2.Equals(slice1))
}

func TestMatchIsValidSucces(t *testing.T) {
	m := Match{
		Kind:     CookieMatchKind,
		Behavior: ExactMatchBehavior,
		From:     Metadatum{Key: "aoeu", Value: ""},
		To:       Metadatum{Key: "snth", Value: "1234"},
	}

	assert.Nil(t, m.IsValid())
}

func TestMatchIsValidFailedFromBadKey(t *testing.T) {
	m := Match{
		Kind:     CookieMatchKind,
		Behavior: ExactMatchBehavior,
		From:     Metadatum{Key: "from]", Value: ""},
		To:       Metadatum{Key: "snth", Value: "1234"},
	}

	assert.NonNil(t, m.IsValid())
}

func TestMatchIsValidFailedFromEmpty(t *testing.T) {
	m := Match{
		Kind:     CookieMatchKind,
		Behavior: ExactMatchBehavior,
		From:     Metadatum{},
		To:       Metadatum{Key: "snth", Value: "1234"},
	}

	assert.NonNil(t, m.IsValid())
}

func TestMatchIsValidSuccessEmptyTo(t *testing.T) {
	m := Match{
		Kind:     CookieMatchKind,
		Behavior: ExactMatchBehavior,
		From:     Metadatum{Key: "aoeu", Value: ""},
		To:       Metadatum{Key: "", Value: ""},
	}

	assert.Nil(t, m.IsValid())
}

func TestMatchIsValidFailedKind(t *testing.T) {
	m := Match{
		Kind:     "snth",
		Behavior: ExactMatchBehavior,
		From:     Metadatum{Key: "aoeu", Value: ""},
		To:       Metadatum{Key: "snth", Value: "1234"},
	}

	assert.NonNil(t, m.IsValid())
}

func TestMatchIsValidFailedToNoKey(t *testing.T) {
	m := Match{
		Kind:     "aoeu",
		Behavior: ExactMatchBehavior,
		From:     Metadatum{Key: "snth", Value: ""},
		To:       Metadatum{Key: "", Value: "snth"},
	}

	assert.NonNil(t, m.IsValid())
}

func TestMatchIsValidFailedBadMatchBehavior(t *testing.T) {
	m := Match{
		Kind:     CookieMatchKind,
		Behavior: "garbage",
		From:     Metadatum{Key: "aoeu", Value: ""},
		To:       Metadatum{Key: "snth", Value: "1234"},
	}

	vm := ValidationMatcher("behavior", `"garbage" is not a valid behavior kind`)
	assert.True(t, vm.Matches(m.IsValid()))
}

func TestMatchIsValidFailedMatchBehaviorRangeEmptyFromValue(t *testing.T) {
	m := Match{
		Kind:     HeaderMatchKind,
		Behavior: RangeMatchBehavior,
		From:     Metadatum{Key: "aoeu", Value: ""},
		To:       Metadatum{Key: "snth", Value: "1234"},
	}

	vm := ValidationMatcher("from.value", `must not be empty if behavior is "range"`)
	assert.True(t, vm.Matches(m.IsValid()))
}

func TestMatchIsValidFailedMatchBehaviorRegexEmptyFromValue(t *testing.T) {
	m := Match{
		Kind:     HeaderMatchKind,
		Behavior: RegexMatchBehavior,
		From:     Metadatum{Key: "aoeu", Value: ""},
		To:       Metadatum{Key: "snth", Value: "1234"},
	}

	vm := ValidationMatcher("from.value", `must not be empty if behavior is "regex"`)
	assert.True(t, vm.Matches(m.IsValid()))
}

func TestMatchIsValidFailedMatchBehaviorRegexInvalidFromValue(t *testing.T) {
	tcs := []struct {
		input, errStr string
	}{
		{
			input:  "\\",
			errStr: "error parsing regexp: trailing backslash at end of expression: ``",
		},
		{
			input:  "[/",
			errStr: "error parsing regexp: missing closing ]: `[/`",
		},
		{
			input:  "^[]",
			errStr: "error parsing regexp: missing closing ]: `[]`",
		},
		{
			input:  "*.",
			errStr: "error parsing regexp: missing argument to repetition operator: `*`",
		},
	}
	for i, tc := range tcs {
		assert.Group(
			fmt.Sprintf("testCases[%d]: regex=[%q]", i, tc.input),
			t,
			func(g *assert.G) {
				m := Match{
					Kind:     HeaderMatchKind,
					Behavior: RegexMatchBehavior,
					From:     Metadatum{Key: "a_key", Value: tc.input},
					To:       Metadatum{Key: "snth", Value: "1234"},
				}

				vm := ValidationMatcher("from.value", tc.errStr)
				assert.True(g, vm.Matches(m.IsValid()))
			},
		)
	}
}

func TestMatchIsValidFailedRegexMatchBehaviorWrongNumberOfCaptureGroups(t *testing.T) {
	for i, br := range []string{
		"([a-zA-Z]+) ([a-zA-Z]+)",
		"(.*),(.*),(.*)",
	} {
		assert.Group(
			fmt.Sprintf("testCases[%d]: regex=[%q]", i, br),
			t,
			func(g *assert.G) {
				m := Match{
					Kind:     CookieMatchKind,
					Behavior: RegexMatchBehavior,
					From:     Metadatum{Key: "a_key", Value: br},
					To:       Metadatum{Key: "snth"},
				}

				vm := ValidationMatcher(
					"from.value",
					`must have exactly one subgroup when to.value is not set`,
				)

				assert.True(g, vm.Matches(m.IsValid()))
			},
		)
	}
}

func TestMatchIsValidSuccessMatchBehaviorRegexSingleSubgroup(t *testing.T) {
	for i, re := range []string{"([a-zA-Z]+)", "[a-zA-Z]+", ".*"} {
		assert.Group(
			fmt.Sprintf("testCases[%d]: regex=[%q]", i, re),
			t,
			func(g *assert.G) {
				m := Match{
					Kind:     CookieMatchKind,
					Behavior: RegexMatchBehavior,
					From:     Metadatum{Key: "a_key", Value: re},
					To:       Metadatum{Key: "snth"},
				}

				assert.Nil(g, m.IsValid())
			},
		)
	}
}

func TestMatchIsValidFailedMatchBehaviorRangeFromValueInvalidFormat(t *testing.T) {
	for i, r := range []string{
		"3,4",
		"1.1,2",
		"[1.1,10)",
		"[-100, +200]",
		"[+-100, +-200)",
		"[+blow, +up)",
		" [ -10asdf0, 200 asdf) ",
	} {
		assert.Group(
			fmt.Sprintf("testCases[%d]: range=[%q]", i, r),
			t,
			func(g *assert.G) {
				m := Match{
					Kind:     HeaderMatchKind,
					Behavior: RangeMatchBehavior,
					From:     Metadatum{Key: "a_key", Value: r},
					To:       Metadatum{Key: "snth"},
				}

				vm := ValidationMatcher(
					"from.value",
					fmt.Sprintf(
						`Invalid range pattern: %q. Must be of the form "[<start>,<end>)"`,
						r,
					),
				)

				assert.True(g, vm.Matches(m.IsValid()))
			},
		)
	}
}

func TestMatchIsvalidFailedMatchBehaviorRangeFromValueInvalid(t *testing.T) {
	tcs := []struct {
		input, errStr string
	}{
		{
			input:  "[-21474836483434324324,100)",
			errStr: `strconv.Atoi: parsing "-21474836483434324324": value out of range`,
		},
		{
			input:  "[100, 21474836483434324324)",
			errStr: `strconv.Atoi: parsing "21474836483434324324": value out of range`,
		},
	}

	for i, tc := range tcs {
		assert.Group(
			fmt.Sprintf("testCases[%d]: range=[%q]", i, tc.input),
			t,
			func(g *assert.G) {
				m := Match{
					Kind:     HeaderMatchKind,
					Behavior: RangeMatchBehavior,
					From:     Metadatum{Key: "a_key", Value: tc.input},
					To:       Metadatum{Key: "snth"},
				}

				vm := ValidationMatcher("from.value", tc.errStr)
				assert.True(g, vm.Matches(m.IsValid()))
			},
		)
	}
}

func TestMatchIsValidFailedMatchBehaviorRangeFromValuesInvalid(t *testing.T) {
	for i, r := range []string{
		"[0,0)",
		"[100,-1)",
		"[-2,-3)",
		"[1,0)",
	} {
		assert.Group(
			fmt.Sprintf("testCases[%d]: range=[%q]", i, r),
			t,
			func(g *assert.G) {
				m := Match{
					Kind:     HeaderMatchKind,
					Behavior: RangeMatchBehavior,
					From:     Metadatum{Key: "a_key", Value: r},
					To:       Metadatum{Key: "snth"},
				}

				vm := ValidationMatcher(
					"from.value",
					"End of range must be greater than start",
				)
				assert.True(g, vm.Matches(m.IsValid()))
			},
		)
	}
}

func TestMatchIsValidSucceedsOnValidRanges(t *testing.T) {
	for i, r := range []string{
		"[-300,100)",
		"[-300,-100)",
		"[100,200)",
	} {
		assert.Group(
			fmt.Sprintf("testCases[%d]: range=[%q]", i, r),
			t,
			func(g *assert.G) {
				m := Match{
					Kind:     HeaderMatchKind,
					Behavior: RangeMatchBehavior,
					From:     Metadatum{Key: "a_key", Value: r},
					To:       Metadatum{Key: "snth"},
				}

				assert.Nil(g, m.IsValid())
			},
		)
	}
}

func TestMatchIsValidQueryMatchWithRegexBehavior(t *testing.T) {
	m := Match{
		Kind:     QueryMatchKind,
		Behavior: RegexMatchBehavior,
		From:     Metadatum{Key: "q1", Value: "(.*)"},
		To:       Metadatum{Key: "stage", Value: "testing"},
	}

	vm := ValidationMatcher(
		"behavior",
		`"query" kind not supported with "regex" behavior`,
	)
	assert.True(t, vm.Matches(m.IsValid()))
}

func TestMatchIsValidCookieMatchWithRangeBehavior(t *testing.T) {
	m := Match{
		Kind:     CookieMatchKind,
		Behavior: RangeMatchBehavior,
		From:     Metadatum{Key: "q1", Value: "[1,4)"},
		To:       Metadatum{Key: "stage", Value: "testing"},
	}

	vm := ValidationMatcher(
		"kind",
		`"cookie" kind not supported with "regex" behavior`,
	)
	assert.True(t, vm.Matches(m.IsValid()))
}

func TestMatchesIsValidSuccess(t *testing.T) {
	m1 := Match{
		Kind:     CookieMatchKind,
		Behavior: ExactMatchBehavior,
		From:     Metadatum{Key: "aoeu", Value: ""},
		To:       Metadatum{Key: "snth", Value: "1234"},
	}

	m2 := Match{
		Kind:     HeaderMatchKind,
		Behavior: ExactMatchBehavior,
		From:     Metadatum{Key: "aoeu", Value: ""},
		To:       Metadatum{Key: "snth", Value: "1234"},
	}

	m3 := Match{
		Kind:     QueryMatchKind,
		Behavior: ExactMatchBehavior,
		From:     Metadatum{Key: "aoeu", Value: ""},
		To:       Metadatum{Key: "snth", Value: "1234"},
	}

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
	m1 := Match{
		Kind:     CookieMatchKind,
		Behavior: ExactMatchBehavior,
		From:     Metadatum{Key: "aoeu", Value: ""},
		To:       Metadatum{Key: "snth", Value: "1234"},
	}
	m2 := Match{
		Kind:     "badmatch",
		Behavior: ExactMatchBehavior,
		From:     Metadatum{Key: "aoeu", Value: ""},
		To:       Metadatum{Key: "snth", Value: "1234"},
	}
	m3 := Match{
		Kind:     QueryMatchKind,
		Behavior: ExactMatchBehavior,
		From:     Metadatum{Key: "aoeu", Value: ""},
		To:       Metadatum{Key: "snth", Value: "1234"},
	}
	m := Matches{m1, m2, m3}

	vm := ValidationMatcher(
		"matches[badmatch:exact:aoeu].kind",
		`"badmatch" is not a valid match kind`,
	)
	assert.True(t, vm.Matches(m.IsValid()))
}

func TestMatchesIsValidDupeMatch(t *testing.T) {
	m1 := Match{
		Kind:     CookieMatchKind,
		Behavior: ExactMatchBehavior,
		From:     Metadatum{Key: "type", Value: "chocolate chip"},
		To:       Metadatum{Key: "texture", Value: "dense"},
	}

	m2 := Match{
		Kind:     CookieMatchKind,
		Behavior: ExactMatchBehavior,
		From:     Metadatum{Key: "type", Value: "amarettie"},
		To:       Metadatum{Key: "texture", Value: "chewy"},
	}

	m := Matches{m1, m2}
	vm := ValidationMatcher(
		"",
		"duplicate match found cookie:exact:type",
	)
	assert.True(t, vm.Matches(m.IsValid()))
}

func TestMatchUnmarshalJSONDefaultsNullMatchBehavior(t *testing.T) {
	bytes := []byte(`{
          "kind": "header",
          "from": {
            "key": "x-version",
            "value": ""
          },
          "to": {
            "key": "version",
            "value": ""
          }
        }`)

	var m Match
	err := m.UnmarshalJSON(bytes)
	assert.Nil(t, err)
	assert.Equal(t, m.Behavior, ExactMatchBehavior)
}

func TestMatchUnmarshalJSONDefaultsEmptyMatchBehavior(t *testing.T) {
	bytes := []byte(`{
          "kind": "header",
          "behavior": "",
          "from": {
            "key": "x-version",
            "value": ""
          },
          "to": {
            "key": "version",
            "value": ""
          }
        }`)

	var m Match
	err := m.UnmarshalJSON(bytes)
	assert.Nil(t, err)
	assert.Equal(t, m.Behavior, ExactMatchBehavior)
}

func TestMatchUnmarshalJSONPassesThroughDefinedBehavior(t *testing.T) {
	bytes := []byte(`{
          "kind": "header",
          "behavior": "regex",
          "from": {
            "key": "x-version",
            "value": ""
          },
          "to": {
            "key": "version",
            "value": ""
          }
        }`)

	var m Match
	err := m.UnmarshalJSON(bytes)
	assert.Nil(t, err)
	assert.Equal(t, m.Behavior, RegexMatchBehavior)
}
