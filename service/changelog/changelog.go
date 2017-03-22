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

// Package changelog includes the filter definition necessary to make adhoc
// queries against the object audit history that we track.
//
// changelog filtering is represented by a Sum-of-Products format. This is
// sufficient to encode any boolean expression. In short this means that a query
// is composed of a series of logical intersections (ANDs) that are unioned
// (ORed) together. We work with the limitation that only individual filters may
// be negated which does not limit expressivity.
package changelog

import (
	"time"

	"github.com/turbinelabs/api"
	"github.com/turbinelabs/api/changetype"
	"github.com/turbinelabs/nonstdlib/ptr"
	tbntime "github.com/turbinelabs/nonstdlib/time"
)

type FilterFn func(Filter) (Filter, error)

// FilterExpr is a form of a filter than can be used to search for change logs.
type FilterExpr interface {
	// Convert the implementing type into a format that is suitable for evaluation.
	AsExpr() FilterOrs

	// ApplyAll applies some function on all Filters within an expression and
	// returns a modified expression or error if the function could not be applied.
	ApplyAll(FilterFn) (FilterExpr, error)
}

// TimeRange describes an inclusive window during which we sholud be looking
// for changes. If a Start or End time is not set then we assume the TimeRange
// provides only one bound and the other is computed to be an hour before or
// after, as apporpriate. Start marks the earlier of the two times and is
// inclusive, end marks the latest of the times and is exclusive. In other
// words queries with a TimeRange will be interpreted as [Start, End).
//
// Start and End represent the start and end of a time range, specified in
// microseconds since the Unix epoch, UTC.
type TimeRange struct {
	Start *int64 `json:"start,omitempty" form:"start"`
	End   *int64 `json:"end,omitempty"   form:"end"`
}

// StartNano returns the range start time, if set, in nanoseconds.
func (tr TimeRange) StartNano() *int64 {
	if tr.Start != nil {
		return ptr.Int64(1000 * *tr.Start)
	}

	return nil
}

// EndNano returns the range end time, if set, in nanoseconds.
func (tr TimeRange) EndNano() *int64 {
	if tr.End != nil {
		return ptr.Int64(1000 * *tr.End)
	}

	return nil
}

// SetStart sets the range start to the specified time in microseconds since
// the Unix epoch UTC.
func (tr *TimeRange) SetStart(t time.Time) {
	if t.IsZero() {
		tr.Start = nil
		return
	}
	tr.Start = ptr.Int64(t.UnixNano() / 1000)
}

// SetEnd sets the range end to the specified time in microseconds since the
// Unix epoch UTC.
func (tr *TimeRange) SetEnd(t time.Time) {
	if t.IsZero() {
		tr.End = nil
		return
	}
	tr.End = ptr.Int64(t.UnixNano() / 1000)
}

// StartTime returns the start time of a range or nil if one is not set.
func (tr TimeRange) StartTime() *time.Time {
	if tr.Start != nil {
		return ptr.Time(tbntime.FromUnixMicro(*tr.Start))
	}

	return nil
}

// EndTime returns the end time of a range or nil if one is not set.
func (tr TimeRange) EndTime() *time.Time {
	if tr.End != nil {
		return ptr.Time(tbntime.FromUnixMicro(*tr.End))
	}

	return nil
}

// FieldFilter describes a specific attribute change on a tracked object. The
// AttributePath is constructed as a dot-separated collection of field names
// rooted with the type of the object being updated. All containers are
// key-indexed and changes within a container are specified using index
// operators [].  For example: a cluster that has an instance running on
// 'smf1-s23' and port 9990 would be have changes rooted on the path
// 'cluster.instance[smf1-s23:9990]' the key construction varies by the type of
// the contained data. As data is added or removed it is tracked on the
// non-indexed path (cluster.instance from our example above) with the new /
// deleted key as the value.
//
// Before, if set, limits results to changes that moved from this value
// After, if set, limits results to changes that set this value
type FieldFilter struct {
	AbsoluteMatchOnly  bool       `json:"absolute_match_only"  form:"absolute_match_only"`
	AttributePath      string     `json:"attribute_path"       form:"attribute_path"`
	ChangeType         ChangeType `json:"change_type"          form:"change_type"`
	AttributeValue     *string    `json:"attribute_value"      form:"attribute_value"`
	ExcludeEmptyValues bool       `json:"exclude_empty_values" form:"exclude_empty_values"`
}

type ChangeType string

var (
	ValueAdded   = ChangeType(changetype.Addition.Name)
	ValueRemoved = ChangeType(changetype.Removal.Name)
)

// Filter collects all the attributes we have available for searching audit logs.
// Additionally the field NegativeMatch may be set to indicate that this filter
// should used to exlude matching change entries. Note that NegativeMatch does
// not apply to TimeRange.
type Filter struct {
	TimeRange     TimeRange   `json:"time_range"     form:"time_range"`
	NegativeMatch bool        `json:"negative_match" form:"negative_match"`
	ObjectType    string      `json:"object_type"    form:"object_type"`
	ObjectKey     string      `json:"object_key"     form:"object_key"`
	ChangeTxn     string      `json:"change_txn"     form:"change_txn"`
	ZoneKey       api.ZoneKey `json:"zone_key"       form:"zone_key"`
	OrgKey        api.OrgKey  `json:"org_key"        form:"org_key"`
	Actor         api.UserKey `json:"actor_key"      form:"actor_key"`
	FieldFilter   `form:"field"`
}

func (f Filter) AsExpr() FilterOrs {
	return NewFilterUnion(f)
}

func (f Filter) ApplyAll(fn FilterFn) (FilterExpr, error) {
	f, err := fn(f)
	if err != nil {
		return nil, err
	}
	return f, nil
}

// FilterOrs represents a collection of logical intersections that will be
// evaluated as a union.
type FilterOrs struct {
	FilterAnds []FilterAnds `json:"or" form:"or"`
}

func (fs FilterOrs) AsExpr() FilterOrs {
	return fs
}

func (fs FilterOrs) ApplyAll(fn FilterFn) (FilterExpr, error) {
	for _, ands := range fs.FilterAnds {
		for j, f := range ands.Filters {
			f, err := fn(f)
			if err != nil {
				return nil, err
			}
			ands.Filters[j] = f
		}
	}

	return fs, nil
}

// FilterAnds represents a collection of filters that will be evaluated by
// ANDing the contents together.
type FilterAnds struct {
	Filters []Filter `json:"and" form:"and"`
}

func (fs FilterAnds) AsExpr() FilterOrs {
	return FilterOrs{[]FilterAnds{fs}}
}

func (fs FilterAnds) ApplyAll(fn FilterFn) (FilterExpr, error) {
	return fs.AsExpr().ApplyAll(fn)
}

// NewFilterIntersection is a convenience function that construts a FilterAnds
// from a set of Filters.
func NewFilterIntersection(fs ...Filter) FilterAnds {
	return FilterAnds{fs}
}

// NewFilterUnion is a convenience function that constructs a statement that ORs
// together a series of Filters.
func NewFilterUnion(fs ...Filter) FilterOrs {
	ands := make([]FilterAnds, len(fs))
	for i, f := range fs {
		ands[i] = NewFilterIntersection(f)
	}
	return FilterOrs{ands}
}
