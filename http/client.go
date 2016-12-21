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
	"errors"
	"net/http"
)

var redirectOverflow = errors.New("Stopped after 5 redirects")

// HeaderPreserving produces an http.Client with CheckRedirect set to:
//
// 1) Pass headers from the initial request to the new request
// 2) Return an error if 5 redirects fail to result in a non 3xx response

func HeaderPreservingClient() *http.Client {
	return &http.Client{CheckRedirect: redirectPolicy}
}

func redirectPolicy(next *http.Request, prev []*http.Request) error {
	if len(prev) > 5 {
		return redirectOverflow
	}

	next.Header = prev[0].Header
	return nil
}
