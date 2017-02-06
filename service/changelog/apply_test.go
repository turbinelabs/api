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
