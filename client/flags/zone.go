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

package flags

//go:generate mockgen -source $GOFILE -destination mock_$GOFILE -package $GOPACKAGE

import (
	"github.com/turbinelabs/api/service"
	tbnflag "github.com/turbinelabs/nonstdlib/flag"
	"github.com/turbinelabs/nonstdlib/flag/usage"
)

// ZoneFromFlags represents command-line flags for specifying a
// Turbine Labs API zone name, which is used to resolve a zone.
type ZoneFromFlags interface {
	Name() string
	Ref() service.ZoneRef
}

// ZoneFromFlagsFlagOptions lets you add various options to NewFromFlags
type ZoneFromFlagsFlagOptions func(*zoneFromFlags) *zoneFromFlags

// ZoneFromFlagsNameOptional allows the caller to specify that the
// zone-name flag is optional
func ZoneFromFlagsNameOptional() ZoneFromFlagsFlagOptions {
	return func(ff *zoneFromFlags) *zoneFromFlags {
		ff.optional = true
		return ff
	}
}

// NewZoneFromFlags configures the necessary command line flags to
// retrieve a zone by zone name.
func NewZoneFromFlags(flagset tbnflag.FlagSet, opts ...ZoneFromFlagsFlagOptions) ZoneFromFlags {
	ff := &zoneFromFlags{}

	for _, fn := range opts {
		ff = fn(ff)
	}

	u := "The name of the API Zone for {{NAME}} requests."
	if !ff.optional {
		u = usage.Required(u)
	}
	flagset.StringVar(&ff.zoneName, "zone-name", "", u)

	return ff
}

type zoneFromFlags struct {
	optional bool
	zoneName string
}

func (ff *zoneFromFlags) Name() string {
	return ff.zoneName
}

func (ff *zoneFromFlags) Ref() service.ZoneRef {
	return service.NewZoneNameZoneRef(ff.zoneName)
}
