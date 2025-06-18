package bytecode_test

import (
	"testing"

	"github.com/MontFerret/ferret/pkg/vm"
)

func TestCollectAggregate(t *testing.T) {
	RunUseCases(t, []UseCase{
		ByteCodeCase(`
LET users = []
FOR u IN users
  COLLECT genderGroup = u.gender
   AGGREGATE minAge = MIN(u.age), maxAge = MAX(u.age)

  RETURN {
    genderGroup,
    minAge,
    maxAge
  }
`, BC{
			I(vm.OpReturn, 0, 7),
		}),
	})
}
