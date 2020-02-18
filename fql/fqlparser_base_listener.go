// Code generated from pkg/parser/antlr/FqlParser.g4 by ANTLR 4.7.1. DO NOT EDIT.

package fql // FqlParser
import "github.com/antlr/antlr4/runtime/Go/antlr"

// BaseFqlParserListener is a complete listener for a parse tree produced by FqlParser.
type BaseFqlParserListener struct{}

var _ FqlParserListener = &BaseFqlParserListener{}

// VisitTerminal is called when a terminal node is visited.
func (s *BaseFqlParserListener) VisitTerminal(node antlr.TerminalNode) {}

// VisitErrorNode is called when an error node is visited.
func (s *BaseFqlParserListener) VisitErrorNode(node antlr.ErrorNode) {}

// EnterEveryRule is called when any rule is entered.
func (s *BaseFqlParserListener) EnterEveryRule(ctx antlr.ParserRuleContext) {}

// ExitEveryRule is called when any rule is exited.
func (s *BaseFqlParserListener) ExitEveryRule(ctx antlr.ParserRuleContext) {}

// EnterProgram is called when production program is entered.
func (s *BaseFqlParserListener) EnterProgram(ctx *ProgramContext) {}

// ExitProgram is called when production program is exited.
func (s *BaseFqlParserListener) ExitProgram(ctx *ProgramContext) {}

// EnterUseExpression is called when production useExpression is entered.
func (s *BaseFqlParserListener) EnterUseExpression(ctx *UseExpressionContext) {}

// ExitUseExpression is called when production useExpression is exited.
func (s *BaseFqlParserListener) ExitUseExpression(ctx *UseExpressionContext) {}

// EnterUseExpressionAs is called when production useExpressionAs is entered.
func (s *BaseFqlParserListener) EnterUseExpressionAs(ctx *UseExpressionAsContext) {}

// ExitUseExpressionAs is called when production useExpressionAs is exited.
func (s *BaseFqlParserListener) ExitUseExpressionAs(ctx *UseExpressionAsContext) {}

// EnterBody is called when production body is entered.
func (s *BaseFqlParserListener) EnterBody(ctx *BodyContext) {}

// ExitBody is called when production body is exited.
func (s *BaseFqlParserListener) ExitBody(ctx *BodyContext) {}
