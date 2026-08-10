package vm

import (
	"context"
	"errors"
	"fmt"
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

func TestSubscribeStreamGroupDoesNotExposePeerCancellation(t *testing.T) {
	setupErr := errors.New("setup failed")
	peerErr := errors.New("peer failed during cancellation")
	failure := &streamGroupSetupObservable{
		subscribeFn: func(context.Context) (runtime.Stream, error) {
			return nil, setupErr
		},
	}
	peer := &streamGroupSetupObservable{
		subscribeFn: func(ctx context.Context) (runtime.Stream, error) {
			<-ctx.Done()
			return nil, fmt.Errorf("peer setup: %w", ctx.Err())
		},
	}
	peerWithError := &streamGroupSetupObservable{
		subscribeFn: func(ctx context.Context) (runtime.Stream, error) {
			<-ctx.Done()
			return nil, errors.Join(peerErr, fmt.Errorf("peer setup: %w", ctx.Err()))
		},
	}

	_, err := subscribeStreamGroup(
		context.Background(),
		runtime.NewArrayWith(failure, peer, peerWithError),
		runtime.NewArrayWith(runtime.NewString("failure"), runtime.NewString("peer"), runtime.NewString("peer-error")),
		runtime.NewArrayWith(runtime.None, runtime.None, runtime.None),
	)
	if !errors.Is(err, setupErr) {
		t.Fatalf("expected setup failure, got %v", err)
	}

	if errors.Is(err, context.Canceled) {
		t.Fatalf("internal fail-fast cancellation escaped setup: %v", err)
	}

	if !errors.Is(err, peerErr) {
		t.Fatalf("genuine peer error was discarded: %v", err)
	}
}

func TestSubscribeStreamGroupPreservesExternalCancellation(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	waiter := &streamGroupSetupObservable{
		subscribeFn: func(ctx context.Context) (runtime.Stream, error) {
			cancel()
			<-ctx.Done()

			return nil, ctx.Err()
		},
	}

	_, err := subscribeStreamGroup(
		ctx,
		runtime.NewArrayWith(waiter),
		runtime.NewArrayWith(runtime.NewString("waiter")),
		runtime.NewArrayWith(runtime.None),
	)
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("expected external cancellation, got %v", err)
	}
}

func TestSubscribeStreamGroupOrdersSetupAndCloseErrorsByDeclaration(t *testing.T) {
	firstSetupErr := errors.New("first setup")
	secondSetupErr := errors.New("second setup")
	firstCloseErr := errors.New("first close")
	secondCloseErr := errors.New("second close")
	firstStream := &streamGroupSetupStream{closeErr: firstCloseErr}
	secondStream := &streamGroupSetupStream{closeErr: secondCloseErr}
	secondStarted := make(chan struct{})
	first := &streamGroupSetupObservable{
		subscribeFn: func(context.Context) (runtime.Stream, error) {
			<-secondStarted

			return firstStream, firstSetupErr
		},
	}
	second := &streamGroupSetupObservable{
		subscribeFn: func(context.Context) (runtime.Stream, error) {
			close(secondStarted)
			return secondStream, secondSetupErr
		},
	}

	_, err := subscribeStreamGroup(
		context.Background(),
		runtime.NewArrayWith(first, second),
		runtime.NewArrayWith(runtime.NewString("first"), runtime.NewString("second")),
		runtime.NewArrayWith(runtime.None, runtime.None),
	)
	want := "first setup\nsecond setup\nfirst close\nsecond close"
	if err == nil || err.Error() != want {
		t.Fatalf("unexpected ordered setup error:\nwant: %q\n got: %q", want, err)
	}
}

func TestSubscribeStreamGroupSetupPanicUsesDeclarationOrder(t *testing.T) {
	secondStarted := make(chan struct{})
	first := &streamGroupSetupObservable{
		subscribeFn: func(context.Context) (runtime.Stream, error) {
			<-secondStarted
			panic("first setup panic")
		},
	}
	second := &streamGroupSetupObservable{
		subscribeFn: func(context.Context) (runtime.Stream, error) {
			close(secondStarted)
			panic("second setup panic")
		},
	}

	got := capturePanic(func() {
		_, _ = subscribeStreamGroup(
			context.Background(),
			runtime.NewArrayWith(first, second),
			runtime.NewArrayWith(runtime.NewString("first"), runtime.NewString("second")),
			runtime.NewArrayWith(runtime.None, runtime.None),
		)
	})
	if got != "first setup panic" {
		t.Fatalf("expected first setup panic, got %v", got)
	}
}

func TestSubscribeStreamGroupSetupPanicClosesEveryPartialStream(t *testing.T) {
	firstStream := &streamGroupSetupStream{closePanic: "first close panic"}
	thirdStream := &streamGroupSetupStream{}
	first := &streamGroupSetupObservable{
		subscribeFn: func(context.Context) (runtime.Stream, error) {
			return firstStream, nil
		},
	}
	second := &streamGroupSetupObservable{
		subscribeFn: func(context.Context) (runtime.Stream, error) {
			panic("setup panic")
		},
	}
	third := &streamGroupSetupObservable{
		subscribeFn: func(context.Context) (runtime.Stream, error) {
			return thirdStream, nil
		},
	}

	got := capturePanic(func() {
		_, _ = subscribeStreamGroup(
			context.Background(),
			runtime.NewArrayWith(first, second, third),
			runtime.NewArrayWith(
				runtime.NewString("first"),
				runtime.NewString("second"),
				runtime.NewString("third"),
			),
			runtime.NewArrayWith(runtime.None, runtime.None, runtime.None),
		)
	})
	if got != "setup panic" {
		t.Fatalf("expected setup panic precedence, got %v", got)
	}

	if got := firstStream.closeCount.Load(); got != 1 {
		t.Fatalf("expected first stream close once, got %d", got)
	}

	if got := thirdStream.closeCount.Load(); got != 1 {
		t.Fatalf("expected cleanup to continue after close panic, got %d", got)
	}
}

func TestSubscribeStreamGroupClosePanicUsesDeclarationOrder(t *testing.T) {
	firstStream := &streamGroupSetupStream{closePanic: "first close panic"}
	secondStream := &streamGroupSetupStream{closePanic: "second close panic"}
	setupErr := errors.New("setup failed")
	first := &streamGroupSetupObservable{
		subscribeFn: func(context.Context) (runtime.Stream, error) {
			return firstStream, setupErr
		},
	}
	second := &streamGroupSetupObservable{
		subscribeFn: func(context.Context) (runtime.Stream, error) {
			return secondStream, setupErr
		},
	}

	got := capturePanic(func() {
		_, _ = subscribeStreamGroup(
			context.Background(),
			runtime.NewArrayWith(first, second),
			runtime.NewArrayWith(runtime.NewString("first"), runtime.NewString("second")),
			runtime.NewArrayWith(runtime.None, runtime.None),
		)
	})
	if got != "first close panic" {
		t.Fatalf("expected first close panic, got %v", got)
	}

	if got := secondStream.closeCount.Load(); got != 1 {
		t.Fatalf("expected second stream to close after first panic, got %d", got)
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

func capturePanic(fn func()) (recovered any) {
	defer func() {
		recovered = recover()
	}()
	fn()

	return nil
}
