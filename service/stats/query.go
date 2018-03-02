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

package stats

import (
	v2 "github.com/turbinelabs/api/service/stats/v2"
	v2querytype "github.com/turbinelabs/api/service/stats/v2/querytype"
	v2timegranularity "github.com/turbinelabs/api/service/stats/v2/timegranularity"
)

// Query is an alias for the V2 stats API Query type.
type Query = v2.Query

// QueryType is an alias for the V2 stats API QueryType type.
type QueryType = v2querytype.QueryType

// QueryTimeSeries is an alias for the V2 stats API QueryTimeSeries type.
type QueryTimeSeries = v2.QueryTimeSeries

// TimeRange is an alias for the V2 stats API TimeRange type.
type TimeRange = v2.TimeRange

// SimpleTimeRange is an alias for the V2 stats API SimpleTimeRange type.
type SimpleTimeRange = v2.SimpleTimeRange

// TimeGranularity is an alis for the V2 stats API TimeGranularity type.
type TimeGranularity = v2timegranularity.TimeGranularity

// ZeroFill is an alias for the V2 stats API ZeroFill type.
type ZeroFill = v2.ZeroFill

// QueryResult is an alias for the V2 stats API QueryResult type.
type QueryResult = v2.QueryResult

// TimeSeries is an alias for the V2 stats API TimeSeries type.
type TimeSeries = v2.TimeSeries

// Point is an alias for the V2 stats API Point type.
type Point = v2.Point
