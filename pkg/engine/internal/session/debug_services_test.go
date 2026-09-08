package session

import (
	"errors"
	"slices"
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/engine/internal/host"
	"github.com/MontFerret/ferret/v2/pkg/engine/internal/resource"
)

func TestDebugServicesCloseOrderAndErrors(t *testing.T) {
	hooks := &host.SessionHooks{}
	resources := resource.NewManager()
	limiter := NewLimiter(1)
	if err := limiter.Acquire(t.Context(), nil); err != nil {
		t.Fatal(err)
	}

	var calls []string
	hookErr, resourceErr := errors.New("hook"), errors.New("resource")
	hooks.OnClose(func() error {
		calls = append(calls, "hook")

		return hookErr
	})
	for _, name := range []resource.Name{"first", "second"} {
		if err := resources.Own(name, func() error {
			calls = append(calls, string(name))
			if len(limiter.ch) != 1 {
				t.Error("permit released before resource cleanup")
			}

			return resourceErr
		}); err != nil {
			t.Fatal(err)
		}
	}

	services := NewDebugServices(DebugServicesConfig{Hooks: hooks}, limiter, resources)
	err := services.Close()
	if !errors.Is(err, hookErr) || !errors.Is(err, resourceErr) {
		t.Fatalf("cleanup lost errors: %v", err)
	}

	if want := []string{"hook", "second", "first"}; !slices.Equal(calls, want) {
		t.Fatalf("cleanup calls = %v, want %v", calls, want)
	}

	if len(limiter.ch) != 0 {
		t.Fatal("cleanup retained the session permit")
	}
}

func TestDebugServicesReleasePreservesSiblingPermits(t *testing.T) {
	limiter := NewLimiter(2)
	for range 2 {
		if err := limiter.Acquire(t.Context(), nil); err != nil {
			t.Fatal(err)
		}
	}

	services := NewDebugServices(DebugServicesConfig{}, limiter, resource.NewManager())
	if len(limiter.ch) != 2 {
		t.Fatal("constructing debug services released a permit")
	}

	if err := services.Close(); err != nil {
		t.Fatal(err)
	}

	if len(limiter.ch) != 1 {
		t.Fatal("cleanup did not release exactly one session permit")
	}

	if err := limiter.Acquire(t.Context(), nil); err != nil {
		t.Fatal(err)
	}

	// A released slot may already belong to another session. Repeated service
	// cleanup must not consume either sibling's permit.
	if err := services.Close(); err != nil {
		t.Fatal(err)
	}

	if len(limiter.ch) != 2 {
		t.Fatal("repeated cleanup released a sibling session permit")
	}

	limiter.Release()
	limiter.Release()
}
