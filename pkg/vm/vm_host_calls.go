package vm

import (
	"context"

	"github.com/MontFerret/ferret/v2/pkg/bytecode"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

func (vm *VM) callv(ctx context.Context, pc int, src1, src2 bytecode.Operand) (runtime.Value, error) {
	reg := vm.registers.Values
	cacheFn := vm.cache.HostFunctions[pc]

	var size int

	if src1 > 0 {
		size = src2.Register() - src1.Register() + 1
	}

	start := int(src1)
	end := int(src1) + size
	args := make([]runtime.Value, size)

	// Iterate over registers starting from src1 and up to the src2
	for i := start; i < end; i++ {
		args[i-start] = reg[i]
	}

	return cacheFn.FnV(ctx, args...)
}

func (vm *VM) call0(ctx context.Context, pc int) (runtime.Value, error) {
	cacheFn := vm.cache.HostFunctions[pc]

	if cacheFn.Fn0 != nil {
		return cacheFn.Fn0(ctx)
	}

	// Fall back to a variadic function call
	return cacheFn.FnV(ctx)
}

func (vm *VM) call1(ctx context.Context, pc int, src1 bytecode.Operand) (runtime.Value, error) {
	reg := vm.registers.Values
	arg := reg[src1]
	cacheFn := vm.cache.HostFunctions[pc]

	if cacheFn.Fn1 != nil {
		return cacheFn.Fn1(ctx, arg)
	}

	// Fall back to a variadic function call
	return cacheFn.FnV(ctx, arg)
}

func (vm *VM) call2(ctx context.Context, pc int, src1, src2 bytecode.Operand) (runtime.Value, error) {
	reg := vm.registers.Values
	cacheFn := vm.cache.HostFunctions[pc]
	arg1 := reg[src1]
	arg2 := reg[src2]

	if cacheFn.Fn2 != nil {
		return cacheFn.Fn2(ctx, arg1, arg2)
	}

	// Fall back to a variadic function call
	return cacheFn.FnV(ctx, arg1, arg2)
}

func (vm *VM) call3(ctx context.Context, pc int, src1 bytecode.Operand) (runtime.Value, error) {
	reg := vm.registers.Values
	cacheFn := vm.cache.HostFunctions[pc]
	arg1 := reg[src1]
	arg2 := reg[src1+1]
	arg3 := reg[src1+2]

	if cacheFn.Fn3 != nil {
		return cacheFn.Fn3(ctx, arg1, arg2, arg3)
	}

	// Fall back to a variadic function call
	return cacheFn.FnV(ctx, arg1, arg2, arg3)
}

func (vm *VM) call4(ctx context.Context, pc int, src1 bytecode.Operand) (runtime.Value, error) {
	reg := vm.registers.Values
	cacheFn := vm.cache.HostFunctions[pc]
	arg1 := reg[src1]
	arg2 := reg[src1+1]
	arg3 := reg[src1+2]
	arg4 := reg[src1+3]

	if cacheFn.Fn4 != nil {
		return cacheFn.Fn4(ctx, arg1, arg2, arg3, arg4)
	}

	// Fall back to a variadic function call
	return cacheFn.FnV(ctx, arg1, arg2, arg3, arg4)
}

func (vm *VM) execHostCall(
	ctx context.Context,
	op bytecode.Opcode,
	pc int,
	dst, src1, src2 bytecode.Operand,
) error {
	var (
		out runtime.Value
		err error
	)

	switch op {
	case bytecode.OpHCall, bytecode.OpProtectedHCall:
		out, err = vm.callv(ctx, pc, src1, src2)
	case bytecode.OpHCall0, bytecode.OpProtectedHCall0:
		out, err = vm.call0(ctx, pc)
	case bytecode.OpHCall1, bytecode.OpProtectedHCall1:
		out, err = vm.call1(ctx, pc, src1)
	case bytecode.OpHCall2, bytecode.OpProtectedHCall2:
		out, err = vm.call2(ctx, pc, src1, src2)
	case bytecode.OpHCall3, bytecode.OpProtectedHCall3:
		out, err = vm.call3(ctx, pc, src1)
	case bytecode.OpHCall4, bytecode.OpProtectedHCall4:
		out, err = vm.call4(ctx, pc, src1)
	default:
		return runtime.Error(runtime.ErrUnexpected, "invalid host call opcode")
	}

	if err := vm.setCallResult(op, dst, out, err); err != nil {
		if vm.unwindToProtected() {
			return nil
		}

		return err
	}

	return nil
}

func (vm *VM) setCallResult(op bytecode.Opcode, dst bytecode.Operand, out runtime.Value, err error) error {
	reg := vm.registers.Values

	if err == nil {
		reg[dst] = out

		return nil
	}

	if isProtectedCall(op) {
		reg[dst] = runtime.None

		return nil
	}

	if catch, ok := vm.tryCatch(vm.pc); ok {
		reg[dst] = runtime.None

		if catch[2] > 0 {
			vm.pc = catch[2]
		}

		return nil
	}

	if vm.unwindToProtected() {
		return nil
	}

	return err
}

func isProtectedCall(op bytecode.Opcode) bool {
	switch op {
	case bytecode.OpProtectedHCall, bytecode.OpProtectedHCall0, bytecode.OpProtectedHCall1, bytecode.OpProtectedHCall2, bytecode.OpProtectedHCall3, bytecode.OpProtectedHCall4,
		bytecode.OpProtectedCall, bytecode.OpProtectedCall0, bytecode.OpProtectedCall1, bytecode.OpProtectedCall2, bytecode.OpProtectedCall3, bytecode.OpProtectedCall4:
		return true
	default:
		return false
	}
}
