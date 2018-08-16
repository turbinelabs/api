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
	"time"

	tbntime "github.com/turbinelabs/nonstdlib/time"
)

type UserKey string
type Users []User
type APIAuthKey string

// A User is an actor of an Org. An API key is set if they're allowed to make
// API calls.
type User struct {
	UserKey    UserKey    `json:"user_key"`
	LoginEmail string     `json:"login_email"`
	Properties Metadata   `json:"properties"`
	APIAuthKey APIAuthKey `json:"api_auth_key,omitempty"`
	OrgKey     OrgKey     `json:"org_key"`
	DeletedAt  *time.Time `json:"deleted_at,omitempty"`
	Checksum
}

func (o User) GetOrgKey() OrgKey     { return o.OrgKey }
func (o User) Key() string           { return string(o.UserKey) }
func (o User) GetChecksum() Checksum { return o.Checksum }

func (u User) IsNil() bool {
	return u.Equals(User{})
}

func (u User) IsValid() *ValidationError {
	scope := func(in string) string { return "user." + in }

	errs := &ValidationError{}

	errCheckKey(string(u.UserKey), errs, scope("user_key"))

	if u.LoginEmail == "" {
		errs.AddNew(ErrorCase{scope("login_email"), "may not be empty"})
	}

	if u.APIAuthKey != "" {
		// can't check for key because alternative auth systems may have a different
		// approach to auth key generation than a simple UUID
		errCheckIndex(string(u.APIAuthKey), errs, scope("api_auth_key"))
	}
	errCheckKey(string(u.OrgKey), errs, scope("org_key"))

	errs.MergePrefixed(UserPropertiesValid(u.Properties), "user")

	return errs.OrNil()
}

func UserPropertiesValid(props Metadata) *ValidationError {
	return MetadataValid(
		"properties",
		props,
		MetadataCheckKeysMatchPattern(AllowedIndexPattern, AllowedIndexPatternMatchFailure),
	)
}

func (u User) Equals(o User) bool {
	return u.UserKey == o.UserKey &&
		u.LoginEmail == o.LoginEmail &&
		u.APIAuthKey == o.APIAuthKey &&
		u.OrgKey == o.OrgKey &&
		tbntime.Equal(u.DeletedAt, o.DeletedAt) &&
		u.Checksum == o.Checksum &&
		u.Properties.Equals(o.Properties)
}
