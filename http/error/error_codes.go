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

package error

// Identifier specifying what kind of error we're reporting.
type ErrorCode string

const (
	DataConstraintErrorCode ErrorCode = "DataConstraintErrorCode" // some integrity constraint rejected the change
	OwningOrgImmutable      ErrorCode = "OwningOrgImmutable"      // an attempt was made to change the org an object belongs to

	DuplicateKeyErrorCode      ErrorCode = "DuplicateClusterKeyCode"    // an object with that key already exists
	InvalidObjectErrorCode     ErrorCode = "InvalidObjectCode"          // attempted to pass on object that did not have valid data
	KeyImmutableErrorCode      ErrorCode = "KeyImmutable"               // can't change the generated object key
	MiscErrorCode              ErrorCode = "MiscCode"                   // misc errors
	NotFoundErrorCode          ErrorCode = "NotFound"                   // the object requested didn't exist
	AuthMethodDeniedCode       ErrorCode = "AuthMethodDenied"           // auth method not allowed
	BadParameterErrorCode      ErrorCode = "BadParameterErrorCode"      // some parameter was bad
	ObjectKeyRequiredErrorCode ErrorCode = "ObjectKeyRequiredErrorCode" // returned when a Get attempt is made with no key
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
