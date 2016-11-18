package api

import (
	"fmt"
)

type ConstraintKey string

/*
	A ClusterConstraint describes a filtering of the Instances in a Cluster
	based on their Metadata. Instances in the keyed cluster with a superset
	of the specified Metadata will be included. The Weight of the
	ClusterConstraint is used to inform selection of one ClusterConstraint
	over another.
*/
type ClusterConstraint struct {
	ConstraintKey ConstraintKey `json:"constraint_key"`
	ClusterKey    ClusterKey    `json:"cluster_key"`
	Metadata      Metadata      `json:"metadata"`
	Properties    Metadata      `json:"properties"`
	Weight        uint32        `json:"weight"`
}

// Checks a ClusterConstraint for validity. A valid constraint will have a non
// empty ConstraintKey and ClusterKey (always), a Weight greater than 0, and
// valid metadata.
func (cc ClusterConstraint) IsValid(precreation bool) *ValidationError {
	ecase := func(f, m string) ErrorCase {
		return ErrorCase{fmt.Sprintf("constraint[%s].%s", string(cc.ConstraintKey), f), m}
	}

	errs := &ValidationError{}

	if cc.ConstraintKey == "" {
		errs.AddNew(ecase("constraint_key", "must not be empty"))
	}

	if cc.ClusterKey == "" {
		errs.AddNew(ecase("cluster_key", "must not be empty"))
	}

	if cc.Weight <= 0 {
		errs.AddNew(ecase("weight", "must be greater than 0"))
	}

	errs.MergePrefixed(
		ConstraintMetadataValid(cc.Metadata),
		fmt.Sprintf("constraint[%s].metadata", string(cc.ConstraintKey)))

	errs.MergePrefixed(
		ConstraintPropertiesValid(cc.Properties),
		fmt.Sprintf("constraint[%s].properties", string(cc.ConstraintKey)))

	return errs.OrNil()
}

func ConstraintMetadataValid(m Metadata) *ValidationError {
	return MetadataValid(m, func(kv Metadatum) *ValidationError {
		ecase := func(f, msg string) ErrorCase {
			return ErrorCase{fmt.Sprintf("metadatum[%s].%s", kv.Key, f), msg}
		}

		errs := &ValidationError{}

		if kv.Key == "" {
			errs.AddNew(ecase("key", "key must not be empty"))
		}
		if kv.Value == "" {
			errs.AddNew(ecase("value", "value must not be empty"))
		}

		return errs.OrNil()
	})
}

func ConstraintPropertiesValid(m Metadata) *ValidationError {
	return MetadataValid(m, func(kv Metadatum) *ValidationError {
		if kv.Key != "" {
			return nil
		}

		err := &ValidationError{}
		err.AddNew(ErrorCase{"property[].key", "key must not be empty"})
		return err
	})
}

/*
	AllConstraints define three different ClusterConstraint slices.
	The Light ClusterConstraint slice is used to determine the Instance to
	which the live request will be sent and from which the response will be
	sent to the caller. The Dark ClusterConstraint slice is used to determine
	an Instance to send a send-and-forget copy of the request to. The Tap
	ClusterConstraint slice is used to determine an Instance to send a copy of
	the request to, comparing the response to the Light response.

        The Dark and Tap ClusterConstraint slices may be empty. The Light
        ClusterConstraint slice must always contain at least one entry.

	TODO: do we need to idenfity/declare which requests are idempotent?
	If Routes are structured properly, this isn't necessary, since you can only
	add Dark/Tap ClusterConstraints for Routes that are safe to call more than
	once for the same input.
*/
type AllConstraints struct {
	Light ClusterConstraints `json:"light"`
	Dark  ClusterConstraints `json:"dark"`
	Tap   ClusterConstraints `json:"tap"`
}

type ClusterConstraints []ClusterConstraint

// Check the Equality of two ClusterConstraint slices. They are (currently)
// equal iff each of the elements of a slice is equal to the corresponding
// element in the other slice. That is: order of the contents matters.
//
// TODO: make this call order agnostic - https://github.com/turbinelabs/tbn/issues/188
func (ccs ClusterConstraints) Equals(o ClusterConstraints) bool {
	if len(ccs) != len(o) {
		return false
	}

	for i, cc := range ccs {
		if !o[i].Equals(cc) {
			return false
		}
	}

	return true
}

// Checks validity of a slice of cluster constraints. This means that each item
// in the slice must be valid and no constraint key may be duplicated.
func (ccs ClusterConstraints) IsValid(precreation bool) *ValidationError {
	errs := &ValidationError{}

	seenKey := map[ConstraintKey]bool{}
	for _, c := range ccs {
		if seenKey[c.ConstraintKey] {
			errs.AddNew(ErrorCase{
				"constraint_key",
				fmt.Sprintf("multiple instances of key %s", string(c.ConstraintKey)),
			})
		}
		seenKey[c.ConstraintKey] = true

		errs.MergePrefixed(c.IsValid(precreation), "")
	}

	return errs.OrNil()
}

func (cc ClusterConstraint) coreEquality(o ClusterConstraint) bool {
	return cc.ConstraintKey == o.ConstraintKey &&
		cc.Weight == o.Weight &&
		cc.ClusterKey == o.ClusterKey
}

// Check equality between two ClusterConstraint objects. For these to be
// considered equal they must share the same ConstraintKey, ClusterKey, Weight,
// and Metadata.
func (cc ClusterConstraint) Equals(o ClusterConstraint) bool {
	return cc.coreEquality(o) &&
		cc.Metadata.Equals(o.Metadata) &&
		cc.Properties.Equals(o.Properties)
}

// Check equality between two AllConstraints objects. The objects are equal
// if each component (Light, Dark, and Tap) is equal.
func (cc AllConstraints) Equals(o AllConstraints) bool {
	return cc.Light.Equals(o.Light) &&
		cc.Dark.Equals(o.Dark) &&
		cc.Tap.Equals(o.Tap)
}

// Check validity of an AllConstraints struct. A valid AllConstraints must have
// at lesat one Light constraint and valid Light, Dark, and Tap constraints.
func (cc AllConstraints) IsValid(precreation bool) *ValidationError {
	errs := &ValidationError{}
	if len(cc.Light) < 1 {
		errs.AddNew(ErrorCase{"light", "must have at least one constraint"})
	}

	errs.MergePrefixed(cc.Light.IsValid(precreation), "light")
	errs.MergePrefixed(cc.Dark.IsValid(precreation), "dark")
	errs.MergePrefixed(cc.Tap.IsValid(precreation), "tap")

	return errs.OrNil()
}