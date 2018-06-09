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

// Package api defines the types used by the Turbine Labs public API.
package api

import (
	"fmt"
	"regexp"
	"sort"
	"strings"
)

// TbnPublicVersion is the currently released version of open-sourced Turbine
// Labs projects
const TbnPublicVersion = "0.17.2"

const (
	// AllowedIndexPatternStr is a regexp pattern describing the acceptable
	// contents of something that may be used as an index in a changelog
	// path.
	AllowedIndexPatternStr = `^[^\[\]]+$`

	// KeyPatternStr is the regexp pattern one of our keys is expected to meet
	KeyPatternStr = `^[0-9a-zA-Z]+(-[0-9a-zA-Z]+)*$`

	// AllowedIndexPatternMatchFailure is a message describing failure to match
	// the pattern for what may be used as an index entry.
	AllowedIndexPatternMatchFailure = "may not contain [ or ] characters"

	// KeyPatternMatchFailure is a message returned when some key did not match
	// the pattern required.
	KeyPatternMatchFailure = "must match pattern: " + KeyPatternStr

	// HeaderNamePatternStr specifies what can be used when a header name is
	// required.
	HeaderNamePatternStr = "^[0-9A-Za-z-]+$"

	// CookieNamePatternStr specifies what can be used when a cooke name is
	// required.
	CookieNamePatternStr = "^[0-9A-Za-z_.-]+$"
)

var (
	// AllowedIndexPattern is the set of characters which are allowed in a
	// value that will be used as an index component of a changelog path.
	AllowedIndexPattern = regexp.MustCompile(AllowedIndexPatternStr)

	// KeyPattern is the pattern that a key must match. This is a more strict
	// form of AllowedIndexPattern.
	KeyPattern = regexp.MustCompile(KeyPatternStr)

	// HeaderNamePattern is the pattern that a string which will be used an a
	// HTTP header.
	HeaderNamePattern = regexp.MustCompile(HeaderNamePatternStr)

	// CookieNamePattern is the pattern that a string must match if it going to
	// be used as a cookie name.
	CookieNamePattern = regexp.MustCompile(CookieNamePatternStr)
)

// ErrorCase represents an error in an API object. It contains both the
// attribute indicated as, approximately, a dot-separated path to the field
// and a description of the error.
type ErrorCase struct {
	Attribute string `json:"attribute"`
	Msg       string `json:"msg"`
}

// ValidationError contains any errors that were found while trying to validate
// an API object.
type ValidationError struct {
	Errors []ErrorCase `json:"errors"`
}

func (ve *ValidationError) Error() string {
	plural := "s"
	if len(ve.Errors) == 1 {
		plural = ""
	}
	msg := fmt.Sprintf("%d validation error%s", len(ve.Errors), plural)

	for _, c := range ve.Errors {
		msg += "; " + fmt.Sprintf("%s: %s", c.Attribute, c.Msg)
	}

	return msg
}

// AddNew appends a new ErrorCase to the set of errors seen by this ValidationError
func (ve *ValidationError) AddNew(c ErrorCase) {
	ve.Errors = append(ve.Errors, c)
}

// OrNil will return a pointer to the ValidationError if any errors have been
// collected. If no errors have been collected it will return nil.
func (ve *ValidationError) OrNil() *ValidationError {
	if len(ve.Errors) == 0 {
		return nil
	}

	return ve
}

// Merge adds the errors collected in o to the errors in this ValidationError.
func (ve *ValidationError) Merge(o *ValidationError) {
	if o == nil {
		return
	}

	for _, e := range o.Errors {
		ve.AddNew(e)
	}
}

// MergePrefixed takes the Errors found in in children and appends them to this
// ValidationError. In the process it attaches the under prefix to the Attribute
// of the error case. The original children error is not modified.
func (ve *ValidationError) MergePrefixed(children *ValidationError, under string) {
	if children == nil {
		return
	}

	c2 := &ValidationError{}
	for _, e := range children.Errors {
		delim := ""
		if e.Attribute != "" && under != "" {
			delim = "."
		}

		c2.Errors = append(
			c2.Errors,
			ErrorCase{fmt.Sprintf("%s%s%s", under, delim, e.Attribute), e.Msg},
		)
	}

	ve.Merge(c2)
}

func errCheckPattern(
	mayBeEmpty bool,
	value string,
	with *ValidationError,
	vs *regexp.Regexp,
	named,
	failureMsg string,
) {
	if !mayBeEmpty && strings.TrimSpace(value) == "" {
		with.AddNew(ErrorCase{named, "may not be empty"})
		return
	}

	if !vs.MatchString(value) {
		if failureMsg == "" {
			failureMsg = fmt.Sprintf("must match %v", vs.String())
		}
		with.AddNew(ErrorCase{named, failureMsg})
	}
}

func errCheckKey(key string, with *ValidationError, named string) {
	errCheckPattern(false, key, with, KeyPattern, named, KeyPatternMatchFailure)
}

func errCheckIndex(v string, with *ValidationError, named string) {
	errCheckPattern(false, v, with, AllowedIndexPattern, named, AllowedIndexPatternMatchFailure)
}

// ValidationErrorsByAttribute implements sort.Interface
type ValidationErrorsByAttribute struct {
	e *ValidationError
}

var _ sort.Interface = ValidationErrorsByAttribute{}

func (eca ValidationErrorsByAttribute) Len() int {
	if eca.e == nil {
		return 0
	}
	return len(eca.e.Errors)
}

func (eca ValidationErrorsByAttribute) Less(i, j int) bool {
	return eca.e.Errors[i].Attribute < eca.e.Errors[j].Attribute
}

func (eca ValidationErrorsByAttribute) Swap(i, j int) {
	eca.e.Errors[i], eca.e.Errors[j] = eca.e.Errors[j], eca.e.Errors[i]
}

type HasZone interface {
	HasKey
	GetZoneKey() ZoneKey
}

func GetZoneKey(i interface{}) (ZoneKey, bool) {
	if z, ok := i.(HasZone); ok {
		return z.GetZoneKey(), true
	}
	return "", false
}

type HasOrg interface {
	HasKey
	GetOrgKey() OrgKey
}

func GetOrgKey(i interface{}) (OrgKey, bool) {
	if o, ok := i.(HasOrg); ok {
		return o.GetOrgKey(), true
	}
	return "", false
}

type HasChecksum interface {
	HasKey
	GetChecksum() Checksum
}

func GetChecksum(i interface{}) (Checksum, bool) {
	if c, ok := i.(HasChecksum); ok {
		return c.GetChecksum(), true
	}
	return Checksum{}, false
}

type HasKey interface {
	Key() string
}

func GetKey(i interface{}) (string, bool) {
	if k, ok := i.(HasKey); ok {
		return k.Key(), true
	}
	return "", false
}
