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

package v2

// Valid stat names
const (
	// Client-facing stats
	Requests  = "requests"
	Responses = "responses"
	Latency   = "latency"

	// Upstream stats
	UpstreamRequests  = "us_requests"
	UpstreamResponses = "us_responses"
	UpstreamLatency   = "us_latency"

	// Proxy stats
	Poll           = "poll"
	Config         = "config"
	ConfigLatency  = "config_latency"
	ConfigInterval = "config_interval"
)

// Valid tag names
const (
	// Client and Upstream tags
	Domain     = "domain"
	RouteKey   = "route"
	Rule       = "rule"
	SharedRule = "shared_rule"
	Method     = "method"
	Upstream   = "upstream"
	Instance   = "instance"
	Constraint = "constraint"

	// Additional Responses and UpstreamResponses tags
	StatusCode = "status_code"

	// Poll tags
	PollResult        = "result"
	PollSuccessResult = "success"
	PollErrorResult   = "error"

	// Config tags and values.
	ConfigState   = "state"
	ConfigValid   = "valid"
	ConfigInvalid = "invalid"
	ConfigType    = "type"
)

// DefaultLimitName specifies the name of the default limits. See
// Payload and Histogram.
const DefaultLimitName = "default"

// Histogram represents count of measurements within predefined
// ranges. In addition it contains a count and sum of all measurements
// and the smallest and largest values measured.
type Histogram struct {
	// Limit specifies which entry in the Payload Limits field
	// should be used for this Histogram. If this field is
	// omitted, it is treated as if the value was "default".
	Limit *string `json:"limit,omitempty"`

	// Buckets must contain one value for every entry in the
	// Payload.Limits for the Payload within which this
	// Histogram is contained. Indexes of this field correspond
	// one-to-one with the indexes in Payload.Limits.
	Buckets []int64 `json:"buckets"`

	// Count is the number of items measured by this histogram:
	// equal to the sum of all entries in Buckets plus any
	// measurements that exceeded the last bucket's limit.
	Count int64 `json:"count"`

	// Sum is the total value of all measurements represented by
	// this histogram. Sum ÷ Count is the true average of the
	// measurements.
	Sum float64 `json:"sum"`

	// Minimum is the smallest measurement used in constructing
	// the histogram.
	Minimum float64 `json:"min"`

	// Minimum is the largest measurement used in constructing the
	// histogram.
	Maximum float64 `json:"max"`
}

// Stat is a named, timestamped data point or histogram.
type Stat struct {
	Name string `json:"name"`

	// Only one of Count, Gauge, or Histogram may be set.
	Count     *float64   `json:"count,omitempty"`
	Gauge     *float64   `json:"gauge,omitempty"`
	Histogram *Histogram `json:"histo,omitempty"`

	// Timestamp is milliseconds since the Unix epoch, UTC. Note
	// the change in units from the old Stats struct.
	Timestamp int64 `json:"timestamp"`

	// Tags contains name/value pairs that apply specifically to
	// this stat. Unknown tags are silently rejected.
	Tags map[string]string `json:"tags,omitempty"`
}

// Payload is a stats payload.
type Payload struct {
	// Source is the source of the measurement. Typically, this is the host that
	// generated the Forward request. See Node.
	Source string `json:"source"`

	// Node identifies the host that generated this payload's measurements. If a host
	// delegates forwarding measurements to an aggregator or other intermediary, Node
	// must be set such that for a given metric and tags, the combination of Source
	// and Node is unique for the host and its intermediaries. The Stats API's
	// behavior is non-determinant when two hosts forward measurements with the same
	// metric, tags, node, and source. If omitted, Node is assumed to be the same as
	// Source (e.g., a single host with a unique Source is generating measurements
	// and forwarding them).
	Node *string `json:"node,omitempty"`

	// Zone is the zone within which the measurement took place.
	Zone string `json:"zone"`

	// Proxy is the name of the proxy generating stats. Multiple
	// proxy instances may share a name.
	Proxy *string `json:"proxy,omitempty"`

	// ProxyVersion is the version identifier of the proxy generating stats.
	ProxyVersion *string `json:"proxy_version,omitempty"`

	// Limits are the upper bound values of each histogram
	// bucket. May be omitted if no histograms are present in the
	// payload.
	//
	// Each set of limits is named and referenced by the
	// Histogram instance that uses it (see Histogram.Limit).
	// As a special case, if there is a limit named "default",
	// Histogram entries may omit a name and the default limits
	// will be used.
	//
	// There must be at least two values in the limits array. The
	// number of values in each Histogram.Buckets field must be
	// the same as the number of values in the corresponding entry
	// in this map. The limit value array must be sorted in
	// ascending order.
	Limits map[string][]float64 `json:"limits,omitempty"`

	// Stats is an array of measurements in this payload.
	Stats []Stat `json:"stats"`
}

// ForwardResult is a JSON-encodable struct that encapsulates the result of
// forwarding metrics.
type ForwardResult struct {
	NumAccepted int `json:"numAccepted"`
}
