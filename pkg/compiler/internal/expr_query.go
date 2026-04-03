package internal

import (
	"github.com/MontFerret/ferret/v2/pkg/bytecode"
	"github.com/MontFerret/ferret/v2/pkg/parser/fql"
)

func (c *ExprCompiler) compileQueryExpression(ctx fql.IQueryExpressionContext) bytecode.Operand {
	return c.queryCompiler.compileQueryExpression(ctx)
}

func (c *ExprCompiler) compileQueryLiteral(ctx fql.IQueryLiteralContext) bytecode.Operand {
	return c.queryCompiler.compileQueryLiteral(ctx)
}
