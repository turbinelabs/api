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

package changelog

// FilterMutator is a function that returns an updated version of the Filter
// provided as input. Note that due to pass-by-value making a shollow copy of
// the Filter it's possibly to mutate the original filter and proper caution
// should be taken if that would be problematic in your use case.
type FilterMutator func(Filter) Filter

// ApplyExpr applies some FilterMutator to all Filters within an expression
// then returns a new expression. The canonical expression format is a
// FilterOrs which is composed of a []FilterAnds. As slices are pointery it's
// possible for ApplyExpr to mutate the originating Filters, FilterAnds, or
// FilterOrs depending on how AsExpr generates the canonical format and how
// the FilterMutator acts.
func ApplyExpr(expr FilterExpr, fn FilterMutator) FilterExpr {
	if expr == nil {
		return nil
	}
	return applyOrs(expr.AsExpr(), fn)
}

func applyOrs(ors FilterOrs, fn FilterMutator) FilterOrs {
	// TODO: build a sample thing that tests if the index deref is needed
	for i, a := range ors.FilterAnds {
		ors.FilterAnds[i] = applyAnds(a, fn)
	}

	return ors
}

func applyAnds(ands FilterAnds, fn FilterMutator) FilterAnds {
	for i, f := range ands.Filters {
		ands.Filters[i] = fn(f)
	}

	return ands
}
