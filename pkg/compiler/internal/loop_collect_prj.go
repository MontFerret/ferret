package internal

import (
	"github.com/MontFerret/ferret/v2/pkg/bytecode"

	"github.com/MontFerret/ferret/v2/pkg/compiler/internal/core"
	parser "github.com/MontFerret/ferret/v2/pkg/parser/diagnostics"
	"github.com/MontFerret/ferret/v2/pkg/parser/fql"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

// initializeProjection handles the projection setup for group variables and counters.
// It processes either a group projection or a counter projection, depending on which is present.
// For group projections, it compiles the variable projection. For counters, it extracts the counter variable name.
// Returns a CollectorProjection object that encapsulates the projection information.
func (c *CollectCompiler) initializeProjection(kv *core.KV, projection fql.ICollectGroupProjectionContext, counter fql.ICollectCounterContext) *core.CollectorProjection {
	// Handle group variable projection
	if projection != nil {
		// Compile the group variable projection and get the variable name
		varName := c.compileGroupVariableProjection(kv, projection)
		// Create a group projection with the variable name
		return core.NewCollectorGroupProjection(varName, projection)
	}

	// Handle counter projection
	if counter != nil {
		// Extract the target variable after INTO
		varName := textOfBindingIdentifier(counter.BindingIdentifier())
		if varName == "" {
			err := c.ctx.Program.Errors.Create(parser.SemanticError, counter, "Missing counter projection variable")
			err.Hint = "Use WITH COUNT INTO <variable>."
			c.ctx.Program.Errors.Add(err)
			return nil
		}

		// Create a count projection with the variable name
		return core.NewCollectorCountProjection(varName, counter)
	}

	// If neither projection nor counter is present, return nil
	return nil
}

// finalizeProjection completes the projection setup by creating and assigning local variables.
// It handles different behaviors based on whether grouping and aggregation are used.
// Returns the register containing the projected value.
func (c *CollectCompiler) finalizeProjection(spec *core.Collector, aggregator bytecode.Operand) bytecode.Operand {
	loop := c.ctx.Function.Loops.Current()
	varName := spec.Projection().VariableName()

	if spec.HasGrouping() || !spec.HasAggregation() {
		// For cases with grouping or without aggregation:
		// Now we need to expand group variables from the dataset
		loop.ValueName = varName
		// Assign the aggregator value to the local variable with the projection name
		if !c.assignLocalOrReport(spec.Projection().Context(), loop.ValueName, core.TypeUnknown, aggregator) {
			if existing, found := c.ctx.Function.Symbols.ResolveBinding(loop.ValueName); found {
				if existing.Storage == core.BindingStorageCell {
					c.ctx.Program.Emitter.EmitStoreCell(existing.Register, ensureOperandRegister(c.ctx, c.facts, aggregator))
				} else {
					c.facts.EmitMoveAuto(existing.Register, aggregator)
				}
			}
		}

		return loop.Value
	}

	// For cases with aggregation but without grouping:
	// Load the value from the aggregator using the projection variable name as key
	key := c.facts.LoadConstant(runtime.String(varName))
	val := c.declareLocalOrReport(spec.Projection().Context(), varName, core.TypeUnknown)
	c.ctx.Program.Emitter.EmitABC(bytecode.OpLoadKey, val, aggregator, key)

	return val
}

// compileGroupVariableProjection processes group variable projections (both default and custom).
// It determines the type of projection (default with identifier or custom with selector)
// and delegates to the appropriate compilation method.
// Returns the variable name for the projection.
func (c *CollectCompiler) compileGroupVariableProjection(kv *core.KV, groupVar fql.ICollectGroupProjectionContext) string {
	// Handle default projection (identifier)
	if identifier := groupVar.BindingIdentifier(); identifier != nil {
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
func (c *CollectCompiler) compileDefaultGroupProjection(kv *core.KV, identifier fql.IBindingIdentifierContext, keeper fql.ICollectGroupProjectionFilterContext) string {
	if keeper == nil {
		// If no filter is provided, project all local variables
		variables := c.ctx.Function.Symbols.LocalVariables()
		// Create a scope projection with all local variables
		scope := core.NewScopeProjection(c.ctx.Function.Registers, c.ctx.Program.Emitter, c.ctx.Function.Symbols, c.ctx.Function.Types, variables)
		// Emit the scope as an object to the value register
		scope.EmitAsObject(kv.Value)
	} else {
		// If a filter is provided, project only the specified variables
		variables := keeper.AllIdentifier()
		resolved := make([]bytecode.Operand, len(variables))
		useTemp := false

		for i, variable := range variables {
			varName := variable.GetText()
			// Resolve the variable from the symbol table
			binding, found := c.ctx.Function.Symbols.ResolveBinding(varName)

			if !found {
				c.ctx.Program.Errors.VariableNotFound(variable.GetSymbol(), varName)
				noneReg := c.ctx.Function.Registers.Allocate()
				c.ctx.Program.Emitter.EmitA(bytecode.OpLoadNone, noneReg)
				c.ctx.Function.Types.Set(noneReg, core.TypeNone)
				resolved[i] = noneReg
				continue
			}

			resolved[i] = c.bindings.LoadBindingValue(binding)

			if binding.Register == kv.Value || resolved[i] == kv.Value {
				useTemp = true
			}
		}

		buildDst := kv.Value

		if useTemp {
			buildDst = c.ctx.Function.Registers.Allocate()
		}

		c.ctx.Program.Emitter.EmitObject(buildDst, len(variables))
		c.ctx.Function.Types.Set(buildDst, core.TypeObject)

		// Process each variable in the filter
		for i, variable := range variables {
			varName := variable.GetText()
			// Store the variable name as a string constant
			keyConst := c.ctx.Function.Symbols.AddConstant(runtime.String(varName))
			// Set the key-value pair in the object directly.
			// If kv.Value is referenced in the projection, buildDst is switched to a temp register.
			c.ctx.Program.Emitter.EmitObjectSetConst(buildDst, keyConst, resolved[i])
		}

		if buildDst != kv.Value {
			c.facts.EmitMoveAuto(kv.Value, buildDst)
		}
	}

	// Return the identifier text as the variable name
	return identifier.GetText()
}

// compileCustomGroupProjection handles custom group projection with a selector expression.
// It compiles the selector expression and moves its result to the value register.
// Returns the selector identifier text as the variable name for the projection.
func (c *CollectCompiler) compileCustomGroupProjection(kv *core.KV, selector fql.ICollectSelectorContext) string {
	// Compile the selector expression
	selectorReg := c.exprs.Compile(selector.Expression())
	// Move the result to the value register
	c.facts.EmitMoveAuto(kv.Value, selectorReg)

	// Return the selector identifier as the variable name
	return textOfBindingIdentifier(selector.BindingIdentifier())
}
