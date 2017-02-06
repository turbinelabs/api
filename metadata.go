package api

import (
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

func metadataCommonEquals(m1, m2 Metadata) bool {
	if len(m1) != len(m2) {
		return false
	}

	m1Map := m1.Map()

	for _, v := range m2 {
		if m2V, m2OK := m1Map[v.Key]; !m2OK || m2V != v.Value {
			return false
		}
	}

	return true
}

func (m Metadata) Equals(o Metadata) bool {
	return metadataCommonEquals(m, o)
}

func (m Metadata) Equivalent(o Metadata) bool {
	return metadataCommonEquals(m, o)
}

// A Metadatum a key/value pair.
type Metadatum struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

// MetadataValid provides a way to check for validity of all Metadatum contained
// within a Metadata. It returns an aggregated list of errors, or nil if none
// were found
func MetadataValid(md Metadata, check func(Metadatum) *ValidationError) *ValidationError {
	errs := &ValidationError{}

	for _, e := range md {
		errs.Merge(check(e))
	}

	return errs.OrNil()
}

func datumCommonEquals(m1, m2 Metadatum) bool {
	return m1.Key == m2.Key && m1.Value == m2.Value
}

// Checks for equality between two Metadatum structs. They will be considered
// equal if the key and value both match.
func (m Metadatum) Equals(o Metadatum) bool {
	return datumCommonEquals(m, o)
}

// Checks for semantic equality between two Metadatum structs. They will be
// considered equal if the key and value both match.
func (m Metadatum) Equivalent(o Metadatum) bool {
	return datumCommonEquals(m, o)
}

// Sort a Metadata by Metadatum Key.
// Eg: sort.Sort(MetadataByKey(metadata))
type MetadataByKey Metadata

var _ sort.Interface = MetadataByKey{}

func (b MetadataByKey) Len() int           { return len(b) }
func (b MetadataByKey) Swap(i, j int)      { b[i], b[j] = b[j], b[i] }
func (b MetadataByKey) Less(i, j int) bool { return b[i].Key < b[j].Key }
