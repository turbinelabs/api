package flags

import (
	"errors"
	"flag"
	"net/http"
	"testing"

	"github.com/golang/mock/gomock"

	clienthttp "github.com/turbinelabs/client/http"
	"github.com/turbinelabs/test/assert"
)

var (
	fakeClient      = &http.Client{}
	fakeEndpoint, _ = clienthttp.NewEndpoint(clienthttp.HTTPS, "localhost", 1234)
)

func TestNewClientFromFlags(t *testing.T) {
	flagset := flag.NewFlagSet("NewClientFromFlags options", flag.PanicOnError)

	ff := NewClientFromFlags(flagset)
	ffImpl := ff.(*clientFromFlags)
	assert.NonNil(t, ffImpl.apiConfigFromFlags)

	assert.NonNil(t, flagset.Lookup("api.key"))
}

func TestNewClientFromFlagsWithSharedAPIKey(t *testing.T) {
	flagset := flag.NewFlagSet("NewClientFromFlags options", flag.PanicOnError)

	apiConfigFromFlags := NewAPIConfigFromFlags(flagset)
	assert.NonNil(t, flagset.Lookup("api.key"))

	ff := NewClientFromFlagsWithSharedAPIConfig(flagset, apiConfigFromFlags)
	ffImpl := ff.(*clientFromFlags)
	assert.NonNil(t, ffImpl.apiConfigFromFlags)
	assert.SameInstance(t, ffImpl.apiConfigFromFlags, apiConfigFromFlags)
}

func TestClientFromFlagsMake(t *testing.T) {
	ctrl := gomock.NewController(assert.Tracing(t))
	defer ctrl.Finish()

	mockApiConfig := NewMockAPIConfigFromFlags(ctrl)

	mockApiConfig.EXPECT().APIKey().Return("api-key")
	mockApiConfig.EXPECT().MakeClient().Return(fakeClient)
	mockApiConfig.EXPECT().MakeEndpoint().Return(fakeEndpoint, nil)

	ff := &clientFromFlags{mockApiConfig}

	svc, err := ff.Make()
	assert.Nil(t, err)
	assert.NonNil(t, svc)
}

func TestClientFromFlagsMakeError(t *testing.T) {
	ctrl := gomock.NewController(assert.Tracing(t))
	defer ctrl.Finish()

	mockApiConfig := NewMockAPIConfigFromFlags(ctrl)

	mockApiConfig.EXPECT().MakeEndpoint().Return(clienthttp.Endpoint{}, errors.New("nope"))

	ff := &clientFromFlags{mockApiConfig}

	svc, err := ff.Make()
	assert.ErrorContains(t, err, "nope")
	assert.Nil(t, svc)
}
