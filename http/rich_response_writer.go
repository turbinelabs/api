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

package http

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/turbinelabs/api/http/envelope"
	httperr "github.com/turbinelabs/api/http/error"
)

type RichResponseWriter struct {
	http.ResponseWriter
}

// this is used if encoding a response fails
func mkLastDitchErrorEnvelope(
	env envelope.Response,
	marshalErr error,
) []byte {
	q := func(i interface{}) string {
		q := fmt.Sprintf("%q", fmt.Sprintf("%+v", i))
		if len(q) < 2 {
			q = ""
		} else {
			q = q[1 : len(q)-1]
		}
		return q
	}

	return []byte(
		fmt.Sprintf(
			`{"error": `+
				`{"message":"failed to encode response object: '%s'; error was: '%s'","code":"%s"}}`,
			q(env),
			q(marshalErr),
			httperr.UnknownEncodingCode))
}

// HasDetails can be implemented by objects passed to WriteEnvelope in order to
// pass some data back in parallel to the result of a called endpoint. For
// example, it can be used to pass pagination context without altering the payload.
type HasDetails interface {
	GetPayload() interface{}
	GetDetails() interface{}
}

// Writes out an error / result pair to a waiting HTTP response. Parameter e
// will be forced into a httperr.Error. If it is already a httperr.Error or nil
// it will not be modified. If it is non-nil but also not already a
// httperr.Error it is encapsulated into an IntervalServerError with the error
// code 'UnknownUnclassifiedCode'.
//
// If the processed e is non-nil the HTTP response code is drawn from the error
// otherwise it is 200.
//
// The actual body of the response is the json marshaled envelope.Response object
// containing the processed e and the result.
//
// Returns the HTTP result code that was sent.
func (rrw RichResponseWriter) WriteEnvelope(e error, result interface{}) int {
	// first ensure we're dealing with a known error type
	err := httperr.FromError(e, httperr.UnknownUnclassifiedCode)
	// and set the result code appropriately
	resultCode := http.StatusOK
	if err != nil {
		resultCode = err.Status
	}

	var (
		payload interface{}
		details interface{}
	)

	if hd, ok := result.(HasDetails); ok {
		payload = hd.GetPayload()
		details = hd.GetDetails()
	} else {
		payload = result
	}

	env := envelope.Response{err, payload, details}

	resultBytes, marshalErr := json.Marshal(env)
	if marshalErr != nil {
		resultBytes = mkLastDitchErrorEnvelope(env, marshalErr)
		resultCode = http.StatusInternalServerError
	}
	rrw.Header().Add("content-type", "application/json")
	rrw.WriteHeader(resultCode)
	rrw.Write(resultBytes)

	return resultCode
}
