package vm

import (
	"fmt"

	"github.com/MontFerret/ferret/v2/pkg/bytecode"
	"github.com/MontFerret/ferret/v2/pkg/internal/hostfunction"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
	"github.com/MontFerret/ferret/v2/pkg/vm/internal/diagnostics"
)

type (
	execPlan struct {
		instructions           []execInstruction
		hostCallDescriptors    []callDescriptor
		udfCallDescriptors     []callDescriptor
		udfTailCallDescriptors []callDescriptor
		paramLoadDescriptors   []paramLoadDescriptor
	}

	paramLoadDescriptor struct {
		PC   int
		Dst  bytecode.Operand
		Slot int
	}
)

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
	paramLoadDesc := make([]paramLoadDescriptor, 0, 4)
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
			bindingID, err := getHostFunctionID(reg[dst])
			if err != nil {
				errs.Add(err, pc, dst)
				continue
			}

			if bindingID < 0 || bindingID >= int64(len(program.Functions.Host)) {
				errs.Add(
					diagnostics.NewInvariantError(
						"invalid host binding id",
						runtime.Errorf(runtime.ErrUnexpected, "host binding id %d at pc %d is outside signature table of size %d", bindingID, pc, len(program.Functions.Host)),
					),
					pc,
					dst,
				)
				continue
			}

			bindingIDIndex := int(bindingID)
			hostFn := program.Functions.Host[bindingIDIndex]
			hostName := hostfunction.CanonicalName(hostFn.Name)
			argCount := callArgCount(src1, src2)
			if argCount != hostFn.ArgCount {
				errs.Add(
					diagnostics.NewInvariantError(
						"host call signature mismatch",
						runtime.Errorf(runtime.ErrUnexpected, "host binding %d %q expects compiled argument count %d, but callsite at pc %d has %d", bindingID, hostName, hostFn.ArgCount, pc, argCount),
					),
					pc,
					dst,
				)
				continue
			}

			descriptor := callDescriptor{
				PC:               pc,
				CallSitePC:       pc - 1,
				DisplayName:      hostName,
				Dst:              dst,
				ID:               bindingIDIndex,
				ArgCount:         hostFn.ArgCount,
				ArgStart:         int(src1),
				RecoveryBoundary: bytecode.IsProtectedCall(op),
			}

			hostCallDesc = append(hostCallDesc, descriptor)
			instructions[pc].InlineSlot = len(hostCallDesc) - 1
		case bytecode.OpLoadParam:
			slot, err := paramSlotAt(len(program.Params), src1, pc)
			if err != nil {
				errs.Add(err, pc, dst)
				continue
			}

			paramLoadDesc = append(paramLoadDesc, paramLoadDescriptor{
				PC:   pc,
				Dst:  dst,
				Slot: slot,
			})
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

		if !skipStaticRegisterInvalidation(op) && dst.IsRegister() {
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

	if len(paramLoadDesc) == 0 {
		paramLoadDesc = nil
	}

	return execPlan{
		instructions:           instructions,
		hostCallDescriptors:    hostCallDesc,
		udfCallDescriptors:     udfCallDesc,
		udfTailCallDescriptors: udfTailCallDesc,
		paramLoadDescriptors:   paramLoadDesc,
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

func paramSlotAt(paramCount int, slot bytecode.Operand, pc int) (int, error) {
	idx := int(slot) - 1
	if idx < 0 || idx >= paramCount {
		return -1, fmt.Errorf("invalid parameter slot %d at pc %d", slot, pc)
	}

	return idx, nil
}

func getHostFunctionID(value runtime.Value) (int64, error) {
	id, ok := value.(runtime.Int)
	if !ok {
		return -1, ErrInvalidFunctionName
	}

	return int64(id), nil
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
