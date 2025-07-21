package internal

import (
	"github.com/antlr4-go/antlr/v4"

	"github.com/MontFerret/ferret/pkg/compiler/internal/core"
	"github.com/MontFerret/ferret/pkg/parser/fql"
	"github.com/MontFerret/ferret/pkg/runtime"
	"github.com/MontFerret/ferret/pkg/vm"
)

type (
	LoopCollectCompiler struct {
		ctx *CompilerContext
	}
)

func NewCollectCompiler(ctx *CompilerContext) *LoopCollectCompiler {
	return &LoopCollectCompiler{ctx: ctx}
}

func (c *LoopCollectCompiler) Compile(ctx fql.ICollectClauseContext) {
	aggregator := ctx.CollectAggregator()
	collectorType, groupSelectors := c.compileCollect(ctx, aggregator != nil)

	// Aggregation loop
	if aggregator != nil {
		c.compileAggregation(aggregator, len(groupSelectors) > 0)
	}

	if len(groupSelectors) > 0 {
		// Now we are defining new variables for the group selectors
		c.compileGroupSelectorVariables(collectorType, groupSelectors, aggregator != nil)
	}
}

func (c *LoopCollectCompiler) compileCollect(ctx fql.ICollectClauseContext, aggregation bool) (core.CollectorType, []fql.ICollectSelectorContext) {
	grouping := ctx.CollectGrouping()
	counter := ctx.CollectCounter()

	// We gather keys and values for the collector.
	kv, groupSelectors := c.initializeGrouping(grouping)
	projectionVarName, collectorType := c.initializeProjection(ctx, kv, counter, grouping != nil)

	// If we use aggregators, we need to collect group items by key
	if aggregation && collectorType != core.CollectorTypeKeyGroup {
		// We need to patch the loop result to be a collector
		collectorType = core.CollectorTypeKeyGroup
	}

	c.finalizeCollector(collectorType, kv)
	loop := c.ctx.Loops.Current()

	// We no longer need KV, so we free registers
	c.ctx.Registers.Free(kv.Key)
	c.ctx.Registers.Free(kv.Value)

	// If we are using a projection, we need to ensure the loop is set to ForInLoop
	if loop.Kind != core.ForInLoop {
		loop.Kind = core.ForInLoop
	}

	if loop.Value == vm.NoopOperand {
		loop.Value = c.ctx.Registers.Allocate(core.Temp)
	}

	if loop.Key == vm.NoopOperand {
		loop.Key = c.ctx.Registers.Allocate(core.Temp)
	}

	// If the projection is used, we allocate a new register for the variable and put the iterator's value into it
	if projectionVarName != "" {
		// Now we need to expand group variables from the dataset
		loop.ValueName = projectionVarName
		c.ctx.Symbols.AssignLocal(loop.ValueName, core.TypeUnknown, loop.Value)
	}

	loop.EmitInitialization(c.ctx.Registers, c.ctx.Emitter, c.ctx.Loops.Depth())

	return collectorType, groupSelectors
}

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

// initializeProjection handles the projection setup for group variables and counters.
// Returns the projection variable name and the appropriate collector type.
func (c *LoopCollectCompiler) initializeProjection(ctx fql.ICollectClauseContext, kv *core.KV, counter fql.ICollectCounterContext, hasGrouping bool) (string, core.CollectorType) {
	projectionVariableName := ""
	collectorType := core.CollectorTypeKey

	// Handle group variable projection
	if groupVar := ctx.CollectGroupVariable(); groupVar != nil {
		projectionVariableName = c.compileGroupVariableProjection(kv, groupVar)
		collectorType = core.CollectorTypeKeyGroup
		return projectionVariableName, collectorType
	}

	// Handle counter projection
	if counter != nil {
		projectionVariableName = counter.Identifier().GetText()
		collectorType = c.determineCounterCollectorType(hasGrouping)
	}

	return projectionVariableName, collectorType
}

// determineCounterCollectorType returns the appropriate collector type for counter operations.
func (c *LoopCollectCompiler) determineCounterCollectorType(hasGrouping bool) core.CollectorType {
	if hasGrouping {
		return core.CollectorTypeKeyCounter
	}

	return core.CollectorTypeCounter
}

// compileGroupVariableProjection processes group variable projections (both default and custom).
func (c *LoopCollectCompiler) compileGroupVariableProjection(kv *core.KV, groupVar fql.ICollectGroupVariableContext) string {
	// Handle default projection (identifier)
	if identifier := groupVar.Identifier(); identifier != nil {
		return c.compileDefaultGroupProjection(kv, identifier, groupVar.CollectGroupVariableKeeper())
	}

	// Handle custom projection (selector expression)
	if selector := groupVar.CollectSelector(); selector != nil {
		return c.compileCustomGroupProjection(kv, selector)
	}

	return ""
}

func (c *LoopCollectCompiler) compileGroupSelectorVariables(collectorType core.CollectorType, selectors []fql.ICollectSelectorContext, isAggregation bool) {
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

func (c *LoopCollectCompiler) compileDefaultGroupProjection(kv *core.KV, identifier antlr.TerminalNode, keeper fql.ICollectGroupVariableKeeperContext) string {
	if keeper == nil {
		variables := c.ctx.Symbols.LocalVariables()
		scope := core.NewScopeProjection(c.ctx.Registers, c.ctx.Emitter, c.ctx.Symbols, variables)
		scope.EmitAsObject(kv.Value)
	} else {
		variables := keeper.AllIdentifier()
		seq := c.ctx.Registers.AllocateSequence(len(variables) * 2)

		for i, j := 0, 0; i < len(variables); i, j = i+1, j+2 {
			varName := variables[i].GetText()
			loadConstantTo(c.ctx, runtime.String(varName), seq[j])

			variable, _, found := c.ctx.Symbols.Resolve(varName)

			if !found {
				panic("variable not found: " + varName)
			}

			c.ctx.Emitter.EmitAB(vm.OpMove, seq[j+1], variable)
		}

		c.ctx.Emitter.EmitAs(vm.OpLoadObject, kv.Value, seq)
		c.ctx.Registers.FreeSequence(seq)
	}

	return identifier.GetText()
}

func (c *LoopCollectCompiler) compileCustomGroupProjection(kv *core.KV, selector fql.ICollectSelectorContext) string {
	selectorReg := c.ctx.ExprCompiler.Compile(selector.Expression())
	c.ctx.Emitter.EmitMove(kv.Value, selectorReg)
	c.ctx.Registers.Free(selectorReg)

	return selector.Identifier().GetText()
}

func (c *LoopCollectCompiler) selectGroupKey(collectorType core.CollectorType, loop *core.Loop) vm.Operand {
	switch collectorType {
	case core.CollectorTypeKeyGroup, core.CollectorTypeKeyCounter:
		return loop.Key
	default:
		return loop.Value
	}
}

func (c *LoopCollectCompiler) finalizeCollector(collectorType core.CollectorType, kv *core.KV) {
	loop := c.ctx.Loops.Current()
	// We replace DataSet initialization with Collector initialization
	dst := loop.PatchDestinationAx(c.ctx.Registers, c.ctx.Emitter, vm.OpDataSetCollector, int(collectorType))
	c.ctx.Emitter.EmitABC(vm.OpPushKV, dst, kv.Key, kv.Value)
	loop.EmitFinalization(c.ctx.Emitter)

	// Move the collector to the next loop source
	c.ctx.Emitter.EmitMove(loop.Src, dst)
}
