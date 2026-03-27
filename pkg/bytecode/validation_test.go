package bytecode

import (
	"errors"
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/file"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

func TestValidateProgram(t *testing.T) {
	tests := []struct {
		target  error
		program *Program
		name    string
	}{
		{
			name:    "unknown_opcode",
			program: withProgramMutation(func(program *Program) { program.Bytecode[0] = NewInstruction(Opcode(255)) }),
			target:  ErrInvalidInstruction,
		},
		{
			name:    "register_out_of_range",
			program: withProgramMutation(func(program *Program) { program.Bytecode[0] = NewInstruction(OpReturn, NewRegister(3)) }),
			target:  ErrInvalidInstruction,
		},
		{
			name: "constant_out_of_range",
			program: withProgramMutation(func(program *Program) {
				program.Bytecode[0] = NewInstruction(OpLoadConst, NewRegister(0), NewConstant(99))
			}),
			target: ErrInvalidInstruction,
		},
		{
			name:    "jump_target_out_of_range",
			program: withProgramMutation(func(program *Program) { program.Bytecode[0] = NewInstruction(OpJump, Operand(99)) }),
			target:  ErrInvalidInstruction,
		},
		{
			name: "invalid_catch_entry",
			program: withProgramMutation(func(program *Program) {
				program.CatchTable = []Catch{{1, 0, 1}}
			}),
			target: ErrInvalidProgram,
		},
		{
			name: "invalid_aggregate_metadata",
			program: withProgramMutation(func(program *Program) {
				program.Metadata.AggregateSelectorSlots = []int{-2, -1}
			}),
			target: ErrInvalidProgram,
		},
		{
			name: "required_constant_type_check",
			program: withProgramMutation(func(program *Program) {
				program.Bytecode[0] = NewInstruction(OpFail, NewConstant(0))
				program.Constants[0] = runtime.NewInt(1)
			}),
			target: ErrInvalidInstruction,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := ValidateProgram(tc.program)
			if !errors.Is(err, tc.target) {
				t.Fatalf("expected %v, got %v", tc.target, err)
			}
		})
	}
}

func validValidationProgram() *Program {
	return &Program{
		Source: file.NewSource("validation.fql", "RETURN 1"),
		Functions: Functions{
			Host: map[string]int{
				"now": 0,
			},
			UserDefined: []UDF{
				{
					Name:        "main",
					DisplayName: "main",
					Entry:       1,
					Registers:   1,
					Params:      0,
				},
			},
		},
		Bytecode: []Instruction{
			NewInstruction(OpLoadConst, NewRegister(0), NewConstant(0)),
			NewInstruction(OpReturn, NewRegister(0)),
		},
		Constants: []runtime.Value{
			runtime.NewString("ok"),
		},
		CatchTable: nil,
		Params:     []string{"input"},
		Metadata: Metadata{
			Labels:                 map[int]string{1: "exit"},
			CompilerVersion:        "test",
			AggregatePlans:         []AggregatePlan{NewAggregatePlan([]runtime.String{runtime.NewString("group")}, []AggregateKind{AggregateCount}, false)},
			AggregateSelectorSlots: []int{-1, -1},
			MatchFailTargets:       []int{-1, -1},
			DebugSpans:             []file.Span{{Start: 0, End: 6}, {Start: 7, End: 8}},
			OptimizationLevel:      1,
		},
		ISAVersion: Version,
		Registers:  3,
	}
}

func withProgramMutation(mutate func(program *Program)) *Program {
	program := validValidationProgram()

	if mutate != nil {
		mutate(program)
	}

	return program
}
