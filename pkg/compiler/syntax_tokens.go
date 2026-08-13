package compiler

import (
	"github.com/antlr4-go/antlr/v4"

	"github.com/MontFerret/ferret/v2/pkg/parser/fql"
	"github.com/MontFerret/ferret/v2/pkg/source"
)

func buildSyntaxTokens(src *source.Source, tokens []antlr.Token) []SyntaxToken {
	byteOffsets := analysisByteOffsets(src)
	out := make([]SyntaxToken, 0, len(tokens))

	for _, token := range tokens {
		kind := syntaxTokenKind(token.GetTokenType())
		if kind == SyntaxTokenKindUnknown && token.GetTokenType() == antlr.TokenEOF {
			continue
		}

		if token.GetTokenType() == fql.FqlLexerWhiteSpaces || token.GetTokenType() == fql.FqlLexerLineTerminator {
			continue
		}

		span := source.Span{Start: token.GetStart(), End: token.GetStop() + 1}
		span = analysisByteSpan(byteOffsets, span)

		if span.Start < 0 || span.End <= span.Start {
			continue
		}

		out = append(out, SyntaxToken{Kind: kind, Span: span})
	}

	return out
}

func syntaxTokenKind(tokenType int) SyntaxTokenKind {
	switch tokenType {
	case fql.FqlLexerMultiLineComment, fql.FqlLexerSingleLineComment:
		return SyntaxTokenKindComment
	case fql.FqlLexerIdentifier, fql.FqlLexerIgnoreIdentifier:
		return SyntaxTokenKindIdentifier
	case fql.FqlLexerNamespaceSegment:
		return SyntaxTokenKindNamespace
	case fql.FqlLexerStringLiteral, fql.FqlLexerBacktickOpen, fql.FqlLexerTemplateChars, fql.FqlLexerBacktickClose:
		return SyntaxTokenKindString
	case fql.FqlLexerIntegerLiteral, fql.FqlLexerFloatLiteral:
		return SyntaxTokenKindNumber
	case fql.FqlLexerDurationLiteral:
		return SyntaxTokenKindDuration
	case fql.FqlLexerGt, fql.FqlLexerLt, fql.FqlLexerEq, fql.FqlLexerGte, fql.FqlLexerLte,
		fql.FqlLexerNeq, fql.FqlLexerMultiAssign, fql.FqlLexerDivAssign, fql.FqlLexerPlusAssign,
		fql.FqlLexerMinusAssign, fql.FqlLexerMulti, fql.FqlLexerDiv, fql.FqlLexerMod,
		fql.FqlLexerPlus, fql.FqlLexerMinus, fql.FqlLexerIncrement, fql.FqlLexerDecrement,
		fql.FqlLexerAnd, fql.FqlLexerOr, fql.FqlLexerRange, fql.FqlLexerArrow,
		fql.FqlLexerDispatchReceive, fql.FqlLexerAssign, fql.FqlLexerCoalesce,
		fql.FqlLexerQuestionMark, fql.FqlLexerRegexNotMatch, fql.FqlLexerRegexMatch,
		fql.FqlLexerTildeQuestion, fql.FqlLexerTilde, fql.FqlLexerEllipsis:
		return SyntaxTokenKindOperator
	case fql.FqlLexerMatch, fql.FqlLexerWhen, fql.FqlLexerFunc, fql.FqlLexerFor,
		fql.FqlLexerReturn, fql.FqlLexerQuery, fql.FqlLexerUsing, fql.FqlLexerWaitfor,
		fql.FqlLexerDispatch, fql.FqlLexerDelete, fql.FqlLexerOptions, fql.FqlLexerTrigger,
		fql.FqlLexerTimeout, fql.FqlLexerEvery, fql.FqlLexerBackoff, fql.FqlLexerJitter,
		fql.FqlLexerExists, fql.FqlLexerCount, fql.FqlLexerValue, fql.FqlLexerOne,
		fql.FqlLexerDistinct, fql.FqlLexerFilter, fql.FqlLexerSort, fql.FqlLexerLimit,
		fql.FqlLexerLet, fql.FqlLexerVar, fql.FqlLexerCollect, fql.FqlLexerSortDirection,
		fql.FqlLexerNone, fql.FqlLexerNull, fql.FqlLexerBooleanLiteral, fql.FqlLexerUse,
		fql.FqlLexerAs, fql.FqlLexerAt, fql.FqlLexerLeast, fql.FqlLexerInto,
		fql.FqlLexerKeep, fql.FqlLexerWith, fql.FqlLexerAll, fql.FqlLexerAny,
		fql.FqlLexerAggregate, fql.FqlLexerEvent, fql.FqlLexerLike, fql.FqlLexerNot,
		fql.FqlLexerIn, fql.FqlLexerDo, fql.FqlLexerWhile:
		return SyntaxTokenKindKeyword
	case fql.FqlLexerColon, fql.FqlLexerSemiColon, fql.FqlLexerDot, fql.FqlLexerComma,
		fql.FqlLexerOpenBracket, fql.FqlLexerCloseBracket, fql.FqlLexerOpenParen,
		fql.FqlLexerCloseParen, fql.FqlLexerOpenBrace, fql.FqlLexerCloseBrace,
		fql.FqlLexerParam, fql.FqlLexerTemplateExprStart, fql.FqlLexerTemplateExprEnd:
		return SyntaxTokenKindPunctuation
	default:
		return SyntaxTokenKindUnknown
	}
}
