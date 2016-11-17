package http

import (
	"errors"
	"net/http"
	"testing"

	"github.com/turbinelabs/api/http/envelope"
	httperr "github.com/turbinelabs/api/http/error"
	httptest "github.com/turbinelabs/test/http"
)

func getRRWTestWriter(t *testing.T) (RichResponseWriter, *httptest.ResponseRecorder) {
	recorder := httptest.NewResponseRecorder(t)
	return RichResponseWriter{recorder}, recorder
}

type testRRWStruct struct {
	Field1 string `json:"field1"`
	Field2 []int  `json:"field2"`
}

type poisonRRWStruct struct{}

func (_ poisonRRWStruct) MarshalJSON() ([]byte, error) {
	return nil, errors.New("w\"at")
}

func (_ poisonRRWStruct) Error() string {
	return "whelp"
}

func TestWriteEnvelopeSimple(t *testing.T) {
	rrw, rec := getRRWTestWriter(t)
	s := testRRWStruct{"aoeu", []int{1, 2, 3, 4}}
	rrw.WriteEnvelope(nil, s)
	rec.AssertBodyJSON(envelope.Response{nil, s})
	rec.AssertStatus(http.StatusOK)
	rec.AssertHeader("content-type", "application/json")
}

func TestWriteEnvelopeNoContent(t *testing.T) {
	rrw, rec := getRRWTestWriter(t)
	var foo *int
	rrw.WriteEnvelope(nil, foo)

	// we do this to ensure that a non-'interface{} nil' payload works as expected
	wantBody := `{"result":null}`
	rec.AssertBody(wantBody)
	rec.AssertStatus(http.StatusOK)
	rec.AssertHeader("content-type", "application/json")
}

func TestWriteEnvelopeBadResult(t *testing.T) {
	rrw, rec := getRRWTestWriter(t)
	rrw.WriteEnvelope(nil, poisonRRWStruct{})
	wantBody := `{"error": {"message":"failed to encode response object: '{Error:<nil> Payload:whelp}'; error was: 'json: error calling MarshalJSON for type http.poisonRRWStruct: w\"at'","code":"UnknownEncodingCode"}}`
	rec.AssertBody(wantBody)
	rec.AssertStatus(http.StatusInternalServerError)
	rec.AssertHeader("content-type", "application/json")
}

func TestWriteEnvelopeInferHTTPErrStatusCode(t *testing.T) {
	rrw, rec := getRRWTestWriter(t)
	err := httperr.New400("some stuff", httperr.UnknownTransportCode)
	s := testRRWStruct{"aosentuh", []int{2, 2, 2}}
	rrw.WriteEnvelope(err, s)
	rec.AssertBodyJSON(envelope.Response{err, s})
	rec.AssertStatus(400)
	rec.AssertHeader("content-type", "application/json")
}

func TestWriteEnvelopeLiftToHTTPErr(t *testing.T) {
	rrw, rec := getRRWTestWriter(t)
	s := testRRWStruct{"asonetuh", nil}
	err := errors.New("whee")
	rrw.WriteEnvelope(err, s)
	wantErr := httperr.New500(err.Error(), httperr.UnknownUnclassifiedCode)
	rec.AssertBodyJSON(envelope.Response{wantErr, s})
	rec.AssertStatus(http.StatusInternalServerError)
	rec.AssertHeader("content-type", "application/json")
}
