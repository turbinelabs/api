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
	"flag"
	"testing"

	tbnflag "github.com/turbinelabs/nonstdlib/flag"
	"github.com/turbinelabs/test/assert"
)

func TestNewAPIAuthKeyFromFlags(t *testing.T) {
	flagset := flag.NewFlagSet("NewAPIAuthFromFlags options", flag.PanicOnError)

	ff := NewAPIAuthKeyFromFlags(flagset)
	ffImpl := ff.(*apiAuthKeyFromFlags)

	flagset.Parse([]string{"-api.key=schlage"})

	assert.Equal(t, ffImpl.apiKey, "schlage")

	theFlag := flagset.Lookup("api.key")
	assert.NonNil(t, theFlag)
	assert.True(t, tbnflag.IsRequired(theFlag))
}

func TestNewPrefixedAPIAuthKeyFromFlags(t *testing.T) {
	flagset := flag.NewFlagSet("NewAPIAuthFromFlags options", flag.PanicOnError)
	prefixedFlagset := tbnflag.NewPrefixedFlagSet(flagset, "test", "test")

	ff := NewPrefixedAPIAuthKeyFromFlags(prefixedFlagset, false)
	ffImpl := ff.(*apiAuthKeyFromFlags)

	flagset.Parse([]string{"-test.key=schlage"})

	assert.Equal(t, ffImpl.apiKey, "schlage")

	theFlag := flagset.Lookup("test.key")
	assert.NonNil(t, theFlag)
	assert.False(t, tbnflag.IsRequired(theFlag))
}

func TestAPIAuthKeyFromFlagsGet(t *testing.T) {
	ff := &apiAuthKeyFromFlags{apiKey: "schlage"}
	assert.Equal(t, ff.Make(), "schlage")
}
