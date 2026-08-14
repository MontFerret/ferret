package compiler

import (
	"strings"

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

		out = append(out, SyntaxToken{
			Kind: kind,
			Word: syntaxWord(token.GetTokenType(), token.GetText()),
			Span: span,
		})
	}

	return out
}

func syntaxWord(tokenType int, text string) SyntaxWord {
	switch tokenType {
	case fql.FqlLexerAggregate:
		return SyntaxWordAggregate
	case fql.FqlLexerAll:
		return SyntaxWordAll
	case fql.FqlLexerAnd:
		if strings.EqualFold(text, "AND") {
			return SyntaxWordAnd
		}
	case fql.FqlLexerAny:
		return SyntaxWordAny
	case fql.FqlLexerAs:
		return SyntaxWordAs
	case fql.FqlLexerAt:
		return SyntaxWordAt
	case fql.FqlLexerBackoff:
		return SyntaxWordBackoff
	case fql.FqlLexerCollect:
		return SyntaxWordCollect
	case fql.FqlLexerCount:
		return SyntaxWordCount
	case fql.FqlLexerDelete:
		return SyntaxWordDelete
	case fql.FqlLexerDispatch:
		return SyntaxWordDispatch
	case fql.FqlLexerDistinct:
		return SyntaxWordDistinct
	case fql.FqlLexerDo:
		return SyntaxWordDo
	case fql.FqlLexerEvent:
		return SyntaxWordEvent
	case fql.FqlLexerEvery:
		return SyntaxWordEvery
	case fql.FqlLexerExists:
		return SyntaxWordExists
	case fql.FqlLexerFilter:
		return SyntaxWordFilter
	case fql.FqlLexerFor:
		return SyntaxWordFor
	case fql.FqlLexerFunc:
		return SyntaxWordFunc
	case fql.FqlLexerIn:
		return SyntaxWordIn
	case fql.FqlLexerInto:
		return SyntaxWordInto
	case fql.FqlLexerJitter:
		return SyntaxWordJitter
	case fql.FqlLexerKeep:
		return SyntaxWordKeep
	case fql.FqlLexerLeast:
		return SyntaxWordLeast
	case fql.FqlLexerLet:
		return SyntaxWordLet
	case fql.FqlLexerLike:
		return SyntaxWordLike
	case fql.FqlLexerLimit:
		return SyntaxWordLimit
	case fql.FqlLexerMatch:
		return SyntaxWordMatch
	case fql.FqlLexerNone:
		return SyntaxWordNone
	case fql.FqlLexerNot:
		if strings.EqualFold(text, "NOT") {
			return SyntaxWordNot
		}
	case fql.FqlLexerNull:
		return SyntaxWordNull
	case fql.FqlLexerOne:
		return SyntaxWordOne
	case fql.FqlLexerOptions:
		return SyntaxWordOptions
	case fql.FqlLexerOr:
		if strings.EqualFold(text, "OR") {
			return SyntaxWordOr
		}
	case fql.FqlLexerQuery:
		return SyntaxWordQuery
	case fql.FqlLexerReturn:
		return SyntaxWordReturn
	case fql.FqlLexerSort:
		return SyntaxWordSort
	case fql.FqlLexerSortDirection:
		switch {
		case strings.EqualFold(text, "ASC"):
			return SyntaxWordAsc
		case strings.EqualFold(text, "DESC"):
			return SyntaxWordDesc
		}
	case fql.FqlLexerTimeout:
		return SyntaxWordTimeout
	case fql.FqlLexerTrigger:
		return SyntaxWordTrigger
	case fql.FqlLexerBooleanLiteral:
		switch {
		case strings.EqualFold(text, "TRUE"):
			return SyntaxWordTrue
		case strings.EqualFold(text, "FALSE"):
			return SyntaxWordFalse
		}
	case fql.FqlLexerUse:
		return SyntaxWordUse
	case fql.FqlLexerUsing:
		return SyntaxWordUsing
	case fql.FqlLexerValue:
		return SyntaxWordValue
	case fql.FqlLexerVar:
		return SyntaxWordVar
	case fql.FqlLexerWaitfor:
		return SyntaxWordWaitfor
	case fql.FqlLexerWhen:
		return SyntaxWordWhen
	case fql.FqlLexerWhile:
		return SyntaxWordWhile
	case fql.FqlLexerWith:
		return SyntaxWordWith
	case fql.FqlLexerIdentifier:
		switch {
		case strings.EqualFold(text, "ON"):
			return SyntaxWordOn
		case strings.EqualFold(text, "ERROR"):
			return SyntaxWordError
		case strings.EqualFold(text, "FAIL"):
			return SyntaxWordFail
		case strings.EqualFold(text, "RETRY"):
			return SyntaxWordRetry
		case strings.EqualFold(text, "DELAY"):
			return SyntaxWordDelay
		}
	}

	return SyntaxWordUnknown
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
