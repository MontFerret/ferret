package internal

import (
	"github.com/MontFerret/ferret/pkg/compiler/internal/core"
	"github.com/MontFerret/ferret/pkg/parser/fql"
	"github.com/MontFerret/ferret/pkg/runtime"
	"github.com/MontFerret/ferret/pkg/vm"
)

// initializeGrouping creates the KeyValue pair for collection, handling both grouping and value setup.
func (c *LoopCollectCompiler) initializeGrouping(grouping fql.ICollectGroupingContext) (*core.KV, []fql.ICollectSelectorContext) {
	var groupSelectors []fql.ICollectSelectorContext

	kv := core.NewKV(vm.NoopOperand, vm.NoopOperand)
	loop := c.ctx.Loops.Current()

	// Handle grouping key if present
	if grouping != nil {
		kv.Key, groupSelectors = c.compileGroupKeys(grouping)
	}

	// Setup value register and emit value from current loop
	if loop.Kind == core.ForInLoop {
		if loop.Value != vm.NoopOperand {
			kv.Value = loop.Value
		} else {
			kv.Value = c.ctx.Registers.Allocate(core.Temp)
			loop.EmitValue(kv.Value, c.ctx.Emitter)
		}
	} else {
		if loop.Key != vm.NoopOperand {
			kv.Value = loop.Key
		} else {
			kv.Value = c.ctx.Registers.Allocate(core.Temp)
			loop.EmitKey(kv.Value, c.ctx.Emitter)
		}
	}

	return kv, groupSelectors
}

// compileGroupKeys compiles the grouping keys from the CollectGroupingContext.
func (c *LoopCollectCompiler) compileGroupKeys(ctx fql.ICollectGroupingContext) (vm.Operand, []fql.ICollectSelectorContext) {
	selectors := ctx.AllCollectSelector()

	if len(selectors) == 0 {
		return vm.NoopOperand, selectors
	}

	var kvKeyReg vm.Operand

	if len(selectors) > 1 {
		// We create a sequence of Registers for the clauses
		// To pack them into an array
		selectorRegs := c.ctx.Registers.AllocateSequence(len(selectors))

		for i, selector := range selectors {
			reg := c.ctx.ExprCompiler.Compile(selector.Expression())
			c.ctx.Emitter.EmitAB(vm.OpMove, selectorRegs[i], reg)
			// Free the register after moving its value to the sequence register
			c.ctx.Registers.Free(reg)
		}

		kvKeyReg = c.ctx.Registers.Allocate(core.Temp)
		c.ctx.Emitter.EmitAs(vm.OpLoadArray, kvKeyReg, selectorRegs)
		c.ctx.Registers.FreeSequence(selectorRegs)
	} else {
		kvKeyReg = c.ctx.ExprCompiler.Compile(selectors[0].Expression())
	}

	return kvKeyReg, selectors
}

func (c *LoopCollectCompiler) compileGrouping(collectorType core.CollectorType, selectors []fql.ICollectSelectorContext) {
	loop := c.ctx.Loops.Current()

	if len(selectors) > 1 {
		variables := make([]vm.Operand, len(selectors))

		for i, selector := range selectors {
			name := selector.Identifier().GetText()

			if variables[i] == vm.NoopOperand {
				variables[i] = c.ctx.Symbols.DeclareLocal(name, core.TypeUnknown)
			}

			reg := c.selectGroupKey(collectorType, loop)

			c.ctx.Emitter.EmitABC(vm.OpLoadIndex, variables[i], reg, loadConstant(c.ctx, runtime.Int(i)))
		}

		// Free the register after moving its value to the variable
		for _, reg := range variables {
			c.ctx.Registers.Free(reg)
		}
	} else {
		// Get the variable name
		name := selectors[0].Identifier().GetText()
		// If we have a single selector, we can just use the loops' register directly
		c.ctx.Symbols.AssignLocal(name, core.TypeUnknown, c.selectGroupKey(collectorType, loop))
	}
}

func (c *LoopCollectCompiler) selectGroupKey(collectorType core.CollectorType, loop *core.Loop) vm.Operand {
	switch collectorType {
	case core.CollectorTypeKeyGroup, core.CollectorTypeKeyCounter:
		return loop.Key
	default:
		return loop.Value
	}
}
