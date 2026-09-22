package bytecode

import (
	"errors"
	"math"
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

func TestAggregateSelectorNativeWidth(t *testing.T) {
	for _, selector := range []runtime.Int{math.MinInt64, math.MinInt32 - 1, 0, math.MaxInt32 + 1, math.MaxInt64} {
		program := &Program{
			ISAVersion: Version,
			Registers:  1,
			Constants:  []runtime.Value{selector},
			Bytecode: []Instruction{
				NewInstruction(OpLoadAggregateKey, NewRegister(0), NewRegister(0), NewConstant(0)),
				NewInstruction(OpReturn, NewRegister(0)),
			},
		}
		err := ValidateProgram(program)
		if selector < math.MinInt || selector > math.MaxInt {
			if !errors.Is(err, ErrInvalidInstruction) {
				t.Errorf("selector %d: %v", selector, err)
			}
		} else if err != nil {
			t.Errorf("representable selector %d: %v", selector, err)
		}
	}
}
