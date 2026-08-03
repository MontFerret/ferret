package math_test

import (
	"context"
	"errors"
	stdmath "math"
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
	"github.com/MontFerret/ferret/v2/pkg/stdlib/math"

	. "github.com/smartystreets/goconvey/convey"
)

func TestRange(t *testing.T) {
	Convey("Should return range of numbers", t, func() {
		out, err := math.Range(context.Background(), runtime.NewInt(1), runtime.NewInt(4))

		So(err, ShouldBeNil)
		So(out.String(), ShouldEqual, "[1,2,3,4]")

		out, err = math.Range(context.Background(),
			runtime.NewInt(1),
			runtime.NewInt(4),
			runtime.NewInt(2))

		So(err, ShouldBeNil)
		So(out.String(), ShouldEqual, "[1,3]")

		out, err = math.Range(context.Background(),
			runtime.NewInt(1),
			runtime.NewInt(4),
			runtime.NewInt(3),
		)

		So(err, ShouldBeNil)
		So(out.String(), ShouldEqual, "[1,4]")

		out, err = math.Range(context.Background(),
			runtime.NewFloat(1.5),
			runtime.NewFloat(2.5),
		)

		So(err, ShouldBeNil)
		So(out.String(), ShouldEqual, "[1.5,2.5]")

		out, err = math.Range(context.Background(),
			runtime.NewFloat(1.5),
			runtime.NewFloat(2.5),
			runtime.NewFloat(0.5),
		)

		So(err, ShouldBeNil)
		So(out.String(), ShouldEqual, "[1.5,2,2.5]")

		out, err = math.Range(context.Background(),
			runtime.NewFloat(-0.75),
			runtime.NewFloat(1.1),
			runtime.NewFloat(0.5),
		)

		So(err, ShouldBeNil)
		So(out.String(), ShouldEqual, "[-0.75,-0.25,0.25,0.75]")
	})

	Convey("Should handle edge cases", t, func() {
		// Same start and end
		out, err := math.Range(context.Background(), runtime.NewInt(5), runtime.NewInt(5))

		So(err, ShouldBeNil)
		So(out.String(), ShouldEqual, "[5]")

		// Zero step should still work if only 2 args provided
		out, err = math.Range(context.Background(), runtime.NewInt(1), runtime.NewInt(3))

		So(err, ShouldBeNil)
		So(out.String(), ShouldEqual, "[1,2,3]")
	})

	Convey("Should return error for invalid arguments", t, func() {
		// Too few arguments
		out, err := math.Range(context.Background(), runtime.NewInt(1))

		So(err, ShouldNotBeNil)
		So(out, ShouldEqual, runtime.None)

		// Non-numeric first argument
		out, err = math.Range(context.Background(), runtime.NewString("invalid"), runtime.NewInt(4))

		So(err, ShouldNotBeNil)
		So(out, ShouldEqual, runtime.None)

		// Non-numeric second argument
		out, err = math.Range(context.Background(), runtime.NewInt(1), runtime.NewString("invalid"))

		So(err, ShouldNotBeNil)
		So(out, ShouldEqual, runtime.None)

		// Non-numeric step argument
		out, err = math.Range(context.Background(), runtime.NewInt(1), runtime.NewInt(4), runtime.NewString("invalid"))

		So(err, ShouldNotBeNil)
		So(out, ShouldEqual, runtime.None)
	})
}

func TestRangeDirectionAndSafety(t *testing.T) {
	ctx := context.Background()

	t.Run("ascending negative endpoints", func(t *testing.T) {
		out, err := math.Range(ctx, runtime.NewInt(-3), runtime.NewInt(-1))
		if err != nil {
			t.Fatal(err)
		}

		if got, want := out.String(), "[-3,-2,-1]"; got != want {
			t.Fatalf("Range() = %s, want %s", got, want)
		}
	})

	t.Run("descending ranges", func(t *testing.T) {
		tests := []struct {
			name string
			want string
			args []runtime.Value
		}{
			{name: "integer", want: "[3,2,1]", args: []runtime.Value{runtime.NewInt(3), runtime.NewInt(1), runtime.NewInt(-1)}},
			{name: "float", want: "[2.5,2,1.5]", args: []runtime.Value{runtime.NewFloat(2.5), runtime.NewFloat(1.5), runtime.NewFloat(-0.5)}},
		}

		for _, test := range tests {
			t.Run(test.name, func(t *testing.T) {
				out, err := math.Range(ctx, test.args...)
				if err != nil {
					t.Fatal(err)
				}

				if got := out.String(); got != test.want {
					t.Fatalf("Range() = %s, want %s", got, test.want)
				}
			})
		}
	})

	t.Run("mismatched direction", func(t *testing.T) {
		tests := [][]runtime.Value{
			{runtime.NewInt(1), runtime.NewInt(3), runtime.NewInt(-1)},
			{runtime.NewInt(3), runtime.NewInt(1), runtime.NewInt(1)},
		}

		for _, args := range tests {
			out, err := math.Range(ctx, args...)
			if err != nil {
				t.Fatal(err)
			}

			if got := out.String(); got != "[]" {
				t.Fatalf("Range(%v) = %s, want []", args, got)
			}
		}
	})

	t.Run("zero step", func(t *testing.T) {
		out, err := math.Range(ctx, runtime.NewInt(1), runtime.NewInt(3), runtime.ZeroInt)
		if !errors.Is(err, runtime.ErrInvalidArgument) {
			t.Fatalf("Range() error = %v, want ErrInvalidArgument", err)
		}
		if pos, ok, _ := runtime.InvalidArgumentDetails(err); !ok || pos != 2 {
			t.Fatalf("Range() argument position = %d, %t, want 2, true", pos, ok)
		}

		if out != runtime.None {
			t.Fatalf("Range() = %v, want None", out)
		}
	})

	t.Run("non-finite arguments", func(t *testing.T) {
		tests := []struct {
			name     string
			args     []runtime.Value
			position int
		}{
			{name: "nan start", args: []runtime.Value{runtime.NewFloat(stdmath.NaN()), runtime.NewInt(1)}, position: 0},
			{name: "infinite end", args: []runtime.Value{runtime.NewInt(0), runtime.NewFloat(stdmath.Inf(1))}, position: 1},
			{name: "nan step", args: []runtime.Value{runtime.NewInt(0), runtime.NewInt(1), runtime.NewFloat(stdmath.NaN())}, position: 2},
			{name: "infinite step", args: []runtime.Value{runtime.NewInt(0), runtime.NewInt(1), runtime.NewFloat(stdmath.Inf(-1))}, position: 2},
		}

		for _, test := range tests {
			t.Run(test.name, func(t *testing.T) {
				out, err := math.Range(ctx, test.args...)
				if !errors.Is(err, runtime.ErrInvalidArgument) {
					t.Fatalf("Range() error = %v, want ErrInvalidArgument", err)
				}
				if pos, ok, _ := runtime.InvalidArgumentDetails(err); !ok || pos != test.position {
					t.Fatalf("Range() argument position = %d, %t, want %d, true", pos, ok, test.position)
				}

				if out != runtime.None {
					t.Fatalf("Range() = %v, want None", out)
				}
			})
		}
	})

	t.Run("step cannot advance", func(t *testing.T) {
		start := 1e20
		end := stdmath.Nextafter(start, stdmath.Inf(1))
		out, err := math.Range(ctx, runtime.NewFloat(start), runtime.NewFloat(end), runtime.NewFloat(1))
		if !errors.Is(err, runtime.ErrInvalidArgument) {
			t.Fatalf("Range() error = %v, want ErrInvalidArgument", err)
		}
		if pos, ok, _ := runtime.InvalidArgumentDetails(err); !ok || pos != 2 {
			t.Fatalf("Range() argument position = %d, %t, want 2, true", pos, ok)
		}

		if out != runtime.None {
			t.Fatalf("Range() = %v, want None", out)
		}
	})

	t.Run("length exceeds array capacity", func(t *testing.T) {
		out, err := math.Range(
			ctx,
			runtime.ZeroFloat,
			runtime.NewFloat(stdmath.MaxFloat64),
			runtime.NewFloat(stdmath.SmallestNonzeroFloat64),
		)
		if !errors.Is(err, runtime.ErrRange) {
			t.Fatalf("Range() error = %v, want ErrRange", err)
		}

		if out != runtime.None {
			t.Fatalf("Range() = %v, want None", out)
		}
	})
}

var rangeBenchmarkResult runtime.Value

func BenchmarkRange(b *testing.B) {
	ctx := context.Background()
	tests := []struct {
		name string
		args []runtime.Value
	}{
		{name: "default_step", args: []runtime.Value{runtime.NewInt(1), runtime.NewInt(100)}},
		{name: "custom_step", args: []runtime.Value{runtime.NewInt(10), runtime.NewInt(100), runtime.NewInt(2)}},
	}

	for _, test := range tests {
		b.Run(test.name, func(b *testing.B) {
			b.ReportAllocs()

			for range b.N {
				result, err := math.Range(ctx, test.args...)
				if err != nil {
					b.Fatal(err)
				}

				rangeBenchmarkResult = result
			}
		})
	}
}
