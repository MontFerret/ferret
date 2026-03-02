package vm

import (
	"context"

	"github.com/MontFerret/ferret/v2/pkg/bytecode"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
	"github.com/MontFerret/ferret/v2/pkg/vm/internal/data"
)

func (vm *VM) execDatasetOps(
	ctx context.Context,
	op bytecode.Opcode,
	inst *data.ExecInstruction,
	dst, src1, src2 bytecode.Operand,
	reg []runtime.Value,
	constants []runtime.Value,
	aggregatePlans []bytecode.AggregatePlan,
) error {
	switch op {
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
				return vm.handleProtectedError(runtime.Errorf(runtime.ErrUnexpected, "invalid aggregate plan"))
			}

			plan := aggregatePlans[planIdx]

			if collectorType == bytecode.CollectorTypeAggregate {
				reg[dst] = data.NewAggregateCollector(plan)
			} else {
				reg[dst] = data.NewGroupedAggregateCollector(plan)
			}

			return nil
		}

		reg[dst] = data.NewCollector(collectorType)
	case bytecode.OpPush:
		ds := reg[dst].(runtime.Appendable)

		if err := ds.Append(ctx, reg[src1]); err != nil {
			return vm.handleError(err)
		}
	case bytecode.OpArrayPush:
		ds := reg[dst].(*runtime.Array)

		_ = ds.Append(ctx, reg[src1])
	case bytecode.OpPushKV:
		tr := reg[dst].(runtime.KeyWritable)

		if err := tr.Set(ctx, reg[src1], reg[src2]); err != nil {
			return vm.handleError(err)
		}
	case bytecode.OpObjectSet:
		obj, ok := reg[dst].(*data.FastObject)
		key := runtime.ToString(reg[src1])
		value := reg[src2]

		if ok {
			_ = obj.Set(ctx, key, value)
			return nil
		}

		_ = reg[dst].(runtime.KeyWritable).Set(ctx, key, value)
	case bytecode.OpObjectSetConst:
		objVal := reg[dst]
		key := runtime.ToString(constants[src1.Constant()])
		value := reg[src2]

		if obj, ok := objVal.(*data.FastObject); ok {
			vm.objectSetConstCached(inst, obj, key, value)
			return nil
		}

		_ = objVal.(runtime.KeyWritable).Set(ctx, key, value)
	}

	return nil
}
