package internal

import (
	"github.com/MontFerret/ferret/pkg/compiler/internal/core"
	"github.com/MontFerret/ferret/pkg/parser/fql"
	"github.com/MontFerret/ferret/pkg/runtime"
	"github.com/MontFerret/ferret/pkg/vm"
	"github.com/antlr4-go/antlr/v4"
)

// initializeProjection handles the projection setup for group variables and counters.
// Returns the projection variable name and the appropriate collector type.
func (c *LoopCollectCompiler) initializeProjection(ctx fql.ICollectClauseContext, kv *core.KV, counter fql.ICollectCounterContext, hasGrouping bool) (string, core.CollectorType) {
	projectionVariableName := ""
	collectorType := core.CollectorTypeKey

	// Handle group variable projection
	if groupVar := ctx.CollectGroupVariable(); groupVar != nil {
		projectionVariableName = c.compileGroupVariableProjection(kv, groupVar)
		collectorType = core.CollectorTypeKeyGroup
		return projectionVariableName, collectorType
	}

	// Handle counter projection
	if counter != nil {
		projectionVariableName = counter.Identifier().GetText()
		collectorType = c.determineCounterCollectorType(hasGrouping)
	}

	return projectionVariableName, collectorType
}

// determineCounterCollectorType returns the appropriate collector type for counter operations.
func (c *LoopCollectCompiler) determineCounterCollectorType(hasGrouping bool) core.CollectorType {
	if hasGrouping {
		return core.CollectorTypeKeyCounter
	}

	return core.CollectorTypeCounter
}

// compileGroupVariableProjection processes group variable projections (both default and custom).
func (c *LoopCollectCompiler) compileGroupVariableProjection(kv *core.KV, groupVar fql.ICollectGroupVariableContext) string {
	// Handle default projection (identifier)
	if identifier := groupVar.Identifier(); identifier != nil {
		return c.compileDefaultGroupProjection(kv, identifier, groupVar.CollectGroupVariableKeeper())
	}

	// Handle custom projection (selector expression)
	if selector := groupVar.CollectSelector(); selector != nil {
		return c.compileCustomGroupProjection(kv, selector)
	}

	return ""
}

func (c *LoopCollectCompiler) compileDefaultGroupProjection(kv *core.KV, identifier antlr.TerminalNode, keeper fql.ICollectGroupVariableKeeperContext) string {
	if keeper == nil {
		variables := c.ctx.Symbols.LocalVariables()
		scope := core.NewScopeProjection(c.ctx.Registers, c.ctx.Emitter, c.ctx.Symbols, variables)
		scope.EmitAsObject(kv.Value)
	} else {
		variables := keeper.AllIdentifier()
		seq := c.ctx.Registers.AllocateSequence(len(variables) * 2)

		for i, j := 0, 0; i < len(variables); i, j = i+1, j+2 {
			varName := variables[i].GetText()
			loadConstantTo(c.ctx, runtime.String(varName), seq[j])

			variable, _, found := c.ctx.Symbols.Resolve(varName)

			if !found {
				panic("variable not found: " + varName)
			}

			c.ctx.Emitter.EmitAB(vm.OpMove, seq[j+1], variable)
		}

		c.ctx.Emitter.EmitAs(vm.OpLoadObject, kv.Value, seq)
		c.ctx.Registers.FreeSequence(seq)
	}

	return identifier.GetText()
}

func (c *LoopCollectCompiler) compileCustomGroupProjection(kv *core.KV, selector fql.ICollectSelectorContext) string {
	selectorReg := c.ctx.ExprCompiler.Compile(selector.Expression())
	c.ctx.Emitter.EmitMove(kv.Value, selectorReg)
	c.ctx.Registers.Free(selectorReg)

	return selector.Identifier().GetText()
}
