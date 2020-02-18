// Code generated from pkg/parser/antlr/FqlParser.g4 by ANTLR 4.7.1. DO NOT EDIT.

package fql // FqlParser
import "github.com/antlr/antlr4/runtime/Go/antlr"

// FqlParserListener is a complete listener for a parse tree produced by FqlParser.
type FqlParserListener interface {
	antlr.ParseTreeListener

	// EnterProgram is called when entering the program production.
	EnterProgram(c *ProgramContext)

	// EnterUseExpression is called when entering the useExpression production.
	EnterUseExpression(c *UseExpressionContext)

	// EnterUseExpressionAs is called when entering the useExpressionAs production.
	EnterUseExpressionAs(c *UseExpressionAsContext)

	// EnterBody is called when entering the body production.
	EnterBody(c *BodyContext)

	// ExitProgram is called when exiting the program production.
	ExitProgram(c *ProgramContext)

	// ExitUseExpression is called when exiting the useExpression production.
	ExitUseExpression(c *UseExpressionContext)

	// ExitUseExpressionAs is called when exiting the useExpressionAs production.
	ExitUseExpressionAs(c *UseExpressionAsContext)

	// ExitBody is called when exiting the body production.
	ExitBody(c *BodyContext)
}
