package handler

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	httperr "github.com/turbinelabs/api/http/error"
)

// RichRequest provides easy access to the body, path args, and query args of
// an underlying net/http.Request
type RichRequest interface {
	// Examines a request's URL.Query() for a parameter matching the specified
	// name.  If the query arg has multiple specified values this will return
	// only the first value.
	//
	// Returns the value (or "") and a bool indicating if the variable was set.
	QueryArgOk(name string) (string, bool)

	// Returns the first value found for name (as QueryArgOk) or "" if no value
	// was found.
	QueryArg(name string) string

	// Returns the first value found in the URL (as QueryArgOk) or the defaultValue
	// if none was found.
	QueryArgOr(name, defaultValue string) string

	// Extract any body available in the request into a []byte. Returns a
	// service.http.error.Error if no body is available or if there are errors
	// reading the body.
	//
	// Important: Once read the body content will no longer be available.
	//
	// Important: This does not have any safety around the size of the memory
	// allocated for consumption of the request body.
	GetBody() ([]byte, error)

	// Consume the body and attempt to unmarshal it as JSON into an interface{}.
	// In addition to the error conditions of GetBody this will also return a
	// service.http.error.Error if there is an error returned from the unmarshal
	// process.
	//
	// Important: Once read the body content will no longer be available.
	//
	// Important: This does not have any safety around the size of the memory
	// allocated for consumption of the request body.
	GetBodyObject(resp interface{}) error

	// Access the underlying request that this is wrapping.
	Underlying() *http.Request
}

func NewRichRequest(request *http.Request) RichRequest {
	return &richRequest{request}
}

type richRequest struct {
	*http.Request
}

var _ RichRequest = &richRequest{}

func (rr *richRequest) Underlying() *http.Request {
	return rr.Request
}

func (rr *richRequest) QueryArgOk(name string) (string, bool) {
	q := rr.URL.Query()
	if v, ok := q[name]; ok && len(v) >= 1 {
		return v[0], true
	} else {
		return "", false
	}
}

func (rr *richRequest) QueryArg(name string) string {
	return rr.QueryArgOr(name, "")
}

func (rr *richRequest) QueryArgOr(name, defaultValue string) string {
	if r, ok := rr.QueryArgOk(name); ok {
		return r
	} else {
		return defaultValue
	}
}

func (rr *richRequest) GetBody() ([]byte, error) {
	if rr.Body == nil {
		return nil, httperr.New400("no body available", httperr.UnknownNoBodyCode)
	}

	b, err := ioutil.ReadAll(rr.Body)
	defer rr.Body.Close()
	if err != nil {
		return nil, httperr.New500("could not read request body", httperr.UnknownTransportCode)
	}
	return b, nil
}

func (rr *richRequest) GetBodyObject(resp interface{}) error {
	b, e := rr.GetBody()
	if e != nil {
		return e
	}

	e = json.Unmarshal(b, resp)
	if e != nil {
		return httperr.New400(
			"error handling JSON content: "+string(b),
			httperr.UnknownDecodingCode)
	}

	return nil
}
