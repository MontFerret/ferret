// Code generated from antlr/FqlParser.g4 by ANTLR 4.13.2. DO NOT EDIT.

package fql // FqlParser
import (
	"fmt"
	"strconv"
	"sync"

	"github.com/antlr4-go/antlr/v4"
)

// Suppress unused import errors
var _ = fmt.Printf
var _ = strconv.Itoa
var _ = sync.Once{}

type FqlParser struct {
	*antlr.BaseParser
}

var FqlParserParserStaticData struct {
	once                   sync.Once
	serializedATN          []int32
	LiteralNames           []string
	SymbolicNames          []string
	RuleNames              []string
	PredictionContextCache *antlr.PredictionContextCache
	atn                    *antlr.ATN
	decisionToDFA          []*antlr.DFA
}

func fqlparserParserInit() {
	staticData := &FqlParserParserStaticData
	staticData.LiteralNames = []string{
		"", "", "", "", "", "':'", "';'", "'.'", "','", "'['", "']'", "'('",
		"')'", "'{'", "'}'", "'>'", "'<'", "'=='", "'>='", "'<='", "'!='", "'*'",
		"'/'", "'%'", "'+'", "'-'", "", "", "", "'='", "'?'", "'!~'", "'=~'",
		"'FOR'", "'RETURN'", "'WAITFOR'", "'OPTIONS'", "'TIMEOUT'", "'DISTINCT'",
		"'FILTER'", "'CURRENT'", "'SORT'", "'LIMIT'", "'LET'", "'COLLECT'",
		"", "'NONE'", "'NULL'", "", "'USE'", "'INTO'", "'KEEP'", "'WITH'", "'ALL'",
		"'ANY'", "'AGGREGATE'", "'EVENT'", "'LIKE'", "", "'IN'", "'DO'", "'WHILE'",
		"'STEP'", "'@'",
	}
	staticData.SymbolicNames = []string{
		"", "MultiLineComment", "SingleLineComment", "WhiteSpaces", "LineTerminator",
		"Colon", "SemiColon", "Dot", "Comma", "OpenBracket", "CloseBracket",
		"OpenParen", "CloseParen", "OpenBrace", "CloseBrace", "Gt", "Lt", "Eq",
		"Gte", "Lte", "Neq", "Multi", "Div", "Mod", "Plus", "Minus", "And",
		"Or", "Range", "Assign", "QuestionMark", "RegexNotMatch", "RegexMatch",
		"For", "Return", "Waitfor", "Options", "Timeout", "Distinct", "Filter",
		"Current", "Sort", "Limit", "Let", "Collect", "SortDirection", "None",
		"Null", "BooleanLiteral", "Use", "Into", "Keep", "With", "All", "Any",
		"Aggregate", "Event", "Like", "Not", "In", "Do", "While", "Step", "Param",
		"Identifier", "IgnoreIdentifier", "StringLiteral", "IntegerLiteral",
		"FloatLiteral", "NamespaceSegment", "UnknownIdentifier",
	}
	staticData.RuleNames = []string{
		"program", "head", "useExpression", "use", "body", "bodyStatement",
		"bodyExpression", "variableDeclaration", "returnExpression", "forExpression",
		"forExpressionSource", "forExpressionClause", "forExpressionStatement",
		"forExpressionBody", "forExpressionReturn", "filterClause", "limitClause",
		"limitClauseValue", "sortClause", "sortClauseExpression", "collectClause",
		"collectSelector", "collectGrouping", "collectAggregateSelector", "collectAggregator",
		"collectGroupProjection", "collectGroupProjectionFilter", "collectCounter",
		"waitForExpression", "waitForEventName", "waitForEventSource", "optionsClause",
		"timeoutClause", "param", "variable", "literal", "arrayLiteral", "objectLiteral",
		"booleanLiteral", "stringLiteral", "floatLiteral", "integerLiteral",
		"noneLiteral", "propertyAssignment", "computedPropertyName", "propertyName",
		"namespaceIdentifier", "namespace", "memberExpression", "memberExpressionSource",
		"functionCallExpression", "functionCall", "functionName", "argumentList",
		"memberExpressionPath", "safeReservedWord", "unsafeReservedWord", "rangeOperator",
		"rangeOperand", "expression", "predicate", "expressionAtom", "arrayOperator",
		"equalityOperator", "inOperator", "likeOperator", "unaryOperator", "regexpOperator",
		"logicalAndOperator", "logicalOrOperator", "multiplicativeOperator",
		"additiveOperator", "errorOperator",
	}
	staticData.PredictionContextCache = antlr.NewPredictionContextCache()
	staticData.serializedATN = []int32{
		4, 1, 70, 684, 2, 0, 7, 0, 2, 1, 7, 1, 2, 2, 7, 2, 2, 3, 7, 3, 2, 4, 7,
		4, 2, 5, 7, 5, 2, 6, 7, 6, 2, 7, 7, 7, 2, 8, 7, 8, 2, 9, 7, 9, 2, 10, 7,
		10, 2, 11, 7, 11, 2, 12, 7, 12, 2, 13, 7, 13, 2, 14, 7, 14, 2, 15, 7, 15,
		2, 16, 7, 16, 2, 17, 7, 17, 2, 18, 7, 18, 2, 19, 7, 19, 2, 20, 7, 20, 2,
		21, 7, 21, 2, 22, 7, 22, 2, 23, 7, 23, 2, 24, 7, 24, 2, 25, 7, 25, 2, 26,
		7, 26, 2, 27, 7, 27, 2, 28, 7, 28, 2, 29, 7, 29, 2, 30, 7, 30, 2, 31, 7,
		31, 2, 32, 7, 32, 2, 33, 7, 33, 2, 34, 7, 34, 2, 35, 7, 35, 2, 36, 7, 36,
		2, 37, 7, 37, 2, 38, 7, 38, 2, 39, 7, 39, 2, 40, 7, 40, 2, 41, 7, 41, 2,
		42, 7, 42, 2, 43, 7, 43, 2, 44, 7, 44, 2, 45, 7, 45, 2, 46, 7, 46, 2, 47,
		7, 47, 2, 48, 7, 48, 2, 49, 7, 49, 2, 50, 7, 50, 2, 51, 7, 51, 2, 52, 7,
		52, 2, 53, 7, 53, 2, 54, 7, 54, 2, 55, 7, 55, 2, 56, 7, 56, 2, 57, 7, 57,
		2, 58, 7, 58, 2, 59, 7, 59, 2, 60, 7, 60, 2, 61, 7, 61, 2, 62, 7, 62, 2,
		63, 7, 63, 2, 64, 7, 64, 2, 65, 7, 65, 2, 66, 7, 66, 2, 67, 7, 67, 2, 68,
		7, 68, 2, 69, 7, 69, 2, 70, 7, 70, 2, 71, 7, 71, 2, 72, 7, 72, 1, 0, 5,
		0, 148, 8, 0, 10, 0, 12, 0, 151, 9, 0, 1, 0, 1, 0, 1, 1, 1, 1, 1, 2, 1,
		2, 1, 3, 1, 3, 1, 3, 1, 4, 5, 4, 163, 8, 4, 10, 4, 12, 4, 166, 9, 4, 1,
		4, 1, 4, 1, 5, 1, 5, 1, 5, 3, 5, 173, 8, 5, 1, 6, 1, 6, 3, 6, 177, 8, 6,
		1, 7, 1, 7, 1, 7, 1, 7, 1, 7, 1, 7, 1, 7, 1, 7, 1, 7, 3, 7, 188, 8, 7,
		1, 8, 1, 8, 3, 8, 192, 8, 8, 1, 8, 1, 8, 1, 9, 1, 9, 1, 9, 1, 9, 3, 9,
		200, 8, 9, 1, 9, 1, 9, 1, 9, 5, 9, 205, 8, 9, 10, 9, 12, 9, 208, 9, 9,
		1, 9, 1, 9, 1, 9, 1, 9, 1, 9, 1, 9, 1, 9, 1, 9, 1, 9, 1, 9, 1, 9, 1, 9,
		1, 9, 5, 9, 223, 8, 9, 10, 9, 12, 9, 226, 9, 9, 1, 9, 1, 9, 1, 9, 1, 9,
		1, 9, 3, 9, 233, 8, 9, 1, 9, 1, 9, 1, 9, 5, 9, 238, 8, 9, 10, 9, 12, 9,
		241, 9, 9, 1, 9, 1, 9, 3, 9, 245, 8, 9, 1, 10, 1, 10, 1, 10, 1, 10, 1,
		10, 1, 10, 1, 10, 3, 10, 254, 8, 10, 1, 11, 1, 11, 1, 11, 1, 11, 3, 11,
		260, 8, 11, 1, 12, 1, 12, 3, 12, 264, 8, 12, 1, 13, 1, 13, 3, 13, 268,
		8, 13, 1, 14, 1, 14, 3, 14, 272, 8, 14, 1, 15, 1, 15, 1, 15, 1, 16, 1,
		16, 1, 16, 1, 16, 3, 16, 281, 8, 16, 1, 17, 1, 17, 1, 17, 1, 17, 1, 17,
		3, 17, 288, 8, 17, 1, 18, 1, 18, 1, 18, 1, 18, 5, 18, 294, 8, 18, 10, 18,
		12, 18, 297, 9, 18, 1, 19, 1, 19, 3, 19, 301, 8, 19, 1, 20, 1, 20, 1, 20,
		1, 20, 1, 20, 1, 20, 1, 20, 1, 20, 1, 20, 1, 20, 1, 20, 1, 20, 1, 20, 1,
		20, 1, 20, 1, 20, 1, 20, 1, 20, 1, 20, 1, 20, 1, 20, 1, 20, 1, 20, 1, 20,
		1, 20, 1, 20, 1, 20, 3, 20, 330, 8, 20, 1, 21, 1, 21, 1, 21, 1, 21, 1,
		22, 1, 22, 1, 22, 5, 22, 339, 8, 22, 10, 22, 12, 22, 342, 9, 22, 1, 23,
		1, 23, 1, 23, 1, 23, 1, 24, 1, 24, 1, 24, 1, 24, 5, 24, 352, 8, 24, 10,
		24, 12, 24, 355, 9, 24, 1, 25, 1, 25, 1, 25, 1, 25, 1, 25, 3, 25, 362,
		8, 25, 3, 25, 364, 8, 25, 1, 26, 1, 26, 1, 26, 1, 26, 5, 26, 370, 8, 26,
		10, 26, 12, 26, 373, 9, 26, 1, 27, 1, 27, 1, 27, 1, 27, 1, 27, 1, 28, 1,
		28, 1, 28, 1, 28, 1, 28, 1, 28, 3, 28, 386, 8, 28, 1, 28, 3, 28, 389, 8,
		28, 1, 28, 3, 28, 392, 8, 28, 1, 29, 1, 29, 1, 29, 1, 29, 1, 29, 3, 29,
		399, 8, 29, 1, 30, 1, 30, 1, 30, 3, 30, 404, 8, 30, 1, 31, 1, 31, 1, 31,
		1, 32, 1, 32, 1, 32, 1, 32, 1, 32, 1, 32, 3, 32, 415, 8, 32, 1, 33, 1,
		33, 1, 33, 1, 33, 3, 33, 421, 8, 33, 1, 34, 1, 34, 3, 34, 425, 8, 34, 1,
		35, 1, 35, 1, 35, 1, 35, 1, 35, 1, 35, 1, 35, 3, 35, 434, 8, 35, 1, 36,
		1, 36, 3, 36, 438, 8, 36, 1, 36, 1, 36, 1, 37, 1, 37, 1, 37, 1, 37, 5,
		37, 446, 8, 37, 10, 37, 12, 37, 449, 9, 37, 1, 37, 3, 37, 452, 8, 37, 3,
		37, 454, 8, 37, 1, 37, 1, 37, 1, 38, 1, 38, 1, 39, 1, 39, 1, 40, 1, 40,
		1, 41, 1, 41, 1, 42, 1, 42, 1, 43, 1, 43, 1, 43, 1, 43, 1, 43, 1, 43, 1,
		43, 1, 43, 1, 43, 3, 43, 477, 8, 43, 1, 44, 1, 44, 1, 44, 1, 44, 1, 45,
		1, 45, 1, 45, 1, 45, 1, 45, 3, 45, 488, 8, 45, 1, 46, 1, 46, 1, 46, 1,
		47, 5, 47, 494, 8, 47, 10, 47, 12, 47, 497, 9, 47, 1, 48, 1, 48, 4, 48,
		501, 8, 48, 11, 48, 12, 48, 502, 1, 49, 1, 49, 1, 49, 1, 49, 1, 49, 3,
		49, 510, 8, 49, 1, 50, 1, 50, 3, 50, 514, 8, 50, 1, 51, 1, 51, 1, 51, 1,
		51, 3, 51, 520, 8, 51, 1, 51, 1, 51, 1, 52, 1, 52, 1, 52, 3, 52, 527, 8,
		52, 1, 53, 1, 53, 1, 53, 5, 53, 532, 8, 53, 10, 53, 12, 53, 535, 9, 53,
		1, 53, 3, 53, 538, 8, 53, 1, 54, 3, 54, 541, 8, 54, 1, 54, 1, 54, 1, 54,
		1, 54, 1, 54, 3, 54, 548, 8, 54, 1, 54, 3, 54, 551, 8, 54, 1, 55, 1, 55,
		1, 56, 1, 56, 1, 57, 1, 57, 1, 57, 1, 57, 1, 58, 1, 58, 1, 58, 3, 58, 564,
		8, 58, 1, 59, 1, 59, 1, 59, 1, 59, 1, 59, 3, 59, 571, 8, 59, 1, 59, 1,
		59, 1, 59, 1, 59, 1, 59, 1, 59, 1, 59, 1, 59, 1, 59, 1, 59, 1, 59, 3, 59,
		584, 8, 59, 1, 59, 1, 59, 5, 59, 588, 8, 59, 10, 59, 12, 59, 591, 9, 59,
		1, 60, 1, 60, 1, 60, 1, 60, 1, 60, 1, 60, 1, 60, 1, 60, 1, 60, 1, 60, 1,
		60, 1, 60, 1, 60, 1, 60, 1, 60, 1, 60, 1, 60, 1, 60, 1, 60, 5, 60, 612,
		8, 60, 10, 60, 12, 60, 615, 9, 60, 1, 61, 1, 61, 1, 61, 1, 61, 1, 61, 1,
		61, 1, 61, 1, 61, 1, 61, 1, 61, 1, 61, 3, 61, 628, 8, 61, 1, 61, 1, 61,
		3, 61, 632, 8, 61, 3, 61, 634, 8, 61, 1, 61, 1, 61, 1, 61, 1, 61, 1, 61,
		1, 61, 1, 61, 1, 61, 1, 61, 1, 61, 1, 61, 1, 61, 5, 61, 648, 8, 61, 10,
		61, 12, 61, 651, 9, 61, 1, 62, 1, 62, 1, 62, 3, 62, 656, 8, 62, 1, 63,
		1, 63, 1, 64, 3, 64, 661, 8, 64, 1, 64, 1, 64, 1, 65, 3, 65, 666, 8, 65,
		1, 65, 1, 65, 1, 66, 1, 66, 1, 67, 1, 67, 1, 68, 1, 68, 1, 69, 1, 69, 1,
		70, 1, 70, 1, 71, 1, 71, 1, 72, 1, 72, 1, 72, 0, 3, 118, 120, 122, 73,
		0, 2, 4, 6, 8, 10, 12, 14, 16, 18, 20, 22, 24, 26, 28, 30, 32, 34, 36,
		38, 40, 42, 44, 46, 48, 50, 52, 54, 56, 58, 60, 62, 64, 66, 68, 70, 72,
		74, 76, 78, 80, 82, 84, 86, 88, 90, 92, 94, 96, 98, 100, 102, 104, 106,
		108, 110, 112, 114, 116, 118, 120, 122, 124, 126, 128, 130, 132, 134, 136,
		138, 140, 142, 144, 0, 10, 1, 0, 64, 65, 1, 0, 46, 47, 5, 0, 26, 27, 36,
		42, 44, 45, 50, 56, 62, 62, 4, 0, 33, 35, 43, 43, 46, 49, 57, 61, 2, 0,
		46, 46, 53, 54, 1, 0, 15, 20, 2, 0, 24, 25, 58, 58, 1, 0, 31, 32, 1, 0,
		21, 23, 1, 0, 24, 25, 727, 0, 149, 1, 0, 0, 0, 2, 154, 1, 0, 0, 0, 4, 156,
		1, 0, 0, 0, 6, 158, 1, 0, 0, 0, 8, 164, 1, 0, 0, 0, 10, 172, 1, 0, 0, 0,
		12, 176, 1, 0, 0, 0, 14, 187, 1, 0, 0, 0, 16, 189, 1, 0, 0, 0, 18, 244,
		1, 0, 0, 0, 20, 253, 1, 0, 0, 0, 22, 259, 1, 0, 0, 0, 24, 263, 1, 0, 0,
		0, 26, 267, 1, 0, 0, 0, 28, 271, 1, 0, 0, 0, 30, 273, 1, 0, 0, 0, 32, 276,
		1, 0, 0, 0, 34, 287, 1, 0, 0, 0, 36, 289, 1, 0, 0, 0, 38, 298, 1, 0, 0,
		0, 40, 329, 1, 0, 0, 0, 42, 331, 1, 0, 0, 0, 44, 335, 1, 0, 0, 0, 46, 343,
		1, 0, 0, 0, 48, 347, 1, 0, 0, 0, 50, 363, 1, 0, 0, 0, 52, 365, 1, 0, 0,
		0, 54, 374, 1, 0, 0, 0, 56, 379, 1, 0, 0, 0, 58, 398, 1, 0, 0, 0, 60, 403,
		1, 0, 0, 0, 62, 405, 1, 0, 0, 0, 64, 408, 1, 0, 0, 0, 66, 420, 1, 0, 0,
		0, 68, 424, 1, 0, 0, 0, 70, 433, 1, 0, 0, 0, 72, 435, 1, 0, 0, 0, 74, 441,
		1, 0, 0, 0, 76, 457, 1, 0, 0, 0, 78, 459, 1, 0, 0, 0, 80, 461, 1, 0, 0,
		0, 82, 463, 1, 0, 0, 0, 84, 465, 1, 0, 0, 0, 86, 476, 1, 0, 0, 0, 88, 478,
		1, 0, 0, 0, 90, 487, 1, 0, 0, 0, 92, 489, 1, 0, 0, 0, 94, 495, 1, 0, 0,
		0, 96, 498, 1, 0, 0, 0, 98, 509, 1, 0, 0, 0, 100, 511, 1, 0, 0, 0, 102,
		515, 1, 0, 0, 0, 104, 526, 1, 0, 0, 0, 106, 528, 1, 0, 0, 0, 108, 550,
		1, 0, 0, 0, 110, 552, 1, 0, 0, 0, 112, 554, 1, 0, 0, 0, 114, 556, 1, 0,
		0, 0, 116, 563, 1, 0, 0, 0, 118, 570, 1, 0, 0, 0, 120, 592, 1, 0, 0, 0,
		122, 633, 1, 0, 0, 0, 124, 652, 1, 0, 0, 0, 126, 657, 1, 0, 0, 0, 128,
		660, 1, 0, 0, 0, 130, 665, 1, 0, 0, 0, 132, 669, 1, 0, 0, 0, 134, 671,
		1, 0, 0, 0, 136, 673, 1, 0, 0, 0, 138, 675, 1, 0, 0, 0, 140, 677, 1, 0,
		0, 0, 142, 679, 1, 0, 0, 0, 144, 681, 1, 0, 0, 0, 146, 148, 3, 2, 1, 0,
		147, 146, 1, 0, 0, 0, 148, 151, 1, 0, 0, 0, 149, 147, 1, 0, 0, 0, 149,
		150, 1, 0, 0, 0, 150, 152, 1, 0, 0, 0, 151, 149, 1, 0, 0, 0, 152, 153,
		3, 8, 4, 0, 153, 1, 1, 0, 0, 0, 154, 155, 3, 4, 2, 0, 155, 3, 1, 0, 0,
		0, 156, 157, 3, 6, 3, 0, 157, 5, 1, 0, 0, 0, 158, 159, 5, 49, 0, 0, 159,
		160, 3, 92, 46, 0, 160, 7, 1, 0, 0, 0, 161, 163, 3, 10, 5, 0, 162, 161,
		1, 0, 0, 0, 163, 166, 1, 0, 0, 0, 164, 162, 1, 0, 0, 0, 164, 165, 1, 0,
		0, 0, 165, 167, 1, 0, 0, 0, 166, 164, 1, 0, 0, 0, 167, 168, 3, 12, 6, 0,
		168, 9, 1, 0, 0, 0, 169, 173, 3, 14, 7, 0, 170, 173, 3, 100, 50, 0, 171,
		173, 3, 56, 28, 0, 172, 169, 1, 0, 0, 0, 172, 170, 1, 0, 0, 0, 172, 171,
		1, 0, 0, 0, 173, 11, 1, 0, 0, 0, 174, 177, 3, 16, 8, 0, 175, 177, 3, 18,
		9, 0, 176, 174, 1, 0, 0, 0, 176, 175, 1, 0, 0, 0, 177, 13, 1, 0, 0, 0,
		178, 179, 5, 43, 0, 0, 179, 180, 7, 0, 0, 0, 180, 181, 5, 29, 0, 0, 181,
		188, 3, 118, 59, 0, 182, 183, 5, 43, 0, 0, 183, 184, 3, 110, 55, 0, 184,
		185, 5, 29, 0, 0, 185, 186, 3, 118, 59, 0, 186, 188, 1, 0, 0, 0, 187, 178,
		1, 0, 0, 0, 187, 182, 1, 0, 0, 0, 188, 15, 1, 0, 0, 0, 189, 191, 5, 34,
		0, 0, 190, 192, 5, 38, 0, 0, 191, 190, 1, 0, 0, 0, 191, 192, 1, 0, 0, 0,
		192, 193, 1, 0, 0, 0, 193, 194, 3, 118, 59, 0, 194, 17, 1, 0, 0, 0, 195,
		196, 5, 33, 0, 0, 196, 199, 7, 0, 0, 0, 197, 198, 5, 8, 0, 0, 198, 200,
		5, 64, 0, 0, 199, 197, 1, 0, 0, 0, 199, 200, 1, 0, 0, 0, 200, 201, 1, 0,
		0, 0, 201, 202, 5, 59, 0, 0, 202, 206, 3, 20, 10, 0, 203, 205, 3, 26, 13,
		0, 204, 203, 1, 0, 0, 0, 205, 208, 1, 0, 0, 0, 206, 204, 1, 0, 0, 0, 206,
		207, 1, 0, 0, 0, 207, 209, 1, 0, 0, 0, 208, 206, 1, 0, 0, 0, 209, 210,
		3, 28, 14, 0, 210, 245, 1, 0, 0, 0, 211, 212, 5, 33, 0, 0, 212, 213, 7,
		0, 0, 0, 213, 214, 5, 29, 0, 0, 214, 215, 3, 118, 59, 0, 215, 216, 5, 61,
		0, 0, 216, 217, 3, 118, 59, 0, 217, 218, 5, 62, 0, 0, 218, 219, 7, 0, 0,
		0, 219, 220, 5, 29, 0, 0, 220, 224, 3, 118, 59, 0, 221, 223, 3, 26, 13,
		0, 222, 221, 1, 0, 0, 0, 223, 226, 1, 0, 0, 0, 224, 222, 1, 0, 0, 0, 224,
		225, 1, 0, 0, 0, 225, 227, 1, 0, 0, 0, 226, 224, 1, 0, 0, 0, 227, 228,
		3, 28, 14, 0, 228, 245, 1, 0, 0, 0, 229, 230, 5, 33, 0, 0, 230, 232, 7,
		0, 0, 0, 231, 233, 5, 60, 0, 0, 232, 231, 1, 0, 0, 0, 232, 233, 1, 0, 0,
		0, 233, 234, 1, 0, 0, 0, 234, 235, 5, 61, 0, 0, 235, 239, 3, 118, 59, 0,
		236, 238, 3, 26, 13, 0, 237, 236, 1, 0, 0, 0, 238, 241, 1, 0, 0, 0, 239,
		237, 1, 0, 0, 0, 239, 240, 1, 0, 0, 0, 240, 242, 1, 0, 0, 0, 241, 239,
		1, 0, 0, 0, 242, 243, 3, 28, 14, 0, 243, 245, 1, 0, 0, 0, 244, 195, 1,
		0, 0, 0, 244, 211, 1, 0, 0, 0, 244, 229, 1, 0, 0, 0, 245, 19, 1, 0, 0,
		0, 246, 254, 3, 100, 50, 0, 247, 254, 3, 72, 36, 0, 248, 254, 3, 74, 37,
		0, 249, 254, 3, 68, 34, 0, 250, 254, 3, 96, 48, 0, 251, 254, 3, 114, 57,
		0, 252, 254, 3, 66, 33, 0, 253, 246, 1, 0, 0, 0, 253, 247, 1, 0, 0, 0,
		253, 248, 1, 0, 0, 0, 253, 249, 1, 0, 0, 0, 253, 250, 1, 0, 0, 0, 253,
		251, 1, 0, 0, 0, 253, 252, 1, 0, 0, 0, 254, 21, 1, 0, 0, 0, 255, 260, 3,
		32, 16, 0, 256, 260, 3, 36, 18, 0, 257, 260, 3, 30, 15, 0, 258, 260, 3,
		40, 20, 0, 259, 255, 1, 0, 0, 0, 259, 256, 1, 0, 0, 0, 259, 257, 1, 0,
		0, 0, 259, 258, 1, 0, 0, 0, 260, 23, 1, 0, 0, 0, 261, 264, 3, 14, 7, 0,
		262, 264, 3, 100, 50, 0, 263, 261, 1, 0, 0, 0, 263, 262, 1, 0, 0, 0, 264,
		25, 1, 0, 0, 0, 265, 268, 3, 24, 12, 0, 266, 268, 3, 22, 11, 0, 267, 265,
		1, 0, 0, 0, 267, 266, 1, 0, 0, 0, 268, 27, 1, 0, 0, 0, 269, 272, 3, 16,
		8, 0, 270, 272, 3, 18, 9, 0, 271, 269, 1, 0, 0, 0, 271, 270, 1, 0, 0, 0,
		272, 29, 1, 0, 0, 0, 273, 274, 5, 39, 0, 0, 274, 275, 3, 118, 59, 0, 275,
		31, 1, 0, 0, 0, 276, 277, 5, 42, 0, 0, 277, 280, 3, 34, 17, 0, 278, 279,
		5, 8, 0, 0, 279, 281, 3, 34, 17, 0, 280, 278, 1, 0, 0, 0, 280, 281, 1,
		0, 0, 0, 281, 33, 1, 0, 0, 0, 282, 288, 3, 82, 41, 0, 283, 288, 3, 66,
		33, 0, 284, 288, 3, 68, 34, 0, 285, 288, 3, 100, 50, 0, 286, 288, 3, 96,
		48, 0, 287, 282, 1, 0, 0, 0, 287, 283, 1, 0, 0, 0, 287, 284, 1, 0, 0, 0,
		287, 285, 1, 0, 0, 0, 287, 286, 1, 0, 0, 0, 288, 35, 1, 0, 0, 0, 289, 290,
		5, 41, 0, 0, 290, 295, 3, 38, 19, 0, 291, 292, 5, 8, 0, 0, 292, 294, 3,
		38, 19, 0, 293, 291, 1, 0, 0, 0, 294, 297, 1, 0, 0, 0, 295, 293, 1, 0,
		0, 0, 295, 296, 1, 0, 0, 0, 296, 37, 1, 0, 0, 0, 297, 295, 1, 0, 0, 0,
		298, 300, 3, 118, 59, 0, 299, 301, 5, 45, 0, 0, 300, 299, 1, 0, 0, 0, 300,
		301, 1, 0, 0, 0, 301, 39, 1, 0, 0, 0, 302, 303, 5, 44, 0, 0, 303, 304,
		3, 44, 22, 0, 304, 305, 3, 48, 24, 0, 305, 306, 3, 50, 25, 0, 306, 330,
		1, 0, 0, 0, 307, 308, 5, 44, 0, 0, 308, 309, 3, 44, 22, 0, 309, 310, 3,
		48, 24, 0, 310, 330, 1, 0, 0, 0, 311, 312, 5, 44, 0, 0, 312, 313, 3, 44,
		22, 0, 313, 314, 3, 50, 25, 0, 314, 330, 1, 0, 0, 0, 315, 316, 5, 44, 0,
		0, 316, 317, 3, 44, 22, 0, 317, 318, 3, 54, 27, 0, 318, 330, 1, 0, 0, 0,
		319, 320, 5, 44, 0, 0, 320, 330, 3, 44, 22, 0, 321, 322, 5, 44, 0, 0, 322,
		323, 3, 48, 24, 0, 323, 324, 3, 50, 25, 0, 324, 330, 1, 0, 0, 0, 325, 326,
		5, 44, 0, 0, 326, 330, 3, 48, 24, 0, 327, 328, 5, 44, 0, 0, 328, 330, 3,
		54, 27, 0, 329, 302, 1, 0, 0, 0, 329, 307, 1, 0, 0, 0, 329, 311, 1, 0,
		0, 0, 329, 315, 1, 0, 0, 0, 329, 319, 1, 0, 0, 0, 329, 321, 1, 0, 0, 0,
		329, 325, 1, 0, 0, 0, 329, 327, 1, 0, 0, 0, 330, 41, 1, 0, 0, 0, 331, 332,
		5, 64, 0, 0, 332, 333, 5, 29, 0, 0, 333, 334, 3, 118, 59, 0, 334, 43, 1,
		0, 0, 0, 335, 340, 3, 42, 21, 0, 336, 337, 5, 8, 0, 0, 337, 339, 3, 42,
		21, 0, 338, 336, 1, 0, 0, 0, 339, 342, 1, 0, 0, 0, 340, 338, 1, 0, 0, 0,
		340, 341, 1, 0, 0, 0, 341, 45, 1, 0, 0, 0, 342, 340, 1, 0, 0, 0, 343, 344,
		5, 64, 0, 0, 344, 345, 5, 29, 0, 0, 345, 346, 3, 100, 50, 0, 346, 47, 1,
		0, 0, 0, 347, 348, 5, 55, 0, 0, 348, 353, 3, 46, 23, 0, 349, 350, 5, 8,
		0, 0, 350, 352, 3, 46, 23, 0, 351, 349, 1, 0, 0, 0, 352, 355, 1, 0, 0,
		0, 353, 351, 1, 0, 0, 0, 353, 354, 1, 0, 0, 0, 354, 49, 1, 0, 0, 0, 355,
		353, 1, 0, 0, 0, 356, 357, 5, 50, 0, 0, 357, 364, 3, 42, 21, 0, 358, 359,
		5, 50, 0, 0, 359, 361, 5, 64, 0, 0, 360, 362, 3, 52, 26, 0, 361, 360, 1,
		0, 0, 0, 361, 362, 1, 0, 0, 0, 362, 364, 1, 0, 0, 0, 363, 356, 1, 0, 0,
		0, 363, 358, 1, 0, 0, 0, 364, 51, 1, 0, 0, 0, 365, 366, 5, 51, 0, 0, 366,
		371, 5, 64, 0, 0, 367, 368, 5, 8, 0, 0, 368, 370, 5, 64, 0, 0, 369, 367,
		1, 0, 0, 0, 370, 373, 1, 0, 0, 0, 371, 369, 1, 0, 0, 0, 371, 372, 1, 0,
		0, 0, 372, 53, 1, 0, 0, 0, 373, 371, 1, 0, 0, 0, 374, 375, 5, 52, 0, 0,
		375, 376, 5, 64, 0, 0, 376, 377, 5, 50, 0, 0, 377, 378, 5, 64, 0, 0, 378,
		55, 1, 0, 0, 0, 379, 380, 5, 35, 0, 0, 380, 381, 5, 56, 0, 0, 381, 382,
		3, 58, 29, 0, 382, 383, 5, 59, 0, 0, 383, 385, 3, 60, 30, 0, 384, 386,
		3, 62, 31, 0, 385, 384, 1, 0, 0, 0, 385, 386, 1, 0, 0, 0, 386, 388, 1,
		0, 0, 0, 387, 389, 3, 30, 15, 0, 388, 387, 1, 0, 0, 0, 388, 389, 1, 0,
		0, 0, 389, 391, 1, 0, 0, 0, 390, 392, 3, 64, 32, 0, 391, 390, 1, 0, 0,
		0, 391, 392, 1, 0, 0, 0, 392, 57, 1, 0, 0, 0, 393, 399, 3, 78, 39, 0, 394,
		399, 3, 68, 34, 0, 395, 399, 3, 66, 33, 0, 396, 399, 3, 100, 50, 0, 397,
		399, 3, 96, 48, 0, 398, 393, 1, 0, 0, 0, 398, 394, 1, 0, 0, 0, 398, 395,
		1, 0, 0, 0, 398, 396, 1, 0, 0, 0, 398, 397, 1, 0, 0, 0, 399, 59, 1, 0,
		0, 0, 400, 404, 3, 100, 50, 0, 401, 404, 3, 68, 34, 0, 402, 404, 3, 96,
		48, 0, 403, 400, 1, 0, 0, 0, 403, 401, 1, 0, 0, 0, 403, 402, 1, 0, 0, 0,
		404, 61, 1, 0, 0, 0, 405, 406, 5, 36, 0, 0, 406, 407, 3, 74, 37, 0, 407,
		63, 1, 0, 0, 0, 408, 414, 5, 37, 0, 0, 409, 415, 3, 82, 41, 0, 410, 415,
		3, 68, 34, 0, 411, 415, 3, 66, 33, 0, 412, 415, 3, 96, 48, 0, 413, 415,
		3, 102, 51, 0, 414, 409, 1, 0, 0, 0, 414, 410, 1, 0, 0, 0, 414, 411, 1,
		0, 0, 0, 414, 412, 1, 0, 0, 0, 414, 413, 1, 0, 0, 0, 415, 65, 1, 0, 0,
		0, 416, 417, 5, 63, 0, 0, 417, 421, 5, 64, 0, 0, 418, 419, 5, 63, 0, 0,
		419, 421, 3, 110, 55, 0, 420, 416, 1, 0, 0, 0, 420, 418, 1, 0, 0, 0, 421,
		67, 1, 0, 0, 0, 422, 425, 5, 64, 0, 0, 423, 425, 3, 110, 55, 0, 424, 422,
		1, 0, 0, 0, 424, 423, 1, 0, 0, 0, 425, 69, 1, 0, 0, 0, 426, 434, 3, 72,
		36, 0, 427, 434, 3, 74, 37, 0, 428, 434, 3, 76, 38, 0, 429, 434, 3, 78,
		39, 0, 430, 434, 3, 80, 40, 0, 431, 434, 3, 82, 41, 0, 432, 434, 3, 84,
		42, 0, 433, 426, 1, 0, 0, 0, 433, 427, 1, 0, 0, 0, 433, 428, 1, 0, 0, 0,
		433, 429, 1, 0, 0, 0, 433, 430, 1, 0, 0, 0, 433, 431, 1, 0, 0, 0, 433,
		432, 1, 0, 0, 0, 434, 71, 1, 0, 0, 0, 435, 437, 5, 9, 0, 0, 436, 438, 3,
		106, 53, 0, 437, 436, 1, 0, 0, 0, 437, 438, 1, 0, 0, 0, 438, 439, 1, 0,
		0, 0, 439, 440, 5, 10, 0, 0, 440, 73, 1, 0, 0, 0, 441, 453, 5, 13, 0, 0,
		442, 447, 3, 86, 43, 0, 443, 444, 5, 8, 0, 0, 444, 446, 3, 86, 43, 0, 445,
		443, 1, 0, 0, 0, 446, 449, 1, 0, 0, 0, 447, 445, 1, 0, 0, 0, 447, 448,
		1, 0, 0, 0, 448, 451, 1, 0, 0, 0, 449, 447, 1, 0, 0, 0, 450, 452, 5, 8,
		0, 0, 451, 450, 1, 0, 0, 0, 451, 452, 1, 0, 0, 0, 452, 454, 1, 0, 0, 0,
		453, 442, 1, 0, 0, 0, 453, 454, 1, 0, 0, 0, 454, 455, 1, 0, 0, 0, 455,
		456, 5, 14, 0, 0, 456, 75, 1, 0, 0, 0, 457, 458, 5, 48, 0, 0, 458, 77,
		1, 0, 0, 0, 459, 460, 5, 66, 0, 0, 460, 79, 1, 0, 0, 0, 461, 462, 5, 68,
		0, 0, 462, 81, 1, 0, 0, 0, 463, 464, 5, 67, 0, 0, 464, 83, 1, 0, 0, 0,
		465, 466, 7, 1, 0, 0, 466, 85, 1, 0, 0, 0, 467, 468, 3, 90, 45, 0, 468,
		469, 5, 5, 0, 0, 469, 470, 3, 118, 59, 0, 470, 477, 1, 0, 0, 0, 471, 472,
		3, 88, 44, 0, 472, 473, 5, 5, 0, 0, 473, 474, 3, 118, 59, 0, 474, 477,
		1, 0, 0, 0, 475, 477, 3, 68, 34, 0, 476, 467, 1, 0, 0, 0, 476, 471, 1,
		0, 0, 0, 476, 475, 1, 0, 0, 0, 477, 87, 1, 0, 0, 0, 478, 479, 5, 9, 0,
		0, 479, 480, 3, 118, 59, 0, 480, 481, 5, 10, 0, 0, 481, 89, 1, 0, 0, 0,
		482, 488, 5, 64, 0, 0, 483, 488, 3, 78, 39, 0, 484, 488, 3, 66, 33, 0,
		485, 488, 3, 110, 55, 0, 486, 488, 3, 112, 56, 0, 487, 482, 1, 0, 0, 0,
		487, 483, 1, 0, 0, 0, 487, 484, 1, 0, 0, 0, 487, 485, 1, 0, 0, 0, 487,
		486, 1, 0, 0, 0, 488, 91, 1, 0, 0, 0, 489, 490, 3, 94, 47, 0, 490, 491,
		5, 64, 0, 0, 491, 93, 1, 0, 0, 0, 492, 494, 5, 69, 0, 0, 493, 492, 1, 0,
		0, 0, 494, 497, 1, 0, 0, 0, 495, 493, 1, 0, 0, 0, 495, 496, 1, 0, 0, 0,
		496, 95, 1, 0, 0, 0, 497, 495, 1, 0, 0, 0, 498, 500, 3, 98, 49, 0, 499,
		501, 3, 108, 54, 0, 500, 499, 1, 0, 0, 0, 501, 502, 1, 0, 0, 0, 502, 500,
		1, 0, 0, 0, 502, 503, 1, 0, 0, 0, 503, 97, 1, 0, 0, 0, 504, 510, 3, 68,
		34, 0, 505, 510, 3, 66, 33, 0, 506, 510, 3, 72, 36, 0, 507, 510, 3, 74,
		37, 0, 508, 510, 3, 102, 51, 0, 509, 504, 1, 0, 0, 0, 509, 505, 1, 0, 0,
		0, 509, 506, 1, 0, 0, 0, 509, 507, 1, 0, 0, 0, 509, 508, 1, 0, 0, 0, 510,
		99, 1, 0, 0, 0, 511, 513, 3, 102, 51, 0, 512, 514, 3, 144, 72, 0, 513,
		512, 1, 0, 0, 0, 513, 514, 1, 0, 0, 0, 514, 101, 1, 0, 0, 0, 515, 516,
		3, 94, 47, 0, 516, 517, 3, 104, 52, 0, 517, 519, 5, 11, 0, 0, 518, 520,
		3, 106, 53, 0, 519, 518, 1, 0, 0, 0, 519, 520, 1, 0, 0, 0, 520, 521, 1,
		0, 0, 0, 521, 522, 5, 12, 0, 0, 522, 103, 1, 0, 0, 0, 523, 527, 5, 64,
		0, 0, 524, 527, 3, 110, 55, 0, 525, 527, 3, 112, 56, 0, 526, 523, 1, 0,
		0, 0, 526, 524, 1, 0, 0, 0, 526, 525, 1, 0, 0, 0, 527, 105, 1, 0, 0, 0,
		528, 533, 3, 118, 59, 0, 529, 530, 5, 8, 0, 0, 530, 532, 3, 118, 59, 0,
		531, 529, 1, 0, 0, 0, 532, 535, 1, 0, 0, 0, 533, 531, 1, 0, 0, 0, 533,
		534, 1, 0, 0, 0, 534, 537, 1, 0, 0, 0, 535, 533, 1, 0, 0, 0, 536, 538,
		5, 8, 0, 0, 537, 536, 1, 0, 0, 0, 537, 538, 1, 0, 0, 0, 538, 107, 1, 0,
		0, 0, 539, 541, 3, 144, 72, 0, 540, 539, 1, 0, 0, 0, 540, 541, 1, 0, 0,
		0, 541, 542, 1, 0, 0, 0, 542, 543, 5, 7, 0, 0, 543, 551, 3, 90, 45, 0,
		544, 545, 3, 144, 72, 0, 545, 546, 5, 7, 0, 0, 546, 548, 1, 0, 0, 0, 547,
		544, 1, 0, 0, 0, 547, 548, 1, 0, 0, 0, 548, 549, 1, 0, 0, 0, 549, 551,
		3, 88, 44, 0, 550, 540, 1, 0, 0, 0, 550, 547, 1, 0, 0, 0, 551, 109, 1,
		0, 0, 0, 552, 553, 7, 2, 0, 0, 553, 111, 1, 0, 0, 0, 554, 555, 7, 3, 0,
		0, 555, 113, 1, 0, 0, 0, 556, 557, 3, 116, 58, 0, 557, 558, 5, 28, 0, 0,
		558, 559, 3, 116, 58, 0, 559, 115, 1, 0, 0, 0, 560, 564, 3, 82, 41, 0,
		561, 564, 3, 68, 34, 0, 562, 564, 3, 66, 33, 0, 563, 560, 1, 0, 0, 0, 563,
		561, 1, 0, 0, 0, 563, 562, 1, 0, 0, 0, 564, 117, 1, 0, 0, 0, 565, 566,
		6, 59, -1, 0, 566, 567, 3, 132, 66, 0, 567, 568, 3, 118, 59, 5, 568, 571,
		1, 0, 0, 0, 569, 571, 3, 120, 60, 0, 570, 565, 1, 0, 0, 0, 570, 569, 1,
		0, 0, 0, 571, 589, 1, 0, 0, 0, 572, 573, 10, 4, 0, 0, 573, 574, 3, 136,
		68, 0, 574, 575, 3, 118, 59, 5, 575, 588, 1, 0, 0, 0, 576, 577, 10, 3,
		0, 0, 577, 578, 3, 138, 69, 0, 578, 579, 3, 118, 59, 4, 579, 588, 1, 0,
		0, 0, 580, 581, 10, 2, 0, 0, 581, 583, 5, 30, 0, 0, 582, 584, 3, 118, 59,
		0, 583, 582, 1, 0, 0, 0, 583, 584, 1, 0, 0, 0, 584, 585, 1, 0, 0, 0, 585,
		586, 5, 5, 0, 0, 586, 588, 3, 118, 59, 3, 587, 572, 1, 0, 0, 0, 587, 576,
		1, 0, 0, 0, 587, 580, 1, 0, 0, 0, 588, 591, 1, 0, 0, 0, 589, 587, 1, 0,
		0, 0, 589, 590, 1, 0, 0, 0, 590, 119, 1, 0, 0, 0, 591, 589, 1, 0, 0, 0,
		592, 593, 6, 60, -1, 0, 593, 594, 3, 122, 61, 0, 594, 613, 1, 0, 0, 0,
		595, 596, 10, 5, 0, 0, 596, 597, 3, 126, 63, 0, 597, 598, 3, 120, 60, 6,
		598, 612, 1, 0, 0, 0, 599, 600, 10, 4, 0, 0, 600, 601, 3, 124, 62, 0, 601,
		602, 3, 120, 60, 5, 602, 612, 1, 0, 0, 0, 603, 604, 10, 3, 0, 0, 604, 605,
		3, 128, 64, 0, 605, 606, 3, 120, 60, 4, 606, 612, 1, 0, 0, 0, 607, 608,
		10, 2, 0, 0, 608, 609, 3, 130, 65, 0, 609, 610, 3, 120, 60, 3, 610, 612,
		1, 0, 0, 0, 611, 595, 1, 0, 0, 0, 611, 599, 1, 0, 0, 0, 611, 603, 1, 0,
		0, 0, 611, 607, 1, 0, 0, 0, 612, 615, 1, 0, 0, 0, 613, 611, 1, 0, 0, 0,
		613, 614, 1, 0, 0, 0, 614, 121, 1, 0, 0, 0, 615, 613, 1, 0, 0, 0, 616,
		617, 6, 61, -1, 0, 617, 634, 3, 100, 50, 0, 618, 634, 3, 114, 57, 0, 619,
		634, 3, 70, 35, 0, 620, 634, 3, 68, 34, 0, 621, 634, 3, 96, 48, 0, 622,
		634, 3, 66, 33, 0, 623, 627, 5, 11, 0, 0, 624, 628, 3, 18, 9, 0, 625, 628,
		3, 56, 28, 0, 626, 628, 3, 118, 59, 0, 627, 624, 1, 0, 0, 0, 627, 625,
		1, 0, 0, 0, 627, 626, 1, 0, 0, 0, 628, 629, 1, 0, 0, 0, 629, 631, 5, 12,
		0, 0, 630, 632, 3, 144, 72, 0, 631, 630, 1, 0, 0, 0, 631, 632, 1, 0, 0,
		0, 632, 634, 1, 0, 0, 0, 633, 616, 1, 0, 0, 0, 633, 618, 1, 0, 0, 0, 633,
		619, 1, 0, 0, 0, 633, 620, 1, 0, 0, 0, 633, 621, 1, 0, 0, 0, 633, 622,
		1, 0, 0, 0, 633, 623, 1, 0, 0, 0, 634, 649, 1, 0, 0, 0, 635, 636, 10, 10,
		0, 0, 636, 637, 3, 140, 70, 0, 637, 638, 3, 122, 61, 11, 638, 648, 1, 0,
		0, 0, 639, 640, 10, 9, 0, 0, 640, 641, 3, 142, 71, 0, 641, 642, 3, 122,
		61, 10, 642, 648, 1, 0, 0, 0, 643, 644, 10, 8, 0, 0, 644, 645, 3, 134,
		67, 0, 645, 646, 3, 122, 61, 9, 646, 648, 1, 0, 0, 0, 647, 635, 1, 0, 0,
		0, 647, 639, 1, 0, 0, 0, 647, 643, 1, 0, 0, 0, 648, 651, 1, 0, 0, 0, 649,
		647, 1, 0, 0, 0, 649, 650, 1, 0, 0, 0, 650, 123, 1, 0, 0, 0, 651, 649,
		1, 0, 0, 0, 652, 655, 7, 4, 0, 0, 653, 656, 3, 128, 64, 0, 654, 656, 3,
		126, 63, 0, 655, 653, 1, 0, 0, 0, 655, 654, 1, 0, 0, 0, 656, 125, 1, 0,
		0, 0, 657, 658, 7, 5, 0, 0, 658, 127, 1, 0, 0, 0, 659, 661, 5, 58, 0, 0,
		660, 659, 1, 0, 0, 0, 660, 661, 1, 0, 0, 0, 661, 662, 1, 0, 0, 0, 662,
		663, 5, 59, 0, 0, 663, 129, 1, 0, 0, 0, 664, 666, 5, 58, 0, 0, 665, 664,
		1, 0, 0, 0, 665, 666, 1, 0, 0, 0, 666, 667, 1, 0, 0, 0, 667, 668, 5, 57,
		0, 0, 668, 131, 1, 0, 0, 0, 669, 670, 7, 6, 0, 0, 670, 133, 1, 0, 0, 0,
		671, 672, 7, 7, 0, 0, 672, 135, 1, 0, 0, 0, 673, 674, 5, 26, 0, 0, 674,
		137, 1, 0, 0, 0, 675, 676, 5, 27, 0, 0, 676, 139, 1, 0, 0, 0, 677, 678,
		7, 8, 0, 0, 678, 141, 1, 0, 0, 0, 679, 680, 7, 9, 0, 0, 680, 143, 1, 0,
		0, 0, 681, 682, 5, 30, 0, 0, 682, 145, 1, 0, 0, 0, 68, 149, 164, 172, 176,
		187, 191, 199, 206, 224, 232, 239, 244, 253, 259, 263, 267, 271, 280, 287,
		295, 300, 329, 340, 353, 361, 363, 371, 385, 388, 391, 398, 403, 414, 420,
		424, 433, 437, 447, 451, 453, 476, 487, 495, 502, 509, 513, 519, 526, 533,
		537, 540, 547, 550, 563, 570, 583, 587, 589, 611, 613, 627, 631, 633, 647,
		649, 655, 660, 665,
	}
	deserializer := antlr.NewATNDeserializer(nil)
	staticData.atn = deserializer.Deserialize(staticData.serializedATN)
	atn := staticData.atn
	staticData.decisionToDFA = make([]*antlr.DFA, len(atn.DecisionToState))
	decisionToDFA := staticData.decisionToDFA
	for index, state := range atn.DecisionToState {
		decisionToDFA[index] = antlr.NewDFA(state, index)
	}
}

// FqlParserInit initializes any static state used to implement FqlParser. By default the
// static state used to implement the parser is lazily initialized during the first call to
// NewFqlParser(). You can call this function if you wish to initialize the static state ahead
// of time.
func FqlParserInit() {
	staticData := &FqlParserParserStaticData
	staticData.once.Do(fqlparserParserInit)
}

// NewFqlParser produces a new parser instance for the optional input antlr.TokenStream.
func NewFqlParser(input antlr.TokenStream) *FqlParser {
	FqlParserInit()
	this := new(FqlParser)
	this.BaseParser = antlr.NewBaseParser(input)
	staticData := &FqlParserParserStaticData
	this.Interpreter = antlr.NewParserATNSimulator(this, staticData.atn, staticData.decisionToDFA, staticData.PredictionContextCache)
	this.RuleNames = staticData.RuleNames
	this.LiteralNames = staticData.LiteralNames
	this.SymbolicNames = staticData.SymbolicNames
	this.GrammarFileName = "FqlParser.g4"

	return this
}

// FqlParser tokens.
const (
	FqlParserEOF               = antlr.TokenEOF
	FqlParserMultiLineComment  = 1
	FqlParserSingleLineComment = 2
	FqlParserWhiteSpaces       = 3
	FqlParserLineTerminator    = 4
	FqlParserColon             = 5
	FqlParserSemiColon         = 6
	FqlParserDot               = 7
	FqlParserComma             = 8
	FqlParserOpenBracket       = 9
	FqlParserCloseBracket      = 10
	FqlParserOpenParen         = 11
	FqlParserCloseParen        = 12
	FqlParserOpenBrace         = 13
	FqlParserCloseBrace        = 14
	FqlParserGt                = 15
	FqlParserLt                = 16
	FqlParserEq                = 17
	FqlParserGte               = 18
	FqlParserLte               = 19
	FqlParserNeq               = 20
	FqlParserMulti             = 21
	FqlParserDiv               = 22
	FqlParserMod               = 23
	FqlParserPlus              = 24
	FqlParserMinus             = 25
	FqlParserAnd               = 26
	FqlParserOr                = 27
	FqlParserRange             = 28
	FqlParserAssign            = 29
	FqlParserQuestionMark      = 30
	FqlParserRegexNotMatch     = 31
	FqlParserRegexMatch        = 32
	FqlParserFor               = 33
	FqlParserReturn            = 34
	FqlParserWaitfor           = 35
	FqlParserOptions           = 36
	FqlParserTimeout           = 37
	FqlParserDistinct          = 38
	FqlParserFilter            = 39
	FqlParserCurrent           = 40
	FqlParserSort              = 41
	FqlParserLimit             = 42
	FqlParserLet               = 43
	FqlParserCollect           = 44
	FqlParserSortDirection     = 45
	FqlParserNone              = 46
	FqlParserNull              = 47
	FqlParserBooleanLiteral    = 48
	FqlParserUse               = 49
	FqlParserInto              = 50
	FqlParserKeep              = 51
	FqlParserWith              = 52
	FqlParserAll               = 53
	FqlParserAny               = 54
	FqlParserAggregate         = 55
	FqlParserEvent             = 56
	FqlParserLike              = 57
	FqlParserNot               = 58
	FqlParserIn                = 59
	FqlParserDo                = 60
	FqlParserWhile             = 61
	FqlParserStep              = 62
	FqlParserParam             = 63
	FqlParserIdentifier        = 64
	FqlParserIgnoreIdentifier  = 65
	FqlParserStringLiteral     = 66
	FqlParserIntegerLiteral    = 67
	FqlParserFloatLiteral      = 68
	FqlParserNamespaceSegment  = 69
	FqlParserUnknownIdentifier = 70
)

// FqlParser rules.
const (
	FqlParserRULE_program                      = 0
	FqlParserRULE_head                         = 1
	FqlParserRULE_useExpression                = 2
	FqlParserRULE_use                          = 3
	FqlParserRULE_body                         = 4
	FqlParserRULE_bodyStatement                = 5
	FqlParserRULE_bodyExpression               = 6
	FqlParserRULE_variableDeclaration          = 7
	FqlParserRULE_returnExpression             = 8
	FqlParserRULE_forExpression                = 9
	FqlParserRULE_forExpressionSource          = 10
	FqlParserRULE_forExpressionClause          = 11
	FqlParserRULE_forExpressionStatement       = 12
	FqlParserRULE_forExpressionBody            = 13
	FqlParserRULE_forExpressionReturn          = 14
	FqlParserRULE_filterClause                 = 15
	FqlParserRULE_limitClause                  = 16
	FqlParserRULE_limitClauseValue             = 17
	FqlParserRULE_sortClause                   = 18
	FqlParserRULE_sortClauseExpression         = 19
	FqlParserRULE_collectClause                = 20
	FqlParserRULE_collectSelector              = 21
	FqlParserRULE_collectGrouping              = 22
	FqlParserRULE_collectAggregateSelector     = 23
	FqlParserRULE_collectAggregator            = 24
	FqlParserRULE_collectGroupProjection       = 25
	FqlParserRULE_collectGroupProjectionFilter = 26
	FqlParserRULE_collectCounter               = 27
	FqlParserRULE_waitForExpression            = 28
	FqlParserRULE_waitForEventName             = 29
	FqlParserRULE_waitForEventSource           = 30
	FqlParserRULE_optionsClause                = 31
	FqlParserRULE_timeoutClause                = 32
	FqlParserRULE_param                        = 33
	FqlParserRULE_variable                     = 34
	FqlParserRULE_literal                      = 35
	FqlParserRULE_arrayLiteral                 = 36
	FqlParserRULE_objectLiteral                = 37
	FqlParserRULE_booleanLiteral               = 38
	FqlParserRULE_stringLiteral                = 39
	FqlParserRULE_floatLiteral                 = 40
	FqlParserRULE_integerLiteral               = 41
	FqlParserRULE_noneLiteral                  = 42
	FqlParserRULE_propertyAssignment           = 43
	FqlParserRULE_computedPropertyName         = 44
	FqlParserRULE_propertyName                 = 45
	FqlParserRULE_namespaceIdentifier          = 46
	FqlParserRULE_namespace                    = 47
	FqlParserRULE_memberExpression             = 48
	FqlParserRULE_memberExpressionSource       = 49
	FqlParserRULE_functionCallExpression       = 50
	FqlParserRULE_functionCall                 = 51
	FqlParserRULE_functionName                 = 52
	FqlParserRULE_argumentList                 = 53
	FqlParserRULE_memberExpressionPath         = 54
	FqlParserRULE_safeReservedWord             = 55
	FqlParserRULE_unsafeReservedWord           = 56
	FqlParserRULE_rangeOperator                = 57
	FqlParserRULE_rangeOperand                 = 58
	FqlParserRULE_expression                   = 59
	FqlParserRULE_predicate                    = 60
	FqlParserRULE_expressionAtom               = 61
	FqlParserRULE_arrayOperator                = 62
	FqlParserRULE_equalityOperator             = 63
	FqlParserRULE_inOperator                   = 64
	FqlParserRULE_likeOperator                 = 65
	FqlParserRULE_unaryOperator                = 66
	FqlParserRULE_regexpOperator               = 67
	FqlParserRULE_logicalAndOperator           = 68
	FqlParserRULE_logicalOrOperator            = 69
	FqlParserRULE_multiplicativeOperator       = 70
	FqlParserRULE_additiveOperator             = 71
	FqlParserRULE_errorOperator                = 72
)

// IProgramContext is an interface to support dynamic dispatch.
type IProgramContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	Body() IBodyContext
	AllHead() []IHeadContext
	Head(i int) IHeadContext

	// IsProgramContext differentiates from other interfaces.
	IsProgramContext()
}

type ProgramContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyProgramContext() *ProgramContext {
	var p = new(ProgramContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = FqlParserRULE_program
	return p
}

func InitEmptyProgramContext(p *ProgramContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = FqlParserRULE_program
}

func (*ProgramContext) IsProgramContext() {}

func NewProgramContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ProgramContext {
	var p = new(ProgramContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = FqlParserRULE_program

	return p
}

func (s *ProgramContext) GetParser() antlr.Parser { return s.parser }

func (s *ProgramContext) Body() IBodyContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IBodyContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IBodyContext)
}

func (s *ProgramContext) AllHead() []IHeadContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(IHeadContext); ok {
			len++
		}
	}

	tst := make([]IHeadContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(IHeadContext); ok {
			tst[i] = t.(IHeadContext)
			i++
		}
	}

	return tst
}

func (s *ProgramContext) Head(i int) IHeadContext {
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IHeadContext); ok {
			if j == i {
				t = ctx.(antlr.RuleContext)
				break
			}
			j++
		}
	}

	if t == nil {
		return nil
	}

	return t.(IHeadContext)
}

func (s *ProgramContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ProgramContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *ProgramContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(FqlParserListener); ok {
		listenerT.EnterProgram(s)
	}
}

func (s *ProgramContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(FqlParserListener); ok {
		listenerT.ExitProgram(s)
	}
}

func (s *ProgramContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case FqlParserVisitor:
		return t.VisitProgram(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *FqlParser) Program() (localctx IProgramContext) {
	localctx = NewProgramContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 0, FqlParserRULE_program)
	var _alt int

	p.EnterOuterAlt(localctx, 1)
	p.SetState(149)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_alt = p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 0, p.GetParserRuleContext())
	if p.HasError() {
		goto errorExit
	}
	for _alt != 2 && _alt != antlr.ATNInvalidAltNumber {
		if _alt == 1 {
			{
				p.SetState(146)
				p.Head()
			}

		}
		p.SetState(151)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_alt = p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 0, p.GetParserRuleContext())
		if p.HasError() {
			goto errorExit
		}
	}
	{
		p.SetState(152)
		p.Body()
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IHeadContext is an interface to support dynamic dispatch.
type IHeadContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	UseExpression() IUseExpressionContext

	// IsHeadContext differentiates from other interfaces.
	IsHeadContext()
}

type HeadContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyHeadContext() *HeadContext {
	var p = new(HeadContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = FqlParserRULE_head
	return p
}

func InitEmptyHeadContext(p *HeadContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = FqlParserRULE_head
}

func (*HeadContext) IsHeadContext() {}

func NewHeadContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *HeadContext {
	var p = new(HeadContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = FqlParserRULE_head

	return p
}

func (s *HeadContext) GetParser() antlr.Parser { return s.parser }

func (s *HeadContext) UseExpression() IUseExpressionContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IUseExpressionContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IUseExpressionContext)
}

func (s *HeadContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *HeadContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *HeadContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(FqlParserListener); ok {
		listenerT.EnterHead(s)
	}
}

func (s *HeadContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(FqlParserListener); ok {
		listenerT.ExitHead(s)
	}
}

func (s *HeadContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case FqlParserVisitor:
		return t.VisitHead(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *FqlParser) Head() (localctx IHeadContext) {
	localctx = NewHeadContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 2, FqlParserRULE_head)
	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(154)
		p.UseExpression()
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IUseExpressionContext is an interface to support dynamic dispatch.
type IUseExpressionContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	Use() IUseContext

	// IsUseExpressionContext differentiates from other interfaces.
	IsUseExpressionContext()
}

type UseExpressionContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyUseExpressionContext() *UseExpressionContext {
	var p = new(UseExpressionContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = FqlParserRULE_useExpression
	return p
}

func InitEmptyUseExpressionContext(p *UseExpressionContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = FqlParserRULE_useExpression
}

func (*UseExpressionContext) IsUseExpressionContext() {}

func NewUseExpressionContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *UseExpressionContext {
	var p = new(UseExpressionContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = FqlParserRULE_useExpression

	return p
}

func (s *UseExpressionContext) GetParser() antlr.Parser { return s.parser }

func (s *UseExpressionContext) Use() IUseContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IUseContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IUseContext)
}

func (s *UseExpressionContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *UseExpressionContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *UseExpressionContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(FqlParserListener); ok {
		listenerT.EnterUseExpression(s)
	}
}

func (s *UseExpressionContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(FqlParserListener); ok {
		listenerT.ExitUseExpression(s)
	}
}

func (s *UseExpressionContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case FqlParserVisitor:
		return t.VisitUseExpression(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *FqlParser) UseExpression() (localctx IUseExpressionContext) {
	localctx = NewUseExpressionContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 4, FqlParserRULE_useExpression)
	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(156)
		p.Use()
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IUseContext is an interface to support dynamic dispatch.
type IUseContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	Use() antlr.TerminalNode
	NamespaceIdentifier() INamespaceIdentifierContext

	// IsUseContext differentiates from other interfaces.
	IsUseContext()
}

type UseContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyUseContext() *UseContext {
	var p = new(UseContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = FqlParserRULE_use
	return p
}

func InitEmptyUseContext(p *UseContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = FqlParserRULE_use
}

func (*UseContext) IsUseContext() {}

func NewUseContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *UseContext {
	var p = new(UseContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = FqlParserRULE_use

	return p
}

func (s *UseContext) GetParser() antlr.Parser { return s.parser }

func (s *UseContext) Use() antlr.TerminalNode {
	return s.GetToken(FqlParserUse, 0)
}

func (s *UseContext) NamespaceIdentifier() INamespaceIdentifierContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(INamespaceIdentifierContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(INamespaceIdentifierContext)
}

func (s *UseContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *UseContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *UseContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(FqlParserListener); ok {
		listenerT.EnterUse(s)
	}
}

func (s *UseContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(FqlParserListener); ok {
		listenerT.ExitUse(s)
	}
}

func (s *UseContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case FqlParserVisitor:
		return t.VisitUse(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *FqlParser) Use() (localctx IUseContext) {
	localctx = NewUseContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 6, FqlParserRULE_use)
	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(158)
		p.Match(FqlParserUse)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(159)
		p.NamespaceIdentifier()
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IBodyContext is an interface to support dynamic dispatch.
type IBodyContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	BodyExpression() IBodyExpressionContext
	AllBodyStatement() []IBodyStatementContext
	BodyStatement(i int) IBodyStatementContext

	// IsBodyContext differentiates from other interfaces.
	IsBodyContext()
}

type BodyContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyBodyContext() *BodyContext {
	var p = new(BodyContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = FqlParserRULE_body
	return p
}

func InitEmptyBodyContext(p *BodyContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = FqlParserRULE_body
}

func (*BodyContext) IsBodyContext() {}

func NewBodyContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *BodyContext {
	var p = new(BodyContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = FqlParserRULE_body

	return p
}

func (s *BodyContext) GetParser() antlr.Parser { return s.parser }

func (s *BodyContext) BodyExpression() IBodyExpressionContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IBodyExpressionContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IBodyExpressionContext)
}

func (s *BodyContext) AllBodyStatement() []IBodyStatementContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(IBodyStatementContext); ok {
			len++
		}
	}

	tst := make([]IBodyStatementContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(IBodyStatementContext); ok {
			tst[i] = t.(IBodyStatementContext)
			i++
		}
	}

	return tst
}

func (s *BodyContext) BodyStatement(i int) IBodyStatementContext {
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IBodyStatementContext); ok {
			if j == i {
				t = ctx.(antlr.RuleContext)
				break
			}
			j++
		}
	}

	if t == nil {
		return nil
	}

	return t.(IBodyStatementContext)
}

func (s *BodyContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *BodyContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *BodyContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(FqlParserListener); ok {
		listenerT.EnterBody(s)
	}
}

func (s *BodyContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(FqlParserListener); ok {
		listenerT.ExitBody(s)
	}
}

func (s *BodyContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case FqlParserVisitor:
		return t.VisitBody(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *FqlParser) Body() (localctx IBodyContext) {
	localctx = NewBodyContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 8, FqlParserRULE_body)
	var _alt int

	p.EnterOuterAlt(localctx, 1)
	p.SetState(164)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_alt = p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 1, p.GetParserRuleContext())
	if p.HasError() {
		goto errorExit
	}
	for _alt != 2 && _alt != antlr.ATNInvalidAltNumber {
		if _alt == 1 {
			{
				p.SetState(161)
				p.BodyStatement()
			}

		}
		p.SetState(166)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_alt = p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 1, p.GetParserRuleContext())
		if p.HasError() {
			goto errorExit
		}
	}
	{
		p.SetState(167)
		p.BodyExpression()
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IBodyStatementContext is an interface to support dynamic dispatch.
type IBodyStatementContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	VariableDeclaration() IVariableDeclarationContext
	FunctionCallExpression() IFunctionCallExpressionContext
	WaitForExpression() IWaitForExpressionContext

	// IsBodyStatementContext differentiates from other interfaces.
	IsBodyStatementContext()
}

type BodyStatementContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyBodyStatementContext() *BodyStatementContext {
	var p = new(BodyStatementContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = FqlParserRULE_bodyStatement
	return p
}

func InitEmptyBodyStatementContext(p *BodyStatementContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = FqlParserRULE_bodyStatement
}

func (*BodyStatementContext) IsBodyStatementContext() {}

func NewBodyStatementContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *BodyStatementContext {
	var p = new(BodyStatementContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = FqlParserRULE_bodyStatement

	return p
}

func (s *BodyStatementContext) GetParser() antlr.Parser { return s.parser }

func (s *BodyStatementContext) VariableDeclaration() IVariableDeclarationContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IVariableDeclarationContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IVariableDeclarationContext)
}

func (s *BodyStatementContext) FunctionCallExpression() IFunctionCallExpressionContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IFunctionCallExpressionContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IFunctionCallExpressionContext)
}

func (s *BodyStatementContext) WaitForExpression() IWaitForExpressionContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IWaitForExpressionContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IWaitForExpressionContext)
}

func (s *BodyStatementContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *BodyStatementContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *BodyStatementContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(FqlParserListener); ok {
		listenerT.EnterBodyStatement(s)
	}
}

func (s *BodyStatementContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(FqlParserListener); ok {
		listenerT.ExitBodyStatement(s)
	}
}

func (s *BodyStatementContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case FqlParserVisitor:
		return t.VisitBodyStatement(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *FqlParser) BodyStatement() (localctx IBodyStatementContext) {
	localctx = NewBodyStatementContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 10, FqlParserRULE_bodyStatement)
	p.SetState(172)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 2, p.GetParserRuleContext()) {
	case 1:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(169)
			p.VariableDeclaration()
		}

	case 2:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(170)
			p.FunctionCallExpression()
		}

	case 3:
		p.EnterOuterAlt(localctx, 3)
		{
			p.SetState(171)
			p.WaitForExpression()
		}

	case antlr.ATNInvalidAltNumber:
		goto errorExit
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IBodyExpressionContext is an interface to support dynamic dispatch.
type IBodyExpressionContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	ReturnExpression() IReturnExpressionContext
	ForExpression() IForExpressionContext

	// IsBodyExpressionContext differentiates from other interfaces.
	IsBodyExpressionContext()
}

type BodyExpressionContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyBodyExpressionContext() *BodyExpressionContext {
	var p = new(BodyExpressionContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = FqlParserRULE_bodyExpression
	return p
}

func InitEmptyBodyExpressionContext(p *BodyExpressionContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = FqlParserRULE_bodyExpression
}

func (*BodyExpressionContext) IsBodyExpressionContext() {}

func NewBodyExpressionContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *BodyExpressionContext {
	var p = new(BodyExpressionContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = FqlParserRULE_bodyExpression

	return p
}

func (s *BodyExpressionContext) GetParser() antlr.Parser { return s.parser }

func (s *BodyExpressionContext) ReturnExpression() IReturnExpressionContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IReturnExpressionContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IReturnExpressionContext)
}

func (s *BodyExpressionContext) ForExpression() IForExpressionContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IForExpressionContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IForExpressionContext)
}

func (s *BodyExpressionContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *BodyExpressionContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *BodyExpressionContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(FqlParserListener); ok {
		listenerT.EnterBodyExpression(s)
	}
}

func (s *BodyExpressionContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(FqlParserListener); ok {
		listenerT.ExitBodyExpression(s)
	}
}

func (s *BodyExpressionContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case FqlParserVisitor:
		return t.VisitBodyExpression(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *FqlParser) BodyExpression() (localctx IBodyExpressionContext) {
	localctx = NewBodyExpressionContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 12, FqlParserRULE_bodyExpression)
	p.SetState(176)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetTokenStream().LA(1) {
	case FqlParserReturn:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(174)
			p.ReturnExpression()
		}

	case FqlParserFor:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(175)
			p.ForExpression()
		}

	default:
		p.SetError(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
		goto errorExit
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IVariableDeclarationContext is an interface to support dynamic dispatch.
type IVariableDeclarationContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// GetId returns the id token.
	GetId() antlr.Token

	// SetId sets the id token.
	SetId(antlr.Token)

	// Getter signatures
	Let() antlr.TerminalNode
	Assign() antlr.TerminalNode
	Expression() IExpressionContext
	Identifier() antlr.TerminalNode
	IgnoreIdentifier() antlr.TerminalNode
	SafeReservedWord() ISafeReservedWordContext

	// IsVariableDeclarationContext differentiates from other interfaces.
	IsVariableDeclarationContext()
}

type VariableDeclarationContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
	id     antlr.Token
}

func NewEmptyVariableDeclarationContext() *VariableDeclarationContext {
	var p = new(VariableDeclarationContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = FqlParserRULE_variableDeclaration
	return p
}

func InitEmptyVariableDeclarationContext(p *VariableDeclarationContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = FqlParserRULE_variableDeclaration
}

func (*VariableDeclarationContext) IsVariableDeclarationContext() {}

func NewVariableDeclarationContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *VariableDeclarationContext {
	var p = new(VariableDeclarationContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = FqlParserRULE_variableDeclaration

	return p
}

func (s *VariableDeclarationContext) GetParser() antlr.Parser { return s.parser }

func (s *VariableDeclarationContext) GetId() antlr.Token { return s.id }

func (s *VariableDeclarationContext) SetId(v antlr.Token) { s.id = v }

func (s *VariableDeclarationContext) Let() antlr.TerminalNode {
	return s.GetToken(FqlParserLet, 0)
}

func (s *VariableDeclarationContext) Assign() antlr.TerminalNode {
	return s.GetToken(FqlParserAssign, 0)
}

func (s *VariableDeclarationContext) Expression() IExpressionContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IExpressionContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IExpressionContext)
}

func (s *VariableDeclarationContext) Identifier() antlr.TerminalNode {
	return s.GetToken(FqlParserIdentifier, 0)
}

func (s *VariableDeclarationContext) IgnoreIdentifier() antlr.TerminalNode {
	return s.GetToken(FqlParserIgnoreIdentifier, 0)
}

func (s *VariableDeclarationContext) SafeReservedWord() ISafeReservedWordContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ISafeReservedWordContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(ISafeReservedWordContext)
}

func (s *VariableDeclarationContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *VariableDeclarationContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *VariableDeclarationContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(FqlParserListener); ok {
		listenerT.EnterVariableDeclaration(s)
	}
}

func (s *VariableDeclarationContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(FqlParserListener); ok {
		listenerT.ExitVariableDeclaration(s)
	}
}

func (s *VariableDeclarationContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case FqlParserVisitor:
		return t.VisitVariableDeclaration(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *FqlParser) VariableDeclaration() (localctx IVariableDeclarationContext) {
	localctx = NewVariableDeclarationContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 14, FqlParserRULE_variableDeclaration)
	var _la int

	p.SetState(187)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 4, p.GetParserRuleContext()) {
	case 1:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(178)
			p.Match(FqlParserLet)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(179)

			var _lt = p.GetTokenStream().LT(1)

			localctx.(*VariableDeclarationContext).id = _lt

			_la = p.GetTokenStream().LA(1)

			if !(_la == FqlParserIdentifier || _la == FqlParserIgnoreIdentifier) {
				var _ri = p.GetErrorHandler().RecoverInline(p)

				localctx.(*VariableDeclarationContext).id = _ri
			} else {
				p.GetErrorHandler().ReportMatch(p)
				p.Consume()
			}
		}
		{
			p.SetState(180)
			p.Match(FqlParserAssign)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(181)
			p.expression(0)
		}

	case 2:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(182)
			p.Match(FqlParserLet)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(183)
			p.SafeReservedWord()
		}
		{
			p.SetState(184)
			p.Match(FqlParserAssign)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(185)
			p.expression(0)
		}

	case antlr.ATNInvalidAltNumber:
		goto errorExit
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IReturnExpressionContext is an interface to support dynamic dispatch.
type IReturnExpressionContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	Return() antlr.TerminalNode
	Expression() IExpressionContext
	Distinct() antlr.TerminalNode

	// IsReturnExpressionContext differentiates from other interfaces.
	IsReturnExpressionContext()
}

type ReturnExpressionContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyReturnExpressionContext() *ReturnExpressionContext {
	var p = new(ReturnExpressionContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = FqlParserRULE_returnExpression
	return p
}

func InitEmptyReturnExpressionContext(p *ReturnExpressionContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = FqlParserRULE_returnExpression
}

func (*ReturnExpressionContext) IsReturnExpressionContext() {}

func NewReturnExpressionContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ReturnExpressionContext {
	var p = new(ReturnExpressionContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = FqlParserRULE_returnExpression

	return p
}

func (s *ReturnExpressionContext) GetParser() antlr.Parser { return s.parser }

func (s *ReturnExpressionContext) Return() antlr.TerminalNode {
	return s.GetToken(FqlParserReturn, 0)
}

func (s *ReturnExpressionContext) Expression() IExpressionContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IExpressionContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IExpressionContext)
}

func (s *ReturnExpressionContext) Distinct() antlr.TerminalNode {
	return s.GetToken(FqlParserDistinct, 0)
}

func (s *ReturnExpressionContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ReturnExpressionContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *ReturnExpressionContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(FqlParserListener); ok {
		listenerT.EnterReturnExpression(s)
	}
}

func (s *ReturnExpressionContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(FqlParserListener); ok {
		listenerT.ExitReturnExpression(s)
	}
}

func (s *ReturnExpressionContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case FqlParserVisitor:
		return t.VisitReturnExpression(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *FqlParser) ReturnExpression() (localctx IReturnExpressionContext) {
	localctx = NewReturnExpressionContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 16, FqlParserRULE_returnExpression)
	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(189)
		p.Match(FqlParserReturn)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	p.SetState(191)
	p.GetErrorHandler().Sync(p)

	if p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 5, p.GetParserRuleContext()) == 1 {
		{
			p.SetState(190)
			p.Match(FqlParserDistinct)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	} else if p.HasError() { // JIM
		goto errorExit
	}
	{
		p.SetState(193)
		p.expression(0)
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IForExpressionContext is an interface to support dynamic dispatch.
type IForExpressionContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// GetValueVariable returns the valueVariable token.
	GetValueVariable() antlr.Token

	// GetCounterVariable returns the counterVariable token.
	GetCounterVariable() antlr.Token

	// GetStepVariable returns the stepVariable token.
	GetStepVariable() antlr.Token

	// SetValueVariable sets the valueVariable token.
	SetValueVariable(antlr.Token)

	// SetCounterVariable sets the counterVariable token.
	SetCounterVariable(antlr.Token)

	// SetStepVariable sets the stepVariable token.
	SetStepVariable(antlr.Token)

	// GetStepInit returns the stepInit rule contexts.
	GetStepInit() IExpressionContext

	// GetStepCondition returns the stepCondition rule contexts.
	GetStepCondition() IExpressionContext

	// GetStepIncrement returns the stepIncrement rule contexts.
	GetStepIncrement() IExpressionContext

	// SetStepInit sets the stepInit rule contexts.
	SetStepInit(IExpressionContext)

	// SetStepCondition sets the stepCondition rule contexts.
	SetStepCondition(IExpressionContext)

	// SetStepIncrement sets the stepIncrement rule contexts.
	SetStepIncrement(IExpressionContext)

	// Getter signatures
	For() antlr.TerminalNode
	In() antlr.TerminalNode
	ForExpressionSource() IForExpressionSourceContext
	ForExpressionReturn() IForExpressionReturnContext
	AllIdentifier() []antlr.TerminalNode
	Identifier(i int) antlr.TerminalNode
	AllIgnoreIdentifier() []antlr.TerminalNode
	IgnoreIdentifier(i int) antlr.TerminalNode
	Comma() antlr.TerminalNode
	AllForExpressionBody() []IForExpressionBodyContext
	ForExpressionBody(i int) IForExpressionBodyContext
	AllAssign() []antlr.TerminalNode
	Assign(i int) antlr.TerminalNode
	While() antlr.TerminalNode
	Step() antlr.TerminalNode
	AllExpression() []IExpressionContext
	Expression(i int) IExpressionContext
	Do() antlr.TerminalNode

	// IsForExpressionContext differentiates from other interfaces.
	IsForExpressionContext()
}

type ForExpressionContext struct {
	antlr.BaseParserRuleContext
	parser          antlr.Parser
	valueVariable   antlr.Token
	counterVariable antlr.Token
	stepInit        IExpressionContext
	stepCondition   IExpressionContext
	stepVariable    antlr.Token
	stepIncrement   IExpressionContext
}

func NewEmptyForExpressionContext() *ForExpressionContext {
	var p = new(ForExpressionContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = FqlParserRULE_forExpression
	return p
}

func InitEmptyForExpressionContext(p *ForExpressionContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = FqlParserRULE_forExpression
}

func (*ForExpressionContext) IsForExpressionContext() {}

func NewForExpressionContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ForExpressionContext {
	var p = new(ForExpressionContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = FqlParserRULE_forExpression

	return p
}

func (s *ForExpressionContext) GetParser() antlr.Parser { return s.parser }

func (s *ForExpressionContext) GetValueVariable() antlr.Token { return s.valueVariable }

func (s *ForExpressionContext) GetCounterVariable() antlr.Token { return s.counterVariable }

func (s *ForExpressionContext) GetStepVariable() antlr.Token { return s.stepVariable }

func (s *ForExpressionContext) SetValueVariable(v antlr.Token) { s.valueVariable = v }

func (s *ForExpressionContext) SetCounterVariable(v antlr.Token) { s.counterVariable = v }

func (s *ForExpressionContext) SetStepVariable(v antlr.Token) { s.stepVariable = v }

func (s *ForExpressionContext) GetStepInit() IExpressionContext { return s.stepInit }

func (s *ForExpressionContext) GetStepCondition() IExpressionContext { return s.stepCondition }

func (s *ForExpressionContext) GetStepIncrement() IExpressionContext { return s.stepIncrement }

func (s *ForExpressionContext) SetStepInit(v IExpressionContext) { s.stepInit = v }

func (s *ForExpressionContext) SetStepCondition(v IExpressionContext) { s.stepCondition = v }

func (s *ForExpressionContext) SetStepIncrement(v IExpressionContext) { s.stepIncrement = v }

func (s *ForExpressionContext) For() antlr.TerminalNode {
	return s.GetToken(FqlParserFor, 0)
}

func (s *ForExpressionContext) In() antlr.TerminalNode {
	return s.GetToken(FqlParserIn, 0)
}

func (s *ForExpressionContext) ForExpressionSource() IForExpressionSourceContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IForExpressionSourceContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IForExpressionSourceContext)
}

func (s *ForExpressionContext) ForExpressionReturn() IForExpressionReturnContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IForExpressionReturnContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IForExpressionReturnContext)
}

func (s *ForExpressionContext) AllIdentifier() []antlr.TerminalNode {
	return s.GetTokens(FqlParserIdentifier)
}

func (s *ForExpressionContext) Identifier(i int) antlr.TerminalNode {
	return s.GetToken(FqlParserIdentifier, i)
}

func (s *ForExpressionContext) AllIgnoreIdentifier() []antlr.TerminalNode {
	return s.GetTokens(FqlParserIgnoreIdentifier)
}

func (s *ForExpressionContext) IgnoreIdentifier(i int) antlr.TerminalNode {
	return s.GetToken(FqlParserIgnoreIdentifier, i)
}

func (s *ForExpressionContext) Comma() antlr.TerminalNode {
	return s.GetToken(FqlParserComma, 0)
}

func (s *ForExpressionContext) AllForExpressionBody() []IForExpressionBodyContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(IForExpressionBodyContext); ok {
			len++
		}
	}

	tst := make([]IForExpressionBodyContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(IForExpressionBodyContext); ok {
			tst[i] = t.(IForExpressionBodyContext)
			i++
		}
	}

	return tst
}

func (s *ForExpressionContext) ForExpressionBody(i int) IForExpressionBodyContext {
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IForExpressionBodyContext); ok {
			if j == i {
				t = ctx.(antlr.RuleContext)
				break
			}
			j++
		}
	}

	if t == nil {
		return nil
	}

	return t.(IForExpressionBodyContext)
}

func (s *ForExpressionContext) AllAssign() []antlr.TerminalNode {
	return s.GetTokens(FqlParserAssign)
}

func (s *ForExpressionContext) Assign(i int) antlr.TerminalNode {
	return s.GetToken(FqlParserAssign, i)
}

func (s *ForExpressionContext) While() antlr.TerminalNode {
	return s.GetToken(FqlParserWhile, 0)
}

func (s *ForExpressionContext) Step() antlr.TerminalNode {
	return s.GetToken(FqlParserStep, 0)
}

func (s *ForExpressionContext) AllExpression() []IExpressionContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(IExpressionContext); ok {
			len++
		}
	}

	tst := make([]IExpressionContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(IExpressionContext); ok {
			tst[i] = t.(IExpressionContext)
			i++
		}
	}

	return tst
}

func (s *ForExpressionContext) Expression(i int) IExpressionContext {
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IExpressionContext); ok {
			if j == i {
				t = ctx.(antlr.RuleContext)
				break
			}
			j++
		}
	}

	if t == nil {
		return nil
	}

	return t.(IExpressionContext)
}

func (s *ForExpressionContext) Do() antlr.TerminalNode {
	return s.GetToken(FqlParserDo, 0)
}

func (s *ForExpressionContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ForExpressionContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *ForExpressionContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(FqlParserListener); ok {
		listenerT.EnterForExpression(s)
	}
}

func (s *ForExpressionContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(FqlParserListener); ok {
		listenerT.ExitForExpression(s)
	}
}

func (s *ForExpressionContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case FqlParserVisitor:
		return t.VisitForExpression(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *FqlParser) ForExpression() (localctx IForExpressionContext) {
	localctx = NewForExpressionContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 18, FqlParserRULE_forExpression)
	var _la int

	var _alt int

	p.SetState(244)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 11, p.GetParserRuleContext()) {
	case 1:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(195)
			p.Match(FqlParserFor)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(196)

			var _lt = p.GetTokenStream().LT(1)

			localctx.(*ForExpressionContext).valueVariable = _lt

			_la = p.GetTokenStream().LA(1)

			if !(_la == FqlParserIdentifier || _la == FqlParserIgnoreIdentifier) {
				var _ri = p.GetErrorHandler().RecoverInline(p)

				localctx.(*ForExpressionContext).valueVariable = _ri
			} else {
				p.GetErrorHandler().ReportMatch(p)
				p.Consume()
			}
		}
		p.SetState(199)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)

		if _la == FqlParserComma {
			{
				p.SetState(197)
				p.Match(FqlParserComma)
				if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
				}
			}
			{
				p.SetState(198)

				var _m = p.Match(FqlParserIdentifier)

				localctx.(*ForExpressionContext).counterVariable = _m
				if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
				}
			}

		}
		{
			p.SetState(201)
			p.Match(FqlParserIn)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(202)
			p.ForExpressionSource()
		}
		p.SetState(206)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_alt = p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 7, p.GetParserRuleContext())
		if p.HasError() {
			goto errorExit
		}
		for _alt != 2 && _alt != antlr.ATNInvalidAltNumber {
			if _alt == 1 {
				{
					p.SetState(203)
					p.ForExpressionBody()
				}

			}
			p.SetState(208)
			p.GetErrorHandler().Sync(p)
			if p.HasError() {
				goto errorExit
			}
			_alt = p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 7, p.GetParserRuleContext())
			if p.HasError() {
				goto errorExit
			}
		}
		{
			p.SetState(209)
			p.ForExpressionReturn()
		}

	case 2:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(211)
			p.Match(FqlParserFor)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(212)

			var _lt = p.GetTokenStream().LT(1)

			localctx.(*ForExpressionContext).valueVariable = _lt

			_la = p.GetTokenStream().LA(1)

			if !(_la == FqlParserIdentifier || _la == FqlParserIgnoreIdentifier) {
				var _ri = p.GetErrorHandler().RecoverInline(p)

				localctx.(*ForExpressionContext).valueVariable = _ri
			} else {
				p.GetErrorHandler().ReportMatch(p)
				p.Consume()
			}
		}
		{
			p.SetState(213)
			p.Match(FqlParserAssign)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(214)

			var _x = p.expression(0)

			localctx.(*ForExpressionContext).stepInit = _x
		}
		{
			p.SetState(215)
			p.Match(FqlParserWhile)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(216)

			var _x = p.expression(0)

			localctx.(*ForExpressionContext).stepCondition = _x
		}
		{
			p.SetState(217)
			p.Match(FqlParserStep)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(218)

			var _lt = p.GetTokenStream().LT(1)

			localctx.(*ForExpressionContext).stepVariable = _lt

			_la = p.GetTokenStream().LA(1)

			if !(_la == FqlParserIdentifier || _la == FqlParserIgnoreIdentifier) {
				var _ri = p.GetErrorHandler().RecoverInline(p)

				localctx.(*ForExpressionContext).stepVariable = _ri
			} else {
				p.GetErrorHandler().ReportMatch(p)
				p.Consume()
			}
		}
		{
			p.SetState(219)
			p.Match(FqlParserAssign)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(220)

			var _x = p.expression(0)

			localctx.(*ForExpressionContext).stepIncrement = _x
		}
		p.SetState(224)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_alt = p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 8, p.GetParserRuleContext())
		if p.HasError() {
			goto errorExit
		}
		for _alt != 2 && _alt != antlr.ATNInvalidAltNumber {
			if _alt == 1 {
				{
					p.SetState(221)
					p.ForExpressionBody()
				}

			}
			p.SetState(226)
			p.GetErrorHandler().Sync(p)
			if p.HasError() {
				goto errorExit
			}
			_alt = p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 8, p.GetParserRuleContext())
			if p.HasError() {
				goto errorExit
			}
		}
		{
			p.SetState(227)
			p.ForExpressionReturn()
		}

	case 3:
		p.EnterOuterAlt(localctx, 3)
		{
			p.SetState(229)
			p.Match(FqlParserFor)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(230)

			var _lt = p.GetTokenStream().LT(1)

			localctx.(*ForExpressionContext).valueVariable = _lt

			_la = p.GetTokenStream().LA(1)

			if !(_la == FqlParserIdentifier || _la == FqlParserIgnoreIdentifier) {
				var _ri = p.GetErrorHandler().RecoverInline(p)

				localctx.(*ForExpressionContext).valueVariable = _ri
			} else {
				p.GetErrorHandler().ReportMatch(p)
				p.Consume()
			}
		}
		p.SetState(232)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)

		if _la == FqlParserDo {
			{
				p.SetState(231)
				p.Match(FqlParserDo)
				if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
				}
			}

		}
		{
			p.SetState(234)
			p.Match(FqlParserWhile)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(235)
			p.expression(0)
		}
		p.SetState(239)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_alt = p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 10, p.GetParserRuleContext())
		if p.HasError() {
			goto errorExit
		}
		for _alt != 2 && _alt != antlr.ATNInvalidAltNumber {
			if _alt == 1 {
				{
					p.SetState(236)
					p.ForExpressionBody()
				}

			}
			p.SetState(241)
			p.GetErrorHandler().Sync(p)
			if p.HasError() {
				goto errorExit
			}
			_alt = p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 10, p.GetParserRuleContext())
			if p.HasError() {
				goto errorExit
			}
		}
		{
			p.SetState(242)
			p.ForExpressionReturn()
		}

	case antlr.ATNInvalidAltNumber:
		goto errorExit
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IForExpressionSourceContext is an interface to support dynamic dispatch.
type IForExpressionSourceContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	FunctionCallExpression() IFunctionCallExpressionContext
	ArrayLiteral() IArrayLiteralContext
	ObjectLiteral() IObjectLiteralContext
	Variable() IVariableContext
	MemberExpression() IMemberExpressionContext
	RangeOperator() IRangeOperatorContext
	Param() IParamContext

	// IsForExpressionSourceContext differentiates from other interfaces.
	IsForExpressionSourceContext()
}

type ForExpressionSourceContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyForExpressionSourceContext() *ForExpressionSourceContext {
	var p = new(ForExpressionSourceContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = FqlParserRULE_forExpressionSource
	return p
}

func InitEmptyForExpressionSourceContext(p *ForExpressionSourceContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = FqlParserRULE_forExpressionSource
}

func (*ForExpressionSourceContext) IsForExpressionSourceContext() {}

func NewForExpressionSourceContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ForExpressionSourceContext {
	var p = new(ForExpressionSourceContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = FqlParserRULE_forExpressionSource

	return p
}

func (s *ForExpressionSourceContext) GetParser() antlr.Parser { return s.parser }

func (s *ForExpressionSourceContext) FunctionCallExpression() IFunctionCallExpressionContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IFunctionCallExpressionContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IFunctionCallExpressionContext)
}

func (s *ForExpressionSourceContext) ArrayLiteral() IArrayLiteralContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IArrayLiteralContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IArrayLiteralContext)
}

func (s *ForExpressionSourceContext) ObjectLiteral() IObjectLiteralContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IObjectLiteralContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IObjectLiteralContext)
}

func (s *ForExpressionSourceContext) Variable() IVariableContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IVariableContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IVariableContext)
}

func (s *ForExpressionSourceContext) MemberExpression() IMemberExpressionContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IMemberExpressionContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IMemberExpressionContext)
}

func (s *ForExpressionSourceContext) RangeOperator() IRangeOperatorContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IRangeOperatorContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IRangeOperatorContext)
}

func (s *ForExpressionSourceContext) Param() IParamContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IParamContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IParamContext)
}

func (s *ForExpressionSourceContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ForExpressionSourceContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *ForExpressionSourceContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(FqlParserListener); ok {
		listenerT.EnterForExpressionSource(s)
	}
}

func (s *ForExpressionSourceContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(FqlParserListener); ok {
		listenerT.ExitForExpressionSource(s)
	}
}

func (s *ForExpressionSourceContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case FqlParserVisitor:
		return t.VisitForExpressionSource(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *FqlParser) ForExpressionSource() (localctx IForExpressionSourceContext) {
	localctx = NewForExpressionSourceContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 20, FqlParserRULE_forExpressionSource)
	p.SetState(253)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 12, p.GetParserRuleContext()) {
	case 1:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(246)
			p.FunctionCallExpression()
		}

	case 2:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(247)
			p.ArrayLiteral()
		}

	case 3:
		p.EnterOuterAlt(localctx, 3)
		{
			p.SetState(248)
			p.ObjectLiteral()
		}

	case 4:
		p.EnterOuterAlt(localctx, 4)
		{
			p.SetState(249)
			p.Variable()
		}

	case 5:
		p.EnterOuterAlt(localctx, 5)
		{
			p.SetState(250)
			p.MemberExpression()
		}

	case 6:
		p.EnterOuterAlt(localctx, 6)
		{
			p.SetState(251)
			p.RangeOperator()
		}

	case 7:
		p.EnterOuterAlt(localctx, 7)
		{
			p.SetState(252)
			p.Param()
		}

	case antlr.ATNInvalidAltNumber:
		goto errorExit
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IForExpressionClauseContext is an interface to support dynamic dispatch.
type IForExpressionClauseContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	LimitClause() ILimitClauseContext
	SortClause() ISortClauseContext
	FilterClause() IFilterClauseContext
	CollectClause() ICollectClauseContext

	// IsForExpressionClauseContext differentiates from other interfaces.
	IsForExpressionClauseContext()
}

type ForExpressionClauseContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyForExpressionClauseContext() *ForExpressionClauseContext {
	var p = new(ForExpressionClauseContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = FqlParserRULE_forExpressionClause
	return p
}

func InitEmptyForExpressionClauseContext(p *ForExpressionClauseContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = FqlParserRULE_forExpressionClause
}

func (*ForExpressionClauseContext) IsForExpressionClauseContext() {}

func NewForExpressionClauseContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ForExpressionClauseContext {
	var p = new(ForExpressionClauseContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = FqlParserRULE_forExpressionClause

	return p
}

func (s *ForExpressionClauseContext) GetParser() antlr.Parser { return s.parser }

func (s *ForExpressionClauseContext) LimitClause() ILimitClauseContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ILimitClauseContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(ILimitClauseContext)
}

func (s *ForExpressionClauseContext) SortClause() ISortClauseContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ISortClauseContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(ISortClauseContext)
}

func (s *ForExpressionClauseContext) FilterClause() IFilterClauseContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IFilterClauseContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IFilterClauseContext)
}

func (s *ForExpressionClauseContext) CollectClause() ICollectClauseContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ICollectClauseContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(ICollectClauseContext)
}

func (s *ForExpressionClauseContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ForExpressionClauseContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *ForExpressionClauseContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(FqlParserListener); ok {
		listenerT.EnterForExpressionClause(s)
	}
}

func (s *ForExpressionClauseContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(FqlParserListener); ok {
		listenerT.ExitForExpressionClause(s)
	}
}

func (s *ForExpressionClauseContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case FqlParserVisitor:
		return t.VisitForExpressionClause(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *FqlParser) ForExpressionClause() (localctx IForExpressionClauseContext) {
	localctx = NewForExpressionClauseContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 22, FqlParserRULE_forExpressionClause)
	p.SetState(259)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetTokenStream().LA(1) {
	case FqlParserLimit:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(255)
			p.LimitClause()
		}

	case FqlParserSort:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(256)
			p.SortClause()
		}

	case FqlParserFilter:
		p.EnterOuterAlt(localctx, 3)
		{
			p.SetState(257)
			p.FilterClause()
		}

	case FqlParserCollect:
		p.EnterOuterAlt(localctx, 4)
		{
			p.SetState(258)
			p.CollectClause()
		}

	default:
		p.SetError(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
		goto errorExit
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IForExpressionStatementContext is an interface to support dynamic dispatch.
type IForExpressionStatementContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	VariableDeclaration() IVariableDeclarationContext
	FunctionCallExpression() IFunctionCallExpressionContext

	// IsForExpressionStatementContext differentiates from other interfaces.
	IsForExpressionStatementContext()
}

type ForExpressionStatementContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyForExpressionStatementContext() *ForExpressionStatementContext {
	var p = new(ForExpressionStatementContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = FqlParserRULE_forExpressionStatement
	return p
}

func InitEmptyForExpressionStatementContext(p *ForExpressionStatementContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = FqlParserRULE_forExpressionStatement
}

func (*ForExpressionStatementContext) IsForExpressionStatementContext() {}

func NewForExpressionStatementContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ForExpressionStatementContext {
	var p = new(ForExpressionStatementContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = FqlParserRULE_forExpressionStatement

	return p
}

func (s *ForExpressionStatementContext) GetParser() antlr.Parser { return s.parser }

func (s *ForExpressionStatementContext) VariableDeclaration() IVariableDeclarationContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IVariableDeclarationContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IVariableDeclarationContext)
}

func (s *ForExpressionStatementContext) FunctionCallExpression() IFunctionCallExpressionContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IFunctionCallExpressionContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IFunctionCallExpressionContext)
}

func (s *ForExpressionStatementContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ForExpressionStatementContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *ForExpressionStatementContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(FqlParserListener); ok {
		listenerT.EnterForExpressionStatement(s)
	}
}

func (s *ForExpressionStatementContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(FqlParserListener); ok {
		listenerT.ExitForExpressionStatement(s)
	}
}

func (s *ForExpressionStatementContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case FqlParserVisitor:
		return t.VisitForExpressionStatement(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *FqlParser) ForExpressionStatement() (localctx IForExpressionStatementContext) {
	localctx = NewForExpressionStatementContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 24, FqlParserRULE_forExpressionStatement)
	p.SetState(263)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 14, p.GetParserRuleContext()) {
	case 1:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(261)
			p.VariableDeclaration()
		}

	case 2:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(262)
			p.FunctionCallExpression()
		}

	case antlr.ATNInvalidAltNumber:
		goto errorExit
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IForExpressionBodyContext is an interface to support dynamic dispatch.
type IForExpressionBodyContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	ForExpressionStatement() IForExpressionStatementContext
	ForExpressionClause() IForExpressionClauseContext

	// IsForExpressionBodyContext differentiates from other interfaces.
	IsForExpressionBodyContext()
}

type ForExpressionBodyContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyForExpressionBodyContext() *ForExpressionBodyContext {
	var p = new(ForExpressionBodyContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = FqlParserRULE_forExpressionBody
	return p
}

func InitEmptyForExpressionBodyContext(p *ForExpressionBodyContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = FqlParserRULE_forExpressionBody
}

func (*ForExpressionBodyContext) IsForExpressionBodyContext() {}

func NewForExpressionBodyContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ForExpressionBodyContext {
	var p = new(ForExpressionBodyContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = FqlParserRULE_forExpressionBody

	return p
}

func (s *ForExpressionBodyContext) GetParser() antlr.Parser { return s.parser }

func (s *ForExpressionBodyContext) ForExpressionStatement() IForExpressionStatementContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IForExpressionStatementContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IForExpressionStatementContext)
}

func (s *ForExpressionBodyContext) ForExpressionClause() IForExpressionClauseContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IForExpressionClauseContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IForExpressionClauseContext)
}

func (s *ForExpressionBodyContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ForExpressionBodyContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *ForExpressionBodyContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(FqlParserListener); ok {
		listenerT.EnterForExpressionBody(s)
	}
}

func (s *ForExpressionBodyContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(FqlParserListener); ok {
		listenerT.ExitForExpressionBody(s)
	}
}

func (s *ForExpressionBodyContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case FqlParserVisitor:
		return t.VisitForExpressionBody(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *FqlParser) ForExpressionBody() (localctx IForExpressionBodyContext) {
	localctx = NewForExpressionBodyContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 26, FqlParserRULE_forExpressionBody)
	p.SetState(267)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 15, p.GetParserRuleContext()) {
	case 1:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(265)
			p.ForExpressionStatement()
		}

	case 2:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(266)
			p.ForExpressionClause()
		}

	case antlr.ATNInvalidAltNumber:
		goto errorExit
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IForExpressionReturnContext is an interface to support dynamic dispatch.
type IForExpressionReturnContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	ReturnExpression() IReturnExpressionContext
	ForExpression() IForExpressionContext

	// IsForExpressionReturnContext differentiates from other interfaces.
	IsForExpressionReturnContext()
}

type ForExpressionReturnContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyForExpressionReturnContext() *ForExpressionReturnContext {
	var p = new(ForExpressionReturnContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = FqlParserRULE_forExpressionReturn
	return p
}

func InitEmptyForExpressionReturnContext(p *ForExpressionReturnContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = FqlParserRULE_forExpressionReturn
}

func (*ForExpressionReturnContext) IsForExpressionReturnContext() {}

func NewForExpressionReturnContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ForExpressionReturnContext {
	var p = new(ForExpressionReturnContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = FqlParserRULE_forExpressionReturn

	return p
}

func (s *ForExpressionReturnContext) GetParser() antlr.Parser { return s.parser }

func (s *ForExpressionReturnContext) ReturnExpression() IReturnExpressionContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IReturnExpressionContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IReturnExpressionContext)
}

func (s *ForExpressionReturnContext) ForExpression() IForExpressionContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IForExpressionContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IForExpressionContext)
}

func (s *ForExpressionReturnContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ForExpressionReturnContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *ForExpressionReturnContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(FqlParserListener); ok {
		listenerT.EnterForExpressionReturn(s)
	}
}

func (s *ForExpressionReturnContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(FqlParserListener); ok {
		listenerT.ExitForExpressionReturn(s)
	}
}

func (s *ForExpressionReturnContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case FqlParserVisitor:
		return t.VisitForExpressionReturn(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *FqlParser) ForExpressionReturn() (localctx IForExpressionReturnContext) {
	localctx = NewForExpressionReturnContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 28, FqlParserRULE_forExpressionReturn)
	p.SetState(271)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetTokenStream().LA(1) {
	case FqlParserReturn:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(269)
			p.ReturnExpression()
		}

	case FqlParserFor:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(270)
			p.ForExpression()
		}

	default:
		p.SetError(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
		goto errorExit
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IFilterClauseContext is an interface to support dynamic dispatch.
type IFilterClauseContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	Filter() antlr.TerminalNode
	Expression() IExpressionContext

	// IsFilterClauseContext differentiates from other interfaces.
	IsFilterClauseContext()
}

type FilterClauseContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyFilterClauseContext() *FilterClauseContext {
	var p = new(FilterClauseContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = FqlParserRULE_filterClause
	return p
}

func InitEmptyFilterClauseContext(p *FilterClauseContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = FqlParserRULE_filterClause
}

func (*FilterClauseContext) IsFilterClauseContext() {}

func NewFilterClauseContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *FilterClauseContext {
	var p = new(FilterClauseContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = FqlParserRULE_filterClause

	return p
}

func (s *FilterClauseContext) GetParser() antlr.Parser { return s.parser }

func (s *FilterClauseContext) Filter() antlr.TerminalNode {
	return s.GetToken(FqlParserFilter, 0)
}

func (s *FilterClauseContext) Expression() IExpressionContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IExpressionContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IExpressionContext)
}

func (s *FilterClauseContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *FilterClauseContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *FilterClauseContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(FqlParserListener); ok {
		listenerT.EnterFilterClause(s)
	}
}

func (s *FilterClauseContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(FqlParserListener); ok {
		listenerT.ExitFilterClause(s)
	}
}

func (s *FilterClauseContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case FqlParserVisitor:
		return t.VisitFilterClause(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *FqlParser) FilterClause() (localctx IFilterClauseContext) {
	localctx = NewFilterClauseContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 30, FqlParserRULE_filterClause)
	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(273)
		p.Match(FqlParserFilter)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(274)
		p.expression(0)
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// ILimitClauseContext is an interface to support dynamic dispatch.
type ILimitClauseContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	Limit() antlr.TerminalNode
	AllLimitClauseValue() []ILimitClauseValueContext
	LimitClauseValue(i int) ILimitClauseValueContext
	Comma() antlr.TerminalNode

	// IsLimitClauseContext differentiates from other interfaces.
	IsLimitClauseContext()
}

type LimitClauseContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyLimitClauseContext() *LimitClauseContext {
	var p = new(LimitClauseContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = FqlParserRULE_limitClause
	return p
}

func InitEmptyLimitClauseContext(p *LimitClauseContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = FqlParserRULE_limitClause
}

func (*LimitClauseContext) IsLimitClauseContext() {}

func NewLimitClauseContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *LimitClauseContext {
	var p = new(LimitClauseContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = FqlParserRULE_limitClause

	return p
}

func (s *LimitClauseContext) GetParser() antlr.Parser { return s.parser }

func (s *LimitClauseContext) Limit() antlr.TerminalNode {
	return s.GetToken(FqlParserLimit, 0)
}

func (s *LimitClauseContext) AllLimitClauseValue() []ILimitClauseValueContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(ILimitClauseValueContext); ok {
			len++
		}
	}

	tst := make([]ILimitClauseValueContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(ILimitClauseValueContext); ok {
			tst[i] = t.(ILimitClauseValueContext)
			i++
		}
	}

	return tst
}

func (s *LimitClauseContext) LimitClauseValue(i int) ILimitClauseValueContext {
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ILimitClauseValueContext); ok {
			if j == i {
				t = ctx.(antlr.RuleContext)
				break
			}
			j++
		}
	}

	if t == nil {
		return nil
	}

	return t.(ILimitClauseValueContext)
}

func (s *LimitClauseContext) Comma() antlr.TerminalNode {
	return s.GetToken(FqlParserComma, 0)
}

func (s *LimitClauseContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *LimitClauseContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *LimitClauseContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(FqlParserListener); ok {
		listenerT.EnterLimitClause(s)
	}
}

func (s *LimitClauseContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(FqlParserListener); ok {
		listenerT.ExitLimitClause(s)
	}
}

func (s *LimitClauseContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case FqlParserVisitor:
		return t.VisitLimitClause(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *FqlParser) LimitClause() (localctx ILimitClauseContext) {
	localctx = NewLimitClauseContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 32, FqlParserRULE_limitClause)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(276)
		p.Match(FqlParserLimit)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(277)
		p.LimitClauseValue()
	}
	p.SetState(280)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	if _la == FqlParserComma {
		{
			p.SetState(278)
			p.Match(FqlParserComma)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(279)
			p.LimitClauseValue()
		}

	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// ILimitClauseValueContext is an interface to support dynamic dispatch.
type ILimitClauseValueContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	IntegerLiteral() IIntegerLiteralContext
	Param() IParamContext
	Variable() IVariableContext
	FunctionCallExpression() IFunctionCallExpressionContext
	MemberExpression() IMemberExpressionContext

	// IsLimitClauseValueContext differentiates from other interfaces.
	IsLimitClauseValueContext()
}

type LimitClauseValueContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyLimitClauseValueContext() *LimitClauseValueContext {
	var p = new(LimitClauseValueContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = FqlParserRULE_limitClauseValue
	return p
}

func InitEmptyLimitClauseValueContext(p *LimitClauseValueContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = FqlParserRULE_limitClauseValue
}

func (*LimitClauseValueContext) IsLimitClauseValueContext() {}

func NewLimitClauseValueContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *LimitClauseValueContext {
	var p = new(LimitClauseValueContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = FqlParserRULE_limitClauseValue

	return p
}

func (s *LimitClauseValueContext) GetParser() antlr.Parser { return s.parser }

func (s *LimitClauseValueContext) IntegerLiteral() IIntegerLiteralContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IIntegerLiteralContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IIntegerLiteralContext)
}

func (s *LimitClauseValueContext) Param() IParamContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IParamContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IParamContext)
}

func (s *LimitClauseValueContext) Variable() IVariableContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IVariableContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IVariableContext)
}

func (s *LimitClauseValueContext) FunctionCallExpression() IFunctionCallExpressionContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IFunctionCallExpressionContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IFunctionCallExpressionContext)
}

func (s *LimitClauseValueContext) MemberExpression() IMemberExpressionContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IMemberExpressionContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IMemberExpressionContext)
}

func (s *LimitClauseValueContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *LimitClauseValueContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *LimitClauseValueContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(FqlParserListener); ok {
		listenerT.EnterLimitClauseValue(s)
	}
}

func (s *LimitClauseValueContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(FqlParserListener); ok {
		listenerT.ExitLimitClauseValue(s)
	}
}

func (s *LimitClauseValueContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case FqlParserVisitor:
		return t.VisitLimitClauseValue(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *FqlParser) LimitClauseValue() (localctx ILimitClauseValueContext) {
	localctx = NewLimitClauseValueContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 34, FqlParserRULE_limitClauseValue)
	p.SetState(287)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 18, p.GetParserRuleContext()) {
	case 1:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(282)
			p.IntegerLiteral()
		}

	case 2:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(283)
			p.Param()
		}

	case 3:
		p.EnterOuterAlt(localctx, 3)
		{
			p.SetState(284)
			p.Variable()
		}

	case 4:
		p.EnterOuterAlt(localctx, 4)
		{
			p.SetState(285)
			p.FunctionCallExpression()
		}

	case 5:
		p.EnterOuterAlt(localctx, 5)
		{
			p.SetState(286)
			p.MemberExpression()
		}

	case antlr.ATNInvalidAltNumber:
		goto errorExit
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// ISortClauseContext is an interface to support dynamic dispatch.
type ISortClauseContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	Sort() antlr.TerminalNode
	AllSortClauseExpression() []ISortClauseExpressionContext
	SortClauseExpression(i int) ISortClauseExpressionContext
	AllComma() []antlr.TerminalNode
	Comma(i int) antlr.TerminalNode

	// IsSortClauseContext differentiates from other interfaces.
	IsSortClauseContext()
}

type SortClauseContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptySortClauseContext() *SortClauseContext {
	var p = new(SortClauseContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = FqlParserRULE_sortClause
	return p
}

func InitEmptySortClauseContext(p *SortClauseContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = FqlParserRULE_sortClause
}

func (*SortClauseContext) IsSortClauseContext() {}

func NewSortClauseContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *SortClauseContext {
	var p = new(SortClauseContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = FqlParserRULE_sortClause

	return p
}

func (s *SortClauseContext) GetParser() antlr.Parser { return s.parser }

func (s *SortClauseContext) Sort() antlr.TerminalNode {
	return s.GetToken(FqlParserSort, 0)
}

func (s *SortClauseContext) AllSortClauseExpression() []ISortClauseExpressionContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(ISortClauseExpressionContext); ok {
			len++
		}
	}

	tst := make([]ISortClauseExpressionContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(ISortClauseExpressionContext); ok {
			tst[i] = t.(ISortClauseExpressionContext)
			i++
		}
	}

	return tst
}

func (s *SortClauseContext) SortClauseExpression(i int) ISortClauseExpressionContext {
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ISortClauseExpressionContext); ok {
			if j == i {
				t = ctx.(antlr.RuleContext)
				break
			}
			j++
		}
	}

	if t == nil {
		return nil
	}

	return t.(ISortClauseExpressionContext)
}

func (s *SortClauseContext) AllComma() []antlr.TerminalNode {
	return s.GetTokens(FqlParserComma)
}

func (s *SortClauseContext) Comma(i int) antlr.TerminalNode {
	return s.GetToken(FqlParserComma, i)
}

func (s *SortClauseContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *SortClauseContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *SortClauseContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(FqlParserListener); ok {
		listenerT.EnterSortClause(s)
	}
}

func (s *SortClauseContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(FqlParserListener); ok {
		listenerT.ExitSortClause(s)
	}
}

func (s *SortClauseContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case FqlParserVisitor:
		return t.VisitSortClause(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *FqlParser) SortClause() (localctx ISortClauseContext) {
	localctx = NewSortClauseContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 36, FqlParserRULE_sortClause)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(289)
		p.Match(FqlParserSort)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(290)
		p.SortClauseExpression()
	}
	p.SetState(295)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	for _la == FqlParserComma {
		{
			p.SetState(291)
			p.Match(FqlParserComma)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(292)
			p.SortClauseExpression()
		}

		p.SetState(297)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// ISortClauseExpressionContext is an interface to support dynamic dispatch.
type ISortClauseExpressionContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	Expression() IExpressionContext
	SortDirection() antlr.TerminalNode

	// IsSortClauseExpressionContext differentiates from other interfaces.
	IsSortClauseExpressionContext()
}

type SortClauseExpressionContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptySortClauseExpressionContext() *SortClauseExpressionContext {
	var p = new(SortClauseExpressionContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = FqlParserRULE_sortClauseExpression
	return p
}

func InitEmptySortClauseExpressionContext(p *SortClauseExpressionContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = FqlParserRULE_sortClauseExpression
}

func (*SortClauseExpressionContext) IsSortClauseExpressionContext() {}

func NewSortClauseExpressionContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *SortClauseExpressionContext {
	var p = new(SortClauseExpressionContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = FqlParserRULE_sortClauseExpression

	return p
}

func (s *SortClauseExpressionContext) GetParser() antlr.Parser { return s.parser }

func (s *SortClauseExpressionContext) Expression() IExpressionContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IExpressionContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IExpressionContext)
}

func (s *SortClauseExpressionContext) SortDirection() antlr.TerminalNode {
	return s.GetToken(FqlParserSortDirection, 0)
}

func (s *SortClauseExpressionContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *SortClauseExpressionContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *SortClauseExpressionContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(FqlParserListener); ok {
		listenerT.EnterSortClauseExpression(s)
	}
}

func (s *SortClauseExpressionContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(FqlParserListener); ok {
		listenerT.ExitSortClauseExpression(s)
	}
}

func (s *SortClauseExpressionContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case FqlParserVisitor:
		return t.VisitSortClauseExpression(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *FqlParser) SortClauseExpression() (localctx ISortClauseExpressionContext) {
	localctx = NewSortClauseExpressionContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 38, FqlParserRULE_sortClauseExpression)
	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(298)
		p.expression(0)
	}
	p.SetState(300)
	p.GetErrorHandler().Sync(p)

	if p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 20, p.GetParserRuleContext()) == 1 {
		{
			p.SetState(299)
			p.Match(FqlParserSortDirection)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	} else if p.HasError() { // JIM
		goto errorExit
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// ICollectClauseContext is an interface to support dynamic dispatch.
type ICollectClauseContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	Collect() antlr.TerminalNode
	CollectGrouping() ICollectGroupingContext
	CollectAggregator() ICollectAggregatorContext
	CollectGroupProjection() ICollectGroupProjectionContext
	CollectCounter() ICollectCounterContext

	// IsCollectClauseContext differentiates from other interfaces.
	IsCollectClauseContext()
}

type CollectClauseContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyCollectClauseContext() *CollectClauseContext {
	var p = new(CollectClauseContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = FqlParserRULE_collectClause
	return p
}

func InitEmptyCollectClauseContext(p *CollectClauseContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = FqlParserRULE_collectClause
}

func (*CollectClauseContext) IsCollectClauseContext() {}

func NewCollectClauseContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *CollectClauseContext {
	var p = new(CollectClauseContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = FqlParserRULE_collectClause

	return p
}

func (s *CollectClauseContext) GetParser() antlr.Parser { return s.parser }

func (s *CollectClauseContext) Collect() antlr.TerminalNode {
	return s.GetToken(FqlParserCollect, 0)
}

func (s *CollectClauseContext) CollectGrouping() ICollectGroupingContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ICollectGroupingContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(ICollectGroupingContext)
}

func (s *CollectClauseContext) CollectAggregator() ICollectAggregatorContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ICollectAggregatorContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(ICollectAggregatorContext)
}

func (s *CollectClauseContext) CollectGroupProjection() ICollectGroupProjectionContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ICollectGroupProjectionContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(ICollectGroupProjectionContext)
}

func (s *CollectClauseContext) CollectCounter() ICollectCounterContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ICollectCounterContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(ICollectCounterContext)
}

func (s *CollectClauseContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *CollectClauseContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *CollectClauseContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(FqlParserListener); ok {
		listenerT.EnterCollectClause(s)
	}
}

func (s *CollectClauseContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(FqlParserListener); ok {
		listenerT.ExitCollectClause(s)
	}
}

func (s *CollectClauseContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case FqlParserVisitor:
		return t.VisitCollectClause(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *FqlParser) CollectClause() (localctx ICollectClauseContext) {
	localctx = NewCollectClauseContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 40, FqlParserRULE_collectClause)
	p.SetState(329)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 21, p.GetParserRuleContext()) {
	case 1:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(302)
			p.Match(FqlParserCollect)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(303)
			p.CollectGrouping()
		}
		{
			p.SetState(304)
			p.CollectAggregator()
		}
		{
			p.SetState(305)
			p.CollectGroupProjection()
		}

	case 2:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(307)
			p.Match(FqlParserCollect)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(308)
			p.CollectGrouping()
		}
		{
			p.SetState(309)
			p.CollectAggregator()
		}

	case 3:
		p.EnterOuterAlt(localctx, 3)
		{
			p.SetState(311)
			p.Match(FqlParserCollect)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(312)
			p.CollectGrouping()
		}
		{
			p.SetState(313)
			p.CollectGroupProjection()
		}

	case 4:
		p.EnterOuterAlt(localctx, 4)
		{
			p.SetState(315)
			p.Match(FqlParserCollect)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(316)
			p.CollectGrouping()
		}
		{
			p.SetState(317)
			p.CollectCounter()
		}

	case 5:
		p.EnterOuterAlt(localctx, 5)
		{
			p.SetState(319)
			p.Match(FqlParserCollect)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(320)
			p.CollectGrouping()
		}

	case 6:
		p.EnterOuterAlt(localctx, 6)
		{
			p.SetState(321)
			p.Match(FqlParserCollect)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(322)
			p.CollectAggregator()
		}
		{
			p.SetState(323)
			p.CollectGroupProjection()
		}

	case 7:
		p.EnterOuterAlt(localctx, 7)
		{
			p.SetState(325)
			p.Match(FqlParserCollect)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(326)
			p.CollectAggregator()
		}

	case 8:
		p.EnterOuterAlt(localctx, 8)
		{
			p.SetState(327)
			p.Match(FqlParserCollect)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(328)
			p.CollectCounter()
		}

	case antlr.ATNInvalidAltNumber:
		goto errorExit
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// ICollectSelectorContext is an interface to support dynamic dispatch.
type ICollectSelectorContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	Identifier() antlr.TerminalNode
	Assign() antlr.TerminalNode
	Expression() IExpressionContext

	// IsCollectSelectorContext differentiates from other interfaces.
	IsCollectSelectorContext()
}

type CollectSelectorContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyCollectSelectorContext() *CollectSelectorContext {
	var p = new(CollectSelectorContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = FqlParserRULE_collectSelector
	return p
}

func InitEmptyCollectSelectorContext(p *CollectSelectorContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = FqlParserRULE_collectSelector
}

func (*CollectSelectorContext) IsCollectSelectorContext() {}

func NewCollectSelectorContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *CollectSelectorContext {
	var p = new(CollectSelectorContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = FqlParserRULE_collectSelector

	return p
}

func (s *CollectSelectorContext) GetParser() antlr.Parser { return s.parser }

func (s *CollectSelectorContext) Identifier() antlr.TerminalNode {
	return s.GetToken(FqlParserIdentifier, 0)
}

func (s *CollectSelectorContext) Assign() antlr.TerminalNode {
	return s.GetToken(FqlParserAssign, 0)
}

func (s *CollectSelectorContext) Expression() IExpressionContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IExpressionContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IExpressionContext)
}

func (s *CollectSelectorContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *CollectSelectorContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *CollectSelectorContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(FqlParserListener); ok {
		listenerT.EnterCollectSelector(s)
	}
}

func (s *CollectSelectorContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(FqlParserListener); ok {
		listenerT.ExitCollectSelector(s)
	}
}

func (s *CollectSelectorContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case FqlParserVisitor:
		return t.VisitCollectSelector(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *FqlParser) CollectSelector() (localctx ICollectSelectorContext) {
	localctx = NewCollectSelectorContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 42, FqlParserRULE_collectSelector)
	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(331)
		p.Match(FqlParserIdentifier)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(332)
		p.Match(FqlParserAssign)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(333)
		p.expression(0)
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// ICollectGroupingContext is an interface to support dynamic dispatch.
type ICollectGroupingContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	AllCollectSelector() []ICollectSelectorContext
	CollectSelector(i int) ICollectSelectorContext
	AllComma() []antlr.TerminalNode
	Comma(i int) antlr.TerminalNode

	// IsCollectGroupingContext differentiates from other interfaces.
	IsCollectGroupingContext()
}

type CollectGroupingContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyCollectGroupingContext() *CollectGroupingContext {
	var p = new(CollectGroupingContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = FqlParserRULE_collectGrouping
	return p
}

func InitEmptyCollectGroupingContext(p *CollectGroupingContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = FqlParserRULE_collectGrouping
}

func (*CollectGroupingContext) IsCollectGroupingContext() {}

func NewCollectGroupingContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *CollectGroupingContext {
	var p = new(CollectGroupingContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = FqlParserRULE_collectGrouping

	return p
}

func (s *CollectGroupingContext) GetParser() antlr.Parser { return s.parser }

func (s *CollectGroupingContext) AllCollectSelector() []ICollectSelectorContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(ICollectSelectorContext); ok {
			len++
		}
	}

	tst := make([]ICollectSelectorContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(ICollectSelectorContext); ok {
			tst[i] = t.(ICollectSelectorContext)
			i++
		}
	}

	return tst
}

func (s *CollectGroupingContext) CollectSelector(i int) ICollectSelectorContext {
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ICollectSelectorContext); ok {
			if j == i {
				t = ctx.(antlr.RuleContext)
				break
			}
			j++
		}
	}

	if t == nil {
		return nil
	}

	return t.(ICollectSelectorContext)
}

func (s *CollectGroupingContext) AllComma() []antlr.TerminalNode {
	return s.GetTokens(FqlParserComma)
}

func (s *CollectGroupingContext) Comma(i int) antlr.TerminalNode {
	return s.GetToken(FqlParserComma, i)
}

func (s *CollectGroupingContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *CollectGroupingContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *CollectGroupingContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(FqlParserListener); ok {
		listenerT.EnterCollectGrouping(s)
	}
}

func (s *CollectGroupingContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(FqlParserListener); ok {
		listenerT.ExitCollectGrouping(s)
	}
}

func (s *CollectGroupingContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case FqlParserVisitor:
		return t.VisitCollectGrouping(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *FqlParser) CollectGrouping() (localctx ICollectGroupingContext) {
	localctx = NewCollectGroupingContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 44, FqlParserRULE_collectGrouping)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(335)
		p.CollectSelector()
	}
	p.SetState(340)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	for _la == FqlParserComma {
		{
			p.SetState(336)
			p.Match(FqlParserComma)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(337)
			p.CollectSelector()
		}

		p.SetState(342)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// ICollectAggregateSelectorContext is an interface to support dynamic dispatch.
type ICollectAggregateSelectorContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	Identifier() antlr.TerminalNode
	Assign() antlr.TerminalNode
	FunctionCallExpression() IFunctionCallExpressionContext

	// IsCollectAggregateSelectorContext differentiates from other interfaces.
	IsCollectAggregateSelectorContext()
}

type CollectAggregateSelectorContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyCollectAggregateSelectorContext() *CollectAggregateSelectorContext {
	var p = new(CollectAggregateSelectorContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = FqlParserRULE_collectAggregateSelector
	return p
}

func InitEmptyCollectAggregateSelectorContext(p *CollectAggregateSelectorContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = FqlParserRULE_collectAggregateSelector
}

func (*CollectAggregateSelectorContext) IsCollectAggregateSelectorContext() {}

func NewCollectAggregateSelectorContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *CollectAggregateSelectorContext {
	var p = new(CollectAggregateSelectorContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = FqlParserRULE_collectAggregateSelector

	return p
}

func (s *CollectAggregateSelectorContext) GetParser() antlr.Parser { return s.parser }

func (s *CollectAggregateSelectorContext) Identifier() antlr.TerminalNode {
	return s.GetToken(FqlParserIdentifier, 0)
}

func (s *CollectAggregateSelectorContext) Assign() antlr.TerminalNode {
	return s.GetToken(FqlParserAssign, 0)
}

func (s *CollectAggregateSelectorContext) FunctionCallExpression() IFunctionCallExpressionContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IFunctionCallExpressionContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IFunctionCallExpressionContext)
}

func (s *CollectAggregateSelectorContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *CollectAggregateSelectorContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *CollectAggregateSelectorContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(FqlParserListener); ok {
		listenerT.EnterCollectAggregateSelector(s)
	}
}

func (s *CollectAggregateSelectorContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(FqlParserListener); ok {
		listenerT.ExitCollectAggregateSelector(s)
	}
}

func (s *CollectAggregateSelectorContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case FqlParserVisitor:
		return t.VisitCollectAggregateSelector(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *FqlParser) CollectAggregateSelector() (localctx ICollectAggregateSelectorContext) {
	localctx = NewCollectAggregateSelectorContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 46, FqlParserRULE_collectAggregateSelector)
	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(343)
		p.Match(FqlParserIdentifier)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(344)
		p.Match(FqlParserAssign)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(345)
		p.FunctionCallExpression()
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// ICollectAggregatorContext is an interface to support dynamic dispatch.
type ICollectAggregatorContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	Aggregate() antlr.TerminalNode
	AllCollectAggregateSelector() []ICollectAggregateSelectorContext
	CollectAggregateSelector(i int) ICollectAggregateSelectorContext
	AllComma() []antlr.TerminalNode
	Comma(i int) antlr.TerminalNode

	// IsCollectAggregatorContext differentiates from other interfaces.
	IsCollectAggregatorContext()
}

type CollectAggregatorContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyCollectAggregatorContext() *CollectAggregatorContext {
	var p = new(CollectAggregatorContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = FqlParserRULE_collectAggregator
	return p
}

func InitEmptyCollectAggregatorContext(p *CollectAggregatorContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = FqlParserRULE_collectAggregator
}

func (*CollectAggregatorContext) IsCollectAggregatorContext() {}

func NewCollectAggregatorContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *CollectAggregatorContext {
	var p = new(CollectAggregatorContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = FqlParserRULE_collectAggregator

	return p
}

func (s *CollectAggregatorContext) GetParser() antlr.Parser { return s.parser }

func (s *CollectAggregatorContext) Aggregate() antlr.TerminalNode {
	return s.GetToken(FqlParserAggregate, 0)
}

func (s *CollectAggregatorContext) AllCollectAggregateSelector() []ICollectAggregateSelectorContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(ICollectAggregateSelectorContext); ok {
			len++
		}
	}

	tst := make([]ICollectAggregateSelectorContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(ICollectAggregateSelectorContext); ok {
			tst[i] = t.(ICollectAggregateSelectorContext)
			i++
		}
	}

	return tst
}

func (s *CollectAggregatorContext) CollectAggregateSelector(i int) ICollectAggregateSelectorContext {
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ICollectAggregateSelectorContext); ok {
			if j == i {
				t = ctx.(antlr.RuleContext)
				break
			}
			j++
		}
	}

	if t == nil {
		return nil
	}

	return t.(ICollectAggregateSelectorContext)
}

func (s *CollectAggregatorContext) AllComma() []antlr.TerminalNode {
	return s.GetTokens(FqlParserComma)
}

func (s *CollectAggregatorContext) Comma(i int) antlr.TerminalNode {
	return s.GetToken(FqlParserComma, i)
}

func (s *CollectAggregatorContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *CollectAggregatorContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *CollectAggregatorContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(FqlParserListener); ok {
		listenerT.EnterCollectAggregator(s)
	}
}

func (s *CollectAggregatorContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(FqlParserListener); ok {
		listenerT.ExitCollectAggregator(s)
	}
}

func (s *CollectAggregatorContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case FqlParserVisitor:
		return t.VisitCollectAggregator(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *FqlParser) CollectAggregator() (localctx ICollectAggregatorContext) {
	localctx = NewCollectAggregatorContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 48, FqlParserRULE_collectAggregator)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(347)
		p.Match(FqlParserAggregate)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(348)
		p.CollectAggregateSelector()
	}
	p.SetState(353)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	for _la == FqlParserComma {
		{
			p.SetState(349)
			p.Match(FqlParserComma)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(350)
			p.CollectAggregateSelector()
		}

		p.SetState(355)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// ICollectGroupProjectionContext is an interface to support dynamic dispatch.
type ICollectGroupProjectionContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	Into() antlr.TerminalNode
	CollectSelector() ICollectSelectorContext
	Identifier() antlr.TerminalNode
	CollectGroupProjectionFilter() ICollectGroupProjectionFilterContext

	// IsCollectGroupProjectionContext differentiates from other interfaces.
	IsCollectGroupProjectionContext()
}

type CollectGroupProjectionContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyCollectGroupProjectionContext() *CollectGroupProjectionContext {
	var p = new(CollectGroupProjectionContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = FqlParserRULE_collectGroupProjection
	return p
}

func InitEmptyCollectGroupProjectionContext(p *CollectGroupProjectionContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = FqlParserRULE_collectGroupProjection
}

func (*CollectGroupProjectionContext) IsCollectGroupProjectionContext() {}

func NewCollectGroupProjectionContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *CollectGroupProjectionContext {
	var p = new(CollectGroupProjectionContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = FqlParserRULE_collectGroupProjection

	return p
}

func (s *CollectGroupProjectionContext) GetParser() antlr.Parser { return s.parser }

func (s *CollectGroupProjectionContext) Into() antlr.TerminalNode {
	return s.GetToken(FqlParserInto, 0)
}

func (s *CollectGroupProjectionContext) CollectSelector() ICollectSelectorContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ICollectSelectorContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(ICollectSelectorContext)
}

func (s *CollectGroupProjectionContext) Identifier() antlr.TerminalNode {
	return s.GetToken(FqlParserIdentifier, 0)
}

func (s *CollectGroupProjectionContext) CollectGroupProjectionFilter() ICollectGroupProjectionFilterContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ICollectGroupProjectionFilterContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(ICollectGroupProjectionFilterContext)
}

func (s *CollectGroupProjectionContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *CollectGroupProjectionContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *CollectGroupProjectionContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(FqlParserListener); ok {
		listenerT.EnterCollectGroupProjection(s)
	}
}

func (s *CollectGroupProjectionContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(FqlParserListener); ok {
		listenerT.ExitCollectGroupProjection(s)
	}
}

func (s *CollectGroupProjectionContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case FqlParserVisitor:
		return t.VisitCollectGroupProjection(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *FqlParser) CollectGroupProjection() (localctx ICollectGroupProjectionContext) {
	localctx = NewCollectGroupProjectionContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 50, FqlParserRULE_collectGroupProjection)
	p.SetState(363)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 25, p.GetParserRuleContext()) {
	case 1:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(356)
			p.Match(FqlParserInto)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(357)
			p.CollectSelector()
		}

	case 2:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(358)
			p.Match(FqlParserInto)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(359)
			p.Match(FqlParserIdentifier)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		p.SetState(361)
		p.GetErrorHandler().Sync(p)

		if p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 24, p.GetParserRuleContext()) == 1 {
			{
				p.SetState(360)
				p.CollectGroupProjectionFilter()
			}

		} else if p.HasError() { // JIM
			goto errorExit
		}

	case antlr.ATNInvalidAltNumber:
		goto errorExit
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// ICollectGroupProjectionFilterContext is an interface to support dynamic dispatch.
type ICollectGroupProjectionFilterContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	Keep() antlr.TerminalNode
	AllIdentifier() []antlr.TerminalNode
	Identifier(i int) antlr.TerminalNode
	AllComma() []antlr.TerminalNode
	Comma(i int) antlr.TerminalNode

	// IsCollectGroupProjectionFilterContext differentiates from other interfaces.
	IsCollectGroupProjectionFilterContext()
}

type CollectGroupProjectionFilterContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyCollectGroupProjectionFilterContext() *CollectGroupProjectionFilterContext {
	var p = new(CollectGroupProjectionFilterContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = FqlParserRULE_collectGroupProjectionFilter
	return p
}

func InitEmptyCollectGroupProjectionFilterContext(p *CollectGroupProjectionFilterContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = FqlParserRULE_collectGroupProjectionFilter
}

func (*CollectGroupProjectionFilterContext) IsCollectGroupProjectionFilterContext() {}

func NewCollectGroupProjectionFilterContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *CollectGroupProjectionFilterContext {
	var p = new(CollectGroupProjectionFilterContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = FqlParserRULE_collectGroupProjectionFilter

	return p
}

func (s *CollectGroupProjectionFilterContext) GetParser() antlr.Parser { return s.parser }

func (s *CollectGroupProjectionFilterContext) Keep() antlr.TerminalNode {
	return s.GetToken(FqlParserKeep, 0)
}

func (s *CollectGroupProjectionFilterContext) AllIdentifier() []antlr.TerminalNode {
	return s.GetTokens(FqlParserIdentifier)
}

func (s *CollectGroupProjectionFilterContext) Identifier(i int) antlr.TerminalNode {
	return s.GetToken(FqlParserIdentifier, i)
}

func (s *CollectGroupProjectionFilterContext) AllComma() []antlr.TerminalNode {
	return s.GetTokens(FqlParserComma)
}

func (s *CollectGroupProjectionFilterContext) Comma(i int) antlr.TerminalNode {
	return s.GetToken(FqlParserComma, i)
}

func (s *CollectGroupProjectionFilterContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *CollectGroupProjectionFilterContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *CollectGroupProjectionFilterContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(FqlParserListener); ok {
		listenerT.EnterCollectGroupProjectionFilter(s)
	}
}

func (s *CollectGroupProjectionFilterContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(FqlParserListener); ok {
		listenerT.ExitCollectGroupProjectionFilter(s)
	}
}

func (s *CollectGroupProjectionFilterContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case FqlParserVisitor:
		return t.VisitCollectGroupProjectionFilter(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *FqlParser) CollectGroupProjectionFilter() (localctx ICollectGroupProjectionFilterContext) {
	localctx = NewCollectGroupProjectionFilterContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 52, FqlParserRULE_collectGroupProjectionFilter)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(365)
		p.Match(FqlParserKeep)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(366)
		p.Match(FqlParserIdentifier)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	p.SetState(371)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	for _la == FqlParserComma {
		{
			p.SetState(367)
			p.Match(FqlParserComma)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(368)
			p.Match(FqlParserIdentifier)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

		p.SetState(373)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// ICollectCounterContext is an interface to support dynamic dispatch.
type ICollectCounterContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	With() antlr.TerminalNode
	AllIdentifier() []antlr.TerminalNode
	Identifier(i int) antlr.TerminalNode
	Into() antlr.TerminalNode

	// IsCollectCounterContext differentiates from other interfaces.
	IsCollectCounterContext()
}

type CollectCounterContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyCollectCounterContext() *CollectCounterContext {
	var p = new(CollectCounterContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = FqlParserRULE_collectCounter
	return p
}

func InitEmptyCollectCounterContext(p *CollectCounterContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = FqlParserRULE_collectCounter
}

func (*CollectCounterContext) IsCollectCounterContext() {}

func NewCollectCounterContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *CollectCounterContext {
	var p = new(CollectCounterContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = FqlParserRULE_collectCounter

	return p
}

func (s *CollectCounterContext) GetParser() antlr.Parser { return s.parser }

func (s *CollectCounterContext) With() antlr.TerminalNode {
	return s.GetToken(FqlParserWith, 0)
}

func (s *CollectCounterContext) AllIdentifier() []antlr.TerminalNode {
	return s.GetTokens(FqlParserIdentifier)
}

func (s *CollectCounterContext) Identifier(i int) antlr.TerminalNode {
	return s.GetToken(FqlParserIdentifier, i)
}

func (s *CollectCounterContext) Into() antlr.TerminalNode {
	return s.GetToken(FqlParserInto, 0)
}

func (s *CollectCounterContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *CollectCounterContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *CollectCounterContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(FqlParserListener); ok {
		listenerT.EnterCollectCounter(s)
	}
}

func (s *CollectCounterContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(FqlParserListener); ok {
		listenerT.ExitCollectCounter(s)
	}
}

func (s *CollectCounterContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case FqlParserVisitor:
		return t.VisitCollectCounter(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *FqlParser) CollectCounter() (localctx ICollectCounterContext) {
	localctx = NewCollectCounterContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 54, FqlParserRULE_collectCounter)
	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(374)
		p.Match(FqlParserWith)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(375)
		p.Match(FqlParserIdentifier)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(376)
		p.Match(FqlParserInto)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(377)
		p.Match(FqlParserIdentifier)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IWaitForExpressionContext is an interface to support dynamic dispatch.
type IWaitForExpressionContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	Waitfor() antlr.TerminalNode
	Event() antlr.TerminalNode
	WaitForEventName() IWaitForEventNameContext
	In() antlr.TerminalNode
	WaitForEventSource() IWaitForEventSourceContext
	OptionsClause() IOptionsClauseContext
	FilterClause() IFilterClauseContext
	TimeoutClause() ITimeoutClauseContext

	// IsWaitForExpressionContext differentiates from other interfaces.
	IsWaitForExpressionContext()
}

type WaitForExpressionContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyWaitForExpressionContext() *WaitForExpressionContext {
	var p = new(WaitForExpressionContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = FqlParserRULE_waitForExpression
	return p
}

func InitEmptyWaitForExpressionContext(p *WaitForExpressionContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = FqlParserRULE_waitForExpression
}

func (*WaitForExpressionContext) IsWaitForExpressionContext() {}

func NewWaitForExpressionContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *WaitForExpressionContext {
	var p = new(WaitForExpressionContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = FqlParserRULE_waitForExpression

	return p
}

func (s *WaitForExpressionContext) GetParser() antlr.Parser { return s.parser }

func (s *WaitForExpressionContext) Waitfor() antlr.TerminalNode {
	return s.GetToken(FqlParserWaitfor, 0)
}

func (s *WaitForExpressionContext) Event() antlr.TerminalNode {
	return s.GetToken(FqlParserEvent, 0)
}

func (s *WaitForExpressionContext) WaitForEventName() IWaitForEventNameContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IWaitForEventNameContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IWaitForEventNameContext)
}

func (s *WaitForExpressionContext) In() antlr.TerminalNode {
	return s.GetToken(FqlParserIn, 0)
}

func (s *WaitForExpressionContext) WaitForEventSource() IWaitForEventSourceContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IWaitForEventSourceContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IWaitForEventSourceContext)
}

func (s *WaitForExpressionContext) OptionsClause() IOptionsClauseContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IOptionsClauseContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IOptionsClauseContext)
}

func (s *WaitForExpressionContext) FilterClause() IFilterClauseContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IFilterClauseContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IFilterClauseContext)
}

func (s *WaitForExpressionContext) TimeoutClause() ITimeoutClauseContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ITimeoutClauseContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(ITimeoutClauseContext)
}

func (s *WaitForExpressionContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *WaitForExpressionContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *WaitForExpressionContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(FqlParserListener); ok {
		listenerT.EnterWaitForExpression(s)
	}
}

func (s *WaitForExpressionContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(FqlParserListener); ok {
		listenerT.ExitWaitForExpression(s)
	}
}

func (s *WaitForExpressionContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case FqlParserVisitor:
		return t.VisitWaitForExpression(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *FqlParser) WaitForExpression() (localctx IWaitForExpressionContext) {
	localctx = NewWaitForExpressionContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 56, FqlParserRULE_waitForExpression)
	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(379)
		p.Match(FqlParserWaitfor)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(380)
		p.Match(FqlParserEvent)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(381)
		p.WaitForEventName()
	}
	{
		p.SetState(382)
		p.Match(FqlParserIn)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(383)
		p.WaitForEventSource()
	}
	p.SetState(385)
	p.GetErrorHandler().Sync(p)

	if p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 27, p.GetParserRuleContext()) == 1 {
		{
			p.SetState(384)
			p.OptionsClause()
		}

	} else if p.HasError() { // JIM
		goto errorExit
	}
	p.SetState(388)
	p.GetErrorHandler().Sync(p)

	if p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 28, p.GetParserRuleContext()) == 1 {
		{
			p.SetState(387)
			p.FilterClause()
		}

	} else if p.HasError() { // JIM
		goto errorExit
	}
	p.SetState(391)
	p.GetErrorHandler().Sync(p)

	if p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 29, p.GetParserRuleContext()) == 1 {
		{
			p.SetState(390)
			p.TimeoutClause()
		}

	} else if p.HasError() { // JIM
		goto errorExit
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IWaitForEventNameContext is an interface to support dynamic dispatch.
type IWaitForEventNameContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	StringLiteral() IStringLiteralContext
	Variable() IVariableContext
	Param() IParamContext
	FunctionCallExpression() IFunctionCallExpressionContext
	MemberExpression() IMemberExpressionContext

	// IsWaitForEventNameContext differentiates from other interfaces.
	IsWaitForEventNameContext()
}

type WaitForEventNameContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyWaitForEventNameContext() *WaitForEventNameContext {
	var p = new(WaitForEventNameContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = FqlParserRULE_waitForEventName
	return p
}

func InitEmptyWaitForEventNameContext(p *WaitForEventNameContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = FqlParserRULE_waitForEventName
}

func (*WaitForEventNameContext) IsWaitForEventNameContext() {}

func NewWaitForEventNameContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *WaitForEventNameContext {
	var p = new(WaitForEventNameContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = FqlParserRULE_waitForEventName

	return p
}

func (s *WaitForEventNameContext) GetParser() antlr.Parser { return s.parser }

func (s *WaitForEventNameContext) StringLiteral() IStringLiteralContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IStringLiteralContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IStringLiteralContext)
}

func (s *WaitForEventNameContext) Variable() IVariableContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IVariableContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IVariableContext)
}

func (s *WaitForEventNameContext) Param() IParamContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IParamContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IParamContext)
}

func (s *WaitForEventNameContext) FunctionCallExpression() IFunctionCallExpressionContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IFunctionCallExpressionContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IFunctionCallExpressionContext)
}

func (s *WaitForEventNameContext) MemberExpression() IMemberExpressionContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IMemberExpressionContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IMemberExpressionContext)
}

func (s *WaitForEventNameContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *WaitForEventNameContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *WaitForEventNameContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(FqlParserListener); ok {
		listenerT.EnterWaitForEventName(s)
	}
}

func (s *WaitForEventNameContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(FqlParserListener); ok {
		listenerT.ExitWaitForEventName(s)
	}
}

func (s *WaitForEventNameContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case FqlParserVisitor:
		return t.VisitWaitForEventName(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *FqlParser) WaitForEventName() (localctx IWaitForEventNameContext) {
	localctx = NewWaitForEventNameContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 58, FqlParserRULE_waitForEventName)
	p.SetState(398)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 30, p.GetParserRuleContext()) {
	case 1:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(393)
			p.StringLiteral()
		}

	case 2:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(394)
			p.Variable()
		}

	case 3:
		p.EnterOuterAlt(localctx, 3)
		{
			p.SetState(395)
			p.Param()
		}

	case 4:
		p.EnterOuterAlt(localctx, 4)
		{
			p.SetState(396)
			p.FunctionCallExpression()
		}

	case 5:
		p.EnterOuterAlt(localctx, 5)
		{
			p.SetState(397)
			p.MemberExpression()
		}

	case antlr.ATNInvalidAltNumber:
		goto errorExit
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IWaitForEventSourceContext is an interface to support dynamic dispatch.
type IWaitForEventSourceContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	FunctionCallExpression() IFunctionCallExpressionContext
	Variable() IVariableContext
	MemberExpression() IMemberExpressionContext

	// IsWaitForEventSourceContext differentiates from other interfaces.
	IsWaitForEventSourceContext()
}

type WaitForEventSourceContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyWaitForEventSourceContext() *WaitForEventSourceContext {
	var p = new(WaitForEventSourceContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = FqlParserRULE_waitForEventSource
	return p
}

func InitEmptyWaitForEventSourceContext(p *WaitForEventSourceContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = FqlParserRULE_waitForEventSource
}

func (*WaitForEventSourceContext) IsWaitForEventSourceContext() {}

func NewWaitForEventSourceContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *WaitForEventSourceContext {
	var p = new(WaitForEventSourceContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = FqlParserRULE_waitForEventSource

	return p
}

func (s *WaitForEventSourceContext) GetParser() antlr.Parser { return s.parser }

func (s *WaitForEventSourceContext) FunctionCallExpression() IFunctionCallExpressionContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IFunctionCallExpressionContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IFunctionCallExpressionContext)
}

func (s *WaitForEventSourceContext) Variable() IVariableContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IVariableContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IVariableContext)
}

func (s *WaitForEventSourceContext) MemberExpression() IMemberExpressionContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IMemberExpressionContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IMemberExpressionContext)
}

func (s *WaitForEventSourceContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *WaitForEventSourceContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *WaitForEventSourceContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(FqlParserListener); ok {
		listenerT.EnterWaitForEventSource(s)
	}
}

func (s *WaitForEventSourceContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(FqlParserListener); ok {
		listenerT.ExitWaitForEventSource(s)
	}
}

func (s *WaitForEventSourceContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case FqlParserVisitor:
		return t.VisitWaitForEventSource(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *FqlParser) WaitForEventSource() (localctx IWaitForEventSourceContext) {
	localctx = NewWaitForEventSourceContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 60, FqlParserRULE_waitForEventSource)
	p.SetState(403)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 31, p.GetParserRuleContext()) {
	case 1:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(400)
			p.FunctionCallExpression()
		}

	case 2:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(401)
			p.Variable()
		}

	case 3:
		p.EnterOuterAlt(localctx, 3)
		{
			p.SetState(402)
			p.MemberExpression()
		}

	case antlr.ATNInvalidAltNumber:
		goto errorExit
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IOptionsClauseContext is an interface to support dynamic dispatch.
type IOptionsClauseContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	Options() antlr.TerminalNode
	ObjectLiteral() IObjectLiteralContext

	// IsOptionsClauseContext differentiates from other interfaces.
	IsOptionsClauseContext()
}

type OptionsClauseContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyOptionsClauseContext() *OptionsClauseContext {
	var p = new(OptionsClauseContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = FqlParserRULE_optionsClause
	return p
}

func InitEmptyOptionsClauseContext(p *OptionsClauseContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = FqlParserRULE_optionsClause
}

func (*OptionsClauseContext) IsOptionsClauseContext() {}

func NewOptionsClauseContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *OptionsClauseContext {
	var p = new(OptionsClauseContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = FqlParserRULE_optionsClause

	return p
}

func (s *OptionsClauseContext) GetParser() antlr.Parser { return s.parser }

func (s *OptionsClauseContext) Options() antlr.TerminalNode {
	return s.GetToken(FqlParserOptions, 0)
}

func (s *OptionsClauseContext) ObjectLiteral() IObjectLiteralContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IObjectLiteralContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IObjectLiteralContext)
}

func (s *OptionsClauseContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *OptionsClauseContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *OptionsClauseContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(FqlParserListener); ok {
		listenerT.EnterOptionsClause(s)
	}
}

func (s *OptionsClauseContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(FqlParserListener); ok {
		listenerT.ExitOptionsClause(s)
	}
}

func (s *OptionsClauseContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case FqlParserVisitor:
		return t.VisitOptionsClause(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *FqlParser) OptionsClause() (localctx IOptionsClauseContext) {
	localctx = NewOptionsClauseContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 62, FqlParserRULE_optionsClause)
	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(405)
		p.Match(FqlParserOptions)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(406)
		p.ObjectLiteral()
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// ITimeoutClauseContext is an interface to support dynamic dispatch.
type ITimeoutClauseContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	Timeout() antlr.TerminalNode
	IntegerLiteral() IIntegerLiteralContext
	Variable() IVariableContext
	Param() IParamContext
	MemberExpression() IMemberExpressionContext
	FunctionCall() IFunctionCallContext

	// IsTimeoutClauseContext differentiates from other interfaces.
	IsTimeoutClauseContext()
}

type TimeoutClauseContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyTimeoutClauseContext() *TimeoutClauseContext {
	var p = new(TimeoutClauseContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = FqlParserRULE_timeoutClause
	return p
}

func InitEmptyTimeoutClauseContext(p *TimeoutClauseContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = FqlParserRULE_timeoutClause
}

func (*TimeoutClauseContext) IsTimeoutClauseContext() {}

func NewTimeoutClauseContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *TimeoutClauseContext {
	var p = new(TimeoutClauseContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = FqlParserRULE_timeoutClause

	return p
}

func (s *TimeoutClauseContext) GetParser() antlr.Parser { return s.parser }

func (s *TimeoutClauseContext) Timeout() antlr.TerminalNode {
	return s.GetToken(FqlParserTimeout, 0)
}

func (s *TimeoutClauseContext) IntegerLiteral() IIntegerLiteralContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IIntegerLiteralContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IIntegerLiteralContext)
}

func (s *TimeoutClauseContext) Variable() IVariableContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IVariableContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IVariableContext)
}

func (s *TimeoutClauseContext) Param() IParamContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IParamContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IParamContext)
}

func (s *TimeoutClauseContext) MemberExpression() IMemberExpressionContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IMemberExpressionContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IMemberExpressionContext)
}

func (s *TimeoutClauseContext) FunctionCall() IFunctionCallContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IFunctionCallContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IFunctionCallContext)
}

func (s *TimeoutClauseContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *TimeoutClauseContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *TimeoutClauseContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(FqlParserListener); ok {
		listenerT.EnterTimeoutClause(s)
	}
}

func (s *TimeoutClauseContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(FqlParserListener); ok {
		listenerT.ExitTimeoutClause(s)
	}
}

func (s *TimeoutClauseContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case FqlParserVisitor:
		return t.VisitTimeoutClause(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *FqlParser) TimeoutClause() (localctx ITimeoutClauseContext) {
	localctx = NewTimeoutClauseContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 64, FqlParserRULE_timeoutClause)
	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(408)
		p.Match(FqlParserTimeout)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	p.SetState(414)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 32, p.GetParserRuleContext()) {
	case 1:
		{
			p.SetState(409)
			p.IntegerLiteral()
		}

	case 2:
		{
			p.SetState(410)
			p.Variable()
		}

	case 3:
		{
			p.SetState(411)
			p.Param()
		}

	case 4:
		{
			p.SetState(412)
			p.MemberExpression()
		}

	case 5:
		{
			p.SetState(413)
			p.FunctionCall()
		}

	case antlr.ATNInvalidAltNumber:
		goto errorExit
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IParamContext is an interface to support dynamic dispatch.
type IParamContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	Param() antlr.TerminalNode
	Identifier() antlr.TerminalNode
	SafeReservedWord() ISafeReservedWordContext

	// IsParamContext differentiates from other interfaces.
	IsParamContext()
}

type ParamContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyParamContext() *ParamContext {
	var p = new(ParamContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = FqlParserRULE_param
	return p
}

func InitEmptyParamContext(p *ParamContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = FqlParserRULE_param
}

func (*ParamContext) IsParamContext() {}

func NewParamContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ParamContext {
	var p = new(ParamContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = FqlParserRULE_param

	return p
}

func (s *ParamContext) GetParser() antlr.Parser { return s.parser }

func (s *ParamContext) Param() antlr.TerminalNode {
	return s.GetToken(FqlParserParam, 0)
}

func (s *ParamContext) Identifier() antlr.TerminalNode {
	return s.GetToken(FqlParserIdentifier, 0)
}

func (s *ParamContext) SafeReservedWord() ISafeReservedWordContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ISafeReservedWordContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(ISafeReservedWordContext)
}

func (s *ParamContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ParamContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *ParamContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(FqlParserListener); ok {
		listenerT.EnterParam(s)
	}
}

func (s *ParamContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(FqlParserListener); ok {
		listenerT.ExitParam(s)
	}
}

func (s *ParamContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case FqlParserVisitor:
		return t.VisitParam(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *FqlParser) Param() (localctx IParamContext) {
	localctx = NewParamContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 66, FqlParserRULE_param)
	p.SetState(420)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 33, p.GetParserRuleContext()) {
	case 1:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(416)
			p.Match(FqlParserParam)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(417)
			p.Match(FqlParserIdentifier)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	case 2:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(418)
			p.Match(FqlParserParam)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(419)
			p.SafeReservedWord()
		}

	case antlr.ATNInvalidAltNumber:
		goto errorExit
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IVariableContext is an interface to support dynamic dispatch.
type IVariableContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	Identifier() antlr.TerminalNode
	SafeReservedWord() ISafeReservedWordContext

	// IsVariableContext differentiates from other interfaces.
	IsVariableContext()
}

type VariableContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyVariableContext() *VariableContext {
	var p = new(VariableContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = FqlParserRULE_variable
	return p
}

func InitEmptyVariableContext(p *VariableContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = FqlParserRULE_variable
}

func (*VariableContext) IsVariableContext() {}

func NewVariableContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *VariableContext {
	var p = new(VariableContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = FqlParserRULE_variable

	return p
}

func (s *VariableContext) GetParser() antlr.Parser { return s.parser }

func (s *VariableContext) Identifier() antlr.TerminalNode {
	return s.GetToken(FqlParserIdentifier, 0)
}

func (s *VariableContext) SafeReservedWord() ISafeReservedWordContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ISafeReservedWordContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(ISafeReservedWordContext)
}

func (s *VariableContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *VariableContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *VariableContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(FqlParserListener); ok {
		listenerT.EnterVariable(s)
	}
}

func (s *VariableContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(FqlParserListener); ok {
		listenerT.ExitVariable(s)
	}
}

func (s *VariableContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case FqlParserVisitor:
		return t.VisitVariable(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *FqlParser) Variable() (localctx IVariableContext) {
	localctx = NewVariableContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 68, FqlParserRULE_variable)
	p.SetState(424)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetTokenStream().LA(1) {
	case FqlParserIdentifier:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(422)
			p.Match(FqlParserIdentifier)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	case FqlParserAnd, FqlParserOr, FqlParserOptions, FqlParserTimeout, FqlParserDistinct, FqlParserFilter, FqlParserCurrent, FqlParserSort, FqlParserLimit, FqlParserCollect, FqlParserSortDirection, FqlParserInto, FqlParserKeep, FqlParserWith, FqlParserAll, FqlParserAny, FqlParserAggregate, FqlParserEvent, FqlParserStep:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(423)
			p.SafeReservedWord()
		}

	default:
		p.SetError(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
		goto errorExit
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// ILiteralContext is an interface to support dynamic dispatch.
type ILiteralContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	ArrayLiteral() IArrayLiteralContext
	ObjectLiteral() IObjectLiteralContext
	BooleanLiteral() IBooleanLiteralContext
	StringLiteral() IStringLiteralContext
	FloatLiteral() IFloatLiteralContext
	IntegerLiteral() IIntegerLiteralContext
	NoneLiteral() INoneLiteralContext

	// IsLiteralContext differentiates from other interfaces.
	IsLiteralContext()
}

type LiteralContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyLiteralContext() *LiteralContext {
	var p = new(LiteralContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = FqlParserRULE_literal
	return p
}

func InitEmptyLiteralContext(p *LiteralContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = FqlParserRULE_literal
}

func (*LiteralContext) IsLiteralContext() {}

func NewLiteralContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *LiteralContext {
	var p = new(LiteralContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = FqlParserRULE_literal

	return p
}

func (s *LiteralContext) GetParser() antlr.Parser { return s.parser }

func (s *LiteralContext) ArrayLiteral() IArrayLiteralContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IArrayLiteralContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IArrayLiteralContext)
}

func (s *LiteralContext) ObjectLiteral() IObjectLiteralContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IObjectLiteralContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IObjectLiteralContext)
}

func (s *LiteralContext) BooleanLiteral() IBooleanLiteralContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IBooleanLiteralContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IBooleanLiteralContext)
}

func (s *LiteralContext) StringLiteral() IStringLiteralContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IStringLiteralContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IStringLiteralContext)
}

func (s *LiteralContext) FloatLiteral() IFloatLiteralContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IFloatLiteralContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IFloatLiteralContext)
}

func (s *LiteralContext) IntegerLiteral() IIntegerLiteralContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IIntegerLiteralContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IIntegerLiteralContext)
}

func (s *LiteralContext) NoneLiteral() INoneLiteralContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(INoneLiteralContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(INoneLiteralContext)
}

func (s *LiteralContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *LiteralContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *LiteralContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(FqlParserListener); ok {
		listenerT.EnterLiteral(s)
	}
}

func (s *LiteralContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(FqlParserListener); ok {
		listenerT.ExitLiteral(s)
	}
}

func (s *LiteralContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case FqlParserVisitor:
		return t.VisitLiteral(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *FqlParser) Literal() (localctx ILiteralContext) {
	localctx = NewLiteralContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 70, FqlParserRULE_literal)
	p.SetState(433)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetTokenStream().LA(1) {
	case FqlParserOpenBracket:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(426)
			p.ArrayLiteral()
		}

	case FqlParserOpenBrace:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(427)
			p.ObjectLiteral()
		}

	case FqlParserBooleanLiteral:
		p.EnterOuterAlt(localctx, 3)
		{
			p.SetState(428)
			p.BooleanLiteral()
		}

	case FqlParserStringLiteral:
		p.EnterOuterAlt(localctx, 4)
		{
			p.SetState(429)
			p.StringLiteral()
		}

	case FqlParserFloatLiteral:
		p.EnterOuterAlt(localctx, 5)
		{
			p.SetState(430)
			p.FloatLiteral()
		}

	case FqlParserIntegerLiteral:
		p.EnterOuterAlt(localctx, 6)
		{
			p.SetState(431)
			p.IntegerLiteral()
		}

	case FqlParserNone, FqlParserNull:
		p.EnterOuterAlt(localctx, 7)
		{
			p.SetState(432)
			p.NoneLiteral()
		}

	default:
		p.SetError(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
		goto errorExit
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IArrayLiteralContext is an interface to support dynamic dispatch.
type IArrayLiteralContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	OpenBracket() antlr.TerminalNode
	CloseBracket() antlr.TerminalNode
	ArgumentList() IArgumentListContext

	// IsArrayLiteralContext differentiates from other interfaces.
	IsArrayLiteralContext()
}

type ArrayLiteralContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyArrayLiteralContext() *ArrayLiteralContext {
	var p = new(ArrayLiteralContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = FqlParserRULE_arrayLiteral
	return p
}

func InitEmptyArrayLiteralContext(p *ArrayLiteralContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = FqlParserRULE_arrayLiteral
}

func (*ArrayLiteralContext) IsArrayLiteralContext() {}

func NewArrayLiteralContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ArrayLiteralContext {
	var p = new(ArrayLiteralContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = FqlParserRULE_arrayLiteral

	return p
}

func (s *ArrayLiteralContext) GetParser() antlr.Parser { return s.parser }

func (s *ArrayLiteralContext) OpenBracket() antlr.TerminalNode {
	return s.GetToken(FqlParserOpenBracket, 0)
}

func (s *ArrayLiteralContext) CloseBracket() antlr.TerminalNode {
	return s.GetToken(FqlParserCloseBracket, 0)
}

func (s *ArrayLiteralContext) ArgumentList() IArgumentListContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IArgumentListContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IArgumentListContext)
}

func (s *ArrayLiteralContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ArrayLiteralContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *ArrayLiteralContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(FqlParserListener); ok {
		listenerT.EnterArrayLiteral(s)
	}
}

func (s *ArrayLiteralContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(FqlParserListener); ok {
		listenerT.ExitArrayLiteral(s)
	}
}

func (s *ArrayLiteralContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case FqlParserVisitor:
		return t.VisitArrayLiteral(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *FqlParser) ArrayLiteral() (localctx IArrayLiteralContext) {
	localctx = NewArrayLiteralContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 72, FqlParserRULE_arrayLiteral)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(435)
		p.Match(FqlParserOpenBracket)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	p.SetState(437)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	if (int64((_la-9)) & ^0x3f) == 0 && ((int64(1)<<(_la-9))&2233785415159480341) != 0 {
		{
			p.SetState(436)
			p.ArgumentList()
		}

	}
	{
		p.SetState(439)
		p.Match(FqlParserCloseBracket)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IObjectLiteralContext is an interface to support dynamic dispatch.
type IObjectLiteralContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	OpenBrace() antlr.TerminalNode
	CloseBrace() antlr.TerminalNode
	AllPropertyAssignment() []IPropertyAssignmentContext
	PropertyAssignment(i int) IPropertyAssignmentContext
	AllComma() []antlr.TerminalNode
	Comma(i int) antlr.TerminalNode

	// IsObjectLiteralContext differentiates from other interfaces.
	IsObjectLiteralContext()
}

type ObjectLiteralContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyObjectLiteralContext() *ObjectLiteralContext {
	var p = new(ObjectLiteralContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = FqlParserRULE_objectLiteral
	return p
}

func InitEmptyObjectLiteralContext(p *ObjectLiteralContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = FqlParserRULE_objectLiteral
}

func (*ObjectLiteralContext) IsObjectLiteralContext() {}

func NewObjectLiteralContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ObjectLiteralContext {
	var p = new(ObjectLiteralContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = FqlParserRULE_objectLiteral

	return p
}

func (s *ObjectLiteralContext) GetParser() antlr.Parser { return s.parser }

func (s *ObjectLiteralContext) OpenBrace() antlr.TerminalNode {
	return s.GetToken(FqlParserOpenBrace, 0)
}

func (s *ObjectLiteralContext) CloseBrace() antlr.TerminalNode {
	return s.GetToken(FqlParserCloseBrace, 0)
}

func (s *ObjectLiteralContext) AllPropertyAssignment() []IPropertyAssignmentContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(IPropertyAssignmentContext); ok {
			len++
		}
	}

	tst := make([]IPropertyAssignmentContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(IPropertyAssignmentContext); ok {
			tst[i] = t.(IPropertyAssignmentContext)
			i++
		}
	}

	return tst
}

func (s *ObjectLiteralContext) PropertyAssignment(i int) IPropertyAssignmentContext {
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IPropertyAssignmentContext); ok {
			if j == i {
				t = ctx.(antlr.RuleContext)
				break
			}
			j++
		}
	}

	if t == nil {
		return nil
	}

	return t.(IPropertyAssignmentContext)
}

func (s *ObjectLiteralContext) AllComma() []antlr.TerminalNode {
	return s.GetTokens(FqlParserComma)
}

func (s *ObjectLiteralContext) Comma(i int) antlr.TerminalNode {
	return s.GetToken(FqlParserComma, i)
}

func (s *ObjectLiteralContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ObjectLiteralContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *ObjectLiteralContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(FqlParserListener); ok {
		listenerT.EnterObjectLiteral(s)
	}
}

func (s *ObjectLiteralContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(FqlParserListener); ok {
		listenerT.ExitObjectLiteral(s)
	}
}

func (s *ObjectLiteralContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case FqlParserVisitor:
		return t.VisitObjectLiteral(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *FqlParser) ObjectLiteral() (localctx IObjectLiteralContext) {
	localctx = NewObjectLiteralContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 74, FqlParserRULE_objectLiteral)
	var _la int

	var _alt int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(441)
		p.Match(FqlParserOpenBrace)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	p.SetState(453)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	if (int64((_la-9)) & ^0x3f) == 0 && ((int64(1)<<(_la-9))&216172782097399809) != 0 {
		{
			p.SetState(442)
			p.PropertyAssignment()
		}
		p.SetState(447)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_alt = p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 37, p.GetParserRuleContext())
		if p.HasError() {
			goto errorExit
		}
		for _alt != 2 && _alt != antlr.ATNInvalidAltNumber {
			if _alt == 1 {
				{
					p.SetState(443)
					p.Match(FqlParserComma)
					if p.HasError() {
						// Recognition error - abort rule
						goto errorExit
					}
				}
				{
					p.SetState(444)
					p.PropertyAssignment()
				}

			}
			p.SetState(449)
			p.GetErrorHandler().Sync(p)
			if p.HasError() {
				goto errorExit
			}
			_alt = p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 37, p.GetParserRuleContext())
			if p.HasError() {
				goto errorExit
			}
		}
		p.SetState(451)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)

		if _la == FqlParserComma {
			{
				p.SetState(450)
				p.Match(FqlParserComma)
				if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
				}
			}

		}

	}
	{
		p.SetState(455)
		p.Match(FqlParserCloseBrace)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IBooleanLiteralContext is an interface to support dynamic dispatch.
type IBooleanLiteralContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	BooleanLiteral() antlr.TerminalNode

	// IsBooleanLiteralContext differentiates from other interfaces.
	IsBooleanLiteralContext()
}

type BooleanLiteralContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyBooleanLiteralContext() *BooleanLiteralContext {
	var p = new(BooleanLiteralContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = FqlParserRULE_booleanLiteral
	return p
}

func InitEmptyBooleanLiteralContext(p *BooleanLiteralContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = FqlParserRULE_booleanLiteral
}

func (*BooleanLiteralContext) IsBooleanLiteralContext() {}

func NewBooleanLiteralContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *BooleanLiteralContext {
	var p = new(BooleanLiteralContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = FqlParserRULE_booleanLiteral

	return p
}

func (s *BooleanLiteralContext) GetParser() antlr.Parser { return s.parser }

func (s *BooleanLiteralContext) BooleanLiteral() antlr.TerminalNode {
	return s.GetToken(FqlParserBooleanLiteral, 0)
}

func (s *BooleanLiteralContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *BooleanLiteralContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *BooleanLiteralContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(FqlParserListener); ok {
		listenerT.EnterBooleanLiteral(s)
	}
}

func (s *BooleanLiteralContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(FqlParserListener); ok {
		listenerT.ExitBooleanLiteral(s)
	}
}

func (s *BooleanLiteralContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case FqlParserVisitor:
		return t.VisitBooleanLiteral(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *FqlParser) BooleanLiteral() (localctx IBooleanLiteralContext) {
	localctx = NewBooleanLiteralContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 76, FqlParserRULE_booleanLiteral)
	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(457)
		p.Match(FqlParserBooleanLiteral)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IStringLiteralContext is an interface to support dynamic dispatch.
type IStringLiteralContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	StringLiteral() antlr.TerminalNode

	// IsStringLiteralContext differentiates from other interfaces.
	IsStringLiteralContext()
}

type StringLiteralContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyStringLiteralContext() *StringLiteralContext {
	var p = new(StringLiteralContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = FqlParserRULE_stringLiteral
	return p
}

func InitEmptyStringLiteralContext(p *StringLiteralContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = FqlParserRULE_stringLiteral
}

func (*StringLiteralContext) IsStringLiteralContext() {}

func NewStringLiteralContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *StringLiteralContext {
	var p = new(StringLiteralContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = FqlParserRULE_stringLiteral

	return p
}

func (s *StringLiteralContext) GetParser() antlr.Parser { return s.parser }

func (s *StringLiteralContext) StringLiteral() antlr.TerminalNode {
	return s.GetToken(FqlParserStringLiteral, 0)
}

func (s *StringLiteralContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *StringLiteralContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *StringLiteralContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(FqlParserListener); ok {
		listenerT.EnterStringLiteral(s)
	}
}

func (s *StringLiteralContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(FqlParserListener); ok {
		listenerT.ExitStringLiteral(s)
	}
}

func (s *StringLiteralContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case FqlParserVisitor:
		return t.VisitStringLiteral(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *FqlParser) StringLiteral() (localctx IStringLiteralContext) {
	localctx = NewStringLiteralContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 78, FqlParserRULE_stringLiteral)
	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(459)
		p.Match(FqlParserStringLiteral)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IFloatLiteralContext is an interface to support dynamic dispatch.
type IFloatLiteralContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	FloatLiteral() antlr.TerminalNode

	// IsFloatLiteralContext differentiates from other interfaces.
	IsFloatLiteralContext()
}

type FloatLiteralContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyFloatLiteralContext() *FloatLiteralContext {
	var p = new(FloatLiteralContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = FqlParserRULE_floatLiteral
	return p
}

func InitEmptyFloatLiteralContext(p *FloatLiteralContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = FqlParserRULE_floatLiteral
}

func (*FloatLiteralContext) IsFloatLiteralContext() {}

func NewFloatLiteralContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *FloatLiteralContext {
	var p = new(FloatLiteralContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = FqlParserRULE_floatLiteral

	return p
}

func (s *FloatLiteralContext) GetParser() antlr.Parser { return s.parser }

func (s *FloatLiteralContext) FloatLiteral() antlr.TerminalNode {
	return s.GetToken(FqlParserFloatLiteral, 0)
}

func (s *FloatLiteralContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *FloatLiteralContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *FloatLiteralContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(FqlParserListener); ok {
		listenerT.EnterFloatLiteral(s)
	}
}

func (s *FloatLiteralContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(FqlParserListener); ok {
		listenerT.ExitFloatLiteral(s)
	}
}

func (s *FloatLiteralContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case FqlParserVisitor:
		return t.VisitFloatLiteral(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *FqlParser) FloatLiteral() (localctx IFloatLiteralContext) {
	localctx = NewFloatLiteralContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 80, FqlParserRULE_floatLiteral)
	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(461)
		p.Match(FqlParserFloatLiteral)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IIntegerLiteralContext is an interface to support dynamic dispatch.
type IIntegerLiteralContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	IntegerLiteral() antlr.TerminalNode

	// IsIntegerLiteralContext differentiates from other interfaces.
	IsIntegerLiteralContext()
}

type IntegerLiteralContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyIntegerLiteralContext() *IntegerLiteralContext {
	var p = new(IntegerLiteralContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = FqlParserRULE_integerLiteral
	return p
}

func InitEmptyIntegerLiteralContext(p *IntegerLiteralContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = FqlParserRULE_integerLiteral
}

func (*IntegerLiteralContext) IsIntegerLiteralContext() {}

func NewIntegerLiteralContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *IntegerLiteralContext {
	var p = new(IntegerLiteralContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = FqlParserRULE_integerLiteral

	return p
}

func (s *IntegerLiteralContext) GetParser() antlr.Parser { return s.parser }

func (s *IntegerLiteralContext) IntegerLiteral() antlr.TerminalNode {
	return s.GetToken(FqlParserIntegerLiteral, 0)
}

func (s *IntegerLiteralContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *IntegerLiteralContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *IntegerLiteralContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(FqlParserListener); ok {
		listenerT.EnterIntegerLiteral(s)
	}
}

func (s *IntegerLiteralContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(FqlParserListener); ok {
		listenerT.ExitIntegerLiteral(s)
	}
}

func (s *IntegerLiteralContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case FqlParserVisitor:
		return t.VisitIntegerLiteral(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *FqlParser) IntegerLiteral() (localctx IIntegerLiteralContext) {
	localctx = NewIntegerLiteralContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 82, FqlParserRULE_integerLiteral)
	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(463)
		p.Match(FqlParserIntegerLiteral)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// INoneLiteralContext is an interface to support dynamic dispatch.
type INoneLiteralContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	Null() antlr.TerminalNode
	None() antlr.TerminalNode

	// IsNoneLiteralContext differentiates from other interfaces.
	IsNoneLiteralContext()
}

type NoneLiteralContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyNoneLiteralContext() *NoneLiteralContext {
	var p = new(NoneLiteralContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = FqlParserRULE_noneLiteral
	return p
}

func InitEmptyNoneLiteralContext(p *NoneLiteralContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = FqlParserRULE_noneLiteral
}

func (*NoneLiteralContext) IsNoneLiteralContext() {}

func NewNoneLiteralContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *NoneLiteralContext {
	var p = new(NoneLiteralContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = FqlParserRULE_noneLiteral

	return p
}

func (s *NoneLiteralContext) GetParser() antlr.Parser { return s.parser }

func (s *NoneLiteralContext) Null() antlr.TerminalNode {
	return s.GetToken(FqlParserNull, 0)
}

func (s *NoneLiteralContext) None() antlr.TerminalNode {
	return s.GetToken(FqlParserNone, 0)
}

func (s *NoneLiteralContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *NoneLiteralContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *NoneLiteralContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(FqlParserListener); ok {
		listenerT.EnterNoneLiteral(s)
	}
}

func (s *NoneLiteralContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(FqlParserListener); ok {
		listenerT.ExitNoneLiteral(s)
	}
}

func (s *NoneLiteralContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case FqlParserVisitor:
		return t.VisitNoneLiteral(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *FqlParser) NoneLiteral() (localctx INoneLiteralContext) {
	localctx = NewNoneLiteralContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 84, FqlParserRULE_noneLiteral)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(465)
		_la = p.GetTokenStream().LA(1)

		if !(_la == FqlParserNone || _la == FqlParserNull) {
			p.GetErrorHandler().RecoverInline(p)
		} else {
			p.GetErrorHandler().ReportMatch(p)
			p.Consume()
		}
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IPropertyAssignmentContext is an interface to support dynamic dispatch.
type IPropertyAssignmentContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	PropertyName() IPropertyNameContext
	Colon() antlr.TerminalNode
	Expression() IExpressionContext
	ComputedPropertyName() IComputedPropertyNameContext
	Variable() IVariableContext

	// IsPropertyAssignmentContext differentiates from other interfaces.
	IsPropertyAssignmentContext()
}

type PropertyAssignmentContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyPropertyAssignmentContext() *PropertyAssignmentContext {
	var p = new(PropertyAssignmentContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = FqlParserRULE_propertyAssignment
	return p
}

func InitEmptyPropertyAssignmentContext(p *PropertyAssignmentContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = FqlParserRULE_propertyAssignment
}

func (*PropertyAssignmentContext) IsPropertyAssignmentContext() {}

func NewPropertyAssignmentContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *PropertyAssignmentContext {
	var p = new(PropertyAssignmentContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = FqlParserRULE_propertyAssignment

	return p
}

func (s *PropertyAssignmentContext) GetParser() antlr.Parser { return s.parser }

func (s *PropertyAssignmentContext) PropertyName() IPropertyNameContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IPropertyNameContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IPropertyNameContext)
}

func (s *PropertyAssignmentContext) Colon() antlr.TerminalNode {
	return s.GetToken(FqlParserColon, 0)
}

func (s *PropertyAssignmentContext) Expression() IExpressionContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IExpressionContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IExpressionContext)
}

func (s *PropertyAssignmentContext) ComputedPropertyName() IComputedPropertyNameContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IComputedPropertyNameContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IComputedPropertyNameContext)
}

func (s *PropertyAssignmentContext) Variable() IVariableContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IVariableContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IVariableContext)
}

func (s *PropertyAssignmentContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *PropertyAssignmentContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *PropertyAssignmentContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(FqlParserListener); ok {
		listenerT.EnterPropertyAssignment(s)
	}
}

func (s *PropertyAssignmentContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(FqlParserListener); ok {
		listenerT.ExitPropertyAssignment(s)
	}
}

func (s *PropertyAssignmentContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case FqlParserVisitor:
		return t.VisitPropertyAssignment(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *FqlParser) PropertyAssignment() (localctx IPropertyAssignmentContext) {
	localctx = NewPropertyAssignmentContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 86, FqlParserRULE_propertyAssignment)
	p.SetState(476)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 40, p.GetParserRuleContext()) {
	case 1:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(467)
			p.PropertyName()
		}
		{
			p.SetState(468)
			p.Match(FqlParserColon)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(469)
			p.expression(0)
		}

	case 2:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(471)
			p.ComputedPropertyName()
		}
		{
			p.SetState(472)
			p.Match(FqlParserColon)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(473)
			p.expression(0)
		}

	case 3:
		p.EnterOuterAlt(localctx, 3)
		{
			p.SetState(475)
			p.Variable()
		}

	case antlr.ATNInvalidAltNumber:
		goto errorExit
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IComputedPropertyNameContext is an interface to support dynamic dispatch.
type IComputedPropertyNameContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	OpenBracket() antlr.TerminalNode
	Expression() IExpressionContext
	CloseBracket() antlr.TerminalNode

	// IsComputedPropertyNameContext differentiates from other interfaces.
	IsComputedPropertyNameContext()
}

type ComputedPropertyNameContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyComputedPropertyNameContext() *ComputedPropertyNameContext {
	var p = new(ComputedPropertyNameContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = FqlParserRULE_computedPropertyName
	return p
}

func InitEmptyComputedPropertyNameContext(p *ComputedPropertyNameContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = FqlParserRULE_computedPropertyName
}

func (*ComputedPropertyNameContext) IsComputedPropertyNameContext() {}

func NewComputedPropertyNameContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ComputedPropertyNameContext {
	var p = new(ComputedPropertyNameContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = FqlParserRULE_computedPropertyName

	return p
}

func (s *ComputedPropertyNameContext) GetParser() antlr.Parser { return s.parser }

func (s *ComputedPropertyNameContext) OpenBracket() antlr.TerminalNode {
	return s.GetToken(FqlParserOpenBracket, 0)
}

func (s *ComputedPropertyNameContext) Expression() IExpressionContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IExpressionContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IExpressionContext)
}

func (s *ComputedPropertyNameContext) CloseBracket() antlr.TerminalNode {
	return s.GetToken(FqlParserCloseBracket, 0)
}

func (s *ComputedPropertyNameContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ComputedPropertyNameContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *ComputedPropertyNameContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(FqlParserListener); ok {
		listenerT.EnterComputedPropertyName(s)
	}
}

func (s *ComputedPropertyNameContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(FqlParserListener); ok {
		listenerT.ExitComputedPropertyName(s)
	}
}

func (s *ComputedPropertyNameContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case FqlParserVisitor:
		return t.VisitComputedPropertyName(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *FqlParser) ComputedPropertyName() (localctx IComputedPropertyNameContext) {
	localctx = NewComputedPropertyNameContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 88, FqlParserRULE_computedPropertyName)
	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(478)
		p.Match(FqlParserOpenBracket)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(479)
		p.expression(0)
	}
	{
		p.SetState(480)
		p.Match(FqlParserCloseBracket)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IPropertyNameContext is an interface to support dynamic dispatch.
type IPropertyNameContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	Identifier() antlr.TerminalNode
	StringLiteral() IStringLiteralContext
	Param() IParamContext
	SafeReservedWord() ISafeReservedWordContext
	UnsafeReservedWord() IUnsafeReservedWordContext

	// IsPropertyNameContext differentiates from other interfaces.
	IsPropertyNameContext()
}

type PropertyNameContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyPropertyNameContext() *PropertyNameContext {
	var p = new(PropertyNameContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = FqlParserRULE_propertyName
	return p
}

func InitEmptyPropertyNameContext(p *PropertyNameContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = FqlParserRULE_propertyName
}

func (*PropertyNameContext) IsPropertyNameContext() {}

func NewPropertyNameContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *PropertyNameContext {
	var p = new(PropertyNameContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = FqlParserRULE_propertyName

	return p
}

func (s *PropertyNameContext) GetParser() antlr.Parser { return s.parser }

func (s *PropertyNameContext) Identifier() antlr.TerminalNode {
	return s.GetToken(FqlParserIdentifier, 0)
}

func (s *PropertyNameContext) StringLiteral() IStringLiteralContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IStringLiteralContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IStringLiteralContext)
}

func (s *PropertyNameContext) Param() IParamContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IParamContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IParamContext)
}

func (s *PropertyNameContext) SafeReservedWord() ISafeReservedWordContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ISafeReservedWordContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(ISafeReservedWordContext)
}

func (s *PropertyNameContext) UnsafeReservedWord() IUnsafeReservedWordContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IUnsafeReservedWordContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IUnsafeReservedWordContext)
}

func (s *PropertyNameContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *PropertyNameContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *PropertyNameContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(FqlParserListener); ok {
		listenerT.EnterPropertyName(s)
	}
}

func (s *PropertyNameContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(FqlParserListener); ok {
		listenerT.ExitPropertyName(s)
	}
}

func (s *PropertyNameContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case FqlParserVisitor:
		return t.VisitPropertyName(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *FqlParser) PropertyName() (localctx IPropertyNameContext) {
	localctx = NewPropertyNameContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 90, FqlParserRULE_propertyName)
	p.SetState(487)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetTokenStream().LA(1) {
	case FqlParserIdentifier:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(482)
			p.Match(FqlParserIdentifier)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	case FqlParserStringLiteral:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(483)
			p.StringLiteral()
		}

	case FqlParserParam:
		p.EnterOuterAlt(localctx, 3)
		{
			p.SetState(484)
			p.Param()
		}

	case FqlParserAnd, FqlParserOr, FqlParserOptions, FqlParserTimeout, FqlParserDistinct, FqlParserFilter, FqlParserCurrent, FqlParserSort, FqlParserLimit, FqlParserCollect, FqlParserSortDirection, FqlParserInto, FqlParserKeep, FqlParserWith, FqlParserAll, FqlParserAny, FqlParserAggregate, FqlParserEvent, FqlParserStep:
		p.EnterOuterAlt(localctx, 4)
		{
			p.SetState(485)
			p.SafeReservedWord()
		}

	case FqlParserFor, FqlParserReturn, FqlParserWaitfor, FqlParserLet, FqlParserNone, FqlParserNull, FqlParserBooleanLiteral, FqlParserUse, FqlParserLike, FqlParserNot, FqlParserIn, FqlParserDo, FqlParserWhile:
		p.EnterOuterAlt(localctx, 5)
		{
			p.SetState(486)
			p.UnsafeReservedWord()
		}

	default:
		p.SetError(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
		goto errorExit
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// INamespaceIdentifierContext is an interface to support dynamic dispatch.
type INamespaceIdentifierContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	Namespace() INamespaceContext
	Identifier() antlr.TerminalNode

	// IsNamespaceIdentifierContext differentiates from other interfaces.
	IsNamespaceIdentifierContext()
}

type NamespaceIdentifierContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyNamespaceIdentifierContext() *NamespaceIdentifierContext {
	var p = new(NamespaceIdentifierContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = FqlParserRULE_namespaceIdentifier
	return p
}

func InitEmptyNamespaceIdentifierContext(p *NamespaceIdentifierContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = FqlParserRULE_namespaceIdentifier
}

func (*NamespaceIdentifierContext) IsNamespaceIdentifierContext() {}

func NewNamespaceIdentifierContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *NamespaceIdentifierContext {
	var p = new(NamespaceIdentifierContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = FqlParserRULE_namespaceIdentifier

	return p
}

func (s *NamespaceIdentifierContext) GetParser() antlr.Parser { return s.parser }

func (s *NamespaceIdentifierContext) Namespace() INamespaceContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(INamespaceContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(INamespaceContext)
}

func (s *NamespaceIdentifierContext) Identifier() antlr.TerminalNode {
	return s.GetToken(FqlParserIdentifier, 0)
}

func (s *NamespaceIdentifierContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *NamespaceIdentifierContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *NamespaceIdentifierContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(FqlParserListener); ok {
		listenerT.EnterNamespaceIdentifier(s)
	}
}

func (s *NamespaceIdentifierContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(FqlParserListener); ok {
		listenerT.ExitNamespaceIdentifier(s)
	}
}

func (s *NamespaceIdentifierContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case FqlParserVisitor:
		return t.VisitNamespaceIdentifier(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *FqlParser) NamespaceIdentifier() (localctx INamespaceIdentifierContext) {
	localctx = NewNamespaceIdentifierContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 92, FqlParserRULE_namespaceIdentifier)
	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(489)
		p.Namespace()
	}
	{
		p.SetState(490)
		p.Match(FqlParserIdentifier)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// INamespaceContext is an interface to support dynamic dispatch.
type INamespaceContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	AllNamespaceSegment() []antlr.TerminalNode
	NamespaceSegment(i int) antlr.TerminalNode

	// IsNamespaceContext differentiates from other interfaces.
	IsNamespaceContext()
}

type NamespaceContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyNamespaceContext() *NamespaceContext {
	var p = new(NamespaceContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = FqlParserRULE_namespace
	return p
}

func InitEmptyNamespaceContext(p *NamespaceContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = FqlParserRULE_namespace
}

func (*NamespaceContext) IsNamespaceContext() {}

func NewNamespaceContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *NamespaceContext {
	var p = new(NamespaceContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = FqlParserRULE_namespace

	return p
}

func (s *NamespaceContext) GetParser() antlr.Parser { return s.parser }

func (s *NamespaceContext) AllNamespaceSegment() []antlr.TerminalNode {
	return s.GetTokens(FqlParserNamespaceSegment)
}

func (s *NamespaceContext) NamespaceSegment(i int) antlr.TerminalNode {
	return s.GetToken(FqlParserNamespaceSegment, i)
}

func (s *NamespaceContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *NamespaceContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *NamespaceContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(FqlParserListener); ok {
		listenerT.EnterNamespace(s)
	}
}

func (s *NamespaceContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(FqlParserListener); ok {
		listenerT.ExitNamespace(s)
	}
}

func (s *NamespaceContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case FqlParserVisitor:
		return t.VisitNamespace(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *FqlParser) Namespace() (localctx INamespaceContext) {
	localctx = NewNamespaceContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 94, FqlParserRULE_namespace)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	p.SetState(495)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	for _la == FqlParserNamespaceSegment {
		{
			p.SetState(492)
			p.Match(FqlParserNamespaceSegment)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

		p.SetState(497)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IMemberExpressionContext is an interface to support dynamic dispatch.
type IMemberExpressionContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	MemberExpressionSource() IMemberExpressionSourceContext
	AllMemberExpressionPath() []IMemberExpressionPathContext
	MemberExpressionPath(i int) IMemberExpressionPathContext

	// IsMemberExpressionContext differentiates from other interfaces.
	IsMemberExpressionContext()
}

type MemberExpressionContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyMemberExpressionContext() *MemberExpressionContext {
	var p = new(MemberExpressionContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = FqlParserRULE_memberExpression
	return p
}

func InitEmptyMemberExpressionContext(p *MemberExpressionContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = FqlParserRULE_memberExpression
}

func (*MemberExpressionContext) IsMemberExpressionContext() {}

func NewMemberExpressionContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *MemberExpressionContext {
	var p = new(MemberExpressionContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = FqlParserRULE_memberExpression

	return p
}

func (s *MemberExpressionContext) GetParser() antlr.Parser { return s.parser }

func (s *MemberExpressionContext) MemberExpressionSource() IMemberExpressionSourceContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IMemberExpressionSourceContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IMemberExpressionSourceContext)
}

func (s *MemberExpressionContext) AllMemberExpressionPath() []IMemberExpressionPathContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(IMemberExpressionPathContext); ok {
			len++
		}
	}

	tst := make([]IMemberExpressionPathContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(IMemberExpressionPathContext); ok {
			tst[i] = t.(IMemberExpressionPathContext)
			i++
		}
	}

	return tst
}

func (s *MemberExpressionContext) MemberExpressionPath(i int) IMemberExpressionPathContext {
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IMemberExpressionPathContext); ok {
			if j == i {
				t = ctx.(antlr.RuleContext)
				break
			}
			j++
		}
	}

	if t == nil {
		return nil
	}

	return t.(IMemberExpressionPathContext)
}

func (s *MemberExpressionContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *MemberExpressionContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *MemberExpressionContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(FqlParserListener); ok {
		listenerT.EnterMemberExpression(s)
	}
}

func (s *MemberExpressionContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(FqlParserListener); ok {
		listenerT.ExitMemberExpression(s)
	}
}

func (s *MemberExpressionContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case FqlParserVisitor:
		return t.VisitMemberExpression(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *FqlParser) MemberExpression() (localctx IMemberExpressionContext) {
	localctx = NewMemberExpressionContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 96, FqlParserRULE_memberExpression)
	var _alt int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(498)
		p.MemberExpressionSource()
	}
	p.SetState(500)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_alt = 1
	for ok := true; ok; ok = _alt != 2 && _alt != antlr.ATNInvalidAltNumber {
		switch _alt {
		case 1:
			{
				p.SetState(499)
				p.MemberExpressionPath()
			}

		default:
			p.SetError(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
			goto errorExit
		}

		p.SetState(502)
		p.GetErrorHandler().Sync(p)
		_alt = p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 43, p.GetParserRuleContext())
		if p.HasError() {
			goto errorExit
		}
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IMemberExpressionSourceContext is an interface to support dynamic dispatch.
type IMemberExpressionSourceContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	Variable() IVariableContext
	Param() IParamContext
	ArrayLiteral() IArrayLiteralContext
	ObjectLiteral() IObjectLiteralContext
	FunctionCall() IFunctionCallContext

	// IsMemberExpressionSourceContext differentiates from other interfaces.
	IsMemberExpressionSourceContext()
}

type MemberExpressionSourceContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyMemberExpressionSourceContext() *MemberExpressionSourceContext {
	var p = new(MemberExpressionSourceContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = FqlParserRULE_memberExpressionSource
	return p
}

func InitEmptyMemberExpressionSourceContext(p *MemberExpressionSourceContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = FqlParserRULE_memberExpressionSource
}

func (*MemberExpressionSourceContext) IsMemberExpressionSourceContext() {}

func NewMemberExpressionSourceContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *MemberExpressionSourceContext {
	var p = new(MemberExpressionSourceContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = FqlParserRULE_memberExpressionSource

	return p
}

func (s *MemberExpressionSourceContext) GetParser() antlr.Parser { return s.parser }

func (s *MemberExpressionSourceContext) Variable() IVariableContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IVariableContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IVariableContext)
}

func (s *MemberExpressionSourceContext) Param() IParamContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IParamContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IParamContext)
}

func (s *MemberExpressionSourceContext) ArrayLiteral() IArrayLiteralContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IArrayLiteralContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IArrayLiteralContext)
}

func (s *MemberExpressionSourceContext) ObjectLiteral() IObjectLiteralContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IObjectLiteralContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IObjectLiteralContext)
}

func (s *MemberExpressionSourceContext) FunctionCall() IFunctionCallContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IFunctionCallContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IFunctionCallContext)
}

func (s *MemberExpressionSourceContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *MemberExpressionSourceContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *MemberExpressionSourceContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(FqlParserListener); ok {
		listenerT.EnterMemberExpressionSource(s)
	}
}

func (s *MemberExpressionSourceContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(FqlParserListener); ok {
		listenerT.ExitMemberExpressionSource(s)
	}
}

func (s *MemberExpressionSourceContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case FqlParserVisitor:
		return t.VisitMemberExpressionSource(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *FqlParser) MemberExpressionSource() (localctx IMemberExpressionSourceContext) {
	localctx = NewMemberExpressionSourceContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 98, FqlParserRULE_memberExpressionSource)
	p.SetState(509)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 44, p.GetParserRuleContext()) {
	case 1:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(504)
			p.Variable()
		}

	case 2:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(505)
			p.Param()
		}

	case 3:
		p.EnterOuterAlt(localctx, 3)
		{
			p.SetState(506)
			p.ArrayLiteral()
		}

	case 4:
		p.EnterOuterAlt(localctx, 4)
		{
			p.SetState(507)
			p.ObjectLiteral()
		}

	case 5:
		p.EnterOuterAlt(localctx, 5)
		{
			p.SetState(508)
			p.FunctionCall()
		}

	case antlr.ATNInvalidAltNumber:
		goto errorExit
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IFunctionCallExpressionContext is an interface to support dynamic dispatch.
type IFunctionCallExpressionContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	FunctionCall() IFunctionCallContext
	ErrorOperator() IErrorOperatorContext

	// IsFunctionCallExpressionContext differentiates from other interfaces.
	IsFunctionCallExpressionContext()
}

type FunctionCallExpressionContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyFunctionCallExpressionContext() *FunctionCallExpressionContext {
	var p = new(FunctionCallExpressionContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = FqlParserRULE_functionCallExpression
	return p
}

func InitEmptyFunctionCallExpressionContext(p *FunctionCallExpressionContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = FqlParserRULE_functionCallExpression
}

func (*FunctionCallExpressionContext) IsFunctionCallExpressionContext() {}

func NewFunctionCallExpressionContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *FunctionCallExpressionContext {
	var p = new(FunctionCallExpressionContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = FqlParserRULE_functionCallExpression

	return p
}

func (s *FunctionCallExpressionContext) GetParser() antlr.Parser { return s.parser }

func (s *FunctionCallExpressionContext) FunctionCall() IFunctionCallContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IFunctionCallContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IFunctionCallContext)
}

func (s *FunctionCallExpressionContext) ErrorOperator() IErrorOperatorContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IErrorOperatorContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IErrorOperatorContext)
}

func (s *FunctionCallExpressionContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *FunctionCallExpressionContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *FunctionCallExpressionContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(FqlParserListener); ok {
		listenerT.EnterFunctionCallExpression(s)
	}
}

func (s *FunctionCallExpressionContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(FqlParserListener); ok {
		listenerT.ExitFunctionCallExpression(s)
	}
}

func (s *FunctionCallExpressionContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case FqlParserVisitor:
		return t.VisitFunctionCallExpression(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *FqlParser) FunctionCallExpression() (localctx IFunctionCallExpressionContext) {
	localctx = NewFunctionCallExpressionContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 100, FqlParserRULE_functionCallExpression)
	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(511)
		p.FunctionCall()
	}
	p.SetState(513)
	p.GetErrorHandler().Sync(p)

	if p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 45, p.GetParserRuleContext()) == 1 {
		{
			p.SetState(512)
			p.ErrorOperator()
		}

	} else if p.HasError() { // JIM
		goto errorExit
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IFunctionCallContext is an interface to support dynamic dispatch.
type IFunctionCallContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	Namespace() INamespaceContext
	FunctionName() IFunctionNameContext
	OpenParen() antlr.TerminalNode
	CloseParen() antlr.TerminalNode
	ArgumentList() IArgumentListContext

	// IsFunctionCallContext differentiates from other interfaces.
	IsFunctionCallContext()
}

type FunctionCallContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyFunctionCallContext() *FunctionCallContext {
	var p = new(FunctionCallContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = FqlParserRULE_functionCall
	return p
}

func InitEmptyFunctionCallContext(p *FunctionCallContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = FqlParserRULE_functionCall
}

func (*FunctionCallContext) IsFunctionCallContext() {}

func NewFunctionCallContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *FunctionCallContext {
	var p = new(FunctionCallContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = FqlParserRULE_functionCall

	return p
}

func (s *FunctionCallContext) GetParser() antlr.Parser { return s.parser }

func (s *FunctionCallContext) Namespace() INamespaceContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(INamespaceContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(INamespaceContext)
}

func (s *FunctionCallContext) FunctionName() IFunctionNameContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IFunctionNameContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IFunctionNameContext)
}

func (s *FunctionCallContext) OpenParen() antlr.TerminalNode {
	return s.GetToken(FqlParserOpenParen, 0)
}

func (s *FunctionCallContext) CloseParen() antlr.TerminalNode {
	return s.GetToken(FqlParserCloseParen, 0)
}

func (s *FunctionCallContext) ArgumentList() IArgumentListContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IArgumentListContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IArgumentListContext)
}

func (s *FunctionCallContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *FunctionCallContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *FunctionCallContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(FqlParserListener); ok {
		listenerT.EnterFunctionCall(s)
	}
}

func (s *FunctionCallContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(FqlParserListener); ok {
		listenerT.ExitFunctionCall(s)
	}
}

func (s *FunctionCallContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case FqlParserVisitor:
		return t.VisitFunctionCall(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *FqlParser) FunctionCall() (localctx IFunctionCallContext) {
	localctx = NewFunctionCallContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 102, FqlParserRULE_functionCall)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(515)
		p.Namespace()
	}
	{
		p.SetState(516)
		p.FunctionName()
	}
	{
		p.SetState(517)
		p.Match(FqlParserOpenParen)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	p.SetState(519)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	if (int64((_la-9)) & ^0x3f) == 0 && ((int64(1)<<(_la-9))&2233785415159480341) != 0 {
		{
			p.SetState(518)
			p.ArgumentList()
		}

	}
	{
		p.SetState(521)
		p.Match(FqlParserCloseParen)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IFunctionNameContext is an interface to support dynamic dispatch.
type IFunctionNameContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	Identifier() antlr.TerminalNode
	SafeReservedWord() ISafeReservedWordContext
	UnsafeReservedWord() IUnsafeReservedWordContext

	// IsFunctionNameContext differentiates from other interfaces.
	IsFunctionNameContext()
}

type FunctionNameContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyFunctionNameContext() *FunctionNameContext {
	var p = new(FunctionNameContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = FqlParserRULE_functionName
	return p
}

func InitEmptyFunctionNameContext(p *FunctionNameContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = FqlParserRULE_functionName
}

func (*FunctionNameContext) IsFunctionNameContext() {}

func NewFunctionNameContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *FunctionNameContext {
	var p = new(FunctionNameContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = FqlParserRULE_functionName

	return p
}

func (s *FunctionNameContext) GetParser() antlr.Parser { return s.parser }

func (s *FunctionNameContext) Identifier() antlr.TerminalNode {
	return s.GetToken(FqlParserIdentifier, 0)
}

func (s *FunctionNameContext) SafeReservedWord() ISafeReservedWordContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ISafeReservedWordContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(ISafeReservedWordContext)
}

func (s *FunctionNameContext) UnsafeReservedWord() IUnsafeReservedWordContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IUnsafeReservedWordContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IUnsafeReservedWordContext)
}

func (s *FunctionNameContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *FunctionNameContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *FunctionNameContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(FqlParserListener); ok {
		listenerT.EnterFunctionName(s)
	}
}

func (s *FunctionNameContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(FqlParserListener); ok {
		listenerT.ExitFunctionName(s)
	}
}

func (s *FunctionNameContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case FqlParserVisitor:
		return t.VisitFunctionName(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *FqlParser) FunctionName() (localctx IFunctionNameContext) {
	localctx = NewFunctionNameContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 104, FqlParserRULE_functionName)
	p.SetState(526)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetTokenStream().LA(1) {
	case FqlParserIdentifier:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(523)
			p.Match(FqlParserIdentifier)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	case FqlParserAnd, FqlParserOr, FqlParserOptions, FqlParserTimeout, FqlParserDistinct, FqlParserFilter, FqlParserCurrent, FqlParserSort, FqlParserLimit, FqlParserCollect, FqlParserSortDirection, FqlParserInto, FqlParserKeep, FqlParserWith, FqlParserAll, FqlParserAny, FqlParserAggregate, FqlParserEvent, FqlParserStep:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(524)
			p.SafeReservedWord()
		}

	case FqlParserFor, FqlParserReturn, FqlParserWaitfor, FqlParserLet, FqlParserNone, FqlParserNull, FqlParserBooleanLiteral, FqlParserUse, FqlParserLike, FqlParserNot, FqlParserIn, FqlParserDo, FqlParserWhile:
		p.EnterOuterAlt(localctx, 3)
		{
			p.SetState(525)
			p.UnsafeReservedWord()
		}

	default:
		p.SetError(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
		goto errorExit
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IArgumentListContext is an interface to support dynamic dispatch.
type IArgumentListContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	AllExpression() []IExpressionContext
	Expression(i int) IExpressionContext
	AllComma() []antlr.TerminalNode
	Comma(i int) antlr.TerminalNode

	// IsArgumentListContext differentiates from other interfaces.
	IsArgumentListContext()
}

type ArgumentListContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyArgumentListContext() *ArgumentListContext {
	var p = new(ArgumentListContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = FqlParserRULE_argumentList
	return p
}

func InitEmptyArgumentListContext(p *ArgumentListContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = FqlParserRULE_argumentList
}

func (*ArgumentListContext) IsArgumentListContext() {}

func NewArgumentListContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ArgumentListContext {
	var p = new(ArgumentListContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = FqlParserRULE_argumentList

	return p
}

func (s *ArgumentListContext) GetParser() antlr.Parser { return s.parser }

func (s *ArgumentListContext) AllExpression() []IExpressionContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(IExpressionContext); ok {
			len++
		}
	}

	tst := make([]IExpressionContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(IExpressionContext); ok {
			tst[i] = t.(IExpressionContext)
			i++
		}
	}

	return tst
}

func (s *ArgumentListContext) Expression(i int) IExpressionContext {
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IExpressionContext); ok {
			if j == i {
				t = ctx.(antlr.RuleContext)
				break
			}
			j++
		}
	}

	if t == nil {
		return nil
	}

	return t.(IExpressionContext)
}

func (s *ArgumentListContext) AllComma() []antlr.TerminalNode {
	return s.GetTokens(FqlParserComma)
}

func (s *ArgumentListContext) Comma(i int) antlr.TerminalNode {
	return s.GetToken(FqlParserComma, i)
}

func (s *ArgumentListContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ArgumentListContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *ArgumentListContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(FqlParserListener); ok {
		listenerT.EnterArgumentList(s)
	}
}

func (s *ArgumentListContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(FqlParserListener); ok {
		listenerT.ExitArgumentList(s)
	}
}

func (s *ArgumentListContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case FqlParserVisitor:
		return t.VisitArgumentList(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *FqlParser) ArgumentList() (localctx IArgumentListContext) {
	localctx = NewArgumentListContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 106, FqlParserRULE_argumentList)
	var _la int

	var _alt int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(528)
		p.expression(0)
	}
	p.SetState(533)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_alt = p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 48, p.GetParserRuleContext())
	if p.HasError() {
		goto errorExit
	}
	for _alt != 2 && _alt != antlr.ATNInvalidAltNumber {
		if _alt == 1 {
			{
				p.SetState(529)
				p.Match(FqlParserComma)
				if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
				}
			}
			{
				p.SetState(530)
				p.expression(0)
			}

		}
		p.SetState(535)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_alt = p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 48, p.GetParserRuleContext())
		if p.HasError() {
			goto errorExit
		}
	}
	p.SetState(537)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	if _la == FqlParserComma {
		{
			p.SetState(536)
			p.Match(FqlParserComma)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IMemberExpressionPathContext is an interface to support dynamic dispatch.
type IMemberExpressionPathContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	Dot() antlr.TerminalNode
	PropertyName() IPropertyNameContext
	ErrorOperator() IErrorOperatorContext
	ComputedPropertyName() IComputedPropertyNameContext

	// IsMemberExpressionPathContext differentiates from other interfaces.
	IsMemberExpressionPathContext()
}

type MemberExpressionPathContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyMemberExpressionPathContext() *MemberExpressionPathContext {
	var p = new(MemberExpressionPathContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = FqlParserRULE_memberExpressionPath
	return p
}

func InitEmptyMemberExpressionPathContext(p *MemberExpressionPathContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = FqlParserRULE_memberExpressionPath
}

func (*MemberExpressionPathContext) IsMemberExpressionPathContext() {}

func NewMemberExpressionPathContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *MemberExpressionPathContext {
	var p = new(MemberExpressionPathContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = FqlParserRULE_memberExpressionPath

	return p
}

func (s *MemberExpressionPathContext) GetParser() antlr.Parser { return s.parser }

func (s *MemberExpressionPathContext) Dot() antlr.TerminalNode {
	return s.GetToken(FqlParserDot, 0)
}

func (s *MemberExpressionPathContext) PropertyName() IPropertyNameContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IPropertyNameContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IPropertyNameContext)
}

func (s *MemberExpressionPathContext) ErrorOperator() IErrorOperatorContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IErrorOperatorContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IErrorOperatorContext)
}

func (s *MemberExpressionPathContext) ComputedPropertyName() IComputedPropertyNameContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IComputedPropertyNameContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IComputedPropertyNameContext)
}

func (s *MemberExpressionPathContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *MemberExpressionPathContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *MemberExpressionPathContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(FqlParserListener); ok {
		listenerT.EnterMemberExpressionPath(s)
	}
}

func (s *MemberExpressionPathContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(FqlParserListener); ok {
		listenerT.ExitMemberExpressionPath(s)
	}
}

func (s *MemberExpressionPathContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case FqlParserVisitor:
		return t.VisitMemberExpressionPath(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *FqlParser) MemberExpressionPath() (localctx IMemberExpressionPathContext) {
	localctx = NewMemberExpressionPathContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 108, FqlParserRULE_memberExpressionPath)
	var _la int

	p.SetState(550)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 52, p.GetParserRuleContext()) {
	case 1:
		p.EnterOuterAlt(localctx, 1)
		p.SetState(540)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)

		if _la == FqlParserQuestionMark {
			{
				p.SetState(539)
				p.ErrorOperator()
			}

		}
		{
			p.SetState(542)
			p.Match(FqlParserDot)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(543)
			p.PropertyName()
		}

	case 2:
		p.EnterOuterAlt(localctx, 2)
		p.SetState(547)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)

		if _la == FqlParserQuestionMark {
			{
				p.SetState(544)
				p.ErrorOperator()
			}
			{
				p.SetState(545)
				p.Match(FqlParserDot)
				if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
				}
			}

		}
		{
			p.SetState(549)
			p.ComputedPropertyName()
		}

	case antlr.ATNInvalidAltNumber:
		goto errorExit
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// ISafeReservedWordContext is an interface to support dynamic dispatch.
type ISafeReservedWordContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	And() antlr.TerminalNode
	Or() antlr.TerminalNode
	Distinct() antlr.TerminalNode
	Filter() antlr.TerminalNode
	Sort() antlr.TerminalNode
	Limit() antlr.TerminalNode
	Collect() antlr.TerminalNode
	SortDirection() antlr.TerminalNode
	Into() antlr.TerminalNode
	Keep() antlr.TerminalNode
	With() antlr.TerminalNode
	All() antlr.TerminalNode
	Any() antlr.TerminalNode
	Aggregate() antlr.TerminalNode
	Event() antlr.TerminalNode
	Timeout() antlr.TerminalNode
	Options() antlr.TerminalNode
	Current() antlr.TerminalNode
	Step() antlr.TerminalNode

	// IsSafeReservedWordContext differentiates from other interfaces.
	IsSafeReservedWordContext()
}

type SafeReservedWordContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptySafeReservedWordContext() *SafeReservedWordContext {
	var p = new(SafeReservedWordContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = FqlParserRULE_safeReservedWord
	return p
}

func InitEmptySafeReservedWordContext(p *SafeReservedWordContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = FqlParserRULE_safeReservedWord
}

func (*SafeReservedWordContext) IsSafeReservedWordContext() {}

func NewSafeReservedWordContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *SafeReservedWordContext {
	var p = new(SafeReservedWordContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = FqlParserRULE_safeReservedWord

	return p
}

func (s *SafeReservedWordContext) GetParser() antlr.Parser { return s.parser }

func (s *SafeReservedWordContext) And() antlr.TerminalNode {
	return s.GetToken(FqlParserAnd, 0)
}

func (s *SafeReservedWordContext) Or() antlr.TerminalNode {
	return s.GetToken(FqlParserOr, 0)
}

func (s *SafeReservedWordContext) Distinct() antlr.TerminalNode {
	return s.GetToken(FqlParserDistinct, 0)
}

func (s *SafeReservedWordContext) Filter() antlr.TerminalNode {
	return s.GetToken(FqlParserFilter, 0)
}

func (s *SafeReservedWordContext) Sort() antlr.TerminalNode {
	return s.GetToken(FqlParserSort, 0)
}

func (s *SafeReservedWordContext) Limit() antlr.TerminalNode {
	return s.GetToken(FqlParserLimit, 0)
}

func (s *SafeReservedWordContext) Collect() antlr.TerminalNode {
	return s.GetToken(FqlParserCollect, 0)
}

func (s *SafeReservedWordContext) SortDirection() antlr.TerminalNode {
	return s.GetToken(FqlParserSortDirection, 0)
}

func (s *SafeReservedWordContext) Into() antlr.TerminalNode {
	return s.GetToken(FqlParserInto, 0)
}

func (s *SafeReservedWordContext) Keep() antlr.TerminalNode {
	return s.GetToken(FqlParserKeep, 0)
}

func (s *SafeReservedWordContext) With() antlr.TerminalNode {
	return s.GetToken(FqlParserWith, 0)
}

func (s *SafeReservedWordContext) All() antlr.TerminalNode {
	return s.GetToken(FqlParserAll, 0)
}

func (s *SafeReservedWordContext) Any() antlr.TerminalNode {
	return s.GetToken(FqlParserAny, 0)
}

func (s *SafeReservedWordContext) Aggregate() antlr.TerminalNode {
	return s.GetToken(FqlParserAggregate, 0)
}

func (s *SafeReservedWordContext) Event() antlr.TerminalNode {
	return s.GetToken(FqlParserEvent, 0)
}

func (s *SafeReservedWordContext) Timeout() antlr.TerminalNode {
	return s.GetToken(FqlParserTimeout, 0)
}

func (s *SafeReservedWordContext) Options() antlr.TerminalNode {
	return s.GetToken(FqlParserOptions, 0)
}

func (s *SafeReservedWordContext) Current() antlr.TerminalNode {
	return s.GetToken(FqlParserCurrent, 0)
}

func (s *SafeReservedWordContext) Step() antlr.TerminalNode {
	return s.GetToken(FqlParserStep, 0)
}

func (s *SafeReservedWordContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *SafeReservedWordContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *SafeReservedWordContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(FqlParserListener); ok {
		listenerT.EnterSafeReservedWord(s)
	}
}

func (s *SafeReservedWordContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(FqlParserListener); ok {
		listenerT.ExitSafeReservedWord(s)
	}
}

func (s *SafeReservedWordContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case FqlParserVisitor:
		return t.VisitSafeReservedWord(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *FqlParser) SafeReservedWord() (localctx ISafeReservedWordContext) {
	localctx = NewSafeReservedWordContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 110, FqlParserRULE_safeReservedWord)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(552)
		_la = p.GetTokenStream().LA(1)

		if !((int64(_la) & ^0x3f) == 0 && ((int64(1)<<_la)&4754736810729406464) != 0) {
			p.GetErrorHandler().RecoverInline(p)
		} else {
			p.GetErrorHandler().ReportMatch(p)
			p.Consume()
		}
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IUnsafeReservedWordContext is an interface to support dynamic dispatch.
type IUnsafeReservedWordContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	Return() antlr.TerminalNode
	None() antlr.TerminalNode
	Null() antlr.TerminalNode
	Let() antlr.TerminalNode
	Use() antlr.TerminalNode
	Waitfor() antlr.TerminalNode
	While() antlr.TerminalNode
	Do() antlr.TerminalNode
	In() antlr.TerminalNode
	Like() antlr.TerminalNode
	Not() antlr.TerminalNode
	For() antlr.TerminalNode
	BooleanLiteral() antlr.TerminalNode

	// IsUnsafeReservedWordContext differentiates from other interfaces.
	IsUnsafeReservedWordContext()
}

type UnsafeReservedWordContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyUnsafeReservedWordContext() *UnsafeReservedWordContext {
	var p = new(UnsafeReservedWordContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = FqlParserRULE_unsafeReservedWord
	return p
}

func InitEmptyUnsafeReservedWordContext(p *UnsafeReservedWordContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = FqlParserRULE_unsafeReservedWord
}

func (*UnsafeReservedWordContext) IsUnsafeReservedWordContext() {}

func NewUnsafeReservedWordContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *UnsafeReservedWordContext {
	var p = new(UnsafeReservedWordContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = FqlParserRULE_unsafeReservedWord

	return p
}

func (s *UnsafeReservedWordContext) GetParser() antlr.Parser { return s.parser }

func (s *UnsafeReservedWordContext) Return() antlr.TerminalNode {
	return s.GetToken(FqlParserReturn, 0)
}

func (s *UnsafeReservedWordContext) None() antlr.TerminalNode {
	return s.GetToken(FqlParserNone, 0)
}

func (s *UnsafeReservedWordContext) Null() antlr.TerminalNode {
	return s.GetToken(FqlParserNull, 0)
}

func (s *UnsafeReservedWordContext) Let() antlr.TerminalNode {
	return s.GetToken(FqlParserLet, 0)
}

func (s *UnsafeReservedWordContext) Use() antlr.TerminalNode {
	return s.GetToken(FqlParserUse, 0)
}

func (s *UnsafeReservedWordContext) Waitfor() antlr.TerminalNode {
	return s.GetToken(FqlParserWaitfor, 0)
}

func (s *UnsafeReservedWordContext) While() antlr.TerminalNode {
	return s.GetToken(FqlParserWhile, 0)
}

func (s *UnsafeReservedWordContext) Do() antlr.TerminalNode {
	return s.GetToken(FqlParserDo, 0)
}

func (s *UnsafeReservedWordContext) In() antlr.TerminalNode {
	return s.GetToken(FqlParserIn, 0)
}

func (s *UnsafeReservedWordContext) Like() antlr.TerminalNode {
	return s.GetToken(FqlParserLike, 0)
}

func (s *UnsafeReservedWordContext) Not() antlr.TerminalNode {
	return s.GetToken(FqlParserNot, 0)
}

func (s *UnsafeReservedWordContext) For() antlr.TerminalNode {
	return s.GetToken(FqlParserFor, 0)
}

func (s *UnsafeReservedWordContext) BooleanLiteral() antlr.TerminalNode {
	return s.GetToken(FqlParserBooleanLiteral, 0)
}

func (s *UnsafeReservedWordContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *UnsafeReservedWordContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *UnsafeReservedWordContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(FqlParserListener); ok {
		listenerT.EnterUnsafeReservedWord(s)
	}
}

func (s *UnsafeReservedWordContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(FqlParserListener); ok {
		listenerT.ExitUnsafeReservedWord(s)
	}
}

func (s *UnsafeReservedWordContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case FqlParserVisitor:
		return t.VisitUnsafeReservedWord(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *FqlParser) UnsafeReservedWord() (localctx IUnsafeReservedWordContext) {
	localctx = NewUnsafeReservedWordContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 112, FqlParserRULE_unsafeReservedWord)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(554)
		_la = p.GetTokenStream().LA(1)

		if !((int64(_la) & ^0x3f) == 0 && ((int64(1)<<_la)&4468635217736761344) != 0) {
			p.GetErrorHandler().RecoverInline(p)
		} else {
			p.GetErrorHandler().ReportMatch(p)
			p.Consume()
		}
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IRangeOperatorContext is an interface to support dynamic dispatch.
type IRangeOperatorContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// GetLeft returns the left rule contexts.
	GetLeft() IRangeOperandContext

	// GetRight returns the right rule contexts.
	GetRight() IRangeOperandContext

	// SetLeft sets the left rule contexts.
	SetLeft(IRangeOperandContext)

	// SetRight sets the right rule contexts.
	SetRight(IRangeOperandContext)

	// Getter signatures
	Range() antlr.TerminalNode
	AllRangeOperand() []IRangeOperandContext
	RangeOperand(i int) IRangeOperandContext

	// IsRangeOperatorContext differentiates from other interfaces.
	IsRangeOperatorContext()
}

type RangeOperatorContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
	left   IRangeOperandContext
	right  IRangeOperandContext
}

func NewEmptyRangeOperatorContext() *RangeOperatorContext {
	var p = new(RangeOperatorContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = FqlParserRULE_rangeOperator
	return p
}

func InitEmptyRangeOperatorContext(p *RangeOperatorContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = FqlParserRULE_rangeOperator
}

func (*RangeOperatorContext) IsRangeOperatorContext() {}

func NewRangeOperatorContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *RangeOperatorContext {
	var p = new(RangeOperatorContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = FqlParserRULE_rangeOperator

	return p
}

func (s *RangeOperatorContext) GetParser() antlr.Parser { return s.parser }

func (s *RangeOperatorContext) GetLeft() IRangeOperandContext { return s.left }

func (s *RangeOperatorContext) GetRight() IRangeOperandContext { return s.right }

func (s *RangeOperatorContext) SetLeft(v IRangeOperandContext) { s.left = v }

func (s *RangeOperatorContext) SetRight(v IRangeOperandContext) { s.right = v }

func (s *RangeOperatorContext) Range() antlr.TerminalNode {
	return s.GetToken(FqlParserRange, 0)
}

func (s *RangeOperatorContext) AllRangeOperand() []IRangeOperandContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(IRangeOperandContext); ok {
			len++
		}
	}

	tst := make([]IRangeOperandContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(IRangeOperandContext); ok {
			tst[i] = t.(IRangeOperandContext)
			i++
		}
	}

	return tst
}

func (s *RangeOperatorContext) RangeOperand(i int) IRangeOperandContext {
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IRangeOperandContext); ok {
			if j == i {
				t = ctx.(antlr.RuleContext)
				break
			}
			j++
		}
	}

	if t == nil {
		return nil
	}

	return t.(IRangeOperandContext)
}

func (s *RangeOperatorContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *RangeOperatorContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *RangeOperatorContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(FqlParserListener); ok {
		listenerT.EnterRangeOperator(s)
	}
}

func (s *RangeOperatorContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(FqlParserListener); ok {
		listenerT.ExitRangeOperator(s)
	}
}

func (s *RangeOperatorContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case FqlParserVisitor:
		return t.VisitRangeOperator(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *FqlParser) RangeOperator() (localctx IRangeOperatorContext) {
	localctx = NewRangeOperatorContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 114, FqlParserRULE_rangeOperator)
	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(556)

		var _x = p.RangeOperand()

		localctx.(*RangeOperatorContext).left = _x
	}
	{
		p.SetState(557)
		p.Match(FqlParserRange)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(558)

		var _x = p.RangeOperand()

		localctx.(*RangeOperatorContext).right = _x
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IRangeOperandContext is an interface to support dynamic dispatch.
type IRangeOperandContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	IntegerLiteral() IIntegerLiteralContext
	Variable() IVariableContext
	Param() IParamContext

	// IsRangeOperandContext differentiates from other interfaces.
	IsRangeOperandContext()
}

type RangeOperandContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyRangeOperandContext() *RangeOperandContext {
	var p = new(RangeOperandContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = FqlParserRULE_rangeOperand
	return p
}

func InitEmptyRangeOperandContext(p *RangeOperandContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = FqlParserRULE_rangeOperand
}

func (*RangeOperandContext) IsRangeOperandContext() {}

func NewRangeOperandContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *RangeOperandContext {
	var p = new(RangeOperandContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = FqlParserRULE_rangeOperand

	return p
}

func (s *RangeOperandContext) GetParser() antlr.Parser { return s.parser }

func (s *RangeOperandContext) IntegerLiteral() IIntegerLiteralContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IIntegerLiteralContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IIntegerLiteralContext)
}

func (s *RangeOperandContext) Variable() IVariableContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IVariableContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IVariableContext)
}

func (s *RangeOperandContext) Param() IParamContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IParamContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IParamContext)
}

func (s *RangeOperandContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *RangeOperandContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *RangeOperandContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(FqlParserListener); ok {
		listenerT.EnterRangeOperand(s)
	}
}

func (s *RangeOperandContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(FqlParserListener); ok {
		listenerT.ExitRangeOperand(s)
	}
}

func (s *RangeOperandContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case FqlParserVisitor:
		return t.VisitRangeOperand(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *FqlParser) RangeOperand() (localctx IRangeOperandContext) {
	localctx = NewRangeOperandContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 116, FqlParserRULE_rangeOperand)
	p.SetState(563)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetTokenStream().LA(1) {
	case FqlParserIntegerLiteral:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(560)
			p.IntegerLiteral()
		}

	case FqlParserAnd, FqlParserOr, FqlParserOptions, FqlParserTimeout, FqlParserDistinct, FqlParserFilter, FqlParserCurrent, FqlParserSort, FqlParserLimit, FqlParserCollect, FqlParserSortDirection, FqlParserInto, FqlParserKeep, FqlParserWith, FqlParserAll, FqlParserAny, FqlParserAggregate, FqlParserEvent, FqlParserStep, FqlParserIdentifier:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(561)
			p.Variable()
		}

	case FqlParserParam:
		p.EnterOuterAlt(localctx, 3)
		{
			p.SetState(562)
			p.Param()
		}

	default:
		p.SetError(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
		goto errorExit
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IExpressionContext is an interface to support dynamic dispatch.
type IExpressionContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// GetTernaryOperator returns the ternaryOperator token.
	GetTernaryOperator() antlr.Token

	// SetTernaryOperator sets the ternaryOperator token.
	SetTernaryOperator(antlr.Token)

	// GetLeft returns the left rule contexts.
	GetLeft() IExpressionContext

	// GetCondition returns the condition rule contexts.
	GetCondition() IExpressionContext

	// GetRight returns the right rule contexts.
	GetRight() IExpressionContext

	// GetOnTrue returns the onTrue rule contexts.
	GetOnTrue() IExpressionContext

	// GetOnFalse returns the onFalse rule contexts.
	GetOnFalse() IExpressionContext

	// SetLeft sets the left rule contexts.
	SetLeft(IExpressionContext)

	// SetCondition sets the condition rule contexts.
	SetCondition(IExpressionContext)

	// SetRight sets the right rule contexts.
	SetRight(IExpressionContext)

	// SetOnTrue sets the onTrue rule contexts.
	SetOnTrue(IExpressionContext)

	// SetOnFalse sets the onFalse rule contexts.
	SetOnFalse(IExpressionContext)

	// Getter signatures
	UnaryOperator() IUnaryOperatorContext
	AllExpression() []IExpressionContext
	Expression(i int) IExpressionContext
	Predicate() IPredicateContext
	LogicalAndOperator() ILogicalAndOperatorContext
	LogicalOrOperator() ILogicalOrOperatorContext
	Colon() antlr.TerminalNode
	QuestionMark() antlr.TerminalNode

	// IsExpressionContext differentiates from other interfaces.
	IsExpressionContext()
}

type ExpressionContext struct {
	antlr.BaseParserRuleContext
	parser          antlr.Parser
	left            IExpressionContext
	condition       IExpressionContext
	right           IExpressionContext
	ternaryOperator antlr.Token
	onTrue          IExpressionContext
	onFalse         IExpressionContext
}

func NewEmptyExpressionContext() *ExpressionContext {
	var p = new(ExpressionContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = FqlParserRULE_expression
	return p
}

func InitEmptyExpressionContext(p *ExpressionContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = FqlParserRULE_expression
}

func (*ExpressionContext) IsExpressionContext() {}

func NewExpressionContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ExpressionContext {
	var p = new(ExpressionContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = FqlParserRULE_expression

	return p
}

func (s *ExpressionContext) GetParser() antlr.Parser { return s.parser }

func (s *ExpressionContext) GetTernaryOperator() antlr.Token { return s.ternaryOperator }

func (s *ExpressionContext) SetTernaryOperator(v antlr.Token) { s.ternaryOperator = v }

func (s *ExpressionContext) GetLeft() IExpressionContext { return s.left }

func (s *ExpressionContext) GetCondition() IExpressionContext { return s.condition }

func (s *ExpressionContext) GetRight() IExpressionContext { return s.right }

func (s *ExpressionContext) GetOnTrue() IExpressionContext { return s.onTrue }

func (s *ExpressionContext) GetOnFalse() IExpressionContext { return s.onFalse }

func (s *ExpressionContext) SetLeft(v IExpressionContext) { s.left = v }

func (s *ExpressionContext) SetCondition(v IExpressionContext) { s.condition = v }

func (s *ExpressionContext) SetRight(v IExpressionContext) { s.right = v }

func (s *ExpressionContext) SetOnTrue(v IExpressionContext) { s.onTrue = v }

func (s *ExpressionContext) SetOnFalse(v IExpressionContext) { s.onFalse = v }

func (s *ExpressionContext) UnaryOperator() IUnaryOperatorContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IUnaryOperatorContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IUnaryOperatorContext)
}

func (s *ExpressionContext) AllExpression() []IExpressionContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(IExpressionContext); ok {
			len++
		}
	}

	tst := make([]IExpressionContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(IExpressionContext); ok {
			tst[i] = t.(IExpressionContext)
			i++
		}
	}

	return tst
}

func (s *ExpressionContext) Expression(i int) IExpressionContext {
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IExpressionContext); ok {
			if j == i {
				t = ctx.(antlr.RuleContext)
				break
			}
			j++
		}
	}

	if t == nil {
		return nil
	}

	return t.(IExpressionContext)
}

func (s *ExpressionContext) Predicate() IPredicateContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IPredicateContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IPredicateContext)
}

func (s *ExpressionContext) LogicalAndOperator() ILogicalAndOperatorContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ILogicalAndOperatorContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(ILogicalAndOperatorContext)
}

func (s *ExpressionContext) LogicalOrOperator() ILogicalOrOperatorContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ILogicalOrOperatorContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(ILogicalOrOperatorContext)
}

func (s *ExpressionContext) Colon() antlr.TerminalNode {
	return s.GetToken(FqlParserColon, 0)
}

func (s *ExpressionContext) QuestionMark() antlr.TerminalNode {
	return s.GetToken(FqlParserQuestionMark, 0)
}

func (s *ExpressionContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ExpressionContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *ExpressionContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(FqlParserListener); ok {
		listenerT.EnterExpression(s)
	}
}

func (s *ExpressionContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(FqlParserListener); ok {
		listenerT.ExitExpression(s)
	}
}

func (s *ExpressionContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case FqlParserVisitor:
		return t.VisitExpression(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *FqlParser) Expression() (localctx IExpressionContext) {
	return p.expression(0)
}

func (p *FqlParser) expression(_p int) (localctx IExpressionContext) {
	var _parentctx antlr.ParserRuleContext = p.GetParserRuleContext()

	_parentState := p.GetState()
	localctx = NewExpressionContext(p, p.GetParserRuleContext(), _parentState)
	var _prevctx IExpressionContext = localctx
	var _ antlr.ParserRuleContext = _prevctx // TODO: To prevent unused variable warning.
	_startState := 118
	p.EnterRecursionRule(localctx, 118, FqlParserRULE_expression, _p)
	var _la int

	var _alt int

	p.EnterOuterAlt(localctx, 1)
	p.SetState(570)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 54, p.GetParserRuleContext()) {
	case 1:
		{
			p.SetState(566)
			p.UnaryOperator()
		}
		{
			p.SetState(567)

			var _x = p.expression(5)

			localctx.(*ExpressionContext).right = _x
		}

	case 2:
		{
			p.SetState(569)
			p.predicate(0)
		}

	case antlr.ATNInvalidAltNumber:
		goto errorExit
	}
	p.GetParserRuleContext().SetStop(p.GetTokenStream().LT(-1))
	p.SetState(589)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_alt = p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 57, p.GetParserRuleContext())
	if p.HasError() {
		goto errorExit
	}
	for _alt != 2 && _alt != antlr.ATNInvalidAltNumber {
		if _alt == 1 {
			if p.GetParseListeners() != nil {
				p.TriggerExitRuleEvent()
			}
			_prevctx = localctx
			p.SetState(587)
			p.GetErrorHandler().Sync(p)
			if p.HasError() {
				goto errorExit
			}

			switch p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 56, p.GetParserRuleContext()) {
			case 1:
				localctx = NewExpressionContext(p, _parentctx, _parentState)
				localctx.(*ExpressionContext).left = _prevctx
				p.PushNewRecursionContext(localctx, _startState, FqlParserRULE_expression)
				p.SetState(572)

				if !(p.Precpred(p.GetParserRuleContext(), 4)) {
					p.SetError(antlr.NewFailedPredicateException(p, "p.Precpred(p.GetParserRuleContext(), 4)", ""))
					goto errorExit
				}
				{
					p.SetState(573)
					p.LogicalAndOperator()
				}
				{
					p.SetState(574)

					var _x = p.expression(5)

					localctx.(*ExpressionContext).right = _x
				}

			case 2:
				localctx = NewExpressionContext(p, _parentctx, _parentState)
				localctx.(*ExpressionContext).left = _prevctx
				p.PushNewRecursionContext(localctx, _startState, FqlParserRULE_expression)
				p.SetState(576)

				if !(p.Precpred(p.GetParserRuleContext(), 3)) {
					p.SetError(antlr.NewFailedPredicateException(p, "p.Precpred(p.GetParserRuleContext(), 3)", ""))
					goto errorExit
				}
				{
					p.SetState(577)
					p.LogicalOrOperator()
				}
				{
					p.SetState(578)

					var _x = p.expression(4)

					localctx.(*ExpressionContext).right = _x
				}

			case 3:
				localctx = NewExpressionContext(p, _parentctx, _parentState)
				localctx.(*ExpressionContext).condition = _prevctx
				p.PushNewRecursionContext(localctx, _startState, FqlParserRULE_expression)
				p.SetState(580)

				if !(p.Precpred(p.GetParserRuleContext(), 2)) {
					p.SetError(antlr.NewFailedPredicateException(p, "p.Precpred(p.GetParserRuleContext(), 2)", ""))
					goto errorExit
				}
				{
					p.SetState(581)

					var _m = p.Match(FqlParserQuestionMark)

					localctx.(*ExpressionContext).ternaryOperator = _m
					if p.HasError() {
						// Recognition error - abort rule
						goto errorExit
					}
				}
				p.SetState(583)
				p.GetErrorHandler().Sync(p)
				if p.HasError() {
					goto errorExit
				}
				_la = p.GetTokenStream().LA(1)

				if (int64((_la-9)) & ^0x3f) == 0 && ((int64(1)<<(_la-9))&2233785415159480341) != 0 {
					{
						p.SetState(582)

						var _x = p.expression(0)

						localctx.(*ExpressionContext).onTrue = _x
					}

				}
				{
					p.SetState(585)
					p.Match(FqlParserColon)
					if p.HasError() {
						// Recognition error - abort rule
						goto errorExit
					}
				}
				{
					p.SetState(586)

					var _x = p.expression(3)

					localctx.(*ExpressionContext).onFalse = _x
				}

			case antlr.ATNInvalidAltNumber:
				goto errorExit
			}

		}
		p.SetState(591)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_alt = p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 57, p.GetParserRuleContext())
		if p.HasError() {
			goto errorExit
		}
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.UnrollRecursionContexts(_parentctx)
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IPredicateContext is an interface to support dynamic dispatch.
type IPredicateContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// GetLeft returns the left rule contexts.
	GetLeft() IPredicateContext

	// GetRight returns the right rule contexts.
	GetRight() IPredicateContext

	// SetLeft sets the left rule contexts.
	SetLeft(IPredicateContext)

	// SetRight sets the right rule contexts.
	SetRight(IPredicateContext)

	// Getter signatures
	ExpressionAtom() IExpressionAtomContext
	EqualityOperator() IEqualityOperatorContext
	AllPredicate() []IPredicateContext
	Predicate(i int) IPredicateContext
	ArrayOperator() IArrayOperatorContext
	InOperator() IInOperatorContext
	LikeOperator() ILikeOperatorContext

	// IsPredicateContext differentiates from other interfaces.
	IsPredicateContext()
}

type PredicateContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
	left   IPredicateContext
	right  IPredicateContext
}

func NewEmptyPredicateContext() *PredicateContext {
	var p = new(PredicateContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = FqlParserRULE_predicate
	return p
}

func InitEmptyPredicateContext(p *PredicateContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = FqlParserRULE_predicate
}

func (*PredicateContext) IsPredicateContext() {}

func NewPredicateContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *PredicateContext {
	var p = new(PredicateContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = FqlParserRULE_predicate

	return p
}

func (s *PredicateContext) GetParser() antlr.Parser { return s.parser }

func (s *PredicateContext) GetLeft() IPredicateContext { return s.left }

func (s *PredicateContext) GetRight() IPredicateContext { return s.right }

func (s *PredicateContext) SetLeft(v IPredicateContext) { s.left = v }

func (s *PredicateContext) SetRight(v IPredicateContext) { s.right = v }

func (s *PredicateContext) ExpressionAtom() IExpressionAtomContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IExpressionAtomContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IExpressionAtomContext)
}

func (s *PredicateContext) EqualityOperator() IEqualityOperatorContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IEqualityOperatorContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IEqualityOperatorContext)
}

func (s *PredicateContext) AllPredicate() []IPredicateContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(IPredicateContext); ok {
			len++
		}
	}

	tst := make([]IPredicateContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(IPredicateContext); ok {
			tst[i] = t.(IPredicateContext)
			i++
		}
	}

	return tst
}

func (s *PredicateContext) Predicate(i int) IPredicateContext {
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IPredicateContext); ok {
			if j == i {
				t = ctx.(antlr.RuleContext)
				break
			}
			j++
		}
	}

	if t == nil {
		return nil
	}

	return t.(IPredicateContext)
}

func (s *PredicateContext) ArrayOperator() IArrayOperatorContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IArrayOperatorContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IArrayOperatorContext)
}

func (s *PredicateContext) InOperator() IInOperatorContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IInOperatorContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IInOperatorContext)
}

func (s *PredicateContext) LikeOperator() ILikeOperatorContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ILikeOperatorContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(ILikeOperatorContext)
}

func (s *PredicateContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *PredicateContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *PredicateContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(FqlParserListener); ok {
		listenerT.EnterPredicate(s)
	}
}

func (s *PredicateContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(FqlParserListener); ok {
		listenerT.ExitPredicate(s)
	}
}

func (s *PredicateContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case FqlParserVisitor:
		return t.VisitPredicate(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *FqlParser) Predicate() (localctx IPredicateContext) {
	return p.predicate(0)
}

func (p *FqlParser) predicate(_p int) (localctx IPredicateContext) {
	var _parentctx antlr.ParserRuleContext = p.GetParserRuleContext()

	_parentState := p.GetState()
	localctx = NewPredicateContext(p, p.GetParserRuleContext(), _parentState)
	var _prevctx IPredicateContext = localctx
	var _ antlr.ParserRuleContext = _prevctx // TODO: To prevent unused variable warning.
	_startState := 120
	p.EnterRecursionRule(localctx, 120, FqlParserRULE_predicate, _p)
	var _alt int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(593)
		p.expressionAtom(0)
	}

	p.GetParserRuleContext().SetStop(p.GetTokenStream().LT(-1))
	p.SetState(613)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_alt = p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 59, p.GetParserRuleContext())
	if p.HasError() {
		goto errorExit
	}
	for _alt != 2 && _alt != antlr.ATNInvalidAltNumber {
		if _alt == 1 {
			if p.GetParseListeners() != nil {
				p.TriggerExitRuleEvent()
			}
			_prevctx = localctx
			p.SetState(611)
			p.GetErrorHandler().Sync(p)
			if p.HasError() {
				goto errorExit
			}

			switch p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 58, p.GetParserRuleContext()) {
			case 1:
				localctx = NewPredicateContext(p, _parentctx, _parentState)
				localctx.(*PredicateContext).left = _prevctx
				p.PushNewRecursionContext(localctx, _startState, FqlParserRULE_predicate)
				p.SetState(595)

				if !(p.Precpred(p.GetParserRuleContext(), 5)) {
					p.SetError(antlr.NewFailedPredicateException(p, "p.Precpred(p.GetParserRuleContext(), 5)", ""))
					goto errorExit
				}
				{
					p.SetState(596)
					p.EqualityOperator()
				}
				{
					p.SetState(597)

					var _x = p.predicate(6)

					localctx.(*PredicateContext).right = _x
				}

			case 2:
				localctx = NewPredicateContext(p, _parentctx, _parentState)
				localctx.(*PredicateContext).left = _prevctx
				p.PushNewRecursionContext(localctx, _startState, FqlParserRULE_predicate)
				p.SetState(599)

				if !(p.Precpred(p.GetParserRuleContext(), 4)) {
					p.SetError(antlr.NewFailedPredicateException(p, "p.Precpred(p.GetParserRuleContext(), 4)", ""))
					goto errorExit
				}
				{
					p.SetState(600)
					p.ArrayOperator()
				}
				{
					p.SetState(601)

					var _x = p.predicate(5)

					localctx.(*PredicateContext).right = _x
				}

			case 3:
				localctx = NewPredicateContext(p, _parentctx, _parentState)
				localctx.(*PredicateContext).left = _prevctx
				p.PushNewRecursionContext(localctx, _startState, FqlParserRULE_predicate)
				p.SetState(603)

				if !(p.Precpred(p.GetParserRuleContext(), 3)) {
					p.SetError(antlr.NewFailedPredicateException(p, "p.Precpred(p.GetParserRuleContext(), 3)", ""))
					goto errorExit
				}
				{
					p.SetState(604)
					p.InOperator()
				}
				{
					p.SetState(605)

					var _x = p.predicate(4)

					localctx.(*PredicateContext).right = _x
				}

			case 4:
				localctx = NewPredicateContext(p, _parentctx, _parentState)
				localctx.(*PredicateContext).left = _prevctx
				p.PushNewRecursionContext(localctx, _startState, FqlParserRULE_predicate)
				p.SetState(607)

				if !(p.Precpred(p.GetParserRuleContext(), 2)) {
					p.SetError(antlr.NewFailedPredicateException(p, "p.Precpred(p.GetParserRuleContext(), 2)", ""))
					goto errorExit
				}
				{
					p.SetState(608)
					p.LikeOperator()
				}
				{
					p.SetState(609)

					var _x = p.predicate(3)

					localctx.(*PredicateContext).right = _x
				}

			case antlr.ATNInvalidAltNumber:
				goto errorExit
			}

		}
		p.SetState(615)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_alt = p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 59, p.GetParserRuleContext())
		if p.HasError() {
			goto errorExit
		}
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.UnrollRecursionContexts(_parentctx)
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IExpressionAtomContext is an interface to support dynamic dispatch.
type IExpressionAtomContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// GetLeft returns the left rule contexts.
	GetLeft() IExpressionAtomContext

	// GetRight returns the right rule contexts.
	GetRight() IExpressionAtomContext

	// SetLeft sets the left rule contexts.
	SetLeft(IExpressionAtomContext)

	// SetRight sets the right rule contexts.
	SetRight(IExpressionAtomContext)

	// Getter signatures
	FunctionCallExpression() IFunctionCallExpressionContext
	RangeOperator() IRangeOperatorContext
	Literal() ILiteralContext
	Variable() IVariableContext
	MemberExpression() IMemberExpressionContext
	Param() IParamContext
	OpenParen() antlr.TerminalNode
	CloseParen() antlr.TerminalNode
	ForExpression() IForExpressionContext
	WaitForExpression() IWaitForExpressionContext
	Expression() IExpressionContext
	ErrorOperator() IErrorOperatorContext
	MultiplicativeOperator() IMultiplicativeOperatorContext
	AllExpressionAtom() []IExpressionAtomContext
	ExpressionAtom(i int) IExpressionAtomContext
	AdditiveOperator() IAdditiveOperatorContext
	RegexpOperator() IRegexpOperatorContext

	// IsExpressionAtomContext differentiates from other interfaces.
	IsExpressionAtomContext()
}

type ExpressionAtomContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
	left   IExpressionAtomContext
	right  IExpressionAtomContext
}

func NewEmptyExpressionAtomContext() *ExpressionAtomContext {
	var p = new(ExpressionAtomContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = FqlParserRULE_expressionAtom
	return p
}

func InitEmptyExpressionAtomContext(p *ExpressionAtomContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = FqlParserRULE_expressionAtom
}

func (*ExpressionAtomContext) IsExpressionAtomContext() {}

func NewExpressionAtomContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ExpressionAtomContext {
	var p = new(ExpressionAtomContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = FqlParserRULE_expressionAtom

	return p
}

func (s *ExpressionAtomContext) GetParser() antlr.Parser { return s.parser }

func (s *ExpressionAtomContext) GetLeft() IExpressionAtomContext { return s.left }

func (s *ExpressionAtomContext) GetRight() IExpressionAtomContext { return s.right }

func (s *ExpressionAtomContext) SetLeft(v IExpressionAtomContext) { s.left = v }

func (s *ExpressionAtomContext) SetRight(v IExpressionAtomContext) { s.right = v }

func (s *ExpressionAtomContext) FunctionCallExpression() IFunctionCallExpressionContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IFunctionCallExpressionContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IFunctionCallExpressionContext)
}

func (s *ExpressionAtomContext) RangeOperator() IRangeOperatorContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IRangeOperatorContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IRangeOperatorContext)
}

func (s *ExpressionAtomContext) Literal() ILiteralContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ILiteralContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(ILiteralContext)
}

func (s *ExpressionAtomContext) Variable() IVariableContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IVariableContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IVariableContext)
}

func (s *ExpressionAtomContext) MemberExpression() IMemberExpressionContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IMemberExpressionContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IMemberExpressionContext)
}

func (s *ExpressionAtomContext) Param() IParamContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IParamContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IParamContext)
}

func (s *ExpressionAtomContext) OpenParen() antlr.TerminalNode {
	return s.GetToken(FqlParserOpenParen, 0)
}

func (s *ExpressionAtomContext) CloseParen() antlr.TerminalNode {
	return s.GetToken(FqlParserCloseParen, 0)
}

func (s *ExpressionAtomContext) ForExpression() IForExpressionContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IForExpressionContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IForExpressionContext)
}

func (s *ExpressionAtomContext) WaitForExpression() IWaitForExpressionContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IWaitForExpressionContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IWaitForExpressionContext)
}

func (s *ExpressionAtomContext) Expression() IExpressionContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IExpressionContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IExpressionContext)
}

func (s *ExpressionAtomContext) ErrorOperator() IErrorOperatorContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IErrorOperatorContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IErrorOperatorContext)
}

func (s *ExpressionAtomContext) MultiplicativeOperator() IMultiplicativeOperatorContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IMultiplicativeOperatorContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IMultiplicativeOperatorContext)
}

func (s *ExpressionAtomContext) AllExpressionAtom() []IExpressionAtomContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(IExpressionAtomContext); ok {
			len++
		}
	}

	tst := make([]IExpressionAtomContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(IExpressionAtomContext); ok {
			tst[i] = t.(IExpressionAtomContext)
			i++
		}
	}

	return tst
}

func (s *ExpressionAtomContext) ExpressionAtom(i int) IExpressionAtomContext {
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IExpressionAtomContext); ok {
			if j == i {
				t = ctx.(antlr.RuleContext)
				break
			}
			j++
		}
	}

	if t == nil {
		return nil
	}

	return t.(IExpressionAtomContext)
}

func (s *ExpressionAtomContext) AdditiveOperator() IAdditiveOperatorContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IAdditiveOperatorContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IAdditiveOperatorContext)
}

func (s *ExpressionAtomContext) RegexpOperator() IRegexpOperatorContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IRegexpOperatorContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IRegexpOperatorContext)
}

func (s *ExpressionAtomContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ExpressionAtomContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *ExpressionAtomContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(FqlParserListener); ok {
		listenerT.EnterExpressionAtom(s)
	}
}

func (s *ExpressionAtomContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(FqlParserListener); ok {
		listenerT.ExitExpressionAtom(s)
	}
}

func (s *ExpressionAtomContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case FqlParserVisitor:
		return t.VisitExpressionAtom(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *FqlParser) ExpressionAtom() (localctx IExpressionAtomContext) {
	return p.expressionAtom(0)
}

func (p *FqlParser) expressionAtom(_p int) (localctx IExpressionAtomContext) {
	var _parentctx antlr.ParserRuleContext = p.GetParserRuleContext()

	_parentState := p.GetState()
	localctx = NewExpressionAtomContext(p, p.GetParserRuleContext(), _parentState)
	var _prevctx IExpressionAtomContext = localctx
	var _ antlr.ParserRuleContext = _prevctx // TODO: To prevent unused variable warning.
	_startState := 122
	p.EnterRecursionRule(localctx, 122, FqlParserRULE_expressionAtom, _p)
	var _alt int

	p.EnterOuterAlt(localctx, 1)
	p.SetState(633)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 62, p.GetParserRuleContext()) {
	case 1:
		{
			p.SetState(617)
			p.FunctionCallExpression()
		}

	case 2:
		{
			p.SetState(618)
			p.RangeOperator()
		}

	case 3:
		{
			p.SetState(619)
			p.Literal()
		}

	case 4:
		{
			p.SetState(620)
			p.Variable()
		}

	case 5:
		{
			p.SetState(621)
			p.MemberExpression()
		}

	case 6:
		{
			p.SetState(622)
			p.Param()
		}

	case 7:
		{
			p.SetState(623)
			p.Match(FqlParserOpenParen)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		p.SetState(627)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}

		switch p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 60, p.GetParserRuleContext()) {
		case 1:
			{
				p.SetState(624)
				p.ForExpression()
			}

		case 2:
			{
				p.SetState(625)
				p.WaitForExpression()
			}

		case 3:
			{
				p.SetState(626)
				p.expression(0)
			}

		case antlr.ATNInvalidAltNumber:
			goto errorExit
		}
		{
			p.SetState(629)
			p.Match(FqlParserCloseParen)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		p.SetState(631)
		p.GetErrorHandler().Sync(p)

		if p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 61, p.GetParserRuleContext()) == 1 {
			{
				p.SetState(630)
				p.ErrorOperator()
			}

		} else if p.HasError() { // JIM
			goto errorExit
		}

	case antlr.ATNInvalidAltNumber:
		goto errorExit
	}
	p.GetParserRuleContext().SetStop(p.GetTokenStream().LT(-1))
	p.SetState(649)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_alt = p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 64, p.GetParserRuleContext())
	if p.HasError() {
		goto errorExit
	}
	for _alt != 2 && _alt != antlr.ATNInvalidAltNumber {
		if _alt == 1 {
			if p.GetParseListeners() != nil {
				p.TriggerExitRuleEvent()
			}
			_prevctx = localctx
			p.SetState(647)
			p.GetErrorHandler().Sync(p)
			if p.HasError() {
				goto errorExit
			}

			switch p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 63, p.GetParserRuleContext()) {
			case 1:
				localctx = NewExpressionAtomContext(p, _parentctx, _parentState)
				localctx.(*ExpressionAtomContext).left = _prevctx
				p.PushNewRecursionContext(localctx, _startState, FqlParserRULE_expressionAtom)
				p.SetState(635)

				if !(p.Precpred(p.GetParserRuleContext(), 10)) {
					p.SetError(antlr.NewFailedPredicateException(p, "p.Precpred(p.GetParserRuleContext(), 10)", ""))
					goto errorExit
				}
				{
					p.SetState(636)
					p.MultiplicativeOperator()
				}
				{
					p.SetState(637)

					var _x = p.expressionAtom(11)

					localctx.(*ExpressionAtomContext).right = _x
				}

			case 2:
				localctx = NewExpressionAtomContext(p, _parentctx, _parentState)
				localctx.(*ExpressionAtomContext).left = _prevctx
				p.PushNewRecursionContext(localctx, _startState, FqlParserRULE_expressionAtom)
				p.SetState(639)

				if !(p.Precpred(p.GetParserRuleContext(), 9)) {
					p.SetError(antlr.NewFailedPredicateException(p, "p.Precpred(p.GetParserRuleContext(), 9)", ""))
					goto errorExit
				}
				{
					p.SetState(640)
					p.AdditiveOperator()
				}
				{
					p.SetState(641)

					var _x = p.expressionAtom(10)

					localctx.(*ExpressionAtomContext).right = _x
				}

			case 3:
				localctx = NewExpressionAtomContext(p, _parentctx, _parentState)
				localctx.(*ExpressionAtomContext).left = _prevctx
				p.PushNewRecursionContext(localctx, _startState, FqlParserRULE_expressionAtom)
				p.SetState(643)

				if !(p.Precpred(p.GetParserRuleContext(), 8)) {
					p.SetError(antlr.NewFailedPredicateException(p, "p.Precpred(p.GetParserRuleContext(), 8)", ""))
					goto errorExit
				}
				{
					p.SetState(644)
					p.RegexpOperator()
				}
				{
					p.SetState(645)

					var _x = p.expressionAtom(9)

					localctx.(*ExpressionAtomContext).right = _x
				}

			case antlr.ATNInvalidAltNumber:
				goto errorExit
			}

		}
		p.SetState(651)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_alt = p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 64, p.GetParserRuleContext())
		if p.HasError() {
			goto errorExit
		}
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.UnrollRecursionContexts(_parentctx)
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IArrayOperatorContext is an interface to support dynamic dispatch.
type IArrayOperatorContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// GetOperator returns the operator token.
	GetOperator() antlr.Token

	// SetOperator sets the operator token.
	SetOperator(antlr.Token)

	// Getter signatures
	All() antlr.TerminalNode
	Any() antlr.TerminalNode
	None() antlr.TerminalNode
	InOperator() IInOperatorContext
	EqualityOperator() IEqualityOperatorContext

	// IsArrayOperatorContext differentiates from other interfaces.
	IsArrayOperatorContext()
}

type ArrayOperatorContext struct {
	antlr.BaseParserRuleContext
	parser   antlr.Parser
	operator antlr.Token
}

func NewEmptyArrayOperatorContext() *ArrayOperatorContext {
	var p = new(ArrayOperatorContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = FqlParserRULE_arrayOperator
	return p
}

func InitEmptyArrayOperatorContext(p *ArrayOperatorContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = FqlParserRULE_arrayOperator
}

func (*ArrayOperatorContext) IsArrayOperatorContext() {}

func NewArrayOperatorContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ArrayOperatorContext {
	var p = new(ArrayOperatorContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = FqlParserRULE_arrayOperator

	return p
}

func (s *ArrayOperatorContext) GetParser() antlr.Parser { return s.parser }

func (s *ArrayOperatorContext) GetOperator() antlr.Token { return s.operator }

func (s *ArrayOperatorContext) SetOperator(v antlr.Token) { s.operator = v }

func (s *ArrayOperatorContext) All() antlr.TerminalNode {
	return s.GetToken(FqlParserAll, 0)
}

func (s *ArrayOperatorContext) Any() antlr.TerminalNode {
	return s.GetToken(FqlParserAny, 0)
}

func (s *ArrayOperatorContext) None() antlr.TerminalNode {
	return s.GetToken(FqlParserNone, 0)
}

func (s *ArrayOperatorContext) InOperator() IInOperatorContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IInOperatorContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IInOperatorContext)
}

func (s *ArrayOperatorContext) EqualityOperator() IEqualityOperatorContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IEqualityOperatorContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IEqualityOperatorContext)
}

func (s *ArrayOperatorContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ArrayOperatorContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *ArrayOperatorContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(FqlParserListener); ok {
		listenerT.EnterArrayOperator(s)
	}
}

func (s *ArrayOperatorContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(FqlParserListener); ok {
		listenerT.ExitArrayOperator(s)
	}
}

func (s *ArrayOperatorContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case FqlParserVisitor:
		return t.VisitArrayOperator(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *FqlParser) ArrayOperator() (localctx IArrayOperatorContext) {
	localctx = NewArrayOperatorContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 124, FqlParserRULE_arrayOperator)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(652)

		var _lt = p.GetTokenStream().LT(1)

		localctx.(*ArrayOperatorContext).operator = _lt

		_la = p.GetTokenStream().LA(1)

		if !((int64(_la) & ^0x3f) == 0 && ((int64(1)<<_la)&27091966508400640) != 0) {
			var _ri = p.GetErrorHandler().RecoverInline(p)

			localctx.(*ArrayOperatorContext).operator = _ri
		} else {
			p.GetErrorHandler().ReportMatch(p)
			p.Consume()
		}
	}
	p.SetState(655)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetTokenStream().LA(1) {
	case FqlParserNot, FqlParserIn:
		{
			p.SetState(653)
			p.InOperator()
		}

	case FqlParserGt, FqlParserLt, FqlParserEq, FqlParserGte, FqlParserLte, FqlParserNeq:
		{
			p.SetState(654)
			p.EqualityOperator()
		}

	default:
		p.SetError(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
		goto errorExit
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IEqualityOperatorContext is an interface to support dynamic dispatch.
type IEqualityOperatorContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	Gt() antlr.TerminalNode
	Lt() antlr.TerminalNode
	Eq() antlr.TerminalNode
	Gte() antlr.TerminalNode
	Lte() antlr.TerminalNode
	Neq() antlr.TerminalNode

	// IsEqualityOperatorContext differentiates from other interfaces.
	IsEqualityOperatorContext()
}

type EqualityOperatorContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyEqualityOperatorContext() *EqualityOperatorContext {
	var p = new(EqualityOperatorContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = FqlParserRULE_equalityOperator
	return p
}

func InitEmptyEqualityOperatorContext(p *EqualityOperatorContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = FqlParserRULE_equalityOperator
}

func (*EqualityOperatorContext) IsEqualityOperatorContext() {}

func NewEqualityOperatorContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *EqualityOperatorContext {
	var p = new(EqualityOperatorContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = FqlParserRULE_equalityOperator

	return p
}

func (s *EqualityOperatorContext) GetParser() antlr.Parser { return s.parser }

func (s *EqualityOperatorContext) Gt() antlr.TerminalNode {
	return s.GetToken(FqlParserGt, 0)
}

func (s *EqualityOperatorContext) Lt() antlr.TerminalNode {
	return s.GetToken(FqlParserLt, 0)
}

func (s *EqualityOperatorContext) Eq() antlr.TerminalNode {
	return s.GetToken(FqlParserEq, 0)
}

func (s *EqualityOperatorContext) Gte() antlr.TerminalNode {
	return s.GetToken(FqlParserGte, 0)
}

func (s *EqualityOperatorContext) Lte() antlr.TerminalNode {
	return s.GetToken(FqlParserLte, 0)
}

func (s *EqualityOperatorContext) Neq() antlr.TerminalNode {
	return s.GetToken(FqlParserNeq, 0)
}

func (s *EqualityOperatorContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *EqualityOperatorContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *EqualityOperatorContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(FqlParserListener); ok {
		listenerT.EnterEqualityOperator(s)
	}
}

func (s *EqualityOperatorContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(FqlParserListener); ok {
		listenerT.ExitEqualityOperator(s)
	}
}

func (s *EqualityOperatorContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case FqlParserVisitor:
		return t.VisitEqualityOperator(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *FqlParser) EqualityOperator() (localctx IEqualityOperatorContext) {
	localctx = NewEqualityOperatorContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 126, FqlParserRULE_equalityOperator)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(657)
		_la = p.GetTokenStream().LA(1)

		if !((int64(_la) & ^0x3f) == 0 && ((int64(1)<<_la)&2064384) != 0) {
			p.GetErrorHandler().RecoverInline(p)
		} else {
			p.GetErrorHandler().ReportMatch(p)
			p.Consume()
		}
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IInOperatorContext is an interface to support dynamic dispatch.
type IInOperatorContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	In() antlr.TerminalNode
	Not() antlr.TerminalNode

	// IsInOperatorContext differentiates from other interfaces.
	IsInOperatorContext()
}

type InOperatorContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyInOperatorContext() *InOperatorContext {
	var p = new(InOperatorContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = FqlParserRULE_inOperator
	return p
}

func InitEmptyInOperatorContext(p *InOperatorContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = FqlParserRULE_inOperator
}

func (*InOperatorContext) IsInOperatorContext() {}

func NewInOperatorContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *InOperatorContext {
	var p = new(InOperatorContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = FqlParserRULE_inOperator

	return p
}

func (s *InOperatorContext) GetParser() antlr.Parser { return s.parser }

func (s *InOperatorContext) In() antlr.TerminalNode {
	return s.GetToken(FqlParserIn, 0)
}

func (s *InOperatorContext) Not() antlr.TerminalNode {
	return s.GetToken(FqlParserNot, 0)
}

func (s *InOperatorContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *InOperatorContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *InOperatorContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(FqlParserListener); ok {
		listenerT.EnterInOperator(s)
	}
}

func (s *InOperatorContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(FqlParserListener); ok {
		listenerT.ExitInOperator(s)
	}
}

func (s *InOperatorContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case FqlParserVisitor:
		return t.VisitInOperator(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *FqlParser) InOperator() (localctx IInOperatorContext) {
	localctx = NewInOperatorContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 128, FqlParserRULE_inOperator)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	p.SetState(660)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	if _la == FqlParserNot {
		{
			p.SetState(659)
			p.Match(FqlParserNot)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	}
	{
		p.SetState(662)
		p.Match(FqlParserIn)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// ILikeOperatorContext is an interface to support dynamic dispatch.
type ILikeOperatorContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	Like() antlr.TerminalNode
	Not() antlr.TerminalNode

	// IsLikeOperatorContext differentiates from other interfaces.
	IsLikeOperatorContext()
}

type LikeOperatorContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyLikeOperatorContext() *LikeOperatorContext {
	var p = new(LikeOperatorContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = FqlParserRULE_likeOperator
	return p
}

func InitEmptyLikeOperatorContext(p *LikeOperatorContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = FqlParserRULE_likeOperator
}

func (*LikeOperatorContext) IsLikeOperatorContext() {}

func NewLikeOperatorContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *LikeOperatorContext {
	var p = new(LikeOperatorContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = FqlParserRULE_likeOperator

	return p
}

func (s *LikeOperatorContext) GetParser() antlr.Parser { return s.parser }

func (s *LikeOperatorContext) Like() antlr.TerminalNode {
	return s.GetToken(FqlParserLike, 0)
}

func (s *LikeOperatorContext) Not() antlr.TerminalNode {
	return s.GetToken(FqlParserNot, 0)
}

func (s *LikeOperatorContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *LikeOperatorContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *LikeOperatorContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(FqlParserListener); ok {
		listenerT.EnterLikeOperator(s)
	}
}

func (s *LikeOperatorContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(FqlParserListener); ok {
		listenerT.ExitLikeOperator(s)
	}
}

func (s *LikeOperatorContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case FqlParserVisitor:
		return t.VisitLikeOperator(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *FqlParser) LikeOperator() (localctx ILikeOperatorContext) {
	localctx = NewLikeOperatorContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 130, FqlParserRULE_likeOperator)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	p.SetState(665)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	if _la == FqlParserNot {
		{
			p.SetState(664)
			p.Match(FqlParserNot)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	}
	{
		p.SetState(667)
		p.Match(FqlParserLike)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IUnaryOperatorContext is an interface to support dynamic dispatch.
type IUnaryOperatorContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	Not() antlr.TerminalNode
	Plus() antlr.TerminalNode
	Minus() antlr.TerminalNode

	// IsUnaryOperatorContext differentiates from other interfaces.
	IsUnaryOperatorContext()
}

type UnaryOperatorContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyUnaryOperatorContext() *UnaryOperatorContext {
	var p = new(UnaryOperatorContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = FqlParserRULE_unaryOperator
	return p
}

func InitEmptyUnaryOperatorContext(p *UnaryOperatorContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = FqlParserRULE_unaryOperator
}

func (*UnaryOperatorContext) IsUnaryOperatorContext() {}

func NewUnaryOperatorContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *UnaryOperatorContext {
	var p = new(UnaryOperatorContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = FqlParserRULE_unaryOperator

	return p
}

func (s *UnaryOperatorContext) GetParser() antlr.Parser { return s.parser }

func (s *UnaryOperatorContext) Not() antlr.TerminalNode {
	return s.GetToken(FqlParserNot, 0)
}

func (s *UnaryOperatorContext) Plus() antlr.TerminalNode {
	return s.GetToken(FqlParserPlus, 0)
}

func (s *UnaryOperatorContext) Minus() antlr.TerminalNode {
	return s.GetToken(FqlParserMinus, 0)
}

func (s *UnaryOperatorContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *UnaryOperatorContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *UnaryOperatorContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(FqlParserListener); ok {
		listenerT.EnterUnaryOperator(s)
	}
}

func (s *UnaryOperatorContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(FqlParserListener); ok {
		listenerT.ExitUnaryOperator(s)
	}
}

func (s *UnaryOperatorContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case FqlParserVisitor:
		return t.VisitUnaryOperator(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *FqlParser) UnaryOperator() (localctx IUnaryOperatorContext) {
	localctx = NewUnaryOperatorContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 132, FqlParserRULE_unaryOperator)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(669)
		_la = p.GetTokenStream().LA(1)

		if !((int64(_la) & ^0x3f) == 0 && ((int64(1)<<_la)&288230376202043392) != 0) {
			p.GetErrorHandler().RecoverInline(p)
		} else {
			p.GetErrorHandler().ReportMatch(p)
			p.Consume()
		}
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IRegexpOperatorContext is an interface to support dynamic dispatch.
type IRegexpOperatorContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	RegexMatch() antlr.TerminalNode
	RegexNotMatch() antlr.TerminalNode

	// IsRegexpOperatorContext differentiates from other interfaces.
	IsRegexpOperatorContext()
}

type RegexpOperatorContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyRegexpOperatorContext() *RegexpOperatorContext {
	var p = new(RegexpOperatorContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = FqlParserRULE_regexpOperator
	return p
}

func InitEmptyRegexpOperatorContext(p *RegexpOperatorContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = FqlParserRULE_regexpOperator
}

func (*RegexpOperatorContext) IsRegexpOperatorContext() {}

func NewRegexpOperatorContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *RegexpOperatorContext {
	var p = new(RegexpOperatorContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = FqlParserRULE_regexpOperator

	return p
}

func (s *RegexpOperatorContext) GetParser() antlr.Parser { return s.parser }

func (s *RegexpOperatorContext) RegexMatch() antlr.TerminalNode {
	return s.GetToken(FqlParserRegexMatch, 0)
}

func (s *RegexpOperatorContext) RegexNotMatch() antlr.TerminalNode {
	return s.GetToken(FqlParserRegexNotMatch, 0)
}

func (s *RegexpOperatorContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *RegexpOperatorContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *RegexpOperatorContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(FqlParserListener); ok {
		listenerT.EnterRegexpOperator(s)
	}
}

func (s *RegexpOperatorContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(FqlParserListener); ok {
		listenerT.ExitRegexpOperator(s)
	}
}

func (s *RegexpOperatorContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case FqlParserVisitor:
		return t.VisitRegexpOperator(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *FqlParser) RegexpOperator() (localctx IRegexpOperatorContext) {
	localctx = NewRegexpOperatorContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 134, FqlParserRULE_regexpOperator)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(671)
		_la = p.GetTokenStream().LA(1)

		if !(_la == FqlParserRegexNotMatch || _la == FqlParserRegexMatch) {
			p.GetErrorHandler().RecoverInline(p)
		} else {
			p.GetErrorHandler().ReportMatch(p)
			p.Consume()
		}
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// ILogicalAndOperatorContext is an interface to support dynamic dispatch.
type ILogicalAndOperatorContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	And() antlr.TerminalNode

	// IsLogicalAndOperatorContext differentiates from other interfaces.
	IsLogicalAndOperatorContext()
}

type LogicalAndOperatorContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyLogicalAndOperatorContext() *LogicalAndOperatorContext {
	var p = new(LogicalAndOperatorContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = FqlParserRULE_logicalAndOperator
	return p
}

func InitEmptyLogicalAndOperatorContext(p *LogicalAndOperatorContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = FqlParserRULE_logicalAndOperator
}

func (*LogicalAndOperatorContext) IsLogicalAndOperatorContext() {}

func NewLogicalAndOperatorContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *LogicalAndOperatorContext {
	var p = new(LogicalAndOperatorContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = FqlParserRULE_logicalAndOperator

	return p
}

func (s *LogicalAndOperatorContext) GetParser() antlr.Parser { return s.parser }

func (s *LogicalAndOperatorContext) And() antlr.TerminalNode {
	return s.GetToken(FqlParserAnd, 0)
}

func (s *LogicalAndOperatorContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *LogicalAndOperatorContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *LogicalAndOperatorContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(FqlParserListener); ok {
		listenerT.EnterLogicalAndOperator(s)
	}
}

func (s *LogicalAndOperatorContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(FqlParserListener); ok {
		listenerT.ExitLogicalAndOperator(s)
	}
}

func (s *LogicalAndOperatorContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case FqlParserVisitor:
		return t.VisitLogicalAndOperator(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *FqlParser) LogicalAndOperator() (localctx ILogicalAndOperatorContext) {
	localctx = NewLogicalAndOperatorContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 136, FqlParserRULE_logicalAndOperator)
	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(673)
		p.Match(FqlParserAnd)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// ILogicalOrOperatorContext is an interface to support dynamic dispatch.
type ILogicalOrOperatorContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	Or() antlr.TerminalNode

	// IsLogicalOrOperatorContext differentiates from other interfaces.
	IsLogicalOrOperatorContext()
}

type LogicalOrOperatorContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyLogicalOrOperatorContext() *LogicalOrOperatorContext {
	var p = new(LogicalOrOperatorContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = FqlParserRULE_logicalOrOperator
	return p
}

func InitEmptyLogicalOrOperatorContext(p *LogicalOrOperatorContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = FqlParserRULE_logicalOrOperator
}

func (*LogicalOrOperatorContext) IsLogicalOrOperatorContext() {}

func NewLogicalOrOperatorContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *LogicalOrOperatorContext {
	var p = new(LogicalOrOperatorContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = FqlParserRULE_logicalOrOperator

	return p
}

func (s *LogicalOrOperatorContext) GetParser() antlr.Parser { return s.parser }

func (s *LogicalOrOperatorContext) Or() antlr.TerminalNode {
	return s.GetToken(FqlParserOr, 0)
}

func (s *LogicalOrOperatorContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *LogicalOrOperatorContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *LogicalOrOperatorContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(FqlParserListener); ok {
		listenerT.EnterLogicalOrOperator(s)
	}
}

func (s *LogicalOrOperatorContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(FqlParserListener); ok {
		listenerT.ExitLogicalOrOperator(s)
	}
}

func (s *LogicalOrOperatorContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case FqlParserVisitor:
		return t.VisitLogicalOrOperator(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *FqlParser) LogicalOrOperator() (localctx ILogicalOrOperatorContext) {
	localctx = NewLogicalOrOperatorContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 138, FqlParserRULE_logicalOrOperator)
	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(675)
		p.Match(FqlParserOr)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IMultiplicativeOperatorContext is an interface to support dynamic dispatch.
type IMultiplicativeOperatorContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	Multi() antlr.TerminalNode
	Div() antlr.TerminalNode
	Mod() antlr.TerminalNode

	// IsMultiplicativeOperatorContext differentiates from other interfaces.
	IsMultiplicativeOperatorContext()
}

type MultiplicativeOperatorContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyMultiplicativeOperatorContext() *MultiplicativeOperatorContext {
	var p = new(MultiplicativeOperatorContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = FqlParserRULE_multiplicativeOperator
	return p
}

func InitEmptyMultiplicativeOperatorContext(p *MultiplicativeOperatorContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = FqlParserRULE_multiplicativeOperator
}

func (*MultiplicativeOperatorContext) IsMultiplicativeOperatorContext() {}

func NewMultiplicativeOperatorContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *MultiplicativeOperatorContext {
	var p = new(MultiplicativeOperatorContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = FqlParserRULE_multiplicativeOperator

	return p
}

func (s *MultiplicativeOperatorContext) GetParser() antlr.Parser { return s.parser }

func (s *MultiplicativeOperatorContext) Multi() antlr.TerminalNode {
	return s.GetToken(FqlParserMulti, 0)
}

func (s *MultiplicativeOperatorContext) Div() antlr.TerminalNode {
	return s.GetToken(FqlParserDiv, 0)
}

func (s *MultiplicativeOperatorContext) Mod() antlr.TerminalNode {
	return s.GetToken(FqlParserMod, 0)
}

func (s *MultiplicativeOperatorContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *MultiplicativeOperatorContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *MultiplicativeOperatorContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(FqlParserListener); ok {
		listenerT.EnterMultiplicativeOperator(s)
	}
}

func (s *MultiplicativeOperatorContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(FqlParserListener); ok {
		listenerT.ExitMultiplicativeOperator(s)
	}
}

func (s *MultiplicativeOperatorContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case FqlParserVisitor:
		return t.VisitMultiplicativeOperator(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *FqlParser) MultiplicativeOperator() (localctx IMultiplicativeOperatorContext) {
	localctx = NewMultiplicativeOperatorContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 140, FqlParserRULE_multiplicativeOperator)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(677)
		_la = p.GetTokenStream().LA(1)

		if !((int64(_la) & ^0x3f) == 0 && ((int64(1)<<_la)&14680064) != 0) {
			p.GetErrorHandler().RecoverInline(p)
		} else {
			p.GetErrorHandler().ReportMatch(p)
			p.Consume()
		}
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IAdditiveOperatorContext is an interface to support dynamic dispatch.
type IAdditiveOperatorContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	Plus() antlr.TerminalNode
	Minus() antlr.TerminalNode

	// IsAdditiveOperatorContext differentiates from other interfaces.
	IsAdditiveOperatorContext()
}

type AdditiveOperatorContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyAdditiveOperatorContext() *AdditiveOperatorContext {
	var p = new(AdditiveOperatorContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = FqlParserRULE_additiveOperator
	return p
}

func InitEmptyAdditiveOperatorContext(p *AdditiveOperatorContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = FqlParserRULE_additiveOperator
}

func (*AdditiveOperatorContext) IsAdditiveOperatorContext() {}

func NewAdditiveOperatorContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *AdditiveOperatorContext {
	var p = new(AdditiveOperatorContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = FqlParserRULE_additiveOperator

	return p
}

func (s *AdditiveOperatorContext) GetParser() antlr.Parser { return s.parser }

func (s *AdditiveOperatorContext) Plus() antlr.TerminalNode {
	return s.GetToken(FqlParserPlus, 0)
}

func (s *AdditiveOperatorContext) Minus() antlr.TerminalNode {
	return s.GetToken(FqlParserMinus, 0)
}

func (s *AdditiveOperatorContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *AdditiveOperatorContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *AdditiveOperatorContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(FqlParserListener); ok {
		listenerT.EnterAdditiveOperator(s)
	}
}

func (s *AdditiveOperatorContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(FqlParserListener); ok {
		listenerT.ExitAdditiveOperator(s)
	}
}

func (s *AdditiveOperatorContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case FqlParserVisitor:
		return t.VisitAdditiveOperator(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *FqlParser) AdditiveOperator() (localctx IAdditiveOperatorContext) {
	localctx = NewAdditiveOperatorContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 142, FqlParserRULE_additiveOperator)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(679)
		_la = p.GetTokenStream().LA(1)

		if !(_la == FqlParserPlus || _la == FqlParserMinus) {
			p.GetErrorHandler().RecoverInline(p)
		} else {
			p.GetErrorHandler().ReportMatch(p)
			p.Consume()
		}
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IErrorOperatorContext is an interface to support dynamic dispatch.
type IErrorOperatorContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	QuestionMark() antlr.TerminalNode

	// IsErrorOperatorContext differentiates from other interfaces.
	IsErrorOperatorContext()
}

type ErrorOperatorContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyErrorOperatorContext() *ErrorOperatorContext {
	var p = new(ErrorOperatorContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = FqlParserRULE_errorOperator
	return p
}

func InitEmptyErrorOperatorContext(p *ErrorOperatorContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = FqlParserRULE_errorOperator
}

func (*ErrorOperatorContext) IsErrorOperatorContext() {}

func NewErrorOperatorContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ErrorOperatorContext {
	var p = new(ErrorOperatorContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = FqlParserRULE_errorOperator

	return p
}

func (s *ErrorOperatorContext) GetParser() antlr.Parser { return s.parser }

func (s *ErrorOperatorContext) QuestionMark() antlr.TerminalNode {
	return s.GetToken(FqlParserQuestionMark, 0)
}

func (s *ErrorOperatorContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ErrorOperatorContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *ErrorOperatorContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(FqlParserListener); ok {
		listenerT.EnterErrorOperator(s)
	}
}

func (s *ErrorOperatorContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(FqlParserListener); ok {
		listenerT.ExitErrorOperator(s)
	}
}

func (s *ErrorOperatorContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case FqlParserVisitor:
		return t.VisitErrorOperator(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *FqlParser) ErrorOperator() (localctx IErrorOperatorContext) {
	localctx = NewErrorOperatorContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 144, FqlParserRULE_errorOperator)
	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(681)
		p.Match(FqlParserQuestionMark)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

func (p *FqlParser) Sempred(localctx antlr.RuleContext, ruleIndex, predIndex int) bool {
	switch ruleIndex {
	case 59:
		var t *ExpressionContext = nil
		if localctx != nil {
			t = localctx.(*ExpressionContext)
		}
		return p.Expression_Sempred(t, predIndex)

	case 60:
		var t *PredicateContext = nil
		if localctx != nil {
			t = localctx.(*PredicateContext)
		}
		return p.Predicate_Sempred(t, predIndex)

	case 61:
		var t *ExpressionAtomContext = nil
		if localctx != nil {
			t = localctx.(*ExpressionAtomContext)
		}
		return p.ExpressionAtom_Sempred(t, predIndex)

	default:
		panic("No predicate with index: " + fmt.Sprint(ruleIndex))
	}
}

func (p *FqlParser) Expression_Sempred(localctx antlr.RuleContext, predIndex int) bool {
	switch predIndex {
	case 0:
		return p.Precpred(p.GetParserRuleContext(), 4)

	case 1:
		return p.Precpred(p.GetParserRuleContext(), 3)

	case 2:
		return p.Precpred(p.GetParserRuleContext(), 2)

	default:
		panic("No predicate with index: " + fmt.Sprint(predIndex))
	}
}

func (p *FqlParser) Predicate_Sempred(localctx antlr.RuleContext, predIndex int) bool {
	switch predIndex {
	case 3:
		return p.Precpred(p.GetParserRuleContext(), 5)

	case 4:
		return p.Precpred(p.GetParserRuleContext(), 4)

	case 5:
		return p.Precpred(p.GetParserRuleContext(), 3)

	case 6:
		return p.Precpred(p.GetParserRuleContext(), 2)

	default:
		panic("No predicate with index: " + fmt.Sprint(predIndex))
	}
}

func (p *FqlParser) ExpressionAtom_Sempred(localctx antlr.RuleContext, predIndex int) bool {
	switch predIndex {
	case 7:
		return p.Precpred(p.GetParserRuleContext(), 10)

	case 8:
		return p.Precpred(p.GetParserRuleContext(), 9)

	case 9:
		return p.Precpred(p.GetParserRuleContext(), 8)

	default:
		panic("No predicate with index: " + fmt.Sprint(predIndex))
	}
}
