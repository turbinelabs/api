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
