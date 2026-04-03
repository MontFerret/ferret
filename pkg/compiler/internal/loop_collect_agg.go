package internal

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/antlr4-go/antlr/v4"

	"github.com/MontFerret/ferret/v2/pkg/diagnostics"

	"github.com/MontFerret/ferret/v2/pkg/bytecode"

	"github.com/MontFerret/ferret/v2/pkg/compiler/internal/core"
	parsing "github.com/MontFerret/ferret/v2/pkg/parser/diagnostics"
	"github.com/MontFerret/ferret/v2/pkg/parser/fql"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

// initializeAggregation processes the aggregation selectors from the COLLECT clause.
// It handles both grouped and global aggregations differently.
// For grouped aggregations, it compiles the selectors and packs them with the loop value.
// For global aggregations, it pushes the selectors directly to the collector.
// Returns a slice of AggregateSelectors that describe the aggregation operations.
func (c *CollectCompiler) initializeAggregation(ctx fql.ICollectAggregatorContext, dst bytecode.Operand, kv *core.KV, withGrouping bool, plan *bytecode.AggregatePlan) *core.CollectorAggregation {
	loop := c.ctx.Loops.Current()
	selectors := ctx.AllCollectAggregateSelector()

	// If we have grouping, we need to pack the selectors into the collector value
	if withGrouping {
		if plan != nil {
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
	aggregateSelectors := c.initializeGlobalAggregationSelectors(selectors, dst, plan != nil)

	return core.NewCollectorAggregation(dst, aggregateSelectors)
}

func (c *CollectCompiler) reportAggregateSemanticError(ctx antlr.ParserRuleContext, message, hint string) {
	var err *diagnostics.Diagnostic

	if ctx != nil {
		err = c.ctx.Errors.Create(parsing.SemanticError, ctx, message)
	} else {
		err = &diagnostics.Diagnostic{
			Source:  c.ctx.Source,
			Kind:    parsing.SemanticError,
			Message: message,
		}
	}

	err.Hint = hint
	c.ctx.Errors.Add(err)
}

// parseAggregateSelector validates and compiles one aggregate selector.
// It returns a normalized selector representation shared by grouped/global initializers.
func (c *CollectCompiler) parseAggregateSelector(selector fql.ICollectAggregateSelectorContext) (*core.CompiledAggregateSelector, bool) {
	name := runtime.String(textOfBindingIdentifier(selector.BindingIdentifier()))
	fcx := selector.FunctionCallExpression()

	if fcx == nil || fcx.FunctionCall() == nil {
		c.reportAggregateSemanticError(
			selector,
			fmt.Sprintf("Invalid aggregate selector for '%s'", name.String()),
			"Use a function call expression, e.g. AGGREGATE total = COUNT(x).",
		)

		return nil, false
	}

	funcName := c.front.Calls.ResolveFunctionName(fcx.FunctionCall())
	args := c.front.Expressions.CompileArgumentList(fcx.FunctionCall().ArgumentList())

	if len(args) == 0 {
		c.reportAggregateSemanticError(
			selector,
			fmt.Sprintf("Aggregate '%s' requires at least one argument", funcName.String()),
			fmt.Sprintf("Provide at least one expression, e.g. %s(x).", funcName.String()),
		)
		return nil, false
	}

	plan := c.front.Recovery.CollectPlan(fcx, core.RecoveryPlanOptions{})

	return core.NewCompiledAggregateSelector(
		name,
		args,
		funcName,
		fcx.ErrorOperator() != nil || hasErrorReturnNoneHandler(plan),
		selector,
	), true
}

// initializeAggregationSelectors encapsulates the common selector compilation flow.
// The caller provides only the mode-specific emission strategy via emit callback.
func (c *CollectCompiler) initializeAggregationSelectors(
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
func (c *CollectCompiler) emitAggregateArgs(dst bytecode.Operand, args core.RegisterSequence, keyForArg func(argIndex int) bytecode.Operand) {
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
func (c *CollectCompiler) initializeGroupedAggregationSelectors(selectors []fql.ICollectAggregateSelectorContext, kv *core.KV, dst bytecode.Operand, fused bool) []*core.AggregateSelector {
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

			c.ctx.Emitter.EmitAggregateGroupUpdate(dst, kv.Key, parsed.Args()[0], i)

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
func (c *CollectCompiler) initializeGlobalAggregationSelectors(selectors []fql.ICollectAggregateSelectorContext, dst bytecode.Operand, planBacked bool) []*core.AggregateSelector {
	return c.initializeAggregationSelectors(selectors, func(i int, parsed *core.CompiledAggregateSelector) bool {
		if planBacked {
			if len(parsed.Args()) != 1 {
				return false
			}

			c.ctx.Emitter.EmitAggregateUpdate(dst, parsed.Args()[0], i)

			return true
		}

		c.emitAggregateArgs(dst, parsed.Args(), func(argIndex int) bytecode.Operand {
			if argIndex >= 0 {
				// Multi-argument selectors are stored as <selectorName>:<argIndex>.
				return c.loadGlobalSelectorKey(parsed.Name(), argIndex)
			}

			return c.front.TypeFacts.LoadConstant(parsed.Name())
		})

		return true
	})
}

func (c *CollectCompiler) buildGlobalAggregatePlan(ctx fql.ICollectAggregatorContext) (*bytecode.AggregatePlan, bool) {
	return c.buildAggregatePlan(ctx.AllCollectAggregateSelector(), false)
}

func (c *CollectCompiler) buildGroupedAggregatePlan(selectors []fql.ICollectAggregateSelectorContext, trackGroupValues bool) (*bytecode.AggregatePlan, bool) {
	return c.buildAggregatePlan(selectors, trackGroupValues)
}

// buildAggregatePlan recognizes selectors that can be handled by VM-native aggregate collectors.
// It intentionally accepts only unprotected calls with exactly one argument.
func (c *CollectCompiler) buildAggregatePlan(selectors []fql.ICollectAggregateSelectorContext, trackGroupValues bool) (*bytecode.AggregatePlan, bool) {
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

		plan := c.front.Recovery.CollectPlan(fce, core.RecoveryPlanOptions{})
		if fce.ErrorOperator() != nil || plan.OnError != nil {
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

		funcName := c.front.Calls.ResolveFunctionName(fce.FunctionCall())
		kind, ok := aggregateKind(funcName)

		if !ok {
			return nil, false
		}

		keys = append(keys, runtime.String(textOfBindingIdentifier(selector.BindingIdentifier())))
		kinds = append(kinds, kind)
	}

	plan := bytecode.NewAggregatePlan(keys, kinds, trackGroupValues)

	return &plan, true
}

func (c *CollectCompiler) shouldFuseGroupedAggregation(grouping fql.ICollectGroupingContext, selectors []fql.ICollectAggregateSelectorContext) bool {
	// Heuristic:
	// - only fuse when there are 3+ aggregates (to amortize setup cost)
	// - only one group key (avoid complex composite/group-array keys)
	if grouping == nil {
		return false
	}

	if len(selectors) < 3 {
		return false
	}

	groupSelectors := grouping.AllCollectSelector()
	return len(groupSelectors) == 1
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
func (c *CollectCompiler) finalizeAggregation(spec *core.Collector) {
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
func (c *CollectCompiler) finalizeGroupedAggregation(spec *core.Collector) {
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
func (c *CollectCompiler) finalizeGlobalAggregation(spec *core.Collector) {
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
		loop.Dst = c.ctx.Loops.RequiredParent(c.ctx.Loops.Depth()).Dst
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
func (c *CollectCompiler) compileGlobalAggregationFuncCalls(spec *core.Collector) {
	aggregator, elseLabel, endLabel := c.emitGlobalAggregationEmptyGuard(spec)
	selectorRegs := c.compileGlobalAggregationSelectors(spec, aggregator)
	projVar := c.compileGlobalAggregationProjection(spec, aggregator)
	c.emitGlobalAggregationEmptyFallback(selectorRegs, projVar, elseLabel, endLabel)
}

func (c *CollectCompiler) emitGlobalAggregationEmptyGuard(spec *core.Collector) (bytecode.Operand, core.Label, core.Label) {
	aggregator := spec.Destination()
	cond := c.ctx.Registers.Allocate()
	c.ctx.Emitter.EmitAB(bytecode.OpLength, cond, aggregator)

	zero := c.front.TypeFacts.LoadConstant(runtime.ZeroInt)
	c.ctx.Emitter.EmitEq(cond, cond, zero)

	elseLabel := c.ctx.Emitter.NewLabel()
	endLabel := c.ctx.Emitter.NewLabel()
	c.ctx.Emitter.EmitJumpIfTrue(cond, elseLabel)

	return aggregator, elseLabel, endLabel
}

func (c *CollectCompiler) compileGlobalAggregationSelectors(spec *core.Collector, aggregator bytecode.Operand) []bytecode.Operand {
	selectors := spec.Aggregation().Selectors()
	selectorRegs := make([]bytecode.Operand, len(selectors))

	if spec.Type() == bytecode.CollectorTypeAggregate {
		c.compilePlanBackedGlobalAggregationSelectors(selectors, aggregator, selectorRegs)
		return selectorRegs
	}

	c.compileGenericGlobalAggregationSelectors(selectors, aggregator, selectorRegs)
	return selectorRegs
}

func (c *CollectCompiler) compilePlanBackedGlobalAggregationSelectors(selectors []*core.AggregateSelector, aggregator bytecode.Operand, selectorRegs []bytecode.Operand) {
	for i, selector := range selectors {
		varReg := c.declareLocalOrReport(selector.Context(), selector.Name().String(), core.TypeUnknown)
		selectorRegs[i] = varReg

		key := c.front.TypeFacts.LoadConstant(selector.Name())
		c.ctx.Emitter.EmitABC(bytecode.OpLoadKey, varReg, aggregator, key)
	}
}

func (c *CollectCompiler) compileGenericGlobalAggregationSelectors(selectors []*core.AggregateSelector, aggregator bytecode.Operand, selectorRegs []bytecode.Operand) {
	for i, selector := range selectors {
		args := c.compileGlobalAggregationSelectorArgs(selector, aggregator)
		result := c.front.Expressions.CompileFunctionCallByNameWith(nil, selector.FuncName(), selector.ProtectedCall(), args)

		varReg := c.declareLocalOrReport(selector.Context(), selector.Name().String(), core.TypeUnknown)
		selectorRegs[i] = varReg
		c.ctx.Emitter.EmitMoveTracked(varReg, result)
	}
}

func (c *CollectCompiler) compileGlobalAggregationSelectorArgs(selector *core.AggregateSelector, aggregator bytecode.Operand) core.RegisterSequence {
	if selector.Args() <= 1 {
		key := c.front.TypeFacts.LoadConstant(selector.Name())
		value := c.ctx.Registers.Allocate()
		c.ctx.Emitter.EmitABC(bytecode.OpLoadKey, value, aggregator, key)
		return core.RegisterSequence{value}
	}

	args := c.ctx.Registers.AllocateSequence(selector.Args())
	for idx, reg := range args {
		argKey := c.loadGlobalSelectorKey(selector.Name(), idx)
		c.ctx.Emitter.EmitABC(bytecode.OpLoadKey, reg, aggregator, argKey)
	}

	return args
}

func (c *CollectCompiler) compileGlobalAggregationProjection(spec *core.Collector, aggregator bytecode.Operand) bytecode.Operand {
	if !spec.HasProjection() {
		return bytecode.NoopOperand
	}

	if projectionState := spec.ProjectionState(); projectionState != bytecode.NoopOperand {
		varName := spec.Projection().VariableName()
		val := c.declareLocalOrReport(spec.Projection().Context(), varName, core.TypeArray)
		c.front.TypeFacts.EmitMoveAuto(val, projectionState)

		return val
	}

	return c.finalizeProjection(spec, aggregator)
}

func (c *CollectCompiler) emitGlobalAggregationEmptyFallback(selectorRegs []bytecode.Operand, projVar bytecode.Operand, elseLabel, endLabel core.Label) {
	c.ctx.Emitter.EmitJump(endLabel)
	c.ctx.Emitter.MarkLabel(elseLabel)

	for _, reg := range selectorRegs {
		c.ctx.Emitter.EmitA(bytecode.OpLoadNone, reg)
		c.ctx.Types.Set(reg, core.TypeNone)
	}

	if projVar != bytecode.NoopOperand {
		c.ctx.Emitter.EmitA(bytecode.OpLoadNone, projVar)
		c.ctx.Types.Set(projVar, core.TypeNone)
	}

	c.ctx.Emitter.MarkLabel(endLabel)
}

func (c *CollectCompiler) compileGroupedAggregationFuncCall(selector *core.AggregateSelector, aggregator bytecode.Operand, idx int, fused bool) {
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

	resArg := c.front.Expressions.CompileFunctionCallByNameWith(nil, selector.FuncName(), selector.ProtectedCall(), args)

	c.ctx.Emitter.EmitMove(valReg, resArg)
}

// loadGlobalSelectorKey creates a key for an aggregation argument by combining the selector name and argument index.
// This is used for global aggregations with multiple arguments to store each argument separately.
// Returns a register containing the key as a string constant.
func (c *CollectCompiler) loadGlobalSelectorKey(selector runtime.String, arg int) bytecode.Operand {
	// Create a key with format "selectorName:argIndex"
	argKey := selector.String() + ":" + strconv.Itoa(arg)
	// Load the key as a string constant
	return c.front.TypeFacts.LoadConstant(runtime.String(argKey))
}

func (c *CollectCompiler) loadSelectorKey(key bytecode.Operand, selector runtime.String, arg int) bytecode.Operand {
	selectorKey := c.ctx.Registers.Allocate()
	selectorName := c.front.TypeFacts.LoadConstant(selector)

	c.ctx.Emitter.EmitABC(bytecode.OpAdd, selectorKey, key, selectorName)

	if arg >= 0 {
		selectorIndex := c.front.TypeFacts.LoadConstant(runtime.String(strconv.Itoa(arg)))
		c.ctx.Emitter.EmitABC(bytecode.OpAdd, selectorKey, selectorKey, selectorIndex)
	}

	return selectorKey
}

func (c *CollectCompiler) loadGroupedAggregateKey(key bytecode.Operand, selectorIdx int) bytecode.Operand {
	aggKey := c.ctx.Registers.Allocate()
	selectorConst := c.ctx.Symbols.AddConstant(runtime.NewInt(selectorIdx))
	c.ctx.Emitter.EmitLoadAggregateKey(aggKey, key, selectorConst)

	return aggKey
}
