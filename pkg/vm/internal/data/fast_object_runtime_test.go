package data_test

import (
	"context"
	"errors"
	"maps"
	goruntime "runtime"
	"strconv"
	"strings"
	"testing"
	"weak"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
	"github.com/MontFerret/ferret/v2/pkg/vm/internal/data"
)

func TestFastObjectGetMissingReturnsNoneWithoutError(t *testing.T) {
	obj := data.NewShapeCache(0, 0).NewFastObject()

	val, err := obj.Get(context.Background(), runtime.NewString("missing"))
	if err != nil {
		t.Fatalf("expected missing fast-object key to return no error, got %v", err)
	}

	if val != runtime.None {
		t.Fatalf("expected runtime.None for missing fast-object key, got %v", val)
	}
}

type promotingKey struct {
	runtime.Value
	object *data.FastObject
	calls  int
}

func (k *promotingKey) String() string {
	k.calls++
	k.object.SetStringCached("new", runtime.Int(23))
	return "new"
}

func TestFastObjectHostKeyPromotesDuringConversion(t *testing.T) {
	for _, method := range []string{"Get", "Lookup", "ContainsKey", "RemoveKey"} {
		t.Run(method, func(t *testing.T) {
			ctx := t.Context()
			fast := data.NewShapeCache(0, 1).NewFastObject()
			fast.SetStringCached("first", runtime.True)
			var object runtime.Map = fast
			key := &promotingKey{Value: runtime.True, object: fast}
			switch method {
			case "Get":
				if v, err := object.Get(ctx, key); v != runtime.Int(23) || err != nil {
					t.Fatalf("Get returned %v, %v", v, err)
				}
			case "Lookup":
				if v, ok, err := object.Lookup(ctx, key); v != runtime.Int(23) || !ok || err != nil {
					t.Fatalf("Lookup returned %v, %t, %v", v, ok, err)
				}
			case "ContainsKey":
				if ok, err := object.ContainsKey(ctx, key); !ok || err != nil {
					t.Fatalf("ContainsKey returned %t, %v", ok, err)
				}
			case "RemoveKey":
				if err := object.RemoveKey(ctx, key); err != nil {
					t.Fatal(err)
				}
				if _, ok, err := object.Lookup(ctx, runtime.String("new")); ok || err != nil {
					t.Fatalf("removed key present: %t, %v", ok, err)
				}
			}
			if key.calls != 1 {
				t.Fatalf("called host String %d times", key.calls)
			}
		})
	}
}

func TestFastObjectDictionaryBackingLifetime(t *testing.T) {
	ctx := t.Context()
	object := data.NewShapeCache(0, 1).NewFastObject()
	object.SetStringCached("keep", runtime.Int(1))
	object.SetStringCached("seed", runtime.True)
	if err := object.RemoveKey(ctx, runtime.String("seed")); err != nil {
		t.Fatal(err)
	}

	alias := *object
	longKey := runtime.String(strings.Repeat("x", 1<<17))
	alias.SetStringCached(string(longKey), runtime.True)
	if err := object.RemoveKey(ctx, longKey); err != nil {
		t.Fatal(err)
	}

	if alias.DebugFieldsPageable() || object.DebugFieldsPageable() {
		t.Fatal("original dictionary lost its key history")
	}

	copied := object.Copy().(*data.FastObject)
	cloned, err := object.Clone(ctx)
	if err != nil {
		t.Fatal(err)
	}

	iter, err := object.Iterate(ctx)
	if err != nil {
		t.Fatal(err)
	}

	object.SetStringCached("keep", runtime.Int(2))
	if err := object.Clear(ctx); err != nil {
		t.Fatal(err)
	}

	object.SetStringCached("keep", runtime.Int(3))
	value, key, err := iter.Next(ctx)
	if err != nil || key != runtime.String("keep") || value != runtime.Int(2) {
		t.Fatalf("iterator after clear: %v=%v, %v", key, value, err)
	}

	value, found, err := alias.Lookup(ctx, runtime.String("keep"))
	if err != nil || !found || value != runtime.Int(2) || alias.DebugFieldsPageable() {
		t.Fatalf("old backing after clear: %v, found %t, pageable %t, %v", value, found, alias.DebugFieldsPageable(), err)
	}

	for _, snapshot := range []*data.FastObject{copied, cloned.(*data.FastObject)} {
		value, found, err := snapshot.Lookup(ctx, runtime.String("keep"))
		if err != nil || !found || value != runtime.Int(1) || !snapshot.DebugFieldsPageable() {
			t.Fatalf("independent dictionary: %v, found %t, pageable %t, %v", value, found, snapshot.DebugFieldsPageable(), err)
		}
	}
}

func TestFastObjectShapeTraversalRetainsSlots(t *testing.T) {
	for _, mutate := range []string{"clear", "promote"} {
		for _, traverse := range []string{"callback", "page"} {
			t.Run(mutate+"/"+traverse, func(t *testing.T) {
				ctx := t.Context()
				obj := data.NewShapeCache(0, 2).NewFastObject()
				obj.SetStringCached("a", runtime.Int(1))
				obj.SetStringCached("b", runtime.Int(2))
				seen := map[string]runtime.Value{}
				visit := func(key string, value runtime.Value) error {
					seen[key] = value
					if key != "a" {
						return nil
					}

					switch mutate {
					case "clear":
						if err := obj.Clear(ctx); err != nil {
							return err
						}
						obj.SetStringCached("fresh", runtime.False)
					case "promote":
						obj.SetStringCached("c", runtime.Int(3))
					}

					return nil
				}

				var err error
				if traverse == "callback" {
					err = obj.ForEach(ctx, func(_ context.Context, value, key runtime.Value) (runtime.Boolean, error) {
						return true, visit(string(key.(runtime.String)), value)
					})
				} else {
					err = obj.ReadDebugFields(0, 2, visit)
				}

				want := map[string]runtime.Value{"a": runtime.Int(1), "b": runtime.Int(2)}
				if err != nil || !maps.Equal(seen, want) {
					t.Fatalf("visited %v: %v, want %v", seen, err, want)
				}
			})
		}
	}
}

func TestFastObjectCallbackUsesCurrentMatchingSlots(t *testing.T) {
	for _, tc := range []struct {
		wantB  runtime.Value
		name   string
		names  []string
		clear  bool
		remove bool
	}{
		{name: "grow_then_remove", names: []string{"c"}, remove: true},
		{name: "grow_then_overwrite", names: []string{"c", "b"}, wantB: runtime.Int(20)},
		{name: "rebuild_same", clear: true, names: []string{"a", "b"}, wantB: runtime.Int(20)},
		{name: "rebuild_wider", clear: true, names: []string{"a", "b", "c"}, wantB: runtime.Int(20)},
		{name: "rebuild_other", clear: true, names: []string{"x", "y"}, wantB: runtime.Int(2)},
	} {
		t.Run(tc.name, func(t *testing.T) {
			object := data.NewShapeCache(16, 0).NewFastObject()
			object.SetStringCached("a", runtime.Int(1))
			object.SetStringCached("b", runtime.Int(2))
			seen := map[string]runtime.Value{}
			err := object.ForEach(t.Context(), func(ctx context.Context, value, key runtime.Value) (runtime.Boolean, error) {
				seen[key.String()] = value
				if key != runtime.String("a") {
					return true, nil
				}
				if tc.clear {
					if err := object.Clear(ctx); err != nil {
						return false, err
					}
				}
				for _, name := range tc.names {
					object.SetStringCached(name, runtime.Int(20))
				}
				if tc.remove {
					return true, object.RemoveKey(ctx, runtime.String("b"))
				}

				return true, nil
			})
			want := map[string]runtime.Value{"a": runtime.Int(1)}
			if tc.wantB != nil {
				want["b"] = tc.wantB
			}
			if err != nil || !maps.Equal(seen, want) {
				t.Fatalf("visited %v, error %v, want %v", seen, err, want)
			}
		})
	}
}

func TestClearedFastDictionaryReleasesReplacement(t *testing.T) {
	for _, threshold := range []int{1, 64, 128} {
		for _, clearOwner := range []bool{false, true} {
			create := func() (*data.FastObject, runtime.Iterator, weak.Pointer[runtime.Array], weak.Pointer[runtime.Array]) {
				owner := data.NewShapeCache(0, threshold).NewFastObject()
				old := runtime.NewArrayWith(runtime.True)
				owner.SetStringCached("old", old)
				for i := range threshold {
					owner.SetStringCached(strconv.Itoa(i), runtime.True)
				}

				if err := owner.RemoveKey(t.Context(), runtime.String("old")); err != nil {
					t.Fatal(err)
				}

				alias := *owner
				cleared, survivor := &alias, owner
				if clearOwner {
					cleared, survivor = survivor, cleared
				}
				iter, err := cleared.Iterate(t.Context())
				if err != nil {
					t.Fatal(err)
				}

				if err := cleared.Clear(t.Context()); err != nil {
					t.Fatal(err)
				}

				child := runtime.NewArrayWith(runtime.False)
				cleared.SetStringCached("new", child)

				return survivor, iter, weak.Make(child), weak.Make(old)
			}

			survivor, iter, child, removed := create()
			goruntime.GC()
			if removed.Value() != nil {
				t.Fatalf("threshold %d: live dictionary retained a removed shape value", threshold)
			}
			if child.Value() != nil {
				t.Fatalf("threshold %d, cleared original %t: old dictionary retained replacement", threshold, clearOwner)
			}
			goruntime.KeepAlive(survivor)
			goruntime.KeepAlive(iter)
		}
	}
}

type mutatingClone struct {
	runtime.Value
	mutate func()
}

func (v *mutatingClone) Clone(context.Context) (runtime.Cloneable, error) {
	v.mutate()

	return runtime.NewArrayWith(runtime.True), nil
}

func TestFastObjectCloneRetainsSourceShape(t *testing.T) {
	for _, change := range []string{"clear", "promote", "remove"} {
		t.Run(change, func(t *testing.T) {
			ctx := t.Context()
			obj := data.NewShapeCache(0, 2).NewFastObject()
			child := &mutatingClone{Value: runtime.True, mutate: func() {
				switch change {
				case "clear":
					if err := obj.Clear(ctx); err != nil {
						t.Fatal(err)
					}
				case "promote":
					obj.SetStringCached("c", runtime.Int(3))
				case "remove":
					if err := obj.RemoveKey(ctx, runtime.String("b")); err != nil {
						t.Fatal(err)
					}
				}
			}}
			obj.SetStringCached("a", child)
			obj.SetStringCached("b", runtime.Int(2))
			shape := obj.Shape()

			value, err := obj.Clone(ctx)
			if err != nil {
				t.Fatal(err)
			}
			cloned := value.(*data.FastObject)
			count, _ := cloned.Length(ctx)
			want := runtime.Int(2)
			if change == "remove" {
				want--
			}
			if cloned.Shape() != shape || count != want {
				t.Fatalf("cloned shape or live count changed: %p/%d, want %p/%d", cloned.Shape(), count, shape, want)
			}
			a, found, err := cloned.Lookup(ctx, runtime.String("a"))
			if _, ok := a.(*runtime.Array); !ok || !found || err != nil {
				t.Fatalf("missing cloned child: %v/%t, %v", a, found, err)
			}
			b, found, err := cloned.Lookup(ctx, runtime.String("b"))
			if err != nil || found != (change != "remove") || (found && b != runtime.Int(2)) {
				t.Fatalf("unexpected second field: %v/%t, %v", b, found, err)
			}
		})
	}
}

type (
	nilCopyChild  struct{ runtime.Value }
	nilCloneChild struct {
		nilCopyChild
		err error
	}
)

func (nilCopyChild) Copy() runtime.Value                                 { return nil }
func (v nilCloneChild) Clone(context.Context) (runtime.Cloneable, error) { return nil, v.err }

func TestObjectClonePreservesNilChildResults(t *testing.T) {
	ctx := t.Context()
	expected := runtime.NewObjectWith(map[string]runtime.Value{"clone": runtime.None, "copy": runtime.None})
	for _, tc := range []struct {
		object runtime.Map
		want   runtime.Value
		name   string
	}{
		{data.NewShapeCache(0, 0).NewFastObject(), runtime.None, "shape"},
		{data.NewShapeCache(0, 1).NewFastObject(), nil, "dictionary"},
		{runtime.NewObject(), nil, "native"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			if err := tc.object.Set(ctx, runtime.String("copy"), nilCopyChild{runtime.True}); err != nil {
				t.Fatal(err)
			}
			if err := tc.object.Set(ctx, runtime.String("clone"), nilCloneChild{nilCopyChild: nilCopyChild{runtime.True}}); err != nil {
				t.Fatal(err)
			}
			first, err := tc.object.(runtime.Cloneable).Clone(ctx)
			if err != nil {
				t.Fatal(err)
			}
			check := func(snapshot runtime.Map) {
				t.Helper()
				if size, err := snapshot.Length(ctx); size != 2 || err != nil {
					t.Fatalf("length %d, %v", size, err)
				}
				keys, err := snapshot.Keys(ctx)
				if err != nil {
					t.Fatal(err)
				}
				if n, err := keys.Length(ctx); n != 2 || err != nil {
					t.Fatalf("key count %d, %v", n, err)
				}
				for _, key := range []runtime.String{"copy", "clone"} {
					if has, err := keys.Contains(ctx, key); !has || err != nil {
						t.Fatalf("keys missing %s: %v", key, err)
					}
					if v, ok, err := snapshot.Lookup(ctx, key); !ok || v != tc.want || err != nil {
						t.Fatalf("%s: got %v, found %t, error %v", key, v, ok, err)
					}
				}
				if eq, err := runtime.EqualValues(ctx, snapshot, expected); !eq || err != nil || snapshot.Hash() != expected.Hash() {
					t.Fatalf("snapshot differs from an object with None children: %v", err)
				}
			}
			check(first.(runtime.Map))
			check(first.Copy().(runtime.Map))
			second, err := first.Clone(ctx)
			if err != nil {
				t.Fatal(err)
			}
			check(second.(runtime.Map))

			failure := errors.New("child clone failed")
			child := nilCloneChild{nilCopyChild: nilCopyChild{runtime.True}, err: failure}
			if err := tc.object.Set(ctx, runtime.String("clone"), child); err != nil {
				t.Fatal(err)
			}
			if partial, err := tc.object.(runtime.Cloneable).Clone(ctx); partial != nil || !errors.Is(err, failure) {
				t.Fatalf("failed clone returned %v, %v", partial, err)
			}
		})
	}
}
