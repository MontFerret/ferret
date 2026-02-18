package vm

import (
	"errors"

	"github.com/MontFerret/ferret/v2/pkg/bytecode"
	"github.com/MontFerret/ferret/v2/pkg/diagnostics"
	"github.com/MontFerret/ferret/v2/pkg/file"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

func toRuntimeError(program *bytecode.Program, pc int, err error) *RuntimeError {
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
		message = "Invalid function name"
		label = "invalid function name"
		hint = "Ensure the function name is valid and does not contain illegal characters"
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
			Note:    note,
			Source:  program.Source,
			Spans:   []diagnostics.ErrorSpan{diagnostics.NewMainErrorSpan(spanAt(program, pc-1), label)},
			Cause:   cause,
		},
	}
}

func spanAt(program *bytecode.Program, pc int) file.Span {
	if program == nil {
		return file.Span{Start: -1, End: -1}
	}

	if pc < 0 || pc >= len(program.Metadata.DebugSpans) {
		return file.Span{Start: -1, End: -1}
	}

	return program.Metadata.DebugSpans[pc]
}
