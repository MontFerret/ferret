package diagnostics

import (
	"errors"
	"fmt"
	"strings"

	"github.com/MontFerret/ferret/v2/pkg/bytecode"
	"github.com/MontFerret/ferret/v2/pkg/diagnostics"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
	"github.com/MontFerret/ferret/v2/pkg/source"
	"github.com/MontFerret/ferret/v2/pkg/vm/internal/frame"
)

const vmStackNotePrefix = "VM stack: "

type (
	// RuntimeError represents a VM execution error with source context.
	RuntimeError struct {
		*diagnostics.Diagnostic
	}

	// RuntimeErrorSet is a specialized diagnostics.Diagnostics type for RuntimeError.
	RuntimeErrorSet struct {
		diagnostics.Diagnostics[*RuntimeError]
	}
)

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

func NewRuntimeErrorSet(size int) *RuntimeErrorSet {
	return &RuntimeErrorSet{
		Diagnostics: *diagnostics.NewDiagnostics[*RuntimeError](size),
	}
}

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
		if wpErrorSet.Size() == 1 {
			wer := wpErrorSet.First()
			return ToRuntimeError(program, wer.PC+1, nil, wer.Cause)
		}

		errs := NewRuntimeErrorSet(5)

		for _, wer := range wpErrorSet.Errors() {
			// warmup PCs are zero-based instruction indices (no pre-increment),
			// while ToRuntimeError expects a post-increment pc (see pc-1 usage)
			errs.Add(ToRuntimeError(program, wer.PC+1, nil, wer.Cause))
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
			Spans:   buildSpans(program, callStack, SpanAt(program, pc-1), ""),
			Cause:   cause,
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
	mainSpan := runtimeErrorSpan(program, pc, err)
	argPos, hasArg, argCause := runtime.InvalidArgumentDetails(err)

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
	case hasArg && (errors.Is(argCause, runtime.ErrInvalidType) || errors.Is(argCause, runtime.ErrInvalidArgumentType)):
		index := argPos + 1
		kind = diagnostics.TypeError
		message = fmt.Sprintf("Invalid argument %d type", index)
		label = fmt.Sprintf("argument %d type mismatch", index)
		hint = fmt.Sprintf("Ensure argument %d matches the expected type", index)
		cause = argCause

		msg, cs := diagnostics.Unwrap(argCause)

		if msg != nil && cs != nil {
			cause = cs
		}
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

		if hasArg {
			index := argPos + 1
			message = fmt.Sprintf("Invalid argument %d", index)
			label = fmt.Sprintf("invalid argument %d", index)
			hint = fmt.Sprintf("Check argument %d", index)
			cause = argCause

			_, cs := diagnostics.Unwrap(cause)
			if cs != nil {
				cause = cs
			}
		} else {
			msg, cs := diagnostics.Unwrap(cause)

			if msg != nil && cs != nil {
				message = diagnostics.FormatMessage(msg.Error())
				cause = cs
			}
		}
	case errors.Is(err, ErrMissedParam):
		kind = UnresolvedSymbol
		message = "Missing parameter"
		label = "missing parameter"
		hint = "Provide all required parameters"
		cause = err
	case errors.Is(err, runtime.ErrInvalidArgumentNumber):
		kind = ArityError
		message = "Wrong number of arguments"
		label = "wrong number of arguments"
		hint = "Check the function signature for the expected argument count"
		cause = err

		_, detail := diagnostics.Unwrap(err)
		if detail != nil {
			s := detail.Error()
			if len(s) > 0 {
				note = strings.ToUpper(s[:1]) + s[1:]
			}
		}
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
			Spans:   buildSpans(program, callStack, mainSpan, label),
			Cause:   cause,
		},
	}
}

func runtimeErrorSpan(program *bytecode.Program, pc int, err error) source.Span {
	mainSpan := SpanAt(program, pc-1)

	argPos, ok, _ := runtime.InvalidArgumentDetails(err)
	if !ok {
		return mainSpan
	}

	argSpan := CallArgumentSpanAt(program, pc-1, argPos)
	if argSpan.Start < 0 || argSpan.End < 0 {
		return mainSpan
	}

	return argSpan
}
