package vm

import (
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/bytecode"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

func TestUnwindToProtected_ReclaimsDiscardedFrameRegisters(t *testing.T) {
	instance := New(&bytecode.Program{
		Registers: 1,
		Functions: bytecode.Functions{
			UserDefined: []bytecode.UDF{
				{Registers: 6},
			},
		},
	})

	lowerRegs := make([]runtime.Value, 2)
	protectedRegs := make([]runtime.Value, 3)
	aboveRegs1 := make([]runtime.Value, 4)
	aboveRegs2 := make([]runtime.Value, 5)
	activeRegs := make([]runtime.Value, 6)

	protectedRegs[1] = runtime.True

	instance.registers.Values = activeRegs
	instance.frames = []callFrame{
		{
			returnPC:   10,
			returnDest: bytecode.NewRegister(0),
			registers:  lowerRegs,
			protected:  false,
		},
		{
			returnPC:   20,
			returnDest: bytecode.NewRegister(1),
			registers:  protectedRegs,
			protected:  true,
		},
		{
			returnPC:   30,
			returnDest: bytecode.NewRegister(0),
			registers:  aboveRegs1,
			protected:  false,
		},
		{
			returnPC:   40,
			returnDest: bytecode.NewRegister(0),
			registers:  aboveRegs2,
			protected:  false,
		},
	}

	before2 := len(instance.regPool.buckets[2])
	before3 := len(instance.regPool.buckets[3])
	before4 := len(instance.regPool.buckets[4])
	before5 := len(instance.regPool.buckets[5])
	before6 := len(instance.regPool.buckets[6])

	if ok := instance.unwindToProtected(); !ok {
		t.Fatal("expected protected unwind to succeed")
	}

	if got, want := instance.pc, 20; got != want {
		t.Fatalf("unexpected pc after unwind: got %d, want %d", got, want)
	}

	if got, want := len(instance.frames), 1; got != want {
		t.Fatalf("unexpected frame depth after unwind: got %d, want %d", got, want)
	}

	if got, want := instance.frames[0].returnPC, 10; got != want {
		t.Fatalf("unexpected surviving frame returnPC: got %d, want %d", got, want)
	}

	if got, want := instance.registers.Values[1], runtime.None; got != want {
		t.Fatalf("expected protected return destination to be reset, got %v", got)
	}

	if got, want := len(instance.regPool.buckets[2]), before2; got != want {
		t.Fatalf("unexpected bucket[2] size: got %d, want %d", got, want)
	}

	if got, want := len(instance.regPool.buckets[3]), before3; got != want {
		t.Fatalf("unexpected bucket[3] size: got %d, want %d", got, want)
	}

	if got, want := len(instance.regPool.buckets[4]), before4+1; got != want {
		t.Fatalf("unexpected bucket[4] size: got %d, want %d", got, want)
	}

	if got, want := len(instance.regPool.buckets[5]), before5+1; got != want {
		t.Fatalf("unexpected bucket[5] size: got %d, want %d", got, want)
	}

	if got, want := len(instance.regPool.buckets[6]), before6+1; got != want {
		t.Fatalf("unexpected bucket[6] size: got %d, want %d", got, want)
	}
}
