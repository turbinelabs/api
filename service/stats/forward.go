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

import (
	"fmt"
)

// A Stat is a named, timestamped, and tagged data point or histogram
type Stat struct {
	Name string `json:"name"`

	// Only one of Value and Histogram may be set.
	Value     *float64   `json:"value,omitempty"`
	Histogram *Histogram `json:"histogram,omitempty"`

	Timestamp int64             `json:"timestamp"` // microseconds since the Unix epoch, UTC
	Tags      map[string]string `json:"tags,omitempty"`
}

func (s Stat) String() string {
	v := "-"
	if s.Value != nil {
		v = fmt.Sprintf("%g (%p)", *s.Value, s.Value)
	}
	return fmt.Sprintf(
		"{Name:%s, Value:%s Histogram:%v, Timestamp:%d, Tags:%s}",
		s.Name,
		v,
		s.Histogram,
		s.Timestamp,
		s.Tags,
	)
}

// A Histogram is a distribution of values into ranges.
type Histogram struct {
	Buckets [][2]float64 `json:"buckets"` // array of [limit, count]
	Count   int64        `json:"count"`
	Sum     float64      `json:"sum"`

	// Non-aggregatable summary fields
	Minimum float64 `json:"min"`
	P50     float64 `json:"p50"`
	P99     float64 `json:"p99"`
	Maximum float64 `json:"max"`
}

// Payload is the payload of a stats update call.
type Payload struct {
	Source string `json:"source"`
	Stats  []Stat `json:"stats"`
}

// ForwardResult is a JSON-encodable struct that encapsulates the result of
// forwarding metrics.
type ForwardResult struct {
	NumAccepted int `json:"numAccepted"`
}
