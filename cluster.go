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
	"sort"
)

// Clusters is a slice of Cluster objects
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
	ClusterKey       ClusterKey        `json:"cluster_key"` // overwritten on create
	ZoneKey          ZoneKey           `json:"zone_key"`
	Name             string            `json:"name"`
	RequireTLS       bool              `json:"require_tls,omitempty"`
	Instances        Instances         `json:"instances"`
	OrgKey           OrgKey            `json:"-"`
	CircuitBreakers  *CircuitBreakers  `json:"circuit_breakers"`
	OutlierDetection *OutlierDetection `json:"outlier_detection"`
	HealthChecks     HealthChecks      `json:"health_checks"`
	Checksum
}

func (o Cluster) GetZoneKey() ZoneKey   { return o.ZoneKey }
func (o Cluster) GetOrgKey() OrgKey     { return o.OrgKey }
func (o Cluster) Key() string           { return string(o.ClusterKey) }
func (o Cluster) GetChecksum() Checksum { return o.Checksum }

func (c Cluster) IsNil() bool {
	return c.Equals(Cluster{})
}

func (c Cluster) Equals(o Cluster) bool {
	coreResp := c.ClusterKey == o.ClusterKey &&
		c.ZoneKey == o.ZoneKey &&
		c.Name == o.Name &&
		c.OrgKey == o.OrgKey &&
		c.RequireTLS == o.RequireTLS &&
		CircuitBreakersPtrEquals(c.CircuitBreakers, o.CircuitBreakers) &&
		OutlierDetectionPtrEquals(c.OutlierDetection, o.OutlierDetection) &&
		c.HealthChecks.Equals(o.HealthChecks) &&
		c.Checksum.Equals(o.Checksum)

	if !coreResp {
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
func (c *Cluster) IsValid() *ValidationError {
	scope := func(i string) string { return "cluster." + i }

	errs := &ValidationError{}

	errCheckKey(string(c.ClusterKey), errs, scope("cluster_key"))
	errCheckKey(string(c.ZoneKey), errs, scope("zone_key"))
	errCheckIndex(c.Name, errs, scope("name"))

	errs.MergePrefixed(c.Instances.IsValid(), "cluster")
	if c.CircuitBreakers != nil {
		errs.MergePrefixed(c.CircuitBreakers.IsValid(), "cluster")
	}

	if c.OutlierDetection != nil {
		errs.MergePrefixed(c.OutlierDetection.IsValid(), "cluster")
	}

	errs.MergePrefixed(c.HealthChecks.IsValid(), "cluster")

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
