package spec

import (
	"context"
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/bytecode"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
	"github.com/MontFerret/ferret/v2/pkg/vm"
)

func TestRunWithAppliesVMOptions(t *testing.T) {
	prog, err := Compile("RETURN PANIC_FN()")
	if err != nil {
		t.Fatalf("compile failed: %v", err)
	}

	defer func() {
		recovered := recover()
		if recovered == nil {
			t.Fatal("expected panic to propagate")
		}

		if got, want := recovered, "boom"; got != want {
			t.Fatalf("unexpected panic: got %v, want %v", got, want)
		}
	}()

	_, _ = RunWith(
		prog,
		[]vm.Option{vm.WithPanicPolicy(vm.PanicPropagate)},
		vm.WithNamespace(Stdlib()),
		vm.WithFunction("PANIC_FN", func(context.Context, ...runtime.Value) (runtime.Value, error) {
			panic("boom")
		}),
	)
}

func constantProgram(value runtime.Value) *bytecode.Program {
	return &bytecode.Program{
		ISAVersion: bytecode.Version,
		Constants:  []runtime.Value{value},
		Bytecode: []bytecode.Instruction{
			bytecode.NewInstruction(bytecode.OpLoadConst, bytecode.NewRegister(0), bytecode.NewConstant(0)),
			bytecode.NewInstruction(bytecode.OpReturn, bytecode.NewRegister(0)),
		},
		Registers: 1,
	}
}
