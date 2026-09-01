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
	src source.Source,
	errors *parserd.ErrorHandler,
	level optimization.Level,
	debugInfo bool,
	recorder *internal.SemanticRecorder,
	syntaxTokens *[]SyntaxToken,
) *Visitor {
	tokenHistory := parserd.NewTokenHistory(64)
	var commonTokens *antlr.CommonTokenStream
	p := parser.New(src.Content(), func(stream antlr.TokenStream) antlr.TokenStream {
		if tokens, ok := stream.(*antlr.CommonTokenStream); ok {
			commonTokens = tokens
		}

		return parserd.NewTrackingTokenStream(stream, tokenHistory)
	})

	defer func() {
		if syntaxTokens == nil || commonTokens == nil {
			return
		}

		commonTokens.Fill()
		*syntaxTokens = buildSyntaxTokens(src, commonTokens.GetAllTokens())
	}()

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
