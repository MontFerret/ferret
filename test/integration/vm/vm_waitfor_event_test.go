package vm_test

import (
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
	"github.com/MontFerret/ferret/v2/pkg/vm"
	. "github.com/MontFerret/ferret/v2/test/spec/exec"
)

func TestWaitforEvent(t *testing.T) {
	matchFirst := NewTestObservable([]runtime.Value{
		NewTestEventType("match"),
		NewTestEventType("other"),
	})
	matchSecond := NewTestObservable([]runtime.Value{
		NewTestEventType("other"),
		NewTestEventType("match"),
	})

	RunSpecs(t, []Spec{
		Error(`LET obj = {}

WAITFOR EVENT "test" IN obj

RETURN NONE`, "Should compile but return an error during execution because the object does not implement the interface"),
		Fn(`LET obs = @obs
WAITFOR EVENT "test" IN obs WHEN .type == "match"
RETURN 1`, ObservableReturnOneAndReads(matchFirst, 1)).Options(vm.WithParams(map[string]runtime.Value{
			"obs": matchFirst,
		})),
		Fn(`LET obs = @obs
WAITFOR EVENT "test" IN obs WHEN .type == "match"
RETURN 1`, ObservableReturnOneAndReads(matchSecond, 2)).Options(vm.WithParams(map[string]runtime.Value{
			"obs": matchSecond,
		})),
	})
}
