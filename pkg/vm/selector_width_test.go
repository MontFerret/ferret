package vm

import (
	"errors"
	"math"
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/bytecode"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
	"github.com/MontFerret/ferret/v2/pkg/vm/internal/data"
)

func TestAggregateSelectorNativeWidth(t *testing.T) {
	for _, selector := range []runtime.Int{math.MinInt64, math.MinInt32 - 1, 0, math.MaxInt32 + 1, math.MaxInt64} {
		program := newTestProgram(1, []runtime.Value{selector},
			bytecode.NewInstruction(bytecode.OpLoadAggregateKey, bytecode.NewRegister(0), bytecode.NewRegister(0), bytecode.NewConstant(0)),
			bytecode.NewInstruction(bytecode.OpReturn, bytecode.NewRegister(0)),
		)
		instance := mustNewVM(t, program, WithPanicPolicy(PanicRecover))
		result, err := instance.Run(t.Context(), NewDefaultEnvironment())
		if result != nil {
			t.Cleanup(func() { _ = result.Close() })
		}
		if selector < math.MinInt || selector > math.MaxInt {
			if !errors.Is(err, runtime.ErrRange) || result != nil {
				t.Errorf("selector %d: result %v, error %v", selector, result, err)
			}
		} else if err != nil || result == nil {
			t.Errorf("representable selector %d: result %v, error %v", selector, result, err)
		} else {
			key, ok := result.Root().(*data.AggregateKey)
			if !ok || runtime.Int(key.SelectorIndex()) != selector {
				t.Errorf("selector changed: %v", result.Root())
			}
		}
	}
}
