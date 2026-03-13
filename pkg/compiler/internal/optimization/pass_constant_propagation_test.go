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

type concatFoldingCase struct {
	wantConst    runtime.Value
	program      func() *bytecode.Program
	name         string
	target       int
	wantModified bool
	wantOpcode   bytecode.Opcode
}

var concatFoldingPositiveCases = []concatFoldingCase{
	{
		name: "folds concat with two constants",
		program: func() *bytecode.Program {
			return &bytecode.Program{
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
			}
		},
		target:       2,
		wantModified: true,
		wantOpcode:   bytecode.OpLoadConst,
		wantConst:    runtime.NewString("sum=3"),
	},
	{
		name: "folds concat with loadnone and scalar values",
		program: func() *bytecode.Program {
			return &bytecode.Program{
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
			}
		},
		target:       3,
		wantModified: true,
		wantOpcode:   bytecode.OpLoadConst,
		wantConst:    runtime.NewString("true7"),
	},
}

var concatFoldingNeutralCases = []concatFoldingCase{
	{
		name: "folds empty concat to empty string",
		program: func() *bytecode.Program {
			return &bytecode.Program{
				Bytecode: []bytecode.Instruction{
					bytecode.NewInstruction(bytecode.OpConcat, bytecode.NewRegister(1), bytecode.NewRegister(1), bytecode.Operand(0)),
					bytecode.NewInstruction(bytecode.OpReturn, bytecode.NewRegister(1)),
				},
			}
		},
		target:       0,
		wantModified: true,
		wantOpcode:   bytecode.OpLoadConst,
		wantConst:    runtime.EmptyString,
	},
	{
		name: "folds single element concat",
		program: func() *bytecode.Program {
			return &bytecode.Program{
				Constants: []runtime.Value{
					runtime.NewString("foo"),
				},
				Bytecode: []bytecode.Instruction{
					bytecode.NewInstruction(bytecode.OpLoadConst, bytecode.NewRegister(2), bytecode.NewConstant(0)),
					bytecode.NewInstruction(bytecode.OpConcat, bytecode.NewRegister(1), bytecode.NewRegister(2), bytecode.Operand(1)),
					bytecode.NewInstruction(bytecode.OpReturn, bytecode.NewRegister(1)),
				},
			}
		},
		target:       1,
		wantModified: true,
		wantOpcode:   bytecode.OpLoadConst,
		wantConst:    runtime.NewString("foo"),
	},
}

var concatFoldingNoFoldCases = []concatFoldingCase{
	{
		name: "keeps concat when inputs are unknown",
		program: func() *bytecode.Program {
			return &bytecode.Program{
				Constants: []runtime.Value{
					runtime.NewString("sum="),
				},
				Bytecode: []bytecode.Instruction{
					bytecode.NewInstruction(bytecode.OpLoadConst, bytecode.NewRegister(1), bytecode.NewConstant(0)),
					bytecode.NewInstruction(bytecode.OpMove, bytecode.NewRegister(2), bytecode.NewRegister(3)),
					bytecode.NewInstruction(bytecode.OpConcat, bytecode.NewRegister(4), bytecode.NewRegister(1), bytecode.Operand(2)),
					bytecode.NewInstruction(bytecode.OpReturn, bytecode.NewRegister(4)),
				},
			}
		},
		target:       2,
		wantModified: false,
		wantOpcode:   bytecode.OpConcat,
	},
}

func runConcatFoldingCase(t *testing.T, tc concatFoldingCase) {
	t.Helper()

	program := tc.program()
	res, err := runConstantPropagation(t, program)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if res.Modified != tc.wantModified {
		t.Fatalf("unexpected modified flag: got %v, want %v", res.Modified, tc.wantModified)
	}

	inst := program.Bytecode[tc.target]
	if inst.Opcode != tc.wantOpcode {
		t.Fatalf("unexpected opcode at %d: got %s, want %s", tc.target, inst.Opcode, tc.wantOpcode)
	}

	if tc.wantConst == nil {
		return
	}

	if !inst.Operands[1].IsConstant() {
		t.Fatalf("expected LOADC to use a constant operand, got %s", inst.Operands[1])
	}

	got := program.Constants[inst.Operands[1].Constant()]
	if err := assertConstEqual(got, tc.wantConst); err != nil {
		t.Fatalf("unexpected folded constant: %v", err)
	}
}

func TestConstantPropagation_ConcatFolding_PositiveCases(t *testing.T) {
	for _, tc := range concatFoldingPositiveCases {
		t.Run(tc.name, func(t *testing.T) {
			runConcatFoldingCase(t, tc)
		})
	}
}

func TestConstantPropagation_ConcatFolding_NeutralCases(t *testing.T) {
	for _, tc := range concatFoldingNeutralCases {
		t.Run(tc.name, func(t *testing.T) {
			runConcatFoldingCase(t, tc)
		})
	}
}

func TestConstantPropagation_ConcatFolding_NoFoldCases(t *testing.T) {
	for _, tc := range concatFoldingNoFoldCases {
		t.Run(tc.name, func(t *testing.T) {
			runConcatFoldingCase(t, tc)
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

type resolveBinaryFoldOperandsCase struct {
	wantLeft  runtime.Value
	wantRight runtime.Value
	state     constState
	program   *bytecode.Program
	name      string
	inst      bytecode.Instruction
	wantOK    bool
}

var resolveBinaryFoldOperandsCases = []resolveBinaryFoldOperandsCase{
	{
		name: "register operands",
		inst: bytecode.NewInstruction(bytecode.OpAdd, bytecode.NewRegister(3), bytecode.NewRegister(1), bytecode.NewRegister(2)),
		state: constState{
			1: runtime.NewInt(10),
			2: runtime.NewInt(5),
		},
		program:   &bytecode.Program{},
		wantLeft:  runtime.NewInt(10),
		wantRight: runtime.NewInt(5),
		wantOK:    true,
	},
	{
		name: "add const right operand",
		inst: bytecode.NewInstruction(bytecode.OpAddConst, bytecode.NewRegister(3), bytecode.NewRegister(1), bytecode.NewConstant(0)),
		state: constState{
			1: runtime.NewInt(2),
		},
		program: &bytecode.Program{
			Constants: []runtime.Value{
				runtime.NewInt(7),
			},
		},
		wantLeft:  runtime.NewInt(2),
		wantRight: runtime.NewInt(7),
		wantOK:    true,
	},
	{
		name: "missing left operand constant",
		inst: bytecode.NewInstruction(bytecode.OpAdd, bytecode.NewRegister(3), bytecode.NewRegister(1), bytecode.NewRegister(2)),
		state: constState{
			2: runtime.NewInt(5),
		},
		program: &bytecode.Program{},
		wantOK:  false,
	},
	{
		name: "add const without constant operand",
		inst: bytecode.NewInstruction(bytecode.OpAddConst, bytecode.NewRegister(3), bytecode.NewRegister(1), bytecode.NewRegister(2)),
		state: constState{
			1: runtime.NewInt(2),
			2: runtime.NewInt(3),
		},
		program: &bytecode.Program{},
		wantOK:  false,
	},
	{
		name:    "non register left operand",
		inst:    bytecode.NewInstruction(bytecode.OpAdd, bytecode.NewRegister(3), bytecode.NewConstant(0), bytecode.NewRegister(2)),
		state:   constState{},
		program: &bytecode.Program{},
		wantOK:  false,
	},
}

func runResolveBinaryFoldOperandsCase(t *testing.T, tc resolveBinaryFoldOperandsCase) {
	t.Helper()

	gotLeft, gotRight, ok := resolveBinaryFoldOperands(tc.inst, tc.state, tc.program)
	if ok != tc.wantOK {
		t.Fatalf("unexpected ok flag: got %v, want %v", ok, tc.wantOK)
	}

	if !tc.wantOK {
		return
	}

	if err := assertConstEqual(gotLeft, tc.wantLeft); err != nil {
		t.Fatalf("unexpected left operand: %v", err)
	}
	if err := assertConstEqual(gotRight, tc.wantRight); err != nil {
		t.Fatalf("unexpected right operand: %v", err)
	}
}

func TestConstantPropagation_ResolveBinaryFoldOperands(t *testing.T) {
	for _, tc := range resolveBinaryFoldOperandsCases {
		t.Run(tc.name, func(t *testing.T) {
			runResolveBinaryFoldOperandsCase(t, tc)
		})
	}
}

func TestConstantPropagation_ResolveConcatFoldOperands(t *testing.T) {
	tests := []struct {
		name      string
		inst      bytecode.Instruction
		wantDst   int
		wantStart int
		wantCount int
		wantOK    bool
	}{
		{
			name:      "register operands",
			inst:      bytecode.NewInstruction(bytecode.OpConcat, bytecode.NewRegister(4), bytecode.NewRegister(2), bytecode.Operand(3)),
			wantDst:   4,
			wantStart: 2,
			wantCount: 3,
			wantOK:    true,
		},
		{
			name:   "missing destination register",
			inst:   bytecode.NewInstruction(bytecode.OpConcat, bytecode.NewConstant(0), bytecode.NewRegister(2), bytecode.Operand(3)),
			wantOK: false,
		},
		{
			name:   "missing start register",
			inst:   bytecode.NewInstruction(bytecode.OpConcat, bytecode.NewRegister(4), bytecode.NewConstant(0), bytecode.Operand(3)),
			wantOK: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			dst, start, count, ok := resolveConcatFoldOperands(tc.inst)
			if ok != tc.wantOK {
				t.Fatalf("unexpected ok flag: got %v, want %v", ok, tc.wantOK)
			}

			if !tc.wantOK {
				return
			}

			if dst != tc.wantDst || start != tc.wantStart || count != tc.wantCount {
				t.Fatalf("unexpected operands: got (%d, %d, %d), want (%d, %d, %d)", dst, start, count, tc.wantDst, tc.wantStart, tc.wantCount)
			}
		})
	}
}

func TestConstantPropagation_BuildConcatFoldConst(t *testing.T) {
	tests := []struct {
		want   runtime.Value
		state  constState
		name   string
		start  int
		count  int
		wantOK bool
	}{
		{
			name: "all constants present",
			state: constState{
				1: runtime.NewString("sum="),
				2: runtime.NewInt(3),
				3: runtime.True,
			},
			start:  1,
			count:  3,
			want:   runtime.NewString("sum=3true"),
			wantOK: true,
		},
		{
			name: "missing range entry",
			state: constState{
				1: runtime.NewString("a"),
				3: runtime.NewString("b"),
			},
			start:  1,
			count:  3,
			wantOK: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, ok := buildConcatFoldConst(tc.state, tc.start, tc.count)
			if ok != tc.wantOK {
				t.Fatalf("unexpected ok flag: got %v, want %v", ok, tc.wantOK)
			}

			if !tc.wantOK {
				return
			}

			if err := assertConstEqual(got, tc.want); err != nil {
				t.Fatalf("unexpected folded concat constant: %v", err)
			}
		})
	}
}
