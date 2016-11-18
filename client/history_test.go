package client

import (
	"encoding/json"
	"fmt"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/turbinelabs/api"
	apihttp "github.com/turbinelabs/api/http"
	"github.com/turbinelabs/api/queryargs"
	"github.com/turbinelabs/api/service"
	"github.com/turbinelabs/api/service/changelog"
	tbntime "github.com/turbinelabs/nonstdlib/time"
	"github.com/turbinelabs/test/assert"
)

var stockCD = []api.ChangeDescription{
	{
		api.ChangeMeta{tbntime.ToUnixMilli(time.Now().UTC()), "some-txn", "org-key", "actor-key", "comment"},
		[]api.ChangeEntry{
			{ObjectKey: "obj-key", Path: "path", Value: "value"},
		},
	},
}

func TestIndex(t *testing.T) {
	filter := changelog.Filter{
		NegativeMatch: true,
		ObjectType:    "foo",
		ObjectKey:     "bar",
	}

	filterStr, err := json.Marshal(filter)
	assert.Nil(t, err)

	verifier := verifyingHandler{
		fn: func(rr apihttp.RichRequest) {
			assert.Equal(t, rr.Underlying().URL.Path, "/v1.0/changelog/adhoc")
			filterGot := rr.QueryArg(queryargs.IndexFilters)
			assert.Equal(t, filterGot, string(filterStr))
		},
		status:   200,
		response: stockCD,
	}

	srv := httptest.NewServer(verifier)
	svc := getAllInterface(srv)
	cds, err := svc.History().Index(filter)
	assert.Nil(t, err)
	assert.DeepEqual(t, cds, stockCD)
}

func TestGraph(t *testing.T) {
	type graphcase struct {
		graphname string
		fn        func(service.All, time.Time, time.Time) ([]api.ChangeDescription, error)
	}

	zt := time.Time{}.UTC()
	tstop := time.Now().UTC()
	tstart := tstop.Add(-1 * time.Hour).UTC()
	tgtKey := "yup-a-key"

	cases := []graphcase{
		{
			"domain",
			func(svc service.All, start, stop time.Time) ([]api.ChangeDescription, error) {
				return svc.History().DomainGraph(api.DomainKey(tgtKey), start, stop)
			},
		},

		{
			"cluster",
			func(svc service.All, start, stop time.Time) ([]api.ChangeDescription, error) {
				return svc.History().ClusterGraph(api.ClusterKey(tgtKey), start, stop)
			},
		},

		{
			"route",
			func(svc service.All, start, stop time.Time) ([]api.ChangeDescription, error) {
				return svc.History().RouteGraph(api.RouteKey(tgtKey), start, stop)
			},
		},
	}

	for _, tc := range cases {
		timeSets := [][]time.Time{
			{tstart, tstop},
			{zt, tstop},
			{tstart, zt},
			{zt, zt},
		}

		for _, ts := range timeSets {
			start := ts[0]
			stop := ts[1]

			assert.Group(
				fmt.Sprintf("test %s graph call (start: %s, stop: %s)", tc.graphname, start, stop),
				t,
				func(t *assert.G) {
					mkCall := func(svc service.All) ([]api.ChangeDescription, error) {
						return tc.fn(svc, start, stop)
					}
					wantPath := fmt.Sprintf("/v1.0/changelog/%s-graph", tc.graphname)

					verifying := verifyingHandler{
						fn: func(rr apihttp.RichRequest) {
							startGot := rr.QueryArg(queryargs.WindowStart)
							stopGot := rr.QueryArg(queryargs.WindowStop)
							assert.Equal(t, rr.Underlying().URL.Path, fmt.Sprintf("%s/%s", wantPath, tgtKey))

							if start.IsZero() {
								assert.Equal(t, startGot, "")
							} else {
								assert.Equal(t, startGot, tbntime.Format(&start))
							}

							if stop.IsZero() {
								assert.Equal(t, stopGot, "")
							} else {
								assert.Equal(t, stopGot, tbntime.Format(&stop))
							}
						},
						status:   200,
						response: stockCD,
					}

					srv := httptest.NewServer(verifying)
					svc := getAllInterface(srv)
					cds, err := mkCall(svc)
					assert.Nil(t, err)
					assert.DeepEqual(t, cds, stockCD)
				},
			)
		}
	}
}