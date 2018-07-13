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
	"fmt"

	tbnstrings "github.com/turbinelabs/nonstdlib/strings"
)

// TracingConfig describes how tracing operations should be applied
// to the given Listener.
type TracingConfig struct {
	// Ingress, when true, specifies that this listener is handling requests from a downstream.
	// When false it indicates that it is handling requests bound to an upstream.
	Ingress bool `json:"ingress"`
	// Each listed header will be added to generated spans as an annotation
	RequestHeadersForTags []string `json:"request_headers_for_tags"`
}

// Equals compares two TraceConfig objects returning true if they are the same.
// RequestHeadersForTags is compared without regard for ordering of its content.
func (tc TracingConfig) Equals(o TracingConfig) bool {
	cmp := func(tcs, os []string) bool {
		s1 := tbnstrings.NewSet(tcs...)
		s2 := tbnstrings.NewSet(os...)
		return s1.Equals(s2)
	}

	return tc.Ingress == o.Ingress &&
		cmp(tc.RequestHeadersForTags, o.RequestHeadersForTags)
}

func (tc TracingConfig) IsValid() *ValidationError {
	errs := &ValidationError{}
	ecase := func(f, m string) ErrorCase {
		return ErrorCase{f, m}
	}
	for _, h := range tc.RequestHeadersForTags {
		if !HeaderNamePattern.MatchString(h) {
			errs.AddNew(ecase(
				"request_headers_for_tags",
				fmt.Sprintf("header %s is not a valid HTTP header name. Must match %s",
					h, HeaderNamePatternStr)))
		}
	}
	return errs.OrNil()
}
