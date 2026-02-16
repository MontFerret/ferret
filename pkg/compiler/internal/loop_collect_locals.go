package internal

import (
	"fmt"

	"github.com/antlr4-go/antlr/v4"

	"github.com/MontFerret/ferret/v2/pkg/compiler/internal/core"
	compilerdiagnostics "github.com/MontFerret/ferret/v2/pkg/compiler/internal/diagnostics"
	"github.com/MontFerret/ferret/v2/pkg/vm"
)

func (c *LoopCollectCompiler) declareLocalOrReport(ctx antlr.ParserRuleContext, name string, typ core.ValueType) vm.Operand {
	reg, ok := c.ctx.Symbols.DeclareLocal(name, typ)
	if ok {
		return reg
	}

	c.reportDuplicateLocal(ctx, name)

	// Keep bytecode emission valid after the diagnostic.
	if existing, _, found := c.ctx.Symbols.Resolve(name); found {
		return existing
	}

	return vm.NoopOperand
}

func (c *LoopCollectCompiler) assignLocalOrReport(ctx antlr.ParserRuleContext, name string, typ core.ValueType, op vm.Operand) bool {
	if c.ctx.Symbols.AssignLocal(name, typ, op) {
		return true
	}

	c.reportDuplicateLocal(ctx, name)
	return false
}

func (c *LoopCollectCompiler) reportDuplicateLocal(ctx antlr.ParserRuleContext, name string) {
	if ctx != nil {
		c.ctx.Errors.VariableNotUnique(ctx, name)
		return
	}

	c.ctx.Errors.Add(compilerdiagnostics.NewError(
		c.ctx.Source,
		compilerdiagnostics.NameError,
		fmt.Sprintf("Variable '%s' is already defined", name),
	))
}
