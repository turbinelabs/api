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
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"strconv"
)

// MatchKind is an Enumeration of the attributes by which a request can be
// matched.
type MatchKind string

// MatchBehavior is an Enumeration of possible ways to match a request attribute.
type MatchBehavior string

const (
	// CookiMatchKind matches against a request's cookies
	CookieMatchKind MatchKind = "cookie"
	// HeaderMatchKind matches against a request's headers
	HeaderMatchKind MatchKind = "header"
	// QueryMatchKind matches against a requests's query parameters
	QueryMatchKind MatchKind = "query"

	// ExactMatchBehavior matches a request attribute with an exact comparison.
	ExactMatchBehavior MatchBehavior = "exact"
	// RegexMatchBehavior matches a request attribute as a regex.
	RegexMatchBehavior MatchBehavior = "regex"
	// RangeBehaviuorKind matches a request attribute as a numeric range.
	RangeMatchBehavior MatchBehavior = "range"
	// PrefixBehaviorkind matches a request attribute as a prefix.
	PrefixMatchBehavior MatchBehavior = "prefix"
	// SuffixBehaviorkind matches a request attribute as a suffix.
	SuffixMatchBehavior MatchBehavior = "suffix"
)

/*
	A Match represents a mapping of a Metadatum from a MatchKind-typed
	request parameter to another Metadatum, with the MatchBehavior dictating
	how the values of the request parameter should be matched.

	Example:

		Match{
			HeaderMatchKind,
			ExactMatchBehavior,
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
			ExactMatchBehavior,
			Metadatum{Key:"X-GitSha"},
			Metadatum{Key:"git-sha"},
		}

	would define a match which looks for the X-GitSha header in a request,
	and if present, adds the value of that header as the value for a "git-sha"
	metadata constraint on Instances in the Cluster defined by the Route.
*/
type Match struct {
	Kind     MatchKind     `json:"kind"`
	Behavior MatchBehavior `json:"behavior"`
	From     Metadatum     `json:"from"`
	To       Metadatum     `json:"to"`
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
			fmt.Sprintf("%q is not a valid match kind", m.Kind)))
	}

	validBehavior := m.Behavior == ExactMatchBehavior ||
		m.Behavior == RegexMatchBehavior ||
		m.Behavior == RangeMatchBehavior ||
		m.Behavior == PrefixMatchBehavior ||
		m.Behavior == SuffixMatchBehavior

	if !validBehavior {
		errs.AddNew(ecase(
			"behavior",
			fmt.Sprintf("%q is not a valid behavior kind", m.Behavior)))
	}

	errCheckIndex(m.From.Key, errs, "from.key")

	if m.To.Value != "" && m.To.Key == "" {
		errs.AddNew(ecase("to.key", "must not be empty if to.value is set"))
	}

	// The only time it's ok to not have a specific matched value is with
	// exact behavior kind, to indicate that all values should be matched.
	if validBehavior && m.From.Value == "" && m.Behavior != ExactMatchBehavior {
		errs.AddNew(
			ecase(
				"from.value",
				fmt.Sprintf("must not be empty if behavior is %q", m.Behavior),
			),
		)
	}

	if m.Kind == QueryMatchKind && m.Behavior == RegexMatchBehavior {
		errs.AddNew(
			ecase(
				"behavior",
				fmt.Sprintf(
					`%q kind not supported with %q behavior`,
					QueryMatchKind,
					RegexMatchBehavior,
				),
			),
		)
	}

	if m.Kind == CookieMatchKind && m.Behavior == RangeMatchBehavior {
		errs.AddNew(
			ecase(
				"kind",
				fmt.Sprintf(
					`%q kind not supported with %q behavior`,
					CookieMatchKind,
					RegexMatchBehavior,
				),
			),
		)
	}

	if m.Behavior == RegexMatchBehavior && m.From.Value != "" {
		if re, e := regexp.Compile(m.From.Value); e != nil {
			errs.AddNew(
				ecase(
					"from.value",
					e.Error(),
				),
			)
		} else if m.To.Value == "" && len(re.SubexpNames()) > 2 {
			errs.AddNew(
				ecase(
					"from.value",
					"must have exactly one subgroup when to.value is not set",
				),
			)
		}
	}

	if m.Behavior == RangeMatchBehavior && m.From.Value != "" {
		_, _, err := ParseRangeBoundaries(m.From.Value)
		if err != nil {
			errs.AddNew(
				ecase(
					"from.value",
					err.Error(),
				),
			)
		}
	}

	return errs.OrNil()
}

func (m Match) Key() string {
	return fmt.Sprintf("%s:%s:%s", string(m.Kind), string(m.Behavior), m.From.Key)
}

// privateMatch is a private type redeclaration that allows us to make use
// of golang's default json unmarshalling while also implementing the
// Unmarshaler interface.
type privateMatch Match

// UnmarshalJSON implements the Unmarshaler method for Match. It defaults an
// empty Behavior field to `ExactBehaviorkind`.
func (m *Match) UnmarshalJSON(bytes []byte) error {
	if m == nil {
		return fmt.Errorf("cannot unmarshal into nil Match")
	}

	if err := json.Unmarshal(bytes, (*privateMatch)(m)); err != nil {
		return err
	}

	if m.Behavior == "" {
		m.Behavior = ExactMatchBehavior
	}

	return nil
}

// Matches is a slice of Match objects
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

// Checks for equality between two match objects. For two Match objects to be
// considered equal they must share a Kind, From, and To.
func (m Match) Equals(o Match) bool {
	return m.Kind == o.Kind &&
		m.Behavior == o.Behavior &&
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

// ParseRangeBoundaries splits out the boundaries contained in a range string.
// Returns any errors along the way.
func ParseRangeBoundaries(s string) (int, int, error) {
	ms := RangeMatchPattern.FindStringSubmatch(s)
	if len(ms) != 3 {
		return 0, 0,
			fmt.Errorf(
				`Invalid range pattern: %q. Must be of the form "[<start>,<end>)"`,
				s,
			)
	}

	start, err := strconv.Atoi(ms[1])
	if err != nil {
		return 0, 0, err
	}

	end, err := strconv.Atoi(ms[2])
	if err != nil {
		return 0, 0, err
	}

	if end <= start {
		return 0, 0, errors.New("End of range must be greater than start")
	}

	return start, end, nil
}
