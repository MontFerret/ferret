package diagnostics

import (
	"errors"
	"fmt"

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
	cause error,
) *RuntimeError {
	return newRuntimeErrorWithSpec(program, nil, runtimeDiagnosticSpec{
		Cause:   cause,
		Kind:    kind,
		Span:    SpanAt(program, pc-1),
		Hint:    hint,
		Label:   label,
		Message: message,
		Note:    note,
	})
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

		errs := NewRuntimeErrorSet(wpErrorSet.Size())

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
	cause := fmt.Errorf("panic: %v", r)

	if err, ok := r.(error); ok {
		cause = err
	}

	return newRuntimeErrorWithSpec(program, callStack, runtimeDiagnosticSpec{
		Cause:   cause,
		Kind:    diagnostics.UnexpectedError,
		Span:    SpanAt(program, pc-1),
		Hint:    "This indicates an internal VM bug; please report it with the query and stack context",
		Label:   "panic occurred during VM execution",
		Message: "unexpected runtime panic",
		Note:    panicValueNote(r),
	})
}

func ToRuntimeError(program *bytecode.Program, pc int, callStack []frame.TraceEntry, err error) *RuntimeError {
	if err == nil {
		return nil
	}

	var runtimeError *RuntimeError

	if errors.As(err, &runtimeError) {
		return runtimeError
	}

	spec := runtimeDiagnosticSpec{
		Cause:   err,
		Kind:    UncaughtError,
		Span:    runtimeErrorSpan(program, pc, err),
		Message: "runtime error",
	}

	var memberErr *MemberAccessError
	var invariantErr *InvariantError
	argPos, hasArg, argCause := runtime.InvalidArgumentDetails(err)

	if hasArg {
		for {
			_, ok, nestedCause := runtime.InvalidArgumentDetails(argCause)
			if !ok {
				break
			}

			argCause = nestedCause
		}
	}

	switch {
	case errors.Is(err, ErrDivisionByZero):
		spec.Kind = DivideByZero
		spec.Message = "division by zero"
		spec.Label = "denominator evaluates to zero"
		spec.Hint = "Ensure the denominator is non-zero before division"
		spec.Cause = ErrDivisionByZero
	case errors.Is(err, ErrModuloByZero):
		spec.Kind = ModuloByZero
		spec.Message = "modulo by zero"
		spec.Label = "divisor evaluates to zero"
		spec.Hint = "Ensure the divisor is non-zero before modulo"
		spec.Cause = ErrModuloByZero
	case errors.As(err, &memberErr):
		spec.Kind = diagnostics.TypeError
		spec.Message = "invalid type"
		spec.Label = memberErr.Label()
		spec.Note = memberErr.Note()
		spec.Hint = memberErr.Hint()
		spec.Cause = runtime.ErrInvalidType
	case hasArg && (errors.Is(argCause, runtime.ErrInvalidType) || errors.Is(argCause, runtime.ErrInvalidArgumentType)):
		index := argPos + 1
		detail, cause := unwrapRuntimeDetail(argCause)

		spec.Kind = diagnostics.TypeError
		spec.Message = "invalid argument type"
		spec.Label = fmt.Sprintf("argument %d has incompatible type", index)
		spec.Note = argumentDetailNote(index, detail)
		spec.Cause = cause
	case errors.Is(err, runtime.ErrInvalidType):
		detail, cause := unwrapRuntimeDetail(err)

		spec.Kind = diagnostics.TypeError
		spec.Message = "invalid type"
		spec.Label = "value has incompatible type"
		spec.Note = detailNote(detail)
		spec.Cause = cause
	case errors.Is(err, runtime.ErrInvalidArgumentType):
		detail, cause := unwrapRuntimeDetail(err)

		spec.Kind = diagnostics.TypeError
		spec.Message = "invalid argument type"
		spec.Label = "argument has incompatible type"
		spec.Note = detailNote(detail)
		spec.Cause = cause
	case errors.Is(err, runtime.ErrInvalidArgumentNumber):
		detail, cause := unwrapRuntimeDetail(err)

		spec.Kind = ArityError
		spec.Message = "invalid number of arguments"
		spec.Note = detailNote(detail)
		spec.Label = arityLabel(spec.Note)
		spec.Cause = cause
	case errors.Is(err, runtime.ErrInvalidArgument):
		spec.Kind = InvalidArgument
		spec.Message = "invalid argument"

		if hasArg {
			index := argPos + 1
			detail, cause := unwrapRuntimeDetail(argCause)

			spec.Label = fmt.Sprintf("argument %d is invalid", index)
			spec.Note = argumentDetailNote(index, detail)
			spec.Cause = cause
		} else {
			detail, cause := unwrapRuntimeDetail(err)

			spec.Label = "invalid argument"
			spec.Note = detailNote(detail)
			spec.Cause = cause
		}
	case errors.Is(err, ErrMissedParam):
		detail, cause := unwrapRuntimeDetail(err)
		name := detailNote(detail)

		spec.Kind = UnresolvedSymbol
		spec.Message = "missing parameter"
		spec.Cause = cause

		if name == "" {
			spec.Label = "parameter was not provided"
			spec.Hint = "Provide all required parameters before executing this query"
			break
		}

		spec.Label = fmt.Sprintf("parameter '%s' was not provided", name)
		spec.Note = fmt.Sprintf("this query requires parameter '%s'", name)
		spec.Hint = fmt.Sprintf("Provide a value for %s before executing this query", name)
	case errors.Is(err, ErrUnresolvedFunction):
		spec.Kind = UnresolvedSymbol
		spec.Message = "unresolved function"
		spec.Label = "function is not registered"
		spec.Note = "function could not be resolved in the current registry"
		spec.Hint = "Register the function before executing this query"
		spec.Cause = ErrUnresolvedFunction
	case errors.Is(err, ErrInvalidFunctionName):
		spec.Kind = UnresolvedSymbol
		spec.Message = "invalid function name"
		spec.Label = "function name is invalid"
		spec.Note = "host call target must resolve to a string function name"
		spec.Hint = "Pass a string function name to the host call"
		spec.Cause = ErrInvalidFunctionName
	case errors.As(err, &invariantErr):
		spec.Kind = diagnostics.UnexpectedError
		spec.Message = "vm invariant violation"
		spec.Label = "internal invariant violated"
		spec.Note = detailNote(invariantErr.Message)
		spec.Hint = "This indicates an internal VM bug; please report it with the query and stack context"
		spec.Cause = invariantErr.Cause
	default:
		detail, cause := unwrapRuntimeDetail(err)

		spec.Message = fallbackRuntimeMessage(err)
		spec.Note = detailNote(detail)
		spec.Cause = cause
	}

	return newRuntimeErrorWithSpec(program, callStack, spec)
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
