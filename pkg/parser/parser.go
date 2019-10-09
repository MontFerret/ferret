//go:generate antlr -Xexact-output-dir -o fql -package fql -visitor -Dlanguage=Go antlr/FqlLexer.g4 antlr/FqlParser.g4
package parser

import (
	"github.com/MontFerret/ferret/pkg/parser/fql"
	resources "github.com/antlr/antlr4/doc/resources"
	"github.com/antlr/antlr4/runtime/Go/antlr"
)

type Parser struct {
	tree *fql.FqlParser
}

func New(query string) *Parser {
	input := antlr.NewInputStream(query)
	// converts tokens to upper case, so now it doesn’t matter
	// in which case the tokens were entered
	upper := resources.NewCaseChangingStream(input, true)

	lexer := fql.NewFqlLexer(upper)
	stream := antlr.NewCommonTokenStream(lexer, antlr.TokenDefaultChannel)

	p := fql.NewFqlParser(stream)
	p.BuildParseTrees = true
	p.AddErrorListener(antlr.NewDiagnosticErrorListener(true))

	return &Parser{tree: p}
}

func (p *Parser) AddErrorListener(listener antlr.ErrorListener) {
	p.tree.AddErrorListener(listener)
}

func (p *Parser) Visit(visitor fql.FqlParserVisitor) interface{} {
	return visitor.VisitProgram(p.tree.Program().(*fql.ProgramContext))
}
