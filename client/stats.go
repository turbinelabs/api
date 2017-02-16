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

package client

//go:generate mockgen -source $GOFILE -destination mock_$GOFILE -package $GOPACKAGE -aux_files statsapi=../service/stats/stats.go

import (
	"bytes"
	"compress/gzip"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	apihttp "github.com/turbinelabs/api/http"
	httperr "github.com/turbinelabs/api/http/error"
	apiheader "github.com/turbinelabs/api/http/header"
	statsapi "github.com/turbinelabs/api/service/stats"
	"github.com/turbinelabs/nonstdlib/executor"
)

const (
	statsClientID string = "tbn-stats-client (v0.1)"

	forwardPath = "/v1.0/stats/forward"
	queryPath   = "/v1.0/stats/query"
	queryArg    = "query"
)

// internalStatsClient is an internal interface for issuing forwarding
// requests with a callback
type internalStatsClient interface {
	statsapi.StatsService
	// Issues a forwarding request for the given payload with the
	// given executor.CallbackFunc.
	ForwardWithCallback(*statsapi.Payload, executor.CallbackFunc) error
}

type httpStatsV1 struct {
	dest           apihttp.Endpoint
	requestHandler apihttp.RequestHandler
	exec           executor.Executor
}

// NewStatsClient returns a blocking implementation of Stats. Each
// invocation of Forward accepts a single Payload, issues a forwarding
// request to a remote stats-server and awaits a response.
func NewStatsClient(
	dest apihttp.Endpoint,
	apiKey string,
	exec executor.Executor,
) (statsapi.StatsService, error) {
	return newInternalStatsClient(dest, apiKey, exec)
}

func newInternalStatsClient(
	dest apihttp.Endpoint,
	apiKey string,
	exec executor.Executor,
) (internalStatsClient, error) {
	// Copy the Endpoint to avoid polluting the original with our
	// headers.
	dest = dest.Copy()

	dest.AddHeader(apiheader.Authorization, apiKey)
	dest.AddHeader(apiheader.ClientID, statsClientID)

	// see encodePayload; payloads are sent as gzipped json
	dest.AddHeader("Content-Type", "application/json")
	dest.AddHeader("Content-Encoding", "gzip")

	return &httpStatsV1{
		dest,
		apihttp.NewRequestHandler(dest.Client()),
		exec,
	}, nil
}

func encodePayload(payload *statsapi.Payload) ([]byte, error) {
	var buffer bytes.Buffer
	gzip := gzip.NewWriter(&buffer)
	encoder := json.NewEncoder(gzip)

	if err := encoder.Encode(payload); err != nil {
		msg := fmt.Sprintf("could not encode stats payload: %+v\n%+v", err, payload)
		return nil, httperr.New400(msg, httperr.UnknownEncodingCode)
	}

	if err := gzip.Close(); err != nil {
		msg := fmt.Sprintf(
			"could not finish encoding stats payload: %+v\n%+v",
			err,
			payload,
		)
		return nil, httperr.New400(msg, httperr.UnknownEncodingCode)
	}

	return buffer.Bytes(), nil
}

func (hs *httpStatsV1) ForwardWithCallback(
	payload *statsapi.Payload,
	cb executor.CallbackFunc,
) error {
	encoded, err := encodePayload(payload)
	if err != nil {
		return err
	}

	hs.exec.Exec(
		func(ctxt context.Context) (interface{}, error) {
			response := &statsapi.ForwardResult{}
			if err := hs.requestHandler.Do(
				func() (*http.Request, error) {
					rdr := bytes.NewReader(encoded)
					req, err := hs.dest.NewRequest(
						"POST",
						forwardPath,
						apihttp.Params{},
						rdr,
					)
					if err != nil {
						return nil, err
					}
					return req.WithContext(ctxt), nil
				},
				response,
			); err != nil {
				return nil, err
			}

			return response, nil
		},
		cb,
	)
	return nil
}

func (hs *httpStatsV1) Forward(payload *statsapi.Payload) (*statsapi.ForwardResult, error) {
	responseChan := make(chan executor.Try, 1)
	defer close(responseChan)

	err := hs.ForwardWithCallback(
		payload,
		func(try executor.Try) {
			responseChan <- try
		},
	)
	if err != nil {
		return nil, err
	}

	try := <-responseChan
	if try.IsError() {
		return nil, try.Error()
	} else {
		return try.Get().(*statsapi.ForwardResult), nil
	}
}

func (hs *httpStatsV1) Close() error {
	return nil
}

func (hs *httpStatsV1) Query(query *statsapi.Query) (*statsapi.QueryResult, error) {
	params := apihttp.Params{}

	if query != nil {
		queryBytes, err := json.Marshal(query)
		if err != nil {
			return nil, httperr.New400(
				fmt.Sprintf("unable to encode query: %v: %s", query, err),
				httperr.UnknownUnclassifiedCode,
			)
		}

		params[queryArg] = string(queryBytes)
	}

	response := &statsapi.QueryResult{}
	reqFn := func() (*http.Request, error) {
		return hs.dest.NewRequest(string(mGET), queryPath, params, nil)
	}

	if err := hs.requestHandler.Do(reqFn, response); err != nil {
		return nil, err
	}

	return response, nil
}
