package parser

import (
	"github.com/MontFerret/ferret/pkg/parser/fql"
	"github.com/antlr/antlr4/runtime/Go/antlr"
)

type Parser struct {
	tree *fql.FqlParser
}

func New(query string) *Parser {
	input := antlr.NewInputStream(query)
	lexer := fql.NewFqlLexer(input)
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
