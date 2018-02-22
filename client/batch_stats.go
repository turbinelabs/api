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

package client

import (
	"errors"
	"log"
	"sync"
	"sync/atomic"
	"time"

	apihttp "github.com/turbinelabs/api/http"
	statsapi "github.com/turbinelabs/api/service/stats"
	"github.com/turbinelabs/nonstdlib/executor"
	tbntime "github.com/turbinelabs/nonstdlib/time"
)

type httpBatchingStatsV2 struct {
	internalStatsClient

	maxDelay time.Duration
	maxSize  int

	batchers map[string]*payloadV2Batcher
	mutex    *sync.RWMutex

	inFlight int32

	logger *log.Logger
}

// NewBatchingStatsV2Client returns a non-blocking implementation of
// StatsServiceV2. Each invocation of ForwardV2 accepts a single
// Payload. The client will return immediately, reporting that all
// stats were successfully sent. Internally, the stats are buffered
// until the buffer contains at least maxSize stats or maxDelay time
// has elapsed since the oldest stats in the buffer were added. At
// that point the buffered stats are forwarded. Failures are logged,
// but not reported to the caller. Separate buffers and deadlines are
// maintained for each unique source and zone combination. In
// addition, the buffering is optimized to assume that the payloads
// proxy, proxy version and named limits do not vary across
// payloads. If the proxy, proxy version, or limits (with the same
// name) vary across payloads, buffers will be flushed prematurely.
func NewBatchingStatsV2Client(
	maxDelay time.Duration,
	maxSize int,
	dest apihttp.Endpoint,
	apiKey string,
	clientApp App,
	exec executor.Executor,
	logger *log.Logger,
) (statsapi.StatsService, error) {
	if maxDelay < time.Second {
		return nil, errors.New("max delay must be at least 1 second")
	}

	if maxSize < 1 {
		return nil, errors.New("max size must be at least 1")
	}

	underlyingStatsClient, err := newInternalStatsClient(dest, v2ForwardPath, apiKey, clientApp, exec)
	if err != nil {
		return nil, err
	}

	return &httpBatchingStatsV2{
		internalStatsClient: underlyingStatsClient,
		maxDelay:            maxDelay,
		maxSize:             maxSize,
		batchers:            map[string]*payloadV2Batcher{},
		mutex:               &sync.RWMutex{},
		logger:              logger,
	}, nil
}

func mkKey(source, zone string) string {
	sourceLen := len(source)
	key := make([]byte, sourceLen+len(zone)+1)
	copy(key, []byte(source))
	key[sourceLen] = '|'
	copy(key[sourceLen+1:], []byte(zone))
	return string(key)
}

func (hs *httpBatchingStatsV2) getBatcher(payload *statsapi.Payload) *payloadV2Batcher {
	key := mkKey(payload.Source, payload.Zone)

	hs.mutex.RLock()
	defer hs.mutex.RUnlock()

	if batcher, ok := hs.batchers[key]; ok {
		return batcher
	}

	batcher := &payloadV2Batcher{
		client: hs,
		source: payload.Source,
		zone:   payload.Zone,
		ch:     make(chan *statsapi.Payload, 10),
	}

	hs.batchers[key] = batcher
	batcher.start()

	return batcher
}

func (hs *httpBatchingStatsV2) ForwardV2(payload *statsapi.Payload) (*statsapi.ForwardResult, error) {
	batcher := hs.getBatcher(payload)

	batcher.ch <- payload

	return &statsapi.ForwardResult{NumAccepted: len(payload.Stats)}, nil
}

func (hs *httpBatchingStatsV2) closeBatchers() {
	hs.mutex.Lock()
	defer hs.mutex.Unlock()

	for _, batcher := range hs.batchers {
		close(batcher.ch)
	}
	hs.batchers = map[string]*payloadV2Batcher{}
}

func (hs *httpBatchingStatsV2) Close() error {
	hs.closeBatchers()

	hs.logger.Print("waiting for final requests to complete")
	start := time.Now()
	for time.Since(start) < 15*time.Second {
		time.Sleep(100 * time.Millisecond)
		if atomic.LoadInt32(&hs.inFlight) == 0 {
			return nil
		}
	}

	return errors.New("timed out waiting for final requests to complete")
}

type payloadV2Batcher struct {
	client       *httpBatchingStatsV2
	source       string
	zone         string
	proxy        *string
	proxyVersion *string
	limits       map[string][]float64
	buffer       []statsapi.Stat
	ch           chan *statsapi.Payload
}

func (b *payloadV2Batcher) start() {
	go b.run(tbntime.NewTimer(0))
}

func (b *payloadV2Batcher) run(timer tbntime.Timer) {
	b.buffer = make([]statsapi.Stat, 0, b.client.maxSize)

	if !timer.Stop() {
		<-timer.C()
	}
	timer.Reset(b.client.maxDelay)
	timerIsLive := true

	for {
		select {
		case <-timer.C():
			if len(b.buffer) > 0 {
				b.flush()
			}
			timerIsLive = false

		case payload, ok := <-b.ch:
			if !ok {
				if len(b.buffer) > 0 {
					b.flush()
				}
				timer.Stop()
				return
			}

			flushed := b.write(payload)
			if len(b.buffer) == 0 {
				timer.Stop()
				timerIsLive = false
			} else if !timerIsLive || flushed {
				timer.Reset(b.client.maxDelay)
			}
		}
	}
}

func limitsRequireFlush(payloadLimitsMap, batchLimitsMap map[string][]float64) bool {
	if len(payloadLimitsMap) == 0 || len(batchLimitsMap) == 0 {
		return false
	}

	// Verify that all existing limits are identical to limits
	// from this payload, matching by name.
	for name, batchLimits := range batchLimitsMap {
		if payloadLimits, found := payloadLimitsMap[name]; found {
			if len(payloadLimits) == len(batchLimits) {
				for idx, v := range payloadLimits {
					if v != batchLimits[idx] {
						return true
					}
				}
			} else {
				return true
			}
		}
	}

	return false
}

func (b *payloadV2Batcher) write(p *statsapi.Payload) bool {
	flushed := false

	proxyRequiresFlush :=
		(b.proxy != nil && p.Proxy != nil && *b.proxy != *p.Proxy) ||
			(b.proxyVersion != nil && p.ProxyVersion != nil && *b.proxyVersion != *p.ProxyVersion)
	if proxyRequiresFlush {
		b.flush()
		flushed = true
	}

	if limitsRequireFlush(p.Limits, b.limits) {
		b.flush()
		flushed = true
	}

	if b.proxy == nil {
		b.proxy = p.Proxy
	}

	if b.proxyVersion == nil {
		b.proxyVersion = p.ProxyVersion
	}

	if len(b.limits) == 0 {
		b.limits = p.Limits
	} else if len(p.Limits) > 0 {
		for name, limits := range p.Limits {
			if _, found := b.limits[name]; !found {
				b.limits[name] = limits
			}
		}
	}

	b.buffer = append(b.buffer, p.Stats...)

	if len(b.buffer) >= b.client.maxSize {
		b.flush()
		flushed = true
	}

	return flushed
}

func (b *payloadV2Batcher) flush() {
	payload := &statsapi.Payload{
		Source:       b.source,
		Zone:         b.zone,
		Proxy:        b.proxy,
		ProxyVersion: b.proxyVersion,
		Limits:       b.limits,
		Stats:        b.buffer,
	}

	b.forward(payload)

	b.proxy = nil
	b.proxyVersion = nil
	b.limits = nil
	b.buffer = b.buffer[0:0]
}

func (b *payloadV2Batcher) forward(payload *statsapi.Payload) {
	atomic.AddInt32(&b.client.inFlight, 1)

	err := b.client.ForwardWithCallback(
		payload,
		func(try executor.Try) {
			atomic.AddInt32(&b.client.inFlight, -1)
			if try.IsError() {
				b.client.logger.Printf(
					"Failed to forward payload: %+v: %s",
					payload,
					try.Error().Error(),
				)
			}
		},
	)
	if err != nil {
		atomic.AddInt32(&b.client.inFlight, -1)

		b.client.logger.Printf(
			"Failed to enqueue request: %+v: %s",
			payload,
			err.Error(),
		)
	}
}
