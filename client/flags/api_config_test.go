package flags

import (
	"flag"
	"testing"

	tbnflag "github.com/turbinelabs/nonstdlib/flag"
	"github.com/turbinelabs/test/assert"
)

func TestNewAPIConfigFromFlags(t *testing.T) {
	flagset := flag.NewFlagSet("NewAPIAuthFromFlags options", flag.PanicOnError)

	ff := NewAPIConfigFromFlags(flagset)
	ffImpl := ff.(*apiConfigFromFlags)

	flagset.Parse([]string{"-api.key=schlage"})

	assert.Equal(t, ffImpl.apiKeyConfig.Make(), "schlage")

	theFlag := flagset.Lookup("api.key")
	assert.NonNil(t, theFlag)
	assert.True(t, tbnflag.IsRequired(theFlag))

	assert.NonNil(t, ffImpl.FromFlags)
}

func TestNewPrefixedAPIConfigFromFlags(t *testing.T) {
	flagset := flag.NewFlagSet("NewAPIAuthFromFlags options", flag.PanicOnError)
	prefixedFlagset := tbnflag.NewPrefixedFlagSet(flagset, "test", "test")

	ff := NewPrefixedAPIConfigFromFlags(
		prefixedFlagset,
		APIConfigSetAPIAuthKeyFromFlags(NewPrefixedAPIAuthKeyFromFlags(prefixedFlagset, false)),
	)
	ffImpl := ff.(*apiConfigFromFlags)

	flagset.Parse([]string{"-test.key=schlage"})

	assert.Equal(t, ffImpl.apiKeyConfig.Make(), "schlage")

	theFlag := flagset.Lookup("test.key")
	assert.NonNil(t, theFlag)
	assert.False(t, tbnflag.IsRequired(theFlag))

	assert.NonNil(t, ffImpl.FromFlags)
}

func TestAPIConfigFromFlagsGet(t *testing.T) {
	ff := &apiConfigFromFlags{apiKeyConfig: &apiAuthKeyFromFlags{"schlage"}}
	assert.Equal(t, ff.APIKey(), "schlage")
}
