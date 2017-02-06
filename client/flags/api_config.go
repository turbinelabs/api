package flags

//go:generate mockgen -source $GOFILE -destination mock_$GOFILE -package $GOPACKAGE -aux_files "clienthttp=$TBN_HOME/client/http/fromflags.go" -imports http=net/http

import (
	"flag"

	clienthttp "github.com/turbinelabs/client/http"
	tbnflag "github.com/turbinelabs/nonstdlib/flag"
)

// APIConfigFromFlags represents command-line flags for specifying an
// API authentication key, host, port and SSL settings for the Turbine
// Labs API.
type APIConfigFromFlags interface {
	clienthttp.FromFlags

	// APIKey Returns the API authentication key from the command line.
	// Equivalent to calling APIAuthKeyFromFlags().Make()
	APIKey() string

	// APIAuthKeyFromFlags returns the underlying APIAuthKeyFromFlags
	// so that it can potentially be shared between APIConfigFromFlags
	// via the APIConfigSetAPIAuthKeyFromFlags APIConfigOption.
	APIAuthKeyFromFlags() APIAuthKeyFromFlags
}

// NewAPIConfigFromFlags configures the necessary command line flags
// and returns an APIConfigFromFlags.
func NewAPIConfigFromFlags(flagset *flag.FlagSet) APIConfigFromFlags {
	return NewPrefixedAPIConfigFromFlags(prefixedFlagSet(flagset))
}

type APIConfigOption func(*apiConfigFromFlags)

// APIConfigSetAPIAuthKeyFromFlags allows the caller to specify a shared
// APIAuthKeyFromFlags, likely obtained via the
// APIConfigFromFlags.APIAuthKeyFromFlags() method.
func APIConfigSetAPIAuthKeyFromFlags(akff APIAuthKeyFromFlags) APIConfigOption {
	return func(ff *apiConfigFromFlags) {
		ff.apiKeyConfig = akff
	}
}

// NewPrefixedAPIConfigFromFlags configures the necessary command
// line flags with a custom prefix and returns an APIConfigFromFlags.
func NewPrefixedAPIConfigFromFlags(
	flagset *tbnflag.PrefixedFlagSet,
	opts ...APIConfigOption,
) APIConfigFromFlags {
	ff := &apiConfigFromFlags{requiredFlag: true}

	for _, applyOpt := range opts {
		applyOpt(ff)
	}

	if ff.apiKeyConfig == nil {
		ff.apiKeyConfig = NewPrefixedAPIAuthKeyFromFlags(flagset, ff.requiredFlag)
	}

	ff.FromFlags = clienthttp.NewFromFlags("api.turbinelabs.io", flagset)

	return ff
}

type apiConfigFromFlags struct {
	clienthttp.FromFlags
	apiKeyConfig APIAuthKeyFromFlags
	requiredFlag bool
}

func (ff *apiConfigFromFlags) APIKey() string {
	return ff.apiKeyConfig.Make()
}

func (ff *apiConfigFromFlags) APIAuthKeyFromFlags() APIAuthKeyFromFlags {
	return ff.apiKeyConfig
}
