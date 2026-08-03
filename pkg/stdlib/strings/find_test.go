package strings_test

import (
	"context"
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/MontFerret/ferret/v2/pkg/runtime"

	"github.com/MontFerret/ferret/v2/pkg/stdlib/strings"
)

var findBenchmarkResult runtime.Value

func TestFindFirst(t *testing.T) {
	Convey("When args are not passed", t, func() {
		Convey("It should return an error", func() {
			var err error
			_, err = strings.FindFirst(context.Background())

			So(err, ShouldBeError)

			_, err = strings.FindFirst(
				context.Background(),
				runtime.NewString("foo"),
			)

			So(err, ShouldBeError)
		})
	})

	Convey("When args are strings", t, func() {
		Convey("FindFirst('foobarbaz', 'ba') should return 3", func() {
			out, _ := strings.FindFirst(
				context.Background(),
				runtime.NewString("foobarbaz"),
				runtime.NewString("ba"),
			)

			So(out, ShouldEqual, 3)
		})

		Convey("FindFirst('foobarbaz', 'ba', 4) should return 6", func() {
			out, _ := strings.FindFirst(
				context.Background(),
				runtime.NewString("foobarbaz"),
				runtime.NewString("ba"),
				runtime.NewInt(4),
			)

			So(out, ShouldEqual, 6)
		})

		Convey("FindFirst('foobarbaz', 'ba', 4) should return -1", func() {
			out, _ := strings.FindFirst(
				context.Background(),
				runtime.NewString("foobarbaz"),
				runtime.NewString("ba"),
				runtime.NewInt(7),
			)

			So(out, ShouldEqual, -1)
		})

		Convey("FindFirst('foobarbaz', 'ba', 0, 3) should return -1", func() {
			out, _ := strings.FindFirst(
				context.Background(),
				runtime.NewString("foobarbaz"),
				runtime.NewString("ba"),
				runtime.NewInt(0),
				runtime.NewInt(3),
			)

			So(out, ShouldEqual, -1)
		})

		Convey("FindFirst('foobarbaz', 'ba', 4, 9) should return 6", func() {
			out, _ := strings.FindFirst(
				context.Background(),
				runtime.NewString("foobarbaz"),
				runtime.NewString("ba"),
				runtime.NewInt(4),
				runtime.NewInt(9),
			)

			So(out, ShouldEqual, 6)
		})
	})
}

func TestFindLast(t *testing.T) {
	Convey("When args are not passed", t, func() {
		Convey("It should return an error", func() {
			var err error
			_, err = strings.FindLast(context.Background())

			So(err, ShouldBeError)

			_, err = strings.FindLast(
				context.Background(),
				runtime.NewString("foo"),
			)

			So(err, ShouldBeError)
		})
	})

	Convey("When args are strings", t, func() {
		Convey("FindLast('foobarbaz', 'ba') should return 6", func() {
			out, _ := strings.FindLast(
				context.Background(),
				runtime.NewString("foobarbaz"),
				runtime.NewString("ba"),
			)

			So(out, ShouldEqual, 6)
		})

		Convey("FindLast('foobarbaz', 'ba', 7) should return -1", func() {
			out, _ := strings.FindLast(
				context.Background(),
				runtime.NewString("foobarbaz"),
				runtime.NewString("ba"),
				runtime.NewInt(7),
			)

			So(out, ShouldEqual, -1)
		})

		Convey("FindLast('foobarbaz', 'ba', 0, 5) should return 3", func() {
			out, _ := strings.FindLast(
				context.Background(),
				runtime.NewString("foobarbaz"),
				runtime.NewString("ba"),
				runtime.NewInt(0),
				runtime.NewInt(5),
			)

			So(out, ShouldEqual, 3)
		})

		Convey("FindLast('foobarbaz', 'ba', 4, 6) should return -1", func() {
			out, _ := strings.FindLast(
				context.Background(),
				runtime.NewString("foobarbaz"),
				runtime.NewString("ba"),
				runtime.NewInt(4),
				runtime.NewInt(6),
			)

			So(out, ShouldEqual, -1)
		})
	})
}

func TestFindUnicodeAndBounds(t *testing.T) {
	tests := []struct {
		name     string
		fn       func(context.Context, ...runtime.Value) (runtime.Value, error)
		args     []runtime.Value
		expected runtime.Int
	}{
		{
			name:     "first two arguments use character position",
			fn:       strings.FindFirst,
			args:     []runtime.Value{runtime.NewString("éaéa"), runtime.NewString("a")},
			expected: runtime.NewInt(1),
		},
		{
			name:     "first three arguments use character bounds",
			fn:       strings.FindFirst,
			args:     []runtime.Value{runtime.NewString("éaéa"), runtime.NewString("a"), runtime.NewInt(2)},
			expected: runtime.NewInt(3),
		},
		{
			name:     "first four arguments use character bounds",
			fn:       strings.FindFirst,
			args:     []runtime.Value{runtime.NewString("éaéaéa"), runtime.NewString("a"), runtime.NewInt(2), runtime.NewInt(6)},
			expected: runtime.NewInt(3),
		},
		{
			name:     "last two arguments use character position",
			fn:       strings.FindLast,
			args:     []runtime.Value{runtime.NewString("éaéa"), runtime.NewString("a")},
			expected: runtime.NewInt(3),
		},
		{
			name:     "last three arguments use character bounds",
			fn:       strings.FindLast,
			args:     []runtime.Value{runtime.NewString("éaéa"), runtime.NewString("a"), runtime.NewInt(2)},
			expected: runtime.NewInt(3),
		},
		{
			name:     "last four arguments use character bounds",
			fn:       strings.FindLast,
			args:     []runtime.Value{runtime.NewString("éaéaéa"), runtime.NewString("a"), runtime.NewInt(2), runtime.NewInt(6)},
			expected: runtime.NewInt(5),
		},
		{
			name:     "first clamps negative start and oversized end",
			fn:       strings.FindFirst,
			args:     []runtime.Value{runtime.NewString("éaéa"), runtime.NewString("a"), runtime.NewInt(-10), runtime.NewInt(100)},
			expected: runtime.NewInt(1),
		},
		{
			name:     "last clamps negative start and oversized end",
			fn:       strings.FindLast,
			args:     []runtime.Value{runtime.NewString("éaéa"), runtime.NewString("a"), runtime.NewInt(-10), runtime.NewInt(100)},
			expected: runtime.NewInt(3),
		},
		{
			name:     "first rejects reversed normalized bounds",
			fn:       strings.FindFirst,
			args:     []runtime.Value{runtime.NewString("éaéa"), runtime.NewString("a"), runtime.NewInt(3), runtime.NewInt(1)},
			expected: runtime.NewInt(-1),
		},
		{
			name:     "last rejects reversed normalized bounds",
			fn:       strings.FindLast,
			args:     []runtime.Value{runtime.NewString("éaéa"), runtime.NewString("a"), runtime.NewInt(3), runtime.NewInt(1)},
			expected: runtime.NewInt(-1),
		},
		{
			name:     "oversized start clamps to end",
			fn:       strings.FindFirst,
			args:     []runtime.Value{runtime.NewString("éa"), runtime.NewString("a"), runtime.NewInt(100)},
			expected: runtime.NewInt(-1),
		},
		{
			name:     "negative end clamps to zero",
			fn:       strings.FindLast,
			args:     []runtime.Value{runtime.NewString("éa"), runtime.NewString("a"), runtime.NewInt(0), runtime.NewInt(-1)},
			expected: runtime.NewInt(-1),
		},
		{
			name:     "first empty search returns effective start",
			fn:       strings.FindFirst,
			args:     []runtime.Value{runtime.NewString("éa"), runtime.EmptyString, runtime.NewInt(1), runtime.NewInt(100)},
			expected: runtime.NewInt(1),
		},
		{
			name:     "last empty search returns effective end",
			fn:       strings.FindLast,
			args:     []runtime.Value{runtime.NewString("éa"), runtime.EmptyString, runtime.NewInt(1), runtime.NewInt(100)},
			expected: runtime.NewInt(2),
		},
		{
			name:     "invalid start falls back to zero",
			fn:       strings.FindFirst,
			args:     []runtime.Value{runtime.NewString("éa"), runtime.NewString("a"), runtime.True},
			expected: runtime.NewInt(1),
		},
		{
			name:     "invalid end falls back to character length",
			fn:       strings.FindLast,
			args:     []runtime.Value{runtime.NewString("éaéa"), runtime.NewString("a"), runtime.ZeroInt, runtime.True},
			expected: runtime.NewInt(3),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result, err := test.fn(context.Background(), test.args...)
			if err != nil {
				t.Fatal(err)
			}

			if result != test.expected {
				t.Fatalf("expected %d, got %s", test.expected, result)
			}
		})
	}
}

func BenchmarkFind(b *testing.B) {
	ctx := context.Background()
	tests := []struct {
		name string
		fn   func(context.Context, ...runtime.Value) (runtime.Value, error)
		args []runtime.Value
	}{
		{
			name: "first_ascii",
			fn:   strings.FindFirst,
			args: []runtime.Value{runtime.NewString("foobarbaz"), runtime.NewString("ba")},
		},
		{
			name: "last_ascii",
			fn:   strings.FindLast,
			args: []runtime.Value{runtime.NewString("foobarbaz"), runtime.NewString("ba")},
		},
		{
			name: "first_unicode_bounded",
			fn:   strings.FindFirst,
			args: []runtime.Value{runtime.NewString("éaéa"), runtime.NewString("a"), runtime.NewInt(0), runtime.NewInt(4)},
		},
		{
			name: "last_unicode_bounded",
			fn:   strings.FindLast,
			args: []runtime.Value{runtime.NewString("éaéa"), runtime.NewString("a"), runtime.NewInt(0), runtime.NewInt(4)},
		},
	}

	for _, test := range tests {
		b.Run(test.name, func(b *testing.B) {
			b.ReportAllocs()

			for b.Loop() {
				result, err := test.fn(ctx, test.args...)
				if err != nil {
					b.Fatal(err)
				}

				findBenchmarkResult = result
			}
		})
	}
}
