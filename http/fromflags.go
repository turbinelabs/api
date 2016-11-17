package http

//go:generate mockgen -source $GOFILE -destination mock_$GOFILE -package $GOPACKAGE

import (
	"crypto/tls"
	"net/http"

	tbnflag "github.com/turbinelabs/nonstdlib/flag"
)

// FromFlags constructs an http.Client and Endpoint from command line
// flags.
type FromFlags interface {
	// Makes an http.Client based on command line flags.
	MakeClient() *http.Client

	// Makes an Endpoint based on command line flags.
	MakeEndpoint() (Endpoint, error)
}

func NewFromFlags(defaultHost string, flagset *tbnflag.PrefixedFlagSet) FromFlags {
	ff := &fromFlags{}

	flagset.StringVar(&ff.host, "host", defaultHost, "The hostname for {{NAME}} requests")
	flagset.IntVar(&ff.port, "port", 443, "The port for {{NAME}} requests")
	flagset.BoolVar(&ff.ssl, "ssl", true, "If true, use SSL for {{NAME}} requests")
	flagset.BoolVar(
		&ff.insecure,
		"insecure",
		false,
		"If true, don't validate server cert when using SSL for {{NAME}} requests",
	)

	return ff
}

type fromFlags struct {
	host     string
	port     int
	ssl      bool
	insecure bool
}

func (ff *fromFlags) MakeClient() *http.Client {
	cl := HeaderPreservingClient()
	if ff.insecure {
		cl.Transport = &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
	}

	return cl
}

func (ff *fromFlags) MakeEndpoint() (Endpoint, error) {
	var protocol Protocol
	if ff.ssl {
		protocol = HTTPS
	} else {
		protocol = HTTP
	}

	return NewEndpoint(protocol, ff.host, ff.port)
}
