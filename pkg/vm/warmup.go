package vm

import (
	"github.com/MontFerret/ferret/v2/pkg/bytecode"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
	"github.com/MontFerret/ferret/v2/pkg/vm/internal/data"
	"github.com/MontFerret/ferret/v2/pkg/vm/internal/diagnostics"
	"github.com/MontFerret/ferret/v2/pkg/vm/internal/mem"
)

func warmup(vm *VM, env *Environment) error {
	if err := ensureHostFunctionsBound(vm, env); err != nil {
		return err
	}

	return ensureRegexpsWarmed(vm)
}

func ensureRegexpsWarmed(vm *VM) error {
	if vm.cache.RegexpsWarmed {
		return nil
	}

	constants := vm.program.Constants
	reg := map[bytecode.Operand]runtime.Value{}
	var warmupErrs diagnostics.WarmupErrorSet

	for pc, inst := range vm.program.Bytecode {
		op := inst.Opcode
		dst, src1, src2 := inst.Operands[0], inst.Operands[1], inst.Operands[2]

		switch op {
		case bytecode.OpLoadConst:
			reg[dst] = constants[src1.Constant()]
		case bytecode.OpMove, bytecode.OpMoveTracked:
			if val, ok := reg[src1]; ok {
				reg[dst] = val
			} else {
				delete(reg, dst)
			}
		case bytecode.OpRegexp:
			if val, ok := reg[src2]; ok {
				r, err := data.DecodeRegexp(val)

				if err != nil {
					warmupErrs.Add(err, pc, dst)
					continue
				}

				pattern := r.String()
				if cached := vm.cache.Regexps[pc]; cached == nil || cached.Pattern != pattern {
					vm.cache.Regexps[pc] = &mem.CachedRegexp{Pattern: pattern, Regexp: r}
				}
			}
		}

		if op != bytecode.OpLoadConst && op != bytecode.OpMove && op != bytecode.OpMoveTracked && dst.IsRegister() {
			delete(reg, dst)
		}
	}

	if warmupErrs.Size() > 0 {
		return &warmupErrs
	}

	vm.cache.RegexpsWarmed = true

	return nil
}

func ensureHostFunctionsBound(vm *VM, env *Environment) error {
	hostCallDescriptors := vm.plan.hostCallDescriptors
	if len(hostCallDescriptors) == 0 {
		return nil
	}

	if vm.cache.FunctionsRef == env.Functions {
		return nil
	}

	var warmupErrs diagnostics.WarmupErrorSet

	for i, descriptor := range hostCallDescriptors {
		if descriptor.ID != i {
			warmupErrs.Add(
				diagnostics.NewInvariantError(
					"invalid host warmup binding id",
					runtime.Errorf(runtime.ErrUnexpected, "invalid host warmup binding id %d at index %d", descriptor.ID, i),
				),
				descriptor.PC,
				descriptor.Dst,
			)

			continue
		}

		if descriptor.ID < 0 || descriptor.ID >= len(vm.cache.HostFunctions) {
			warmupErrs.Add(
				diagnostics.NewInvariantError(
					"invalid host warmup slot",
					runtime.Errorf(runtime.ErrUnexpected, "invalid host warmup slot %d at pc %d", descriptor.ID, descriptor.PC),
				),
				descriptor.PC,
				descriptor.Dst,
			)

			continue
		}

		vm.cache.HostFunctions[descriptor.ID] = mem.CachedHostFunction{}
		cachedFn, err := warmupBindHostCall(descriptor, env.Functions)

		if err != nil {
			warmupErrs.Add(err, descriptor.PC, descriptor.Dst)
			continue
		}

		cachedFn.Bound = true
		vm.cache.HostFunctions[descriptor.ID] = cachedFn
	}

	if warmupErrs.Size() > 0 {
		return &warmupErrs
	}

	vm.cache.FunctionsRef = env.Functions

	return nil
}

func resolveHostFn[T runtime.FunctionConstraint](
	primary func(name string) (T, bool),
	fallback runtime.FunctionCollection[runtime.Function],
	setter func(*mem.CachedHostFunction, T),
	fnName string,
) (mem.CachedHostFunction, error) {
	if fn, ok := primary(fnName); ok {
		var c mem.CachedHostFunction
		setter(&c, fn)
		return c, nil
	}

	if fallback != nil {
		if fnv, ok := fallback.Get(fnName); ok {
			return mem.CachedHostFunction{FnV: fnv}, nil
		}
	}

	return mem.CachedHostFunction{}, ErrUnresolvedFunction
}

func warmupBindHostCall(descriptor callDescriptor, functions *runtime.Functions) (mem.CachedHostFunction, error) {
	if descriptor.DisplayName == "" {
		return mem.CachedHostFunction{}, ErrInvalidFunctionName
	}

	argCount := descriptor.ArgCount

	switch argCount {
	case 0:
		return resolveHostFn(functions.A0().Get, functions.Var(), func(f *mem.CachedHostFunction, fn runtime.Function0) { f.Fn0 = fn }, descriptor.DisplayName)
	case 1:
		return resolveHostFn(functions.A1().Get, functions.Var(), func(f *mem.CachedHostFunction, fn runtime.Function1) { f.Fn1 = fn }, descriptor.DisplayName)
	case 2:
		return resolveHostFn(functions.A2().Get, functions.Var(), func(f *mem.CachedHostFunction, fn runtime.Function2) { f.Fn2 = fn }, descriptor.DisplayName)
	case 3:
		return resolveHostFn(functions.A3().Get, functions.Var(), func(f *mem.CachedHostFunction, fn runtime.Function3) { f.Fn3 = fn }, descriptor.DisplayName)
	case 4:
		return resolveHostFn(functions.A4().Get, functions.Var(), func(f *mem.CachedHostFunction, fn runtime.Function4) { f.Fn4 = fn }, descriptor.DisplayName)
	default:
		return resolveHostFn(functions.Var().Get, functions.Var(), func(f *mem.CachedHostFunction, fn runtime.Function) { f.FnV = fn }, descriptor.DisplayName)
	}
}
