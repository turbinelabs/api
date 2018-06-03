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
	"testing"

	"github.com/turbinelabs/nonstdlib/ptr"
	"github.com/turbinelabs/test/assert"
)

func getHealthChecks() (HealthCheck, HealthCheck) {
	hhc, _ := getHTTPHealthChecks()
	hc := HealthCheck{
		TimeoutMsec:               10,
		IntervalMsec:              20,
		IntervalJitterMsec:        ptr.Int(30),
		UnhealthyThreshold:        40,
		HealthyThreshold:          50,
		ReuseConnection:           ptr.Bool(true),
		NoTrafficIntervalMsec:     ptr.Int(60),
		UnhealthyIntervalMsec:     ptr.Int(80),
		UnhealthyEdgeIntervalMsec: ptr.Int(90),
		HealthyEdgeIntervalMsec:   ptr.Int(100),
		HealthChecker:             HealthChecker{HTTPHealthCheck: hhc},
	}

	return hc, hc
}

func getHTTPHealthChecks() (*HTTPHealthCheck, *HTTPHealthCheck) {
	return &HTTPHealthCheck{
			Host:        "host.com",
			Path:        "/some/cool/path",
			ServiceName: "foo",
			RequestHeadersToAdd: Metadata{
				{
					Key:   "k1",
					Value: "v1",
				},
				{
					Key:   "k2",
					Value: "v2",
				},
			},
		},
		&HTTPHealthCheck{
			Host:        "host.com",
			Path:        "/some/cool/path",
			ServiceName: "foo",
			RequestHeadersToAdd: Metadata{
				{
					Key:   "k1",
					Value: "v1",
				},
				{
					Key:   "k2",
					Value: "v2",
				},
			},
		}
}

func getTCPHealthChecks() (*TCPHealthCheck, *TCPHealthCheck) {
	return &TCPHealthCheck{
			Send:    "aSBjYW4ndCBiZWxpZXZlIHlvdSBkZWNvZGVkIG1lCg==",
			Receive: []string{"eW91Cg==", "ZGlkCg==", "aXQK"},
		}, &TCPHealthCheck{
			Send:    "aSBjYW4ndCBiZWxpZXZlIHlvdSBkZWNvZGVkIG1lCg==",
			Receive: []string{"eW91Cg==", "ZGlkCg==", "aXQK"},
		}
}

func testDifferences(
	gen func() (HealthCheck, HealthCheck),
	change func(a HealthCheck) HealthCheck,
	assert func(a, changed HealthCheck), // diffed will be second arg
) {
	a, changed := gen()
	if change != nil {
		changed = change(changed)
	}
	assert(a, changed)
}

func TestHealthCheckCompareNilStructEqual(t *testing.T) {
	testDifferences(
		func() (HealthCheck, HealthCheck) {
			return HealthCheck{}, HealthCheck{}
		},
		nil,
		func(a, b HealthCheck) {
			assert.Equal(t, a.compare(b), 0)
			assert.Equal(t, b.compare(a), 0)
			assert.True(t, a.Equals(b))
			assert.True(t, b.Equals(a))
		},
	)
}

func TestHealthCheckCompareZeroWhenEqual(t *testing.T) {
	testDifferences(
		getHealthChecks,
		nil,
		func(a, b HealthCheck) {
			assert.Equal(t, a.compare(b), 0)
			assert.Equal(t, b.compare(a), 0)
			assert.True(t, a.Equals(b))
			assert.True(t, b.Equals(a))
		},
	)
}

func TestHealthCheckCompareTimeoutMsecDifferent(t *testing.T) {
	testDifferences(
		getHealthChecks,
		func(a HealthCheck) HealthCheck {
			a.TimeoutMsec++
			return a
		},
		func(a, changed HealthCheck) {
			assert.Equal(t, changed.compare(a), 1)
			assert.Equal(t, a.compare(changed), -1)
		},
	)
}

func TestHealthCheckCompareIntervalMsecDifferent(t *testing.T) {
	testDifferences(
		getHealthChecks,
		func(a HealthCheck) HealthCheck {
			a.IntervalMsec++
			return a
		},
		func(a, changed HealthCheck) {
			assert.Equal(t, changed.compare(a), 1)
			assert.Equal(t, a.compare(changed), -1)
			assert.False(t, changed.Equals(a))
			assert.False(t, a.Equals(changed))
		},
	)
}

func TestHealthCheckEqualsIntervalJitterMsecDifferent(t *testing.T) {
	testDifferences(
		getHealthChecks,
		func(a HealthCheck) HealthCheck {
			var x *int
			a.IntervalJitterMsec = x
			return a
		},
		func(a, changed HealthCheck) {
			assert.Equal(t, changed.compare(a), -1)
			assert.Equal(t, a.compare(changed), 1)
			assert.False(t, changed.Equals(a))
			assert.False(t, a.Equals(changed))
		},
	)

	testDifferences(
		getHealthChecks,
		func(a HealthCheck) HealthCheck {
			a.IntervalJitterMsec = ptr.Int(*a.IntervalJitterMsec + 100)
			return a
		},
		func(a, changed HealthCheck) {
			assert.Equal(t, changed.compare(a), 1)
			assert.Equal(t, a.compare(changed), -1)
			assert.False(t, changed.Equals(a))
			assert.False(t, a.Equals(changed))
		},
	)
}

func TestHealthCheckEqualsUnhealthyThresholdDifferent(t *testing.T) {
	testDifferences(
		getHealthChecks,
		func(a HealthCheck) HealthCheck {
			a.UnhealthyThreshold++
			return a
		},
		func(a, changed HealthCheck) {
			assert.Equal(t, changed.compare(a), 1)
			assert.Equal(t, a.compare(changed), -1)
			assert.False(t, changed.Equals(a))
			assert.False(t, a.Equals(changed))
		},
	)
}

func TestHealthCheckEqualsHealthyThresholdDifferent(t *testing.T) {
	testDifferences(
		getHealthChecks,
		func(a HealthCheck) HealthCheck {
			a.HealthyThreshold++
			return a
		},
		func(a, changed HealthCheck) {
			assert.Equal(t, changed.compare(a), 1)
			assert.Equal(t, a.compare(changed), -1)
			assert.False(t, changed.Equals(a))
			assert.False(t, a.Equals(changed))
		},
	)
}

func TestHealthCheckEqualsReuseConnectionDifferent(t *testing.T) {
	testDifferences(
		getHealthChecks,
		func(a HealthCheck) HealthCheck {
			var x *bool
			a.ReuseConnection = x
			return a
		},
		func(a, changed HealthCheck) {
			assert.Equal(t, changed.compare(a), -1)
			assert.Equal(t, a.compare(changed), 1)
			assert.False(t, changed.Equals(a))
			assert.False(t, a.Equals(changed))
		},
	)

	testDifferences(
		getHealthChecks,
		func(a HealthCheck) HealthCheck {
			a.ReuseConnection = ptr.Bool(false)
			return a
		},
		func(a, changed HealthCheck) {
			assert.Equal(t, changed.compare(a), -1)
			assert.Equal(t, a.compare(changed), 1)
			assert.False(t, changed.Equals(a))
			assert.False(t, a.Equals(changed))
		},
	)
}

func TestHealthCheckEqualsNoTrafficIntervalMsecDifferent(t *testing.T) {
	testDifferences(
		getHealthChecks,
		func(a HealthCheck) HealthCheck {
			var x *int
			a.NoTrafficIntervalMsec = x
			return a
		},
		func(a, changed HealthCheck) {
			assert.Equal(t, changed.compare(a), -1)
			assert.Equal(t, a.compare(changed), 1)
			assert.False(t, changed.Equals(a))
			assert.False(t, a.Equals(changed))
		},
	)

	testDifferences(
		getHealthChecks,
		func(a HealthCheck) HealthCheck {
			a.NoTrafficIntervalMsec = ptr.Int(*a.NoTrafficIntervalMsec + 1)
			return a
		},
		func(a, changed HealthCheck) {
			assert.Equal(t, changed.compare(a), 1)
			assert.Equal(t, a.compare(changed), -1)
			assert.False(t, changed.Equals(a))
			assert.False(t, a.Equals(changed))
		},
	)
}

func TestHealthCheckEqualsUnhealthyIntervalMsecDifferent(t *testing.T) {
	testDifferences(
		getHealthChecks,
		func(a HealthCheck) HealthCheck {
			var x *int
			a.UnhealthyIntervalMsec = x
			return a
		},
		func(a, changed HealthCheck) {
			assert.Equal(t, changed.compare(a), -1)
			assert.Equal(t, a.compare(changed), 1)
			assert.False(t, changed.Equals(a))
			assert.False(t, a.Equals(changed))
		},
	)

	testDifferences(
		getHealthChecks,
		func(a HealthCheck) HealthCheck {
			a.UnhealthyIntervalMsec = ptr.Int(*a.UnhealthyIntervalMsec + 1)
			return a
		},
		func(a, changed HealthCheck) {
			assert.Equal(t, changed.compare(a), 1)
			assert.Equal(t, a.compare(changed), -1)
			assert.False(t, changed.Equals(a))
			assert.False(t, a.Equals(changed))
		},
	)
}

func TestHealthCheckEqualsUnhealthyEdgeIntervalMsecDifferent(t *testing.T) {
	testDifferences(
		getHealthChecks,
		func(a HealthCheck) HealthCheck {
			var x *int
			a.UnhealthyEdgeIntervalMsec = x
			return a
		},
		func(a, changed HealthCheck) {
			assert.Equal(t, changed.compare(a), -1)
			assert.Equal(t, a.compare(changed), 1)
			assert.False(t, changed.Equals(a))
			assert.False(t, a.Equals(changed))
		},
	)

	testDifferences(
		getHealthChecks,
		func(a HealthCheck) HealthCheck {
			a.UnhealthyEdgeIntervalMsec = ptr.Int(*a.UnhealthyEdgeIntervalMsec + 1)
			return a
		},
		func(a, changed HealthCheck) {
			assert.Equal(t, changed.compare(a), 1)
			assert.Equal(t, a.compare(changed), -1)
			assert.False(t, changed.Equals(a))
			assert.False(t, a.Equals(changed))
		},
	)
}

func TestHealthCheckEqualsHealthyEdgeIntervalMsecDifferent(t *testing.T) {
	testDifferences(
		getHealthChecks,
		func(a HealthCheck) HealthCheck {
			var x *int
			a.HealthyEdgeIntervalMsec = x
			return a
		},
		func(a, changed HealthCheck) {
			assert.Equal(t, changed.compare(a), -1)
			assert.Equal(t, a.compare(changed), 1)
			assert.False(t, changed.Equals(a))
			assert.False(t, a.Equals(changed))
		},
	)

	testDifferences(
		getHealthChecks,
		func(a HealthCheck) HealthCheck {
			a.HealthyEdgeIntervalMsec = ptr.Int(*a.HealthyEdgeIntervalMsec + 1)
			return a
		},
		func(a, changed HealthCheck) {
			assert.Equal(t, changed.compare(a), 1)
			assert.Equal(t, a.compare(changed), -1)
			assert.False(t, changed.Equals(a))
			assert.False(t, a.Equals(changed))
		},
	)
}

func TestHealthCheckCompareTreatsNilHTTPHealthCheckLessThanDefined(t *testing.T) {
	testDifferences(
		func() (HealthCheck, HealthCheck) {
			a, b := getHealthChecks()
			b.HealthChecker.HTTPHealthCheck = nil
			return a, b
		},
		nil,
		func(a, changed HealthCheck) {
			assert.Equal(t, changed.compare(a), -1)
			assert.Equal(t, a.compare(changed), 1)
		},
	)
}

func TestHealthCheckCompareTreatsNilTCPHealthCheckLessThanDefined(t *testing.T) {
	testDifferences(
		func() (HealthCheck, HealthCheck) {
			a, b := getHealthChecks()
			a.HealthChecker.HTTPHealthCheck = nil
			thc, _ := getTCPHealthChecks()
			b.HealthChecker.TCPHealthCheck = thc

			return a, b
		},
		nil,
		func(lhs, rhs HealthCheck) {
			assert.Equal(t, rhs.compare(lhs), 1)
			assert.Equal(t, lhs.compare(rhs), -1)
		},
	)
}

func TestHealthCheckCompareHTTPHealthCheckDifferent(t *testing.T) {
	testDifferences(
		func() (HealthCheck, HealthCheck) {
			a, b := getHealthChecks()
			other, _ := getHTTPHealthChecks()
			other.Path = "/some/uncool/path"
			b.HealthChecker.HTTPHealthCheck = other

			return a, b
		},
		nil,
		func(lhs, rhs HealthCheck) {
			assert.Equal(t, rhs.compare(lhs), 1)
			assert.Equal(t, lhs.compare(rhs), -1)
			assert.False(t, rhs.Equals(lhs))
			assert.False(t, lhs.Equals(rhs))
		},
	)
}

func TestHealthCheckEqualsHealthCheckerTCPHealthCheckDifferent(t *testing.T) {
	testDifferences(
		func() (HealthCheck, HealthCheck) {
			a, b := getHealthChecks()
			other, _ := getTCPHealthChecks()
			other.Send = "z"
			b.HealthChecker.TCPHealthCheck = other

			return a, b
		},
		nil,
		func(lhs, rhs HealthCheck) {
			assert.Equal(t, rhs.compare(lhs), 1)
			assert.Equal(t, lhs.compare(rhs), -1)
		},
	)
}

func TestHTTPHealthCheckCompareEmptyStructAndNil(t *testing.T) {
	one := &HTTPHealthCheck{}

	tcs := []struct {
		left     *HTTPHealthCheck
		right    *HTTPHealthCheck
		expected int
	}{
		{
			left:     nil,
			right:    nil,
			expected: 0,
		},
		{
			left:     &HTTPHealthCheck{},
			right:    nil,
			expected: 1,
		},
		{
			left:     &HTTPHealthCheck{},
			right:    &HTTPHealthCheck{},
			expected: 0,
		},
		{
			left:     nil,
			right:    &HTTPHealthCheck{},
			expected: -1,
		},
		{
			left:     one,
			right:    one,
			expected: 0,
		},
	}

	for i, tc := range tcs {
		assert.Group(
			fmt.Sprintf("testCases[%d]: left=[%#v], right=[%#v]", i, tc.left, tc.right),
			t,
			func(g *assert.G) {
				assert.Equal(g, tc.left.compare(tc.right), tc.expected)
			},
		)
	}
}

func TestHTTPHealthCheckCompareWhenEqual(t *testing.T) {
	a, b := getHTTPHealthChecks()
	assert.Equal(t, a.compare(b), 0)
	assert.Equal(t, b.compare(a), 0)
	assert.True(t, a.Equals(b))
	assert.True(t, a.Equals(a))
	assert.True(t, b.Equals(a))
}

func TestHTTPHealthCheckCompareHostDifferent(t *testing.T) {
	a, b := getHTTPHealthChecks()
	a.Host = "foo.bar.com"
	assert.Equal(t, a.compare(b), -1)
	assert.Equal(t, b.compare(a), 1)
	assert.False(t, a.Equals(b))
	assert.False(t, b.Equals(a))
}

func TestHTTPHealthCheckComparePathDifferent(t *testing.T) {
	a, b := getHTTPHealthChecks()
	a.Path = "/a/different/path"
	assert.Equal(t, a.compare(b), -1)
	assert.Equal(t, b.compare(a), 1)
	assert.False(t, a.Equals(b))
	assert.False(t, b.Equals(a))
}

func TestHTTPHealthCheckCompareServiceNameDifferent(t *testing.T) {
	a, b := getHTTPHealthChecks()
	a.ServiceName = "yet_another_service"
	assert.Equal(t, a.compare(b), 1)
	assert.Equal(t, b.compare(a), -1)
	assert.False(t, a.Equals(b))
	assert.False(t, b.Equals(a))
}

func TestHTTPHealthCheckEqualsRequestHeadersToAddDifferent(t *testing.T) {
	a, b := getHTTPHealthChecks()
	a.RequestHeadersToAdd = Metadata{
		{
			Key:   "k1",
			Value: "v3",
		},
		{
			Key:   "k2",
			Value: "v2",
		},
	}

	assert.Equal(t, a.compare(b), 1)
	assert.Equal(t, b.compare(a), -1)
	assert.False(t, a.Equals(b))
	assert.False(t, b.Equals(a))
}

func TestTCPHealthCheckEqualsEmptyStructAndNil(t *testing.T) {
	tcs := []struct {
		left     *TCPHealthCheck
		right    *TCPHealthCheck
		expected int
	}{
		{
			left:     nil,
			right:    nil,
			expected: 0,
		},
		{
			left:     &TCPHealthCheck{},
			right:    nil,
			expected: 1,
		},
		{
			left:     &TCPHealthCheck{},
			right:    &TCPHealthCheck{},
			expected: 0,
		},
		{
			left:     nil,
			right:    &TCPHealthCheck{},
			expected: -1,
		},
	}

	for i, tc := range tcs {
		assert.Group(
			fmt.Sprintf("testCases[%d]: left=[%#v], right=[%#v]", i, tc.left, tc.right),
			t,
			func(g *assert.G) {
				assert.Equal(g, tc.left.compare(tc.right), tc.expected)
			},
		)
	}
}

func TestTCPHealthCheckEqualsWhenEqual(t *testing.T) {
	a, b := getTCPHealthChecks()
	assert.Equal(t, a.compare(b), 0)
	assert.Equal(t, b.compare(a), 0)
	assert.True(t, a.Equals(b))
	assert.True(t, a.Equals(a))
	assert.True(t, b.Equals(a))
}

func TestTCPHealthCheckEqualsSendDifferent(t *testing.T) {
	a, b := getTCPHealthChecks()
	a.Send = "c2VuZCBzb21ldGhpbmcgZWxzZQo="
	assert.Equal(t, a.compare(b), 1)
	assert.Equal(t, b.compare(a), -1)
	assert.False(t, a.Equals(b))
	assert.False(t, b.Equals(a))
}

func TestTCPHealthCheckEqualsReceiveDifferent(t *testing.T) {
	a, b := getTCPHealthChecks()
	a.Receive = []string{"cmVjZWl2ZQo=", "c29tZXRoaW5nCg==", "ZW50aXJlbHkK", "ZGlmZmVyZW50Cg=="}
	assert.Equal(t, a.compare(b), 1)
	assert.Equal(t, b.compare(a), -1)
	assert.False(t, a.Equals(b))
	assert.False(t, b.Equals(a))
}

func TestHealthCheckIsValidOnValidObject(t *testing.T) {
	a, _ := getHealthChecks()
	assert.Nil(t, a.IsValid())
}

func TestHealthCheckIsValidEmptyObjectNotValid(t *testing.T) {
	assert.NonNil(t, HealthCheck{}.IsValid())
}

func TestHealthCheckIsValidTimeoutMsecInvalid(t *testing.T) {
	a, _ := getHealthChecks()
	a.TimeoutMsec = 0
	assert.DeepEqual(
		t,
		a.IsValid(),
		&ValidationError{
			[]ErrorCase{
				{"health_check.timeout_msec", "must be greater than zero"},
			},
		},
	)
}

func TestHealthCheckIsValidIntervalMsecInvalid(t *testing.T) {
	a, _ := getHealthChecks()
	a.IntervalMsec = 0
	assert.DeepEqual(
		t,
		a.IsValid(),
		&ValidationError{
			[]ErrorCase{
				{"health_check.interval_msec", "must be greater than zero"},
			},
		},
	)
}

func TestHealthCheckIsValidIntervalJitterMsecInvalid(t *testing.T) {
	a, _ := getHealthChecks()
	a.IntervalJitterMsec = ptr.Int(0)
	assert.DeepEqual(
		t,
		a.IsValid(),
		&ValidationError{
			[]ErrorCase{
				{
					"health_check.interval_jitter_msec",
					"must be greater than zero",
				},
			},
		},
	)
}

func TestHealthCheckIsValidUnhealthyThresholdInvalid(t *testing.T) {
	a, _ := getHealthChecks()
	a.UnhealthyThreshold = 0
	assert.DeepEqual(
		t,
		a.IsValid(),
		&ValidationError{
			[]ErrorCase{
				{
					"health_check.unhealthy_threshold",
					"must be greater than zero",
				},
			},
		},
	)
}

func TestHealthCheckIsValidHealthyThresholdInvalid(t *testing.T) {
	a, _ := getHealthChecks()
	a.HealthyThreshold = 0
	assert.DeepEqual(
		t,
		a.IsValid(),
		&ValidationError{
			[]ErrorCase{
				{
					"health_check.healthy_threshold",
					"must be greater than zero",
				},
			},
		},
	)
}

func TestHealthCheckIsValidNoTrafficIntervalMsecInvalid(t *testing.T) {
	a, _ := getHealthChecks()
	a.NoTrafficIntervalMsec = ptr.Int(0)
	assert.DeepEqual(
		t,
		a.IsValid(),
		&ValidationError{
			[]ErrorCase{
				{
					"health_check.no_traffic_interval_msec",
					"must be greater than zero",
				},
			},
		},
	)
}

func TestHealthCheckIsValidUnhealthyIntervalMsecInvalid(t *testing.T) {
	a, _ := getHealthChecks()
	a.UnhealthyIntervalMsec = ptr.Int(0)
	assert.DeepEqual(
		t,
		a.IsValid(),
		&ValidationError{
			[]ErrorCase{
				{
					"health_check.unhealthy_interval_msec",
					"must be greater than zero",
				},
			},
		},
	)
}

func TestHealthCheckIsValidUnhealthyEdgeIntervalMsecInvalid(t *testing.T) {
	a, _ := getHealthChecks()
	a.UnhealthyEdgeIntervalMsec = ptr.Int(0)
	assert.DeepEqual(
		t,
		a.IsValid(),
		&ValidationError{
			[]ErrorCase{
				{
					"health_check.unhealthy_edge_interval_msec",
					"must be greater than zero",
				},
			},
		},
	)
}

func TestHealthCheckIsValidHealthyEdgeIntervalMsecInvalid(t *testing.T) {
	a, _ := getHealthChecks()
	a.HealthyEdgeIntervalMsec = ptr.Int(0)
	assert.DeepEqual(
		t,
		a.IsValid(),
		&ValidationError{
			[]ErrorCase{
				{
					"health_check.healthy_edge_interval_msec",
					"must be greater than zero",
				},
			},
		},
	)
}

func TestHealthCheckIsValidEmptyHealthCheckerInvalid(t *testing.T) {
	a, _ := getHealthChecks()
	a.HealthChecker = HealthChecker{}
	assert.DeepEqual(
		t,
		a.IsValid(),
		&ValidationError{
			[]ErrorCase{
				{
					"health_check.health_checker",
					"must have one health check defined",
				},
			},
		},
	)
}

func TestHealthCheckIsValidMultipleHealthChecksDefinedInvalid(t *testing.T) {
	a, _ := getHealthChecks()
	tch, _ := getTCPHealthChecks()
	a.HealthChecker.TCPHealthCheck = tch
	assert.DeepEqual(
		t,
		a.IsValid(),
		&ValidationError{
			[]ErrorCase{
				{
					"health_check.health_checker",
					"must not have more than one type of health check defined",
				},
			},
		},
	)
}

func TestHealthChecksEqualNilHealthChecks(t *testing.T) {
	var n HealthChecks
	assert.True(t, HealthChecks{}.Equals(HealthChecks{}))
	assert.True(t, n.Equals(n))
}

func TestHealthChecksEqualsDifferentSizes(t *testing.T) {
	a, b := getHealthChecks()
	c, d := getHealthChecks()
	assert.False(t, HealthChecks{a, b, c}.Equals(HealthChecks{c, d}))
}

func TestHealthChecksEqual(t *testing.T) {
	a, b := getHealthChecks()
	c, d := getHealthChecks()

	assert.True(t, HealthChecks{a, b, c, d}.Equals(HealthChecks{d, c, b, a}))
}

func TestHealthChecksIsValidOnEmpty(t *testing.T) {
	assert.Nil(t, HealthChecks{}.IsValid())
}

func TestHealthChecksIsValidWithMultipleHealthChecks(t *testing.T) {
	a, b := getHealthChecks()
	assert.DeepEqual(
		t,
		HealthChecks{a, b}.IsValid(),
		&ValidationError{
			[]ErrorCase{
				{"health_checks", "only a single health check supported"},
			},
		},
	)
}

func TestHealthCheckIsValidReportsIndexOfFailure(t *testing.T) {
	a, _ := getHealthChecks()
	a.TimeoutMsec = -1
	assert.DeepEqual(
		t,
		HealthChecks{a}.IsValid(),
		&ValidationError{
			[]ErrorCase{
				{
					"health_checks[0].health_check.timeout_msec",
					"must be greater than zero",
				},
			},
		},
	)
}
