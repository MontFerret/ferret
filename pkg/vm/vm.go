package vm

import (
	"context"
	"io"

	"github.com/MontFerret/ferret/pkg/runtime"
	"github.com/MontFerret/ferret/pkg/vm/internal/data"
	"github.com/MontFerret/ferret/pkg/vm/internal/mem"
	"github.com/MontFerret/ferret/pkg/vm/internal/operators"
)

type VM struct {
	registers *mem.RegisterFile
	cache     *mem.Cache
	env       *Environment
	program   *Program
	pc        int
}

func New(program *Program) *VM {
	vm := new(VM)
	vm.program = program
	vm.registers = mem.NewRegisterFile(program.Registers)
	vm.cache = mem.NewCache()

	return vm
}

func (vm *VM) Run(ctx context.Context, env *Environment) (runtime.Value, error) {
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

	bytecode := vm.program.Bytecode
	constants := vm.program.Constants
	reg := vm.registers.Values
loop:
	for vm.pc < len(bytecode) {
		inst := bytecode[vm.pc]
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
			reg[dst] = operators.Add(ctx, reg[src1], reg[src2])
		case OpSub:
			reg[dst] = operators.Subtract(ctx, reg[src1], reg[src2])
		case OpMulti:
			reg[dst] = operators.Multiply(ctx, reg[src1], reg[src2])
		case OpDiv:
			reg[dst] = operators.Divide(ctx, reg[src1], reg[src2])
		case OpMod:
			reg[dst] = operators.Modulus(ctx, reg[src1], reg[src2])
		case OpIncr:
			reg[dst] = operators.Increment(ctx, reg[dst])
		case OpDecr:
			reg[dst] = operators.Decrement(ctx, reg[dst])
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
			// TODO: Add caching to avoid recompilation
			r, err := data.ToRegexp(reg[src2])

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
			var size int

			if src1 > 0 {
				size = src2.Register() - src1.Register() + 1
			}

			arr := runtime.NewArray(size)
			start := int(src1)
			end := int(src1) + size

			for i := start; i < end; i++ {
				_ = arr.Add(ctx, reg[i])
			}

			reg[dst] = arr
		case OpLoadObject:
			obj := runtime.NewObject()
			var args int

			if src1 > 0 {
				args = src2.Register() - src1.Register() + 1
			}

			start := int(src1)
			end := int(src1) + args

			for i := start; i < end; i += 2 {
				key := reg[i]
				value := reg[i+1]

				_ = obj.Set(ctx, runtime.ToString(key), value)
			}

			reg[dst] = obj
		case OpLoadIndex, OpLoadIndexOptional:
			src := reg[src1]
			arg := reg[src2]
			out, err := vm.loadIndex(ctx, src, arg)

			if err := vm.setOrOptional(dst, out, err, op == OpLoadIndexOptional); err != nil {
				return nil, err
			}

		case OpLoadKey, OpLoadKeyOptional:
			src := reg[src1]
			arg := reg[src2]
			out, err := vm.loadKey(ctx, src, arg)

			if err := vm.setOrOptional(dst, out, err, op == OpLoadKeyOptional); err != nil {
				return nil, err
			}

		case OpLoadProperty, OpLoadPropertyOptional:
			src := reg[src1]
			prop := reg[src2]

			var out runtime.Value
			var err error

			switch getter := prop.(type) {
			case runtime.String:
				out, err = vm.loadKey(ctx, src, getter)
			case runtime.Float, runtime.Int:
				out, err = vm.loadIndex(ctx, src, getter)
			default:
				out, err = vm.loadKey(ctx, src, runtime.ToString(prop))
			}

			if err := vm.setOrOptional(dst, out, err, op == OpLoadPropertyOptional); err != nil {
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
			ds := reg[dst].(runtime.List)

			if err := ds.Add(ctx, reg[src1]); err != nil {
				if _, catch := vm.tryCatch(vm.pc); catch {
					continue
				} else {
					return nil, err
				}
			}
		case OpPushKV:
			tr := reg[dst].(data.Transformer)

			if err := tr.Add(ctx, reg[src1], reg[src2]); err != nil {
				if _, catch := vm.tryCatch(vm.pc); catch {
					continue
				}

				return nil, err
			}
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
