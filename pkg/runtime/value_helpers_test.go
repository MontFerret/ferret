package runtime_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"math"
	"reflect"
	"strconv"
	"strings"
	"testing"
	"time"
	"unsafe"

	encodingjson "github.com/MontFerret/ferret/v2/pkg/encoding/json"
	"github.com/MontFerret/ferret/v2/pkg/runtime"

	. "github.com/smartystreets/goconvey/convey"
)

type CustomValue struct {
	properties map[runtime.Value]runtime.Value
}

type DummyStruct struct{}

type valueOfStruct struct {
	Exported string
	hidden   int
}

type valueOfFieldPayload struct {
	Items []any
}

func TestIsNil(t *testing.T) {
	Convey("Should match", t, func() {
		// nil == invalid
		t := runtime.IsNil(nil)

		So(t, ShouldBeTrue)

		a := []string{}
		t = runtime.IsNil(a)

		So(t, ShouldBeFalse)

		b := make([]string, 1)
		t = runtime.IsNil(b)

		So(t, ShouldBeFalse)

		c := make(map[string]string)
		t = runtime.IsNil(c)

		So(t, ShouldBeFalse)

		var s struct {
			Test string
		}
		t = runtime.IsNil(s)

		So(t, ShouldBeFalse)

		f := func() {}
		t = runtime.IsNil(f)

		So(t, ShouldBeFalse)

		i := DummyStruct{}
		t = runtime.IsNil(i)

		So(t, ShouldBeFalse)

		ch := make(chan string)
		t = runtime.IsNil(ch)

		So(t, ShouldBeFalse)

		var y unsafe.Pointer
		var vy int
		y = unsafe.Pointer(&vy)
		t = runtime.IsNil(y)

		So(t, ShouldBeFalse)
	})
}

func (t *CustomValue) MarshalJSON() ([]byte, error) {
	return nil, runtime.ErrNotImplemented
}

func (t *CustomValue) String() string {
	return ""
}

func (t *CustomValue) Unwrap() interface{} {
	return t
}

func (t *CustomValue) Hash() uint64 {
	return 0
}

func (t *CustomValue) Copy() runtime.Value {
	return runtime.None
}

func TestHelpers(t *testing.T) {
	Convey("Helpers", t, func() {
		Convey("Parse", func() {
			Convey("It should parse values", func() {
				inputs := []struct {
					Parsed runtime.Value
					Raw    interface{}
				}{
					{Parsed: runtime.NewInt(1), Raw: int(1)},
					{Parsed: runtime.NewInt(1), Raw: int8(1)},
					{Parsed: runtime.NewInt(1), Raw: int16(1)},
					{Parsed: runtime.NewInt(1), Raw: int32(1)},
					{Parsed: runtime.NewInt(1), Raw: int64(1)},
				}

				for _, input := range inputs {
					out := runtime.Parse(input.Raw)

					expectedType := reflect.TypeOf(input.Parsed)
					actualType := reflect.TypeOf(out)

					So(actualType, ShouldEqual, expectedType)
					So(out, ShouldEqual, input.Parsed)
				}
			})
		})

		Convey("ToBoolean", func() {
			Convey("Should convert values", func() {
				inputs := [][]runtime.Value{
					{
						runtime.None,
						runtime.False,
					},
					{
						runtime.True,
						runtime.True,
					},
					{
						runtime.False,
						runtime.False,
					},
					{
						runtime.NewInt(1),
						runtime.True,
					},
					{
						runtime.NewInt(0),
						runtime.False,
					},
					{
						runtime.NewFloat(1),
						runtime.True,
					},
					{
						runtime.NewFloat(0),
						runtime.False,
					},
					{
						runtime.NewString("Foo"),
						runtime.True,
					},
					{
						runtime.EmptyString,
						runtime.False,
					},
					{
						runtime.NewCurrentDateTime(),
						runtime.True,
					},
					{
						runtime.NewArray(1),
						runtime.True,
					},
					{
						runtime.NewObject(),
						runtime.True,
					},
					{
						runtime.NewBinary([]byte("")),
						runtime.True,
					},
				}

				for _, pair := range inputs {
					actual := runtime.ToBoolean(pair[0])
					expected := pair[1]

					So(actual, ShouldEqual, expected)
				}
			})
		})

		ctx := context.Background()

		Convey("ToFloat", func() {
			Convey("Should convert Int", func() {
				input := runtime.NewInt(100)
				output, err := runtime.ToFloat(ctx, input)

				So(err, ShouldBeNil)
				So(output, ShouldEqual, runtime.NewFloat(100))
			})

			Convey("Should convert Float", func() {
				input := runtime.NewFloat(100)
				output, err := runtime.ToFloat(ctx, input)

				So(err, ShouldBeNil)
				So(output, ShouldEqual, runtime.NewFloat(100))
			})

			Convey("Should convert String", func() {
				input := runtime.NewString("100.1")
				output, err := runtime.ToFloat(ctx, input)

				So(err, ShouldBeNil)
				So(output, ShouldEqual, runtime.NewFloat(100.1))

				output2, err := runtime.ToFloat(ctx, runtime.NewString("foobar"))
				So(err, ShouldNotBeNil)
				So(output2, ShouldEqual, runtime.ZeroFloat)
			})

			Convey("Should convert Boolean", func() {
				out, err := runtime.ToFloat(ctx, runtime.True)
				So(err, ShouldBeNil)
				So(out, ShouldEqual, runtime.NewFloat(1))

				out, err = runtime.ToFloat(ctx, runtime.False)
				So(err, ShouldBeNil)
				So(out, ShouldEqual, runtime.NewFloat(0))
			})

			Convey("Should convert Array with single item", func() {
				out, err := runtime.ToFloat(ctx, runtime.NewArrayWith(runtime.NewFloat(1)))
				So(err, ShouldBeNil)
				So(out, ShouldEqual, runtime.NewFloat(1))
			})

			Convey("Should convert Array with multiple items", func() {
				arg := runtime.NewArrayWith(runtime.NewFloat(1), runtime.NewFloat(1))

				out, err := runtime.ToFloat(ctx, arg)
				So(err, ShouldBeNil)
				So(out, ShouldEqual, runtime.NewFloat(2))
			})

			Convey("Should convert DateTime", func() {
				dt := runtime.NewCurrentDateTime()
				ts := dt.Unix()

				out, err := runtime.ToFloat(ctx, dt)
				So(err, ShouldBeNil)
				So(out, ShouldEqual, runtime.NewFloat(float64(ts)))
			})

			Convey("Should NOT convert other types", func() {
				inputs := []runtime.Value{
					runtime.NewObject(),
					runtime.NewBinary([]byte("")),
				}

				for _, input := range inputs {
					out, err := runtime.ToFloat(ctx, input)
					So(err, ShouldNotBeNil)
					So(out, ShouldEqual, runtime.ZeroFloat)
				}
			})
		})

		Convey("ToInt", func() {
			Convey("Should convert Int", func() {
				input := runtime.NewInt(100)
				output, err := runtime.ToInt(ctx, input)

				So(err, ShouldBeNil)
				So(output, ShouldEqual, runtime.NewInt(100))
			})

			Convey("Should convert Float", func() {
				input := runtime.NewFloat(100.1)
				output, err := runtime.ToInt(ctx, input)
				So(err, ShouldBeNil)
				So(output, ShouldEqual, runtime.NewInt(100))
			})

			Convey("Should convert String", func() {
				input := runtime.NewString("100")
				output, err := runtime.ToInt(ctx, input)

				So(err, ShouldBeNil)
				So(output, ShouldEqual, runtime.NewInt(100))

				output2, err := runtime.ToInt(ctx, runtime.NewString("foobar"))
				So(err, ShouldNotBeNil)
				So(output2, ShouldEqual, runtime.ZeroInt)
			})

			Convey("Should convert Boolean", func() {
				out, err := runtime.ToInt(ctx, runtime.True)
				So(err, ShouldBeNil)
				So(out, ShouldEqual, runtime.NewInt(1))

				out, err = runtime.ToInt(ctx, runtime.False)
				So(err, ShouldBeNil)
				So(out, ShouldEqual, runtime.NewInt(0))
			})

			Convey("Should convert Array with single item", func() {
				out, err := runtime.ToInt(ctx, runtime.NewArrayWith(runtime.NewInt(1)))
				So(err, ShouldBeNil)
				So(out, ShouldEqual, runtime.NewInt(1))
			})

			Convey("Should convert Array with multiple items", func() {
				arg := runtime.NewArrayWith(runtime.NewFloat(1), runtime.NewFloat(1))
				out, err := runtime.ToInt(ctx, arg)

				So(err, ShouldBeNil)

				So(out, ShouldEqual, runtime.NewFloat(2))
			})

			Convey("Should convert DateTime", func() {
				dt := runtime.NewCurrentDateTime()
				ts := dt.Unix()
				out, err := runtime.ToInt(ctx, dt)

				So(err, ShouldBeNil)
				So(out, ShouldEqual, runtime.NewInt(int(ts)))
			})

			Convey("Should NOT convert other types", func() {
				inputs := []runtime.Value{
					runtime.NewObject(),
					runtime.NewBinary([]byte("")),
				}

				for _, input := range inputs {
					out, err := runtime.ToInt(ctx, input)
					So(err, ShouldNotBeNil)
					So(out, ShouldEqual, runtime.ZeroInt)
				}
			})
		})

		Convey("ToList", func() {
			Convey("Should convert primitives", func() {
				dt := runtime.NewCurrentDateTime()

				inputs := [][]runtime.Value{
					{
						runtime.None,
						runtime.NewArray(0),
					},
					{
						runtime.True,
						runtime.NewArrayWith(runtime.True),
					},
					{
						runtime.NewInt(1),
						runtime.NewArrayWith(runtime.NewInt(1)),
					},
					{
						runtime.NewFloat(1),
						runtime.NewArrayWith(runtime.NewFloat(1)),
					},
					{
						runtime.NewString("foo"),
						runtime.NewArrayWith(runtime.NewString("foo")),
					},
					{
						dt,
						runtime.NewArrayWith(dt),
					},
				}

				for _, pairs := range inputs {
					actual, err := runtime.ToList(context.Background(), pairs[0])
					So(err, ShouldBeNil)
					expected := pairs[1]

					So(actual.Compare(expected), ShouldEqual, 0)
				}
			})
		})

		Convey("Unmarshal", func() {
			Convey("Should deserialize object", func() {
				input := map[string]interface{}{
					"foo": []string{
						"bar",
						"qaz",
					},
				}
				json1, err := json.Marshal(input)

				So(err, ShouldBeNil)

				val, err := encodingjson.Default.Decode(json1)

				So(err, ShouldBeNil)

				json2, err := encodingjson.Default.Encode(val)

				So(err, ShouldBeNil)
				So(json2, ShouldResemble, json1)
			})
		})
	})
}

func TestValueOf_Success(t *testing.T) {
	now := time.Date(2024, time.March, 30, 12, 0, 0, 0, time.UTC)
	passthrough := runtime.NewArrayWith(runtime.NewInt(1))
	pointedInt := 42
	var nilIntPtr *int

	tests := []struct {
		input  any
		want   runtime.Value
		sameAs runtime.Value
		name   string
	}{
		{
			name:  "nil",
			input: nil,
			want:  runtime.None,
		},
		{
			name:   "runtime value pass through",
			input:  passthrough,
			want:   passthrough,
			sameAs: passthrough,
		},
		{
			name:  "slice of runtime values",
			input: []runtime.Value{runtime.NewInt(1), runtime.NewString("two")},
			want:  runtime.NewArrayWith(runtime.NewInt(1), runtime.NewString("two")),
		},
		{
			name:  "bool",
			input: true,
			want:  runtime.True,
		},
		{
			name:  "string",
			input: "ferret",
			want:  runtime.NewString("ferret"),
		},
		{
			name:  "int",
			input: int(1),
			want:  runtime.NewInt(1),
		},
		{
			name:  "int8",
			input: int8(2),
			want:  runtime.NewInt(2),
		},
		{
			name:  "int16",
			input: int16(3),
			want:  runtime.NewInt(3),
		},
		{
			name:  "int32",
			input: int32(4),
			want:  runtime.NewInt(4),
		},
		{
			name:  "int64",
			input: int64(5),
			want:  runtime.NewInt(5),
		},
		{
			name:  "uint",
			input: uint(6),
			want:  runtime.NewInt(6),
		},
		{
			name:  "uint8",
			input: uint8(7),
			want:  runtime.NewInt(7),
		},
		{
			name:  "uint16",
			input: uint16(8),
			want:  runtime.NewInt(8),
		},
		{
			name:  "uint32",
			input: uint32(9),
			want:  runtime.NewInt(9),
		},
		{
			name:  "uint64",
			input: uint64(10),
			want:  runtime.NewInt(10),
		},
		{
			name:  "float32",
			input: float32(1.5),
			want:  runtime.NewFloat(1.5),
		},
		{
			name:  "float64",
			input: float64(2.5),
			want:  runtime.NewFloat(2.5),
		},
		{
			name:  "time.Time",
			input: now,
			want:  runtime.NewDateTime(now),
		},
		{
			name:  "slice of any",
			input: []any{1, "two", true},
			want:  runtime.NewArrayWith(runtime.NewInt(1), runtime.NewString("two"), runtime.True),
		},
		{
			name:  "map string any",
			input: map[string]any{"ok": true},
			want:  runtime.NewObjectWith(map[string]runtime.Value{"ok": runtime.True}),
		},
		{
			name:  "map any any",
			input: map[any]any{7: "seven"},
			want:  runtime.NewObjectWith(map[string]runtime.Value{"7": runtime.NewString("seven")}),
		},
		{
			name:  "bytes",
			input: []byte("bin"),
			want:  runtime.NewBinary([]byte("bin")),
		},
		{
			name:  "pointer dereference",
			input: &pointedInt,
			want:  runtime.NewInt(pointedInt),
		},
		{
			name:  "nil pointer",
			input: nilIntPtr,
			want:  runtime.None,
		},
		{
			name:  "reflect slice fallback",
			input: []string{"a", "b"},
			want:  runtime.NewArrayWith(runtime.NewString("a"), runtime.NewString("b")),
		},
		{
			name:  "reflect array fallback",
			input: [2]int{1, 2},
			want:  runtime.NewArrayWith(runtime.NewInt(1), runtime.NewInt(2)),
		},
		{
			name:  "reflect map fallback",
			input: map[int]string{1: "one"},
			want:  runtime.NewObjectWith(map[string]runtime.Value{"1": runtime.NewString("one")}),
		},
		{
			name:  "struct exported fields only",
			input: valueOfStruct{Exported: "visible", hidden: 10},
			want:  runtime.NewObjectWith(map[string]runtime.Value{"Exported": runtime.NewString("visible")}),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := runtime.ValueOf(tt.input)
			if err != nil {
				t.Fatalf("ValueOf(%T) returned unexpected error: %v", tt.input, err)
			}

			if tt.sameAs != nil && got != tt.sameAs {
				t.Fatalf("expected ValueOf(%T) to return the same instance", tt.input)
			}

			assertRuntimeValueEqual(t, got, tt.want)
		})
	}
}

func TestValueOf_Error(t *testing.T) {
	keyFailureChan := make(chan int)

	tests := []struct {
		name       string
		input      any
		target     error
		substrings []string
	}{
		{
			name:       "unsupported type",
			input:      make(chan int),
			target:     runtime.ErrInvalidType,
			substrings: []string{"cannot parse type chan int"},
		},
		{
			name:       "uint64 overflow",
			input:      uint64(math.MaxInt64) + 1,
			target:     runtime.ErrRange,
			substrings: []string{"invalid integer", "exceeds runtime.Int range"},
		},
		{
			name:       "slice any nested context",
			input:      []any{1, make(chan int)},
			target:     runtime.ErrInvalidType,
			substrings: []string{"cannot parse type chan int", "at index 1"},
		},
		{
			name:       "map string any nested context",
			input:      map[string]any{"items": []any{1, make(chan int)}},
			target:     runtime.ErrInvalidType,
			substrings: []string{"cannot parse type chan int", "at index 1", `at key "items"`},
		},
		{
			name:       "map any any nested context",
			input:      map[any]any{"items": []any{make(chan int)}},
			target:     runtime.ErrInvalidType,
			substrings: []string{"cannot parse type chan int", "at index 0", "at key items"},
		},
		{
			name:       "reflect slice nested context",
			input:      []chan int{make(chan int)},
			target:     runtime.ErrInvalidType,
			substrings: []string{"cannot parse type chan int", "at index 0"},
		},
		{
			name:       "reflect map value context",
			input:      map[int]any{7: make(chan int)},
			target:     runtime.ErrInvalidType,
			substrings: []string{"cannot parse type chan int", "at key 7"},
		},
		{
			name:       "reflect map key context",
			input:      map[chan int]string{keyFailureChan: "value"},
			target:     runtime.ErrInvalidType,
			substrings: []string{"cannot parse type chan int", "at key "},
		},
		{
			name:       "struct field context",
			input:      valueOfFieldPayload{Items: []any{make(chan int)}},
			target:     runtime.ErrInvalidType,
			substrings: []string{"cannot parse type chan int", "at index 0", `at field "Items"`},
		},
	}

	if strconv.IntSize == 64 {
		tests = append(tests, struct {
			name       string
			input      any
			target     error
			substrings []string
		}{
			name:       "uint overflow",
			input:      ^uint(0),
			target:     runtime.ErrRange,
			substrings: []string{"invalid integer", "exceeds runtime.Int range"},
		})
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := runtime.ValueOf(tt.input)
			assertWrappedError(t, err, tt.target, tt.substrings...)
		})
	}
}

func assertRuntimeValueEqual(t *testing.T, actual, expected runtime.Value) {
	t.Helper()

	if reflect.TypeOf(actual) != reflect.TypeOf(expected) {
		t.Fatalf("unexpected value type: got %T, want %T", actual, expected)
	}

	switch expectedValue := expected.(type) {
	case *runtime.Array:
		actualValue := actual.(*runtime.Array)
		if actualValue.Compare(expectedValue) != 0 {
			t.Fatalf("unexpected array value: got %s, want %s", actualValue, expectedValue)
		}
	case *runtime.Object:
		actualValue := actual.(*runtime.Object)
		if actualValue.Compare(expectedValue) != 0 {
			t.Fatalf("unexpected object value: got %s, want %s", actualValue, expectedValue)
		}
	case runtime.DateTime:
		actualValue := actual.(runtime.DateTime)
		if actualValue.Compare(expectedValue) != 0 {
			t.Fatalf("unexpected datetime value: got %s, want %s", actualValue, expectedValue)
		}
	case runtime.Binary:
		actualValue := actual.(runtime.Binary)
		if !bytes.Equal([]byte(actualValue), []byte(expectedValue)) {
			t.Fatalf("unexpected binary value: got %v, want %v", []byte(actualValue), []byte(expectedValue))
		}
	default:
		if actual != expected {
			t.Fatalf("unexpected value: got %v, want %v", actual, expected)
		}
	}
}

func assertWrappedError(t *testing.T, err error, target error, substrings ...string) {
	t.Helper()

	if err == nil {
		t.Fatal("expected ValueOf to fail")
	}

	if !errors.Is(err, target) {
		t.Fatalf("expected error to match %v, got %v", target, err)
	}

	errMsg := err.Error()
	start := 0

	for _, substring := range substrings {
		idx := strings.Index(errMsg[start:], substring)
		if idx == -1 {
			t.Fatalf("expected error %q to include substring %q after offset %d", errMsg, substring, start)
		}

		start += idx + len(substring)
	}
}
