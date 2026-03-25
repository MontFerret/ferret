package optimization_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/bytecode"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
	"github.com/MontFerret/ferret/v2/pkg/vm"
)

func TestOpConcat(t *testing.T) {
	testCases := []struct {
		name     string
		expected runtime.String
		values   []runtime.Value
		startReg int
		count    int
	}{
		{
			name:     "count zero returns empty string",
			values:   nil,
			startReg: 1,
			count:    0,
			expected: runtime.EmptyString,
		},
		{
			name: "count one coerces to string",
			values: []runtime.Value{
				runtime.NewInt(42),
			},
			startReg: 1,
			count:    1,
			expected: runtime.NewString("42"),
		},
		{
			name: "count two fast-path skips empty prefix",
			values: []runtime.Value{
				runtime.EmptyString,
				runtime.NewString("bar"),
			},
			startReg: 1,
			count:    2,
			expected: runtime.NewString("bar"),
		},
		{
			name: "count three fast-path handles trailing empties",
			values: []runtime.Value{
				runtime.NewString("a"),
				runtime.EmptyString,
				runtime.EmptyString,
			},
			startReg: 1,
			count:    3,
			expected: runtime.NewString("a"),
		},
		{
			name: "builder path with all empty parts returns empty string",
			values: []runtime.Value{
				runtime.None,
				runtime.EmptyString,
				runtime.None,
				runtime.EmptyString,
			},
			startReg: 1,
			count:    4,
			expected: runtime.EmptyString,
		},
		{
			name: "builder path with unicode and mixed scalar types",
			values: []runtime.Value{
				runtime.NewString("пр"),
				runtime.NewInt(7),
				runtime.True,
				runtime.NewString("!"),
			},
			startReg: 5,
			count:    4,
			expected: runtime.NewString("пр7true!"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			program := buildConcatProgram(tc.values, tc.startReg, tc.count)

			instance, err := vm.New(program)
			if err != nil {
				t.Fatalf("unexpected constructor error: %v", err)
			}

			out, err := instance.Run(context.Background(), vm.NewDefaultEnvironment())
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			defer func() {
				_ = out.Close()
			}()

			got, ok := out.Root().(runtime.String)
			if !ok {
				t.Fatalf("expected runtime.String, got %T", out.Root())
			}

			if got != tc.expected {
				t.Fatalf("unexpected result: got %q, want %q", got, tc.expected)
			}
		})
	}
}

func TestOpAddConst(t *testing.T) {
	testCases := []struct {
		left     runtime.Value
		right    runtime.Value
		expected runtime.Value
		name     string
	}{
		{
			name:     "int plus int",
			left:     runtime.NewInt(1),
			right:    runtime.NewInt(2),
			expected: runtime.NewInt(3),
		},
		{
			name:     "float plus int",
			left:     runtime.NewFloat(1.5),
			right:    runtime.NewInt(2),
			expected: runtime.NewFloat(3.5),
		},
		{
			name:     "int plus string",
			left:     runtime.NewInt(1),
			right:    runtime.NewString("x"),
			expected: runtime.NewString("1x"),
		},
		{
			name:     "string plus int",
			left:     runtime.NewString("x"),
			right:    runtime.NewInt(1),
			expected: runtime.NewString("x1"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			program := &bytecode.Program{
				ISAVersion: bytecode.Version,
				Constants: []runtime.Value{
					tc.left,
					tc.right,
				},
				Bytecode: []bytecode.Instruction{
					bytecode.NewInstruction(bytecode.OpLoadConst, bytecode.NewRegister(1), bytecode.NewConstant(0)),
					bytecode.NewInstruction(bytecode.OpAddConst, bytecode.NewRegister(2), bytecode.NewRegister(1), bytecode.NewConstant(1)),
					bytecode.NewInstruction(bytecode.OpReturn, bytecode.NewRegister(2)),
				},
				Registers: 3,
			}

			instance, err := vm.New(program)
			if err != nil {
				t.Fatalf("unexpected constructor error: %v", err)
			}

			out, err := instance.Run(context.Background(), vm.NewDefaultEnvironment())
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			defer func() {
				_ = out.Close()
			}()

			if err := assertRuntimeValueEqual(out.Root(), tc.expected); err != nil {
				t.Fatalf("unexpected result: %v", err)
			}
		})
	}
}

func buildConcatProgram(values []runtime.Value, startReg, count int) *bytecode.Program {
	instructions := make([]bytecode.Instruction, 0, len(values)+2)

	for i := 0; i < len(values); i++ {
		instructions = append(instructions, bytecode.NewInstruction(
			bytecode.OpLoadConst,
			bytecode.NewRegister(startReg+i),
			bytecode.NewConstant(i),
		))
	}

	dst := startReg + len(values) + 1

	instructions = append(instructions, bytecode.NewInstruction(
		bytecode.OpConcat,
		bytecode.NewRegister(dst),
		bytecode.NewRegister(startReg),
		bytecode.Operand(count),
	))
	instructions = append(instructions, bytecode.NewInstruction(bytecode.OpReturn, bytecode.NewRegister(dst)))

	return &bytecode.Program{
		ISAVersion: bytecode.Version,
		Constants:  values,
		Bytecode:   instructions,
		Registers:  dst + 1,
	}
}

func assertRuntimeValueEqual(actual, expected runtime.Value) error {
	switch want := expected.(type) {
	case runtime.String:
		got, ok := actual.(runtime.String)
		if !ok {
			return fmt.Errorf("expected runtime.String, got %T", actual)
		}
		if got != want {
			return fmt.Errorf("expected %q, got %q", want, got)
		}
		return nil
	case runtime.Int:
		got, ok := actual.(runtime.Int)
		if !ok {
			return fmt.Errorf("expected runtime.Int, got %T", actual)
		}
		if got != want {
			return fmt.Errorf("expected %v, got %v", want, got)
		}
		return nil
	case runtime.Float:
		got, ok := actual.(runtime.Float)
		if !ok {
			return fmt.Errorf("expected runtime.Float, got %T", actual)
		}
		if got != want {
			return fmt.Errorf("expected %v, got %v", want, got)
		}
		return nil
	case runtime.Boolean:
		got, ok := actual.(runtime.Boolean)
		if !ok {
			return fmt.Errorf("expected runtime.Boolean, got %T", actual)
		}
		if got != want {
			return fmt.Errorf("expected %v, got %v", want, got)
		}
		return nil
	default:
		return fmt.Errorf("unsupported expected type %T", expected)
	}
}
