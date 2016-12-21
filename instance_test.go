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
