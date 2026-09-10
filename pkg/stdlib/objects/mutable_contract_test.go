package objects_test

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
	"github.com/MontFerret/ferret/v2/pkg/stdlib/objects"
	objectcases "github.com/MontFerret/ferret/v2/test/spec/objects"
)

func TestMutableTransformationFixtures(t *testing.T) {
	ctx := context.Background()
	for _, deep := range []bool{false, true} {
		for _, listForm := range []bool{false, true} {
			for _, test := range objectcases.MergeCases() {
				t.Run(fmt.Sprintf("merge/%s/deep=%v/list=%v", test.Name, deep, listForm), func(t *testing.T) {
					target := &hostMap{Map: test.Sources[0]}
					args := []runtime.Value{target}
					var sources []runtime.Value
					var before []runtime.Value
					for _, source := range test.Sources[1:] {
						sources = append(sources, &hostMap{Map: source})
						copied, err := source.Clone(ctx)
						if err != nil {
							t.Fatal(err)
						}

						before = append(before, copied)
					}

					if listForm {
						args = append(args, runtime.NewArrayWith(sources...))
					} else {
						args = append(args, sources...)
					}

					merge, want := objects.MergeMutable, test.Shallow
					if deep {
						merge, want = objects.MergeDeepMutable, test.Deep
					}

					got, err := merge(ctx, args...)
					if err != nil || got != target {
						t.Fatalf("identity=%v err=%v", got == target, err)
					}

					assertObjectValue(t, target, want)
					for index, source := range test.Sources[1:] {
						assertObjectValue(t, source, before[index])
					}
				})
			}
		}
	}

	for _, keep := range []bool{false, true} {
		for _, listForm := range []bool{false, true} {
			for _, test := range objectcases.FilterCases() {
				if len(test.Keys) == 0 && !listForm {
					continue
				}

				t.Run(fmt.Sprintf("filter/%s/keep=%v/list=%v", test.Name, keep, listForm), func(t *testing.T) {
					target := &hostMap{Map: test.Source}
					args := []runtime.Value{target}
					var keys []runtime.Value
					for _, key := range test.Keys {
						keys = append(keys, key)
					}

					if listForm {
						args = append(args, runtime.NewArrayWith(keys...))
					} else {
						args = append(args, keys...)
					}

					apply, want := objects.OmitKeysMutable, test.Omitted
					if keep {
						apply, want = objects.KeepKeysMutable, test.Kept
					}

					got, err := apply(ctx, args...)
					if err != nil || got != target {
						t.Fatalf("identity=%v err=%v", got == target, err)
					}

					assertObjectValue(t, target, want)
				})
			}
		}
	}
}

func TestMutableMergesNoSourceDoesNotCopyOrWrite(t *testing.T) {
	ctx := context.Background()
	for _, merge := range []runtime.Function{objects.MergeMutable, objects.MergeDeepMutable} {
		for _, form := range []string{"target", "empty list", "empty map"} {
			writes := 0
			unexpected := errors.New("unexpected copy or mutation")
			nested := runtime.NewObjectWith(map[string]runtime.Value{"a": runtime.Int(1)})
			target := &hostMap{Map: runtime.NewObjectWith(map[string]runtime.Value{"nested": nested}),
				cloneErr: unexpected, emptyErr: unexpected, setErr: unexpected, setCalls: &writes}
			args := []runtime.Value{target}
			switch form {
			case "empty list":
				args = append(args, runtime.NewArray(0))
			case "empty map":
				args = append(args, runtime.NewObject())
			}

			got, err := merge(ctx, args...)
			if err != nil || got != target || writes != 0 {
				t.Fatalf("%s: identity=%v writes=%d err=%v", form, got == target, writes, err)
			}

			retained, err := target.Get(ctx, runtime.String("nested"))
			if err != nil || retained != nested {
				t.Fatalf("nested identity changed: %v", err)
			}
		}
	}
}

func TestMutableShallowOperationsRetainExistingValues(t *testing.T) {
	ctx := context.Background()
	for name, apply := range map[string]runtime.Function{"merge": objects.MergeMutable, "keep": objects.KeepKeysMutable, "omit": objects.OmitKeysMutable} {
		nested := runtime.NewObject()
		target := &hostMap{Map: runtime.NewObjectWith(map[string]runtime.Value{"nested": nested, "remove": runtime.True}),
			cloneErr: errors.New("target must not be cloned"), emptyErr: errors.New("target must not be replaced")}
		var argument runtime.Value = runtime.String("remove")
		if name == "merge" {
			argument = runtime.NewObjectWith(map[string]runtime.Value{"added": runtime.True})
		} else if name == "keep" {
			argument = runtime.String("nested")
		}

		got, err := apply(ctx, target, argument)
		if err != nil || got != target {
			t.Fatalf("identity=%v err=%v", got == target, err)
		}

		retained, err := target.Get(ctx, runtime.String("nested"))
		if err != nil || retained != nested {
			t.Fatalf("retained value was copied: %v", err)
		}
	}
}

func TestMutableArgumentsValidateBeforeMutation(t *testing.T) {
	ctx := context.Background()
	for name, apply := range map[string]runtime.Function{
		"merge": objects.MergeMutable, "deep": objects.MergeDeepMutable,
		"keep": objects.KeepKeysMutable, "omit": objects.OmitKeysMutable,
	} {
		t.Run(name, func(t *testing.T) {
			merge := name == "merge" || name == "deep"
			for _, listForm := range []bool{false, true} {
				writes, removals := 0, 0
				target := &hostMap{Map: runtime.NewObjectWith(map[string]runtime.Value{"a": runtime.Int(1)}), setCalls: &writes, removeCalls: &removals}
				var valid runtime.Value = runtime.String("a")
				if merge {
					valid = runtime.NewObjectWith(map[string]runtime.Value{"a": runtime.Int(2)})
				}

				args := []runtime.Value{target, valid, runtime.Int(1)}
				position := "position 3"
				if listForm {
					args = []runtime.Value{target, runtime.NewArrayWith(valid, runtime.Int(1))}
					position = "position 2"
				}

				got, err := apply(ctx, args...)
				if got != runtime.None || !errors.Is(err, runtime.ErrInvalidType) || !strings.Contains(err.Error(), position) || writes != 0 || removals != 0 {
					t.Fatalf("result=%v writes=%d removals=%d err=%v", got, writes, removals, err)
				}

				if listForm && !strings.Contains(err.Error(), "item 1") {
					t.Fatalf("missing list index: %v", err)
				}
			}

			for _, target := range []runtime.Value{runtime.Int(1), runtime.NewArray(0), &readableObject{Value: runtime.None}} {
				got, err := apply(ctx, target, runtime.String("a"))
				if got != runtime.None || !errors.Is(err, runtime.ErrInvalidType) || !strings.Contains(err.Error(), "position 1") {
					t.Fatalf("target capability: %v, %v", got, err)
				}
			}

			got, err := apply(ctx)
			if got != runtime.None || !errors.Is(err, runtime.ErrInvalidArgumentNumber) {
				t.Fatalf("empty args: %v, %v", got, err)
			}

			if !merge {
				got, err = apply(ctx, runtime.NewObject())
				if got != runtime.None || !errors.Is(err, runtime.ErrInvalidArgumentNumber) {
					t.Fatalf("missing keys: %v, %v", got, err)
				}
			}

			canceled, cancel := context.WithCancel(ctx)
			cancel()
			got, err = apply(canceled, runtime.NewObject(), runtime.NewArray(0))
			if got != runtime.None || !errors.Is(err, context.Canceled) {
				t.Fatalf("cancellation: %v, %v", got, err)
			}
		})
	}
}

func TestMutableHostFailures(t *testing.T) {
	ctx := context.Background()
	sentinel := errors.New("host failure")
	basic := runtime.NewObjectWith(map[string]runtime.Value{"a": runtime.Int(1)})
	for name, apply := range map[string]runtime.Function{
		"merge": objects.MergeMutable, "deep": objects.MergeDeepMutable,
		"keep": objects.KeepKeysMutable, "omit": objects.OmitKeysMutable,
	} {
		t.Run(name, func(t *testing.T) {
			for _, failure := range []error{sentinel, runtime.ErrNotSupported} {
				target := &hostMap{Map: basic, setErr: failure, removeErr: failure, cloneErr: errors.New("no fallback")}
				var argument runtime.Value = runtime.String("a")
				position := "position 1"
				if name == "merge" || name == "deep" {
					argument = basic
					position = "position 2"
				} else if name == "keep" {
					argument = runtime.NewArray(0)
				}

				got, err := apply(ctx, target, argument)
				if got != runtime.None || !errors.Is(err, failure) || !strings.Contains(err.Error(), position) {
					t.Fatalf("result=%v err=%v", got, err)
				}
			}
		})
	}

	for _, merge := range []runtime.Function{objects.MergeMutable, objects.MergeDeepMutable} {
		for _, listForm := range []bool{false, true} {
			bad := &hostMap{Map: basic, walkErr: sentinel}
			var source runtime.Value = bad
			if listForm {
				source = runtime.NewArrayWith(bad)
			}

			got, err := merge(ctx, runtime.NewObject(), source)
			if got != runtime.None || !errors.Is(err, sentinel) || !strings.Contains(err.Error(), "position 2") {
				t.Fatalf("source error: %v, %v", got, err)
			}
		}
	}
}

func TestMutableMergePreservesTraversalAndCleanupErrors(t *testing.T) {
	ctx := context.Background()
	for _, merge := range []runtime.Function{objects.MergeMutable, objects.MergeDeepMutable} {
		for _, failWrite := range []bool{false, true} {
			walkErr, closeErr, writeErr := errors.New("walk"), errors.New("close"), errors.New("write")
			stream := &entryStream{Value: runtime.None, entries: []runtime.Value{runtime.Int(1)}, nextErr: walkErr, closeErr: closeErr}
			source := &hostMap{Map: runtime.NewObject(), walk: func(ctx context.Context, fn runtime.KeyReadablePredicate) error {
				return runtime.ForEach(ctx, stream, func(ctx context.Context, value, _ runtime.Value) (runtime.Boolean, error) {
					return fn(ctx, value, runtime.String("a"))
				})
			}}
			target := &hostMap{Map: runtime.NewObject()}
			wantErr := walkErr
			if failWrite {
				target.setErr, wantErr = writeErr, writeErr
			}

			got, err := merge(ctx, target, source)
			if got != runtime.None || !errors.Is(err, wantErr) || !errors.Is(err, closeErr) {
				t.Fatalf("combined errors: %v, %v", got, err)
			}

			if stream.iteratorCloses != 1 || stream.sourceCloses != 0 {
				t.Fatalf("iterator closes=%d source closes=%d", stream.iteratorCloses, stream.sourceCloses)
			}

			if !failWrite {
				assertObjectValue(t, target, runtime.NewObjectWith(map[string]runtime.Value{"a": runtime.Int(1)}))
			}
		}
	}
}

func TestMutableFiltersSnapshotAndPartiallyApply(t *testing.T) {
	ctx := context.Background()
	sentinel := errors.New("second removal")
	for _, keep := range []bool{false, true} {
		reading := false
		backing := runtime.NewObjectWith(map[string]runtime.Value{"a": runtime.Int(1), "b": runtime.Int(2), "c": runtime.Int(3)})
		base := &hostMap{Map: backing, readingKeys: &reading, keys: &hostList{List: runtime.NewArrayWith(runtime.String("a"), runtime.String("b"), runtime.String("c")), reading: &reading}}
		target := &fallibleRemovalMap{hostMap: base, failure: sentinel}
		apply := objects.OmitKeysMutable
		keys := runtime.NewArrayWith(runtime.String("a"), runtime.String("b"))
		if keep {
			apply = objects.KeepKeysMutable
			keys = runtime.NewArrayWith(runtime.String("c"))
		}

		got, err := apply(ctx, target, keys)
		if got != runtime.None || !errors.Is(err, sentinel) {
			t.Fatalf("partial removal: %v, %v", got, err)
		}

		assertObjectValue(t, backing, runtime.NewObjectWith(map[string]runtime.Value{"b": runtime.Int(2), "c": runtime.Int(3)}))
	}
}

func TestMutableAccessAndNormalizationErrors(t *testing.T) {
	ctx := context.Background()
	sentinel := errors.New("host access")
	for _, apply := range []runtime.Function{objects.KeepKeysMutable, objects.OmitKeysMutable} {
		for _, target := range []runtime.Map{
			&hostMap{Map: runtime.NewObject(), keysErr: sentinel},
			&hostMap{Map: runtime.NewObject(), keys: &hostList{List: runtime.NewArray(0), walkErr: sentinel}},
		} {
			got, err := apply(ctx, target, runtime.NewArray(0))
			if got != runtime.None || !errors.Is(err, sentinel) || !strings.Contains(err.Error(), "position 1") {
				t.Fatalf("key access: %v, %v", got, err)
			}
		}
	}

	for _, apply := range []runtime.Function{objects.MergeMutable, objects.MergeDeepMutable, objects.KeepKeysMutable, objects.OmitKeysMutable} {
		writes, removals := 0, 0
		target := &hostMap{Map: runtime.NewObject(), setCalls: &writes, removeCalls: &removals}
		got, err := apply(ctx, target, &hostList{List: runtime.NewArray(0), walkErr: sentinel})
		if got != runtime.None || !errors.Is(err, sentinel) || !strings.Contains(err.Error(), "position 2") || writes != 0 || removals != 0 {
			t.Fatalf("list traversal: result=%v writes=%d removals=%d err=%v", got, writes, removals, err)
		}
	}

	target := &hostMap{Map: runtime.NewObject(), lookupErr: sentinel}
	got, err := objects.MergeDeepMutable(ctx, target, runtime.NewObjectWith(map[string]runtime.Value{"branch": runtime.NewObject()}))
	if got != runtime.None || !errors.Is(err, sentinel) {
		t.Fatalf("target lookup: %v, %v", got, err)
	}
}

func TestMutableMergesPartialWriteAndCancellation(t *testing.T) {
	for _, merge := range []runtime.Function{objects.MergeMutable, objects.MergeDeepMutable} {
		for _, canceled := range []bool{false, true} {
			ctx, cancel := context.WithCancel(context.Background())
			sentinel := errors.New("second write")
			target := &fallibleWriteMap{hostMap: &hostMap{Map: runtime.NewObject()}, failure: sentinel}
			source := &hostMap{Map: runtime.NewObject(), walk: func(ctx context.Context, fn runtime.KeyReadablePredicate) error {
				if _, err := fn(ctx, runtime.Int(1), runtime.String("first")); err != nil {
					return err
				}

				if canceled {
					cancel()
				}

				_, err := fn(ctx, runtime.Int(2), runtime.String("second"))

				return err
			}}
			got, err := merge(ctx, target, source)
			cancel()
			wantErr := sentinel
			if canceled {
				wantErr = context.Canceled
			}

			if got != runtime.None || !errors.Is(err, wantErr) {
				t.Fatalf("partial write: %v, %v", got, err)
			}

			assertObjectValue(t, target, runtime.NewObjectWith(map[string]runtime.Value{"first": runtime.Int(1)}))
		}
	}
}
