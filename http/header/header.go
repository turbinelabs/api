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

// package header defines constants for HTTP headers used by the Turbine Labs API
package header

const (
	// APIKey is the caller's API key. Used to associate a request with a
	// specific user and, by extension, organization.
	APIKey = "X-Turbine-Api-Key"

	// Standard HTTP authorization header.
	Authorization = "Authorization"

	// ClientID indicates which Turbine client is making a request (go http
	// client vs js web client etc.)
	ClientID = "X-Turbine-Api-Clientid"
)

var (
	headers = []string{
		APIKey,
		Authorization,
		ClientID,
	}
)
