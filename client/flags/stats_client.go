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
	"log"
	"time"

	"github.com/turbinelabs/api/client"
	statsapi "github.com/turbinelabs/api/service/stats"
	"github.com/turbinelabs/nonstdlib/executor"
	tbnflag "github.com/turbinelabs/nonstdlib/flag"
)

//go:generate mockgen -source $GOFILE -destination mock_$GOFILE -package $GOPACKAGE

const (
	DefaultMaxBatchDelay = 1 * time.Second
	DefaultMaxBatchSize  = 100
)

// StatsClientFromFlags validates and constructs a a statsapi.StatsService from command line
// flags.
type StatsClientFromFlags interface {
	Validate() error

	// Make constructs a statsapi.StatsService using the given Logger.
	Make(*log.Logger) (statsapi.StatsService, error)

	// Make constructs a statsapi.StatsServiceV2 using the given Logger.
	MakeV2(*log.Logger) (statsapi.StatsServiceV2, error)

	// APIKey returns the API Key used to construct the statsapi.StatsService.
	APIKey() string
}

// StatsClientOption represents an option passed to NewStatsClientFromFlags.
type StatsClientOption func(*statsClientFromFlags)

// StatsClientWithAPIConfigFromFlags configures
// NewStatsClientFromFlags to use a shared APIConfigFromFlags rather
// than creating its own.
func StatsClientWithAPIConfigFromFlags(apiConfigFromFlags APIConfigFromFlags) StatsClientOption {
	return func(ff *statsClientFromFlags) {
		ff.apiConfigFromFlags = apiConfigFromFlags
	}
}

// StatsClientWithExecutorFromFlags configures NewStatsClientFromFlags
// to use a shared executor.FromFlags rather than creating its own.
func StatsClientWithExecutorFromFlags(execFromFlags executor.FromFlags) StatsClientOption {
	return func(ff *statsClientFromFlags) {
		ff.execFromFlags = execFromFlags
	}
}

func NewStatsClientFromFlags(
	clientApp client.App,
	pfs tbnflag.FlagSet,
	options ...StatsClientOption,
) StatsClientFromFlags {
	ff := &statsClientFromFlags{clientApp: clientApp}

	for _, option := range options {
		option(ff)
	}

	if ff.apiConfigFromFlags == nil {
		ff.apiConfigFromFlags = NewAPIConfigFromFlags(pfs)
	}

	if ff.execFromFlags == nil {
		ff.execFromFlags = executor.NewFromFlags(pfs)
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
	clientApp          client.App
	apiConfigFromFlags APIConfigFromFlags
	execFromFlags      executor.FromFlags
	useBatching        bool
	maxBatchDelay      time.Duration
	maxBatchSize       int

	cachedClient   statsapi.StatsService
	cachedV2Client statsapi.StatsServiceV2
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
	logger *log.Logger,
) (statsapi.StatsService, error) {
	if ff.cachedClient != nil {
		return ff.cachedClient, nil
	}

	endpoint, err := ff.apiConfigFromFlags.MakeEndpoint()
	if err != nil {
		return nil, err
	}

	exec := ff.execFromFlags.Make(logger)

	var stats statsapi.StatsService
	if ff.useBatching {
		stats, err = client.NewBatchingStatsClient(
			ff.maxBatchDelay,
			ff.maxBatchSize,
			endpoint,
			ff.apiConfigFromFlags.APIKey(),
			ff.clientApp,
			exec,
			logger,
		)
	} else {
		stats, err = client.NewStatsClient(
			endpoint,
			ff.apiConfigFromFlags.APIKey(),
			ff.clientApp,
			exec,
		)
	}

	if err != nil {
		return nil, err
	}

	ff.cachedClient = stats

	return stats, nil
}

func (ff *statsClientFromFlags) MakeV2(
	logger *log.Logger,
) (statsapi.StatsServiceV2, error) {
	if ff.cachedV2Client != nil {
		return ff.cachedV2Client, nil
	}

	endpoint, err := ff.apiConfigFromFlags.MakeEndpoint()
	if err != nil {
		return nil, err
	}

	exec := ff.execFromFlags.Make(logger)

	var stats statsapi.StatsServiceV2
	if ff.useBatching {
		stats, err = client.NewBatchingStatsV2Client(
			ff.maxBatchDelay,
			ff.maxBatchSize,
			endpoint,
			ff.apiConfigFromFlags.APIKey(),
			ff.clientApp,
			exec,
			logger,
		)
	} else {
		stats, err = client.NewStatsV2Client(
			endpoint,
			ff.apiConfigFromFlags.APIKey(),
			ff.clientApp,
			exec,
		)
	}

	if err != nil {
		return nil, err
	}

	ff.cachedV2Client = stats

	return stats, nil
}

func (ff *statsClientFromFlags) APIKey() string {
	return ff.apiConfigFromFlags.APIKey()
}
