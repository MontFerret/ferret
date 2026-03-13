package formatter_test

import (
	"fmt"
	"strings"
	"testing"

	convey "github.com/smartystreets/goconvey/convey"

	"github.com/MontFerret/ferret/v2/pkg/file"
	"github.com/MontFerret/ferret/v2/pkg/formatter"
)

func Case(expression, expected string, desc ...string) UseCase {
	normalized := expected

	if strings.HasSuffix(normalized, "\n\n") {
		// already ends with two newlines; nothing to do
	} else if strings.HasSuffix(normalized, "\n") {
		// ends with a single newline; add one more to make it two
		normalized += "\n"
	} else {
		// no trailing newline; add two
		normalized += "\n\n"
	}

	return NewCase(expression, normalized, convey.ShouldEqual, desc...)
}

func SkipCase(expression, expected string, desc ...string) UseCase {
	return Skip(Case(expression, expected, desc...))
}

func RunUseCases(t *testing.T, useCases []UseCase) {
	RunUseCasesWith(t, formatter.New(), useCases)
}

func RunUseCasesWith(t *testing.T, f *formatter.Formatter, useCases []UseCase) {
	for _, useCase := range useCases {
		name := useCase.String()
		skip := useCase.Skip

		t.Run("Formatter Test: "+name, func(t *testing.T) {
			if skip {
				t.Skip()

				return
			}

			convey.Convey(useCase.Expression, t, func() {
				var out strings.Builder

				err := f.Format(&out, file.NewSource("Test case", useCase.Expression))

				if useCase.DebugOutput {
					fmt.Println("Input:")
					fmt.Println(useCase.Expression)
					fmt.Println("")
					fmt.Println("Output:")
					fmt.Println(out.String())
					fmt.Println("")
				}

				convey.So(err, convey.ShouldBeNil)

				actual := out.String()

				for _, assertion := range useCase.Assertions {
					convey.So(actual, assertion, useCase.Expected)
				}
			})
		})
	}
}
