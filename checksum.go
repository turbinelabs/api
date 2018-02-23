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

package api

// A type commonly embedded in other domain objects to ensure modifications
// are being on an underlying object in the expected state.
type Checksum struct {
	Checksum string `json:"checksum"` // may be overwritten
}

func (c *Checksum) IsNil() bool {
	return c.Equals(Checksum{})
}

// An empty checksum is equivalent to an unset checksum.
func (c *Checksum) IsEmpty() bool {
	return len(c.Checksum) == 0
}

func (c Checksum) Equals(o Checksum) bool {
	return c.Checksum == o.Checksum
}
