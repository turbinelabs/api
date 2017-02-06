package flags

//go:generate mockgen -source $GOFILE -destination mock_$GOFILE -package $GOPACKAGE

import (
	"flag"

	tbnflag "github.com/turbinelabs/nonstdlib/flag"
)

// APIAuthKeyFromFlags represents command-line flags for specifying an
// API authentication key for the Turbine Labs API.
type APIAuthKeyFromFlags interface {
	// Returns the API authentication key from the command line.
	Make() string
}

// NewAPIAuthKeyFromFlags configures the necessary command line flags
// and returns an APIAuthKeyFromFlags.
func NewAPIAuthKeyFromFlags(flagset *flag.FlagSet) APIAuthKeyFromFlags {
	return NewPrefixedAPIAuthKeyFromFlags(prefixedFlagSet(flagset), true)
}

// NewPrefixedAPIAuthKeyFromFlags configures the necessary command
// line flags with a custom prefix and returns an APIAuthKeyFromFlags.
func NewPrefixedAPIAuthKeyFromFlags(
	flagset *tbnflag.PrefixedFlagSet,
	requiredFlag bool,
) APIAuthKeyFromFlags {
	ff := &apiAuthKeyFromFlags{}

	usage := "The auth key for {{NAME}} requests"
	if requiredFlag {
		usage = tbnflag.Required(usage)
	}

	flagset.StringVar(&ff.apiKey, "key", "", usage)

	return ff
}

type apiAuthKeyFromFlags struct {
	apiKey string
}

func (ff *apiAuthKeyFromFlags) Make() string {
	return ff.apiKey
}
