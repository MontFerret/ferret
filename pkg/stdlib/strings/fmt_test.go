package strings_test

import (
	"context"
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
	"github.com/MontFerret/ferret/pkg/stdlib/strings"
)

type testCase struct {
	Name      string
	Expected  string
	Format    string
	Args      []core.Value
	ShouldErr bool
}

func TestFmt(t *testing.T) {
	tcs := []*testCase{
		&testCase{
			Name:     `FMT("{}", 1) return "1"`,
			Expected: "1",
			Format:   "{}",
			Args: []core.Value{
				values.NewInt(1),
			},
		},
		&testCase{
			Name:     `FMT("{1} {} {0} {}", 1, 2) return "2 1 1 2"`,
			Expected: "2 1 1 2",
			Format:   "{1} {} {0} {}",
			Args: []core.Value{
				values.NewInt(1),
				values.NewInt(2),
			},
		},
		&testCase{
			Name:     `FMT("{1} {} {0} {} {}", 1, 2, 3) return "2 1 1 2 3"`,
			Expected: "2 1 1 2 3",
			Format:   "{1} {} {0} {} {}",
			Args: []core.Value{
				values.NewInt(1),
				values.NewInt(2),
				values.NewInt(3),
			},
		},
		&testCase{
			Name:     `FMT("{2}{1} {0}", "World!", ",", "Hello") return "Hello, World!"`,
			Expected: "Hello, World!",
			Format:   "{2}{1} {0}",
			Args: []core.Value{
				values.NewString("World!"),
				values.NewString(","),
				values.NewString("Hello"),
			},
		},
		&testCase{
			Name:     `FMT({}, {key:"value"}) return "{"key":"value"}"`,
			Expected: `{"key":"value"}`,
			Format:   "{}",
			Args: []core.Value{
				values.NewObjectWith(
					values.NewObjectProperty(
						"key", values.NewString("value"),
					),
				),
			},
		},
		&testCase{
			Name:     `FMT({}, {key:"value"}) return "{"key":"value"}"`,
			Expected: `{"key":"value","yek":"eulav"}`,
			Format:   "{}",
			Args: []core.Value{
				values.NewObjectWith(
					values.NewObjectProperty(
						"key", values.NewString("value"),
					),
					values.NewObjectProperty(
						"yek", values.NewString("eulav"),
					),
				),
			},
		},
		&testCase{
			Name:     `FMT("string") return "string"`,
			Expected: "string",
			Format:   "string",
		},
		&testCase{
			Name:     `FMT("string") return "string"`,
			Expected: "string",
			Format:   "string",
			Args: []core.Value{
				values.NewInt(1),
			},
		},
		&testCase{
			Name:      `FMT("{}") return error`,
			Format:    "{}",
			Args:      []core.Value{},
			ShouldErr: true,
		},
		&testCase{
			Name:   `FMT("{1}", 10) return error`,
			Format: "{1}",
			Args: []core.Value{
				values.NewInt(10),
			},
			ShouldErr: true,
		},
		&testCase{
			Name:   `FMT("{1} {} {0} {}", 1, 2, 3) return error`,
			Format: "{1} {} {0} {}",
			Args: []core.Value{
				values.NewInt(1),
				values.NewInt(2),
				values.NewInt(3),
			},
			ShouldErr: true,
		},
	}

	for _, tc := range tcs {
		tc.Do(t)
	}
}

func (tc *testCase) Do(t *testing.T) {
	Convey(tc.Name, t, func() {
		var expected core.Value = values.NewString(tc.Expected)

		args := []core.Value{values.NewString(tc.Format)}
		args = append(args, tc.Args...)

		formatted, err := strings.Fmt(context.Background(), args...)

		if tc.ShouldErr {
			So(err, ShouldBeError)
			expected = values.None
		} else {
			So(err, ShouldBeNil)
		}

		So(formatted, ShouldEqual, expected)
	})
}
