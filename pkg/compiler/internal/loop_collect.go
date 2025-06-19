package internal

import (
	"github.com/antlr4-go/antlr/v4"

	"github.com/MontFerret/ferret/pkg/compiler/internal/core"
	"github.com/MontFerret/ferret/pkg/parser/fql"
	"github.com/MontFerret/ferret/pkg/runtime"
	"github.com/MontFerret/ferret/pkg/vm"
)

type LoopCollectCompiler struct {
	ctx *CompilerContext
}

func NewCollectCompiler(ctx *CompilerContext) *LoopCollectCompiler {
	return &LoopCollectCompiler{ctx: ctx}
}

func (cc *LoopCollectCompiler) Compile(ctx fql.ICollectClauseContext) {
	aggregator := ctx.CollectAggregator()
	kv, groupSelectors := cc.compileCollect(ctx, aggregator != nil)

	// Aggregation loop
	if aggregator != nil {
		cc.compileAggregation(aggregator, len(groupSelectors) > 0)
	}

	if len(groupSelectors) > 0 {
		// Now we are defining new variables for the group selectors
		cc.compileGroupSelectorVariables(groupSelectors, kv, aggregator != nil)
	}
}

func (cc *LoopCollectCompiler) compileCollect(ctx fql.ICollectClauseContext, aggregation bool) (*core.KV, []fql.ICollectSelectorContext) {
	grouping := ctx.CollectGrouping()
	counter := ctx.CollectCounter()

	if grouping == nil && counter == nil {
		return core.NewKV(vm.NoopOperand, vm.NoopOperand), nil
	}

	loop := cc.ctx.Loops.Current()

	kv, groupSelectors := cc.initializeCollector(grouping)
	projectionVariableName, collectorType := cc.initializeProjection(ctx, loop, kv, counter, grouping != nil)

	// If we use aggregators, we need to collect group items by key
	if aggregation && collectorType != core.CollectorTypeKeyGroup {
		// We need to patch the loop result to be a collector
		collectorType = core.CollectorTypeKeyGroup
	}

	cc.finalizeCollector(loop, collectorType, kv)

	// If the projection is used, we allocate a new register for the variable and put the iterator's value into it
	if projectionVariableName != "" {
		// Now we need to expand group variables from the dataset
		loop.DeclareValueVar(projectionVariableName, cc.ctx.Symbols)
		loop.EmitInitialization(cc.ctx.Registers, cc.ctx.Emitter)

		loop.EmitKey(kv.Value, cc.ctx.Emitter)
	} else {
		loop.EmitInitialization(cc.ctx.Registers, cc.ctx.Emitter)
		loop.EmitKey(kv.Key, cc.ctx.Emitter)
	}

	return kv, groupSelectors
}

// initializeKeyValue creates the KeyValue pair for collection, handling both grouping and value setup.
func (cc *LoopCollectCompiler) initializeCollector(grouping fql.ICollectGroupingContext) (*core.KV, []fql.ICollectSelectorContext) {
	var groupSelectors []fql.ICollectSelectorContext

	kv := core.NewKV(vm.NoopOperand, vm.NoopOperand)
	loop := cc.ctx.Loops.Current()

	// Handle grouping key if present
	if grouping != nil {
		keyReg, selectors := cc.compileGrouping(grouping)
		kv.Key = keyReg
		groupSelectors = selectors
	}

	// Setup value register and emit value from current loop
	kv.Value = cc.ctx.Registers.Allocate(core.Temp)
	loop.EmitValue(kv.Value, cc.ctx.Emitter)

	return kv, groupSelectors
}

func (cc *LoopCollectCompiler) finalizeCollector(loop *core.Loop, collectorType core.CollectorType, kv *core.KV) {
	// We replace DataSet initialization with Collector initialization
	cc.ctx.Emitter.PatchSwapAx(loop.DstPos, vm.OpDataSetCollector, loop.Dst, int(collectorType))
	cc.ctx.Emitter.EmitABC(vm.OpPushKV, loop.Dst, kv.Key, kv.Value)
	loop.EmitFinalization(cc.ctx.Emitter)

	cc.ctx.Emitter.EmitMove(loop.Src, loop.Dst)

	cc.ctx.Registers.Free(loop.Value)
	cc.ctx.Registers.Free(loop.Key)
	loop.Value = kv.Value
	loop.Key = vm.NoopOperand
}

// initializeProjection handles the projection setup for group variables and counters.
// Returns the projection variable name and the appropriate collector type.
func (cc *LoopCollectCompiler) initializeProjection(ctx fql.ICollectClauseContext, loop *core.Loop, kv *core.KV, counter fql.ICollectCounterContext, hasGrouping bool) (string, core.CollectorType) {
	projectionVariableName := ""
	collectorType := core.CollectorTypeKey

	// Handle group variable projection
	if groupVar := ctx.CollectGroupVariable(); groupVar != nil {
		projectionVariableName = cc.compileGroupVariableProjection(loop, kv, groupVar)
		collectorType = core.CollectorTypeKeyGroup
		return projectionVariableName, collectorType
	}

	// Handle counter projection
	if counter != nil {
		projectionVariableName = counter.Identifier().GetText()
		collectorType = cc.determineCounterCollectorType(hasGrouping)
	}

	return projectionVariableName, collectorType
}

// determineCounterCollectorType returns the appropriate collector type for counter operations.
func (cc *LoopCollectCompiler) determineCounterCollectorType(hasGrouping bool) core.CollectorType {
	if hasGrouping {
		return core.CollectorTypeKeyCounter
	}

	return core.CollectorTypeCounter
}

func (cc *LoopCollectCompiler) compileGrouping(ctx fql.ICollectGroupingContext) (vm.Operand, []fql.ICollectSelectorContext) {
	selectors := ctx.AllCollectSelector()

	if len(selectors) == 0 {
		return vm.NoopOperand, selectors
	}

	var kvKeyReg vm.Operand

	if len(selectors) > 1 {
		// We create a sequence of Registers for the clauses
		// To pack them into an array
		selectorRegs := cc.ctx.Registers.AllocateSequence(len(selectors))

		for i, selector := range selectors {
			reg := cc.ctx.ExprCompiler.Compile(selector.Expression())
			cc.ctx.Emitter.EmitAB(vm.OpMove, selectorRegs[i], reg)
			// Free the register after moving its value to the sequence register
			cc.ctx.Registers.Free(reg)
		}

		kvKeyReg = cc.ctx.Registers.Allocate(core.Temp)
		cc.ctx.Emitter.EmitAs(vm.OpList, kvKeyReg, selectorRegs)
		cc.ctx.Registers.FreeSequence(selectorRegs)
	} else {
		kvKeyReg = cc.ctx.ExprCompiler.Compile(selectors[0].Expression())
	}

	return kvKeyReg, selectors
}

// compileGroupVariableProjection processes group variable projections (both default and custom).
func (cc *LoopCollectCompiler) compileGroupVariableProjection(loop *core.Loop, kv *core.KV, groupVar fql.ICollectGroupVariableContext) string {
	// Handle default projection (identifier)
	if identifier := groupVar.Identifier(); identifier != nil {
		return cc.compileDefaultGroupProjection(loop, kv, identifier, groupVar.CollectGroupVariableKeeper())
	}

	// Handle custom projection (selector expression)
	if selector := groupVar.CollectSelector(); selector != nil {
		return cc.compileCustomGroupProjection(loop, kv, selector)
	}

	return ""
}

func (cc *LoopCollectCompiler) compileGroupSelectorVariables(selectors []fql.ICollectSelectorContext, kv *core.KV, isAggregation bool) {
	if len(selectors) > 1 {
		variables := make([]vm.Operand, len(selectors))

		for i, selector := range selectors {
			name := selector.Identifier().GetText()

			if variables[i] == vm.NoopOperand {
				variables[i] = cc.ctx.Symbols.DeclareLocal(name)
			}

			reg := kv.Value

			if isAggregation {
				reg = kv.Key
			}

			cc.ctx.Emitter.EmitABC(vm.OpLoadIndex, variables[i], reg, loadConstant(cc.ctx, runtime.Int(i)))
		}

		// Free the register after moving its value to the variable
		for _, reg := range variables {
			cc.ctx.Registers.Free(reg)
		}
	} else {
		// Get the variable name
		name := selectors[0].Identifier().GetText()
		// Define a variable for each selector
		varReg := cc.ctx.Symbols.DeclareLocal(name)
		reg := cc.selectGroupKey(isAggregation, kv)

		// If we have a single selector, we can just move the value
		cc.ctx.Emitter.EmitAB(vm.OpMove, varReg, reg)
	}
}

func (cc *LoopCollectCompiler) compileDefaultGroupProjection(loop *core.Loop, kv *core.KV, identifier antlr.TerminalNode, keeper fql.ICollectGroupVariableKeeperContext) string {
	if keeper == nil {
		seq := cc.ctx.Registers.AllocateSequence(2) // Key and Value for Map

		// TODO: Review this. It's quite a questionable ArrangoDB feature of wrapping group items by a nested object
		// We will keep it for now for backward compatibility.
		loadConstantTo(cc.ctx, runtime.String(loop.ValueName), seq[0]) // Map key
		cc.ctx.Emitter.EmitAB(vm.OpMove, seq[1], kv.Value)             // Map value
		cc.ctx.Emitter.EmitAs(vm.OpMap, kv.Value, seq)

		cc.ctx.Registers.FreeSequence(seq)
	} else {
		variables := keeper.AllIdentifier()
		seq := cc.ctx.Registers.AllocateSequence(len(variables) * 2)

		for i, j := 0, 0; i < len(variables); i, j = i+1, j+2 {
			varName := variables[i].GetText()
			loadConstantTo(cc.ctx, runtime.String(varName), seq[j])

			variable, _, found := cc.ctx.Symbols.Resolve(varName)

			if !found {
				panic("variable not found: " + varName)
			}

			cc.ctx.Emitter.EmitAB(vm.OpMove, seq[j+1], variable)
		}

		cc.ctx.Emitter.EmitAs(vm.OpMap, kv.Value, seq)
		cc.ctx.Registers.FreeSequence(seq)
	}

	return identifier.GetText()
}

func (cc *LoopCollectCompiler) compileCustomGroupProjection(_ *core.Loop, kv *core.KV, selector fql.ICollectSelectorContext) string {
	selectorReg := cc.ctx.ExprCompiler.Compile(selector.Expression())
	cc.ctx.Emitter.EmitMove(kv.Value, selectorReg)
	cc.ctx.Registers.Free(selectorReg)

	return selector.Identifier().GetText()
}

func (cc *LoopCollectCompiler) selectGroupKey(isAggregation bool, kv *core.KV) vm.Operand {
	if isAggregation {
		return kv.Key
	}

	return kv.Value
}
