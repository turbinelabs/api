package flags

//go:generate mockgen -source $GOFILE -destination mock_$GOFILE -package $GOPACKAGE

import (
	"flag"

	"github.com/turbinelabs/api/client"
	"github.com/turbinelabs/api/service"
)

// ServiceFromFlags represents command-line flags specifying
// configuration of a service.All backed by the Turbine Labs API.
type ClientFromFlags interface {
	// Make produces a service.All from the provided flags, or an
	// error.
	Make() (service.All, error)
}

// NewServiceFromFlags creates a ServiceFromFlags, which configures
// the necessary flags to construct a service.All instance.
func NewClientFromFlags(flagset *flag.FlagSet) ClientFromFlags {
	return NewClientFromFlagsWithSharedAPIConfig(flagset, nil)
}

// NewServiceFromFlagsWithSharedAPIConfig creates a ServiceFromFlags,
// which configures the necessary flags to construct a service.All
// instance. The given APIConfigFromFlags is used to obtain the API
// auth key.
func NewClientFromFlagsWithSharedAPIConfig(
	flagset *flag.FlagSet,
	apiConfigFromFlags APIConfigFromFlags,
) ClientFromFlags {
	ff := &clientFromFlags{}

	if apiConfigFromFlags == nil {
		ff.apiConfigFromFlags = NewAPIConfigFromFlags(flagset)
	} else {
		ff.apiConfigFromFlags = apiConfigFromFlags
	}

	return ff
}

type clientFromFlags struct {
	apiConfigFromFlags APIConfigFromFlags
}

func (ff *clientFromFlags) Make() (service.All, error) {
	endpoint, err := ff.apiConfigFromFlags.MakeEndpoint()
	if err != nil {
		return nil, err
	}

	apiKey := ff.apiConfigFromFlags.APIKey()
	c := ff.apiConfigFromFlags.MakeClient()

	return client.NewAll(endpoint, apiKey, c)
}
