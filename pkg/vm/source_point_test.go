package vm

import (
	"context"
	"reflect"
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/bytecode"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
	"github.com/MontFerret/ferret/v2/pkg/source"
)

func TestRunTreatsSourcePointsAsNoopWithoutObserver(t *testing.T) {
	instance := mustNewVM(t, sourcePointTestProgram())
	result := mustRunResult(t, instance, nil)

	if got := mustResultRootAndClose(t, result); got != runtime.NewInt(1) {
		t.Fatalf("unexpected result: %v", got)
	}
}

func TestRunNotifiesSourcePointObserver(t *testing.T) {
	instance := mustNewVM(t, sourcePointTestProgram())
	observer := &recordingSourcePointObserver{}
	instance.sourcePointObserver = observer

	result := mustRunResult(t, instance, nil)
	_ = mustResultRootAndClose(t, result)

	if len(observer.states) != 2 {
		t.Fatalf("expected two observed source points, got %#v", observer.states)
	}
	for pointID, state := range observer.states {
		if state.pointID != pointID || state.pc != pointID*2 || state.depth != 0 {
			t.Fatalf("unexpected source point state %d: %#v", pointID, state)
		}
	}
}

func TestBuildExecPlanDoesNotTreatSourcePointIDAsRegister(t *testing.T) {
	program := newTestProgram(
		1,
		[]runtime.Value{runtime.NewString("TEST")},
		bytecode.NewInstruction(bytecode.OpLoadConst, bytecode.NewRegister(0), bytecode.NewConstant(0)),
		bytecode.NewInstruction(bytecode.OpSourcePoint, bytecode.Operand(0)),
		bytecode.NewInstruction(bytecode.OpHCall, bytecode.NewRegister(0)),
		bytecode.NewInstruction(bytecode.OpReturn, bytecode.NewRegister(0)),
	)

	if _, err := New(program); err != nil {
		t.Fatalf("source point cleared host function tracking: %v", err)
	}
}

func TestRegexpWarmupDoesNotTreatSourcePointIDAsRegister(t *testing.T) {
	const regexpPC = 2

	program := newTestProgram(
		3,
		[]runtime.Value{runtime.NewString("value")},
		bytecode.NewInstruction(bytecode.OpLoadConst, bytecode.NewRegister(2), bytecode.NewConstant(0)),
		bytecode.NewInstruction(bytecode.OpSourcePoint, bytecode.Operand(2)),
		bytecode.NewInstruction(bytecode.OpRegexp, bytecode.NewRegister(1), bytecode.NewRegister(0), bytecode.NewRegister(2)),
		bytecode.NewInstruction(bytecode.OpReturn, bytecode.NewRegister(1)),
	)
	instance := mustNewVM(t, program)

	if err := ensureRegexpsWarmed(instance); err != nil {
		t.Fatal(err)
	}
	if instance.cache.Regexps[regexpPC] == nil {
		t.Fatal("expected regexp after source point to be warmed")
	}
}

func TestNewDebugExecutionObservesSourcePointsWithoutMutatingPlan(t *testing.T) {
	instance := mustNewVM(t, sourcePointTestProgram())
	original := append([]execInstruction(nil), instance.plan.instructions...)

	execution, err := NewDebugExecution(instance, nil)
	if err != nil {
		t.Fatal(err)
	}
	defer execution.Close()

	if !reflect.DeepEqual(instance.plan.instructions, original) {
		t.Fatalf("debug execution mutated instruction plan: got %#v, want %#v", instance.plan.instructions, original)
	}
	if instance.sourcePointObserver == nil {
		t.Fatal("expected debug execution to install a source point observer")
	}
	if _, err := NewDebugExecution(instance, nil); err == nil {
		t.Fatal("expected duplicate debug execution configuration to fail")
	}

	event, err := execution.Start(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if event.Reason != DebugStopEntry || event.Point == nil || event.Point.PC != 0 {
		t.Fatalf("unexpected entry event: %#v", event)
	}

	event, err = execution.Resume(context.Background(), DebugResumeStep, nil)
	if err != nil {
		t.Fatal(err)
	}
	if event.Reason != DebugStopStep || event.Point == nil || event.Point.PC != 2 {
		t.Fatalf("unexpected step event: %#v", event)
	}

	event, err = execution.Resume(context.Background(), DebugResumeContinue, nil)
	if err != nil {
		t.Fatal(err)
	}
	if event.Reason != DebugStopCompleted {
		t.Fatalf("unexpected completion event: %#v", event)
	}
}

func sourcePointTestProgram() *bytecode.Program {
	return &bytecode.Program{
		ISAVersion: bytecode.Version,
		Registers:  2,
		Source:     source.New("source_point.fql", "LET x = 1\nRETURN x"),
		Bytecode: []bytecode.Instruction{
			bytecode.NewInstruction(bytecode.OpSourcePoint, bytecode.Operand(0)),
			bytecode.NewInstruction(bytecode.OpLoadConst, bytecode.NewRegister(1), bytecode.NewConstant(0)),
			bytecode.NewInstruction(bytecode.OpSourcePoint, bytecode.Operand(1)),
			bytecode.NewInstruction(bytecode.OpReturn, bytecode.NewRegister(1)),
		},
		Constants: []runtime.Value{runtime.NewInt(1)},
		Metadata: bytecode.Metadata{
			DebugPoints: []bytecode.DebugPoint{
				{PC: 0, Span: source.Span{Start: 0, End: 9}, FunctionID: -1},
				{PC: 2, Span: source.Span{Start: 10, End: 18}, FunctionID: -1},
			},
		},
	}
}
