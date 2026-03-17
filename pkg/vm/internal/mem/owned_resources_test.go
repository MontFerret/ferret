package mem

import (
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

func TestOwnedResourcesTrackIsIdempotent(t *testing.T) {
	owned := OwnedResources{}
	closer := newTestCloser("dup")

	owned.Track(closer)
	owned.Track(closer)

	if !owned.Owns(closer) {
		t.Fatal("expected tracked closer to remain owned")
	}

	if got, want := len(owned.closers), 1; got != want {
		t.Fatalf("expected one owned closer after duplicate tracking, got %d", got)
	}
}

func TestOwnedResourcesExtractManyTransfersUniqueClosers(t *testing.T) {
	source := OwnedResources{}
	transferred := newTestCloser("transferred")
	discarded := newTestCloser("discarded")

	source.Track(transferred)
	source.Track(transferred)
	source.Track(discarded)

	var dst OwnedResources
	source.ExtractMany([]runtime.Value{transferred, transferred}, &dst)

	if !dst.Owns(transferred) {
		t.Fatal("expected transferred closer to become owned by destination")
	}

	if got, want := len(dst.closers), 1; got != want {
		t.Fatalf("expected one transferred closer in destination, got %d", got)
	}

	deferred := DeferredClosers{}
	source.DrainTo(&deferred)

	if got, want := deferred.set.Len(), 1; got != want {
		t.Fatalf("expected only discarded closer to remain deferred, got %d", got)
	}

	if err := deferred.CloseAll(); err != nil {
		t.Fatalf("expected deferred close to succeed, got %v", err)
	}

	if got := transferred.closed; got != 0 {
		t.Fatalf("expected transferred closer to remain open, got %d closes", got)
	}

	if got := discarded.closed; got != 1 {
		t.Fatalf("expected discarded closer to close once, got %d closes", got)
	}
}

func TestOwnedResourcesReleaseIsTerminalForStaleAliases(t *testing.T) {
	owned := OwnedResources{}
	closer := newTestCloser("dup")

	owned.Track(closer)
	owned.Track(closer)

	released, ok := owned.Release(closer)
	if !ok {
		t.Fatal("expected release to succeed for tracked closer")
	}

	if released != closer {
		t.Fatalf("expected release to return original closer, got %v", released)
	}

	if owned.Owns(closer) {
		t.Fatal("expected release to retire ownership for the closer")
	}

	deferred := DeferredClosers{}
	owned.Discard(closer, &deferred)
	owned.DrainTo(&deferred)

	if got := deferred.set.Len(); got != 0 {
		t.Fatalf("expected release to keep stale aliases out of deferred cleanup, got %d closers", got)
	}

	if err := released.Close(); err != nil {
		t.Fatalf("expected released closer to close successfully, got %v", err)
	}

	if err := deferred.CloseAll(); err != nil {
		t.Fatalf("expected deferred close to succeed, got %v", err)
	}

	if got := closer.closed; got != 1 {
		t.Fatalf("expected closer to close exactly once, got %d closes", got)
	}
}
