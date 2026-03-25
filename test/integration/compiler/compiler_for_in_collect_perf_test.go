package compiler_test

import (
	"fmt"
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/bytecode"
	"github.com/MontFerret/ferret/v2/test/spec"
)

func collectProjectionValueDefOpcode(prog *bytecode.Program) (bytecode.Opcode, error) {
	pushKVIndex, ok := findFirstOpcodeIndex(prog.Bytecode, bytecode.OpPushKV)
	if !ok {
		return 0, fmt.Errorf("expected OpPushKV in bytecode")
	}

	valueReg := prog.Bytecode[pushKVIndex].Operands[2].Register()
	valueDef, ok := lastRegisterDefOpcodeBefore(prog.Bytecode, pushKVIndex, valueReg)
	if !ok {
		return 0, fmt.Errorf("expected defining opcode for collect projection register R%d", valueReg)
	}

	return valueDef, nil
}

func TestCollectProjectionO0UsesPlainMoveForAllVarsScopePacking(t *testing.T) {
	RunSpecs(t, []spec.Spec{
		ProgramCheck(`
LET users = [{ age: 1 }]

FOR i IN users
	COLLECT g = i.age INTO groups
	RETURN groups
`, func(prog *bytecode.Program) error {
			got, err := collectProjectionValueDefOpcode(prog)
			if err != nil {
				return err
			}
			if got != bytecode.OpMove {
				return fmt.Errorf("expected all-vars collect projection handoff to use MOVE, got %s", got.String())
			}

			return nil
		}, "all-vars scope packing uses move"),
	})
}

func TestCollectProjectionO0UsesPlainMoveForKeepProjectionObject(t *testing.T) {
	RunSpecs(t, []spec.Spec{
		ProgramCheck(`
LET users = [{ age: 1 }]

FOR i IN users
	COLLECT g = i.age INTO groups KEEP i
	RETURN groups
`, func(prog *bytecode.Program) error {
			got, err := collectProjectionValueDefOpcode(prog)
			if err != nil {
				return err
			}
			if got != bytecode.OpMove {
				return fmt.Errorf("expected KEEP collect projection handoff to use MOVE, got %s", got.String())
			}

			return nil
		}, "keep projection uses move"),
	})
}

func TestCollectProjectionO0UsesPlainMoveForTypedCustomProjection(t *testing.T) {
	RunSpecs(t, []spec.Spec{
		ProgramCheck(`
LET users = ["alice"]

FOR i IN users
	COLLECT g = i INTO groups = i + "1"
	RETURN groups
`, func(prog *bytecode.Program) error {
			got, err := collectProjectionValueDefOpcode(prog)
			if err != nil {
				return err
			}
			if got != bytecode.OpMove {
				return fmt.Errorf("expected typed custom collect projection handoff to use MOVE, got %s", got.String())
			}

			return nil
		}, "typed custom projection uses move"),
	})
}

func TestCollectProjectionO0KeepsTrackedMoveForUnknownCustomProjection(t *testing.T) {
	RunSpecs(t, []spec.Spec{
		ProgramCheck(`
FUNC project(v) => v

LET users = [1]

FOR i IN users
	COLLECT g = i INTO groups = project(i)
	RETURN groups
`, func(prog *bytecode.Program) error {
			got, err := collectProjectionValueDefOpcode(prog)
			if err != nil {
				return err
			}
			if got != bytecode.OpMoveTracked {
				return fmt.Errorf("expected unknown custom collect projection handoff to use MOVET, got %s", got.String())
			}

			return nil
		}, "unknown custom projection keeps movet"),
	})
}

func TestCollectAggregateGlobalPlanUsesAggregateUpdateOpcode(t *testing.T) {
	RunSpecs(t, []spec.Spec{
		ProgramCheck(`
LET users = [1, 2, 3]

FOR u IN users
	COLLECT AGGREGATE total = SUM(u)
	RETURN total
`, func(prog *bytecode.Program) error {
			if !hasOpcode(prog.Bytecode, bytecode.OpAggregateUpdate) {
				return fmt.Errorf("expected plan-backed global aggregation to use OpAggregateUpdate")
			}

			if hasOpcode(prog.Bytecode, bytecode.OpPushKV) {
				return fmt.Errorf("expected plan-backed global aggregation to avoid generic PushKV writes")
			}

			return nil
		}, "global aggregate uses aggregate update"),
	})
}

func TestCollectAggregateGlobalIntoUsesProjectionBufferArrayPush(t *testing.T) {
	RunSpecs(t, []spec.Spec{
		ProgramCheck(`
LET users = [1, 2, 3]

FOR u IN users
	COLLECT AGGREGATE total = SUM(u) INTO groups
	RETURN groups
`, func(prog *bytecode.Program) error {
			if !hasOpcode(prog.Bytecode, bytecode.OpAggregateUpdate) {
				return fmt.Errorf("expected global aggregate INTO to use OpAggregateUpdate")
			}

			if !hasOpcode(prog.Bytecode, bytecode.OpArrayPush) {
				return fmt.Errorf("expected global aggregate INTO to append projection rows into a hidden array")
			}

			if hasOpcode(prog.Bytecode, bytecode.OpPushKV) {
				return fmt.Errorf("expected global aggregate INTO to avoid pushing projection rows into the aggregate collector")
			}

			return nil
		}, "global aggregate into uses projection buffer"),
	})
}

func TestCollectProjectionCountUsesDedicatedCounterIncrement(t *testing.T) {
	RunSpecs(t, []spec.Spec{
		ProgramCheck(`
FOR i IN 1..10
	COLLECT WITH COUNT INTO total
	RETURN total
`, func(prog *bytecode.Program) error {
			if !hasOpcode(prog.Bytecode, bytecode.OpCounterInc) {
				return fmt.Errorf("expected COLLECT WITH COUNT INTO to use OpCounterInc")
			}

			if hasOpcode(prog.Bytecode, bytecode.OpPushKV) {
				return fmt.Errorf("expected COLLECT WITH COUNT INTO to avoid generic PushKV collector writes")
			}

			return nil
		}, "collect with count uses counter increment"),
	})
}
