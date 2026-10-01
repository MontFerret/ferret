package diagnostics

import (
	"strings"

	"github.com/antlr4-go/antlr/v4"

	"github.com/MontFerret/ferret/v2/pkg/parser/fql"
)

type (
	arrayEntryValidation struct {
		parser   *fql.FqlParser
		stream   *arrayEntryStream
		errors   *arrayEntryErrors
		recovery *arrayEntryRecovery
	}

	// A bounded cursor over already buffered tokens. The original parser's
	// cursor and token coordinates remain untouched, including during prediction.
	arrayEntryStream struct {
		*antlr.CommonTokenStream
		eof antlr.Token
		end int
	}

	arrayEntryErrors struct {
		*antlr.DefaultErrorListener
		failed bool
	}

	arrayEntryRecovery struct {
		*fql.BaseFqlParserListener
		rules      []arrayEntryRecoveryRule
		incomplete bool
	}

	arrayEntryRecoveryRule struct {
		context antlr.ParserRuleContext
		missing uint8
	}
)

const (
	arrayEntryRecoveryCondition uint8 = 1 << iota
	arrayEntryRecoveryAction
	arrayEntryRecoveryOperand
	arrayEntryRecoveryCount
	arrayEntryRecoveryDelay
	arrayEntryRecoveryBackoff
	arrayEntryRecoveryFallback
)

func newArrayEntryValidation(tokens *antlr.CommonTokenStream) *arrayEntryValidation {
	cursor := *tokens

	return &arrayEntryValidation{
		stream: &arrayEntryStream{CommonTokenStream: &cursor},
		errors: &arrayEntryErrors{DefaultErrorListener: antlr.NewDefaultErrorListener()},
	}
}

func (v *arrayEntryValidation) parse(start, end int) (*arrayEntryStream, bool) {
	v.stream.resetRange(start, end)
	v.errors.failed = false
	if v.parser == nil {
		v.parser = fql.NewFqlParser(v.stream)
		v.parser.BuildParseTrees = false
		v.parser.RemoveErrorListeners()
		v.parser.AddErrorListener(v.errors)
	} else {
		v.parser.SetTokenStream(v.stream)
	}

	// Some recovery productions intentionally accept partial tails so that the
	// compiler can diagnose them. Syntax acceptance alone is not completeness.
	checkRecovery := v.hasRecoveryTokens(start, end)
	if checkRecovery {
		if v.recovery == nil {
			v.recovery = &arrayEntryRecovery{BaseFqlParserListener: &fql.BaseFqlParserListener{}}
		}

		v.recovery.rules = v.recovery.rules[:0]
		v.recovery.incomplete = false
		v.parser.AddParseListener(v.recovery)
	}

	v.parser.ArrayEntry()
	if checkRecovery {
		v.parser.RemoveParseListener(v.recovery)
	}

	if v.errors.failed {
		// Failed rules may leave grammar-owned context flags unsettled. Only
		// reuse a parser after an entry parsed without any error recovery.
		v.parser = nil
	}

	return v.stream, !v.errors.failed && (!checkRecovery || !v.recovery.incomplete)
}

func (v *arrayEntryValidation) hasRecoveryTokens(start, end int) bool {
	for i := start; i < end; i++ {
		token := v.stream.CommonTokenStream.Get(i)
		if token.GetTokenType() == fql.FqlLexerIdentifier && strings.EqualFold(token.GetText(), "ON") {
			return true
		}
	}

	return false
}

func (s *arrayEntryStream) resetRange(start, end int) {
	boundary := s.CommonTokenStream.Get(end)
	eof := antlr.NewCommonToken(&antlr.TokenSourceCharStreamPair{}, antlr.TokenEOF, antlr.TokenDefaultChannel, boundary.GetStart(), boundary.GetStart()-1)
	eof.SetTokenIndex(end)
	eof.SetText("<EOF>")
	s.eof = eof
	s.end = end
	s.CommonTokenStream.Seek(start)
}

func (s *arrayEntryStream) LT(k int) antlr.Token {
	token := s.CommonTokenStream.LT(k)
	if token != nil && token.GetTokenIndex() >= s.end {
		return s.eof
	}

	return token
}

func (s *arrayEntryStream) LA(k int) int {
	token := s.LT(k)
	if token == nil {
		return 0
	}

	return token.GetTokenType()
}

func (s *arrayEntryStream) Consume() {
	if s.Index() < s.end {
		s.CommonTokenStream.Consume()
	}
}

func (s *arrayEntryStream) Seek(index int) {
	s.CommonTokenStream.Seek(min(index, s.end))
}

func (s *arrayEntryStream) Get(index int) antlr.Token {
	if index >= s.end {
		return s.eof
	}

	return s.CommonTokenStream.Get(index)
}

func (s *arrayEntryStream) Size() int {
	return s.end + 1
}

func (e *arrayEntryErrors) SyntaxError(antlr.Recognizer, interface{}, int, int, string, antlr.RecognitionException) {
	e.failed = true
}

func (r *arrayEntryRecovery) EnterEveryRule(ctx antlr.ParserRuleContext) {
	// Record direct grammar children without constructing an entry parse tree.
	if len(r.rules) > 0 {
		parent := &r.rules[len(r.rules)-1]
		if ctx.GetParent() == parent.context {
			switch ctx.GetRuleIndex() {
			case fql.FqlParserRULE_recoveryCondition:
				parent.missing &^= arrayEntryRecoveryCondition
			case fql.FqlParserRULE_recoveryAction:
				parent.missing &^= arrayEntryRecoveryAction
			case fql.FqlParserRULE_recoveryReturnExpr:
				parent.missing &^= arrayEntryRecoveryOperand
			case fql.FqlParserRULE_recoveryRetryCount:
				parent.missing &^= arrayEntryRecoveryCount
			case fql.FqlParserRULE_recoveryRetryDelayValue:
				parent.missing &^= arrayEntryRecoveryDelay
			case fql.FqlParserRULE_recoveryRetryBackoffKind:
				parent.missing &^= arrayEntryRecoveryBackoff
			case fql.FqlParserRULE_recoveryRetryFinalAction:
				parent.missing &^= arrayEntryRecoveryFallback
			}
		}
	}

	var missing uint8
	switch ctx.GetRuleIndex() {
	case fql.FqlParserRULE_recoveryTail:
		missing = arrayEntryRecoveryCondition | arrayEntryRecoveryAction
	case fql.FqlParserRULE_recoveryAction, fql.FqlParserRULE_recoveryRetryFinalAction:
		if ctx.GetStart().GetTokenType() == fql.FqlLexerReturn {
			missing = arrayEntryRecoveryOperand
		}
	case fql.FqlParserRULE_recoveryRetryAction:
		missing = arrayEntryRecoveryCount
	case fql.FqlParserRULE_recoveryRetryDelayClause:
		if ctx.GetStart().GetTokenType() != fql.FqlLexerBackoff {
			missing = arrayEntryRecoveryDelay
		}
	case fql.FqlParserRULE_recoveryRetryBackoffClause:
		missing = arrayEntryRecoveryBackoff
	case fql.FqlParserRULE_recoveryRetryOrClause, fql.FqlParserRULE_recoveryActionOrClause:
		missing = arrayEntryRecoveryFallback
	}

	if missing != 0 {
		r.rules = append(r.rules, arrayEntryRecoveryRule{context: ctx, missing: missing})
	}
}

func (r *arrayEntryRecovery) ExitEveryRule(ctx antlr.ParserRuleContext) {
	if len(r.rules) == 0 || r.rules[len(r.rules)-1].context != ctx {
		return
	}

	last := len(r.rules) - 1
	if r.rules[last].missing != 0 {
		r.incomplete = true
	}

	r.rules[last] = arrayEntryRecoveryRule{}
	r.rules = r.rules[:last]
}
