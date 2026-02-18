package runtime_test

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

type encodeParams struct {
	Name    string `ferret:"name"`
	Age     int    `json:"age"`
	City    string
	Ignored string `ferret:"-"`
	private string `ferret:"private"`
}

type encodeInner struct {
	Value  string `ferret:"value"`
	Hidden string
}

type encodeOuter struct {
	Inner *encodeInner `ferret:"inner"`
	Count int          `ferret:"count"`
}

func TestEncode(t *testing.T) {
	Convey("Should encode tagged fields only", t, func() {
		input := encodeParams{
			Name:    "Alice",
			Age:     30,
			City:    "Paris",
			Ignored: "skip",
			private: "secret",
		}

		out := runtime.Encode(input)

		expected := runtime.NewObjectWith(
			runtime.NewObjectProperty("name", runtime.NewString("Alice")),
			runtime.NewObjectProperty("age", runtime.NewInt(30)),
		)

		So(out, ShouldResemble, expected)
	})

	Convey("Should encode nested tagged structs", t, func() {
		input := encodeOuter{
			Inner: &encodeInner{
				Value:  "ok",
				Hidden: "skip",
			},
			Count: 2,
		}

		out := runtime.Encode(input)

		expected := runtime.NewObjectWith(
			runtime.NewObjectProperty("inner", runtime.NewObjectWith(
				runtime.NewObjectProperty("value", runtime.NewString("ok")),
			)),
			runtime.NewObjectProperty("count", runtime.NewInt(2)),
		)

		So(out, ShouldResemble, expected)
	})
}

func TestEncodeByKey(t *testing.T) {
	type SomeValue struct {
		StrProp           string     `ferret:"strProp"`
		IntProp           int        `ferret:"intProp"`
		SliceProp         []int      `ferret:"sliceProp"`
		PointerProp       *SomeValue `ferret:"pointerProp"`
		JsonTag           string     `ferret:"jsonTag"`
		JsonAndRuntimeTag string     `ferret:"ferretTag" json:"jsonFerretTag"`

		UntaggedProp string

		privateStrProp string `ferret:"privateStrProp"`
	}

	type testCase struct {
		Name     string
		Input    SomeValue
		Field    string
		Expected runtime.Value
	}

	testCases := []testCase{
		{
			Name: "string",
			Input: SomeValue{
				StrProp: "test",
			},
			Field:    "strProp",
			Expected: runtime.String("test"),
		},
		{
			Name: "int",
			Input: SomeValue{
				IntProp: 99,
			},
			Field:    "intProp",
			Expected: runtime.Int(99),
		},
		{
			Name: "slice",
			Input: SomeValue{
				SliceProp: []int{1, 2, 3},
			},
			Field:    "sliceProp",
			Expected: runtime.NewArrayWith(runtime.Int(1), runtime.Int(2), runtime.Int(3)),
		},
		{
			Name: "pointer",
			Input: SomeValue{
				PointerProp: &SomeValue{
					StrProp: "test",
				},
			},
			Field: "pointerProp",
			Expected: runtime.NewObjectWith(
				runtime.NewObjectProperty("strProp", runtime.String("test")),
				runtime.NewObjectProperty("intProp", runtime.Int(0)),
				runtime.NewObjectProperty("sliceProp", runtime.NewArray(0)),
				runtime.NewObjectProperty("pointerProp", runtime.None),
				runtime.NewObjectProperty("jsonTag", runtime.EmptyString),
				runtime.NewObjectProperty("ferretTag", runtime.EmptyString),
			),
		},
		{
			Name: "json tag",
			Input: SomeValue{
				JsonTag: "json value",
			},
			Field:    "jsonTag",
			Expected: runtime.String("json value"),
		},
		{
			Name: "json and runtime tag. ferret tag should take precedence",
			Input: SomeValue{
				JsonAndRuntimeTag: "json and runtime value",
			},
			Field:    "ferretTag",
			Expected: runtime.String("json and runtime value"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			Convey(tc.Name, t, func() {
				actual, err := runtime.EncodeField(t.Context(), tc.Input, runtime.String(tc.Field))

				So(err, ShouldBeNil)
				So(actual, ShouldResemble, tc.Expected)
			})
		})

	}

	Convey("Struct", t, func() {
		Convey("should not read a private tagged field from a struct", func() {
			sv := SomeValue{
				privateStrProp: "hello world",
			}

			res, err := runtime.EncodeField(nil, sv, runtime.String("privateStrProp"))

			So(res, ShouldEqual, runtime.None)
			So(err, ShouldBeNil)
		})

		Convey("should not read a non-tagged field from a struct", func() {
			sv := SomeValue{
				privateStrProp: "hello world",
			}
			actual, err := runtime.EncodeField(nil, sv, runtime.String("UntaggedProp"))

			So(err, ShouldBeNil)
			So(actual, ShouldEqual, runtime.None)
		})
	})
}
