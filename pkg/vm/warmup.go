package vm

import (
	"github.com/MontFerret/ferret/v2/pkg/bytecode"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
	"github.com/MontFerret/ferret/v2/pkg/vm/internal/data"
	"github.com/MontFerret/ferret/v2/pkg/vm/internal/diagnostic"
	"github.com/MontFerret/ferret/v2/pkg/vm/internal/mem"
)

type hostCallWarmupDescriptor struct {
	FnName    string
	PC        int
	ID        int
	Dst       bytecode.Operand
	ArgCount  int
	HasFnName bool
}

func warmup(vm *VM, state *execState, env *Environment) error {
	if err := state.bindParams(env); err != nil {
		return err
	}

	return warmupShared(vm, env)
}

func warmupShared(vm *VM, env *Environment) error {
	warmupRegexps(vm)

	if len(vm.hostWarmups) == 0 {
		vm.cache.HostFunctionsWarmed = true
		return nil
	}

	if vm.cache.HostFunctionsWarmed && vm.cache.FunctionsRef == env.Functions {
		return nil
	}

	hash := env.Functions.Hash()

	if vm.cache.HostFunctionsWarmed && vm.cache.FuncHash == hash {
		vm.cache.FunctionsRef = env.Functions
		return nil
	}

	var warmupErrs diagnostic.WarmupErrorSet

	functions := env.Functions
	for _, descriptor := range vm.hostWarmups {
		if descriptor.ID < 0 || descriptor.ID >= len(vm.cache.HostFunctions) {
			warmupErrs.Add(
				diagnostic.NewInvariantError(
					"invalid host warmup slot",
					runtime.Errorf(runtime.ErrUnexpected, "invalid host warmup slot %d at pc %d", descriptor.ID, descriptor.PC),
				),
				descriptor.PC,
				descriptor.Dst,
			)
			continue
		}

		vm.cache.HostFunctions[descriptor.ID] = mem.CachedHostFunction{}

		cachedFn, err := warmupBindHostCall(descriptor, functions)
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

	vm.cache.FuncHash = hash
	vm.cache.FunctionsRef = env.Functions
	vm.cache.HostFunctionsWarmed = true

	return nil
}

func warmupRegexps(vm *VM) {
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

func warmupArgCount(src1, src2 bytecode.Operand) int {
	argCount := 0

	if src1.IsRegister() && src2.IsRegister() {
		start := src1.Register()
		end := src2.Register()

		if start > 0 && end >= start {
			argCount = end - start + 1
		}
	}

	return argCount
}

func warmupBindHostCall(descriptor hostCallWarmupDescriptor, functions *runtime.Functions) (mem.CachedHostFunction, error) {
	if !descriptor.HasFnName {
		return mem.CachedHostFunction{}, ErrInvalidFunctionName
	}

	argCount := descriptor.ArgCount

	switch argCount {
	case 0:
		return resolveHostFn(functions.A0().Get, functions.Var(), func(f *mem.CachedHostFunction, fn runtime.Function0) { f.Fn0 = fn }, descriptor.FnName)
	case 1:
		return resolveHostFn(functions.A1().Get, functions.Var(), func(f *mem.CachedHostFunction, fn runtime.Function1) { f.Fn1 = fn }, descriptor.FnName)
	case 2:
		return resolveHostFn(functions.A2().Get, functions.Var(), func(f *mem.CachedHostFunction, fn runtime.Function2) { f.Fn2 = fn }, descriptor.FnName)
	case 3:
		return resolveHostFn(functions.A3().Get, functions.Var(), func(f *mem.CachedHostFunction, fn runtime.Function3) { f.Fn3 = fn }, descriptor.FnName)
	case 4:
		return resolveHostFn(functions.A4().Get, functions.Var(), func(f *mem.CachedHostFunction, fn runtime.Function4) { f.Fn4 = fn }, descriptor.FnName)
	default:
		return resolveHostFn(functions.Var().Get, functions.Var(), func(f *mem.CachedHostFunction, fn runtime.Function) { f.FnV = fn }, descriptor.FnName)
	}
}

func inlineHostCallIDs(instructions []data.ExecInstruction, descriptors []hostCallWarmupDescriptor) {
	for pc := range instructions {
		op := instructions[pc].Opcode

		if op == bytecode.OpHCall || op == bytecode.OpProtectedHCall {
			instructions[pc].InlineSlot = -1
		}
	}

	for id := range descriptors {
		descriptors[id].ID = id
		pc := descriptors[id].PC
		if pc < 0 || pc >= len(instructions) {
			continue
		}

		op := instructions[pc].Opcode
		if op != bytecode.OpHCall && op != bytecode.OpProtectedHCall {
			continue
		}

		instructions[pc].InlineSlot = id
	}
}

func buildHostWarmupDescriptors(program *bytecode.Program) []hostCallWarmupDescriptor {
	if program == nil || len(program.Bytecode) == 0 {
		return nil
	}

	constants := program.Constants
	reg := map[bytecode.Operand]runtime.Value{}
	descriptors := make([]hostCallWarmupDescriptor, 0, 8)

	for pc, inst := range program.Bytecode {
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
		case bytecode.OpHCall, bytecode.OpProtectedHCall:
			descriptor := hostCallWarmupDescriptor{
				PC:       pc,
				Dst:      dst,
				ArgCount: warmupArgCount(src1, src2),
			}

			fnName, err := resolveHostFnName(reg, dst)
			if err == nil {
				descriptor.FnName = fnName
				descriptor.HasFnName = true
			}

			descriptors = append(descriptors, descriptor)
		}

		if op != bytecode.OpLoadConst && op != bytecode.OpMove && dst.IsRegister() {
			delete(reg, dst)
		}
	}

	if len(descriptors) == 0 {
		return nil
	}

	return descriptors
}
