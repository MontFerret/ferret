package vm

import (
	"context"
	"io"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
	"github.com/MontFerret/ferret/v2/pkg/vm/internal/data"
	"github.com/MontFerret/ferret/v2/pkg/vm/internal/mem"
	"github.com/MontFerret/ferret/v2/pkg/vm/internal/operators"
)

type VM struct {
	registers               *mem.RegisterFile
	cache                   *mem.Cache
	env                     *Environment
	program                 *Program
	fastObjectDictThreshold int
	bytecode                []Instruction
	pc                      int
}

func New(program *Program) *VM {
	return NewWithOptions(program)
}

func NewWithOptions(program *Program, opts ...VMOption) *VM {
	cfg := defaultVMConfig()

	for _, opt := range opts {
		opt(&cfg)
	}

	vm := new(VM)
	vm.program = program
	vm.registers = mem.NewRegisterFile(program.Registers)
	vm.cache = mem.NewCache(len(program.Bytecode), cfg.shapeCacheLimit)
	vm.fastObjectDictThreshold = cfg.fastObjectDictThreshold
	vm.bytecode = make([]Instruction, len(program.Bytecode))
	copy(vm.bytecode, program.Bytecode)

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

	bytecode := vm.bytecode
	constants := vm.program.Constants
	reg := vm.registers.Values
	shapeCache := vm.cache.ShapeCache
loop:
	for vm.pc < len(bytecode) {
		inst := &bytecode[vm.pc]
		op := inst.Opcode
		dst, src1, src2 := inst.Operands[0], inst.Operands[1], inst.Operands[2]
		vm.pc++

		switch op {
		case OpLoadNone:
			reg[dst] = runtime.None
		case OpLoadZero:
			reg[dst] = runtime.ZeroInt
		case OpLoadBool:
			reg[dst] = runtime.Boolean(src1 == 1)
		case OpMove:
			reg[dst] = reg[src1]
		case OpLoadConst:
			reg[dst] = constants[src1.Constant()]
		case OpLoadParam:
			name := constants[src1.Constant()]
			reg[dst] = vm.env.Params[name.String()]
		case OpJump:
			vm.pc = int(dst)
		case OpJumpIfFalse:
			if !runtime.ToBoolean(reg[src1]) {
				vm.pc = int(dst)
			}
		case OpJumpIfTrue:
			if runtime.ToBoolean(reg[src1]) {
				vm.pc = int(dst)
			}
		case OpAdd:
			reg[dst] = runtime.Add(ctx, reg[src1], reg[src2])
		case OpSub:
			reg[dst] = runtime.Subtract(ctx, reg[src1], reg[src2])
		case OpMulti:
			reg[dst] = runtime.Multiply(ctx, reg[src1], reg[src2])
		case OpDiv:
			if err := vm.checkDivisionByZero(ctx, reg[src1], reg[src2]); err != nil {
				return nil, err
			}
			reg[dst] = runtime.Divide(ctx, reg[src1], reg[src2])
		case OpMod:
			if err := vm.checkModuloByZero(ctx, reg[src2]); err != nil {
				return nil, err
			}
			reg[dst] = runtime.Modulus(ctx, reg[src1], reg[src2])
		case OpIncr:
			reg[dst] = runtime.Increment(ctx, reg[dst])
		case OpDecr:
			reg[dst] = runtime.Decrement(ctx, reg[dst])
		case OpCastBool:
			reg[dst] = runtime.ToBoolean(reg[src1])
		case OpNegate:
			reg[dst] = operators.Negate(reg[src1])
		case OpFlipPositive:
			reg[dst] = operators.Positive(reg[src1])
		case OpFlipNegative:
			reg[dst] = operators.Negative(reg[src1])
		case OpCmp:
			reg[dst] = operators.Compare(ctx, reg[src1], reg[src2])
		case OpNot:
			reg[dst] = !runtime.ToBoolean(reg[src1])
		case OpEq:
			reg[dst] = operators.Equals(ctx, reg[src1], reg[src2])
		case OpNe:
			reg[dst] = operators.NotEquals(ctx, reg[src1], reg[src2])
		case OpGt:
			reg[dst] = operators.GreaterThan(ctx, reg[src1], reg[src2])
		case OpLt:
			reg[dst] = operators.LessThan(ctx, reg[src1], reg[src2])
		case OpGte:
			reg[dst] = operators.GreaterThanOrEqual(ctx, reg[src1], reg[src2])
		case OpLte:
			reg[dst] = operators.LessThanOrEqual(ctx, reg[src1], reg[src2])
		case OpIn:
			reg[dst] = operators.Contains(ctx, reg[src2], reg[src1])
		case OpLike:
			res, err := operators.Like(reg[src1], reg[src2])

			if err == nil {
				reg[dst] = res
			} else {
				return nil, err
			}
		case OpRegexp:
			r, err := vm.regexpCached(vm.pc-1, reg[src2])

			if err == nil {
				reg[dst] = r.Match(reg[src1])
			} else if _, catch := vm.tryCatch(vm.pc); catch {
				reg[dst] = runtime.False
			} else {
				return nil, err
			}
		case OpAllEq, OpAllNe, OpAllGt, OpAllGte, OpAllLt, OpAllLte, OpAllIn:
			cmp := operators.ComparatorFromByte(int(op) - int(OpAllEq))
			res, err := operators.ArrayAll(ctx, cmp, reg[src1], reg[src2])

			if err := vm.setOrTryCatch(dst, res, err); err != nil {
				return nil, err
			}
		case OpAnyEq, OpAnyNe, OpAnyGt, OpAnyGte, OpAnyLt, OpAnyLte, OpAnyIn:
			cmp := operators.ComparatorFromByte(int(op) - int(OpAnyEq))
			res, err := operators.ArrayAny(ctx, cmp, reg[src1], reg[src2])

			if err := vm.setOrTryCatch(dst, res, err); err != nil {
				return nil, err
			}
		case OpNoneEq, OpNoneNe, OpNoneGt, OpNoneGte, OpNoneLt, OpNoneLte, OpNoneIn:
			cmp := operators.ComparatorFromByte(int(op) - int(OpNoneEq))
			res, err := operators.ArrayNone(ctx, cmp, reg[src1], reg[src2])

			if err := vm.setOrTryCatch(dst, res, err); err != nil {
				return nil, err
			}
		case OpLoadArray:
			reg[dst] = runtime.NewArray(int(src1))
		case OpLoadObject:
			reg[dst] = data.NewFastObjectOf(shapeCache, vm.fastObjectDictThreshold, int(src1))
		case OpLoadIndex, OpLoadIndexOptional:
			src := reg[src1]
			arg := reg[src2]
			out, err := vm.loadIndex(ctx, src, arg)

			if err := vm.setOrOptional(dst, out, err, op == OpLoadIndexOptional); err != nil {
				return nil, err
			}

		case OpLoadIndexConst, OpLoadIndexOptionalConst:
			src := reg[src1]
			arg := constants[src2.Constant()]
			out, err := vm.loadIndex(ctx, src, arg)

			if err := vm.setOrOptional(dst, out, err, op == OpLoadIndexOptionalConst); err != nil {
				return nil, err
			}

		case OpLoadKey, OpLoadKeyOptional:
			src := reg[src1]
			arg := reg[src2]
			out, err := vm.loadKeyCached(ctx, vm.pc-1, src, arg)

			if err := vm.setOrOptional(dst, out, err, op == OpLoadKeyOptional); err != nil {
				return nil, err
			}

		case OpLoadKeyConst, OpLoadKeyOptionalConst:
			src := reg[src1]
			arg := constants[src2.Constant()]
			out, err := vm.loadKeyConstCached(ctx, vm.pc-1, inst, src, arg)

			if err := vm.setOrOptional(dst, out, err, op == OpLoadKeyOptionalConst); err != nil {
				return nil, err
			}

		case OpLoadPropertyConst, OpLoadPropertyOptionalConst:
			src := reg[src1]
			prop := constants[src2.Constant()]

			var out runtime.Value
			var err error

			switch getter := prop.(type) {
			case runtime.String:
				out, err = vm.loadKeyConstCached(ctx, vm.pc-1, inst, src, getter)
			case runtime.Float, runtime.Int:
				out, err = vm.loadIndex(ctx, src, getter)
			default:
				out, err = vm.loadKeyConstCached(ctx, vm.pc-1, inst, src, runtime.ToString(prop))
			}

			if err := vm.setOrOptional(dst, out, err, op == OpLoadPropertyOptionalConst); err != nil {
				return nil, err
			}

		case OpLoadProperty, OpLoadPropertyOptional:
			src := reg[src1]
			prop := reg[src2]

			var out runtime.Value
			var err error

			switch getter := prop.(type) {
			case runtime.String:
				out, err = vm.loadKeyCached(ctx, vm.pc-1, src, getter)
			case runtime.Float, runtime.Int:
				out, err = vm.loadIndex(ctx, src, getter)
			default:
				out, err = vm.loadKeyCached(ctx, vm.pc-1, src, runtime.ToString(prop))
			}

			if err := vm.setOrOptional(dst, out, err, op == OpLoadPropertyOptional); err != nil {
				return nil, err
			}
		case OpApplyQuery:
			src := reg[src1]

			if src1.IsConstant() {
				src = constants[src1.Constant()]
			}

			var arg runtime.Value

			if src2.IsConstant() {
				arg = constants[src2.Constant()]
			} else {
				arg = reg[src2]
			}

			var query runtime.Query

			switch value := arg.(type) {
			case runtime.Query:
				query = value
			case *runtime.Array:
				length, err := value.Length(ctx)
				if err != nil {
					if err := vm.setOrTryCatch(dst, runtime.None, err); err != nil {
						return nil, err
					}

					break
				}

				if length < 2 {
					if err := vm.setOrTryCatch(dst, runtime.None, runtime.TypeErrorOf(arg, runtime.TypeQuery)); err != nil {
						return nil, err
					}

					break
				}

				kindVal, err := value.Get(ctx, runtime.NewInt(0))
				if err != nil {
					if err := vm.setOrTryCatch(dst, runtime.None, err); err != nil {
						return nil, err
					}

					break
				}

				payloadVal, err := value.Get(ctx, runtime.NewInt(1))
				if err != nil {
					if err := vm.setOrTryCatch(dst, runtime.None, err); err != nil {
						return nil, err
					}

					break
				}

				var paramsVal runtime.Value = runtime.None
				if length > 2 {
					paramsVal, err = value.Get(ctx, runtime.NewInt(2))
					if err != nil {
						if err := vm.setOrTryCatch(dst, runtime.None, err); err != nil {
							return nil, err
						}

						break
					}
				}

				kind, err := runtime.CastString(kindVal)
				if err != nil {
					if err := vm.setOrTryCatch(dst, runtime.None, runtime.TypeErrorOf(kindVal, runtime.TypeString)); err != nil {
						return nil, err
					}

					break
				}

				payload := runtime.EmptyString
				if payloadVal != runtime.None {
					payload, err = runtime.CastString(payloadVal)
					if err != nil {
						if err := vm.setOrTryCatch(dst, runtime.None, runtime.TypeErrorOf(payloadVal, runtime.TypeString, runtime.TypeNone)); err != nil {
							return nil, err
						}

						break
					}
				}

				query = runtime.NewQuery(kind, payload)
				query.Params = paramsVal
			default:
				if err := vm.setOrTryCatch(dst, runtime.None, runtime.TypeErrorOf(arg, runtime.TypeQuery, runtime.TypeArray)); err != nil {
					return nil, err
				}

				break
			}

			queryable, ok := src.(runtime.Queryable)

			if !ok {
				if err := vm.setOrTryCatch(dst, runtime.None, runtime.TypeErrorOf(src, runtime.TypeQueryable)); err != nil {
					return nil, err
				}

				break
			}

			res, err := queryable.Query(ctx, query)

			if err := vm.setOrTryCatch(dst, res, err); err != nil {
				return nil, err
			}

		case OpCall, OpProtectedCall:
			out, err := vm.callv(ctx, vm.pc-1, src1, src2)

			if err := vm.setCallResult(op, dst, out, err); err != nil {
				return nil, err
			}

		case OpCall0, OpProtectedCall0:
			out, err := vm.call0(ctx, vm.pc-1)

			if err := vm.setCallResult(op, dst, out, err); err != nil {
				return nil, err
			}

		case OpCall1, OpProtectedCall1:
			out, err := vm.call1(ctx, vm.pc-1, src1)

			if err := vm.setCallResult(op, dst, out, err); err != nil {
				return nil, err
			}

		case OpCall2, OpProtectedCall2:
			out, err := vm.call2(ctx, vm.pc-1, src1, src2)

			if err := vm.setCallResult(op, dst, out, err); err != nil {
				return nil, err
			}

		case OpCall3, OpProtectedCall3:
			out, err := vm.call3(ctx, vm.pc-1, src1)

			if err := vm.setCallResult(op, dst, out, err); err != nil {
				return nil, err
			}

		case OpCall4, OpProtectedCall4:
			out, err := vm.call4(ctx, vm.pc-1, src1)

			if err := vm.setCallResult(op, dst, out, err); err != nil {
				return nil, err
			}

		case OpExists:
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
		case OpLength:
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
		case OpType:
			reg[dst] = runtime.String(runtime.Reflect(reg[src1]))
		case OpFlatten:
			depth := src2.Register()

			if depth < 1 {
				depth = 1
			}

			res, err := operators.Flatten(ctx, reg[src1], depth)

			if err := vm.setOrTryCatch(dst, res, err); err != nil {
				return nil, err
			}
		case OpClose:
			val, ok := reg[dst].(io.Closer)
			reg[dst] = runtime.None

			if ok {
				err := val.Close()

				if err != nil {
					if _, catch := vm.tryCatch(vm.pc); !catch {
						return nil, err
					}
				}
			}
		case OpLoadRange:
			res, err := operators.ToRange(ctx, reg[src1], reg[src2])

			if err == nil {
				reg[dst] = res
			} else {
				return nil, err
			}
		case OpDataSet:
			reg[dst] = data.NewDataSet(src1 == 1)
		case OpDataSetSorter:
			reg[dst] = data.NewSorter(runtime.SortDirection(src1))
		case OpDataSetMultiSorter:
			encoded := src1.Register()
			count := src2.Register()

			reg[dst] = data.NewMultiSorter(runtime.DecodeSortDirections(encoded, count))
		case OpDataSetCollector:
			reg[dst] = data.NewCollector(data.CollectorType(src1))
		case OpPush:
			ds := reg[dst].(runtime.Appendable)

			if err := ds.Append(ctx, reg[src1]); err != nil {
				if _, catch := vm.tryCatch(vm.pc); catch {
					continue
				} else {
					return nil, err
				}
			}
		case OpArrayPush:
			ds := reg[dst].(*runtime.Array)

			_ = ds.Append(ctx, reg[src1])
		case OpPushKV:
			tr := reg[dst].(runtime.KeyWritable)

			if err := tr.Set(ctx, reg[src1], reg[src2]); err != nil {
				if _, catch := vm.tryCatch(vm.pc); catch {
					continue
				}

				return nil, err
			}
		case OpObjectSet:
			obj, ok := reg[dst].(*data.FastObject)
			key := runtime.ToString(reg[src1])
			value := reg[src2]

			if ok {
				_ = obj.Set(ctx, key, value)
				break
			}

			_ = reg[dst].(runtime.KeyWritable).Set(ctx, key, value)
		case OpObjectSetConst:
			objVal := reg[dst]
			key := runtime.ToString(constants[src1.Constant()])
			value := reg[src2]

			if obj, ok := objVal.(*data.FastObject); ok {
				vm.objectSetConstCached(inst, obj, key, value)
				break
			}

			_ = objVal.(runtime.KeyWritable).Set(ctx, key, value)
		case OpIter:
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
		case OpIterNext:
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
		case OpIterValue:
			iterator := reg[src1].(*data.Iterator)
			reg[dst] = iterator.Value()
		case OpIterKey:
			iterator := reg[src1].(*data.Iterator)
			reg[dst] = iterator.Key()
		case OpIterSkip:
			state := runtime.ToIntSafe(ctx, reg[src1])
			threshold := runtime.ToIntSafe(ctx, reg[src2])

			if state < threshold {
				state++
				reg[src1] = state
				vm.pc = int(dst)
			}
		case OpIterLimit:
			state := runtime.ToIntSafe(ctx, reg[src1])
			threshold := runtime.ToIntSafe(ctx, reg[src2])

			if state < threshold {
				state++
				reg[src1] = state
			} else {
				vm.pc = int(dst)
			}
		case OpStream:
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
		case OpStreamIter:
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
		case OpSleep:
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
		case OpRand:
			reg[dst] = runtime.NewFloat(runtime.RandomDefault())
		case OpReturn:
			reg[NoopOperand] = reg[dst]

			break loop
		default:
			// TODO: Return an error or ignore unknown opcodes?
			continue
		}
	}

	return reg[NoopOperand], nil
}
