package internal

import (
	"strconv"

	"github.com/MontFerret/ferret/pkg/compiler/internal/core"
	"github.com/MontFerret/ferret/pkg/parser/fql"
	"github.com/MontFerret/ferret/pkg/runtime"
	"github.com/MontFerret/ferret/pkg/vm"
)

func (c *LoopCollectCompiler) compileAggregation(ctx fql.ICollectAggregatorContext, isGrouped bool) {
	if isGrouped {
		c.compileGroupedAggregation(ctx)
	} else {
		c.compileGlobalAggregation(ctx)
	}
}

func (c *LoopCollectCompiler) compileGroupedAggregation(ctx fql.ICollectAggregatorContext) {
	parentLoop := c.ctx.Loops.Current()
	// We need to allocate a temporary accumulator to store aggregation results
	selectors := ctx.AllCollectAggregateSelector()
	accumulator := c.ctx.Registers.Allocate(core.Temp)
	c.ctx.Emitter.EmitAx(vm.OpDataSetCollector, accumulator, int(core.CollectorTypeKeyGroup))

	loop := c.ctx.Loops.NewForInLoop(core.TemporalLoop, false)
	loop.Src = c.ctx.Registers.Allocate(core.Temp)

	// Now we iterate over the grouped items
	parentLoop.EmitValue(loop.Src, c.ctx.Emitter)

	// Nested scope for aggregators
	c.ctx.Symbols.EnterScope()
	loop.DeclareValueVar(parentLoop.ValueName, c.ctx.Symbols)
	loop.EmitInitialization(c.ctx.Registers, c.ctx.Emitter, c.ctx.Loops.Depth())

	// Add value selectors to the accumulators
	argsPkg := c.compileAggregationFuncArgs(selectors, accumulator)

	loop.EmitFinalization(c.ctx.Emitter)
	c.ctx.Symbols.ExitScope()

	// Now we can iterate over the selectors and execute the aggregation functions by passing the accumulators
	// And define variables for each accumulator result
	c.compileAggregationFuncCall(selectors, accumulator, argsPkg)
	c.ctx.Registers.Free(accumulator)
}

func (c *LoopCollectCompiler) compileGlobalAggregation(ctx fql.ICollectAggregatorContext) {
	parentLoop := c.ctx.Loops.Current()
	// we create a custom collector for aggregators
	dst := parentLoop.PatchDestinationAx(c.ctx.Registers, c.ctx.Emitter, vm.OpDataSetCollector, int(core.CollectorTypeKeyGroup))
	// Nested scope for aggregators
	c.ctx.Symbols.EnterScope()
	// Now we add value selectors to the collector
	selectors := ctx.AllCollectAggregateSelector()
	argsPkg := c.compileAggregationFuncArgs(selectors, dst)
	parentLoop.EmitFinalization(c.ctx.Emitter)
	c.ctx.Loops.Pop()
	c.ctx.Symbols.ExitScope()

	// Now we can iterate over the grouped items
	zero := c.ctx.Registers.Allocate(core.Temp)
	c.ctx.Emitter.EmitA(vm.OpLoadZero, zero)
	// We move the aggregator to a temporary register to access it later from the new loop
	aggregator := c.ctx.Registers.Allocate(core.Temp)
	c.ctx.Emitter.EmitAB(vm.OpMove, aggregator, dst)

	if parentLoop.Dst != dst && !parentLoop.Allocate {
		c.ctx.Registers.Free(dst)
	}

	// NewForLoop new loop with 1 iteration only
	c.ctx.Symbols.EnterScope()
	c.ctx.Emitter.EmitABC(vm.OpLoadRange, parentLoop.Src, zero, zero)
	loop := c.ctx.Loops.NewForInLoop(core.TemporalLoop, parentLoop.Distinct)
	loop.Src = parentLoop.Src
	loop.Dst = parentLoop.Dst
	loop.Allocate = parentLoop.Allocate
	c.ctx.Loops.Push(loop)
	loop.EmitInitialization(c.ctx.Registers, c.ctx.Emitter, c.ctx.Loops.Depth())

	// We just need to take the grouped values and call aggregation functions using them as args
	c.compileAggregationFuncCall(selectors, aggregator, argsPkg)
	c.ctx.Registers.Free(aggregator)
}

func (c *LoopCollectCompiler) compileAggregationFuncArgs(selectors []fql.ICollectAggregateSelectorContext, collector vm.Operand) []int {
	argsPkg := make([]int, len(selectors))

	for i := 0; i < len(selectors); i++ {
		selector := selectors[i]
		fcx := selector.FunctionCallExpression()
		args := c.ctx.ExprCompiler.CompileArgumentList(fcx.FunctionCall().ArgumentList())

		if len(args) == 0 {
			// TODO: Better error handling
			panic("No arguments provided for the function call in the aggregate selector")
		}

		aggrKeyReg := loadConstant(c.ctx, runtime.Int(i))
		// we keep information about the args - whether we need to unpack them or not
		argsPkg[i] = len(args)

		if len(args) > 1 {
			for y, arg := range args {
				argKeyReg := c.loadAggregationArgKey(i, y)
				c.ctx.Emitter.EmitABC(vm.OpPushKV, collector, argKeyReg, arg)
				c.ctx.Registers.Free(argKeyReg)
			}
		} else {
			c.ctx.Emitter.EmitABC(vm.OpPushKV, collector, aggrKeyReg, args[0])
		}

		c.ctx.Registers.Free(aggrKeyReg)
		c.ctx.Registers.FreeSequence(args)
	}

	return argsPkg
}

func (c *LoopCollectCompiler) compileAggregationFuncCall(selectors []fql.ICollectAggregateSelectorContext, accumulator vm.Operand, argsPkg []int) {
	// Gets the number of records in the accumulator
	cond := c.ctx.Registers.Allocate(core.Temp)
	c.ctx.Emitter.EmitAB(vm.OpLength, cond, accumulator)
	zero := loadConstant(c.ctx, runtime.ZeroInt)
	// Check if the number equals to zero
	c.ctx.Emitter.EmitEq(cond, cond, zero)
	c.ctx.Registers.Free(zero)
	elseLabel := c.ctx.Emitter.NewLabel()
	endLabel := c.ctx.Emitter.NewLabel()

	// We skip the key retrieval and function call of there are no records in the accumulator
	c.ctx.Emitter.EmitJumpIfTrue(cond, elseLabel)

	selectorVarRegs := make([]vm.Operand, len(selectors))

	for i, selector := range selectors {
		argsNum := argsPkg[i]

		var args core.RegisterSequence

		// We need to unpack arguments
		if argsNum > 1 {
			args = c.ctx.Registers.AllocateSequence(argsNum)

			for y, reg := range args {
				argKeyReg := c.loadAggregationArgKey(i, y)
				c.ctx.Emitter.EmitABC(vm.OpLoadKey, reg, accumulator, argKeyReg)
				c.ctx.Registers.Free(argKeyReg)
			}
		} else {
			key := loadConstant(c.ctx, runtime.Int(i))
			value := c.ctx.Registers.Allocate(core.Temp)
			c.ctx.Emitter.EmitABC(vm.OpLoadKey, value, accumulator, key)
			args = core.RegisterSequence{value}
			c.ctx.Registers.Free(key)
		}

		fcx := selector.FunctionCallExpression()
		result := c.ctx.ExprCompiler.CompileFunctionCallWith(fcx.FunctionCall(), fcx.ErrorOperator() != nil, args)

		// We define the variable for the selector result in the upper scope
		// Since this temporary scope is only for aggregators and will be closed after the aggregation
		selectorVarName := selector.Identifier().GetText()
		varReg := c.ctx.Symbols.DeclareLocal(selectorVarName)
		selectorVarRegs[i] = varReg
		c.ctx.Emitter.EmitAB(vm.OpMove, varReg, result)
		c.ctx.Registers.Free(result)
	}

	c.ctx.Emitter.EmitJump(endLabel)
	c.ctx.Emitter.MarkLabel(elseLabel)

	for _, varReg := range selectorVarRegs {
		c.ctx.Emitter.EmitA(vm.OpLoadNone, varReg)
	}

	c.ctx.Emitter.MarkLabel(endLabel)
	c.ctx.Registers.Free(cond)
}

func (c *LoopCollectCompiler) loadAggregationArgKey(selector int, arg int) vm.Operand {
	argKey := strconv.Itoa(selector) + ":" + strconv.Itoa(arg)
	return loadConstant(c.ctx, runtime.String(argKey))
}
