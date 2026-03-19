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

type sliceValueCloser []int

func (sliceValueCloser) Close() error {
	return nil
}

func (v sliceValueCloser) String() string {
	return "slice-closer"
}

func (v sliceValueCloser) Hash() uint64 {
	return uint64(len(v))
}

func (v sliceValueCloser) Copy() runtime.Value {
	cp := append(sliceValueCloser(nil), v...)
	return cp
}

type sliceResourceCloser []int

func (sliceResourceCloser) Close() error {
	return nil
}

func (v sliceResourceCloser) ResourceID() uint64 {
	return 42
}

func (v sliceResourceCloser) String() string {
	return "slice-resource"
}

func (v sliceResourceCloser) Hash() uint64 {
	return uint64(len(v))
}

func (v sliceResourceCloser) Copy() runtime.Value {
	cp := append(sliceResourceCloser(nil), v...)
	return cp
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

func TestResourceKeyOfRejectsNonComparablePlainClosers(t *testing.T) {
	closer := sliceValueCloser{1, 2, 3}

	if !CanTrackValue(closer) {
		t.Fatal("expected non-comparable plain closer value to reach ResourceKeyOf checks")
	}

	key, resolved, ok := ResourceKeyOf(closer)
	if ok {
		t.Fatalf("expected non-comparable plain closer to have no resource key, got %#v", key)
	}

	if resolved != nil {
		t.Fatalf("expected no resolved closer for non-comparable plain closer, got %v", resolved)
	}
}

func TestOwnedResourcesSkipNonComparablePlainClosers(t *testing.T) {
	closer := sliceValueCloser{1, 2, 3}
	var owned OwnedResources
	var deferred DeferredClosers

	owned.Track(closer)
	if owned.Owns(closer) {
		t.Fatal("expected non-comparable plain closer not to be tracked as owned")
	}

	if owned.Extract(closer) {
		t.Fatal("expected extract to report non-comparable plain closer as untracked")
	}

	if released, ok := owned.Release(closer); ok || released != nil {
		t.Fatalf("expected release to skip untracked closer, got closer=%v ok=%v", released, ok)
	}

	owned.Discard(closer, &deferred)
	owned.DrainTo(&deferred)

	if !owned.Empty() {
		t.Fatal("expected owned resources to remain empty for untracked closer")
	}

	if !deferred.Empty() {
		t.Fatal("expected deferred closers to remain empty for untracked closer")
	}
}

func TestResourceKeyOfTracksNonComparableResourcesByID(t *testing.T) {
	resource := sliceResourceCloser{1, 2, 3}

	key, resolved, ok := ResourceKeyOf(resource)
	if !ok {
		t.Fatal("expected runtime.Resource path to track non-comparable resource")
	}

	if key != (ResourceKey{ID: resource.ResourceID()}) {
		t.Fatalf("expected resource key by ID, got %#v", key)
	}

	got, ok := resolved.(sliceResourceCloser)
	if !ok {
		t.Fatalf("expected resolved closer to keep sliceResourceCloser type, got %T", resolved)
	}

	if got.ResourceID() != resource.ResourceID() {
		t.Fatalf("expected resolved closer to keep resource ID %d, got %d", resource.ResourceID(), got.ResourceID())
	}

	var owned OwnedResources
	owned.Track(resource)

	if !owned.Owns(resource) {
		t.Fatal("expected non-comparable runtime.Resource to remain trackable")
	}
}
