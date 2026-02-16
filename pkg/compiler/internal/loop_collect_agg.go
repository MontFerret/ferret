package internal

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/antlr4-go/antlr/v4"

	"github.com/MontFerret/ferret/v2/pkg/bytecode"

	"github.com/MontFerret/ferret/v2/pkg/compiler/internal/core"
	compilerdiagnostics "github.com/MontFerret/ferret/v2/pkg/compiler/internal/diagnostics"
	"github.com/MontFerret/ferret/v2/pkg/parser/fql"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

// initializeAggregation processes the aggregation selectors from the COLLECT clause.
// It handles both grouped and global aggregations differently.
// For grouped aggregations, it compiles the selectors and packs them with the loop value.
// For global aggregations, it pushes the selectors directly to the collector.
// Returns a slice of AggregateSelectors that describe the aggregation operations.
func (c *LoopCollectCompiler) initializeAggregation(ctx fql.ICollectAggregatorContext, dst bytecode.Operand, kv *core.KV, withGrouping bool, groupedPlan *bytecode.AggregatePlan) *core.CollectorAggregation {
	loop := c.ctx.Loops.Current()
	selectors := ctx.AllCollectAggregateSelector()

	// If we have grouping, we need to pack the selectors into the collector value
	if withGrouping {
		if groupedPlan != nil {
			// Fused grouped aggregation writes directly into the primary collector.
			aggregateSelectors := c.initializeGroupedAggregationSelectors(selectors, kv, dst, true)
			return core.NewCollectorAggregationFused(dst, aggregateSelectors)
		}

		aggregator := c.ctx.Registers.Allocate()

		// We create a separate collector for aggregation in grouped mode
		c.ctx.Emitter.InsertAx(loop.StartLabel(), bytecode.OpDataSetCollector, aggregator, int(bytecode.CollectorTypeKeyGroup))

		// Compile selectors for grouped aggregation
		aggregateSelectors := c.initializeGroupedAggregationSelectors(selectors, kv, aggregator, false)

		return core.NewCollectorAggregation(aggregator, aggregateSelectors)
	}

	// For global aggregation, we just push the selectors into the global collector
	aggregateSelectors := c.initializeGlobalAggregationSelectors(selectors, dst)

	return core.NewCollectorAggregation(dst, aggregateSelectors)
}

func (c *LoopCollectCompiler) reportAggregateSemanticError(ctx antlr.ParserRuleContext, message, hint string) {
	var err *compilerdiagnostics.CompilationError

	if ctx != nil {
		err = c.ctx.Errors.Create(compilerdiagnostics.SemanticError, ctx, message)
	} else {
		err = compilerdiagnostics.NewError(c.ctx.Source, compilerdiagnostics.SemanticError, message)
	}

	err.Hint = hint
	c.ctx.Errors.Add(err)
}

// parseAggregateSelector validates and compiles one aggregate selector.
// It returns a normalized selector representation shared by grouped/global initializers.
func (c *LoopCollectCompiler) parseAggregateSelector(selector fql.ICollectAggregateSelectorContext) (*core.CompiledAggregateSelector, bool) {
	name := runtime.String(selector.Identifier().GetText())
	fcx := selector.FunctionCallExpression()

	if fcx == nil || fcx.FunctionCall() == nil {
		c.reportAggregateSemanticError(
			selector,
			fmt.Sprintf("Invalid aggregate selector for '%s'", name.String()),
			"Use a function call expression, e.g. AGGREGATE total = COUNT(x).",
		)

		return nil, false
	}

	funcName := getFunctionName(fcx.FunctionCall(), c.ctx.UseAliases)
	args := c.ctx.ExprCompiler.CompileArgumentList(fcx.FunctionCall().ArgumentList())

	if len(args) == 0 {
		c.reportAggregateSemanticError(
			selector,
			fmt.Sprintf("Aggregate '%s' requires at least one argument", funcName.String()),
			fmt.Sprintf("Provide at least one expression, e.g. %s(x).", funcName.String()),
		)
		return nil, false
	}

	return core.NewCompiledAggregateSelector(
		name,
		args,
		funcName,
		fcx.ErrorOperator() != nil,
		selector,
	), true
}

// initializeAggregationSelectors encapsulates the common selector compilation flow.
// The caller provides only the mode-specific emission strategy via emit callback.
func (c *LoopCollectCompiler) initializeAggregationSelectors(
	selectors []fql.ICollectAggregateSelectorContext,
	emit func(int, *core.CompiledAggregateSelector) bool,
) []*core.AggregateSelector {
	wrappedSelectors := make([]*core.AggregateSelector, 0, len(selectors))

	for i, selector := range selectors {
		parsed, ok := c.parseAggregateSelector(selector)

		if !ok {
			return nil
		}

		if emit != nil && !emit(i, parsed) {
			return nil
		}

		wrappedSelectors = append(
			wrappedSelectors,
			core.NewAggregateSelector(parsed.Name(), len(parsed.Args()), parsed.FuncName(), parsed.ProtectedCall(), parsed.Context()),
		)
	}

	return wrappedSelectors
}

// emitAggregateArgs emits selector arguments as collector key/value pairs.
// keyForArg receives:
//
//	-1 for a single-argument selector
//	>=0 for each argument index in multi-argument selectors.
func (c *LoopCollectCompiler) emitAggregateArgs(dst bytecode.Operand, args core.RegisterSequence, keyForArg func(argIndex int) bytecode.Operand) {
	if len(args) > 1 {
		for i, arg := range args {
			c.ctx.Emitter.EmitPushKV(dst, keyForArg(i), arg)
		}

		return
	}

	c.ctx.Emitter.EmitPushKV(dst, keyForArg(-1), args[0])
}

// initializeGroupedAggregationSelectors processes aggregation selectors for grouped aggregations.
// It compiles each selector's function call expression and arguments, and creates AggregateSelector objects.
// For selectors with multiple arguments, it packs them into an array.
// Returns a slice of AggregateSelectors that describe the aggregation operations.
func (c *LoopCollectCompiler) initializeGroupedAggregationSelectors(selectors []fql.ICollectAggregateSelectorContext, kv *core.KV, dst bytecode.Operand, fused bool) []*core.AggregateSelector {
	return c.initializeAggregationSelectors(selectors, func(i int, parsed *core.CompiledAggregateSelector) bool {
		if fused {
			// Fused grouped mode routes directly to the VM grouped aggregate collector,
			// so each selector must map to exactly one argument.
			if len(parsed.Args()) != 1 {
				c.reportAggregateSemanticError(
					parsed.Context(),
					fmt.Sprintf("Aggregate '%s' requires exactly one argument in fused grouped mode", parsed.FuncName().String()),
					"Use a single argument in each aggregate function call when grouped aggregation fusion is enabled.",
				)

				return false
			}

			key := c.loadGroupedAggregateKey(kv.Key, i)
			c.ctx.Emitter.EmitPushKV(dst, key, parsed.Args()[0])

			return true
		}

		c.emitAggregateArgs(dst, parsed.Args(), func(argIndex int) bytecode.Operand {
			return c.loadSelectorKey(kv.Key, parsed.Name(), argIndex)
		})

		return true
	})
}

// initializeGlobalAggregationSelectors processes aggregation selectors for global (non-grouped) aggregations.
// It compiles each selector's function call expression and arguments, and pushes them directly to the collector.
// For selectors with multiple arguments, it uses indexed keys to store each argument separately.
// Returns a slice of AggregateSelectors that describe the aggregation operations.
func (c *LoopCollectCompiler) initializeGlobalAggregationSelectors(selectors []fql.ICollectAggregateSelectorContext, dst bytecode.Operand) []*core.AggregateSelector {
	return c.initializeAggregationSelectors(selectors, func(_ int, parsed *core.CompiledAggregateSelector) bool {
		c.emitAggregateArgs(dst, parsed.Args(), func(argIndex int) bytecode.Operand {
			if argIndex >= 0 {
				// Multi-argument selectors are stored as <selectorName>:<argIndex>.
				return c.loadGlobalSelectorKey(parsed.Name(), argIndex)
			}

			return loadConstant(c.ctx, parsed.Name())
		})

		return true
	})
}

func (c *LoopCollectCompiler) buildGlobalAggregatePlan(ctx fql.ICollectAggregatorContext) (*bytecode.AggregatePlan, bool) {
	return c.buildAggregatePlan(ctx.AllCollectAggregateSelector())
}

func (c *LoopCollectCompiler) buildGroupedAggregatePlan(selectors []fql.ICollectAggregateSelectorContext) (*bytecode.AggregatePlan, bool) {
	return c.buildAggregatePlan(selectors)
}

// buildAggregatePlan recognizes selectors that can be handled by VM-native aggregate collectors.
// It intentionally accepts only unprotected calls with exactly one argument.
func (c *LoopCollectCompiler) buildAggregatePlan(selectors []fql.ICollectAggregateSelectorContext) (*bytecode.AggregatePlan, bool) {
	if len(selectors) == 0 {
		return nil, false
	}

	keys := make([]runtime.String, 0, len(selectors))
	kinds := make([]bytecode.AggregateKind, 0, len(selectors))

	for _, selector := range selectors {
		fce := selector.FunctionCallExpression()
		if fce == nil || fce.FunctionCall() == nil {
			return nil, false
		}

		if fce.ErrorOperator() != nil {
			return nil, false
		}

		args := fce.FunctionCall().ArgumentList()
		if args == nil {
			return nil, false
		}

		exps := args.AllExpression()
		if len(exps) != 1 {
			return nil, false
		}

		funcName := getFunctionName(fce.FunctionCall(), c.ctx.UseAliases)
		kind, ok := aggregateKind(funcName)

		if !ok {
			return nil, false
		}

		keys = append(keys, runtime.String(selector.Identifier().GetText()))
		kinds = append(kinds, kind)
	}

	plan := bytecode.NewAggregatePlan(keys, kinds)

	return &plan, true
}

func (c *LoopCollectCompiler) shouldFuseGroupedAggregation(grouping fql.ICollectGroupingContext, selectors []fql.ICollectAggregateSelectorContext) bool {
	// Heuristic:
	// - only fuse when there are 3+ aggregates (to amortize setup cost)
	// - only one group key (avoid complex composite/group-array keys)
	// - only simple group expressions (identifier/param/simple member path/string literal)
	if grouping == nil {
		return false
	}

	if len(selectors) < 3 {
		return false
	}

	groupSelectors := grouping.AllCollectSelector()
	if len(groupSelectors) != 1 {
		return false
	}

	return isSimpleGroupExpression(groupSelectors[0].Expression())
}

func isSimpleGroupExpression(expr fql.IExpressionContext) bool {
	if expr == nil {
		return false
	}

	predicate := expr.Predicate()
	if predicate == nil {
		return false
	}

	atom := predicate.ExpressionAtom()
	if atom == nil {
		return false
	}

	// Disallow operators and complex expression forms.
	if atom.AdditiveOperator() != nil || atom.MultiplicativeOperator() != nil || atom.RegexpOperator() != nil {
		return false
	}

	if atom.FunctionCallExpression() != nil || atom.RangeOperator() != nil || atom.DispatchExpression() != nil || atom.ForExpression() != nil || atom.WaitForExpression() != nil {
		return false
	}

	if lit := atom.Literal(); lit != nil {
		return lit.StringLiteral() != nil ||
			lit.IntegerLiteral() != nil ||
			lit.FloatLiteral() != nil ||
			lit.BooleanLiteral() != nil ||
			lit.NoneLiteral() != nil
	}

	if atom.Variable() != nil || atom.Param() != nil {
		return true
	}

	if me := atom.MemberExpression(); me != nil {
		source := me.MemberExpressionSource()
		if source == nil || (source.Variable() == nil && source.Param() == nil) {
			return false
		}

		return isSimpleMemberPathChain(me.AllMemberExpressionPath())
	}

	return false
}

func aggregateKind(name runtime.String) (bytecode.AggregateKind, bool) {
	fn := strings.ToUpper(name.String())
	if strings.Contains(fn, runtime.NamespaceSeparator) {
		parts := strings.Split(fn, runtime.NamespaceSeparator)
		fn = parts[len(parts)-1]
	}

	switch fn {
	case "COUNT":
		return bytecode.AggregateCount, true
	case "SUM":
		return bytecode.AggregateSum, true
	case "MIN":
		return bytecode.AggregateMin, true
	case "MAX":
		return bytecode.AggregateMax, true
	case "AVERAGE":
		return bytecode.AggregateAverage, true
	default:
		return 0, false
	}
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
// In fused mode, aggregate values are loaded directly from the primary grouped collector.
// In non-fused mode, selectors are evaluated from the auxiliary per-group collector.
func (c *LoopCollectCompiler) finalizeGroupedAggregation(spec *core.Collector) {
	loop := c.ctx.Loops.Current()
	aggregator := spec.Aggregation().State()

	if spec.Aggregation().IsFused() {
		// In fused mode the collector register can be reused for the output dataset,
		// so load aggregate values from the loop source instead.
		aggregator = loop.Src
	}

	for i, selector := range spec.Aggregation().Selectors() {
		c.compileGroupedAggregationFuncCall(selector, aggregator, i, spec.Aggregation().IsFused())
	}
}

// finalizeGlobalAggregation handles global (non-grouped) aggregation operations.
// It creates a new loop with a single iteration to process the aggregation results.
// This approach allows the aggregation to be processed in a consistent way with other operations.
func (c *LoopCollectCompiler) finalizeGlobalAggregation(spec *core.Collector) {
	// At this point, the previous loop is finalized, so we can pop it and free its registers
	prevLoop := c.ctx.Loops.Pop()

	// Create a new loop with 1 iteration only to process the aggregation
	loop := c.ctx.Loops.NewLoop(core.ForInLoop, core.NormalLoop, prevLoop.Distinct)
	c.ctx.Loops.Push(loop)

	// Set up the loop source to be a range from 0 to 0 (one iteration)
	loop.Src = c.ctx.Registers.Allocate()
	zero := c.ctx.Registers.Allocate()
	c.ctx.Emitter.EmitA(bytecode.OpLoadZero, zero)
	c.ctx.Emitter.EmitABC(bytecode.OpLoadRange, loop.Src, zero, zero)

	// Inherit allocation flag from previous loop
	loop.Allocate = prevLoop.Allocate

	// If not allocating, use the parent loop's destination
	if !loop.Allocate {
		parent := c.ctx.Loops.FindParent(c.ctx.Loops.Depth())
		loop.Dst = parent.Dst
	}

	// Initialize the loop
	loop.EmitInitialization(c.ctx.Registers, c.ctx.Emitter)

	// Evaluate aggregate outputs and projection locals inside the synthetic one-iteration loop.
	c.compileGlobalAggregationFuncCalls(spec)
}

// compileGlobalAggregationFuncCalls processes the aggregation function calls for the selectors.
// It loads the arguments from the aggregator, calls the aggregation functions,
// and assigns the results to local variables.
// It also handles the case where there are no records in the aggregator by loading NONE values.
func (c *LoopCollectCompiler) compileGlobalAggregationFuncCalls(spec *core.Collector) {
	// Gets the number of records in the accumulator
	aggregator := spec.Destination()
	cond := c.ctx.Registers.Allocate()
	c.ctx.Emitter.EmitAB(bytecode.OpLength, cond, aggregator)

	zero := loadConstant(c.ctx, runtime.ZeroInt)
	// Check if the number equals to zero
	c.ctx.Emitter.EmitEq(cond, cond, zero)

	elseLabel := c.ctx.Emitter.NewLabel()
	endLabel := c.ctx.Emitter.NewLabel()

	// We skip the key retrieval and function call if there are no records in the accumulator
	c.ctx.Emitter.EmitJumpIfTrue(cond, elseLabel)

	selectors := spec.Aggregation().Selectors()
	selectorVarRegs := make([]bytecode.Operand, len(selectors))

	if spec.Type() == bytecode.CollectorTypeAggregate {
		// Plan-backed aggregate collector already stores final values by selector name.
		for i, selector := range selectors {
			selectorVarName := selector.Name()
			varReg := c.declareLocalOrReport(selector.Context(), selectorVarName.String(), core.TypeUnknown)
			selectorVarRegs[i] = varReg

			key := loadConstant(c.ctx, selector.Name())
			c.ctx.Emitter.EmitABC(bytecode.OpLoadKey, varReg, aggregator, key)
		}
	} else {
		// Generic collector path: reconstruct selector args and evaluate aggregate functions in bytecode.
		// Process each aggregation selector
		for i, selector := range selectors {
			var args core.RegisterSequence

			// We need to unpack arguments from the aggregator
			if selector.Args() > 1 {
				// For multiple arguments, allocate a sequence and load each argument by its indexed key
				args = c.ctx.Registers.AllocateSequence(selector.Args())

				for y, reg := range args {
					argKeyReg := c.loadGlobalSelectorKey(selector.Name(), y)
					c.ctx.Emitter.EmitABC(bytecode.OpLoadKey, reg, aggregator, argKeyReg)
				}
			} else {
				// For a single argument, load it directly using the selector name as key
				key := loadConstant(c.ctx, selector.Name())
				value := c.ctx.Registers.Allocate()
				c.ctx.Emitter.EmitABC(bytecode.OpLoadKey, value, aggregator, key)
				args = core.RegisterSequence{value}
			}

			// Call the aggregation function with the loaded arguments
			result := c.ctx.ExprCompiler.CompileFunctionCallByNameWith(selector.FuncName(), selector.ProtectedCall(), args)

			// Declare a local variable for the aggregation result
			selectorVarName := selector.Name()
			varReg := c.declareLocalOrReport(selector.Context(), selectorVarName.String(), core.TypeUnknown)
			selectorVarRegs[i] = varReg

			// Move the function result to the variable
			c.ctx.Emitter.EmitAB(bytecode.OpMove, varReg, result)
		}
	}

	var projVar bytecode.Operand

	// If the projection is used, we allocate a new register for the variable and put the iterator's value into it
	if spec.HasProjection() {
		projVar = c.finalizeProjection(spec, aggregator)
	}

	// Skip the else block (for empty aggregator)
	c.ctx.Emitter.EmitJump(endLabel)
	c.ctx.Emitter.MarkLabel(elseLabel)

	// If there are no records in the aggregator, load NONE values for all variables
	for _, varReg := range selectorVarRegs {
		c.ctx.Emitter.EmitA(bytecode.OpLoadNone, varReg)
		c.ctx.Types.Set(varReg, core.TypeNone)
	}

	if projVar != bytecode.NoopOperand {
		c.ctx.Emitter.EmitA(bytecode.OpLoadNone, projVar)
		c.ctx.Types.Set(projVar, core.TypeNone)
	}

	c.ctx.Emitter.MarkLabel(endLabel)
}

func (c *LoopCollectCompiler) compileGroupedAggregationFuncCall(selector *core.AggregateSelector, aggregator bytecode.Operand, idx int, fused bool) {
	loop := c.ctx.Loops.Current()
	// Declare a local variable with the selector name
	valReg := c.declareLocalOrReport(selector.Context(), selector.Name().String(), core.TypeUnknown)

	if fused {
		// Fused grouped collector exposes aggregate values under tagged keys.
		key := c.loadGroupedAggregateKey(loop.Key, idx)
		c.ctx.Emitter.EmitABC(bytecode.OpLoadKey, valReg, aggregator, key)

		return
	}

	var args core.RegisterSequence

	// We need to unpack arguments from the aggregator
	if selector.Args() > 1 {
		// For multiple arguments, allocate a sequence and load each argument by its indexed key
		args = c.ctx.Registers.AllocateSequence(selector.Args())

		for y, reg := range args {
			key := c.loadSelectorKey(loop.Key, selector.Name(), y)
			c.ctx.Emitter.EmitABC(bytecode.OpLoadKey, reg, aggregator, key)
		}
	} else {
		// For a single argument, load it directly using the selector name as key
		key := c.loadSelectorKey(loop.Key, selector.Name(), -1)
		value := c.ctx.Registers.Allocate()
		c.ctx.Emitter.EmitABC(bytecode.OpLoadKey, value, aggregator, key)
		args = core.RegisterSequence{value}
	}

	resArg := c.ctx.ExprCompiler.CompileFunctionCallByNameWith(selector.FuncName(), selector.ProtectedCall(), args)

	c.ctx.Emitter.EmitMove(valReg, resArg)
}

// loadGlobalSelectorKey creates a key for an aggregation argument by combining the selector name and argument index.
// This is used for global aggregations with multiple arguments to store each argument separately.
// Returns a register containing the key as a string constant.
func (c *LoopCollectCompiler) loadGlobalSelectorKey(selector runtime.String, arg int) bytecode.Operand {
	// Create a key with format "selectorName:argIndex"
	argKey := selector.String() + ":" + strconv.Itoa(arg)
	// Load the key as a string constant
	return loadConstant(c.ctx, runtime.String(argKey))
}

func (c *LoopCollectCompiler) loadSelectorKey(key bytecode.Operand, selector runtime.String, arg int) bytecode.Operand {
	selectorKey := c.ctx.Registers.Allocate()
	selectorName := loadConstant(c.ctx, selector)

	c.ctx.Emitter.EmitABC(bytecode.OpAdd, selectorKey, key, selectorName)

	if arg >= 0 {
		selectorIndex := loadConstant(c.ctx, runtime.String(strconv.Itoa(arg)))
		c.ctx.Emitter.EmitABC(bytecode.OpAdd, selectorKey, selectorKey, selectorIndex)
	}

	return selectorKey
}

func (c *LoopCollectCompiler) loadGroupedAggregateKey(key bytecode.Operand, selectorIdx int) bytecode.Operand {
	aggKey := c.ctx.Registers.Allocate()
	// Use a tagged array key to avoid collisions with user-provided group keys.
	// Format: [AggregateKeyMarker, <groupKey>, <selectorIdx>]
	c.ctx.Emitter.EmitArray(aggKey, 3)
	c.ctx.Emitter.EmitArrayPush(aggKey, loadConstant(c.ctx, bytecode.AggregateKeyMarker))
	c.ctx.Emitter.EmitArrayPush(aggKey, key)
	c.ctx.Emitter.EmitArrayPush(aggKey, loadConstant(c.ctx, runtime.Int(selectorIdx)))

	return aggKey
}
