package api

import (
	"testing"

	"github.com/turbinelabs/test/assert"
)

func getRetryPolicy() (RetryPolicy, RetryPolicy) {
	return RetryPolicy{1, 30, 60},
		RetryPolicy{1, 30, 60}
}

func TestRetryPolicyEquals(t *testing.T) {
	p1, p2 := getRetryPolicy()
	assert.True(t, p1.Equals(p2))
	assert.True(t, p2.Equals(p1))
}

func TestRetryPolicyEqualsNumRetriesChange(t *testing.T) {
	p1, p2 := getRetryPolicy()
	p2.NumRetries++
	assert.False(t, p1.Equals(p2))
	assert.False(t, p2.Equals(p1))
}

func TestRetryPolicyEqualsPerTryTimeoutMsecChange(t *testing.T) {
	p1, p2 := getRetryPolicy()
	p2.PerTryTimeoutMsec++
	assert.False(t, p1.Equals(p2))
	assert.False(t, p2.Equals(p1))
}

func TestRetryPolicyEqualsTimeoutMsecChange(t *testing.T) {
	p1, p2 := getRetryPolicy()
	p2.TimeoutMsec++
	assert.False(t, p1.Equals(p2))
	assert.False(t, p2.Equals(p1))
}

func TestRetryPolicyIsValid(t *testing.T) {
	p, _ := getRetryPolicy()
	assert.Nil(t, p.IsValid())
}

func TestRetryPolicyIsValidBadNumRetries(t *testing.T) {
	p, _ := getRetryPolicy()
	p.NumRetries = -1
	assert.DeepEqual(t, p.IsValid(), &ValidationError{[]ErrorCase{
		{"retry_policy.num_retries", "must not be negative"},
	}})
}

func TestRetryPolicyIsValidBadPerTryTimeoutMsec(t *testing.T) {
	p, _ := getRetryPolicy()
	p.PerTryTimeoutMsec = -1
	assert.DeepEqual(t, p.IsValid(), &ValidationError{[]ErrorCase{
		{"retry_policy.per_try_timeout_msec", "must not be negative"},
	}})
}

func TestRetryPolicyIsValidBadTimeoutMsec(t *testing.T) {
	p, _ := getRetryPolicy()
	p.TimeoutMsec = -1
	assert.DeepEqual(t, p.IsValid(), &ValidationError{[]ErrorCase{
		{"retry_policy.timeout_msec", "must not be negative"},
	}})
}
