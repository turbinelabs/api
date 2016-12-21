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

//go:generate mockgen -source $GOFILE -destination mock_$GOFILE -package $GOPACKAGE -aux_files "apihttp=../../http/fromflags.go"

import (
	"flag"

	apihttp "github.com/turbinelabs/api/http"
	tbnflag "github.com/turbinelabs/nonstdlib/flag"
)

// APIConfigFromFlags represents command-line flags for specifying an
// API authentication key, host, port and SSL settings for the Turbine
// Labs API.
type APIConfigFromFlags interface {
	apihttp.FromFlags

	// APIKey Returns the API authentication key from the command line.
	// Equivalent to calling APIAuthKeyFromFlags().Make()
	APIKey() string

	// APIAuthKeyFromFlags returns the underlying APIAuthKeyFromFlags
	// so that it can potentially be shared between APIConfigFromFlags
	// via the APIConfigSetAPIAuthKeyFromFlags APIConfigOption.
	APIAuthKeyFromFlags() APIAuthKeyFromFlags
}

// NewAPIConfigFromFlags configures the necessary command line flags
// and returns an APIConfigFromFlags.
func NewAPIConfigFromFlags(flagset *flag.FlagSet) APIConfigFromFlags {
	return NewPrefixedAPIConfigFromFlags(prefixedFlagSet(flagset))
}

type APIConfigOption func(*apiConfigFromFlags)

// APIConfigSetAPIAuthKeyFromFlags allows the caller to specify a shared
// APIAuthKeyFromFlags, likely obtained via the
// APIConfigFromFlags.APIAuthKeyFromFlags() method.
func APIConfigSetAPIAuthKeyFromFlags(akff APIAuthKeyFromFlags) APIConfigOption {
	return func(ff *apiConfigFromFlags) {
		ff.apiKeyConfig = akff
	}
}

// NewPrefixedAPIConfigFromFlags configures the necessary command
// line flags with a custom prefix and returns an APIConfigFromFlags.
func NewPrefixedAPIConfigFromFlags(
	flagset *tbnflag.PrefixedFlagSet,
	opts ...APIConfigOption,
) APIConfigFromFlags {
	ff := &apiConfigFromFlags{requiredFlag: true}

	for _, applyOpt := range opts {
		applyOpt(ff)
	}

	if ff.apiKeyConfig == nil {
		ff.apiKeyConfig = NewPrefixedAPIAuthKeyFromFlags(flagset, ff.requiredFlag)
	}

	ff.FromFlags = apihttp.NewFromFlags("api.turbinelabs.io", flagset)

	return ff
}

type apiConfigFromFlags struct {
	apihttp.FromFlags
	apiKeyConfig APIAuthKeyFromFlags
	requiredFlag bool
}

func (ff *apiConfigFromFlags) APIKey() string {
	return ff.apiKeyConfig.Make()
}

func (ff *apiConfigFromFlags) APIAuthKeyFromFlags() APIAuthKeyFromFlags {
	return ff.apiKeyConfig
}
