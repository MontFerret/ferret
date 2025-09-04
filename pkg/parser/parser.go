//go:generate antlr -Xexact-output-dir -o fql -package fql -visitor -Dlanguage=Go antlr/FqlLexer.g4 antlr/FqlParser.g4
package parser

import (
	antlr "github.com/antlr4-go/antlr/v4"
	"regexp"

	"github.com/MontFerret/ferret/pkg/parser/fql"
)

type Parser struct {
	tree *fql.FqlParser
}

// preprocessStepSyntax transforms STEP variable++ and STEP variable-- syntax
// into STEP variable = variable + 1 and STEP variable = variable - 1
func preprocessStepSyntax(query string) string {
	// Pattern to match STEP followed by identifier++ or identifier--
	// This pattern looks for:
	// - STEP keyword (case insensitive)
	// - whitespace
	// - identifier (letters, numbers, underscores)
	// - ++ or --
	incrementPattern := regexp.MustCompile(`(?i)\bSTEP\s+([a-zA-Z_][a-zA-Z0-9_]*)\s*\+\+`)
	decrementPattern := regexp.MustCompile(`(?i)\bSTEP\s+([a-zA-Z_][a-zA-Z0-9_]*)\s*--`)
	
	// Replace i++ with i = i + 1
	query = incrementPattern.ReplaceAllString(query, "STEP $1 = $1 + 1")
	
	// Replace i-- with i = i - 1
	query = decrementPattern.ReplaceAllString(query, "STEP $1 = $1 - 1")
	
	return query
}

func New(query string, tr ...TokenStreamTransformer) *Parser {
	// Preprocess the query to expand ++ and -- syntax
	query = preprocessStepSyntax(query)
	
	input := antlr.NewInputStream(query)
	// converts tokens to upper case, so now it doesn't matter
	// in which case the tokens were entered
	upper := newCaseChangingStream(input, true)
	lexer := fql.NewFqlLexer(upper)

	var stream antlr.TokenStream
	stream = antlr.NewCommonTokenStream(lexer, antlr.TokenDefaultChannel)

	// Apply all transformations to the token stream
	for _, transform := range tr {
		stream = transform(stream)
	}

	p := fql.NewFqlParser(stream)
	p.BuildParseTrees = true

	return &Parser{tree: p}
}

func (p *Parser) GetLiteralNames() []string {
	return p.tree.GetLiteralNames()[:]
}

func (p *Parser) AddErrorListener(listener antlr.ErrorListener) {
	p.tree.AddErrorListener(listener)
}

func (p *Parser) RemoveErrorListeners() {
	p.tree.RemoveErrorListeners()
}

func (p *Parser) Visit(visitor fql.FqlParserVisitor) interface{} {
	return visitor.VisitProgram(p.tree.Program().(*fql.ProgramContext))
}

func (p *Parser) Walk(listener fql.FqlParserListener) {
	antlr.ParseTreeWalkerDefault.Walk(listener, p.tree.Program())
}
