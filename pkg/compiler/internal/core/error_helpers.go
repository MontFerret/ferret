package core

import (
	"github.com/antlr4-go/antlr/v4"

	"github.com/MontFerret/ferret/pkg/file"
)

func LocationFromRuleContext(ctx antlr.ParserRuleContext) *file.Location {
	start := ctx.GetStart()
	stop := ctx.GetStop()

	// Defensive: avoid nil dereference
	if start == nil || stop == nil {
		return file.EmptyLocation()
	}

	return file.NewLocation(
		start.GetLine(),
		start.GetColumn()+1,
		start.GetStart(),
		stop.GetStop(),
	)
}

func LocationFromToken(token antlr.Token) *file.Location {
	if token == nil {
		return file.EmptyLocation()
	}

	return file.NewLocation(
		token.GetLine(),
		token.GetColumn()+1,
		token.GetStart(),
		token.GetStop(),
	)
}
