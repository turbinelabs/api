package api

import (
	"testing"

	"github.com/turbinelabs/test/assert"
)

func getCohortSeed() (CohortSeed, CohortSeed) {
	return CohortSeed{CohortSeedHeader, "x-cohort-seed", true},
		CohortSeed{CohortSeedHeader, "x-cohort-seed", true}
}

func TestCohortSeedEquals(t *testing.T) {
	c1, c2 := getCohortSeed()
	assert.True(t, c1.Equals(c2))
	assert.True(t, c2.Equals(c1))
}

func TestCohortSeedEqualsNameChange(t *testing.T) {
	c1, c2 := getCohortSeed()
	c2.Name = "aosnetuh"
	assert.False(t, c1.Equals(c2))
	assert.False(t, c1.Equals(c2))
}

func TestCohortSeedEqualsTypeChange(t *testing.T) {
	c1, c2 := getCohortSeed()
	c2.Type = "aosnetuh"
	assert.False(t, c1.Equals(c2))
	assert.False(t, c1.Equals(c2))
}

func TestCohortSeedEqualsUseZeroValueSeedChange(t *testing.T) {
	c1, c2 := getCohortSeed()
	c2.UseZeroValueSeed = !c1.UseZeroValueSeed
	assert.False(t, c1.Equals(c2))
	assert.False(t, c1.Equals(c2))
}

func TestCohortSeedIsValid(t *testing.T) {
	c, _ := getCohortSeed()
	assert.Nil(t, c.IsValid())
}

func TestCohortSeedIsValidBadType(t *testing.T) {
	c, _ := getCohortSeed()
	c.Type = "WHEE"
	assert.DeepEqual(t, c.IsValid(), &ValidationError{[]ErrorCase{
		{"cohort_seed.type", `"WHEE" is not a valid seed type`},
	}})
}

func TestCohortSeedIsValidBadName(t *testing.T) {
	c, _ := getCohortSeed()
	c.Name = ""
	assert.DeepEqual(t, c.IsValid(), &ValidationError{[]ErrorCase{
		{"cohort_seed.name", "may not be empty"},
	}})
}
