package internal

import (
	"github.com/MontFerret/ferret/v2/pkg/bytecode"
	"github.com/MontFerret/ferret/v2/pkg/compiler/internal/core"
	"github.com/MontFerret/ferret/v2/pkg/parser/fql"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

// CollectCompiler handles the compilation of COLLECT clauses in FQL queries.
// It transforms COLLECT operations into VM instructions for data aggregation and grouping.
type CollectCompiler struct {
	ctx      *CompilationSession
	bindings *BindingCompiler
	calls    *CallResolver
	exprs    *ExprCompiler
	recovery *RecoveryCompiler
	facts    *TypeFacts
}

// NewCollectCompiler creates a new instance of CollectCompiler with the given compiler context.
func NewCollectCompiler(ctx *CompilationSession) *CollectCompiler {
	return &CollectCompiler{ctx: ctx}
}

func (c *CollectCompiler) bind(bindings *BindingCompiler, calls *CallResolver, exprs *ExprCompiler, recovery *RecoveryCompiler, facts *TypeFacts) {
	if c == nil {
		return
	}

	c.bindings = bindings
	c.calls = calls
	c.exprs = exprs
	c.recovery = recovery
	c.facts = facts
}

// Compile processes a COLLECT clause from the FQL AST and generates the appropriate VM instructions.
// It first compiles the collector specification and then compiles the loop operations based on that spec.
func (c *CollectCompiler) Compile(ctx fql.ICollectClauseContext) {
	scope := c.compileCollector(ctx)

	c.compileLoop(scope)
}

// compileCollector processes the COLLECT clause components and creates a Collector.
// This function handles the initialization of grouping, aggregation, and projection operations,
// and sets up the appropriate collector type based on the COLLECT clause structure.
func (c *CollectCompiler) compileCollector(ctx fql.ICollectClauseContext) *core.Collector {
	// Extract all components of the COLLECT clause
	groupingCtx := ctx.CollectGrouping()
	projectionCtx := ctx.CollectGroupProjection()
	counterCtx := ctx.CollectCounter()
	aggregationCtx := ctx.CollectAggregator()

	// We gather keys and values for the collector.
	kv, groupSelectors := c.initializeGrouping(groupingCtx)

	// Determine the collector type based on the presence of different COLLECT components
	collectorType := core.DetermineCollectorType(len(groupSelectors) > 0, aggregationCtx != nil, projectionCtx != nil, counterCtx != nil)
	useAggregateCollector := false
	aggregatePlanIndex := 0
	var globalAggregatePlan *bytecode.AggregatePlan
	var groupedAggregatePlan *bytecode.AggregatePlan
	useGroupedAggregateCollector := false

	// Fast path for global aggregation:
	// if every selector is plan-compatible, we can use the VM aggregate collector directly.
	if aggregationCtx != nil && len(groupSelectors) == 0 {
		if plan, ok := c.buildGlobalAggregatePlan(aggregationCtx); ok {
			collectorType = bytecode.CollectorTypeAggregate
			useAggregateCollector = true
			globalAggregatePlan = plan
			aggregatePlanIndex = c.ctx.AddAggregatePlan(plan)
		}
	}

	// Fast path for grouped aggregation:
	// enable fused grouped collector only for a narrow, predictable shape.
	if aggregationCtx != nil && groupingCtx != nil {
		selectors := aggregationCtx.AllCollectAggregateSelector()

		if c.shouldFuseGroupedAggregation(groupingCtx, selectors) {
			if plan, ok := c.buildGroupedAggregatePlan(selectors, projectionCtx != nil); ok {
				groupedAggregatePlan = plan
				useGroupedAggregateCollector = true
				collectorType = bytecode.CollectorTypeAggregateGroup
			}
		}
	}

	// We replace DataSet initialization with Collector initialization
	loop := c.ctx.Loops.Current()
	var dst bytecode.Operand

	// Aggregate collectors store their plan index in the collector opcode arguments.
	// Generic collectors do not require an aggregate plan operand.
	if useAggregateCollector {
		dst = loop.PatchDestinationAxy(c.ctx.Registers, c.ctx.Emitter, bytecode.OpDataSetCollector, int(collectorType), aggregatePlanIndex)
	} else if useGroupedAggregateCollector {
		planIdx := c.ctx.AddAggregatePlan(groupedAggregatePlan)
		dst = loop.PatchDestinationAxy(c.ctx.Registers, c.ctx.Emitter, bytecode.OpDataSetCollector, int(collectorType), planIdx)
	} else {
		dst = loop.PatchDestinationAx(c.ctx.Registers, c.ctx.Emitter, bytecode.OpDataSetCollector, int(collectorType))
	}

	var aggregation *core.CollectorAggregation

	// Initialize aggregationCtx if present in the COLLECT clause
	if aggregationCtx != nil {
		plan := globalAggregatePlan
		if len(groupSelectors) > 0 {
			plan = groupedAggregatePlan
		}

		aggregation = c.initializeAggregation(aggregationCtx, dst, kv, len(groupSelectors) > 0, plan)
	}

	// Initialize projectionCtx for group variables or counters
	projection := c.initializeProjection(kv, projectionCtx, counterCtx)
	projectionState := bytecode.NoopOperand

	if useAggregateCollector && projection != nil {
		projectionState = c.insertGlobalAggregateProjectionBuffer(loop)
	}

	// Create the collector specification with all components
	spec := core.NewCollector(collectorType, dst, projection, projectionState, groupSelectors, aggregation)

	// Finalize the collector setup
	c.finalizeCollector(dst, kv, spec)

	return spec
}

// finalizeCollector completes the collector setup by pushing key-value pairs to the collector
// and emitting finalization instructions for the current loop.
// The behavior varies based on whether grouping and aggregation are used.
func (c *CollectCompiler) finalizeCollector(dst bytecode.Operand, kv *core.KV, spec *core.Collector) {
	loop := c.ctx.Loops.Current()

	// Fused grouped aggregate collectors without INTO only need direct aggregate updates.
	if spec.HasGrouping() {
		if spec.Type() != bytecode.CollectorTypeAggregateGroup || spec.HasProjection() {
			c.ctx.Emitter.EmitPushKV(dst, kv.Key, kv.Value)
		}
	} else if spec.Type() == bytecode.CollectorTypeCounter {
		c.ctx.Emitter.EmitCounterInc(dst)
	} else if !spec.HasAggregation() {
		c.ctx.Emitter.EmitPushKV(dst, kv.Key, kv.Value)
	} else if spec.ProjectionState() != bytecode.NoopOperand {
		c.ctx.Emitter.EmitArrayPush(spec.ProjectionState(), kv.Value)
	} else if spec.HasProjection() {
		// For projection without grouping but with aggregation, use the projection variable name as the key
		key := c.facts.LoadConstant(runtime.String(spec.Projection().VariableName()))
		c.ctx.Emitter.EmitPushKV(dst, key, kv.Value)
	}

	// Emit finalization instructions for the current loop
	loop.EmitFinalization(c.ctx.Emitter)
}

func (c *CollectCompiler) insertGlobalAggregateProjectionBuffer(loop *core.Loop) bytecode.Operand {
	buf := c.ctx.Registers.Allocate()
	c.ctx.Emitter.InsertAx(loop.StartLabel(), bytecode.OpLoadArray, buf, 8)
	c.ctx.Types.Set(buf, core.TypeArray)

	return buf
}

// compileLoop compiles a loop construct by configuring its kind, registers, initialization, and processing based on the specification.
func (c *CollectCompiler) compileLoop(spec *core.Collector) {
	loop := c.ctx.Loops.Current()

	// If we are using a projection, we need to ensure the loop is set to ForInLoop
	if loop.Kind != core.ForInLoop {
		loop.Kind = core.ForInLoop
	}

	// Ensure loop value register is allocated
	if loop.Value == bytecode.NoopOperand {
		loop.Value = c.ctx.Registers.Allocate()
	}

	// Ensure loop key register is allocated
	if loop.Key == bytecode.NoopOperand {
		loop.Key = c.ctx.Registers.Allocate()
	}

	// Determine if we need to initialize the loop
	// We only need to initialize if we have grouping or if we don't have aggregation
	// (aggregate-only global mode is finalized in a synthetic loop later).
	doInit := spec.HasGrouping() || !spec.HasAggregation()

	if doInit {
		if loop.Allocate {
			// Move the collector to the next loop source
			c.ctx.Emitter.EmitMove(loop.Src, spec.Destination())
		} else {
			// We do not control the source of the loop, so we just set it to the destination
			loop.Src = spec.Destination()
		}

		loop.EmitInitialization(c.ctx.Registers, c.ctx.Emitter)

		if spec.HasProjection() && loop.Value.IsRegister() {
			if spec.Projection().IsCounted() {
				c.ctx.Types.Set(loop.Value, core.TypeInt)
			} else {
				c.ctx.Types.Set(loop.Value, core.TypeArray)
			}
		}
	}

	// Process aggregation if present
	if spec.HasAggregation() {
		c.finalizeAggregation(spec)
	}

	// Process grouping if present
	if spec.HasGrouping() {
		c.finalizeGrouping(spec)
	}

	// Process projection if present
	if spec.HasProjection() {
		// Global aggregate-only projection is finalized in compileGlobalAggregationFuncCalls.
		if spec.HasGrouping() || !spec.HasAggregation() {
			c.finalizeProjection(spec, loop.Value)
		}
	}
}
