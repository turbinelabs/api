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
func (m Match) IsValid(precreation bool) *ValidationError {
	ecase := func(f, msg string) ErrorCase {
		return ErrorCase{fmt.Sprintf("match[%s].%s", m.Key(), f), msg}
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

	if m.From.Key == "" {
		errs.AddNew(ecase("from.key", "must not be empty"))
	}

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
func (m Matches) IsValid(precreation bool) *ValidationError {
	errs := &ValidationError{}

	for _, e := range m {
		errs.MergePrefixed(e.IsValid(precreation), "")
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
