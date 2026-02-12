package optimization

import (
	"context"
	"fmt"

	"github.com/MontFerret/ferret/pkg/runtime"
	"github.com/MontFerret/ferret/pkg/vm"
)

const ConstantPropagationPassName = "constant-propagation"

// ConstantPropagationPass performs a simple constant propagation and folding.
// It is conservative across control-flow merges: a register is considered constant
// only if all predecessors agree on the same constant value.
type ConstantPropagationPass struct{}

// NewConstantPropagationPass creates a new constant propagation pass.
func NewConstantPropagationPass() Pass {
	return &ConstantPropagationPass{}
}

func (p *ConstantPropagationPass) Name() string {
	return ConstantPropagationPassName
}

func (p *ConstantPropagationPass) Requires() []string {
	return []string{}
}

func (p *ConstantPropagationPass) Run(ctx *PassContext) (*PassResult, error) {
	if ctx == nil || ctx.Program == nil || ctx.CFG == nil {
		return &PassResult{Modified: false}, nil
	}
	if len(ctx.Program.Bytecode) == 0 {
		return &PassResult{Modified: false}, nil
	}

	constIndex := buildConstIndex(ctx.Program.Constants)
	in := make(map[int]constState)
	out := make(map[int]constState)

	modified := false
	changed := true
	bg := context.Background()

	for changed {
		changed = false

		for _, block := range ctx.CFG.Blocks {
			if block == ctx.CFG.Exit {
				continue
			}

			inState := meetPreds(block, out)
			prevIn := in[block.ID]
			if !statesEqual(prevIn, inState) {
				in[block.ID] = inState
				changed = true
			}

			state := copyState(inState)
			for i := 0; i < len(block.Instructions); i++ {
				inst := block.Instructions[i]
				if applyConstFolding(&inst, state, ctx.Program, constIndex, bg) {
					block.Instructions[i] = inst
					ctx.Program.Bytecode[block.Start+i] = inst
					modified = true
					changed = true
				}
			}

			if !statesEqual(out[block.ID], state) {
				out[block.ID] = state
				changed = true
			}
		}
	}

	return &PassResult{
		Modified: modified,
		Metadata: map[string]any{},
	}, nil
}

type constState map[int]runtime.Value

func copyState(in constState) constState {
	if len(in) == 0 {
		return constState{}
	}
	out := make(constState, len(in))
	for k, v := range in {
		out[k] = v
	}
	return out
}

func statesEqual(a, b constState) bool {
	if len(a) != len(b) {
		return false
	}
	for k, v := range a {
		if bv, ok := b[k]; !ok || !constEqual(v, bv) {
			return false
		}
	}
	return true
}

func meetPreds(block *BasicBlock, out map[int]constState) constState {
	if len(block.Predecessors) == 0 {
		return constState{}
	}

	var base constState
	for _, pred := range block.Predecessors {
		if pred == nil {
			continue
		}
		base = out[pred.ID]
		break
	}
	if len(base) == 0 {
		return constState{}
	}

	result := make(constState, len(base))
	for reg, val := range base {
		keep := true
		for _, pred := range block.Predecessors {
			ps := out[pred.ID]
			pv, ok := ps[reg]
			if !ok || !constEqual(val, pv) {
				keep = false
				break
			}
		}
		if keep {
			result[reg] = val
		}
	}

	return result
}

func applyConstFolding(inst *vm.Instruction, state constState, program *vm.Program, constIndex map[string]int, bg context.Context) bool {
	modified := false
	defs := defsOnly(*inst)
	newConsts := make(map[int]runtime.Value)

	switch inst.Opcode {
	case vm.OpLoadConst:
		if inst.Operands[0].IsRegister() {
			val := program.Constants[inst.Operands[1].Constant()]
			if isSimpleConst(val) {
				newConsts[inst.Operands[0].Register()] = val
			}
		}
	case vm.OpLoadNone:
		if inst.Operands[0].IsRegister() {
			newConsts[inst.Operands[0].Register()] = runtime.None
		}
	case vm.OpLoadZero:
		if inst.Operands[0].IsRegister() {
			newConsts[inst.Operands[0].Register()] = runtime.ZeroInt
		}
	case vm.OpLoadBool:
		if inst.Operands[0].IsRegister() {
			val := runtime.Boolean(inst.Operands[1] == 1)
			newConsts[inst.Operands[0].Register()] = val
		}
	case vm.OpMove:
		if inst.Operands[0].IsRegister() && inst.Operands[1].IsRegister() {
			if val, ok := state[inst.Operands[1].Register()]; ok {
				newConsts[inst.Operands[0].Register()] = val
				modified = replaceWithConstLoad(inst, inst.Operands[0].Register(), val, program, constIndex) || modified
			}
		}
	case vm.OpIncr, vm.OpDecr:
		if inst.Operands[0].IsRegister() {
			reg := inst.Operands[0].Register()
			if val, ok := state[reg]; ok {
				var out runtime.Value
				if inst.Opcode == vm.OpIncr {
					out = increment(bg, val)
				} else {
					out = decrement(bg, val)
				}
				if isSimpleConst(out) {
					newConsts[reg] = out
					modified = replaceWithConstLoad(inst, reg, out, program, constIndex) || modified
				}
			}
		}
	case vm.OpCastBool, vm.OpNegate, vm.OpFlipPositive, vm.OpFlipNegative, vm.OpNot:
		if inst.Operands[0].IsRegister() && inst.Operands[1].IsRegister() {
			src := inst.Operands[1].Register()
			if val, ok := state[src]; ok {
				out, ok := foldUnary(inst.Opcode, val, bg)
				if ok && isSimpleConst(out) {
					dst := inst.Operands[0].Register()
					newConsts[dst] = out
					modified = replaceWithConstLoad(inst, dst, out, program, constIndex) || modified
				}
			}
		}
	case vm.OpAdd, vm.OpSub, vm.OpMulti, vm.OpDiv, vm.OpMod,
		vm.OpCmp, vm.OpEq, vm.OpNe, vm.OpGt, vm.OpLt, vm.OpGte, vm.OpLte:
		if inst.Operands[0].IsRegister() && inst.Operands[1].IsRegister() && inst.Operands[2].IsRegister() {
			left, lok := state[inst.Operands[1].Register()]
			right, rok := state[inst.Operands[2].Register()]
			if lok && rok {
				out, ok := foldBinary(inst.Opcode, left, right, bg)
				if ok && isSimpleConst(out) {
					dst := inst.Operands[0].Register()
					newConsts[dst] = out
					modified = replaceWithConstLoad(inst, dst, out, program, constIndex) || modified
				}
			}
		}
	}

	for _, reg := range defs {
		if val, ok := newConsts[reg]; ok {
			state[reg] = val
		} else {
			delete(state, reg)
		}
	}

	return modified
}

func defsOnly(inst vm.Instruction) []int {
	_, defs := instructionUseDef(inst)
	return defs
}

func isSimpleConst(val runtime.Value) bool {
	if val == nil {
		return false
	}
	if val == runtime.None {
		return true
	}
	switch val.(type) {
	case runtime.Int, runtime.Float, runtime.String, runtime.Boolean:
		return true
	default:
		return false
	}
}

func constEqual(a, b runtime.Value) bool {
	if a == b {
		return true
	}
	switch av := a.(type) {
	case runtime.Int:
		if bv, ok := b.(runtime.Int); ok {
			return av == bv
		}
	case runtime.Float:
		if bv, ok := b.(runtime.Float); ok {
			return av == bv
		}
	case runtime.String:
		if bv, ok := b.(runtime.String); ok {
			return av == bv
		}
	case runtime.Boolean:
		if bv, ok := b.(runtime.Boolean); ok {
			return av == bv
		}
	}
	return false
}

func foldUnary(op vm.Opcode, val runtime.Value, bg context.Context) (runtime.Value, bool) {
	switch op {
	case vm.OpCastBool:
		return runtime.ToBoolean(val), true
	case vm.OpNot:
		return runtime.Boolean(!runtime.ToBoolean(val)), true
	case vm.OpNegate:
		return negate(val), true
	case vm.OpFlipPositive:
		return positive(val), true
	case vm.OpFlipNegative:
		return negative(val), true
	default:
		return nil, false
	}
}

func foldBinary(op vm.Opcode, left, right runtime.Value, bg context.Context) (runtime.Value, bool) {
	switch op {
	case vm.OpAdd:
		return runtime.Add(bg, left, right), true
	case vm.OpSub:
		return runtime.Subtract(bg, left, right), true
	case vm.OpMulti:
		return runtime.Multiply(bg, left, right), true
	case vm.OpDiv:
		lv := runtime.ToNumberOnly(bg, left)
		if _, ok := lv.(runtime.Int); ok {
			rv := runtime.ToNumberOnly(bg, right)
			if ri, ok := rv.(runtime.Int); ok && ri == 0 {
				return nil, false
			}
		}
		return runtime.Divide(bg, left, right), true
	case vm.OpMod:
		if r, _ := runtime.ToInt(bg, right); r == 0 {
			return nil, false
		}
		return runtime.Modulus(bg, left, right), true
	case vm.OpCmp:
		return compare(bg, left, right), true
	case vm.OpEq:
		return equals(bg, left, right), true
	case vm.OpNe:
		return notEquals(bg, left, right), true
	case vm.OpGt:
		return greaterThan(bg, left, right), true
	case vm.OpLt:
		return lessThan(bg, left, right), true
	case vm.OpGte:
		return greaterThanOrEqual(bg, left, right), true
	case vm.OpLte:
		return lessThanOrEqual(bg, left, right), true
	default:
		return nil, false
	}
}

func buildConstIndex(constants []runtime.Value) map[string]int {
	index := make(map[string]int, len(constants))
	for i, val := range constants {
		if key, ok := constKey(val); ok {
			index[key] = i
		}
	}
	return index
}

func constKey(val runtime.Value) (string, bool) {
	if val == runtime.None {
		return "none", true
	}
	switch v := val.(type) {
	case runtime.Int:
		return fmt.Sprintf("i:%s", v.String()), true
	case runtime.Float:
		return fmt.Sprintf("f:%s", v.String()), true
	case runtime.String:
		return fmt.Sprintf("s:%s", v.String()), true
	case runtime.Boolean:
		if v {
			return "b:true", true
		}
		return "b:false", true
	default:
		return "", false
	}
}

func replaceWithConstLoad(inst *vm.Instruction, dst int, val runtime.Value, program *vm.Program, constIndex map[string]int) bool {
	newInst := buildConstLoad(dst, val, program, constIndex)
	if inst.Opcode == newInst.Opcode && inst.Operands == newInst.Operands {
		return false
	}
	*inst = newInst
	return true
}

func buildConstLoad(dst int, val runtime.Value, program *vm.Program, constIndex map[string]int) vm.Instruction {
	if val == runtime.None {
		return vm.NewInstruction(vm.OpLoadNone, vm.NewRegister(dst))
	}
	switch v := val.(type) {
	case runtime.Boolean:
		if v {
			return vm.NewInstruction(vm.OpLoadBool, vm.NewRegister(dst), vm.Operand(1))
		}
		return vm.NewInstruction(vm.OpLoadBool, vm.NewRegister(dst), vm.Operand(0))
	case runtime.Int:
		if v == 0 {
			return vm.NewInstruction(vm.OpLoadZero, vm.NewRegister(dst))
		}
	}

	key, ok := constKey(val)
	if !ok {
		return vm.NewInstruction(vm.OpLoadConst, vm.NewRegister(dst), vm.NewConstant(appendConst(program, constIndex, val)))
	}

	if idx, ok := constIndex[key]; ok {
		return vm.NewInstruction(vm.OpLoadConst, vm.NewRegister(dst), vm.NewConstant(idx))
	}

	return vm.NewInstruction(vm.OpLoadConst, vm.NewRegister(dst), vm.NewConstant(appendConst(program, constIndex, val)))
}

func appendConst(program *vm.Program, constIndex map[string]int, val runtime.Value) int {
	idx := len(program.Constants)
	program.Constants = append(program.Constants, val)
	if key, ok := constKey(val); ok {
		constIndex[key] = idx
	}
	return idx
}

func compare(_ context.Context, left, right runtime.Value) runtime.Int {
	return runtime.Int(runtime.CompareValues(nil, right, left))
}

func equals(_ context.Context, left, right runtime.Value) runtime.Boolean {
	return runtime.CompareValues(nil, left, right) == 0
}

func notEquals(_ context.Context, left, right runtime.Value) runtime.Boolean {
	return runtime.CompareValues(nil, left, right) != 0
}

func greaterThan(_ context.Context, left, right runtime.Value) runtime.Boolean {
	return runtime.CompareValues(nil, left, right) > 0
}

func greaterThanOrEqual(_ context.Context, left, right runtime.Value) runtime.Boolean {
	return runtime.CompareValues(nil, left, right) >= 0
}

func lessThan(_ context.Context, left, right runtime.Value) runtime.Boolean {
	return runtime.CompareValues(nil, left, right) < 0
}

func lessThanOrEqual(_ context.Context, left, right runtime.Value) runtime.Boolean {
	return runtime.CompareValues(nil, left, right) <= 0
}

func negate(input runtime.Value) runtime.Value {
	switch val := input.(type) {
	case runtime.Int:
		return -val
	case runtime.Float:
		return -val
	case runtime.Boolean:
		return !val
	default:
		return runtime.None
	}
}

func negative(input runtime.Value) runtime.Value {
	switch value := input.(type) {
	case runtime.Int:
		return -value
	case runtime.Float:
		return -value
	default:
		return runtime.None
	}
}

func positive(input runtime.Value) runtime.Value {
	switch value := input.(type) {
	case runtime.Int:
		return +value
	case runtime.Float:
		return +value
	default:
		return runtime.None
	}
}

func increment(ctx context.Context, input runtime.Value) runtime.Value {
	left := runtime.ToNumberOnly(ctx, input)
	switch value := left.(type) {
	case runtime.Int:
		return value + 1
	case runtime.Float:
		return value + 1
	default:
		return runtime.None
	}
}

func decrement(ctx context.Context, input runtime.Value) runtime.Value {
	left := runtime.ToNumberOnly(ctx, input)
	switch value := left.(type) {
	case runtime.Int:
		return value - 1
	case runtime.Float:
		return value - 1
	default:
		return runtime.None
	}
}
