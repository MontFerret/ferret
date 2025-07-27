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
		"'/'", "'%'", "'+'", "'-'", "'--'", "'++'", "", "", "", "'='", "'?'",
		"'!~'", "'=~'", "'FOR'", "'RETURN'", "'WAITFOR'", "'OPTIONS'", "'TIMEOUT'",
		"'DISTINCT'", "'FILTER'", "'CURRENT'", "'SORT'", "'LIMIT'", "'LET'",
		"'COLLECT'", "", "'NONE'", "'NULL'", "", "'USE'", "'INTO'", "'KEEP'",
		"'WITH'", "'ALL'", "'ANY'", "'AGGREGATE'", "'EVENT'", "'LIKE'", "",
		"'IN'", "'DO'", "'WHILE'", "'@'",
	}
	staticData.SymbolicNames = []string{
		"", "MultiLineComment", "SingleLineComment", "WhiteSpaces", "LineTerminator",
		"Colon", "SemiColon", "Dot", "Comma", "OpenBracket", "CloseBracket",
		"OpenParen", "CloseParen", "OpenBrace", "CloseBrace", "Gt", "Lt", "Eq",
		"Gte", "Lte", "Neq", "Multi", "Div", "Mod", "Plus", "Minus", "MinusMinus",
		"PlusPlus", "And", "Or", "Range", "Assign", "QuestionMark", "RegexNotMatch",
		"RegexMatch", "For", "Return", "Waitfor", "Options", "Timeout", "Distinct",
		"Filter", "Current", "Sort", "Limit", "Let", "Collect", "SortDirection",
		"None", "Null", "BooleanLiteral", "Use", "Into", "Keep", "With", "All",
		"Any", "Aggregate", "Event", "Like", "Not", "In", "Do", "While", "Param",
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
		4, 1, 71, 662, 2, 0, 7, 0, 2, 1, 7, 1, 2, 2, 7, 2, 2, 3, 7, 3, 2, 4, 7,
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
		1, 9, 1, 9, 1, 9, 1, 9, 1, 9, 3, 9, 215, 8, 9, 1, 9, 1, 9, 1, 9, 5, 9,
		220, 8, 9, 10, 9, 12, 9, 223, 9, 9, 1, 9, 1, 9, 3, 9, 227, 8, 9, 1, 10,
		1, 10, 1, 10, 1, 10, 1, 10, 1, 10, 1, 10, 3, 10, 236, 8, 10, 1, 11, 1,
		11, 1, 11, 1, 11, 3, 11, 242, 8, 11, 1, 12, 1, 12, 3, 12, 246, 8, 12, 1,
		13, 1, 13, 3, 13, 250, 8, 13, 1, 14, 1, 14, 3, 14, 254, 8, 14, 1, 15, 1,
		15, 1, 15, 1, 16, 1, 16, 1, 16, 1, 16, 3, 16, 263, 8, 16, 1, 17, 1, 17,
		1, 17, 1, 17, 1, 17, 3, 17, 270, 8, 17, 1, 18, 1, 18, 1, 18, 1, 18, 5,
		18, 276, 8, 18, 10, 18, 12, 18, 279, 9, 18, 1, 19, 1, 19, 3, 19, 283, 8,
		19, 1, 20, 1, 20, 1, 20, 1, 20, 1, 20, 1, 20, 1, 20, 1, 20, 1, 20, 1, 20,
		1, 20, 1, 20, 1, 20, 1, 20, 1, 20, 1, 20, 1, 20, 1, 20, 1, 20, 1, 20, 1,
		20, 1, 20, 1, 20, 3, 20, 308, 8, 20, 1, 21, 1, 21, 1, 21, 1, 21, 1, 22,
		1, 22, 1, 22, 5, 22, 317, 8, 22, 10, 22, 12, 22, 320, 9, 22, 1, 23, 1,
		23, 1, 23, 1, 23, 1, 24, 1, 24, 1, 24, 1, 24, 5, 24, 330, 8, 24, 10, 24,
		12, 24, 333, 9, 24, 1, 25, 1, 25, 1, 25, 1, 25, 1, 25, 3, 25, 340, 8, 25,
		3, 25, 342, 8, 25, 1, 26, 1, 26, 1, 26, 1, 26, 5, 26, 348, 8, 26, 10, 26,
		12, 26, 351, 9, 26, 1, 27, 1, 27, 1, 27, 1, 27, 1, 27, 1, 28, 1, 28, 1,
		28, 1, 28, 1, 28, 1, 28, 3, 28, 364, 8, 28, 1, 28, 3, 28, 367, 8, 28, 1,
		28, 3, 28, 370, 8, 28, 1, 29, 1, 29, 1, 29, 1, 29, 1, 29, 3, 29, 377, 8,
		29, 1, 30, 1, 30, 1, 30, 3, 30, 382, 8, 30, 1, 31, 1, 31, 1, 31, 1, 32,
		1, 32, 1, 32, 1, 32, 1, 32, 1, 32, 3, 32, 393, 8, 32, 1, 33, 1, 33, 1,
		33, 1, 33, 3, 33, 399, 8, 33, 1, 34, 1, 34, 3, 34, 403, 8, 34, 1, 35, 1,
		35, 1, 35, 1, 35, 1, 35, 1, 35, 1, 35, 3, 35, 412, 8, 35, 1, 36, 1, 36,
		3, 36, 416, 8, 36, 1, 36, 1, 36, 1, 37, 1, 37, 1, 37, 1, 37, 5, 37, 424,
		8, 37, 10, 37, 12, 37, 427, 9, 37, 1, 37, 3, 37, 430, 8, 37, 3, 37, 432,
		8, 37, 1, 37, 1, 37, 1, 38, 1, 38, 1, 39, 1, 39, 1, 40, 1, 40, 1, 41, 1,
		41, 1, 42, 1, 42, 1, 43, 1, 43, 1, 43, 1, 43, 1, 43, 1, 43, 1, 43, 1, 43,
		1, 43, 3, 43, 455, 8, 43, 1, 44, 1, 44, 1, 44, 1, 44, 1, 45, 1, 45, 1,
		45, 1, 45, 1, 45, 3, 45, 466, 8, 45, 1, 46, 1, 46, 1, 46, 1, 47, 5, 47,
		472, 8, 47, 10, 47, 12, 47, 475, 9, 47, 1, 48, 1, 48, 4, 48, 479, 8, 48,
		11, 48, 12, 48, 480, 1, 49, 1, 49, 1, 49, 1, 49, 1, 49, 3, 49, 488, 8,
		49, 1, 50, 1, 50, 3, 50, 492, 8, 50, 1, 51, 1, 51, 1, 51, 1, 51, 3, 51,
		498, 8, 51, 1, 51, 1, 51, 1, 52, 1, 52, 1, 52, 3, 52, 505, 8, 52, 1, 53,
		1, 53, 1, 53, 5, 53, 510, 8, 53, 10, 53, 12, 53, 513, 9, 53, 1, 53, 3,
		53, 516, 8, 53, 1, 54, 3, 54, 519, 8, 54, 1, 54, 1, 54, 1, 54, 1, 54, 1,
		54, 3, 54, 526, 8, 54, 1, 54, 3, 54, 529, 8, 54, 1, 55, 1, 55, 1, 56, 1,
		56, 1, 57, 1, 57, 1, 57, 1, 57, 1, 58, 1, 58, 1, 58, 3, 58, 542, 8, 58,
		1, 59, 1, 59, 1, 59, 1, 59, 1, 59, 3, 59, 549, 8, 59, 1, 59, 1, 59, 1,
		59, 1, 59, 1, 59, 1, 59, 1, 59, 1, 59, 1, 59, 1, 59, 1, 59, 3, 59, 562,
		8, 59, 1, 59, 1, 59, 5, 59, 566, 8, 59, 10, 59, 12, 59, 569, 9, 59, 1,
		60, 1, 60, 1, 60, 1, 60, 1, 60, 1, 60, 1, 60, 1, 60, 1, 60, 1, 60, 1, 60,
		1, 60, 1, 60, 1, 60, 1, 60, 1, 60, 1, 60, 1, 60, 1, 60, 5, 60, 590, 8,
		60, 10, 60, 12, 60, 593, 9, 60, 1, 61, 1, 61, 1, 61, 1, 61, 1, 61, 1, 61,
		1, 61, 1, 61, 1, 61, 1, 61, 1, 61, 3, 61, 606, 8, 61, 1, 61, 1, 61, 3,
		61, 610, 8, 61, 3, 61, 612, 8, 61, 1, 61, 1, 61, 1, 61, 1, 61, 1, 61, 1,
		61, 1, 61, 1, 61, 1, 61, 1, 61, 1, 61, 1, 61, 5, 61, 626, 8, 61, 10, 61,
		12, 61, 629, 9, 61, 1, 62, 1, 62, 1, 62, 3, 62, 634, 8, 62, 1, 63, 1, 63,
		1, 64, 3, 64, 639, 8, 64, 1, 64, 1, 64, 1, 65, 3, 65, 644, 8, 65, 1, 65,
		1, 65, 1, 66, 1, 66, 1, 67, 1, 67, 1, 68, 1, 68, 1, 69, 1, 69, 1, 70, 1,
		70, 1, 71, 1, 71, 1, 72, 1, 72, 1, 72, 0, 3, 118, 120, 122, 73, 0, 2, 4,
		6, 8, 10, 12, 14, 16, 18, 20, 22, 24, 26, 28, 30, 32, 34, 36, 38, 40, 42,
		44, 46, 48, 50, 52, 54, 56, 58, 60, 62, 64, 66, 68, 70, 72, 74, 76, 78,
		80, 82, 84, 86, 88, 90, 92, 94, 96, 98, 100, 102, 104, 106, 108, 110, 112,
		114, 116, 118, 120, 122, 124, 126, 128, 130, 132, 134, 136, 138, 140, 142,
		144, 0, 10, 1, 0, 65, 66, 1, 0, 48, 49, 4, 0, 28, 29, 38, 44, 46, 47, 52,
		58, 4, 0, 35, 37, 45, 45, 48, 51, 59, 63, 2, 0, 48, 48, 55, 56, 1, 0, 15,
		20, 2, 0, 24, 25, 60, 60, 1, 0, 33, 34, 1, 0, 21, 23, 1, 0, 24, 25, 702,
		0, 149, 1, 0, 0, 0, 2, 154, 1, 0, 0, 0, 4, 156, 1, 0, 0, 0, 6, 158, 1,
		0, 0, 0, 8, 164, 1, 0, 0, 0, 10, 172, 1, 0, 0, 0, 12, 176, 1, 0, 0, 0,
		14, 187, 1, 0, 0, 0, 16, 189, 1, 0, 0, 0, 18, 226, 1, 0, 0, 0, 20, 235,
		1, 0, 0, 0, 22, 241, 1, 0, 0, 0, 24, 245, 1, 0, 0, 0, 26, 249, 1, 0, 0,
		0, 28, 253, 1, 0, 0, 0, 30, 255, 1, 0, 0, 0, 32, 258, 1, 0, 0, 0, 34, 269,
		1, 0, 0, 0, 36, 271, 1, 0, 0, 0, 38, 280, 1, 0, 0, 0, 40, 307, 1, 0, 0,
		0, 42, 309, 1, 0, 0, 0, 44, 313, 1, 0, 0, 0, 46, 321, 1, 0, 0, 0, 48, 325,
		1, 0, 0, 0, 50, 341, 1, 0, 0, 0, 52, 343, 1, 0, 0, 0, 54, 352, 1, 0, 0,
		0, 56, 357, 1, 0, 0, 0, 58, 376, 1, 0, 0, 0, 60, 381, 1, 0, 0, 0, 62, 383,
		1, 0, 0, 0, 64, 386, 1, 0, 0, 0, 66, 398, 1, 0, 0, 0, 68, 402, 1, 0, 0,
		0, 70, 411, 1, 0, 0, 0, 72, 413, 1, 0, 0, 0, 74, 419, 1, 0, 0, 0, 76, 435,
		1, 0, 0, 0, 78, 437, 1, 0, 0, 0, 80, 439, 1, 0, 0, 0, 82, 441, 1, 0, 0,
		0, 84, 443, 1, 0, 0, 0, 86, 454, 1, 0, 0, 0, 88, 456, 1, 0, 0, 0, 90, 465,
		1, 0, 0, 0, 92, 467, 1, 0, 0, 0, 94, 473, 1, 0, 0, 0, 96, 476, 1, 0, 0,
		0, 98, 487, 1, 0, 0, 0, 100, 489, 1, 0, 0, 0, 102, 493, 1, 0, 0, 0, 104,
		504, 1, 0, 0, 0, 106, 506, 1, 0, 0, 0, 108, 528, 1, 0, 0, 0, 110, 530,
		1, 0, 0, 0, 112, 532, 1, 0, 0, 0, 114, 534, 1, 0, 0, 0, 116, 541, 1, 0,
		0, 0, 118, 548, 1, 0, 0, 0, 120, 570, 1, 0, 0, 0, 122, 611, 1, 0, 0, 0,
		124, 630, 1, 0, 0, 0, 126, 635, 1, 0, 0, 0, 128, 638, 1, 0, 0, 0, 130,
		643, 1, 0, 0, 0, 132, 647, 1, 0, 0, 0, 134, 649, 1, 0, 0, 0, 136, 651,
		1, 0, 0, 0, 138, 653, 1, 0, 0, 0, 140, 655, 1, 0, 0, 0, 142, 657, 1, 0,
		0, 0, 144, 659, 1, 0, 0, 0, 146, 148, 3, 2, 1, 0, 147, 146, 1, 0, 0, 0,
		148, 151, 1, 0, 0, 0, 149, 147, 1, 0, 0, 0, 149, 150, 1, 0, 0, 0, 150,
		152, 1, 0, 0, 0, 151, 149, 1, 0, 0, 0, 152, 153, 3, 8, 4, 0, 153, 1, 1,
		0, 0, 0, 154, 155, 3, 4, 2, 0, 155, 3, 1, 0, 0, 0, 156, 157, 3, 6, 3, 0,
		157, 5, 1, 0, 0, 0, 158, 159, 5, 51, 0, 0, 159, 160, 3, 92, 46, 0, 160,
		7, 1, 0, 0, 0, 161, 163, 3, 10, 5, 0, 162, 161, 1, 0, 0, 0, 163, 166, 1,
		0, 0, 0, 164, 162, 1, 0, 0, 0, 164, 165, 1, 0, 0, 0, 165, 167, 1, 0, 0,
		0, 166, 164, 1, 0, 0, 0, 167, 168, 3, 12, 6, 0, 168, 9, 1, 0, 0, 0, 169,
		173, 3, 14, 7, 0, 170, 173, 3, 100, 50, 0, 171, 173, 3, 56, 28, 0, 172,
		169, 1, 0, 0, 0, 172, 170, 1, 0, 0, 0, 172, 171, 1, 0, 0, 0, 173, 11, 1,
		0, 0, 0, 174, 177, 3, 16, 8, 0, 175, 177, 3, 18, 9, 0, 176, 174, 1, 0,
		0, 0, 176, 175, 1, 0, 0, 0, 177, 13, 1, 0, 0, 0, 178, 179, 5, 45, 0, 0,
		179, 180, 7, 0, 0, 0, 180, 181, 5, 31, 0, 0, 181, 188, 3, 118, 59, 0, 182,
		183, 5, 45, 0, 0, 183, 184, 3, 110, 55, 0, 184, 185, 5, 31, 0, 0, 185,
		186, 3, 118, 59, 0, 186, 188, 1, 0, 0, 0, 187, 178, 1, 0, 0, 0, 187, 182,
		1, 0, 0, 0, 188, 15, 1, 0, 0, 0, 189, 191, 5, 36, 0, 0, 190, 192, 5, 40,
		0, 0, 191, 190, 1, 0, 0, 0, 191, 192, 1, 0, 0, 0, 192, 193, 1, 0, 0, 0,
		193, 194, 3, 118, 59, 0, 194, 17, 1, 0, 0, 0, 195, 196, 5, 35, 0, 0, 196,
		199, 7, 0, 0, 0, 197, 198, 5, 8, 0, 0, 198, 200, 5, 65, 0, 0, 199, 197,
		1, 0, 0, 0, 199, 200, 1, 0, 0, 0, 200, 201, 1, 0, 0, 0, 201, 202, 5, 61,
		0, 0, 202, 206, 3, 20, 10, 0, 203, 205, 3, 26, 13, 0, 204, 203, 1, 0, 0,
		0, 205, 208, 1, 0, 0, 0, 206, 204, 1, 0, 0, 0, 206, 207, 1, 0, 0, 0, 207,
		209, 1, 0, 0, 0, 208, 206, 1, 0, 0, 0, 209, 210, 3, 28, 14, 0, 210, 227,
		1, 0, 0, 0, 211, 212, 5, 35, 0, 0, 212, 214, 7, 0, 0, 0, 213, 215, 5, 62,
		0, 0, 214, 213, 1, 0, 0, 0, 214, 215, 1, 0, 0, 0, 215, 216, 1, 0, 0, 0,
		216, 217, 5, 63, 0, 0, 217, 221, 3, 118, 59, 0, 218, 220, 3, 26, 13, 0,
		219, 218, 1, 0, 0, 0, 220, 223, 1, 0, 0, 0, 221, 219, 1, 0, 0, 0, 221,
		222, 1, 0, 0, 0, 222, 224, 1, 0, 0, 0, 223, 221, 1, 0, 0, 0, 224, 225,
		3, 28, 14, 0, 225, 227, 1, 0, 0, 0, 226, 195, 1, 0, 0, 0, 226, 211, 1,
		0, 0, 0, 227, 19, 1, 0, 0, 0, 228, 236, 3, 100, 50, 0, 229, 236, 3, 72,
		36, 0, 230, 236, 3, 74, 37, 0, 231, 236, 3, 68, 34, 0, 232, 236, 3, 96,
		48, 0, 233, 236, 3, 114, 57, 0, 234, 236, 3, 66, 33, 0, 235, 228, 1, 0,
		0, 0, 235, 229, 1, 0, 0, 0, 235, 230, 1, 0, 0, 0, 235, 231, 1, 0, 0, 0,
		235, 232, 1, 0, 0, 0, 235, 233, 1, 0, 0, 0, 235, 234, 1, 0, 0, 0, 236,
		21, 1, 0, 0, 0, 237, 242, 3, 32, 16, 0, 238, 242, 3, 36, 18, 0, 239, 242,
		3, 30, 15, 0, 240, 242, 3, 40, 20, 0, 241, 237, 1, 0, 0, 0, 241, 238, 1,
		0, 0, 0, 241, 239, 1, 0, 0, 0, 241, 240, 1, 0, 0, 0, 242, 23, 1, 0, 0,
		0, 243, 246, 3, 14, 7, 0, 244, 246, 3, 100, 50, 0, 245, 243, 1, 0, 0, 0,
		245, 244, 1, 0, 0, 0, 246, 25, 1, 0, 0, 0, 247, 250, 3, 24, 12, 0, 248,
		250, 3, 22, 11, 0, 249, 247, 1, 0, 0, 0, 249, 248, 1, 0, 0, 0, 250, 27,
		1, 0, 0, 0, 251, 254, 3, 16, 8, 0, 252, 254, 3, 18, 9, 0, 253, 251, 1,
		0, 0, 0, 253, 252, 1, 0, 0, 0, 254, 29, 1, 0, 0, 0, 255, 256, 5, 41, 0,
		0, 256, 257, 3, 118, 59, 0, 257, 31, 1, 0, 0, 0, 258, 259, 5, 44, 0, 0,
		259, 262, 3, 34, 17, 0, 260, 261, 5, 8, 0, 0, 261, 263, 3, 34, 17, 0, 262,
		260, 1, 0, 0, 0, 262, 263, 1, 0, 0, 0, 263, 33, 1, 0, 0, 0, 264, 270, 3,
		82, 41, 0, 265, 270, 3, 66, 33, 0, 266, 270, 3, 68, 34, 0, 267, 270, 3,
		100, 50, 0, 268, 270, 3, 96, 48, 0, 269, 264, 1, 0, 0, 0, 269, 265, 1,
		0, 0, 0, 269, 266, 1, 0, 0, 0, 269, 267, 1, 0, 0, 0, 269, 268, 1, 0, 0,
		0, 270, 35, 1, 0, 0, 0, 271, 272, 5, 43, 0, 0, 272, 277, 3, 38, 19, 0,
		273, 274, 5, 8, 0, 0, 274, 276, 3, 38, 19, 0, 275, 273, 1, 0, 0, 0, 276,
		279, 1, 0, 0, 0, 277, 275, 1, 0, 0, 0, 277, 278, 1, 0, 0, 0, 278, 37, 1,
		0, 0, 0, 279, 277, 1, 0, 0, 0, 280, 282, 3, 118, 59, 0, 281, 283, 5, 47,
		0, 0, 282, 281, 1, 0, 0, 0, 282, 283, 1, 0, 0, 0, 283, 39, 1, 0, 0, 0,
		284, 285, 5, 46, 0, 0, 285, 286, 3, 44, 22, 0, 286, 287, 3, 48, 24, 0,
		287, 288, 3, 50, 25, 0, 288, 308, 1, 0, 0, 0, 289, 290, 5, 46, 0, 0, 290,
		291, 3, 44, 22, 0, 291, 292, 3, 48, 24, 0, 292, 308, 1, 0, 0, 0, 293, 294,
		5, 46, 0, 0, 294, 295, 3, 44, 22, 0, 295, 296, 3, 50, 25, 0, 296, 308,
		1, 0, 0, 0, 297, 298, 5, 46, 0, 0, 298, 308, 3, 44, 22, 0, 299, 300, 5,
		46, 0, 0, 300, 301, 3, 48, 24, 0, 301, 302, 3, 50, 25, 0, 302, 308, 1,
		0, 0, 0, 303, 304, 5, 46, 0, 0, 304, 308, 3, 48, 24, 0, 305, 306, 5, 46,
		0, 0, 306, 308, 3, 54, 27, 0, 307, 284, 1, 0, 0, 0, 307, 289, 1, 0, 0,
		0, 307, 293, 1, 0, 0, 0, 307, 297, 1, 0, 0, 0, 307, 299, 1, 0, 0, 0, 307,
		303, 1, 0, 0, 0, 307, 305, 1, 0, 0, 0, 308, 41, 1, 0, 0, 0, 309, 310, 5,
		65, 0, 0, 310, 311, 5, 31, 0, 0, 311, 312, 3, 118, 59, 0, 312, 43, 1, 0,
		0, 0, 313, 318, 3, 42, 21, 0, 314, 315, 5, 8, 0, 0, 315, 317, 3, 42, 21,
		0, 316, 314, 1, 0, 0, 0, 317, 320, 1, 0, 0, 0, 318, 316, 1, 0, 0, 0, 318,
		319, 1, 0, 0, 0, 319, 45, 1, 0, 0, 0, 320, 318, 1, 0, 0, 0, 321, 322, 5,
		65, 0, 0, 322, 323, 5, 31, 0, 0, 323, 324, 3, 100, 50, 0, 324, 47, 1, 0,
		0, 0, 325, 326, 5, 57, 0, 0, 326, 331, 3, 46, 23, 0, 327, 328, 5, 8, 0,
		0, 328, 330, 3, 46, 23, 0, 329, 327, 1, 0, 0, 0, 330, 333, 1, 0, 0, 0,
		331, 329, 1, 0, 0, 0, 331, 332, 1, 0, 0, 0, 332, 49, 1, 0, 0, 0, 333, 331,
		1, 0, 0, 0, 334, 335, 5, 52, 0, 0, 335, 342, 3, 42, 21, 0, 336, 337, 5,
		52, 0, 0, 337, 339, 5, 65, 0, 0, 338, 340, 3, 52, 26, 0, 339, 338, 1, 0,
		0, 0, 339, 340, 1, 0, 0, 0, 340, 342, 1, 0, 0, 0, 341, 334, 1, 0, 0, 0,
		341, 336, 1, 0, 0, 0, 342, 51, 1, 0, 0, 0, 343, 344, 5, 53, 0, 0, 344,
		349, 5, 65, 0, 0, 345, 346, 5, 8, 0, 0, 346, 348, 5, 65, 0, 0, 347, 345,
		1, 0, 0, 0, 348, 351, 1, 0, 0, 0, 349, 347, 1, 0, 0, 0, 349, 350, 1, 0,
		0, 0, 350, 53, 1, 0, 0, 0, 351, 349, 1, 0, 0, 0, 352, 353, 5, 54, 0, 0,
		353, 354, 5, 65, 0, 0, 354, 355, 5, 52, 0, 0, 355, 356, 5, 65, 0, 0, 356,
		55, 1, 0, 0, 0, 357, 358, 5, 37, 0, 0, 358, 359, 5, 58, 0, 0, 359, 360,
		3, 58, 29, 0, 360, 361, 5, 61, 0, 0, 361, 363, 3, 60, 30, 0, 362, 364,
		3, 62, 31, 0, 363, 362, 1, 0, 0, 0, 363, 364, 1, 0, 0, 0, 364, 366, 1,
		0, 0, 0, 365, 367, 3, 30, 15, 0, 366, 365, 1, 0, 0, 0, 366, 367, 1, 0,
		0, 0, 367, 369, 1, 0, 0, 0, 368, 370, 3, 64, 32, 0, 369, 368, 1, 0, 0,
		0, 369, 370, 1, 0, 0, 0, 370, 57, 1, 0, 0, 0, 371, 377, 3, 78, 39, 0, 372,
		377, 3, 68, 34, 0, 373, 377, 3, 66, 33, 0, 374, 377, 3, 100, 50, 0, 375,
		377, 3, 96, 48, 0, 376, 371, 1, 0, 0, 0, 376, 372, 1, 0, 0, 0, 376, 373,
		1, 0, 0, 0, 376, 374, 1, 0, 0, 0, 376, 375, 1, 0, 0, 0, 377, 59, 1, 0,
		0, 0, 378, 382, 3, 100, 50, 0, 379, 382, 3, 68, 34, 0, 380, 382, 3, 96,
		48, 0, 381, 378, 1, 0, 0, 0, 381, 379, 1, 0, 0, 0, 381, 380, 1, 0, 0, 0,
		382, 61, 1, 0, 0, 0, 383, 384, 5, 38, 0, 0, 384, 385, 3, 74, 37, 0, 385,
		63, 1, 0, 0, 0, 386, 392, 5, 39, 0, 0, 387, 393, 3, 82, 41, 0, 388, 393,
		3, 68, 34, 0, 389, 393, 3, 66, 33, 0, 390, 393, 3, 96, 48, 0, 391, 393,
		3, 102, 51, 0, 392, 387, 1, 0, 0, 0, 392, 388, 1, 0, 0, 0, 392, 389, 1,
		0, 0, 0, 392, 390, 1, 0, 0, 0, 392, 391, 1, 0, 0, 0, 393, 65, 1, 0, 0,
		0, 394, 395, 5, 64, 0, 0, 395, 399, 5, 65, 0, 0, 396, 397, 5, 64, 0, 0,
		397, 399, 3, 110, 55, 0, 398, 394, 1, 0, 0, 0, 398, 396, 1, 0, 0, 0, 399,
		67, 1, 0, 0, 0, 400, 403, 5, 65, 0, 0, 401, 403, 3, 110, 55, 0, 402, 400,
		1, 0, 0, 0, 402, 401, 1, 0, 0, 0, 403, 69, 1, 0, 0, 0, 404, 412, 3, 72,
		36, 0, 405, 412, 3, 74, 37, 0, 406, 412, 3, 76, 38, 0, 407, 412, 3, 78,
		39, 0, 408, 412, 3, 80, 40, 0, 409, 412, 3, 82, 41, 0, 410, 412, 3, 84,
		42, 0, 411, 404, 1, 0, 0, 0, 411, 405, 1, 0, 0, 0, 411, 406, 1, 0, 0, 0,
		411, 407, 1, 0, 0, 0, 411, 408, 1, 0, 0, 0, 411, 409, 1, 0, 0, 0, 411,
		410, 1, 0, 0, 0, 412, 71, 1, 0, 0, 0, 413, 415, 5, 9, 0, 0, 414, 416, 3,
		106, 53, 0, 415, 414, 1, 0, 0, 0, 415, 416, 1, 0, 0, 0, 416, 417, 1, 0,
		0, 0, 417, 418, 5, 10, 0, 0, 418, 73, 1, 0, 0, 0, 419, 431, 5, 13, 0, 0,
		420, 425, 3, 86, 43, 0, 421, 422, 5, 8, 0, 0, 422, 424, 3, 86, 43, 0, 423,
		421, 1, 0, 0, 0, 424, 427, 1, 0, 0, 0, 425, 423, 1, 0, 0, 0, 425, 426,
		1, 0, 0, 0, 426, 429, 1, 0, 0, 0, 427, 425, 1, 0, 0, 0, 428, 430, 5, 8,
		0, 0, 429, 428, 1, 0, 0, 0, 429, 430, 1, 0, 0, 0, 430, 432, 1, 0, 0, 0,
		431, 420, 1, 0, 0, 0, 431, 432, 1, 0, 0, 0, 432, 433, 1, 0, 0, 0, 433,
		434, 5, 14, 0, 0, 434, 75, 1, 0, 0, 0, 435, 436, 5, 50, 0, 0, 436, 77,
		1, 0, 0, 0, 437, 438, 5, 67, 0, 0, 438, 79, 1, 0, 0, 0, 439, 440, 5, 69,
		0, 0, 440, 81, 1, 0, 0, 0, 441, 442, 5, 68, 0, 0, 442, 83, 1, 0, 0, 0,
		443, 444, 7, 1, 0, 0, 444, 85, 1, 0, 0, 0, 445, 446, 3, 90, 45, 0, 446,
		447, 5, 5, 0, 0, 447, 448, 3, 118, 59, 0, 448, 455, 1, 0, 0, 0, 449, 450,
		3, 88, 44, 0, 450, 451, 5, 5, 0, 0, 451, 452, 3, 118, 59, 0, 452, 455,
		1, 0, 0, 0, 453, 455, 3, 68, 34, 0, 454, 445, 1, 0, 0, 0, 454, 449, 1,
		0, 0, 0, 454, 453, 1, 0, 0, 0, 455, 87, 1, 0, 0, 0, 456, 457, 5, 9, 0,
		0, 457, 458, 3, 118, 59, 0, 458, 459, 5, 10, 0, 0, 459, 89, 1, 0, 0, 0,
		460, 466, 5, 65, 0, 0, 461, 466, 3, 78, 39, 0, 462, 466, 3, 66, 33, 0,
		463, 466, 3, 110, 55, 0, 464, 466, 3, 112, 56, 0, 465, 460, 1, 0, 0, 0,
		465, 461, 1, 0, 0, 0, 465, 462, 1, 0, 0, 0, 465, 463, 1, 0, 0, 0, 465,
		464, 1, 0, 0, 0, 466, 91, 1, 0, 0, 0, 467, 468, 3, 94, 47, 0, 468, 469,
		5, 65, 0, 0, 469, 93, 1, 0, 0, 0, 470, 472, 5, 70, 0, 0, 471, 470, 1, 0,
		0, 0, 472, 475, 1, 0, 0, 0, 473, 471, 1, 0, 0, 0, 473, 474, 1, 0, 0, 0,
		474, 95, 1, 0, 0, 0, 475, 473, 1, 0, 0, 0, 476, 478, 3, 98, 49, 0, 477,
		479, 3, 108, 54, 0, 478, 477, 1, 0, 0, 0, 479, 480, 1, 0, 0, 0, 480, 478,
		1, 0, 0, 0, 480, 481, 1, 0, 0, 0, 481, 97, 1, 0, 0, 0, 482, 488, 3, 68,
		34, 0, 483, 488, 3, 66, 33, 0, 484, 488, 3, 72, 36, 0, 485, 488, 3, 74,
		37, 0, 486, 488, 3, 102, 51, 0, 487, 482, 1, 0, 0, 0, 487, 483, 1, 0, 0,
		0, 487, 484, 1, 0, 0, 0, 487, 485, 1, 0, 0, 0, 487, 486, 1, 0, 0, 0, 488,
		99, 1, 0, 0, 0, 489, 491, 3, 102, 51, 0, 490, 492, 3, 144, 72, 0, 491,
		490, 1, 0, 0, 0, 491, 492, 1, 0, 0, 0, 492, 101, 1, 0, 0, 0, 493, 494,
		3, 94, 47, 0, 494, 495, 3, 104, 52, 0, 495, 497, 5, 11, 0, 0, 496, 498,
		3, 106, 53, 0, 497, 496, 1, 0, 0, 0, 497, 498, 1, 0, 0, 0, 498, 499, 1,
		0, 0, 0, 499, 500, 5, 12, 0, 0, 500, 103, 1, 0, 0, 0, 501, 505, 5, 65,
		0, 0, 502, 505, 3, 110, 55, 0, 503, 505, 3, 112, 56, 0, 504, 501, 1, 0,
		0, 0, 504, 502, 1, 0, 0, 0, 504, 503, 1, 0, 0, 0, 505, 105, 1, 0, 0, 0,
		506, 511, 3, 118, 59, 0, 507, 508, 5, 8, 0, 0, 508, 510, 3, 118, 59, 0,
		509, 507, 1, 0, 0, 0, 510, 513, 1, 0, 0, 0, 511, 509, 1, 0, 0, 0, 511,
		512, 1, 0, 0, 0, 512, 515, 1, 0, 0, 0, 513, 511, 1, 0, 0, 0, 514, 516,
		5, 8, 0, 0, 515, 514, 1, 0, 0, 0, 515, 516, 1, 0, 0, 0, 516, 107, 1, 0,
		0, 0, 517, 519, 3, 144, 72, 0, 518, 517, 1, 0, 0, 0, 518, 519, 1, 0, 0,
		0, 519, 520, 1, 0, 0, 0, 520, 521, 5, 7, 0, 0, 521, 529, 3, 90, 45, 0,
		522, 523, 3, 144, 72, 0, 523, 524, 5, 7, 0, 0, 524, 526, 1, 0, 0, 0, 525,
		522, 1, 0, 0, 0, 525, 526, 1, 0, 0, 0, 526, 527, 1, 0, 0, 0, 527, 529,
		3, 88, 44, 0, 528, 518, 1, 0, 0, 0, 528, 525, 1, 0, 0, 0, 529, 109, 1,
		0, 0, 0, 530, 531, 7, 2, 0, 0, 531, 111, 1, 0, 0, 0, 532, 533, 7, 3, 0,
		0, 533, 113, 1, 0, 0, 0, 534, 535, 3, 116, 58, 0, 535, 536, 5, 30, 0, 0,
		536, 537, 3, 116, 58, 0, 537, 115, 1, 0, 0, 0, 538, 542, 3, 82, 41, 0,
		539, 542, 3, 68, 34, 0, 540, 542, 3, 66, 33, 0, 541, 538, 1, 0, 0, 0, 541,
		539, 1, 0, 0, 0, 541, 540, 1, 0, 0, 0, 542, 117, 1, 0, 0, 0, 543, 544,
		6, 59, -1, 0, 544, 545, 3, 132, 66, 0, 545, 546, 3, 118, 59, 5, 546, 549,
		1, 0, 0, 0, 547, 549, 3, 120, 60, 0, 548, 543, 1, 0, 0, 0, 548, 547, 1,
		0, 0, 0, 549, 567, 1, 0, 0, 0, 550, 551, 10, 4, 0, 0, 551, 552, 3, 136,
		68, 0, 552, 553, 3, 118, 59, 5, 553, 566, 1, 0, 0, 0, 554, 555, 10, 3,
		0, 0, 555, 556, 3, 138, 69, 0, 556, 557, 3, 118, 59, 4, 557, 566, 1, 0,
		0, 0, 558, 559, 10, 2, 0, 0, 559, 561, 5, 32, 0, 0, 560, 562, 3, 118, 59,
		0, 561, 560, 1, 0, 0, 0, 561, 562, 1, 0, 0, 0, 562, 563, 1, 0, 0, 0, 563,
		564, 5, 5, 0, 0, 564, 566, 3, 118, 59, 3, 565, 550, 1, 0, 0, 0, 565, 554,
		1, 0, 0, 0, 565, 558, 1, 0, 0, 0, 566, 569, 1, 0, 0, 0, 567, 565, 1, 0,
		0, 0, 567, 568, 1, 0, 0, 0, 568, 119, 1, 0, 0, 0, 569, 567, 1, 0, 0, 0,
		570, 571, 6, 60, -1, 0, 571, 572, 3, 122, 61, 0, 572, 591, 1, 0, 0, 0,
		573, 574, 10, 5, 0, 0, 574, 575, 3, 126, 63, 0, 575, 576, 3, 120, 60, 6,
		576, 590, 1, 0, 0, 0, 577, 578, 10, 4, 0, 0, 578, 579, 3, 124, 62, 0, 579,
		580, 3, 120, 60, 5, 580, 590, 1, 0, 0, 0, 581, 582, 10, 3, 0, 0, 582, 583,
		3, 128, 64, 0, 583, 584, 3, 120, 60, 4, 584, 590, 1, 0, 0, 0, 585, 586,
		10, 2, 0, 0, 586, 587, 3, 130, 65, 0, 587, 588, 3, 120, 60, 3, 588, 590,
		1, 0, 0, 0, 589, 573, 1, 0, 0, 0, 589, 577, 1, 0, 0, 0, 589, 581, 1, 0,
		0, 0, 589, 585, 1, 0, 0, 0, 590, 593, 1, 0, 0, 0, 591, 589, 1, 0, 0, 0,
		591, 592, 1, 0, 0, 0, 592, 121, 1, 0, 0, 0, 593, 591, 1, 0, 0, 0, 594,
		595, 6, 61, -1, 0, 595, 612, 3, 100, 50, 0, 596, 612, 3, 114, 57, 0, 597,
		612, 3, 70, 35, 0, 598, 612, 3, 68, 34, 0, 599, 612, 3, 96, 48, 0, 600,
		612, 3, 66, 33, 0, 601, 605, 5, 11, 0, 0, 602, 606, 3, 18, 9, 0, 603, 606,
		3, 56, 28, 0, 604, 606, 3, 118, 59, 0, 605, 602, 1, 0, 0, 0, 605, 603,
		1, 0, 0, 0, 605, 604, 1, 0, 0, 0, 606, 607, 1, 0, 0, 0, 607, 609, 5, 12,
		0, 0, 608, 610, 3, 144, 72, 0, 609, 608, 1, 0, 0, 0, 609, 610, 1, 0, 0,
		0, 610, 612, 1, 0, 0, 0, 611, 594, 1, 0, 0, 0, 611, 596, 1, 0, 0, 0, 611,
		597, 1, 0, 0, 0, 611, 598, 1, 0, 0, 0, 611, 599, 1, 0, 0, 0, 611, 600,
		1, 0, 0, 0, 611, 601, 1, 0, 0, 0, 612, 627, 1, 0, 0, 0, 613, 614, 10, 10,
		0, 0, 614, 615, 3, 140, 70, 0, 615, 616, 3, 122, 61, 11, 616, 626, 1, 0,
		0, 0, 617, 618, 10, 9, 0, 0, 618, 619, 3, 142, 71, 0, 619, 620, 3, 122,
		61, 10, 620, 626, 1, 0, 0, 0, 621, 622, 10, 8, 0, 0, 622, 623, 3, 134,
		67, 0, 623, 624, 3, 122, 61, 9, 624, 626, 1, 0, 0, 0, 625, 613, 1, 0, 0,
		0, 625, 617, 1, 0, 0, 0, 625, 621, 1, 0, 0, 0, 626, 629, 1, 0, 0, 0, 627,
		625, 1, 0, 0, 0, 627, 628, 1, 0, 0, 0, 628, 123, 1, 0, 0, 0, 629, 627,
		1, 0, 0, 0, 630, 633, 7, 4, 0, 0, 631, 634, 3, 128, 64, 0, 632, 634, 3,
		126, 63, 0, 633, 631, 1, 0, 0, 0, 633, 632, 1, 0, 0, 0, 634, 125, 1, 0,
		0, 0, 635, 636, 7, 5, 0, 0, 636, 127, 1, 0, 0, 0, 637, 639, 5, 60, 0, 0,
		638, 637, 1, 0, 0, 0, 638, 639, 1, 0, 0, 0, 639, 640, 1, 0, 0, 0, 640,
		641, 5, 61, 0, 0, 641, 129, 1, 0, 0, 0, 642, 644, 5, 60, 0, 0, 643, 642,
		1, 0, 0, 0, 643, 644, 1, 0, 0, 0, 644, 645, 1, 0, 0, 0, 645, 646, 5, 59,
		0, 0, 646, 131, 1, 0, 0, 0, 647, 648, 7, 6, 0, 0, 648, 133, 1, 0, 0, 0,
		649, 650, 7, 7, 0, 0, 650, 135, 1, 0, 0, 0, 651, 652, 5, 28, 0, 0, 652,
		137, 1, 0, 0, 0, 653, 654, 5, 29, 0, 0, 654, 139, 1, 0, 0, 0, 655, 656,
		7, 8, 0, 0, 656, 141, 1, 0, 0, 0, 657, 658, 7, 9, 0, 0, 658, 143, 1, 0,
		0, 0, 659, 660, 5, 32, 0, 0, 660, 145, 1, 0, 0, 0, 67, 149, 164, 172, 176,
		187, 191, 199, 206, 214, 221, 226, 235, 241, 245, 249, 253, 262, 269, 277,
		282, 307, 318, 331, 339, 341, 349, 363, 366, 369, 376, 381, 392, 398, 402,
		411, 415, 425, 429, 431, 454, 465, 473, 480, 487, 491, 497, 504, 511, 515,
		518, 525, 528, 541, 548, 561, 565, 567, 589, 591, 605, 609, 611, 625, 627,
		633, 638, 643,
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
	FqlParserMinusMinus        = 26
	FqlParserPlusPlus          = 27
	FqlParserAnd               = 28
	FqlParserOr                = 29
	FqlParserRange             = 30
	FqlParserAssign            = 31
	FqlParserQuestionMark      = 32
	FqlParserRegexNotMatch     = 33
	FqlParserRegexMatch        = 34
	FqlParserFor               = 35
	FqlParserReturn            = 36
	FqlParserWaitfor           = 37
	FqlParserOptions           = 38
	FqlParserTimeout           = 39
	FqlParserDistinct          = 40
	FqlParserFilter            = 41
	FqlParserCurrent           = 42
	FqlParserSort              = 43
	FqlParserLimit             = 44
	FqlParserLet               = 45
	FqlParserCollect           = 46
	FqlParserSortDirection     = 47
	FqlParserNone              = 48
	FqlParserNull              = 49
	FqlParserBooleanLiteral    = 50
	FqlParserUse               = 51
	FqlParserInto              = 52
	FqlParserKeep              = 53
	FqlParserWith              = 54
	FqlParserAll               = 55
	FqlParserAny               = 56
	FqlParserAggregate         = 57
	FqlParserEvent             = 58
	FqlParserLike              = 59
	FqlParserNot               = 60
	FqlParserIn                = 61
	FqlParserDo                = 62
	FqlParserWhile             = 63
	FqlParserParam             = 64
	FqlParserIdentifier        = 65
	FqlParserIgnoreIdentifier  = 66
	FqlParserStringLiteral     = 67
	FqlParserIntegerLiteral    = 68
	FqlParserFloatLiteral      = 69
	FqlParserNamespaceSegment  = 70
	FqlParserUnknownIdentifier = 71
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

	// SetValueVariable sets the valueVariable token.
	SetValueVariable(antlr.Token)

	// SetCounterVariable sets the counterVariable token.
	SetCounterVariable(antlr.Token)

	// Getter signatures
	For() antlr.TerminalNode
	In() antlr.TerminalNode
	ForExpressionSource() IForExpressionSourceContext
	ForExpressionReturn() IForExpressionReturnContext
	AllIdentifier() []antlr.TerminalNode
	Identifier(i int) antlr.TerminalNode
	IgnoreIdentifier() antlr.TerminalNode
	Comma() antlr.TerminalNode
	AllForExpressionBody() []IForExpressionBodyContext
	ForExpressionBody(i int) IForExpressionBodyContext
	While() antlr.TerminalNode
	Expression() IExpressionContext
	Do() antlr.TerminalNode

	// IsForExpressionContext differentiates from other interfaces.
	IsForExpressionContext()
}

type ForExpressionContext struct {
	antlr.BaseParserRuleContext
	parser          antlr.Parser
	valueVariable   antlr.Token
	counterVariable antlr.Token
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

func (s *ForExpressionContext) SetValueVariable(v antlr.Token) { s.valueVariable = v }

func (s *ForExpressionContext) SetCounterVariable(v antlr.Token) { s.counterVariable = v }

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

func (s *ForExpressionContext) IgnoreIdentifier() antlr.TerminalNode {
	return s.GetToken(FqlParserIgnoreIdentifier, 0)
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

func (s *ForExpressionContext) While() antlr.TerminalNode {
	return s.GetToken(FqlParserWhile, 0)
}

func (s *ForExpressionContext) Expression() IExpressionContext {
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

	p.SetState(226)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 10, p.GetParserRuleContext()) {
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

			localctx.(*ForExpressionContext).counterVariable = _lt

			_la = p.GetTokenStream().LA(1)

			if !(_la == FqlParserIdentifier || _la == FqlParserIgnoreIdentifier) {
				var _ri = p.GetErrorHandler().RecoverInline(p)

				localctx.(*ForExpressionContext).counterVariable = _ri
			} else {
				p.GetErrorHandler().ReportMatch(p)
				p.Consume()
			}
		}
		p.SetState(214)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)

		if _la == FqlParserDo {
			{
				p.SetState(213)
				p.Match(FqlParserDo)
				if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
				}
			}

		}
		{
			p.SetState(216)
			p.Match(FqlParserWhile)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(217)
			p.expression(0)
		}
		p.SetState(221)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_alt = p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 9, p.GetParserRuleContext())
		if p.HasError() {
			goto errorExit
		}
		for _alt != 2 && _alt != antlr.ATNInvalidAltNumber {
			if _alt == 1 {
				{
					p.SetState(218)
					p.ForExpressionBody()
				}

			}
			p.SetState(223)
			p.GetErrorHandler().Sync(p)
			if p.HasError() {
				goto errorExit
			}
			_alt = p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 9, p.GetParserRuleContext())
			if p.HasError() {
				goto errorExit
			}
		}
		{
			p.SetState(224)
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
	p.SetState(235)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 11, p.GetParserRuleContext()) {
	case 1:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(228)
			p.FunctionCallExpression()
		}

	case 2:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(229)
			p.ArrayLiteral()
		}

	case 3:
		p.EnterOuterAlt(localctx, 3)
		{
			p.SetState(230)
			p.ObjectLiteral()
		}

	case 4:
		p.EnterOuterAlt(localctx, 4)
		{
			p.SetState(231)
			p.Variable()
		}

	case 5:
		p.EnterOuterAlt(localctx, 5)
		{
			p.SetState(232)
			p.MemberExpression()
		}

	case 6:
		p.EnterOuterAlt(localctx, 6)
		{
			p.SetState(233)
			p.RangeOperator()
		}

	case 7:
		p.EnterOuterAlt(localctx, 7)
		{
			p.SetState(234)
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
	p.SetState(241)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetTokenStream().LA(1) {
	case FqlParserLimit:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(237)
			p.LimitClause()
		}

	case FqlParserSort:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(238)
			p.SortClause()
		}

	case FqlParserFilter:
		p.EnterOuterAlt(localctx, 3)
		{
			p.SetState(239)
			p.FilterClause()
		}

	case FqlParserCollect:
		p.EnterOuterAlt(localctx, 4)
		{
			p.SetState(240)
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
	p.SetState(245)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 13, p.GetParserRuleContext()) {
	case 1:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(243)
			p.VariableDeclaration()
		}

	case 2:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(244)
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
	p.SetState(249)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 14, p.GetParserRuleContext()) {
	case 1:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(247)
			p.ForExpressionStatement()
		}

	case 2:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(248)
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
	p.SetState(253)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetTokenStream().LA(1) {
	case FqlParserReturn:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(251)
			p.ReturnExpression()
		}

	case FqlParserFor:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(252)
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
		p.SetState(255)
		p.Match(FqlParserFilter)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(256)
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
		p.SetState(258)
		p.Match(FqlParserLimit)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(259)
		p.LimitClauseValue()
	}
	p.SetState(262)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	if _la == FqlParserComma {
		{
			p.SetState(260)
			p.Match(FqlParserComma)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(261)
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
	p.SetState(269)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 17, p.GetParserRuleContext()) {
	case 1:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(264)
			p.IntegerLiteral()
		}

	case 2:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(265)
			p.Param()
		}

	case 3:
		p.EnterOuterAlt(localctx, 3)
		{
			p.SetState(266)
			p.Variable()
		}

	case 4:
		p.EnterOuterAlt(localctx, 4)
		{
			p.SetState(267)
			p.FunctionCallExpression()
		}

	case 5:
		p.EnterOuterAlt(localctx, 5)
		{
			p.SetState(268)
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
		p.SetState(271)
		p.Match(FqlParserSort)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(272)
		p.SortClauseExpression()
	}
	p.SetState(277)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	for _la == FqlParserComma {
		{
			p.SetState(273)
			p.Match(FqlParserComma)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(274)
			p.SortClauseExpression()
		}

		p.SetState(279)
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
		p.SetState(280)
		p.expression(0)
	}
	p.SetState(282)
	p.GetErrorHandler().Sync(p)

	if p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 19, p.GetParserRuleContext()) == 1 {
		{
			p.SetState(281)
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
	p.SetState(307)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 20, p.GetParserRuleContext()) {
	case 1:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(284)
			p.Match(FqlParserCollect)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(285)
			p.CollectGrouping()
		}
		{
			p.SetState(286)
			p.CollectAggregator()
		}
		{
			p.SetState(287)
			p.CollectGroupProjection()
		}

	case 2:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(289)
			p.Match(FqlParserCollect)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(290)
			p.CollectGrouping()
		}
		{
			p.SetState(291)
			p.CollectAggregator()
		}

	case 3:
		p.EnterOuterAlt(localctx, 3)
		{
			p.SetState(293)
			p.Match(FqlParserCollect)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(294)
			p.CollectGrouping()
		}
		{
			p.SetState(295)
			p.CollectGroupProjection()
		}

	case 4:
		p.EnterOuterAlt(localctx, 4)
		{
			p.SetState(297)
			p.Match(FqlParserCollect)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(298)
			p.CollectGrouping()
		}

	case 5:
		p.EnterOuterAlt(localctx, 5)
		{
			p.SetState(299)
			p.Match(FqlParserCollect)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(300)
			p.CollectAggregator()
		}
		{
			p.SetState(301)
			p.CollectGroupProjection()
		}

	case 6:
		p.EnterOuterAlt(localctx, 6)
		{
			p.SetState(303)
			p.Match(FqlParserCollect)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(304)
			p.CollectAggregator()
		}

	case 7:
		p.EnterOuterAlt(localctx, 7)
		{
			p.SetState(305)
			p.Match(FqlParserCollect)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(306)
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
		p.SetState(309)
		p.Match(FqlParserIdentifier)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(310)
		p.Match(FqlParserAssign)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(311)
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
		p.SetState(313)
		p.CollectSelector()
	}
	p.SetState(318)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	for _la == FqlParserComma {
		{
			p.SetState(314)
			p.Match(FqlParserComma)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(315)
			p.CollectSelector()
		}

		p.SetState(320)
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
		p.SetState(321)
		p.Match(FqlParserIdentifier)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(322)
		p.Match(FqlParserAssign)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(323)
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
		p.SetState(325)
		p.Match(FqlParserAggregate)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(326)
		p.CollectAggregateSelector()
	}
	p.SetState(331)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	for _la == FqlParserComma {
		{
			p.SetState(327)
			p.Match(FqlParserComma)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(328)
			p.CollectAggregateSelector()
		}

		p.SetState(333)
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
	p.SetState(341)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 24, p.GetParserRuleContext()) {
	case 1:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(334)
			p.Match(FqlParserInto)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(335)
			p.CollectSelector()
		}

	case 2:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(336)
			p.Match(FqlParserInto)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(337)
			p.Match(FqlParserIdentifier)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		p.SetState(339)
		p.GetErrorHandler().Sync(p)

		if p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 23, p.GetParserRuleContext()) == 1 {
			{
				p.SetState(338)
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
		p.SetState(343)
		p.Match(FqlParserKeep)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(344)
		p.Match(FqlParserIdentifier)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	p.SetState(349)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	for _la == FqlParserComma {
		{
			p.SetState(345)
			p.Match(FqlParserComma)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(346)
			p.Match(FqlParserIdentifier)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

		p.SetState(351)
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
		p.SetState(352)
		p.Match(FqlParserWith)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(353)
		p.Match(FqlParserIdentifier)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(354)
		p.Match(FqlParserInto)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(355)
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
		p.SetState(357)
		p.Match(FqlParserWaitfor)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(358)
		p.Match(FqlParserEvent)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(359)
		p.WaitForEventName()
	}
	{
		p.SetState(360)
		p.Match(FqlParserIn)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(361)
		p.WaitForEventSource()
	}
	p.SetState(363)
	p.GetErrorHandler().Sync(p)

	if p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 26, p.GetParserRuleContext()) == 1 {
		{
			p.SetState(362)
			p.OptionsClause()
		}

	} else if p.HasError() { // JIM
		goto errorExit
	}
	p.SetState(366)
	p.GetErrorHandler().Sync(p)

	if p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 27, p.GetParserRuleContext()) == 1 {
		{
			p.SetState(365)
			p.FilterClause()
		}

	} else if p.HasError() { // JIM
		goto errorExit
	}
	p.SetState(369)
	p.GetErrorHandler().Sync(p)

	if p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 28, p.GetParserRuleContext()) == 1 {
		{
			p.SetState(368)
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
	p.SetState(376)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 29, p.GetParserRuleContext()) {
	case 1:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(371)
			p.StringLiteral()
		}

	case 2:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(372)
			p.Variable()
		}

	case 3:
		p.EnterOuterAlt(localctx, 3)
		{
			p.SetState(373)
			p.Param()
		}

	case 4:
		p.EnterOuterAlt(localctx, 4)
		{
			p.SetState(374)
			p.FunctionCallExpression()
		}

	case 5:
		p.EnterOuterAlt(localctx, 5)
		{
			p.SetState(375)
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
	p.SetState(381)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 30, p.GetParserRuleContext()) {
	case 1:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(378)
			p.FunctionCallExpression()
		}

	case 2:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(379)
			p.Variable()
		}

	case 3:
		p.EnterOuterAlt(localctx, 3)
		{
			p.SetState(380)
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
		p.SetState(383)
		p.Match(FqlParserOptions)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(384)
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
		p.SetState(386)
		p.Match(FqlParserTimeout)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	p.SetState(392)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 31, p.GetParserRuleContext()) {
	case 1:
		{
			p.SetState(387)
			p.IntegerLiteral()
		}

	case 2:
		{
			p.SetState(388)
			p.Variable()
		}

	case 3:
		{
			p.SetState(389)
			p.Param()
		}

	case 4:
		{
			p.SetState(390)
			p.MemberExpression()
		}

	case 5:
		{
			p.SetState(391)
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
	p.SetState(398)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 32, p.GetParserRuleContext()) {
	case 1:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(394)
			p.Match(FqlParserParam)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(395)
			p.Match(FqlParserIdentifier)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	case 2:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(396)
			p.Match(FqlParserParam)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(397)
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
	p.SetState(402)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetTokenStream().LA(1) {
	case FqlParserIdentifier:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(400)
			p.Match(FqlParserIdentifier)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	case FqlParserAnd, FqlParserOr, FqlParserOptions, FqlParserTimeout, FqlParserDistinct, FqlParserFilter, FqlParserCurrent, FqlParserSort, FqlParserLimit, FqlParserCollect, FqlParserSortDirection, FqlParserInto, FqlParserKeep, FqlParserWith, FqlParserAll, FqlParserAny, FqlParserAggregate, FqlParserEvent:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(401)
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
	p.SetState(411)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetTokenStream().LA(1) {
	case FqlParserOpenBracket:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(404)
			p.ArrayLiteral()
		}

	case FqlParserOpenBrace:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(405)
			p.ObjectLiteral()
		}

	case FqlParserBooleanLiteral:
		p.EnterOuterAlt(localctx, 3)
		{
			p.SetState(406)
			p.BooleanLiteral()
		}

	case FqlParserStringLiteral:
		p.EnterOuterAlt(localctx, 4)
		{
			p.SetState(407)
			p.StringLiteral()
		}

	case FqlParserFloatLiteral:
		p.EnterOuterAlt(localctx, 5)
		{
			p.SetState(408)
			p.FloatLiteral()
		}

	case FqlParserIntegerLiteral:
		p.EnterOuterAlt(localctx, 6)
		{
			p.SetState(409)
			p.IntegerLiteral()
		}

	case FqlParserNone, FqlParserNull:
		p.EnterOuterAlt(localctx, 7)
		{
			p.SetState(410)
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
		p.SetState(413)
		p.Match(FqlParserOpenBracket)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	p.SetState(415)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	if (int64((_la-9)) & ^0x3f) == 0 && ((int64(1)<<(_la-9))&4467570830286094357) != 0 {
		{
			p.SetState(414)
			p.ArgumentList()
		}

	}
	{
		p.SetState(417)
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
		p.SetState(419)
		p.Match(FqlParserOpenBrace)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	p.SetState(431)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	if (int64((_la-9)) & ^0x3f) == 0 && ((int64(1)<<(_la-9))&432345564162031617) != 0 {
		{
			p.SetState(420)
			p.PropertyAssignment()
		}
		p.SetState(425)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_alt = p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 36, p.GetParserRuleContext())
		if p.HasError() {
			goto errorExit
		}
		for _alt != 2 && _alt != antlr.ATNInvalidAltNumber {
			if _alt == 1 {
				{
					p.SetState(421)
					p.Match(FqlParserComma)
					if p.HasError() {
						// Recognition error - abort rule
						goto errorExit
					}
				}
				{
					p.SetState(422)
					p.PropertyAssignment()
				}

			}
			p.SetState(427)
			p.GetErrorHandler().Sync(p)
			if p.HasError() {
				goto errorExit
			}
			_alt = p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 36, p.GetParserRuleContext())
			if p.HasError() {
				goto errorExit
			}
		}
		p.SetState(429)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)

		if _la == FqlParserComma {
			{
				p.SetState(428)
				p.Match(FqlParserComma)
				if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
				}
			}

		}

	}
	{
		p.SetState(433)
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
		p.SetState(435)
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
		p.SetState(437)
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
		p.SetState(439)
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
		p.SetState(441)
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
		p.SetState(443)
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
	p.SetState(454)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 39, p.GetParserRuleContext()) {
	case 1:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(445)
			p.PropertyName()
		}
		{
			p.SetState(446)
			p.Match(FqlParserColon)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(447)
			p.expression(0)
		}

	case 2:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(449)
			p.ComputedPropertyName()
		}
		{
			p.SetState(450)
			p.Match(FqlParserColon)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(451)
			p.expression(0)
		}

	case 3:
		p.EnterOuterAlt(localctx, 3)
		{
			p.SetState(453)
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
		p.SetState(456)
		p.Match(FqlParserOpenBracket)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(457)
		p.expression(0)
	}
	{
		p.SetState(458)
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
	p.SetState(465)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetTokenStream().LA(1) {
	case FqlParserIdentifier:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(460)
			p.Match(FqlParserIdentifier)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	case FqlParserStringLiteral:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(461)
			p.StringLiteral()
		}

	case FqlParserParam:
		p.EnterOuterAlt(localctx, 3)
		{
			p.SetState(462)
			p.Param()
		}

	case FqlParserAnd, FqlParserOr, FqlParserOptions, FqlParserTimeout, FqlParserDistinct, FqlParserFilter, FqlParserCurrent, FqlParserSort, FqlParserLimit, FqlParserCollect, FqlParserSortDirection, FqlParserInto, FqlParserKeep, FqlParserWith, FqlParserAll, FqlParserAny, FqlParserAggregate, FqlParserEvent:
		p.EnterOuterAlt(localctx, 4)
		{
			p.SetState(463)
			p.SafeReservedWord()
		}

	case FqlParserFor, FqlParserReturn, FqlParserWaitfor, FqlParserLet, FqlParserNone, FqlParserNull, FqlParserBooleanLiteral, FqlParserUse, FqlParserLike, FqlParserNot, FqlParserIn, FqlParserDo, FqlParserWhile:
		p.EnterOuterAlt(localctx, 5)
		{
			p.SetState(464)
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
		p.SetState(467)
		p.Namespace()
	}
	{
		p.SetState(468)
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
	p.SetState(473)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	for _la == FqlParserNamespaceSegment {
		{
			p.SetState(470)
			p.Match(FqlParserNamespaceSegment)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

		p.SetState(475)
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
		p.SetState(476)
		p.MemberExpressionSource()
	}
	p.SetState(478)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_alt = 1
	for ok := true; ok; ok = _alt != 2 && _alt != antlr.ATNInvalidAltNumber {
		switch _alt {
		case 1:
			{
				p.SetState(477)
				p.MemberExpressionPath()
			}

		default:
			p.SetError(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
			goto errorExit
		}

		p.SetState(480)
		p.GetErrorHandler().Sync(p)
		_alt = p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 42, p.GetParserRuleContext())
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
	p.SetState(487)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 43, p.GetParserRuleContext()) {
	case 1:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(482)
			p.Variable()
		}

	case 2:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(483)
			p.Param()
		}

	case 3:
		p.EnterOuterAlt(localctx, 3)
		{
			p.SetState(484)
			p.ArrayLiteral()
		}

	case 4:
		p.EnterOuterAlt(localctx, 4)
		{
			p.SetState(485)
			p.ObjectLiteral()
		}

	case 5:
		p.EnterOuterAlt(localctx, 5)
		{
			p.SetState(486)
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
		p.SetState(489)
		p.FunctionCall()
	}
	p.SetState(491)
	p.GetErrorHandler().Sync(p)

	if p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 44, p.GetParserRuleContext()) == 1 {
		{
			p.SetState(490)
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
		p.SetState(493)
		p.Namespace()
	}
	{
		p.SetState(494)
		p.FunctionName()
	}
	{
		p.SetState(495)
		p.Match(FqlParserOpenParen)
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

	if (int64((_la-9)) & ^0x3f) == 0 && ((int64(1)<<(_la-9))&4467570830286094357) != 0 {
		{
			p.SetState(496)
			p.ArgumentList()
		}

	}
	{
		p.SetState(499)
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
	p.SetState(504)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetTokenStream().LA(1) {
	case FqlParserIdentifier:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(501)
			p.Match(FqlParserIdentifier)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	case FqlParserAnd, FqlParserOr, FqlParserOptions, FqlParserTimeout, FqlParserDistinct, FqlParserFilter, FqlParserCurrent, FqlParserSort, FqlParserLimit, FqlParserCollect, FqlParserSortDirection, FqlParserInto, FqlParserKeep, FqlParserWith, FqlParserAll, FqlParserAny, FqlParserAggregate, FqlParserEvent:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(502)
			p.SafeReservedWord()
		}

	case FqlParserFor, FqlParserReturn, FqlParserWaitfor, FqlParserLet, FqlParserNone, FqlParserNull, FqlParserBooleanLiteral, FqlParserUse, FqlParserLike, FqlParserNot, FqlParserIn, FqlParserDo, FqlParserWhile:
		p.EnterOuterAlt(localctx, 3)
		{
			p.SetState(503)
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
		p.SetState(506)
		p.expression(0)
	}
	p.SetState(511)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_alt = p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 47, p.GetParserRuleContext())
	if p.HasError() {
		goto errorExit
	}
	for _alt != 2 && _alt != antlr.ATNInvalidAltNumber {
		if _alt == 1 {
			{
				p.SetState(507)
				p.Match(FqlParserComma)
				if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
				}
			}
			{
				p.SetState(508)
				p.expression(0)
			}

		}
		p.SetState(513)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_alt = p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 47, p.GetParserRuleContext())
		if p.HasError() {
			goto errorExit
		}
	}
	p.SetState(515)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	if _la == FqlParserComma {
		{
			p.SetState(514)
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

	p.SetState(528)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 51, p.GetParserRuleContext()) {
	case 1:
		p.EnterOuterAlt(localctx, 1)
		p.SetState(518)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)

		if _la == FqlParserQuestionMark {
			{
				p.SetState(517)
				p.ErrorOperator()
			}

		}
		{
			p.SetState(520)
			p.Match(FqlParserDot)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(521)
			p.PropertyName()
		}

	case 2:
		p.EnterOuterAlt(localctx, 2)
		p.SetState(525)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)

		if _la == FqlParserQuestionMark {
			{
				p.SetState(522)
				p.ErrorOperator()
			}
			{
				p.SetState(523)
				p.Match(FqlParserDot)
				if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
				}
			}

		}
		{
			p.SetState(527)
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
		p.SetState(530)
		_la = p.GetTokenStream().LA(1)

		if !((int64(_la) & ^0x3f) == 0 && ((int64(1)<<_la)&572203169208074240) != 0) {
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
		p.SetState(532)
		_la = p.GetTokenStream().LA(1)

		if !((int64(_la) & ^0x3f) == 0 && ((int64(1)<<_la)&-572203202762506240) != 0) {
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
		p.SetState(534)

		var _x = p.RangeOperand()

		localctx.(*RangeOperatorContext).left = _x
	}
	{
		p.SetState(535)
		p.Match(FqlParserRange)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(536)

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
	p.SetState(541)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetTokenStream().LA(1) {
	case FqlParserIntegerLiteral:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(538)
			p.IntegerLiteral()
		}

	case FqlParserAnd, FqlParserOr, FqlParserOptions, FqlParserTimeout, FqlParserDistinct, FqlParserFilter, FqlParserCurrent, FqlParserSort, FqlParserLimit, FqlParserCollect, FqlParserSortDirection, FqlParserInto, FqlParserKeep, FqlParserWith, FqlParserAll, FqlParserAny, FqlParserAggregate, FqlParserEvent, FqlParserIdentifier:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(539)
			p.Variable()
		}

	case FqlParserParam:
		p.EnterOuterAlt(localctx, 3)
		{
			p.SetState(540)
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
	p.SetState(548)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 53, p.GetParserRuleContext()) {
	case 1:
		{
			p.SetState(544)
			p.UnaryOperator()
		}
		{
			p.SetState(545)

			var _x = p.expression(5)

			localctx.(*ExpressionContext).right = _x
		}

	case 2:
		{
			p.SetState(547)
			p.predicate(0)
		}

	case antlr.ATNInvalidAltNumber:
		goto errorExit
	}
	p.GetParserRuleContext().SetStop(p.GetTokenStream().LT(-1))
	p.SetState(567)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_alt = p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 56, p.GetParserRuleContext())
	if p.HasError() {
		goto errorExit
	}
	for _alt != 2 && _alt != antlr.ATNInvalidAltNumber {
		if _alt == 1 {
			if p.GetParseListeners() != nil {
				p.TriggerExitRuleEvent()
			}
			_prevctx = localctx
			p.SetState(565)
			p.GetErrorHandler().Sync(p)
			if p.HasError() {
				goto errorExit
			}

			switch p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 55, p.GetParserRuleContext()) {
			case 1:
				localctx = NewExpressionContext(p, _parentctx, _parentState)
				localctx.(*ExpressionContext).left = _prevctx
				p.PushNewRecursionContext(localctx, _startState, FqlParserRULE_expression)
				p.SetState(550)

				if !(p.Precpred(p.GetParserRuleContext(), 4)) {
					p.SetError(antlr.NewFailedPredicateException(p, "p.Precpred(p.GetParserRuleContext(), 4)", ""))
					goto errorExit
				}
				{
					p.SetState(551)
					p.LogicalAndOperator()
				}
				{
					p.SetState(552)

					var _x = p.expression(5)

					localctx.(*ExpressionContext).right = _x
				}

			case 2:
				localctx = NewExpressionContext(p, _parentctx, _parentState)
				localctx.(*ExpressionContext).left = _prevctx
				p.PushNewRecursionContext(localctx, _startState, FqlParserRULE_expression)
				p.SetState(554)

				if !(p.Precpred(p.GetParserRuleContext(), 3)) {
					p.SetError(antlr.NewFailedPredicateException(p, "p.Precpred(p.GetParserRuleContext(), 3)", ""))
					goto errorExit
				}
				{
					p.SetState(555)
					p.LogicalOrOperator()
				}
				{
					p.SetState(556)

					var _x = p.expression(4)

					localctx.(*ExpressionContext).right = _x
				}

			case 3:
				localctx = NewExpressionContext(p, _parentctx, _parentState)
				localctx.(*ExpressionContext).condition = _prevctx
				p.PushNewRecursionContext(localctx, _startState, FqlParserRULE_expression)
				p.SetState(558)

				if !(p.Precpred(p.GetParserRuleContext(), 2)) {
					p.SetError(antlr.NewFailedPredicateException(p, "p.Precpred(p.GetParserRuleContext(), 2)", ""))
					goto errorExit
				}
				{
					p.SetState(559)

					var _m = p.Match(FqlParserQuestionMark)

					localctx.(*ExpressionContext).ternaryOperator = _m
					if p.HasError() {
						// Recognition error - abort rule
						goto errorExit
					}
				}
				p.SetState(561)
				p.GetErrorHandler().Sync(p)
				if p.HasError() {
					goto errorExit
				}
				_la = p.GetTokenStream().LA(1)

				if (int64((_la-9)) & ^0x3f) == 0 && ((int64(1)<<(_la-9))&4467570830286094357) != 0 {
					{
						p.SetState(560)

						var _x = p.expression(0)

						localctx.(*ExpressionContext).onTrue = _x
					}

				}
				{
					p.SetState(563)
					p.Match(FqlParserColon)
					if p.HasError() {
						// Recognition error - abort rule
						goto errorExit
					}
				}
				{
					p.SetState(564)

					var _x = p.expression(3)

					localctx.(*ExpressionContext).onFalse = _x
				}

			case antlr.ATNInvalidAltNumber:
				goto errorExit
			}

		}
		p.SetState(569)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_alt = p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 56, p.GetParserRuleContext())
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
		p.SetState(571)
		p.expressionAtom(0)
	}

	p.GetParserRuleContext().SetStop(p.GetTokenStream().LT(-1))
	p.SetState(591)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_alt = p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 58, p.GetParserRuleContext())
	if p.HasError() {
		goto errorExit
	}
	for _alt != 2 && _alt != antlr.ATNInvalidAltNumber {
		if _alt == 1 {
			if p.GetParseListeners() != nil {
				p.TriggerExitRuleEvent()
			}
			_prevctx = localctx
			p.SetState(589)
			p.GetErrorHandler().Sync(p)
			if p.HasError() {
				goto errorExit
			}

			switch p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 57, p.GetParserRuleContext()) {
			case 1:
				localctx = NewPredicateContext(p, _parentctx, _parentState)
				localctx.(*PredicateContext).left = _prevctx
				p.PushNewRecursionContext(localctx, _startState, FqlParserRULE_predicate)
				p.SetState(573)

				if !(p.Precpred(p.GetParserRuleContext(), 5)) {
					p.SetError(antlr.NewFailedPredicateException(p, "p.Precpred(p.GetParserRuleContext(), 5)", ""))
					goto errorExit
				}
				{
					p.SetState(574)
					p.EqualityOperator()
				}
				{
					p.SetState(575)

					var _x = p.predicate(6)

					localctx.(*PredicateContext).right = _x
				}

			case 2:
				localctx = NewPredicateContext(p, _parentctx, _parentState)
				localctx.(*PredicateContext).left = _prevctx
				p.PushNewRecursionContext(localctx, _startState, FqlParserRULE_predicate)
				p.SetState(577)

				if !(p.Precpred(p.GetParserRuleContext(), 4)) {
					p.SetError(antlr.NewFailedPredicateException(p, "p.Precpred(p.GetParserRuleContext(), 4)", ""))
					goto errorExit
				}
				{
					p.SetState(578)
					p.ArrayOperator()
				}
				{
					p.SetState(579)

					var _x = p.predicate(5)

					localctx.(*PredicateContext).right = _x
				}

			case 3:
				localctx = NewPredicateContext(p, _parentctx, _parentState)
				localctx.(*PredicateContext).left = _prevctx
				p.PushNewRecursionContext(localctx, _startState, FqlParserRULE_predicate)
				p.SetState(581)

				if !(p.Precpred(p.GetParserRuleContext(), 3)) {
					p.SetError(antlr.NewFailedPredicateException(p, "p.Precpred(p.GetParserRuleContext(), 3)", ""))
					goto errorExit
				}
				{
					p.SetState(582)
					p.InOperator()
				}
				{
					p.SetState(583)

					var _x = p.predicate(4)

					localctx.(*PredicateContext).right = _x
				}

			case 4:
				localctx = NewPredicateContext(p, _parentctx, _parentState)
				localctx.(*PredicateContext).left = _prevctx
				p.PushNewRecursionContext(localctx, _startState, FqlParserRULE_predicate)
				p.SetState(585)

				if !(p.Precpred(p.GetParserRuleContext(), 2)) {
					p.SetError(antlr.NewFailedPredicateException(p, "p.Precpred(p.GetParserRuleContext(), 2)", ""))
					goto errorExit
				}
				{
					p.SetState(586)
					p.LikeOperator()
				}
				{
					p.SetState(587)

					var _x = p.predicate(3)

					localctx.(*PredicateContext).right = _x
				}

			case antlr.ATNInvalidAltNumber:
				goto errorExit
			}

		}
		p.SetState(593)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_alt = p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 58, p.GetParserRuleContext())
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
	p.SetState(611)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 61, p.GetParserRuleContext()) {
	case 1:
		{
			p.SetState(595)
			p.FunctionCallExpression()
		}

	case 2:
		{
			p.SetState(596)
			p.RangeOperator()
		}

	case 3:
		{
			p.SetState(597)
			p.Literal()
		}

	case 4:
		{
			p.SetState(598)
			p.Variable()
		}

	case 5:
		{
			p.SetState(599)
			p.MemberExpression()
		}

	case 6:
		{
			p.SetState(600)
			p.Param()
		}

	case 7:
		{
			p.SetState(601)
			p.Match(FqlParserOpenParen)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		p.SetState(605)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}

		switch p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 59, p.GetParserRuleContext()) {
		case 1:
			{
				p.SetState(602)
				p.ForExpression()
			}

		case 2:
			{
				p.SetState(603)
				p.WaitForExpression()
			}

		case 3:
			{
				p.SetState(604)
				p.expression(0)
			}

		case antlr.ATNInvalidAltNumber:
			goto errorExit
		}
		{
			p.SetState(607)
			p.Match(FqlParserCloseParen)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		p.SetState(609)
		p.GetErrorHandler().Sync(p)

		if p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 60, p.GetParserRuleContext()) == 1 {
			{
				p.SetState(608)
				p.ErrorOperator()
			}

		} else if p.HasError() { // JIM
			goto errorExit
		}

	case antlr.ATNInvalidAltNumber:
		goto errorExit
	}
	p.GetParserRuleContext().SetStop(p.GetTokenStream().LT(-1))
	p.SetState(627)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_alt = p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 63, p.GetParserRuleContext())
	if p.HasError() {
		goto errorExit
	}
	for _alt != 2 && _alt != antlr.ATNInvalidAltNumber {
		if _alt == 1 {
			if p.GetParseListeners() != nil {
				p.TriggerExitRuleEvent()
			}
			_prevctx = localctx
			p.SetState(625)
			p.GetErrorHandler().Sync(p)
			if p.HasError() {
				goto errorExit
			}

			switch p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 62, p.GetParserRuleContext()) {
			case 1:
				localctx = NewExpressionAtomContext(p, _parentctx, _parentState)
				localctx.(*ExpressionAtomContext).left = _prevctx
				p.PushNewRecursionContext(localctx, _startState, FqlParserRULE_expressionAtom)
				p.SetState(613)

				if !(p.Precpred(p.GetParserRuleContext(), 10)) {
					p.SetError(antlr.NewFailedPredicateException(p, "p.Precpred(p.GetParserRuleContext(), 10)", ""))
					goto errorExit
				}
				{
					p.SetState(614)
					p.MultiplicativeOperator()
				}
				{
					p.SetState(615)

					var _x = p.expressionAtom(11)

					localctx.(*ExpressionAtomContext).right = _x
				}

			case 2:
				localctx = NewExpressionAtomContext(p, _parentctx, _parentState)
				localctx.(*ExpressionAtomContext).left = _prevctx
				p.PushNewRecursionContext(localctx, _startState, FqlParserRULE_expressionAtom)
				p.SetState(617)

				if !(p.Precpred(p.GetParserRuleContext(), 9)) {
					p.SetError(antlr.NewFailedPredicateException(p, "p.Precpred(p.GetParserRuleContext(), 9)", ""))
					goto errorExit
				}
				{
					p.SetState(618)
					p.AdditiveOperator()
				}
				{
					p.SetState(619)

					var _x = p.expressionAtom(10)

					localctx.(*ExpressionAtomContext).right = _x
				}

			case 3:
				localctx = NewExpressionAtomContext(p, _parentctx, _parentState)
				localctx.(*ExpressionAtomContext).left = _prevctx
				p.PushNewRecursionContext(localctx, _startState, FqlParserRULE_expressionAtom)
				p.SetState(621)

				if !(p.Precpred(p.GetParserRuleContext(), 8)) {
					p.SetError(antlr.NewFailedPredicateException(p, "p.Precpred(p.GetParserRuleContext(), 8)", ""))
					goto errorExit
				}
				{
					p.SetState(622)
					p.RegexpOperator()
				}
				{
					p.SetState(623)

					var _x = p.expressionAtom(9)

					localctx.(*ExpressionAtomContext).right = _x
				}

			case antlr.ATNInvalidAltNumber:
				goto errorExit
			}

		}
		p.SetState(629)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_alt = p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 63, p.GetParserRuleContext())
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
		p.SetState(630)

		var _lt = p.GetTokenStream().LT(1)

		localctx.(*ArrayOperatorContext).operator = _lt

		_la = p.GetTokenStream().LA(1)

		if !((int64(_la) & ^0x3f) == 0 && ((int64(1)<<_la)&108367866033602560) != 0) {
			var _ri = p.GetErrorHandler().RecoverInline(p)

			localctx.(*ArrayOperatorContext).operator = _ri
		} else {
			p.GetErrorHandler().ReportMatch(p)
			p.Consume()
		}
	}
	p.SetState(633)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetTokenStream().LA(1) {
	case FqlParserNot, FqlParserIn:
		{
			p.SetState(631)
			p.InOperator()
		}

	case FqlParserGt, FqlParserLt, FqlParserEq, FqlParserGte, FqlParserLte, FqlParserNeq:
		{
			p.SetState(632)
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
		p.SetState(635)
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
	p.SetState(638)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	if _la == FqlParserNot {
		{
			p.SetState(637)
			p.Match(FqlParserNot)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	}
	{
		p.SetState(640)
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
	p.SetState(643)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	if _la == FqlParserNot {
		{
			p.SetState(642)
			p.Match(FqlParserNot)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	}
	{
		p.SetState(645)
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
		p.SetState(647)
		_la = p.GetTokenStream().LA(1)

		if !((int64(_la) & ^0x3f) == 0 && ((int64(1)<<_la)&1152921504657178624) != 0) {
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
		p.SetState(649)
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
		p.SetState(651)
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
		p.SetState(653)
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
		p.SetState(655)
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
		p.SetState(657)
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
		p.SetState(659)
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
