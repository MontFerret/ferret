package strings_test

import (
	"context"
	"testing"

	"github.com/MontFerret/ferret/pkg/runtime/core"

	"github.com/MontFerret/ferret/pkg/runtime/values"
	"github.com/MontFerret/ferret/pkg/stdlib/strings"
	. "github.com/smartystreets/goconvey/convey"
)

type testCase struct {
	Name      string
	Expected  string
	Args      []core.Value
	ShouldErr bool
}

func (tc *testCase) Do(t *testing.T) {
	Convey(tc.Name, t, func() {
		expected := values.NewString(tc.Expected)

		formatted, err := strings.Fmt(context.Background(), tc.Args...)

		if tc.ShouldErr {
			So(err, ShouldBeError)
		} else {
			So(err, ShouldBeNil)
		}

		So(formatted, ShouldEqual, expected)
	})
}

func TestFmt(t *testing.T) {
	tcs := []*testCase{
		&testCase{
			Name:     `Should return "1"`,
			Expected: "1",
			Args: []core.Value{
				values.NewString("{}"),
				values.NewInt(1),
			},
		},
		&testCase{
			Name:     `Should return "2 1 1 2"`,
			Expected: "2 1 1 2",
			Args: []core.Value{
				values.NewString("{1} {} {0} {}"),
				values.NewInt(1),
				values.NewInt(2),
			},
		},
		&testCase{
			Name:     `Should return "Hello, World!"`,
			Expected: "Hello, World!",
			Args: []core.Value{
				values.NewString("{2}{1} {0}"),
				values.NewString("World!"),
				values.NewString(","),
				values.NewString("Hello"),
			},
		},
		&testCase{
			Name:     `Should return error [1]`,
			Expected: "{1}",
			Args: []core.Value{
				values.NewInt(10),
			},
			ShouldErr: true,
		},
	}

	for _, tc := range tcs {
		tc.Do(t)
	}
}
