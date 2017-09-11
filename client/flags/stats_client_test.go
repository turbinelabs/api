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

func testStatsClientFromFlagsDelegatesToAPIConfigFromFlags(t *testing.T) {
	ctrl := gomock.NewController(assert.Tracing(t))
	defer ctrl.Finish()

	fs := tbnflag.NewTestFlagSet()

	mockExec := executor.NewMockExecutor(ctrl)

	mockExecFromFlags := executor.NewMockFromFlags(ctrl)
	mockExecFromFlags.EXPECT().Make(gomock.Any()).Return(mockExec)

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
		StatsClientWithExecutorFromFlags(mockExecFromFlags),
	)
	assert.NonNil(t, ff)

	fs.Parse([]string{
		"-pfix.batch=false",
	})

	ffImpl := ff.(*statsClientFromFlags)
	assert.False(t, ffImpl.useBatching)
	assert.Equal(t, ffImpl.maxBatchDelay, DefaultMaxBatchDelay)
	assert.Equal(t, ffImpl.maxBatchSize, DefaultMaxBatchSize)

	statsClient, err := ff.Make(log.NewNoopLogger())
	assert.NonNil(t, statsClient)
	assert.Nil(t, err)
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

	statsClient, err = ff.Make(log.NewNoopLogger())
	assert.Nil(t, statsClient)
	assert.NonNil(t, err)
}

func TestStatsClientFromFlagsDelegatesToAPIConfigFromFlags(t *testing.T) {
	testStatsClientFromFlagsDelegatesToAPIConfigFromFlags(t)
}

func testStatsClientFromFlagsCachesClient(t *testing.T) {
	ctrl := gomock.NewController(assert.Tracing(t))
	defer ctrl.Finish()

	pfs := tbnflag.NewTestFlagSet().Scope("pfix", "")

	mockExec := executor.NewMockExecutor(ctrl)

	mockExecFromFlags := executor.NewMockFromFlags(ctrl)
	mockExecFromFlags.EXPECT().Make(gomock.Any()).AnyTimes().Return(mockExec)

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
	statsClient, err = ff.Make(log.NewNoopLogger())
	assert.NonNil(t, statsClient)
	assert.Nil(t, err)

	statsClient2, err = ff.Make(log.NewNoopLogger())
	assert.NonNil(t, statsClient2)
	assert.Nil(t, err)

	assert.SameInstance(t, statsClient2, statsClient)
}

func TestStatsClientFromFlagsCachesClient(t *testing.T) {
	testStatsClientFromFlagsCachesClient(t)
}

func testStatsClientFromFlagsCreatesBatchingClient(t *testing.T) {
	ctrl := gomock.NewController(assert.Tracing(t))
	defer ctrl.Finish()

	fs := tbnflag.NewTestFlagSet()

	mockExec := executor.NewMockExecutor(ctrl)

	mockExecFromFlags := executor.NewMockFromFlags(ctrl)
	mockExecFromFlags.EXPECT().Make(gomock.Any()).Return(mockExec)

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
		StatsClientWithExecutorFromFlags(mockExecFromFlags),
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

	statsClient, err := ff.Make(log.NewNoopLogger())
	assert.NonNil(t, statsClient)
	assert.Nil(t, err)
}

func TestStatsClientFromFlagsCreatesBatchingClient(t *testing.T) {
	testStatsClientFromFlagsCreatesBatchingClient(t)
}
