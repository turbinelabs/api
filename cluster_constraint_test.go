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

	assert.Nil(t, cc.IsValid(true))
	assert.Nil(t, cc.IsValid(false))
}

func TestClusterConstraintIsValidWithoutConstraintKeyFailure(t *testing.T) {
	cc := getValidClusterConstraint()
	cc.ConstraintKey = ""

	assert.NonNil(t, cc.IsValid(true))
	assert.NonNil(t, cc.IsValid(false))
}

func TestClusterConstraintIsValidCkeyfailure(t *testing.T) {
	cc := getValidClusterConstraint()
	cc.ClusterKey = ""

	assert.NonNil(t, cc.IsValid(true))
	assert.NonNil(t, cc.IsValid(false))
}

func TestClusterConstraintIsValidMetadataFailure(t *testing.T) {
	cc := getValidClusterConstraint()
	cc.Metadata = Metadata{{"key", "value"}, {"", ""}}

	assert.NonNil(t, cc.IsValid(true))
	assert.NonNil(t, cc.IsValid(false))
}

func TestClusterConstraintIsValidWeightFailure(t *testing.T) {
	cc := getValidClusterConstraint()
	cc.Weight = 0

	assert.NonNil(t, cc.IsValid(true))
	assert.NonNil(t, cc.IsValid(false))
}

// ClusterConstarintSlice.IsValid
func TestClusterConstraintsIsValidSuccess(t *testing.T) {
	cc1 := ClusterConstraint{"cckey1", "ckey1", Metadata{{"key", "value"}, {"k", "v"}}, Metadata{}, 234}
	cc2 := ClusterConstraint{"cckey2", "ckey2", Metadata{{"key", "value"}, {"k", "v"}}, Metadata{}, 123}
	ccs := ClusterConstraints{cc1, cc2}

	assert.Nil(t, ccs.IsValid(true))
	assert.Nil(t, ccs.IsValid(false))
}

func TestClusterConstraintsIsValidFailureOnDuplicateConstraintKeys(t *testing.T) {
	cc1 := ClusterConstraint{"cckey1", "ckey", Metadata{{"key", "value"}, {"k", "v"}}, Metadata{}, 234}
	cc2 := ClusterConstraint{"cckey2", "ckey2", Metadata{{"key", "value"}, {"k", "v"}}, Metadata{}, 123}
	cc2.ConstraintKey = "cckey1"
	ccs := ClusterConstraints{cc1, cc2}

	assert.NonNil(t, ccs.IsValid(true))
	assert.NonNil(t, ccs.IsValid(false))
}

func TestClusterConstraintsIsValidInvalidContents(t *testing.T) {
	cc1 := ClusterConstraint{"cckey1", "ckey", Metadata{{"key", "value"}, {"k", "v"}}, Metadata{}, 234}
	cc2 := ClusterConstraint{"cckey2", "ckey2", Metadata{{"key", "value"}, {"k", "v"}}, Metadata{}, 0}
	ccs := ClusterConstraints{cc1, cc2}

	assert.NonNil(t, ccs.IsValid(true))
	assert.NonNil(t, ccs.IsValid(false))
}

func TestClusterConstraintsIsValidEmpty(t *testing.T) {
	ccs := ClusterConstraints{}

	assert.Nil(t, ccs.IsValid(true))
	assert.Nil(t, ccs.IsValid(false))
}

// ClusterConstraints.IsValid
func TestClusterConstraintsIsValidSucces(t *testing.T) {
	cc1 := ClusterConstraint{"cckey1", "ckey1", Metadata{{"key", "value"}, {"k", "v"}}, Metadata{}, 234}
	cc2 := ClusterConstraint{"cckey2", "ckey2", Metadata{{"key", "value"}, {"k", "v"}}, Metadata{}, 123}
	ccs := ClusterConstraints{cc1, cc2}

	set := AllConstraints{ccs, ccs, ccs}

	assert.Nil(t, set.IsValid(true))
	assert.Nil(t, set.IsValid(false))
}

func TestClusterConstraintsIsValidFailsWithBadLight(t *testing.T) {
	cc1 := ClusterConstraint{"cckey1", "ckey", Metadata{{"key", "value"}, {"k", "v"}}, Metadata{}, 234}
	cc2 := ClusterConstraint{"cckey2", "ckey2", Metadata{{"key", "value"}, {"k", "v"}}, Metadata{}, 123}
	ccbad := ClusterConstraint{"cc-bad", "ckey2", Metadata{{"key", "value"}, {"k", "v"}}, Metadata{}, 0}

	ccs := ClusterConstraints{cc1, cc2}
	ccsBad := ClusterConstraints{cc1, cc2, ccbad}

	set := AllConstraints{ccs, ccs, ccs}
	set.Light = ccsBad

	assert.NonNil(t, set.IsValid(true))
	assert.NonNil(t, set.IsValid(false))
}

func TestClusterConstraintsIsValidFailsWithBadDark(t *testing.T) {
	cc1 := ClusterConstraint{"cckey1", "ckey", Metadata{{"key", "value"}, {"k", "v"}}, Metadata{}, 234}
	cc2 := ClusterConstraint{"cckey2", "ckey2", Metadata{{"key", "value"}, {"k", "v"}}, Metadata{}, 123}
	ccbad := ClusterConstraint{"cc-bad", "ckey2", Metadata{{"key", "value"}, {"k", "v"}}, Metadata{}, 0}

	ccs := ClusterConstraints{cc1, cc2}
	ccsBad := ClusterConstraints{cc1, cc2, ccbad}

	set := AllConstraints{ccs, ccs, ccs}
	set.Dark = ccsBad

	assert.NonNil(t, set.IsValid(true))
	assert.NonNil(t, set.IsValid(false))
}

func TestClusterConstraintsIsValidFailsWithBadTap(t *testing.T) {
	cc1 := ClusterConstraint{"cckey1", "ckey", Metadata{{"key", "value"}, {"k", "v"}}, Metadata{}, 234}
	cc2 := ClusterConstraint{"cckey2", "ckey2", Metadata{{"key", "value"}, {"k", "v"}}, Metadata{}, 123}
	ccbad := ClusterConstraint{"cc-bad", "ckey2", Metadata{{"key", "value"}, {"k", "v"}}, Metadata{}, 0}

	ccs := ClusterConstraints{cc1, cc2}
	ccsBad := ClusterConstraints{cc1, cc2, ccbad}

	set := AllConstraints{ccs, ccs, ccs}
	set.Tap = ccsBad

	assert.NonNil(t, set.IsValid(true))
	assert.NonNil(t, set.IsValid(false))
}

func TestClusterConstraintsIsValidFailsWithZeroLight(t *testing.T) {
	cc1 := ClusterConstraint{"cckey1", "ckey", Metadata{{"key", "value"}, {"k", "v"}}, Metadata{}, 234}
	cc2 := ClusterConstraint{"cckey2", "ckey2", Metadata{{"key", "value"}, {"k", "v"}}, Metadata{}, 123}
	ccs := ClusterConstraints{cc1, cc2}

	set := AllConstraints{ccs, ccs, ccs}
	set.Light = ClusterConstraints{}

	assert.NonNil(t, set.IsValid(true))
	assert.NonNil(t, set.IsValid(false))
}

func TestConstraintMetadataValid(t *testing.T) {
	mgood := Metadatum{"key", "value"}
	m := Metadata{mgood, mgood, mgood}
	assert.Nil(t, ConstraintMetadataValid(m))
}

func TestConstraintMetadataValidFailed(t *testing.T) {
	mgood := Metadatum{"key", "value"}
	mbad1 := Metadatum{"key-2", ""}
	mbad2 := Metadatum{"", "value-2"}
	m := Metadata{mgood, mbad1, mbad2, mgood}
	errs := ConstraintMetadataValid(m)
	sort.Sort(ValidationErrorsByAttribute{errs})
	assert.DeepEqual(t, errs, &ValidationError{[]ErrorCase{
		{"metadatum[].key", "key must not be empty"},
		{"metadatum[key-2].value", "value must not be empty"},
	}})
}

func TestConstraintPropertiesValid(t *testing.T) {
	mgood := Metadatum{"key", "value"}
	mgood2 := Metadatum{"key", ""}
	m := Metadata{mgood, mgood2}
	assert.Nil(t, ConstraintPropertiesValid(m))
}

func TestConstraintPropertiesValidFailed(t *testing.T) {
	mgood := Metadatum{"key", "value"}
	mbad := Metadatum{"", "value"}
	m := Metadata{mbad, mgood}
	errs := ConstraintPropertiesValid(m)
	assert.DeepEqual(t, errs, &ValidationError{[]ErrorCase{
		{"property[].key", "key must not be empty"},
	}})
}
