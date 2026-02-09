package vm

import (
	"context"
	"errors"
	"fmt"

	"github.com/MontFerret/ferret/pkg/diagnostics"
	"github.com/MontFerret/ferret/pkg/runtime"
)

func (vm *VM) wrapRuntimeError(err error) error {
	if err == nil {
		return nil
	}

	var runtimeError *RuntimeError

	if errors.As(err, &runtimeError) {
		return err
	}

	var wpErrorSet *warmupErrorSet

	if errors.As(err, &wpErrorSet) {
		errs := NewRuntimeErrorSet(5)

		for _, wer := range wpErrorSet.Errors {
			// warmup PCs are zero-based instruction indices (no pre-increment),
			// while toRuntimeError expects a post-increment pc (see pc-1 usage)
			errs.Add(toRuntimeError(vm.program, wer.PC+1, wer.Err))
		}

		return errs
	}

	return toRuntimeError(vm.program, vm.pc, err)
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
			Spans:   []diagnostics.ErrorSpan{diagnostics.NewMainErrorSpan(spanAt(vm.program, vm.pc-1), "")},
			Cause:   cause,
		},
	}
}

func (vm *VM) newRuntimeError(kind diagnostics.Kind, message, label, hint, note string) *RuntimeError {
	return &RuntimeError{
		Diagnostic: &diagnostics.Diagnostic{
			Kind:    kind,
			Message: message,
			Hint:    hint,
			Note:    note,
			Source:  vm.program.Source,
			Spans:   []diagnostics.ErrorSpan{diagnostics.NewMainErrorSpan(spanAt(vm.program, vm.pc-1), label)},
		},
	}
}

func (vm *VM) checkDivisionByZero(ctx context.Context, left, right runtime.Value) error {
	l := runtime.ToNumberOnly(ctx, left)
	if _, ok := l.(runtime.Int); !ok {
		return nil
	}

	r := runtime.ToNumberOnly(ctx, right)
	if rv, ok := r.(runtime.Int); ok && rv == 0 {
		return vm.newRuntimeError(
			DivideByZero,
			"Division by zero",
			"attempt to divide by zero",
			"Ensure the denominator is non-zero before division",
			"Add a conditional check before dividing",
		)
	}

	return nil
}

func (vm *VM) checkModuloByZero(ctx context.Context, right runtime.Value) error {
	rv, _ := runtime.ToInt(ctx, right)
	if rv == 0 {
		return vm.newRuntimeError(
			ModuloByZero,
			"Modulo by zero",
			"attempt to take modulo by zero",
			"Ensure the divisor is non-zero before modulo",
			"Add a conditional check before modulo",
		)
	}

	return nil
}
