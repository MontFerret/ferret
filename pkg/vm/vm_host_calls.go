package vm

import (
	"context"

	"github.com/MontFerret/ferret/v2/pkg/bytecode"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
	"github.com/MontFerret/ferret/v2/pkg/vm/internal/mem"
)

func hostCallArgs(reg []runtime.Value, src1, src2 bytecode.Operand) []runtime.Value {
	if !src1.IsRegister() || !src2.IsRegister() {
		return nil
	}

	start := src1.Register()
	end := src2.Register()

	if start <= 0 || end < start {
		return nil
	}

	size := end - start + 1
	args := make([]runtime.Value, size)

	for i := 0; i < size; i++ {
		args[i] = reg[start+i]
	}

	return args
}

func callCachedHostFunction(
	ctx context.Context,
	cacheFn *mem.CachedHostFunction,
	args []runtime.Value,
) (runtime.Value, error) {
	switch len(args) {
	case 0:
		if cacheFn.Fn0 != nil {
			return cacheFn.Fn0(ctx)
		}
	case 1:
		if cacheFn.Fn1 != nil {
			return cacheFn.Fn1(ctx, args[0])
		}
	case 2:
		if cacheFn.Fn2 != nil {
			return cacheFn.Fn2(ctx, args[0], args[1])
		}
	case 3:
		if cacheFn.Fn3 != nil {
			return cacheFn.Fn3(ctx, args[0], args[1], args[2])
		}
	case 4:
		if cacheFn.Fn4 != nil {
			return cacheFn.Fn4(ctx, args[0], args[1], args[2], args[3])
		}
	}

	return cacheFn.FnV(ctx, args...)
}

func (vm *VM) execHostCall(
	ctx context.Context,
	op bytecode.Opcode,
	pc int,
	dst, src1, src2 bytecode.Operand,
) error {
	if op != bytecode.OpHCall && op != bytecode.OpProtectedHCall {
		return runtime.Error(runtime.ErrUnexpected, "invalid host call opcode")
	}

	cacheFn := vm.cache.HostFunctions[pc]
	args := hostCallArgs(vm.registers.Values, src1, src2)
	out, err := callCachedHostFunction(ctx, cacheFn, args)

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
	case bytecode.OpProtectedHCall, bytecode.OpProtectedCall:
		return true
	default:
		return false
	}
}
