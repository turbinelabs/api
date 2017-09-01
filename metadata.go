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

// Produce a string/string Map from Metadata.
func (metadata *Metadata) Map() map[string]string {
	result := make(map[string]string)
	for _, metadatum := range *metadata {
		result[metadatum.Key] = metadatum.Value
	}

	return result
}

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

// A Metadatum a key/value pair.
type Metadatum struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

// MetadataValid provides a way to check for validity of all Metadatum contained
// within a Metadata. It returns an aggregated list of errors, or nil if none
// were found
func MetadataValid(
	container string,
	md Metadata,
	check func(Metadatum) *ValidationError,
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

		errs.MergePrefixed(check(e), fmt.Sprintf("%s[%v]", container, e.Key))
		seenKey[e.Key] = true
	}

	return errs.OrNil()
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
