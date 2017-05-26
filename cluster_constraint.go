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
	ResponseData  ResponseData  `json:"response_data"`
	Weight        uint32        `json:"weight"`
}

// Checks a ClusterConstraint for validity. A valid constraint will have a non
// empty ConstraintKey and ClusterKey (always), a Weight greater than 0, and
// valid metadata.
func (cc ClusterConstraint) IsValid() *ValidationError {
	errs := &ValidationError{}

	errCheckKey(string(cc.ConstraintKey), errs, "constraint_key")
	errCheckKey(string(cc.ClusterKey), errs, "cluster_key")

	if cc.Weight <= 0 {
		errs.AddNew(ErrorCase{"weight", "must be greater than 0"})
	}

	errs.MergePrefixed(ConstraintMetadataValid(cc.Metadata), "")
	errs.MergePrefixed(ConstraintPropertiesValid(cc.Properties), "")
	errs.MergePrefixed(cc.ResponseData.IsValid(), "response_data")

	return errs.OrNil()
}

func ConstraintMetadataValid(m Metadata) *ValidationError {
	return MetadataValid("metadata", m, func(kv Metadatum) *ValidationError {
		ecase := func(f, msg string) ErrorCase {
			return ErrorCase{f, msg}
		}

		errs := &ValidationError{}

		if kv.Key == "" {
			errs.AddNew(ecase("key", "must not be empty"))
		}
		if kv.Value == "" {
			errs.AddNew(ecase("value", "must not be empty"))
		}

		return errs.OrNil()
	})
}

func ConstraintPropertiesValid(m Metadata) *ValidationError {
	return MetadataValid("properties", m, func(kv Metadatum) *ValidationError {
		if kv.Key != "" {
			return nil
		}

		err := &ValidationError{}
		err.AddNew(ErrorCase{"key", "must not be empty"})
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

	TODO: do we need to identify/declare which requests are idempotent?
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

func (ccs ClusterConstraints) AsMap() map[ConstraintKey]ClusterConstraint {
	m := map[ConstraintKey]ClusterConstraint{}
	for _, cc := range ccs {
		m[cc.ConstraintKey] = cc
	}

	return m
}

// Checks validity of a slice of cluster constraints. This means that each item
// in the slice must be valid and no constraint key may be duplicated.
func (ccs ClusterConstraints) IsValid(container string) *ValidationError {
	errs := &ValidationError{}

	seenKey := map[ConstraintKey]bool{}
	for _, c := range ccs {
		if seenKey[c.ConstraintKey] {
			errs.AddNew(ErrorCase{
				container,
				fmt.Sprintf("multiple instances of key %s", string(c.ConstraintKey)),
			})
		}
		seenKey[c.ConstraintKey] = true

		errs.MergePrefixed(c.IsValid(), fmt.Sprintf("%s[%v]", container, c.ConstraintKey))
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
		cc.Properties.Equals(o.Properties) &&
		cc.ResponseData.Equals(o.ResponseData)
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
func (cc AllConstraints) IsValid(container string) *ValidationError {
	errs := &ValidationError{}
	if len(cc.Light) < 1 {
		errs.AddNew(ErrorCase{"", "must have at least one light constraint"})
	}

	errs.MergePrefixed(cc.Light.IsValid("light"), container)
	errs.MergePrefixed(cc.Dark.IsValid("dark"), container)
	errs.MergePrefixed(cc.Tap.IsValid("tap"), container)

	return errs.OrNil()
}
