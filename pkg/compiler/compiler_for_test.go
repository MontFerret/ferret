package compiler_test

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/MontFerret/ferret/pkg/compiler"
)

func TestFor(t *testing.T) {
	RunUseCases(t, compiler.New(), []UseCase{
		{
			"FOR i IN 1..5 RETURN i",
			[]any{1, 2, 3, 4, 5},
			ShouldEqualJSON,
		},
		{
			`FOR i IN 1..5
				                           LET x = i
				                           PRINT(x)
											RETURN i
				`,
			[]any{1, 2, 3, 4, 5},
			ShouldEqualJSON,
		},
		{
			`FOR val, counter IN 1..5
		                          LET x = val
		                          PRINT(counter)
									LET y = counter
									RETURN [x, y]
		`,
			[]any{[]any{1, 0}, []any{2, 1}, []any{3, 2}, []any{4, 3}, []any{5, 4}},
			ShouldEqualJSON,
		},
		{
			`FOR i IN [] RETURN i
		`,
			[]any{},
			ShouldEqualJSON,
		},
		{
			`FOR i IN [1, 2, 3] RETURN i
		`,
			[]any{1, 2, 3},
			ShouldEqualJSON,
		},
		// Should not allocate memory if NONE is a return statement
		//{
		//	`FOR i IN 0..100
		//		RETURN NONE`,
		//	[]any{},
		//	ShouldEqualJSON,
		//},
		{
			`FOR i, k IN [1, 2, 3] RETURN k`,
			[]any{0, 1, 2},
			ShouldEqualJSON,
		},
	})
}
