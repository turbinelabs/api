package handler

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

// Writes out an error / result pair to a waiting HTTP response. Parameter e
// will be forced into a httperr.Error. If it is already a httperr.Error or nil
// it will not be modified. If it is non-nil but also not already a
// httperr.Error it is encapsulated into an IntervalServerError with the error
// code 'UnknownUnclassifiedCode'.
//
// If the processed e is non-nil the HTTP response code is drawn from tha error
// otherwise it is 200.
//
// The actual body of the response is the json marshaled envelope.Response object
// containing the processed e and the result.
//
// Returns the HTTP result code that was sent.
func (rrw RichResponseWriter) WriteEnvelope(e error, result interface{}) int {
	// first ensure we're dealing with a known error type
	err := httperr.FromError(e, httperr.UnknownUnclassifiedCode)
	env := envelope.Response{err, result}
	resultCode := http.StatusOK

	if err != nil {
		resultCode = err.Status
	}

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
