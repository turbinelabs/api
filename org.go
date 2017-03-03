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

type OrgKey string
type Orgs []Org

// An Org is a Turbine customer. It is composed of users and is ultimately the
// entity that owns all other objects in our system: Clusters, Routes, etc.
type Org struct {
	OrgKey       OrgKey `json:"org_key"`
	Name         string `json:"name"`
	ContactEmail string `json:"contact_email"`
	Checksum
}

func (o Org) IsNil() bool {
	return o.Equals(Org{})
}

func (o Org) IsValid() *ValidationError {
	scope := func(i string) string { return "org." + i }

	errs := &ValidationError{}

	errCheckKey(string(o.OrgKey), errs, scope("org_key"))
	errCheckIndex(o.Name, errs, scope("name"))

	if o.ContactEmail == "" {
		errs.AddNew(ErrorCase{scope("login_email"), "must not be empty"})
	}

	return errs.OrNil()
}

func (o Org) Equals(ot Org) bool {
	return o.OrgKey == ot.OrgKey &&
		o.Name == ot.Name &&
		o.ContactEmail == ot.ContactEmail &&
		o.Checksum == ot.Checksum
}
