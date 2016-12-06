package flags

import (
	"errors"
	"flag"
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

	fs := flag.NewFlagSet("stats-client", flag.PanicOnError)
	pfs := tbnflag.NewPrefixedFlagSet(fs, "pfix", "")

	apiConfigFromFlags := NewMockAPIConfigFromFlags(ctrl)

	ff := NewStatsClientFromFlags(pfs, StatsClientWithAPIConfigFromFlags(apiConfigFromFlags))
	assert.NonNil(t, ff)
	assert.SameInstance(t, apiConfigFromFlags, ff.(*statsClientFromFlags).apiConfigFromFlags)

	apiConfigFromFlags.EXPECT().Validate().Return(nil)
	assert.Nil(t, ff.Validate())

	expectedErr := errors.New("boom")
	apiConfigFromFlags.EXPECT().Validate().Return(expectedErr)
	assert.Equal(t, ff.Validate(), expectedErr)
}

func TestStatsClientFromFlagsDelegatesToAPIConfigFromFlags(t *testing.T) {
	ctrl := gomock.NewController(assert.Tracing(t))
	defer ctrl.Finish()

	fs := flag.NewFlagSet("stats-client", flag.PanicOnError)
	pfs := tbnflag.NewPrefixedFlagSet(fs, "pfix", "")

	mockExec := executor.NewMockExecutor(ctrl)

	endpoint, err := apihttp.NewEndpoint(apihttp.HTTPS, "example.com", 538)
	assert.Nil(t, err)
	assert.NonNil(t, endpoint)

	apiConfigFromFlags := NewMockAPIConfigFromFlags(ctrl)
	apiConfigFromFlags.EXPECT().MakeEndpoint().Return(endpoint, nil)
	apiConfigFromFlags.EXPECT().APIKey().Return("OTAY")

	ff := NewStatsClientFromFlags(pfs, StatsClientWithAPIConfigFromFlags(apiConfigFromFlags))
	assert.NonNil(t, ff)

	fs.Parse([]string{
		"-pfix.batch=false",
	})

	ffImpl := ff.(*statsClientFromFlags)
	assert.False(t, ffImpl.useBatching)
	assert.Equal(t, ffImpl.maxBatchDelay, DefaultMaxBatchDelay)
	assert.Equal(t, ffImpl.maxBatchSize, DefaultMaxBatchSize)

	statsClient, err := ff.Make(mockExec, log.NewNoopLogger())
	assert.NonNil(t, statsClient)
	assert.Nil(t, err)

	expectedErr := errors.New("no endpoints for you!")
	apiConfigFromFlags.EXPECT().
		MakeEndpoint().
		Return(apihttp.Endpoint{}, expectedErr)

	fs = flag.NewFlagSet("stats-client", flag.PanicOnError)
	pfs = tbnflag.NewPrefixedFlagSet(fs, "pfix", "")
	ff = NewStatsClientFromFlags(pfs, StatsClientWithAPIConfigFromFlags(apiConfigFromFlags))
	assert.NonNil(t, ff)

	statsClient, err = ff.Make(mockExec, log.NewNoopLogger())
	assert.Nil(t, statsClient)
	assert.NonNil(t, err)
}

func TestStatsClientFromFlagsCachesClient(t *testing.T) {
	ctrl := gomock.NewController(assert.Tracing(t))
	defer ctrl.Finish()

	fs := flag.NewFlagSet("stats-client", flag.PanicOnError)
	pfs := tbnflag.NewPrefixedFlagSet(fs, "pfix", "")

	mockExec := executor.NewMockExecutor(ctrl)
	otherMockExec := executor.NewMockExecutor(ctrl)

	endpoint, err := apihttp.NewEndpoint(apihttp.HTTPS, "example.com", 538)
	assert.Nil(t, err)
	assert.NonNil(t, endpoint)

	apiConfigFromFlags := NewMockAPIConfigFromFlags(ctrl)
	apiConfigFromFlags.EXPECT().MakeEndpoint().Return(endpoint, nil)
	apiConfigFromFlags.EXPECT().APIKey().Return("OTAY")

	ff := NewStatsClientFromFlags(pfs, StatsClientWithAPIConfigFromFlags(apiConfigFromFlags))
	assert.NonNil(t, ff)

	statsClient, err := ff.Make(mockExec, log.NewNoopLogger())
	assert.NonNil(t, statsClient)
	assert.Nil(t, err)

	statsClient2, err := ff.Make(otherMockExec, log.NewNoopLogger())
	assert.NonNil(t, statsClient2)
	assert.Nil(t, err)

	assert.SameInstance(t, statsClient2, statsClient)
}

func TestStatsClientFromFlagsCreatesBatchingClient(t *testing.T) {
	ctrl := gomock.NewController(assert.Tracing(t))
	defer ctrl.Finish()

	fs := flag.NewFlagSet("stats-client", flag.PanicOnError)
	pfs := tbnflag.NewPrefixedFlagSet(fs, "pfix", "")

	mockExec := executor.NewMockExecutor(ctrl)

	endpoint, err := apihttp.NewEndpoint(apihttp.HTTPS, "example.com", 538)
	assert.Nil(t, err)
	assert.NonNil(t, endpoint)

	apiConfigFromFlags := NewMockAPIConfigFromFlags(ctrl)
	apiConfigFromFlags.EXPECT().MakeEndpoint().Return(endpoint, nil)
	apiConfigFromFlags.EXPECT().APIKey().Return("OTAY")

	ff := NewStatsClientFromFlags(pfs, StatsClientWithAPIConfigFromFlags(apiConfigFromFlags))
	assert.NonNil(t, ff)

	fs.Parse([]string{
		"-pfix.max-batch-delay=5s",
		"-pfix.max-batch-size=99",
	})

	ffImpl := ff.(*statsClientFromFlags)
	assert.True(t, ffImpl.useBatching)
	assert.Equal(t, ffImpl.maxBatchDelay, 5*time.Second)
	assert.Equal(t, ffImpl.maxBatchSize, 99)

	statsClient, err := ff.Make(mockExec, log.NewNoopLogger())
	assert.NonNil(t, statsClient)
	assert.Nil(t, err)
}

func TestStatsClientFromFlagsValidatesBatchingClient(t *testing.T) {
	ctrl := gomock.NewController(assert.Tracing(t))
	defer ctrl.Finish()

	fs := flag.NewFlagSet("stats-client", flag.PanicOnError)
	pfs := tbnflag.NewPrefixedFlagSet(fs, "pfix", "")

	apiConfigFromFlags := NewMockAPIConfigFromFlags(ctrl)

	ff := NewStatsClientFromFlags(pfs, StatsClientWithAPIConfigFromFlags(apiConfigFromFlags))
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