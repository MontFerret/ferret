package mem

import (
	"context"
	"io"
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
	"github.com/MontFerret/ferret/v2/pkg/vm/internal/data"
)

func TestCanTrackValueRejectsCommonScalars(t *testing.T) {
	values := []runtime.Value{
		nil,
		runtime.None,
		runtime.True,
		runtime.NewInt(1),
		runtime.NewFloat(1.5),
		runtime.NewString("x"),
		runtime.NewArray(0),
		runtime.NewObject(),
		data.NewDataSet(false),
		data.NewFastObject(nil, 0),
		data.NewKV(runtime.NewString("k"), runtime.NewInt(1)),
		data.NewAggregateKey(runtime.NewString("g"), 0),
		data.NewIterator(&plainRuntimeIterator{}),
	}

	for _, val := range values {
		if CanTrackValue(val) {
			t.Fatalf("expected %T to skip ownership tracking", val)
		}

		if _, _, ok := ResourceKeyOf(val); ok {
			t.Fatalf("expected %T to have no resource key", val)
		}
	}
}

func TestIteratorTrackingSkipsPlainIteratorsButTracksClosableOnes(t *testing.T) {
	plain := data.NewIterator(&plainRuntimeIterator{})
	if CanTrackValue(plain) {
		t.Fatal("expected plain iterator wrapper to skip ownership tracking")
	}

	if _, _, ok := ResourceKeyOf(plain); ok {
		t.Fatal("expected plain iterator wrapper to have no resource key")
	}

	var owned OwnedResources
	var deferred DeferredClosers

	owned.Track(plain)
	owned.Discard(plain, &deferred)

	if !owned.Empty() {
		t.Fatal("expected plain iterator wrapper not to enter owned resources")
	}

	if !deferred.Empty() {
		t.Fatal("expected plain iterator wrapper not to enter deferred closers")
	}

	src := &closableRuntimeIterator{}
	closable := data.WrapIterator(src)
	if !CanTrackValue(closable) {
		t.Fatal("expected closable iterator wrapper to remain ownership-trackable")
	}

	key, resolved, ok := ResourceKeyOf(closable)
	if !ok {
		t.Fatal("expected closable iterator wrapper to resolve a resource key")
	}

	owned.TrackResolved(key, resolved)
	if !owned.OwnsKey(key) {
		t.Fatal("expected closable iterator wrapper to be tracked as owned")
	}

	owned.Discard(closable, &deferred)
	if deferred.Empty() {
		t.Fatal("expected closable iterator wrapper to be deferred on discard")
	}

	if err := deferred.CloseAll(); err != nil {
		t.Fatalf("close deferred closers: %v", err)
	}

	if src.closed != 1 {
		t.Fatalf("expected wrapped closable iterator source to close once, got %d", src.closed)
	}
}

type plainRuntimeIterator struct{}

func (*plainRuntimeIterator) Next(context.Context) (runtime.Value, runtime.Value, error) {
	return runtime.None, runtime.None, io.EOF
}

type closableRuntimeIterator struct {
	closed int
}

func (*closableRuntimeIterator) Next(context.Context) (runtime.Value, runtime.Value, error) {
	return runtime.None, runtime.None, io.EOF
}

func (it *closableRuntimeIterator) Close() error {
	it.closed++
	return nil
}

func TestOwnedResourcesResolvedHelpers(t *testing.T) {
	owned := OwnedResources{}
	closer := newTestCloser("tracked")
	key, resolved, ok := ResourceKeyOf(closer)
	if !ok {
		t.Fatal("expected closer to resolve to a resource key")
	}

	owned.TrackResolved(key, resolved)
	if !owned.OwnsKey(key) {
		t.Fatal("expected resolved key to be owned")
	}

	if !owned.ExtractByKey(key) {
		t.Fatal("expected extract by key to succeed")
	}

	if owned.OwnsKey(key) {
		t.Fatal("expected extract by key to retire ownership")
	}
}
