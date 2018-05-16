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
	"testing"

	"github.com/turbinelabs/nonstdlib/ptr"
	"github.com/turbinelabs/test/assert"
)

func getClusters() (Cluster, Cluster) {
	ia := Instance{"Host", 1234, nil}
	ib := Instance{"Host2", 1234, nil}
	ic := Instance{"Host3", 1234, nil}
	i := Instances{ia, ib, ic}
	cb := CircuitBreakers{ptr.Int(1), ptr.Int(2), ptr.Int(3), ptr.Int(4)}
	od := OutlierDetection{
		IntervalMsec:                       ptr.Int(1),
		BaseEjectionTimeMsec:               ptr.Int(2),
		MaxEjectionPercent:                 ptr.Int(3),
		Consecutive5xx:                     ptr.Int(4),
		EnforcingConsecutive5xx:            ptr.Int(5),
		EnforcingSuccessRate:               ptr.Int(6),
		SuccessRateMinimumHosts:            ptr.Int(7),
		SuccessRateRequestVolume:           ptr.Int(8),
		SuccessRateStdevFactor:             ptr.Int(9),
		ConsecutiveGatewayFailure:          ptr.Int(10),
		EnforcingConsecutiveGatewayFailure: ptr.Int(11),
	}

	c := Cluster{"ckey", "zkey", "name", true, i, "okey1", &cb, &od, Checksum{}}
	return c, c
}

func TestClusterEqualsOrgVaries(t *testing.T) {
	c1, c2 := getClusters()
	c2.OrgKey = "okey2"

	assert.False(t, c2.Equals(c1))
	assert.False(t, c1.Equals(c2))
}

func TestClusterEqualsSuccess(t *testing.T) {
	c1, c2 := getClusters()

	assert.True(t, c2.Equals(c1))
	assert.True(t, c1.Equals(c2))
}

func TestClusterEqualsTLSTrueFalse(t *testing.T) {
	c1, c2 := getClusters()
	c1.RequireTLS = true
	c2.RequireTLS = false

	assert.False(t, c1.Equals(c2))
	assert.False(t, c2.Equals(c1))
}

func TestClusterKeyVaries(t *testing.T) {
	c1, c2 := getClusters()
	c2.ClusterKey = "ckey2"

	assert.False(t, c1.Equals(c2))
	assert.False(t, c2.Equals(c1))
}

func TestClusterZoneKeyVaries(t *testing.T) {
	c1, c2 := getClusters()
	c2.ZoneKey = "zkey2"

	assert.False(t, c1.Equals(c2))
	assert.False(t, c2.Equals(c1))
}

func TestClusterNameVaries(t *testing.T) {
	c1, c2 := getClusters()
	c2.Name = "name2"

	assert.False(t, c1.Equals(c2))
	assert.False(t, c2.Equals(c1))
}

func TestClusterInstancesZeroNil(t *testing.T) {
	c1, c2 := getClusters()
	c1.Instances = nil
	c2.Instances = Instances{}

	assert.True(t, c1.Equals(c2))
	assert.True(t, c2.Equals(c1))
}

func TestClusterInstancesVaries(t *testing.T) {
	c1, c2 := getClusters()
	c2.Instances = c1.Instances[1:]

	assert.False(t, c1.Equals(c2))
	assert.False(t, c2.Equals(c1))
}

func TestClusterChecksumVaries(t *testing.T) {
	c1, c2 := getClusters()
	c2.Checksum = Checksum{"csum2"}

	assert.False(t, c1.Equals(c2))
	assert.False(t, c2.Equals(c1))
}

func TestClusterEqualInstanceOrderVaries(t *testing.T) {
	c1, c2 := getClusters()
	i := c1.Instances
	c2.Instances = Instances{i[2], i[0], i[1]}

	assert.True(t, c1.Equals(c2))
	assert.True(t, c2.Equals(c1))
}

func TestClusterEqualCircuitBreakersVaries(t *testing.T) {
	c1, c2 := getClusters()
	c2.CircuitBreakers = &CircuitBreakers{ptr.Int(4), ptr.Int(3), ptr.Int(2), ptr.Int(1)}

	assert.False(t, c1.Equals(c2))
	assert.False(t, c2.Equals(c1))
}

func TestClusterEqualOutlierDetectionVaries(t *testing.T) {
	c1, c2 := getClusters()
	c2.OutlierDetection = &OutlierDetection{EnforcingSuccessRate: ptr.Int(100)}

	assert.False(t, c1.Equals(c2))
	assert.False(t, c2.Equals(c1))
}

func TestClustersGroupBy(t *testing.T) {
	c := Clusters{
		Cluster{Name: "a", ClusterKey: "a"},
		Cluster{Name: "b", ClusterKey: "b"},
		Cluster{Name: "c", ClusterKey: "c"},
		Cluster{Name: "d", ClusterKey: "d"},
	}

	want := map[string]Clusters{
		"a": {{Name: "a", ClusterKey: "a"}},
		"b": {{Name: "b", ClusterKey: "b"}},
		"c": {{Name: "c", ClusterKey: "c"}},
		"d": {{Name: "d", ClusterKey: "d"}},
	}

	assert.DeepEqual(t, c.GroupBy(func(c Cluster) string { return c.Name }), want)
}

func TestClustersGroupByCollection(t *testing.T) {
	c := Clusters{
		Cluster{Name: "a", ClusterKey: "a"},
		Cluster{Name: "b", ClusterKey: "b"},
		Cluster{Name: "c", ClusterKey: "c"},
		Cluster{Name: "d", ClusterKey: "d"},
	}

	want := map[string]Clusters{"a": c}

	assert.DeepEqual(t, c.GroupBy(func(c Cluster) string { return "a" }), want)
}

func mkTestC() *Cluster {
	return &Cluster{
		ClusterKey: "ck-1",
		Name:       "a cluster name",
		ZoneKey:    "zk-1",
		Instances: Instances{
			{"foo", 9090, MetadataFromMap(map[string]string{"key1": "value1"})},
			{"bar", 9090, MetadataFromMap(map[string]string{"key1": "value1", "key2": "value2"})},
		},
		OrgKey:           "ok-1",
		CircuitBreakers:  &CircuitBreakers{ptr.Int(1), ptr.Int(2), ptr.Int(3), ptr.Int(4)},
		OutlierDetection: &OutlierDetection{EnforcingConsecutive5xx: ptr.Int(100)},
		Checksum:         Checksum{"ck-1"},
	}
}

func TestClusterIsValid(t *testing.T) {
	assert.Nil(t, mkTestC().IsValid())
}

func TestClusterIsValidNoKey(t *testing.T) {
	c := mkTestC()
	c.ClusterKey = ""
	assert.NonNil(t, c.IsValid())
}

func TestClusterKeyIsValidBadKey(t *testing.T) {
	c := mkTestC()
	c.ClusterKey = "a-bad-key!"
	assert.NonNil(t, c.IsValid())
}

func TestClusterKeyIsValidBadZoneKey(t *testing.T) {
	c := mkTestC()
	c.ZoneKey = "a-bad-key!"
	assert.NonNil(t, c.IsValid())
}

func TestClusterIsValidBadName(t *testing.T) {
	c := mkTestC()
	c.Name = "aoeu[]',.p"
	assert.NonNil(t, c.IsValid())
}

func TestClusterIsValidBadInstances(t *testing.T) {
	c := mkTestC()
	c.Instances = append(c.Instances, c.Instances[0])
	assert.NonNil(t, c.IsValid())
}

func TestClusterIsValidBadCircuitBreakers(t *testing.T) {
	c := mkTestC()
	c.CircuitBreakers.MaxConnections = ptr.Int(-1)
	assert.NonNil(t, c.IsValid())
}

func TestClusterIsValidBadOutlierDetection(t *testing.T) {
	c := mkTestC()
	c.OutlierDetection.EnforcingSuccessRate = ptr.Int(-1)
	assert.NonNil(t, c.IsValid())
}
