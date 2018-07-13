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
	"testing"

	"github.com/turbinelabs/test/assert"
)

func getTracingConfigs() (TracingConfig, TracingConfig) {
	tc := TracingConfig{
		Ingress:               true,
		RequestHeadersForTags: []string{"x-foo", "x-bar"},
	}

	return tc, tc
}

func TestTracingConfigEquals(t *testing.T) {
	tc1, tc2 := getTracingConfigs()

	assert.True(t, tc1.Equals(tc2))
	assert.True(t, tc2.Equals(tc1))
}

func TestTracingConfigEqualsDiffIngress(t *testing.T) {
	tc1, tc2 := getTracingConfigs()
	tc2.Ingress = false
	assert.False(t, tc1.Equals(tc2))
	assert.False(t, tc2.Equals(tc1))
}

func TestTracingConfigEqualsDiffRequestHeaders(t *testing.T) {
	tc1, tc2 := getTracingConfigs()
	tc2.RequestHeadersForTags = []string{"x-foo"}
	assert.False(t, tc1.Equals(tc2))
	assert.False(t, tc2.Equals(tc1))
}

func TestTracingConfigEqualsDiffRequestHeadersOrder(t *testing.T) {
	tc1, tc2 := getTracingConfigs()
	tc2.RequestHeadersForTags = []string{"x-bar", "x-foo"}
	assert.True(t, tc1.Equals(tc2))
	assert.True(t, tc2.Equals(tc1))
}

func mkTestTC() TracingConfig {
	return TracingConfig{
		Ingress:               true,
		RequestHeadersForTags: []string{"x-foo", "x-bar"},
	}
}

func TestTracingConfigIsValid(t *testing.T) {
	tc := mkTestTC()
	assert.Nil(t, tc.IsValid())
}

func TestTracingConfigIsValidBadHeader(t *testing.T) {
	tc := mkTestTC()
	badHeader := "---%%%-"
	tc.RequestHeadersForTags = []string{badHeader}
	gotErr := tc.IsValid()
	assert.DeepEqual(t, gotErr, &ValidationError{[]ErrorCase{
		{"request_headers_for_tags", fmt.Sprintf("header %s is not a valid HTTP header name. Must match %s", badHeader, HeaderNamePattern)},
	}})
}
