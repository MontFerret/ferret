package collections_test

import (
	"context"
	"errors"
	"math"
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
	"github.com/MontFerret/ferret/v2/pkg/stdlib/collections"
)

func TestCountingIterableCapabilities(t *testing.T) {
	for name, call := range map[string]func(context.Context, runtime.Value) (runtime.Value, error){
		"count": collections.Count, "distinct": collections.CountDistinct,
	} {
		t.Run(name, func(t *testing.T) {
			for _, source := range []runtime.Value{runtime.NewRange(1, 3), runtime.NewRange(3, 1), &scanSource{values: []runtime.Value{runtime.Int(1), runtime.Int(2), runtime.Int(3)}}} {
				got, err := call(t.Context(), source)
				if err != nil || got != runtime.Int(3) {
					t.Fatalf("%T: got %v, %v", source, got, err)
				}
			}

			for _, source := range []runtime.Value{runtime.String(""), runtime.String("abc"), runtime.Int(1), runtime.True, runtime.None, nil} {
				got, err := call(t.Context(), source)
				if got != runtime.ZeroInt || !errors.Is(err, runtime.ErrInvalidType) {
					t.Fatalf("invalid %T: got %v, %v", source, got, err)
				}
			}
		})
	}

	if _, err := collections.Count(t.Context(), runtime.NewRange(math.MinInt64, math.MaxInt64)); !errors.Is(err, runtime.ErrRange) {
		t.Fatalf("range overflow: %v", err)
	}
}

func TestCountMeasurementNeverTraverses(t *testing.T) {
	failure := errors.New("length failed")
	for _, stage := range []string{"success", "failure", "cancel", "pre-cancel"} {
		t.Run(stage, func(t *testing.T) {
			ctx, cancel := context.WithCancel(t.Context())
			defer cancel()
			source := &measuredSource{scanSource: &scanSource{expected: ctx}, length: 42}
			var want error
			switch stage {
			case "failure":
				source.lengthErr, want = failure, failure
			case "cancel":
				source.onLength, want = cancel, context.Canceled
			case "pre-cancel":
				cancel()
				want = context.Canceled
			}

			got, err := collections.Count(ctx, source)
			if !errors.Is(err, want) || (want == nil && got != runtime.Int(42)) || (want != nil && got != runtime.ZeroInt) {
				t.Fatalf("got %v, %v; want error %v", got, err, want)
			}

			calls := 1
			if stage == "pre-cancel" {
				calls = 0
			}

			if source.lengthCalls != calls || source.iterations != 0 || source.closes != 0 || source.sourceCloses != 0 {
				t.Fatalf("unexpected host work: %+v", source)
			}
		})
	}
}

func TestOwnedScans(t *testing.T) {
	primary, cleanup := errors.New("traversal"), errors.New("close")
	for name, call := range map[string]func(context.Context, runtime.Value) (runtime.Value, error){
		"count":    collections.Count,
		"distinct": collections.CountDistinct,
		"includes": func(ctx context.Context, source runtime.Value) (runtime.Value, error) {
			return collections.Includes(ctx, source, runtime.Int(2))
		},
	} {
		for _, stage := range []string{"empty", "success", "create", "next", "next+close", "close", "pre-cancel", "cancel", "cancel+close", "cancel-eof", "cancel-close"} {
			t.Run(name+"/"+stage, func(t *testing.T) {
				ctx, cancel := context.WithCancel(t.Context())
				defer cancel()
				source := &scanSource{values: []runtime.Value{runtime.Int(1), runtime.Int(2), runtime.Int(2)}, expected: ctx, failAt: 1}
				wantCloses, wantIterations := 1, 1
				var wantPrimary, wantClose error
				switch stage {
				case "empty":
					source.values = nil
				case "create":
					source.iterateErr, wantPrimary, wantCloses = primary, primary, 0
				case "next":
					source.nextErr, wantPrimary = primary, primary
				case "next+close":
					source.nextErr, source.closeErr, wantPrimary, wantClose = primary, cleanup, primary, cleanup
				case "close":
					source.closeErr, wantClose = cleanup, cleanup
				case "pre-cancel":
					cancel()
					wantPrimary, wantCloses, wantIterations = context.Canceled, 0, 0
				case "cancel", "cancel+close":
					source.onNext = func(index int) {
						if index == 1 {
							cancel()
						}
					}

					wantPrimary = context.Canceled
					if stage == "cancel+close" {
						source.closeErr, wantClose = cleanup, cleanup
					}
				case "cancel-eof":
					source.values = nil
					source.onNext = func(int) { cancel() }
					wantPrimary = context.Canceled
				case "cancel-close":
					source.onClose, source.closeErr, wantPrimary, wantClose = cancel, cleanup, context.Canceled, cleanup
				}

				got, err := call(ctx, source)
				if (wantPrimary != nil && !errors.Is(err, wantPrimary)) || (wantClose != nil && !errors.Is(err, wantClose)) {
					t.Fatalf("lost failure: %v; want %v, %v", err, wantPrimary, wantClose)
				}

				if wantPrimary == nil && wantClose == nil {
					want := runtime.Value(runtime.Int(3))
					if name == "distinct" {
						want = runtime.Int(2)
					}

					if name == "includes" {
						want = runtime.True
					}

					if stage == "empty" {
						want = runtime.ZeroInt
						if name == "includes" {
							want = runtime.False
						}
					}

					if err != nil || got != want {
						t.Fatalf("got %v, %v; want %v", got, err, want)
					}
				} else if got != runtime.ZeroInt && got != runtime.False {
					t.Fatalf("partial result: %v", got)
				}

				if source.closes != wantCloses || source.iterations != wantIterations || source.sourceCloses != 0 {
					t.Fatalf("ownership: iterations %d, closes %d, source closes %d", source.iterations, source.closes, source.sourceCloses)
				}

				if (stage == "cancel" || stage == "cancel+close" || (name == "includes" && stage == "success")) && source.nexts != 2 {
					t.Fatalf("did not stop promptly: %d next calls", source.nexts)
				}
			})
		}
	}
}

func TestIncludesDelegatesContains(t *testing.T) {
	failure := errors.New("membership failed")
	for _, stage := range []string{"success", "false", "failure", "cancel", "pre-cancel"} {
		t.Run(stage, func(t *testing.T) {
			ctx, cancel := context.WithCancel(t.Context())
			defer cancel()
			source := &containableSource{scanSource: &scanSource{expected: ctx}, contains: true}
			var want error
			switch stage {
			case "false":
				source.contains = false
			case "failure":
				source.containsErr, want = failure, failure
			case "cancel":
				source.onContains, want = cancel, context.Canceled
			case "pre-cancel":
				cancel()
				want = context.Canceled
			}

			got, err := collections.Includes(ctx, source, runtime.Int(2))
			if !errors.Is(err, want) || (want == nil && got != source.contains) {
				t.Fatalf("got %v, %v", got, err)
			}

			calls := 1
			if stage == "pre-cancel" {
				calls = 0
			}

			if source.containsCalls != calls || source.iterations != 0 || source.sourceCloses != 0 || source.closes != 0 {
				t.Fatalf("incorrect dispatch: %+v", source)
			}
		})
	}
}

func TestScanEqualityAndBorrowedValues(t *testing.T) {
	primary, cleanup := errors.New("equality"), errors.New("close")
	for _, operation := range []string{"count", "distinct", "includes"} {
		for _, stage := range []string{"success", "equality", "cancel-equality"} {
			if operation == "count" && stage != "success" {
				continue
			}

			t.Run(operation+"/"+stage, func(t *testing.T) {
				ctx, cancel := context.WithCancel(t.Context())
				defer cancel()
				first, second := &scanValue{}, &scanValue{}
				source := &scanSource{values: []runtime.Value{first, second}, expected: ctx}
				var want error
				if stage == "equality" {
					first.equalErr, want, source.closeErr = primary, primary, cleanup
				}

				if stage == "cancel-equality" {
					first.onEqual, want, source.closeErr = cancel, context.Canceled, cleanup
				}

				var got runtime.Value
				var err error
				switch operation {
				case "count":
					got, err = collections.Count(ctx, source)
				case "distinct":
					got, err = collections.CountDistinct(ctx, source)
				case "includes":
					got, err = collections.Includes(ctx, source, first)
				}

				if want == nil {
					if err != nil {
						t.Fatal(err)
					}

					if operation == "count" && got != runtime.Int(2) || operation == "distinct" && got != runtime.Int(1) || operation == "includes" && got != runtime.True {
						t.Fatalf("got %v", got)
					}
				} else if !errors.Is(err, want) || !errors.Is(err, cleanup) {
					t.Fatalf("lost errors: %v", err)
				}

				if source.closes != 1 || source.sourceCloses != 0 || first.closes+second.closes+first.copies+second.copies+first.clones+second.clones != 0 {
					t.Fatal("borrowed values were owned or iterator leaked")
				}
			})
		}
	}
}

func TestMapValuesAndStringCoercion(t *testing.T) {
	ctx := t.Context()
	object := runtime.NewObjectWith(map[string]runtime.Value{"key": runtime.Int(7), "other": runtime.Float(7)})
	got, err := collections.CountDistinct(ctx, object)
	if err != nil || got != runtime.Int(1) {
		t.Fatalf("distinct map values: %v, %v", got, err)
	}

	for _, tc := range []struct{ source, needle, want runtime.Value }{
		{object, runtime.String("key"), runtime.False}, {object, runtime.Int(7), runtime.True},
		{runtime.String("123"), runtime.Int(123), runtime.True},
	} {
		got, err := collections.Includes(ctx, tc.source, tc.needle)
		if err != nil || got != tc.want {
			t.Fatalf("membership: %v, %v", got, err)
		}
	}
}

func TestIncludesClosesAfterNonemptyExhaustion(t *testing.T) {
	source := &scanSource{values: []runtime.Value{runtime.Int(1), runtime.Int(2)}}
	got, err := collections.Includes(t.Context(), source, runtime.Int(3))
	if err != nil || got != runtime.False || source.iterations != 1 || source.nexts != 3 || source.closes != 1 || source.sourceCloses != 0 {
		t.Fatalf("exhaustion: %v, %v, %+v", got, err, source)
	}
}
