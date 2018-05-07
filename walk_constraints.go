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

// ConstraintVisitor is a function that operates on a ClusterConstraint, and
// the enclosing Rule, if any.
type ConstraintVisitor func(ClusterConstraint, *Rule)

// WalkConstraintsFromSharedRules applies the vistior function to all the
// ClusterConstraints referenced by the given SharedRules
func WalkConstraintsFromSharedRules(sr SharedRules, visit ConstraintVisitor) {
	WalkConstraintsFromAllConstraints(sr.Default, nil, visit)
	WalkConstraintsFromRules(sr.Rules, visit)
}

// WalkConstraintsFromRules applies the vistior function to all the
// ClusterConstraints, along with the enclosing rule, referenced by the given
// Rules.
func WalkConstraintsFromRules(rs Rules, visit ConstraintVisitor) {
	for _, r := range rs {
		WalkConstraintsFromAllConstraints(r.Constraints, &r, visit)
	}
}

// WalkConstraintsFromAllConstraints applies the vistior function to all the
// constraints referenced by the given ClusterConstraints, and the given rule
// pointer.
func WalkConstraintsFromAllConstraints(ac AllConstraints, rule *Rule, visit ConstraintVisitor) {
	WalkConstraints(ac.Light, rule, visit)
	WalkConstraints(ac.Dark, rule, visit)
	WalkConstraints(ac.Tap, rule, visit)
}

// WalkConstraints applies the vistior function to all the elements referenced
// by the given ClusterConstraints, and the given rule pointer.
func WalkConstraints(ccf ClusterConstraints, rule *Rule, visit ConstraintVisitor) {
	for _, cc := range ccf {
		visit(cc, rule)
	}
}
