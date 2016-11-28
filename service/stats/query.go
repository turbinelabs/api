package stats

import (
	"github.com/turbinelabs/api"
	"github.com/turbinelabs/api/service/stats/querytype"
	"github.com/turbinelabs/api/service/stats/timegranularity"
)

type TimeRange struct {
	// Start and End represent the start and end of a time range,
	// specified in microseconds since the Unix epoch, UTC. End
	// takes precedence over Duration.
	Start *int64 `json:"start,omitempty" form:"start"`
	End   *int64 `json:"end,omitempty" form:"end"`

	// Duration specifies how long a time span of stats data to
	// return in microseconds. End takes precedence over
	// Duration. If Start is specified, Duration sets the end of
	// the time span (e.g. from Start for Duration
	// microseconds). If Start is not specified, Duration sets the
	// start of the time span that many microseconds into the past
	// (e.g., Duration microseconds ago, until now).
	Duration *int64 `json:"duration,omitempty" form:"duration"`

	// Granularity specifies how much time each data point
	// represents. If absent, it defaults to "seconds". Valid
	// values are "seconds", "minutes", or "hours".
	Granularity timegranularity.TimeGranularity `json:"granularity" form:"granularity"`
}

type QueryTimeSeries struct {
	// Specifies a name for this timeseries query. It may be used
	// to assist in identifying the corresponding data in the
	// response object.
	Name string `json:"name,omitempty" form:"name"`

	// Specifies the type of data returned. Required.
	QueryType querytype.QueryType `json:"query_type" form:"query_type"`

	// Specifies the domain host for which stats are returned. The
	// host may be just a domain name (e.g., "example.com"), or a
	// domain name and port (e.g., "example.com:443"). The former
	// aggregates stats across all ports serving the domain. If
	// DomainHost is not specified, stats are aggregated across
	// all domains.
	DomainHost *string `json:"domain_host,omitempty" form:"domain_host"`

	// Specifies the RouteKey for which stats are returned. If
	// not specified, stats are aggregated across routes.
	RouteKey *api.RouteKey `json:"route_key,omitempty" form:"route_key"`

	// Specifies the SharedRule name for which stats are
	// returned. If not specified, stats are aggregated across
	// shared rules.
	SharedRuleName *string `json:"shared_rule_name,omitempty" form:"shared_rule_name"`

	// Specifies the RuleKey for which stats are returned.
	// Requires that a RouteKey or SharedRuleName is given. If not
	// specified, stats are aggregated across rules.
	RuleKey *api.RuleKey `json:"rule_key,omitempty" form:"rule_key"`

	// Specifies the HTTP method for which stats are returned. If
	// not specified, stats are aggregated across methods.
	Method *string `json:"method,omitempty" form:"method"`

	// Specifies the Cluster name for which stats are returned. If
	// not specified, stats are aggregated across clusters.
	ClusterName *string `json:"cluster_name,omitempty" form:"cluster_name"`

	// Specifies the Instance keys (host:port) for which stats are
	// returned. If empty, stats are aggregated across all
	// instances. If one ore more instances are given, stats are
	// aggregated across only those instances.
	InstanceKeys []string `json:"instance_keys,omitempty" form:"instance_keys"`
}

type Query struct {
	// Specifies the zone name for which stats are
	// queried. Required.
	ZoneName string `json:"zone_name" form:"zone_name"`

	// Specifies the time range of the query. Defaults to the last
	// hour.
	TimeRange TimeRange `json:"time_range" form:"time_range"`

	// Specifies one or more queries to execute against the given
	// zone and time range.
	TimeSeries []QueryTimeSeries `json:"timeseries" form:"timeseries"`
}

type Point struct {
	// A data point.
	Value float64 `json:"value"`

	// Collection timestamp in microseconds since the Unix epoch,
	// UTC. N.B. that the actual resolution of the timestamp may
	// be less granular than microseconds.
	//
	// Microsecond resolution timestamps with an epoch of
	// 1970-01-01 00:00:00 reach 2^53 - 1, the maximum integer
	// exactly representable in Javascript, some time in 2255:
	// (2^53 - 1) / (86400 * 1000 * 1000)
	//     = 10249.99 days / 365.24
	//     = 285.42 years
	Timestamp int64 `json:"timestamp"`
}

type TimeSeries struct {
	// The QueryTimeSeries object corresponding to the data
	// points.
	Query QueryTimeSeries `json:"query"`

	// The data points that represent the time series.
	Points []Point `json:"points"`
}

type QueryResult struct {
	// The TimeRange used to issue this query. The object is
	// normalized such that all of its fields are set and
	// consistent.
	TimeRange TimeRange `json:"time_range"`

	// Represents the timeseries returned by the query. The order
	// of returned TimeSeries values matches the order of the
	// original QueryTimeSeries values in the request.
	TimeSeries []TimeSeries `json:"timeseries"`
}
