package vm

import (
	"context"
	"errors"
	"io"

	"github.com/MontFerret/ferret/v2/pkg/bytecode"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
	"github.com/MontFerret/ferret/v2/pkg/vm/internal/data"
)

func (vm *VM) execIterOps(
	ctx context.Context,
	op bytecode.Opcode,
	dst, src1, src2 bytecode.Operand,
	reg []runtime.Value,
) error {
	switch op {
	case bytecode.OpIter:
		input := reg[src1]

		switch src := input.(type) {
		case runtime.Iterable:
			iterator, err := src.Iterate(ctx)

			if err != nil {
				if vm.unwindToProtected() {
					return nil
				}

				return err
			}

			reg[dst] = data.NewIterator(iterator)
		default:
			if _, catch := vm.tryCatch(vm.pc); catch {
				// Fall back to an empty iterator
				reg[dst] = data.NoopIter
			} else {
				if vm.unwindToProtected() {
					return nil
				}

				return runtime.TypeErrorOf(src, runtime.TypeIterable)
			}
		}
	case bytecode.OpIterNext:
		iterator := reg[src1].(*data.Iterator)
		if err := iterator.Next(ctx); err != nil {
			if errors.Is(err, io.EOF) {
				vm.pc = int(dst)
			} else {
				if vm.unwindToProtected() {
					return nil
				}

				return err
			}
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
	}

	return nil
}
