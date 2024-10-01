package compiler_test

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestForLimit(t *testing.T) {
	RunUseCases(t, []UseCase{
		{
			`
			FOR i IN [ 1, 2, 3, 4, 1, 3 ]
				LIMIT 2
				RETURN i
		`,
			[]any{1, 2},
			ShouldEqualJSON,
		},
	})
}
