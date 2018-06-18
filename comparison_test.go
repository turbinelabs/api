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

func TestCompareInts(t *testing.T) {
	assert.Equal(t, compareInts(1, 2), -1)
	assert.Equal(t, compareInts(2, 2), 0)
	assert.Equal(t, compareInts(3, 2), 1)
}

func TestCompareStrings(t *testing.T) {
	tcs := []struct {
		left, right []string
		expected    int
	}{
		{
			left:     nil,
			right:    nil,
			expected: 0,
		},
		{
			left:     nil,
			right:    []string{"a"},
			expected: -1,
		},
		{
			left:     []string{"a"},
			right:    nil,
			expected: 1,
		},
		{
			left:     []string{"a", "b", "c", "d", "z"},
			right:    []string{"a", "b", "c", "d", "a"},
			expected: 1,
		},
		{
			left:     []string{"a", "b", "c"},
			right:    []string{"x", "y", "z"},
			expected: -1,
		},
		{
			left:     []string{"a", "b", "c"},
			right:    []string{"a", "b", "c"},
			expected: 0,
		},
	}

	for i, tc := range tcs {
		assert.Group(
			fmt.Sprintf("testCases[%d]: left=[%#v], right=[%#v]", i, tc.left, tc.right),
			t,
			func(g *assert.G) {
				assert.Equal(g, compareStrings(tc.left, tc.right), tc.expected)
			},
		)
	}
}
