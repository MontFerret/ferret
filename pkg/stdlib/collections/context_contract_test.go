package collections_test

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
	"github.com/MontFerret/ferret/v2/pkg/stdlib/collections"
)

func TestCollectionContextStates(t *testing.T) {
	for _, kind := range []string{"Nil", "Open", "Closed", "Deadline", "Layered"} {
		for _, operation := range []string{"count", "distinct", "includes", "reverse"} {
			t.Run(kind+"/"+operation, func(t *testing.T) {
				ctx, cancel := collectionTestContext(kind)
				defer cancel()
				values := []runtime.Value{runtime.Int(1), runtime.Int(2), runtime.Int(2)}
				source := &scanSource{expected: ctx, values: values}
				list := &configuredList{Array: runtime.NewArrayWith(values...), expected: ctx, calls: make(map[string]int)}
				var got, want runtime.Value
				var err error
				switch operation {
				case "count":
					got, err = collections.Count(ctx, source)
					want = runtime.Int(3)
				case "distinct":
					got, err = collections.CountDistinct(ctx, source)
					want = runtime.Int(2)
				case "includes":
					got, err = collections.Includes(ctx, source, runtime.Int(2))
					want = runtime.True
				case "reverse":
					got, err = collections.Reverse(ctx, list)
					want = list.created
				}

				if expected := ctx.Err(); expected != nil {
					if !errors.Is(err, expected) || source.iterations != 0 || len(list.calls) != 0 || ctx.doneCalls != 0 {
						t.Fatalf("pre-canceled work: %v, %v, %+v, %+v", got, err, source, list)
					}

					return
				}

				if err != nil || got != want || ctx.doneCalls > 1 {
					t.Fatalf("got %v, %v; want %v; Done calls %d", got, err, want, ctx.doneCalls)
				}

				if operation == "reverse" {
					first, err := list.created.At(ctx, 0)
					if err != nil || first != runtime.Int(2) || list.created.closes != 0 || list.closes != 0 {
						t.Fatal("reversal order or ownership changed")
					}

					if err := list.created.Close(); err != nil {
						t.Fatal(err)
					}
				} else if source.closes != 1 || source.sourceCloses != 0 {
					t.Fatal("scan ownership changed")
				}
			})
		}
	}
}

func TestCollectionFastPathsDoNotAcquireDone(t *testing.T) {
	for _, operation := range []string{"count", "contains", "includes-string", "reverse-string", "reverse-empty", "measurable-only"} {
		t.Run(operation, func(t *testing.T) {
			ctx, cancel := collectionTestContext("Open")
			defer cancel()
			var err error
			switch operation {
			case "count":
				_, err = collections.Count(ctx, &measuredSource{scanSource: &scanSource{expected: ctx}, length: 2})
			case "contains":
				_, err = collections.Includes(ctx, &containableSource{scanSource: &scanSource{expected: ctx}, contains: true}, runtime.Int(2))
			case "includes-string":
				_, err = collections.Includes(ctx, runtime.String("123"), runtime.Int(123))
			case "reverse-string":
				_, err = collections.Reverse(ctx, runtime.String("abc"))
			case "reverse-empty":
				_, err = collections.Reverse(ctx, runtime.NewArray(0))
			case "measurable-only":
				source := &measurableOnly{}
				got, failure := collections.Count(ctx, source)
				if got != runtime.ZeroInt || !errors.Is(failure, runtime.ErrInvalidType) || !strings.Contains(failure.Error(), "position 1") || source.lengthCalls != 0 {
					t.Fatalf("measurable-only dispatch: %v, %v, %+v", got, failure, source)
				}
			}

			if err != nil || ctx.doneCalls != 0 {
				t.Fatalf("fast path: %v, Done calls %d", err, ctx.doneCalls)
			}
		})
	}
}

func TestLayeredScanCancellationPreservesCauses(t *testing.T) {
	primary, cleanup := errors.New("host operation"), errors.New("iterator cleanup")
	for _, operation := range []string{"count", "distinct", "includes"} {
		for _, stage := range []string{"create", "next", "equality", "eof", "close"} {
			if operation == "count" && stage == "equality" {
				continue
			}

			t.Run(operation+"/"+stage, func(t *testing.T) {
				ctx, cancel := collectionTestContext("Layered")
				defer cancel()
				first, second := &scanValue{}, &scanValue{}
				source := &scanSource{expected: ctx, values: []runtime.Value{first, second}, closeErr: cleanup}
				wantCloses := 1
				wantPrimary := true
				switch stage {
				case "create":
					source.onIterate, source.iterateErr = cancel, primary
					wantCloses = 0
				case "next":
					source.onNext = func(int) { cancel() }
					source.nextErr = primary
				case "equality":
					first.onEqual, first.equalErr = cancel, primary
				case "eof":
					source.values = nil
					source.onNext = func(int) { cancel() }
					wantPrimary = false
				case "close":
					source.onClose = cancel
					wantPrimary = false
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

				if !errors.Is(err, context.Canceled) || wantPrimary && !errors.Is(err, primary) || wantCloses == 1 && !errors.Is(err, cleanup) {
					t.Fatalf("lost failure: %v, %v", got, err)
				}

				if got != runtime.ZeroInt && got != runtime.False || source.closes != wantCloses || source.sourceCloses != 0 || first.closes+second.closes != 0 {
					t.Fatalf("partial result or ownership changed: %v, %+v", got, source)
				}
			})
		}
	}
}

func TestReverseLayeredCancellationBoundaries(t *testing.T) {
	cleanup := errors.New("destination cleanup")
	for _, stage := range []string{"Length", "New", "At", "Append"} {
		t.Run(stage, func(t *testing.T) {
			ctx, cancel := collectionTestContext("Layered")
			defer cancel()
			value := &scanValue{}
			source := &configuredList{Array: runtime.NewArrayWith(value), expected: ctx, calls: make(map[string]int), closeErr: cleanup}
			source.after = func(operation string) {
				if operation == stage {
					cancel()
				}
			}

			got, err := collections.Reverse(ctx, source)
			if got != runtime.None || !errors.Is(err, context.Canceled) || source.closes != 0 || value.closes+value.copies+value.clones != 0 {
				t.Fatalf("cancellation: %v, %v", got, err)
			}

			if stage == "Length" {
				if source.created != nil || source.calls["New"] != 0 {
					t.Fatal("constructed after cancellation")
				}
			} else if source.created == nil || source.created.closes != 1 || !errors.Is(err, cleanup) {
				t.Fatal("incomplete destination cleanup lost")
			}

			if stage == "New" && source.calls["At"] != 0 || stage == "At" && source.created.calls["Append"] != 0 {
				t.Fatal("host operation after cancellation")
			}
		})
	}
}

func collectionTestContext(kind string) (*doneCountingContext, context.CancelFunc) {
	ctx := context.Background()
	cancel := func() {}
	switch kind {
	case "Open", "Layered", "Closed":
		ctx, cancel = context.WithCancel(ctx)
		if kind == "Closed" {
			cancel()
		}
	case "Deadline":
		ctx, cancel = context.WithDeadline(ctx, time.Now().Add(-time.Second))
	}

	if kind == "Layered" {
		type key int

		for i := range 4 {
			ctx = context.WithValue(ctx, key(i), i)
		}
	}

	return &doneCountingContext{Context: ctx}, cancel
}
