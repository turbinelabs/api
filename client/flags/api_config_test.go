/*
Copyright 2017 Turbine Labs, Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package flags

import (
	"io/ioutil"
	"net/http"
	"os"
	"testing"

	"github.com/golang/mock/gomock"
	"golang.org/x/oauth2"

	"github.com/turbinelabs/api/client/tokencache"
	apihttp "github.com/turbinelabs/api/http"
	tbnflag "github.com/turbinelabs/nonstdlib/flag"
	"github.com/turbinelabs/nonstdlib/flag/usage"
	"github.com/turbinelabs/test/assert"
)

func TestNewAPIConfigFromFlags(t *testing.T) {
	flagset := tbnflag.NewTestFlagSet()

	ff := NewAPIConfigFromFlags(flagset)
	ffImpl := ff.(*apiConfigFromFlags)

	flagset.Parse([]string{"-key=schlage"})

	assert.Equal(t, ffImpl.apiKeyConfig.Make(), "schlage")

	theFlag := flagset.Unwrap().Lookup("key")
	assert.NonNil(t, theFlag)
	assert.True(t, usage.IsSensitive(theFlag))
	assert.True(t, usage.IsRequired(theFlag))

	assert.NonNil(t, ffImpl.clientFromFlags)
}

func TestNewAPIConfigFromFlagsWithPrefix(t *testing.T) {
	flagset := tbnflag.NewTestFlagSet()
	apiScopedFlagSet := flagset.Scope("api", "test")
	ff := NewAPIConfigFromFlags(
		apiScopedFlagSet,
		APIConfigSetAPIAuthKeyFromFlags(
			NewAPIAuthKeyFromFlags(apiScopedFlagSet, APIAuthKeyFlagsOptional()),
		),
	)
	ffImpl := ff.(*apiConfigFromFlags)

	flagset.Parse([]string{"-api.key=schlage"})

	assert.Equal(t, ffImpl.apiKeyConfig.Make(), "schlage")

	theFlag := flagset.Unwrap().Lookup("api.key")
	assert.NonNil(t, theFlag)
	assert.False(t, usage.IsRequired(theFlag))

	assert.NonNil(t, ffImpl.clientFromFlags)
}

func TestAPIConfigFromFlagsGet(t *testing.T) {
	ff := &apiConfigFromFlags{
		apiKeyConfig: &apiAuthKeyFromFlags{optional: false, apiKey: "schlage"},
	}
	assert.Equal(t, ff.APIKey(), "schlage")
}

func TestNewAPIConfigFromFlagsWithAuthTokenImpactOnKeyFlag(t *testing.T) {
	fset := tbnflag.NewTestFlagSet()
	NewAPIConfigFromFlags(fset, APIConfigMayUseAuthToken(tokencache.NewStaticPath("nonsense")))

	keyFlag := fset.Unwrap().Lookup("key")
	assert.NonNil(t, keyFlag)
	assert.False(t, usage.IsRequired(keyFlag))
	assert.True(t, usage.IsSensitive(keyFlag))
}

func TestNewAPIConfigFromFlagsValidateWithAuthToken(t *testing.T) {
	fset := tbnflag.NewTestFlagSet()
	ff := NewAPIConfigFromFlags(fset, APIConfigMayUseAuthToken(tokencache.NewStaticPath("")))
	fset.Parse(nil)
	assert.Nil(t, ff.Validate())
}

func TestNewAPIConfigFromFlagsValidateWithAuthTokenExpired(t *testing.T) {
	file, err := ioutil.TempFile("", "token-cache")
	assert.Nil(t, err)
	path := file.Name()
	defer func() { os.Remove(path) }()

	tcBytes := []byte(`{
       "Username":"testing",
       "ProviderURL":"https://login.turbinelabs.io/auth/realms/turbine-labs"
     }`,
	)
	assert.Nil(t, ioutil.WriteFile(path, tcBytes, 0600))

	fset := tbnflag.NewTestFlagSet()
	ff := NewAPIConfigFromFlags(fset, APIConfigMayUseAuthToken(tokencache.NewStaticPath(path)))

	fset.Parse(nil)
	assert.Nil(t, ff.Validate())
	_, gotErr := ff.MakeEndpoint()
	assert.ErrorContains(t, gotErr, "session has timed out")
}

func TestNewAPIConfigFromFlagsValidateWithAuthTokenBadProvider(t *testing.T) {
	file, err := ioutil.TempFile("", "token-cache")
	assert.Nil(t, err)
	path := file.Name()
	defer func() { os.Remove(path) }()

	tcBytes := []byte(`{
       "Username":"testing",
       "ProviderURL":"http://www.example.com",
       "Token": {}
     }`,
	)
	assert.Nil(t, ioutil.WriteFile(path, tcBytes, 0600))

	fset := tbnflag.NewTestFlagSet()
	ff := NewAPIConfigFromFlags(fset, APIConfigMayUseAuthToken(tokencache.NewStaticPath(path)))
	ffImpl := ff.(*apiConfigFromFlags)

	fset.Parse(nil)
	assert.Nil(t, ff.Validate())
	_, gotErr := ff.MakeEndpoint()
	assert.ErrorContains(t, gotErr, "unable to construct OIDC client config")
	assert.NonNil(t, ffImpl.tokenCache)
	assert.DeepEqual(t, ffImpl.oauth2Config, oauth2.Config{})
}

func TestNewAPIConfigFromFlagsMakeEndpoint(t *testing.T) {
	ctrl := gomock.NewController(assert.Tracing(t))
	defer ctrl.Finish()

	file, err := ioutil.TempFile("", "token-cache")
	assert.Nil(t, err)
	path := file.Name()
	defer func() { os.Remove(path) }()

	user := "testing"
	provider := "https://login.turbinelabs.io/auth/realms/turbine-labs"
	tcBytes := []byte(`{
       "Username":"testing",
       "ProviderURL":"https://login.turbinelabs.io/auth/realms/turbine-labs",
       "Token": {}
     }`,
	)
	assert.Nil(t, ioutil.WriteFile(path, tcBytes, 0600))

	fset := tbnflag.NewTestFlagSet()
	ff := NewAPIConfigFromFlags(fset, APIConfigMayUseAuthToken(tokencache.NewStaticPath(path)))
	ffImpl := ff.(*apiConfigFromFlags)

	mockFF := apihttp.NewMockFromFlags(ctrl)
	ffImpl.clientFromFlags = mockFF

	client := &http.Client{}
	ep := apihttp.Endpoint{}
	ep.SetClient(client)
	gomock.InOrder(
		mockFF.EXPECT().Validate().Return(nil),
		mockFF.EXPECT().MakeEndpoint().Return(ep, nil),
	)

	fset.Parse(nil)

	assert.Nil(t, ff.Validate())
	gotEp, err := ff.MakeEndpoint()
	assert.Nil(t, err)
	assert.NotSameInstance(t, gotEp.Client(), client)

	assert.NonNil(t, ffImpl.oauth2Config)
	if assert.NonNil(t, ffImpl.tokenCache) {
		assert.Equal(t, ffImpl.tokenCache.Username, user)
		assert.Equal(t, ffImpl.tokenCache.ProviderURL, provider)
	}
}

func TestNewAPIConfigFromFlagsMakeEndpointKeyOverride(t *testing.T) {
	ctrl := gomock.NewController(assert.Tracing(t))
	defer ctrl.Finish()

	file, err := ioutil.TempFile("", "token-cache")
	assert.Nil(t, err)
	path := file.Name()
	defer func() { os.Remove(path) }()

	tcBytes := []byte(`{
       "Username":"testing",
       "ProviderURL":"https://login.turbinelabs.io/auth/realms/turbine-labs",
       "Token": {}
     }`,
	)
	assert.Nil(t, ioutil.WriteFile(path, tcBytes, 0600))

	fset := tbnflag.NewTestFlagSet()
	ff := NewAPIConfigFromFlags(fset, APIConfigMayUseAuthToken(tokencache.NewStaticPath(path)))
	ffImpl := ff.(*apiConfigFromFlags)

	mockFF := apihttp.NewMockFromFlags(ctrl)
	ffImpl.clientFromFlags = mockFF

	client := &http.Client{}
	ep := apihttp.Endpoint{}
	ep.SetClient(client)
	gomock.InOrder(
		mockFF.EXPECT().Validate().Return(nil),
		mockFF.EXPECT().MakeEndpoint().Return(ep, nil),
	)

	fset.Parse([]string{"-key=wheee"})

	assert.Nil(t, ff.Validate())
	gotEp, err := ff.MakeEndpoint()
	assert.Nil(t, err)
	assert.SameInstance(t, gotEp.Client(), client)
}
