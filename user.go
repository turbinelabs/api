package api

import (
	"fmt"
	"time"

	tbntime "github.com/turbinelabs/nonstdlib/time"
	tbnauth "github.com/turbinelabs/server/auth"
)

type UserKey string
type Users []User

// A User is an actor of an Org. An API key is set if they're allowed to make
// API calls.
type User struct {
	UserKey    UserKey            `json:"user_key"`
	LoginEmail string             `json:"login_email"`
	APIAuthKey tbnauth.APIAuthKey `json:"api_auth_key"`
	OrgKey     OrgKey             `json:"org_key"`
	DeletedAt  *time.Time         `json:"deleted_at"`
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
