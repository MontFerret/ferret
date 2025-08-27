package strings_test

import (
	"context"
	"testing"

	"github.com/MontFerret/ferret/pkg/runtime"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/MontFerret/ferret/pkg/stdlib/strings"
)

type testCase struct {
	Name      string
	Expected  string
	Format    string
	Args      []runtime.Value
	ShouldErr bool
}

func TestFmt(t *testing.T) {
	tcs := []*testCase{
		&testCase{
			Name:     `FMT("{}", 1) return "1"`,
			Expected: "1",
			Format:   "{}",
			Args: []runtime.Value{
				runtime.NewInt(1),
			},
		},
		&testCase{
			Name:     `FMT("{1} {} {0} {}", 1, 2) return "2 1 1 2"`,
			Expected: "2 1 1 2",
			Format:   "{1} {} {0} {}",
			Args: []runtime.Value{
				runtime.NewInt(1),
				runtime.NewInt(2),
			},
		},
		&testCase{
			Name:     `FMT("{1} {} {0} {} {}", 1, 2, 3) return "2 1 1 2 3"`,
			Expected: "2 1 1 2 3",
			Format:   "{1} {} {0} {} {}",
			Args: []runtime.Value{
				runtime.NewInt(1),
				runtime.NewInt(2),
				runtime.NewInt(3),
			},
		},
		&testCase{
			Name:     `FMT("{2}{1} {0}", "World!", ",", "Hello") return "Hello, World!"`,
			Expected: "Hello, World!",
			Format:   "{2}{1} {0}",
			Args: []runtime.Value{
				runtime.NewString("World!"),
				runtime.NewString(","),
				runtime.NewString("Hello"),
			},
		},
		&testCase{
			Name:     `FMT({}, {key:"value"}) return "{"key":"value"}"`,
			Expected: `{"key":"value"}`,
			Format:   "{}",
			Args: []runtime.Value{
				runtime.NewObjectWith(
					runtime.NewObjectProperty(
						"key", runtime.NewString("value"),
					),
				),
			},
		},
		&testCase{
			Name:     `FMT({}, {key:"value"}) return "{"key":"value"}"`,
			Expected: `{"key":"value","yek":"eulav"}`,
			Format:   "{}",
			Args: []runtime.Value{
				runtime.NewObjectWith(
					runtime.NewObjectProperty(
						"key", runtime.NewString("value"),
					),
					runtime.NewObjectProperty(
						"yek", runtime.NewString("eulav"),
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
			Args: []runtime.Value{
				runtime.NewInt(1),
			},
		},
		&testCase{
			Name:      `FMT("{}") return error`,
			Format:    "{}",
			Args:      []runtime.Value{},
			ShouldErr: true,
		},
		&testCase{
			Name:   `FMT("{1}", 10) return error`,
			Format: "{1}",
			Args: []runtime.Value{
				runtime.NewInt(10),
			},
			ShouldErr: true,
		},
		&testCase{
			Name:   `FMT("{1} {} {0} {}", 1, 2, 3) return error`,
			Format: "{1} {} {0} {}",
			Args: []runtime.Value{
				runtime.NewInt(1),
				runtime.NewInt(2),
				runtime.NewInt(3),
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
		var expected runtime.Value = runtime.NewString(tc.Expected)

		args := []runtime.Value{runtime.NewString(tc.Format)}
		args = append(args, tc.Args...)

		formatted, err := strings.Fmt(context.Background(), args...)

		if tc.ShouldErr {
			So(err, ShouldBeError)
			expected = runtime.None
		} else {
			So(err, ShouldBeNil)
		}

		So(formatted, ShouldEqual, expected)
	})
}
