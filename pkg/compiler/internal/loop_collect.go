package internal

import (
	"github.com/MontFerret/ferret/pkg/compiler/internal/core"
	"github.com/MontFerret/ferret/pkg/parser/fql"
	"github.com/MontFerret/ferret/pkg/runtime"
	"github.com/MontFerret/ferret/pkg/vm"
)

// LoopCollectCompiler handles the compilation of COLLECT clauses in FQL queries.
// It transforms COLLECT operations into VM instructions for data aggregation and grouping.
type LoopCollectCompiler struct {
	ctx *CompilerContext
}

// NewCollectCompiler creates a new instance of LoopCollectCompiler with the given compiler context.
func NewCollectCompiler(ctx *CompilerContext) *LoopCollectCompiler {
	return &LoopCollectCompiler{ctx: ctx}
}

// Compile processes a COLLECT clause from the FQL AST and generates the appropriate VM instructions.
// It first compiles the collector specification and then compiles the loop operations based on that spec.
func (c *LoopCollectCompiler) Compile(ctx fql.ICollectClauseContext) {
	scope := c.compileCollector(ctx)

	c.compileLoop(scope)
}

// compileCollector processes the COLLECT clause components and creates a CollectorSpec.
// This function handles the initialization of grouping, aggregation, and projection operations,
// and sets up the appropriate collector type based on the COLLECT clause structure.
func (c *LoopCollectCompiler) compileCollector(ctx fql.ICollectClauseContext) *core.CollectorSpec {
	// Extract all components of the COLLECT clause
	grouping := ctx.CollectGrouping()
	projection := ctx.CollectGroupProjection()
	counter := ctx.CollectCounter()
	aggregation := ctx.CollectAggregator()

	// We gather keys and values for the collector.
	kv, groupSelectors := c.initializeGrouping(grouping)

	// Determine the collector type based on the presence of different COLLECT components
	collectorType := core.DetermineCollectorType(len(groupSelectors) > 0, aggregation != nil, projection != nil, counter != nil)

	// We replace DataSet initialization with Collector initialization
	loop := c.ctx.Loops.Current()
	dst := loop.PatchDestinationAx(c.ctx.Registers, c.ctx.Emitter, vm.OpDataSetCollector, int(collectorType))

	var aggregationSelectors []*core.AggregateSelector

	// Initialize aggregation if present in the COLLECT clause
	if aggregation != nil {
		aggregationSelectors = c.initializeAggregation(aggregation, dst, kv, len(groupSelectors) > 0)
	}

	// Initialize projection for group variables or counters
	groupProjection := c.initializeProjection(kv, projection, counter)

	// Create the collector specification with all components
	spec := core.NewCollectorSpec(collectorType, dst, groupProjection, groupSelectors, aggregationSelectors)

	// Finalize the collector setup
	c.finalizeCollector(dst, kv, spec)

	// We no longer need KV, so we free registers
	c.ctx.Registers.Free(kv.Key)
	c.ctx.Registers.Free(kv.Value)

	return spec
}

// finalizeCollector completes the collector setup by pushing key-value pairs to the collector
// and emitting finalization instructions for the current loop.
// The behavior varies based on whether grouping and aggregation are used.
func (c *LoopCollectCompiler) finalizeCollector(dst vm.Operand, kv *core.KV, spec *core.CollectorSpec) {
	loop := c.ctx.Loops.Current()

	// If we do not use grouping but use aggregation, we do not need to push the key and value
	// because they are already pushed by the global aggregation.
	if spec.HasGrouping() || !spec.HasAggregation() {
		c.ctx.Emitter.EmitPushKV(dst, kv.Key, kv.Value)
	} else if spec.HasProjection() {
		// For projection without grouping but with aggregation, use the projection variable name as the key
		key := loadConstant(c.ctx, runtime.String(spec.Projection().VariableName()))
		c.ctx.Emitter.EmitPushKV(dst, key, kv.Value)
		c.ctx.Registers.Free(key)
	}

	// Emit finalization instructions for the current loop
	loop.EmitFinalization(c.ctx.Emitter)
}

// compileLoop processes the loop operations based on the collector specification.
// It handles different combinations of grouping, aggregation, and projection operations,
// ensuring that the appropriate VM instructions are generated for each case.
func (c *LoopCollectCompiler) compileLoop(spec *core.CollectorSpec) {
	loop := c.ctx.Loops.Current()

	// If we are using a projection, we need to ensure the loop is set to ForInLoop
	if loop.Kind != core.ForInLoop {
		loop.Kind = core.ForInLoop
	}

	// Ensure loop value register is allocated
	if loop.Value == vm.NoopOperand {
		loop.Value = c.ctx.Registers.Allocate(core.Temp)
	}

	// Ensure loop key register is allocated
	if loop.Key == vm.NoopOperand {
		loop.Key = c.ctx.Registers.Allocate(core.Temp)
	}

	// Determine if we need to initialize the loop
	// We only need to initialize if we have grouping or if we don't have aggregation
	doInit := spec.HasGrouping() || !spec.HasAggregation()

	if doInit {
		// Move the collector to the next loop source
		c.ctx.Emitter.EmitMove(loop.Src, spec.Destination())
		loop.EmitInitialization(c.ctx.Registers, c.ctx.Emitter, c.ctx.Loops.Depth())
	}

	// Process aggregation if present
	if spec.HasAggregation() {
		c.unpackGroupedValues(spec)
		c.compileAggregation(spec)
	}

	// Process grouping if present
	if spec.HasGrouping() {
		c.compileGrouping(spec)
	}

	// We finalize projection only if we have a projection and no aggregation
	// Because if we have aggregation, we finalize it in the compileAggregation method.
	if spec.HasProjection() && !spec.HasAggregation() {
		c.finalizeProjection(spec, loop.Value)
	}
}
