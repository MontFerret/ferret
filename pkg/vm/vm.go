package vm

import (
	"context"
	"io"

	"github.com/MontFerret/ferret/v2/pkg/bytecode"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
	"github.com/MontFerret/ferret/v2/pkg/vm/internal/data"
	"github.com/MontFerret/ferret/v2/pkg/vm/internal/frame"
	"github.com/MontFerret/ferret/v2/pkg/vm/internal/mem"
	"github.com/MontFerret/ferret/v2/pkg/vm/internal/operators"
)

type VM struct {
	registers               *mem.RegisterFile
	cache                   *mem.Cache
	env                     *Environment
	program                 *bytecode.Program
	runSafetyMode           RunSafetyMode
	fastObjectDictThreshold int
	instructions            []data.ExecInstruction
	catchByPC               []int
	pc                      int
	frames                  frame.CallStack
}

func New(program *bytecode.Program) *VM {
	return NewWithOptions(program)
}

func NewWithOptions(program *bytecode.Program, opts ...Option) *VM {
	cfg := newOptions(opts)

	vm := &VM{
		registers:               mem.NewRegisterFile(program.Registers),
		cache:                   mem.NewCache(len(program.Bytecode), cfg.shapeCacheLimit),
		program:                 program,
		runSafetyMode:           cfg.runSafetyMode,
		fastObjectDictThreshold: cfg.fastObjectDictThreshold,
		instructions:            buildExecInstructions(program.Bytecode),
		catchByPC:               buildCatchByPC(len(program.Bytecode), program.CatchTable),
	}

	vm.frames.Init(maxUDFRegisters(program.Functions.UserDefined))

	return vm
}

func buildExecInstructions(code []bytecode.Instruction) []data.ExecInstruction {
	instructions := make([]data.ExecInstruction, len(code))

	for i := range code {
		instructions[i] = data.ExecInstruction{
			Instruction: code[i],
		}
	}

	return instructions
}

func maxUDFRegisters(udfs []bytecode.UDF) int {
	maxUDFRegs := 0

	for i := range udfs {
		if udfs[i].Registers > maxUDFRegs {
			maxUDFRegs = udfs[i].Registers
		}
	}

	return maxUDFRegs
}

func buildCatchByPC(bytecodeLen int, catches []bytecode.Catch) []int {
	if bytecodeLen <= 0 {
		return nil
	}

	catchByPC := make([]int, bytecodeLen)
	for i := range catchByPC {
		catchByPC[i] = -1
	}

	for i, pair := range catches {
		start, end := pair[0], pair[1]
		if start < 0 {
			start = 0
		}
		if end >= bytecodeLen {
			end = bytecodeLen - 1
		}
		for pc := start; pc <= end; pc++ {
			if catchByPC[pc] == -1 {
				catchByPC[pc] = i
			}
		}
	}

	return catchByPC
}

func (vm *VM) Run(ctx context.Context, env *Environment) (runtime.Value, error) {
	switch vm.runSafetyMode {
	case RunSafetyFast:
		return vm.runFast(ctx, env)
	default:
		return vm.runStrict(ctx, env)
	}
}

func (vm *VM) runStrict(ctx context.Context, env *Environment) (result runtime.Value, err error) {
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

func (vm *VM) runFast(ctx context.Context, env *Environment) (runtime.Value, error) {
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
loop:
	for vm.pc < len(instructions) {
		reg := vm.registers.Values
		inst := &instructions[vm.pc]
		op := inst.Opcode
		dst, src1, src2 := inst.Operands[0], inst.Operands[1], inst.Operands[2]
		vm.pc++

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
			name := constants[src1.Constant()]
			reg[dst] = vm.env.Params[name.String()]
		case bytecode.OpLoadArray:
			reg[dst] = runtime.NewArray(int(src1))
		case bytecode.OpLoadObject:
			reg[dst] = data.NewFastObjectOf(shapeCache, vm.fastObjectDictThreshold, int(src1))
		case bytecode.OpJump:
			vm.pc = int(dst)
		case bytecode.OpJumpIfFalse:
			if !runtime.ToBoolean(reg[src1]) {
				vm.pc = int(dst)
			}
		case bytecode.OpJumpIfTrue:
			if runtime.ToBoolean(reg[src1]) {
				vm.pc = int(dst)
			}
		case bytecode.OpJumpIfNone:
			if reg[src1] == runtime.None {
				vm.pc = int(dst)
			}
		case bytecode.OpJumpIfNe:
			if operators.NotEquals(ctx, reg[src1], reg[src2]) {
				vm.pc = int(dst)
			}
		case bytecode.OpJumpIfNeConst:
			if operators.NotEquals(ctx, reg[src1], constants[src2.Constant()]) {
				vm.pc = int(dst)
			}
		case bytecode.OpJumpIfEq:
			if operators.Equals(ctx, reg[src1], reg[src2]) {
				vm.pc = int(dst)
			}
		case bytecode.OpJumpIfEqConst:
			if operators.Equals(ctx, reg[src1], constants[src2.Constant()]) {
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
			res, err := operators.ToRange(ctx, reg[src1], reg[src2])

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
			if err := vm.applyQuery(ctx, reg, src1, constants, src2, dst); err != nil {
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
			vm.concatStrings(reg, dst, src1, src2)
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
		case bytecode.OpLike:
			res, err := operators.Like(reg[src1], reg[src2])

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
			cmp := operators.ComparatorFromByte(int(op) - int(bytecode.OpAllEq))
			res, err := operators.ArrayAll(ctx, cmp, reg[src1], reg[src2])

			if err := vm.setOrTryCatch(dst, res, err); err != nil {
				return nil, err
			}
		case bytecode.OpAnyEq, bytecode.OpAnyNe, bytecode.OpAnyGt, bytecode.OpAnyGte, bytecode.OpAnyLt, bytecode.OpAnyLte, bytecode.OpAnyIn:
			cmp := operators.ComparatorFromByte(int(op) - int(bytecode.OpAnyEq))
			res, err := operators.ArrayAny(ctx, cmp, reg[src1], reg[src2])

			if err := vm.setOrTryCatch(dst, res, err); err != nil {
				return nil, err
			}
		case bytecode.OpNoneEq, bytecode.OpNoneNe, bytecode.OpNoneGt, bytecode.OpNoneGte, bytecode.OpNoneLt, bytecode.OpNoneLte, bytecode.OpNoneIn:
			cmp := operators.ComparatorFromByte(int(op) - int(bytecode.OpNoneEq))
			res, err := operators.ArrayNone(ctx, cmp, reg[src1], reg[src2])

			if err := vm.setOrTryCatch(dst, res, err); err != nil {
				return nil, err
			}
		case bytecode.OpFlatten:
			depth := src2.Register()

			if depth < 1 {
				depth = 1
			}

			res, err := operators.Flatten(ctx, reg[src1], depth)

			if err := vm.setOrTryCatch(dst, res, err); err != nil {
				return nil, err
			}
		case bytecode.OpHCall, bytecode.OpProtectedHCall:
			if err := vm.execHostCall(ctx, op, vm.pc-1, dst, src1, src2); err != nil {
				return nil, err
			}
		case bytecode.OpCall, bytecode.OpProtectedCall, bytecode.OpTailCall:
			if err := vm.execUdfCall(op, dst, src1, src2); err != nil {
				return nil, err
			}
			continue
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
		case bytecode.OpReturn:
			retVal := reg[dst]

			if vm.returnToCaller(retVal) {
				continue
			}

			reg[bytecode.NoopOperand] = retVal
			break loop
		default:
			return nil, runtime.Errorf(runtime.ErrUnexpected, "unknown opcode %d at pc %d", op, vm.pc-1)
		}
	}

	return vm.registers.Values[bytecode.NoopOperand], nil
}
