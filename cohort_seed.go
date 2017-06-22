package api

import (
	"fmt"
	"strings"
)

// CohortSeedType indicates how a CohortSeed sources its data.
type CohortSeedType string

const (
	// CohortSeedHeader specifies that seed data will be drawn from a request header.
	CohortSeedHeader CohortSeedType = "header"

	// CohortSeedCookie specifies that seed data will be taken from a request cookie.
	CohortSeedCookie CohortSeedType = "cookie"

	// CohortSeedQuery specifies that seed data will be taken from a path query parameter.
	CohortSeedQuery CohortSeedType = "query"
)

// CohortSeed identifies a request attribute that will be used to group a
// request into cohort. When a cohort member is mapped to a subset of backend
// instances it will be consistently sent to the same subset each time.
type CohortSeed struct {
	Type CohortSeedType `json:"type"` // Type indicates what kind of attribute Name references.
	Name string         `json:"name"` // Name specifies the source of cohort seed data.

	// UseZeroValueSeed controls whether or not missing / unset data is included
	// in a stable cohort.
	//
	// If false (default value): when a request contains no seed value, it will be
	// processed normally and handled by a randomly selected backend.
	//
	// If true: when a request is processed that has no value for the seed then we
	// will derive a seed from the empty string. This means that all requests
	// without a seed value will be serviced by the same backend. There are
	// problematic edge cases here (misspelling the seed source and having all
	// traffic land on a backend intended to take 1%, etc) so use with caution.
	UseZeroValueSeed bool `json:"use_zero_value_seed"`
}

func (cs CohortSeed) Equals(o CohortSeed) bool {
	return cs == o
}

func (c CohortSeed) IsValid() *ValidationError {
	errs := &ValidationError{}

	switch c.Type {
	case CohortSeedHeader, CohortSeedCookie, CohortSeedQuery:
	default:
		errs.AddNew(ErrorCase{"cohort_seed.type", fmt.Sprintf("%q is not a valid seed type", c.Type)})
	}

	if strings.TrimSpace(c.Name) == "" {
		errs.AddNew(ErrorCase{"cohort_seed.name", "may not be empty"})
	}

	return errs.OrNil()
}

func CohortSeedPtrEquals(cp1, cp2 *CohortSeed) bool {
	switch {
	case cp1 == nil && cp2 == nil:
		return true
	case cp1 == nil || cp2 == nil:
		return false
	default:
		return cp1.Equals(*cp2)
	}
}
