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
func (c *LoopCollectCompiler) initializeAggregation(ctx fql.ICollectAggregatorContext, dst vm.Operand, kv *core.KV, withGrouping bool) []*core.AggregateSelector {
	selectors := ctx.AllCollectAggregateSelector()
	var compiledSelectors []*core.AggregateSelector

	// If we have grouping, we need to pack the selectors into the collector value
	if withGrouping {
		// Compile selectors for grouped aggregation
		compiledSelectors = c.compileGroupedAggregationSelectors(selectors)

		// Pack the selectors into the collector value along with the loop value
		c.packGroupedValues(kv, compiledSelectors)
	} else {
		// For global aggregation, we just push the selectors into the global collector
		compiledSelectors = c.compileGlobalAggregationSelectors(selectors, dst)
	}

	return compiledSelectors
}

// packGroupedValues combines the loop value with aggregation selector values into a single array.
// This is used for grouped aggregations to keep all values together for each group.
// The first element of the array is the loop value, followed by the aggregation selector values.
func (c *LoopCollectCompiler) packGroupedValues(kv *core.KV, selectors []*core.AggregateSelector) {
	// Allocate a sequence of registers for the array elements
	// We need one extra register for the loop value (hence +1)
	seq := c.ctx.Registers.AllocateSequence(len(selectors) + 1)

	// Move the loop value to the first position in the sequence
	c.ctx.Emitter.EmitMove(seq[0], kv.Value)

	// Move each selector value to its position in the sequence
	for i, selector := range selectors {
		c.ctx.Emitter.EmitMove(seq[i+1], selector.Register())
		// Free the selector register as we no longer need it
		c.ctx.Registers.Free(selector.Register())
	}

	// Create an array from the sequence and store it in the kv.Value register
	// This replaces the original loop value with an array containing both
	// the loop value and all selector values
	c.ctx.Emitter.EmitArray(kv.Value, seq)
}

// compileGroupedAggregationSelectors processes aggregation selectors for grouped aggregations.
// It compiles each selector's function call expression and arguments, and creates AggregateSelector objects.
// For selectors with multiple arguments, it packs them into an array.
// Returns a slice of AggregateSelectors that describe the aggregation operations.
func (c *LoopCollectCompiler) compileGroupedAggregationSelectors(selectors []fql.ICollectAggregateSelectorContext) []*core.AggregateSelector {
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

		var selectorArg vm.Operand

		if len(args) > 1 {
			// For multiple arguments, pack them into an array
			selectorArg = c.ctx.Registers.Allocate(core.Temp)
			c.ctx.Emitter.EmitArray(selectorArg, args)
			c.ctx.Registers.FreeSequence(args)
		} else {
			// For a single argument, use it directly
			selectorArg = args[0]
		}

		// Get the function name and check if it's a protected call (with TRY)
		fce := selector.FunctionCallExpression()
		funcName := getFunctionName(fce.FunctionCall())
		isProtected := fce.ErrorOperator() != nil

		// Create an AggregateSelector with all the information needed to process it later
		wrappedSelectors = append(wrappedSelectors, core.NewAggregateSelector(name, len(args), funcName, isProtected, selectorArg))
	}

	return wrappedSelectors
}

// compileGlobalAggregationSelectors processes aggregation selectors for global (non-grouped) aggregations.
// It compiles each selector's function call expression and arguments, and pushes them directly to the collector.
// For selectors with multiple arguments, it uses indexed keys to store each argument separately.
// Returns a slice of AggregateSelectors that describe the aggregation operations.
func (c *LoopCollectCompiler) compileGlobalAggregationSelectors(selectors []fql.ICollectAggregateSelectorContext, dst vm.Operand) []*core.AggregateSelector {
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
				key := c.loadAggregationArgKey(name, y)
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
		wrappedSelectors = append(wrappedSelectors, core.NewAggregateSelector(name, len(args), funcName, isProtected, vm.NoopOperand))

		// Free the argument registers
		c.ctx.Registers.FreeSequence(args)
	}

	return wrappedSelectors
}

// unpackGroupedValues extracts values from the packed array created during grouped aggregation.
// It loads the loop value from index 0 and each aggregation selector value from subsequent indices.
// This is only needed for grouped aggregations, so it returns early if there's no grouping.
func (c *LoopCollectCompiler) unpackGroupedValues(spec *core.CollectorSpec) {
	// Skip if there's no grouping
	if !spec.HasGrouping() {
		return
	}

	loop := c.ctx.Loops.Current()
	// Allocate a temporary register for the loop value
	valReg := c.ctx.Registers.Allocate(core.Temp)

	// Load the original loop value from index 0 of the array
	loadIndex(c.ctx, valReg, loop.Value, 0)

	// Load each aggregation selector value from its index in the array
	for i, selector := range spec.AggregationSelectors() {
		loadIndex(c.ctx, selector.Register(), loop.Value, i+1)
	}

	// Free the temporary register
	c.ctx.Registers.Free(valReg)
}

// compileAggregation processes the aggregation operations based on the collector specification.
// It delegates to either grouped or global aggregation compilation based on whether grouping is used.
func (c *LoopCollectCompiler) compileAggregation(spec *core.CollectorSpec) {
	if spec.HasGrouping() {
		// For aggregations with grouping
		c.compileGroupedAggregation(spec)
	} else {
		// For global aggregations without grouping
		c.compileGlobalAggregation(spec)
	}
}

// compileGroupedAggregation handles grouped aggregation operations.
// This function is currently commented out in the original code, likely because
// the functionality is implemented differently or is being refactored.
// The commented code shows the intended approach for handling grouped aggregations.
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

// compileGlobalAggregation handles global (non-grouped) aggregation operations.
// It creates a new loop with a single iteration to process the aggregation results.
// This approach allows the aggregation to be processed in a consistent way with other operations.
func (c *LoopCollectCompiler) compileGlobalAggregation(spec *core.CollectorSpec) {
	// At this point, the previous loop is finalized, so we can pop it and free its registers
	prevLoop := c.ctx.Loops.Pop()
	c.ctx.Registers.Free(prevLoop.Key)
	c.ctx.Registers.Free(prevLoop.Value)
	c.ctx.Registers.Free(prevLoop.Src)

	// Create a new loop with 1 iteration only to process the aggregation
	c.ctx.Symbols.EnterScope()
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

	// Process the aggregation function calls using the values from the previous loop's collector
	c.compileAggregationFuncCalls(spec, prevLoop.Dst)

	// Free the previous loop's destination register
	c.ctx.Registers.Free(prevLoop.Dst)
}

// compileAggregationFuncCalls processes the aggregation function calls for the selectors.
// It loads the arguments from the aggregator, calls the aggregation functions,
// and assigns the results to local variables.
// It also handles the case where there are no records in the aggregator by loading NONE values.
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

	// We skip the key retrieval and function call if there are no records in the accumulator
	c.ctx.Emitter.EmitJumpIfTrue(cond, elseLabel)

	selectors := spec.AggregationSelectors()
	selectorVarRegs := make([]vm.Operand, len(selectors))

	// Process each aggregation selector
	for i, selector := range selectors {
		var args core.RegisterSequence

		// We need to unpack arguments from the aggregator
		if selector.Args() > 1 {
			// For multiple arguments, allocate a sequence and load each argument by its indexed key
			args = c.ctx.Registers.AllocateSequence(selector.Args())

			for y, reg := range args {
				argKeyReg := c.loadAggregationArgKey(selector.Name(), y)
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

// compileAggregationFuncCall processes a single aggregation function call.
// It declares a local variable for the aggregation result and loads the value from the selector register.
// This is used for grouped aggregations where the selector values are stored in registers.
func (c *LoopCollectCompiler) compileAggregationFuncCall(selector *core.AggregateSelector) {
	// Declare a local variable with the selector name
	varReg := c.ctx.Symbols.DeclareLocal(selector.Name().String(), core.TypeUnknown)
	// Load the value from index 1 of the selector register (index 0 is the original value)
	loadIndex(c.ctx, varReg, selector.Register(), 1)
}

// loadAggregationArgKey creates a key for an aggregation argument by combining the selector name and argument index.
// This is used for global aggregations with multiple arguments to store each argument separately.
// Returns a register containing the key as a string constant.
func (c *LoopCollectCompiler) loadAggregationArgKey(selector runtime.String, arg int) vm.Operand {
	// Create a key with format "selectorName:argIndex"
	argKey := selector.String() + ":" + strconv.Itoa(arg)
	// Load the key as a string constant
	return loadConstant(c.ctx, runtime.String(argKey))
}
