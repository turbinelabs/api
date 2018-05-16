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

func getOutlierDetections() (OutlierDetection, OutlierDetection) {
	od := OutlierDetection{
		IntervalMsec:                       ptr.Int(10),
		BaseEjectionTimeMsec:               ptr.Int(20),
		MaxEjectionPercent:                 ptr.Int(30),
		Consecutive5xx:                     ptr.Int(40),
		EnforcingConsecutive5xx:            ptr.Int(50),
		EnforcingSuccessRate:               ptr.Int(60),
		SuccessRateMinimumHosts:            ptr.Int(70),
		SuccessRateRequestVolume:           ptr.Int(80),
		SuccessRateStdevFactor:             ptr.Int(90),
		ConsecutiveGatewayFailure:          ptr.Int(100),
		EnforcingConsecutiveGatewayFailure: ptr.Int(10),
	}

	return od, od
}

func TestOutlierDetectionNilsAreEqual(t *testing.T) {
	a := OutlierDetection{}
	b := OutlierDetection{}

	assert.True(t, a.Equals(b))
	assert.True(t, b.Equals(a))
}

func TestOutlierDetectionIntervalMsecDifferent(t *testing.T) {
	for _, v := range []*int{nil, ptr.Int(100)} {
		a, b := getOutlierDetections()
		a.IntervalMsec = v

		assert.False(t, a.Equals(b))
		assert.False(t, b.Equals(a))
	}
}

func TestOutlierDetectionBaseEjectionTimeMsecDifferent(t *testing.T) {
	for _, v := range []*int{nil, ptr.Int(200)} {
		a, b := getOutlierDetections()
		a.BaseEjectionTimeMsec = v

		assert.False(t, a.Equals(b))
		assert.False(t, b.Equals(a))
	}
}

func TestOutlierDetectionMaxEjectionPercentDifferent(t *testing.T) {
	for _, v := range []*int{nil, ptr.Int(100)} {
		a, b := getOutlierDetections()
		a.MaxEjectionPercent = v

		assert.False(t, a.Equals(b))
		assert.False(t, b.Equals(a))
	}
}

func TestOutlierDetectionConsecutive5xxDifferent(t *testing.T) {
	for _, v := range []*int{nil, ptr.Int(400)} {
		a, b := getOutlierDetections()
		a.Consecutive5xx = v

		assert.False(t, a.Equals(b))
		assert.False(t, b.Equals(a))
	}
}

func TestOutlierDetectionEnforcingConsecutive5xxDifferent(t *testing.T) {
	for _, v := range []*int{nil, ptr.Int(100)} {
		a, b := getOutlierDetections()
		a.EnforcingConsecutive5xx = v

		assert.False(t, a.Equals(b))
		assert.False(t, b.Equals(a))
	}
}

func TestOutlierDetectionEnforcingSuccessRateDifferent(t *testing.T) {
	for _, v := range []*int{nil, ptr.Int(100)} {
		a, b := getOutlierDetections()
		a.EnforcingSuccessRate = v

		assert.False(t, a.Equals(b))
		assert.False(t, b.Equals(a))
	}
}

func TestOutlierDetectionSuccessRateMinimumHostsDifferent(t *testing.T) {
	for _, v := range []*int{nil, ptr.Int(400)} {
		a, b := getOutlierDetections()
		a.SuccessRateMinimumHosts = v

		assert.False(t, a.Equals(b))
		assert.False(t, b.Equals(a))
	}
}

func TestOutlierDetectionSuccessRateRequestVolumeDifferent(t *testing.T) {
	for _, v := range []*int{nil, ptr.Int(400)} {
		a, b := getOutlierDetections()
		a.SuccessRateRequestVolume = v

		assert.False(t, a.Equals(b))
		assert.False(t, b.Equals(a))
	}
}

func TestOutlierDetectionSuccessRateStdevFactorDifferent(t *testing.T) {
	for _, v := range []*int{nil, ptr.Int(400)} {
		a, b := getOutlierDetections()
		a.SuccessRateStdevFactor = v

		assert.False(t, a.Equals(b))
		assert.False(t, b.Equals(a))
	}
}
func TestOutlierDetectionConsecutiveGatewayFailureDifferent(t *testing.T) {
	for _, v := range []*int{nil, ptr.Int(400)} {
		a, b := getOutlierDetections()
		a.ConsecutiveGatewayFailure = v

		assert.False(t, a.Equals(b))
		assert.False(t, b.Equals(a))
	}
}

func TestOutlierDetectionEnforcingConsecutiveGatewayFailureDifferent(t *testing.T) {
	for _, v := range []*int{nil, ptr.Int(100)} {
		a, b := getOutlierDetections()
		a.EnforcingConsecutiveGatewayFailure = v

		assert.False(t, a.Equals(b))
		assert.False(t, b.Equals(a))
	}
}

func TestOutlierDetectionSamePointerValuesAreEqual(t *testing.T) {
	a, _ := getOutlierDetections()
	b := OutlierDetection{
		IntervalMsec:                       ptr.Int(10),
		BaseEjectionTimeMsec:               ptr.Int(20),
		MaxEjectionPercent:                 ptr.Int(30),
		Consecutive5xx:                     ptr.Int(40),
		EnforcingConsecutive5xx:            ptr.Int(50),
		EnforcingSuccessRate:               ptr.Int(60),
		SuccessRateMinimumHosts:            ptr.Int(70),
		SuccessRateRequestVolume:           ptr.Int(80),
		SuccessRateStdevFactor:             ptr.Int(90),
		ConsecutiveGatewayFailure:          ptr.Int(100),
		EnforcingConsecutiveGatewayFailure: ptr.Int(10),
	}

	assert.True(t, a.Equals(b))
	assert.True(t, b.Equals(a))
}

func TestOutlierDetectionIsValidOnValidObject(t *testing.T) {
	a, _ := getOutlierDetections()
	assert.Nil(t, a.IsValid())
}

func TestOutlierDetectionIsValidOnEmptyObject(t *testing.T) {
	a := OutlierDetection{}
	assert.Nil(t, a.IsValid())
}

func TestOutlierDetectionIsValidIntervalMsecSetToZero(t *testing.T) {
	a, _ := getOutlierDetections()
	a.IntervalMsec = ptr.Int(0)
	assert.DeepEqual(
		t,
		a.IsValid(),
		&ValidationError{
			[]ErrorCase{
				{"outlier_detection.interval_msec", "must be greater than zero"},
			},
		},
	)
}

func TestOutlierDetectionIsValidBaseEjectionTimeMsecNegative(t *testing.T) {
	a, _ := getOutlierDetections()
	a.BaseEjectionTimeMsec = ptr.Int(-1)
	assert.DeepEqual(
		t,
		a.IsValid(),
		&ValidationError{
			[]ErrorCase{
				{"outlier_detection.base_ejection_time_msec", "must not be negative"},
			},
		},
	)
}

func TestOutlierDetectionIsValidMaxEjectionPercentNegative(t *testing.T) {
	a, _ := getOutlierDetections()
	a.MaxEjectionPercent = ptr.Int(-2)
	assert.DeepEqual(
		t,
		a.IsValid(),
		&ValidationError{
			[]ErrorCase{
				{"outlier_detection.max_ejection_percent", "must not be negative"},
			},
		},
	)
}

func TestOutlierDetectionIsValidMaxEjectionPercentGreaterThan100(t *testing.T) {
	a, _ := getOutlierDetections()
	a.MaxEjectionPercent = ptr.Int(102)
	assert.DeepEqual(
		t,
		a.IsValid(),
		&ValidationError{
			[]ErrorCase{
				{
					"outlier_detection.max_ejection_percent",
					"must be less than or equal to 100",
				},
			},
		},
	)
}

func TestOutlierDetectionIsValidConsecutive5xxNegative(t *testing.T) {
	a, _ := getOutlierDetections()
	a.Consecutive5xx = ptr.Int(-2)
	assert.DeepEqual(
		t,
		a.IsValid(),
		&ValidationError{
			[]ErrorCase{
				{"outlier_detection.consecutive_5xx", "must not be negative"},
			},
		},
	)
}

func TestOutlierDetectionIsValidEnforcingConsecutive5xxNegative(t *testing.T) {
	a, _ := getOutlierDetections()
	a.EnforcingConsecutive5xx = ptr.Int(-2)
	assert.DeepEqual(
		t,
		a.IsValid(),
		&ValidationError{
			[]ErrorCase{
				{"outlier_detection.enforcing_consecutive_5xx", "must not be negative"},
			},
		},
	)
}

func TestOutlierDetectionIsValidEnforcingConsecutive5xxGreaterThan100(t *testing.T) {
	a, _ := getOutlierDetections()
	a.EnforcingConsecutive5xx = ptr.Int(101)
	assert.DeepEqual(
		t,
		a.IsValid(),
		&ValidationError{
			[]ErrorCase{
				{
					"outlier_detection.enforcing_consecutive_5xx",
					"must be less than or equal to 100"},
			},
		},
	)
}

func TestOutlierDetectionIsValidEnforcingSuccessRateNegative(t *testing.T) {
	a, _ := getOutlierDetections()
	a.EnforcingSuccessRate = ptr.Int(-2)
	assert.DeepEqual(
		t,
		a.IsValid(),
		&ValidationError{
			[]ErrorCase{
				{"outlier_detection.enforcing_success_rate", "must not be negative"},
			},
		},
	)
}

func TestOutlierDetectionIsValidEnforcingSuccessRateGreaterThan100(t *testing.T) {
	a, _ := getOutlierDetections()
	a.EnforcingSuccessRate = ptr.Int(102)
	assert.DeepEqual(
		t,
		a.IsValid(),
		&ValidationError{
			[]ErrorCase{
				{
					"outlier_detection.enforcing_success_rate",
					"must be less than or equal to 100",
				},
			},
		},
	)
}

func TestOutlierDetectionIsValidSuccessRateMinimumHostsNegative(t *testing.T) {
	a, _ := getOutlierDetections()
	a.SuccessRateMinimumHosts = ptr.Int(-2)
	assert.DeepEqual(
		t,
		a.IsValid(),
		&ValidationError{
			[]ErrorCase{
				{
					"outlier_detection.success_rate_minimum_hosts",
					"must not be negative",
				},
			},
		},
	)
}

func TestOutlierDetectionIsValidSuccessRateRequestVolumeSetToZero(t *testing.T) {
	a, _ := getOutlierDetections()
	a.SuccessRateRequestVolume = ptr.Int(0)
	assert.DeepEqual(
		t,
		a.IsValid(),
		&ValidationError{
			[]ErrorCase{
				{
					"outlier_detection.success_rate_request_volume",
					"must be greater than 0",
				},
			},
		},
	)
}

func TestOutlierDetectionIsValidSuccessRateStdevFactorNegative(t *testing.T) {
	a, _ := getOutlierDetections()
	a.SuccessRateStdevFactor = ptr.Int(-1)
	assert.DeepEqual(
		t,
		a.IsValid(),
		&ValidationError{
			[]ErrorCase{
				{
					"outlier_detection.success_rate_stdev_factor",
					"must not be negative",
				},
			},
		},
	)
}
func TestOutlierDetectionIsValidConsecutiveGatewayFailureNegative(t *testing.T) {
	a, _ := getOutlierDetections()
	a.ConsecutiveGatewayFailure = ptr.Int(-3)
	assert.DeepEqual(
		t,
		a.IsValid(),
		&ValidationError{
			[]ErrorCase{
				{
					"outlier_detection.consecutive_gateway_failure",
					"must not be negative",
				},
			},
		},
	)
}
func TestOutlierDetectionIsValidEnforcingConsecutiveGatewayFailureNegative(t *testing.T) {
	a, _ := getOutlierDetections()
	a.EnforcingConsecutiveGatewayFailure = ptr.Int(-1)
	assert.DeepEqual(
		t,
		a.IsValid(),
		&ValidationError{
			[]ErrorCase{
				{
					"outlier_detection.enforcing_consecutive_gateway_failure",
					"must not be negative",
				},
			},
		},
	)
}
func TestOutlierDetectionIsValidEnforcingConsecutiveGatewayFailureGreaterThan100(t *testing.T) {
	a, _ := getOutlierDetections()
	a.EnforcingConsecutiveGatewayFailure = ptr.Int(101)
	assert.DeepEqual(
		t,
		a.IsValid(),
		&ValidationError{
			[]ErrorCase{
				{
					"outlier_detection.enforcing_consecutive_gateway_failure",
					"must be less than or equal to 100",
				},
			},
		},
	)
}
