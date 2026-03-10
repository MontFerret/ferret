package diagnostic

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/MontFerret/ferret/v2/pkg/bytecode"
	"github.com/MontFerret/ferret/v2/pkg/diagnostics"
	"github.com/MontFerret/ferret/v2/pkg/file"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
	"github.com/MontFerret/ferret/v2/pkg/vm/internal/frame"
)

const vmStackNotePrefix = "VM stack: "

func WrapRuntimeError(program *bytecode.Program, pc int, callStack []frame.TraceEntry, err error) error {
	if err == nil {
		return nil
	}

	var runtimeError *RuntimeError

	if errors.As(err, &runtimeError) {
		return attachCallStack(runtimeError, program, callStack)
	}

	var wpErrorSet *WarmupErrorSet

	if errors.As(err, &wpErrorSet) {
		if len(wpErrorSet.Errors) == 1 {
			wer := wpErrorSet.Errors[0]
			return ToRuntimeError(program, wer.PC+1, nil, wer.Err)
		}

		errs := NewRuntimeErrorSet(5)

		for _, wer := range wpErrorSet.Errors {
			// warmup PCs are zero-based instruction indices (no pre-increment),
			// while ToRuntimeError expects a post-increment pc (see pc-1 usage)
			errs.Add(ToRuntimeError(program, wer.PC+1, nil, wer.Err))
		}

		return errs
	}

	return ToRuntimeError(program, pc, callStack, err)
}

func RuntimeErrorFromPanic(program *bytecode.Program, pc int, callStack []frame.TraceEntry, r any) error {
	message := "unexpected runtime panic"
	cause := fmt.Errorf("panic: %v", r)

	if err, ok := r.(error); ok {
		cause = err
	}

	return &RuntimeError{
		Diagnostic: &diagnostics.Diagnostic{
			Kind:    diagnostics.UnexpectedError,
			Message: fmt.Sprintf("%s. %s", message, cause.Error()),
			Source:  program.Source,
			Note:    stackNote(callStack),
			Spans:   buildSpans(program, pc, callStack, ""),
			Cause:   cause,
		},
	}
}

func NewRuntimeError(
	program *bytecode.Program,
	pc int,
	kind diagnostics.Kind,
	message string,
	label string,
	hint string,
	note string,
) *RuntimeError {
	return &RuntimeError{
		Diagnostic: &diagnostics.Diagnostic{
			Kind:    kind,
			Message: message,
			Hint:    hint,
			Note:    note,
			Source:  program.Source,
			Spans:   []diagnostics.ErrorSpan{diagnostics.NewMainErrorSpan(SpanAt(program, pc-1), label)},
		},
	}
}

func ToRuntimeError(program *bytecode.Program, pc int, callStack []frame.TraceEntry, err error) *RuntimeError {
	if err == nil {
		return nil
	}

	var runtimeError *RuntimeError

	if errors.As(err, &runtimeError) {
		return runtimeError
	}

	kind := diagnostics.Unknown
	message := ""
	label := ""
	hint := ""
	note := ""
	var cause error
	var memberErr *MemberAccessError
	var invariantErr *InvariantError

	switch {
	case errors.Is(err, ErrDivisionByZero):
		kind = DivideByZero
		message = "Division by zero"
		label = "attempt to divide by zero"
		hint = "Ensure the denominator is non-zero before division"
		note = "Add a conditional check before dividing"
	case errors.Is(err, ErrModuloByZero):
		kind = ModuloByZero
		message = "Modulo by zero"
		label = "attempt to take modulo by zero"
		hint = "Ensure the divisor is non-zero before modulo"
		note = "Add a conditional check before modulo"
	case errors.As(err, &memberErr):
		kind = diagnostics.TypeError
		message = memberErr.Error()
		if message != "" {
			message = strings.ToUpper(message[:1]) + message[1:]
		}
		label = memberErr.Label()
		hint = memberErr.Hint()
	case errors.Is(err, runtime.ErrInvalidType):
		kind = diagnostics.TypeError
		message = "Invalid type"
		label = "type mismatch"
		hint = "Ensure the value has the expected type"
		cause = err

		msg, cs := diagnostics.Unwrap(err)

		if msg != nil && cs != nil {
			message = diagnostics.FormatMessage(msg.Error())
			cause = cs
		}
	case errors.Is(err, runtime.ErrInvalidArgumentType):
		kind = diagnostics.TypeError
		message = "Invalid argument type"
		hint = "Ensure the argument types match the function signature"

		msg, cs := diagnostics.Unwrap(err)

		if msg != nil && cs != nil {
			message = diagnostics.FormatMessage(msg.Error())
			cause = cs
		}
	case errors.Is(err, runtime.ErrInvalidArgumentNumber):
		kind = ArityError
		message = "Invalid number of arguments"
		hint = "Check the function signature for the expected argument count"
		cause = err
	case errors.Is(err, runtime.ErrInvalidArgument):
		kind = ArityError
		message = "Invalid argument"
		hint = "Check the function arguments"
		cause = err

		msg, cs := diagnostics.Unwrap(err)

		if msg != nil && cs != nil {
			message = diagnostics.FormatMessage(msg.Error())
			cause = cs
		}
	case errors.Is(err, ErrMissedParam):
		kind = UnresolvedSymbol
		message = "Missing parameter"
		label = "missing parameter"
		hint = "Provide all required parameters"
		cause = err
	case errors.Is(err, ErrUnresolvedFunction):
		kind = UnresolvedSymbol
		message = "Unresolved function"
		label = "unresolved function"
		hint = "Ensure the function is registered and accessible in the current context"
		note = "Add the function to the registry if it's missing"
	case errors.Is(err, ErrInvalidFunctionName):
		kind = UnresolvedSymbol
		message = "Invalid function name"
		label = "invalid function name"
		hint = "Ensure the function name is valid and does not contain illegal characters"
	case errors.As(err, &invariantErr):
		kind = diagnostics.UnexpectedError
		message = "VM invariant violation"
		if invariantErr.Message != "" {
			message = invariantErr.Message
		}

		label = "internal invariant violated"
		hint = "This indicates an internal VM bug; please report it with the query and stack context"
		cause = invariantErr.Cause
	default:
		kind = UncaughtError
		msg, cs := diagnostics.Unwrap(err)

		if msg != nil && cs != nil {
			message = diagnostics.FormatMessage(msg.Error())
			cause = cs
		} else {
			message = "Runtime Error"
			cause = err
		}
	}

	return &RuntimeError{
		Diagnostic: &diagnostics.Diagnostic{
			Kind:    kind,
			Message: message,
			Hint:    hint,
			Note:    appendStackNote(note, callStack),
			Source:  program.Source,
			Spans:   buildSpans(program, pc, callStack, label),
			Cause:   cause,
		},
	}
}

func attachCallStack(err *RuntimeError, program *bytecode.Program, callStack []frame.TraceEntry) *RuntimeError {
	if err == nil || len(callStack) == 0 || err.Diagnostic == nil {
		return err
	}

	if !hasStackTraceLabel(err.Spans) {
		err.Spans = append(callStackSpans(program, callStack), err.Spans...)
	}

	err.Note = appendStackNote(err.Note, callStack)

	return err
}

func hasStackTraceLabel(spans []diagnostics.ErrorSpan) bool {
	for _, span := range spans {
		if strings.HasPrefix(span.Label, "called from ") {
			return true
		}
	}

	return false
}

func callStackSpans(program *bytecode.Program, callStack []frame.TraceEntry) []diagnostics.ErrorSpan {
	if len(callStack) == 0 {
		return nil
	}

	spans := make([]diagnostics.ErrorSpan, 0, len(callStack))

	for i, entry := range callStack {
		span := SpanAt(program, entry.CallSitePC)
		if span.Start < 0 || span.End < 0 {
			continue
		}

		label := fmt.Sprintf("called from (#%d)", i+1)
		if name := strings.TrimSpace(entry.FnName); name != "" {
			label = fmt.Sprintf("called from %s (#%d)", name, i+1)
		}

		spans = append(spans, diagnostics.NewSecondaryErrorSpan(span, label))
	}

	return spans
}

func buildSpans(program *bytecode.Program, pc int, callStack []frame.TraceEntry, label string) []diagnostics.ErrorSpan {
	spans := callStackSpans(program, callStack)
	spans = append(spans, diagnostics.NewMainErrorSpan(SpanAt(program, pc-1), label))

	return spans
}

func stackNote(callStack []frame.TraceEntry) string {
	if len(callStack) == 0 {
		return ""
	}

	names := make([]string, 0, len(callStack))

	// callStack is nearest -> farthest; render as outer -> ... -> inner
	for i := len(callStack) - 1; i >= 0; i-- {
		entry := callStack[i]
		name := strings.TrimSpace(entry.FnName)
		if name == "" {
			if entry.FnID < 0 {
				continue
			}

			name = fmt.Sprintf("#%d", entry.FnID)
		}

		names = append(names, name)
	}

	if len(names) == 0 {
		return ""
	}

	return vmStackNotePrefix + strings.Join(names, " -> ")
}

func appendStackNote(note string, callStack []frame.TraceEntry) string {
	stack := stackNote(callStack)
	if stack == "" {
		return note
	}

	if strings.Contains(note, vmStackNotePrefix) || strings.Contains(note, stack) {
		return note
	}

	if note == "" {
		return stack
	}

	return note + "\n" + stack
}

func SpanAt(program *bytecode.Program, pc int) file.Span {
	if program == nil {
		return file.Span{Start: -1, End: -1}
	}

	if pc < 0 || pc >= len(program.Metadata.DebugSpans) {
		return file.Span{Start: -1, End: -1}
	}

	return program.Metadata.DebugSpans[pc]
}

func CheckDivisionByZero(
	ctx context.Context,
	program *bytecode.Program,
	pc int,
	left runtime.Value,
	right runtime.Value,
) error {
	l := runtime.ToNumberOnly(ctx, left)
	if _, ok := l.(runtime.Int); !ok {
		return nil
	}

	r := runtime.ToNumberOnly(ctx, right)
	if rv, ok := r.(runtime.Int); ok && rv == 0 {
		return NewRuntimeError(
			program,
			pc,
			DivideByZero,
			"Division by zero",
			"attempt to divide by zero",
			"Ensure the denominator is non-zero before division",
			"Add a conditional check before dividing",
		)
	}

	return nil
}

func CheckModuloByZero(
	ctx context.Context,
	program *bytecode.Program,
	pc int,
	right runtime.Value,
) error {
	rv, err := runtime.ToInt(ctx, right)
	if err != nil {
		// Keep modulo diagnostics type-safe for invalid string inputs like "x".
		if _, ok := right.(runtime.String); ok {
			return runtime.TypeErrorOf(right, runtime.TypeInt)
		}

		return err
	}

	if rv == 0 {
		return NewRuntimeError(
			program,
			pc,
			ModuloByZero,
			"Modulo by zero",
			"attempt to take modulo by zero",
			"Ensure the divisor is non-zero before modulo",
			"Add a conditional check before modulo",
		)
	}

	return nil
}
