package optimization

import (
	"fmt"
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/bytecode"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

func runConstantPropagation(t *testing.T, program *bytecode.Program) (*PassResult, error) {
	t.Helper()

	builder := NewBuilder(program)
	cfg, err := builder.Build()
	if err != nil {
		return nil, err
	}

	pass := NewConstantPropagationPass()

	return pass.Run(&PassContext{
		Program: program,
		CFG:     cfg,
	})
}

func TestConstantPropagation_ConcatFoldingCases(t *testing.T) {
	testCases := []struct {
		name         string
		program      *bytecode.Program
		target       int
		wantModified bool
		wantOpcode   bytecode.Opcode
		wantConst    runtime.Value
	}{
		{
			name: "folds concat with two constants",
			program: &bytecode.Program{
				Constants: []runtime.Value{
					runtime.NewString("sum="),
					runtime.NewInt(3),
				},
				Bytecode: []bytecode.Instruction{
					bytecode.NewInstruction(bytecode.OpLoadConst, bytecode.NewRegister(1), bytecode.NewConstant(0)),
					bytecode.NewInstruction(bytecode.OpLoadConst, bytecode.NewRegister(2), bytecode.NewConstant(1)),
					bytecode.NewInstruction(bytecode.OpConcat, bytecode.NewRegister(3), bytecode.NewRegister(1), bytecode.Operand(2)),
					bytecode.NewInstruction(bytecode.OpReturn, bytecode.NewRegister(3)),
				},
			},
			target:       2,
			wantModified: true,
			wantOpcode:   bytecode.OpLoadConst,
			wantConst:    runtime.NewString("sum=3"),
		},
		{
			name: "folds empty concat to empty string",
			program: &bytecode.Program{
				Bytecode: []bytecode.Instruction{
					bytecode.NewInstruction(bytecode.OpConcat, bytecode.NewRegister(1), bytecode.NewRegister(1), bytecode.Operand(0)),
					bytecode.NewInstruction(bytecode.OpReturn, bytecode.NewRegister(1)),
				},
			},
			target:       0,
			wantModified: true,
			wantOpcode:   bytecode.OpLoadConst,
			wantConst:    runtime.EmptyString,
		},
		{
			name: "folds single element concat",
			program: &bytecode.Program{
				Constants: []runtime.Value{
					runtime.NewString("foo"),
				},
				Bytecode: []bytecode.Instruction{
					bytecode.NewInstruction(bytecode.OpLoadConst, bytecode.NewRegister(2), bytecode.NewConstant(0)),
					bytecode.NewInstruction(bytecode.OpConcat, bytecode.NewRegister(1), bytecode.NewRegister(2), bytecode.Operand(1)),
					bytecode.NewInstruction(bytecode.OpReturn, bytecode.NewRegister(1)),
				},
			},
			target:       1,
			wantModified: true,
			wantOpcode:   bytecode.OpLoadConst,
			wantConst:    runtime.NewString("foo"),
		},
		{
			name: "folds concat with loadnone and scalar values",
			program: &bytecode.Program{
				Constants: []runtime.Value{
					runtime.True,
					runtime.NewInt(7),
				},
				Bytecode: []bytecode.Instruction{
					bytecode.NewInstruction(bytecode.OpLoadNone, bytecode.NewRegister(1)),
					bytecode.NewInstruction(bytecode.OpLoadConst, bytecode.NewRegister(2), bytecode.NewConstant(0)),
					bytecode.NewInstruction(bytecode.OpLoadConst, bytecode.NewRegister(3), bytecode.NewConstant(1)),
					bytecode.NewInstruction(bytecode.OpConcat, bytecode.NewRegister(4), bytecode.NewRegister(1), bytecode.Operand(3)),
					bytecode.NewInstruction(bytecode.OpReturn, bytecode.NewRegister(4)),
				},
			},
			target:       3,
			wantModified: true,
			wantOpcode:   bytecode.OpLoadConst,
			wantConst:    runtime.NewString("true7"),
		},
		{
			name: "keeps concat when inputs are unknown",
			program: &bytecode.Program{
				Constants: []runtime.Value{
					runtime.NewString("sum="),
				},
				Bytecode: []bytecode.Instruction{
					bytecode.NewInstruction(bytecode.OpLoadConst, bytecode.NewRegister(1), bytecode.NewConstant(0)),
					bytecode.NewInstruction(bytecode.OpMove, bytecode.NewRegister(2), bytecode.NewRegister(3)),
					bytecode.NewInstruction(bytecode.OpConcat, bytecode.NewRegister(4), bytecode.NewRegister(1), bytecode.Operand(2)),
					bytecode.NewInstruction(bytecode.OpReturn, bytecode.NewRegister(4)),
				},
			},
			target:       2,
			wantModified: false,
			wantOpcode:   bytecode.OpConcat,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			res, err := runConstantPropagation(t, tc.program)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if res.Modified != tc.wantModified {
				t.Fatalf("unexpected modified flag: got %v, want %v", res.Modified, tc.wantModified)
			}

			inst := tc.program.Bytecode[tc.target]
			if inst.Opcode != tc.wantOpcode {
				t.Fatalf("unexpected opcode at %d: got %s, want %s", tc.target, inst.Opcode, tc.wantOpcode)
			}

			if tc.wantConst == nil {
				return
			}

			if !inst.Operands[1].IsConstant() {
				t.Fatalf("expected LOADC to use a constant operand, got %s", inst.Operands[1])
			}

			got := tc.program.Constants[inst.Operands[1].Constant()]
			if err := assertConstEqual(got, tc.wantConst); err != nil {
				t.Fatalf("unexpected folded constant: %v", err)
			}
		})
	}
}

func assertConstEqual(actual, expected runtime.Value) error {
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
		if expected == runtime.None {
			if actual != runtime.None {
				return fmt.Errorf("expected runtime.None, got %T", actual)
			}
			return nil
		}

		return fmt.Errorf("unsupported expected constant type %T", expected)
	}
}
