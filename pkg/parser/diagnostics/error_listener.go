package diagnostics

import (
	"github.com/antlr4-go/antlr/v4"

	"github.com/MontFerret/ferret/v2/pkg/source"

	"github.com/MontFerret/ferret/v2/pkg/diagnostics"
)

type ErrorListener struct {
	*antlr.DiagnosticErrorListener
	src     *source.Source
	handler *ErrorHandler
	history *TokenHistory
}

func NewErrorListener(src *source.Source, handler *ErrorHandler, history *TokenHistory) antlr.ErrorListener {
	return &ErrorListener{
		DiagnosticErrorListener: antlr.NewDiagnosticErrorListener(false),
		src:                     src,
		handler:                 handler,
		history:                 history,
	}
}

func (d *ErrorListener) ReportAttemptingFullContext(recognizer antlr.Parser, dfa *antlr.DFA, startIndex, stopIndex int, conflictingAlts *antlr.BitSet, configs *antlr.ATNConfigSet) {
}

func (d *ErrorListener) ReportContextSensitivity(recognizer antlr.Parser, dfa *antlr.DFA, startIndex, stopIndex, prediction int, configs *antlr.ATNConfigSet) {
}

func (d *ErrorListener) ReportAmbiguity(recognizer antlr.Parser, dfa *antlr.DFA, startIndex, stopIndex int, exact bool, ambigAlts *antlr.BitSet, configs *antlr.ATNConfigSet) {
	if recognizer == nil || dfa == nil || d.DiagnosticErrorListener == nil {
		return
	}

	if ctx := recognizer.GetParserRuleContext(); ctx != nil {
		for _, rule := range recognizer.GetRuleInvocationStack(ctx) {
			if rule == "expressionAtom" {
				return
			}
		}
	}

	d.DiagnosticErrorListener.ReportAmbiguity(recognizer, dfa, startIndex, stopIndex, exact, ambigAlts, configs)
}

func (d *ErrorListener) SyntaxError(_ antlr.Recognizer, offendingSymbol interface{}, line, column int, msg string, e antlr.RecognitionException) {
	var offending antlr.Token

	// Get offending token
	if tok, ok := offendingSymbol.(antlr.Token); ok {
		offending = tok
	}

	if !d.handler.HasErrorOnLine(line) {
		if err := d.parseError(msg, offending); err != nil {
			d.handler.Add(err)
		}
	}
}

func (d *ErrorListener) parseError(msg string, offending antlr.Token) *diagnostics.Diagnostic {
	span := spanFromTokenSafe(offending, d.src)

	err := &diagnostics.Diagnostic{
		Kind:    SyntaxError,
		Source:  d.src,
		Message: "Syntax error: " + msg,
		Hint:    "Check your syntax. Did you forget to write something?",
		Spans: []diagnostics.ErrorSpan{
			{Span: span, Main: true},
		},
	}

	node := analyzedTokenNode(d.history, offending)

	AnalyzeSyntaxError(d.src, err, node)

	return err
}

func analyzedTokenNode(history *TokenHistory, offending antlr.Token) *TokenNode {
	var node *TokenNode
	if history != nil {
		node = cloneTokenChain(history.Last())
	}

	if node == nil {
		if offending == nil {
			return nil
		}

		return &TokenNode{token: offending}
	}

	if offending == nil || sameToken(node.token, offending) {
		return node
	}

	offendingNode := &TokenNode{token: offending, prev: node}
	node.next = offendingNode

	return node
}

func cloneTokenChain(head *TokenNode) *TokenNode {
	if head == nil {
		return nil
	}

	clonedHead := &TokenNode{token: head.token}
	clonedPrev := clonedHead

	for curr := head.Prev(); curr != nil; curr = curr.Prev() {
		cloned := &TokenNode{token: curr.token}
		clonedPrev.prev = cloned
		cloned.next = clonedPrev
		clonedPrev = cloned
	}

	return clonedHead
}

func sameToken(a, b antlr.Token) bool {
	if a == nil || b == nil {
		return a == b
	}

	return a.GetTokenType() == b.GetTokenType() &&
		a.GetStart() == b.GetStart() &&
		a.GetStop() == b.GetStop()
}
