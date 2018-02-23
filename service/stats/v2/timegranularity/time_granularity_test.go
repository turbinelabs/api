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

// This file was automatically generated by doc.go from ../../enum_test.template.
// Any changes will be lost if this file is regenerated.

package timegranularity

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/turbinelabs/test/assert"
)

type testStruct struct {
	TimeGranularity TimeGranularity `json:"time_granularity"`
}

func assertHasAllValues(t testing.TB, i interface{}) {
	mappedValues := map[TimeGranularity]struct{}{}

	switch m := i.(type) {
	case map[TimeGranularity]string:
		for qt := range m {
			mappedValues[qt] = struct{}{}
		}
	case map[string]TimeGranularity:
		for _, qt := range m {
			mappedValues[qt] = struct{}{}
		}
	default:
		t.Fatalf("cannot validate type %T", i)
		return
	}

	ForEach(func(v TimeGranularity) {
		assert.Group(
			fmt.Sprintf("TimeGranularity %s", v.String()),
			t,
			func(g *assert.G) {
				_, ok := mappedValues[v]
				assert.True(g, ok)
			},
		)
	})
}

func TestTimeGranularityString(t *testing.T) {
	assert.Equal(
		t,
		Hours.String(),
		"hours",
	)
	assert.MatchesRegex(t, Unknown.String(), `unknown\([0-9]+\)`)
	assert.Equal(t, TimeGranularity(100).String(), "unknown(100)")
}

func TestIsValid(t *testing.T) {
	invalid := []TimeGranularity{
		TimeGranularity(minTimeGranularity - 1),
		TimeGranularity(maxTimeGranularity + 1),
	}

	for _, qt := range invalid {
		assert.False(t, IsValid(qt))
	}

	ForEach(func(qt TimeGranularity) {
		assert.True(t, IsValid(qt))
	})
}

func TestFromName(t *testing.T) {
	validValues := map[TimeGranularity]string{
		Minutes: strMinutes,
		Hours:   strHours,
	}
	assertHasAllValues(t, validValues)

	for expectedQt, name := range validValues {
		qt := FromName(name)
		assert.Equal(t, qt, expectedQt)
	}

	invalidValues := []string{"bob", "unknown", "1"}

	for _, name := range invalidValues {
		qt := FromName(name)
		assert.Equal(t, qt, Unknown)
	}
}

func TestTimeGranularityMarshalJSON(t *testing.T) {
	vals := map[TimeGranularity]string{
		Minutes: strMinutes,
		Hours:   strHours,
	}
	assertHasAllValues(t, vals)

	for v, name := range vals {
		bytes, err := v.MarshalJSON()
		assert.Nil(t, err)
		expected := []byte(fmt.Sprintf(`"%s"`, name))
		assert.DeepEqual(t, bytes, expected)
	}
}

func TestTimeGranularityMarshalJSONUnknown(t *testing.T) {
	unknownValues := []TimeGranularity{
		Unknown,
		TimeGranularity(maxTimeGranularity + 1),
	}

	for _, unknownTimeGranularity := range unknownValues {
		bytes, err := unknownTimeGranularity.MarshalJSON()
		assert.Nil(t, bytes)
		assert.ErrorContains(t, err, "cannot marshal unknown TimeGranularity")
	}
}

func TestTimeGranularityMarshalJSONNil(t *testing.T) {
	var v *TimeGranularity

	bytes, err := v.MarshalJSON()
	assert.ErrorContains(t, err, "cannot marshal unknown")
	assert.Nil(t, bytes)
}

func TestTimeGranularityUnmarshalJSON(t *testing.T) {
	quoted := func(s string) string {
		return fmt.Sprintf(`"%s"`, s)
	}

	vals := map[string]TimeGranularity{
		quoted(strMinutes): Minutes,
		quoted(strHours):   Hours,
	}
	assertHasAllValues(t, vals)

	for data, expected := range vals {
		var v TimeGranularity

		err := v.UnmarshalJSON([]byte(data))
		assert.Nil(t, err)
		assert.Equal(t, v, expected)
	}
}

func TestTimeGranularityUnmarshalJSONUnknown(t *testing.T) {
	unknownVals := []string{`"unknown"`, `"nope"`}

	for _, unknownName := range unknownVals {
		var v TimeGranularity

		err := v.UnmarshalJSON([]byte(unknownName))
		assert.ErrorContains(t, err, "cannot unmarshal unknown")
	}
}

func TestTimeGranularityUnmarshalJSONNil(t *testing.T) {
	var v *TimeGranularity

	err := v.UnmarshalJSON([]byte(`"hours"`))
	assert.ErrorContains(t, err, "cannot unmarshal into nil TimeGranularity")
}

func TestTimeGranularityUnmarshalJSONInvalid(t *testing.T) {
	invalidNames := []string{``, `"`, `x`, `xx`, `"x`, `x"`, `'something'`}

	for _, invalidName := range invalidNames {
		var v TimeGranularity

		err := v.UnmarshalJSON([]byte(invalidName))
		assert.ErrorContains(t, err, "cannot unmarshal invalid JSON")
	}
}

func TestTimeGranularityUnmarshalForm(t *testing.T) {
	vals := map[string]TimeGranularity{
		strMinutes: Minutes,
		strHours:   Hours,
	}
	assertHasAllValues(t, vals)

	for data, expected := range vals {
		var v TimeGranularity

		err := v.UnmarshalForm(data)
		assert.Nil(t, err)
		assert.Equal(t, v, expected)
	}
}

func TestTimeGranularityUnmarshalFormUnknown(t *testing.T) {
	unknownNames := []string{`unknown`, `nope`}

	for _, unknownName := range unknownNames {
		var v TimeGranularity

		err := v.UnmarshalForm(unknownName)
		assert.ErrorContains(t, err, "cannot unmarshal unknown TimeGranularity")
	}
}

func TestTimeGranularityUnmarshalFormNil(t *testing.T) {
	var v *TimeGranularity

	err := v.UnmarshalForm(`hours`)
	assert.ErrorContains(t, err, "cannot unmarshal into nil TimeGranularity")
}

func TestTimeGranularityRoundTripStruct(t *testing.T) {
	expected := testStruct{Hours}

	bytes, err := json.Marshal(&expected)
	assert.Nil(t, err)
	assert.NonNil(t, bytes)
	assert.Equal(
		t,
		string(bytes),
		`{"time_granularity":"hours"}`,
	)

	var ts testStruct
	err = json.Unmarshal(bytes, &ts)
	assert.Nil(t, err)
	assert.Equal(t, ts, expected)
}
