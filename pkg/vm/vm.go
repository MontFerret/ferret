package vm

import (
	"context"
	"io"

	"github.com/MontFerret/ferret/v2/pkg/bytecode"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
	"github.com/MontFerret/ferret/v2/pkg/vm/internal/data"
	"github.com/MontFerret/ferret/v2/pkg/vm/internal/frame"
	"github.com/MontFerret/ferret/v2/pkg/vm/internal/mem"
)

type VM struct {
	options      options
	registers    *mem.RegisterFile
	cache        *mem.Cache
	scratch      *mem.Scratch
	env          *Environment
	program      *bytecode.Program
	instructions []data.ExecInstruction
	catchByPC    []int
	pc           int
	frames       frame.CallStack
}

func New(program *bytecode.Program) *VM {
	return NewWith(program)
}

func NewWith(program *bytecode.Program, opts ...Option) *VM {
	o := newOptions(opts)

	vm := &VM{
		registers:    mem.NewRegisterFile(program.Registers),
		cache:        mem.NewCache(len(program.Bytecode), o.shapeCacheLimit),
		scratch:      mem.NewScratch(len(program.Params)),
		program:      program,
		options:      o,
		instructions: buildExecInstructions(program.Bytecode),
		catchByPC:    buildCatchByPC(len(program.Bytecode), program.CatchTable),
	}

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
	defer func() {
		if r := recover(); r != nil {
			err = vm.runtimeErrorFromPanic(r)
			result = nil

			return
		}

		if err != nil {
			err = vm.wrapRuntimeError(err)
		}
	}()

	return vm.runCore(ctx, env)
}

func (vm *VM) runUnchecked(ctx context.Context, env *Environment) (runtime.Value, error) {
	result, err := vm.runCore(ctx, env)

	if err != nil {
		return nil, vm.wrapRuntimeError(err)
	}

	return result, nil
}

func (vm *VM) runCore(ctx context.Context, env *Environment) (runtime.Value, error) {
	if env == nil {
		env = noopEnv
	}

	if err := validate(env, vm.program); err != nil {
		return nil, err
	}

	if err := vm.warmup(env); err != nil {
		return nil, err
	}

	if vm.registers.IsDirty() {
		vm.registers.Reset()
	}

	vm.registers.MarkDirty()
	vm.env = env
	vm.pc = 0
	vm.frames.Reset()

	instructions := vm.instructions
	constants := vm.program.Constants
	aggregatePlans := vm.program.Metadata.AggregatePlans
	shapeCache := vm.cache.ShapeCache
	paramSlots := vm.scratch.Params
loop:
	for vm.pc < len(instructions) {
		inst := &instructions[vm.pc]
		op := inst.Opcode
		dst, src1, src2 := inst.Operands[0], inst.Operands[1], inst.Operands[2]
		reg := vm.registers.Values
		vm.pc++

		switch op {
		case bytecode.OpReturn:
			retVal := reg[dst]

			if vm.returnToCaller(retVal) {
				continue
			}

			reg[bytecode.NoopOperand] = retVal

			break loop
		case bytecode.OpJump:
			vm.pc = int(dst)
		case bytecode.OpJumpIfFalse:
			if !coerceBool(reg[src1]) {
				vm.pc = int(dst)
			}
		case bytecode.OpJumpIfTrue:
			if coerceBool(reg[src1]) {
				vm.pc = int(dst)
			}
		case bytecode.OpJumpIfNone:
			if reg[src1] == runtime.None {
				vm.pc = int(dst)
			}
		case bytecode.OpJumpIfNe:
			if ne(ctx, reg[src1], reg[src2]) {
				vm.pc = int(dst)
			}
		case bytecode.OpJumpIfNeConst:
			if ne(ctx, reg[src1], constants[src2.Constant()]) {
				vm.pc = int(dst)
			}
		case bytecode.OpJumpIfEq:
			if eq(ctx, reg[src1], reg[src2]) {
				vm.pc = int(dst)
			}
		case bytecode.OpJumpIfEqConst:
			if eq(ctx, reg[src1], constants[src2.Constant()]) {
				vm.pc = int(dst)
			}
		case bytecode.OpJumpIfMissingProperty:
			obj, ok := reg[src1].(runtime.Map)
			if !ok {
				vm.pc = int(dst)
				continue
			}

			key, ok := reg[src2].(runtime.String)
			if !ok {
				vm.pc = int(dst)
				continue
			}

			has, err := obj.ContainsKey(ctx, key)
			if err != nil && vm.handleProtectedError(err) != nil {
				return nil, err
			}

			if !has {
				vm.pc = int(dst)
			}
		case bytecode.OpJumpIfMissingPropertyConst:
			obj, ok := reg[src1].(runtime.Map)
			if !ok {
				vm.pc = int(dst)

				continue
			}

			key, ok := constants[src2.Constant()].(runtime.String)
			if !ok {
				vm.pc = int(dst)
				continue
			}

			has, err := obj.ContainsKey(ctx, key)
			if err != nil && vm.handleProtectedError(err) != nil {
				return nil, err
			}

			if !has {
				vm.pc = int(dst)
			}
		case bytecode.OpFail:
			if !dst.IsConstant() {
				if err := vm.handleError(runtime.Error(runtime.ErrInvalidOperation, "FAIL expects a constant string message")); err != nil {
					return nil, err
				}

				continue
			}

			idx := dst.Constant()
			if idx < 0 || idx >= len(constants) {
				if err := vm.handleError(runtime.Error(runtime.ErrInvalidOperation, "FAIL expects a valid constant string message")); err != nil {
					return nil, err
				}

				continue
			}

			msg, ok := constants[idx].(runtime.String)
			if !ok {
				if err := vm.handleError(runtime.TypeErrorOf(constants[idx], runtime.TypeString)); err != nil {
					return nil, err
				}

				continue
			}

			if err := vm.handleError(runtime.Error(runtime.ErrInvalidOperation, msg.String())); err != nil {
				return nil, err
			}

			continue
		case bytecode.OpHCall, bytecode.OpProtectedHCall:
			cacheFn := vm.cache.HostFunctions[vm.pc-1]
			out, err := callCachedHostFunction(ctx, cacheFn, vm.registers.Values, src1, src2)

			if err := vm.setCallResult(op, dst, out, err); err != nil {
				if vm.unwindToProtected() {
					continue
				}

				return nil, err
			}
		case bytecode.OpCall, bytecode.OpProtectedCall:
			if err := vm.callUdf(op, dst, src1, src2); err != nil {
				if err := vm.setCallResult(op, dst, runtime.None, err); err != nil {
					if vm.unwindToProtected() {
						continue
					}

					return nil, err
				}
			}
		case bytecode.OpTailCall:
			if err := vm.tailCallUdf(dst, src1, src2); err != nil {
				if vm.unwindToProtected() {
					continue
				}

				return nil, err
			}
		case bytecode.OpDispatch:
			dispatcher, eventName, payload, options, err := vm.castDispatchArgs(ctx, reg[dst], reg[src1], reg[src2])

			if err != nil {
				if err := vm.setOrTryCatch(dst, runtime.None, err); err != nil {
					return nil, err
				}

				continue
			}

			out, err := dispatcher.Dispatch(ctx, runtime.DispatchEvent{
				Name:    eventName,
				Payload: payload,
				Options: options,
			})

			if out == nil {
				out = runtime.None
			}

			if err := vm.setOrTryCatch(dst, out, err); err != nil {
				return nil, err
			}
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
			reg[dst] = data.NewFastObjectOf(shapeCache, vm.options.fastObjectDictThreshold, int(src1))

		case bytecode.OpExists:
			val := reg[src1]

			if val == runtime.None {
				reg[dst] = runtime.False
				continue
			}

			if measurable, ok := val.(runtime.Measurable); ok {
				length, err := measurable.Length(ctx)

				if err != nil {
					if err := vm.handleErrorWithCatch(err, func() {
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
					if err := vm.handleErrorWithCatch(err, func() {
						length = 0
					}); err != nil {
						return nil, err
					}
				}

				reg[dst] = length
				continue
			}

			if err := vm.handleErrorWithCatch(runtime.TypeErrorOf(reg[src1],
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
		case bytecode.OpClose:
			val, ok := reg[dst].(io.Closer)
			reg[dst] = runtime.None

			if ok {
				closeErr := val.Close()

				if closeErr != nil {
					if err := vm.handleError(closeErr); err != nil {
						return nil, err
					}

					continue
				}
			}
		case bytecode.OpLoadRange:
			res, err := ToRange(ctx, reg[src1], reg[src2])

			if err == nil {
				reg[dst] = res
				break
			}

			if err := vm.handleProtectedError(err); err != nil {
				return nil, err
			}

			continue
		case bytecode.OpLoadIndex, bytecode.OpLoadIndexOptional:
			src := reg[src1]
			optional := op == bytecode.OpLoadIndexOptional
			arg := reg[src2]

			if err := vm.loadIndexAndSet(ctx, dst, src, arg, optional); err != nil {
				if err := vm.handleProtectedError(err); err != nil {
					return nil, err
				}

				continue
			}
		case bytecode.OpLoadIndexConst, bytecode.OpLoadIndexOptionalConst:
			src := reg[src1]
			optional := op == bytecode.OpLoadIndexOptionalConst
			arg := constants[src2.Constant()]

			if err := vm.loadIndexAndSet(ctx, dst, src, arg, optional); err != nil {
				if err := vm.handleProtectedError(err); err != nil {
					return nil, err
				}

				continue
			}
		case bytecode.OpLoadKey, bytecode.OpLoadKeyOptional:
			src := reg[src1]
			optional := op == bytecode.OpLoadKeyOptional
			arg := reg[src2]

			if err := vm.loadKeyAndSet(ctx, dst, vm.pc-1, src, arg, optional); err != nil {
				if err := vm.handleProtectedError(err); err != nil {
					return nil, err
				}

				continue
			}
		case bytecode.OpLoadKeyConst, bytecode.OpLoadKeyOptionalConst:
			src := reg[src1]
			optional := op == bytecode.OpLoadKeyOptionalConst
			arg := constants[src2.Constant()]

			if err := vm.loadKeyConstAndSet(ctx, dst, vm.pc-1, inst, src, arg, optional); err != nil {
				if err := vm.handleProtectedError(err); err != nil {
					return nil, err
				}

				continue
			}
		case bytecode.OpLoadPropertyConst, bytecode.OpLoadPropertyOptionalConst:
			src := reg[src1]
			optional := op == bytecode.OpLoadPropertyOptionalConst
			prop := constants[src2.Constant()]

			if err := vm.loadPropertyConstAndSet(ctx, dst, vm.pc-1, inst, src, prop, optional); err != nil {
				if err := vm.handleProtectedError(err); err != nil {
					return nil, err
				}

				continue
			}
		case bytecode.OpLoadProperty, bytecode.OpLoadPropertyOptional:
			src := reg[src1]
			optional := op == bytecode.OpLoadPropertyOptional
			prop := reg[src2]

			if err := vm.loadPropertyAndSet(ctx, dst, vm.pc-1, src, prop, optional); err != nil {
				if err := vm.handleProtectedError(err); err != nil {
					return nil, err
				}

				continue
			}
		case bytecode.OpApplyQuery:
			src := readOperandValue(reg, constants, src1)
			descriptor := readOperandValue(reg, constants, src2)
			out, err := ApplyQuery(ctx, src, descriptor)

			if err := vm.setOrTryCatch(dst, out, err); err != nil {
				if err := vm.handleProtectedError(err); err != nil {
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
			if err := vm.checkDivisionByZero(ctx, reg[src1], reg[src2]); err != nil {
				if err := vm.handleProtectedError(err); err != nil {
					return nil, err
				}

				continue
			}
			reg[dst] = runtime.Divide(ctx, reg[src1], reg[src2])
		case bytecode.OpMod:
			if err := vm.checkModuloByZero(ctx, reg[src2]); err != nil {
				if err := vm.handleProtectedError(err); err != nil {
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
			reg[dst] = coerceBool(reg[src1])
		case bytecode.OpNegate:
			reg[dst] = Negate(reg[src1])
		case bytecode.OpFlipPositive:
			reg[dst] = Positive(reg[src1])
		case bytecode.OpFlipNegative:
			reg[dst] = Negative(reg[src1])
		case bytecode.OpCmp:
			reg[dst] = cmp(ctx, reg[src1], reg[src2])
		case bytecode.OpNot:
			reg[dst] = !coerceBool(reg[src1])
		case bytecode.OpEq:
			reg[dst] = eq(ctx, reg[src1], reg[src2])
		case bytecode.OpNe:
			reg[dst] = ne(ctx, reg[src1], reg[src2])
		case bytecode.OpGt:
			reg[dst] = gt(ctx, reg[src1], reg[src2])
		case bytecode.OpLt:
			reg[dst] = lt(ctx, reg[src1], reg[src2])
		case bytecode.OpGte:
			reg[dst] = gte(ctx, reg[src1], reg[src2])
		case bytecode.OpLte:
			reg[dst] = lte(ctx, reg[src1], reg[src2])
		case bytecode.OpIn:
			reg[dst] = contains(ctx, reg[src2], reg[src1])
		case bytecode.OpLike:
			res, err := Like(reg[src1], reg[src2])

			if err == nil {
				reg[dst] = res
				break
			}

			if err := vm.handleProtectedError(err); err != nil {
				return nil, err
			}

			continue
		case bytecode.OpRegexp:
			r, err := vm.regexpCached(vm.pc-1, reg[src2])

			if err == nil {
				reg[dst] = r.Match(reg[src1])
			} else {
				if err := vm.handleErrorWithCatch(err, func() {
					reg[dst] = runtime.False
				}); err != nil {
					return nil, err
				}

				continue
			}
		case bytecode.OpAllEq, bytecode.OpAllNe, bytecode.OpAllGt, bytecode.OpAllGte, bytecode.OpAllLt, bytecode.OpAllLte, bytecode.OpAllIn:
			cmp := ComparatorFromByte(int(op) - int(bytecode.OpAllEq))
			res, err := ArrayAll(ctx, cmp, reg[src1], reg[src2])

			if err := vm.setOrTryCatch(dst, res, err); err != nil {
				return nil, err
			}
		case bytecode.OpAnyEq, bytecode.OpAnyNe, bytecode.OpAnyGt, bytecode.OpAnyGte, bytecode.OpAnyLt, bytecode.OpAnyLte, bytecode.OpAnyIn:
			cmp := ComparatorFromByte(int(op) - int(bytecode.OpAnyEq))
			res, err := ArrayAny(ctx, cmp, reg[src1], reg[src2])

			if err := vm.setOrTryCatch(dst, res, err); err != nil {
				return nil, err
			}
		case bytecode.OpNoneEq, bytecode.OpNoneNe, bytecode.OpNoneGt, bytecode.OpNoneGte, bytecode.OpNoneLt, bytecode.OpNoneLte, bytecode.OpNoneIn:
			cmp := ComparatorFromByte(int(op) - int(bytecode.OpNoneEq))
			res, err := ArrayNone(ctx, cmp, reg[src1], reg[src2])

			if err := vm.setOrTryCatch(dst, res, err); err != nil {
				return nil, err
			}
		case bytecode.OpFlatten:
			depth := src2.Register()

			if depth < 1 {
				depth = 1
			}

			res, err := Flatten(ctx, reg[src1], depth)

			if err := vm.setOrTryCatch(dst, res, err); err != nil {
				return nil, err
			}
		case bytecode.OpDataSet, bytecode.OpDataSetCollector, bytecode.OpDataSetSorter, bytecode.OpDataSetMultiSorter,
			bytecode.OpPush, bytecode.OpPushKV, bytecode.OpArrayPush, bytecode.OpObjectSet, bytecode.OpObjectSetConst:
			if err := vm.execDatasetOps(ctx, op, inst, dst, src1, src2, reg, constants, aggregatePlans); err != nil {
				return nil, err
			}
		case bytecode.OpIter, bytecode.OpIterNext, bytecode.OpIterValue, bytecode.OpIterKey, bytecode.OpIterLimit, bytecode.OpIterSkip:
			if err := vm.execIterOps(ctx, op, dst, src1, src2, reg); err != nil {
				return nil, err
			}
		case bytecode.OpStream:
			if err := vm.execStreamOp(ctx, dst, src1, src2, reg); err != nil {
				return nil, err
			}
		case bytecode.OpStreamIter:
			if err := vm.execStreamIterOp(ctx, dst, src1, src2, reg); err != nil {
				return nil, err
			}
		case bytecode.OpSleep:
			dur, err := runtime.ToInt(ctx, reg[dst])

			if err != nil {
				if err := vm.handleError(err); err != nil {
					return nil, err
				}

				continue
			}

			if err := data.Sleep(ctx, dur); err != nil {
				if err := vm.handleProtectedError(err); err != nil {
					return nil, err
				}

				continue
			}
		case bytecode.OpRand:
			reg[dst] = runtime.NewFloat(runtime.RandomDefault())

		default:
			return nil, runtime.Errorf(runtime.ErrUnexpected, "unknown opcode %d at pc %d", op, vm.pc-1)
		}
	}

	return vm.registers.Values[bytecode.NoopOperand], nil
}

func (vm *VM) unwindToProtected() bool {
	registers, pc, ok := vm.frames.UnwindToProtectedFrame(vm.registers.Values)
	if !ok {
		return false
	}

	vm.registers.Values = registers
	vm.pc = pc
	return true
}

func (vm *VM) returnToCaller(retVal runtime.Value) bool {
	registers, pc, ok := vm.frames.ReturnToCaller(vm.registers.Values, retVal)
	if !ok {
		return false
	}

	vm.registers.Values = registers
	vm.pc = pc
	return true
}
