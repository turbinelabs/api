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

package changelog

import (
	"testing"

	"github.com/turbinelabs/test/assert"
)

func TestApplyExpr(t *testing.T) {
	var expr FilterExpr
	expr = FilterOrs{[]FilterAnds{
		{[]Filter{{}, {}, {}, {}}},
		{[]Filter{{}}},
		{},
	}}

	newExpr := ApplyExpr(expr, func(f Filter) Filter {
		f.ObjectType = "aoeu"
		return f
	})

	ors := newExpr.AsExpr()

	assert.Equal(t, len(ors.FilterAnds), 3)

	ands1 := ors.FilterAnds[0]
	ands2 := ors.FilterAnds[1]
	ands3 := ors.FilterAnds[2]

	assert.Equal(t, len(ands1.Filters), 4)
	assert.Equal(t, len(ands2.Filters), 1)
	assert.Equal(t, len(ands3.Filters), 0)

	for _, a := range ors.FilterAnds {
		for _, f := range a.Filters {
			assert.DeepEqual(t, f, Filter{ObjectType: "aoeu"})
		}
	}
}
