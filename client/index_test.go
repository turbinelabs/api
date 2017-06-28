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

package client

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/turbinelabs/api"
	apihttp "github.com/turbinelabs/api/http"
	"github.com/turbinelabs/api/http/envelope"
	httperr "github.com/turbinelabs/api/http/error"
	"github.com/turbinelabs/api/queryargs"
	"github.com/turbinelabs/api/service"
	"github.com/turbinelabs/test/assert"
	"github.com/turbinelabs/test/check"
)

type MkSvcDoIndexFn func(*testing.T, interface{}, *httptest.Server) (interface{}, error)
type AssertEqualityFn func(testing.TB, interface{}, interface{}) bool

type indexTestCase struct {
	callSvc        MkSvcDoIndexFn   // function that calls index on a service
	assertEquality AssertEqualityFn // function that can assert equality between two responses
	urlPrefix      string           // what url prefix should we expect
	filters        interface{}      // what filter we're passing, if any
	responseObj    interface{}      // will be rendered as json and set as the response body
	wantResp       interface{}      // what the API call should produce
	wantErr        error            // what error, if any, the API call should produce
}

func assertIndexURL(t *testing.T, inurl *url.URL, prefix, filterStr string) {
	if !assertURLPrefix(t, inurl.Path, prefix) {
		// if it's not a cluster url then we can just give up b/c we're about to
		// panic otherwise
		return
	}

	// the path should be fully consumed by the cluster url
	remainder := stripURLPrefix(inurl.Path, prefix)
	assert.Equal(t, len(remainder), 0)

	// by default we don't want anything
	wantStr := ""
	if filterStr != "" {
		// but if the encoded filters are non-nil then we want it to be included as
		// an arg for the correct parameter
		wantStr = queryargs.IndexFilters + "=" + url.QueryEscape(filterStr)
	}
	assert.Equal(t, inurl.RawQuery, wantStr)
}

func (tc indexTestCase) run(t *testing.T) {
	wantFilterStr := ""

	if !check.IsNil(tc.filters) {
		filterB, e := json.Marshal(tc.filters)
		if !assert.Nil(t, e) {
			log.Fatal(e)
		}
		wantFilterStr = string(filterB)
	}

	verifier := verifyingHandler{
		fn: func(rr apihttp.RichRequest) {
			assertIndexURL(t, rr.Underlying().URL, tc.urlPrefix, wantFilterStr)
		},
		status:   http.StatusOK,
		response: tc.responseObj,
	}

	s := httptest.NewServer(verifier)
	defer s.Close()

	if tc.wantErr != nil {
		if e, ok := tc.wantErr.(*httperr.Error); ok {
			e.Message = strings.Replace(e.Message, "{{URL}}", s.URL, -1)
		}
	}

	got, gotErr := tc.callSvc(t, tc.filters, s)
	assert.DeepEqual(t, gotErr, tc.wantErr)

	if tc.wantResp == nil {
		assert.Nil(t, got)
		return
	}

	tc.assertEquality(t, got, tc.wantResp)
}

func callClusterIndex(
	t *testing.T, filters interface{}, server *httptest.Server,
) (interface{}, error) {
	arg, ok := filters.([]service.ClusterFilter)
	assert.True(t, ok)
	svc := getAllInterface(server).Cluster()
	return svc.Index(arg...)
}

type clusterIndexTest struct {
	filters     []service.ClusterFilter
	responseObj interface{}
	wantResp    interface{}
	wantErr     error
}

func (tc clusterIndexTest) run(t *testing.T) {
	indexTestCase{
		callSvc:        callClusterIndex,
		assertEquality: assert.HasSameElements,
		urlPrefix:      clusterCommonURL,
		filters:        tc.filters,
		responseObj:    tc.responseObj,
		wantResp:       tc.wantResp,
		wantErr:        tc.wantErr,
	}.run(t)
}

type domainIndexTest struct {
	filters     []service.DomainFilter
	responseObj interface{}
	wantResp    interface{}
	wantErr     error
}

func callDomainIndex(
	t *testing.T, filters interface{}, server *httptest.Server,
) (interface{}, error) {
	arg, ok := filters.([]service.DomainFilter)
	assert.True(t, ok)
	svc := getAllInterface(server).Domain()
	return svc.Index(arg...)
}

func (tc domainIndexTest) run(t *testing.T) {
	indexTestCase{
		callSvc:        callDomainIndex,
		assertEquality: assert.HasSameElements,
		urlPrefix:      domainCommonURL,
		filters:        tc.filters,
		responseObj:    tc.responseObj,
		wantResp:       tc.wantResp,
		wantErr:        tc.wantErr,
	}.run(t)
}

type proxyIndexTest struct {
	filters     []service.ProxyFilter
	responseObj interface{}
	wantResp    interface{}
	wantErr     error
}

func callProxyIndex(
	t *testing.T, filters interface{}, server *httptest.Server,
) (interface{}, error) {
	arg, ok := filters.([]service.ProxyFilter)
	assert.True(t, ok)
	svc := getAllInterface(server).Proxy()
	return svc.Index(arg...)
}

func (tc proxyIndexTest) run(t *testing.T) {
	indexTestCase{
		callSvc:        callProxyIndex,
		assertEquality: assert.HasSameElements,
		urlPrefix:      proxyCommonURL,
		filters:        tc.filters,
		responseObj:    tc.responseObj,
		wantResp:       tc.wantResp,
		wantErr:        tc.wantErr,
	}.run(t)
}

type routeIndexTest struct {
	filters     []service.RouteFilter
	responseObj interface{}
	wantResp    interface{}
	wantErr     error
}

func callRouteIndex(
	t *testing.T, filters interface{}, server *httptest.Server,
) (interface{}, error) {
	arg, ok := filters.([]service.RouteFilter)
	assert.True(t, ok)
	svc := getAllInterface(server).Route()
	return svc.Index(arg...)
}

func (tc routeIndexTest) run(t *testing.T) {
	indexTestCase{
		callSvc:        callRouteIndex,
		assertEquality: assert.HasSameElements,
		urlPrefix:      routeCommonURL,
		filters:        tc.filters,
		responseObj:    tc.responseObj,
		wantResp:       tc.wantResp,
		wantErr:        tc.wantErr,
	}.run(t)
}

type sharedRulesIndexTest struct {
	filters     []service.SharedRulesFilter
	responseObj interface{}
	wantResp    interface{}
	wantErr     error
}

func callSharedRulesIndex(
	t *testing.T, filters interface{}, server *httptest.Server,
) (interface{}, error) {
	arg, ok := filters.([]service.SharedRulesFilter)
	assert.True(t, ok)
	svc := getAllInterface(server).SharedRules()
	return svc.Index(arg...)
}

func (tc sharedRulesIndexTest) run(t *testing.T) {
	indexTestCase{
		callSvc:        callSharedRulesIndex,
		assertEquality: assert.HasSameElements,
		urlPrefix:      sharedRulesCommonURL,
		filters:        tc.filters,
		responseObj:    tc.responseObj,
		wantResp:       tc.wantResp,
		wantErr:        tc.wantErr,
	}.run(t)
}

type userIndexTest struct {
	filters     []service.UserFilter
	responseObj interface{}
	wantResp    interface{}
	wantErr     error
}

func callUserIndex(
	t *testing.T, filters interface{}, server *httptest.Server,
) (interface{}, error) {
	arg, ok := filters.([]service.UserFilter)
	assert.True(t, ok)
	admin := getAdminInterface(server).User()
	return admin.Index(arg...)
}

func (tc userIndexTest) run(t *testing.T) {
	indexTestCase{
		callSvc:        callUserIndex,
		assertEquality: assert.HasSameElements,
		urlPrefix:      userCommonURL,
		filters:        tc.filters,
		responseObj:    tc.responseObj,
		wantResp:       tc.wantResp,
		wantErr:        tc.wantErr,
	}.run(t)
}

type zoneIndexTest struct {
	filters     []service.ZoneFilter
	responseObj interface{}
	wantResp    interface{}
	wantErr     error
}

func callZoneIndex(
	t *testing.T, filters interface{}, server *httptest.Server,
) (interface{}, error) {
	arg, ok := filters.([]service.ZoneFilter)
	assert.True(t, ok)
	svc := getAllInterface(server).Zone()
	return svc.Index(arg...)
}

func (tc zoneIndexTest) run(t *testing.T) {
	indexTestCase{
		callSvc:        callZoneIndex,
		assertEquality: assert.HasSameElements,
		urlPrefix:      zoneCommonURL,
		filters:        tc.filters,
		responseObj:    tc.responseObj,
		wantResp:       tc.wantResp,
		wantErr:        tc.wantErr,
	}.run(t)
}

func TestIndexSimple(t *testing.T) {
	clusterIndexTest{
		responseObj: fixtures.PublicClusterSlice,
		wantResp:    fixtures.PublicClusterSlice,
	}.run(t)

	domainIndexTest{
		responseObj: fixtures.PublicDomainSlice,
		wantResp:    fixtures.PublicDomainSlice,
	}.run(t)

	routeIndexTest{
		responseObj: fixtures.PublicRouteSlice,
		wantResp:    fixtures.PublicRouteSlice,
	}.run(t)

	sharedRulesIndexTest{
		responseObj: fixtures.PublicSharedRulesSlice,
		wantResp:    fixtures.PublicSharedRulesSlice,
	}.run(t)

	proxyIndexTest{
		responseObj: fixtures.PublicProxySlice,
		wantResp:    fixtures.PublicProxySlice,
	}.run(t)

	userIndexTest{
		responseObj: fixtures.PublicUserSlice,
		wantResp:    fixtures.PublicUserSlice,
	}.run(t)

	zoneIndexTest{
		responseObj: fixtures.PublicZoneSlice,
		wantResp:    fixtures.PublicZoneSlice,
	}.run(t)
}

func TestIndexOneSimpleFilter(t *testing.T) {
	clusterIndexTest{
		filters:     []service.ClusterFilter{{Name: "bob"}},
		responseObj: fixtures.PublicClusterSlice,
		wantResp:    fixtures.PublicClusterSlice,
	}.run(t)

	domainIndexTest{
		filters:     []service.DomainFilter{{Name: "bob"}},
		responseObj: fixtures.PublicDomainSlice,
		wantResp:    fixtures.PublicDomainSlice,
	}.run(t)

	routeIndexTest{
		filters:     []service.RouteFilter{{Path: "bob"}},
		responseObj: fixtures.PublicRouteSlice,
		wantResp:    fixtures.PublicRouteSlice,
	}.run(t)

	sharedRulesIndexTest{
		filters:     []service.SharedRulesFilter{{ZoneKey: api.ZoneKey("bob")}},
		responseObj: fixtures.PublicSharedRulesSlice,
		wantResp:    fixtures.PublicSharedRulesSlice,
	}.run(t)

	proxyIndexTest{
		filters:     []service.ProxyFilter{{Name: "bob"}},
		responseObj: fixtures.PublicProxySlice,
		wantResp:    fixtures.PublicProxySlice,
	}.run(t)

	userIndexTest{
		filters:     []service.UserFilter{{LoginEmail: "bob"}},
		responseObj: fixtures.PublicUserSlice,
		wantResp:    fixtures.PublicUserSlice,
	}.run(t)

	zoneIndexTest{
		filters:     []service.ZoneFilter{{Name: "bob"}},
		responseObj: fixtures.PublicZoneSlice,
		wantResp:    fixtures.PublicZoneSlice,
	}.run(t)
}

func TestIndexComplexFilterRequiringEscaping(t *testing.T) {
	clusterIndexTest{
		filters:     []service.ClusterFilter{{Name: "\"'%?^&aoeu=snth"}},
		responseObj: fixtures.PublicClusterSlice,
		wantResp:    fixtures.PublicClusterSlice,
	}.run(t)

	domainIndexTest{
		filters:     []service.DomainFilter{{Name: "\"'%?^&aoeu=snth"}},
		responseObj: fixtures.PublicDomainSlice,
		wantResp:    fixtures.PublicDomainSlice,
	}.run(t)

	routeIndexTest{
		filters:     []service.RouteFilter{{Path: "\"'%?^&aoeu=snth"}},
		responseObj: fixtures.PublicRouteSlice,
		wantResp:    fixtures.PublicRouteSlice,
	}.run(t)

	sharedRulesIndexTest{
		filters:     []service.SharedRulesFilter{{ZoneKey: api.ZoneKey("\"'%?^&aoeu=snth")}},
		responseObj: fixtures.PublicSharedRulesSlice,
		wantResp:    fixtures.PublicSharedRulesSlice,
	}.run(t)

	proxyIndexTest{
		filters:     []service.ProxyFilter{{Name: "\"'%?^&aoeu=snth"}},
		responseObj: fixtures.PublicProxySlice,
		wantResp:    fixtures.PublicProxySlice,
	}.run(t)

	userIndexTest{
		filters:     []service.UserFilter{{LoginEmail: "\"'%?^&aoeu=snth"}},
		responseObj: fixtures.PublicUserSlice,
		wantResp:    fixtures.PublicUserSlice,
	}.run(t)

	zoneIndexTest{
		filters:     []service.ZoneFilter{{Name: "\"'%?^&aoeu=snth"}},
		responseObj: fixtures.PublicZoneSlice,
		wantResp:    fixtures.PublicZoneSlice,
	}.run(t)
}

func TestIndexWrapsWeirdResponses(t *testing.T) {
	wantErr := func(path string) *httperr.Error {
		return httperr.New500(
			fmt.Sprintf(
				"got malformed response for {{URL}}%s; unmarshal error: "+
					"'invalid character 'w' looking for beginning of value' "+
					"- content: 'wtf'",
				path,
			),
			httperr.UnknownDecodingCode,
		)
	}

	clusterIndexTest{
		responseObj: "wtf",
		wantErr:     wantErr("/v1.0/cluster"),
	}.run(t)

	domainIndexTest{
		responseObj: "wtf",
		wantErr:     wantErr("/v1.0/domain"),
	}.run(t)

	proxyIndexTest{
		responseObj: "wtf",
		wantErr:     wantErr("/v1.0/proxy"),
	}.run(t)

	routeIndexTest{
		responseObj: "wtf",
		wantErr:     wantErr("/v1.0/route"),
	}.run(t)

	sharedRulesIndexTest{
		responseObj: "wtf",
		wantErr:     wantErr("/v1.0/shared_rules"),
	}.run(t)

	userIndexTest{
		responseObj: "wtf",
		wantErr:     wantErr("/v1.0/admin/user"),
	}.run(t)

	zoneIndexTest{
		responseObj: "wtf",
		wantErr:     wantErr("/v1.0/zone"),
	}.run(t)
}

func TestIndexHandlesRealErrors(t *testing.T) {
	msg := "aoeuaoeu"
	code := httperr.ErrorCode("snthsnth")
	sentErr := httperr.New400(msg, code)

	clusterIndexTest{
		responseObj: envelope.Response{sentErr, nil},
		wantErr:     sentErr,
	}.run(t)

	domainIndexTest{
		responseObj: envelope.Response{sentErr, nil},
		wantErr:     sentErr,
	}.run(t)

	routeIndexTest{
		responseObj: envelope.Response{sentErr, nil},
		wantErr:     sentErr,
	}.run(t)

	sharedRulesIndexTest{
		responseObj: envelope.Response{sentErr, nil},
		wantErr:     sentErr,
	}.run(t)

	proxyIndexTest{
		responseObj: envelope.Response{sentErr, nil},
		wantErr:     sentErr,
	}.run(t)

	userIndexTest{
		responseObj: envelope.Response{sentErr, nil},
		wantErr:     sentErr,
	}.run(t)

	zoneIndexTest{
		responseObj: envelope.Response{sentErr, nil},
		wantErr:     sentErr,
	}.run(t)
}
