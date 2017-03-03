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
	"sort"
	"testing"

	"github.com/turbinelabs/test/assert"
)

func TestClusterConstraintsEqualsSuccess(t *testing.T) {
	cc1 := ClusterConstraint{"cckey1", "ckey1", Metadata{{"key", "value"}, {"key2", "value2"}}, Metadata{{"state", "released"}}, 1234}
	cc2 := ClusterConstraint{"cckey2", "ckey2", Metadata{{"key-2", "value-2"}}, Metadata{{"stata", "testing"}}, 1234}

	slice1 := ClusterConstraints{cc1, cc2}
	slice2 := ClusterConstraints{cc1, cc2}

	assert.True(t, slice2.Equals(slice1))
	assert.True(t, slice1.Equals(slice2))
}

func TestClusterConstraintsEqualsFailureMetadata(t *testing.T) {
	cc1 := ClusterConstraint{"cckey1", "ckey1", Metadata{{"key", "value"}, {"key2", "value2"}}, nil, 1234}
	cc2a := ClusterConstraint{"cckey2", "ckey2", Metadata{}, nil, 1234}
	cc2b := ClusterConstraint{"cckey2", "ckey2", Metadata{{"key-2", "value-2"}}, nil, 1234}

	slice1 := ClusterConstraints{cc1, cc2a}
	slice2 := ClusterConstraints{cc1, cc2b}

	assert.False(t, slice1.Equals(slice2))
	assert.False(t, slice2.Equals(slice1))
}

func TestClusterConstraintsEqualsFailureProperties(t *testing.T) {
	cc1 := ClusterConstraint{"cckey1", "ckey1", nil, Metadata{{"state", "released"}}, 1234}
	cc2a := ClusterConstraint{"cckey2", "ckey2", nil, Metadata{{"state", "released"}}, 1234}
	cc2b := ClusterConstraint{"cckey2", "ckey2", nil, Metadata{{"state", "releasing"}}, 1234}

	slice1 := ClusterConstraints{cc1, cc2a}
	slice2 := ClusterConstraints{cc1, cc2b}

	assert.False(t, slice1.Equals(slice2))
	assert.False(t, slice2.Equals(slice1))
}

func TestClusterConstraintsEqualsLengthMismatch(t *testing.T) {
	cc1 := ClusterConstraint{"cckey1", "ckey1", Metadata{{"key", "value"}, {"key2", "value2"}}, nil, 1234}
	cc2 := ClusterConstraint{"cckey2", "ckey2", Metadata{{"key", "value"}, {"key2", "value2"}}, nil, 1234}

	slice1 := ClusterConstraints{cc1, cc2}
	slice2 := ClusterConstraints{}

	assert.False(t, slice1.Equals(slice2))
	assert.False(t, slice2.Equals(slice1))
}

func getClusterConstraintPair() (ClusterConstraint, ClusterConstraint) {
	cc1 := ClusterConstraint{"cckey1", "ckey1", Metadata{{"key", "value"}, {"key2", "value2"}}, nil, 1234}
	cc2 := ClusterConstraint{"cckey1", "ckey1", Metadata{{"key", "value"}, {"key2", "value2"}}, nil, 1234}
	return cc1, cc2
}

func TestClusterConstraintEqualsSuccess(t *testing.T) {
	cc1, cc2 := getClusterConstraintPair()

	assert.True(t, cc1.Equals(cc2))
	assert.True(t, cc2.Equals(cc1))
}

func TestClusterConstraintEqualsConstraintKeyVaries(t *testing.T) {
	cc1, cc2 := getClusterConstraintPair()
	cc2.ConstraintKey = "cckey2"

	assert.False(t, cc1.Equals(cc2))
	assert.False(t, cc2.Equals(cc1))
}

func TestClusterConstraintEqualsClusterVaries(t *testing.T) {
	cc1, cc2 := getClusterConstraintPair()
	cc2.ClusterKey = "ckey2"

	assert.False(t, cc1.Equals(cc2))
	assert.False(t, cc2.Equals(cc1))
}

func TestClusterConstraintEqualsWeightVaries(t *testing.T) {
	cc1, cc2 := getClusterConstraintPair()
	cc2.Weight = 1235

	assert.False(t, cc1.Equals(cc2))
	assert.False(t, cc2.Equals(cc1))
}

func TestClusterConstraintEqualsMetadataVaries(t *testing.T) {
	cc1, cc2 := getClusterConstraintPair()
	cc2.Metadata = Metadata{{"key-variation", "value"}, {"key2", "value2"}}

	assert.False(t, cc1.Equals(cc2))
	assert.False(t, cc2.Equals(cc1))
}

// ClusterConstraint.IsValid
func getValidClusterConstraint() ClusterConstraint {
	return ClusterConstraint{"cckey1", "ckey1", Metadata{{"key", "value"}, {"k", "v"}}, nil, 1234}
}

func TestClusterConstraintIsValidSuccess(t *testing.T) {
	cc := getValidClusterConstraint()

	assert.Nil(t, cc.IsValid())
}

func TestClusterConstraintIsValidWithoutConstraintKeyFailure(t *testing.T) {
	cc := getValidClusterConstraint()
	cc.ConstraintKey = ""

	assert.NonNil(t, cc.IsValid())
}

func TestClusterConstraintIsValidCkeyfailure(t *testing.T) {
	cc := getValidClusterConstraint()
	cc.ClusterKey = ""

	assert.NonNil(t, cc.IsValid())
}

func TestClusterConstraintIsValidMetadataFailure(t *testing.T) {
	cc := getValidClusterConstraint()
	cc.Metadata = Metadata{{"key", "value"}, {"", ""}}

	assert.NonNil(t, cc.IsValid())
}

func TestClusterConstraintIsValidWeightFailure(t *testing.T) {
	cc := getValidClusterConstraint()
	cc.Weight = 0

	assert.NonNil(t, cc.IsValid())
}

// ClusterConstarintSlice.IsValid
func TestClusterConstraintsIsValidSuccess(t *testing.T) {
	cc1 := ClusterConstraint{"cckey1", "ckey1", Metadata{{"key", "value"}, {"k", "v"}}, Metadata{}, 234}
	cc2 := ClusterConstraint{"cckey2", "ckey2", Metadata{{"key", "value"}, {"k", "v"}}, Metadata{}, 123}
	ccs := ClusterConstraints{cc1, cc2}

	assert.Nil(t, ccs.IsValid("test"))
}

func TestClusterConstraintsIsValidFailureOnDuplicateConstraintKeys(t *testing.T) {
	cc1 := ClusterConstraint{"cckey1", "ckey", Metadata{{"key", "value"}, {"k", "v"}}, Metadata{}, 234}
	cc2 := ClusterConstraint{"cckey2", "ckey2", Metadata{{"key", "value"}, {"k", "v"}}, Metadata{}, 123}
	cc2.ConstraintKey = "cckey1"
	ccs := ClusterConstraints{cc1, cc2}

	assert.NonNil(t, ccs.IsValid("test"))
}

func TestClusterConstraintsIsValidInvalidContents(t *testing.T) {
	cc1 := ClusterConstraint{"cckey1", "ckey", Metadata{{"key", "value"}, {"k", "v"}}, Metadata{}, 234}
	cc2 := ClusterConstraint{"cckey2", "ckey2", Metadata{{"key", "value"}, {"k", "v"}}, Metadata{}, 0}
	ccs := ClusterConstraints{cc1, cc2}

	assert.NonNil(t, ccs.IsValid("test"))
}

func TestClusterConstraintsIsValidEmpty(t *testing.T) {
	ccs := ClusterConstraints{}

	assert.Nil(t, ccs.IsValid("test"))
}

// ClusterConstraints.IsValid
func TestClusterConstraintsIsValidSucces(t *testing.T) {
	cc1 := ClusterConstraint{"cckey1", "ckey1", Metadata{{"key", "value"}, {"k", "v"}}, Metadata{}, 234}
	cc2 := ClusterConstraint{"cckey2", "ckey2", Metadata{{"key", "value"}, {"k", "v"}}, Metadata{}, 123}
	ccs := ClusterConstraints{cc1, cc2}

	set := AllConstraints{ccs, ccs, ccs}

	assert.Nil(t, set.IsValid("test"))
}

func TestClusterConstraintsIsValidFailsWithBadLight(t *testing.T) {
	cc1 := ClusterConstraint{"cckey1", "ckey", Metadata{{"key", "value"}, {"k", "v"}}, Metadata{}, 234}
	cc2 := ClusterConstraint{"cckey2", "ckey2", Metadata{{"key", "value"}, {"k", "v"}}, Metadata{}, 123}
	ccbad := ClusterConstraint{"cc-bad", "ckey2", Metadata{{"key", "value"}, {"k", "v"}}, Metadata{}, 0}

	ccs := ClusterConstraints{cc1, cc2}
	ccsBad := ClusterConstraints{cc1, cc2, ccbad}

	set := AllConstraints{ccs, ccs, ccs}
	set.Light = ccsBad

	assert.NonNil(t, set.IsValid("test"))
}

func TestClusterConstraintsIsValidFailsWithBadDark(t *testing.T) {
	cc1 := ClusterConstraint{"cckey1", "ckey", Metadata{{"key", "value"}, {"k", "v"}}, Metadata{}, 234}
	cc2 := ClusterConstraint{"cckey2", "ckey2", Metadata{{"key", "value"}, {"k", "v"}}, Metadata{}, 123}
	ccbad := ClusterConstraint{"cc-bad", "ckey2", Metadata{{"key", "value"}, {"k", "v"}}, Metadata{}, 0}

	ccs := ClusterConstraints{cc1, cc2}
	ccsBad := ClusterConstraints{cc1, cc2, ccbad}

	set := AllConstraints{ccs, ccs, ccs}
	set.Dark = ccsBad

	assert.NonNil(t, set.IsValid("test"))
}

func TestClusterConstraintsIsValidFailsWithBadTap(t *testing.T) {
	cc1 := ClusterConstraint{"cckey1", "ckey", Metadata{{"key", "value"}, {"k", "v"}}, Metadata{}, 234}
	cc2 := ClusterConstraint{"cckey2", "ckey2", Metadata{{"key", "value"}, {"k", "v"}}, Metadata{}, 123}
	ccbad := ClusterConstraint{"cc-bad", "ckey2", Metadata{{"key", "value"}, {"k", "v"}}, Metadata{}, 0}

	ccs := ClusterConstraints{cc1, cc2}
	ccsBad := ClusterConstraints{cc1, cc2, ccbad}

	set := AllConstraints{ccs, ccs, ccs}
	set.Tap = ccsBad

	assert.NonNil(t, set.IsValid("test"))
}

func TestClusterConstraintsIsValidFailsWithZeroLight(t *testing.T) {
	cc1 := ClusterConstraint{"cckey1", "ckey", Metadata{{"key", "value"}, {"k", "v"}}, Metadata{}, 234}
	cc2 := ClusterConstraint{"cckey2", "ckey2", Metadata{{"key", "value"}, {"k", "v"}}, Metadata{}, 123}
	ccs := ClusterConstraints{cc1, cc2}

	set := AllConstraints{ccs, ccs, ccs}
	set.Light = ClusterConstraints{}

	assert.NonNil(t, set.IsValid("test"))
}

func TestConstraintMetadataValid(t *testing.T) {
	m := Metadata{
		Metadatum{"key", "value"},
		Metadatum{"key2", "value"},
	}
	assert.Nil(t, ConstraintMetadataValid(m))
}

func TestConstraintMetadataValidFailedDupes(t *testing.T) {
	mgood := Metadatum{"key", "value"}
	m := Metadata{mgood, mgood, mgood}
	assert.DeepEqual(t, ConstraintMetadataValid(m), &ValidationError{[]ErrorCase{
		{"metadata", "duplicate metadata key 'key'"},
		{"metadata", "duplicate metadata key 'key'"},
	}})
}

func TestConstraintMetadataValidFailed(t *testing.T) {
	mgood := Metadatum{"key", "value"}
	mbad1 := Metadatum{"key-2", ""}
	mbad2 := Metadatum{"", "value-2"}
	m := Metadata{mgood, mbad1, mbad2}
	errs := ConstraintMetadataValid(m)
	sort.Sort(ValidationErrorsByAttribute{errs})
	assert.DeepEqual(t, errs, &ValidationError{[]ErrorCase{
		{"metadata[].key", "must not be empty"},
		{"metadata[key-2].value", "must not be empty"},
	}})
}

func TestConstraintPropertiesValid(t *testing.T) {
	mgood := Metadatum{"key", "value"}
	mgood2 := Metadatum{"key2", ""}
	m := Metadata{mgood, mgood2}
	assert.Nil(t, ConstraintPropertiesValid(m))
}

func TestConstraintPropertiesValidFailed(t *testing.T) {
	mgood := Metadatum{"key", "value"}
	mbad := Metadatum{"", "value"}
	m := Metadata{mbad, mgood}
	errs := ConstraintPropertiesValid(m)
	assert.DeepEqual(t, errs, &ValidationError{[]ErrorCase{
		{"properties[].key", "must not be empty"},
	}})
}

func TestConstraintPropertiesValidFailedDupes(t *testing.T) {
	mgood := Metadatum{"key", "value"}
	m := Metadata{mgood, mgood}
	errs := ConstraintPropertiesValid(m)
	assert.DeepEqual(t, errs, &ValidationError{[]ErrorCase{
		{"properties", "duplicate properties key 'key'"},
	}})
}

func getTestCC() ClusterConstraint {
	md := Metadata{{"md-key", "value"}}
	props := Metadata{{"p-key", "p-value"}}

	return ClusterConstraint{
		"cck-1",
		"ck-1",
		md,
		props,
		100,
	}
}

func TestClusterConstraintIsValid(t *testing.T) {
	cc := getTestCC()
	assert.Nil(t, cc.IsValid())
}

func TestClusterConstraintIsValidBadKey(t *testing.T) {
	cc := getTestCC()
	cc.ConstraintKey = "bad-key!"
	assert.NonNil(t, cc.IsValid())
}

func TestClusterConstraintIsValidNoKey(t *testing.T) {
	cc := getTestCC()
	cc.ConstraintKey = ""
	assert.NonNil(t, cc.IsValid())
}

func TestClusterConstraintIsValidBadClusterKey(t *testing.T) {
	cc := getTestCC()
	cc.ClusterKey = "1234-@-test"
	assert.NonNil(t, cc.IsValid())
}

func TestClusterConstraintIsValidNoClusterKey(t *testing.T) {
	cc := getTestCC()
	cc.ClusterKey = ""
	assert.NonNil(t, cc.IsValid())
}

func TestClusterConstraintIsValidBadWeight(t *testing.T) {
	cc := getTestCC()
	cc.Weight = 0
	assert.NonNil(t, cc.IsValid())
}

func TestClusterConstraintIsValidBadMetadata(t *testing.T) {
	cc := getTestCC()
	cc.Metadata[0].Key = ""
	assert.DeepEqual(t, cc.IsValid(), &ValidationError{[]ErrorCase{
		{"metadata[].key", "must not be empty"},
	}})
}

func TestClusterConstraintIsValidBadProperty(t *testing.T) {
	cc := getTestCC()
	cc.Properties[0].Key = ""
	assert.DeepEqual(t, cc.IsValid(), &ValidationError{[]ErrorCase{
		{"properties[].key", "must not be empty"},
	}})
}

func TestClusterConstraintIsValidNoMetadataNoProperty(t *testing.T) {
	cc := getTestCC()
	cc.Metadata = nil
	cc.Properties = nil
	assert.Nil(t, cc.IsValid())
}

func getTestCCS() ClusterConstraints {
	cc1 := getTestCC()
	cc2 := getTestCC()
	cc2.ConstraintKey = "cck-2"
	return ClusterConstraints{cc1, cc2}
}

func TestClusterConstraintsIsValid(t *testing.T) {
	ccs := getTestCCS()
	assert.Nil(t, ccs.IsValid("ccs"))
}

func TestClusterConstraintsIsValidDupes(t *testing.T) {
	ccs := getTestCCS()
	ccs = append(ccs, ccs[0])
	assert.DeepEqual(t, ccs.IsValid("ccs"), &ValidationError{[]ErrorCase{
		{"ccs", "multiple instances of key cck-1"},
	}})
}

func TestClusterConstraintsIsValidBadConstraint(t *testing.T) {
	ccs := getTestCCS()
	ccs[0].Weight = 0
	assert.DeepEqual(t, ccs.IsValid("ccs"), &ValidationError{[]ErrorCase{
		{"ccs[cck-1].weight", "must be greater than 0"},
	}})
}

func getTestAC() AllConstraints {
	return AllConstraints{
		getTestCCS(),
		getTestCCS(),
		getTestCCS(),
	}
}

func TestAllConstraintsIsValid(t *testing.T) {
	acs := getTestAC()
	assert.Nil(t, acs.IsValid("ac"))
}

func TestAllConstraintsIsValidBad(t *testing.T) {
	acs := getTestAC()
	acs.Light[0].Weight = 0
	acs.Dark[0].Weight = 0
	acs.Tap[0].Weight = 0
	assert.DeepEqual(t, acs.IsValid("ac"), &ValidationError{[]ErrorCase{
		{"ac.light[cck-1].weight", "must be greater than 0"},
		{"ac.dark[cck-1].weight", "must be greater than 0"},
		{"ac.tap[cck-1].weight", "must be greater than 0"},
	}})
}
