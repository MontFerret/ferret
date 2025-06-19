package internal

import (
	"github.com/MontFerret/ferret/pkg/compiler/internal/core"
	"github.com/MontFerret/ferret/pkg/parser/fql"
	"github.com/MontFerret/ferret/pkg/runtime"
	"github.com/MontFerret/ferret/pkg/vm"
)

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
	loop := cc.ctx.Loops.CreateFor(core.TemporalLoop, cc.ctx.Registers.Allocate(core.Temp), false)

	// Now we iterate over the grouped items
	parentLoop.EmitValue(loop.Src, cc.ctx.Emitter)

	// Nested scope for aggregators
	cc.ctx.Symbols.EnterScope()
	loop.DeclareValueVar(parentLoop.ValueName, cc.ctx.Symbols)
	loop.EmitInitialization(cc.ctx.Registers, cc.ctx.Emitter)

	// Add value selectors to the accumulators
	cc.collectAggregationFuncArgs(selectors, func(i int, resultReg vm.Operand) {
		cc.ctx.Emitter.EmitAB(vm.OpPush, accums[i], resultReg)
	})

	loop.EmitFinalization(cc.ctx.Emitter)
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
	// we create a custom collector for aggregators
	cc.ctx.Emitter.PatchSwapAx(parentLoop.DstPos, vm.OpDataSetCollector, parentLoop.Dst, int(core.CollectorTypeKeyGroup))

	// Nested scope for aggregators
	cc.ctx.Symbols.EnterScope()
	// Now we add value selectors to the accumulators
	selectors := c.AllCollectAggregateSelector()
	cc.collectAggregationFuncArgs(selectors, func(i int, resultReg vm.Operand) {
		aggrKeyName := selectors[i].Identifier().GetText()
		aggrKeyReg := loadConstant(cc.ctx, runtime.String(aggrKeyName))
		cc.ctx.Emitter.EmitABC(vm.OpPushKV, parentLoop.Dst, aggrKeyReg, resultReg)
		cc.ctx.Registers.Free(aggrKeyReg)
	})

	parentLoop.EmitFinalization(cc.ctx.Emitter)
	cc.ctx.Loops.Pop()
	cc.ctx.Symbols.ExitScope()

	// Now we can iterate over the grouped items
	zero := cc.ctx.Registers.Allocate(core.Temp)
	cc.ctx.Emitter.EmitA(vm.OpLoadZero, zero)
	// We move the aggregator to a temporary register to access it later from the new loop
	aggregator := cc.ctx.Registers.Allocate(core.Temp)
	cc.ctx.Emitter.EmitAB(vm.OpMove, aggregator, parentLoop.Dst)

	// CreateFor new loop with 1 iteration only
	cc.ctx.Symbols.EnterScope()
	cc.ctx.Emitter.EmitABC(vm.OpRange, parentLoop.Src, zero, zero)
	loop := cc.ctx.Loops.CreateFor(core.TemporalLoop, parentLoop.Src, parentLoop.Distinct)
	loop.Dst = parentLoop.Dst
	loop.Allocate = true
	cc.ctx.Loops.Push(loop)
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
