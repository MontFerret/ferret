package frame

import (
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

func TestPoolGet_NewWindowInitializesSlotsWithNone(t *testing.T) {
	var p Pool
	p.Init(4)

	reg := p.Get(3)
	if got, want := len(reg), 3; got != want {
		t.Fatalf("unexpected register window size: got %d, want %d", got, want)
	}

	for i := range reg {
		if got := reg[i]; got != runtime.None {
			t.Fatalf("expected slot %d to be runtime.None, got %v", i, got)
		}
	}
}

func TestPoolPutGet_ReusedWindowIsNormalizedToNone(t *testing.T) {
	var p Pool
	p.Init(4)

	reg := make([]runtime.Value, 3)
	reg[0] = runtime.NewInt(1)
	reg[1] = runtime.NewInt(2)
	reg[2] = runtime.NewInt(3)

	p.Put(reg)
	reused := p.Get(3)

	if &reused[0] != &reg[0] {
		t.Fatal("expected pooled register window reuse")
	}

	for i := range reused {
		if got := reused[i]; got != runtime.None {
			t.Fatalf("expected reused slot %d to be runtime.None, got %v", i, got)
		}
	}
}
