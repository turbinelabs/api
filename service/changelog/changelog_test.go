package changelog

import (
	"errors"
	"testing"
	"time"

	"github.com/turbinelabs/api"
	"github.com/turbinelabs/nonstdlib/ptr"
	tbntime "github.com/turbinelabs/nonstdlib/time"
	"github.com/turbinelabs/test/assert"
)

var (
	f1 = Filter{NegativeMatch: true, ObjectType: "cluster", ObjectKey: "whee"}
	f2 = Filter{ObjectType: "cluster", Actor: api.UserKey("bob")}
	f3 = Filter{FieldFilter: FieldFilter{AttributePath: "a.b.c"}}
)

func TestFilterAsExpr(t *testing.T) {
	fe := f1.AsExpr()
	assert.DeepEqual(t, fe, FilterOrs{[]FilterAnds{{[]Filter{f1}}}})
}

func TestFilterAndsAsExpr(t *testing.T) {
	fas := FilterAnds{[]Filter{f1, f2, f3}}
	assert.DeepEqual(t, fas.AsExpr(), FilterOrs{[]FilterAnds{fas}})
}

func TestFilterOrsAsExpr(t *testing.T) {
	fos := FilterOrs{[]FilterAnds{{[]Filter{f1, f2}}, {[]Filter{f3}}}}
	assert.DeepEqual(t, fos.AsExpr(), fos)
}

func TestNewFilterIntersection(t *testing.T) {
	fas := NewFilterIntersection(f1, f2, f3)
	assert.DeepEqual(t, fas, FilterAnds{[]Filter{f1, f2, f3}})
}

func TestNewFilterUnion(t *testing.T) {
	fus := NewFilterUnion(f1, f2, f3)
	assert.DeepEqual(t, fus, FilterOrs{
		[]FilterAnds{{[]Filter{f1}}, {[]Filter{f2}}, {[]Filter{f3}}},
	})
}

var ok = api.OrgKey("whee")

func setOrgKey(f Filter) (Filter, error) {
	f.OrgKey = ok
	return f, nil
}

var fnErr = errors.New("boom")

func errOut(_ Filter) (Filter, error) {
	return Filter{}, fnErr
}

func verifyOrgKey(t *testing.T, got FilterExpr, gotErr error) {
	for _, ands := range got.AsExpr().FilterAnds {
		for _, f := range ands.Filters {
			assert.Equal(t, f.OrgKey, ok)
		}
	}
	assert.Nil(t, gotErr)
}

func verifyGotErr(t *testing.T, got FilterExpr, gotErr error) {
	assert.Nil(t, got)
	assert.Equal(t, gotErr, fnErr)
}

func TestApplyAllOrs(t *testing.T) {
	fs := FilterOrs{
		[]FilterAnds{
			{[]Filter{f1, f2}},
			{[]Filter{f3}},
		},
	}

	got, gotErr := fs.ApplyAll(setOrgKey)

	verifyOrgKey(t, got, gotErr)
}

func TestApplyAllOrsErr(t *testing.T) {
	fs := FilterOrs{
		[]FilterAnds{
			{[]Filter{f1, f2}},
			{[]Filter{f3}},
		},
	}

	got, gotErr := fs.ApplyAll(errOut)
	verifyGotErr(t, got, gotErr)
}

func TestApplyAllAnds(t *testing.T) {
	fs := FilterAnds{[]Filter{f1, f2}}

	got, gotErr := fs.ApplyAll(setOrgKey)
	verifyOrgKey(t, got, gotErr)
}

func TestApplyAllAndsErr(t *testing.T) {
	fs := FilterAnds{[]Filter{f1, f2}}

	got, gotErr := fs.ApplyAll(errOut)
	verifyGotErr(t, got, gotErr)
}

func TestApplyAllFilter(t *testing.T) {
	got, gotErr := f1.ApplyAll(setOrgKey)
	verifyOrgKey(t, got, gotErr)
}

func TestApplyAllFilterErr(t *testing.T) {
	got, gotErr := f1.ApplyAll(errOut)
	verifyGotErr(t, got, gotErr)
}

type testStandarizeTime struct {
	hasWindowStart bool
	hasWindowEnd   bool
	hasRangeStart  bool
	hasRangeEnd    bool
	newStart       bool
	newEnd         bool
}

func TestTimeRangeSetStart(t *testing.T) {
	tr := TimeRange{}
	n := time.Now()
	tr.SetStart(n)
	assert.DeepEqual(t, tr.Start, ptr.Int64(tbntime.ToUnixMicro(n)))
}

func TestTimeRangeSetEnd(t *testing.T) {
	tr := TimeRange{}
	n := time.Now()
	tr.SetEnd(n)
	assert.DeepEqual(t, tr.End, ptr.Int64(tbntime.ToUnixMicro(n)))
}

func TestTimeRangeStartNano(t *testing.T) {
	tr := TimeRange{}
	n := time.Now()
	tr.SetStart(n)
	assert.DeepEqual(t, tr.StartNano(), ptr.Int64(n.UnixNano()/1000*1000))
}

func TestTimeRangeEndNano(t *testing.T) {
	tr := TimeRange{}
	n := time.Now()
	tr.SetEnd(n)
	assert.DeepEqual(t, tr.EndNano(), ptr.Int64(n.UnixNano()/1000*1000))
}

func TestTimeRangeStartNanoNil(t *testing.T) {
	tr := TimeRange{}
	assert.Nil(t, tr.StartNano())
}

func TestTimeRangeEndNanoNil(t *testing.T) {
	tr := TimeRange{}
	assert.Nil(t, tr.EndNano())
}

func TestTimeRangeStartTime(t *testing.T) {
	tn := time.Now()

	tr := TimeRange{}
	tr.SetStart(tn)

	assert.DeepEqual(t, tr.StartTime(), ptr.Time(tbntime.TruncUnixMicro(tn)))
}

func TestTimeRangeStartTimeNil(t *testing.T) {
	ts := TimeRange{nil, nil}.StartTime()
	assert.Nil(t, ts)
}

func TestTimeRangeEndTime(t *testing.T) {
	tn := time.Now()

	tr := TimeRange{}
	tr.SetEnd(tn)

	assert.DeepEqual(t, tr.EndTime(), ptr.Time(tbntime.TruncUnixMicro(tn)))
}

func TestTimeRangeEndTimeNil(t *testing.T) {
	ts := TimeRange{nil, nil}.EndTime()
	assert.Nil(t, ts)
}
