package internal

import (
	"github.com/MontFerret/ferret/v2/pkg/bytecode"
	"github.com/MontFerret/ferret/v2/pkg/compiler/internal/core"
	"github.com/MontFerret/ferret/v2/pkg/parser/fql"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

// initializeGrouping creates the KeyValue pair for collection, handling both grouping and value setup.
// It processes the grouping context to extract group keys and sets up the value register based on the loop type.
// Returns a KV struct containing key and value registers, and a slice of CollectSelectors for grouping.
func (c *CollectCompiler) initializeGrouping(grouping fql.ICollectGroupingContext) (*core.KV, []*core.CollectSelector) {
	var groupSelectors []*core.CollectSelector

	// Initialize key-value pair with no-op operands
	kv := core.NewKV(bytecode.NoopOperand, bytecode.NoopOperand)
	loop := c.ctx.Loops.Current()

	// Handle grouping key if present
	if grouping != nil {
		kv.Key, groupSelectors = c.compileGroupKeys(grouping)
	}

	// Setup value register and emit value from current loop
	// The behavior differs based on the loop type (FOR IN vs FOR)
	if loop.Kind == core.ForInLoop {
		if loop.Value != bytecode.NoopOperand {
			// Reuse existing value register if available
			kv.Value = loop.Value
		} else {
			// Allocate new register and emit value instruction
			kv.Value = c.ctx.Registers.Allocate()
			loop.EmitValue(kv.Value, c.ctx.Emitter)
		}
	} else {
		if loop.Key != bytecode.NoopOperand {
			// For non-ForInLoop, use key as value if available
			kv.Value = loop.Key
		} else {
			// Allocate new register and emit key instruction
			kv.Value = c.ctx.Registers.Allocate()
			loop.EmitKey(kv.Value, c.ctx.Emitter)
		}
	}

	return kv, groupSelectors
}

// compileGroupKeys compiles the grouping keys from the CollectGroupingContext.
// It processes the selectors in the grouping context and creates the appropriate VM instructions.
// For multiple selectors, it creates an array of values. For a single selector, it uses the value directly.
// Returns the register containing the key value and a slice of CollectSelectors for later use.
func (c *CollectCompiler) compileGroupKeys(ctx fql.ICollectGroupingContext) (bytecode.Operand, []*core.CollectSelector) {
	selectors := ctx.AllCollectSelector()

	// If no selectors are present, return no-op operand
	if len(selectors) == 0 {
		return bytecode.NoopOperand, nil
	}

	var kvKeyReg bytecode.Operand
	var collectSelectors []*core.CollectSelector

	if len(selectors) > 1 {
		// Handle multiple selectors by creating an array
		collectSelectors = make([]*core.CollectSelector, len(selectors))
		kvKeyReg = c.ctx.Registers.Allocate()
		c.ctx.Emitter.EmitArray(kvKeyReg, len(selectors))

		// Process each selector expression and push into the array
		for i, selector := range selectors {
			reg := c.exprs.Compile(selector.Expression())
			c.ctx.Emitter.EmitArrayPush(kvKeyReg, reg)

			// Create a CollectSelector for each selector with its identifier
			collectSelectors[i] = core.NewCollectSelector(runtime.String(textOfBindingIdentifier(selector.BindingIdentifier())), selector)
		}
	} else {
		// Handle single selector case - simpler, no need for array
		selector := selectors[0]
		kvKeyReg = c.exprs.Compile(selector.Expression())
		collectSelectors = []*core.CollectSelector{core.NewCollectSelector(runtime.String(textOfBindingIdentifier(selector.BindingIdentifier())), selector)}
	}

	return kvKeyReg, collectSelectors
}

// finalizeGrouping processes the group selectors and creates local variables for them.
// It handles both multiple selectors (as array elements) and single selectors differently.
func (c *CollectCompiler) finalizeGrouping(spec *core.Collector) {
	loop := c.ctx.Loops.Current()
	// Depending on collector mode, grouped keys are emitted either as iterator key or value.
	// selectGroupKey keeps this branching in one place.
	groupKeyReg := c.selectGroupKey(spec.Type(), loop)

	if len(spec.GroupSelectors()) > 1 {
		// Handle multiple group selectors.
		for i, selector := range spec.GroupSelectors() {
			name := selector.Name()
			reg := c.declareLocalOrReport(selector.Context(), name.String(), core.TypeUnknown)

			// Load the value at index i from the group key array into the local variable.
			c.ctx.Emitter.EmitABC(bytecode.OpLoadIndex, reg, groupKeyReg, c.facts.LoadConstant(runtime.Int(i)))
		}
	} else {
		// Handle single group selector - simpler case
		// Get the variable name
		name := spec.GroupSelectors()[0].Name()
		// If we have a single selector, we can just use the loop register directly
		c.assignLocalOrReport(spec.GroupSelectors()[0].Context(), name.String(), core.TypeUnknown, groupKeyReg)
	}
}

// selectGroupKey determines which register (key or value) to use based on the collector type.
// Different collector types require different registers to be used as the group key.
func (c *CollectCompiler) selectGroupKey(collectorType bytecode.CollectorType, loop *core.Loop) bytecode.Operand {
	switch collectorType {
	case bytecode.CollectorTypeKeyGroup, bytecode.CollectorTypeKeyCounter, bytecode.CollectorTypeAggregateGroup:
		// For key-based collectors, use the key register
		return loop.Key
	default:
		// For other collector types, use the value register
		return loop.Value
	}
}
