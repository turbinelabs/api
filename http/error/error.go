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

// Package error contains an error definition intended to be serialized and sent over
// the wire between the api-server and clients.
package error

import (
	"fmt"
	"net/http"
)

// Error represents an HTTP API error serialized between the api-server and its
// clients.
type Error struct {
	Message string      `json:"message"`           // Human readable message
	Code    ErrorCode   `json:"code"`              // This provides more detail into the type of error
	Status  int         `json:"-"`                 // what http response code this should be (or was) sent with
	Details interface{} `json:"details,omitempty"` // any error specific information
}

func (e *Error) Error() string {
	hd := e.Details != nil
	return fmt.Sprintf("{Message: %s, Code: %s, Detailed: %v}", e.Message, e.Code, hd)
}

// NewDetailed400 constructs a new Error with status code 400, bad request.
func NewDetailed400(msg string, code ErrorCode, details interface{}) *Error {
	return &Error{msg, code, http.StatusBadRequest, details}
}

// New400 constructs a new Error with status code 400, bad request.
func New400(msg string, code ErrorCode) *Error {
	return NewDetailed400(msg, code, nil)
}

// NewDetailed500 constructs a new Error with status code 500, internal server error.
func NewDetailed500(msg string, code ErrorCode, details interface{}) *Error {
	return &Error{msg, code, http.StatusInternalServerError, details}
}

// New500 constructs a new Error with status code 500, internal server error.
func New500(msg string, code ErrorCode) *Error {
	return NewDetailed500(msg, code, nil)
}

// NewDetailed501 constructs a new Error with status code 501, not implemented.
func NewDetailed501(msg string, code ErrorCode, details interface{}) *Error {
	return &Error{msg, code, http.StatusNotImplemented, details}
}

// New501 constructs a new Error with status code 501, not implemented.
func New501(msg string, code ErrorCode) *Error {
	return NewDetailed501(msg, code, nil)
}

// NewDetailed404 constructs a new Error with status code 404, not found.
func NewDetailed404(msg string, code ErrorCode, details interface{}) *Error {
	return &Error{msg, code, http.StatusNotFound, details}
}

// New404 constructs a new Error with status code 404, not found.
func New404(msg string, code ErrorCode) *Error {
	return NewDetailed404(msg, code, nil)
}

// NewDetailed409 constructs a new Error with status code 409, status conflict.
func NewDetailed409(msg string, code ErrorCode, details interface{}) *Error {
	return &Error{msg, code, http.StatusConflict, details}
}

// New409 constructs a new Error with status code 409, status conflict.
func New409(msg string, code ErrorCode) *Error {
	return NewDetailed409(msg, code, nil)
}

// AuthorizationError creates a generic authorization Error with status code 403,
// forbidden.
func AuthorizationError() *Error {
	return &Error{
		"not authorized to make this request",
		UnknownUnauthorizedCode,
		http.StatusForbidden,
		nil,
	}
}

// AuthorizationMethodDeniedError creates a new Error indicating that the
// authorization method was denied, with status code 403, forbidden.
func AuthorizationMethodDeniedError() *Error {
	return &Error{
		"authorization method denied",
		AuthMethodDeniedCode,
		http.StatusForbidden,
		nil,
	}
}

// FromErrorDefaultType returns an *Error based on the provided error object. If err
// is already an *Error it is returned unchanged. Otherwise, a new *Error will be
// created with the specified code and status. The Message will be err.Error().
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

// FromError returns an *Error based on the provided error object. If err is already
// an *Error it is returned unchanged. Otherwise, a new *Error will be created with
// the status code 500, internal server error. The Message will be err.Error().
func FromError(err error, code ErrorCode) *Error {
	return FromErrorDefaultType(err, code, 500)
}
