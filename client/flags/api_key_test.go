package flags

import (
	"flag"
	"testing"

	tbnflag "github.com/turbinelabs/nonstdlib/flag"
	"github.com/turbinelabs/test/assert"
)

func TestNewAPIAuthKeyFromFlags(t *testing.T) {
	flagset := flag.NewFlagSet("NewAPIAuthFromFlags options", flag.PanicOnError)

	ff := NewAPIAuthKeyFromFlags(flagset)
	ffImpl := ff.(*apiAuthKeyFromFlags)

	flagset.Parse([]string{"-api.key=schlage"})

	assert.Equal(t, ffImpl.apiKey, "schlage")

	theFlag := flagset.Lookup("api.key")
	assert.NonNil(t, theFlag)
	assert.True(t, tbnflag.IsRequired(theFlag))
}

func TestNewPrefixedAPIAuthKeyFromFlags(t *testing.T) {
	flagset := flag.NewFlagSet("NewAPIAuthFromFlags options", flag.PanicOnError)
	prefixedFlagset := tbnflag.NewPrefixedFlagSet(flagset, "test", "test")

	ff := NewPrefixedAPIAuthKeyFromFlags(prefixedFlagset, false)
	ffImpl := ff.(*apiAuthKeyFromFlags)

	flagset.Parse([]string{"-test.key=schlage"})

	assert.Equal(t, ffImpl.apiKey, "schlage")

	theFlag := flagset.Lookup("test.key")
	assert.NonNil(t, theFlag)
	assert.False(t, tbnflag.IsRequired(theFlag))
}

func TestAPIAuthKeyFromFlagsGet(t *testing.T) {
	ff := &apiAuthKeyFromFlags{apiKey: "schlage"}
	assert.Equal(t, ff.Make(), "schlage")
}
