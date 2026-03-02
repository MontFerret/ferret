package vm

import (
	"github.com/MontFerret/ferret/v2/pkg/bytecode"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
	"github.com/MontFerret/ferret/v2/pkg/vm/internal/data"
	"github.com/MontFerret/ferret/v2/pkg/vm/internal/diagnostic"
	"github.com/MontFerret/ferret/v2/pkg/vm/internal/mem"
)

func (vm *VM) warmup(env *Environment) error {
	vm.warmupRegexps()

	hash := env.Functions.Hash()

	if vm.cache.FuncHash == hash || hash == 0 {
		return nil
	}

	errors := &diagnostic.WarmupErrorSet{}
	constants := vm.program.Constants
	functions := env.Functions
	reg := map[bytecode.Operand]runtime.Value{}

	for pc, inst := range vm.program.Bytecode {
		op := inst.Opcode
		dst, src1 := inst.Operands[0], inst.Operands[1]

		switch op {
		case bytecode.OpLoadConst:
			reg[dst] = constants[src1.Constant()]
		case bytecode.OpMove:
			reg[dst] = reg[src1]
		case bytecode.OpHCall, bytecode.OpProtectedHCall,
			bytecode.OpHCall0, bytecode.OpProtectedHCall0,
			bytecode.OpHCall1, bytecode.OpProtectedHCall1,
			bytecode.OpHCall2, bytecode.OpProtectedHCall2,
			bytecode.OpHCall3, bytecode.OpProtectedHCall3,
			bytecode.OpHCall4, bytecode.OpProtectedHCall4:
			warmupResolveHostCall(pc, op, dst, reg, functions, vm.cache.HostFunctions, errors)
		default:
			continue
		}
	}

	if errors.Size() > 0 {
		return errors
	}

	vm.cache.FuncHash = env.Functions.Hash()

	return nil
}

func (vm *VM) warmupRegexps() {
	if vm.cache.RegexpsWarmed {
		return
	}

	constants := vm.program.Constants
	reg := map[bytecode.Operand]runtime.Value{}

	for pc, inst := range vm.program.Bytecode {
		op := inst.Opcode
		dst, src1, src2 := inst.Operands[0], inst.Operands[1], inst.Operands[2]

		switch op {
		case bytecode.OpLoadConst:
			reg[dst] = constants[src1.Constant()]
		case bytecode.OpMove:
			if val, ok := reg[src1]; ok {
				reg[dst] = val
			} else {
				delete(reg, dst)
			}
		case bytecode.OpRegexp:
			if val, ok := reg[src2]; ok {
				r, err := data.ToRegexp(val)

				if err == nil {
					pattern := r.String()
					if cached := vm.cache.Regexps[pc]; cached == nil || cached.Pattern != pattern {
						vm.cache.Regexps[pc] = &mem.CachedRegexp{Pattern: pattern, Regexp: r}
					}
				}
			}
		}

		if op != bytecode.OpLoadConst && op != bytecode.OpMove && dst.IsRegister() {
			delete(reg, dst)
		}
	}

	vm.cache.RegexpsWarmed = true
}

func resolveHostFnName(reg map[bytecode.Operand]runtime.Value, dst bytecode.Operand) (string, error) {
	val, ok := reg[dst]

	if ok {
		fnName, ok := val.(runtime.String)

		if ok {
			return fnName.String(), nil
		}
	}

	return "", ErrInvalidFunctionName
}

func resolveHostFn[T runtime.FunctionConstraint](
	primary func(name string) (T, bool),
	fallback runtime.FunctionCollection[runtime.Function],
	setter func(*mem.CachedHostFunction, T),
	fnName string,
) (*mem.CachedHostFunction, error) {
	if fn, ok := primary(fnName); ok {
		c := &mem.CachedHostFunction{}
		setter(c, fn)
		return c, nil
	}

	if fallback != nil {
		if fnv, ok := fallback.Get(fnName); ok {
			return &mem.CachedHostFunction{FnV: fnv}, nil
		}
	}

	return nil, ErrUnresolvedFunction
}

func resolveHostFnAndCache[T runtime.FunctionConstraint](
	pc int,
	dst bytecode.Operand,
	reg map[bytecode.Operand]runtime.Value,
	get func(string) (T, bool),
	fallback runtime.FunctionCollection[runtime.Function],
	assign func(*mem.CachedHostFunction, T),
	funcs []*mem.CachedHostFunction,
	errList *diagnostic.WarmupErrorSet,
) {
	fnName, err := resolveHostFnName(reg, dst)

	if err != nil {
		errList.Add(err, pc, dst)
		return
	}

	fn, err := resolveHostFn(get, fallback, assign, fnName)

	if err != nil {
		errList.Add(err, pc, dst)
		return
	}

	funcs[pc] = fn
}

func resolveHostCall[T runtime.FunctionConstraint](
	pc int,
	dst bytecode.Operand,
	reg map[bytecode.Operand]runtime.Value,
	get func(string) (T, bool),
	functions *runtime.Functions,
	assign func(*mem.CachedHostFunction, T),
	funcs []*mem.CachedHostFunction,
	errList *diagnostic.WarmupErrorSet,
) {
	resolveHostFnAndCache(pc, dst, reg, get, functions.Var(), assign, funcs, errList)
}

func warmupResolveHostCall(
	pc int,
	op bytecode.Opcode,
	dst bytecode.Operand,
	reg map[bytecode.Operand]runtime.Value,
	functions *runtime.Functions,
	funcs []*mem.CachedHostFunction,
	errors *diagnostic.WarmupErrorSet,
) {
	switch op {
	case bytecode.OpHCall, bytecode.OpProtectedHCall:
		resolveHostCall(pc, dst, reg, functions.Var().Get, functions, func(f *mem.CachedHostFunction, fn runtime.Function) { f.FnV = fn }, funcs, errors)
	case bytecode.OpHCall0, bytecode.OpProtectedHCall0:
		resolveHostCall(pc, dst, reg, functions.A0().Get, functions, func(f *mem.CachedHostFunction, fn runtime.Function0) { f.Fn0 = fn }, funcs, errors)
	case bytecode.OpHCall1, bytecode.OpProtectedHCall1:
		resolveHostCall(pc, dst, reg, functions.A1().Get, functions, func(f *mem.CachedHostFunction, fn runtime.Function1) { f.Fn1 = fn }, funcs, errors)
	case bytecode.OpHCall2, bytecode.OpProtectedHCall2:
		resolveHostCall(pc, dst, reg, functions.A2().Get, functions, func(f *mem.CachedHostFunction, fn runtime.Function2) { f.Fn2 = fn }, funcs, errors)
	case bytecode.OpHCall3, bytecode.OpProtectedHCall3:
		resolveHostCall(pc, dst, reg, functions.A3().Get, functions, func(f *mem.CachedHostFunction, fn runtime.Function3) { f.Fn3 = fn }, funcs, errors)
	case bytecode.OpHCall4, bytecode.OpProtectedHCall4:
		resolveHostCall(pc, dst, reg, functions.A4().Get, functions, func(f *mem.CachedHostFunction, fn runtime.Function4) { f.Fn4 = fn }, funcs, errors)
	}
}
