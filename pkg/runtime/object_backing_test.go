package runtime

import (
	"context"
	"errors"
	"hash/maphash"
	"io"
	goruntime "runtime"
	"slices"
	"strconv"
	"testing"
	"weak"
)

func TestObjectBackingAliasesAndClear(t *testing.T) {
	for _, hint := range []int{0, 16, 256} {
		for _, clearOwner := range []bool{false, true} {
			t.Run(strconv.Itoa(hint)+"/clear_owner_"+strconv.FormatBool(clearOwner), func(t *testing.T) {
				obj := NewObjectOf(hint)
				alias := *obj
				obj.SetString("", True)
				iter, _ := alias.Iterate(t.Context())
				for i := range 20 {
					alias.SetString(strconv.Itoa(i), Int(i))
				}
				alias.SetString("", False)
				cleared, survivor := &alias, obj
				if clearOwner {
					cleared, survivor = obj, &alias
				}

				for range 3 {
					if err := cleared.Clear(t.Context()); err != nil {
						t.Fatal(err)
					}
					cleared.SetString("new", True)
					if !cleared.DebugFieldsPageable() {
						t.Fatal("fresh backing lost its inspection budget")
					}
					if v, ok, _ := survivor.Lookup(t.Context(), String("19")); !ok || v != Int(19) {
						t.Fatal("old properties disappeared", v)
					}
					if _, ok, _ := survivor.Lookup(t.Context(), String("new")); ok {
						t.Fatal("new property appeared in detached backing")
					}
				}
				v, k, err := iter.Next(t.Context())
				if v != False || k != String("") || err != nil {
					t.Fatalf("iterator lost its backing or empty key: %v, %v, %v", v, k, err)
				}
				if _, _, err := iter.Next(t.Context()); !errors.Is(err, io.EOF) {
					t.Fatal("iterator included a later insertion", err)
				}
			})
		}
	}
}

func TestObjectBackingReleasesRemovedValues(t *testing.T) {
	for _, count := range []int{1, 2, 9, 128, 256} {
		t.Run(strconv.Itoa(count), func(t *testing.T) {
			create := func() (*Object, weak.Pointer[Array]) {
				obj := NewObject()
				alias := *obj
				value := NewArrayWith(True)
				ref := weak.Make(value)
				obj.SetString("", value)
				for i := 1; i < count; i++ {
					alias.SetString(strconv.Itoa(i), Int(i))
				}
				if err := alias.RemoveKey(t.Context(), String("")); err != nil {
					t.Fatal(err)
				}

				return obj, ref
			}
			obj, ref := create()
			goruntime.GC()
			if ref.Value() != nil {
				t.Fatal("removed value retained by live backing")
			}
			goruntime.KeepAlive(obj)
		})
	}
}

func TestObjectBackingNilEmptyKey(t *testing.T) {
	obj := NewObjectWith(map[string]Value{"": nil})
	alias := *obj
	alias.SetString("second", True)
	for _, v := range []*Object{obj, obj.Copy().(*Object)} {
		got, present, err := v.Lookup(t.Context(), String(""))
		if got != nil || !present || err != nil {
			t.Fatalf("raw nil property: %v %t %v", got, present, err)
		}
	}
}

func TestObjectBackingSparseInspection(t *testing.T) {
	const size = 1_000_000
	for _, reserved := range []bool{true, false} {
		t.Run("reserved_"+strconv.FormatBool(reserved), func(t *testing.T) {
			hint := 0
			if reserved {
				hint = size
			}
			obj := NewObjectOf(hint)
			alias := *obj
			alias.SetString("keep", True)
			if !reserved {
				for i := 1; i < size; i++ {
					alias.SetString(strconv.Itoa(i), False)
				}
				for i := 1; i < size; i++ {
					_ = obj.RemoveKey(t.Context(), String(strconv.Itoa(i)))
				}
			}
			if n, _ := obj.Length(t.Context()); n != 1 || alias.DebugFieldsPageable() {
				t.Fatal("sparse backing reported pageable or wrong size", n)
			}
			err := obj.ReadDebugFields(0, 1, func(string, Value) error {
				t.Fatal("visited sparse backing")

				return nil
			})
			if !errors.Is(err, errors.ErrUnsupported) {
				t.Fatal(err)
			}
			cp := obj.Copy().(*Object)
			cloned, err := alias.Clone(t.Context())
			if err != nil {
				t.Fatal(err)
			}
			for _, v := range []*Object{cp, cloned.(*Object)} {
				v.SetString("a", False)
				var names []string
				err := v.ReadDebugFields(0, 2, func(key string, _ Value) error {
					names = append(names, key)

					return nil
				})
				if err != nil || !slices.Equal(names, []string{"a", "keep"}) {
					t.Fatalf("independent backing: keys %q, error %v", names, err)
				}
			}
		})
	}
}

func TestObjectFingerprintCollisions(t *testing.T) {
	// Find three distinct names with the same tag for this process's hash seed.
	var buckets [256][]string
	var keys []string
	for i := 0; len(keys) < 3; i++ {
		key := strconv.Itoa(i)
		tag := byte(maphash.String(objectKeySeed, key))
		buckets[tag] = append(buckets[tag], key)
		if len(buckets[tag]) == 3 {
			keys = buckets[tag]
		}
	}

	obj := NewObjectOf(16)
	alias := *obj
	for i, k := range keys {
		obj.SetString(k, Int(i))
	}
	alias.SetString(keys[0], True)
	if err := alias.RemoveKey(t.Context(), String(keys[1])); err != nil {
		t.Fatal(err)
	}
	if err := obj.RemoveKey(t.Context(), String(keys[2])); err != nil {
		t.Fatal(err)
	}
	if v, ok, _ := alias.Lookup(t.Context(), String(keys[0])); !ok || v != True {
		t.Fatalf("lost colliding property: %v", v)
	}

	alias.SetString(keys[1], False)
	if v, ok, _ := obj.Lookup(t.Context(), String(keys[2])); ok || v != None {
		t.Fatalf("removed collision returned %v", v)
	}
	names := []string{keys[0], keys[1]}
	slices.Sort(names)
	var got []string
	err := alias.ReadDebugFields(0, 2, func(key string, _ Value) error {
		got = append(got, key)

		return nil
	})
	if err != nil || !slices.Equal(got, names) {
		t.Fatalf("collision page: %q, %v", got, err)
	}
}

func TestObjectForEachKeepsBackingAcrossMutation(t *testing.T) {
	for _, hint := range []int{0, 16, 256} {
		t.Run(strconv.Itoa(hint), func(t *testing.T) {
			obj := NewObjectOf(hint)
			alias := *obj
			names := []string{"a", "b", "c"}
			for _, k := range names {
				obj.SetString(k, True)
			}
			visited := map[string]bool{}
			err := obj.ForEach(t.Context(), func(_ context.Context, v, k Value) (Boolean, error) {
				key := string(k.(String))
				if !slices.Contains(names, key) {
					return True, nil
				}
				if visited[key] {
					t.Fatalf("repeated original field %q", key)
				}
				if len(visited) == 0 {
					_ = obj.Clear(t.Context())
					obj.SetString("replacement", None)
					for _, other := range names {
						if other != key {
							alias.SetString(other, False)
						}
					}
					for i := 0; i < 20; i++ {
						alias.SetString(strconv.Itoa(i), True)
					}
				} else if v != False {
					t.Fatalf("stale value for %q: %v", key, v)
				}
				visited[key] = true
				return True, nil
			})
			if err != nil || len(visited) != len(names) {
				t.Fatalf("original fields visited: %v, %v", visited, err)
			}
		})
	}
}
