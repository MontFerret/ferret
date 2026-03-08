package vm

import (
	"strings"

	"github.com/MontFerret/ferret/v2/pkg/bytecode"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
	"github.com/MontFerret/ferret/v2/pkg/vm/internal/data"
	"github.com/MontFerret/ferret/v2/pkg/vm/internal/diagnostic"
	"github.com/MontFerret/ferret/v2/pkg/vm/internal/mem"
)

func (vm *VM) warmup(env *Environment) error {
	if err := vm.bindParams(env); err != nil {
		return err
	}

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
		dst, src1, src2 := inst.Operands[0], inst.Operands[1], inst.Operands[2]

		switch op {
		case bytecode.OpLoadConst:
			reg[dst] = constants[src1.Constant()]
		case bytecode.OpMove:
			reg[dst] = reg[src1]
		case bytecode.OpHCall, bytecode.OpProtectedHCall:
			warmupResolveHostCall(pc, op, dst, src1, src2, reg, functions, vm.cache.HostFunctions, errors)
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

func (vm *VM) bindParams(env *Environment) error {
	required := vm.program.Params

	if len(required) == 0 {
		vm.cache.Params = vm.cache.Params[:0]
		return nil
	}

	if len(vm.cache.Params) < len(required) {
		vm.cache.Params = make([]runtime.Value, len(required))
	} else {
		vm.cache.Params = vm.cache.Params[:len(required)]
	}

	var missedParams []string

	for idx, name := range required {
		val, exists := env.Params[name]
		if !exists {
			if missedParams == nil {
				missedParams = make([]string, 0, len(required))
			}

			missedParams = append(missedParams, "@"+name)
			continue
		}

		vm.cache.Params[idx] = val
	}

	if len(missedParams) > 0 {
		return runtime.Error(ErrMissedParam, strings.Join(missedParams, ", "))
	}

	return nil
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
	src1 bytecode.Operand,
	src2 bytecode.Operand,
	reg map[bytecode.Operand]runtime.Value,
	functions *runtime.Functions,
	funcs []*mem.CachedHostFunction,
	errors *diagnostic.WarmupErrorSet,
) {
	if op != bytecode.OpHCall && op != bytecode.OpProtectedHCall {
		return
	}

	argCount := 0

	if src1.IsRegister() && src2.IsRegister() {
		start := src1.Register()
		end := src2.Register()

		if start > 0 && end >= start {
			argCount = end - start + 1
		}
	}

	switch argCount {
	case 0:
		resolveHostCall(pc, dst, reg, functions.A0().Get, functions, func(f *mem.CachedHostFunction, fn runtime.Function0) { f.Fn0 = fn }, funcs, errors)
	case 1:
		resolveHostCall(pc, dst, reg, functions.A1().Get, functions, func(f *mem.CachedHostFunction, fn runtime.Function1) { f.Fn1 = fn }, funcs, errors)
	case 2:
		resolveHostCall(pc, dst, reg, functions.A2().Get, functions, func(f *mem.CachedHostFunction, fn runtime.Function2) { f.Fn2 = fn }, funcs, errors)
	case 3:
		resolveHostCall(pc, dst, reg, functions.A3().Get, functions, func(f *mem.CachedHostFunction, fn runtime.Function3) { f.Fn3 = fn }, funcs, errors)
	case 4:
		resolveHostCall(pc, dst, reg, functions.A4().Get, functions, func(f *mem.CachedHostFunction, fn runtime.Function4) { f.Fn4 = fn }, funcs, errors)
	default:
		resolveHostCall(pc, dst, reg, functions.Var().Get, functions, func(f *mem.CachedHostFunction, fn runtime.Function) { f.FnV = fn }, funcs, errors)
	}
}
