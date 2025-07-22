package internal

import (
	"github.com/MontFerret/ferret/pkg/compiler/internal/core"
	"github.com/MontFerret/ferret/pkg/parser/fql"
	"github.com/MontFerret/ferret/pkg/vm"
)

type (
	LoopCollectCompiler struct {
		ctx *CompilerContext
	}

	collectorScope struct {
		Type                 core.CollectorType
		Projection           string
		GroupSelectors       []fql.ICollectSelectorContext
		AggregationSelectors []*aggregateSelector
	}
)

func NewCollectCompiler(ctx *CompilerContext) *LoopCollectCompiler {
	return &LoopCollectCompiler{ctx: ctx}
}

func (c *LoopCollectCompiler) Compile(ctx fql.ICollectClauseContext) {
	scope := c.compileCollector(ctx)

	c.compileLoop(scope)
}

func (c *LoopCollectCompiler) compileCollector(ctx fql.ICollectClauseContext) *collectorScope {
	grouping := ctx.CollectGrouping()
	counter := ctx.CollectCounter()
	aggregation := ctx.CollectAggregator()

	// We gather keys and values for the collector.
	kv, groupSelectors := c.initializeGrouping(grouping)
	projectionVarName, collectorType := c.initializeProjection(ctx, kv, counter, grouping != nil)

	// If we use aggregators, we need to collect group items by key
	if aggregation != nil && collectorType != core.CollectorTypeKeyGroup {
		// We need to patch the loop result to be a collector
		collectorType = core.CollectorTypeKeyGroup
	}

	loop := c.ctx.Loops.Current()
	// We replace DataSet initialization with Collector initialization
	dst := loop.PatchDestinationAx(c.ctx.Registers, c.ctx.Emitter, vm.OpDataSetCollector, int(collectorType))

	var aggregationSelectors []*aggregateSelector

	// Fuse aggregation loop
	if aggregation != nil {
		aggregationSelectors = c.initializeAggregation(aggregation, dst, kv, len(aggregationSelectors) > 0)
	}

	c.finalizeCollector(dst, kv, len(groupSelectors) > 0, aggregation != nil)

	// We no longer need KV, so we free registers
	c.ctx.Registers.Free(kv.Key)
	c.ctx.Registers.Free(kv.Value)

	return &collectorScope{collectorType, projectionVarName, groupSelectors, aggregationSelectors}
}

func (c *LoopCollectCompiler) finalizeCollector(dst vm.Operand, kv *core.KV, withGrouping bool, withAggregation bool) {
	loop := c.ctx.Loops.Current()

	// If we do not use grouping but use aggregation, we do not need to push the key and value
	// because they are already pushed by the global aggregation.
	push := withGrouping || !withAggregation

	if push {
		c.ctx.Emitter.EmitABC(vm.OpPushKV, dst, kv.Key, kv.Value)
	}

	loop.EmitFinalization(c.ctx.Emitter)

	// Move the collector to the next loop source
	c.ctx.Emitter.EmitMove(loop.Src, dst)
}

func (c *LoopCollectCompiler) compileLoop(scope *collectorScope) {
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

	withGrouping := len(scope.GroupSelectors) > 0
	withAggregation := len(scope.AggregationSelectors) > 0
	doInit := withGrouping || !withAggregation

	if doInit {
		loop.EmitInitialization(c.ctx.Registers, c.ctx.Emitter, c.ctx.Loops.Depth())
	}

	if withAggregation {
		c.unpackGroupedValues(scope.AggregationSelectors, withGrouping)
		c.compileAggregation(scope.AggregationSelectors, withGrouping)
	}

	// If the projection is used, we allocate a new register for the variable and put the iterator's value into it
	if scope.Projection != "" {
		// Now we need to expand group variables from the dataset
		loop.ValueName = scope.Projection
		c.ctx.Symbols.AssignLocal(loop.ValueName, core.TypeUnknown, loop.Value)
	}

	if withGrouping {
		c.compileGrouping(scope.Type, scope.GroupSelectors)
	}
}
