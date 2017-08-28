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

//go:generate $TBN_HOME/scripts/mockgen_internal.sh -type internalStatsClient -source $GOFILE -destination mock_$GOFILE -package $GOPACKAGE -aux_files statsapi=../service/stats/stats.go

import (
	"bytes"
	"compress/gzip"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/turbinelabs/api"
	apihttp "github.com/turbinelabs/api/http"
	httperr "github.com/turbinelabs/api/http/error"
	apiheader "github.com/turbinelabs/api/http/header"
	statsapi "github.com/turbinelabs/api/service/stats"
	"github.com/turbinelabs/nonstdlib/executor"
)

const (
	v1ForwardPath = "/v1.0/stats/forward"
	v2ForwardPath = "/v2.0/stats/forward"
	v1QueryPath   = "/v1.0/stats/query"

	queryArg = "query"
)

var (
	methodNotAllowedV1Err = errors.New("cannot use ForwardV2 on v1 client")
	methodNotAllowedV2Err = errors.New("cannot use Forward on v2 client")
)

// internalStatsClient is an internal interface for issuing forwarding
// requests with a callback
type internalStatsClient interface {
	statsapi.StatsService

	// Issues a forwarding request for the given payload with the
	// given executor.CallbackFunc.
	ForwardWithCallback(interface{}, executor.CallbackFunc) error
}

type httpStats struct {
	isV2           bool
	dest           apihttp.Endpoint
	forwardPath    string
	requestHandler apihttp.RequestHandler
	exec           executor.Executor
}

// NewStatsClient returns a blocking implementation of StatsService. Each
// invocation of Forward accepts a single Payload, issues a forwarding
// request to a remote stats-server and awaits a response.
func NewStatsClient(
	dest apihttp.Endpoint,
	apiKey string,
	clientApp App,
	exec executor.Executor,
) (statsapi.StatsService, error) {
	return newInternalStatsClient(dest, v1ForwardPath, apiKey, clientApp, exec)
}

// NewStatsV2Client returns a blocking implementation of StatsServiceV2. Each
// invocation of ForwardV2 accepts a single PayloadV2, issues a forwarding
// request to a remote stats-server and awaits a response.
func NewStatsV2Client(
	dest apihttp.Endpoint,
	apiKey string,
	clientApp App,
	exec executor.Executor,
) (statsapi.StatsServiceV2, error) {
	client, err := newInternalStatsClient(dest, v2ForwardPath, apiKey, clientApp, exec)
	if client != nil {
		client.isV2 = true
	}
	return client, err
}

func newInternalStatsClient(
	dest apihttp.Endpoint,
	forwardPath string,
	apiKey string,
	clientApp App,
	exec executor.Executor,
) (*httpStats, error) {
	// Copy the Endpoint to avoid polluting the original with our
	// headers.
	dest = dest.Copy()

	dest.AddHeader(apiheader.Authorization, apiKey)
	dest.AddHeader(apiheader.ClientType, clientType)
	dest.AddHeader(apiheader.ClientVersion, api.TbnPublicVersion)
	dest.AddHeader(apiheader.ClientApp, string(clientApp))

	// see encodePayload; payloads are sent as gzipped json
	dest.AddHeader("Content-Type", "application/json")
	dest.AddHeader("Content-Encoding", "gzip")

	return &httpStats{
		dest:           dest,
		forwardPath:    forwardPath,
		requestHandler: apihttp.NewRequestHandler(dest.Client()),
		exec:           exec,
	}, nil
}

func encodePayload(payload interface{}) ([]byte, error) {
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

func (hs *httpStats) ForwardWithCallback(
	payload interface{},
	cb executor.CallbackFunc,
) error {
	// TODO: remove type check and revert signature to struct pointer when v1 is removed.
	if hs.isV2 {
		if _, ok := payload.(*statsapi.PayloadV2); !ok {
			return fmt.Errorf("invalid payload type: %T, expected *PayloadV2", payload)
		}
	} else if _, ok := payload.(*statsapi.Payload); !ok {
		return fmt.Errorf("invalid payload type: %T, expected *Payload", payload)
	}

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
					req, err := hs.dest.NewRequest("POST", hs.forwardPath, apihttp.Params{}, rdr)
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

func (hs *httpStats) doForward(payload interface{}) (*statsapi.ForwardResult, error) {
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
	}
	return try.Get().(*statsapi.ForwardResult), nil
}

func (hs *httpStats) Forward(payload *statsapi.Payload) (*statsapi.ForwardResult, error) {
	if hs.isV2 {
		return nil, methodNotAllowedV2Err
	}
	return hs.doForward(payload)
}

func (hs *httpStats) ForwardV2(payload *statsapi.PayloadV2) (*statsapi.ForwardResult, error) {
	if !hs.isV2 {
		return nil, methodNotAllowedV1Err
	}
	return hs.doForward(payload)
}

func (hs *httpStats) Close() error {
	return nil
}

func (hs *httpStats) Query(query *statsapi.Query) (*statsapi.QueryResult, error) {
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
		return hs.dest.NewRequest(string(mGET), v1QueryPath, params, nil)
	}

	if err := hs.requestHandler.Do(reqFn, response); err != nil {
		return nil, err
	}

	return response, nil
}
