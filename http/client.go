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
