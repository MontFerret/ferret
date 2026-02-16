package compiler_test

import (
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/bytecode"
)

func TestForWhileDistinct(t *testing.T) {
	RunUseCases(t, []UseCase{
		SkipByteCodeCase(`
			LET departments = []
			LET genders = []

			FOR i WHILE UNTIL(LENGTH(departments))
				FOR j WHILE UNTIL(LENGTH(genders))
					LET dept = departments[i]
					LET gender = genders[j]
					RETURN DISTINCT { department: dept, gender }
`, BC{
			I(bytecode.OpReturn, 0, 7),
		}),
	})
}
