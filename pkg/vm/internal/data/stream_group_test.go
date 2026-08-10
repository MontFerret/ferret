package data

import (
	"context"
	"errors"
	"io"
	"testing"
	"time"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

func TestStreamGroupIteratorUsesOneFixedDeadline(t *testing.T) {
	stream := newStreamGroupTestStream(1)
	iterator := NewStreamGroupValue([]runtime.Stream{stream}).(*StreamGroupValue).Iterate(runtime.Duration(50 * time.Millisecond))

	started := time.Now()
	time.AfterFunc(30*time.Millisecond, func() {
		stream.messages <- runtime.NewValueMessage(runtime.NewString("event"))
	})
	if err := iterator.Next(context.Background()); err != nil {
		t.Fatalf("first event: %v", err)
	}

	secondStarted := time.Now()
	err := iterator.Next(context.Background())
	if !errors.Is(err, runtime.ErrTimeout) {
		t.Fatalf("expected fixed deadline timeout, got %v", err)
	}
	if secondElapsed := time.Since(secondStarted); secondElapsed >= 40*time.Millisecond {
		t.Fatalf("second wait restarted the deadline: %s", secondElapsed)
	}
	if total := time.Since(started); total >= 70*time.Millisecond {
		t.Fatalf("group deadline exceeded expected bound: %s", total)
	}

	if err := iterator.(interface{ Close() error }).Close(); err != nil {
		t.Fatalf("close iterator: %v", err)
	}
}

func TestStreamGroupIteratorReportsClosedArmAndContinues(t *testing.T) {
	first := newStreamGroupTestStream(0)
	second := newStreamGroupTestStream(1)
	close(first.messages)

	iterator := NewStreamGroupValue([]runtime.Stream{first, second}).(*StreamGroupValue).
		Iterate(runtime.Duration(time.Second)).(*StreamGroupIterator)

	if err := iterator.Next(context.Background()); err != nil {
		t.Fatalf("closed arm notification: %v", err)
	}
	if iterator.Key() != runtime.None {
		t.Fatalf("expected None closure key, got %v", iterator.Key())
	}
	if got := iterator.Value(); got != runtime.NewInt(0) {
		t.Fatalf("expected closed arm index 0, got %v", got)
	}
	if got := first.closeCount.Load(); got != 1 {
		t.Fatalf("expected exhausted arm to close once, got %d", got)
	}
	if got := second.closeCount.Load(); got != 0 {
		t.Fatalf("remaining arm closed before group completion: %d", got)
	}

	second.messages <- runtime.NewValueMessage(runtime.NewString("second"))
	if err := iterator.Next(context.Background()); err != nil {
		t.Fatalf("remaining arm event: %v", err)
	}
	if got := iterator.Key(); got != runtime.NewInt(1) {
		t.Fatalf("expected remaining arm index 1, got %v", got)
	}
	if got := iterator.Value(); got != runtime.NewString("second") {
		t.Fatalf("expected remaining arm event, got %v", got)
	}

	if err := iterator.Close(); err != nil {
		t.Fatalf("close iterator: %v", err)
	}
	if got := first.closeCount.Load(); got != 1 {
		t.Fatalf("group cleanup closed exhausted arm %d times", got)
	}
	if got := second.closeCount.Load(); got != 1 {
		t.Fatalf("expected remaining arm cleanup, got %d closes", got)
	}
	select {
	case <-iterator.done:
	default:
		t.Fatal("iterator workers remained active after close")
	}
}

func TestStreamGroupIteratorReportsEveryClosureBeforeEOF(t *testing.T) {
	first := newStreamGroupTestStream(0)
	second := newStreamGroupTestStream(0)
	close(first.messages)

	iterator := NewStreamGroupValue([]runtime.Stream{first, second}).(*StreamGroupValue).
		Iterate(runtime.Duration(time.Second)).(*StreamGroupIterator)

	if err := iterator.Next(context.Background()); err != nil {
		t.Fatalf("first closure: %v", err)
	}
	if iterator.Key() != runtime.None || iterator.Value() != runtime.NewInt(0) {
		t.Fatalf("unexpected first closure notification: key=%v value=%v", iterator.Key(), iterator.Value())
	}

	close(second.messages)
	if err := iterator.Next(context.Background()); err != nil {
		t.Fatalf("second closure: %v", err)
	}
	if iterator.Key() != runtime.None || iterator.Value() != runtime.NewInt(1) {
		t.Fatalf("unexpected second closure notification: key=%v value=%v", iterator.Key(), iterator.Value())
	}

	if err := iterator.Next(context.Background()); !errors.Is(err, io.EOF) {
		t.Fatalf("expected EOF after every closure notification, got %v", err)
	}
	if err := iterator.Close(); err != nil {
		t.Fatalf("close iterator: %v", err)
	}
	if got := first.closeCount.Load(); got != 1 {
		t.Fatalf("expected first arm to close once, got %d", got)
	}
	if got := second.closeCount.Load(); got != 1 {
		t.Fatalf("expected second arm to close once, got %d", got)
	}
}

func TestStreamGroupIteratorSuppressesClosureFromCompletedArm(t *testing.T) {
	first := newStreamGroupTestStream(1)
	second := newStreamGroupTestStream(1)
	first.messages <- runtime.NewValueMessage(runtime.NewString("first"))
	close(first.messages)

	iterator := NewStreamGroupValue([]runtime.Stream{first, second}).(*StreamGroupValue).
		Iterate(runtime.Duration(time.Second)).(*StreamGroupIterator)
	if err := iterator.Next(context.Background()); err != nil {
		t.Fatalf("first event: %v", err)
	}
	if err := iterator.ArmDone(0); err != nil {
		t.Fatalf("complete first arm: %v", err)
	}

	second.messages <- runtime.NewValueMessage(runtime.NewString("second"))
	if err := iterator.Next(context.Background()); err != nil {
		t.Fatalf("second event after completed-arm closure: %v", err)
	}
	if got := iterator.Key(); got != runtime.NewInt(1) {
		t.Fatalf("expected second arm, got key %v", got)
	}

	if err := iterator.Close(); err != nil {
		t.Fatalf("close iterator: %v", err)
	}
	if got := first.closeCount.Load(); got != 1 {
		t.Fatalf("expected completed arm to close once, got %d", got)
	}
	if got := second.closeCount.Load(); got != 1 {
		t.Fatalf("expected remaining arm to close once, got %d", got)
	}
}

func TestStreamGroupIteratorIgnoresCompletedArmErrors(t *testing.T) {
	first := newStreamGroupTestStream(2)
	second := newStreamGroupTestStream(1)
	first.messages <- runtime.NewValueMessage(runtime.NewString("first"))
	first.messages <- runtime.NewErrorMessage(errors.New("late error"))

	iterator := NewStreamGroupValue([]runtime.Stream{first, second}).(*StreamGroupValue).
		Iterate(runtime.Duration(time.Second)).(*StreamGroupIterator)
	if err := iterator.Next(context.Background()); err != nil {
		t.Fatalf("first event: %v", err)
	}
	if got := iterator.Key(); got != runtime.NewInt(0) {
		t.Fatalf("expected first arm, got key %v", got)
	}
	if err := iterator.ArmDone(0); err != nil {
		t.Fatalf("complete first arm: %v", err)
	}

	second.messages <- runtime.NewValueMessage(runtime.NewString("second"))
	if err := iterator.Next(context.Background()); err != nil {
		t.Fatalf("second event after completed-arm error: %v", err)
	}
	if got := iterator.Key(); got != runtime.NewInt(1) {
		t.Fatalf("expected second arm, got key %v", got)
	}

	if err := iterator.Close(); err != nil {
		t.Fatalf("close iterator: %v", err)
	}
	if got := first.closeCount.Load(); got != 1 {
		t.Fatalf("expected completed arm to close once, got %d", got)
	}
	if got := second.closeCount.Load(); got != 1 {
		t.Fatalf("expected remaining arm to close once, got %d", got)
	}
	if err := iterator.Close(); err != nil {
		t.Fatalf("repeat close: %v", err)
	}
	if got := first.closeCount.Load(); got != 1 {
		t.Fatalf("repeated close closed first arm %d times", got)
	}
}

func TestStreamGroupIteratorRejectsInvalidArmIndexesWithoutClosingArm(t *testing.T) {
	stream := newStreamGroupTestStream(1)
	iterator := NewStreamGroupValue([]runtime.Stream{stream}).(*StreamGroupValue).
		Iterate(runtime.Duration(time.Second)).(*StreamGroupIterator)
	t.Cleanup(func() {
		if err := iterator.Close(); err != nil {
			t.Errorf("close iterator: %v", err)
		}
	})

	for _, index := range []runtime.Int{-1, 1, runtime.NewInt64(1 << 32)} {
		err := iterator.ArmDone(index)
		if !errors.Is(err, runtime.ErrInvalidArgument) {
			t.Fatalf("expected invalid argument for arm index %d, got %v", index, err)
		}

		if !iterator.active[0] || iterator.activeCount != 1 {
			t.Fatalf("arm index %d changed active state: active=%v count=%d", index, iterator.active, iterator.activeCount)
		}

		if got := stream.closeCount.Load(); got != 0 {
			t.Fatalf("arm index %d closed stream %d times", index, got)
		}
	}
}

func TestStreamGroupIteratorPropagatesReadPanicByDeclarationOrder(t *testing.T) {
	first := newStreamGroupTestStream(0)
	first.readPanic = "first read panic"
	second := newStreamGroupTestStream(0)
	second.readPanic = "second read panic"
	iterator := NewStreamGroupValue([]runtime.Stream{first, second}).(*StreamGroupValue).
		Iterate(runtime.Duration(time.Second)).(*StreamGroupIterator)

	got := captureStreamGroupPanic(func() {
		_ = iterator.Next(context.Background())
	})
	if got != "first read panic" {
		t.Fatalf("expected first read panic, got %v", got)
	}

	if err := iterator.Close(); err != nil {
		t.Fatalf("close iterator: %v", err)
	}

	if first.closeCount.Load() != 1 || second.closeCount.Load() != 1 {
		t.Fatalf("expected both streams closed once, got first=%d second=%d", first.closeCount.Load(), second.closeCount.Load())
	}
}

func TestStreamGroupIteratorClosePropagatesPendingReadPanic(t *testing.T) {
	stream := newStreamGroupTestStream(0)
	stream.readPanic = "read panic"
	iterator := NewStreamGroupValue([]runtime.Stream{stream}).(*StreamGroupValue).
		Iterate(runtime.Duration(time.Second)).(*StreamGroupIterator)
	iterator.start(context.Background())
	<-iterator.done

	got := captureStreamGroupPanic(func() {
		_ = iterator.Close()
	})
	if got != "read panic" {
		t.Fatalf("expected pending read panic, got %v", got)
	}

	if got := stream.closeCount.Load(); got != 1 {
		t.Fatalf("expected stream close once, got %d", got)
	}
}

func TestStreamGroupIteratorPropagatesActiveArmError(t *testing.T) {
	stream := newStreamGroupTestStream(1)
	want := errors.New("stream failed")
	stream.messages <- runtime.NewErrorMessage(want)
	iterator := NewStreamGroupValue([]runtime.Stream{stream}).(*StreamGroupValue).
		Iterate(runtime.Duration(time.Second)).(*StreamGroupIterator)

	if err := iterator.Next(context.Background()); !errors.Is(err, want) {
		t.Fatalf("expected active arm error, got %v", err)
	}
	if err := iterator.Close(); err != nil {
		t.Fatalf("close iterator: %v", err)
	}
}

func captureStreamGroupPanic(fn func()) (recovered any) {
	defer func() {
		recovered = recover()
	}()
	fn()

	return nil
}
