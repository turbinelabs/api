package api

import (
	"fmt"
)

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

func (o Org) IsValid(precreation bool) *ValidationError {
	ecase := func(f, m string) ErrorCase {
		return ErrorCase{fmt.Sprintf("org.%s", f), m}
	}

	errs := &ValidationError{}

	validOrgKey := (precreation || o.OrgKey != "")
	validName := o.Name != ""
	validEmail := o.ContactEmail != ""

	if !validOrgKey {
		errs.AddNew(ecase("org_key", "must not be empty"))
	}

	if !validName {
		errs.AddNew(ecase("name", "must not be empty"))
	}

	if !validEmail {
		errs.AddNew(ecase("login_email", "must not be empty"))
	}

	return errs.OrNil()
}

func (o Org) Equals(ot Org) bool {
	return o.OrgKey == ot.OrgKey &&
		o.Name == ot.Name &&
		o.ContactEmail == ot.ContactEmail &&
		o.Checksum == ot.Checksum
}
