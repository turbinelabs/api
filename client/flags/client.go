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

package flags

//go:generate mockgen -source $GOFILE -destination mock_$GOFILE -package $GOPACKAGE

import (
	"github.com/turbinelabs/api/client"
	"github.com/turbinelabs/api/service"
	tbnflag "github.com/turbinelabs/nonstdlib/flag"
)

// ClientFromFlags represents command-line flags specifying
// configuration of a service.All backed by the Turbine Labs API.
type ClientFromFlags interface {
	Validate() error

	// Make produces a service.All from the provided flags, or an
	// error.
	Make() (service.All, error)

	// MakeAdmin produces a service.Admin from the provided flags,
	// or an error.
	MakeAdmin() (service.Admin, error)
}

// NewClientFromFlags creates a ServiceFromFlags, which configures the
// necessary flags to construct a service.All instance.
func NewClientFromFlags(clientApp client.App, flagset tbnflag.FlagSet) ClientFromFlags {
	return NewClientFromFlagsWithSharedAPIConfig(clientApp, flagset, nil)
}

// NewClientFromFlagsWithSharedAPIConfig creates a ClientFromFlags,
// which configures the necessary flags to construct a service.All
// instance. The given APIConfigFromFlags is used to obtain the API
// auth key.
func NewClientFromFlagsWithSharedAPIConfig(
	clientApp client.App,
	flagset tbnflag.FlagSet,
	apiConfigFromFlags APIConfigFromFlags,
) ClientFromFlags {
	ff := &clientFromFlags{clientApp: clientApp}

	if apiConfigFromFlags == nil {
		ff.apiConfigFromFlags = NewAPIConfigFromFlags(flagset)
	} else {
		ff.apiConfigFromFlags = apiConfigFromFlags
	}

	return ff
}

type clientFromFlags struct {
	clientApp          client.App
	apiConfigFromFlags APIConfigFromFlags
}

func (ff *clientFromFlags) Validate() error {
	return ff.apiConfigFromFlags.Validate()
}

func (ff *clientFromFlags) Make() (service.All, error) {
	endpoint, err := ff.apiConfigFromFlags.MakeEndpoint()
	if err != nil {
		return nil, err
	}

	apiKey := ff.apiConfigFromFlags.APIKey()

	return client.NewAll(endpoint, apiKey, ff.clientApp)
}

func (ff *clientFromFlags) MakeAdmin() (service.Admin, error) {
	endpoint, err := ff.apiConfigFromFlags.MakeEndpoint()
	if err != nil {
		return nil, err
	}

	apiKey := ff.apiConfigFromFlags.APIKey()

	return client.NewAdmin(endpoint, apiKey, ff.clientApp)
}
