/*
	Collection of all query args that might be passed to the Turbine Labs API.
*/
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

	// WindowStop specifices the end of a bounded time frame.
	WindowStop = "end"

	// Should the lookup include deleted objects?
	IncludeDeleted = "include_deleted"
)
