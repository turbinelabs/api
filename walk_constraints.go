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
