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

package api

import (
	"fmt"
	"strings"
	"time"

	tbntime "github.com/turbinelabs/nonstdlib/time"
)

// ResponseData is a collection of annotations that should be applied to
// responses when handling a request.
type ResponseData struct {
	// Headers are HTTP headers that will be attached to a response.
	Headers []HeaderDatum `json:"headers,omitempty"`

	// Cookies are attached via 'Set-Cookie' header.
	Cookies []CookieDatum `json:"cookies,omitempty"`
}

// Equals checks if two ResponseData objects are semantically equivalent. They
// are considered equal iff they contain the same set of Headers and Cookies.
// A change in slice order is not considered a difference in ResponseData objects.
//
// A HeaderDatum is identified by its Name attribute via case-insensitive
// comparison.  This means if ResponseData 1 has Header "x-foo" and
// ResponseData 2 has header "X-Foo" that are otherwise equal then are
// considered the same.  CookieDatum is identified by its name attribute
// via case sensitive comparison.
func (rd ResponseData) Equals(o ResponseData) bool {
	if len(rd.Headers) != len(o.Headers) {
		return false
	}
	if len(rd.Cookies) != len(o.Cookies) {
		return false
	}

	hdrs := map[string]HeaderDatum{}
	checkedHdr := map[string]bool{}
	for _, hdr := range rd.Headers {
		hdrs[hdr.Name] = hdr
	}

	cks := map[string]CookieDatum{}
	checkedCk := map[string]bool{}
	for _, ck := range rd.Cookies {
		cks[ck.Name] = ck
	}

	for _, hdr := range o.Headers {
		n := hdr.Name
		old, has := hdrs[n]
		if checkedHdr[n] || !has || !old.Equals(hdr) {
			return false
		}
		checkedHdr[n] = true
	}

	for _, ck := range o.Cookies {
		n := ck.Name
		old, has := cks[n]
		if checkedCk[n] || !has || !old.Equals(ck) {
			return false
		}
		checkedCk[n] = true
	}

	return true
}

func (rd ResponseData) IsValid() *ValidationError {
	errs := &ValidationError{}

	hdrSeen := map[string]int{}
	hdrIsSeen := func(s string) bool {
		canonical := strings.ToLower(s)
		hdrSeen[canonical]++
		return hdrSeen[canonical] > 1
	}

	for _, hdr := range rd.Headers {
		if hdrIsSeen(hdr.Name) {
			errs.AddNew(ErrorCase{"headers", fmt.Sprintf("Header %q exported multiple times", hdr.Name)})
		}
		parent := fmt.Sprintf("headers[%v]", hdr.Name)
		errs.MergePrefixed(hdr.IsValid(), parent)
	}

	ckSeen := map[string]int{}
	ckIsSeen := func(s string) bool {
		ckSeen[s]++
		return ckSeen[s] > 1
	}

	for _, ck := range rd.Cookies {
		if ckIsSeen(ck.Name) {
			errs.AddNew(ErrorCase{"cookies", fmt.Sprintf("Cookie %q exported multiple times", ck.Name)})
		}

		parent := fmt.Sprintf("cookies[%v]", ck.Name)
		errs.MergePrefixed(ck.IsValid(), parent)
	}

	return errs.OrNil()
}

func (rd ResponseData) Len() int {
	return len(rd.Headers) + len(rd.Cookies)
}

// ResponseDatum represents the set of information necessary to determine
// how to name and produce the value that should be attached to a response
// and under what condictions the data should be sent back.
type ResponseDatum struct {
	// Name of the data being sent back to the requesting client.
	Name string `json:"name"`

	// Value is either a literal value or a reference to metadatum on the server
	// that handles a request.
	Value string `json:"value"`

	// ValueIsLiteral, if set, menas that Value will be treated as a literal
	// instead of a reference to be resolved as the key of a metadatum set on
	// the server handling a request.
	ValueIsLiteral bool `json:"value_is_literal,omitempty"`

	// AlwaysSend, if set, will cause a ResponseDatum to always be returned with
	// a request. If false the annotation will only be set on response codes of
	// 2xx or 3xx.
	AlwaysSend bool `json:"always_send"`
}

// HeaderDatum represents a header that should be attached to a response to
// some requset served by the object containing a ResponseData config. Some
// points to note are that HeaderDatum are not case sensitive with respect to
// their Name value which impacts equality checks.
type HeaderDatum struct {
	ResponseDatum
}

// Equals compares two HeaderDatum objects. A HeaderDatum is determined to be
// equal if the name (case insensitive check), value, ValueIsLiteral, and
// AlwaysSend attributes are equal.
func (hd HeaderDatum) Equals(o HeaderDatum) bool {
	return strings.ToLower(hd.Name) == strings.ToLower(o.Name) &&
		hd.Value == o.Value &&
		hd.ValueIsLiteral == o.ValueIsLiteral &&
		hd.AlwaysSend == o.AlwaysSend
}

// IsValid ensures that HeaderDatum attributes have reasonable values:
//
//   - Name must not be empty
//   - Name must be a valid header
//   - Value may not be empty
func (hd HeaderDatum) IsValid() *ValidationError {
	errs := &ValidationError{}

	errCheckPattern(false, hd.Name, errs, HeaderNamePattern, "name", "")
	if strings.TrimSpace(hd.Value) == "" {
		errs.AddNew(ErrorCase{"value", "may not be empty"})
	}

	return errs.OrNil()
}

// SameSiteType allows specification for the 'SameSite' annotation on a cookie
// response. This allows some control over when the cookie is sent to a server
// see https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Set-Cookie for
// details.
type SameSiteType string

const (
	// SameSiteUnset represents the default value and will not impact the cookie
	// annotation set.
	SameSiteUnset SameSiteType = ""

	// SameSiteStrict causes 'SameSite=Strict' to be passed back with a cookie.
	SameSiteStrict SameSiteType = "Strict"

	// SameSiteLax causes 'SameSite=Lax' to be passed back with a cookie.
	SameSiteLax SameSiteType = "Lax"
)

// CookieDatum represents a cookie that should be attached to the response to
// some requset served by the object containing a ResponseData config. A
// CookieDatum's Name is case sensitive.
type CookieDatum struct {
	ResponseDatum

	// ExpiresInSec indicates how long a cookie will be valid (in seconds) or
	// indicates that a cookie should be expired if set to 0. Specifically for
	// values > 0 this becomes a 'Max-Age' cookie annotation and for 0 'Expires'
	// is set to the unix epoch, UTC.
	ExpiresInSec *uint `json:"expires_in_sec"`

	// Domain specifies hosts to which the cookie will be sent.
	Domain string `json:"domain"`

	// Path indicates a URL path that must be met in a requset for the cookie to
	// be sent to the server.
	Path string `json:"path"`

	// Secure will only be sent to a server when a request is made via HTTPS.
	Secure bool `json:"secure"`

	// HttpOnly cookies are not available via javascript throught Document.cookie.
	HttpOnly bool `json:"http_only"`

	// SameSiteType specifies how a cookie should be treated when a request is being
	// made across site boundaries (e.g. from another domain, used to help protect
	// against CSRF).
	SameSite SameSiteType `json:"same_site"`
}

// Equals returns true if all CookieDatum attributes are the same.
func (cd CookieDatum) Equals(o CookieDatum) bool {
	expTimeEq := (cd.ExpiresInSec == nil && o.ExpiresInSec == nil) ||
		(cd.ExpiresInSec != nil && o.ExpiresInSec != nil && *cd.ExpiresInSec == *o.ExpiresInSec)

	return cd.Name == o.Name &&
		cd.Value == o.Value &&
		cd.ValueIsLiteral == o.ValueIsLiteral &&
		cd.AlwaysSend == o.AlwaysSend &&
		expTimeEq &&
		cd.Domain == o.Domain &&
		cd.Path == o.Path &&
		cd.Secure == o.Secure &&
		cd.HttpOnly == o.HttpOnly &&
		cd.SameSite == o.SameSite
}

// IsValid ensures reasonable values for CookieDatum attributes:
//
//   - Name may not be empty and must be a valid cookie name
//   - SameSite must be one of the defined SameSite values
func (cd CookieDatum) IsValid() *ValidationError {
	errs := &ValidationError{}

	errCheckPattern(false, cd.Name, errs, CookieNamePattern, "name", "")
	ss := cd.SameSite
	if ss != SameSiteUnset && ss != SameSiteStrict && ss != SameSiteLax {
		errs.AddNew(ErrorCase{"same_site", fmt.Sprintf("%q is not a valid value", ss)})
	}

	return errs.OrNil()
}

// Annotation returns a string that is attached to the cookie returned to
// specify how it should be treated by the browser based on its configuration.
func (c CookieDatum) Annotation() string {
	auxStr := []string{}
	add := func(s ...string) {
		auxStr = append(auxStr, strings.Join(s, ""))
	}

	if c.ExpiresInSec != nil {
		if *c.ExpiresInSec == 0 {
			add("Expires=", time.Time{}.Format(tbntime.CookieFormat))
		} else {
			add("Max-Age=", fmt.Sprintf("%v", *c.ExpiresInSec))
		}
	}

	if c.Domain != "" {
		add("Domain=", c.Domain)
	}

	if c.Path != "" {
		add("Path=", c.Path)
	}

	if c.Secure {
		add("Secure")
	}

	if c.HttpOnly {
		add("HttpOnly")
	}

	if c.SameSite != SameSiteUnset {
		add("SameSite=", string(c.SameSite))
	}

	return strings.Join(auxStr, "; ")
}