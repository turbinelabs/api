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

type AccessTokenKey string
type AccessTokens []AccessToken

// AccessTokens are personal access tokens owned by users, further specialized
// by org
type AccessToken struct {
	// AccessTokenKey is the lookup key for this token.
	AccessTokenKey AccessTokenKey `json:"access_token_key"`

	// Description is a summary of how this token will be used. It may not be
	// empty and must be less than 255 characters.
	Description string `json:"description"`

	// SignedToken is a field that is set only once when the token is created and
	// may be passed by a request to authorize it. This may be revoked and should
	// be treated as carefully as a password.
	SignedToken string `json:"signed_token"`

	UserKey   UserKey    `json:"user_key"`
	OrgKey    OrgKey     `json:"-"`
	CreatedAt *time.Time `json:"created_at"`
	Checksum
}

func (t AccessToken) Key() string           { return string(t.AccessTokenKey) }
func (t AccessToken) GetUserKey() UserKey   { return t.UserKey }
func (t AccessToken) GetOrgKey() OrgKey     { return t.OrgKey }
func (t AccessToken) GetChecksum() Checksum { return t.Checksum }

func (t AccessToken) IsNil() bool {
	return t.Equals(AccessToken{})
}

const DescriptionLen = 255

var (
	tooLongErr  = fmt.Sprintf("must be less than %v characters", DescriptionLen)
	tooShortErr = "must have a description"
)

func (t AccessToken) IsValid() *ValidationError {
	scope := func(in string) string { return "access_token." + in }

	errs := &ValidationError{}

	errCheckKey(string(t.AccessTokenKey), errs, scope("access_token_key"))
	errCheckKey(string(t.UserKey), errs, scope("user_key"))
	errCheckKey(string(t.OrgKey), errs, scope("org_key"))
	if strings.TrimSpace(t.Description) == "" {
		errs.AddNew(ErrorCase{scope("description"), tooShortErr})
	}
	if len(t.Description) > DescriptionLen {
		errs.AddNew(ErrorCase{scope("description"), tooLongErr})
	}

	return errs.OrNil()
}

func (t AccessToken) Equals(o AccessToken) bool {
	return t.AccessTokenKey == o.AccessTokenKey &&
		t.Description == o.Description &&
		t.UserKey == o.UserKey &&
		t.OrgKey == o.OrgKey &&
		t.SignedToken == o.SignedToken &&
		tbntime.Equal(t.CreatedAt, o.CreatedAt) &&
		t.Checksum == o.Checksum
}
