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

//go:generate mockgen -source $GOFILE -destination mock_$GOFILE -package $GOPACKAGE --write_package_comment=false

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"strings"

	tbnflag "github.com/turbinelabs/nonstdlib/flag"
	tbnstrings "github.com/turbinelabs/nonstdlib/strings"
)

// FromFlags constructs an Endpoint from command line flags.
type FromFlags interface {
	// Validates the command line flags for an Endpoint.
	Validate() error

	// Makes an Endpoint based on command line flags.
	MakeEndpoint() (Endpoint, error)
}

func NewFromFlags(defaultHostPort string, flagset tbnflag.FlagSet) FromFlags {
	ff := &fromFlags{
		headers: tbnflag.NewStrings(),
	}

	flagset.HostPortVar(
		&ff.hostPort,
		"host",
		ff.initialHostPort(defaultHostPort),
		"The address (`host:port`) for {{NAME}} requests. If no port is given, it defaults to port 443 if --{{PREFIX}}ssl is true and port 80 otherwise.",
	)
	flagset.BoolVar(&ff.ssl, "ssl", true, "If true, use SSL for {{NAME}} requests")
	flagset.BoolVar(
		&ff.insecure,
		"insecure",
		false,
		"If true, don't validate server cert when using SSL for {{NAME}} requests",
	)

	flagset.Var(
		&ff.headers,
		"header",
		"Specifies a custom `header` to send with every {{NAME}} request. Headers are given as name:value pairs. Leading and trailing whitespace will be stripped from the name and value. For multiple headers, this flag may be repeated or multiple headers can be delimited with commas.",
	)

	return ff
}

type header string

func (h header) split() (string, string, error) {
	name, value := tbnstrings.SplitFirstColon(string(h))
	if name == "" || value == "" {
		return "", "", fmt.Errorf("invalid header: %s", string(h))
	}

	return strings.TrimSpace(name), strings.TrimSpace(value), nil
}

type fromFlags struct {
	hostPort tbnflag.HostPort
	ssl      bool
	insecure bool
	headers  tbnflag.Strings
}

func (ff *fromFlags) makeClient() *http.Client {
	cl := HeaderPreservingClient()
	if ff.insecure {
		cl.Transport = &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
	}

	return cl
}

func (ff *fromFlags) Validate() error {
	for _, hs := range ff.headers.Strings {
		if _, _, err := header(hs).split(); err != nil {
			return err
		}
	}

	return nil
}

func (ff *fromFlags) MakeEndpoint() (Endpoint, error) {
	var protocol Protocol
	if ff.ssl {
		protocol = HTTPS
	} else {
		protocol = HTTP
	}

	e, err := NewEndpoint(protocol, ff.hostPort.Addr())
	if err != nil {
		return e, err
	}

	e.SetClient(ff.makeClient())

	for _, hs := range ff.headers.Strings {
		hdr, value, err := header(hs).split()
		if err != nil {
			return e, err
		}

		e.AddHeader(hdr, value)
	}

	return e, nil
}

func (ff *fromFlags) initialHostPort(defaultHostPort string) tbnflag.HostPort {
	defaultPortFunc := func() int {
		if ff.ssl {
			return 443
		}
		return 80

	}
	return tbnflag.NewHostPortWithDefaultPort(defaultHostPort, defaultPortFunc)
}
