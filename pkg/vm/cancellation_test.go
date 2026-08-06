package vm

import (
	"context"
	"errors"
	"io"
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/bytecode"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
	"github.com/MontFerret/ferret/v2/pkg/vm/internal/data"
)

func TestRunChecksCancellationBeforeFirstInstruction(t *testing.T) {
	program := newHostCallProgram(hostCallSpec{name: "CALL"})
	instance := mustNewVM(t, program)
	t.Cleanup(func() { _ = instance.Close() })

	calls := 0
	env := mustNewEnvironment(t, WithFunction("CALL", func(context.Context, ...runtime.Value) (runtime.Value, error) {
		calls++
		return runtime.None, nil
	}))
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	result, err := instance.Run(ctx, env)
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("run error = %v, want context.Canceled", err)
	}
	if result != nil {
		t.Fatal("canceled run returned a result")
	}
	if calls != 0 {
		t.Fatalf("host calls = %d, want zero", calls)
	}
}

func TestRunChecksCancellationBetweenInstructions(t *testing.T) {
	program := newHostCallProgram(
		hostCallSpec{name: "CANCEL"},
		hostCallSpec{name: "AFTER"},
	)
	instance := mustNewVM(t, program)
	t.Cleanup(func() { _ = instance.Close() })

	ctx, cancel := context.WithCancel(context.Background())
	cancelCompleted := false
	afterCalls := 0
	env := mustNewEnvironment(t,
		WithFunction("CANCEL", func(context.Context, ...runtime.Value) (runtime.Value, error) {
			cancel()
			cancelCompleted = true
			return runtime.NewInt(1), nil
		}),
		WithFunction("AFTER", func(context.Context, ...runtime.Value) (runtime.Value, error) {
			afterCalls++
			return runtime.NewInt(2), nil
		}),
	)

	result, err := instance.Run(ctx, env)
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("run error = %v, want context.Canceled", err)
	}
	if result != nil {
		t.Fatal("canceled run returned a result")
	}
	if !cancelCompleted {
		t.Fatal("canceling instruction did not finish")
	}
	if afterCalls != 0 {
		t.Fatalf("post-cancellation host calls = %d, want zero", afterCalls)
	}
}

func TestArrayDistinctDoesNotPollCancellation(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	result, err := arrayDistinct(ctx, runtime.NewArrayOf([]runtime.Value{
		runtime.NewInt(1),
		runtime.NewInt(1),
		runtime.NewInt(2),
	}))
	if err != nil {
		t.Fatalf("arrayDistinct() error = %v, want nil", err)
	}
	length, err := result.Length(ctx)
	if err != nil || length != 2 {
		t.Fatalf("distinct length = %d, %v, want 2, nil", length, err)
	}
}

func TestRunDefersCancellationAcrossCheapInstructions(t *testing.T) {
	probe := &cancellationProbeValue{}
	array := runtime.NewArrayWith(runtime.ZeroInt)
	program := newTestProgram(
		6,
		[]runtime.Value{
			runtime.NewString("prefix-"),
			probe,
			runtime.NewInt(1),
			array,
			runtime.ZeroInt,
		},
		bytecode.NewInstruction(bytecode.OpLoadConst, bytecode.NewRegister(0), bytecode.NewConstant(0)),
		bytecode.NewInstruction(bytecode.OpLoadConst, bytecode.NewRegister(1), bytecode.NewConstant(1)),
		bytecode.NewInstruction(bytecode.OpAdd, bytecode.NewRegister(2), bytecode.NewRegister(0), bytecode.NewRegister(1)),
		bytecode.NewInstruction(bytecode.OpLoadConst, bytecode.NewRegister(3), bytecode.NewConstant(2)),
		bytecode.NewInstruction(bytecode.OpEq, bytecode.NewRegister(4), bytecode.NewRegister(3), bytecode.NewRegister(3)),
		bytecode.NewInstruction(bytecode.OpLoadConst, bytecode.NewRegister(5), bytecode.NewConstant(3)),
		bytecode.NewInstruction(bytecode.OpSetIndexConst, bytecode.NewRegister(5), bytecode.NewConstant(4), bytecode.NewRegister(3)),
		bytecode.NewInstruction(bytecode.OpReturn, bytecode.NewRegister(2)),
	)

	err := runPreCanceledProgram(t, program)
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("run error = %v, want context.Canceled", err)
	}
	if probe.stringCalls != 1 {
		t.Fatalf("host string calls = %d, want 1", probe.stringCalls)
	}
	value, err := array.At(context.Background(), runtime.ZeroInt)
	if err != nil {
		t.Fatalf("read mutated array: %v", err)
	}
	if value != runtime.NewInt(1) {
		t.Fatalf("array value = %v, want 1", value)
	}
}

func TestRunDoesNotCheckCancellationAtForwardJump(t *testing.T) {
	probe := &cancellationProbeValue{}
	program := newTestProgram(
		3,
		[]runtime.Value{runtime.NewString("prefix-"), probe},
		bytecode.NewInstruction(bytecode.OpLoadConst, bytecode.NewRegister(0), bytecode.NewConstant(0)),
		bytecode.NewInstruction(bytecode.OpLoadConst, bytecode.NewRegister(1), bytecode.NewConstant(1)),
		bytecode.NewInstruction(bytecode.OpJump, bytecode.Operand(4)),
		bytecode.NewInstruction(bytecode.OpReturn, bytecode.NewRegister(0)),
		bytecode.NewInstruction(bytecode.OpAdd, bytecode.NewRegister(2), bytecode.NewRegister(0), bytecode.NewRegister(1)),
		bytecode.NewInstruction(bytecode.OpReturn, bytecode.NewRegister(2)),
	)

	err := runPreCanceledProgram(t, program)
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("run error = %v, want context.Canceled", err)
	}
	if probe.stringCalls != 1 {
		t.Fatalf("host string calls = %d, want 1", probe.stringCalls)
	}
}

func TestRunChecksCancellationAtBackwardJump(t *testing.T) {
	program := newTestProgram(
		1,
		nil,
		bytecode.NewInstruction(bytecode.OpLoadZero, bytecode.NewRegister(0)),
		bytecode.NewInstruction(bytecode.OpJump, bytecode.Operand(1)),
		bytecode.NewInstruction(bytecode.OpReturn, bytecode.NewRegister(0)),
	)

	err := runPreCanceledProgram(t, program)
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("run error = %v, want context.Canceled", err)
	}
}

func TestRunChecksEveryTakenBackwardEdge(t *testing.T) {
	tests := map[string]*bytecode.Program{
		"jump false": newTestProgram(
			1,
			nil,
			bytecode.NewInstruction(bytecode.OpLoadBool, bytecode.NewRegister(0), bytecode.Operand(0)),
			bytecode.NewInstruction(bytecode.OpJumpIfFalse, bytecode.Operand(1), bytecode.NewRegister(0)),
			bytecode.NewInstruction(bytecode.OpReturn, bytecode.NewRegister(0)),
		),
		"jump true": newTestProgram(
			1,
			nil,
			bytecode.NewInstruction(bytecode.OpLoadBool, bytecode.NewRegister(0), bytecode.Operand(1)),
			bytecode.NewInstruction(bytecode.OpJumpIfTrue, bytecode.Operand(1), bytecode.NewRegister(0)),
			bytecode.NewInstruction(bytecode.OpReturn, bytecode.NewRegister(0)),
		),
		"jump none": newTestProgram(
			1,
			nil,
			bytecode.NewInstruction(bytecode.OpLoadNone, bytecode.NewRegister(0)),
			bytecode.NewInstruction(bytecode.OpJumpIfNone, bytecode.Operand(1), bytecode.NewRegister(0)),
			bytecode.NewInstruction(bytecode.OpReturn, bytecode.NewRegister(0)),
		),
		"jump equal": comparisonJumpProgram(bytecode.OpJumpIfEq, runtime.NewInt(1), runtime.NewInt(1)),
		"jump not equal": comparisonJumpProgram(
			bytecode.OpJumpIfNe,
			runtime.NewInt(1),
			runtime.NewInt(2),
		),
		"jump equal const": comparisonConstJumpProgram(
			bytecode.OpJumpIfEqConst,
			runtime.NewInt(1),
			runtime.NewInt(1),
		),
		"jump not equal const": comparisonConstJumpProgram(
			bytecode.OpJumpIfNeConst,
			runtime.NewInt(1),
			runtime.NewInt(2),
		),
		"missing property":       missingPropertyJumpProgram(false),
		"missing property const": missingPropertyJumpProgram(true),
		"iterator skip": newTestProgram(
			2,
			[]runtime.Value{runtime.ZeroInt, runtime.NewInt(1)},
			bytecode.NewInstruction(bytecode.OpLoadConst, bytecode.NewRegister(0), bytecode.NewConstant(0)),
			bytecode.NewInstruction(bytecode.OpLoadConst, bytecode.NewRegister(1), bytecode.NewConstant(1)),
			bytecode.NewInstruction(bytecode.OpIterSkip, bytecode.Operand(2), bytecode.NewRegister(0), bytecode.NewRegister(1)),
			bytecode.NewInstruction(bytecode.OpReturn, bytecode.NewRegister(0)),
		),
		"iterator limit": newTestProgram(
			2,
			[]runtime.Value{runtime.NewInt(1), runtime.NewInt(1)},
			bytecode.NewInstruction(bytecode.OpLoadConst, bytecode.NewRegister(0), bytecode.NewConstant(0)),
			bytecode.NewInstruction(bytecode.OpLoadConst, bytecode.NewRegister(1), bytecode.NewConstant(1)),
			bytecode.NewInstruction(bytecode.OpIterLimit, bytecode.Operand(2), bytecode.NewRegister(0), bytecode.NewRegister(1)),
			bytecode.NewInstruction(bytecode.OpReturn, bytecode.NewRegister(0)),
		),
	}
	tests["match property"] = matchPropertyJumpProgram()

	for name, program := range tests {
		t.Run(name, func(t *testing.T) {
			err := runPreCanceledProgram(t, program)
			if !errors.Is(err, context.Canceled) {
				t.Fatalf("run error = %v, want context.Canceled", err)
			}
		})
	}
}

func TestRunChecksCancellationBeforeExternalCapability(t *testing.T) {
	probe := &cancellationProbeValue{}
	tests := []struct {
		calls func() int
		name  string
		inst  bytecode.Instruction
	}{
		{name: "load index", inst: bytecode.NewInstruction(bytecode.OpLoadIndex, bytecode.NewRegister(4), bytecode.NewRegister(0), bytecode.NewRegister(2)), calls: func() int { return probe.getCalls }},
		{name: "load key", inst: bytecode.NewInstruction(bytecode.OpLoadKey, bytecode.NewRegister(4), bytecode.NewRegister(0), bytecode.NewRegister(1)), calls: func() int { return probe.getCalls }},
		{name: "load property", inst: bytecode.NewInstruction(bytecode.OpLoadProperty, bytecode.NewRegister(4), bytecode.NewRegister(0), bytecode.NewRegister(1)), calls: func() int { return probe.getCalls }},
		{name: "set index", inst: bytecode.NewInstruction(bytecode.OpSetIndex, bytecode.NewRegister(0), bytecode.NewRegister(2), bytecode.NewRegister(3)), calls: func() int { return probe.setCalls }},
		{name: "set key", inst: bytecode.NewInstruction(bytecode.OpSetKey, bytecode.NewRegister(0), bytecode.NewRegister(1), bytecode.NewRegister(3)), calls: func() int { return probe.setCalls }},
		{name: "set property", inst: bytecode.NewInstruction(bytecode.OpSetProperty, bytecode.NewRegister(0), bytecode.NewRegister(1), bytecode.NewRegister(3)), calls: func() int { return probe.setCalls }},
		{name: "object set", inst: bytecode.NewInstruction(bytecode.OpObjectSet, bytecode.NewRegister(0), bytecode.NewRegister(1), bytecode.NewRegister(3)), calls: func() int { return probe.setCalls }},
		{name: "delete key", inst: bytecode.NewInstruction(bytecode.OpDeleteKey, bytecode.NewRegister(0), bytecode.NewRegister(1)), calls: func() int { return probe.removeCalls }},
		{name: "delete property", inst: bytecode.NewInstruction(bytecode.OpDeleteProperty, bytecode.NewRegister(0), bytecode.NewRegister(2)), calls: func() int { return probe.removeCalls }},
		{name: "append", inst: bytecode.NewInstruction(bytecode.OpPush, bytecode.NewRegister(0), bytecode.NewRegister(3)), calls: func() int { return probe.appendCalls }},
		{name: "append key value", inst: bytecode.NewInstruction(bytecode.OpPushKV, bytecode.NewRegister(0), bytecode.NewRegister(1), bytecode.NewRegister(3)), calls: func() int { return probe.setCalls }},
		{name: "length", inst: bytecode.NewInstruction(bytecode.OpLength, bytecode.NewRegister(4), bytecode.NewRegister(0)), calls: func() int { return probe.lengthCalls }},
		{name: "exists", inst: bytecode.NewInstruction(bytecode.OpExists, bytecode.NewRegister(4), bytecode.NewRegister(0)), calls: func() int { return probe.lengthCalls }},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := runPreCanceledProgram(t, externalCapabilityProgram(probe, test.inst))
			if !errors.Is(err, context.Canceled) {
				t.Fatalf("run error = %v, want context.Canceled", err)
			}
			if calls := test.calls(); calls != 0 {
				t.Fatalf("external capability calls = %d, want 0", calls)
			}
		})
	}
}

func TestRunChecksCancellationBeforeHostComparison(t *testing.T) {
	probe := &cancellationProbeValue{}
	program := newTestProgram(
		2,
		[]runtime.Value{probe},
		bytecode.NewInstruction(bytecode.OpLoadConst, bytecode.NewRegister(0), bytecode.NewConstant(0)),
		bytecode.NewInstruction(bytecode.OpEq, bytecode.NewRegister(1), bytecode.NewRegister(0), bytecode.NewRegister(0)),
		bytecode.NewInstruction(bytecode.OpReturn, bytecode.NewRegister(1)),
	)

	err := runPreCanceledProgram(t, program)
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("run error = %v, want context.Canceled", err)
	}
	if probe.equalCalls != 0 {
		t.Fatalf("host equality calls = %d, want 0", probe.equalCalls)
	}
}

func TestRunChecksCancellationBeforeIteratorCreation(t *testing.T) {
	probe := &cancellationProbeValue{}
	program := newTestProgram(
		2,
		[]runtime.Value{probe},
		bytecode.NewInstruction(bytecode.OpLoadConst, bytecode.NewRegister(0), bytecode.NewConstant(0)),
		bytecode.NewInstruction(bytecode.OpIter, bytecode.NewRegister(1), bytecode.NewRegister(0)),
		bytecode.NewInstruction(bytecode.OpReturn, bytecode.NewRegister(1)),
	)

	err := runPreCanceledProgram(t, program)
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("run error = %v, want context.Canceled", err)
	}
	if probe.iterateCalls != 0 {
		t.Fatalf("host iterate calls = %d, want 0", probe.iterateCalls)
	}
}

func TestRunChecksCancellationBeforeIteratorAdvancement(t *testing.T) {
	for _, op := range []bytecode.Opcode{bytecode.OpIterNext, bytecode.OpIterNextTimeout} {
		t.Run(op.String(), func(t *testing.T) {
			probe := &cancellationProbeValue{}
			iterator := data.NewIterator(probe)
			program := newTestProgram(
				2,
				[]runtime.Value{iterator},
				bytecode.NewInstruction(bytecode.OpLoadConst, bytecode.NewRegister(0), bytecode.NewConstant(0)),
				bytecode.NewInstruction(op, bytecode.Operand(2), bytecode.NewRegister(0), bytecode.NewRegister(1)),
				bytecode.NewInstruction(bytecode.OpReturn, bytecode.NewRegister(0)),
			)

			err := runPreCanceledProgram(t, program)
			if !errors.Is(err, context.Canceled) {
				t.Fatalf("run error = %v, want context.Canceled", err)
			}
			if probe.nextCalls != 0 {
				t.Fatalf("host next calls = %d, want 0", probe.nextCalls)
			}
		})
	}
}

func TestRunChecksCancellationTakenDuringIteratorBackwardEdge(t *testing.T) {
	for _, test := range []struct {
		nextErr error
		name    string
		op      bytecode.Opcode
	}{
		{name: "next eof", op: bytecode.OpIterNext, nextErr: io.EOF},
		{name: "timeout eof", op: bytecode.OpIterNextTimeout, nextErr: io.EOF},
		{name: "timeout", op: bytecode.OpIterNextTimeout, nextErr: runtime.ErrTimeout},
	} {
		t.Run(test.name, func(t *testing.T) {
			ctx, cancel := context.WithCancel(context.Background())
			probe := &cancellationProbeValue{cancel: cancel, nextErr: test.nextErr}
			iterator := data.NewIterator(probe)
			program := newTestProgram(
				5,
				[]runtime.Value{runtime.NewString("prefix-"), probe, iterator},
				bytecode.NewInstruction(bytecode.OpLoadConst, bytecode.NewRegister(0), bytecode.NewConstant(0)),
				bytecode.NewInstruction(bytecode.OpLoadConst, bytecode.NewRegister(1), bytecode.NewConstant(1)),
				bytecode.NewInstruction(bytecode.OpLoadConst, bytecode.NewRegister(3), bytecode.NewConstant(2)),
				bytecode.NewInstruction(bytecode.OpAdd, bytecode.NewRegister(2), bytecode.NewRegister(0), bytecode.NewRegister(1)),
				bytecode.NewInstruction(test.op, bytecode.Operand(3), bytecode.NewRegister(3), bytecode.NewRegister(4)),
				bytecode.NewInstruction(bytecode.OpReturn, bytecode.NewRegister(2)),
			)
			instance := mustNewVM(t, program)
			t.Cleanup(func() { _ = instance.Close() })

			result, err := instance.Run(ctx, NewDefaultEnvironment())
			if result != nil {
				_ = result.Close()
			}
			if !errors.Is(err, context.Canceled) {
				t.Fatalf("run error = %v, want context.Canceled", err)
			}
			if probe.stringCalls != 1 {
				t.Fatalf("host string calls = %d, want 1", probe.stringCalls)
			}
			if probe.nextCalls != 1 {
				t.Fatalf("iterator next calls = %d, want 1", probe.nextCalls)
			}
		})
	}
}

func TestRunChecksCancellationBeforeSleep(t *testing.T) {
	program := newTestProgram(
		1,
		[]runtime.Value{runtime.NewDuration(0)},
		bytecode.NewInstruction(bytecode.OpLoadConst, bytecode.NewRegister(0), bytecode.NewConstant(0)),
		bytecode.NewInstruction(bytecode.OpSleep, bytecode.NewRegister(0)),
		bytecode.NewInstruction(bytecode.OpReturn, bytecode.NewRegister(0)),
	)

	err := runPreCanceledProgram(t, program)
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("run error = %v, want context.Canceled", err)
	}
}

func TestRunChecksCancellationBeforeExternalOperations(t *testing.T) {
	probe := &cancellationProbeValue{}
	tests := []struct {
		program *bytecode.Program
		calls   func() int
		name    string
	}{
		{name: "dispatch", program: eventOperationProgram(bytecode.OpDispatch, probe), calls: func() int { return probe.dispatchCalls }},
		{name: "stream", program: eventOperationProgram(bytecode.OpStream, probe), calls: func() int { return probe.subscribeCalls }},
		{name: "query", program: queryOperationProgram(bytecode.OpQuery, probe), calls: func() int { return probe.queryCalls }},
		{name: "query exists", program: queryOperationProgram(bytecode.OpQueryExists, probe), calls: func() int { return probe.queryCalls }},
		{name: "query count", program: queryOperationProgram(bytecode.OpQueryCount, probe), calls: func() int { return probe.queryCalls }},
		{name: "query one", program: queryOperationProgram(bytecode.OpQueryOne, probe), calls: func() int { return probe.queryCalls }},
		{name: "close", program: closeOperationProgram(probe), calls: func() int { return probe.closeCalls }},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := runPreCanceledProgram(t, test.program)
			if !errors.Is(err, context.Canceled) {
				t.Fatalf("run error = %v, want context.Canceled", err)
			}
			if calls := test.calls(); calls != 0 {
				t.Fatalf("external operation calls = %d, want 0", calls)
			}
		})
	}
}

func TestRunChecksCancellationBeforeCollectionOpcode(t *testing.T) {
	probe := &cancellationProbeValue{}
	left := runtime.NewArrayWith(probe, probe)
	tests := map[string]bytecode.Instruction{
		"flatten":  bytecode.NewInstruction(bytecode.OpFlatten, bytecode.NewRegister(2), bytecode.NewRegister(0), bytecode.NewRegister(1)),
		"distinct": bytecode.NewInstruction(bytecode.OpDistinct, bytecode.NewRegister(2), bytecode.NewRegister(0)),
		"in":       bytecode.NewInstruction(bytecode.OpIn, bytecode.NewRegister(2), bytecode.NewRegister(1), bytecode.NewRegister(0)),
		"all":      bytecode.NewInstruction(bytecode.OpAllEq, bytecode.NewRegister(2), bytecode.NewRegister(0), bytecode.NewRegister(1)),
		"any":      bytecode.NewInstruction(bytecode.OpAnyEq, bytecode.NewRegister(2), bytecode.NewRegister(0), bytecode.NewRegister(1)),
		"none":     bytecode.NewInstruction(bytecode.OpNoneEq, bytecode.NewRegister(2), bytecode.NewRegister(0), bytecode.NewRegister(1)),
	}

	for name, operation := range tests {
		t.Run(name, func(t *testing.T) {
			program := newTestProgram(
				3,
				[]runtime.Value{left, probe},
				bytecode.NewInstruction(bytecode.OpLoadConst, bytecode.NewRegister(0), bytecode.NewConstant(0)),
				bytecode.NewInstruction(bytecode.OpLoadConst, bytecode.NewRegister(1), bytecode.NewConstant(1)),
				operation,
				bytecode.NewInstruction(bytecode.OpReturn, bytecode.NewRegister(2)),
			)

			err := runPreCanceledProgram(t, program)
			if !errors.Is(err, context.Canceled) {
				t.Fatalf("run error = %v, want context.Canceled", err)
			}
			if probe.equalCalls != 0 {
				t.Fatalf("host equality calls = %d, want 0", probe.equalCalls)
			}
		})
	}
}

func TestCollectionOpcodeCompletesBeforeNextSafepoint(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	probe := &cancellationProbeValue{cancel: cancel}
	program := newTestProgram(
		2,
		[]runtime.Value{runtime.NewArrayWith(probe, probe)},
		bytecode.NewInstruction(bytecode.OpLoadConst, bytecode.NewRegister(0), bytecode.NewConstant(0)),
		bytecode.NewInstruction(bytecode.OpDistinct, bytecode.NewRegister(1), bytecode.NewRegister(0)),
		bytecode.NewInstruction(bytecode.OpReturn, bytecode.NewRegister(1)),
	)
	instance := mustNewVM(t, program)
	t.Cleanup(func() { _ = instance.Close() })

	result, err := instance.Run(ctx, NewDefaultEnvironment())
	if result != nil {
		_ = result.Close()
	}
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("run error = %v, want context.Canceled", err)
	}
	if probe.equalCalls != 1 {
		t.Fatalf("host equality calls = %d, want 1", probe.equalCalls)
	}
}

func TestExecutionCancellationRecognizesDeadline(t *testing.T) {
	err := runtime.Errorf(runtime.ErrUnexpected, "wrapped: %v", context.DeadlineExceeded)
	if isExecutionCancellation(err) {
		t.Fatal("non-wrapping formatting unexpectedly preserved deadline identity")
	}

	err = errors.Join(errors.New("host failure"), context.DeadlineExceeded)
	if !isExecutionCancellation(err) {
		t.Fatal("wrapped deadline was not recognized")
	}
}

func TestCapabilitySafepointClassification(t *testing.T) {
	probe := &cancellationProbeValue{}
	for name, value := range map[string]runtime.Value{
		"array":       runtime.NewArray(0),
		"binary":      runtime.NewBinary(nil),
		"dataset":     data.NewDataSet(false),
		"fast object": data.NewFastObject(nil, 0),
		"int":         runtime.NewInt(1),
		"object":      runtime.NewObject(),
		"range":       runtime.NewRange(0, 1),
		"string":      runtime.NewString("value"),
	} {
		t.Run(name, func(t *testing.T) {
			if isExternalCapabilityReceiver(value) {
				t.Fatalf("%T classified as external", value)
			}
		})
	}
	if !isExternalCapabilityReceiver(probe) {
		t.Fatal("host capability classified as VM-owned")
	}
	if comparisonNeedsSafepoint(runtime.NewInt(1), runtime.NewFloat(1)) {
		t.Fatal("native scalar comparison classified as a safepoint")
	}
	if !comparisonNeedsSafepoint(runtime.NewArray(0), runtime.NewArray(0)) {
		t.Fatal("collection comparison was not classified as a safepoint")
	}
}

func runPreCanceledProgram(t *testing.T, program *bytecode.Program) error {
	t.Helper()

	instance := mustNewVM(t, program)
	t.Cleanup(func() { _ = instance.Close() })
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	result, err := instance.Run(ctx, NewDefaultEnvironment())
	if result != nil {
		_ = result.Close()
	}

	return err
}

func comparisonJumpProgram(op bytecode.Opcode, left, right runtime.Value) *bytecode.Program {
	return newTestProgram(
		2,
		[]runtime.Value{left, right},
		bytecode.NewInstruction(bytecode.OpLoadConst, bytecode.NewRegister(0), bytecode.NewConstant(0)),
		bytecode.NewInstruction(bytecode.OpLoadConst, bytecode.NewRegister(1), bytecode.NewConstant(1)),
		bytecode.NewInstruction(op, bytecode.Operand(2), bytecode.NewRegister(0), bytecode.NewRegister(1)),
		bytecode.NewInstruction(bytecode.OpReturn, bytecode.NewRegister(0)),
	)
}

func comparisonConstJumpProgram(op bytecode.Opcode, left, right runtime.Value) *bytecode.Program {
	return newTestProgram(
		1,
		[]runtime.Value{left, right},
		bytecode.NewInstruction(bytecode.OpLoadConst, bytecode.NewRegister(0), bytecode.NewConstant(0)),
		bytecode.NewInstruction(op, bytecode.Operand(1), bytecode.NewRegister(0), bytecode.NewConstant(1)),
		bytecode.NewInstruction(bytecode.OpReturn, bytecode.NewRegister(0)),
	)
}

func missingPropertyJumpProgram(constant bool) *bytecode.Program {
	object := runtime.NewObject()
	key := runtime.NewString("missing")
	if constant {
		return newTestProgram(
			1,
			[]runtime.Value{object, key},
			bytecode.NewInstruction(bytecode.OpLoadConst, bytecode.NewRegister(0), bytecode.NewConstant(0)),
			bytecode.NewInstruction(bytecode.OpJumpIfMissingPropertyConst, bytecode.Operand(1), bytecode.NewRegister(0), bytecode.NewConstant(1)),
			bytecode.NewInstruction(bytecode.OpReturn, bytecode.NewRegister(0)),
		)
	}

	return newTestProgram(
		2,
		[]runtime.Value{object, key},
		bytecode.NewInstruction(bytecode.OpLoadConst, bytecode.NewRegister(0), bytecode.NewConstant(0)),
		bytecode.NewInstruction(bytecode.OpLoadConst, bytecode.NewRegister(1), bytecode.NewConstant(1)),
		bytecode.NewInstruction(bytecode.OpJumpIfMissingProperty, bytecode.Operand(2), bytecode.NewRegister(0), bytecode.NewRegister(1)),
		bytecode.NewInstruction(bytecode.OpReturn, bytecode.NewRegister(0)),
	)
}

func matchPropertyJumpProgram() *bytecode.Program {
	program := newTestProgram(
		2,
		[]runtime.Value{runtime.NewObject(), runtime.NewString("missing")},
		bytecode.NewInstruction(bytecode.OpLoadConst, bytecode.NewRegister(0), bytecode.NewConstant(0)),
		bytecode.NewInstruction(bytecode.OpMatchLoadPropertyConst, bytecode.NewRegister(1), bytecode.NewRegister(0), bytecode.NewConstant(1)),
		bytecode.NewInstruction(bytecode.OpReturn, bytecode.NewRegister(1)),
	)
	program.Metadata.MatchFailTargets = []int{0, 1, 0}

	return program
}

func eventOperationProgram(op bytecode.Opcode, target runtime.Value) *bytecode.Program {
	return newTestProgram(
		3,
		[]runtime.Value{target, runtime.NewString("event"), runtime.None},
		bytecode.NewInstruction(bytecode.OpLoadConst, bytecode.NewRegister(0), bytecode.NewConstant(0)),
		bytecode.NewInstruction(bytecode.OpLoadConst, bytecode.NewRegister(1), bytecode.NewConstant(1)),
		bytecode.NewInstruction(bytecode.OpLoadConst, bytecode.NewRegister(2), bytecode.NewConstant(2)),
		bytecode.NewInstruction(op, bytecode.NewRegister(0), bytecode.NewRegister(1), bytecode.NewRegister(2)),
		bytecode.NewInstruction(bytecode.OpReturn, bytecode.NewRegister(0)),
	)
}

func queryOperationProgram(op bytecode.Opcode, target runtime.Value) *bytecode.Program {
	descriptor := runtime.NewArrayWith(runtime.EmptyString, runtime.EmptyString, runtime.None, runtime.None)

	return newTestProgram(
		3,
		[]runtime.Value{target, descriptor},
		bytecode.NewInstruction(bytecode.OpLoadConst, bytecode.NewRegister(0), bytecode.NewConstant(0)),
		bytecode.NewInstruction(bytecode.OpLoadConst, bytecode.NewRegister(1), bytecode.NewConstant(1)),
		bytecode.NewInstruction(op, bytecode.NewRegister(2), bytecode.NewRegister(0), bytecode.NewRegister(1)),
		bytecode.NewInstruction(bytecode.OpReturn, bytecode.NewRegister(2)),
	)
}

func closeOperationProgram(target runtime.Value) *bytecode.Program {
	return newTestProgram(
		1,
		[]runtime.Value{target},
		bytecode.NewInstruction(bytecode.OpLoadConst, bytecode.NewRegister(0), bytecode.NewConstant(0)),
		bytecode.NewInstruction(bytecode.OpClose, bytecode.NewRegister(0)),
		bytecode.NewInstruction(bytecode.OpReturn, bytecode.NewRegister(0)),
	)
}

func externalCapabilityProgram(target runtime.Value, operation bytecode.Instruction) *bytecode.Program {
	return newTestProgram(
		5,
		[]runtime.Value{target, runtime.NewString("key"), runtime.ZeroInt, runtime.NewInt(1)},
		bytecode.NewInstruction(bytecode.OpLoadConst, bytecode.NewRegister(0), bytecode.NewConstant(0)),
		bytecode.NewInstruction(bytecode.OpLoadConst, bytecode.NewRegister(1), bytecode.NewConstant(1)),
		bytecode.NewInstruction(bytecode.OpLoadConst, bytecode.NewRegister(2), bytecode.NewConstant(2)),
		bytecode.NewInstruction(bytecode.OpLoadConst, bytecode.NewRegister(3), bytecode.NewConstant(3)),
		operation,
		bytecode.NewInstruction(bytecode.OpReturn, bytecode.NewRegister(4)),
	)
}
