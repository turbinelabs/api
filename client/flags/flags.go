package flags

import (
	"flag"

	tbnflag "github.com/turbinelabs/nonstdlib/flag"
)

const (
	flagPrefix = "api"
	moduleName = "API"
)

func prefixedFlagSet(flagset *flag.FlagSet) *tbnflag.PrefixedFlagSet {
	return tbnflag.NewPrefixedFlagSet(flagset, flagPrefix, moduleName)
}
