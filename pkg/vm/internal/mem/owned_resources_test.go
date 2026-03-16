package mem

import (
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

func TestOwnedResourcesExtractManyPreservesDuplicateAliasCounts(t *testing.T) {
	owned := OwnedResources{}
	closer := newTestCloser("dup")

	owned.Track(closer)
	owned.Track(closer)

	var dst OwnedResources
	owned.ExtractMany([]runtime.Value{closer, closer}, &dst)

	if got, want := dst.closers[closer], 2; got != want {
		t.Fatalf("unexpected transferred alias count: got %d, want %d", got, want)
	}

	if owned.closers != nil {
		t.Fatalf("expected source ownership to be empty after transfer, got %+v", owned.closers)
	}
}

func TestOwnedResourcesExtractManyRemovesTransferredClosersFromSourceDrain(t *testing.T) {
	source := OwnedResources{}
	transferred := newTestCloser("transferred")
	discarded := newTestCloser("discarded")

	source.Track(transferred)
	source.Track(transferred)
	source.Track(discarded)

	var dst OwnedResources
	source.ExtractMany([]runtime.Value{transferred, transferred}, &dst)

	deferred := DeferredClosers{}
	source.DrainTo(&deferred)

	if got, want := len(deferred.closers), 1; got != want {
		t.Fatalf("expected only discarded closer to remain deferred, got %d", got)
	}

	if got, want := dst.closers[transferred], 2; got != want {
		t.Fatalf("unexpected transferred alias count after drain: got %d, want %d", got, want)
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

func TestOwnedResourcesReleaseIsTerminalForDuplicateAliases(t *testing.T) {
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
		t.Fatal("expected release to retire ownership for all aliases")
	}

	deferred := DeferredClosers{}
	owned.Discard(closer, &deferred)
	owned.DrainTo(&deferred)

	if got := len(deferred.closers); got != 0 {
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
