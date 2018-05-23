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

import (
	"errors"
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
}

func TestInstanceMetadataZeroNil(t *testing.T) {
	i1 := Instance{"Host", 1234, Metadata{}}
	i2 := Instance{"Host", 1234, nil}

	assert.True(t, i2.Equals(i1))
	assert.True(t, i1.Equals(i2))
}

func TestHostVaries(t *testing.T) {
	i1 := Instance{"Host", 1234, nil}
	i2 := Instance{"Host2", 1234, nil}

	assert.False(t, i2.Equals(i1))
	assert.False(t, i1.Equals(i2))
}

func TestPortVaries(t *testing.T) {
	i1 := Instance{"Host", 1234, nil}
	i2 := Instance{"Host", 1235, nil}

	assert.False(t, i2.Equals(i1))
	assert.False(t, i1.Equals(i2))
}

func TestInstanceMatches(t *testing.T) {
	ma := Metadatum{"Key", "Value"}
	mb := Metadatum{"Key2", "Value"}
	i1 := Instance{"Host", 1234, Metadata{ma, mb}}
	i2 := Instance{"Host", 1234, Metadata{mb, ma}}

	assert.True(t, i2.Equals(i1))
	assert.True(t, i1.Equals(i1))
}

// Instances
func TestInstancesZeroNil(t *testing.T) {
	i1 := Instances{}
	var i2 Instances

	assert.True(t, i1.Equals(i2))
	assert.True(t, i2.Equals(i1))
}

func TestInstancesZeroZero(t *testing.T) {
	i1 := Instances{}
	i2 := Instances{}

	assert.True(t, i1.Equals(i2))
	assert.True(t, i2.Equals(i1))
}

func TestInstancesOutOfOrder(t *testing.T) {
	ia := Instance{"Host", 8080, nil}
	ib := Instance{"Host2", 80, nil}
	i1 := Instances{ia, ib}
	i2 := Instances{ib, ia}

	assert.True(t, i1.Equals(i2))
	assert.True(t, i2.Equals(i1))
}

func TestInstancesExtraElement(t *testing.T) {
	ia := Instance{"Host", 8080, nil}
	ib := Instance{"Host2", 80, nil}
	ic := Instance{"Host3", 8081, nil}
	i1 := Instances{ia, ib, ic}
	i2 := Instances{ib, ia}

	assert.False(t, i1.Equals(i2))
	assert.False(t, i2.Equals(i1))
}

func mkTestI() Instance {
	return Instance{
		Host: "host-name",
		Port: 30080,
		Metadata: MetadataFromMap(map[string]string{
			"key":  "value",
			"key2": "value",
		}),
	}
}

func TestInstanceIsValid(t *testing.T) {
	assert.Nil(t, mkTestI().IsValid())
}

func TestInstanceIsValidBadHost(t *testing.T) {
	i := mkTestI()
	i.Host = "some-bad-!"
	assert.NonNil(t, i.IsValid())
}

func TestInstanceIsValidNoHost(t *testing.T) {
	i := mkTestI()
	i.Host = ""
	assert.NonNil(t, i.IsValid())
}

func TestInstanceIsValidBadPort(t *testing.T) {
	i := mkTestI()
	i.Port = 0
	assert.NonNil(t, i.IsValid())
}

func TestInstanceIsValidBadMetadata(t *testing.T) {
	i := mkTestI()
	i.Metadata = append(i.Metadata, i.Metadata[0])
	assert.NonNil(t, i.IsValid())
}

func mkTestIMD() Metadata {
	return Metadata{
		{"key", "value"},
		{"key2", "value"},
	}
}

func TestInstanceMetadataIsValid(t *testing.T) {
	im := mkTestIMD()
	assert.Nil(t, InstanceMetadataIsValid(im))
}

func TestInstanceMetadataIsValidBadKey(t *testing.T) {
	im := mkTestIMD()
	im[0].Key = "aoeu]"
	assert.NonNil(t, InstanceMetadataIsValid(im))
}

func TestInstanceMetadataIsValidDupes(t *testing.T) {
	im := mkTestIMD()
	im = append(im, im[0])
	assert.NonNil(t, InstanceMetadataIsValid(im))
}

var selectFixture = Instances{
	{Host: "host", Port: 0},
	{Host: "host", Port: 1},
	{Host: "host", Port: 2},
	{Host: "host", Port: 3},
}

func doTestSelect(
	t *testing.T,
	fn func(Instance) (bool, error),
	want Instances,
	wantErr error,
) {
	got, gotErr := selectFixture.Select(fn)
	assert.DeepEqual(t, got, want)
	assert.DeepEqual(t, gotErr, wantErr)
}

func TestSelectNone(t *testing.T) {
	doTestSelect(
		t,
		func(i Instance) (bool, error) {
			return false, nil
		},
		Instances{},
		nil,
	)
}

func TestSelectSome(t *testing.T) {
	doTestSelect(
		t,
		func(i Instance) (bool, error) {
			return i.Port == 1 || i.Port == 3, nil
		},
		Instances{selectFixture[1], selectFixture[3]},
		nil,
	)
}

func TestSelectOrror(t *testing.T) {
	e := errors.New("aoesntuh")

	doTestSelect(
		t,
		func(i Instance) (bool, error) {
			return i.Port == 1 || i.Port == 3, e
		},
		nil,
		e,
	)
}

func doTestMatchesMetadata(
	t *testing.T,
	mdMap map[string]string,
	meets bool,
) {
	inst := Instance{
		Metadata: MetadataFromMap(map[string]string{
			"k1": "v1",
			"k2": "v2",
			"k3": "v3",
		}),
	}

	var md Metadata = nil
	if mdMap != nil {
		md = MetadataFromMap(mdMap)
	}

	assert.Equal(t, inst.MatchesMetadata(md), meets)
}

func TestMatchesMetadataNil(t *testing.T) {
	doTestMatchesMetadata(
		t,
		nil,
		true,
	)
}

func TestMatchesMetadataEmpty(t *testing.T) {
	doTestMatchesMetadata(
		t,
		map[string]string{},
		true,
	)
}

func TestMatchesMetadataFalse(t *testing.T) {
	doTestMatchesMetadata(
		t,
		map[string]string{
			"aoeu": "snth",
		},
		false,
	)
}

func TestMatchesMetadataSubsetTrue(t *testing.T) {
	doTestMatchesMetadata(
		t,
		map[string]string{
			"k2": "v2",
		},
		true,
	)
}

func TestMatchesMetadataSubsetFalse(t *testing.T) {
	doTestMatchesMetadata(
		t,
		map[string]string{
			"k1": "v1",
			"k2": "v2-different",
		},
		false,
	)
}

func TestMatchesMetadataSupersetFalse(t *testing.T) {
	doTestMatchesMetadata(
		t,
		map[string]string{
			"k1": "v1",
			"k2": "v2",
			"k3": "v3",
			"k4": "v4",
		},
		false,
	)
}

func TestMeetsConstraintKeyCaseSensitivityFalse(t *testing.T) {
	doTestMatchesMetadata(
		t,
		map[string]string{
			"K2": "v2",
		},
		false,
	)
}

func TestMatchesMetadataValueCaseSensitivityFalse(t *testing.T) {
	doTestMatchesMetadata(
		t,
		map[string]string{
			"k2": "V2",
		},
		false,
	)
}
