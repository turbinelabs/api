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

// Package envelope contains the Response envelope used by the Turbine Labs
// public API to encapsulate server behavior.
package envelope

import httperr "github.com/turbinelabs/api/http/error"

// Response is constructed at API render time to enable a predictable way to
// transmit error and request payload to a HTTP client. It is received by the
// HTTP client and unpacked into the appropriate types depending on the call
// being made.
type Response struct {
	Error   *httperr.Error `json:"error,omitempty"`
	Payload interface{}    `json:"result,omitempty"`
}
