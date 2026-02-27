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
		case bytecode.OpCall, bytecode.OpProtectedCall:
			resolveFnAndCache(
				pc, dst, reg,
				functions.Var().Get,
				functions.Var(),
				func(f *mem.CachedFunction, fn runtime.Function) { f.FnV = fn },
				vm.cache.Functions,
				errors,
			)
		case bytecode.OpCall0, bytecode.OpProtectedCall0:
			resolveFnAndCache(
				pc, dst, reg,
				functions.A0().Get,
				functions.Var(),
				func(f *mem.CachedFunction, fn runtime.Function0) { f.Fn0 = fn },
				vm.cache.Functions,
				errors,
			)
		case bytecode.OpCall1, bytecode.OpProtectedCall1:
			resolveFnAndCache(
				pc, dst, reg,
				functions.A1().Get,
				functions.Var(),
				func(f *mem.CachedFunction, fn runtime.Function1) { f.Fn1 = fn },
				vm.cache.Functions,
				errors,
			)
		case bytecode.OpCall2, bytecode.OpProtectedCall2:
			resolveFnAndCache(
				pc, dst, reg,
				functions.A2().Get,
				functions.Var(),
				func(f *mem.CachedFunction, fn runtime.Function2) { f.Fn2 = fn },
				vm.cache.Functions,
				errors,
			)
		case bytecode.OpCall3, bytecode.OpProtectedCall3:
			resolveFnAndCache(
				pc, dst, reg,
				functions.A3().Get,
				functions.Var(),
				func(f *mem.CachedFunction, fn runtime.Function3) { f.Fn3 = fn },
				vm.cache.Functions,
				errors,
			)
		case bytecode.OpCall4, bytecode.OpProtectedCall4:
			resolveFnAndCache(
				pc, dst, reg,
				functions.A4().Get,
				functions.Var(),
				func(f *mem.CachedFunction, fn runtime.Function4) { f.Fn4 = fn },
				vm.cache.Functions,
				errors,
			)
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

func resolveFnName(reg map[bytecode.Operand]runtime.Value, dst bytecode.Operand) (string, error) {
	val, ok := reg[dst]

	if ok {
		fnName, ok := val.(runtime.String)

		if ok {
			return fnName.String(), nil
		}
	}

	return "", ErrInvalidFunctionName
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

	return nil, ErrUnresolvedFunction
}

func resolveFnAndCache[T runtime.FunctionConstraint](
	pc int,
	dst bytecode.Operand,
	reg map[bytecode.Operand]runtime.Value,
	get func(string) (T, bool),
	fallback runtime.FunctionCollection[runtime.Function],
	assign func(*mem.CachedFunction, T),
	funcs map[int]*mem.CachedFunction,
	errList *diagnostic.WarmupErrorSet,
) {
	fnName, err := resolveFnName(reg, dst)

	if err != nil {
		errList.Add(err, pc, dst)
		return
	}

	fn, err := resolveFn(get, fallback, assign, fnName)

	if err != nil {
		errList.Add(err, pc, dst)
		return
	}

	funcs[pc] = fn
}
