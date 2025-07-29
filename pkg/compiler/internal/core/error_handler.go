package core

import (
	"fmt"

	"github.com/antlr4-go/antlr/v4"

	"github.com/MontFerret/ferret/pkg/file"
)

type ErrorHandler struct {
	src       *file.Source
	errors    []*CompilationError
	threshold int
}

func ParserLocation(ctx antlr.ParserRuleContext) *file.Location {
	start := ctx.GetStart()
	stop := ctx.GetStop()

	// Defensive: avoid nil dereference
	if start == nil || stop == nil {
		return file.NewLocation(0, 0, 0, 0)
	}

	return file.NewLocation(
		start.GetLine(),
		start.GetColumn()+1,
		start.GetStart(),
		stop.GetStop(),
	)
}

func NewErrorHandler(src *file.Source, threshold int) *ErrorHandler {
	if threshold <= 0 {
		threshold = 10
	}

	return &ErrorHandler{
		src:       src,
		errors:    make([]*CompilationError, 0),
		threshold: threshold,
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

	if err.Source == nil {
		err.Source = h.src
	}

	h.errors = append(h.errors, err)

	if len(h.errors) >= h.threshold {
		h.errors = append(h.errors, &CompilationError{
			Message: "Too many errors",
			Kind:    SemanticError,
			Hint:    "Too many errors encountered during compilation.",
		})

		panic(h.Unwrap())
	}
}

func (h *ErrorHandler) UnexpectedToken(ctx antlr.ParserRuleContext) {
	h.Add(&CompilationError{
		Message:  fmt.Sprintf("Unexpected token '%s'", ctx.GetText()),
		Source:   h.src,
		Location: ParserLocation(ctx),
		Kind:     SyntaxError,
	})
}

func (h *ErrorHandler) VariableNotUnique(ctx antlr.ParserRuleContext, name string) {
	// TODO: Add information where the variable was defined
	h.Add(&CompilationError{
		Message:  fmt.Sprintf("Variable '%s' is already defined", name),
		Source:   h.src,
		Location: ParserLocation(ctx),
		Kind:     NameError,
	})
}

func (h *ErrorHandler) VariableNotFound(ctx antlr.ParserRuleContext, name string) {
	h.Add(&CompilationError{
		Message:  fmt.Sprintf("Variable '%s' is not defined", name),
		Source:   h.src,
		Location: ParserLocation(ctx),
		Kind:     NameError,
	})
}
