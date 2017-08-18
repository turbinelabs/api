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
	"testing"
	"time"

	"github.com/golang/mock/gomock"

	apihttp "github.com/turbinelabs/api/http"
	"github.com/turbinelabs/nonstdlib/executor"
	tbnflag "github.com/turbinelabs/nonstdlib/flag"
	"github.com/turbinelabs/test/assert"
	"github.com/turbinelabs/test/log"
)

func TestStatsClientFromFlagsValidatesNormalClient(t *testing.T) {
	ctrl := gomock.NewController(assert.Tracing(t))
	defer ctrl.Finish()

	fs := tbnflag.NewTestFlagSet()

	apiConfigFromFlags := NewMockAPIConfigFromFlags(ctrl)

	ff := NewStatsClientFromFlags(
		"app",
		fs.Scope("pfix", ""),
		StatsClientWithAPIConfigFromFlags(apiConfigFromFlags),
	)
	assert.NonNil(t, ff)
	assert.SameInstance(t, apiConfigFromFlags, ff.(*statsClientFromFlags).apiConfigFromFlags)

	apiConfigFromFlags.EXPECT().Validate().Return(nil)
	assert.Nil(t, ff.Validate())

	expectedErr := errors.New("boom")
	apiConfigFromFlags.EXPECT().Validate().Return(expectedErr)
	assert.Equal(t, ff.Validate(), expectedErr)
}

func TestStatsClientFromFlagsValidatesBatchingClient(t *testing.T) {
	ctrl := gomock.NewController(assert.Tracing(t))
	defer ctrl.Finish()

	fs := tbnflag.NewTestFlagSet()

	apiConfigFromFlags := NewMockAPIConfigFromFlags(ctrl)

	ff := NewStatsClientFromFlags(
		"app",
		fs.Scope("pfix", ""),
		StatsClientWithAPIConfigFromFlags(apiConfigFromFlags),
	)
	assert.NonNil(t, ff)

	fs.Parse([]string{
		"-pfix.batch=true",
		"-pfix.max-batch-delay=0s",
	})

	assert.ErrorContains(t, ff.Validate(), "max-batch-delay")

	fs.Parse([]string{
		"-pfix.batch=true",
		"-pfix.max-batch-delay=1s",
		"-pfix.max-batch-size=0",
	})

	assert.ErrorContains(t, ff.Validate(), "max-batch-size")

	fs.Parse([]string{
		"-pfix.batch=true",
		"-pfix.max-batch-delay=1s",
		"-pfix.max-batch-size=1",
	})

	expectedErr := errors.New("boom")
	apiConfigFromFlags.EXPECT().Validate().Return(expectedErr)
	assert.Equal(t, ff.Validate(), expectedErr)

	apiConfigFromFlags.EXPECT().Validate().Return(nil)
	assert.Nil(t, ff.Validate())
}

func testStatsClientFromFlagsDelegatesToAPIConfigFromFlags(t *testing.T, useV2 bool) {
	ctrl := gomock.NewController(assert.Tracing(t))
	defer ctrl.Finish()

	fs := tbnflag.NewTestFlagSet()

	mockExec := executor.NewMockExecutor(ctrl)

	endpoint, err := apihttp.NewEndpoint(apihttp.HTTPS, "example.com:538")
	assert.Nil(t, err)
	assert.NonNil(t, endpoint)

	apiConfigFromFlags := NewMockAPIConfigFromFlags(ctrl)
	apiConfigFromFlags.EXPECT().MakeEndpoint().Return(endpoint, nil)
	apiConfigFromFlags.EXPECT().APIKey().Return("OTAY")

	ff := NewStatsClientFromFlags(
		"app",
		fs.Scope("pfix", ""),
		StatsClientWithAPIConfigFromFlags(apiConfigFromFlags),
	)
	assert.NonNil(t, ff)

	fs.Parse([]string{
		"-pfix.batch=false",
	})

	ffImpl := ff.(*statsClientFromFlags)
	assert.False(t, ffImpl.useBatching)
	assert.Equal(t, ffImpl.maxBatchDelay, DefaultMaxBatchDelay)
	assert.Equal(t, ffImpl.maxBatchSize, DefaultMaxBatchSize)

	if useV2 {
		statsClient, err := ff.MakeV2(mockExec, log.NewNoopLogger())
		assert.NonNil(t, statsClient)
		assert.Nil(t, err)
	} else {
		statsClient, err := ff.Make(mockExec, log.NewNoopLogger())
		assert.NonNil(t, statsClient)
		assert.Nil(t, err)
	}
	expectedErr := errors.New("no endpoints for you")
	apiConfigFromFlags.EXPECT().
		MakeEndpoint().
		Return(apihttp.Endpoint{}, expectedErr)

	fs = tbnflag.NewTestFlagSet()
	ff = NewStatsClientFromFlags(
		"app",
		fs.Scope("pfix", ""),
		StatsClientWithAPIConfigFromFlags(apiConfigFromFlags),
	)
	assert.NonNil(t, ff)

	if useV2 {
		statsClient, err := ff.MakeV2(mockExec, log.NewNoopLogger())
		assert.Nil(t, statsClient)
		assert.NonNil(t, err)
	} else {
		statsClient, err := ff.Make(mockExec, log.NewNoopLogger())
		assert.Nil(t, statsClient)
		assert.NonNil(t, err)
	}
}

func TestStatsClientFromFlagsDelegatesToAPIConfigFromFlags(t *testing.T) {
	testStatsClientFromFlagsDelegatesToAPIConfigFromFlags(t, false)
}

func TestStatsClientFromFlagsDelegatesToAPIConfigFromFlagsV2(t *testing.T) {
	testStatsClientFromFlagsDelegatesToAPIConfigFromFlags(t, true)
}

func testStatsClientFromFlagsCachesClient(t *testing.T, useV2 bool) {
	ctrl := gomock.NewController(assert.Tracing(t))
	defer ctrl.Finish()

	pfs := tbnflag.NewTestFlagSet().Scope("pfix", "")

	mockExec := executor.NewMockExecutor(ctrl)
	otherMockExec := executor.NewMockExecutor(ctrl)

	endpoint, err := apihttp.NewEndpoint(apihttp.HTTPS, "example.com:538")
	assert.Nil(t, err)
	assert.NonNil(t, endpoint)

	apiConfigFromFlags := NewMockAPIConfigFromFlags(ctrl)
	apiConfigFromFlags.EXPECT().MakeEndpoint().Return(endpoint, nil)
	apiConfigFromFlags.EXPECT().APIKey().Return("OTAY")

	ff := NewStatsClientFromFlags(
		"app",
		pfs,
		StatsClientWithAPIConfigFromFlags(apiConfigFromFlags),
	)
	assert.NonNil(t, ff)

	var statsClient, statsClient2 interface{}
	if useV2 {
		statsClient, err = ff.MakeV2(mockExec, log.NewNoopLogger())
	} else {
		statsClient, err = ff.Make(mockExec, log.NewNoopLogger())
	}
	assert.NonNil(t, statsClient)
	assert.Nil(t, err)

	if useV2 {
		statsClient2, err = ff.MakeV2(otherMockExec, log.NewNoopLogger())
	} else {
		statsClient2, err = ff.Make(otherMockExec, log.NewNoopLogger())
	}
	assert.NonNil(t, statsClient2)
	assert.Nil(t, err)

	assert.SameInstance(t, statsClient2, statsClient)
}

func TestStatsClientFromFlagsCachesClient(t *testing.T) {
	testStatsClientFromFlagsCachesClient(t, false)
}

func TestStatsClientFromFlagsCachesClientV2(t *testing.T) {
	testStatsClientFromFlagsCachesClient(t, true)
}

func testStatsClientFromFlagsCreatesBatchingClient(t *testing.T, useV2 bool) {
	ctrl := gomock.NewController(assert.Tracing(t))
	defer ctrl.Finish()

	fs := tbnflag.NewTestFlagSet()

	mockExec := executor.NewMockExecutor(ctrl)

	endpoint, err := apihttp.NewEndpoint(apihttp.HTTPS, "example.com:538")
	assert.Nil(t, err)
	assert.NonNil(t, endpoint)

	apiConfigFromFlags := NewMockAPIConfigFromFlags(ctrl)
	apiConfigFromFlags.EXPECT().MakeEndpoint().Return(endpoint, nil)
	apiConfigFromFlags.EXPECT().APIKey().Return("OTAY")

	ff := NewStatsClientFromFlags(
		"app",
		fs.Scope("pfix", ""),
		StatsClientWithAPIConfigFromFlags(apiConfigFromFlags),
	)
	assert.NonNil(t, ff)

	fs.Parse([]string{
		"-pfix.max-batch-delay=5s",
		"-pfix.max-batch-size=99",
	})

	ffImpl := ff.(*statsClientFromFlags)
	assert.True(t, ffImpl.useBatching)
	assert.Equal(t, ffImpl.maxBatchDelay, 5*time.Second)
	assert.Equal(t, ffImpl.maxBatchSize, 99)

	if useV2 {
		statsClient, err := ff.MakeV2(mockExec, log.NewNoopLogger())
		assert.NonNil(t, statsClient)
		assert.Nil(t, err)
	} else {
		statsClient, err := ff.Make(mockExec, log.NewNoopLogger())
		assert.NonNil(t, statsClient)
		assert.Nil(t, err)
	}
}

func TestStatsClientFromFlagsCreatesBatchingClient(t *testing.T) {
	testStatsClientFromFlagsCreatesBatchingClient(t, false)
}

func TestStatsClientFromFlagsCreatesBatchingClientV2(t *testing.T) {
	testStatsClientFromFlagsCreatesBatchingClient(t, true)
}
