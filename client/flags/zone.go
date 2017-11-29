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

// NewZoneFromFlags configures the necessary command line flags to
// retrieve a zone by zone name.
func NewZoneFromFlags(flagset tbnflag.FlagSet) ZoneFromFlags {
	ff := &zoneFromFlags{}

	flagset.StringVar(
		&ff.zoneName,
		"zone-name",
		"",
		usage.Required("The name of the API Zone for {{NAME}} requests."),
	)

	return ff
}

type zoneFromFlags struct {
	zoneName string
}

func (ff *zoneFromFlags) Name() string {
	return ff.zoneName
}

func (ff *zoneFromFlags) Ref() service.ZoneRef {
	return service.NewZoneNameZoneRef(ff.zoneName)
}
