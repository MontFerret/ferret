package mem

import (
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

func TestNewScratch_InitializesSlotsWithNone(t *testing.T) {
	s := NewScratch(3)

	if got, want := len(s.Params), 3; got != want {
		t.Fatalf("unexpected params size: got %d, want %d", got, want)
	}

	for i := range s.Params {
		if got := s.Params[i]; got != runtime.None {
			t.Fatalf("expected slot %d to be runtime.None, got %v", i, got)
		}
	}
}

func TestScratchResizeParams_GrowBeyondCapacityInitializesNewSlotsWithNone(t *testing.T) {
	s := NewScratch(1)
	s.Params[0] = runtime.NewInt(7)

	s.ResizeParams(4)

	if got, want := len(s.Params), 4; got != want {
		t.Fatalf("unexpected params size after growth: got %d, want %d", got, want)
	}

	if got, want := s.Params[0], runtime.NewInt(7); got != want {
		t.Fatalf("unexpected preserved slot value: got %v, want %v", got, want)
	}

	for i := 1; i < len(s.Params); i++ {
		if s.Params[i] == nil {
			t.Fatalf("expected slot %d to not be nil", i)
		}

		if got := s.Params[i]; got != runtime.None {
			t.Fatalf("expected slot %d to be runtime.None, got %v", i, got)
		}
	}
}

func TestScratchResizeParams_GrowWithinCapacityResetsExposedSlotsToNone(t *testing.T) {
	s := NewScratch(3)
	s.Params[1] = runtime.NewInt(11)
	s.Params[2] = runtime.NewInt(13)

	s.ResizeParams(1)
	s.ResizeParams(3)

	if got, want := len(s.Params), 3; got != want {
		t.Fatalf("unexpected params size after regrowth: got %d, want %d", got, want)
	}

	for i := 1; i < len(s.Params); i++ {
		if s.Params[i] == nil {
			t.Fatalf("expected slot %d to not be nil", i)
		}

		if got := s.Params[i]; got != runtime.None {
			t.Fatalf("expected slot %d to be runtime.None after regrowth, got %v", i, got)
		}
	}
}
