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
