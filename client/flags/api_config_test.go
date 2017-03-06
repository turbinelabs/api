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

import (
	"testing"

	tbnflag "github.com/turbinelabs/nonstdlib/flag"
	"github.com/turbinelabs/test/assert"
)

func TestNewAPIConfigFromFlags(t *testing.T) {
	flagset := tbnflag.NewTestFlagSet()

	ff := NewAPIConfigFromFlags(flagset)
	ffImpl := ff.(*apiConfigFromFlags)

	flagset.Parse([]string{"-key=schlage"})

	assert.Equal(t, ffImpl.apiKeyConfig.Make(), "schlage")

	theFlag := flagset.Unwrap().Lookup("key")
	assert.NonNil(t, theFlag)
	assert.True(t, tbnflag.IsRequired(theFlag))

	assert.NonNil(t, ffImpl.FromFlags)
}

func TestNewAPIConfigFromFlagsWithPrefix(t *testing.T) {
	flagset := tbnflag.NewTestFlagSet()
	apiScopedFlagSet := flagset.Scope("api", "test")
	ff := NewAPIConfigFromFlags(
		apiScopedFlagSet,
		APIConfigSetAPIAuthKeyFromFlags(
			NewAPIAuthKeyFromFlags(apiScopedFlagSet, APIAuthKeyFlagsOptional()),
		),
	)
	ffImpl := ff.(*apiConfigFromFlags)

	flagset.Parse([]string{"-api.key=schlage"})

	assert.Equal(t, ffImpl.apiKeyConfig.Make(), "schlage")

	theFlag := flagset.Unwrap().Lookup("api.key")
	assert.NonNil(t, theFlag)
	assert.False(t, tbnflag.IsRequired(theFlag))

	assert.NonNil(t, ffImpl.FromFlags)
}

func TestAPIConfigFromFlagsGet(t *testing.T) {
	ff := &apiConfigFromFlags{
		apiKeyConfig: &apiAuthKeyFromFlags{optional: false, apiKey: "schlage"},
	}
	assert.Equal(t, ff.APIKey(), "schlage")
}
