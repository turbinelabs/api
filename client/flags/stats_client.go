package flags

import (
	"errors"
	"log"
	"time"

	"github.com/turbinelabs/api/client"
	statsapi "github.com/turbinelabs/api/service/stats"
	"github.com/turbinelabs/nonstdlib/executor"
	"github.com/turbinelabs/nonstdlib/flag"
)

//go:generate mockgen -source $GOFILE -destination mock_$GOFILE -package $GOPACKAGE

const (
	DefaultMaxBatchDelay = 1 * time.Second
	DefaultMaxBatchSize  = 100
)

// FromFlags validates and constructs a a statsapi.StatsService from command line
// flags.
type StatsClientFromFlags interface {
	Validate() error

	// Constructs a statsapi.StatsService using the given Executor and Logger.
	Make(executor.Executor, *log.Logger) (statsapi.StatsService, error)

	// Returns the API Key used to construct the statsapi.StatsService.
	APIKey() string
}

type StatsClientOption func(*statsClientFromFlags)

func StatsClientWithAPIConfigFromFlags(apiConfigFromFlags APIConfigFromFlags) StatsClientOption {
	return func(ff *statsClientFromFlags) {
		ff.apiConfigFromFlags = apiConfigFromFlags
	}
}

func NewStatsClientFromFlags(pfs *flag.PrefixedFlagSet, options ...StatsClientOption) StatsClientFromFlags {
	ff := &statsClientFromFlags{}

	for _, option := range options {
		option(ff)
	}

	if ff.apiConfigFromFlags == nil {
		ff.apiConfigFromFlags = NewPrefixedAPIConfigFromFlags(pfs)
	}

	pfs.BoolVar(
		&ff.useBatching,
		"batch",
		true,
		"If true, {{NAME}} requests are batched together for performance.",
	)

	pfs.DurationVar(
		&ff.maxBatchDelay,
		"max-batch-delay",
		DefaultMaxBatchDelay,
		"If batching is enabled, the maximum amount of time requests are held before transmission",
	)

	pfs.IntVar(
		&ff.maxBatchSize,
		"max-batch-size",
		DefaultMaxBatchSize,
		"If batching is enabled, the maximum number of requests that will be combined.",
	)

	return ff

}

type statsClientFromFlags struct {
	apiConfigFromFlags APIConfigFromFlags
	useBatching        bool
	maxBatchDelay      time.Duration
	maxBatchSize       int

	cachedClient statsapi.StatsService
}

func (ff *statsClientFromFlags) Validate() error {
	if ff.useBatching {
		if ff.maxBatchDelay < 1*time.Second {
			return errors.New(
				"max-batch-delay may not be less than 1 second",
			)
		}

		if ff.maxBatchSize < 1 {
			return errors.New(
				"max-batch-size may not be less than 1",
			)
		}
	}

	return ff.apiConfigFromFlags.Validate()
}

func (ff *statsClientFromFlags) Make(
	exec executor.Executor,
	logger *log.Logger,
) (statsapi.StatsService, error) {
	if ff.cachedClient != nil {
		return ff.cachedClient, nil
	}

	endpoint, err := ff.apiConfigFromFlags.MakeEndpoint()
	if err != nil {
		return nil, err
	}

	var stats statsapi.StatsService
	if ff.useBatching {
		stats, err = client.NewBatchingStatsClient(
			ff.maxBatchDelay,
			ff.maxBatchSize,
			endpoint,
			ff.apiConfigFromFlags.APIKey(),
			exec,
			logger,
		)
	} else {
		stats, err = client.NewStatsClient(endpoint, ff.apiConfigFromFlags.APIKey(), exec)
	}

	if err != nil {
		return nil, err
	}

	ff.cachedClient = stats

	return stats, nil
}

func (ff *statsClientFromFlags) APIKey() string {
	return ff.apiConfigFromFlags.APIKey()
}
