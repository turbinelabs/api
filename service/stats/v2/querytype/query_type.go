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

// This file was automatically generated by doc.go from ../../enum.template.
// Any changes will be lost if this file is regenerated.

package querytype

import (
	"encoding/json"
	"fmt"
)

// QueryType is an enumeration.
type QueryType int

// Defined values of QueryType.
const (
	Unknown QueryType = iota
	Requests
	Responses
	Success
	Error
	Failure
	LatencyP50
	LatencyP99
	SuccessRate
	ResponsesForCode
	DownstreamRequests
	DownstreamResponses
	DownstreamSuccess
	DownstreamError
	DownstreamFailure
	DownstreamLatencyP50
	DownstreamLatencyP99
	DownstreamSuccessRate
	DownstreamResponsesForCode
)

var _dummy = QueryType(0)
var _ json.Marshaler = &_dummy
var _ json.Unmarshaler = &_dummy

const (
	strUnknown                    = "unknown"
	strRequests                   = "requests"
	strResponses                  = "responses"
	strSuccess                    = "success"
	strError                      = "error"
	strFailure                    = "failure"
	strLatencyP50                 = "latency_p50"
	strLatencyP99                 = "latency_p99"
	strSuccessRate                = "success_rate"
	strResponsesForCode           = "responses_for_code"
	strDownstreamRequests         = "downstream_requests"
	strDownstreamResponses        = "downstream_responses"
	strDownstreamSuccess          = "downstream_success"
	strDownstreamError            = "downstream_error"
	strDownstreamFailure          = "downstream_failure"
	strDownstreamLatencyP50       = "downstream_latency_p50"
	strDownstreamLatencyP99       = "downstream_latency_p99"
	strDownstreamSuccessRate      = "downstream_success_rate"
	strDownstreamResponsesForCode = "downstream_responses_for_code"
)

var queryTypeNames = [...]string{
	strUnknown,
	strRequests,
	strResponses,
	strSuccess,
	strError,
	strFailure,
	strLatencyP50,
	strLatencyP99,
	strSuccessRate,
	strResponsesForCode,
	strDownstreamRequests,
	strDownstreamResponses,
	strDownstreamSuccess,
	strDownstreamError,
	strDownstreamFailure,
	strDownstreamLatencyP50,
	strDownstreamLatencyP99,
	strDownstreamSuccessRate,
	strDownstreamResponsesForCode,
}

const minQueryType = QueryType(1)

var maxQueryType = QueryType(len(queryTypeNames) - 1)

// IsValid returns a boolean indicating whether the given QueryType is defined
// (valid) or not.
func IsValid(i QueryType) bool {
	return i >= minQueryType && i <= maxQueryType
}

// FromName converts a QueryType name into a QueryType.
// Returns Unknown if the name is not known. Names are case sensitive.
func FromName(s string) QueryType {
	for idx, name := range queryTypeNames {
		candidate := QueryType(idx)
		if IsValid(candidate) && name == s {
			return candidate
		}
	}

	return Unknown
}

// ForEach invokes the given function for each valid QueryType.
func ForEach(f func(QueryType)) {
	for i := int(minQueryType); i <= int(maxQueryType); i++ {
		tg := QueryType(i)
		f(tg)
	}
}

// String return this QueryType's string representation.
func (i QueryType) String() string {
	if !IsValid(i) {
		return fmt.Sprintf("unknown(%d)", i)
	}
	return queryTypeNames[i]
}

// MarshalJSON converts this QueryType to a quoted JSON string. Returns an
// error if the QueryType is nil or invalid.
func (i *QueryType) MarshalJSON() ([]byte, error) {
	if i == nil {
		return nil, fmt.Errorf("cannot marshal unknown QueryType (nil)")
	}

	qt := *i
	if !IsValid(qt) {
		return nil, fmt.Errorf("cannot marshal unknown QueryType (%d)", qt)
	}

	name := queryTypeNames[qt]
	b := make([]byte, 0, len(name)+2)
	b = append(b, '"')
	b = append(b, name...)
	return append(b, '"'), nil
}

// UnmarshalJSON converts a quoted JSON string into a QueryType. Returns an
// error if the receiver is nil, the JSON is not a quoted string, or if the
// string does not represent a valid QueryType. Otherwise, the receiver's value
// is set to the QueryType represented by the string.
func (i *QueryType) UnmarshalJSON(bytes []byte) error {
	if i == nil {
		return fmt.Errorf("cannot unmarshal into nil QueryType")
	}

	length := len(bytes)
	if length <= 2 || bytes[0] != '"' || bytes[length-1] != '"' {
		return fmt.Errorf("cannot unmarshal invalid JSON: %q", string(bytes))
	}

	unmarshalName := string(bytes[1 : length-1])

	qt := FromName(unmarshalName)
	if qt == Unknown {
		return fmt.Errorf(
			"cannot unmarshal unknown QueryType %q",
			unmarshalName,
		)
	}

	*i = qt
	return nil
}

// UnmarshalForm converts a string into a QueryType. Returns an error if the
// receiver is nil or if the string does not represent a valid QueryType.
// Otherwise, the receiver's value is set to the QueryType represented by the
// string.
func (i *QueryType) UnmarshalForm(value string) error {
	if i == nil {
		return fmt.Errorf("cannot unmarshal into nil QueryType")
	}

	qt := FromName(value)
	if qt == Unknown {
		return fmt.Errorf("cannot unmarshal unknown QueryType %q", value)
	}

	*i = qt
	return nil
}