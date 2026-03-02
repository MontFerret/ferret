package vm

import (
	"context"
	"io"

	"github.com/MontFerret/ferret/v2/pkg/bytecode"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
	"github.com/MontFerret/ferret/v2/pkg/vm/internal/data"
	"github.com/MontFerret/ferret/v2/pkg/vm/internal/mem"
	"github.com/MontFerret/ferret/v2/pkg/vm/internal/operators"
)

type VM struct {
	registers               *mem.RegisterFile
	cache                   *mem.Cache
	env                     *Environment
	program                 *bytecode.Program
	fastObjectDictThreshold int
	instructions            []data.ExecInstruction
	catchByPC               []int
	pc                      int
	frames                  []callFrame
	regPool                 regPool
}

func New(program *bytecode.Program) *VM {
	return NewWithOptions(program)
}

func NewWithOptions(program *bytecode.Program, opts ...VMOption) *VM {
	cfg := defaultVMConfig()

	for _, opt := range opts {
		opt(&cfg)
	}

	vm := new(VM)
	vm.program = program
	vm.registers = mem.NewRegisterFile(program.Registers)
	vm.cache = mem.NewCache(len(program.Bytecode), cfg.shapeCacheLimit)
	vm.fastObjectDictThreshold = cfg.fastObjectDictThreshold
	vm.instructions = make([]data.ExecInstruction, len(program.Bytecode))
	maxUdfRegs := 0
	for i := range program.Functions.UserDefined {
		if program.Functions.UserDefined[i].Registers > maxUdfRegs {
			maxUdfRegs = program.Functions.UserDefined[i].Registers
		}
	}
	vm.regPool.init(maxUdfRegs)

	for i := range program.Bytecode {
		vm.instructions[i] = data.ExecInstruction{
			Instruction: program.Bytecode[i],
		}
	}

	if bytecodeLen := len(program.Bytecode); bytecodeLen > 0 {
		vm.catchByPC = make([]int, bytecodeLen)
		for i := range vm.catchByPC {
			vm.catchByPC[i] = -1
		}
		for i, pair := range program.CatchTable {
			start, end := pair[0], pair[1]
			if start < 0 {
				start = 0
			}
			if end >= bytecodeLen {
				end = bytecodeLen - 1
			}
			for pc := start; pc <= end; pc++ {
				if vm.catchByPC[pc] == -1 {
					vm.catchByPC[pc] = i
				}
			}
		}
	}

	return vm
}

func (vm *VM) Run(ctx context.Context, env *Environment) (result runtime.Value, err error) {
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
	vm.frames = vm.frames[:0]

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
		case bytecode.OpExists:
			val := reg[src1]

			if val == runtime.None {
				reg[dst] = runtime.False
				continue
			}

			if measurable, ok := val.(runtime.Measurable); ok {
				length, err := measurable.Length(ctx)

				if err != nil {
					if _, catch := vm.tryCatch(vm.pc); catch {
						reg[dst] = runtime.False

						continue
					}

					if vm.unwindToProtected() {
						continue
					}

					return nil, err
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
					if _, catch := vm.tryCatch(vm.pc); catch {
						length = 0
					} else {
						if vm.unwindToProtected() {
							continue
						}

						return nil, err
					}
				}

				reg[dst] = length
				continue
			}

			if _, catch := vm.tryCatch(vm.pc); catch {
				reg[dst] = runtime.ZeroInt
				continue
			}

			if vm.unwindToProtected() {
				continue
			}

			return runtime.None, runtime.TypeErrorOf(reg[src1],
				runtime.TypeString,
				runtime.TypeList,
				runtime.TypeMap,
				runtime.TypeBinary,
				runtime.TypeMeasurable,
			)
		case bytecode.OpType:
			reg[dst] = runtime.NewString(runtime.TypeName(runtime.TypeOf(reg[src1])))
		case bytecode.OpClose:
			val, ok := reg[dst].(io.Closer)
			reg[dst] = runtime.None

			if ok {
				closeErr := val.Close()

				if closeErr != nil {
					if _, catch := vm.tryCatch(vm.pc); !catch {
						if vm.unwindToProtected() {
							continue
						}

						return nil, closeErr
					}
				}
			}
		case bytecode.OpLoadRange:
			res, err := operators.ToRange(ctx, reg[src1], reg[src2])

			if err == nil {
				reg[dst] = res
				break
			}

			if vm.unwindToProtected() {
				continue
			}

			return nil, err
		case bytecode.OpLoadIndex, bytecode.OpLoadIndexOptional:
			src := reg[src1]
			optional := op == bytecode.OpLoadIndexOptional
			arg := reg[src2]

			if err := vm.loadIndexAndSet(ctx, dst, src, arg, optional); err != nil {
				if vm.unwindToProtected() {
					continue
				}

				return nil, err
			}
		case bytecode.OpLoadIndexConst, bytecode.OpLoadIndexOptionalConst:
			src := reg[src1]
			optional := op == bytecode.OpLoadIndexOptionalConst
			arg := constants[src2.Constant()]

			if err := vm.loadIndexAndSet(ctx, dst, src, arg, optional); err != nil {
				if vm.unwindToProtected() {
					continue
				}

				return nil, err
			}
		case bytecode.OpLoadKey, bytecode.OpLoadKeyOptional:
			src := reg[src1]
			optional := op == bytecode.OpLoadKeyOptional
			arg := reg[src2]

			if err := vm.loadKeyAndSet(ctx, dst, vm.pc-1, src, arg, optional); err != nil {
				if vm.unwindToProtected() {
					continue
				}

				return nil, err
			}
		case bytecode.OpLoadKeyConst, bytecode.OpLoadKeyOptionalConst:
			src := reg[src1]
			optional := op == bytecode.OpLoadKeyOptionalConst
			arg := constants[src2.Constant()]

			if err := vm.loadKeyConstAndSet(ctx, dst, vm.pc-1, inst, src, arg, optional); err != nil {
				if vm.unwindToProtected() {
					continue
				}

				return nil, err
			}
		case bytecode.OpLoadPropertyConst, bytecode.OpLoadPropertyOptionalConst:
			src := reg[src1]
			optional := op == bytecode.OpLoadPropertyOptionalConst
			prop := constants[src2.Constant()]

			if err := vm.loadPropertyConstAndSet(ctx, dst, vm.pc-1, inst, src, prop, optional); err != nil {
				if vm.unwindToProtected() {
					continue
				}

				return nil, err
			}
		case bytecode.OpLoadProperty, bytecode.OpLoadPropertyOptional:
			src := reg[src1]
			optional := op == bytecode.OpLoadPropertyOptional
			prop := reg[src2]

			if err := vm.loadPropertyAndSet(ctx, dst, vm.pc-1, src, prop, optional); err != nil {
				if vm.unwindToProtected() {
					continue
				}

				return nil, err
			}
		case bytecode.OpApplyQuery:
			if err := vm.applyQuery(ctx, reg, src1, constants, src2, dst); err != nil {
				if vm.unwindToProtected() {
					continue
				}

				return nil, err
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
				if vm.unwindToProtected() {
					continue
				}

				return nil, err
			}
			reg[dst] = runtime.Divide(ctx, reg[src1], reg[src2])
		case bytecode.OpMod:
			if err := vm.checkModuloByZero(ctx, reg[src2]); err != nil {
				if vm.unwindToProtected() {
					continue
				}

				return nil, err
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

			if vm.unwindToProtected() {
				continue
			}

			return nil, err
		case bytecode.OpRegexp:
			r, err := vm.regexpCached(vm.pc-1, reg[src2])

			if err == nil {
				reg[dst] = r.Match(reg[src1])
			} else if _, catch := vm.tryCatch(vm.pc); catch {
				reg[dst] = runtime.False
			} else {
				if vm.unwindToProtected() {
					continue
				}

				return nil, err
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
		case bytecode.OpJumpIfNe, bytecode.OpJumpIfNeConst, bytecode.OpJumpIfEq, bytecode.OpJumpIfEqConst,
			bytecode.OpJumpIfMissingProperty, bytecode.OpJumpIfMissingPropertyConst,
			bytecode.OpReturn:
			done, err := vm.execControlOps(ctx, op, dst, src1, src2, reg, constants)
			if err != nil {
				return nil, err
			}
			if done {
				break loop
			}
		case bytecode.OpHCall, bytecode.OpProtectedHCall,
			bytecode.OpHCall0, bytecode.OpProtectedHCall0,
			bytecode.OpHCall1, bytecode.OpProtectedHCall1,
			bytecode.OpHCall2, bytecode.OpProtectedHCall2,
			bytecode.OpHCall3, bytecode.OpProtectedHCall3,
			bytecode.OpHCall4, bytecode.OpProtectedHCall4:
			if err := vm.execHostCall(ctx, op, vm.pc-1, dst, src1, src2); err != nil {
				return nil, err
			}
		case bytecode.OpCall, bytecode.OpProtectedCall,
			bytecode.OpCall0, bytecode.OpProtectedCall0,
			bytecode.OpCall1, bytecode.OpProtectedCall1,
			bytecode.OpCall2, bytecode.OpProtectedCall2,
			bytecode.OpCall3, bytecode.OpProtectedCall3,
			bytecode.OpCall4, bytecode.OpProtectedCall4,
			bytecode.OpTailCall, bytecode.OpTailCall0, bytecode.OpTailCall1, bytecode.OpTailCall2, bytecode.OpTailCall3, bytecode.OpTailCall4:
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
		case bytecode.OpSleep:
			dur, err := runtime.ToInt(ctx, reg[dst])

			if err != nil {
				if _, catch := vm.tryCatch(vm.pc); catch {
					continue
				}

				if vm.unwindToProtected() {
					continue
				}

				return nil, err
			}

			if err := data.Sleep(ctx, dur); err != nil {
				if vm.unwindToProtected() {
					continue
				}

				return nil, err
			}
		case bytecode.OpRand:
			reg[dst] = runtime.NewFloat(runtime.RandomDefault())
		default:
			// TODO: Return an error or ignore unknown opcodes?
			continue
		}
	}

	return vm.registers.Values[bytecode.NoopOperand], nil
}
