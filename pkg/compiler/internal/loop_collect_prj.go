package internal

import (
	"strings"

	"github.com/antlr4-go/antlr/v4"

	"github.com/MontFerret/ferret/pkg/compiler/internal/core"
	"github.com/MontFerret/ferret/pkg/parser/fql"
	"github.com/MontFerret/ferret/pkg/runtime"
	"github.com/MontFerret/ferret/pkg/vm"
)

// initializeProjection handles the projection setup for group variables and counters.
// It processes either a group projection or a counter projection, depending on which is present.
// For group projections, it compiles the variable projection. For counters, it extracts the counter variable name.
// Returns a CollectorProjection object that encapsulates the projection information.
func (c *LoopCollectCompiler) initializeProjection(kv *core.KV, projection fql.ICollectGroupProjectionContext, counter fql.ICollectCounterContext) *core.CollectorProjection {
	// Handle group variable projection
	if projection != nil {
		// Compile the group variable projection and get the variable name
		varName := c.compileGroupVariableProjection(kv, projection)
		// Create a group projection with the variable name
		return core.NewCollectorGroupProjection(varName)
	}

	// Handle counter projection
	if counter != nil {
		// Extract the counter variable name from the context
		// Extract the target variable (the second Identifier after INTO)
		varName := counter.Identifier(1).GetText()

		// Optional: validate that the first Identifier is actually "COUNT"
		if strings.ToUpper(counter.Identifier(0).GetText()) != "COUNT" {
			panic("counter identifier must be COUNT")
		}

		// Create a count projection with the variable name
		return core.NewCollectorCountProjection(varName)
	}

	// If neither projection nor counter is present, return nil
	return nil
}

// finalizeProjection completes the projection setup by creating and assigning local variables.
// It handles different behaviors based on whether grouping and aggregation are used.
// Returns the register containing the projected value.
func (c *LoopCollectCompiler) finalizeProjection(spec *core.Collector, aggregator vm.Operand) vm.Operand {
	loop := c.ctx.Loops.Current()
	varName := spec.Projection().VariableName()

	if spec.HasGrouping() || !spec.HasAggregation() {
		// For cases with grouping or without aggregation:
		// Now we need to expand group variables from the dataset
		loop.ValueName = varName
		// Assign the aggregator value to the local variable with the projection name
		c.ctx.Symbols.AssignLocal(loop.ValueName, core.TypeUnknown, aggregator)

		return loop.Value
	}

	// For cases with aggregation but without grouping:
	// Load the value from the aggregator using the projection variable name as key
	key := loadConstant(c.ctx, runtime.String(varName))
	val := c.ctx.Symbols.DeclareLocal(varName, core.TypeUnknown)
	c.ctx.Emitter.EmitABC(vm.OpLoadKey, val, aggregator, key)
	c.ctx.Registers.Free(key)

	return val
}

// compileGroupVariableProjection processes group variable projections (both default and custom).
// It determines the type of projection (default with identifier or custom with selector)
// and delegates to the appropriate compilation method.
// Returns the variable name for the projection.
func (c *LoopCollectCompiler) compileGroupVariableProjection(kv *core.KV, groupVar fql.ICollectGroupProjectionContext) string {
	// Handle default projection (identifier)
	if identifier := groupVar.Identifier(); identifier != nil {
		// Default projection uses an identifier and optional filter
		return c.compileDefaultGroupProjection(kv, identifier, groupVar.CollectGroupProjectionFilter())
	}

	// Handle custom projection (selector expression)
	if selector := groupVar.CollectSelector(); selector != nil {
		// Custom projection uses a selector expression
		return c.compileCustomGroupProjection(kv, selector)
	}

	// Return empty string if neither type of projection is present
	return ""
}

// compileDefaultGroupProjection handles the default group projection with an identifier.
// It can either project all local variables (when keeper is nil) or only specific variables (when keeper is provided).
// Returns the identifier text as the variable name for the projection.
func (c *LoopCollectCompiler) compileDefaultGroupProjection(kv *core.KV, identifier antlr.TerminalNode, keeper fql.ICollectGroupProjectionFilterContext) string {
	if keeper == nil {
		// If no filter is provided, project all local variables
		variables := c.ctx.Symbols.LocalVariables()
		// Create a scope projection with all local variables
		scope := core.NewScopeProjection(c.ctx.Registers, c.ctx.Emitter, c.ctx.Symbols, variables)
		// Emit the scope as an object to the value register
		scope.EmitAsObject(kv.Value)
	} else {
		// If a filter is provided, project only the specified variables
		variables := keeper.AllIdentifier()
		// Allocate a sequence of registers for key-value pairs (hence *2)
		seq := c.ctx.Registers.AllocateSequence(len(variables) * 2)

		// Process each variable in the filter
		for i, j := 0, 0; i < len(variables); i, j = i+1, j+2 {
			varName := variables[i].GetText()
			// Load the variable name as a string constant to the key register
			loadConstantTo(c.ctx, runtime.String(varName), seq[j])

			// Resolve the variable from the symbol table
			variable, _, found := c.ctx.Symbols.Resolve(varName)

			if !found {
				panic("variable not found: " + varName)
			}

			// Move the variable value to the value register
			c.ctx.Emitter.EmitAB(vm.OpMove, seq[j+1], variable)
		}

		// Create an object from the key-value pairs
		c.ctx.Emitter.EmitAs(vm.OpLoadObject, kv.Value, seq)
		// Free the sequence registers
		c.ctx.Registers.FreeSequence(seq)
	}

	// Return the identifier text as the variable name
	return identifier.GetText()
}

// compileCustomGroupProjection handles custom group projection with a selector expression.
// It compiles the selector expression and moves its result to the value register.
// Returns the selector identifier text as the variable name for the projection.
func (c *LoopCollectCompiler) compileCustomGroupProjection(kv *core.KV, selector fql.ICollectSelectorContext) string {
	// Compile the selector expression
	selectorReg := c.ctx.ExprCompiler.Compile(selector.Expression())
	// Move the result to the value register
	c.ctx.Emitter.EmitMove(kv.Value, selectorReg)
	// Free the temporary register
	c.ctx.Registers.Free(selectorReg)

	// Return the selector identifier as the variable name
	return selector.Identifier().GetText()
}
