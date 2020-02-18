// Code generated from pkg/parser/antlr/FqlParser.g4 by ANTLR 4.7.1. DO NOT EDIT.

package fql // FqlParser
import (
	"fmt"
	"reflect"
	"strconv"

	"github.com/antlr/antlr4/runtime/Go/antlr"
)

// Suppress unused import errors
var _ = fmt.Printf
var _ = reflect.Copy
var _ = strconv.Itoa

var parserATN = []uint16{
	3, 24715, 42794, 33075, 47597, 16764, 15335, 30598, 22884, 3, 66, 20, 4,
	2, 9, 2, 4, 3, 9, 3, 4, 4, 9, 4, 4, 5, 9, 5, 3, 2, 3, 2, 3, 3, 3, 3, 3,
	4, 3, 4, 3, 4, 3, 5, 3, 5, 3, 5, 2, 2, 6, 2, 4, 6, 8, 2, 2, 2, 15, 2, 10,
	3, 2, 2, 2, 4, 12, 3, 2, 2, 2, 6, 14, 3, 2, 2, 2, 8, 17, 3, 2, 2, 2, 10,
	11, 5, 8, 5, 2, 11, 3, 3, 2, 2, 2, 12, 13, 5, 6, 4, 2, 13, 5, 3, 2, 2,
	2, 14, 15, 7, 49, 2, 2, 15, 16, 7, 50, 2, 2, 16, 7, 3, 2, 2, 2, 17, 18,
	5, 4, 3, 2, 18, 9, 3, 2, 2, 2, 2,
}
var deserializer = antlr.NewATNDeserializer(nil)
var deserializedATN = deserializer.DeserializeFromUInt16(parserATN)

var literalNames = []string{
	"", "", "", "", "", "':'", "';'", "'.'", "','", "'['", "']'", "'('", "')'",
	"'{'", "'}'", "'>'", "'<'", "'=='", "'>='", "'<='", "'!='", "'*'", "'/'",
	"'%'", "'+'", "'-'", "'--'", "'++'", "", "", "", "'='", "'?'", "'!~'",
	"'=~'", "'FOR'", "'RETURN'", "'DISTINCT'", "'FILTER'", "'SORT'", "'LIMIT'",
	"'LET'", "'COLLECT'", "", "'NONE'", "'NULL'", "", "'USE'", "'AS'", "'INTO'",
	"'KEEP'", "'WITH'", "'COUNT'", "'ALL'", "'ANY'", "'AGGREGATE'", "'LIKE'",
	"", "'IN'", "'@'",
}
var symbolicNames = []string{
	"", "MultiLineComment", "SingleLineComment", "WhiteSpaces", "LineTerminator",
	"Colon", "SemiColon", "Dot", "Comma", "OpenBracket", "CloseBracket", "OpenParen",
	"CloseParen", "OpenBrace", "CloseBrace", "Gt", "Lt", "Eq", "Gte", "Lte",
	"Neq", "Multi", "Div", "Mod", "Plus", "Minus", "MinusMinus", "PlusPlus",
	"And", "Or", "Range", "Assign", "QuestionMark", "RegexNotMatch", "RegexMatch",
	"For", "Return", "Distinct", "Filter", "Sort", "Limit", "Let", "Collect",
	"SortDirection", "None", "Null", "BooleanLiteral", "Use", "As", "Into",
	"Keep", "With", "Count", "All", "Any", "Aggregate", "Like", "Not", "In",
	"Param", "Identifier", "StringLiteral", "IntegerLiteral", "FloatLiteral",
	"NamespaceSegment",
}

var ruleNames = []string{
	"program", "useExpression", "useExpressionAs", "body",
}
var decisionToDFA = make([]*antlr.DFA, len(deserializedATN.DecisionToState))

func init() {
	for index, ds := range deserializedATN.DecisionToState {
		decisionToDFA[index] = antlr.NewDFA(ds, index)
	}
}

type FqlParser struct {
	*antlr.BaseParser
}

func NewFqlParser(input antlr.TokenStream) *FqlParser {
	this := new(FqlParser)

	this.BaseParser = antlr.NewBaseParser(input)

	this.Interpreter = antlr.NewParserATNSimulator(this, deserializedATN, decisionToDFA, antlr.NewPredictionContextCache())
	this.RuleNames = ruleNames
	this.LiteralNames = literalNames
	this.SymbolicNames = symbolicNames
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
	FqlParserDistinct          = 37
	FqlParserFilter            = 38
	FqlParserSort              = 39
	FqlParserLimit             = 40
	FqlParserLet               = 41
	FqlParserCollect           = 42
	FqlParserSortDirection     = 43
	FqlParserNone              = 44
	FqlParserNull              = 45
	FqlParserBooleanLiteral    = 46
	FqlParserUse               = 47
	FqlParserAs                = 48
	FqlParserInto              = 49
	FqlParserKeep              = 50
	FqlParserWith              = 51
	FqlParserCount             = 52
	FqlParserAll               = 53
	FqlParserAny               = 54
	FqlParserAggregate         = 55
	FqlParserLike              = 56
	FqlParserNot               = 57
	FqlParserIn                = 58
	FqlParserParam             = 59
	FqlParserIdentifier        = 60
	FqlParserStringLiteral     = 61
	FqlParserIntegerLiteral    = 62
	FqlParserFloatLiteral      = 63
	FqlParserNamespaceSegment  = 64
)

// FqlParser rules.
const (
	FqlParserRULE_program         = 0
	FqlParserRULE_useExpression   = 1
	FqlParserRULE_useExpressionAs = 2
	FqlParserRULE_body            = 3
)

// IProgramContext is an interface to support dynamic dispatch.
type IProgramContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsProgramContext differentiates from other interfaces.
	IsProgramContext()
}

type ProgramContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyProgramContext() *ProgramContext {
	var p = new(ProgramContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = FqlParserRULE_program
	return p
}

func (*ProgramContext) IsProgramContext() {}

func NewProgramContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ProgramContext {
	var p = new(ProgramContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = FqlParserRULE_program

	return p
}

func (s *ProgramContext) GetParser() antlr.Parser { return s.parser }

func (s *ProgramContext) Body() IBodyContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IBodyContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IBodyContext)
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

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(8)
		p.Body()
	}

	return localctx
}

// IUseExpressionContext is an interface to support dynamic dispatch.
type IUseExpressionContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsUseExpressionContext differentiates from other interfaces.
	IsUseExpressionContext()
}

type UseExpressionContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyUseExpressionContext() *UseExpressionContext {
	var p = new(UseExpressionContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = FqlParserRULE_useExpression
	return p
}

func (*UseExpressionContext) IsUseExpressionContext() {}

func NewUseExpressionContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *UseExpressionContext {
	var p = new(UseExpressionContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = FqlParserRULE_useExpression

	return p
}

func (s *UseExpressionContext) GetParser() antlr.Parser { return s.parser }

func (s *UseExpressionContext) UseExpressionAs() IUseExpressionAsContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IUseExpressionAsContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IUseExpressionAsContext)
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
	p.EnterRule(localctx, 2, FqlParserRULE_useExpression)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(10)
		p.UseExpressionAs()
	}

	return localctx
}

// IUseExpressionAsContext is an interface to support dynamic dispatch.
type IUseExpressionAsContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsUseExpressionAsContext differentiates from other interfaces.
	IsUseExpressionAsContext()
}

type UseExpressionAsContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyUseExpressionAsContext() *UseExpressionAsContext {
	var p = new(UseExpressionAsContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = FqlParserRULE_useExpressionAs
	return p
}

func (*UseExpressionAsContext) IsUseExpressionAsContext() {}

func NewUseExpressionAsContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *UseExpressionAsContext {
	var p = new(UseExpressionAsContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = FqlParserRULE_useExpressionAs

	return p
}

func (s *UseExpressionAsContext) GetParser() antlr.Parser { return s.parser }

func (s *UseExpressionAsContext) Use() antlr.TerminalNode {
	return s.GetToken(FqlParserUse, 0)
}

func (s *UseExpressionAsContext) As() antlr.TerminalNode {
	return s.GetToken(FqlParserAs, 0)
}

func (s *UseExpressionAsContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *UseExpressionAsContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *UseExpressionAsContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(FqlParserListener); ok {
		listenerT.EnterUseExpressionAs(s)
	}
}

func (s *UseExpressionAsContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(FqlParserListener); ok {
		listenerT.ExitUseExpressionAs(s)
	}
}

func (s *UseExpressionAsContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case FqlParserVisitor:
		return t.VisitUseExpressionAs(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *FqlParser) UseExpressionAs() (localctx IUseExpressionAsContext) {
	localctx = NewUseExpressionAsContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 4, FqlParserRULE_useExpressionAs)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(12)
		p.Match(FqlParserUse)
	}
	{
		p.SetState(13)
		p.Match(FqlParserAs)
	}

	return localctx
}

// IBodyContext is an interface to support dynamic dispatch.
type IBodyContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsBodyContext differentiates from other interfaces.
	IsBodyContext()
}

type BodyContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyBodyContext() *BodyContext {
	var p = new(BodyContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = FqlParserRULE_body
	return p
}

func (*BodyContext) IsBodyContext() {}

func NewBodyContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *BodyContext {
	var p = new(BodyContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = FqlParserRULE_body

	return p
}

func (s *BodyContext) GetParser() antlr.Parser { return s.parser }

func (s *BodyContext) UseExpression() IUseExpressionContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IUseExpressionContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IUseExpressionContext)
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
	p.EnterRule(localctx, 6, FqlParserRULE_body)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(15)
		p.UseExpression()
	}

	return localctx
}
