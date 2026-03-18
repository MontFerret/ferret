package mem

import (
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
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
