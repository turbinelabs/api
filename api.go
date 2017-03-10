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

// Package api defines the types used by the Turbine Labs public API.
package api

import (
	"fmt"
	"regexp"
	"sort"
	"strings"
)

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
)

var (
	// AllowedIndexPattern is the set of characters which are allowed in a
	// value that will be used as an index component of a changelog path.
	AllowedIndexPattern = regexp.MustCompile(AllowedIndexPatternStr)

	// KeyPattern is the pattern that a key must match. This is a more strict
	// form of AllowedIndexPattern.
	KeyPattern = regexp.MustCompile(KeyPatternStr)
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

func errCheckKey(key string, with *ValidationError, named string) {
	if strings.TrimSpace(key) == "" {
		with.AddNew(ErrorCase{named, "may not be empty"})
	} else if !KeyPattern.MatchString(key) {
		with.AddNew(ErrorCase{named, KeyPatternMatchFailure})
	}
}

func errCheckIndex(v string, with *ValidationError, named string) {
	if strings.TrimSpace(v) == "" {
		with.AddNew(ErrorCase{named, "may not be empty"})
	} else if !AllowedIndexPattern.MatchString(v) {
		with.AddNew(ErrorCase{named, AllowedIndexPatternMatchFailure})
	}
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
