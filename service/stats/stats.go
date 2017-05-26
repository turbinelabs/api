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

// Package stats defines the interfaces representing the portion of the
// Turbine Labs public API prefixed by /v1.0/stats
package stats

//go:generate mockgen -source $GOFILE -destination mock_$GOFILE -package $GOPACKAGE

import (
	"fmt"
	"strings"
	"time"

	"github.com/turbinelabs/nonstdlib/ptr"
	tbntime "github.com/turbinelabs/nonstdlib/time"
	"github.com/turbinelabs/stats"
)

var truePtr = ptr.Bool(true)

// StatsService forwards stats data to a remote stats-server.
type StatsService interface {
	// Forward the given stats payload.
	Forward(*Payload) (*ForwardResult, error)

	// Query for stats
	Query(*Query) (*QueryResult, error)

	// Closes the client and releases any resources it created.
	Close() error
}

// Creates a stats.Stats that uses this StatsService to forward
// arbitrary stats with the given source.
func AsStats(svc StatsService, source string) stats.Stats {
	return &asStats{
		svc:    svc,
		source: source,
	}
}

type asStats struct {
	svc    StatsService
	source string
	scope  string
	tags   map[string]string
}

var _ stats.Stats = &asStats{}

func (s *asStats) toTagMap(tags []stats.Tag) map[string]string {
	if len(tags) > 0 {
		tagsMap := make(map[string]string, len(tags)+len(s.tags))
		for _, tag := range tags {
			tagsMap[tag.K] = tag.V
		}
		for k, v := range s.tags {
			tagsMap[k] = v
		}
		return tagsMap
	}

	return s.tags
}

func (s *asStats) statName(stat string) string {
	if s.scope != "" {
		return fmt.Sprintf("%s/%s", s.scope, stat)
	}
	return stat
}

func (s *asStats) Count(stat string, value float64, tags ...stats.Tag) {
	payload := &Payload{
		Source: s.source,
		Stats: []Stat{
			{
				Name:      s.statName(stat),
				Value:     &value,
				Timestamp: tbntime.ToUnixMicro(time.Now()),
				Tags:      s.toTagMap(tags),
			},
		},
	}

	s.svc.Forward(payload)
}

func (s *asStats) Gauge(stat string, value float64, tags ...stats.Tag) {
	payload := &Payload{
		Source: s.source,
		Stats: []Stat{
			{
				Name:      s.statName(stat),
				Value:     &value,
				IsGauge:   truePtr,
				Timestamp: tbntime.ToUnixMicro(time.Now()),
				Tags:      s.toTagMap(tags),
			},
		},
	}

	s.svc.Forward(payload)
}

func (s *asStats) Histogram(stat string, value float64, tags ...stats.Tag) {
	payload := &Payload{
		Source: s.source,
		Stats: []Stat{
			{
				Name:      s.statName(stat),
				Value:     &value,
				IsGauge:   truePtr,
				Timestamp: tbntime.ToUnixMicro(time.Now()),
				Tags:      s.toTagMap(tags),
			},
		},
	}

	s.svc.Forward(payload)
}

func (s *asStats) Timing(stat string, value time.Duration, tags ...stats.Tag) {
	s.Histogram(stat, value.Seconds(), tags...)
}

func (s *asStats) AddTags(tags ...stats.Tag) {
	if s.tags == nil {
		s.tags = make(map[string]string, len(tags))
	}

	for _, tag := range tags {
		s.tags[tag.K] = tag.V
	}
}

func (s *asStats) Scope(scope string, scopes ...string) stats.Stats {
	newScopes := scope
	if len(scopes) > 0 {
		newScopes = fmt.Sprintf("%s/%s", scope, strings.Join(scopes, "/"))
	}

	if s.scope != "" {
		newScopes = fmt.Sprintf("%s/%s", s.scope, newScopes)
	}

	return &asStats{
		svc:    s.svc,
		source: s.source,
		scope:  newScopes,
	}
}

func (s *asStats) Close() error {
	return s.svc.Close()
}
