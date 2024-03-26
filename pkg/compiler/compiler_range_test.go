package compiler_test

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestRange(t *testing.T) {
	RunUseCases(t, []UseCase{
		{
			"RETURN 1..10",
			[]any{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
			ShouldEqualJSON,
		},
		{
			"RETURN 10..1",
			[]any{10, 9, 8, 7, 6, 5, 4, 3, 2, 1},
			ShouldEqualJSON,
		},
		{
			`
		LET start = 1
		LET end = 10
		RETURN start..end
		`,
			[]any{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
			ShouldEqualJSON,
		},
		//{
		//	`
		//LET start = @start
		//LET end = @end
		//RETURN start..end
		//`,
		//	[]any{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
		//	ShouldEqualJSON,
		//},
	})
}
