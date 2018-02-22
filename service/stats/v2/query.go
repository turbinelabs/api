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

import (
	"github.com/turbinelabs/api"
	"github.com/turbinelabs/api/service/stats/v2/querytype"
)

// Query represents a query in the v2.0 stats API.
type Query struct {
	// Specifies the time range of the query. Defaults to the last hour.
	TimeRange TimeRange `json:"time_range" form:"time_range"`

	// ZeroFill, if set, controls how the stats API fills in values for timestamps
	// that have no value. If all values are filled then an additional EmptySeries
	// field will be set on the response TimeSeriesV2.
	ZeroFill *ZeroFill `json:"zero_fill,omitempty" form:"zero_fill"`

	// Filters specifies named filters that can be referenced from QueryTimeSeriesV2
	// via the FilterName field.
	Filters map[string]QueryFilter `json:"filters,omitempty" form:"filters"`

	// Specifies one or more queries to execute.
	TimeSeries []QueryTimeSeries `json:"timeseries" form:"timeseries"`
}

// QueryTimeSeries represents a query time series in the v2.0 stats API.
type QueryTimeSeries struct {
	// Specifies a name for this time series query. It may be used to assist in
	// identifying the corresponding data in the response object.
	Name string `json:"name,omitempty" form:"name"`

	// Specifies the type of data returned. Required. Some values of QueryType are
	// not supported in the v2.0 API.
	QueryType querytype.QueryType `json:"query_type" form:"query_type"`

	// TimeRangeOverride is a way to specify that this QueryTimeSeries should limit
	// the window that data is returned for beyond the global TimeRange specified in
	// Query. It will share the same granularity as the parent TimeRange. If either
	// Start or End exceed the parent range the query will be rejected. If only one
	// Start or End value is specified then the other will be inferred from the
	// parent range. If both Start and End are nil the query will be rejected.
	TimeRangeOverride *SimpleTimeRange `json:"time_range,omitempty" form:"time_range"`

	// ZeroFillDefault allows a query to override the default value specified when
	// populating series points if there is no data. If unspecified 0 will be used.
	ZeroFillDefault *float64 `json:"zero_fill_default" form:"zero_fill_default"`

	// FilterName allows a pre-defined filter to be selected. The Filter field is
	// ignored if FilterName is set.
	FilterName *string `json:"filter_name,omitempty" form:"filter_name"`

	// Filter specifies a set of filter parameters for this QueryTimeSeries. This
	// field is ignored if FilterName is set.
	Filter *QueryFilter `json:"filter,omitempty" form:"filter"`
}

// QueryFilter represents a filter applied a QueryTimeSeries in the v2.0 stats
// API.
type QueryFilter struct {
	// ZoneName specifies the zone name for which stats are queried. Required for
	// upstream QueryTypes. Optional for downstream QueryTypes.
	ZoneName *string `json:"zone_name,omitempty" form:"zone_name"`

	// DomainHost specifies the domain host for which stats are returned. The host
	// may be just a domain name (e.g., "example.com"), or a domain name and port
	// (e.g., "example.com:443"). The former aggregates stats across all ports
	// serving the domain. If DomainHost is not specified, stats are aggregated
	// across all domains.
	DomainHost *string `json:"domain_host,omitempty" form:"domain_host"`

	// RouteKey specifies the RouteKey for which stats are returned. If not
	// specified, stats are aggregated across routes.
	RouteKey *api.RouteKey `json:"route_key,omitempty" form:"route_key"`

	// SharedRuleName specifies the SharedRule name for which stats are returned. If
	// not specified, stats are aggregated across shared rules.
	SharedRuleName *string `json:"shared_rule_name,omitempty" form:"shared_rule_name"`

	// RuleKey specifies the RuleKey for which stats are returned. If set, a RouteKey
	// or SharedRuleName must also be given. If not specified, stats are aggregated
	// across rules.
	RuleKey *api.RuleKey `json:"rule_key,omitempty" form:"rule_key"`

	// Method specifies the HTTP method for which stats are returned. If not
	// specified, stats are aggregated across methods.
	Method *string `json:"method,omitempty" form:"method"`

	// ClusterName specifies the Cluster name for which stats are returned. If not
	// specified, stats are aggregated across clusters.
	ClusterName *string `json:"cluster_name,omitempty" form:"cluster_name"`

	// InstanceKeys specifies the Instance keys (host:port) for which stats are
	// returned. If empty, stats are aggregated across all instances. If one or more
	// instances are given, stats are aggregated across only those instances.
	InstanceKeys []string `json:"instance_keys,omitempty" form:"instance_keys"`

	// StatusCodes filters responses by status code for the ResponsesForCode and
	// DownstreamResponsesForCode QueryTypes only. Valid values are exact HTTP status
	// codes (e.g., "404", "503") or a single numeric digit followed by "xx" to
	// represent a grouping of status codes (e.g., "2xx" for success or "5xx" for
	// server errors). Results are aggregated across all status codes given. Required
	// when the QueryType is ResponsesForyCode or DownstreamResposnesForCode. Ignored
	// for all other QueryTypes.
	StatusCodes []string `json:"status_codes,omitempty" form:"status_codes"`
}

// Point represents a data point in a timeseries result in the both the v2.0 stats
// API. Note that the the definition of Timestamp varies across versions.
type Point struct {
	// A data point.
	Value float64 `json:"value"`

	// Collection timestamp since the Unix epoch, UTC. In the v2.0 API these values
	// are in units of seconds. Note that the actual resolution of the timestamp may
	// be less granular than the unit.
	Timestamp int64 `json:"timestamp"`
}

// TimeSeries represents a result time series in the v2.0 stats API.
type TimeSeries struct {
	// The QueryTimeSeries object corresponding to the data points.
	Query QueryTimeSeries `json:"query"`

	// EmptySeries is true if the Query ZeroFill field was set to Full and all
	// Points were filled with the QueryTimeSeries ZeroFillDefault.
	EmptySeries *bool `json:"empty_series,omitempty"`

	// The data points that represent the time series.
	Points []Point `json:"points"`
}

// QueryResult represents a query result in the v2.0 stats API.
type QueryResult struct {
	// The TimeRange used to issue this query. The object is normalized such that all
	// of its fields are set and consistent.
	TimeRange TimeRange `json:"time_range"`

	// Represents the timeseries returned by the query. The order of returned
	// TimeSeries values matches the order of the original QueryTimeSeries values
	// in the request.
	TimeSeries []TimeSeries `json:"timeseries"`
}
