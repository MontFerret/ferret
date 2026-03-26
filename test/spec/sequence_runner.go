package spec

import (
	"fmt"
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/bytecode"
	"github.com/MontFerret/ferret/v2/pkg/compiler"
	"github.com/MontFerret/ferret/v2/pkg/diagnostics"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
	"github.com/MontFerret/ferret/v2/pkg/vm"
)

type SequenceRunner struct {
	Name     string
	Compiler *compiler.Compiler
	Env      []vm.EnvironmentOption
}

func NewSequenceRunner(suite string, opts ...compiler.Option) *SequenceRunner {
	return &SequenceRunner{
		Name:     suite,
		Compiler: compiler.New(opts...),
	}
}

func (r *SequenceRunner) Run(t *testing.T, sequences []Sequence) {
	t.Helper()

	std := Stdlib()

	for _, sequence := range sequences {
		suiteName := sequence.Base.SuiteName(r.Name)

		t.Run(suiteName, func(t *testing.T) {
			if sequence.Base.SkipInfo.Active {
				t.Skip(sequence.Base.SkipInfo.Reason)
			}

			var prog *bytecode.Program

			defer func() {
				if recovered := recover(); recovered != nil {
					PrintDebug(t, suiteName, prog)
					t.Fatalf("panic: %v", recovered)
				}
			}()

			prog, err := sequence.Base.Input.ResolveProgram(suiteName, r.Compiler)
			if err != nil {
				if sequence.Base.DebugOutput {
					PrintError(t, err)
				}

				t.Fatalf("unexpected compilation error:\n%s", diagnostics.Format(err))
			}

			instance, err := vm.NewWith(prog, sequence.VM...)
			if err != nil {
				t.Fatalf("unexpected VM constructor error: %v", err)
			}
			defer func() {
				_ = instance.Close()
			}()

			for i, step := range sequence.Steps {
				stepName := step.Name
				if stepName == "" {
					stepName = fmt.Sprintf("Step %d", i+1)
				}

				t.Run(stepName, func(t *testing.T) {
					r.runStep(t, suiteName, prog, std, instance, sequence, step)
				})
			}

			if sequence.Base.DebugOutput {
				PrintDebug(t, suiteName, prog)
			}
		})
	}
}

func (r *SequenceRunner) runStep(
	t *testing.T,
	suiteName string,
	prog *bytecode.Program,
	std runtime.Namespace,
	instance *vm.VM,
	sequence Sequence,
	step SequenceStep,
) {
	t.Helper()

	if countStepExpectations(step) != 1 {
		t.Fatal("sequence step must define exactly one of result, error, or panic expectation")
	}

	env, err := r.newStepEnvironment(std, sequence, step)
	if err != nil {
		t.Fatalf("unexpected environment error: %v", err)
	}

	if step.Panic.Defined() {
		r.runPanicStep(t, instance, env, step)
		return
	}

	defer func() {
		if recovered := recover(); recovered != nil {
			PrintDebug(t, suiteName, prog)
			t.Fatalf("panic: %v", recovered)
		}
	}()

	actual, err := ExecInstance(instance, false, env)
	if err != nil {
		if step.Error.Defined() {
			step.Error.Assert(t, err)
			return
		}

		t.Fatalf("unexpected runtime error: %v", err)
	}

	if step.Error.Defined() {
		t.Fatal("expected runtime error, got none")
	}

	step.Result.Assert(t, actual)
}

func (r *SequenceRunner) newStepEnvironment(ns runtime.Namespace, sequence Sequence, step SequenceStep) (*vm.Environment, error) {
	if step.EnvFactory != nil {
		return step.EnvFactory()
	}

	opts := []vm.EnvironmentOption{
		vm.WithNamespace(ns),
	}

	if len(r.Env) > 0 {
		opts = append(opts, r.Env...)
	}

	if len(sequence.Env) > 0 {
		opts = append(opts, sequence.Env...)
	}

	if len(step.Env) > 0 {
		opts = append(opts, step.Env...)
	}

	return vm.NewEnvironment(opts)
}

func (r *SequenceRunner) runPanicStep(t *testing.T, instance *vm.VM, env *vm.Environment, step SequenceStep) {
	t.Helper()

	defer func() {
		recovered := recover()
		if recovered == nil {
			t.Fatal("expected panic, got none")
		}

		step.Panic.Assert(t, recovered)
	}()

	if _, err := ExecInstance(instance, false, env); err != nil {
		t.Fatalf("unexpected runtime error: %v", err)
	}
}

func countStepExpectations(step SequenceStep) int {
	count := 0

	if step.Result.Defined() {
		count++
	}

	if step.Error.Defined() {
		count++
	}

	if step.Panic.Defined() {
		count++
	}

	return count
}
