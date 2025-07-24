package internal

import (
	"github.com/MontFerret/ferret/pkg/compiler/internal/core"
	"github.com/MontFerret/ferret/pkg/parser/fql"
	"github.com/MontFerret/ferret/pkg/runtime"
	"github.com/MontFerret/ferret/pkg/vm"
)

// initializeGrouping creates the KeyValue pair for collection, handling both grouping and value setup.
func (c *LoopCollectCompiler) initializeGrouping(grouping fql.ICollectGroupingContext) (*core.KV, []*core.CollectSelector) {
	var groupSelectors []*core.CollectSelector

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
func (c *LoopCollectCompiler) compileGroupKeys(ctx fql.ICollectGroupingContext) (vm.Operand, []*core.CollectSelector) {
	selectors := ctx.AllCollectSelector()

	if len(selectors) == 0 {
		return vm.NoopOperand, nil
	}

	var kvKeyReg vm.Operand
	var collectSelectors []*core.CollectSelector

	if len(selectors) > 1 {
		// We create a sequence of Registers for the clauses
		// To pack them into an array
		collectSelectors = make([]*core.CollectSelector, len(selectors))
		selectorRegs := c.ctx.Registers.AllocateSequence(len(selectors))

		for i, selector := range selectors {
			reg := c.ctx.ExprCompiler.Compile(selector.Expression())
			c.ctx.Emitter.EmitAB(vm.OpMove, selectorRegs[i], reg)
			// Free the register after moving its value to the sequence register
			c.ctx.Registers.Free(reg)

			collectSelectors[i] = core.NewCollectSelector(runtime.String(selector.Identifier().GetText()))
		}

		kvKeyReg = c.ctx.Registers.Allocate(core.Temp)
		c.ctx.Emitter.EmitAs(vm.OpLoadArray, kvKeyReg, selectorRegs)
		c.ctx.Registers.FreeSequence(selectorRegs)
	} else {
		selector := selectors[0]
		kvKeyReg = c.ctx.ExprCompiler.Compile(selector.Expression())
		collectSelectors = []*core.CollectSelector{core.NewCollectSelector(runtime.String(selector.Identifier().GetText()))}
	}

	return kvKeyReg, collectSelectors
}

func (c *LoopCollectCompiler) compileGrouping(spec *core.CollectorSpec) {
	loop := c.ctx.Loops.Current()

	if len(spec.GroupSelectors()) > 1 {
		variables := make([]vm.Operand, len(spec.GroupSelectors()))

		for i, selector := range spec.GroupSelectors() {
			name := selector.Name()

			if variables[i] == vm.NoopOperand {
				variables[i] = c.ctx.Symbols.DeclareLocal(name.String(), core.TypeUnknown)
			}

			reg := c.selectGroupKey(spec.Type(), loop)

			c.ctx.Emitter.EmitABC(vm.OpLoadIndex, variables[i], reg, loadConstant(c.ctx, runtime.Int(i)))
		}

		// Free the register after moving its value to the variable
		for _, reg := range variables {
			c.ctx.Registers.Free(reg)
		}
	} else {
		// Get the variable name
		name := spec.GroupSelectors()[0].Name()
		// If we have a single selector, we can just use the loops' register directly
		c.ctx.Symbols.AssignLocal(name.String(), core.TypeUnknown, c.selectGroupKey(spec.Type(), loop))
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
