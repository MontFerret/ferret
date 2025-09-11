package vm

import (
	"context"

	"github.com/MontFerret/ferret/pkg/runtime"
)

func (vm *VM) tryCatch(pos int) (Catch, bool) {
	for _, pair := range vm.program.CatchTable {
		if pos >= pair[0] && pos <= pair[1] {
			return pair, true
		}
	}

	return Catch{}, false
}

func (vm *VM) callv(ctx context.Context, pc int, src1, src2 Operand) (runtime.Value, error) {
	reg := vm.registers.Values
	cacheFn := vm.cache.Functions[pc]

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
	cacheFn := vm.cache.Functions[pc]

	if cacheFn.Fn0 != nil {
		return cacheFn.Fn0(ctx)
	}

	// Fall back to a variadic function call
	return cacheFn.FnV(ctx)
}

func (vm *VM) call1(ctx context.Context, pc int, src1 Operand) (runtime.Value, error) {
	reg := vm.registers.Values
	arg := reg[src1]
	cacheFn := vm.cache.Functions[pc]

	if cacheFn.Fn1 != nil {
		return cacheFn.Fn1(ctx, arg)
	}

	// Fall back to a variadic function call
	return cacheFn.FnV(ctx, arg)
}

func (vm *VM) call2(ctx context.Context, pc int, src1, src2 Operand) (runtime.Value, error) {
	reg := vm.registers.Values
	cacheFn := vm.cache.Functions[pc]
	arg1 := reg[src1]
	arg2 := reg[src2]

	if cacheFn.Fn2 != nil {
		return cacheFn.Fn2(ctx, arg1, arg2)
	}

	// Fall back to a variadic function call
	return cacheFn.FnV(ctx, arg1, arg2)
}

func (vm *VM) call3(ctx context.Context, pc int, src1 Operand) (runtime.Value, error) {
	reg := vm.registers.Values
	cacheFn := vm.cache.Functions[pc]
	arg1 := reg[src1]
	arg2 := reg[src1+1]
	arg3 := reg[src1+2]

	if cacheFn.Fn3 != nil {
		return cacheFn.Fn3(ctx, arg1, arg2, arg3)
	}

	// Fall back to a variadic function call
	return cacheFn.FnV(ctx, arg1, arg2, arg3)
}

func (vm *VM) call4(ctx context.Context, pc int, src1 Operand) (runtime.Value, error) {
	reg := vm.registers.Values
	cacheFn := vm.cache.Functions[pc]
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

func (vm *VM) loadIndex(ctx context.Context, src, arg runtime.Value) (runtime.Value, error) {
	indexed, ok := src.(runtime.Indexed)

	if !ok {
		return nil, runtime.TypeErrorOf(src, runtime.TypeIndexed)
	}

	var idx runtime.Int
	var err error

	switch v := arg.(type) {
	case runtime.Int:
		idx = v
	case runtime.Float:
		// Convert float to int, rounding down
		idx = runtime.Int(v)
	default:
		err = runtime.TypeErrorOf(arg, runtime.TypeInt, runtime.TypeFloat)
	}

	if err != nil {
		return nil, err
	}

	return indexed.Get(ctx, idx)
}

func (vm *VM) loadKey(ctx context.Context, src, arg runtime.Value) (runtime.Value, error) {
	keyed, ok := src.(runtime.Keyed)

	if !ok {
		return nil, runtime.TypeErrorOf(src, runtime.TypeKeyed)
	}

	out, err := keyed.Get(ctx, arg)

	if err != nil {
		return nil, err
	}

	return out, nil
}

func (vm *VM) castSubscribeArgs(dst, eventName, opts runtime.Value) (runtime.Observable, runtime.String, runtime.Map, error) {
	observable, ok := dst.(runtime.Observable)

	if !ok {
		return nil, "", nil, runtime.TypeErrorOf(dst, runtime.TypeObservable)
	}

	eventNameStr, ok := eventName.(runtime.String)

	if !ok {
		return nil, "", nil, runtime.TypeErrorOf(eventName, runtime.TypeString)
	}

	var options runtime.Map

	if opts != nil && opts != runtime.None {
		m, ok := opts.(runtime.Map)

		if !ok {
			return nil, "", nil, runtime.TypeErrorOf(opts, runtime.TypeMap)
		}

		options = m
	}

	return observable, eventNameStr, options, nil
}

func (vm *VM) setOrTryCatch(dst Operand, val runtime.Value, err error) error {
	reg := vm.registers.Values

	if err == nil {
		reg[dst] = val

		return nil
	}

	if _, catch := vm.tryCatch(vm.pc); catch {
		reg[dst] = runtime.None

		return nil
	}

	return err
}

func (vm *VM) setCallResult(op Opcode, dst Operand, out runtime.Value, err error) error {
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

	return err
}

func (vm *VM) setOrOptional(dst Operand, val runtime.Value, err error, optional bool) error {
	if err == nil {
		vm.registers.Values[dst] = val

		return nil
	}

	if optional {
		vm.registers.Values[dst] = runtime.None

		return nil
	}

	return err
}

func isProtectedCall(op Opcode) bool {
	switch op {
	case OpProtectedCall, OpProtectedCall0, OpProtectedCall1, OpProtectedCall2, OpProtectedCall3, OpProtectedCall4:
		return true
	default:
		return false
	}
}
