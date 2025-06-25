package internal

import (
	"strconv"

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
	// We need to allocate a temporary accumulator to store aggregation results
	selectors := c.AllCollectAggregateSelector()
	accumulator := cc.ctx.Registers.Allocate(core.Temp)
	cc.ctx.Emitter.EmitAx(vm.OpDataSetCollector, accumulator, int(core.CollectorTypeKeyGroup))

	loop := cc.ctx.Loops.CreateFor(core.TemporalLoop, cc.ctx.Registers.Allocate(core.Temp), false)

	// Now we iterate over the grouped items
	parentLoop.EmitValue(loop.Src, cc.ctx.Emitter)

	// Nested scope for aggregators
	cc.ctx.Symbols.EnterScope()
	loop.DeclareValueVar(parentLoop.ValueName, cc.ctx.Symbols)
	loop.EmitInitialization(cc.ctx.Registers, cc.ctx.Emitter)

	// Add value selectors to the accumulators
	argsPkg := cc.compileAggregationFuncArgs(selectors, accumulator)

	loop.EmitFinalization(cc.ctx.Emitter)
	cc.ctx.Symbols.ExitScope()

	// Now we can iterate over the selectors and execute the aggregation functions by passing the accumulators
	// And define variables for each accumulator result
	cc.compileAggregationFuncCall(selectors, accumulator, argsPkg)
	cc.ctx.Registers.Free(accumulator)
}

func (cc *LoopCollectCompiler) compileGlobalAggregation(c fql.ICollectAggregatorContext) {
	parentLoop := cc.ctx.Loops.Current()
	// we create a custom collector for aggregators
	cc.ctx.Emitter.PatchSwapAx(parentLoop.Pos, vm.OpDataSetCollector, parentLoop.Dst, int(core.CollectorTypeKeyGroup))

	// Nested scope for aggregators
	cc.ctx.Symbols.EnterScope()
	// Now we add value selectors to the collector
	selectors := c.AllCollectAggregateSelector()
	argsPkg := cc.compileAggregationFuncArgs(selectors, parentLoop.Dst)

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
	cc.compileAggregationFuncCall(selectors, aggregator, argsPkg)
	cc.ctx.Registers.Free(aggregator)
}

func (cc *LoopCollectCompiler) compileAggregationFuncArgs(selectors []fql.ICollectAggregateSelectorContext, collector vm.Operand) []int {
	argsPkg := make([]int, len(selectors))

	for i := 0; i < len(selectors); i++ {
		selector := selectors[i]
		fcx := selector.FunctionCallExpression()
		args := cc.ctx.ExprCompiler.CompileArgumentList(fcx.FunctionCall().ArgumentList())

		if len(args) == 0 {
			// TODO: Better error handling
			panic("No arguments provided for the function call in the aggregate selector")
		}

		aggrKeyReg := loadConstant(cc.ctx, runtime.Int(i))
		// we keep information about the args - whether we need to unpack them or not
		argsPkg[i] = len(args)

		if len(args) > 1 {
			for y, arg := range args {
				argKeyReg := cc.loadAggregationArgKey(i, y)
				cc.ctx.Emitter.EmitABC(vm.OpPushKV, collector, argKeyReg, arg)
				cc.ctx.Registers.Free(argKeyReg)
			}
		} else {
			cc.ctx.Emitter.EmitABC(vm.OpPushKV, collector, aggrKeyReg, args[0])
		}

		cc.ctx.Registers.Free(aggrKeyReg)
		cc.ctx.Registers.FreeSequence(args)
	}

	return argsPkg
}

func (cc *LoopCollectCompiler) compileAggregationFuncCall(selectors []fql.ICollectAggregateSelectorContext, accumulator vm.Operand, argsPkg []int) {
	for i, selector := range selectors {
		argsNum := argsPkg[i]

		var args core.RegisterSequence

		// We need to unpack arguments
		if argsNum > 1 {
			args = cc.ctx.Registers.AllocateSequence(argsNum)

			for y, reg := range args {
				argKeyReg := cc.loadAggregationArgKey(i, y)
				cc.ctx.Emitter.EmitABC(vm.OpLoadKey, reg, accumulator, argKeyReg)

				cc.ctx.Registers.Free(argKeyReg)
			}
		} else {
			key := loadConstant(cc.ctx, runtime.Int(i))
			value := cc.ctx.Registers.Allocate(core.Temp)
			cc.ctx.Emitter.EmitABC(vm.OpLoadKey, value, accumulator, key)
			args = core.RegisterSequence{value}
			cc.ctx.Registers.Free(key)
		}

		fcx := selector.FunctionCallExpression()
		result := cc.ctx.ExprCompiler.CompileFunctionCallWith(fcx.FunctionCall(), fcx.ErrorOperator() != nil, args)

		// We define the variable for the selector result in the upper scope
		// Since this temporary scope is only for aggregators and will be closed after the aggregation
		selectorVarName := selector.Identifier().GetText()
		varReg := cc.ctx.Symbols.DeclareLocal(selectorVarName)
		cc.ctx.Emitter.EmitAB(vm.OpMove, varReg, result)
		cc.ctx.Registers.Free(result)
	}
}

func (cc *LoopCollectCompiler) loadAggregationArgKey(selector int, arg int) vm.Operand {
	argKey := strconv.Itoa(selector) + ":" + strconv.Itoa(arg)
	return loadConstant(cc.ctx, runtime.String(argKey))
}
