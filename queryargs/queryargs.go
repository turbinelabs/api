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

// Package queryargs is a collection of all query arguments that might be passed
// to the Turbine Labs API.
package queryargs

const (
	// On mutating requests this can carry a comment describing why the change
	// was made.
	ChangeComment string = "comment"

	// When mutating an object a checksum is required; this query arg holds the
	// expected value.
	Checksum = "checksum"

	// Index handlers can take JSON encoded filters as one of two ways to provide
	// index query configuration. If present the JSON encoding takes precedence
	IndexFilters = "filters"

	// Duration specifies how much time we should include when the recent changes
	// endpoint is requested.
	Duration = "duration"

	// WindowStart indicates a start of some bounded time frame.
	WindowStart = "start"

	// WindowStop specifies the end of a bounded time frame.
	WindowStop = "end"

	// Should the lookup include deleted objects?
	IncludeDeleted = "include_deleted"
)
