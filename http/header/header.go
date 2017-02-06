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
	ClientID = "X-Turbine-Api-ClientId"
)
