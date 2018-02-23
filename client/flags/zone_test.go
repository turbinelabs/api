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

import (
	"testing"

	tbnflag "github.com/turbinelabs/nonstdlib/flag"
	"github.com/turbinelabs/test/assert"
)

func TestNewZoneFromFlags(t *testing.T) {
	flagset := tbnflag.NewTestFlagSet()
	ff := NewZoneFromFlags(flagset)
	flagset.Parse([]string{"-zone-name=red-sector-a"})
	assert.Equal(t, ff.Name(), "red-sector-a")
}
