package api

import (
	"testing"

	"github.com/turbinelabs/test/assert"
)

func getClusters() (Cluster, Cluster) {
	ia := Instance{"Host", 1234, nil}
	ib := Instance{"Host2", 1234, nil}
	ic := Instance{"Host3", 1234, nil}
	i := Instances{ia, ib, ic}
	c := Cluster{"ckey", "zkey", "name", i, "okey1", Checksum{}}
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
