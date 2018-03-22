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

// CircuitBreakers provides limits on various parameters to protect clusters
// against sudden surges in traffic.
type CircuitBreakers struct {
	// MaxConnections is the maximum number of connections that will be
	// established to all instances in a cluster within a proxy.
	// If set to 0, no new connections will be created. If not specified,
	// defaults to 1024.
	MaxConnections *int `json:"max_connections"`

	// MaxPendingRequests is the maximum number of requests that will be
	// queued while waiting on a connection pool to a cluster within a proxy.
	// If set to 0, no requests will be queued. If not specified,
	// defaults to 1024.
	MaxPendingRequests *int `json:"max_pending_requests"`

	// MaxRetries is the maximum number of retries that can be outstanding
	// to all instances in a cluster within a proxy. If set to 0, requests
	// will not be retried. If not specified, defaults to 3.
	MaxRetries *int `json:"max_retries"`

	// MaxRequests is the maximum number of requests that can be outstanding
	// to all instances in a cluster within a proxy. Only applicable to
	// HTTP/2 traffic since HTTP/1.1 clusters are governed by the maximum
	// connections circuit breaker. If set to 0, no requests will be made.
	// If not specified, defaults to 1024.
	MaxRequests *int `json:"max_requests"`
}

// Equals compares two CircuitBreakers for equality
func (cb CircuitBreakers) Equals(o CircuitBreakers) bool {
	return compareIntPtrs(cb.MaxConnections, o.MaxConnections) &&
		compareIntPtrs(cb.MaxPendingRequests, o.MaxPendingRequests) &&
		compareIntPtrs(cb.MaxRetries, o.MaxRetries) &&
		compareIntPtrs(cb.MaxRequests, o.MaxRequests)
}

func compareIntPtrs(a, b *int) bool {
	if a == nil && b == nil {
		return true
	}

	if a == nil || b == nil {
		return false
	}

	return *a == *b

}

// IsValid checks for the validity of contained fields.
func (cb CircuitBreakers) IsValid() *ValidationError {
	scope := func(s string) string { return "circuit_breakers." + s }

	errs := &ValidationError{}
	if cb.MaxConnections != nil && *cb.MaxConnections < 0 {
		errs.AddNew(ErrorCase{scope("max_connections"), "must not be negative"})
	}
	if cb.MaxPendingRequests != nil && *cb.MaxPendingRequests < 0 {
		errs.AddNew(ErrorCase{scope("max_pending_requests"), "must not be negative"})
	}
	if cb.MaxRetries != nil && *cb.MaxRetries < 0 {
		errs.AddNew(ErrorCase{scope("max_retries"), "must not be negative"})
	}

	if cb.MaxRequests != nil && *cb.MaxRequests < 0 {
		errs.AddNew(ErrorCase{scope("max_requests"), "must not be negative"})
	}

	return errs.OrNil()

}

// CircuitBreakersPtrEquals provides a way to compare two CircuitBreakers
// pointers
func CircuitBreakersPtrEquals(cb1, cb2 *CircuitBreakers) bool {
	switch {
	case cb1 == nil && cb2 == nil:
		return true
	case cb1 == nil || cb2 == nil:
		return false
	default:
		return cb1.Equals(*cb2)
	}
}
