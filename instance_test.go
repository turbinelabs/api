package api

import (
	"testing"

	"github.com/turbinelabs/test/assert"
)

// Instance
func TestInstanceMetadataVaries(t *testing.T) {
	ma := Metadatum{"Key", "Value"}
	mb := Metadatum{"Key2", "Value"}
	i1 := Instance{"Host", 1234, Metadata{ma, mb}}
	i2 := Instance{"Host", 1234, Metadata{mb}}

	assert.False(t, i1.Equals(i2))
	assert.False(t, i2.Equals(i1))
	assert.False(t, i1.Equivalent(i2))
	assert.False(t, i2.Equivalent(i1))
}

func TestInstanceMetadataZeroNil(t *testing.T) {
	i1 := Instance{"Host", 1234, Metadata{}}
	i2 := Instance{"Host", 1234, nil}

	assert.True(t, i2.Equals(i1))
	assert.True(t, i1.Equals(i2))
	assert.True(t, i2.Equivalent(i1))
	assert.True(t, i1.Equivalent(i2))
}

func TestHostVaries(t *testing.T) {
	i1 := Instance{"Host", 1234, nil}
	i2 := Instance{"Host2", 1234, nil}

	assert.False(t, i2.Equals(i1))
	assert.False(t, i1.Equals(i2))
	assert.False(t, i2.Equivalent(i1))
	assert.False(t, i1.Equivalent(i2))
}

func TestPortVaries(t *testing.T) {
	i1 := Instance{"Host", 1234, nil}
	i2 := Instance{"Host", 1235, nil}

	assert.False(t, i2.Equals(i1))
	assert.False(t, i1.Equals(i2))
	assert.False(t, i2.Equivalent(i1))
	assert.False(t, i1.Equivalent(i2))
}

func TestInstanceMatches(t *testing.T) {
	ma := Metadatum{"Key", "Value"}
	mb := Metadatum{"Key2", "Value"}
	i1 := Instance{"Host", 1234, Metadata{ma, mb}}
	i2 := Instance{"Host", 1234, Metadata{mb, ma}}

	assert.True(t, i2.Equals(i1))
	assert.True(t, i1.Equals(i1))
	assert.True(t, i2.Equivalent(i1))
	assert.True(t, i1.Equivalent(i1))
}

// Instances
func TestInstancesZeroNil(t *testing.T) {
	i1 := Instances{}
	var i2 Instances

	assert.True(t, i1.Equals(i2))
	assert.True(t, i2.Equals(i1))
	assert.True(t, i1.Equivalent(i2))
	assert.True(t, i2.Equivalent(i1))
}

func TestInstancesZeroZero(t *testing.T) {
	i1 := Instances{}
	i2 := Instances{}

	assert.True(t, i1.Equals(i2))
	assert.True(t, i2.Equals(i1))
	assert.True(t, i1.Equivalent(i2))
	assert.True(t, i2.Equivalent(i1))
}

func TestInstancesOutOfOrder(t *testing.T) {
	ia := Instance{"Host", 8080, nil}
	ib := Instance{"Host2", 80, nil}
	i1 := Instances{ia, ib}
	i2 := Instances{ib, ia}

	assert.True(t, i1.Equals(i2))
	assert.True(t, i2.Equals(i1))
	assert.True(t, i1.Equivalent(i2))
	assert.True(t, i2.Equivalent(i1))
}

func TestInstancesExtraElement(t *testing.T) {
	ia := Instance{"Host", 8080, nil}
	ib := Instance{"Host2", 80, nil}
	ic := Instance{"Host3", 8081, nil}
	i1 := Instances{ia, ib, ic}
	i2 := Instances{ib, ia}

	assert.False(t, i1.Equals(i2))
	assert.False(t, i2.Equals(i1))
	assert.False(t, i1.Equivalent(i2))
	assert.False(t, i2.Equivalent(i1))
}
