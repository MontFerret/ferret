package math_test

import (
	"context"
	"errors"
	"fmt"
	stdmath "math"
	"slices"
	"strings"
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
	"github.com/MontFerret/ferret/v2/pkg/stdlib/math"
)

func TestClamp(t *testing.T) {
	negativeZero := runtime.Float(stdmath.Copysign(0, -1))
	positiveInf := runtime.Float(stdmath.Inf(1))
	negativeInf := runtime.Float(stdmath.Inf(-1))
	nan := runtime.Float(stdmath.Float64frombits(0x7ff8000000001234))
	negativeNaN := runtime.Float(stdmath.Float64frombits(0xfff8000000004321))
	for _, test := range []struct {
		value, min, max, want runtime.Value
		name                  string
	}{
		{name: "inside", value: runtime.Int(5), min: runtime.Int(0), max: runtime.Int(10), want: runtime.Int(5)},
		{name: "below", value: runtime.Int(-1), min: runtime.Int(0), max: runtime.Int(10), want: runtime.Int(0)},
		{name: "above", value: runtime.Int(20), min: runtime.Int(0), max: runtime.Int(10), want: runtime.Int(10)},
		{name: "at minimum", value: runtime.Int(0), min: runtime.Int(0), max: runtime.Int(10), want: runtime.Int(0)},
		{name: "at maximum", value: runtime.Int(10), min: runtime.Int(0), max: runtime.Int(10), want: runtime.Int(10)},
		{name: "equal bounds", value: runtime.Int(20), min: runtime.Int(3), max: runtime.Int(3), want: runtime.Int(3)},
		{name: "float inside", value: runtime.Float(0.25), min: runtime.Float(0.1), max: runtime.Float(0.5), want: runtime.Float(0.25)},
		{name: "float below", value: runtime.Float(0.05), min: runtime.Float(0.1), max: runtime.Float(0.5), want: runtime.Float(0.1)},
		{name: "float above", value: runtime.Float(0.75), min: runtime.Float(0.1), max: runtime.Float(0.5), want: runtime.Float(0.5)},
		{name: "mixed Int value", value: runtime.Int(5), min: runtime.Float(0.5), max: runtime.Float(10.5), want: runtime.Int(5)},
		{name: "mixed Float value", value: runtime.Float(5.5), min: runtime.Int(0), max: runtime.Int(10), want: runtime.Float(5.5)},
		{name: "mixed Int minimum", value: runtime.Float(-0.5), min: runtime.Int(0), max: runtime.Float(10.5), want: runtime.Int(0)},
		{name: "mixed Float minimum", value: runtime.Int(-1), min: runtime.Float(0.5), max: runtime.Int(10), want: runtime.Float(0.5)},
		{name: "mixed Int maximum", value: runtime.Float(20.5), min: runtime.Float(0.5), max: runtime.Int(10), want: runtime.Int(10)},
		{name: "mixed Float maximum", value: runtime.Int(20), min: runtime.Int(0), max: runtime.Float(10.5), want: runtime.Float(10.5)},
		{name: "Int equal to Float minimum", value: runtime.Int(1), min: runtime.Float(1), max: runtime.Float(10), want: runtime.Int(1)},
		{name: "Float equal to Int maximum", value: runtime.Float(10), min: runtime.Int(0), max: runtime.Int(10), want: runtime.Float(10)},
		{name: "below equal mixed bounds", value: runtime.Int(2), min: runtime.Int(3), max: runtime.Float(3), want: runtime.Int(3)},
		{name: "above equal mixed bounds", value: runtime.Int(4), min: runtime.Int(3), max: runtime.Float(3), want: runtime.Float(3)},
		{name: "Int equal to mixed bounds", value: runtime.Int(3), min: runtime.Float(3), max: runtime.Int(3), want: runtime.Int(3)},
		{name: "Float equal to mixed bounds", value: runtime.Float(3), min: runtime.Int(3), max: runtime.Float(3), want: runtime.Float(3)},
		{name: "large Int value", value: runtime.Int(9007199254740993), min: runtime.Int(0), max: runtime.Int(stdmath.MaxInt64), want: runtime.Int(9007199254740993)},
		{name: "large Int minimum", value: runtime.Float(9007199254740992), min: runtime.Int(9007199254740993), max: runtime.Int(stdmath.MaxInt64), want: runtime.Int(9007199254740993)},
		{name: "large Int maximum", value: runtime.Float(9007199254740994), min: runtime.Int(0), max: runtime.Int(9007199254740993), want: runtime.Int(9007199254740993)},
		{name: "large Float maximum", value: runtime.Int(9007199254740993), min: runtime.Int(0), max: runtime.Float(9007199254740992), want: runtime.Float(9007199254740992)},
		{name: "Int64 maximum", value: positiveInf, min: runtime.Int(0), max: runtime.Int(stdmath.MaxInt64), want: runtime.Int(stdmath.MaxInt64)},
		{name: "Int64 minimum", value: negativeInf, min: runtime.Int(stdmath.MinInt64), max: runtime.Int(0), want: runtime.Int(stdmath.MinInt64)},
		{name: "infinite bounds", value: runtime.Int(5), min: negativeInf, max: positiveInf, want: runtime.Int(5)},
		{name: "positive infinity", value: positiveInf, min: runtime.Int(0), max: runtime.Int(10), want: runtime.Int(10)},
		{name: "negative infinity", value: negativeInf, min: runtime.Int(0), max: runtime.Int(10), want: runtime.Int(0)},
		{name: "equal positive infinities", value: runtime.Int(5), min: positiveInf, max: positiveInf, want: positiveInf},
		{name: "equal negative infinities", value: runtime.Int(5), min: negativeInf, max: negativeInf, want: negativeInf},
		{name: "NaN value", value: nan, min: runtime.Int(0), max: runtime.Int(10), want: nan},
		{name: "negative NaN value", value: negativeNaN, min: runtime.Int(0), max: runtime.Int(10), want: negativeNaN},
		{name: "NaN and positive infinity", value: nan, min: positiveInf, max: positiveInf, want: nan},
		{name: "NaN and negative infinity", value: nan, min: negativeInf, max: negativeInf, want: nan},
		{name: "negative zero", value: negativeZero, min: runtime.Int(-1), max: runtime.Int(1), want: negativeZero},
		{name: "positive zero minimum", value: negativeZero, min: runtime.Float(0), max: runtime.Int(1), want: negativeZero},
		{name: "negative zero maximum", value: runtime.Float(0), min: runtime.Int(-1), max: negativeZero, want: runtime.Float(0)},
		{name: "selected negative zero minimum", value: runtime.Int(-1), min: negativeZero, max: runtime.Int(1), want: negativeZero},
		{name: "selected negative zero maximum", value: runtime.Int(1), min: runtime.Int(-1), max: negativeZero, want: negativeZero},
		{name: "equal signed zero bounds", value: negativeZero, min: runtime.Float(0), max: negativeZero, want: negativeZero},
	} {
		t.Run(test.name, func(t *testing.T) {
			got, err := math.Clamp(t.Context(), test.value, test.min, test.max)
			if err != nil {
				t.Fatal(err)
			}

			if expected, ok := test.want.(runtime.Float); ok {
				actual, ok := got.(runtime.Float)
				if !ok || stdmath.Float64bits(float64(actual)) != stdmath.Float64bits(float64(expected)) {
					t.Fatalf("result = %v (%T), want Float %v with identical bits", got, got, expected)
				}
			} else if got != test.want {
				t.Fatalf("result = %v (%T), want %v (%T)", got, got, test.want, test.want)
			}
		})
	}
}

func TestClampInvalidBounds(t *testing.T) {
	for _, test := range []struct {
		min      runtime.Value
		max      runtime.Value
		name     string
		position int
	}{
		{name: "reversed integers", min: runtime.Int(10), max: runtime.Int(0), position: 2},
		{name: "reversed floats", min: runtime.Float(0.5), max: runtime.Float(0.1), position: 2},
		{name: "reversed large integers", min: runtime.Int(9007199254740993), max: runtime.Int(9007199254740992), position: 2},
		{name: "reversed mixed bounds", min: runtime.Int(9007199254740993), max: runtime.Float(9007199254740992), position: 2},
		{name: "reversed negative mixed bounds", min: runtime.Float(-9007199254740992), max: runtime.Int(-9007199254740993), position: 2},
		{name: "reversed infinities", min: runtime.Float(stdmath.Inf(1)), max: runtime.Float(stdmath.Inf(-1)), position: 2},
		{name: "NaN minimum", min: runtime.NaN(), max: runtime.Int(10), position: 2},
		{name: "NaN maximum", min: runtime.Int(0), max: runtime.NaN(), position: 3},
		{name: "NaN bounds", min: runtime.NaN(), max: runtime.NaN(), position: 2},
	} {
		t.Run(test.name, func(t *testing.T) {
			for _, value := range []runtime.Value{runtime.Int(5), runtime.NaN()} {
				got, err := math.Clamp(t.Context(), value, test.min, test.max)
				assertScalarMathError(t, got, err, runtime.ErrInvalidArgument, test.position)
			}
		})
	}
}

func TestSign(t *testing.T) {
	testUnaryMath(t, math.Sign, []unaryMathCase{
		{runtime.Int(-12), runtime.Int(-1)},
		{runtime.Float(-8.5), runtime.Int(-1)},
		{runtime.Int(0), runtime.Int(0)},
		{runtime.Float(0), runtime.Int(0)},
		{runtime.Float(stdmath.Copysign(0, -1)), runtime.Int(0)},
		{runtime.Int(12), runtime.Int(1)},
		{runtime.Float(8.5), runtime.Int(1)},
		{runtime.Int(stdmath.MinInt64), runtime.Int(-1)},
		{runtime.Int(stdmath.MaxInt64), runtime.Int(1)},
		{runtime.Float(stdmath.SmallestNonzeroFloat64), runtime.Int(1)},
		{runtime.Float(-stdmath.SmallestNonzeroFloat64), runtime.Int(-1)},
		{runtime.Float(stdmath.Inf(1)), runtime.Int(1)},
		{runtime.Float(stdmath.Inf(-1)), runtime.Int(-1)},
	})

	got, err := math.Sign(t.Context(), runtime.NaN())
	assertScalarMathError(t, got, err, runtime.ErrInvalidArgument, 1)
}

func TestTrunc(t *testing.T) {
	negativeZero := runtime.Float(stdmath.Copysign(0, -1))
	testUnaryMath(t, math.Trunc, []unaryMathCase{
		{runtime.Float(3.9), runtime.Float(3)},
		{runtime.Float(-3.9), runtime.Float(-3)},
		{runtime.Int(3), runtime.Int(3)},
		{runtime.Int(-3), runtime.Int(-3)},
		{runtime.Float(3), runtime.Float(3)},
		{runtime.Int(9007199254740993), runtime.Int(9007199254740993)},
		{runtime.Int(stdmath.MaxInt64), runtime.Int(stdmath.MaxInt64)},
		{runtime.Int(stdmath.MinInt64), runtime.Int(stdmath.MinInt64)},
		{runtime.Float(stdmath.MaxFloat64), runtime.Float(stdmath.MaxFloat64)},
		{runtime.Int(0), runtime.Int(0)},
		{runtime.Float(0), runtime.Float(0)},
		{runtime.Float(0.5), runtime.Float(0)},
		{runtime.Float(-0.5), negativeZero},
		{negativeZero, negativeZero},
		{runtime.NaN(), runtime.NaN()},
		{runtime.Float(stdmath.Inf(1)), runtime.Float(stdmath.Inf(1))},
		{runtime.Float(stdmath.Inf(-1)), runtime.Float(stdmath.Inf(-1))},
	})
}

func TestCbrt(t *testing.T) {
	negativeZero := runtime.Float(stdmath.Copysign(0, -1))
	testUnaryMath(t, math.Cbrt, []unaryMathCase{
		{runtime.Int(27), runtime.Float(3)},
		{runtime.Int(-8), runtime.Float(-2)},
		{runtime.Int(0), runtime.Float(0)},
		{runtime.Float(0.125), runtime.Float(0.5)},
		{runtime.Int(2), runtime.Float(1.2599210498948732)},
		{negativeZero, negativeZero},
		{runtime.NaN(), runtime.NaN()},
		{runtime.Float(stdmath.Inf(1)), runtime.Float(stdmath.Inf(1))},
		{runtime.Float(stdmath.Inf(-1)), runtime.Float(stdmath.Inf(-1))},
	})
}

func TestHypot(t *testing.T) {
	for _, test := range []struct {
		x, y runtime.Value
		want runtime.Float
	}{
		{runtime.Int(3), runtime.Int(4), 5},
		{runtime.Int(0), runtime.Int(0), 0},
		{runtime.Int(-3), runtime.Float(-4), 5},
		{runtime.Float(3e200), runtime.Float(4e200), 5e200},
		{runtime.Float(3e-200), runtime.Float(4e-200), 5e-200},
		{runtime.Float(stdmath.Copysign(0, -1)), runtime.Float(0), 0},
		{runtime.NaN(), runtime.Int(0), runtime.NaN()},
		{runtime.Int(0), runtime.NaN(), runtime.NaN()},
		{runtime.Float(stdmath.Inf(-1)), runtime.Int(0), runtime.Float(stdmath.Inf(1))},
		{runtime.NaN(), runtime.Float(stdmath.Inf(1)), runtime.Float(stdmath.Inf(1))},
		{runtime.Float(stdmath.Inf(-1)), runtime.NaN(), runtime.Float(stdmath.Inf(1))},
	} {
		t.Run(fmt.Sprintf("%v/%v", test.x, test.y), func(t *testing.T) {
			got, err := math.Hypot(t.Context(), test.x, test.y)
			if err != nil {
				t.Fatal(err)
			}

			assertScalarMathResult(t, got, test.want)
		})
	}
}

func TestLog1p(t *testing.T) {
	negativeZero := runtime.Float(stdmath.Copysign(0, -1))
	testUnaryMath(t, math.Log1p, []unaryMathCase{
		{runtime.Int(0), runtime.Float(0)},
		{runtime.Float(1e-16), runtime.Float(1e-16)},
		{runtime.Float(-1e-16), runtime.Float(-1e-16)},
		{runtime.Int(1), runtime.Float(stdmath.Ln2)},
		{runtime.Int(-1), runtime.Float(stdmath.Inf(-1))},
		{runtime.Float(-1.5), runtime.NaN()},
		{negativeZero, negativeZero},
		{runtime.NaN(), runtime.NaN()},
		{runtime.Float(stdmath.Inf(1)), runtime.Float(stdmath.Inf(1))},
		{runtime.Float(stdmath.Inf(-1)), runtime.NaN()},
	})
}

func TestExpm1(t *testing.T) {
	negativeZero := runtime.Float(stdmath.Copysign(0, -1))
	testUnaryMath(t, math.Expm1, []unaryMathCase{
		{runtime.Int(0), runtime.Float(0)},
		{runtime.Float(1e-16), runtime.Float(1e-16)},
		{runtime.Float(-1e-16), runtime.Float(-1e-16)},
		{runtime.Int(1), runtime.Float(stdmath.E - 1)},
		{negativeZero, negativeZero},
		{runtime.NaN(), runtime.NaN()},
		{runtime.Float(stdmath.Inf(1)), runtime.Float(stdmath.Inf(1))},
		{runtime.Float(stdmath.Inf(-1)), runtime.Float(-1)},
	})
}

func TestE(t *testing.T) {
	got, err := math.E(t.Context())
	if err != nil {
		t.Fatal(err)
	}

	if want := runtime.Float(stdmath.E); got != want {
		t.Fatalf("result = %v (%T), want Float %v", got, got, want)
	}
}

func TestScalarPrimitivesRejectNonNumbers(t *testing.T) {
	type invocation struct {
		call func(context.Context, ...runtime.Value) (runtime.Value, error)
		name string
		args []runtime.Value
	}

	functions := []invocation{
		{name: "clamp", args: []runtime.Value{runtime.Int(5), runtime.Int(0), runtime.Int(10)}, call: func(ctx context.Context, args ...runtime.Value) (runtime.Value, error) {
			return math.Clamp(ctx, args[0], args[1], args[2])
		}},
		{name: "hypot", args: []runtime.Value{runtime.Int(3), runtime.Int(4)}, call: func(ctx context.Context, args ...runtime.Value) (runtime.Value, error) {
			return math.Hypot(ctx, args[0], args[1])
		}},
	}
	for name, fn := range map[string]runtime.Function1{"sign": math.Sign, "trunc": math.Trunc, "cbrt": math.Cbrt, "log1p": math.Log1p, "expm1": math.Expm1} {
		functions = append(functions, invocation{name: name, args: []runtime.Value{runtime.Int(1)}, call: func(ctx context.Context, args ...runtime.Value) (runtime.Value, error) {
			return fn(ctx, args[0])
		}})
	}

	for _, fn := range functions {
		for pos := range fn.args {
			for _, invalid := range []runtime.Value{runtime.None, runtime.String("2"), runtime.True, runtime.NewArray(0), runtime.NewObject(), runtime.Duration(1)} {
				t.Run(fmt.Sprintf("%s/%d/%T", fn.name, pos, invalid), func(t *testing.T) {
					args := slices.Clone(fn.args)
					args[pos] = invalid
					got, err := fn.call(t.Context(), args...)
					assertScalarMathError(t, got, err, runtime.ErrInvalidType, pos+1)
					if !strings.Contains(err.Error(), "Int") || !strings.Contains(err.Error(), "Float") {
						t.Fatalf("error = %v, want expected numeric types", err)
					}
				})
			}
		}
	}

	got, err := math.Clamp(t.Context(), runtime.None, runtime.True, runtime.String("2"))
	assertScalarMathError(t, got, err, runtime.ErrInvalidType, 1)
	got, err = math.Hypot(t.Context(), runtime.None, runtime.True)
	assertScalarMathError(t, got, err, runtime.ErrInvalidType, 1)
}

type unaryMathCase struct {
	input runtime.Value
	want  runtime.Value
}

func testUnaryMath(t *testing.T, fn runtime.Function1, cases []unaryMathCase) {
	t.Helper()

	for index, test := range cases {
		t.Run(fmt.Sprintf("%d/%T/%v", index, test.input, test.input), func(t *testing.T) {
			got, err := fn(t.Context(), test.input)
			if err != nil {
				t.Fatal(err)
			}

			assertScalarMathResult(t, got, test.want)
		})
	}
}

func assertScalarMathResult(t *testing.T, got, want runtime.Value) {
	t.Helper()

	expected, floating := want.(runtime.Float)
	if !floating {
		if got != want {
			t.Fatalf("result = %v (%T), want %v (%T)", got, got, want, want)
		}

		return
	}

	actual, ok := got.(runtime.Float)
	if !ok {
		t.Fatalf("result = %v (%T), want Float %v", got, got, want)
	}

	x, y := float64(actual), float64(expected)
	if stdmath.IsNaN(y) {
		if !stdmath.IsNaN(x) {
			t.Fatalf("result = %v, want NaN", got)
		}

		return
	}

	if y == 0 || stdmath.IsInf(y, 0) {
		if x != y || stdmath.Signbit(x) != stdmath.Signbit(y) {
			t.Fatalf("result = %v, want %v with matching sign", got, want)
		}

		return
	}

	// Relative error must stay small even near zero; accepting zero would hide
	// cancellation in log(1+x)/exp(x)-1 and underflow in a naive hypot.
	if stdmath.IsNaN(x) || stdmath.IsInf(x, 0) || stdmath.Abs((x-y)/y) > 1e-14 {
		t.Fatalf("result = %.17g, want %.17g", x, y)
	}
}

func assertScalarMathError(t *testing.T, got runtime.Value, err, want error, position int) {
	t.Helper()

	if got != runtime.None || !errors.Is(err, want) || !strings.Contains(err.Error(), fmt.Sprintf("position %d", position)) {
		t.Fatalf("result = %v, error = %v; want None and %v at position %d", got, err, want, position)
	}
}
