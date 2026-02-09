package diagnostics

import (
	"fmt"

	"github.com/antlr4-go/antlr/v4"

	"github.com/MontFerret/ferret/pkg/diagnostics"

	"github.com/MontFerret/ferret/pkg/file"
)

type ErrorHandler struct {
	src             *file.Source
	errors          *CompilationErrorSet
	linesWithErrors map[int]bool
	threshold       int
}

func NewErrorHandler(src *file.Source, threshold int) *ErrorHandler {
	if threshold <= 0 {
		threshold = 10
	}

	return &ErrorHandler{
		src:             src,
		errors:          NewCompilationErrorSet(5),
		linesWithErrors: make(map[int]bool),
		threshold:       threshold,
	}
}

func (h *ErrorHandler) Errors() *CompilationErrorSet {
	return h.errors
}

func (h *ErrorHandler) HasErrors() bool {
	return h.errors.Size() > 0
}

func (h *ErrorHandler) Unwrap() error {
	if h.errors.Size() == 0 {
		return nil
	}

	if h.errors.Size() == 1 {
		return h.errors.First()
	}

	return h.errors
}

func (h *ErrorHandler) Add(err *CompilationError) {
	if err == nil {
		return
	}

	// If the number of errors exceeds the threshold, we stop adding new errors
	if h.errors.Size() > h.threshold {
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

	h.errors.Add(err)

	if h.errors.Size() == h.threshold {
		h.errors.Add(&CompilationError{
			Diagnostic: &diagnostics.Diagnostic{
				Message: "Too many errors",
				Kind:    SemanticError,
				Hint:    "Too many errors encountered during compilation.",
			},
		})
	}
}

func (h *ErrorHandler) Create(kind diagnostics.Kind, ctx antlr.ParserRuleContext, msg string) *CompilationError {
	return &CompilationError{
		Diagnostic: &diagnostics.Diagnostic{
			Kind:    kind,
			Source:  h.src,
			Spans:   []diagnostics.ErrorSpan{diagnostics.NewMainErrorSpan(SpanFromRuleContext(ctx), "")},
			Message: msg,
		},
	}
}

func (h *ErrorHandler) HasErrorOnLine(line int) bool {
	return h.linesWithErrors[line]
}

func (h *ErrorHandler) VariableNotUnique(ctx antlr.ParserRuleContext, name string) {
	// TODO: Add information where the variable was defined
	h.Add(&CompilationError{
		Diagnostic: &diagnostics.Diagnostic{
			Message: fmt.Sprintf("Variable '%s' is already defined", name),
			Source:  h.src,
			Spans:   []diagnostics.ErrorSpan{diagnostics.NewMainErrorSpan(SpanFromRuleContext(ctx), "")},
			Kind:    NameError,
		},
	})
}

func (h *ErrorHandler) VariableNotFound(token antlr.Token, name string) {
	h.Add(&CompilationError{
		Diagnostic: &diagnostics.Diagnostic{
			Message: fmt.Sprintf("Variable '%s' is not defined", name),
			Source:  h.src,
			Spans:   []diagnostics.ErrorSpan{diagnostics.NewMainErrorSpan(SpanFromToken(token), "undefined variable")},
			Kind:    NameError,
			Hint:    "Did you forget to declare it?",
		},
	})
}

func (h *ErrorHandler) MissingReturnValue(ctx antlr.ParserRuleContext) {
	h.Add(&CompilationError{
		Diagnostic: &diagnostics.Diagnostic{
			Message: fmt.Sprintf("Expected expression after '%s'", ctx.GetText()),
			Hint:    "Did you forget to provide a value to return?",
			Source:  h.src,
			Spans: []diagnostics.ErrorSpan{
				diagnostics.NewMainErrorSpan(SpanFromRuleContext(ctx), "missing return value")},
			Kind: SyntaxError,
		},
	})
}

func (h *ErrorHandler) InvalidRegexExpression(ctx antlr.ParserRuleContext, expression string) {
	h.Add(&CompilationError{
		Diagnostic: &diagnostics.Diagnostic{
			Message: fmt.Sprintf("Invalid regular expression: %s", expression),
			Hint:    "Check the syntax of the regular expression.",
			Source:  h.src,
			Spans: []diagnostics.ErrorSpan{
				diagnostics.NewMainErrorSpan(SpanFromRuleContext(ctx), "invalid regex")},
			Kind: SyntaxError,
		},
	})
}

func (h *ErrorHandler) InvalidToken(token antlr.Token) {
	h.Add(&CompilationError{
		Diagnostic: &diagnostics.Diagnostic{
			Message: fmt.Sprintf("Invalid token: %s", token),
			Hint:    "Check the syntax of the literal.",
			Source:  h.src,
			Spans: []diagnostics.ErrorSpan{
				diagnostics.NewMainErrorSpan(SpanFromToken(token), "invalid token")},
			Kind: SyntaxError,
		},
	})
}
