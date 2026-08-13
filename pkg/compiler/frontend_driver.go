package compiler

import (
	"github.com/antlr4-go/antlr/v4"

	"github.com/MontFerret/ferret/v2/pkg/compiler/internal"
	"github.com/MontFerret/ferret/v2/pkg/compiler/internal/optimization"
	"github.com/MontFerret/ferret/v2/pkg/parser"
	parserd "github.com/MontFerret/ferret/v2/pkg/parser/diagnostics"
	"github.com/MontFerret/ferret/v2/pkg/source"
)

func runFrontend(
	src *source.Source,
	errors *parserd.ErrorHandler,
	level optimization.Level,
	debugInfo bool,
	recorder *internal.SemanticRecorder,
) *Visitor {
	tokenHistory := parserd.NewTokenHistory(64)
	p := parser.New(src.Content(), func(stream antlr.TokenStream) antlr.TokenStream {
		return parserd.NewTrackingTokenStream(stream, tokenHistory)
	})

	p.RemoveErrorListeners()
	p.AddErrorListener(parserd.NewErrorListener(src, errors, tokenHistory))
	p.Program()

	if errors.HasErrors() {
		return nil
	}

	visitor := newVisitor(src, errors, level, recorder)
	visitor.Session.Program.DebugInfo = debugInfo
	p.Visit(visitor)

	return visitor
}
