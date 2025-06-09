package vm_test

import (
	. "github.com/MontFerret/ferret/test/integration/base"
	"testing"
)

func TestWaitforEvent(t *testing.T) {
	RunUseCases(t, []UseCase{
		CaseRuntimeError(`LET obj = {}

WAITFOR EVENT "test" IN obj

RETURN NONE`, "Should compile but return an error during execution because the object does not implement the interface"),
	})
}
