/*
Copyright 2017-2018 Turbine Labs, Inc.

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

import (
	"github.com/turbinelabs/api/service/stats/v2/timegranularity"
)

// TimeRange specifies a range of time over which stats queries are performed.
type TimeRange struct {
	// SimpleTimeRange specifies the window this time range covers.
	SimpleTimeRange

	// Duration specifies how long a time span of stats data to return. In the v1.0
	// API, it is specified in microseconds. In the v2.0 API, it is specified in
	// seconds. End takes precedence over Duration. If Start is specified, Duration
	// sets the end of the time span (e.g. from Start for a period of Duration). If
	// Start is not specified, Duration sets the start of the time span that period
	// into the past (e.g., a period lasting Duration, until now).
	Duration *int64 `json:"duration,omitempty" form:"duration"`

	// Granularity specifies how much time each data point represents. If absent, it
	// defaults to "minutes". Valid values are minutes" or "hours".
	Granularity timegranularity.TimeGranularity `json:"granularity" form:"granularity"`
}

// SimpleTimeRange represents the start and end of a time range.In the v2.0 API,
// times are specified in seconds since the Unix epoch, UTC.
type SimpleTimeRange struct {
	// Start indicates when data should begin being reported.
	Start *int64 `json:"start,omitempty" form:"start"`

	// End specifies when data is no longer desired.
	End *int64 `json:"end,omitempty" form:"end"`
}
