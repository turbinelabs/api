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
)

type RedirectType string

const (
	// PermanentRedirect will be handled by returned a HTTP response code 301
	// with the new URL that should be retrieved.
	PermanentRedirect RedirectType = "permanent"

	// TemporaryRedirect will be handled by returned a HTTP response code 302
	// with the new URL that should be retrieved.
	TemporaryRedirect RedirectType = "temporary"
)

var NamePattern = regexp.MustCompile("^[0-9a-zA-Z_-]+$")

// Redirects is a collection of Domain redirect definitions
type Redirects []Redirect

// AsMap converts an ordered slice of Redirects into a map where each redirect
// is indexed by its name.
func (rs Redirects) AsMap() map[string]Redirect {
	m := map[string]Redirect{}
	for _, r := range rs {
		m[r.Name] = r
	}
	return m
}

// Keys returns a list of the Name attributes from a Redirects; order is the same
// as the slice.
func (rs Redirects) Keys() []string {
	ks := make([]string, len(rs))
	for i, r := range rs {
		ks[i] = r.Name
	}

	return ks
}

// Equals checks that rs and o are the same redirect slices. Because redirect
// application depends on order we verify that the slices have the same order.
func (rs Redirects) Equals(o Redirects) bool {
	if len(rs) != len(o) {
		return false
	}

	for i, r := range rs {
		if !r.Equals(o[i]) {
			return false
		}
	}

	return true
}

// IsValid verifies that no two Redirect entries have the same name and that each
// Redirect definition contained is valid.
func (rs Redirects) IsValid() *ValidationError {
	errs := &ValidationError{}

	keySeen := map[string]bool{}
	for _, r := range rs {
		if keySeen[r.Name] {
			errs.AddNew(ErrorCase{
				"",
				fmt.Sprintf(
					"name must be unique, multiple redirects found called '%v'", r.Name),
			})
		}

		keySeen[r.Name] = true
		errs.MergePrefixed(r.IsValid(), "")
	}

	return errs.OrNil()
}

// Redirect specifies how URLs within this domain may need to be rewritten.
// Each Redirect has a name, a regex that matches the requested URL, a to
// indicating how the url should be rewritten, and a flag to indicate how the
// redirect will be handled by the proxying layer.
//
// From may include capture groups which may be referenced by "$<group number>".
//
//   Example:
//     Redirect{
//       Name:        "force-https",
//       From:        "(.*)",
//       To:          "https://www.example.com/$1",
//       RedirectType: PermanentRedirect,
//     }
type Redirect struct {
	Name         string       `json:"name"`
	From         string       `json:"from"`
	To           string       `json:"to"`
	RedirectType RedirectType `json:"redirect_type"`
}

func (r Redirect) Equals(o Redirect) bool {
	return r == o
}

// IsValid checks the validity of a Redirect; we currently verify that a
// redirect:
//
//   * has a non-empty name matching NamePattern
//   * contains a valid regex in From
//   * contains a non-empty to
//   * has a valid redirect type
//
// It is worth noting that no attempt is made to verify the capture group
// mapping into the 'To' field or ensure that the 'To' field is even a valid
// URL after mapping is done.
func (r Redirect) IsValid() *ValidationError {
	errs := &ValidationError{}
	ecase := func(f, m string) ErrorCase {
		return ErrorCase{
			fmt.Sprintf("redirects[%s].%s", r.Name, f), m}
	}

	if r.Name == "" {
		errs.AddNew(ecase("name", "must not be empty"))
	} else {
		if !NamePattern.Match([]byte(r.Name)) {
			errs.AddNew(ecase(
				"name",
				fmt.Sprintf("must match %s", NamePattern.String()),
			))
		}
	}

	if r.From == "" {
		errs.AddNew(ecase("from", "must not be empty"))
	}

	if _, e := regexp.Compile(r.From); e != nil {
		errs.AddNew(ecase("from", fmt.Sprintf("invalid url match expression '%v'", e)))
	}

	if r.To == "" {
		errs.AddNew(ecase("to", "must not be empty"))
	}

	switch r.RedirectType {
	case PermanentRedirect, TemporaryRedirect:
	default:
		errs.AddNew(ecase("type", fmt.Sprintf("'%s' is an invalid redirect type", r.RedirectType)))
	}

	return errs.OrNil()
}
