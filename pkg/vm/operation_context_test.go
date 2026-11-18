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

	admissions := 0
	admitted := func() { admissions++ }
	event, err = execution.Resume(nil, DebugResumeContinue, nil, admitted)
	if event != nil || !errors.Is(err, runtime.ErrInvalidArgument) {
		t.Fatalf("resume event=%+v error=%v", event, err)
	}

	if admissions != 0 {
		t.Fatal("rejected resume admitted")
	}

	event, err = execution.Resume(t.Context(), DebugResumeContinue, nil, admitted)
	if err != nil || event.Reason != DebugStopCompleted || admissions != 1 {
		t.Fatalf("resume after rejected context event=%+v error=%v", event, err)
	}
}
