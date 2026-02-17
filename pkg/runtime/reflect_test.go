package runtime

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestReadValueByKey(t *testing.T) {
	type SomeValue struct {
		StrProp     string     `ferret:"strProp"`
		IntProp     int        `ferret:"intProp"`
		SliceProp   []int      `ferret:"sliceProp"`
		PointerProp *SomeValue `ferret:"pointerProp"`

		UntaggedProp string

		privateStrProp string `ferret:"privateStrProp"`
	}

	type testCase struct {
		Name     string
		Input    SomeValue
		Field    string
		Expected Value
	}

	testCases := []testCase{
		{
			Name: "string",
			Input: SomeValue{
				StrProp: "test",
			},
			Field:    "strProp",
			Expected: String("test"),
		},
		{
			Name: "int",
			Input: SomeValue{
				IntProp: 99,
			},
			Field:    "intProp",
			Expected: Int(99),
		},
		{
			Name: "slice",
			Input: SomeValue{
				SliceProp: []int{1, 2, 3},
			},
			Field:    "sliceProp",
			Expected: NewArrayWith(Int(1), Int(2), Int(3)),
		},
		{
			Name: "pointer",
			Input: SomeValue{
				PointerProp: &SomeValue{
					StrProp: "test",
				},
			},
			Field:    "pointerProp",
			Expected: NewObjectWith(),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			Convey(tc.Name, t, func() {
				actual, err := ReadValueByKey(t.Context(), tc.Input, String(tc.Field))

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

			_, err := ReadValueByKey(nil, sv, String("privateStrProp"))

			So(err, ShouldNotBeNil)
		})

		Convey("should not read a non-tagged field from a struct", func() {
			sv := SomeValue{
				privateStrProp: "hello world",
			}
			actual, err := ReadValueByKey(nil, sv, String("UntaggedProp"))

			So(err, ShouldBeNil)
			So(actual, ShouldEqual, None)
		})
	})
}
