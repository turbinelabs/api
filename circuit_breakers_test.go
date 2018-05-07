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

func getCircuitBreakers() (CircuitBreakers, CircuitBreakers) {
	cb := CircuitBreakers{
		MaxConnections:     ptr.Int(10),
		MaxPendingRequests: ptr.Int(20),
		MaxRetries:         ptr.Int(30),
		MaxRequests:        ptr.Int(40),
	}

	return cb, cb
}

func TestCircuitBreakersNilsAreEqual(t *testing.T) {
	a := CircuitBreakers{}
	b := CircuitBreakers{}

	assert.True(t, a.Equals(b))
	assert.True(t, b.Equals(a))
}

func TestCircuitBreakersMaxConnsDifferent(t *testing.T) {
	for _, v := range []*int{nil, ptr.Int(100)} {
		a, b := getCircuitBreakers()
		a.MaxConnections = v

		assert.False(t, a.Equals(b))
		assert.False(t, b.Equals(a))
	}
}

func TestCircuitBreakersMaxPendingRequestsDifferent(t *testing.T) {
	for _, v := range []*int{nil, ptr.Int(200)} {
		a, b := getCircuitBreakers()
		a.MaxPendingRequests = v

		assert.False(t, a.Equals(b))
		assert.False(t, b.Equals(a))
	}
}

func TestCircuitBreakersMaxRetriesDifferent(t *testing.T) {
	for _, v := range []*int{nil, ptr.Int(300)} {
		a, b := getCircuitBreakers()
		a.MaxRetries = v

		assert.False(t, a.Equals(b))
		assert.False(t, b.Equals(a))
	}
}

func TestCircuitBreakersMaxRequestsDifferent(t *testing.T) {
	for _, v := range []*int{nil, ptr.Int(400)} {
		a, b := getCircuitBreakers()
		a.MaxRequests = v

		assert.False(t, a.Equals(b))
		assert.False(t, b.Equals(a))
	}
}

func TestCircuitBreakersSamePointerValuesAreEqual(t *testing.T) {
	a, _ := getCircuitBreakers()
	b := CircuitBreakers{
		MaxConnections:     ptr.Int(10),
		MaxPendingRequests: ptr.Int(20),
		MaxRetries:         ptr.Int(30),
		MaxRequests:        ptr.Int(40),
	}

	assert.True(t, a.Equals(b))
	assert.True(t, b.Equals(a))
}

func TestCircuitBreakersIsValidOnValidObject(t *testing.T) {
	a, _ := getCircuitBreakers()
	assert.Nil(t, a.IsValid())
}

func TestCircuitBreakersIsValidOnEmptyObject(t *testing.T) {
	a := CircuitBreakers{}
	assert.Nil(t, a.IsValid())
}

func TestCircuitBreakersIsValidNegativeMaxConnections(t *testing.T) {
	a, _ := getCircuitBreakers()
	a.MaxConnections = ptr.Int(-199)
	assert.DeepEqual(
		t,
		a.IsValid(),
		&ValidationError{
			[]ErrorCase{
				{"circuit_breakers.max_connections", "must not be negative"},
			},
		},
	)
}

func TestCircuitBreakersIsValidNegativeMaxPendingRequests(t *testing.T) {
	a, _ := getCircuitBreakers()
	a.MaxPendingRequests = ptr.Int(-1)
	assert.DeepEqual(
		t,
		a.IsValid(),
		&ValidationError{
			[]ErrorCase{
				{"circuit_breakers.max_pending_requests", "must not be negative"},
			},
		},
	)
}

func TestCircuitBreakersIsValidNegativeMaxRetries(t *testing.T) {
	a, _ := getCircuitBreakers()
	a.MaxRetries = ptr.Int(-2)
	assert.DeepEqual(
		t,
		a.IsValid(),
		&ValidationError{
			[]ErrorCase{
				{"circuit_breakers.max_retries", "must not be negative"},
			},
		},
	)
}

func TestCircuitBreakersIsValidNegativeMaxRequests(t *testing.T) {
	a, _ := getCircuitBreakers()
	a.MaxRequests = ptr.Int(-3)
	assert.DeepEqual(
		t,
		a.IsValid(),
		&ValidationError{
			[]ErrorCase{
				{"circuit_breakers.max_requests", "must not be negative"},
			},
		},
	)
}
