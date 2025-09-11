package vm

import (
	"github.com/MontFerret/ferret/pkg/runtime"
	"github.com/MontFerret/ferret/pkg/vm/internal/mem"
)

func (vm *VM) warmup(env *Environment) error {
	hash := env.Functions.Hash()

	if vm.cache.FuncHash == hash || hash == 0 {
		return nil
	}

	errors := make([]error, 0)
	constants := vm.program.Constants
	functions := env.Functions
	reg := map[Operand]runtime.Value{}

	for pc, inst := range vm.program.Bytecode {
		op := inst.Opcode
		dst, src1 := inst.Operands[0], inst.Operands[1]

		switch op {
		case OpLoadConst:
			reg[dst] = constants[src1.Constant()]
		case OpMove:
			reg[dst] = reg[src1]
		case OpCall, OpProtectedCall:
			resolveFnAndCache(
				pc, dst, reg,
				functions.FV().Get,
				functions.FV(),
				func(f *mem.CachedFunction, fn runtime.Function) { f.FnV = fn },
				vm.cache.Functions,
				&errors,
			)
		case OpCall0, OpProtectedCall0:
			resolveFnAndCache(
				pc, dst, reg,
				functions.F0().Get,
				functions.FV(),
				func(f *mem.CachedFunction, fn runtime.Function0) { f.Fn0 = fn },
				vm.cache.Functions,
				&errors,
			)
		case OpCall1, OpProtectedCall1:
			resolveFnAndCache(
				pc, dst, reg,
				functions.F1().Get,
				functions.FV(),
				func(f *mem.CachedFunction, fn runtime.Function1) { f.Fn1 = fn },
				vm.cache.Functions,
				&errors,
			)
		case OpCall2, OpProtectedCall2:
			resolveFnAndCache(
				pc, dst, reg,
				functions.F2().Get,
				functions.FV(),
				func(f *mem.CachedFunction, fn runtime.Function2) { f.Fn2 = fn },
				vm.cache.Functions,
				&errors,
			)
		case OpCall3, OpProtectedCall3:
			resolveFnAndCache(
				pc, dst, reg,
				functions.F3().Get,
				functions.FV(),
				func(f *mem.CachedFunction, fn runtime.Function3) { f.Fn3 = fn },
				vm.cache.Functions,
				&errors,
			)
		case OpCall4, OpProtectedCall4:
			resolveFnAndCache(
				pc, dst, reg,
				functions.F4().Get,
				functions.FV(),
				func(f *mem.CachedFunction, fn runtime.Function4) { f.Fn4 = fn },
				vm.cache.Functions,
				&errors,
			)
		default:
			continue
		}
	}

	if len(errors) > 0 {
		return runtime.Errorsf("failed to warm up the VM", errors...)
	}

	vm.cache.FuncHash = env.Functions.Hash()

	return nil
}

func resolveFnName(reg map[Operand]runtime.Value, pc int, dst Operand) (string, error) {
	val, ok := reg[dst]

	if ok {
		fnName, ok := val.(runtime.String)

		if ok {
			return fnName.String(), nil
		}
	}

	return "", runtime.Errorf(ErrInvalidFunctionName, "at pc=%d, dst=R%d", pc, dst)
}

func resolveFn[T runtime.FunctionConstraint](
	primary func(name string) (T, bool),
	fallback runtime.FunctionCollection[runtime.Function],
	setter func(*mem.CachedFunction, T),
	fnName string,
) (*mem.CachedFunction, error) {
	if fn, ok := primary(fnName); ok {
		c := &mem.CachedFunction{}
		setter(c, fn)
		return c, nil
	}

	if fallback != nil {
		if fnv, ok := fallback.Get(fnName); ok {
			return &mem.CachedFunction{FnV: fnv}, nil
		}
	}

	return nil, runtime.Error(ErrUnresolvedFunction, fnName)
}

func resolveFnAndCache[T runtime.FunctionConstraint](
	pc int,
	dst Operand,
	reg map[Operand]runtime.Value,
	get func(string) (T, bool),
	fallback runtime.FunctionCollection[runtime.Function],
	assign func(*mem.CachedFunction, T),
	funcs map[int]*mem.CachedFunction,
	errList *[]error,
) {
	fnName, err := resolveFnName(reg, pc, dst)

	if err != nil {
		*errList = append(*errList, runtime.Errorf(ErrInvalidFunctionName, "at %d", pc))
		return
	}

	fn, err := resolveFn(get, fallback, assign, fnName)

	if err != nil {
		*errList = append(*errList, err)
		return
	}

	funcs[pc] = fn
}
