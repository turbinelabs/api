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

package v1

// ZeroFill controls whether and how missing data is filled in query results.
type ZeroFill string

const (
	// None is the default mode and doesn't fill any values in. This is the
	// default behavior.
	None ZeroFill = "none"

	// Partial fills in only series that are partially complete
	Partial ZeroFill = "partial"

	// Full fills in all series even if there was no data initially. If this is
	// the case EmptySeries will be set on the TimeSeries.
	Full ZeroFill = "full"
)

func (zf ZeroFill) IsNone() bool {
	return zf == None || !(zf.IsPartial() || zf.IsFull())
}

func (zf ZeroFill) IsPartial() bool {
	return zf == Partial
}

func (zf ZeroFill) IsFull() bool {
	return zf == Full
}
