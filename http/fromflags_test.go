package http

import (
	"flag"
	"net/http"
	"testing"

	tbnflag "github.com/turbinelabs/nonstdlib/flag"
	"github.com/turbinelabs/test/assert"
)

func TestNewFromFlags(t *testing.T) {
	flagset := tbnflag.NewPrefixedFlagSet(
		flag.NewFlagSet("api/http options", flag.PanicOnError),
		"api",
		"API",
	)

	ff := NewFromFlags("api.turbinelabs.io", flagset)
	ffImpl := ff.(*fromFlags)

	assert.Equal(t, ffImpl.host, "api.turbinelabs.io")

	flagset.Parse([]string{
		"-api.host=example.com",
		"-api.port=999",
		"-api.ssl=false",
		"-api.insecure=true",
	})

	assert.Equal(t, ffImpl.host, "example.com")
	assert.Equal(t, ffImpl.port, 999)
	assert.False(t, ffImpl.ssl)
	assert.True(t, ffImpl.insecure)
}

func TestFromFlagsMakeClient(t *testing.T) {
	ff := &fromFlags{}
	client := ff.MakeClient()
	assert.Nil(t, client.Transport)
}

func TestFromFlagsMakeClientInsecure(t *testing.T) {
	ff := &fromFlags{insecure: true}
	client := ff.MakeClient()

	switch transport := client.Transport.(type) {
	case *http.Transport:
		assert.True(t, transport.TLSClientConfig.InsecureSkipVerify)

	default:
		t.Fatal("bad type")
	}
}

func TestFromFlagsMakeEndpoint(t *testing.T) {
	ff := &fromFlags{
		host: "example.com",
		port: 80,
	}

	e, err := ff.MakeEndpoint()
	assert.Nil(t, err)
	assert.Equal(t, e.urlBase.String(), "http://example.com:80")

	ff.port = 443
	ff.ssl = true

	e, err = ff.MakeEndpoint()
	assert.Nil(t, err)
	assert.Equal(t, e.urlBase.String(), "https://example.com:443")

	ff.host = "not a domain"

	_, err = ff.MakeEndpoint()
	assert.NonNil(t, err)
}
