package vm

import (
	"context"

	"github.com/MontFerret/ferret/v2/pkg/bytecode"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
	"github.com/MontFerret/ferret/v2/pkg/vm/internal/data"
	"github.com/MontFerret/ferret/v2/pkg/vm/internal/frame"
	"github.com/MontFerret/ferret/v2/pkg/vm/internal/mem"
	"github.com/MontFerret/ferret/v2/pkg/vm/internal/operators"
)

type VM struct {
	options      options
	cache        *mem.Cache
	program      *bytecode.Program
	instructions []data.ExecInstruction
	catchByPC    []int
	frames       frame.CallStack
	exec         execState
}

func New(program *bytecode.Program) *VM {
	return NewWith(program)
}

func NewWith(program *bytecode.Program, opts ...Option) *VM {
	o := newOptions(opts)

	vm := &VM{
		cache:        mem.NewCache(len(program.Bytecode), o.shapeCacheLimit),
		program:      program,
		options:      o,
		instructions: buildExecInstructions(program.Bytecode),
		catchByPC:    buildCatchByPC(len(program.Bytecode), program.CatchTable),
	}

	vm.exec = execState{
		vm:        vm,
		registers: mem.NewRegisterFile(program.Registers),
		scratch:   mem.NewScratch(len(program.Params)),
	}
	vm.exec.errors = errorHandler{state: &vm.exec}
	vm.frames.Init(maxUDFRegisters(program.Functions.UserDefined))

	return vm
}

func (vm *VM) Run(ctx context.Context, env *Environment) (runtime.Value, error) {
	switch vm.options.panicPolicy {
	case PanicPropagate:
		return vm.runUnchecked(ctx, env)
	default:
		return vm.runRecovered(ctx, env)
	}
}

func (vm *VM) runRecovered(ctx context.Context, env *Environment) (result runtime.Value, err error) {
	exec := &vm.exec

	defer func() {
		if r := recover(); r != nil {
			err = exec.runtimeErrorFromPanic(r)
			result = nil

			return
		}

		if err != nil {
			err = exec.wrapRuntimeError(err)
		}
	}()

	return exec.runCore(ctx, env)
}

func (vm *VM) runUnchecked(ctx context.Context, env *Environment) (runtime.Value, error) {
	exec := &vm.exec
	result, err := exec.runCore(ctx, env)

	if err != nil {
		return nil, exec.wrapRuntimeError(err)
	}

	return result, nil
}

func (exec *execState) runCore(ctx context.Context, env *Environment) (runtime.Value, error) {
	if env == nil {
		env = noopEnv
	}

	if err := validate(env, exec.vm.program); err != nil {
		return nil, err
	}

	if err := exec.warmup(env); err != nil {
		return nil, err
	}

	if exec.registers.IsDirty() {
		exec.registers.Reset()
	}

	exec.registers.MarkDirty()
	exec.env = env
	exec.pc = 0
	exec.vm.frames.Reset()

	instructions := exec.vm.instructions
	constants := exec.vm.program.Constants
	shapeCache := exec.vm.cache.ShapeCache
	paramSlots := exec.scratch.Params
loop:
	for exec.pc < len(instructions) {
		reg := exec.registers.Values
		inst := &instructions[exec.pc]
		op := inst.Opcode
		dst, src1, src2 := inst.Operands[0], inst.Operands[1], inst.Operands[2]
		exec.pc++

		switch op {
		case bytecode.OpLoadNone:
			reg[dst] = runtime.None
		case bytecode.OpLoadZero:
			reg[dst] = runtime.ZeroInt
		case bytecode.OpLoadBool:
			reg[dst] = runtime.Boolean(src1 == 1)
		case bytecode.OpMove:
			reg[dst] = reg[src1]
		case bytecode.OpLoadConst:
			reg[dst] = constants[src1.Constant()]
		case bytecode.OpLoadParam:
			reg[dst] = paramSlots[int(src1)-1]
		case bytecode.OpLoadArray:
			reg[dst] = runtime.NewArray(int(src1))
		case bytecode.OpLoadObject:
			reg[dst] = data.NewFastObjectOf(shapeCache, exec.vm.options.fastObjectDictThreshold, int(src1))
		case bytecode.OpJump:
			exec.pc = int(dst)
		case bytecode.OpJumpIfFalse:
			if !runtime.ToBoolean(reg[src1]) {
				exec.pc = int(dst)
			}
		case bytecode.OpJumpIfTrue:
			if runtime.ToBoolean(reg[src1]) {
				exec.pc = int(dst)
			}
		case bytecode.OpJumpIfNone:
			if reg[src1] == runtime.None {
				exec.pc = int(dst)
			}
		case bytecode.OpJumpIfNe:
			if operators.NotEquals(ctx, reg[src1], reg[src2]) {
				exec.pc = int(dst)
			}
		case bytecode.OpJumpIfNeConst:
			if operators.NotEquals(ctx, reg[src1], constants[src2.Constant()]) {
				exec.pc = int(dst)
			}
		case bytecode.OpJumpIfEq:
			if operators.Equals(ctx, reg[src1], reg[src2]) {
				exec.pc = int(dst)
			}
		case bytecode.OpJumpIfEqConst:
			if operators.Equals(ctx, reg[src1], constants[src2.Constant()]) {
				exec.pc = int(dst)
			}
		case bytecode.OpJumpIfMissingProperty:
			obj, ok := reg[src1].(runtime.Map)
			if !ok {
				exec.pc = int(dst)
				continue
			}

			key, ok := reg[src2].(runtime.String)
			if !ok {
				exec.pc = int(dst)
				continue
			}

			has, err := obj.ContainsKey(ctx, key)
			if err != nil && exec.errors.protected(err) != nil {
				return nil, err
			}

			if !has {
				exec.pc = int(dst)
			}
		case bytecode.OpJumpIfMissingPropertyConst:
			obj, ok := reg[src1].(runtime.Map)
			if !ok {
				exec.pc = int(dst)

				continue
			}

			key, ok := constants[src2.Constant()].(runtime.String)
			if !ok {
				exec.pc = int(dst)
				continue
			}

			has, err := obj.ContainsKey(ctx, key)
			if err != nil && exec.errors.protected(err) != nil {
				return nil, err
			}

			if !has {
				exec.pc = int(dst)
			}
		case bytecode.OpExists:
			val := reg[src1]

			if val == runtime.None {
				reg[dst] = runtime.False
				continue
			}

			if measurable, ok := val.(runtime.Measurable); ok {
				length, err := measurable.Length(ctx)

				if err != nil {
					if err := exec.errors.handleWithCatch(err, func() {
						reg[dst] = runtime.False
					}); err != nil {
						return nil, err
					}

					continue
				}

				reg[dst] = runtime.NewBoolean(length != 0)

				continue
			}

			reg[dst] = runtime.True
		case bytecode.OpLength:
			val, ok := reg[src1].(runtime.Measurable)

			if ok {
				length, err := val.Length(ctx)

				if err != nil {
					if err := exec.errors.handleWithCatch(err, func() {
						length = 0
					}); err != nil {
						return nil, err
					}
				}

				reg[dst] = length
				continue
			}

			if err := exec.errors.handleWithCatch(runtime.TypeErrorOf(reg[src1],
				runtime.TypeString,
				runtime.TypeList,
				runtime.TypeMap,
				runtime.TypeBinary,
				runtime.TypeMeasurable,
			), func() {
				reg[dst] = runtime.ZeroInt
			}); err != nil {
				return runtime.None, err
			}

			continue
		case bytecode.OpType:
			reg[dst] = runtime.NewString(runtime.TypeName(runtime.TypeOf(reg[src1])))
		case bytecode.OpLoadIndex, bytecode.OpLoadIndexOptional:
			src := reg[src1]
			optional := op == bytecode.OpLoadIndexOptional
			arg := reg[src2]

			if err := exec.loadIndexAndSet(ctx, dst, src, arg, optional); err != nil {
				if err := exec.errors.protected(err); err != nil {
					return nil, err
				}

				continue
			}
		case bytecode.OpLoadIndexConst, bytecode.OpLoadIndexOptionalConst:
			src := reg[src1]
			optional := op == bytecode.OpLoadIndexOptionalConst
			arg := constants[src2.Constant()]

			if err := exec.loadIndexAndSet(ctx, dst, src, arg, optional); err != nil {
				if err := exec.errors.protected(err); err != nil {
					return nil, err
				}

				continue
			}
		case bytecode.OpLoadKey, bytecode.OpLoadKeyOptional:
			src := reg[src1]
			optional := op == bytecode.OpLoadKeyOptional
			arg := reg[src2]

			if err := exec.loadKeyAndSet(ctx, dst, exec.pc-1, src, arg, optional); err != nil {
				if err := exec.errors.protected(err); err != nil {
					return nil, err
				}

				continue
			}
		case bytecode.OpLoadKeyConst, bytecode.OpLoadKeyOptionalConst:
			src := reg[src1]
			optional := op == bytecode.OpLoadKeyOptionalConst
			arg := constants[src2.Constant()]

			if err := exec.loadKeyConstAndSet(ctx, dst, exec.pc-1, inst, src, arg, optional); err != nil {
				if err := exec.errors.protected(err); err != nil {
					return nil, err
				}

				continue
			}
		case bytecode.OpLoadPropertyConst, bytecode.OpLoadPropertyOptionalConst:
			src := reg[src1]
			optional := op == bytecode.OpLoadPropertyOptionalConst
			prop := constants[src2.Constant()]

			if err := exec.loadPropertyConstAndSet(ctx, dst, exec.pc-1, inst, src, prop, optional); err != nil {
				if err := exec.errors.protected(err); err != nil {
					return nil, err
				}

				continue
			}
		case bytecode.OpLoadProperty, bytecode.OpLoadPropertyOptional:
			src := reg[src1]
			optional := op == bytecode.OpLoadPropertyOptional
			prop := reg[src2]

			if err := exec.loadPropertyAndSet(ctx, dst, exec.pc-1, src, prop, optional); err != nil {
				if err := exec.errors.protected(err); err != nil {
					return nil, err
				}

				continue
			}
		case bytecode.OpAdd:
			reg[dst] = runtime.Add(ctx, reg[src1], reg[src2])
		case bytecode.OpAddConst:
			reg[dst] = runtime.Add(ctx, reg[src1], constants[src2.Constant()])
		case bytecode.OpConcat:
			concatStrings(reg, dst, src1, src2)
		case bytecode.OpSub:
			reg[dst] = runtime.Subtract(ctx, reg[src1], reg[src2])
		case bytecode.OpMulti:
			reg[dst] = runtime.Multiply(ctx, reg[src1], reg[src2])
		case bytecode.OpDiv:
			if err := exec.checkDivisionByZero(ctx, reg[src1], reg[src2]); err != nil {
				if err := exec.errors.protected(err); err != nil {
					return nil, err
				}

				continue
			}
			reg[dst] = runtime.Divide(ctx, reg[src1], reg[src2])
		case bytecode.OpMod:
			if err := exec.checkModuloByZero(ctx, reg[src2]); err != nil {
				if err := exec.errors.protected(err); err != nil {
					return nil, err
				}

				continue
			}
			reg[dst] = runtime.Modulus(ctx, reg[src1], reg[src2])
		case bytecode.OpIncr:
			reg[dst] = runtime.Increment(ctx, reg[dst])
		case bytecode.OpDecr:
			reg[dst] = runtime.Decrement(ctx, reg[dst])
		case bytecode.OpCastBool:
			reg[dst] = runtime.ToBoolean(reg[src1])
		case bytecode.OpNegate:
			reg[dst] = operators.Negate(reg[src1])
		case bytecode.OpFlipPositive:
			reg[dst] = operators.Positive(reg[src1])
		case bytecode.OpFlipNegative:
			reg[dst] = operators.Negative(reg[src1])
		case bytecode.OpCmp:
			reg[dst] = operators.Compare(ctx, reg[src1], reg[src2])
		case bytecode.OpNot:
			reg[dst] = !runtime.ToBoolean(reg[src1])
		case bytecode.OpEq:
			reg[dst] = operators.Equals(ctx, reg[src1], reg[src2])
		case bytecode.OpNe:
			reg[dst] = operators.NotEquals(ctx, reg[src1], reg[src2])
		case bytecode.OpGt:
			reg[dst] = operators.GreaterThan(ctx, reg[src1], reg[src2])
		case bytecode.OpLt:
			reg[dst] = operators.LessThan(ctx, reg[src1], reg[src2])
		case bytecode.OpGte:
			reg[dst] = operators.GreaterThanOrEqual(ctx, reg[src1], reg[src2])
		case bytecode.OpLte:
			reg[dst] = operators.LessThanOrEqual(ctx, reg[src1], reg[src2])
		case bytecode.OpIn:
			reg[dst] = operators.Contains(ctx, reg[src2], reg[src1])
		case bytecode.OpAllEq, bytecode.OpAllNe, bytecode.OpAllGt, bytecode.OpAllGte, bytecode.OpAllLt, bytecode.OpAllLte, bytecode.OpAllIn:
			cmp := operators.ComparatorFromByte(int(op) - int(bytecode.OpAllEq))
			res, err := operators.ArrayAll(ctx, cmp, reg[src1], reg[src2])

			if err := exec.errors.setOrCatch(dst, res, err); err != nil {
				return nil, err
			}
		case bytecode.OpAnyEq, bytecode.OpAnyNe, bytecode.OpAnyGt, bytecode.OpAnyGte, bytecode.OpAnyLt, bytecode.OpAnyLte, bytecode.OpAnyIn:
			cmp := operators.ComparatorFromByte(int(op) - int(bytecode.OpAnyEq))
			res, err := operators.ArrayAny(ctx, cmp, reg[src1], reg[src2])

			if err := exec.errors.setOrCatch(dst, res, err); err != nil {
				return nil, err
			}
		case bytecode.OpNoneEq, bytecode.OpNoneNe, bytecode.OpNoneGt, bytecode.OpNoneGte, bytecode.OpNoneLt, bytecode.OpNoneLte, bytecode.OpNoneIn:
			cmp := operators.ComparatorFromByte(int(op) - int(bytecode.OpNoneEq))
			res, err := operators.ArrayNone(ctx, cmp, reg[src1], reg[src2])

			if err := exec.errors.setOrCatch(dst, res, err); err != nil {
				return nil, err
			}
		case bytecode.OpHCall, bytecode.OpProtectedHCall:
			if err := exec.execHostCall(ctx, op, exec.pc-1, dst, src1, src2); err != nil {
				return nil, err
			}
		case bytecode.OpCall, bytecode.OpProtectedCall, bytecode.OpTailCall:
			if err := exec.execUdfCall(op, dst, src1, src2); err != nil {
				return nil, err
			}
			continue
		case bytecode.OpDataSet, bytecode.OpDataSetCollector, bytecode.OpDataSetSorter, bytecode.OpDataSetMultiSorter,
			bytecode.OpPush, bytecode.OpPushKV, bytecode.OpArrayPush, bytecode.OpObjectSet, bytecode.OpObjectSetConst:
			if err := exec.execDatasetOps(ctx, op, inst, dst, src1, src2); err != nil {
				return nil, err
			}
		case bytecode.OpIter, bytecode.OpIterNext, bytecode.OpIterValue, bytecode.OpIterKey, bytecode.OpIterLimit, bytecode.OpIterSkip:
			if err := exec.execIterOps(ctx, op, dst, src1, src2); err != nil {
				return nil, err
			}
		case bytecode.OpRand:
			reg[dst] = runtime.NewFloat(runtime.RandomDefault())
		case bytecode.OpReturn:
			retVal := reg[dst]

			if exec.returnToCaller(retVal) {
				continue
			}

			reg[bytecode.NoopOperand] = retVal
			break loop
		default:
			if handled, err := exec.execColdOps(ctx, op, dst, src1, src2); err != nil {
				return nil, err
			} else if handled {
				continue
			}

			return nil, runtime.Errorf(runtime.ErrUnexpected, "unknown opcode %d at pc %d", op, exec.pc-1)
		}
	}

	return exec.registers.Values[bytecode.NoopOperand], nil
}

func (exec *execState) unwindToProtected() bool {
	registers, pc, ok := exec.vm.frames.UnwindToProtectedFrame(exec.registers.Values)
	if !ok {
		return false
	}

	exec.registers.Values = registers
	exec.pc = pc
	return true
}

func (exec *execState) returnToCaller(retVal runtime.Value) bool {
	registers, pc, ok := exec.vm.frames.ReturnToCaller(exec.registers.Values, retVal)
	if !ok {
		return false
	}

	exec.registers.Values = registers
	exec.pc = pc
	return true
}
