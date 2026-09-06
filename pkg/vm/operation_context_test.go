package vm

import (
	"errors"
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

func TestOperationContextsRejectNilBeforeStartingWork(t *testing.T) {
	instance := mustNewVM(t, sourcePointTestProgram())
	result, err := instance.Run(nil, nil)
	if result != nil || !errors.Is(err, runtime.ErrInvalidArgument) {
		t.Fatalf("run result=%v error=%v", result, err)
	}

	execution, err := NewDebugExecution(instance, nil)
	if err != nil {
		t.Fatal(err)
	}

	t.Cleanup(func() { _ = execution.Close() })

	event, err := execution.Start(nil)
	if event != nil || !errors.Is(err, runtime.ErrInvalidArgument) {
		t.Fatalf("start event=%+v error=%v", event, err)
	}

	event, err = execution.Start(t.Context())
	if err != nil || event.Reason != DebugStopEntry {
		t.Fatalf("start after rejected context event=%+v error=%v", event, err)
	}

	event, err = execution.Resume(nil, DebugResumeContinue, nil)
	if event != nil || !errors.Is(err, runtime.ErrInvalidArgument) {
		t.Fatalf("resume event=%+v error=%v", event, err)
	}

	event, err = execution.Resume(t.Context(), DebugResumeContinue, nil)
	if err != nil || event.Reason != DebugStopCompleted {
		t.Fatalf("resume after rejected context event=%+v error=%v", event, err)
	}
}
