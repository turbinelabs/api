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
	v1 "github.com/turbinelabs/api/service/stats/v1"
	v1querytype "github.com/turbinelabs/api/service/stats/v1/querytype"
	v1timegranularity "github.com/turbinelabs/api/service/stats/v1/timegranularity"
)

// Query is an alias for the V1 stats API Query type.
type Query = v1.Query

// QueryType is an alias for the V1 stats API QueryType type.
type QueryType = v1querytype.QueryType

// QueryTimeSeries is an alias for the V1 stats API QueryTimeSeries type.
type QueryTimeSeries = v1.QueryTimeSeries

// TimeRange is an alias for the V1 stats API TimeRange type.
type TimeRange = v1.TimeRange

// SimpleTimeRange is an alias for the V1 stats API SimpleTimeRange type.
type SimpleTimeRange = v1.SimpleTimeRange

// TimeGranularity is an alis for the V1 stats API TimeGranularity type.
type TimeGranularity = v1timegranularity.TimeGranularity

// ZeroFill is an alias for the V1 stats API ZeroFill type.
type ZeroFill = v1.ZeroFill

// QueryResult is an alias for the V1 stats API QueryResult type.
type QueryResult = v1.QueryResult

// TimeSeries is an alias for the V1 stats API TimeSeries type.
type TimeSeries = v1.TimeSeries

// Point is an alias for the V1 stats API Point type.
type Point = v1.Point
