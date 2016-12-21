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
	APIAuthKey APIAuthKey `json:"api_auth_key"`
	OrgKey     OrgKey     `json:"org_key"`
	DeletedAt  *time.Time `json:"deleted_at"`
	Checksum
}

func (u User) IsNil() bool {
	return u.Equals(User{})
}

func (u User) IsValid(precreation bool) *ValidationError {
	ecase := func(f, m string) ErrorCase {
		return ErrorCase{fmt.Sprintf("user.%s", f), m}
	}

	errs := &ValidationError{}

	keyValid := precreation || u.UserKey != ""
	if !keyValid {
		errs.AddNew(ecase("user_key", "must not be empty"))
	}

	if u.LoginEmail == "" {
		errs.AddNew(ecase("login_email", "must not be empty"))
	}

	if u.APIAuthKey == "" {
		errs.AddNew(ecase("api_auth_key", "must not be empty"))
	}

	if u.OrgKey == "" {
		errs.AddNew(ecase("org_key", "must not be empty"))
	}

	return errs.OrNil()
}

func (u User) Equals(o User) bool {
	return u.UserKey == o.UserKey &&
		u.LoginEmail == o.LoginEmail &&
		u.APIAuthKey == o.APIAuthKey &&
		u.OrgKey == o.OrgKey &&
		tbntime.Equal(u.DeletedAt, o.DeletedAt) &&
		u.Checksum == o.Checksum
}
