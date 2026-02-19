package internal

import (
	"github.com/antlr4-go/antlr/v4"

	"github.com/MontFerret/ferret/v2/pkg/parser/fql"
)

type programFormatter struct {
	*engine
}

func (f *programFormatter) bodyFirstElement(ctx *fql.BodyContext) antlr.ParserRuleContext {
	if ctx == nil {
		return nil
	}

	stmts := ctx.AllBodyStatement()

	if len(stmts) > 0 {
		return stmts[0].(antlr.ParserRuleContext)
	}

	if expr := ctx.BodyExpression(); expr != nil {
		return expr.(antlr.ParserRuleContext)
	}

	return nil
}

func (f *programFormatter) bodyLastElement(ctx *fql.BodyContext) antlr.ParserRuleContext {
	if ctx == nil {
		return nil
	}

	if expr := ctx.BodyExpression(); expr != nil {
		return expr.(antlr.ParserRuleContext)
	}

	stmts := ctx.AllBodyStatement()

	if len(stmts) > 0 {
		return stmts[len(stmts)-1].(antlr.ParserRuleContext)
	}

	return nil
}

func (f *programFormatter) formatProgram(ctx *fql.ProgramContext) {
	if ctx == nil {
		return
	}

	heads := ctx.AllHead()
	var first antlr.ParserRuleContext

	if len(heads) > 0 {
		first = heads[0].(antlr.ParserRuleContext)
	} else if body := ctx.Body(); body != nil {
		first = f.bodyFirstElement(body.(*fql.BodyContext))
	}

	f.trivia.emitLeading(first)

	for i, head := range heads {
		if i > 0 {
			f.trivia.emitBetween(heads[i-1].(antlr.ParserRuleContext), head.(antlr.ParserRuleContext))
		}

		f.formatHead(head.(*fql.HeadContext))
	}

	if ctx.Body() != nil {
		body := ctx.Body().(*fql.BodyContext)

		if len(heads) > 0 {
			f.trivia.emitBetween(heads[len(heads)-1].(antlr.ParserRuleContext), f.bodyFirstElement(body))
		}

		f.formatBody(body)
		last := f.bodyLastElement(body)
		if last == nil && len(heads) > 0 {
			last = heads[len(heads)-1].(antlr.ParserRuleContext)
		}

		f.trivia.emitTrailing(last)
	} else if len(heads) > 0 {
		f.trivia.emitTrailing(heads[len(heads)-1].(antlr.ParserRuleContext))
	}
}

func (f *programFormatter) formatHead(ctx *fql.HeadContext) {
	if ctx == nil {
		return
	}

	if useExpr := ctx.UseExpression(); useExpr != nil {
		f.formatUseExpression(useExpr.(*fql.UseExpressionContext))
	}
}

func (f *programFormatter) formatUseExpression(ctx *fql.UseExpressionContext) {
	if ctx == nil {
		return
	}

	if use := ctx.Use(); use != nil {
		f.formatUse(use.(*fql.UseContext))
	}
}

func (f *programFormatter) formatUse(ctx *fql.UseContext) {
	if ctx == nil {
		return
	}

	f.writeKeyword("USE")
	f.p.space()

	if ns := ctx.NamespaceIdentifier(); ns != nil {
		f.p.write(ns.GetText())
	}

	f.p.space()
	f.writeKeyword("AS")
	f.p.space()

	if alias := ctx.GetAlias(); alias != nil {
		f.p.write(alias.GetText())
	}
}

func (f *programFormatter) formatBody(ctx *fql.BodyContext) {
	if ctx == nil {
		return
	}

	stmts := ctx.AllBodyStatement()
	for i, stmt := range stmts {
		f.statement.formatBodyStatement(stmt.(*fql.BodyStatementContext))

		if i < len(stmts)-1 {
			f.trivia.emitBetween(stmt.(antlr.ParserRuleContext), stmts[i+1].(antlr.ParserRuleContext))
		}
	}

	if expr := ctx.BodyExpression(); expr != nil {
		if len(stmts) > 0 {
			f.trivia.emitBetween(stmts[len(stmts)-1].(antlr.ParserRuleContext), expr.(antlr.ParserRuleContext))
		}

		f.statement.formatBodyExpression(expr.(*fql.BodyExpressionContext))
	}
}
