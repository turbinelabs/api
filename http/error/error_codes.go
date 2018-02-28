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

package error

// ErrorCode is an identifier specifying what kind of error we're reporting.
type ErrorCode string

const (
	// DataConstraintErrorCode indicates some integrity constraint rejected the change.
	DataConstraintErrorCode ErrorCode = "DataConstraintErrorCode"

	// OwningOrgImmutable indicates an attempt was made to change the org to which an
	// object belongs.
	OwningOrgImmutable ErrorCode = "OwningOrgImmutable"

	// DuplicateKeyErrorCode indicates an object with the given key already exists.
	DuplicateKeyErrorCode ErrorCode = "DuplicateClusterKeyCode"

	// InvalidObjectErrorCode indicates an attempt to pass on invalid object.
	InvalidObjectErrorCode ErrorCode = "InvalidObjectCode"

	// KeyImmutableErrorCode indicates an attempt was made to change an object's key.
	KeyImmutableErrorCode ErrorCode = "KeyImmutable"

	// MiscErrorCode indicates a miscellaneous error.
	MiscErrorCode ErrorCode = "MiscCode"

	// NotFoundErrorCode indicates the request object does not exist.
	NotFoundErrorCode ErrorCode = "NotFound"

	// AuthMethodDeniedCode indicates an authorization method was denied.
	AuthMethodDeniedCode ErrorCode = "AuthMethodDenied"

	// BadParameterErrorCode indicates an invalid parameter was passed.
	BadParameterErrorCode ErrorCode = "BadParameterErrorCode"

	// ObjectKeyRequiredErrorCode indicates a Get attempt with no key was made.
	ObjectKeyRequiredErrorCode ErrorCode = "ObjectKeyRequiredErrorCode"
)

const (
	// UnknownDecodingCode indicates there was a problem decoding something.
	UnknownDecodingCode ErrorCode = "UnknownDecodingCode"

	// UnknownEncodingCode indicates there was a problem encoding something.
	UnknownEncodingCode ErrorCode = "UnknownEncodingCode"

	// UnknownNoBodyCode indicates the server expected some content from the request
	// body but could not find it.
	UnknownNoBodyCode ErrorCode = "UnknownNoBodyCode"

	// UnknownTransportCode indicates an error manipulating HTTP response/request.
	UnknownTransportCode ErrorCode = "UnknownTransportCode"

	// UnknownUnauthorizedCode indicates that authorization for this request failed.
	UnknownUnauthorizedCode ErrorCode = "UnknownUnathorizedCode"

	// UnknownUnclassifiedCode indicates an unclassified failure.
	UnknownUnclassifiedCode ErrorCode = "UnknownUnclassifiedCode"

	// UnknownModificationConflict indicates an attempt to save data failed because
	// the request was made against an out-of-date copy.
	UnknownModificationConflict ErrorCode = "UnknownModificationConflict"
)
