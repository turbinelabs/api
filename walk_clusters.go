package api

// WalkClusterKeysFromSharedRules applies the vistior function to all the
// ClusterKeys referenced by the given SharedRules
func WalkClusterKeysFromSharedRules(sr SharedRules, visit func(ClusterKey)) {
	WalkClustersFromAllConstraints(sr.Default, visit)
	WalkClusterKeysFromRules(sr.Rules, visit)
}

// WalkClusterKeysFromRules applies the vistior function to all the ClusterKeys
// referenced by the given Rules
func WalkClusterKeysFromRules(rs Rules, visit func(ClusterKey)) {
	for _, r := range rs {
		WalkClustersFromAllConstraints(r.Constraints, visit)
	}
}

// WalkClustersFromAllConstraints applies the vistior function to all
// ClusterKeys referenced by the given AllConstraints
func WalkClustersFromAllConstraints(ac AllConstraints, visit func(ClusterKey)) {
	WalkClustersFromConstraints(ac.Light, visit)
	WalkClustersFromConstraints(ac.Dark, visit)
	WalkClustersFromConstraints(ac.Tap, visit)
}

// WalkClustersFromConstraints applies the vistior function to all the
// ClusterKeys referenced by the given ClusterConstraints
func WalkClustersFromConstraints(ccf ClusterConstraints, visit func(ClusterKey)) {
	for _, cc := range ccf {
		visit(cc.ClusterKey)
	}
}
