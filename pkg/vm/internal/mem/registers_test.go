package mem

import (
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

func TestNewRegisterFile_InitializesSlotsWithNone(t *testing.T) {
	rf := NewRegisterFile(4)

	if got, want := len(rf), 4; got != want {
		t.Fatalf("unexpected register size: got %d, want %d", got, want)
	}

	for i := range rf {
		if got := rf[i]; got != runtime.None {
			t.Fatalf("expected register %d to be runtime.None, got %v", i, got)
		}
	}
}

func TestRegisterFileReset_ResetsToNoneAndClearsDirty(t *testing.T) {
	rf := NewRegisterFile(3)
	rf[0] = runtime.NewInt(1)
	rf[1] = runtime.NewInt(2)
	rf[2] = runtime.NewInt(3)

	rf.Reset()

	for i := range rf {
		if got := rf[i]; got != runtime.None {
			t.Fatalf("expected register %d to reset to runtime.None, got %v", i, got)
		}
	}
}

func TestRegisterFileReset_ScrubsWithoutClosingValues(t *testing.T) {
	rf := NewRegisterFile(1)
	closer := newTestCloser("register")
	rf[0] = closer

	rf.Reset()

	if got := closer.closed; got != 0 {
		t.Fatalf("expected reset to avoid closing register value, got %d closes", got)
	}

	if got := rf[0]; got != runtime.None {
		t.Fatalf("expected register to reset to runtime.None, got %v", got)
	}
}
