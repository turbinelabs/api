package api

import (
	"fmt"
)

type ZoneKey string

type Zones []Zone

// A Zone is a logical grouping that various other objects can be associated
// with. A single zone will often map to some geographic region.
type Zone struct {
	ZoneKey ZoneKey `json:"zone_key"`
	Name    string  `json:"name"`
	OrgKey  OrgKey  `json:"-"`
	Checksum
}

func (z Zone) IsNil() bool {
	return z.Equals(Zone{})
}

func (z Zone) Equals(o Zone) bool {
	return z.ZoneKey == o.ZoneKey &&
		z.Name == o.Name &&
		z.OrgKey == o.OrgKey &&
		z.Checksum.Equals(o.Checksum) &&
		z.OrgKey == o.OrgKey
}

func (z Zone) IsValid(precreation bool) *ValidationError {
	ecase := func(f, m string) ErrorCase {
		return ErrorCase{fmt.Sprintf("zone.%s", f), m}
	}

	errs := &ValidationError{}

	keyValid := precreation || z.ZoneKey != ""
	if !keyValid {
		errs.AddNew(ecase("zone_key", "must not be empty"))
	}

	if z.OrgKey == "" {
		errs.AddNew(ecase("org_key", "must not be empty"))
	}

	if z.Name == "" {
		errs.AddNew(ecase("name", "must not be empty"))
	}

	return errs.OrNil()
}
