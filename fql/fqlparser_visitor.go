// Code generated from pkg/parser/antlr/FqlParser.g4 by ANTLR 4.7.1. DO NOT EDIT.

package fql // FqlParser
import "github.com/antlr/antlr4/runtime/Go/antlr"

// A complete Visitor for a parse tree produced by FqlParser.
type FqlParserVisitor interface {
	antlr.ParseTreeVisitor

	// Visit a parse tree produced by FqlParser#program.
	VisitProgram(ctx *ProgramContext) interface{}

	// Visit a parse tree produced by FqlParser#useExpression.
	VisitUseExpression(ctx *UseExpressionContext) interface{}

	// Visit a parse tree produced by FqlParser#useExpressionAs.
	VisitUseExpressionAs(ctx *UseExpressionAsContext) interface{}

	// Visit a parse tree produced by FqlParser#body.
	VisitBody(ctx *BodyContext) interface{}
}
