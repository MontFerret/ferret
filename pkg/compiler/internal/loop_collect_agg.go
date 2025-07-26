package internal

import (
	"strconv"

	"github.com/MontFerret/ferret/pkg/compiler/internal/core"
	"github.com/MontFerret/ferret/pkg/parser/fql"
	"github.com/MontFerret/ferret/pkg/runtime"
	"github.com/MontFerret/ferret/pkg/vm"
)

// initializeAggregation processes the aggregation selectors from the COLLECT clause.
// It handles both grouped and global aggregations differently.
// For grouped aggregations, it compiles the selectors and packs them with the loop value.
// For global aggregations, it pushes the selectors directly to the collector.
// Returns a slice of AggregateSelectors that describe the aggregation operations.
func (c *LoopCollectCompiler) initializeAggregation(ctx fql.ICollectAggregatorContext, dst vm.Operand, kv *core.KV, withGrouping bool) *core.CollectorAggregation {
	loop := c.ctx.Loops.Current()
	selectors := ctx.AllCollectAggregateSelector()

	// If we have grouping, we need to pack the selectors into the collector value
	if withGrouping {
		// TODO: We need to figure out how to free the aggregator register later
		aggregator := c.ctx.Registers.Allocate(core.State)
		// We create a separate collector for aggregation in grouped mode
		c.ctx.Emitter.InsertAx(loop.StartLabel, vm.OpDataSetCollector, aggregator, int(core.CollectorTypeKeyGroup))

		// Compile selectors for grouped aggregation
		aggregateSelectors := c.initializeGroupedAggregationSelectors(selectors, kv, aggregator)

		return core.NewCollectorAggregation(aggregator, aggregateSelectors)
	}

	// For global aggregation, we just push the selectors into the global collector
	aggregateSelectors := c.initializeGlobalAggregationSelectors(selectors, dst)

	return core.NewCollectorAggregation(dst, aggregateSelectors)
}

// initializeGroupedAggregationSelectors processes aggregation selectors for grouped aggregations.
// It compiles each selector's function call expression and arguments, and creates AggregateSelector objects.
// For selectors with multiple arguments, it packs them into an array.
// Returns a slice of AggregateSelectors that describe the aggregation operations.
func (c *LoopCollectCompiler) initializeGroupedAggregationSelectors(selectors []fql.ICollectAggregateSelectorContext, kv *core.KV, dst vm.Operand) []*core.AggregateSelector {
	wrappedSelectors := make([]*core.AggregateSelector, len(selectors))

	for i := 0; i < len(selectors); i++ {
		selector := selectors[i]
		// Get the variable name for this aggregation result
		name := runtime.String(selector.Identifier().GetText())
		// Get the function call expression
		fcx := selector.FunctionCallExpression()
		// Compile the function arguments
		args := c.ctx.ExprCompiler.CompileArgumentList(fcx.FunctionCall().ArgumentList())

		if len(args) == 0 {
			// TODO: Better error handling
			panic("No arguments provided for the function call in the aggregate selector")
		}

		if len(args) > 1 {
			// For multiple arguments, push each one with an indexed key
			for y := 0; y < len(args); y++ {
				// Create a key with format "name:index"
				key := c.loadSelectorKey(kv.Key, name, y)
				// Push the key-value pair to the collector
				c.ctx.Emitter.EmitPushKV(dst, key, args[y])
				c.ctx.Registers.Free(key)
			}
		} else {
			// For a single argument, use the selector name as the key
			key := c.loadSelectorKey(kv.Key, name, -1)
			// Push the key-value pair to the collector
			c.ctx.Emitter.EmitPushKV(dst, key, args[0])
			c.ctx.Registers.Free(key)
		}

		// Get the function name and check if it's a protected call (with TRY)
		fce := selector.FunctionCallExpression()
		funcName := getFunctionName(fce.FunctionCall())
		isProtected := fce.ErrorOperator() != nil

		// Create an AggregateSelector with all the information needed to process it later
		wrappedSelectors[i] = core.NewAggregateSelector(name, len(args), funcName, isProtected)

		// Free the argument registers
		c.ctx.Registers.FreeSequence(args)
	}

	return wrappedSelectors
}

// initializeGlobalAggregationSelectors processes aggregation selectors for global (non-grouped) aggregations.
// It compiles each selector's function call expression and arguments, and pushes them directly to the collector.
// For selectors with multiple arguments, it uses indexed keys to store each argument separately.
// Returns a slice of AggregateSelectors that describe the aggregation operations.
func (c *LoopCollectCompiler) initializeGlobalAggregationSelectors(selectors []fql.ICollectAggregateSelectorContext, dst vm.Operand) []*core.AggregateSelector {
	wrappedSelectors := make([]*core.AggregateSelector, 0, len(selectors))

	for i := 0; i < len(selectors); i++ {
		selector := selectors[i]
		// Get the variable name for this aggregation result
		name := runtime.String(selector.Identifier().GetText())
		// Get the function call expression
		fcx := selector.FunctionCallExpression()
		// Compile the function arguments
		args := c.ctx.ExprCompiler.CompileArgumentList(fcx.FunctionCall().ArgumentList())

		if len(args) == 0 {
			// TODO: Better error handling
			panic("No arguments provided for the function call in the aggregate selector")
		}

		if len(args) > 1 {
			// For multiple arguments, push each one with an indexed key
			for y := 0; y < len(args); y++ {
				// Create a key with format "name:index"
				key := c.loadGlobalSelectorKey(name, y)
				// Push the key-value pair to the collector
				c.ctx.Emitter.EmitPushKV(dst, key, args[y])
				c.ctx.Registers.Free(key)
			}
		} else {
			// For a single argument, use the selector name as the key
			key := loadConstant(c.ctx, name)
			// Push the key-value pair to the collector
			c.ctx.Emitter.EmitPushKV(dst, key, args[0])
			c.ctx.Registers.Free(key)
		}

		// Get the function name and check if it's a protected call (with TRY)
		fce := selector.FunctionCallExpression()
		funcName := getFunctionName(fce.FunctionCall())
		isProtected := fce.ErrorOperator() != nil

		// For global aggregation, we don't need to store the register in the selector
		// as the values are already pushed to the collector
		wrappedSelectors = append(wrappedSelectors, core.NewAggregateSelector(name, len(args), funcName, isProtected))

		// Free the argument registers
		c.ctx.Registers.FreeSequence(args)
	}

	return wrappedSelectors
}

// finalizeAggregation processes the aggregation operations based on the collector specification.
// It delegates to either grouped or global aggregation compilation based on whether grouping is used.
func (c *LoopCollectCompiler) finalizeAggregation(spec *core.Collector) {
	if spec.HasGrouping() {
		// For aggregations with grouping
		c.finalizeGroupedAggregation(spec)
	} else {
		// For global aggregations without grouping
		c.finalizeGlobalAggregation(spec)
	}
}

// finalizeGroupedAggregation handles grouped aggregation operations.
// This function is currently commented out in the original code, likely because
// the functionality is implemented differently or is being refactored.
// The commented code shows the intended approach for handling grouped aggregations.
func (c *LoopCollectCompiler) finalizeGroupedAggregation(spec *core.Collector) {
	for i, selector := range spec.Aggregation().Selectors() {
		c.compileGroupedAggregationFuncCall(selector, spec.Aggregation().State(), i)
	}
}

// finalizeGlobalAggregation handles global (non-grouped) aggregation operations.
// It creates a new loop with a single iteration to process the aggregation results.
// This approach allows the aggregation to be processed in a consistent way with other operations.
func (c *LoopCollectCompiler) finalizeGlobalAggregation(spec *core.Collector) {
	// At this point, the previous loop is finalized, so we can pop it and free its registers
	prevLoop := c.ctx.Loops.Pop()
	c.ctx.Registers.Free(prevLoop.Key)
	c.ctx.Registers.Free(prevLoop.Value)
	c.ctx.Registers.Free(prevLoop.Src)

	// Create a new loop with 1 iteration only to process the aggregation
	loop := c.ctx.Loops.NewLoop(core.ForInLoop, core.NormalLoop, prevLoop.Distinct)
	c.ctx.Loops.Push(loop)

	// Set up the loop source to be a range from 0 to 0 (one iteration)
	loop.Src = c.ctx.Registers.Allocate(core.Temp)
	zero := c.ctx.Registers.Allocate(core.Temp)
	c.ctx.Emitter.EmitA(vm.OpLoadZero, zero)
	c.ctx.Emitter.EmitABC(vm.OpLoadRange, loop.Src, zero, zero)

	// Inherit allocation flag from previous loop
	loop.Allocate = prevLoop.Allocate

	// If not allocating, use the parent loop's destination
	if !loop.Allocate {
		parent := c.ctx.Loops.FindParent(c.ctx.Loops.Depth())
		loop.Dst = parent.Dst
	}

	// Initialize the loop
	loop.EmitInitialization(c.ctx.Registers, c.ctx.Emitter, c.ctx.Loops.Depth())

	c.compileGlobalAggregationFuncCalls(spec)
}

// compileGlobalAggregationFuncCalls processes the aggregation function calls for the selectors.
// It loads the arguments from the aggregator, calls the aggregation functions,
// and assigns the results to local variables.
// It also handles the case where there are no records in the aggregator by loading NONE values.
func (c *LoopCollectCompiler) compileGlobalAggregationFuncCalls(spec *core.Collector) {
	// Gets the number of records in the accumulator
	aggregator := spec.Destination()
	cond := c.ctx.Registers.Allocate(core.Temp)
	c.ctx.Emitter.EmitAB(vm.OpLength, cond, aggregator)
	zero := loadConstant(c.ctx, runtime.ZeroInt)
	// Check if the number equals to zero
	c.ctx.Emitter.EmitEq(cond, cond, zero)
	c.ctx.Registers.Free(zero)
	elseLabel := c.ctx.Emitter.NewLabel()
	endLabel := c.ctx.Emitter.NewLabel()

	// We skip the key retrieval and function call if there are no records in the accumulator
	c.ctx.Emitter.EmitJumpIfTrue(cond, elseLabel)

	selectors := spec.Aggregation().Selectors()
	selectorVarRegs := make([]vm.Operand, len(selectors))

	// Process each aggregation selector
	for i, selector := range selectors {
		var args core.RegisterSequence

		// We need to unpack arguments from the aggregator
		if selector.Args() > 1 {
			// For multiple arguments, allocate a sequence and load each argument by its indexed key
			args = c.ctx.Registers.AllocateSequence(selector.Args())

			for y, reg := range args {
				argKeyReg := c.loadGlobalSelectorKey(selector.Name(), y)
				c.ctx.Emitter.EmitABC(vm.OpLoadKey, reg, aggregator, argKeyReg)
				c.ctx.Registers.Free(argKeyReg)
			}
		} else {
			// For a single argument, load it directly using the selector name as key
			key := loadConstant(c.ctx, selector.Name())
			value := c.ctx.Registers.Allocate(core.Temp)
			c.ctx.Emitter.EmitABC(vm.OpLoadKey, value, aggregator, key)
			args = core.RegisterSequence{value}
			c.ctx.Registers.Free(key)
		}

		// Call the aggregation function with the loaded arguments
		result := c.ctx.ExprCompiler.CompileFunctionCallByNameWith(selector.FuncName(), selector.ProtectedCall(), args)

		// Declare a local variable for the aggregation result
		selectorVarName := selector.Name()
		varReg := c.ctx.Symbols.DeclareLocal(selectorVarName.String(), core.TypeUnknown)
		selectorVarRegs[i] = varReg
		// Move the function result to the variable
		c.ctx.Emitter.EmitAB(vm.OpMove, varReg, result)
		c.ctx.Registers.Free(result)
	}

	var projVar vm.Operand

	// If the projection is used, we allocate a new register for the variable and put the iterator's value into it
	if spec.HasProjection() {
		projVar = c.finalizeProjection(spec, aggregator)
	}

	// Skip the else block (for empty aggregator)
	c.ctx.Emitter.EmitJump(endLabel)
	c.ctx.Emitter.MarkLabel(elseLabel)

	// If there are no records in the aggregator, load NONE values for all variables
	for _, varReg := range selectorVarRegs {
		c.ctx.Emitter.EmitA(vm.OpLoadNone, varReg)
	}

	if projVar != vm.NoopOperand {
		c.ctx.Emitter.EmitA(vm.OpLoadNone, projVar)
	}

	c.ctx.Emitter.MarkLabel(endLabel)
	c.ctx.Registers.Free(cond)
}

func (c *LoopCollectCompiler) compileGroupedAggregationFuncCall(selector *core.AggregateSelector, aggregator vm.Operand, idx int) {
	loop := c.ctx.Loops.Current()
	// Declare a local variable with the selector name
	valReg := c.ctx.Symbols.DeclareLocal(selector.Name().String(), core.TypeUnknown)

	var args core.RegisterSequence

	// We need to unpack arguments from the aggregator
	if selector.Args() > 1 {
		// For multiple arguments, allocate a sequence and load each argument by its indexed key
		args = c.ctx.Registers.AllocateSequence(selector.Args())

		for y, reg := range args {
			key := c.loadSelectorKey(loop.Key, selector.Name(), y)
			c.ctx.Emitter.EmitABC(vm.OpLoadKey, reg, aggregator, key)
			c.ctx.Registers.Free(key)
		}
	} else {
		// For a single argument, load it directly using the selector name as key
		key := c.loadSelectorKey(loop.Key, selector.Name(), -1)
		value := c.ctx.Registers.Allocate(core.Temp)
		c.ctx.Emitter.EmitABC(vm.OpLoadKey, value, aggregator, key)
		args = core.RegisterSequence{value}
		c.ctx.Registers.Free(key)
	}

	resArg := c.ctx.ExprCompiler.CompileFunctionCallByNameWith(selector.FuncName(), selector.ProtectedCall(), args)

	c.ctx.Emitter.EmitMove(valReg, resArg)
}

// loadGlobalSelectorKey creates a key for an aggregation argument by combining the selector name and argument index.
// This is used for global aggregations with multiple arguments to store each argument separately.
// Returns a register containing the key as a string constant.
func (c *LoopCollectCompiler) loadGlobalSelectorKey(selector runtime.String, arg int) vm.Operand {
	// Create a key with format "selectorName:argIndex"
	argKey := selector.String() + ":" + strconv.Itoa(arg)
	// Load the key as a string constant
	return loadConstant(c.ctx, runtime.String(argKey))
}

func (c *LoopCollectCompiler) loadSelectorKey(key vm.Operand, selector runtime.String, arg int) vm.Operand {
	selectorKey := c.ctx.Registers.Allocate(core.Temp)
	selectorName := loadConstant(c.ctx, selector)

	c.ctx.Emitter.EmitABC(vm.OpAdd, selectorKey, key, selectorName)

	if arg >= 0 {
		selectorIndex := loadConstant(c.ctx, runtime.String(strconv.Itoa(arg)))
		c.ctx.Emitter.EmitABC(vm.OpAdd, selectorKey, selectorKey, selectorIndex)
		c.ctx.Registers.Free(selectorIndex)
	}

	c.ctx.Registers.Free(selectorName)

	return selectorKey
}
