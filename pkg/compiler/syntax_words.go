package compiler

type (
	// SyntaxWord identifies a canonical FQL language word independently of
	// generated lexer token numbers. Zero means the token is not a language word.
	SyntaxWord uint8

	// SyntaxWordCategory classifies the role of a canonical FQL language word.
	SyntaxWordCategory uint8

	// SyntaxWordInfo describes canonical FQL language-word metadata.
	SyntaxWordInfo struct {
		Spelling string
		Word     SyntaxWord
		Category SyntaxWordCategory
	}
)

const (
	SyntaxWordUnknown SyntaxWord = iota
	SyntaxWordAggregate
	SyntaxWordAll
	SyntaxWordAnd
	SyntaxWordAny
	SyntaxWordAs
	SyntaxWordAsc
	SyntaxWordAt
	SyntaxWordBackoff
	SyntaxWordCollect
	SyntaxWordCount
	SyntaxWordDelay
	SyntaxWordDelete
	SyntaxWordDesc
	SyntaxWordDispatch
	SyntaxWordDistinct
	SyntaxWordDo
	SyntaxWordError
	SyntaxWordEvent
	SyntaxWordEvery
	SyntaxWordExists
	SyntaxWordFail
	SyntaxWordFalse
	SyntaxWordFilter
	SyntaxWordFor
	SyntaxWordFunc
	SyntaxWordIn
	SyntaxWordInto
	SyntaxWordJitter
	SyntaxWordKeep
	SyntaxWordLeast
	SyntaxWordLet
	SyntaxWordLike
	SyntaxWordLimit
	SyntaxWordMatch
	SyntaxWordNone
	SyntaxWordNot
	SyntaxWordNull
	SyntaxWordOn
	SyntaxWordOne
	SyntaxWordOptions
	SyntaxWordOr
	SyntaxWordQuery
	SyntaxWordRetry
	SyntaxWordReturn
	SyntaxWordSort
	SyntaxWordTimeout
	SyntaxWordTrigger
	SyntaxWordTrue
	SyntaxWordUse
	SyntaxWordUsing
	SyntaxWordValue
	SyntaxWordVar
	SyntaxWordWaitfor
	SyntaxWordWhen
	SyntaxWordWhile
	SyntaxWordWith
)

const (
	SyntaxWordCategoryUnknown SyntaxWordCategory = iota
	SyntaxWordCategoryKeyword
	SyntaxWordCategoryOperator
	SyntaxWordCategoryLiteral
	SyntaxWordCategoryContextual
)

var syntaxWords = []SyntaxWordInfo{
	{Word: SyntaxWordAggregate, Spelling: "AGGREGATE", Category: SyntaxWordCategoryKeyword},
	{Word: SyntaxWordAll, Spelling: "ALL", Category: SyntaxWordCategoryKeyword},
	{Word: SyntaxWordAnd, Spelling: "AND", Category: SyntaxWordCategoryOperator},
	{Word: SyntaxWordAny, Spelling: "ANY", Category: SyntaxWordCategoryKeyword},
	{Word: SyntaxWordAs, Spelling: "AS", Category: SyntaxWordCategoryKeyword},
	{Word: SyntaxWordAsc, Spelling: "ASC", Category: SyntaxWordCategoryKeyword},
	{Word: SyntaxWordAt, Spelling: "AT", Category: SyntaxWordCategoryKeyword},
	{Word: SyntaxWordBackoff, Spelling: "BACKOFF", Category: SyntaxWordCategoryKeyword},
	{Word: SyntaxWordCollect, Spelling: "COLLECT", Category: SyntaxWordCategoryKeyword},
	{Word: SyntaxWordCount, Spelling: "COUNT", Category: SyntaxWordCategoryKeyword},
	{Word: SyntaxWordDelay, Spelling: "DELAY", Category: SyntaxWordCategoryContextual},
	{Word: SyntaxWordDelete, Spelling: "DELETE", Category: SyntaxWordCategoryKeyword},
	{Word: SyntaxWordDesc, Spelling: "DESC", Category: SyntaxWordCategoryKeyword},
	{Word: SyntaxWordDispatch, Spelling: "DISPATCH", Category: SyntaxWordCategoryKeyword},
	{Word: SyntaxWordDistinct, Spelling: "DISTINCT", Category: SyntaxWordCategoryKeyword},
	{Word: SyntaxWordDo, Spelling: "DO", Category: SyntaxWordCategoryKeyword},
	{Word: SyntaxWordError, Spelling: "ERROR", Category: SyntaxWordCategoryContextual},
	{Word: SyntaxWordEvent, Spelling: "EVENT", Category: SyntaxWordCategoryKeyword},
	{Word: SyntaxWordEvery, Spelling: "EVERY", Category: SyntaxWordCategoryKeyword},
	{Word: SyntaxWordExists, Spelling: "EXISTS", Category: SyntaxWordCategoryKeyword},
	{Word: SyntaxWordFail, Spelling: "FAIL", Category: SyntaxWordCategoryContextual},
	{Word: SyntaxWordFalse, Spelling: "FALSE", Category: SyntaxWordCategoryLiteral},
	{Word: SyntaxWordFilter, Spelling: "FILTER", Category: SyntaxWordCategoryKeyword},
	{Word: SyntaxWordFor, Spelling: "FOR", Category: SyntaxWordCategoryKeyword},
	{Word: SyntaxWordFunc, Spelling: "FUNC", Category: SyntaxWordCategoryKeyword},
	{Word: SyntaxWordIn, Spelling: "IN", Category: SyntaxWordCategoryOperator},
	{Word: SyntaxWordInto, Spelling: "INTO", Category: SyntaxWordCategoryKeyword},
	{Word: SyntaxWordJitter, Spelling: "JITTER", Category: SyntaxWordCategoryKeyword},
	{Word: SyntaxWordKeep, Spelling: "KEEP", Category: SyntaxWordCategoryKeyword},
	{Word: SyntaxWordLeast, Spelling: "LEAST", Category: SyntaxWordCategoryKeyword},
	{Word: SyntaxWordLet, Spelling: "LET", Category: SyntaxWordCategoryKeyword},
	{Word: SyntaxWordLike, Spelling: "LIKE", Category: SyntaxWordCategoryOperator},
	{Word: SyntaxWordLimit, Spelling: "LIMIT", Category: SyntaxWordCategoryKeyword},
	{Word: SyntaxWordMatch, Spelling: "MATCH", Category: SyntaxWordCategoryKeyword},
	{Word: SyntaxWordNone, Spelling: "NONE", Category: SyntaxWordCategoryLiteral},
	{Word: SyntaxWordNot, Spelling: "NOT", Category: SyntaxWordCategoryOperator},
	{Word: SyntaxWordNull, Spelling: "NULL", Category: SyntaxWordCategoryLiteral},
	{Word: SyntaxWordOn, Spelling: "ON", Category: SyntaxWordCategoryContextual},
	{Word: SyntaxWordOne, Spelling: "ONE", Category: SyntaxWordCategoryKeyword},
	{Word: SyntaxWordOptions, Spelling: "OPTIONS", Category: SyntaxWordCategoryKeyword},
	{Word: SyntaxWordOr, Spelling: "OR", Category: SyntaxWordCategoryOperator},
	{Word: SyntaxWordQuery, Spelling: "QUERY", Category: SyntaxWordCategoryKeyword},
	{Word: SyntaxWordRetry, Spelling: "RETRY", Category: SyntaxWordCategoryContextual},
	{Word: SyntaxWordReturn, Spelling: "RETURN", Category: SyntaxWordCategoryKeyword},
	{Word: SyntaxWordSort, Spelling: "SORT", Category: SyntaxWordCategoryKeyword},
	{Word: SyntaxWordTimeout, Spelling: "TIMEOUT", Category: SyntaxWordCategoryKeyword},
	{Word: SyntaxWordTrigger, Spelling: "TRIGGER", Category: SyntaxWordCategoryKeyword},
	{Word: SyntaxWordTrue, Spelling: "TRUE", Category: SyntaxWordCategoryLiteral},
	{Word: SyntaxWordUse, Spelling: "USE", Category: SyntaxWordCategoryKeyword},
	{Word: SyntaxWordUsing, Spelling: "USING", Category: SyntaxWordCategoryKeyword},
	{Word: SyntaxWordValue, Spelling: "VALUE", Category: SyntaxWordCategoryKeyword},
	{Word: SyntaxWordVar, Spelling: "VAR", Category: SyntaxWordCategoryKeyword},
	{Word: SyntaxWordWaitfor, Spelling: "WAITFOR", Category: SyntaxWordCategoryKeyword},
	{Word: SyntaxWordWhen, Spelling: "WHEN", Category: SyntaxWordCategoryKeyword},
	{Word: SyntaxWordWhile, Spelling: "WHILE", Category: SyntaxWordCategoryKeyword},
	{Word: SyntaxWordWith, Spelling: "WITH", Category: SyntaxWordCategoryKeyword},
}

// SyntaxWords returns canonical FQL language words in spelling order.
// The returned slice is a defensive copy.
func SyntaxWords() []SyntaxWordInfo {
	return append([]SyntaxWordInfo(nil), syntaxWords...)
}
