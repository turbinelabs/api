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

/*
	The http package provides common HTTP client building blocks for
	Turbine Labs servers.

	Because the Turbine api server requires various headers on a
	request if we use the default http.Client as our method of
	making requests we won't be able to gracefully handle
	redirects.

	This package provides a http.Client with a customized (but
	still trivial) redirect policy.
*/
package http
