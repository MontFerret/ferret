package mem

import (
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

func TestNewRegisterFile_InitializesSlotsWithNone(t *testing.T) {
	rf := NewRegisterFile(4)

	if got, want := rf.Size(), 4; got != want {
		t.Fatalf("unexpected register size: got %d, want %d", got, want)
	}

	for i := range rf.Values {
		if got := rf.Values[i]; got != runtime.None {
			t.Fatalf("expected register %d to be runtime.None, got %v", i, got)
		}
	}
}

func TestRegisterFileReset_ResetsToNoneAndClearsDirty(t *testing.T) {
	rf := NewRegisterFile(3)
	rf.Set(0, runtime.NewInt(1))
	rf.Set(1, runtime.NewInt(2))
	rf.Set(2, runtime.NewInt(3))
	rf.MarkDirty()

	rf.Reset()

	if rf.IsDirty() {
		t.Fatal("expected register file to be clean after reset")
	}

	for i := range rf.Values {
		if got := rf.Values[i]; got != runtime.None {
			t.Fatalf("expected register %d to reset to runtime.None, got %v", i, got)
		}
	}
}
