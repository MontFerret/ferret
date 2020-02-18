// Code generated from pkg/parser/antlr/FqlParser.g4 by ANTLR 4.7.1. DO NOT EDIT.

package fql // FqlParser
import "github.com/antlr/antlr4/runtime/Go/antlr"

type BaseFqlParserVisitor struct {
	*antlr.BaseParseTreeVisitor
}

func (v *BaseFqlParserVisitor) VisitProgram(ctx *ProgramContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseFqlParserVisitor) VisitUseExpression(ctx *UseExpressionContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseFqlParserVisitor) VisitUseExpressionAs(ctx *UseExpressionAsContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseFqlParserVisitor) VisitBody(ctx *BodyContext) interface{} {
	return v.VisitChildren(ctx)
}
