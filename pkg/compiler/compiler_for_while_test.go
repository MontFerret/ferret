package compiler_test

import (
	"context"
	"testing"

	"github.com/MontFerret/ferret/pkg/runtime"
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/MontFerret/ferret/pkg/compiler"
)

func TestForWhile(t *testing.T) {
	var counter int64
	RunUseCases(t, compiler.New(), []UseCase{
		//{
		//	"FOR i WHILE false RETURN i",
		//	[]any{},
		//	ShouldEqualJSON,
		//},
		{
			"FOR i WHILE UNTIL(5) RETURN i",
			[]any{0, 1, 2, 3, 4},
			ShouldEqualJSON,
		},
	}, runtime.WithFunctions(map[string]core.Function{
		"UNTIL": func(ctx context.Context, args ...core.Value) (core.Value, error) {
			if counter < int64(values.ToInt(args[0])) {
				counter++

				return values.True, nil
			}

			return values.False, nil
		},
	}))
}
