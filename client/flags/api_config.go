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
	"fmt"

	"golang.org/x/net/context"
	"golang.org/x/oauth2"

	"github.com/turbinelabs/api/client/tokencache"
	apihttp "github.com/turbinelabs/api/http"
	tbnflag "github.com/turbinelabs/nonstdlib/flag"
	"github.com/turbinelabs/nonstdlib/log/console"
)

// APIConfigFromFlags represents command-line flags for specifying
// API authentication, host, port and SSL settings for the Turbine
// Labs API.
type APIConfigFromFlags interface {
	// apihttp.FromFlags constructs an API endpoint
	apihttp.FromFlags

	// APIKey Returns the API authentication key from the command line.
	// Equivalent to calling APIAuthKeyFromFlags().Make()
	APIKey() string

	// APIAuthKeyFromFlags returns the underlying APIAuthKeyFromFlags
	// so that it can potentially be shared between APIConfigFromFlags
	// via the APIConfigSetAPIAuthKeyFromFlags APIConfigOption.
	APIAuthKeyFromFlags() APIAuthKeyFromFlags
}

// APIConfigOption represents an option passed to
// NewAPIConfigFromFlags.
type APIConfigOption func(*apiConfigFromFlags)

// APIConfigSetAPIAuthKeyFromFlags allows the caller to specify a shared
// APIAuthKeyFromFlags, likely obtained via the
// APIConfigFromFlags.APIAuthKeyFromFlags() method.
func APIConfigSetAPIAuthKeyFromFlags(akff APIAuthKeyFromFlags) APIConfigOption {
	return func(ff *apiConfigFromFlags) {
		ff.apiKeyConfig = akff
	}
}

// APIConfigMayUseAuthToken indicates that the API Config can use an auth
// token to authenticate instead of relying on an API key. As a result if
// this is set API Key becomes optional. If an API Key is set it takes
// precedence over an auth token found from the TokenCache.
//
// If an APIAuthKeyFromFlags is provided via APIConfigSetAPIAuthKeyFromFlags
// it may still be considered required depending on the APIAuthKeyFromFlags
// construction.
//
// Parameters:
//   cachePath - this is a file where cached authed token should be read from
func APIConfigMayUseAuthToken(cachePath tokencache.PathFromFlags) APIConfigOption {
	return func(ff *apiConfigFromFlags) {
		ff.mayUseAuthToken = true
		ff.cachePath = cachePath
	}
}

// NewAPIConfigFromFlags configures the necessary command line flags
// and returns an APIConfigFromFlags.
func NewAPIConfigFromFlags(
	flagset tbnflag.FlagSet,
	opts ...APIConfigOption,
) APIConfigFromFlags {
	ff := &apiConfigFromFlags{}

	for _, applyOpt := range opts {
		applyOpt(ff)
	}

	if ff.apiKeyConfig == nil {
		opts := []APIAuthKeyOption{}
		if ff.mayUseAuthToken {
			opts = append(opts, APIAuthKeyFlagsOptional())
		}
		ff.apiKeyConfig = NewAPIAuthKeyFromFlags(flagset, opts...)
	}

	ff.clientFromFlags = apihttp.NewFromFlags("api.turbinelabs.io", flagset)
	return ff
}

type apiConfigFromFlags struct {
	clientFromFlags apihttp.FromFlags
	apiKeyConfig    APIAuthKeyFromFlags

	mayUseAuthToken bool
	cachePath       tokencache.PathFromFlags

	oauth2Config oauth2.Config
	tokenCache   tokencache.TokenCache
}

func (ff *apiConfigFromFlags) APIKey() string {
	return ff.apiKeyConfig.Make()
}

func (ff *apiConfigFromFlags) APIAuthKeyFromFlags() APIAuthKeyFromFlags {
	return ff.apiKeyConfig
}

func (ff *apiConfigFromFlags) Validate() error {
	err := ff.clientFromFlags.Validate()
	if err != nil {
		return err
	}

	if ff.mayUseAuthToken && ff.APIKey() != "" {
		console.Info().Println("Preferring API Key for authentication")
	}

	return nil
}

func (ff *apiConfigFromFlags) MakeEndpoint() (apihttp.Endpoint, error) {
	ep, err := ff.clientFromFlags.MakeEndpoint()
	if err != nil {
		return apihttp.Endpoint{}, err
	}

	// If an API Key is present it takes precedence and is assumed valid
	if ff.APIKey() != "" {
		return ep, nil
	}

	tc, err := tokencache.NewFromFile(ff.cachePath.CachePath())
	if err != nil {
		return apihttp.Endpoint{}, err
	}

	if tc.Expired() {
		return apihttp.Endpoint{}, fmt.Errorf("your session has timed out, please login again")
	}
	ff.tokenCache = tc

	cfg, err := tokencache.ToOAuthConfig(tc)
	if err != nil {
		return apihttp.Endpoint{}, fmt.Errorf("unable to construct OIDC client config: %v", err)
	}
	ff.oauth2Config = cfg

	// otherwise use the token from cache
	ctx := context.Background()
	client := ff.oauth2Config.Client(ctx, ff.tokenCache.Token)
	ep.SetClient(apihttp.MakeHeaderPreserving(client))

	return ep, nil
}
