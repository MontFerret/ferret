package diagnostics

import (
	"fmt"

	"github.com/antlr4-go/antlr/v4"

	"github.com/MontFerret/ferret/pkg/file"
)

type ErrorHandler struct {
	src             *file.Source
	errors          []*CompilationError
	linesWithErrors map[int]bool
	threshold       int
}

func NewErrorHandler(src *file.Source, threshold int) *ErrorHandler {
	if threshold <= 0 {
		threshold = 10
	}

	return &ErrorHandler{
		src:             src,
		errors:          make([]*CompilationError, 0),
		linesWithErrors: make(map[int]bool),
		threshold:       threshold,
	}
}

func (h *ErrorHandler) Errors() []*CompilationError {
	return h.errors
}

func (h *ErrorHandler) HasErrors() bool {
	return len(h.errors) > 0
}

func (h *ErrorHandler) Unwrap() error {
	if len(h.errors) == 0 {
		return nil
	}

	if len(h.errors) == 1 {
		return h.errors[0]
	}

	return NewMultiCompilationError(h.errors)
}

func (h *ErrorHandler) Add(err *CompilationError) {
	if err == nil {
		return
	}

	// If the number of errors exceeds the threshold, we stop adding new errors
	if len(h.errors) > h.threshold {
		return
	}

	if err.Source == nil {
		err.Source = h.src
	}

	for _, span := range err.Spans {
		if err.Source != nil {
			line, _ := err.Source.LocationAt(span.Span)
			h.linesWithErrors[line] = true
		}
	}

	h.errors = append(h.errors, err)

	if len(h.errors) == h.threshold {
		h.errors = append(h.errors, &CompilationError{
			Message: "Too many errors",
			Kind:    SemanticError,
			Hint:    "Too many errors encountered during compilation.",
		})
	}
}

func (h *ErrorHandler) HasErrorOnLine(line int) bool {
	return h.linesWithErrors[line]
}

func (h *ErrorHandler) VariableNotUnique(ctx antlr.ParserRuleContext, name string) {
	// TODO: Add information where the variable was defined
	h.Add(&CompilationError{
		Message: fmt.Sprintf("Variable '%s' is already defined", name),
		Source:  h.src,
		Spans:   []ErrorSpan{NewMainErrorSpan(SpanFromRuleContext(ctx), "")},
		Kind:    NameError,
	})
}

func (h *ErrorHandler) VariableNotFound(token antlr.Token, name string) {
	h.Add(&CompilationError{
		Message: fmt.Sprintf("Variable '%s' is not defined", name),
		Source:  h.src,
		Spans:   []ErrorSpan{NewMainErrorSpan(SpanFromToken(token), "")},
		Kind:    NameError,
	})
}

func (h *ErrorHandler) MissingReturnValue(ctx antlr.ParserRuleContext) {
	//span := spanFromTokenSafe(offending.Token(), src)
	//err.Message = fmt.Sprintf("Expected expression after '%s'", offending)
	//err.Hint = "Did you forget to provide a value to return?"
	//err.Spans = []ErrorSpan{
	//	NewMainErrorSpan(span, "missing return value"),
	//}

	h.Add(&CompilationError{
		Message: fmt.Sprintf("Expected expression after '%s'", ctx.GetText()),
		Hint:    "Did you forget to provide a value to return?",
		Source:  h.src,
		Spans: []ErrorSpan{
			NewMainErrorSpan(SpanFromRuleContext(ctx), "missing return value")},
		Kind: SyntaxError,
	})
}

func (h *ErrorHandler) InvalidRegexExpression(ctx antlr.ParserRuleContext, expression string) {
	h.Add(&CompilationError{
		Message: fmt.Sprintf("Invalid regular expression: %s", expression),
		Hint:    "Check the syntax of the regular expression.",
		Source:  h.src,
		Spans: []ErrorSpan{
			NewMainErrorSpan(SpanFromRuleContext(ctx), "invalid regex")},
		Kind: SyntaxError,
	})
}
