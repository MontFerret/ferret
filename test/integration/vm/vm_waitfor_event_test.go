package vm_test

import (
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
	"github.com/MontFerret/ferret/v2/pkg/vm"
	"github.com/MontFerret/ferret/v2/test/spec"
	. "github.com/MontFerret/ferret/v2/test/spec/exec"
	. "github.com/MontFerret/ferret/v2/test/spec/mock"
)

func TestWaitforEvent(t *testing.T) {
	matchFirst := NewObservable([]runtime.Value{
		NewTestEventType("match"),
		NewTestEventType("other"),
	})
	matchSecond := NewObservable([]runtime.Value{
		NewTestEventType("other"),
		NewTestEventType("match"),
	})
	blocking := NewBlockingObservable()

	RunSpecs(t, []spec.Spec{
		Error(`LET obj = {}

WAITFOR EVENT "test" IN obj

RETURN NONE`, "Should compile but return an error during execution because the object does not implement the interface"),
		S(`LET obj = {}

WAITFOR EVENT "test" IN obj ON ERROR RETURN NONE

RETURN 1`, 1, "Statement suppression should continue after WAITFOR EVENT runtime failure"),
		S(`LET obj = {}

LET status = WAITFOR EVENT "test" IN obj TIMEOUT 1ms ON TIMEOUT RETURN "timeout" ON ERROR RETURN "error"

RETURN status`, "error", "WAITFOR EVENT should choose ON ERROR for runtime failures even when ON TIMEOUT is present"),
		S(`LET obj = {}

LET status = (WAITFOR EVENT "test" IN obj TIMEOUT 1ms) ON TIMEOUT RETURN "timeout" ON ERROR RETURN "error"

RETURN status`, "error", "Grouped WAITFOR EVENT should choose ON ERROR for runtime failures"),
		Fn(`LET obs = @obs
WAITFOR EVENT "test" IN obs WHEN .type == "match"
RETURN 1`, ObservableReturnOneAndReads(matchFirst, 1)).Env(vm.WithParams(map[string]runtime.Value{
			"obs": matchFirst,
		})),
		Fn(`LET obs = @obs
WAITFOR EVENT "test" IN obs WHEN .type == "match"
RETURN 1`, ObservableReturnOneAndReads(matchSecond, 2)).Env(vm.WithParams(map[string]runtime.Value{
			"obs": matchSecond,
		})),
		S(`LET obs = @obs

LET status = WAITFOR EVENT "test" IN obs TIMEOUT 1ms ON TIMEOUT RETURN "timeout" ON ERROR RETURN "error"

RETURN status`, "timeout", "WAITFOR EVENT should choose ON TIMEOUT when the stream times out").Env(vm.WithParams(map[string]runtime.Value{
			"obs": blocking,
		})),
		S(`LET obs = @obs

LET status = (WAITFOR EVENT "test" IN obs TIMEOUT 1ms) ON TIMEOUT RETURN "timeout" ON ERROR RETURN "error"

RETURN status`, "timeout", "Grouped WAITFOR EVENT should choose ON TIMEOUT when the stream times out").Env(vm.WithParams(map[string]runtime.Value{
			"obs": blocking,
		})),
		spec.NewSpec(`LET obs = @obs

RETURN WAITFOR EVENT "test" IN obs TIMEOUT 1ms ON ERROR RETURN "error"`, "WAITFOR EVENT timeout should not be caught by ON ERROR").Expect().ExecError(
			ShouldBeRuntimeError,
			&ExpectedRuntimeError{Contains: []string{runtime.ErrTimeout.Error()}},
		).Env(vm.WithParams(map[string]runtime.Value{
			"obs": blocking,
		})),
		spec.NewSpec(`LET obs = @obs

RETURN (WAITFOR EVENT "test" IN obs TIMEOUT 1ms) ON ERROR RETURN "error"`, "Grouped WAITFOR EVENT timeout should not be caught by ON ERROR").Expect().ExecError(
			ShouldBeRuntimeError,
			&ExpectedRuntimeError{Contains: []string{runtime.ErrTimeout.Error()}},
		).Env(vm.WithParams(map[string]runtime.Value{
			"obs": blocking,
		})),
	})
}
