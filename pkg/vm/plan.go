package vm

import (
	"fmt"

	"github.com/MontFerret/ferret/v2/pkg/bytecode"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
	"github.com/MontFerret/ferret/v2/pkg/vm/internal/diagnostics"
)

type execPlan struct {
	instructions           []execInstruction
	hostCallDescriptors    []callDescriptor
	udfCallDescriptors     []callDescriptor
	udfTailCallDescriptors []callDescriptor
}

func buildExecPlan(program *bytecode.Program) (execPlan, error) {
	if program == nil || len(program.Bytecode) == 0 {
		return execPlan{}, nil
	}

	instructions := make([]execInstruction, len(program.Bytecode))
	aggregateSelectorSlots := program.Metadata.AggregateSelectorSlots
	matchFailTargets := program.Metadata.MatchFailTargets
	constants := program.Constants
	udfs := program.Functions.UserDefined
	reg := map[bytecode.Operand]runtime.Value{}
	hostCallDesc := make([]callDescriptor, 0, 4)
	udfCallDesc := make([]callDescriptor, 0, 4)
	udfTailCallDesc := make([]callDescriptor, 0, 4)
	errs := diagnostics.NewInitializationErrorSet(4)

	if len(aggregateSelectorSlots) > 0 && len(aggregateSelectorSlots) != len(program.Bytecode) {
		errs.Add(
			fmt.Errorf("aggregate selector slot metadata length %d does not match bytecode length %d", len(aggregateSelectorSlots), len(program.Bytecode)),
			0,
			bytecode.NoopOperand,
		)
	}

	if len(matchFailTargets) > 0 && len(matchFailTargets) != len(program.Bytecode) {
		errs.Add(
			fmt.Errorf("match fail target metadata length %d does not match bytecode length %d", len(matchFailTargets), len(program.Bytecode)),
			0,
			bytecode.NoopOperand,
		)
	}

	for pc, inst := range program.Bytecode {
		instructions[pc] = execInstruction{
			Instruction: inst,
			InlineSlot:  -1,
		}

		op := inst.Opcode
		dst, src1, src2 := inst.Operands[0], inst.Operands[1], inst.Operands[2]

		switch op {
		case bytecode.OpLoadConst:
			reg[dst] = constants[src1.Constant()]
		case bytecode.OpMove, bytecode.OpMoveTracked:
			if val, ok := reg[src1]; ok {
				reg[dst] = val
			} else {
				delete(reg, dst)
			}
		case bytecode.OpCall, bytecode.OpProtectedCall:
			fnID, err := getUDFID(reg[dst])

			if err != nil {
				errs.Add(err, pc, dst)
				continue
			}

			udf := udfs[fnID]
			descriptor := callDescriptor{
				PC:               pc,
				CallSitePC:       pc - 1,
				DisplayName:      udf.DisplayName,
				Dst:              dst,
				ID:               fnID,
				ArgCount:         callArgCount(src1, src2),
				ArgStart:         int(src1),
				RecoveryBoundary: bytecode.IsProtectedCall(op),
			}

			udfCallDesc = append(udfCallDesc, descriptor)
			instructions[pc].InlineSlot = len(udfCallDesc) - 1
		case bytecode.OpTailCall:
			fnID, err := getUDFID(reg[dst])

			if err != nil {
				errs.Add(err, pc, dst)
				continue
			}

			udf := udfs[fnID]
			descriptor := callDescriptor{
				PC:          pc,
				CallSitePC:  pc - 1,
				DisplayName: udf.DisplayName,
				Dst:         dst,
				ID:          fnID,
				ArgCount:    callArgCount(src1, src2),
				ArgStart:    int(src1),
			}

			udfTailCallDesc = append(udfTailCallDesc, descriptor)
			instructions[pc].InlineSlot = len(udfTailCallDesc) - 1
		case bytecode.OpHCall, bytecode.OpProtectedHCall:
			descriptor := callDescriptor{
				PC:               pc,
				CallSitePC:       pc - 1,
				Dst:              dst,
				ID:               len(hostCallDesc),
				ArgCount:         callArgCount(src1, src2),
				ArgStart:         int(src1),
				RecoveryBoundary: bytecode.IsProtectedCall(op),
			}

			fnName, err := resolveHostFnName(reg, dst)

			if err != nil {
				errs.Add(err, pc, dst)
				continue
			}

			descriptor.DisplayName = fnName
			hostCallDesc = append(hostCallDesc, descriptor)
			instructions[pc].InlineSlot = descriptor.ID
		case bytecode.OpAggregateUpdate, bytecode.OpAggregateGroupUpdate:
			slot, err := aggregateSelectorSlotAt(aggregateSelectorSlots, pc)
			if err != nil {
				errs.Add(err, pc, dst)
				continue
			}

			instructions[pc].InlineSlot = slot
		case bytecode.OpMatchLoadPropertyConst:
			target, err := matchFailTargetAt(matchFailTargets, len(program.Bytecode), pc)
			if err != nil {
				errs.Add(err, pc, dst)
				continue
			}

			instructions[pc].InlineSlot = target
		}

		if op != bytecode.OpLoadConst && op != bytecode.OpMove && op != bytecode.OpMoveTracked && dst.IsRegister() {
			delete(reg, dst)
		}
	}

	if errs.Size() > 0 {
		return execPlan{}, errs
	}

	if len(hostCallDesc) == 0 {
		hostCallDesc = nil
	}

	if len(udfCallDesc) == 0 {
		udfCallDesc = nil
	}

	if len(udfTailCallDesc) == 0 {
		udfTailCallDesc = nil
	}

	return execPlan{
		instructions:           instructions,
		hostCallDescriptors:    hostCallDesc,
		udfCallDescriptors:     udfCallDesc,
		udfTailCallDescriptors: udfTailCallDesc,
	}, nil
}

func aggregateSelectorSlotAt(slots []int, pc int) (int, error) {
	if pc < 0 || pc >= len(slots) {
		return -1, fmt.Errorf("invalid aggregate selector slot metadata at pc %d", pc)
	}

	slot := slots[pc]
	if slot < 0 {
		return -1, fmt.Errorf("invalid aggregate selector slot at pc %d", pc)
	}

	return slot, nil
}

func matchFailTargetAt(targets []int, bytecodeLen, pc int) (int, error) {
	if pc < 0 || pc >= len(targets) {
		return -1, fmt.Errorf("invalid match fail target metadata at pc %d", pc)
	}

	target := targets[pc]
	if target < 0 || target >= bytecodeLen {
		return -1, fmt.Errorf("invalid match fail target at pc %d", pc)
	}

	return target, nil
}

func resolveHostFnName(reg map[bytecode.Operand]runtime.Value, dst bytecode.Operand) (string, error) {
	val, ok := reg[dst]

	if ok {
		fnName, ok := val.(runtime.String)

		if ok {
			return fnName.String(), nil
		}
	}

	return "", ErrInvalidFunctionName
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
