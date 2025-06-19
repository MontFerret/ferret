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
	kvKeyReg, kvValReg, groupSelectors := cc.compileCollect(ctx, aggregator != nil)

	// Aggregation loop
	if aggregator != nil {
		cc.compileAggregation(aggregator, len(groupSelectors) > 0)
	}

	if len(groupSelectors) > 0 {
		// Now we are defining new variables for the group selectors
		cc.compileGroupSelectorVariables(groupSelectors, kvKeyReg, kvValReg, aggregator != nil)
	}
}

func (cc *LoopCollectCompiler) compileCollect(ctx fql.ICollectClauseContext, aggregation bool) (vm.Operand, vm.Operand, []fql.ICollectSelectorContext) {
	var kvKeyReg, kvValReg vm.Operand
	var groupSelectors []fql.ICollectSelectorContext
	grouping := ctx.CollectGrouping()
	counter := ctx.CollectCounter()

	if grouping == nil && counter == nil {
		return kvKeyReg, kvValReg, groupSelectors
	}

	loop := cc.ctx.Loops.Current()

	if grouping != nil {
		kvKeyReg, groupSelectors = cc.compileGrouping(grouping)
	}

	kvValReg = cc.ctx.Registers.Allocate(core.Temp)
	loop.EmitValue(kvValReg, cc.ctx.Emitter)

	var projectionVariableName string
	collectorType := core.CollectorTypeKey

	// If we have a collect group variable, we need to project it
	if groupVar := ctx.CollectGroupVariable(); groupVar != nil {
		// Projection can be either a default projection (identifier) or a custom projection (selector expression)
		if identifier := groupVar.Identifier(); identifier != nil {
			projectionVariableName = cc.compileDefaultGroupProjection(loop, kvValReg, identifier, groupVar.CollectGroupVariableKeeper())
		} else if selector := groupVar.CollectSelector(); selector != nil {
			projectionVariableName = cc.compileCustomGroupProjection(loop, kvValReg, selector)
		}

		collectorType = core.CollectorTypeKeyGroup
	} else if counter != nil {
		projectionVariableName = counter.Identifier().GetText()

		if grouping != nil {
			collectorType = core.CollectorTypeKeyCounter
		} else {
			collectorType = core.CollectorTypeCounter
		}
	}

	// If we use aggregators, we need to collect group items by key
	if aggregation && collectorType != core.CollectorTypeKeyGroup {
		// We need to patch the loop result to be a collector
		collectorType = core.CollectorTypeKeyGroup
	}

	// We replace DataSet initialization with Collector initialization
	cc.ctx.Emitter.PatchSwapAx(loop.ResultPos, vm.OpDataSetCollector, loop.Result, int(collectorType))
	cc.ctx.Emitter.EmitABC(vm.OpPushKV, loop.Result, kvKeyReg, kvValReg)
	loop.EmitFinalization(cc.ctx.Emitter)

	cc.ctx.Emitter.EmitMove(loop.Src, loop.Result)

	cc.ctx.Registers.Free(loop.Value)
	cc.ctx.Registers.Free(loop.Key)
	loop.Value = kvValReg
	loop.Key = vm.NoopOperand

	// If the projection is used, we allocate a new register for the variable and put the iterator's value into it
	if projectionVariableName != "" {
		// Now we need to expand group variables from the dataset
		loop.DeclareValueVar(projectionVariableName, cc.ctx.Symbols)
		loop.EmitInitialization(cc.ctx.Registers, cc.ctx.Emitter)

		//cc.ctx.LoopCompiler.EmitLoopBegin(loop)

		loop.EmitKey(kvValReg, cc.ctx.Emitter)
		loop.BindValueVar(cc.ctx.Emitter)
	} else {
		//cc.ctx.LoopCompiler.EmitLoopBegin(loop)
		loop.EmitInitialization(cc.ctx.Registers, cc.ctx.Emitter)

		loop.EmitKey(kvKeyReg, cc.ctx.Emitter)
		//loop.EmitValue(kvValReg, cc.ctx.Emitter)
	}

	return kvKeyReg, kvValReg, groupSelectors
}

func (cc *LoopCollectCompiler) compileAggregation(c fql.ICollectAggregatorContext, isGrouped bool) {
	if isGrouped {
		cc.compileGroupedAggregation(c)
	} else {
		cc.compileGlobalAggregation(c)
	}
}

func (cc *LoopCollectCompiler) compileGroupedAggregation(c fql.ICollectAggregatorContext) {
	parentLoop := cc.ctx.Loops.Current()
	// We need to allocate a temporary accumulators to store aggregation results
	selectors := c.AllCollectAggregateSelector()
	accums := cc.initAggrAccumulators(selectors)
	loop := cc.ctx.Loops.Create(core.TemporalLoop, core.ForLoop, false)
	loop.Src = cc.ctx.Registers.Allocate(core.Temp)

	// Now we iterate over the grouped items
	parentLoop.EmitValue(loop.Src, cc.ctx.Emitter)
	loop.EmitInitialization(cc.ctx.Registers, cc.ctx.Emitter)
	loop.EmitNext(cc.ctx.Emitter)
	loop.ValueName = parentLoop.ValueName

	// Nested scope for aggregators
	cc.ctx.Symbols.EnterScope()
	loop.DeclareValueVar(loop.ValueName, cc.ctx.Symbols)
	loop.BindValueVar(cc.ctx.Emitter)

	// Now we add value selectors to the accumulators
	cc.collectAggregationFuncArgs(selectors, func(i int, resultReg vm.Operand) {
		cc.ctx.Emitter.EmitAB(vm.OpPush, accums[i], resultReg)
	})

	// Now we can iterate over the grouped items
	loop.EmitFinalization(cc.ctx.Emitter)
	// Now close the aggregators scope
	cc.ctx.Symbols.ExitScope()

	// Now we can iterate over the selectors and execute the aggregation functions by passing the accumulators
	// And define variables for each accumulator result
	cc.compileAggregationFuncCall(selectors, func(i int, _ string) core.RegisterSequence {
		return core.RegisterSequence{accums[i]}
	}, func(i int) {})

	// Free the registers for accumulators
	for _, reg := range accums {
		cc.ctx.Registers.Free(reg)
	}

	// Free the register for the iterator value
	// cc.ctx.Registers.Free(aggrIterVal)
}

func (cc *LoopCollectCompiler) compileGlobalAggregation(c fql.ICollectAggregatorContext) {
	parentLoop := cc.ctx.Loops.Current()
	loop := parentLoop
	// we create a custom collector for aggregators
	cc.ctx.Emitter.PatchSwapAx(loop.ResultPos, vm.OpDataSetCollector, loop.Result, int(core.CollectorTypeKeyGroup))

	// Nested scope for aggregators
	cc.ctx.Symbols.EnterScope()
	// Now we add value selectors to the accumulators
	selectors := c.AllCollectAggregateSelector()
	cc.collectAggregationFuncArgs(selectors, func(i int, resultReg vm.Operand) {
		aggrKeyName := selectors[i].Identifier().GetText()
		aggrKeyReg := loadConstant(cc.ctx, runtime.String(aggrKeyName))
		cc.ctx.Emitter.EmitABC(vm.OpPushKV, loop.Result, aggrKeyReg, resultReg)
		cc.ctx.Registers.Free(aggrKeyReg)
	})

	loop.EmitFinalization(cc.ctx.Emitter)
	cc.ctx.Symbols.ExitScope()

	parentLoop.ValueName = ""
	parentLoop.KeyName = ""

	// Now we can iterate over the grouped items
	zero := loadConstant(cc.ctx, runtime.Int(0))
	one := loadConstant(cc.ctx, runtime.Int(1))
	// We move the aggregator to a temporary register to access it later from the new loop
	aggregator := cc.ctx.Registers.Allocate(core.Temp)
	cc.ctx.Emitter.EmitAB(vm.OpMove, aggregator, loop.Result)

	// Create new loop with 1 iteration only
	cc.ctx.Symbols.EnterScope()
	cc.ctx.Emitter.EmitABC(vm.OpRange, loop.Src, zero, one)
	cc.ctx.Emitter.EmitAb(vm.OpDataSet, loop.Result, loop.Distinct)
	loop.EmitInitialization(cc.ctx.Registers, cc.ctx.Emitter)

	// We just need to take the grouped values and call aggregation functions using them as args
	var key vm.Operand
	var value vm.Operand
	cc.compileAggregationFuncCall(selectors, func(i int, selectorVarName string) core.RegisterSequence {
		// We execute the function call with the accumulator as an argument
		key = loadConstant(cc.ctx, runtime.String(selectorVarName))
		value = cc.ctx.Registers.Allocate(core.Temp)
		cc.ctx.Emitter.EmitABC(vm.OpLoadKey, value, aggregator, key)

		return core.RegisterSequence{value}
	}, func(_ int) {
		cc.ctx.Registers.Free(value)
		cc.ctx.Registers.Free(key)
	})

	cc.ctx.Registers.Free(aggregator)

	// Free the register for the iterator value
	// cc.ctx.Registers.Free(aggrIterVal)
}

func (cc *LoopCollectCompiler) collectAggregationFuncArgs(selectors []fql.ICollectAggregateSelectorContext, collector func(int, vm.Operand)) {
	for i := 0; i < len(selectors); i++ {
		selector := selectors[i]
		fcx := selector.FunctionCallExpression()
		args := cc.ctx.ExprCompiler.CompileArgumentList(fcx.FunctionCall().ArgumentList())

		if len(args) == 0 {
			// TODO: Better error handling
			panic("No arguments provided for the function call in the aggregate selector")
		}

		if len(args) > 1 {
			// TODO: Better error handling
			panic("Too many arguments")
		}

		resultReg := args[0]
		collector(i, resultReg)
		cc.ctx.Registers.Free(resultReg)
	}
}

func (cc *LoopCollectCompiler) compileAggregationFuncCall(selectors []fql.ICollectAggregateSelectorContext, provider func(int, string) core.RegisterSequence, cleanup func(int)) {
	for i, selector := range selectors {
		fcx := selector.FunctionCallExpression()
		// We won't make any checks here, as we already did it before
		selectorVarName := selector.Identifier().GetText()

		result := cc.ctx.ExprCompiler.CompileFunctionCallWith(fcx.FunctionCall(), fcx.ErrorOperator() != nil, provider(i, selectorVarName))

		// We define the variable for the selector result in the upper scope
		// Since this temporary scope is only for aggregators and will be closed after the aggregation
		varReg := cc.ctx.Symbols.DeclareLocal(selectorVarName)
		cc.ctx.Emitter.EmitAB(vm.OpMove, varReg, result)
		cc.ctx.Registers.Free(result)

		cleanup(i)
	}
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

func (cc *LoopCollectCompiler) compileGroupSelectorVariables(selectors []fql.ICollectSelectorContext, kvKeyReg, kvValReg vm.Operand, isAggregation bool) {
	if len(selectors) > 1 {
		variables := make([]vm.Operand, len(selectors))

		for i, selector := range selectors {
			name := selector.Identifier().GetText()

			if variables[i] == vm.NoopOperand {
				variables[i] = cc.ctx.Symbols.DeclareLocal(name)
			}

			reg := kvValReg

			if isAggregation {
				reg = kvKeyReg
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
		reg := cc.selectGroupKey(isAggregation, kvKeyReg, kvValReg)

		// If we have a single selector, we can just move the value
		cc.ctx.Emitter.EmitAB(vm.OpMove, varReg, reg)
	}
}

func (cc *LoopCollectCompiler) compileDefaultGroupProjection(loop *core.Loop, kvValReg vm.Operand, identifier antlr.TerminalNode, keeper fql.ICollectGroupVariableKeeperContext) string {
	if keeper == nil {
		seq := cc.ctx.Registers.AllocateSequence(2) // Key and Value for Map

		// TODO: Review this. It's quite a questionable ArrangoDB feature of wrapping group items by a nested object
		// We will keep it for now for backward compatibility.
		loadConstantTo(cc.ctx, runtime.String(loop.ValueName), seq[0]) // Map key
		cc.ctx.Emitter.EmitAB(vm.OpMove, seq[1], kvValReg)             // Map value
		cc.ctx.Emitter.EmitAs(vm.OpMap, kvValReg, seq)

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

		cc.ctx.Emitter.EmitAs(vm.OpMap, kvValReg, seq)
		cc.ctx.Registers.FreeSequence(seq)
	}

	return identifier.GetText()
}

func (cc *LoopCollectCompiler) compileCustomGroupProjection(_ *core.Loop, kvValReg vm.Operand, selector fql.ICollectSelectorContext) string {
	selectorReg := cc.ctx.ExprCompiler.Compile(selector.Expression())
	cc.ctx.Emitter.EmitMove(kvValReg, selectorReg)
	cc.ctx.Registers.Free(selectorReg)

	return selector.Identifier().GetText()
}

func (cc *LoopCollectCompiler) selectGroupKey(isAggregation bool, kvKeyReg, kvValReg vm.Operand) vm.Operand {
	if isAggregation {
		return kvKeyReg
	}

	return kvValReg
}

func (cc *LoopCollectCompiler) initAggrAccumulators(selectors []fql.ICollectAggregateSelectorContext) []vm.Operand {
	accums := make([]vm.Operand, len(selectors))

	// First of all, we allocate registers for accumulators
	accums = make([]vm.Operand, len(selectors))

	// We need to allocate a register for each accumulator
	for i := 0; i < len(selectors); i++ {
		reg := cc.ctx.Registers.Allocate(core.Temp)
		accums[i] = reg

		// TODO: Select persistent List type, we do not know how many items we will have
		cc.ctx.Emitter.EmitA(vm.OpList, reg)
	}

	return accums
}

func (cc *LoopCollectCompiler) emitPushToAggrAccumulators(accums []vm.Operand, selectors []fql.ICollectAggregateSelectorContext, loop *core.Loop) {
	for i, selector := range selectors {
		fcx := selector.FunctionCallExpression()
		args := cc.ctx.ExprCompiler.CompileArgumentList(fcx.FunctionCall().ArgumentList())

		if len(args) != 1 {
			panic("aggregate function must have exactly one argument")
		}

		cc.ctx.Emitter.EmitAB(vm.OpPush, accums[i], args[0])
		cc.ctx.Registers.Free(args[0])
	}
}
