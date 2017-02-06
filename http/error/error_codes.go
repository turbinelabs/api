package error

// Identifier specifying what kind of error we're reporting.
type ErrorCode string

const (
	DataConstraintErrorCode ErrorCode = "DataConstraintErrorCode" // some integrity constraint rejected the change
	OwningOrgImmutable      ErrorCode = "OwningOrgImmutable"      // an attempt was made to change the org an object belongs to

	DuplicateKeyErrorCode  ErrorCode = "DuplicateClusterKeyCode" // an object with that key already exists
	InvalidObjectErrorCode ErrorCode = "InvalidObjectCode"       // attempted to pass on object that did not have valid data
	KeyImmutableErrorCode  ErrorCode = "KeyImmutable"            // can't change the generated object key
	MiscErrorCode          ErrorCode = "MiscCode"                // misc errors
	NotFoundErrorCode      ErrorCode = "NotFound"                // the object requested didn't exist
	AuthMethodDeniedCode   ErrorCode = "AuthMethodDenied"        // auth method not allowed
	BadParameterErrorCode  ErrorCode = "BadParameterErrorCode"   // some parameter was bad
)

const (
	// there was a problem decoding something
	UnknownDecodingCode ErrorCode = "UnknownDecodingCode"

	// there was a problem encoding something
	UnknownEncodingCode ErrorCode = "UnknownEncodingCode"

	// expected some content from the request body but could not find it
	UnknownNoBodyCode ErrorCode = "UnknownNoBodyCode"

	// Something involving manipulating HTTP response/request
	UnknownTransportCode ErrorCode = "UnknownTransportCode"

	// authorization for this request failed
	UnknownUnauthorizedCode ErrorCode = "UnknownUnathorizedCode"

	// unclassified failure
	UnknownUnclassifiedCode ErrorCode = "UnknownUnclassifiedCode"

	// attempted to save data but it was modified & you were working with an old copy
	UnknownModificationConflict ErrorCode = "UnknownModificationConflict"
)
