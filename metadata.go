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
	"regexp"
	"sort"
)

// Metadata is a vector of Metadatums.
type Metadata []Metadatum

// MetadataFromMap converts a map[string]string into a Metadata
func MetadataFromMap(m map[string]string) Metadata {
	meta := Metadata{}
	for k, v := range m {
		meta = append(meta, Metadatum{k, v})
	}

	return meta
}

// Map produces a string/string Map from Metadata.
func (m Metadata) Map() map[string]string {
	result := map[string]string{}
	for _, metadatum := range m {
		result[metadatum.Key] = metadatum.Value
	}

	return result
}

// Equals checks equality with another Metadata
func (m Metadata) Equals(o Metadata) bool {
	if len(m) != len(o) {
		return false
	}

	mMap := m.Map()
	oMap := o.Map()

	for k, v := range mMap {
		if ov, ok := oMap[k]; !ok || ov != v {
			return false
		}
	}

	return true
}

// Compare compares the receiver to another Metadata.
// It returns a value > 0 if the receiver is greater,
// < 0 if the receiver is lessor, and 0 if they are equal.
// Both receiver and target are sorted by key as a side
// effect.
func (m Metadata) Compare(o Metadata) int {
	sort.Sort(MetadataByKey(m))
	sort.Sort(MetadataByKey(o))

	if len(m) > len(o) {
		return 1
	}
	if len(m) < len(o) {
		return -1
	}
	for idx := range m {
		if m[idx].Key > o[idx].Key {
			return 1
		}
		if m[idx].Key < o[idx].Key {
			return -1
		}

		if m[idx].Value > o[idx].Value {
			return 1
		}
		if m[idx].Value < o[idx].Value {
			return -1
		}
	}
	return 0
}

// A Metadatum a key/value pair.
type Metadatum struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type MetadataCheck func(Metadatum) *ValidationError

// MetadataValid provides a way to check for validity of all Metadatum contained
// within a Metadata, for various definitions of validity. It returns an
// aggregated list of errors from all checks, or nil if none were found.
func MetadataValid(
	container string,
	md Metadata,
	checks ...MetadataCheck,
) *ValidationError {
	errs := &ValidationError{}

	seenKey := map[string]bool{}
	for _, e := range md {
		if seenKey[e.Key] {
			errs.AddNew(ErrorCase{
				container,
				fmt.Sprintf("duplicate %v key '%v'", container, e.Key),
			})
		}
		seenKey[e.Key] = true

		for _, check := range checks {
			errs.MergePrefixed(check(e), fmt.Sprintf("%s[%v]", container, e.Key))
		}
	}

	return errs.OrNil()
}

var (
	// MetadataCheckNonEmptyKeys produces an error if the Metadatum has an empty
	// Key
	MetadataCheckNonEmptyKeys = MetadataCheck(func(kv Metadatum) *ValidationError {
		if kv.Key == "" {
			return &ValidationError{
				[]ErrorCase{{"key", "must not be empty"}},
			}
		}
		return nil
	})

	// MetadataCheckNonEmptyValues produces an error if the Metadatum has an empty
	// Value
	MetadataCheckNonEmptyValues = MetadataCheck(func(kv Metadatum) *ValidationError {
		if kv.Value == "" {
			return &ValidationError{
				[]ErrorCase{{"value", "must not be empty"}},
			}
		}
		return nil
	})
)

// MetadataCheckKeysMatchPattern produces an error with the given error
// string, if the Metadatum fails to match the given pattern.
func MetadataCheckKeysMatchPattern(pattern *regexp.Regexp, errStr string) MetadataCheck {
	return func(kv Metadatum) *ValidationError {
		if !pattern.MatchString(kv.Key) {
			return &ValidationError{
				[]ErrorCase{{"key", errStr}},
			}
		}
		return nil
	}
}

// Equals checks for equality between two Metadatum structs. They will be
// considered equal if the key and value both match.
func (m Metadatum) Equals(o Metadatum) bool {
	return m.Key == o.Key && m.Value == o.Value
}

// MetadataByKey implements sort.Interface to allow sorting of Metadata by
// Metadatum Key. Eg: sort.Sort(MetadataByKey(metadata))
type MetadataByKey Metadata

var _ sort.Interface = MetadataByKey{}

func (b MetadataByKey) Len() int           { return len(b) }
func (b MetadataByKey) Swap(i, j int)      { b[i], b[j] = b[j], b[i] }
func (b MetadataByKey) Less(i, j int) bool { return b[i].Key < b[j].Key }
