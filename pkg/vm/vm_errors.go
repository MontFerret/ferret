package vm

import (
	"context"
	"errors"
	"fmt"

	"github.com/MontFerret/ferret/pkg/diagnostics"
	"github.com/MontFerret/ferret/pkg/file"
	"github.com/MontFerret/ferret/pkg/runtime"
	"github.com/MontFerret/ferret/pkg/vm/internal/data"
)

func (vm *VM) wrapRuntimeError(err error) error {
	if err == nil {
		return nil
	}

	var runtimeError *RuntimeError

	if errors.As(err, &runtimeError) {
		return err
	}

	kind := diagnostics.Unknown
	message := err.Error()
	label := ""
	hint := ""
	note := ""

	switch {
	case errors.Is(err, ErrDivisionByZero):
		kind = DivideByZero
		message = "Division by zero"
		label = "attempt to divide by zero"
		hint = "Ensure the denominator is non-zero before division"
		note = "Add a conditional check before dividing"
	case errors.Is(err, ErrModuloByZero):
		kind = ModuloByZero
		message = "modulo by zero"
		label = "attempt to take modulo by zero"
		hint = "Ensure the divisor is non-zero before modulo"
		note = "Add a conditional check before modulo"
	case errors.Is(err, runtime.ErrInvalidType):
		kind = diagnostics.TypeError
		message = "invalid type"
		label = "type mismatch"
		hint = "Ensure the value has the expected type"
		note = err.Error()
	case errors.Is(err, runtime.ErrInvalidArgumentType):
		kind = diagnostics.TypeError
		message = "invalid argument type"
		hint = "Ensure the argument types match the function signature"
		note = err.Error()
	case errors.Is(err, runtime.ErrInvalidArgumentNumber):
		kind = ArityError
		message = "invalid number of arguments"
		hint = "Check the function signature for the expected argument count"
		note = err.Error()
	case errors.Is(err, runtime.ErrInvalidArgument):
		kind = ArityError
		message = "invalid argument"
		hint = "Check the function arguments"
		note = err.Error()
	case errors.Is(err, ErrMissedParam):
		kind = UnresolvedSymbol
		message = "missing parameter"
		hint = "Provide all required parameters"
		note = err.Error()
	case errors.Is(err, ErrFunctionNotFound):
		kind = UnresolvedSymbol
		message = "function not found"
		hint = "Ensure the function is registered and accessible in the current context"
		note = "Add the function to the registry if it's missing"
	case errors.Is(err, ErrUnresolvedFunction):
		kind = UnresolvedSymbol
		message = "Unresolved function"
		hint = "Ensure the function is registered and accessible in the current context"
		note = "Add the function to the registry if it's missing"
	case errors.Is(err, ErrInvalidFunctionName):
		message = "invalid function name"
		hint = "Ensure the function name is valid and does not contain illegal characters"
	}

	return &RuntimeError{
		Diagnostic: &diagnostics.Diagnostic{
			Kind:    kind,
			Message: message,
			Hint:    hint,
			Note:    note,
			Source:  vm.program.Source,
			Spans:   []diagnostics.ErrorSpan{diagnostics.NewMainErrorSpan(vm.spanAt(vm.pc-1), label)},
			Cause:   err,
		},
	}
}

func (vm *VM) runtimeErrorFromPanic(r any) error {
	message := "unexpected runtime panic"
	cause := fmt.Errorf("panic: %v", r)

	if err, ok := r.(error); ok {
		cause = err
	}

	return &RuntimeError{
		Diagnostic: &diagnostics.Diagnostic{
			Kind:    diagnostics.UnexpectedError,
			Message: fmt.Sprintf("%s. %s", message, cause.Error()),
			Source:  vm.program.Source,
			Spans:   []diagnostics.ErrorSpan{diagnostics.NewMainErrorSpan(vm.spanAt(vm.pc-1), "")},
			Cause:   cause,
		},
	}
}

func (vm *VM) newRuntimeError(kind diagnostics.Kind, message, label, hint, note string, cause error) *RuntimeError {
	return &RuntimeError{
		Diagnostic: &diagnostics.Diagnostic{
			Kind:    kind,
			Message: message,
			Hint:    hint,
			Note:    note,
			Source:  vm.program.Source,
			Spans:   []diagnostics.ErrorSpan{diagnostics.NewMainErrorSpan(vm.spanAt(vm.pc-1), label)},
			Cause:   cause,
		},
	}
}

func (vm *VM) spanAt(pc int) file.Span {
	if vm == nil || vm.program == nil {
		return file.Span{Start: -1, End: -1}
	}

	if pc < 0 || pc >= len(vm.program.DebugSpans) {
		return file.Span{Start: -1, End: -1}
	}

	return vm.program.DebugSpans[pc]
}

func (vm *VM) checkDivisionByZero(ctx context.Context, left, right runtime.Value) error {
	l := data.ToNumberOnly(ctx, left)
	if _, ok := l.(runtime.Int); !ok {
		return nil
	}

	r := data.ToNumberOnly(ctx, right)
	if rv, ok := r.(runtime.Int); ok && rv == 0 {
		return vm.newRuntimeError(
			DivideByZero,
			"division by zero",
			"attempt to divide by zero",
			"Ensure the denominator is non-zero before division",
			"Add a conditional check before dividing",
			ErrDivisionByZero,
		)
	}

	return nil
}

func (vm *VM) checkModuloByZero(ctx context.Context, right runtime.Value) error {
	rv, _ := runtime.ToInt(ctx, right)
	if rv == 0 {
		return vm.newRuntimeError(
			ModuloByZero,
			"modulo by zero",
			"attempt to take modulo by zero",
			"Ensure the divisor is non-zero before modulo",
			"Add a conditional check before modulo",
			ErrModuloByZero,
		)
	}

	return nil
}
