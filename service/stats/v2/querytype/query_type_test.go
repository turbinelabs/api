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

// This file was automatically generated by doc.go from github.com/turbinelabs/api/service/stats/enum_test.template.
// Any changes will be lost if this file is regenerated.

package querytype

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/turbinelabs/test/assert"
)

type testStruct struct {
	QueryType QueryType `json:"query_type"`
}

func assertHasAllValues(t testing.TB, i interface{}) {
	mappedValues := map[QueryType]struct{}{}

	switch m := i.(type) {
	case map[QueryType]string:
		for qt := range m {
			mappedValues[qt] = struct{}{}
		}
	case map[string]QueryType:
		for _, qt := range m {
			mappedValues[qt] = struct{}{}
		}
	default:
		t.Fatalf("cannot validate type %T", i)
		return
	}

	ForEach(func(v QueryType) {
		assert.Group(
			fmt.Sprintf("QueryType %s", v.String()),
			t,
			func(g *assert.G) {
				_, ok := mappedValues[v]
				assert.True(g, ok)
			},
		)
	})
}

func TestQueryTypeString(t *testing.T) {
	assert.Equal(
		t,
		Requests.String(),
		"requests",
	)
	assert.MatchesRegex(t, Unknown.String(), `^unknown\([0-9]+\)$`)
	assert.Equal(t, QueryType(100).String(), "unknown(100)")
}

func TestIsValid(t *testing.T) {
	invalid := []QueryType{
		QueryType(minQueryType - 1),
		QueryType(maxQueryType + 1),
	}

	for _, qt := range invalid {
		assert.False(t, IsValid(qt))
	}

	ForEach(func(qt QueryType) {
		assert.True(t, IsValid(qt))
	})
}

func TestFromName(t *testing.T) {
	validValues := map[QueryType]string{
		Requests:                   strRequests,
		Responses:                  strResponses,
		Success:                    strSuccess,
		Error:                      strError,
		Failure:                    strFailure,
		LatencyP50:                 strLatencyP50,
		LatencyP99:                 strLatencyP99,
		SuccessRate:                strSuccessRate,
		ResponsesForCode:           strResponsesForCode,
		DownstreamRequests:         strDownstreamRequests,
		DownstreamResponses:        strDownstreamResponses,
		DownstreamSuccess:          strDownstreamSuccess,
		DownstreamError:            strDownstreamError,
		DownstreamFailure:          strDownstreamFailure,
		DownstreamLatencyP50:       strDownstreamLatencyP50,
		DownstreamLatencyP99:       strDownstreamLatencyP99,
		DownstreamSuccessRate:      strDownstreamSuccessRate,
		DownstreamResponsesForCode: strDownstreamResponsesForCode,
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

func TestQueryTypeMarshalJSON(t *testing.T) {
	vals := map[QueryType]string{
		Requests:                   strRequests,
		Responses:                  strResponses,
		Success:                    strSuccess,
		Error:                      strError,
		Failure:                    strFailure,
		LatencyP50:                 strLatencyP50,
		LatencyP99:                 strLatencyP99,
		SuccessRate:                strSuccessRate,
		ResponsesForCode:           strResponsesForCode,
		DownstreamRequests:         strDownstreamRequests,
		DownstreamResponses:        strDownstreamResponses,
		DownstreamSuccess:          strDownstreamSuccess,
		DownstreamError:            strDownstreamError,
		DownstreamFailure:          strDownstreamFailure,
		DownstreamLatencyP50:       strDownstreamLatencyP50,
		DownstreamLatencyP99:       strDownstreamLatencyP99,
		DownstreamSuccessRate:      strDownstreamSuccessRate,
		DownstreamResponsesForCode: strDownstreamResponsesForCode,
	}
	assertHasAllValues(t, vals)

	for v, name := range vals {
		bytes, err := v.MarshalJSON()
		assert.Nil(t, err)
		expected := []byte(fmt.Sprintf(`"%s"`, name))
		assert.DeepEqual(t, bytes, expected)
	}
}

func TestQueryTypeMarshalJSONUnknown(t *testing.T) {
	unknownValues := []QueryType{
		Unknown,
		QueryType(maxQueryType + 1),
	}

	for _, unknownQueryType := range unknownValues {
		bytes, err := unknownQueryType.MarshalJSON()
		assert.Nil(t, bytes)
		assert.ErrorContains(t, err, "cannot marshal unknown QueryType")
	}
}

func TestQueryTypeMarshalJSONNil(t *testing.T) {
	var v *QueryType

	bytes, err := v.MarshalJSON()
	assert.ErrorContains(t, err, "cannot marshal unknown")
	assert.Nil(t, bytes)
}

func TestQueryTypeUnmarshalJSON(t *testing.T) {
	quoted := func(s string) string {
		return fmt.Sprintf(`"%s"`, s)
	}

	vals := map[string]QueryType{
		quoted(strRequests):                   Requests,
		quoted(strResponses):                  Responses,
		quoted(strSuccess):                    Success,
		quoted(strError):                      Error,
		quoted(strFailure):                    Failure,
		quoted(strLatencyP50):                 LatencyP50,
		quoted(strLatencyP99):                 LatencyP99,
		quoted(strSuccessRate):                SuccessRate,
		quoted(strResponsesForCode):           ResponsesForCode,
		quoted(strDownstreamRequests):         DownstreamRequests,
		quoted(strDownstreamResponses):        DownstreamResponses,
		quoted(strDownstreamSuccess):          DownstreamSuccess,
		quoted(strDownstreamError):            DownstreamError,
		quoted(strDownstreamFailure):          DownstreamFailure,
		quoted(strDownstreamLatencyP50):       DownstreamLatencyP50,
		quoted(strDownstreamLatencyP99):       DownstreamLatencyP99,
		quoted(strDownstreamSuccessRate):      DownstreamSuccessRate,
		quoted(strDownstreamResponsesForCode): DownstreamResponsesForCode,
	}
	assertHasAllValues(t, vals)

	for data, expected := range vals {
		var v QueryType

		err := v.UnmarshalJSON([]byte(data))
		assert.Nil(t, err)
		assert.Equal(t, v, expected)
	}
}

func TestQueryTypeUnmarshalJSONUnknown(t *testing.T) {
	unknownVals := []string{`"unknown"`, `"nope"`}

	for _, unknownName := range unknownVals {
		var v QueryType

		err := v.UnmarshalJSON([]byte(unknownName))
		assert.ErrorContains(t, err, "cannot unmarshal unknown")
	}
}

func TestQueryTypeUnmarshalJSONNil(t *testing.T) {
	var v *QueryType

	err := v.UnmarshalJSON([]byte(`"requests"`))
	assert.ErrorContains(t, err, "cannot unmarshal into nil QueryType")
}

func TestQueryTypeUnmarshalJSONInvalid(t *testing.T) {
	invalidNames := []string{``, `"`, `x`, `xx`, `"x`, `x"`, `'something'`}

	for _, invalidName := range invalidNames {
		var v QueryType

		err := v.UnmarshalJSON([]byte(invalidName))
		assert.ErrorContains(t, err, "cannot unmarshal invalid JSON")
	}
}

func TestQueryTypeUnmarshalForm(t *testing.T) {
	vals := map[string]QueryType{
		strRequests:                   Requests,
		strResponses:                  Responses,
		strSuccess:                    Success,
		strError:                      Error,
		strFailure:                    Failure,
		strLatencyP50:                 LatencyP50,
		strLatencyP99:                 LatencyP99,
		strSuccessRate:                SuccessRate,
		strResponsesForCode:           ResponsesForCode,
		strDownstreamRequests:         DownstreamRequests,
		strDownstreamResponses:        DownstreamResponses,
		strDownstreamSuccess:          DownstreamSuccess,
		strDownstreamError:            DownstreamError,
		strDownstreamFailure:          DownstreamFailure,
		strDownstreamLatencyP50:       DownstreamLatencyP50,
		strDownstreamLatencyP99:       DownstreamLatencyP99,
		strDownstreamSuccessRate:      DownstreamSuccessRate,
		strDownstreamResponsesForCode: DownstreamResponsesForCode,
	}
	assertHasAllValues(t, vals)

	for data, expected := range vals {
		var v QueryType

		err := v.UnmarshalForm(data)
		assert.Nil(t, err)
		assert.Equal(t, v, expected)
	}
}

func TestQueryTypeUnmarshalFormUnknown(t *testing.T) {
	unknownNames := []string{`unknown`, `nope`}

	for _, unknownName := range unknownNames {
		var v QueryType

		err := v.UnmarshalForm(unknownName)
		assert.ErrorContains(t, err, "cannot unmarshal unknown QueryType")
	}
}

func TestQueryTypeUnmarshalFormNil(t *testing.T) {
	var v *QueryType

	err := v.UnmarshalForm(`requests`)
	assert.ErrorContains(t, err, "cannot unmarshal into nil QueryType")
}

func TestQueryTypeRoundTripStruct(t *testing.T) {
	expected := testStruct{Requests}

	bytes, err := json.Marshal(&expected)
	assert.Nil(t, err)
	assert.NonNil(t, bytes)
	assert.Equal(
		t,
		string(bytes),
		`{"query_type":"requests"}`,
	)

	var ts testStruct
	err = json.Unmarshal(bytes, &ts)
	assert.Nil(t, err)
	assert.Equal(t, ts, expected)
}
