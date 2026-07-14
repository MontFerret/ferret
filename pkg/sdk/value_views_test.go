package sdk_test

import (
	"context"
	"errors"
	"io"
	"math"
	"strconv"
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
	"github.com/MontFerret/ferret/v2/pkg/sdk"
)

type codecRecord struct {
	ID int
}

var (
	_ runtime.Value          = (*sdk.HostValue[int])(nil)
	_ runtime.Iterable       = (*sdk.IterableValue[*runtime.Array])(nil)
	_ runtime.IndexReadable  = (*sdk.SliceView[int])(nil)
	_ runtime.IndexLookup    = (*sdk.SliceView[int])(nil)
	_ runtime.IndexWritable  = (*sdk.SliceView[int])(nil)
	_ runtime.IndexRemovable = (*sdk.SliceView[int])(nil)
	_ runtime.ValueRemovable = (*sdk.SliceView[int])(nil)
	_ runtime.Iterable       = (*sdk.SliceView[int])(nil)
	_ runtime.Measurable     = (*sdk.SliceView[int])(nil)
	_ runtime.Sortable       = (*sdk.SliceView[int])(nil)
	_ runtime.KeyReadable    = (*sdk.MapView[string, int])(nil)
	_ runtime.KeyLookup      = (*sdk.MapView[string, int])(nil)
	_ runtime.KeyWritable    = (*sdk.MapView[string, int])(nil)
	_ runtime.KeyRemovable   = (*sdk.MapView[string, int])(nil)
	_ runtime.ValueRemovable = (*sdk.MapView[string, int])(nil)
	_ runtime.Iterable       = (*sdk.MapView[string, int])(nil)
	_ runtime.Measurable     = (*sdk.MapView[string, int])(nil)
)

func TestHostValueIdentityAndCapabilities(t *testing.T) {
	typeName := runtime.NewType("sdk_test", "opaque", func(runtime.Value) bool { return true })
	value := sdk.NewHostValueWithType(typeName, 42)
	copy, ok := value.Copy().(*sdk.HostValue[int])
	if !ok {
		t.Fatalf("copy has type %T", value.Copy())
	}

	if value.Hash() == 0 || copy.Hash() != value.Hash() {
		t.Fatalf("identity was not preserved: original=%d copy=%d", value.Hash(), copy.Hash())
	}
	if runtime.TypeName(copy.Type()) != runtime.TypeName(typeName) || copy.Target() != 42 {
		t.Fatalf("copy lost type or target: type=%s target=%d", copy.Type(), copy.Target())
	}
	if sdk.NewHostValue(42).Hash() == sdk.NewHostValue(42).Hash() {
		t.Fatal("distinct host values unexpectedly share identity")
	}

	assertDoesNotImplement[runtime.Iterable](t, value, "runtime.Iterable")
	assertDoesNotImplement[runtime.Queryable](t, value, "runtime.Queryable")
	assertDoesNotImplement[runtime.KeyReadable](t, value, "runtime.KeyReadable")
}

func TestCapabilityAccurateIterableValues(t *testing.T) {
	array := runtime.NewArrayWith(runtime.NewInt(1))
	iterable := sdk.NewIterableValue(array)
	assertDoesNotImplement[runtime.Queryable](t, iterable, "runtime.Queryable")
	assertDoesNotImplement[runtime.KeyReadable](t, iterable, "runtime.KeyReadable")

	iterator, err := array.Iterate(t.Context())
	if err != nil {
		t.Fatalf("iterate array: %v", err)
	}
	value := sdk.NewIteratorValue(iterator)
	if _, ok := any(value).(runtime.Iterator); !ok {
		t.Fatal("iterator value does not implement runtime.Iterator")
	}
	if _, ok := any(value).(runtime.Iterable); !ok {
		t.Fatal("iterator value does not implement runtime.Iterable")
	}
	assertDoesNotImplement[runtime.Queryable](t, value, "runtime.Queryable")

	item, index, err := value.Next(t.Context())
	if err != nil || item != runtime.NewInt(1) || index != runtime.NewInt(0) {
		t.Fatalf("unexpected iterator result: value=%v index=%v err=%v", item, index, err)
	}
}

func TestSliceViewBoundsMutationAndCopy(t *testing.T) {
	data := []int{3, 1, 2}
	view := sdk.NewSliceView(data)

	if _, err := view.At(t.Context(), -1); err == nil {
		t.Fatal("expected negative index error")
	}
	if _, err := view.At(t.Context(), 3); err == nil {
		t.Fatal("expected upper-bound error")
	}
	if _, err := view.At(t.Context(), runtime.Int(math.MaxInt64)); err == nil {
		t.Fatal("expected maximum index error")
	}
	if value, found, err := view.LookupAt(t.Context(), 10); err != nil || found || value != runtime.None {
		t.Fatalf("unexpected missing lookup: value=%v found=%v err=%v", value, found, err)
	}

	if err := view.SetAt(t.Context(), 1, runtime.NewInt(4)); err != nil {
		t.Fatalf("set: %v", err)
	}
	if data[1] != 4 {
		t.Fatalf("backing slice was not updated: %v", data)
	}
	if err := view.SortAsc(t.Context()); err != nil {
		t.Fatalf("sort: %v", err)
	}
	if data[0] != 2 || data[1] != 3 || data[2] != 4 {
		t.Fatalf("sort did not update backing slice: %v", data)
	}

	removed, err := view.RemoveAt(t.Context(), 1)
	if err != nil || removed != runtime.NewInt(3) {
		t.Fatalf("remove: value=%v err=%v", removed, err)
	}
	if got := view.Target(); len(got) != 2 || got[0] != 2 || got[1] != 4 {
		t.Fatalf("unexpected view after remove: %v", got)
	}

	copy, ok := view.Copy().(*sdk.SliceView[int])
	if !ok {
		t.Fatalf("copy has type %T", view.Copy())
	}
	if copy.Hash() != view.Hash() || runtime.TypeName(copy.Type()) != runtime.TypeName(view.Type()) {
		t.Fatal("slice copy lost identity or type")
	}
	if err := copy.SetAt(t.Context(), 0, runtime.NewInt(9)); err != nil || view.Target()[0] != 9 {
		t.Fatalf("copy does not share live backing storage: %v, %v", view.Target(), err)
	}

	assertDoesNotImplement[runtime.List](t, view, "runtime.List")
	assertDoesNotImplement[runtime.Queryable](t, view, "runtime.Queryable")
}

func TestMapViewMutationIterationAndCopy(t *testing.T) {
	data := map[string]int{"one": 1}
	view := sdk.NewMapView(data)

	if err := view.Set(t.Context(), runtime.NewString("two"), runtime.NewInt(2)); err != nil {
		t.Fatalf("set: %v", err)
	}
	if data["two"] != 2 {
		t.Fatalf("backing map was not updated: %v", data)
	}
	if value, found, err := view.Lookup(t.Context(), runtime.NewString("two")); err != nil || !found || value != runtime.NewInt(2) {
		t.Fatalf("lookup: value=%v found=%v err=%v", value, found, err)
	}
	if err := view.RemoveKey(t.Context(), runtime.NewString("one")); err != nil {
		t.Fatalf("remove key: %v", err)
	}
	if _, exists := data["one"]; exists {
		t.Fatalf("key remains in backing map: %v", data)
	}

	iterator, err := view.Iterate(t.Context())
	if err != nil {
		t.Fatalf("iterate: %v", err)
	}
	value, key, err := iterator.Next(t.Context())
	if err != nil || value != runtime.NewInt(2) || key != runtime.NewString("two") {
		t.Fatalf("unexpected iterator result: value=%v key=%v err=%v", value, key, err)
	}

	var decoded map[string]int
	if err := sdk.Decode(t.Context(), view, &decoded); err != nil {
		t.Fatalf("decode map view: %v", err)
	}
	if decoded["two"] != 2 {
		t.Fatalf("unexpected decoded map: %v", decoded)
	}

	copy, ok := view.Copy().(*sdk.MapView[string, int])
	if !ok {
		t.Fatalf("copy has type %T", view.Copy())
	}
	if copy.Hash() != view.Hash() || runtime.TypeName(copy.Type()) != runtime.TypeName(view.Type()) {
		t.Fatal("map copy lost identity or type")
	}
	if err := copy.Set(t.Context(), runtime.NewString("three"), runtime.NewInt(3)); err != nil || data["three"] != 3 {
		t.Fatalf("copy does not share backing map: %v, %v", data, err)
	}

	assertDoesNotImplement[runtime.Map](t, view, "runtime.Map")
	assertDoesNotImplement[runtime.Queryable](t, view, "runtime.Queryable")
}

func TestNilCollectionViewsPreserveHostType(t *testing.T) {
	if got := runtime.TypeName(sdk.NewSliceView[int](nil).Type()); got == runtime.TypeNone.Name() {
		t.Fatalf("nil slice view lost its host type: %s", got)
	}
	if got := runtime.TypeName(sdk.NewMapView[string, int](nil).Type()); got == runtime.TypeNone.Name() {
		t.Fatalf("nil map view lost its host type: %s", got)
	}
}

func TestCollectionViewsUseCodecs(t *testing.T) {
	codecErr := errors.New("cannot encode record")
	codec := sdk.NewCodec[codecRecord](
		func(_ context.Context, record codecRecord) (runtime.Value, error) {
			if record.ID < 0 {
				return runtime.None, codecErr
			}

			return runtime.NewString(strconv.Itoa(record.ID)), nil
		},
		func(_ context.Context, value runtime.Value) (codecRecord, error) {
			text, err := runtime.CastArg[runtime.String](value, 0)
			if err != nil {
				return codecRecord{}, err
			}
			id, err := strconv.Atoi(text.String())
			return codecRecord{ID: id}, err
		})

	data := []codecRecord{{ID: 1}}
	view := sdk.NewSliceViewWithEncoding(data, codec)
	value, err := view.At(t.Context(), 0)
	if err != nil || value != runtime.NewString("1") {
		t.Fatalf("custom encode: value=%v err=%v", value, err)
	}
	if err := view.SetAt(t.Context(), 0, runtime.NewString("7")); err != nil || data[0].ID != 7 {
		t.Fatalf("custom decode: data=%v err=%v", data, err)
	}

	view = sdk.NewSliceViewWithEncoding([]codecRecord{{ID: -1}}, codec)
	_, err = view.At(t.Context(), 0)
	if !errors.Is(err, codecErr) {
		t.Fatalf("expected codec error, got %v", err)
	}
}

func TestNativeIteratorsEncodeValues(t *testing.T) {
	sliceIterator := sdk.NewSliceIterator([]int{3})
	value, key, err := sliceIterator.Next(t.Context())
	if err != nil || value != runtime.NewInt(3) || key != runtime.NewInt(0) {
		t.Fatalf("slice iterator: value=%v key=%v err=%v", value, key, err)
	}

	mapIterator := sdk.NewMapIterator(map[string]int{"three": 3})
	value, key, err = mapIterator.Next(t.Context())
	if err != nil || value != runtime.NewInt(3) || key != runtime.NewString("three") {
		t.Fatalf("map iterator: value=%v key=%v err=%v", value, key, err)
	}

	codecErr := errors.New("encode failed")
	failing := sdk.NewMapIteratorWithEncoding(
		map[string]int{"bad": 2},
		sdk.DefaultCodec[string](),
		sdk.NewCodec[int](func(context.Context, int) (runtime.Value, error) {
			return runtime.None, codecErr
		}, nil),
	)
	_, _, err = failing.Next(t.Context())
	if !errors.Is(err, codecErr) {
		t.Fatalf("expected conversion error, got %v", err)
	}

	_, _, err = sliceIterator.Next(t.Context())
	if !errors.Is(err, io.EOF) {
		t.Fatalf("expected EOF, got %v", err)
	}
}

func assertDoesNotImplement[T any](t *testing.T, value any, name string) {
	t.Helper()
	if _, ok := value.(T); ok {
		t.Fatalf("%T unexpectedly implements %s", value, name)
	}
}
