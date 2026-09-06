package objects_test

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
	"github.com/MontFerret/ferret/v2/pkg/stdlib/objects"
	objectcases "github.com/MontFerret/ferret/v2/test/spec/objects"
)

func TestImmutableTransformationFixtures(t *testing.T) {
	ctx := context.Background()
	for _, deep := range []bool{false, true} {
		for _, listForm := range []bool{false, true} {
			for _, test := range objectcases.MergeCases() {
				args := make([]runtime.Value, len(test.Sources))
				for i, source := range test.Sources {
					args[i] = &hostMap{Map: source}
				}

				if listForm {
					args = []runtime.Value{runtime.NewArrayWith(args...)}
				}

				merge, want := objects.Merge, test.Shallow
				if deep {
					merge, want = objects.MergeDeep, test.Deep
				}

				got, err := merge(ctx, args...)
				if err != nil {
					t.Fatal(err)
				}

				assertObjectValue(t, got, want)
			}
		}
	}

	for _, keep := range []bool{false, true} {
		for _, test := range objectcases.FilterCases() {
			keys := runtime.NewArray(0)
			for _, key := range test.Keys {
				if err := keys.Append(ctx, key); err != nil {
					t.Fatal(err)
				}
			}

			apply, want := objects.OmitKeys, test.Omitted
			if keep {
				apply, want = objects.KeepKeys, test.Kept
			}

			before, err := test.Source.Clone(ctx)
			if err != nil {
				t.Fatal(err)
			}

			got, err := apply(ctx, &hostMap{Map: test.Source}, keys)
			if err != nil {
				t.Fatal(err)
			}

			assertObjectValue(t, got, want)
			assertObjectValue(t, test.Source, before)
		}
	}
}

func TestImmutableNestedValuesAndCopyContract(t *testing.T) {
	ctx := context.Background()
	for _, name := range []string{"values", "entries", "keep", "omit", "merge", "deep", "zip", "from_entries"} {
		t.Run(name, func(t *testing.T) {
			copies, clones := 0, 0
			copyValue := &copyOnlyValue{Value: runtime.String("copy"), copies: &copies, number: 1}
			cloneLeaf := &cloneValue{Value: runtime.String("clone"), clones: &clones, number: 1}
			nested := runtime.NewObjectWith(map[string]runtime.Value{"copy": copyValue, "clone": cloneLeaf, "array": runtime.NewArrayWith(runtime.Int(1))})
			src := runtime.NewObjectWith(map[string]runtime.Value{"key": nested})
			var result runtime.Value
			var err error
			switch name {
			case "values":
				result, err = objects.Values(ctx, src)
			case "entries":
				result, err = objects.Entries(ctx, src)
			case "keep":
				result, err = objects.KeepKeys(ctx, src, runtime.String("key"))
			case "omit":
				result, err = objects.OmitKeys(ctx, src, runtime.String("missing"))
			case "merge":
				result, err = objects.Merge(ctx, src)
			case "deep":
				result, err = objects.MergeDeep(ctx, src, runtime.NewObjectWith(map[string]runtime.Value{"key": runtime.NewObjectWith(map[string]runtime.Value{"new": runtime.Int(2)})}))
			case "zip":
				result, err = objects.Zip(ctx, runtime.NewArrayWith(runtime.String("key")), runtime.NewArrayWith(nested))
			case "from_entries":
				result, err = objects.FromEntries(ctx, runtime.NewArrayWith(runtime.NewArrayWith(runtime.String("key"), nested)))
			}

			if err != nil {
				t.Fatal(err)
			}

			if copies != 1 || clones != 1 {
				t.Fatalf("copies=%d clones=%d, want one each", copies, clones)
			}

			var copied runtime.Value
			if name == "values" {
				copied, err = result.(runtime.List).At(ctx, 0)
			} else if name == "entries" {
				var pair runtime.Value
				pair, err = result.(runtime.List).At(ctx, 0)
				if err == nil {
					copied, err = pair.(runtime.List).At(ctx, 1)
				}
			} else {
				copied, err = result.(runtime.Map).Get(ctx, runtime.String("key"))
			}

			if err != nil {
				t.Fatal(err)
			}

			copyValue.number = 2
			copiedMap := copied.(runtime.Map)
			copiedLeaf, err := copiedMap.Get(ctx, runtime.String("copy"))
			if err != nil || copiedLeaf.(*copyOnlyValue).number != 1 {
				t.Fatalf("source alias: %v", err)
			}
			copiedLeaf.(*copyOnlyValue).number = 3
			if copyValue.number != 2 {
				t.Fatal("result aliases source")
			}

			cloneLeaf.number = 2
			copiedClone, err := copiedMap.Get(ctx, runtime.String("clone"))
			if err != nil || copiedClone.(*cloneValue).number != 1 {
				t.Fatalf("clone aliases source: %v", err)
			}

			copiedClone.(*cloneValue).number = 3
			if cloneLeaf.number != 2 {
				t.Fatal("cloned result aliases source")
			}

			if err := nested.Set(ctx, runtime.String("source_change"), runtime.True); err != nil {
				t.Fatal(err)
			}

			found, err := copiedMap.ContainsKey(ctx, runtime.String("source_change"))
			if err != nil || found {
				t.Fatalf("source map aliases result: %v", err)
			}

			if err := copiedMap.Set(ctx, runtime.String("changed"), runtime.True); err != nil {
				t.Fatal(err)
			}
			found, err = nested.ContainsKey(ctx, runtime.String("changed"))
			if err != nil || found {
				t.Fatalf("nested map aliases source: %v", err)
			}

			sourceArray, err := nested.Get(ctx, runtime.String("array"))
			if err != nil {
				t.Fatal(err)
			}

			if err := sourceArray.(runtime.List).Append(ctx, runtime.Int(2)); err != nil {
				t.Fatal(err)
			}
			resultArray, err := copiedMap.Get(ctx, runtime.String("array"))
			if err != nil {
				t.Fatal(err)
			}
			length, err := resultArray.(runtime.List).Length(ctx)
			if err != nil || length != 1 {
				t.Fatalf("nested list aliases source: %v", err)
			}

			if err := resultArray.(runtime.List).Append(ctx, runtime.Int(3)); err != nil {
				t.Fatal(err)
			}

			length, err = sourceArray.(runtime.List).Length(ctx)
			if err != nil || length != 2 {
				t.Fatalf("result list aliases source: %v", err)
			}
		})
	}
}

func TestFromEntriesSinglePassAndTupleCapabilities(t *testing.T) {
	ctx := context.Background()
	stream := &entryStream{Value: runtime.None, entries: []runtime.Value{
		&entryPair{Value: runtime.None, length: 2, items: [2]runtime.Value{runtime.String("a"), runtime.Int(1)}},
		runtime.NewArrayWith(runtime.String("a"), runtime.Int(2)),
	}}
	result, err := objects.FromEntries(ctx, stream)
	if err != nil {
		t.Fatal(err)
	}

	assertObjectValue(t, result, runtime.NewObjectWith(map[string]runtime.Value{"a": runtime.Int(2)}))
	if stream.iterations != 1 || stream.iteratorCloses != 1 || stream.sourceCloses != 0 {
		t.Fatalf("iterator lifecycle = %+v", stream)
	}
}

func TestObjectHostErrors(t *testing.T) {
	ctx := context.Background()
	sentinel := errors.New("host failure")
	other := errors.New("close failure")
	basic := runtime.NewObjectWith(map[string]runtime.Value{"a": runtime.Int(1)})
	keys := runtime.NewArrayWith(runtime.String("a"))
	vals := runtime.NewArrayWith(runtime.Int(1))
	var clones int
	badValue := &cloneValue{Value: runtime.String("bad"), clones: &clones, cloneErr: sentinel}
	cases := []struct {
		name     string
		run      func() (runtime.Value, error)
		argument string
	}{
		{"keys", func() (runtime.Value, error) { return objects.Keys(ctx, &hostMap{Map: basic, keysErr: sentinel}) }, "position 1"},
		{"key list length", func() (runtime.Value, error) {
			return objects.Keys(ctx, &hostMap{Map: basic, keys: &hostList{List: keys, lengthErr: sentinel}})
		}, "position 1"},
		{"values", func() (runtime.Value, error) { return objects.Values(ctx, &hostMap{Map: basic, valuesErr: sentinel}) }, "position 1"},
		{"values walk", func() (runtime.Value, error) {
			return objects.Values(ctx, &hostMap{Map: basic, values: &hostList{List: vals, walkErr: sentinel}})
		}, "position 1"},
		{"entries", func() (runtime.Value, error) { return objects.Entries(ctx, &hostMap{Map: basic, walkErr: sentinel}) }, "position 1"},
		{"has key", func() (runtime.Value, error) {
			return objects.HasKey(ctx, &hostMap{Map: basic, containsErr: sentinel}, runtime.String("a"))
		}, "position 1"},
		{"zip key length", func() (runtime.Value, error) {
			return objects.Zip(ctx, &hostList{List: keys, lengthErr: sentinel}, vals)
		}, "position 1"},
		{"zip value length", func() (runtime.Value, error) {
			return objects.Zip(ctx, keys, &hostList{List: vals, lengthErr: sentinel})
		}, "position 2"},
		{"zip key access", func() (runtime.Value, error) {
			return objects.Zip(ctx, &hostList{List: keys, accessErr: sentinel}, vals)
		}, "position 1"},
		{"zip value access", func() (runtime.Value, error) {
			return objects.Zip(ctx, keys, &hostList{List: vals, accessErr: sentinel})
		}, "position 2"},
		{"zip clone", func() (runtime.Value, error) { return objects.Zip(ctx, keys, runtime.NewArrayWith(badValue)) }, "position 2"},
		{"merge empty", func() (runtime.Value, error) { return objects.Merge(ctx, &hostMap{Map: basic, emptyErr: sentinel}) }, "position 1"},
		{"merge source", func() (runtime.Value, error) {
			return objects.Merge(ctx, basic, &hostMap{Map: basic, walkErr: sentinel})
		}, "position 2"},
		{"deep source", func() (runtime.Value, error) {
			return objects.MergeDeep(ctx, runtime.NewArrayWith(basic, &hostMap{Map: basic, walkErr: sentinel}))
		}, "position 1"},
		{"merge destination", func() (runtime.Value, error) {
			return objects.Merge(ctx, &hostMap{Map: basic, empty: &hostMap{Map: runtime.NewObject(), setErr: sentinel}})
		}, "position 1"},
		{"deep lookup", func() (runtime.Value, error) {
			return objects.MergeDeep(ctx, &hostMap{Map: runtime.NewObjectWith(map[string]runtime.Value{"a": basic}), empty: &hostMap{Map: runtime.NewObject(), lookupErr: sentinel}})
		}, "position 1"},
		{"merge list", func() (runtime.Value, error) {
			return objects.MergeDeep(ctx, &hostList{List: runtime.NewArrayWith(basic), walkErr: sentinel})
		}, "position 1"},
		{"filter clone", func() (runtime.Value, error) {
			return objects.KeepKeys(ctx, &hostMap{Map: basic, cloneErr: sentinel}, runtime.String("a"))
		}, "position 1"},
		{"filter list", func() (runtime.Value, error) {
			return objects.OmitKeys(ctx, basic, &hostList{List: keys, walkErr: sentinel})
		}, "position 2"},
		{"filter remove", func() (runtime.Value, error) {
			return objects.OmitKeys(ctx, &hostMap{Map: basic, cloned: &hostMap{Map: runtime.NewObjectWith(map[string]runtime.Value{"a": runtime.Int(1)}), removeErr: sentinel}}, runtime.String("a"))
		}, "position 1"},
	}

	for _, test := range cases {
		t.Run(test.name, func(t *testing.T) {
			result, err := test.run()
			if result != runtime.None || !errors.Is(err, sentinel) || !strings.Contains(err.Error(), test.argument) {
				t.Fatalf("result=%v error=%v; want None and %s with host cause", result, err, test.argument)
			}
		})
	}

	for _, stage := range []string{"iterate", "next", "length", "key", "value", "clone", "close", "both"} {
		t.Run(stage, func(t *testing.T) {
			pair := &entryPair{Value: runtime.None, length: 2, items: [2]runtime.Value{runtime.String("a"), runtime.Int(1)}}
			stream := &entryStream{Value: runtime.None, entries: []runtime.Value{pair}}
			switch stage {
			case "iterate":
				stream.iterateErr = sentinel
			case "next":
				stream.nextErr = sentinel
			case "length":
				pair.lengthErr = sentinel
			case "key":
				pair.accessErr = sentinel
			case "value":
				pair.accessErr, pair.failAt = sentinel, 1
			case "clone":
				pair.items[1] = badValue
			case "close":
				stream.closeErr = sentinel
			case "both":
				stream.nextErr, stream.closeErr = sentinel, other
			}
			result, err := objects.FromEntries(ctx, stream)
			if result != runtime.None || !errors.Is(err, sentinel) {
				t.Fatalf("result=%v error=%v", result, err)
			}

			if stage == "both" && !errors.Is(err, other) {
				t.Fatalf("lost close error: %v", err)
			}

			wantCloses := 1
			if stage == "iterate" {
				wantCloses = 0
			}

			if stream.iterations != 1 || stream.iteratorCloses != wantCloses || stream.sourceCloses != 0 {
				t.Fatalf("lifecycle: %+v", stream)
			}
		})
	}
}

func TestObjectCancellationAndArgumentPositions(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	for _, run := range []func() (runtime.Value, error){
		func() (runtime.Value, error) { return objects.Keys(ctx, runtime.NewObject()) },
		func() (runtime.Value, error) { return objects.Values(ctx, runtime.NewObject()) },
		func() (runtime.Value, error) { return objects.Entries(ctx, runtime.NewObject()) },
		func() (runtime.Value, error) { return objects.HasKey(ctx, runtime.NewObject(), runtime.String("a")) },
		func() (runtime.Value, error) { return objects.Merge(ctx, runtime.NewArray(0)) },
		func() (runtime.Value, error) { return objects.MergeDeep(ctx, runtime.NewArray(0)) },
		func() (runtime.Value, error) { return objects.KeepKeys(ctx, runtime.NewObject(), runtime.NewArray(0)) },
		func() (runtime.Value, error) { return objects.OmitKeys(ctx, runtime.NewObject(), runtime.NewArray(0)) },
		func() (runtime.Value, error) { return objects.Zip(ctx, runtime.NewArray(0), runtime.NewArray(0)) },
		func() (runtime.Value, error) { return objects.FromEntries(ctx, &entryStream{Value: runtime.None}) },
	} {
		result, err := run()
		if result != runtime.None || !errors.Is(err, context.Canceled) {
			t.Fatalf("cancellation: %v, %v", result, err)
		}
	}

	for _, apply := range []func(context.Context, ...runtime.Value) (runtime.Value, error){objects.KeepKeys, objects.OmitKeys} {
		_, err := apply(context.Background(), runtime.NewObject(), runtime.String("ok"), runtime.Int(1))
		if err == nil || !strings.Contains(err.Error(), "position 3") {
			t.Fatalf("key position: %v", err)
		}
	}

	for _, apply := range []func(context.Context, ...runtime.Value) (runtime.Value, error){objects.Merge, objects.MergeDeep} {
		_, err := apply(context.Background(), runtime.NewObject(), runtime.Int(1))
		if err == nil || !strings.Contains(err.Error(), "position 2") {
			t.Fatalf("map position: %v", err)
		}
	}
}

func assertObjectValue(t *testing.T, got, want runtime.Value) {
	t.Helper()
	equal, err := runtime.EqualValues(context.Background(), got, want)
	if err != nil || !equal {
		t.Fatalf("got %v, want %v: %v", got, want, err)
	}
}

func TestKeyFiltersSnapshotLiveViews(t *testing.T) {
	ctx := context.Background()
	for _, keep := range []bool{false, true} {
		reading := false
		backing := runtime.NewObjectWith(map[string]runtime.Value{"a": runtime.Int(1), "b": runtime.Int(2)})
		view := &hostList{List: runtime.NewArrayWith(runtime.String("a"), runtime.String("b")), reading: &reading}
		dst := &hostMap{Map: backing, keys: view, readingKeys: &reading}
		apply, want := runtime.OmitMapKeys, runtime.NewObjectWith(map[string]runtime.Value{"b": runtime.Int(2)})
		if keep {
			apply, want = runtime.KeepMapKeys, runtime.NewObjectWith(map[string]runtime.Value{"a": runtime.Int(1)})
		}

		if err := apply(ctx, dst, runtime.String("a")); err != nil {
			t.Fatal(err)
		}

		assertObjectValue(t, backing, want)
	}
}

func TestMergesDoNotDependOnEquality(t *testing.T) {
	ctx := context.Background()
	src := &hostMap{Map: runtime.NewObjectWith(map[string]runtime.Value{"a": runtime.Int(1)}), equalityErr: errors.New("equality must not be called")}
	for _, merge := range []func(context.Context, ...runtime.Value) (runtime.Value, error){objects.Merge, objects.MergeDeep} {
		result, err := merge(ctx, src)
		if err != nil {
			t.Fatal(err)
		}

		assertObjectValue(t, result, src.Map)
	}
}

func TestAllImmutableMapCopiesPropagateCloneFailure(t *testing.T) {
	ctx := context.Background()
	sentinel := errors.New("clone failed")
	var clones int
	src := runtime.NewObjectWith(map[string]runtime.Value{"bad": &cloneValue{Value: runtime.None, clones: &clones, cloneErr: sentinel}})
	for _, run := range []func() (runtime.Value, error){
		func() (runtime.Value, error) { return objects.Values(ctx, src) },
		func() (runtime.Value, error) { return objects.Entries(ctx, src) },
		func() (runtime.Value, error) { return objects.KeepKeys(ctx, src, runtime.String("missing")) },
		func() (runtime.Value, error) { return objects.OmitKeys(ctx, src, runtime.String("bad")) },
		func() (runtime.Value, error) { return objects.Merge(ctx, src) },
		func() (runtime.Value, error) { return objects.MergeDeep(ctx, src) },
	} {
		result, err := run()
		if result != runtime.None || !errors.Is(err, sentinel) {
			t.Fatalf("clone failure: %v, %v", result, err)
		}
	}
}

func TestFromEntriesStrictPairShapeAndCombinedFailure(t *testing.T) {
	ctx := context.Background()
	for _, pair := range []runtime.Value{
		runtime.NewObject(), runtime.String("ab"), runtime.None,
		&entryPair{Value: runtime.None, length: 1},
		&entryPair{Value: runtime.None, length: 3},
		&entryPair{Value: runtime.None, length: 2, items: [2]runtime.Value{runtime.Int(1), runtime.True}},
	} {
		stream := &entryStream{Value: runtime.None, entries: []runtime.Value{pair}}
		result, err := objects.FromEntries(ctx, stream)
		if err == nil || result != runtime.None || stream.iteratorCloses != 1 {
			t.Fatalf("invalid pair: %v, %v", result, err)
		}
	}

	cleanup := errors.New("close failed")
	stream := &entryStream{Value: runtime.None, entries: []runtime.Value{runtime.NewArray(0)}, closeErr: cleanup}
	result, err := objects.FromEntries(ctx, stream)
	if result != runtime.None || !errors.Is(err, cleanup) || !errors.Is(err, runtime.ErrInvalidArgument) {
		t.Fatalf("lost predicate or cleanup error: %v, %v", result, err)
	}
}

func TestEntriesUsesOneAssociatedTraversal(t *testing.T) {
	ctx := context.Background()
	calls := 0
	sentinel := errors.New("separate key/value access is forbidden")
	src := &hostMap{
		Map:     runtime.NewObjectWith(map[string]runtime.Value{"a": runtime.Int(1), "b": runtime.Int(2)}),
		keysErr: sentinel, valuesErr: sentinel, walkCalls: &calls,
	}
	entries, err := objects.Entries(ctx, src)
	if err != nil {
		t.Fatal(err)
	}

	if calls != 1 {
		t.Fatalf("map traversed %d times", calls)
	}

	roundTrip, err := objects.FromEntries(ctx, entries)
	if err != nil {
		t.Fatal(err)
	}

	assertObjectValue(t, roundTrip, src.Map)
}
