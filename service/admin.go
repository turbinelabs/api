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

package service

//go:generate mockgen -source $GOFILE -destination mock_$GOFILE -package $GOPACKAGE --write_package_comment=false

import (
	"reflect"
	"time"

	"github.com/turbinelabs/api"
)

/*
	Admin defines the interface for the public JSON/REST administrative
	API.

	See All for a discussion of the methodology behind this interface.
*/
type Admin interface {
	User() User

	// AccessToken returns an interface to interact with the access tokens
	// for the user who is making an authenticated request.
	AccessToken() AccessToken
}

type User interface {
	// GET /v1.0/admin/user
	//
	// Index returns all Users to which the given filters apply. All non-empty
	// fields in a filter must apply for the filter to apply. Any User to which
	// any filter applies is included in the result.
	//
	// If no filters are supplied, all Users are returned.
	Index(filters ...UserFilter) (api.Users, error)

	// GET /v1.0/admin/user/<string:userKey>[?include_deleted]
	//
	// Get returns a User for the given UserKey. If the User does not
	// exist, an error is returned.
	Get(userKey api.UserKey) (api.User, error)

	// POST /v1.0/admin/user
	//
	// Create creates the given User. User Names must be unique for a given
	// ZoneKey. If a UserKey is specified in the User, it is ignored and
	// replaced in the result with the authoritative UserKey.
	Create(user api.User) (api.User, error)

	// PUT /v1.0/admin/user/<string:userKey>
	//
	// Modify modifies the given User. User Names must be unique for a given
	// ZoneKey. The given User Checksum must match the existing Checksum.
	Modify(user api.User) (api.User, error)

	// DELETE /v1.0/admin/user/<string:userKey>?checksum=<string:checksum>
	//
	// Delete marks the User corresponding to the given UserKey as deleted.
	// The given User Checksum must match the existing Checksum. The
	// timestamp of the DB operation is used as a deletion time.
	Delete(userKey api.UserKey, checksum api.Checksum) error
}

type UserFilter struct {
	UserKey       api.UserKey    `json:"user_key"`
	LoginEmail    string         `json:"login_email"`
	APIAuthKey    api.APIAuthKey `json:"api_auth_key"`
	OrgKey        api.OrgKey     `json:"org_key"`
	Active        *bool          `json:"active"`
	DeletedBefore *time.Time     `json:"deleted_before"`
	DeletedAfter  *time.Time     `json:"deleted_after"`
}

func (uf UserFilter) IsNil() bool {
	return uf.Equals(UserFilter{})
}

func (uf UserFilter) Equals(o UserFilter) bool {
	activeEquals := reflect.DeepEqual(uf.Active, o.Active)
	dbEquals := reflect.DeepEqual(uf.DeletedBefore, o.DeletedBefore)
	daEquals := reflect.DeepEqual(uf.DeletedAfter, o.DeletedAfter)

	return uf.UserKey == o.UserKey &&
		uf.LoginEmail == o.LoginEmail &&
		uf.APIAuthKey == o.APIAuthKey &&
		uf.OrgKey == o.OrgKey &&
		activeEquals &&
		dbEquals &&
		daEquals
}

type Org interface {
	// GET /v1.0/admin/org
	//
	// Index returns all Orgs to which the given filters apply. All non-empty
	// fields in a filter must apply for the filter to apply. Any Org to which
	// any filter applies is included in the result.
	//
	// If no filters are supplied, all Orgs are returned.
	Index(filters ...OrgFilter) (api.Orgs, error)

	// GET /v1.0/admin/org/<string:orgKey>[?include_deleted]
	//
	// Get returns an Org for the given OrgKey. If the Org does not
	// exist, an error is returned.
	Get(orgKey api.OrgKey) (api.Org, error)

	// POST /v1.0/admin/org
	//
	// Create creates the given Org. Org Names must be unique for a given
	// ZoneKey. If a OrgKey is specified in the Org, it is ignored and
	// replaced in the result with the authoritative OrgKey.
	Create(org api.Org) (api.Org, error)

	// PUT /v1.0/admin/org/<string:orgKey>
	//
	// Modify modifies the given Org. Org Names must be unique for a given
	// ZoneKey. The given Org Checksum must match the existing Checksum.
	Modify(org api.Org) (api.Org, error)

	// DELETE /v1.0/admin/org/<string:orgKey>?checksum=<checksum>
	//
	// Delete completely removes the Org data from the database.
	// If the checksum does not match no action is taken and an error
	// is returned.
	Delete(orgKey api.OrgKey, checksum api.Checksum) error
}

type OrgFilter struct {
	OrgKey       api.OrgKey `json:"org_key"`
	Name         string     `json:"name"`
	ContactEmail string     `json:"contact_email"`
}

func (of OrgFilter) IsNil() bool {
	return of.Equals(OrgFilter{})
}

func (of OrgFilter) Equals(o OrgFilter) bool {
	return of == o
}

type AccessToken interface {
	// GET /v1.0/admin/user/self/access_token
	//
	// Index returns all AccessTokens to which the given filters apply. All non-empty
	// fields in a filter must apply for the filter to apply. Any AccessToken to which
	// any filter applies is included in the result.
	//
	// If no filters are supplied, all AccessTokens are returned.
	Index(filters ...AccessTokenFilter) (api.AccessTokens, error)

	// GET /v1.0/admin/user/self/access_token/<string:AccessTokenKey>
	//
	// Get returns an AccessToken for the given AccessTokenKey. If the AccessToken
	// does not exist, an error is returned.
	Get(key api.AccessTokenKey) (api.AccessToken, error)

	// POST /v1.0/admin/user/self/access_token
	//
	// Create creates the given AccessToken. AccessToken description should explain
	// the intended use of the key that will be issued. Any fields other than
	// Description specified in the POSTed AccessToken are ignored.
	//
	// The response will have the SignedToken field populated and is the only chance a
	// caller will have to save this value. An AccessToken may be revoked by calling
	// Delete below.
	Create(token api.AccessToken) (api.AccessToken, error)

	// DELETE /v1.0/admin/user/self/access_token/<string:AccessTokenKey>?checksum=<checksum>
	//
	// Delete completely removes the AccessToken data from the database.
	// If the checksum does not match no action is taken and an error
	// is returned.
	Delete(key api.AccessTokenKey, checksum api.Checksum) error
}
