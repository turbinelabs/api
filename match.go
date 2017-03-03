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
)

// MatchKind is an Enumeration of possible ways to match a request
type MatchKind string

const (
	CookieMatchKind MatchKind = "cookie" // matches a cookie
	HeaderMatchKind           = "header" // matches a header
	QueryMatchKind            = "query"  // matches a query variable
)

func MatchKindFromString(s string) (MatchKind, error) {
	m := MatchKind(s)
	switch m {
	case CookieMatchKind, HeaderMatchKind, QueryMatchKind:
		return m, nil
	}
	return MatchKind(""), fmt.Errorf("unknown MatchKind: %s", s)
}

/*
	A Match is represents a mapping of a Metadatum from a MatchKind-typed
	request parameter to another Metadatum.

	Example:

		Match{
			HeaderMatchKind,
			Metadatum{"X-SwVersion", "1.0"},
			Metadatum{"sunset", "true"},
		}

	would define a match which looks for a specific value of "1.0" for the
	X-SwVersion header in a request, and if present, adds a specific key/value
	constraint of sunset=true.

	Values can be omitted to map any value for the specified key.

	Example:

		Match{
			HeaderMatchKind,
			Metadatum{Key:"X-GitSha"},
			Metadatum{Key:"git-sha"},
		}

	would define a match which looks for the X-GitSha header in a request,
	and if present, adds the value of that header as the value for a "git-sha"
	metadata constraint on Instances in the Cluster defined by the Route.
*/
type Match struct {
	Kind MatchKind `json:"kind"`
	From Metadatum `json:"from"`
	To   Metadatum `json:"to"`
}

// Check this Match for validity. A valid match requires a valid matchkind,
// a From datum with Key set, and either an empty To datum or one with both
// Key and Value set.
func (m Match) IsValid() *ValidationError {
	ecase := func(f, msg string) ErrorCase {
		return ErrorCase{f, msg}
	}

	errs := &ValidationError{}

	validKind := m.Kind == CookieMatchKind ||
		m.Kind == HeaderMatchKind ||
		m.Kind == QueryMatchKind

	if !validKind {
		errs.AddNew(ecase(
			"kind",
			fmt.Sprintf("%s is not a valid match kind", string(m.Kind))))
	}

	errCheckIndex(m.From.Key, errs, "from.key")

	if m.To.Value != "" && m.To.Key == "" {
		errs.AddNew(ecase("to.key", "must not be empty if value is set"))
	}

	return errs.OrNil()
}

func (m Match) Key() string {
	return fmt.Sprintf("%s:%s", string(m.Kind), m.From.Key)
}

type Matches []Match

// Checks validity of a slice of Match objects. For the slice to be valid each
// entry must be valid.
func (m Matches) IsValid() *ValidationError {
	errs := &ValidationError{}

	seenMatch := map[string]bool{}
	for _, e := range m {
		if seenMatch[e.Key()] {
			errs.AddNew(ErrorCase{"", "duplicate match found " + e.Key()})
		}
		errs.MergePrefixed(e.IsValid(), fmt.Sprintf("matches[%v]", e.Key()))
		seenMatch[e.Key()] = true
	}

	return errs.OrNil()
}

// Checks for equalit between two match objects. For two Match objects to be
// considered equal they must share a Kind, From, and To.
func (m Match) Equals(o Match) bool {
	return m.Kind == o.Kind &&
		m.From.Equals(o.From) &&
		m.To.Equals(o.To)
}

// Checks a slice of Match objects for equality against another. This comparison
// is order agnostic.
func (m Matches) Equals(o Matches) bool {
	if len(m) != len(o) {
		return false
	}

	hasMatch := make(map[Match]bool)
	for _, e := range m {
		hasMatch[e] = true
	}

	for _, e := range o {
		if !hasMatch[e] {
			return false
		}
	}

	return true
}
