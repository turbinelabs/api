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

type RuleKey string

/*
	A Rule defines a mapping from a list of Methods and Matches to an AllConstraints
	struct. A Rule applies to a request if one of the Methods and all of the Matches
	apply.

	If a Rule applies, the constraints inferred from the Matches should be merged
	with each of the ClusterConstraints, which are then used to find a live
	Instance. The ClusterConstraints are randomly shuffled using their weights
	to affect the distribution. Each ClusterConstraint is examined to find a
	matching Instance, until one is found.

	TODO: Another approach to consider is a default ClusterConstraint, and
	then a ranked list of ClusterConstraint, each of which is attempted based
	on an assigned probability. EG, "1% of the time, do the canary. 1% of the
	remaining time (or if there is no canary), try this, etc, using the default
	as a fallback." The downsides are that the distribution gets less and less
	accurate as you go down the list. The upside is that you don't have to shuffle
	the list.
*/
type Rule struct {
	RuleKey     RuleKey        `json:"rule_key"`
	Methods     []string       `json:"methods"`
	Matches     Matches        `json:"matches"`
	Constraints AllConstraints `json:"constraints"`
}

type Rules []Rule

// Checks for equality with another Rule slice. Slices will be equal if each
// element i is Equal to ith element of the other slice and the slices are
// of the same length.
func (r Rules) Equals(o Rules) bool {
	if len(r) != len(o) {
		return false
	}

	for i, e := range r {
		if !e.Equals(o[i]) {
			return false
		}
	}

	return true
}

func (rs Rules) AsMap() map[RuleKey]Rule {
	m := map[RuleKey]Rule{}
	for _, r := range rs {
		m[r.RuleKey] = r
	}
	return m
}

// Check for validity of a slice of Rule objects. A valid rule is one that is
// composed only of valid Rule structs.
func (r Rules) IsValid() *ValidationError {
	errs := &ValidationError{}

	seenKey := map[RuleKey]bool{}
	for _, r := range r {
		if seenKey[r.RuleKey] {
			errs.AddNew(ErrorCase{
				"rules", fmt.Sprintf("multiple instances of key %s", string(r.RuleKey)),
			})
		}
		seenKey[r.RuleKey] = true

		errs.MergePrefixed(r.IsValid(), fmt.Sprintf("rules[%v]", r.RuleKey))
	}

	return errs.OrNil()
}
func (r Rule) methodsEqual(o Rule) bool {
	if len(r.Methods) != len(o.Methods) {
		return false
	}

	m := make(map[string]bool)
	for _, e := range r.Methods {
		m[e] = true
	}

	for _, e := range o.Methods {
		if !m[e] {
			return false
		}
	}

	return true
}

// Checks for equality between two Rules. Rules are equal if the rule key,
// methods, constraints, and matches are all equal.
func (r Rule) Equals(o Rule) bool {
	if r.RuleKey != o.RuleKey {
		return false
	}

	if !r.methodsEqual(o) {
		return false
	}

	if !r.Constraints.Equals(o.Constraints) {
		return false
	}

	if !r.Matches.Equals(o.Matches) {
		return false
	}

	return true
}

var validMethod map[string]bool = map[string]bool{
	"GET":    true,
	"PUT":    true,
	"POST":   true,
	"DELETE": true,
}

// Checks this rule for validity. A rule is considered valid if it has a RuleKey,
// at least one valid HTTP method (GET, PUT, POST, DELETE), the defined
// matches are valid, and the Constraints are valid.
func (r Rule) IsValid() *ValidationError {
	ecase := func(f, m string) ErrorCase {
		return ErrorCase{f, m}
	}

	errs := &ValidationError{}

	errCheckKey(string(r.RuleKey), errs, "rule_key")

	for _, m := range r.Methods {
		if !validMethod[m] {
			errs.AddNew(ecase(
				"methods",
				fmt.Sprintf("%s is not a valid method", m),
			))
		}
	}

	if len(r.Methods) == 0 && len(r.Matches) == 0 {
		errs.AddNew(ecase("", "at least one method or match must be present"))
	}

	errs.MergePrefixed(r.Matches.IsValid(), "")
	errs.MergePrefixed(r.Constraints.IsValid("constraints"), "")

	return errs.OrNil()
}
