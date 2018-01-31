package api

// RetryPolicy specifies the number of times to retry a request and how long to
// wait before timing out.
type RetryPolicy struct {
	// Number of times to retry an upstream request. Note that the initial
	// connection attempt is not included in this number, hence 0 means initial
	// attempt and no retries, and 1 means initial attempt plus one retry.
	NumRetries int `json:"num_retries"`
	// Time limit in milliseconds for a single attempt.
	PerTryTimeoutMsec int `json:"per_try_timeout_msec"`
	// Total time limit in milliseconds for all attempts (including the initial
	// attempt).
	TimeoutMsec int `json:"timeout_msec"`
}

// Checks for exact equality between this retry policy and another. Exact
// equality means each field must be equal (== or Equal, as appropriate) to the
// corresponding field in the parameter.
func (p RetryPolicy) Equals(o RetryPolicy) bool {
	return p == o
}

// Convenience function for calling Equals when you have pointers to two retry
// policies.
func RetryPolicyEquals(a, b *RetryPolicy) bool {
	if a == b {
		return true
	}
	if a == nil || b == nil {
		return false
	}
	return a.Equals(*b)
}

// Checks validity of a retry policy. For a retry policy to be valid it must
// have non-negative NumRetries, PerTryTimeoutMsec, and TimeoutMsec.
func (p RetryPolicy) IsValid() *ValidationError {
	scope := func(s string) string { return "retry_policy." + s }

	errs := &ValidationError{}
	if p.NumRetries < 0 {
		errs.AddNew(ErrorCase{scope("num_retries"), "must not be negative"})
	}
	if p.PerTryTimeoutMsec < 0 {
		errs.AddNew(ErrorCase{scope("per_try_timeout_msec"), "must not be negative"})
	}
	if p.TimeoutMsec < 0 {
		errs.AddNew(ErrorCase{scope("timeout_msec"), "must not be negative"})
	}

	return errs.OrNil()
}
