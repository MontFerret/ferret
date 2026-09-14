package bytecode

import (
	"errors"
	"testing"
)

func TestValidateAggregateReduce(t *testing.T) {
	for _, test := range []struct {
		name           string
		dst, src, kind Operand
		valid          bool
	}{
		{"count", 1, 2, Operand(AggregateCount), true},
		{"average", 1, 2, Operand(AggregateAverage), true},
		{"aliased registers", 1, 1, Operand(AggregateSum), true},
		{"invalid destination", 3, 1, Operand(AggregateSum), false},
		{"invalid source", 1, 3, Operand(AggregateSum), false},
		{"constant source", 1, NewConstant(0), Operand(AggregateSum), false},
		{"negative kind", 1, 2, -1, false},
		{"unknown kind", 1, 2, Operand(AggregateAverage + 1), false},
	} {
		t.Run(test.name, func(t *testing.T) {
			program := withProgramMutation(func(p *Program) { p.Bytecode[0] = NewInstruction(OpAggregateReduce, test.dst, test.src, test.kind) })
			err := ValidateProgram(program)
			if test.valid && err != nil {
				t.Fatal(err)
			}

			if !test.valid && !errors.Is(err, ErrInvalidInstruction) {
				t.Fatalf("error = %v, want invalid instruction", err)
			}
		})
	}
}
