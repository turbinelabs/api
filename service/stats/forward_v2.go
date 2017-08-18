/*
Copyright 2017 Turbine Labs, Inc.

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

package stats

// Valid stat names.
const (
	Requests  = "requests"
	Responses = "responses"
	Latency   = "latency"

	UpstreamRequests  = "us_requests"
	UpstreamResponses = "us_responses"
	UpstreamLatency   = "us_latency"

	InvalidConfig = "invalid_config"
	PollSuccess   = "poll_success"
)

// Valid tag names.
const (
	Domain     = "domain"      // valid for Upstream* stats only
	RouteKey   = "route"       // "
	Rule       = "rule"        // "
	SharedRule = "shared_rule" // "
	Method     = "method"      // "
	Upstream   = "upstream"    // "
	Instance   = "instance"    // "
	StatusCode = "status_code" // valid for Responses and UpstreamResponses only
)

// DefaultLimitName specifies the name of the default limits. See
// PayloadV2 and HistogramV2.
const DefaultLimitName = "default"

// HistogramV2 represents count of measurements within predefined
// ranges. In addition it contains a count and sum of all measurements
// and the smallest and largest values measured.
type HistogramV2 struct {
	// Limit specifies which entry in the PayloadV2 Limits field
	// should be used for this Histogram. If this field is
	// omitted, it is treated as if the value was "default".
	Limit *string `json:"limit,omitempty"`

	// Buckets must contain one value for every entry in the
	// PayloadV2.Limits for the PayloadV2 within which this
	// HistogramV2 is contained. Indexes of this field correspond
	// one-to-one with the indexes in PayloadV2.Limits.
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

// StatV2 is a named, timestamped data point or histogram.
type StatV2 struct {
	Name string `json:"name"`

	// Only one of Count, Gauge, or Histogram may be set.
	Count     *float64     `json:"count,omitempty"`
	Gauge     *float64     `json:"gauge,omitempty"`
	Histogram *HistogramV2 `json:"histo,omitempty"`

	// Timestamp is milliseconds since the Unix epoch, UTC. Note
	// the change in units from the old Stats struct.
	Timestamp int64 `json:"timestamp"`

	// Tags contains name/value pairs that apply specifically to
	// this stat. Unknown tags are silently rejected.
	Tags map[string]string `json:"tags,omitempty"`
}

// PayloadV2 is a stats payload.
type PayloadV2 struct {
	// Source is the source of the measurement.
	Source string `json:"source"`

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
	// HistogramV2 instance that uses it (see HistogramV2.Limit).
	// As a special case, if there is a limit named "default",
	// HistogramV2 entries may omit a name and the default limits
	// will be used.
	//
	// There must be at least two values in the limits array. The
	// number of values in each HistogramV2.Buckets field must be
	// the same as the number of values in the corresponding entry
	// in this map. The limit value array must be sorted in
	// ascending order.
	Limits map[string][]float64 `json:"limits,omitempty"`

	// Stats is an array of measurements in this payload.
	Stats []StatV2 `json:"stats"`
}
