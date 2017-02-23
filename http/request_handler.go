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

package http

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/turbinelabs/api/http/envelope"
	httperr "github.com/turbinelabs/api/http/error"
)

type RequestHandler struct {
	client *http.Client
}

func NewRequestHandler(client *http.Client) RequestHandler {
	return RequestHandler{client}
}

func getBody(response *http.Response) ([]byte, *httperr.Error) {
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, httperr.New400(
			"unable to process server response: "+err.Error(),
			httperr.UnknownTransportCode)
	}

	return body, nil
}

func expectsNoPayload(response *http.Response) error {
	if response.StatusCode == http.StatusOK {
		return nil
	}

	if response.Body == nil {
		err := httperr.Error{
			"error response with no additional information",
			httperr.UnknownNoBodyCode,
			response.StatusCode,
			nil}
		return &err
	}

	bodyBytes, err := getBody(response)
	if err != nil {
		return err
	}

	env := envelope.Response{}
	unmarshalErr := json.Unmarshal(bodyBytes, &env)
	if unmarshalErr != nil {
		return mkUnmarshalErr(bodyBytes, unmarshalErr)
	}

	if env.Error != nil {
		env.Error.Status = response.StatusCode
	}

	return env.Error
}

func expectsPayload(response *http.Response, payloadDest interface{}) error {
	if response.Body == nil {
		return httperr.New500(
			fmt.Sprintf(
				"expected payload but response (%d) included no content",
				response.StatusCode),
			httperr.UnknownNoBodyCode)
	}

	bodyBytes, err := getBody(response)
	if err != nil {
		return err
	}

	var rawResponse json.RawMessage
	env := envelope.Response{Payload: &rawResponse}
	unmarshalErr := json.Unmarshal(bodyBytes, &env)
	if unmarshalErr != nil {
		return mkUnmarshalErr(bodyBytes, unmarshalErr)
	}

	if env.Error != nil {
		env.Error.Status = response.StatusCode
		return env.Error
	}

	unmarshalErr = json.Unmarshal(rawResponse, payloadDest)
	if unmarshalErr != nil {
		return mkUnmarshalErr([]byte(rawResponse), unmarshalErr)
	}

	return nil
}

// Given a request and response container make the request and populate the
// response object. If the server returns an error (an encoded service.error)
// or there are problems decoding the response return an error.
func (rh RequestHandler) Do(
	mkReq func() (*http.Request, error),
	response interface{},
) error {
	req, err := mkReq()
	if err != nil {
		return fmt.Errorf("could not create request: %s", err.Error())
	}

	// make HTTP request
	resp, err := rh.client.Do(req)

	// if there was a problem with actually making the request bail indicating
	// something was wrong with the server (this is, admittedly, a guess without
	// further introspection but we'll let it stand for now).
	if err != nil {
		return fmt.Errorf("could not successfully make request: %s", err.Error())
	}

	if response == nil {
		return expectsNoPayload(resp)
	}

	return expectsPayload(resp, response)
}

func mkUnmarshalErr(content []byte, underlying error) *httperr.Error {
	return httperr.New500(
		fmt.Sprintf(
			"got malformed response; unmarshal error: '%s' - content: '%s'",
			underlying.Error(),
			string(content)),
		httperr.UnknownDecodingCode)
}
