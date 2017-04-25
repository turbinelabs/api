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

// Package timegranularity defines the TimeGranularity pseudo-enum
package timegranularity

import (
	"encoding/json"
	"fmt"
)

// TimeGranularity represents the granularity of stats query results.
type TimeGranularity int

const (
	// Seconds specifies per-second query results. Currently
	// not supported.
	Seconds TimeGranularity = iota

	// Minutes specifies per-minute query results.
	Minutes

	// Hours specifies per-minute query results.
	Hours

	// Unknown represents an unknown time granularity.
	Unknown
)

var _dummyGranularity = Seconds
var _ json.Marshaler = &_dummyGranularity
var _ json.Unmarshaler = &_dummyGranularity

const (
	seconds = "seconds"
	minutes = "minutes"
	hours   = "hours"
	unknown = "unknown"
)

var granularityNames = [...]string{
	seconds,
	minutes,
	hours,
}

var maxTimeGranularity = TimeGranularity(len(granularityNames) - 1)

// IsValid tests a TimeGranularity and reports whether it is valid
// value or not.
func IsValid(i TimeGranularity) bool {
	return i >= 0 && i <= maxTimeGranularity
}

// FromName converts a string into a TimeGranularity, returning
// Unknown if the string is not a valid TimeGranularity name.
func FromName(s string) TimeGranularity {
	for idx, name := range granularityNames {
		if name == s {
			return TimeGranularity(idx)
		}
	}

	return Unknown
}

// ForEach iterates over the value values and invokes f for each one.
func ForEach(f func(TimeGranularity)) {
	for i := 0; i <= int(maxTimeGranularity); i++ {
		tg := TimeGranularity(i)
		f(tg)
	}
}

// String converts the TimeGranularity into a string.
func (tg TimeGranularity) String() string {
	if !IsValid(tg) {
		return fmt.Sprintf("unknown(%d)", tg)
	}
	return granularityNames[tg]
}

// MarshalJSON converts the TimeGranularity into a quoted JSON string.
func (tg *TimeGranularity) MarshalJSON() ([]byte, error) {
	if tg == nil {
		return nil, fmt.Errorf("cannot marshal unknown time granularity (nil)")
	}

	timeGran := *tg
	if !IsValid(timeGran) {
		return nil, fmt.Errorf("cannot marshal unknown time granularity (%d)", timeGran)
	}

	name := granularityNames[timeGran]
	b := make([]byte, 0, len(name)+2)
	b = append(b, '"')
	b = append(b, name...)
	return append(b, '"'), nil
}

// UnmarshalJSON parses a quoted JSON string and updates the value of
// the target TimeGranularity.
func (tg *TimeGranularity) UnmarshalJSON(bytes []byte) error {
	if tg == nil {
		return fmt.Errorf("cannot unmarshal into nil TimeGranularity")
	}

	length := len(bytes)
	if length <= 2 || bytes[0] != '"' || bytes[length-1] != '"' {
		return fmt.Errorf("cannot unmarshal invalid JSON: `%s`", string(bytes))
	}

	unmarshalName := string(bytes[1 : length-1])
	timeGran := FromName(unmarshalName)
	if timeGran == Unknown {
		return fmt.Errorf("cannot unmarshal unknown time granularity `%s`", unmarshalName)
	}

	*tg = timeGran
	return nil
}

// UnmarshalForm parses a form value string and updates the value of
// the target TimeGranularity.
func (tg *TimeGranularity) UnmarshalForm(value string) error {
	if tg == nil {
		return fmt.Errorf("cannot unmarshal into nil TimeGranularity")
	}

	timeGran := FromName(value)
	if timeGran == Unknown {
		return fmt.Errorf("cannot unmarshal unknown time granularity `%s`", value)
	}

	*tg = timeGran
	return nil

}
