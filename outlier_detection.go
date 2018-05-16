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

import "github.com/turbinelabs/nonstdlib/ptr"

// OutlierDetection is a form of passive health checking that dynamically
// determines whether instances in a cluster are performing unlike others
// and preemptively removes them from a load balancing set.
type OutlierDetection struct {
	// The time interval between ejection analysis sweeps. This can result in
	// both new ejections due to success rate outlier detection as well as
	// hosts being returned to service. Defaults to 10s and must be greater
	// than 0.
	IntervalMsec *int `json:"interval_msec"`

	// The base time that a host is ejected for. The real time is equal to
	// the base time multiplied by the number of times the host has been
	// ejected. Defaults to 30s. Setting this to 0 means that no host will be
	// ejected for longer than `interval_msec`.
	BaseEjectionTimeMsec *int `json:"base_ejection_time_msec"`

	// The maximum % of an upstream cluster that can be ejected due to
	// outlier detection. Defaults to 10% but will always eject at least one
	// host.
	MaxEjectionPercent *int `json:"max_ejection_percent"`

	// The number of consecutive 5xx responses before a consecutive 5xx ejection
	// occurs. Defaults to 5. Setting this to 0 effectively turns off the
	// consecutive 5xx detector.
	Consecutive5xx *int `json:"consecutive_5xx"`

	// The % chance that a host will be actually ejected when an outlier status
	// is detected through consecutive 5xx. This setting can be used to disable
	// ejection or to ramp it up slowly. Defaults to 100.
	EnforcingConsecutive5xx *int `json:"enforcing_consecutive_5xx"`

	// The % chance that a host will be actually ejected when an outlier status
	// is detected through success rate statistics. This setting can be used to
	// disable ejection or to ramp it up slowly. Defaults to 100.
	EnforcingSuccessRate *int `json:"enforcing_success_rate"`

	// The number of hosts in a cluster that must have enough request volume to
	// detect success rate outliers. If the number of hosts is less than this
	// setting, outlier detection via success rate statistics is not performed
	// for any host in the cluster. Defaults to 5. Setting this to 0 effectively
	// triggers the success rate detector regardless of the number of valid hosts
	// during an interval (as determined by `success_rate_request_volume`).
	SuccessRateMinimumHosts *int `json:"success_rate_minimum_hosts"`

	// The minimum number of total requests that must be collected in one
	// interval (as defined by the interval duration) to include this host
	// in success rate based outlier detection. If the volume is lower than this
	// setting, outlier detection via success rate statistics is not performed
	// for that host. Defaults to 100.
	SuccessRateRequestVolume *int `json:"success_rate_request_volume"`

	// This factor is used to determine the ejection threshold for success rate
	// outlier ejection. The ejection threshold is the difference between the
	// mean success rate, and the product of this factor and the standard
	// deviation of the mean success rate: mean - (stdev *
	// success_rate_stdev_factor). This factor is divided by a thousand to get a
	// double. That is, if the desired factor is 1.9, the runtime value should
	// be 1900. Defaults to 1900. Setting this to 0 effectively turns off the
	// success rate detector.
	SuccessRateStdevFactor *int `json:"success_rate_stdev_factor"`

	// The number of consecutive gateway failures (502, 503, 504 status or
	// connection errors that are mapped to one of those status codes) before a
	// consecutive gateway failure ejection occurs. Defaults to 5.
	ConsecutiveGatewayFailure *int `json:"consecutive_gateway_failure"`

	// The % chance that a host will be actually ejected when an outlier status
	// is detected through consecutive gateway failures. This setting can be
	// used to disable ejection or to ramp it up slowly. Defaults to 0.
	EnforcingConsecutiveGatewayFailure *int `json:"enforcing_consecutive_gateway_failure"`
}

// Equals compares two OutlierDetections for equality
func (od OutlierDetection) Equals(o OutlierDetection) bool {
	return ptr.IntEqual(od.IntervalMsec, o.IntervalMsec) &&
		ptr.IntEqual(od.BaseEjectionTimeMsec, o.BaseEjectionTimeMsec) &&
		ptr.IntEqual(od.MaxEjectionPercent, o.MaxEjectionPercent) &&
		ptr.IntEqual(od.Consecutive5xx, o.Consecutive5xx) &&
		ptr.IntEqual(od.EnforcingConsecutive5xx, o.EnforcingConsecutive5xx) &&
		ptr.IntEqual(od.EnforcingSuccessRate, o.EnforcingSuccessRate) &&
		ptr.IntEqual(od.SuccessRateMinimumHosts, o.SuccessRateMinimumHosts) &&
		ptr.IntEqual(od.SuccessRateRequestVolume, o.SuccessRateRequestVolume) &&
		ptr.IntEqual(od.SuccessRateStdevFactor, o.SuccessRateStdevFactor) &&
		ptr.IntEqual(od.ConsecutiveGatewayFailure, o.ConsecutiveGatewayFailure) &&
		ptr.IntEqual(
			od.EnforcingConsecutiveGatewayFailure,
			o.EnforcingConsecutiveGatewayFailure,
		)
}

// IsValid checks for the validity of contained fields.
func (od OutlierDetection) IsValid() *ValidationError {
	scope := func(s string) string { return "outlier_detection." + s }

	errs := &ValidationError{}
	// While CDS allows this to be set to 0, doing so will cause envoy to get stuck in a
	// loop executing its ejection analysis.
	if od.IntervalMsec != nil && *od.IntervalMsec < 1 {
		errs.AddNew(ErrorCase{scope("interval_msec"), "must be greater than zero"})
	}

	if od.BaseEjectionTimeMsec != nil && *od.BaseEjectionTimeMsec < 0 {
		errs.AddNew(ErrorCase{scope("base_ejection_time_msec"), "must not be negative"})
	}

	if od.MaxEjectionPercent != nil && *od.MaxEjectionPercent < 0 {
		errs.AddNew(
			ErrorCase{
				scope("max_ejection_percent"),
				"must not be negative",
			},
		)
	}

	if od.MaxEjectionPercent != nil && *od.MaxEjectionPercent > 100 {
		errs.AddNew(
			ErrorCase{
				scope("max_ejection_percent"),
				"must be less than or equal to 100",
			},
		)
	}

	if od.Consecutive5xx != nil && *od.Consecutive5xx < 0 {
		errs.AddNew(ErrorCase{scope("consecutive_5xx"), "must not be negative"})
	}

	if od.EnforcingConsecutive5xx != nil && *od.EnforcingConsecutive5xx < 0 {
		errs.AddNew(
			ErrorCase{
				scope("enforcing_consecutive_5xx"),
				"must not be negative",
			},
		)
	}

	if od.EnforcingConsecutive5xx != nil && *od.EnforcingConsecutive5xx > 100 {
		errs.AddNew(
			ErrorCase{
				scope("enforcing_consecutive_5xx"),
				"must be less than or equal to 100",
			},
		)
	}

	if od.EnforcingSuccessRate != nil && *od.EnforcingSuccessRate < 0 {
		errs.AddNew(
			ErrorCase{
				scope("enforcing_success_rate"),
				"must not be negative",
			},
		)
	}

	if od.EnforcingSuccessRate != nil && *od.EnforcingSuccessRate > 100 {
		errs.AddNew(
			ErrorCase{
				scope("enforcing_success_rate"),
				"must be less than or equal to 100",
			},
		)
	}

	if od.SuccessRateMinimumHosts != nil && *od.SuccessRateMinimumHosts < 0 {
		errs.AddNew(
			ErrorCase{
				scope("success_rate_minimum_hosts"),
				"must not be negative",
			},
		)
	}

	if od.SuccessRateRequestVolume != nil && *od.SuccessRateRequestVolume < 1 {
		errs.AddNew(
			ErrorCase{
				scope("success_rate_request_volume"),
				"must be greater than 0",
			},
		)
	}

	if od.SuccessRateStdevFactor != nil && *od.SuccessRateStdevFactor < 0 {
		errs.AddNew(
			ErrorCase{
				scope("success_rate_stdev_factor"),
				"must not be negative",
			},
		)
	}

	if od.ConsecutiveGatewayFailure != nil && *od.ConsecutiveGatewayFailure < 0 {
		errs.AddNew(
			ErrorCase{
				scope("consecutive_gateway_failure"),
				"must not be negative",
			},
		)
	}

	if od.EnforcingConsecutiveGatewayFailure != nil && *od.EnforcingConsecutiveGatewayFailure < 0 {
		errs.AddNew(
			ErrorCase{
				scope("enforcing_consecutive_gateway_failure"),
				"must not be negative",
			},
		)
	}

	if od.EnforcingConsecutiveGatewayFailure != nil && *od.EnforcingConsecutiveGatewayFailure > 100 {
		errs.AddNew(
			ErrorCase{
				scope("enforcing_consecutive_gateway_failure"),
				"must be less than or equal to 100",
			},
		)
	}

	return errs.OrNil()

}

// OutlierDetectionPtrEquals provides a way to compare two OutlierDetection
// pointers
func OutlierDetectionPtrEquals(od1, od2 *OutlierDetection) bool {
	switch {
	case od1 == nil && od2 == nil:
		return true
	case od1 == nil || od2 == nil:
		return false
	default:
		return od1.Equals(*od2)
	}
}
