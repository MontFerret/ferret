package vm_test

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/bytecode"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
	"github.com/MontFerret/ferret/v2/pkg/vm"
	"github.com/MontFerret/ferret/v2/test/spec"
	"github.com/MontFerret/ferret/v2/test/spec/assert"
)

type integrationQueryStub struct {
	result  runtime.List
	err     error
	queries []runtime.Query
}

func (q *integrationQueryStub) Query(_ context.Context, query runtime.Query) (runtime.List, error) {
	q.queries = append(q.queries, query)

	if q.err != nil {
		return nil, q.err
	}

	if q.result != nil {
		return q.result, nil
	}

	return runtime.NewArray(0), nil
}

func (q *integrationQueryStub) String() string {
	return "integration-query-stub"
}

func (q *integrationQueryStub) Hash() uint64 {
	return runtime.NewString("integration-query-stub").Hash()
}

func (q *integrationQueryStub) Copy() runtime.Value {
	return q
}

func programWithApplyQueryConstSource(src runtime.Value) *bytecode.Program {
	return &bytecode.Program{
		ISAVersion: bytecode.Version,
		Registers:  1,
		Bytecode: []bytecode.Instruction{
			bytecode.NewInstruction(bytecode.OpQuery, bytecode.NewRegister(0), bytecode.NewConstant(0), bytecode.NewConstant(1)),
			bytecode.NewInstruction(bytecode.OpReturn, bytecode.NewRegister(0)),
		},
		Constants: []runtime.Value{
			src,
			runtime.NewArrayWith(
				runtime.NewString("css"),
				runtime.NewString(".items"),
				runtime.None,
			),
		},
	}
}

func TestApplyQueryConstantSource_Strict(t *testing.T) {
	stub := &integrationQueryStub{
		result: runtime.NewArrayWith(runtime.NewString("ok")),
	}

	runProgramSpecs(t, []spec.Spec{
		spec.NewWith(spec.NewProgramInput(programWithApplyQueryConstSource(stub)), "strict constant source").
			Expect().Exec(assert.NewUnaryAssertion(func(actual any) error {
			if err := assert.Equal(actual, []any{"ok"}); err != nil {
				return err
			}

			if len(stub.queries) != 1 {
				return fmt.Errorf("unexpected query count: got %d, want 1", len(stub.queries))
			}

			if got, want := stub.queries[0].Kind, runtime.NewString("css"); got != want {
				return fmt.Errorf("unexpected query kind: got %q, want %q", got, want)
			}

			if got, want := stub.queries[0].Payload, runtime.NewString(".items"); got != want {
				return fmt.Errorf("unexpected query payload: got %q, want %q", got, want)
			}

			return nil
		})),
	})
}

func TestApplyQueryConstantSource_FastMode_NoPanic(t *testing.T) {
	stub := &integrationQueryStub{
		result: runtime.NewArrayWith(runtime.NewString("ok")),
	}

	runProgramSpecs(t, []spec.Spec{
		spec.NewWith(spec.NewProgramInput(programWithApplyQueryConstSource(stub)), "fast constant source").
			VM(vm.WithPanicPolicy(vm.PanicPropagate)).
			Expect().Exec(assert.NewUnaryAssertion(func(actual any) error {
			if err := assert.Equal(actual, []any{"ok"}); err != nil {
				return err
			}

			if len(stub.queries) != 1 {
				return fmt.Errorf("unexpected query count: got %d, want 1", len(stub.queries))
			}

			return nil
		})),
	})
}

func TestApplyQueryConstantSource_NonQueryable_NoPanicTypeError(t *testing.T) {
	runProgramSpecs(t, []spec.Spec{
		spec.NewWith(spec.NewProgramInput(programWithApplyQueryConstSource(runtime.NewInt(1))), "strict non-queryable").
			Expect().ExecError(assert.NewUnaryAssertion(func(actual any) error {
			err, ok := actual.(error)
			if !ok || err == nil {
				return errors.New("expected type error")
			}

			var rtErr *vm.RuntimeError
			if !errors.As(err, &rtErr) {
				return fmt.Errorf("expected runtime error, got %T", err)
			}

			if !strings.Contains(strings.ToLower(rtErr.Message), "invalid type") {
				return fmt.Errorf("expected invalid type error, got %q", rtErr.Message)
			}

			return nil
		})),
		spec.NewWith(spec.NewProgramInput(programWithApplyQueryConstSource(runtime.NewInt(1))), "fast non-queryable").
			VM(vm.WithPanicPolicy(vm.PanicPropagate)).
			Expect().ExecError(assert.NewUnaryAssertion(func(actual any) error {
			err, ok := actual.(error)
			if !ok || err == nil {
				return errors.New("expected type error")
			}

			var rtErr *vm.RuntimeError
			if !errors.As(err, &rtErr) {
				return fmt.Errorf("expected runtime error, got %T", err)
			}

			if !strings.Contains(strings.ToLower(rtErr.Message), "invalid type") {
				return fmt.Errorf("expected invalid type error, got %q", rtErr.Message)
			}

			return nil
		})),
	})
}
