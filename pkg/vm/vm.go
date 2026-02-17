package vm

import (
	"context"
	"io"
	"strings"

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
	pc                      int
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

	for i := range program.Bytecode {
		vm.instructions[i] = data.ExecInstruction{
			Instruction: program.Bytecode[i],
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

	instructions := vm.instructions
	constants := vm.program.Constants
	aggregatePlans := vm.program.Metadata.AggregatePlans
	reg := vm.registers.Values
	shapeCache := vm.cache.ShapeCache
loop:
	for vm.pc < len(instructions) {
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
		case bytecode.OpAdd:
			reg[dst] = runtime.Add(ctx, reg[src1], reg[src2])
		case bytecode.OpAddConst:
			reg[dst] = runtime.Add(ctx, reg[src1], constants[src2.Constant()])
		case bytecode.OpConcat:
			start := int(src1)
			count := int(src2)

			if count <= 0 {
				reg[dst] = runtime.EmptyString
				break
			}

			if count == 1 {
				reg[dst] = runtime.ToString(reg[start])
				break
			}

			if count == 2 {
				s1 := runtime.ToString(reg[start])
				s2 := runtime.ToString(reg[start+1])

				if s1 == runtime.EmptyString {
					reg[dst] = s2
					break
				}

				if s2 == runtime.EmptyString {
					reg[dst] = s1
					break
				}

				reg[dst] = runtime.NewString(string(s1) + string(s2))
				break
			}

			if count == 3 {
				s1 := runtime.ToString(reg[start])
				s2 := runtime.ToString(reg[start+1])
				s3 := runtime.ToString(reg[start+2])

				if s1 == runtime.EmptyString {
					if s2 == runtime.EmptyString {
						reg[dst] = s3
						break
					}
					if s3 == runtime.EmptyString {
						reg[dst] = s2
						break
					}
				} else if s2 == runtime.EmptyString {
					if s3 == runtime.EmptyString {
						reg[dst] = s1
						break
					}
				}

				reg[dst] = runtime.NewString(string(s1) + string(s2) + string(s3))
				break
			}

			parts := make([]runtime.String, count)
			totalLen := 0

			for i := 0; i < count; i++ {
				s := runtime.ToString(reg[start+i])
				parts[i] = s
				totalLen += len(s)
			}

			if totalLen == 0 {
				reg[dst] = runtime.EmptyString
				break
			}

			var b strings.Builder
			b.Grow(totalLen)

			for i := 0; i < count; i++ {
				if parts[i] == runtime.EmptyString {
					continue
				}

				b.WriteString(string(parts[i]))
			}

			reg[dst] = runtime.NewString(b.String())
		case bytecode.OpSub:
			reg[dst] = runtime.Subtract(ctx, reg[src1], reg[src2])
		case bytecode.OpMulti:
			reg[dst] = runtime.Multiply(ctx, reg[src1], reg[src2])
		case bytecode.OpDiv:
			if err := vm.checkDivisionByZero(ctx, reg[src1], reg[src2]); err != nil {
				return nil, err
			}
			reg[dst] = runtime.Divide(ctx, reg[src1], reg[src2])
		case bytecode.OpMod:
			if err := vm.checkModuloByZero(ctx, reg[src2]); err != nil {
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
			} else {
				return nil, err
			}
		case bytecode.OpRegexp:
			r, err := vm.regexpCached(vm.pc-1, reg[src2])

			if err == nil {
				reg[dst] = r.Match(reg[src1])
			} else if _, catch := vm.tryCatch(vm.pc); catch {
				reg[dst] = runtime.False
			} else {
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
		case bytecode.OpLoadArray:
			reg[dst] = runtime.NewArray(int(src1))
		case bytecode.OpLoadObject:
			reg[dst] = data.NewFastObjectOf(shapeCache, vm.fastObjectDictThreshold, int(src1))
		case bytecode.OpLoadIndex, bytecode.OpLoadIndexOptional:
			src := reg[src1]
			optional := op == bytecode.OpLoadIndexOptional
			arg := reg[src2]

			if err := vm.loadIndexAndSet(ctx, dst, src, arg, optional); err != nil {
				return nil, err
			}
		case bytecode.OpLoadIndexConst, bytecode.OpLoadIndexOptionalConst:
			src := reg[src1]
			optional := op == bytecode.OpLoadIndexOptionalConst
			arg := constants[src2.Constant()]

			if err := vm.loadIndexAndSet(ctx, dst, src, arg, optional); err != nil {
				return nil, err
			}
		case bytecode.OpLoadKey, bytecode.OpLoadKeyOptional:
			src := reg[src1]
			optional := op == bytecode.OpLoadKeyOptional
			arg := reg[src2]

			if err := vm.loadKeyAndSet(ctx, dst, vm.pc-1, src, arg, optional); err != nil {
				return nil, err
			}
		case bytecode.OpLoadKeyConst, bytecode.OpLoadKeyOptionalConst:
			src := reg[src1]
			optional := op == bytecode.OpLoadKeyOptionalConst
			arg := constants[src2.Constant()]

			if err := vm.loadKeyConstAndSet(ctx, dst, vm.pc-1, inst, src, arg, optional); err != nil {
				return nil, err
			}
		case bytecode.OpLoadPropertyConst, bytecode.OpLoadPropertyOptionalConst:
			src := reg[src1]
			optional := op == bytecode.OpLoadPropertyOptionalConst
			prop := constants[src2.Constant()]

			if err := vm.loadPropertyConstAndSet(ctx, dst, vm.pc-1, inst, src, prop, optional); err != nil {
				return nil, err
			}
		case bytecode.OpLoadProperty, bytecode.OpLoadPropertyOptional:
			src := reg[src1]
			optional := op == bytecode.OpLoadPropertyOptional
			prop := reg[src2]

			if err := vm.loadPropertyAndSet(ctx, dst, vm.pc-1, src, prop, optional); err != nil {
				return nil, err
			}
		case bytecode.OpApplyQuery:
			if err := vm.applyQuery(ctx, reg, src1, constants, src2, dst); err != nil {
				return nil, err
			}
		case bytecode.OpCall, bytecode.OpProtectedCall:
			out, err := vm.callv(ctx, vm.pc-1, src1, src2)

			if err := vm.setCallResult(op, dst, out, err); err != nil {
				return nil, err
			}
		case bytecode.OpCall0, bytecode.OpProtectedCall0:
			out, err := vm.call0(ctx, vm.pc-1)

			if err := vm.setCallResult(op, dst, out, err); err != nil {
				return nil, err
			}
		case bytecode.OpCall1, bytecode.OpProtectedCall1:
			out, err := vm.call1(ctx, vm.pc-1, src1)

			if err := vm.setCallResult(op, dst, out, err); err != nil {
				return nil, err
			}
		case bytecode.OpCall2, bytecode.OpProtectedCall2:
			out, err := vm.call2(ctx, vm.pc-1, src1, src2)

			if err := vm.setCallResult(op, dst, out, err); err != nil {
				return nil, err
			}
		case bytecode.OpCall3, bytecode.OpProtectedCall3:
			out, err := vm.call3(ctx, vm.pc-1, src1)

			if err := vm.setCallResult(op, dst, out, err); err != nil {
				return nil, err
			}
		case bytecode.OpCall4, bytecode.OpProtectedCall4:
			out, err := vm.call4(ctx, vm.pc-1, src1)

			if err := vm.setCallResult(op, dst, out, err); err != nil {
				return nil, err
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
						return nil, err
					}
				}

				reg[dst] = length
			} else if _, catch := vm.tryCatch(vm.pc); catch {
				reg[dst] = runtime.ZeroInt
			} else {
				return runtime.None, runtime.TypeErrorOf(reg[src1],
					runtime.TypeString,
					runtime.TypeList,
					runtime.TypeMap,
					runtime.TypeBinary,
					runtime.TypeMeasurable,
				)
			}
		case bytecode.OpType:
			reg[dst] = runtime.String(runtime.Reflect(reg[src1]))
		case bytecode.OpFlatten:
			depth := src2.Register()

			if depth < 1 {
				depth = 1
			}

			res, err := operators.Flatten(ctx, reg[src1], depth)

			if err := vm.setOrTryCatch(dst, res, err); err != nil {
				return nil, err
			}
		case bytecode.OpClose:
			val, ok := reg[dst].(io.Closer)
			reg[dst] = runtime.None

			if ok {
				closeErr := val.Close()

				if closeErr != nil {
					if _, catch := vm.tryCatch(vm.pc); !catch {
						return nil, closeErr
					}
				}
			}
		case bytecode.OpLoadRange:
			res, err := operators.ToRange(ctx, reg[src1], reg[src2])

			if err == nil {
				reg[dst] = res
			} else {
				return nil, err
			}
		case bytecode.OpDataSet:
			reg[dst] = data.NewDataSet(src1 == 1)
		case bytecode.OpDataSetSorter:
			reg[dst] = data.NewSorter(runtime.SortDirection(src1))
		case bytecode.OpDataSetMultiSorter:
			encoded := src1.Register()
			count := src2.Register()

			reg[dst] = data.NewMultiSorter(runtime.DecodeSortDirections(encoded, count))
		case bytecode.OpDataSetCollector:
			collectorType := bytecode.CollectorType(src1)

			if collectorType == bytecode.CollectorTypeAggregate || collectorType == bytecode.CollectorTypeAggregateGroup {
				planIdx := int(src2)

				if planIdx < 0 || planIdx >= len(aggregatePlans) {
					return nil, runtime.Errorf(runtime.ErrUnexpected, "invalid aggregate plan")
				}

				plan := aggregatePlans[planIdx]

				if collectorType == bytecode.CollectorTypeAggregate {
					reg[dst] = data.NewAggregateCollector(plan)
				} else {
					reg[dst] = data.NewGroupedAggregateCollector(plan)
				}

				break
			}

			reg[dst] = data.NewCollector(collectorType)
		case bytecode.OpPush:
			ds := reg[dst].(runtime.Appendable)

			if err := ds.Append(ctx, reg[src1]); err != nil {
				if _, catch := vm.tryCatch(vm.pc); catch {
					continue
				} else {
					return nil, err
				}
			}
		case bytecode.OpArrayPush:
			ds := reg[dst].(*runtime.Array)

			_ = ds.Append(ctx, reg[src1])
		case bytecode.OpPushKV:
			tr := reg[dst].(runtime.KeyWritable)

			if err := tr.Set(ctx, reg[src1], reg[src2]); err != nil {
				if _, catch := vm.tryCatch(vm.pc); catch {
					continue
				}

				return nil, err
			}
		case bytecode.OpObjectSet:
			obj, ok := reg[dst].(*data.FastObject)
			key := runtime.ToString(reg[src1])
			value := reg[src2]

			if ok {
				_ = obj.Set(ctx, key, value)
				break
			}

			_ = reg[dst].(runtime.KeyWritable).Set(ctx, key, value)
		case bytecode.OpObjectSetConst:
			objVal := reg[dst]
			key := runtime.ToString(constants[src1.Constant()])
			value := reg[src2]

			if obj, ok := objVal.(*data.FastObject); ok {
				vm.objectSetConstCached(inst, obj, key, value)
				break
			}

			_ = objVal.(runtime.KeyWritable).Set(ctx, key, value)
		case bytecode.OpIter:
			input := reg[src1]

			switch src := input.(type) {
			case runtime.Iterable:
				iterator, err := src.Iterate(ctx)

				if err != nil {
					return nil, err
				}

				reg[dst] = data.NewIterator(iterator)
			default:
				if _, catch := vm.tryCatch(vm.pc); catch {
					// Fall back to an empty iterator
					reg[dst] = data.NoopIter
				} else {
					return nil, runtime.TypeErrorOf(src, runtime.TypeIterable)
				}
			}
		case bytecode.OpIterNext:
			iterator := reg[src1].(*data.Iterator)
			hasNext, err := iterator.HasNext(ctx)

			if err != nil {
				return nil, err
			}

			if hasNext {
				if err := iterator.Next(ctx); err != nil {
					return nil, err
				}
			} else {
				vm.pc = int(dst)
			}
		case bytecode.OpIterValue:
			iterator := reg[src1].(*data.Iterator)
			reg[dst] = iterator.Value()
		case bytecode.OpIterKey:
			iterator := reg[src1].(*data.Iterator)
			reg[dst] = iterator.Key()
		case bytecode.OpIterSkip:
			state := runtime.ToIntSafe(ctx, reg[src1])
			threshold := runtime.ToIntSafe(ctx, reg[src2])

			if state < threshold {
				state++
				reg[src1] = state
				vm.pc = int(dst)
			}
		case bytecode.OpIterLimit:
			state := runtime.ToIntSafe(ctx, reg[src1])
			threshold := runtime.ToIntSafe(ctx, reg[src2])

			if state < threshold {
				state++
				reg[src1] = state
			} else {
				vm.pc = int(dst)
			}
		case bytecode.OpStream:
			observable, eventName, options, err := vm.castSubscribeArgs(reg[dst], reg[src1], reg[src2])

			if err != nil {
				if _, catch := vm.tryCatch(vm.pc); catch {
					continue
				} else {
					return nil, err
				}
			}

			stream, err := observable.Subscribe(ctx, runtime.Subscription{
				EventName: eventName,
				Options:   options,
			})

			if err != nil {
				if _, catch := vm.tryCatch(vm.pc); catch {
					continue
				} else {
					return nil, err
				}
			}

			reg[dst] = data.NewStreamValue(stream)
		case bytecode.OpStreamIter:
			stream := reg[src1].(*data.StreamValue)

			var timeout runtime.Int

			if reg[src2] != nil && reg[src2] != runtime.None {
				t, err := runtime.CastInt(reg[src2])

				if err != nil {
					if _, catch := vm.tryCatch(vm.pc); catch {
						continue
					} else {
						return nil, err
					}
				}

				timeout = t
			}

			reg[dst] = stream.Iterate(timeout)
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
				} else {
					return nil, err
				}
			}

			if err := data.Sleep(ctx, dur); err != nil {
				return nil, err
			}
		case bytecode.OpRand:
			reg[dst] = runtime.NewFloat(runtime.RandomDefault())
		case bytecode.OpReturn:
			reg[bytecode.NoopOperand] = reg[dst]

			break loop
		default:
			// TODO: Return an error or ignore unknown opcodes?
			continue
		}
	}

	return reg[bytecode.NoopOperand], nil
}
