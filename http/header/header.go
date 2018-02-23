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

// Package header defines constants for HTTP headers used by the Turbine Labs API
package header

const (
	// APIKey is the caller's API key. Used to associate a request with a
	// specific user and, by extension, organization.
	APIKey = "X-Turbine-Api-Key"

	// Authorization is the Standard HTTP authorization header.
	Authorization = "Authorization"

	// ClientApp is the name of the application calling the API
	ClientApp = "X-Tbn-Api-Client-App"

	// ClientType indicates which kind of API client is making a request (go http
	// client vs js web client etc.)
	ClientType = "X-Tbn-Api-Client-Type"

	// ClientVersion indicates the version of the API client
	ClientVersion = "X-Tbn-Api-Client-Version"
)

var (
	headers = []string{
		APIKey,
		Authorization,
		ClientApp,
		ClientType,
		ClientVersion,
	}
)

// Headers returns the list of headers used for the API.
func Headers() []string {
	return headers
}
