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
	"github.com/turbinelabs/nonstdlib/flag/usage"
	"github.com/turbinelabs/test/assert"
)

func TestNewAPIAuthKeyFromFlags(t *testing.T) {
	flagset := tbnflag.NewTestFlagSet()

	ff := NewAPIAuthKeyFromFlags(flagset)
	ffImpl := ff.(*apiAuthKeyFromFlags)

	flagset.Parse([]string{"-key=schlage"})

	assert.Equal(t, ffImpl.apiKey, "schlage")

	theFlag := flagset.Unwrap().Lookup("key")
	assert.NonNil(t, theFlag)
	assert.True(t, usage.IsRequired(theFlag))
}

func TestNewAPIAuthKeyFromFlagsOptionalWithPrefix(t *testing.T) {
	flagset := tbnflag.NewTestFlagSet()

	ff := NewAPIAuthKeyFromFlags(flagset.Scope("test", "test"), APIAuthKeyFlagsOptional())
	ffImpl := ff.(*apiAuthKeyFromFlags)

	flagset.Parse([]string{"-test.key=schlage"})

	assert.Equal(t, ffImpl.apiKey, "schlage")

	theFlag := flagset.Unwrap().Lookup("test.key")
	assert.NonNil(t, theFlag)
	assert.False(t, usage.IsRequired(theFlag))
}

func TestAPIAuthKeyFromFlagsGet(t *testing.T) {
	ff := &apiAuthKeyFromFlags{apiKey: "schlage"}
	assert.Equal(t, ff.Make(), "schlage")
}
