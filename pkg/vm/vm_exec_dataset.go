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

	}

	return nil
}
