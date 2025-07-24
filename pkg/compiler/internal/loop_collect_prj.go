package internal

import (
	"github.com/antlr4-go/antlr/v4"

	"github.com/MontFerret/ferret/pkg/compiler/internal/core"
	"github.com/MontFerret/ferret/pkg/parser/fql"
	"github.com/MontFerret/ferret/pkg/runtime"
	"github.com/MontFerret/ferret/pkg/vm"
)

// initializeProjection handles the projection setup for group variables and counters.
// Returns the projection variable name and the appropriate collector type.
func (c *LoopCollectCompiler) initializeProjection(ctx fql.ICollectClauseContext, kv *core.KV, counter fql.ICollectCounterContext) *core.CollectorProjection {
	// Handle group variable projection
	if groupVar := ctx.CollectGroupProjection(); groupVar != nil {
		varName := c.compileGroupVariableProjection(kv, groupVar)
		return core.NewCollectorGroupProjection(varName)
	}

	// Handle counter projection
	if counter != nil {
		varName := counter.Identifier().GetText()

		return core.NewCollectorCountProjection(varName)
	}

	return nil
}

func (c *LoopCollectCompiler) finalizeProjection(spec *core.CollectorSpec) {
	loop := c.ctx.Loops.Current()
	varName := spec.Projection().VariableName()

	if spec.HasGrouping() || !spec.HasAggregation() {
		// Now we need to expand group variables from the dataset
		loop.ValueName = varName
		c.ctx.Symbols.AssignLocal(loop.ValueName, core.TypeUnknown, loop.Value)
	} else {
		key := loadConstant(c.ctx, runtime.String(varName))
		val := c.ctx.Symbols.DeclareLocal(varName, core.TypeUnknown)
		c.ctx.Emitter.EmitABC(vm.OpLoadKey, val, loop.Dst, key)
		c.ctx.Registers.Free(key)
	}
}

// compileGroupVariableProjection processes group variable projections (both default and custom).
func (c *LoopCollectCompiler) compileGroupVariableProjection(kv *core.KV, groupVar fql.ICollectGroupProjectionContext) string {
	// Handle default projection (identifier)
	if identifier := groupVar.Identifier(); identifier != nil {
		return c.compileDefaultGroupProjection(kv, identifier, groupVar.CollectGroupProjectionFilter())
	}

	// Handle custom projection (selector expression)
	if selector := groupVar.CollectSelector(); selector != nil {
		return c.compileCustomGroupProjection(kv, selector)
	}

	return ""
}

func (c *LoopCollectCompiler) compileDefaultGroupProjection(kv *core.KV, identifier antlr.TerminalNode, keeper fql.ICollectGroupProjectionFilterContext) string {
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
