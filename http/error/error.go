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
	Package error contains an error definition intended to be serialized
	and sent over the wire between the api-server and clients.

	An error can be associated with either one of the top level domain objects
	(Clusters, Routes, etc.) or not (Unknown). Examples of each would include
	business validation failures (domain object category) and errors processing
	the HTTP request object (unknown categary). This association and the specific
	error is captured, By convention, in the code name: {Type}{Error}Code.

	A service.http.error.Error is created in the api-server by processing a
	server.model.error.Error and converting it to something client appropriate.
	See server.handler docs for more details on this process.
*/
package error

import (
	"fmt"
	"net/http"
)

// on the server side we will map (as noted above) from server.error.Error to a
// domain-specific service.http.error.Error; on the client side we'll just
// decode into this.
type Error struct {
	Message string      `json:"message"`           // Human readable message
	Code    ErrorCode   `json:"code"`              // This provides more detail into the type of error
	Status  int         `json:"-"`                 // what http response code this should be (or was) sent with
	Details interface{} `json:"details,omitempty"` // any error specific information
}

// Construct a new error code that will be sent to a client
func NewDetailed400(msg string, code ErrorCode, details interface{}) *Error {
	return &Error{msg, code, http.StatusBadRequest, details}
}

func New400(msg string, code ErrorCode) *Error {
	return NewDetailed400(msg, code, nil)
}

// Construct a new HTTP 500 (InternalServerError) response
func NewDetailed500(msg string, code ErrorCode, details interface{}) *Error {
	return &Error{msg, code, http.StatusInternalServerError, details}
}

func New500(msg string, code ErrorCode) *Error {
	return NewDetailed500(msg, code, nil)
}

// Construct a new HTTP 501 (NotImplemented) response
func NewDetailed501(msg string, code ErrorCode, details interface{}) *Error {
	return &Error{msg, code, http.StatusNotImplemented, details}
}

func New501(msg string, code ErrorCode) *Error {
	return NewDetailed501(msg, code, nil)
}

// Construct a new HTTP 404 (NotFound) response
func NewDetailed404(msg string, code ErrorCode, details interface{}) *Error {
	return &Error{msg, code, http.StatusNotFound, details}
}

func New404(msg string, code ErrorCode) *Error {
	return NewDetailed404(msg, code, nil)
}

// Construct a new HTTP 409 (Conflict) response
func NewDetailed409(msg string, code ErrorCode, details interface{}) *Error {
	return &Error{msg, code, http.StatusConflict, details}
}

func New409(msg string, code ErrorCode) *Error {
	return NewDetailed409(msg, code, nil)
}

func AuthorizationError() *Error {
	return &Error{
		"not authorized to make this request",
		UnknownUnauthorizedCode,
		http.StatusForbidden,
		nil,
	}
}

func AuthorizationMethodDeniedError() *Error {
	return &Error{
		"authorization method denied",
		AuthMethodDeniedCode,
		http.StatusForbidden,
		nil,
	}
}

// FromErrorDefaultType returns an *Error based on the provided base error
// object. If err is already an *Error it will be cast and returned unchanged,
// otherwise a new *Error will be created with the specifid code and status.
// The Message will be err.Error().
func FromErrorDefaultType(err error, code ErrorCode, status int) *Error {
	if err == nil {
		return nil
	}

	switch httpErr := err.(type) {
	case *Error:
		return httpErr
	default:
		return &Error{
			err.Error(),
			code,
			status,
			nil,
		}
	}
}

// Generically converts a go error into an Error. If err is already an
// Error, it is returned directly (as a pointer). If not, a new Error
// with status 500 (internal server error), the given code, and the
// original error message is returned.
func FromError(err error, code ErrorCode) *Error {
	return FromErrorDefaultType(err, code, 500)
}

func (e *Error) Error() string {
	hd := e.Details != nil
	return fmt.Sprintf("{Message: %s, Code: %s, Detailed: %v}", e.Message, e.Code, hd)
}
