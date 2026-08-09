package vm

import (
	"context"
	"errors"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
	"github.com/MontFerret/ferret/v2/test/spec/mock"
)

func TestSubscribeStreamGroupEstablishesSubscriptionsConcurrently(t *testing.T) {
	var started atomic.Int32
	gate := make(chan struct{})
	gateOnce := &sync.Once{}
	first := newStreamGroupSetupObservable(&started, gate, gateOnce, 2, nil)
	second := newStreamGroupSetupObservable(&started, gate, gateOnce, 2, nil)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	streams, err := subscribeStreamGroup(
		ctx,
		runtime.NewArrayWith(first, second),
		runtime.NewArrayWith(runtime.NewString("first"), runtime.NewString("second")),
		runtime.NewArrayWith(runtime.None, runtime.None),
	)
	if err != nil {
		t.Fatalf("expected concurrent stream setup, got %v", err)
	}
	if got := started.Load(); got != 2 {
		t.Fatalf("expected both subscriptions to start, got %d", got)
	}
	if err := closeStreamGroupSetup(streams); err != nil {
		t.Fatalf("close streams: %v", err)
	}
}

func TestSubscribeStreamGroupClosesPartialSetup(t *testing.T) {
	success := mock.NewTriggerObservable()
	var started atomic.Int32
	gate := make(chan struct{})
	close(gate)
	failure := newStreamGroupSetupObservable(&started, gate, &sync.Once{}, 2, errors.New("setup failed"))

	_, err := subscribeStreamGroup(
		context.Background(),
		runtime.NewArrayWith(success, failure),
		runtime.NewArrayWith(runtime.NewString("success"), runtime.NewString("failure")),
		runtime.NewArrayWith(runtime.None, runtime.None),
	)
	if !errors.Is(err, failure.fail) {
		t.Fatalf("expected setup failure, got %v", err)
	}
	if got := success.CloseCount(); got != 1 {
		t.Fatalf("expected partial subscription to close once, got %d", got)
	}
}

func TestSubscribeStreamGroupValidatesAllDescriptorsBeforeSetup(t *testing.T) {
	success := mock.NewTriggerObservable()
	_, err := subscribeStreamGroup(
		context.Background(),
		runtime.NewArrayWith(success, runtime.NewString("invalid")),
		runtime.NewArrayWith(runtime.NewString("success"), runtime.NewString("failure")),
		runtime.NewArrayWith(runtime.None, runtime.None),
	)
	if err == nil {
		t.Fatal("expected descriptor validation failure")
	}
	if got := success.SubscribeCount(); got != 0 {
		t.Fatalf("expected validation before setup, got %d subscriptions", got)
	}
}
