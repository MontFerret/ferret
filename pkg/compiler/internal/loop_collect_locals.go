package internal

import (
	"fmt"

	"github.com/antlr4-go/antlr/v4"

	"github.com/MontFerret/ferret/v2/pkg/bytecode"
	"github.com/MontFerret/ferret/v2/pkg/compiler/internal/core"
	parser "github.com/MontFerret/ferret/v2/pkg/parser/diagnostics"
	"github.com/MontFerret/ferret/v2/pkg/parser/fql"
)

func (c *CollectCompiler) declareLocalOrReport(ctx antlr.ParserRuleContext, name string, typ core.ValueType) bytecode.Operand {
	bindingCtx := c.collectBindingContext(ctx)
	id := bindingIDFromRule(bindingCtx)
	reg, ok := c.ctx.Function.Symbols.DeclareLocalWithOptions(name, typ, core.BindingOptions{ID: id})

	if ok {
		if c.ctx.Program.Semantics != nil {
			c.recordCollectBinding(ctx, bindingCtx, id, name, typ)
		}

		return reg
	}

	c.reportDuplicateLocal(ctx, name)

	// Keep bytecode emission valid after the diagnostic.
	if existing, _, found := c.ctx.Function.Symbols.Resolve(name); found {
		return existing
	}

	return bytecode.NoopOperand
}

func (c *CollectCompiler) assignLocalOrReport(ctx antlr.ParserRuleContext, name string, typ core.ValueType, op bytecode.Operand) bool {
	bindingCtx := c.collectBindingContext(ctx)
	id := bindingIDFromRule(bindingCtx)

	if c.ctx.Function.Symbols.AssignLocalWithOptions(name, typ, op, core.BindingOptions{ID: id}) {
		if c.ctx.Program.Semantics != nil {
			c.recordCollectBinding(ctx, bindingCtx, id, name, typ)
		}

		return true
	}

	c.reportDuplicateLocal(ctx, name)
	return false
}

func (c *CollectCompiler) recordCollectBinding(
	declaration antlr.ParserRuleContext,
	selection antlr.ParserRuleContext,
	id core.BindingID,
	name string,
	typ core.ValueType,
) {
	if c == nil || c.ctx == nil || c.ctx.Program.Semantics == nil || declaration == nil || selection == nil {
		return
	}

	declarationSpan := parser.SpanFromRuleContext(declaration)
	c.ctx.Program.Semantics.RecordBinding(
		id,
		name,
		SemanticSymbolCollectBinding,
		declarationSpan,
		parser.SpanFromRuleContext(selection),
		false,
		typ,
		c.collectClauseActivation(declaration, declarationSpan.End),
		c.ctx.Program.Semantics.CurrentFunctionSymbol(),
		0,
	)
}

func (c *CollectCompiler) collectClauseActivation(ctx antlr.ParserRuleContext, fallback int) int {
	for node := antlr.Tree(ctx); node != nil; node = node.GetParent() {
		if clause, ok := node.(fql.ICollectClauseContext); ok {
			return parser.SpanFromRuleContext(clause.(antlr.ParserRuleContext)).End
		}
	}

	return fallback
}

func (c *CollectCompiler) collectBindingContext(ctx antlr.ParserRuleContext) antlr.ParserRuleContext {
	switch typed := ctx.(type) {
	case fql.ICollectSelectorContext:
		return c.bindingIdentifierRule(typed.BindingIdentifier())
	case fql.ICollectAggregateSelectorContext:
		return c.bindingIdentifierRule(typed.BindingIdentifier())
	case fql.ICollectGroupProjectionContext:
		if id := typed.BindingIdentifier(); id != nil {
			return c.bindingIdentifierRule(id)
		}

		if selector := typed.CollectSelector(); selector != nil {
			return c.bindingIdentifierRule(selector.BindingIdentifier())
		}
	case fql.ICollectCounterContext:
		return c.bindingIdentifierRule(typed.BindingIdentifier())
	}

	return ctx
}

func (c *CollectCompiler) bindingIdentifierRule(ctx fql.IBindingIdentifierContext) antlr.ParserRuleContext {
	if ctx == nil {
		return nil
	}

	return ctx.(antlr.ParserRuleContext)
}

func (c *CollectCompiler) reportDuplicateLocal(ctx antlr.ParserRuleContext, name string) {
	if ctx != nil {
		c.ctx.Program.Errors.VariableNotUnique(ctx, name)
		return
	}

	c.ctx.Program.Errors.Add(parser.NewError(
		c.ctx.Program.Source,
		parser.NameError,
		fmt.Sprintf("Variable '%s' is already defined", name),
	))
}
