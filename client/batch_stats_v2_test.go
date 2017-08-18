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
	"sync"
	"testing"
	"time"

	"github.com/golang/mock/gomock"

	statsapi "github.com/turbinelabs/api/service/stats"
	"github.com/turbinelabs/nonstdlib/executor"
	"github.com/turbinelabs/nonstdlib/ptr"
	tbntime "github.com/turbinelabs/nonstdlib/time"
	"github.com/turbinelabs/test/assert"
	"github.com/turbinelabs/test/log"
)

func payloadV2OfSize(s int) *statsapi.PayloadV2 {
	switch s {
	case 0:
		return &statsapi.PayloadV2{Source: sourceString1, Zone: zoneString1}

	case 1:
		p := *payloadV2
		return &p

	default:
		a := make([]statsapi.StatV2, s)
		for i := 0; i < s; i++ {
			a[i] = payloadV2.Stats[0]
		}
		return &statsapi.PayloadV2{Source: sourceString1, Zone: zoneString1, Stats: a}
	}
}

type batcherV2Test struct {
	expectedPayloadSizes   []int                 // the output batches: N payloads of given sizes
	expectedCustomPayloads []*statsapi.PayloadV2 // or else specific expected batches

	numForwards    int                   // number of payloads passed to the batcher
	forwardedSize  int                   // size of each payload
	customPayloads []*statsapi.PayloadV2 // or else specific input payloads

	closeAfterLastPayload bool

	maxDelay time.Duration // batcher timer setting
	maxSize  int           // batcher max payload size setting

	timerBehavior func(*tbntime.MockTimer)
}

func (bt batcherV2Test) run(t *testing.T) {
	ctrl := gomock.NewController(assert.Tracing(t))
	defer ctrl.Finish()

	cbfChan := make(chan executor.CallbackFunc, 2)

	mockUnderlyingStatsClient := newMockInternalStatsClient(ctrl)

	if bt.expectedCustomPayloads == nil {
		for _, payloadSize := range bt.expectedPayloadSizes {
			expectedPayload := payloadV2OfSize(payloadSize)
			mockUnderlyingStatsClient.EXPECT().
				ForwardWithCallback(expectedPayload, gomock.Any()).
				Do(func(_ *statsapi.PayloadV2, cb executor.CallbackFunc) { cbfChan <- cb }).
				Return(nil)
		}
	} else {
		for _, expectedPayload := range bt.expectedCustomPayloads {
			mockUnderlyingStatsClient.EXPECT().
				ForwardWithCallback(expectedPayload, gomock.Any()).
				Do(func(_ *statsapi.PayloadV2, cb executor.CallbackFunc) { cbfChan <- cb }).
				Return(nil)
		}
	}

	batcher := &payloadV2Batcher{
		client: &httpBatchingStatsV2{
			internalStatsClient: mockUnderlyingStatsClient,
			maxDelay:            bt.maxDelay,
			maxSize:             bt.maxSize,
		},
		source: sourceString1,
		zone:   zoneString1,
		ch:     make(chan *statsapi.PayloadV2, 2*bt.maxSize),
	}

	mockTimer := tbntime.NewMockTimer(ctrl)

	bt.timerBehavior(mockTimer)

	wg := &sync.WaitGroup{}
	wg.Add(1)
	defer wg.Wait()

	go func() {
		defer wg.Done()
		batcher.run(mockTimer)
	}()

	if !bt.closeAfterLastPayload {
		defer close(batcher.ch)
	}

	if bt.customPayloads == nil {
		for i := 0; i < bt.numForwards; i++ {
			batcher.ch <- payloadV2OfSize(bt.forwardedSize)
		}
	} else {
		for _, payload := range bt.customPayloads {
			batcher.ch <- payload
		}
	}

	if bt.closeAfterLastPayload {
		close(batcher.ch)
	}

	if bt.expectedCustomPayloads == nil {
		for range bt.expectedPayloadSizes {
			<-cbfChan
		}
	} else {
		for range bt.expectedCustomPayloads {
			<-cbfChan
		}
	}
}

func TestNewBatchingStatsV2Client(t *testing.T) {
	ctrl := gomock.NewController(assert.Tracing(t))
	defer ctrl.Finish()

	exec := executor.NewMockExecutor(ctrl)
	client, err := NewBatchingStatsV2Client(
		time.Second,
		100,
		endpoint,
		clientTestAPIKey,
		clientTestApp,
		exec,
		log.NewNoopLogger(),
	)
	assert.NonNil(t, client)
	assert.Nil(t, err)

	clientImpl, ok := client.(*httpBatchingStatsV2)
	assert.True(t, ok)

	assert.NonNil(t, clientImpl.internalStatsClient)
	underlyingImpl, ok := clientImpl.internalStatsClient.(*httpStats)
	assert.True(t, ok)
	assert.NotDeepEqual(t, underlyingImpl.dest, endpoint)
	assert.SameInstance(t, underlyingImpl.exec, exec)

	assert.Equal(t, clientImpl.maxDelay, time.Second)
	assert.Equal(t, clientImpl.maxSize, 100)
	assert.NonNil(t, clientImpl.batchers)
	assert.Equal(t, len(clientImpl.batchers), 0)
	assert.NonNil(t, clientImpl.mutex)
}

func TestNewBatchingStatsV2ClientValidation(t *testing.T) {
	ctrl := gomock.NewController(assert.Tracing(t))
	defer ctrl.Finish()

	exec := executor.NewMockExecutor(ctrl)
	log := log.NewNoopLogger()

	client, err := NewBatchingStatsV2Client(
		999*time.Millisecond,
		1,
		endpoint,
		clientTestAPIKey,
		clientTestApp,
		exec,
		log,
	)
	assert.Nil(t, client)
	assert.ErrorContains(t, err, "max delay must be at least 1 second")

	client, err = NewBatchingStatsV2Client(
		time.Second,
		0,
		endpoint,
		clientTestAPIKey,
		clientTestApp,
		exec,
		log,
	)
	assert.Nil(t, client)
	assert.ErrorContains(t, err, "max size must be at least 1")
}

func TestHttpBatchingStatsV2GetBatcher(t *testing.T) {
	client := &httpBatchingStatsV2{
		batchers: map[string]*payloadV2Batcher{},
		mutex:    &sync.RWMutex{},
	}

	batcher := client.getBatcher(&statsapi.PayloadV2{Source: "s", Zone: "z"})
	defer close(batcher.ch)

	assert.NonNil(t, batcher)
	assert.SameInstance(t, batcher.client, client)
	assert.Equal(t, batcher.source, "s")
	assert.Equal(t, batcher.zone, "z")
	assert.NonNil(t, batcher.ch)

	batcher2 := client.getBatcher(&statsapi.PayloadV2{Source: "s", Zone: "z"})
	assert.SameInstance(t, batcher2, batcher)
}

func TestHttpBatchingStatsV2Forward(t *testing.T) {
	client := &httpBatchingStatsV2{
		batchers: map[string]*payloadV2Batcher{},
		mutex:    &sync.RWMutex{},
	}

	expectedPayload := payloadV2OfSize(3)

	result, err := client.ForwardV2(expectedPayload)
	assert.NonNil(t, result)
	assert.Nil(t, err)
	assert.Equal(t, result.NumAccepted, 3)

	batcher, ok := client.batchers[expectedPayload.Source+"|"+expectedPayload.Zone]
	assert.True(t, ok)
	assert.NonNil(t, batcher)
	defer close(batcher.ch)

	select {
	case payload := <-batcher.ch:
		assert.SameInstance(t, payload, expectedPayload)

	default:
		assert.Failed(t, "payload not enqueued in batcher's channel")
	}
}

func TestHttpBatchingStatsV2Close(t *testing.T) {
	client := &httpBatchingStatsV2{
		batchers: map[string]*payloadV2Batcher{},
		mutex:    &sync.RWMutex{},
	}

	client.getBatcher(&statsapi.PayloadV2{Source: "this-source", Zone: "zone"})
	client.getBatcher(&statsapi.PayloadV2{Source: "that-source", Zone: "zone"})
	assert.Equal(t, len(client.batchers), 2)

	ch1 := client.batchers["this-source|zone"].ch
	ch2 := client.batchers["that-source|zone"].ch

	assert.Nil(t, client.Close())
	assert.Equal(t, len(client.batchers), 0)

	select {
	case _, ok := <-ch1:
		assert.False(t, ok)
	default:
		assert.Failed(t, "expected closed channel ch1, saw empty channel")
	}

	select {
	case _, ok := <-ch2:
		assert.False(t, ok)
	default:
		assert.Failed(t, "expected closed channel ch2, saw empty channel")
	}
}

func TestPayloadV2BatcherRunSendsBatchBySize(t *testing.T) {
	batcherV2Test{
		expectedPayloadSizes: []int{5},
		numForwards:          5,
		forwardedSize:        1,
		maxDelay:             time.Second,
		maxSize:              5,
		timerBehavior: func(mockTimer *tbntime.MockTimer) {
			emptyTimeChan := make(chan time.Time, 1)

			gomock.InOrder(
				mockTimer.EXPECT().Stop().Return(true),
				mockTimer.EXPECT().Reset(1*time.Second).Return(false),
				mockTimer.EXPECT().C().Times(5).Return(emptyTimeChan),
				mockTimer.EXPECT().Stop().Return(true),
				mockTimer.EXPECT().C().Return(emptyTimeChan),
				mockTimer.EXPECT().Stop().Return(false),
			)
		},
	}.run(t)
}

func TestPayloadV2BatcherRunSendsBatchBySizeOnFirstCall(t *testing.T) {
	batcherV2Test{
		expectedPayloadSizes: []int{5},
		numForwards:          1,
		forwardedSize:        5,
		maxDelay:             time.Second,
		maxSize:              3,
		timerBehavior: func(mockTimer *tbntime.MockTimer) {
			emptyTimeChan := make(chan time.Time, 1)

			gomock.InOrder(
				mockTimer.EXPECT().Stop().Return(true),
				mockTimer.EXPECT().Reset(1*time.Second).Return(false),
				mockTimer.EXPECT().C().Times(1).Return(emptyTimeChan),
				mockTimer.EXPECT().Stop().Return(true),
				mockTimer.EXPECT().C().Return(emptyTimeChan),
				mockTimer.EXPECT().Stop().Return(false),
			)
		},
	}.run(t)
}

func TestPayloadV2BatcherRunSendsBatchByDelay(t *testing.T) {
	batcherV2Test{
		expectedPayloadSizes: []int{5},
		numForwards:          5,
		forwardedSize:        1,
		maxDelay:             time.Second,
		maxSize:              50,
		timerBehavior: func(mockTimer *tbntime.MockTimer) {
			emptyTimeChan := make(chan time.Time, 1)
			deadlineTimeChan := make(chan time.Time, 1)
			deadlineTimeChan <- time.Now()

			gomock.InOrder(
				mockTimer.EXPECT().Stop().Return(true),
				mockTimer.EXPECT().Reset(1*time.Second).Return(false),
				mockTimer.EXPECT().C().Times(5).Return(emptyTimeChan),
				mockTimer.EXPECT().C().Return(deadlineTimeChan),
				mockTimer.EXPECT().C().Return(emptyTimeChan),
				mockTimer.EXPECT().Stop().Return(false),
			)
		},
	}.run(t)
}

func TestPayloadV2BatcherRunResetsTimer(t *testing.T) {
	batcherV2Test{
		expectedPayloadSizes: []int{5, 1},
		numForwards:          6,
		forwardedSize:        1,
		maxDelay:             time.Second,
		maxSize:              5,
		timerBehavior: func(mockTimer *tbntime.MockTimer) {
			emptyTimeChan := make(chan time.Time, 1)
			deadlineTimeChan := make(chan time.Time, 1)
			deadlineTimeChan <- time.Now()

			gomock.InOrder(
				mockTimer.EXPECT().Stop().Return(true),
				mockTimer.EXPECT().Reset(1*time.Second).Return(false),
				mockTimer.EXPECT().C().Times(5).Return(emptyTimeChan),
				mockTimer.EXPECT().Stop().Return(true),
				mockTimer.EXPECT().C().Times(1).Return(emptyTimeChan),
				mockTimer.EXPECT().Reset(1*time.Second).Return(false),
				mockTimer.EXPECT().C().Return(deadlineTimeChan),
				mockTimer.EXPECT().C().Return(emptyTimeChan),
				mockTimer.EXPECT().Stop().Return(false),
			)
		},
	}.run(t)
}

func TestPayloadV2BatcherRunSendsOnProxyChange(t *testing.T) {
	payload1 := payloadV2OfSize(1)
	payload1.Proxy = ptr.String("p1")
	payload2 := payloadV2OfSize(1)
	payload2.Proxy = ptr.String("p2")

	batcherV2Test{
		expectedCustomPayloads: []*statsapi.PayloadV2{payload1, payload2},
		customPayloads:         []*statsapi.PayloadV2{payload1, payload2},
		maxDelay:               time.Second,
		maxSize:                5,
		timerBehavior: func(mockTimer *tbntime.MockTimer) {
			emptyTimeChan := make(chan time.Time, 1)
			deadlineTimeChan := make(chan time.Time, 1)
			deadlineTimeChan <- time.Now()

			gomock.InOrder(
				mockTimer.EXPECT().Stop().Return(true),
				mockTimer.EXPECT().Reset(1*time.Second).Return(false),
				mockTimer.EXPECT().C().Times(2).Return(emptyTimeChan),

				mockTimer.EXPECT().Reset(1*time.Second).Return(false),
				mockTimer.EXPECT().C().Return(deadlineTimeChan),
				mockTimer.EXPECT().C().Return(emptyTimeChan),
				mockTimer.EXPECT().Stop().Return(false),
			)
		},
	}.run(t)
}

func TestPayloadV2BatcherRunSendsOnProxyVersionChange(t *testing.T) {
	payload1 := payloadV2OfSize(1)
	payload1.ProxyVersion = ptr.String("v1")
	payload2 := payloadV2OfSize(1)
	payload2.ProxyVersion = ptr.String("v2")

	batcherV2Test{
		expectedCustomPayloads: []*statsapi.PayloadV2{payload1, payload2},
		customPayloads:         []*statsapi.PayloadV2{payload1, payload2},
		maxDelay:               time.Second,
		maxSize:                5,
		timerBehavior: func(mockTimer *tbntime.MockTimer) {
			emptyTimeChan := make(chan time.Time, 1)
			deadlineTimeChan := make(chan time.Time, 1)
			deadlineTimeChan <- time.Now()

			gomock.InOrder(
				mockTimer.EXPECT().Stop().Return(true),
				mockTimer.EXPECT().Reset(1*time.Second).Return(false),
				mockTimer.EXPECT().C().Times(2).Return(emptyTimeChan),
				mockTimer.EXPECT().Reset(1*time.Second).Return(false),
				mockTimer.EXPECT().C().Return(deadlineTimeChan),
				mockTimer.EXPECT().C().Return(emptyTimeChan),
				mockTimer.EXPECT().Stop().Return(false),
			)
		},
	}.run(t)
}

func TestPayloadV2BatcherRunSendsMergedProxyData(t *testing.T) {
	payload1 := payloadV2OfSize(1)
	payload1.Proxy = ptr.String("p1")
	payload2 := payloadV2OfSize(1)
	payload2.ProxyVersion = ptr.String("v2")
	payload3 := payloadV2OfSize(1)
	payload3.Proxy = ptr.String("p1")
	payload3.ProxyVersion = ptr.String("v2")
	payload4 := payloadV2OfSize(1)

	combined := payloadV2OfSize(4)
	combined.Proxy = ptr.String("p1")
	combined.ProxyVersion = ptr.String("v2")

	batcherV2Test{
		expectedCustomPayloads: []*statsapi.PayloadV2{combined},
		customPayloads:         []*statsapi.PayloadV2{payload1, payload2, payload3, payload4},
		maxDelay:               time.Second,
		maxSize:                4,
		timerBehavior: func(mockTimer *tbntime.MockTimer) {
			emptyTimeChan := make(chan time.Time, 1)
			deadlineTimeChan := make(chan time.Time, 1)
			deadlineTimeChan <- time.Now()

			gomock.InOrder(
				mockTimer.EXPECT().Stop().Return(true),
				mockTimer.EXPECT().Reset(1*time.Second).Return(false),
				mockTimer.EXPECT().C().Times(4).Return(emptyTimeChan),
				mockTimer.EXPECT().Stop().Return(true),
				mockTimer.EXPECT().C().Times(1).Return(emptyTimeChan),
				mockTimer.EXPECT().Stop().Return(true),
			)
		},
	}.run(t)
}

func TestPayloadV2BatcherRunSendsOnDefaultLimitsChange(t *testing.T) {
	payload1 := payloadV2OfSize(1)
	payload1.Limits = map[string][]float64{
		"default": {0.1, 0.2, 0.4, 0.8},
		"l1":      {0.1, 0.2, 0.4, 0.8},
	}
	payload2 := payloadV2OfSize(1)
	payload2.Limits = map[string][]float64{
		"default": {1, 2, 4, 8},
	}

	batcherV2Test{
		expectedCustomPayloads: []*statsapi.PayloadV2{payload1, payload2},
		customPayloads:         []*statsapi.PayloadV2{payload1, payload2},
		maxDelay:               time.Second,
		maxSize:                5,
		timerBehavior: func(mockTimer *tbntime.MockTimer) {
			emptyTimeChan := make(chan time.Time, 1)
			deadlineTimeChan := make(chan time.Time, 1)
			deadlineTimeChan <- time.Now()

			gomock.InOrder(
				mockTimer.EXPECT().Stop().Return(true),
				mockTimer.EXPECT().Reset(1*time.Second).Return(false),
				mockTimer.EXPECT().C().Times(2).Return(emptyTimeChan),
				mockTimer.EXPECT().Reset(1*time.Second).Return(false),
				mockTimer.EXPECT().C().Return(deadlineTimeChan),
				mockTimer.EXPECT().C().Return(emptyTimeChan),
				mockTimer.EXPECT().Stop().Return(false),
			)
		},
	}.run(t)
}

func TestPayloadV2BatcherRunSendsOnLimitsChange(t *testing.T) {
	payload1 := payloadV2OfSize(1)
	payload1.Limits = map[string][]float64{
		"default": {0.1, 0.2, 0.4, 0.8},
		"l1":      {0.1, 0.2, 0.4, 0.8},
	}
	payload2 := payloadV2OfSize(1)
	payload2.Limits = map[string][]float64{
		"default": {0.1, 0.2, 0.4, 0.8},
		"l1":      {1, 2, 4, 8},
	}

	batcherV2Test{
		expectedCustomPayloads: []*statsapi.PayloadV2{payload1, payload2},
		customPayloads:         []*statsapi.PayloadV2{payload1, payload2},
		maxDelay:               time.Second,
		maxSize:                5,
		timerBehavior: func(mockTimer *tbntime.MockTimer) {
			emptyTimeChan := make(chan time.Time, 1)
			deadlineTimeChan := make(chan time.Time, 1)
			deadlineTimeChan <- time.Now()

			gomock.InOrder(
				mockTimer.EXPECT().Stop().Return(true),
				mockTimer.EXPECT().Reset(1*time.Second).Return(false),
				mockTimer.EXPECT().C().Times(2).Return(emptyTimeChan),
				mockTimer.EXPECT().Reset(1*time.Second).Return(false),
				mockTimer.EXPECT().C().Return(deadlineTimeChan),
				mockTimer.EXPECT().C().Return(emptyTimeChan),
				mockTimer.EXPECT().Stop().Return(false),
			)
		},
	}.run(t)
}

func TestPayloadV2BatcherRunSendsOnLimitsSizeChange(t *testing.T) {
	payload1 := payloadV2OfSize(1)
	payload1.Limits = map[string][]float64{
		"default": {0.1, 0.2, 0.4, 0.8},
	}
	payload2 := payloadV2OfSize(1)
	payload2.Limits = map[string][]float64{
		"default": {0.1, 0.2, 0.4},
	}

	batcherV2Test{
		expectedCustomPayloads: []*statsapi.PayloadV2{payload1, payload2},
		customPayloads:         []*statsapi.PayloadV2{payload1, payload2},
		maxDelay:               time.Second,
		maxSize:                5,
		timerBehavior: func(mockTimer *tbntime.MockTimer) {
			emptyTimeChan := make(chan time.Time, 1)
			deadlineTimeChan := make(chan time.Time, 1)
			deadlineTimeChan <- time.Now()

			gomock.InOrder(
				mockTimer.EXPECT().Stop().Return(true),
				mockTimer.EXPECT().Reset(1*time.Second).Return(false),
				mockTimer.EXPECT().C().Times(2).Return(emptyTimeChan),
				mockTimer.EXPECT().Reset(1*time.Second).Return(false),
				mockTimer.EXPECT().C().Return(deadlineTimeChan),
				mockTimer.EXPECT().C().Return(emptyTimeChan),
				mockTimer.EXPECT().Stop().Return(false),
			)
		},
	}.run(t)
}

func TestPayloadV2BatcherRunSendsMergedLimits(t *testing.T) {
	payload1 := payloadV2OfSize(1)
	payload1.Limits = map[string][]float64{
		"default": {0.1, 0.2, 0.4, 0.8},
	}
	payload2 := payloadV2OfSize(1)
	payload2.Limits = map[string][]float64{
		"l1": {1, 2, 4, 8},
	}
	payload3 := payloadV2OfSize(1)

	combined := payloadV2OfSize(3)
	combined.Limits = map[string][]float64{
		"default": {0.1, 0.2, 0.4, 0.8},
		"l1":      {1, 2, 4, 8},
	}

	batcherV2Test{
		expectedCustomPayloads: []*statsapi.PayloadV2{combined},
		customPayloads:         []*statsapi.PayloadV2{payload1, payload2, payload3},
		maxDelay:               time.Second,
		maxSize:                3,
		timerBehavior: func(mockTimer *tbntime.MockTimer) {
			emptyTimeChan := make(chan time.Time, 1)
			deadlineTimeChan := make(chan time.Time, 1)
			deadlineTimeChan <- time.Now()

			gomock.InOrder(
				mockTimer.EXPECT().Stop().Return(true),
				mockTimer.EXPECT().Reset(1*time.Second).Return(false),
				mockTimer.EXPECT().C().Times(3).Return(emptyTimeChan),
				mockTimer.EXPECT().Stop().Return(true),
				mockTimer.EXPECT().C().Times(1).Return(emptyTimeChan),
				mockTimer.EXPECT().Stop().Return(true),
			)
		},
	}.run(t)
}

func TestPayloadV2BatcherRunSendsOnClose(t *testing.T) {
	batcherV2Test{
		expectedPayloadSizes:  []int{3},
		numForwards:           1,
		forwardedSize:         3,
		closeAfterLastPayload: true,
		maxDelay:              time.Second,
		maxSize:               5,
		timerBehavior: func(mockTimer *tbntime.MockTimer) {
			emptyTimeChan := make(chan time.Time, 1)
			deadlineTimeChan := make(chan time.Time, 1)
			deadlineTimeChan <- time.Now()

			gomock.InOrder(
				mockTimer.EXPECT().Stop().Return(true),
				mockTimer.EXPECT().Reset(1*time.Second).Return(false),
				mockTimer.EXPECT().C().Times(2).Return(emptyTimeChan),
				mockTimer.EXPECT().Stop().Return(false),
			)
		},
	}.run(t)
}

func TestBatchingStatsV2ClientQuery(t *testing.T) {
	ctrl := gomock.NewController(assert.Tracing(t))
	defer ctrl.Finish()

	query := &statsapi.Query{}
	want := &statsapi.QueryResult{}

	mockClient := newMockInternalStatsClient(ctrl)
	mockClient.EXPECT().Query(query).Return(want, nil)

	client := &httpBatchingStatsV2{internalStatsClient: mockClient}

	got, gotErr := client.Query(query)
	assert.Equal(t, got, want)
	assert.Nil(t, gotErr)
}

func TestBatchingStatsV2ClientQueryErr(t *testing.T) {
	ctrl := gomock.NewController(assert.Tracing(t))
	defer ctrl.Finish()

	wantErr := errors.New("Gah!")
	query := &statsapi.Query{}

	mockClient := newMockInternalStatsClient(ctrl)
	mockClient.EXPECT().Query(query).Return(nil, wantErr)

	client := &httpBatchingStatsV2{internalStatsClient: mockClient}

	got, gotErr := client.Query(query)
	assert.Nil(t, got)
	assert.Equal(t, gotErr, wantErr)
}
