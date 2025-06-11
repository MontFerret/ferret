package vm_test

import (
	"testing"

	. "github.com/MontFerret/ferret/test/integration/base"
)

func TestArrayAllOperator(t *testing.T) {
	RunUseCases(t, []UseCase{
		// TODO: Implement
		SkipCase("RETURN [1,2,3] ALL IN [1,2,3]", true, "All elements are in"),
	})
}
