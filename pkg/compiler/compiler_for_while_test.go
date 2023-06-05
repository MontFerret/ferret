package compiler_test

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/MontFerret/ferret/pkg/compiler"
)

func TestForWhile(t *testing.T) {
	RunUseCases(t, compiler.New(), []UseCase{
		{
			"FOR i WHILE false RETURN i",
			[]any{},
			ShouldEqualJSON,
		},
	})
}
