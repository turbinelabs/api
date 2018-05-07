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

package api

import (
	"fmt"
	"strings"
	"testing"

	"github.com/turbinelabs/nonstdlib/ptr"
	"github.com/turbinelabs/test/assert"
)

func mkHD() HeaderDatum {
	return HeaderDatum{
		ResponseDatum{
			"X-Header-Name",
			"Header-Value",
			false,
			true,
		},
	}
}

func mkCD() CookieDatum {
	return CookieDatum{
		ResponseDatum{
			"Cookie-name",
			"cookie-value",
			true,
			false,
		},
		nil,
		"domain.com",
		"/path/foo/bar",
		true,
		false,
		SameSiteLax,
	}
}

func mkRD() ResponseData {
	h1 := mkHD()
	h2 := mkHD()
	h2.Name += "-header2"
	c1 := mkCD()
	c2 := mkCD()
	c2.Name += "-cookie2"

	return ResponseData{
		[]HeaderDatum{h1, h2},
		[]CookieDatum{c1, c2},
	}
}

type headerDatumTest struct {
	t     *testing.T
	vs    HeaderDatum
	equal bool
}

func (hdt headerDatumTest) run() {
	assert.Equal(hdt.t, mkHD().Equals(hdt.vs), hdt.equal)
}

func TestHeaderDatumEquals(t *testing.T) {
	headerDatumTest{t, mkHD(), true}.run()
}

func TestHeaderDatumNameChanged(t *testing.T) {
	hd := mkHD()
	hd.Name += "snth"
	headerDatumTest{t, hd, false}.run()
}

func TestHeaderDatumNameChangedToUpperCase(t *testing.T) {
	hd := mkHD()
	hd.Name = strings.ToUpper(hd.Name)
	headerDatumTest{t, hd, true}.run()
}

func TestHeaderDatumValueChanged(t *testing.T) {
	hd := mkHD()
	hd.Value = strings.ToLower(hd.Value)
	headerDatumTest{t, hd, false}.run()
}

func TestHeaderDatumLiteralChanged(t *testing.T) {
	hd := mkHD()
	hd.ValueIsLiteral = !hd.ValueIsLiteral
	headerDatumTest{t, hd, false}.run()
}

func TestHeaderDatumAlwaysChanged(t *testing.T) {
	hd := mkHD()
	hd.AlwaysSend = !hd.AlwaysSend
	headerDatumTest{t, hd, false}.run()
}

func TestHeaderDatumCanonicalName(t *testing.T) {
	hd := mkHD()
	assert.Equal(t, hd.CanonicalName(), strings.ToLower(hd.Name))
}

type cookieDatumTest struct {
	t     *testing.T
	init  CookieDatum
	vs    CookieDatum
	equal bool
}

func (cd cookieDatumTest) run() {
	assert.Equal(cd.t, cd.init.Equals(cd.vs), cd.equal)
}

func testCk(t *testing.T, c CookieDatum, equalAry ...bool) {
	equal := false
	if 1 <= len(equalAry) {
		equal = equalAry[0]
	}

	cookieDatumTest{t, mkCD(), c, equal}.run()
}

func TestCookieDatumEqual(t *testing.T) {
	testCk(t, mkCD(), true)
}

func TestCookieDatumChangeName(t *testing.T) {
	ck := mkCD()
	ck.Name += "-aoeu"
	testCk(t, ck)
}

func TestCookieDatumChangeValue(t *testing.T) {
	ck := mkCD()
	ck.Value += "-1234"
	testCk(t, ck)
}

func TestCookieDatumChangeLiteral(t *testing.T) {
	ck := mkCD()
	ck.ValueIsLiteral = !ck.ValueIsLiteral
	testCk(t, ck)
}

func TestCookieDatumChangeAlways(t *testing.T) {
	ck := mkCD()
	ck.AlwaysSend = !ck.AlwaysSend
	testCk(t, ck)
}

func TestCookieDatumNameToLower(t *testing.T) {
	ck := mkCD()
	ck.Name = strings.ToLower(ck.Name)
	testCk(t, ck)
}

func TestCookieDatumDomain(t *testing.T) {
	ck := mkCD()
	ck.Domain += ".io"
	testCk(t, ck)
}

func TestCookieDatumPath(t *testing.T) {
	ck := mkCD()
	ck.Path = "/snth"
	testCk(t, ck)
}

func TestCookieDatumSecure(t *testing.T) {
	ck := mkCD()
	ck.Secure = !ck.Secure
	testCk(t, ck)
}

func TestCookieDatumHttpOnly(t *testing.T) {
	ck := mkCD()
	ck.HttpOnly = !ck.HttpOnly
	testCk(t, ck)
}

func TestCookieDatumSameSite(t *testing.T) {
	ck := mkCD()
	ck.SameSite = SameSiteStrict
	testCk(t, ck)
}

func TestCookieDatumExpiresInSecNilSome(t *testing.T) {
	cka := mkCD()
	cka.ExpiresInSec = nil

	ckb := mkCD()
	ckb.ExpiresInSec = ptr.Uint(60)

	cookieDatumTest{t, cka, ckb, false}.run()
}

func TestCookieDatumExpiresInSecNilNil(t *testing.T) {
	cka := mkCD()
	cka.ExpiresInSec = nil

	ckb := mkCD()
	ckb.ExpiresInSec = nil

	cookieDatumTest{t, cka, ckb, true}.run()
}

func TestCookieDatumExpiresInSecSomeNil(t *testing.T) {
	cka := mkCD()
	cka.ExpiresInSec = ptr.Uint(2023)

	ckb := mkCD()
	ckb.ExpiresInSec = nil

	cookieDatumTest{t, cka, ckb, false}.run()
}

func TestCookieDatumExpiresInSecSomeSome(t *testing.T) {
	cka := mkCD()
	cka.ExpiresInSec = ptr.Uint(30)

	ckb := mkCD()
	ckb.ExpiresInSec = ptr.Uint(90)

	cookieDatumTest{t, cka, ckb, false}.run()
}

type responseDataCase struct {
	t     *testing.T
	rd1   ResponseData
	rd2   ResponseData
	equal bool
}

func (rdc responseDataCase) run() {
	assert.Equal(rdc.t, rdc.rd1.Equals(rdc.rd2), rdc.equal)
}

func testRD(t *testing.T, rd ResponseData, equalAry ...bool) {
	equal := false
	if len(equalAry) >= 1 {
		equal = equalAry[0]
	}
	responseDataCase{t, mkRD(), rd, equal}.run()
}

func TestResponseDataEquals(t *testing.T) {
	testRD(t, mkRD(), true)
}

func TestResponseDataEqualsOrderHeaders(t *testing.T) {
	rd := mkRD()
	rd.Headers[0], rd.Headers[1] = rd.Headers[1], rd.Headers[0]
	testRD(t, rd, true)
}

func TestResponseDataEqualsOrderCookies(t *testing.T) {
	rd := mkRD()
	rd.Cookies[0], rd.Cookies[1] = rd.Cookies[1], rd.Cookies[0]
	testRD(t, rd, true)
}

func TestResponseDataChangeHeader(t *testing.T) {
	rd := mkRD()
	rd.Headers[1].Value += "new-value"
	testRD(t, rd)
}

func TestResponseDataChangeCookie(t *testing.T) {
	rd := mkRD()
	rd.Cookies[1].Value += "new-value"
	testRD(t, rd)
}

func TestHeaderDatumIsValid(t *testing.T) {
	assert.Nil(t, mkHD().IsValid())
}

func TestHeaderDatumIsValidNoName(t *testing.T) {
	hd := mkHD()
	hd.Name = ""
	assert.DeepEqual(t, hd.IsValid(), &ValidationError{[]ErrorCase{
		{"name", "may not be empty"},
	}})
}

func TestHeaderDatumIsValidNoValue(t *testing.T) {
	hd := mkHD()
	hd.Value = ""
	assert.DeepEqual(t, hd.IsValid(), &ValidationError{[]ErrorCase{
		{"value", "may not be empty"},
	}})
}

func TestHeaderDatumIsValidBadName(t *testing.T) {
	hd := mkHD()
	hd.Name = "x-header_foo"
	assert.DeepEqual(t, hd.IsValid(), &ValidationError{[]ErrorCase{
		{"name", fmt.Sprintf("must match %v", HeaderNamePatternStr)},
	}})
}

func TestHeaderDatumIsValidWeirdName(t *testing.T) {
	hd := mkHD()
	hd.Value = "x-header_foo.:!@#$ bar"
	assert.Nil(t, hd.IsValid())
}

func cdOK(t *testing.T, cd CookieDatum) {
	assert.Nil(t, cd.IsValid())
}

func TestCookieDatumIsValid(t *testing.T) {
	cdOK(t, mkCD())
}

func TestCookieDatumIsValidNoName(t *testing.T) {
	cd := mkCD()
	cd.Name = ""
	assert.DeepEqual(t, cd.IsValid(), &ValidationError{[]ErrorCase{
		{"name", "may not be empty"},
	}})
}

func TestCookieDatumIsValidBadName(t *testing.T) {
	cd := mkCD()
	cd.Name = "cookie name foo"
	assert.DeepEqual(t, cd.IsValid(), &ValidationError{[]ErrorCase{
		{"name", fmt.Sprintf("must match %v", CookieNamePatternStr)},
	}})
}

func TestCookieDatumIsValidWeirdName(t *testing.T) {
	cd := mkCD()
	cd.Name = "cookie.name"
	cdOK(t, cd)
}

func TestCookieDatumIsValidNoValue(t *testing.T) {
	cd := mkCD()
	cd.Value = ""
	cdOK(t, cd)
}

func TestCookieDatumIsValidNoDomain(t *testing.T) {
	cd := mkCD()
	cd.Domain = ""
	cdOK(t, cd)
}

func TestCookieDatumIsValidNoPath(t *testing.T) {
	cd := mkCD()
	cd.Path = ""
	cdOK(t, cd)
}

func TestCookieDatumIsValidBadSameSite(t *testing.T) {
	cd := mkCD()
	cd.SameSite = "snth"
	assert.DeepEqual(t, cd.IsValid(), &ValidationError{[]ErrorCase{
		{"same_site", fmt.Sprintf("%q is not a valid value", cd.SameSite)},
	}})
}

func TestCookieDatumIsValidGoodSameSite(t *testing.T) {
	cd := mkCD()
	for _, ss := range []SameSiteType{SameSiteUnset, SameSiteLax, SameSiteStrict} {
		cd.SameSite = ss
		cdOK(t, cd)
	}
}

func TestResponseDataIsValid(t *testing.T) {
	assert.Nil(t, mkRD().IsValid())
}

func TestResponseDataIsValidBadHeader(t *testing.T) {
	rd := mkRD()
	rd.Headers[1].Name = ""
	assert.DeepEqual(t, rd.IsValid(), &ValidationError{[]ErrorCase{
		{"headers[].name", "may not be empty"},
	}})
}

func TestResponseDataIsValidDuplicateHeaders(t *testing.T) {
	rd := mkRD()
	rd.Headers = append(rd.Headers, rd.Headers[1])
	n := rd.Headers[1].Name
	assert.DeepEqual(t, rd.IsValid(), &ValidationError{[]ErrorCase{
		{"headers", `Header "` + n + `" exported multiple times`},
	}})
}

func TestResponseDataIsValidDuplicateHeadersCaseDiffers(t *testing.T) {
	rd := mkRD()
	rd.Headers = append(rd.Headers, rd.Headers[1])
	rd.Headers[2].Name = strings.ToUpper(rd.Headers[2].Name)
	n := rd.Headers[2].Name

	assert.DeepEqual(t, rd.IsValid(), &ValidationError{[]ErrorCase{
		{"headers", `Header "` + n + `" exported multiple times`},
	}})
}

func TestResponseDataIsValidBadCookie(t *testing.T) {
	rd := mkRD()
	rd.Cookies[0].SameSite = "foo"
	n := rd.Cookies[0].Name
	assert.DeepEqual(t, rd.IsValid(), &ValidationError{[]ErrorCase{
		{"cookies[" + n + "].same_site", `"foo" is not a valid value`},
	}})
}

func TestResponseDataIsValidDuplicateCokies(t *testing.T) {
	rd := mkRD()
	rd.Cookies = append(rd.Cookies, rd.Cookies[1])
	n := rd.Cookies[1].Name
	assert.DeepEqual(t, rd.IsValid(), &ValidationError{[]ErrorCase{
		{"cookies", `Cookie "` + n + `" exported multiple times`},
	}})
}

func TestResponseDataIsValidDuplicateCokiesDifferntCase(t *testing.T) {
	rd := mkRD()
	rd.Cookies = append(rd.Cookies, rd.Cookies[1])
	rd.Cookies[2].Name = strings.ToUpper(rd.Cookies[2].Name)
	assert.Nil(t, rd.IsValid())
}

func TestResponseDataIsValidEmpty(t *testing.T) {
	rd := ResponseData{}
	assert.Nil(t, rd.IsValid())
}

func TestResponseDataMergeFromTrivial(t *testing.T) {
	rd := mkRD()
	rd2 := mkRD()

	assert.DeepEqual(t, rd.MergeFrom(rd2), mkRD())
}

func TestResponseDataMergeFromHeaders(t *testing.T) {
	recv := ResponseData{
		Headers: []HeaderDatum{
			{ResponseDatum{Name: "x-a", Value: "v"}},
			{ResponseDatum{Name: "x-b", Value: "v"}},
			{ResponseDatum{Name: "x-c", Value: "v"}},
		},
	}

	over := ResponseData{
		Headers: []HeaderDatum{
			{ResponseDatum{Name: "X-0", Value: "v-override"}},
			{ResponseDatum{Name: "X-B", Value: "v-override"}},
			{ResponseDatum{Name: "X-D", Value: "v-override"}},
		},
	}

	expected := ResponseData{
		Headers: []HeaderDatum{
			{ResponseDatum{Name: "x-a", Value: "v"}},
			{ResponseDatum{Name: "X-B", Value: "v-override"}},
			{ResponseDatum{Name: "x-c", Value: "v"}},
			{ResponseDatum{Name: "X-0", Value: "v-override"}},
			{ResponseDatum{Name: "X-D", Value: "v-override"}},
		},
	}

	assert.DeepEqual(t, recv.MergeFrom(over), expected)
}

func TestResponseDataMergeFromCookies(t *testing.T) {
	recv := ResponseData{
		Cookies: []CookieDatum{
			{ResponseDatum: ResponseDatum{Name: "x-a", Value: "v"}},
			{ResponseDatum: ResponseDatum{Name: "x-b", Value: "v"}},
			{ResponseDatum: ResponseDatum{Name: "x-c", Value: "v"}},
		},
	}

	over := ResponseData{
		Cookies: []CookieDatum{
			{ResponseDatum: ResponseDatum{Name: "x-0", Value: "v-override"}},
			{ResponseDatum: ResponseDatum{Name: "x-b", Value: "v-override"}},
			{ResponseDatum: ResponseDatum{Name: "X-C", Value: "v-override"}},
			{ResponseDatum: ResponseDatum{Name: "x-d", Value: "v-override"}},
		},
	}

	expected := ResponseData{
		Cookies: []CookieDatum{
			{ResponseDatum: ResponseDatum{Name: "x-a", Value: "v"}},
			{ResponseDatum: ResponseDatum{Name: "x-b", Value: "v-override"}},
			{ResponseDatum: ResponseDatum{Name: "x-c", Value: "v"}},
			{ResponseDatum: ResponseDatum{Name: "x-0", Value: "v-override"}},
			{ResponseDatum: ResponseDatum{Name: "X-C", Value: "v-override"}},
			{ResponseDatum: ResponseDatum{Name: "x-d", Value: "v-override"}},
		},
	}

	assert.DeepEqual(t, recv.MergeFrom(over), expected)
}
