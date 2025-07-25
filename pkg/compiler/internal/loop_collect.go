package internal

import (
	"github.com/MontFerret/ferret/pkg/compiler/internal/core"
	"github.com/MontFerret/ferret/pkg/parser/fql"
	"github.com/MontFerret/ferret/pkg/runtime"
	"github.com/MontFerret/ferret/pkg/vm"
)

type LoopCollectCompiler struct {
	ctx *CompilerContext
}

func NewCollectCompiler(ctx *CompilerContext) *LoopCollectCompiler {
	return &LoopCollectCompiler{ctx: ctx}
}

func (c *LoopCollectCompiler) Compile(ctx fql.ICollectClauseContext) {
	scope := c.compileCollector(ctx)

	c.compileLoop(scope)
}

func (c *LoopCollectCompiler) compileCollector(ctx fql.ICollectClauseContext) *core.CollectorSpec {
	grouping := ctx.CollectGrouping()
	projection := ctx.CollectGroupProjection()
	counter := ctx.CollectCounter()
	aggregation := ctx.CollectAggregator()

	// We gather keys and values for the collector.
	kv, groupSelectors := c.initializeGrouping(grouping)

	collectorType := core.DetermineCollectorType(len(groupSelectors) > 0, aggregation != nil, projection != nil, counter != nil)
	// We replace DataSet initialization with Collector initialization
	loop := c.ctx.Loops.Current()
	dst := loop.PatchDestinationAx(c.ctx.Registers, c.ctx.Emitter, vm.OpDataSetCollector, int(collectorType))

	var aggregationSelectors []*core.AggregateSelector

	// Fuse aggregation loop
	if aggregation != nil {
		aggregationSelectors = c.initializeAggregation(aggregation, dst, kv, len(groupSelectors) > 0)
	}

	groupProjection := c.initializeProjection(kv, projection, counter)

	spec := core.NewCollectorSpec(collectorType, dst, groupProjection, groupSelectors, aggregationSelectors)

	c.finalizeCollector(dst, kv, spec)

	// We no longer need KV, so we free registers
	c.ctx.Registers.Free(kv.Key)
	c.ctx.Registers.Free(kv.Value)

	return spec
}

func (c *LoopCollectCompiler) finalizeCollector(dst vm.Operand, kv *core.KV, spec *core.CollectorSpec) {
	loop := c.ctx.Loops.Current()

	// If we do not use grouping but use aggregation, we do not need to push the key and value
	// because they are already pushed by the global aggregation.
	if spec.HasGrouping() || !spec.HasAggregation() {
		c.ctx.Emitter.EmitPushKV(dst, kv.Key, kv.Value)
	} else if spec.HasProjection() {
		key := loadConstant(c.ctx, runtime.String(spec.Projection().VariableName()))
		c.ctx.Emitter.EmitPushKV(dst, key, kv.Value)
		c.ctx.Registers.Free(key)
	}

	loop.EmitFinalization(c.ctx.Emitter)
}

func (c *LoopCollectCompiler) compileLoop(spec *core.CollectorSpec) {
	loop := c.ctx.Loops.Current()

	// If we are using a projection, we need to ensure the loop is set to ForInLoop
	if loop.Kind != core.ForInLoop {
		loop.Kind = core.ForInLoop
	}

	if loop.Value == vm.NoopOperand {
		loop.Value = c.ctx.Registers.Allocate(core.Temp)
	}

	if loop.Key == vm.NoopOperand {
		loop.Key = c.ctx.Registers.Allocate(core.Temp)
	}

	doInit := spec.HasGrouping() || !spec.HasAggregation()

	if doInit {
		// Move the collector to the next loop source
		c.ctx.Emitter.EmitMove(loop.Src, spec.Destination())
		loop.EmitInitialization(c.ctx.Registers, c.ctx.Emitter, c.ctx.Loops.Depth())
	}

	if spec.HasAggregation() {
		c.unpackGroupedValues(spec)
		c.compileAggregation(spec)
	}

	if spec.HasGrouping() {
		c.compileGrouping(spec)
	}

	if spec.HasProjection() && !spec.HasAggregation() {
		c.finalizeProjection(spec, loop.Value)
	}
}
