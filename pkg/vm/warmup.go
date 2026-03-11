package vm

import (
	"github.com/MontFerret/ferret/v2/pkg/bytecode"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
	"github.com/MontFerret/ferret/v2/pkg/vm/internal/data"
	"github.com/MontFerret/ferret/v2/pkg/vm/internal/diagnostic"
	"github.com/MontFerret/ferret/v2/pkg/vm/internal/mem"
)

type hostCallBindingDescriptor struct {
	FnName    string
	ID        int
	ArgCount  int
	HasFnName bool
}

type hostCallsiteWarmup struct {
	PC        int
	Dst       bytecode.Operand
	BindingID int
}

func warmup(vm *VM, state *execState, env *Environment) error {
	if err := state.bindParams(env); err != nil {
		return err
	}

	return warmupShared(vm, env)
}

func warmupShared(vm *VM, env *Environment) error {
	warmupRegexps(vm)

	if len(vm.hostBindings) == 0 {
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
	bindingErrors := make([]error, len(vm.hostBindings))

	for i, descriptor := range vm.hostBindings {
		if descriptor.ID != i {
			bindingErrors[i] = diagnostic.NewInvariantError(
				"invalid host warmup binding id",
				runtime.Errorf(runtime.ErrUnexpected, "invalid host warmup binding id %d at index %d", descriptor.ID, i),
			)
			continue
		}

		vm.cache.HostFunctions[descriptor.ID] = mem.CachedHostFunction{}

		cachedFn, err := warmupBindHostCall(descriptor, functions)
		if err != nil {
			bindingErrors[descriptor.ID] = err
			continue
		}

		cachedFn.Bound = true
		vm.cache.HostFunctions[descriptor.ID] = cachedFn
	}

	for _, site := range vm.hostWarmupSites {
		bindingID := site.BindingID

		if bindingID < 0 || bindingID >= len(bindingErrors) {
			warmupErrs.Add(
				diagnostic.NewInvariantError(
					"invalid host warmup slot",
					runtime.Errorf(runtime.ErrUnexpected, "invalid host warmup slot %d at pc %d", bindingID, site.PC),
				),
				site.PC,
				site.Dst,
			)
			continue
		}

		if err := bindingErrors[bindingID]; err != nil {
			warmupErrs.Add(err, site.PC, site.Dst)
		}
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

func warmupBindHostCall(descriptor hostCallBindingDescriptor, functions *runtime.Functions) (mem.CachedHostFunction, error) {
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

type hostBindingKey struct {
	FnName    string
	BindClass uint8
}

const hostBindClassVarArg uint8 = 5

func hostBindClass(argCount int) uint8 {
	if argCount >= 0 && argCount <= 4 {
		return uint8(argCount)
	}

	return hostBindClassVarArg
}

func buildExecPlan(program *bytecode.Program) ([]data.ExecInstruction, []hostCallBindingDescriptor, []hostCallsiteWarmup) {
	if program == nil || len(program.Bytecode) == 0 {
		return nil, nil, nil
	}

	instructions := make([]data.ExecInstruction, len(program.Bytecode))
	constants := program.Constants
	reg := map[bytecode.Operand]runtime.Value{}
	hostBindings := make([]hostCallBindingDescriptor, 0, 8)
	hostWarmupSites := make([]hostCallsiteWarmup, 0, 8)
	hostBindingByKey := map[hostBindingKey]int{}

	for pc, inst := range program.Bytecode {
		instructions[pc] = data.ExecInstruction{
			Instruction: inst,
		}

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
			descriptor := hostCallBindingDescriptor{
				ArgCount: warmupArgCount(src1, src2),
			}

			fnName, err := resolveHostFnName(reg, dst)
			if err == nil {
				descriptor.FnName = fnName
				descriptor.HasFnName = true

				key := hostBindingKey{
					FnName:    fnName,
					BindClass: hostBindClass(descriptor.ArgCount),
				}
				if id, ok := hostBindingByKey[key]; ok {
					descriptor.ID = id
				} else {
					descriptor.ID = len(hostBindings)
					hostBindingByKey[key] = descriptor.ID
					hostBindings = append(hostBindings, descriptor)
				}
			} else {
				descriptor.ID = len(hostBindings)
				hostBindings = append(hostBindings, descriptor)
			}

			instructions[pc].InlineSlot = descriptor.ID
			hostWarmupSites = append(hostWarmupSites, hostCallsiteWarmup{
				PC:        pc,
				Dst:       dst,
				BindingID: descriptor.ID,
			})
		}

		if op != bytecode.OpLoadConst && op != bytecode.OpMove && dst.IsRegister() {
			delete(reg, dst)
		}
	}

	if len(hostBindings) == 0 {
		hostBindings = nil
		hostWarmupSites = nil
	}

	return instructions, hostBindings, hostWarmupSites
}
