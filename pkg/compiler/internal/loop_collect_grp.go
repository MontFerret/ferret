package internal

import (
	"github.com/MontFerret/ferret/pkg/compiler/internal/core"
	"github.com/MontFerret/ferret/pkg/parser/fql"
	"github.com/MontFerret/ferret/pkg/runtime"
	"github.com/MontFerret/ferret/pkg/vm"
)

// initializeGrouping creates the KeyValue pair for collection, handling both grouping and value setup.
// It processes the grouping context to extract group keys and sets up the value register based on the loop type.
// Returns a KV struct containing key and value registers, and a slice of CollectSelectors for grouping.
func (c *LoopCollectCompiler) initializeGrouping(grouping fql.ICollectGroupingContext) (*core.KV, []*core.CollectSelector) {
	var groupSelectors []*core.CollectSelector

	// Initialize key-value pair with no-op operands
	kv := core.NewKV(vm.NoopOperand, vm.NoopOperand)
	loop := c.ctx.Loops.Current()

	// Handle grouping key if present
	if grouping != nil {
		kv.Key, groupSelectors = c.compileGroupKeys(grouping)
	}

	// Setup value register and emit value from current loop
	// The behavior differs based on the loop type (FOR IN vs FOR)
	if loop.Kind == core.ForInLoop {
		if loop.Value != vm.NoopOperand {
			// Reuse existing value register if available
			kv.Value = loop.Value
		} else {
			// Allocate new register and emit value instruction
			kv.Value = c.ctx.Registers.Allocate(core.Temp)
			loop.EmitValue(kv.Value, c.ctx.Emitter)
		}
	} else {
		if loop.Key != vm.NoopOperand {
			// For non-ForInLoop, use key as value if available
			kv.Value = loop.Key
		} else {
			// Allocate new register and emit key instruction
			kv.Value = c.ctx.Registers.Allocate(core.Temp)
			loop.EmitKey(kv.Value, c.ctx.Emitter)
		}
	}

	return kv, groupSelectors
}

// compileGroupKeys compiles the grouping keys from the CollectGroupingContext.
// It processes the selectors in the grouping context and creates the appropriate VM instructions.
// For multiple selectors, it creates an array of values. For a single selector, it uses the value directly.
// Returns the register containing the key value and a slice of CollectSelectors for later use.
func (c *LoopCollectCompiler) compileGroupKeys(ctx fql.ICollectGroupingContext) (vm.Operand, []*core.CollectSelector) {
	selectors := ctx.AllCollectSelector()

	// If no selectors are present, return no-op operand
	if len(selectors) == 0 {
		return vm.NoopOperand, nil
	}

	var kvKeyReg vm.Operand
	var collectSelectors []*core.CollectSelector

	if len(selectors) > 1 {
		// Handle multiple selectors by creating an array
		// We create a sequence of Registers for the clauses
		// To pack them into an array
		collectSelectors = make([]*core.CollectSelector, len(selectors))
		selectorRegs := c.ctx.Registers.AllocateSequence(len(selectors))

		// Process each selector expression and store in sequence
		for i, selector := range selectors {
			reg := c.ctx.ExprCompiler.Compile(selector.Expression())
			c.ctx.Emitter.EmitAB(vm.OpMove, selectorRegs[i], reg)
			// Free the register after moving its value to the sequence register
			c.ctx.Registers.Free(reg)

			// Create a CollectSelector for each selector with its identifier
			collectSelectors[i] = core.NewCollectSelector(runtime.String(selector.Identifier().GetText()))
		}

		// Create an array from the sequence of registers
		kvKeyReg = c.ctx.Registers.Allocate(core.Temp)
		c.ctx.Emitter.EmitAs(vm.OpLoadArray, kvKeyReg, selectorRegs)
		c.ctx.Registers.FreeSequence(selectorRegs)
	} else {
		// Handle single selector case - simpler, no need for array
		selector := selectors[0]
		kvKeyReg = c.ctx.ExprCompiler.Compile(selector.Expression())
		collectSelectors = []*core.CollectSelector{core.NewCollectSelector(runtime.String(selector.Identifier().GetText()))}
	}

	return kvKeyReg, collectSelectors
}

// finalizeGrouping processes the group selectors and creates local variables for them.
// It handles both multiple selectors (as array elements) and single selectors differently.
func (c *LoopCollectCompiler) finalizeGrouping(spec *core.Collector) {
	loop := c.ctx.Loops.Current()

	if len(spec.GroupSelectors()) > 1 {
		// Handle multiple group selectors
		variables := make([]vm.Operand, len(spec.GroupSelectors()))

		// Process each selector and create a local variable for it
		for i, selector := range spec.GroupSelectors() {
			name := selector.Name()

			// Declare a local variable for the selector if not already done
			if variables[i] == vm.NoopOperand {
				// TODO: Handle error if the variable already exists
				reg, _ := c.ctx.Symbols.DeclareLocal(name.String(), core.TypeUnknown)
				variables[i] = reg
			}

			// Get the appropriate register (key or value) based on collector type
			reg := c.selectGroupKey(spec.Type(), loop)

			// Load the value at index i from the array in reg into the variable
			c.ctx.Emitter.EmitABC(vm.OpLoadIndex, variables[i], reg, loadConstant(c.ctx, runtime.Int(i)))
		}

		// Free the register after moving its value to the variable
		for _, reg := range variables {
			c.ctx.Registers.Free(reg)
		}
	} else {
		// Handle single group selector - simpler case
		// Get the variable name
		name := spec.GroupSelectors()[0].Name()
		// If we have a single selector, we can just use the loops' register directly
		c.ctx.Symbols.AssignLocal(name.String(), core.TypeUnknown, c.selectGroupKey(spec.Type(), loop))
	}
}

// selectGroupKey determines which register (key or value) to use based on the collector type.
// Different collector types require different registers to be used as the group key.
func (c *LoopCollectCompiler) selectGroupKey(collectorType core.CollectorType, loop *core.Loop) vm.Operand {
	switch collectorType {
	case core.CollectorTypeKeyGroup, core.CollectorTypeKeyCounter:
		// For key-based collectors, use the key register
		return loop.Key
	default:
		// For other collector types, use the value register
		return loop.Value
	}
}
