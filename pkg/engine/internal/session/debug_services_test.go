package session

import (
	"errors"
	"slices"
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/engine/internal/host"
	"github.com/MontFerret/ferret/v2/pkg/engine/internal/resource"
	"github.com/MontFerret/ferret/v2/pkg/vm"
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

	release := NewPermitRelease(limiter, nil)
	services := &DebugServices{
		Hooks: hooks, Resources: resources,
		ReleasePermit: func(instance *vm.VM) {
			calls = append(calls, "permit")
			if instance != nil {
				t.Error("debug services attempted to return a borrowed VM")
			}

			release(instance)
		},
	}
	err := services.Close()
	if !errors.Is(err, hookErr) || !errors.Is(err, resourceErr) {
		t.Fatalf("cleanup lost errors: %v", err)
	}

	if want := []string{"hook", "second", "first", "permit"}; !slices.Equal(calls, want) {
		t.Fatalf("cleanup calls = %v, want %v", calls, want)
	}

	if len(limiter.ch) != 0 || services.ReleasePermit != nil {
		t.Fatal("cleanup retained the session permit")
	}
}
