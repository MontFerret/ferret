package diagnostics

import (
	"fmt"
	"strings"

	"github.com/antlr4-go/antlr/v4"

	"github.com/MontFerret/ferret/v2/pkg/diagnostics"
	"github.com/MontFerret/ferret/v2/pkg/parser/fql"
	"github.com/MontFerret/ferret/v2/pkg/source"
)

type (
	arrayLiteralAnalysis struct {
		tokens     *antlr.CommonTokenStream
		results    map[int]arrayLiteralProblem
		validation *arrayEntryValidation
	}

	arrayLiteralProblem struct {
		token     antlr.Token
		message   string
		hint      string
		label     string
		highlight bool
	}

	arrayLiteralFrame struct {
		commas []int
		open   int
	}
)

func newArrayLiteralAnalysis(stream antlr.TokenStream) *arrayLiteralAnalysis {
	if tracking, ok := stream.(*TrackingTokenStream); ok {
		stream = tracking.TokenStream
	}

	tokens, ok := stream.(*antlr.CommonTokenStream)
	if !ok {
		return nil
	}

	return &arrayLiteralAnalysis{tokens: tokens, results: make(map[int]arrayLiteralProblem)}
}

func (a *arrayLiteralAnalysis) match(src source.Source, diagnostic *diagnostics.Diagnostic, ctx antlr.ParserRuleContext) bool {
	if !isNoAlternative(diagnostic.Message) && !isMissing(diagnostic.Message) && !isMismatched(diagnostic.Message) {
		return false
	}

	// An expression atom beginning with '[' is a literal, even when prediction
	// fails before entering arrayLiteral. Indexes, keys and bindings do not
	// establish this context. Never infer a literal from an unmatched '[' alone.
	for ctx != nil {
		start := ctx.GetStart()
		if (ctx.GetRuleIndex() == fql.FqlParserRULE_expressionAtom || ctx.GetRuleIndex() == fql.FqlParserRULE_arrayLiteral) &&
			start != nil && start.GetTokenType() == fql.FqlLexerOpenBracket {
			open := start.GetTokenIndex()
			problem, found := a.results[open]
			if !found {
				a.tokens.Fill()
				problem = a.analyze(open)
				a.results[open] = problem
			}

			if problem.token == nil {
				return false
			}

			diagnostic.Message = problem.message
			diagnostic.Hint = problem.hint
			span := insertionSpanAfterToken(problem.token, src)
			if problem.highlight {
				span = spanFromTokenSafe(problem.token, src)
			}

			diagnostic.Spans = []diagnostics.ErrorSpan{diagnostics.NewMainErrorSpan(span, problem.label)}

			return true
		}

		ctx, _ = ctx.GetParent().(antlr.ParserRuleContext)
	}

	return false
}

func (a *arrayLiteralAnalysis) analyze(open int) arrayLiteralProblem {
	stack := []arrayLiteralFrame{{open: open}}
	end := a.tokens.Size() - 1
	last := open
	closed := false
	var mismatched antlr.Token
scan:
	for i := open + 1; i < a.tokens.Size(); i++ {
		token := a.tokens.Get(i)
		if token.GetChannel() != antlr.TokenDefaultChannel {
			continue
		}

		typeID := token.GetTokenType()
		if typeID == antlr.TokenEOF {
			break
		}

		// Preserve the existing missing-delimiter diagnosis before a following
		// statement, but do not truncate recovery actions or keyword-named calls.
		if len(stack) == 1 && (typeID == fql.FqlLexerReturn || typeID == fql.FqlLexerLet || typeID == fql.FqlLexerVar) &&
			a.tokens.Get(last).GetTokenType() != fql.FqlLexerDot &&
			(typeID != fql.FqlLexerReturn || !a.isRecoveryReturn(last)) &&
			a.next(i+1) < a.tokens.Size() && a.tokens.Get(a.next(i+1)).GetTokenType() != fql.FqlLexerOpenParen {
			if a.hasClosingBracket(i) {
				return arrayLiteralProblem{}
			}

			end = i

			break
		}

		switch typeID {
		case fql.FqlLexerOpenBracket, fql.FqlLexerOpenParen, fql.FqlLexerOpenBrace, fql.FqlLexerBacktickOpen, fql.FqlLexerTemplateExprStart:
			stack = append(stack, arrayLiteralFrame{open: i})
		case fql.FqlLexerCloseBracket, fql.FqlLexerCloseParen, fql.FqlLexerCloseBrace, fql.FqlLexerBacktickClose, fql.FqlLexerTemplateExprEnd:
			opener := a.tokens.Get(stack[len(stack)-1].open)
			matches := delimitersMatch(opener.GetText(), token.GetText()) ||
				(opener.GetTokenType() == fql.FqlLexerBacktickOpen && typeID == fql.FqlLexerBacktickClose) ||
				(opener.GetTokenType() == fql.FqlLexerTemplateExprStart && typeID == fql.FqlLexerTemplateExprEnd)
			if !matches {
				mismatched = token
				end = i

				break scan
			}

			if len(stack) == 1 {
				closed = true
				end = i

				break scan
			}

			stack = stack[:len(stack)-1]
		case fql.FqlLexerComma:
			frame := &stack[len(stack)-1]
			frame.commas = append(frame.commas, i)
		}

		last = i
	}

	// An unfinished inner object, call, grouping or template owns the error.
	for _, frame := range stack {
		if a.tokens.Get(frame.open).GetTokenType() != fql.FqlLexerOpenBracket {
			return arrayLiteralProblem{}
		}
	}

	for index, frame := range stack {
		start := a.next(frame.open + 1)
		for commaIndex, comma := range frame.commas {
			if start == comma {
				if commaIndex == 0 {
					return arrayLiteralProblem{
						token: a.tokens.Get(frame.open), message: "Expected a valid list of values",
						hint: "Did you forget to provide a value?", label: "missing value",
					}
				}

				return arrayLiteralProblem{
					token: a.tokens.Get(frame.commas[commaIndex-1]), message: "Expected expression after ','",
					hint: "Did you forget to provide a value?", label: "missing value",
				}
			}

			if problem, valid := a.entry(start, comma); !valid {
				return problem
			}

			start = a.next(comma + 1)
		}

		if index+1 < len(stack) {
			// Descend only when the unfinished array is itself the next entry.
			// In particular, 'value[' is indexing and an earlier malformed entry
			// must not disappear behind a later unmatched opener.
			if start != stack[index+1].open {
				return arrayLiteralProblem{}
			}

			continue
		}

		if start < end {
			if problem, valid := a.entry(start, end); !valid {
				return problem
			}
		}
	}

	if mismatched != nil {
		return arrayLiteralProblem{
			token: mismatched, message: fmt.Sprintf("Unexpected closing delimiter '%s'", mismatched.GetText()),
			hint: "Check that opening and closing delimiters match.", label: "unexpected closing delimiter", highlight: true,
		}
	}

	if closed {
		return arrayLiteralProblem{}
	}

	return arrayLiteralProblem{
		token: a.tokens.Get(last), message: "Unclosed array literal",
		hint: "Add ']' to close the array.", label: "expected ']'",
	}
}

func (a *arrayLiteralAnalysis) next(index int) int {
	for index < a.tokens.Size() && a.tokens.Get(index).GetChannel() != antlr.TokenDefaultChannel {
		index++
	}

	return index
}

func (a *arrayLiteralAnalysis) isRecoveryReturn(previous int) bool {
	typeID := a.tokens.Get(previous).GetTokenType()
	if typeID == fql.FqlLexerOr {
		// The entry grammar decides whether this is a retry fallback or an
		// incomplete ordinary expression; neither is a statement boundary.
		return true
	}

	if typeID != fql.FqlLexerIdentifier && typeID != fql.FqlLexerTimeout {
		return false
	}

	for previous--; previous >= 0; previous-- {
		token := a.tokens.Get(previous)
		if token.GetChannel() == antlr.TokenDefaultChannel {
			return token.GetTokenType() == fql.FqlLexerIdentifier && strings.EqualFold(token.GetText(), "ON")
		}
	}

	return false
}

// The suffix has not been visited by the structural scan yet. A later closer
// makes a proposed statement boundary ambiguous; do not truncate such an array.
func (a *arrayLiteralAnalysis) hasClosingBracket(start int) bool {
	depth := 1
	for i := start; i < a.tokens.Size(); i++ {
		switch a.tokens.Get(i).GetTokenType() {
		case fql.FqlLexerOpenBracket:
			depth++
		case fql.FqlLexerCloseBracket:
			depth--
			if depth == 0 {
				return true
			}
		}
	}

	return false
}

func (a *arrayLiteralAnalysis) entry(start, end int) (arrayLiteralProblem, bool) {
	stream, valid := a.parseEntry(start, end)
	if valid && stream.LA(1) == antlr.TokenEOF {
		return arrayLiteralProblem{}, true
	}

	if valid && isArrayItemStartToken(stream.LT(1)) {
		return arrayLiteralProblem{
			token: stream.LT(-1), message: missingArrayItemCommaMessage,
			hint: missingArrayItemCommaHint, label: "missing comma",
		}, false
	}

	last := stream.CommonTokenStream.Get(end - 1)
	for last.GetChannel() != antlr.TokenDefaultChannel && last.GetTokenIndex() > start {
		last = a.tokens.Get(last.GetTokenIndex() - 1)
	}

	if start == last.GetTokenIndex() && last.GetTokenType() == fql.FqlLexerEllipsis {
		return arrayLiteralProblem{
			token: last, message: "Expected expression after '...' in array spread",
			hint: "Provide an Array expression or none after '...'.", label: "missing expression",
		}, false
	}

	switch last.GetTokenType() {
	case fql.FqlLexerPlus, fql.FqlLexerMinus, fql.FqlLexerMulti, fql.FqlLexerDiv, fql.FqlLexerMod:
		// Verify the complete left operand before attributing failure to the
		// final operator; a plausible final token alone is not enough evidence.
		if last.GetTokenIndex() > start {
			left, complete := a.parseEntry(start, last.GetTokenIndex())
			if complete && left.LA(1) == antlr.TokenEOF {
				return arrayLiteralProblem{
					token: last, message: fmt.Sprintf("Expected right-hand expression after '%s'", last.GetText()),
					hint: "Provide an expression after the operator.", label: "missing expression",
				}, false
			}
		}
	}

	return arrayLiteralProblem{}, false
}

func (a *arrayLiteralAnalysis) parseEntry(start, end int) (*arrayEntryStream, bool) {
	if a.validation == nil {
		a.validation = newArrayEntryValidation(a.tokens)
	}

	return a.validation.parse(start, end)
}
