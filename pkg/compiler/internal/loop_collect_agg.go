package internal

import (
	"strconv"

	"github.com/MontFerret/ferret/pkg/compiler/internal/core"
	"github.com/MontFerret/ferret/pkg/parser/fql"
	"github.com/MontFerret/ferret/pkg/runtime"
	"github.com/MontFerret/ferret/pkg/vm"
)

func (c *LoopCollectCompiler) initializeAggregation(ctx fql.ICollectAggregatorContext, dst vm.Operand, kv *core.KV, withGrouping bool) []*core.AggregateSelector {
	selectors := ctx.AllCollectAggregateSelector()
	var compiledSelectors []*core.AggregateSelector

	// if we have grouping, we need to pack the selectors into the collector value
	if withGrouping {
		compiledSelectors = c.compileGroupedAggregationSelectors(selectors)

		// Pack the selectors into the collector value
		c.packGroupedValues(kv, compiledSelectors)
	} else {
		// We just push the selectors into the global collector
		compiledSelectors = c.compileGlobalAggregationSelectors(selectors, dst)
	}

	return compiledSelectors
}

func (c *LoopCollectCompiler) packGroupedValues(kv *core.KV, selectors []*core.AggregateSelector) {
	// We need to add the loop value to the array
	seq := c.ctx.Registers.AllocateSequence(len(selectors) + 1)
	c.ctx.Emitter.EmitMove(seq[0], kv.Value)

	for i, selector := range selectors {
		c.ctx.Emitter.EmitMove(seq[i+1], selector.Register())
		c.ctx.Registers.Free(selector.Register())
	}

	// Now we need to wrap the selectors into a single array with the loop value
	c.ctx.Emitter.EmitArray(kv.Value, seq)
}

func (c *LoopCollectCompiler) compileGroupedAggregationSelectors(selectors []fql.ICollectAggregateSelectorContext) []*core.AggregateSelector {
	wrappedSelectors := make([]*core.AggregateSelector, 0, len(selectors))

	for i := 0; i < len(selectors); i++ {
		selector := selectors[i]
		name := runtime.String(selector.Identifier().GetText())
		fcx := selector.FunctionCallExpression()
		args := c.ctx.ExprCompiler.CompileArgumentList(fcx.FunctionCall().ArgumentList())

		if len(args) == 0 {
			// TODO: Better error handling
			panic("No arguments provided for the function call in the aggregate selector")
		}

		var selectorArg vm.Operand

		if len(args) > 1 {
			// We pack multiple arguments into an array
			selectorArg = c.ctx.Registers.Allocate(core.Temp)
			c.ctx.Emitter.EmitArray(selectorArg, args)
			c.ctx.Registers.FreeSequence(args)
		} else {
			// We can use a single argument directly
			selectorArg = args[0]
		}

		fce := selector.FunctionCallExpression()
		funcName := getFunctionName(fce.FunctionCall())
		isProtected := fce.ErrorOperator() != nil

		// Collect information about the selector to unpack it later
		wrappedSelectors = append(wrappedSelectors, core.NewAggregateSelector(name, len(args), funcName, isProtected, selectorArg))
	}

	return wrappedSelectors
}

func (c *LoopCollectCompiler) compileGlobalAggregationSelectors(selectors []fql.ICollectAggregateSelectorContext, dst vm.Operand) []*core.AggregateSelector {
	wrappedSelectors := make([]*core.AggregateSelector, 0, len(selectors))

	for i := 0; i < len(selectors); i++ {
		selector := selectors[i]
		name := runtime.String(selector.Identifier().GetText())
		fcx := selector.FunctionCallExpression()
		args := c.ctx.ExprCompiler.CompileArgumentList(fcx.FunctionCall().ArgumentList())

		if len(args) == 0 {
			// TODO: Better error handling
			panic("No arguments provided for the function call in the aggregate selector")
		}

		if len(args) > 1 {
			for y := 0; y < len(args); i++ {
				key := c.loadAggregationArgKey(name, y)
				c.ctx.Emitter.EmitPushKV(dst, key, args[y])
				c.ctx.Registers.Free(key)
			}
		} else {
			// We can use a single argument directly
			key := loadConstant(c.ctx, name)
			c.ctx.Emitter.EmitPushKV(dst, key, args[0])
			c.ctx.Registers.Free(key)
		}

		fce := selector.FunctionCallExpression()
		funcName := getFunctionName(fce.FunctionCall())
		isProtected := fce.ErrorOperator() != nil

		// Collect information about the selector to unpack it later
		wrappedSelectors = append(wrappedSelectors, core.NewAggregateSelector(name, len(args), funcName, isProtected, vm.NoopOperand))

		c.ctx.Registers.FreeSequence(args)
	}

	return wrappedSelectors
}

func (c *LoopCollectCompiler) unpackGroupedValues(spec *core.CollectorSpec) {
	if !spec.HasGrouping() {
		return
	}

	loop := c.ctx.Loops.Current()
	valReg := c.ctx.Registers.Allocate(core.Temp)

	loadIndex(c.ctx, valReg, loop.Value, 0)

	for i, selector := range spec.AggregationSelectors() {
		loadIndex(c.ctx, selector.Register(), loop.Value, i+1)
	}

	c.ctx.Registers.Free(valReg)
}

func (c *LoopCollectCompiler) compileAggregation(spec *core.CollectorSpec) {
	if spec.HasGrouping() {
		c.compileGroupedAggregation(spec)
	} else {
		c.compileGlobalAggregation(spec)
	}
}

func (c *LoopCollectCompiler) compileGroupedAggregation(spec *core.CollectorSpec) {
	//parentLoop := c.ctx.Loops.Current()
	//// We need to allocate a temporary accumulator to store aggregation results
	//selectors := ctx.AllCollectAggregateSelector()
	//accumulator := c.ctx.Registers.Allocate(core.Temp)
	//c.ctx.Emitter.EmitAx(vm.OpDataSetCollector, accumulator, int(core.CollectorTypeKeyGroup))
	//
	//loop := c.ctx.Loops.NewForInLoop(core.TemporalLoop, false)
	//loop.Src = c.ctx.Registers.Allocate(core.Temp)
	//
	//// Now we iterate over the grouped items
	//parentLoop.EmitValue(loop.Src, c.ctx.Emitter)
	//
	//// Nested scope for aggregators
	//c.ctx.Symbols.EnterScope()
	//loop.DeclareValueVar(parentLoop.ValueName, c.ctx.Symbols)
	//loop.EmitInitialization(c.ctx.Registers, c.ctx.Emitter, c.ctx.Loops.Depth())
	//
	//// Add value selectors to the accumulators
	//argsPkg := c.compileGroupedAggregationSelectors(selectors, accumulator)
	//
	//loop.EmitFinalization(c.ctx.Emitter)
	//c.ctx.Symbols.ExitScope()
	//
	//// Now we can iterate over the selectors and execute the aggregation functions by passing the accumulators
	//// And define variables for each accumulator result
	//c.compileAggregationFuncCalls(selectors, accumulator, argsPkg)
	//c.ctx.Registers.Free(accumulator)
}

func (c *LoopCollectCompiler) compileGlobalAggregation(spec *core.CollectorSpec) {
	// At this point, it's finalized.
	prevLoop := c.ctx.Loops.Pop()
	c.ctx.Registers.Free(prevLoop.Key)
	c.ctx.Registers.Free(prevLoop.Value)
	c.ctx.Registers.Free(prevLoop.Src)

	// NewForLoop new loop with 1 iteration only
	c.ctx.Symbols.EnterScope()
	loop := c.ctx.Loops.NewLoop(core.ForInLoop, core.NormalLoop, prevLoop.Distinct)
	c.ctx.Loops.Push(loop)

	loop.Src = c.ctx.Registers.Allocate(core.Temp)
	zero := c.ctx.Registers.Allocate(core.Temp)
	c.ctx.Emitter.EmitA(vm.OpLoadZero, zero)
	c.ctx.Emitter.EmitABC(vm.OpLoadRange, loop.Src, zero, zero)
	loop.Allocate = prevLoop.Allocate

	if !loop.Allocate {
		parent := c.ctx.Loops.FindParent(c.ctx.Loops.Depth())
		loop.Dst = parent.Dst
	}

	loop.EmitInitialization(c.ctx.Registers, c.ctx.Emitter, c.ctx.Loops.Depth())

	// We just need to take the grouped values and call aggregation functions using them as args
	c.compileAggregationFuncCalls(spec, prevLoop.Dst)

	c.ctx.Registers.Free(prevLoop.Dst)
}

func (c *LoopCollectCompiler) compileAggregationFuncCalls(spec *core.CollectorSpec, aggregator vm.Operand) {
	// Gets the number of records in the accumulator
	cond := c.ctx.Registers.Allocate(core.Temp)
	c.ctx.Emitter.EmitAB(vm.OpLength, cond, aggregator)
	zero := loadConstant(c.ctx, runtime.ZeroInt)
	// Check if the number equals to zero
	c.ctx.Emitter.EmitEq(cond, cond, zero)
	c.ctx.Registers.Free(zero)
	elseLabel := c.ctx.Emitter.NewLabel()
	endLabel := c.ctx.Emitter.NewLabel()

	// We skip the key retrieval and function call of there are no records in the accumulator
	c.ctx.Emitter.EmitJumpIfTrue(cond, elseLabel)

	selectors := spec.AggregationSelectors()
	selectorVarRegs := make([]vm.Operand, len(selectors))

	for i, selector := range selectors {
		var args core.RegisterSequence

		// We need to unpack arguments
		if selector.Args() > 1 {
			args = c.ctx.Registers.AllocateSequence(selector.Args())

			for y, reg := range args {
				argKeyReg := c.loadAggregationArgKey(selector.Name(), y)
				c.ctx.Emitter.EmitABC(vm.OpLoadKey, reg, aggregator, argKeyReg)
				c.ctx.Registers.Free(argKeyReg)
			}
		} else {
			key := loadConstant(c.ctx, selector.Name())
			value := c.ctx.Registers.Allocate(core.Temp)
			c.ctx.Emitter.EmitABC(vm.OpLoadKey, value, aggregator, key)
			args = core.RegisterSequence{value}
			c.ctx.Registers.Free(key)
		}

		result := c.ctx.ExprCompiler.CompileFunctionCallByNameWith(selector.FuncName(), selector.ProtectedCall(), args)

		// We define the variable for the selector result in the upper scope
		// Since this temporary scope is only for aggregators and will be closed after the aggregation
		selectorVarName := selector.Name()
		varReg := c.ctx.Symbols.DeclareLocal(selectorVarName.String(), core.TypeUnknown)
		selectorVarRegs[i] = varReg
		c.ctx.Emitter.EmitAB(vm.OpMove, varReg, result)
		c.ctx.Registers.Free(result)
	}

	var projVar vm.Operand

	// If the projection is used, we allocate a new register for the variable and put the iterator's value into it
	if spec.HasProjection() {
		projVar = c.finalizeProjection(spec, aggregator)
	}

	c.ctx.Emitter.EmitJump(endLabel)
	c.ctx.Emitter.MarkLabel(elseLabel)

	for _, varReg := range selectorVarRegs {
		c.ctx.Emitter.EmitA(vm.OpLoadNone, varReg)
	}

	if projVar != vm.NoopOperand {
		c.ctx.Emitter.EmitA(vm.OpLoadNone, projVar)
	}

	c.ctx.Emitter.MarkLabel(endLabel)
	c.ctx.Registers.Free(cond)
}

func (c *LoopCollectCompiler) compileAggregationFuncCall(selector *core.AggregateSelector) {
	varReg := c.ctx.Symbols.DeclareLocal(selector.Name().String(), core.TypeUnknown)
	loadIndex(c.ctx, varReg, selector.Register(), 1)
}

func (c *LoopCollectCompiler) loadAggregationArgKey(selector runtime.String, arg int) vm.Operand {
	argKey := selector.String() + ":" + strconv.Itoa(arg)
	return loadConstant(c.ctx, runtime.String(argKey))
}
