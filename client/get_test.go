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
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/turbinelabs/api"
	apihttp "github.com/turbinelabs/api/http"
	"github.com/turbinelabs/api/http/envelope"
	httperr "github.com/turbinelabs/api/http/error"
	"github.com/turbinelabs/test/assert"
)

type MkSvcDoGetFn func(string, *httptest.Server) (interface{}, error)
type EqualityCheckFn func(interface{}, interface{}) bool
type AssertURLFn func(*testing.T, *url.URL, string)

type getTestCase struct {
	svcCall       MkSvcDoGetFn    // makes the appropriate get call
	checkEquality EqualityCheckFn // compare the object returned and the object we want
	key           string          // which cluster are we getting?
	responseObj   interface{}     // will be rendered as json and set as the response body
	wantResp      interface{}     // what the API call should produce
	wantErr       error           // what error, if any, the API call should produce
	assertURL     AssertURLFn     // verify the url is correct
}

func (tc getTestCase) run(t *testing.T) {
	verifier := verifyingHandler{
		fn: func(rr apihttp.RichRequest) {
			tc.assertURL(t, rr.Underlying().URL, tc.key)
		},
		status:   http.StatusOK,
		response: tc.responseObj,
	}

	s := httptest.NewServer(verifier)
	defer s.Close()

	got, gotErr := tc.svcCall(tc.key, s)

	assert.DeepEqual(t, gotErr, tc.wantErr)
	assert.True(t, tc.checkEquality(got, tc.wantResp))
}

func mkAssertGetURL(prefix string) func(*testing.T, *url.URL, string) {
	return func(t *testing.T, inurl *url.URL, key string) {
		if !assertURLPrefix(t, inurl.Path, prefix) {
			return
		}

		// the path should have the desired key after the object base
		remainder := stripURLPrefix(inurl.Path, prefix)
		wantKeyStr := "/" + key
		assert.Equal(t, remainder, wantKeyStr)

		// no parameters expected on get requests
		assert.Equal(t, inurl.RawQuery, "")
	}
}

type clusterGetTest struct {
	clusterKey  api.ClusterKey
	responseObj interface{}
	wantResp    api.Cluster
	wantErr     error
}

type domainGetTest struct {
	domainKey   api.DomainKey
	responseObj interface{}
	wantResp    api.Domain
	wantErr     error
}

type proxyGetTest struct {
	proxyKey    api.ProxyKey
	responseObj interface{}
	wantResp    api.Proxy
	wantErr     error
}

type routeGetTest struct {
	routeKey    api.RouteKey
	responseObj interface{}
	wantResp    api.Route
	wantErr     error
}

type sharedRulesGetTest struct {
	sharedRulesKey api.SharedRulesKey
	responseObj    interface{}
	wantResp       api.SharedRules
	wantErr        error
}

type userGetTest struct {
	userKey     api.UserKey
	responseObj interface{}
	wantResp    api.User
	wantErr     error
}

type zoneGetTest struct {
	zoneKey     api.ZoneKey
	responseObj interface{}
	wantResp    api.Zone
	wantErr     error
}

func (tc clusterGetTest) run(t *testing.T) {
	callEquals := func(i1, i2 interface{}) bool {
		c1, ok := i1.(api.Cluster)
		assert.True(t, ok)

		c2, ok := i2.(api.Cluster)
		assert.True(t, ok)

		return c1.Equals(c2)
	}

	svcCall := func(key string, server *httptest.Server) (interface{}, error) {
		svc := getAllInterface(server).Cluster()
		return svc.Get(api.ClusterKey(key))
	}

	getTestCase{
		svcCall:       svcCall,
		checkEquality: callEquals,
		key:           string(tc.clusterKey),
		responseObj:   tc.responseObj,
		wantResp:      tc.wantResp,
		wantErr:       tc.wantErr,
		assertURL:     mkAssertGetURL(clusterCommonURL),
	}.run(t)
}

func (tc domainGetTest) run(t *testing.T) {
	callEquals := func(i1, i2 interface{}) bool {
		c1, ok := i1.(api.Domain)
		assert.True(t, ok)

		c2, ok := i2.(api.Domain)
		assert.True(t, ok)

		return c1.Equals(c2)
	}

	svcCall := func(key string, server *httptest.Server) (interface{}, error) {
		svc := getAllInterface(server).Domain()
		return svc.Get(api.DomainKey(key))
	}

	getTestCase{
		svcCall:       svcCall,
		checkEquality: callEquals,
		key:           string(tc.domainKey),
		responseObj:   tc.responseObj,
		wantResp:      tc.wantResp,
		wantErr:       tc.wantErr,
		assertURL:     mkAssertGetURL(domainCommonURL),
	}.run(t)
}

func (tc proxyGetTest) run(t *testing.T) {
	callEquals := func(i1, i2 interface{}) bool {
		c1, ok := i1.(api.Proxy)
		assert.True(t, ok)

		c2, ok := i2.(api.Proxy)
		assert.True(t, ok)

		return c1.Equals(c2)
	}

	svcCall := func(key string, server *httptest.Server) (interface{}, error) {
		svc := getAllInterface(server).Proxy()
		return svc.Get(api.ProxyKey(key))
	}

	getTestCase{
		svcCall:       svcCall,
		checkEquality: callEquals,
		key:           string(tc.proxyKey),
		responseObj:   tc.responseObj,
		wantResp:      tc.wantResp,
		wantErr:       tc.wantErr,
		assertURL:     mkAssertGetURL(proxyCommonURL),
	}.run(t)
}

func (tc userGetTest) run(t *testing.T) {
	callEquals := func(i1, i2 interface{}) bool {
		c1, ok := i1.(api.User)
		assert.True(t, ok)

		c2, ok := i2.(api.User)
		assert.True(t, ok)

		return c1.Equals(c2)
	}

	svcCall := func(key string, server *httptest.Server) (interface{}, error) {
		svc := getAdminInterface(server).User()
		return svc.Get(api.UserKey(key))
	}

	getTestCase{
		svcCall:       svcCall,
		checkEquality: callEquals,
		key:           string(tc.userKey),
		responseObj:   tc.responseObj,
		wantResp:      tc.wantResp,
		wantErr:       tc.wantErr,
		assertURL:     mkAssertGetURL(userCommonURL),
	}.run(t)
}

func (tc zoneGetTest) run(t *testing.T) {
	callEquals := func(i1, i2 interface{}) bool {
		c1, ok := i1.(api.Zone)
		assert.True(t, ok)

		c2, ok := i2.(api.Zone)
		assert.True(t, ok)

		return c1.Equals(c2)
	}

	svcCall := func(key string, server *httptest.Server) (interface{}, error) {
		svc := getAllInterface(server).Zone()
		return svc.Get(api.ZoneKey(key))
	}

	getTestCase{
		svcCall:       svcCall,
		checkEquality: callEquals,
		key:           string(tc.zoneKey),
		responseObj:   tc.responseObj,
		wantResp:      tc.wantResp,
		wantErr:       tc.wantErr,
		assertURL:     mkAssertGetURL(zoneCommonURL),
	}.run(t)
}

func (tc routeGetTest) run(t *testing.T) {
	callEquals := func(i1, i2 interface{}) bool {
		c1, ok := i1.(api.Route)
		assert.True(t, ok)

		c2, ok := i2.(api.Route)
		assert.True(t, ok)

		return c1.Equals(c2)
	}

	svcCall := func(key string, server *httptest.Server) (interface{}, error) {
		svc := getAllInterface(server).Route()
		return svc.Get(api.RouteKey(key))
	}

	getTestCase{
		svcCall:       svcCall,
		checkEquality: callEquals,
		key:           string(tc.routeKey),
		responseObj:   tc.responseObj,
		wantResp:      tc.wantResp,
		wantErr:       tc.wantErr,
		assertURL:     mkAssertGetURL(routeCommonURL),
	}.run(t)
}

func (tc sharedRulesGetTest) run(t *testing.T) {
	callEquals := func(i1, i2 interface{}) bool {
		c1, ok := i1.(api.SharedRules)
		assert.True(t, ok)

		c2, ok := i2.(api.SharedRules)
		assert.True(t, ok)

		return c1.Equals(c2)
	}

	svcCall := func(key string, server *httptest.Server) (interface{}, error) {
		svc := getAllInterface(server).SharedRules()
		return svc.Get(api.SharedRulesKey(key))
	}

	getTestCase{
		svcCall:       svcCall,
		checkEquality: callEquals,
		key:           string(tc.sharedRulesKey),
		responseObj:   tc.responseObj,
		wantResp:      tc.wantResp,
		wantErr:       tc.wantErr,
		assertURL:     mkAssertGetURL(sharedRulesCommonURL),
	}.run(t)
}

func TestGetNoKey(t *testing.T) {
	e := func(t string) error {
		return httperr.New400(
			fmt.Sprintf("%sKey is a required parameter", t),
			"ObjectKeyRequiredErrorCode",
		)
	}

	clusterGetTest{wantErr: e("Cluster")}.run(t)
	domainGetTest{wantErr: e("Domain")}.run(t)
	routeGetTest{wantErr: e("Route")}.run(t)
	sharedRulesGetTest{wantErr: e("SharedRules")}.run(t)
	proxyGetTest{wantErr: e("Proxy")}.run(t)
	userGetTest{wantErr: e("User")}.run(t)
	zoneGetTest{wantErr: e("Zone")}.run(t)
}

func TestGetSimple(t *testing.T) {
	rc := fixtures.Cluster1
	rc.OrgKey = ""
	clusterGetTest{
		clusterKey:  fixtures.ClusterKey1,
		responseObj: fixtures.Cluster1,
		wantResp:    rc,
	}.run(t)

	rd := fixtures.Domain1
	rd.OrgKey = ""
	domainGetTest{
		domainKey:   fixtures.DomainKey1,
		responseObj: fixtures.Domain1,
		wantResp:    rd,
	}.run(t)

	rr := fixtures.Route1
	rr.OrgKey = ""
	routeGetTest{
		routeKey:    fixtures.RouteKey1,
		responseObj: fixtures.Route1,
		wantResp:    rr,
	}.run(t)

	rs := fixtures.SharedRules1
	rs.OrgKey = ""
	sharedRulesGetTest{
		sharedRulesKey: fixtures.SharedRulesKey1,
		responseObj:    fixtures.SharedRules1,
		wantResp:       rs,
	}.run(t)

	rp := fixtures.Proxy1
	rp.OrgKey = ""
	proxyGetTest{
		proxyKey:    fixtures.ProxyKey1,
		responseObj: fixtures.Proxy1,
		wantResp:    rp,
	}.run(t)

	rz := fixtures.Zone1
	rz.OrgKey = ""
	zoneGetTest{
		zoneKey:     fixtures.ZoneKey1,
		responseObj: fixtures.Zone1,
		wantResp:    rz,
	}.run(t)
}

func TestGetUserReturnsOrg(t *testing.T) {
	ru := fixtures.User1
	userGetTest{
		userKey:     fixtures.UserKey1,
		responseObj: fixtures.User1,
		wantResp:    ru,
	}.run(t)
}

func TestGetWrapsWeirdResponses(t *testing.T) {
	wantErr := httperr.New500(
		"got malformed response; unmarshal error: 'invalid character 'w' looking "+
			"for beginning of value' - content: 'wtf'",
		"UnknownDecodingCode")

	clusterGetTest{
		clusterKey:  fixtures.ClusterKey1,
		responseObj: "wtf",
		wantResp:    api.Cluster{},
		wantErr:     wantErr,
	}.run(t)

	domainGetTest{
		domainKey:   fixtures.DomainKey1,
		responseObj: "wtf",
		wantResp:    api.Domain{},
		wantErr:     wantErr,
	}.run(t)

	routeGetTest{
		routeKey:    fixtures.RouteKey1,
		responseObj: "wtf",
		wantResp:    api.Route{},
		wantErr:     wantErr,
	}.run(t)

	sharedRulesGetTest{
		sharedRulesKey: fixtures.SharedRulesKey1,
		responseObj:    "wtf",
		wantResp:       api.SharedRules{},
		wantErr:        wantErr,
	}.run(t)

	proxyGetTest{
		proxyKey:    fixtures.ProxyKey1,
		responseObj: "wtf",
		wantResp:    api.Proxy{},
		wantErr:     wantErr,
	}.run(t)

	userGetTest{
		userKey:     fixtures.UserKey1,
		responseObj: "wtf",
		wantResp:    api.User{},
		wantErr:     wantErr,
	}.run(t)

	zoneGetTest{
		zoneKey:     fixtures.ZoneKey1,
		responseObj: "wtf",
		wantResp:    api.Zone{},
		wantErr:     wantErr,
	}.run(t)
}

func TestGetHandlesRealErrors(t *testing.T) {
	msg := "aoeuaoeu"
	code := httperr.ErrorCode("snthsnth")
	sentErr := httperr.New400(msg, code)

	clusterGetTest{
		clusterKey:  fixtures.ClusterKey1,
		responseObj: envelope.Response{sentErr, nil},
		wantResp:    api.Cluster{},
		wantErr:     sentErr,
	}.run(t)

	domainGetTest{
		domainKey:   fixtures.DomainKey1,
		responseObj: envelope.Response{sentErr, nil},
		wantResp:    api.Domain{},
		wantErr:     sentErr,
	}.run(t)

	routeGetTest{
		routeKey:    fixtures.RouteKey1,
		responseObj: envelope.Response{sentErr, nil},
		wantResp:    api.Route{},
		wantErr:     sentErr,
	}.run(t)

	sharedRulesGetTest{
		sharedRulesKey: fixtures.SharedRulesKey1,
		responseObj:    envelope.Response{sentErr, nil},
		wantResp:       api.SharedRules{},
		wantErr:        sentErr,
	}.run(t)

	proxyGetTest{
		proxyKey:    fixtures.ProxyKey1,
		responseObj: envelope.Response{sentErr, nil},
		wantResp:    api.Proxy{},
		wantErr:     sentErr,
	}.run(t)

	userGetTest{
		userKey:     fixtures.UserKey1,
		responseObj: envelope.Response{sentErr, nil},
		wantResp:    api.User{},
		wantErr:     sentErr,
	}.run(t)

	zoneGetTest{
		zoneKey:     fixtures.ZoneKey1,
		responseObj: envelope.Response{sentErr, nil},
		wantResp:    api.Zone{},
		wantErr:     sentErr,
	}.run(t)
}
