package vm

import (
	"context"
	"errors"
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/bytecode"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

func TestAggregateReduce(t *testing.T) {
	for _, test := range []struct {
		want  runtime.Value
		empty runtime.Value
		kind  bytecode.AggregateKind
	}{
		{kind: bytecode.AggregateCount, want: runtime.Int(5), empty: runtime.ZeroInt},
		{kind: bytecode.AggregateSum, want: runtime.Float(4), empty: runtime.ZeroFloat},
		{kind: bytecode.AggregateMin, want: runtime.Float(1), empty: runtime.None},
		{kind: bytecode.AggregateMax, want: runtime.Float(3), empty: runtime.None},
		{kind: bytecode.AggregateAverage, want: runtime.Float(2), empty: runtime.ZeroFloat},
	} {
		for _, empty := range []bool{false, true} {
			for _, alias := range []bool{false, true} {
				source := &reductionList{}
				want := test.empty
				if !empty {
					source.values = []runtime.Value{runtime.Int(1), runtime.None, runtime.String("2"), runtime.True, runtime.Float(3)}
					want = test.want
				}

				program := aggregateReductionProgram(test.kind, false)
				if alias {
					program.Bytecode[1].Operands[0] = 1
					program.Bytecode[2].Operands[0] = 1
				}

				instance := mustNewVM(t, program)
				env := mustNewEnvironment(t, WithParam("values", source))
				result, err := instance.Run(t.Context(), env)
				if err != nil {
					t.Fatal(err)
				}

				if got := result.Root(); got != want {
					t.Fatalf("kind %d: got %v (%T), want %v (%T)", test.kind, got, got, want, want)
				}

				if err := result.Close(); err != nil {
					t.Fatal(err)
				}

				if err := instance.Close(); err != nil {
					t.Fatal(err)
				}

				if source.closed != 0 {
					t.Fatal("borrowed list was closed")
				}
			}
		}
	}
}

func TestAggregateReduceErrorsAndRecovery(t *testing.T) {
	sentinel := errors.New("host read failed")
	for _, protected := range []bool{false, true} {
		for _, cancelWithError := range []bool{false, true} {
			ctx, cancel := context.WithCancel(t.Context())
			source := &reductionList{values: []runtime.Value{runtime.Int(1)}, err: sentinel}
			if cancelWithError {
				source.cancel = cancel
			}

			instance := mustNewVM(t, aggregateReductionProgram(bytecode.AggregateSum, protected))
			env := mustNewEnvironment(t, WithParam("values", source))
			result, err := instance.Run(ctx, env)
			cancel()
			if protected && !cancelWithError {
				if err != nil {
					t.Fatal(err)
				}

				if result.Root() != runtime.None {
					t.Fatalf("recovery returned partial result %v", result.Root())
				}
			} else if !protected && !errors.Is(err, sentinel) {
				t.Fatalf("error = %v, want original host error", err)
			} else if protected && cancelWithError && err == nil {
				t.Fatal("cancellation was swallowed")
			}

			if result != nil {
				if err := result.Close(); err != nil {
					t.Fatal(err)
				}
			}

			if err := instance.Close(); err != nil {
				t.Fatal(err)
			}

			if source.closed != 0 {
				t.Fatal("borrowed list was closed")
			}
		}
	}

	ctx, cancel := context.WithCancel(t.Context())
	source := &reductionList{values: []runtime.Value{runtime.Int(1)}, cancel: cancel}
	instance := mustNewVM(t, aggregateReductionProgram(bytecode.AggregateSum, false))
	defer instance.Close()
	_, err := instance.Run(ctx, mustNewEnvironment(t, WithParam("values", source)))
	cancel()
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("error = %v, want cancellation", err)
	}
}

func TestAggregateReduceInvalidProgram(t *testing.T) {
	for _, kind := range []bytecode.AggregateKind{-1, bytecode.AggregateAverage + 1} {
		if _, err := NewWith(aggregateReductionProgram(kind, false)); err == nil {
			t.Fatal("invalid kind accepted")
		}
	}

	instance := mustNewVM(t, aggregateReductionProgram(bytecode.AggregateSum, true))
	defer instance.Close()
	_, err := instance.Run(t.Context(), mustNewEnvironment(t, WithParam("values", runtime.Int(1))))
	var invariant *RuntimeError
	if !errors.As(err, &invariant) || invariant.Message != "vm invariant violation" {
		t.Fatalf("error = %v, want unsuppressed invariant", err)
	}
}

func aggregateReductionProgram(kind bytecode.AggregateKind, protected bool) *bytecode.Program {
	p := &bytecode.Program{
		ISAVersion: bytecode.Version, Registers: 3, Params: []string{"values"},
		Bytecode: []bytecode.Instruction{
			bytecode.NewInstruction(bytecode.OpLoadParam, 1, 1),
			bytecode.NewInstruction(bytecode.OpAggregateReduce, 2, 1, bytecode.Operand(kind)),
			bytecode.NewInstruction(bytecode.OpReturn, 2),
		},
	}

	if protected {
		p.CatchTable = []bytecode.Catch{{1, 1, 2}}
	}

	return p
}
