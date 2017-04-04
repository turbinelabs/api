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
	tbnflag "github.com/turbinelabs/nonstdlib/flag"
	"github.com/turbinelabs/nonstdlib/flag/usage"
)

// APIAuthKeyFromFlags represents command-line flags for specifying an
// API authentication key for the Turbine Labs API.
type APIAuthKeyFromFlags interface {
	// Returns the API authentication key from the command line.
	Make() string
}

// APIAuthKeyOption represents an option passed to
// NewAPIAuthKeyFromFlags.
type APIAuthKeyOption func(*apiAuthKeyFromFlags)

// APIAuthKeyFlagsOptional allows the caller to specify that the flags
// created by APIAuthKeyFromFlags are optional.
func APIAuthKeyFlagsOptional() APIAuthKeyOption {
	return func(ff *apiAuthKeyFromFlags) {
		ff.optional = true
	}
}

// NewAPIAuthKeyFromFlags configures the necessary command line flags
// and returns an APIAuthKeyFromFlags.
func NewAPIAuthKeyFromFlags(flagset tbnflag.FlagSet, opts ...APIAuthKeyOption) APIAuthKeyFromFlags {
	ff := &apiAuthKeyFromFlags{}

	for _, apply := range opts {
		apply(ff)
	}

	u := usage.New("The auth key for {{NAME}} requests").SetSensitive()
	if !ff.optional {
		u = u.SetRequired()
	}

	flagset.StringVar(&ff.apiKey, "key", "", u.String())

	return ff
}

type apiAuthKeyFromFlags struct {
	optional bool
	apiKey   string
}

func (ff *apiAuthKeyFromFlags) Make() string {
	return ff.apiKey
}
