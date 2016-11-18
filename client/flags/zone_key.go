package flags

//go:generate mockgen -source $GOFILE -destination mock_$GOFILE -package $GOPACKAGE

import (
	"flag"
	"fmt"

	"github.com/turbinelabs/api"
	"github.com/turbinelabs/api/service"
	tbnflag "github.com/turbinelabs/nonstdlib/flag"
)

// ZoneKeyFromFlags represents command-line flags for specifying a
// Turbine Labs API zone name, which is used to resolves a zone key.
type ZoneKeyFromFlags interface {
	// Given a valid service.All instance, looks up the zone name
	// given on the command line and returns the corresponding
	// zone key or an error.
	Get(service.All) (api.ZoneKey, error)

	ZoneName() string
}

// NewZoneKeyFromFlags configures the necessary command line flags to
// retrieve a zone key by zone name.
func NewZoneKeyFromFlags(flagset *flag.FlagSet) ZoneKeyFromFlags {
	return NewZoneKeyFromFlagsWithPrefix(prefixedFlagSet(flagset))
}

// NewZoneKeyFromFlagsWithPrefix configures the necessary command line
// flags to retrieve a zone key by zone name with a custom command
// line flag prefix.
func NewZoneKeyFromFlagsWithPrefix(flagset *tbnflag.PrefixedFlagSet) ZoneKeyFromFlags {
	ff := &zoneKeyFromFlags{}

	flagset.StringVar(
		&ff.zoneName,
		"zone-name",
		"",
		tbnflag.Required("The name of the API Zone for {{NAME}} requests."),
	)

	return ff
}

type zoneKeyFromFlags struct {
	zoneName string
}

func (ff *zoneKeyFromFlags) Get(svc service.All) (api.ZoneKey, error) {
	zs, err := svc.Zone().Index(service.ZoneFilter{Name: ff.zoneName})
	if err != nil {
		return "", fmt.Errorf("error finding Zone with name %s: %s", ff.zoneName, err)
	}
	if len(zs) == 0 {
		return "", fmt.Errorf("Zone with name %s does not exist", ff.zoneName)
	}

	return zs[0].ZoneKey, nil

}

func (ff *zoneKeyFromFlags) ZoneName() string {
	return ff.zoneName
}