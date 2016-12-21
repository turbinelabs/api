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
	"errors"
	"flag"
	"net/http"
	"testing"

	"github.com/golang/mock/gomock"

	apihttp "github.com/turbinelabs/api/http"
	"github.com/turbinelabs/test/assert"
)

var (
	fakeClient      = &http.Client{}
	fakeEndpoint, _ = apihttp.NewEndpoint(apihttp.HTTPS, "localhost", 1234)
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

	mockApiConfig.EXPECT().MakeEndpoint().Return(apihttp.Endpoint{}, errors.New("nope"))

	ff := &clientFromFlags{mockApiConfig}

	svc, err := ff.Make()
	assert.ErrorContains(t, err, "nope")
	assert.Nil(t, svc)
}
