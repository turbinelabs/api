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
	"fmt"
	"sort"

	"github.com/turbinelabs/nonstdlib/arrays"
	"github.com/turbinelabs/nonstdlib/ptr"
)

// HealthCheck configures the parameters to do health checking against instances
// in a cluster.
type HealthCheck struct {
	// TimeoutMsec is the time to wait for a health check response. If the
	// timeout is reached without a response, the health check attempt will
	// be considered a failure. This is a required field and must be greater
	// than 0.
	TimeoutMsec int `json:"timeout_msec"`

	// IntervalMsec is the interval between health checks. Note that the
	// first round of health checks will occur during startup before any
	// traffic is routed to a cluster. This means that the
	// `NoTrafficIntervalMsec` value will be used as the first interval of
	// health checks.
	IntervalMsec int `json:"interval_msec"`

	// IntervalJitterMsec is an optional jitter amount that is added to each
	// interval value calculated by the proxy. If not specified, defaults
	// to 0.
	IntervalJitterMsec *int `json:"interval_jitter_msec,omitempty"`

	// UnhealthyThreshold is the number of unhealthy health checks required
	// before a host is marked unhealthy. Note that for *http* health
	// checking if a host responds with 503 this threshold is ignored and
	// the host is considered unhealthy immediately.
	UnhealthyThreshold int `json:"unhealthy_threshold"`

	// HealthyThreshold is the number of healthy health checks required
	// before a host is marked healthy. Note that during startup, only a
	// single successful health check is required to mark a host healthy.
	HealthyThreshold int `json:"healthy_threshold"`

	// ReuseConnection determines whether to reuse a health check connection
	// between health checks. Default is true.
	ReuseConnection *bool `json:"reuse_connection,omitempty"`

	// NoTrafficIntervalMsec is a special health check interval that is
	// used when a cluster has never had traffic routed to it. This lower
	// interval allows cluster information to be kept up to date, without
	// sending a potentially large amount of active health checking traffic
	// for no reason. Once a cluster has been used for traffic routing,
	// The proxy will shift back to using the standard health check interval
	// that is defined. Note that this interval takes precedence over any
	// other. Defaults to 60s.
	NoTrafficIntervalMsec *int `json:"no_traffic_interval_msec,omitempty"`

	// UnhealthyIntervalMsec is a health check interval that is used for
	// hosts that are marked as unhealthy. As soon as the host is marked as
	// healthy, the proxy will shift back to using the standard health check
	// interval that is defined. This defaults to the same value as
	// IntervalMsec if not specified.
	UnhealthyIntervalMsec *int `json:"unhealthy_interval_msec,omitempty"`

	// UnhealthyEdgeIntervalMsec is a special health check interval that
	// is used for the first health check right after a host is marked as
	// unhealthy. For subsequent health checks the proxy will shift back to
	// using either "unhealthy interval" if present or the standard
	// health check interval that is defined. Defaults to the same value as
	// UnhealthIntervalMsec if not specified.
	UnhealthyEdgeIntervalMsec *int `json:"unhealthy_edge_interval_msec,omitempty"`

	// HealthyEdgeIntervalMsec is a special health check interval that is
	// used for the first health check right after a host is marked as
	// healthy. For subsequent health checks the proxy will shift back to
	// using the standard health check interval that is defined. Defaults
	// to the same value as IntervalMsec if not specified
	HealthyEdgeIntervalMsec *int `json:"health_edge_interval_msec,omitempty"`

	// HealthChecker defines the type of health checking to use.
	HealthChecker HealthChecker `json:"health_checker"`
}

// HealthChecks is a slice of HealthCheck objects. Currently, the proxy only
// supports a single health check per cluster
type HealthChecks []HealthCheck

// HealthChecksByType implements sort.Interface to allow sorting by
// health check type
type HealthChecksByType HealthChecks

var _ sort.Interface = HealthChecksByType{}

func (h HealthChecksByType) Len() int      { return len(h) }
func (h HealthChecksByType) Swap(i, j int) { h[i], h[j] = h[j], h[i] }

func (h HealthChecksByType) Less(i, j int) bool {
	return h[i].compare(h[j]) == -1
}

// Equals checks two HealthChecks for equality.
func (hcs HealthChecks) Equals(hcs2 HealthChecks) bool {
	if len(hcs) != len(hcs2) {
		return false
	}

	sort.Sort(HealthChecksByType(hcs))
	sort.Sort(HealthChecksByType(hcs2))

	for i := range hcs {
		if !hcs[i].Equals(hcs2[i]) {
			return false
		}
	}

	return true
}

// IsValid confirms a HealthChecks instance is valid.
func (hcs HealthChecks) IsValid() *ValidationError {
	errs := &ValidationError{}
	if len(hcs) > 1 {
		errs.AddNew(
			ErrorCase{
				"health_checks",
				"only a single health check supported",
			},
		)
	}

	for i, hc := range hcs {
		errs.MergePrefixed(hc.IsValid(), fmt.Sprintf("health_checks[%d]", i))
	}

	return errs.OrNil()
}

// Equals compares two HealthChecks for equality
func (hc HealthCheck) Equals(o HealthCheck) bool {
	return hc.compare(o) == 0
}

func (hc HealthCheck) compare(ohc HealthCheck) int {
	if cmp := compareInts(hc.TimeoutMsec, ohc.TimeoutMsec); cmp != 0 {
		return cmp
	}

	if cmp := compareInts(hc.IntervalMsec, ohc.IntervalMsec); cmp != 0 {
		return cmp
	}

	if cmp := ptr.CompareInts(hc.IntervalJitterMsec, ohc.IntervalJitterMsec); cmp != 0 {
		return cmp
	}

	if cmp := ptr.CompareInts(hc.NoTrafficIntervalMsec, ohc.NoTrafficIntervalMsec); cmp != 0 {
		return cmp
	}

	if cmp := ptr.CompareInts(hc.UnhealthyIntervalMsec, ohc.UnhealthyIntervalMsec); cmp != 0 {
		return cmp
	}

	if cmp := ptr.CompareInts(hc.UnhealthyEdgeIntervalMsec, ohc.UnhealthyEdgeIntervalMsec); cmp != 0 {
		return cmp
	}

	if cmp := ptr.CompareInts(hc.HealthyEdgeIntervalMsec, ohc.HealthyEdgeIntervalMsec); cmp != 0 {
		return cmp
	}

	if cmp := compareInts(hc.UnhealthyThreshold, ohc.UnhealthyThreshold); cmp != 0 {
		return cmp
	}

	if cmp := compareInts(hc.HealthyThreshold, ohc.HealthyThreshold); cmp != 0 {
		return cmp
	}

	if cmp := ptr.CompareBools(hc.ReuseConnection, ohc.ReuseConnection); cmp != 0 {
		return cmp
	}

	return hc.HealthChecker.compare(ohc.HealthChecker)
}

// IsValid checks that a HealthCheck object is valid.
func (hc HealthCheck) IsValid() *ValidationError {
	errs := &ValidationError{}
	if hc.TimeoutMsec < 1 {
		errs.AddNew(ErrorCase{"timeout_msec", "must be greater than zero"})
	}

	if hc.IntervalMsec < 1 {
		errs.AddNew(ErrorCase{"interval_msec", "must be greater than zero"})
	}

	if v, ok := ptr.IntValueOk(hc.IntervalJitterMsec); ok && v < 1 {
		errs.AddNew(
			ErrorCase{
				"interval_jitter_msec",
				"must be greater than zero",
			},
		)
	}

	if hc.UnhealthyThreshold < 1 {
		errs.AddNew(ErrorCase{"unhealthy_threshold", "must be greater than zero"})
	}

	if hc.HealthyThreshold < 1 {
		errs.AddNew(ErrorCase{"healthy_threshold", "must be greater than zero"})
	}

	if v, ok := ptr.IntValueOk(hc.NoTrafficIntervalMsec); ok && v < 1 {
		errs.AddNew(
			ErrorCase{
				"no_traffic_interval_msec",
				"must be greater than zero",
			},
		)
	}

	if v, ok := ptr.IntValueOk(hc.UnhealthyIntervalMsec); ok && v < 1 {
		errs.AddNew(
			ErrorCase{
				"unhealthy_interval_msec",
				"must be greater than zero",
			},
		)
	}

	if v, ok := ptr.IntValueOk(hc.UnhealthyEdgeIntervalMsec); ok && v < 1 {
		errs.AddNew(
			ErrorCase{
				"unhealthy_edge_interval_msec",
				"must be greater than zero",
			},
		)
	}

	if v, ok := ptr.IntValueOk(hc.HealthyEdgeIntervalMsec); ok && v < 1 {
		errs.AddNew(
			ErrorCase{
				"healthy_edge_interval_msec",
				"must be greater than zero",
			},
		)
	}

	errs.Merge(hc.HealthChecker.IsValid())
	return errs.OrNil()
}

// HealthChecker is a union type where only a single field can be defined.
type HealthChecker struct {
	// HTTPHealthCheck defines the parameters for http health checking.
	HTTPHealthCheck *HTTPHealthCheck `json:"http_health_check,omitempty"`

	// TCPHealthCheck defines the parameters for tcp health checking.
	TCPHealthCheck *TCPHealthCheck `json:"tcp_health_check,omitempty"`
}

// Equals checks two HealthChecker objects for equality
func (hc HealthChecker) Equals(ohc HealthChecker) bool {
	return hc.compare(ohc) == 0
}

func (hc HealthChecker) compare(ohc HealthChecker) int {
	if cmp := hc.HTTPHealthCheck.compare(ohc.HTTPHealthCheck); cmp != 0 {
		return cmp
	}

	return hc.TCPHealthCheck.compare(ohc.TCPHealthCheck)
}

// IsValid checks a HealthChecker object for validity.
func (hc HealthChecker) IsValid() *ValidationError {
	errs := &ValidationError{}

	switch {
	case hc.HTTPHealthCheck == nil && hc.TCPHealthCheck == nil:
		errs.AddNew(
			ErrorCase{
				"health_checker",
				"must have one health check defined",
			},
		)

	case hc.HTTPHealthCheck != nil && hc.TCPHealthCheck != nil:
		errs.AddNew(
			ErrorCase{
				"health_checker",
				"must not have more than one type of health check defined",
			},
		)
	}

	return errs.OrNil()
}

// HTTPHealthCheck configures the http health check endpoint for a cluster.
type HTTPHealthCheck struct {
	// Host defines the value of the host header in the HTTP health check
	// request. If left empty (default value), the name of the cluster being
	// health checked will be used.
	Host string `json:"host"`

	// Path specifies the HTTP path that will be requested during health
	// checking.
	Path string `json:"path"`

	// ServiceName is an optional service name parameter which is used to
	// validate the identity of the health checked cluster. This is done by
	// comparing the `X-Envoy-Upstream-Healthchecked-Cluster` header to
	// this value.
	ServiceName string `json:"service_name"`

	// RequestHeadersToAdd specifies a list of HTTP headers that should be
	// added to each request that is sent to the health checked cluster.
	RequestHeadersToAdd Metadata `json:"request_headers_to_add,omitempty"`
}

// Equals checks for equality between two HTTPHealthCheck pointers
func (hhc *HTTPHealthCheck) Equals(ohc *HTTPHealthCheck) bool {
	return hhc.compare(ohc) == 0
}

// Treats nil as being less than defined pointers.
func (hhc *HTTPHealthCheck) compare(ohc *HTTPHealthCheck) int {
	switch {
	case hhc == nil && ohc == nil:
		return 0
	case hhc == nil && ohc != nil:
		return -1
	case hhc != nil && ohc == nil:
		return 1
	case hhc != nil && ohc != nil:
		if hhc.Host < ohc.Host {
			return -1
		}
		if hhc.Host > ohc.Host {
			return 1
		}

		if hhc.Path < ohc.Path {
			return -1
		}

		if hhc.Path > ohc.Path {
			return 1
		}

		if hhc.ServiceName < ohc.ServiceName {
			return -1
		}

		if hhc.ServiceName > ohc.ServiceName {
			return 1
		}

		return hhc.RequestHeadersToAdd.Compare(ohc.RequestHeadersToAdd)
	}

	return 0
}

// TCPHealthCheck configures the tcp health checker for each instance in a
// cluster.
type TCPHealthCheck struct {
	// Send is a base64 encoded string representing an array of bytes to be
	// sent in health check requests. Leaving this field empty implies a
	// connect-only health check.
	Send string `json:"send"`

	// Receive is an array of base64 encoded strings, each representing
	// array of bytes that is expected in health check responses. When
	// checking the response, "fuzzy" matching is performed such that each
	// binary block must be found, and in the order specified, but not
	// necessarily contiguously.
	Receive []string `json:"receive,omitempty"`
}

// Equals checks for equality between two TCPHealthCheck pointers
func (thc *TCPHealthCheck) Equals(othc *TCPHealthCheck) bool {
	return thc.compare(othc) == 0
}

// Treats nil as being less than defined pointers
func (thc *TCPHealthCheck) compare(othc *TCPHealthCheck) int {
	switch {
	case thc == nil && othc == nil:
		return 0
	case thc == nil && othc != nil:
		return -1
	case thc != nil && othc == nil:
		return 1
	case thc != nil && othc != nil:
		if thc.Send < othc.Send {
			return -1
		}

		if thc.Send > othc.Send {
			return 1
		}

		return arrays.CompareStringSlices(thc.Receive, othc.Receive)
	}

	return 0
}
