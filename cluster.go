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
	"sort"
)

type Clusters []Cluster

func (c Clusters) GroupBy(fn func(Cluster) string) map[string]Clusters {
	if c == nil {
		return nil
	}

	result := map[string]Clusters{}

	for _, cl := range c {
		k := fn(cl)
		result[k] = append(result[k], cl)
	}

	return result
}

type ClusterKey string

// A Cluster is a named list of Instances within a zone
type Cluster struct {
	ClusterKey ClusterKey `json:"cluster_key"` // overwritten on create
	ZoneKey    ZoneKey    `json:"zone_key"`
	Name       string     `json:"name"`
	Instances  Instances  `json:"instances"`
	OrgKey     OrgKey     `json:"-"`
	Checksum
}

func (c Cluster) IsNil() bool {
	return c.Equals(Cluster{})
}

func (c Cluster) Equals(o Cluster) bool {
	coreResp := c.ClusterKey == o.ClusterKey &&
		c.ZoneKey == o.ZoneKey &&
		c.Name == o.Name &&
		c.OrgKey == o.OrgKey

	if !coreResp {
		return false
	}

	if !c.Checksum.Equals(o.Checksum) {
		return false
	}

	// special case this to treat [] and nil as equal values
	if c.Instances == nil || len(c.Instances) == 0 {
		return o.Instances == nil || len(o.Instances) == 0
	}

	return c.Instances.Equals(o.Instances)
}

// Checks the data set on a Cluster and returns whether or not sufficient
// information is available.
//
// If this is being checked before cluster creation (as indicated by the
// precreation param) then a cluster key is not required.
func (c *Cluster) IsValid(precreation bool) *ValidationError {
	ecase := func(f, m string) ErrorCase {
		return ErrorCase{fmt.Sprintf("cluster[%s].%s", string(c.ClusterKey), f), m}
	}

	errs := &ValidationError{}

	validClusterKey := c.ClusterKey != "" || precreation
	validZoneKey := c.ZoneKey != ""
	validName := c.Name != ""

	if !validClusterKey {
		errs.AddNew(ecase("cluster_key", "must not be empty"))
	}

	if !validZoneKey {
		errs.AddNew(ecase("zone_key", "must not be empty"))
	}

	if !validName {
		errs.AddNew(ecase("name", "must not be empty"))
	}

	errs.MergePrefixed(
		c.Instances.IsValid(precreation),
		fmt.Sprintf("cluster[%s].instances", string(c.ClusterKey)))

	return errs.OrNil()
}

// Sort a slice of Clusters by ClusterKey.
// Eg: sort.Sort(ClusterByClusterKey(clusters))
type ClusterByClusterKey []Cluster

var _ sort.Interface = ClusterByClusterKey{}

func (b ClusterByClusterKey) Len() int           { return len(b) }
func (b ClusterByClusterKey) Swap(i, j int)      { b[i], b[j] = b[j], b[i] }
func (b ClusterByClusterKey) Less(i, j int) bool { return b[i].ClusterKey < b[j].ClusterKey }

// Sort a slice of Clusters by Name.
// Eg: sort.Sort(ClusterByName(clusters))
type ClusterByName []Cluster

var _ sort.Interface = ClusterByName{}

func (b ClusterByName) Len() int           { return len(b) }
func (b ClusterByName) Swap(i, j int)      { b[i], b[j] = b[j], b[i] }
func (b ClusterByName) Less(i, j int) bool { return b[i].Name < b[j].Name }
