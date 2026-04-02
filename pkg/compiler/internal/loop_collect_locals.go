package internal

import (
	"fmt"

	"github.com/antlr4-go/antlr/v4"

	"github.com/MontFerret/ferret/v2/pkg/bytecode"
	"github.com/MontFerret/ferret/v2/pkg/compiler/internal/core"
	parser "github.com/MontFerret/ferret/v2/pkg/parser/diagnostics"
)

func (c *CollectCompiler) declareLocalOrReport(ctx antlr.ParserRuleContext, name string, typ core.ValueType) bytecode.Operand {
	reg, ok := c.ctx.Symbols.DeclareLocal(name, typ)
	if ok {
		return reg
	}

	c.reportDuplicateLocal(ctx, name)

	// Keep bytecode emission valid after the diagnostic.
	if existing, _, found := c.ctx.Symbols.Resolve(name); found {
		return existing
	}

	return bytecode.NoopOperand
}

func (c *CollectCompiler) assignLocalOrReport(ctx antlr.ParserRuleContext, name string, typ core.ValueType, op bytecode.Operand) bool {
	if c.ctx.Symbols.AssignLocal(name, typ, op) {
		return true
	}

	c.reportDuplicateLocal(ctx, name)
	return false
}

func (c *CollectCompiler) reportDuplicateLocal(ctx antlr.ParserRuleContext, name string) {
	if ctx != nil {
		c.ctx.Errors.VariableNotUnique(ctx, name)
		return
	}

	c.ctx.Errors.Add(parser.NewError(
		c.ctx.Source,
		parser.NameError,
		fmt.Sprintf("Variable '%s' is already defined", name),
	))
}
