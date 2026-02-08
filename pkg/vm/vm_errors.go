package vm

import (
	"context"
	"errors"
	"fmt"

	"github.com/MontFerret/ferret/pkg/file"
	"github.com/MontFerret/ferret/pkg/runtime"
	"github.com/MontFerret/ferret/pkg/vm/internal/data"
)

func (vm *VM) wrapRuntimeError(err error) error {
	if err == nil {
		return nil
	}

	if _, ok := err.(*RuntimeError); ok {
		return err
	}

	message := err.Error()
	label := ""
	help := ""
	note := ""

	switch {
	case errors.Is(err, ErrDivisionByZero):
		message = "division by zero"
		label = "attempt to divide by zero"
		help = "ensure the denominator is non-zero before division"
		note = "add a conditional check before dividing"
	case errors.Is(err, ErrModuloByZero):
		message = "modulo by zero"
		label = "attempt to take modulo by zero"
		help = "ensure the divisor is non-zero before modulo"
		note = "add a conditional check before modulo"
	case errors.Is(err, runtime.ErrInvalidType):
		message = "invalid type"
		label = "type mismatch"
		help = "ensure the value has the expected type"
		note = err.Error()
	case errors.Is(err, runtime.ErrInvalidArgumentType):
		message = "invalid argument type"
		help = "ensure the argument types match the function signature"
		note = err.Error()
	case errors.Is(err, runtime.ErrInvalidArgumentNumber):
		message = "invalid number of arguments"
		help = "check the function signature for the expected argument count"
		note = err.Error()
	case errors.Is(err, runtime.ErrInvalidArgument):
		message = "invalid argument"
		help = "check the function arguments"
		note = err.Error()
	case errors.Is(err, ErrMissedParam):
		message = "missing parameter"
		help = "provide all required parameters"
		note = err.Error()
	case errors.Is(err, ErrFunctionNotFound):
		message = "function not found"
		note = err.Error()
	case errors.Is(err, ErrUnresolvedFunction):
		message = "unresolved function"
		note = err.Error()
	case errors.Is(err, ErrInvalidFunctionName):
		message = "invalid function name"
		note = err.Error()
	}

	return &RuntimeError{
		Message: message,
		Hint:    help,
		Note:    note,
		Label:   label,
		Source:  vm.program.Source,
		Span:    vm.spanAt(vm.pc - 1),
		Cause:   err,
	}
}

func (vm *VM) runtimeErrorFromPanic(r any) error {
	message := "unexpected runtime panic"
	cause := fmt.Errorf("panic: %v", r)

	if err, ok := r.(error); ok {
		cause = err
	}

	return &RuntimeError{
		Message: message,
		Note:    cause.Error(),
		Source:  vm.program.Source,
		Span:    vm.spanAt(vm.pc - 1),
		Cause:   cause,
	}
}

func (vm *VM) newRuntimeError(message, label, help, note string, cause error) *RuntimeError {
	return &RuntimeError{
		Message: message,
		Hint:    help,
		Note:    note,
		Label:   label,
		Source:  vm.program.Source,
		Span:    vm.spanAt(vm.pc - 1),
		Cause:   cause,
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
			"division by zero",
			"attempt to divide by zero",
			"ensure the denominator is non-zero before division",
			"add a conditional check before dividing",
			ErrDivisionByZero,
		)
	}

	return nil
}

func (vm *VM) checkModuloByZero(ctx context.Context, right runtime.Value) error {
	rv, _ := runtime.ToInt(ctx, right)
	if rv == 0 {
		return vm.newRuntimeError(
			"modulo by zero",
			"attempt to take modulo by zero",
			"ensure the divisor is non-zero before modulo",
			"add a conditional check before modulo",
			ErrModuloByZero,
		)
	}

	return nil
}
