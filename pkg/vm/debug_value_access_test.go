package vm

import (
	"context"
	"errors"
	"fmt"
	"math"
	"strconv"
	"strings"
	"testing"

	"github.com/MontFerret/api/debugger"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
	"github.com/MontFerret/ferret/v2/pkg/vm/internal/data"
)

type debugMetadataValue struct {
	info runtime.DebugInfo
}

var _ runtime.DebugInspectable = debugMetadataValue{}

func (v debugMetadataValue) String() string {
	return "debug metadata value"
}

func (v debugMetadataValue) Hash() uint64 {
	return 0
}

func (v debugMetadataValue) Copy() runtime.Value {
	return v
}

func (v debugMetadataValue) DebugInfo() runtime.DebugInfo {
	return v.info
}

func TestDebugMillionElementBufferedPage(t *testing.T) {
	const size = 1_000_000

	for _, kind := range []string{"array", "dataset"} {
		t.Run(kind, func(t *testing.T) {
			value := bufferedDebugFixture(t, kind, size)
			access := NewDebugValueAccess()
			t.Run("describe", func(t *testing.T) {
				d := access.Describe(value)
				if d.TypeName != "Array" || d.Length != size || d.Children != indexedChildren(size) {
					t.Fatalf("description: %+v", d)
				}
			})
			t.Run("tail", func(t *testing.T) {
				page, err := access.ReadChildren(t.Context(), value, debugger.VariablePage{Filter: debugger.VariableFilterIndexed, Start: size - 20, Count: 20})
				if err != nil || len(page) != 20 {
					t.Fatalf("page length %d, error %v", len(page), err)
				}

				for i, child := range page {
					pos := size - 20 + i
					if child.Position != int64(pos) || child.Name != strconv.Itoa(pos) || child.Value != runtime.Int(pos) {
						t.Fatalf("child %d: %+v", i, child)
					}
				}
			})
		})
	}
}

func bufferedDebugFixture(t testing.TB, kind string, size int) runtime.Value {
	t.Helper()
	if kind == "array" {
		values := make([]runtime.Value, size)

		for i := range values {
			values[i] = runtime.Int(i)
		}

		return runtime.NewArrayOf(values)
	}

	set := data.NewDataSet(false)

	for i := range size {
		if err := set.Append(t.Context(), runtime.Int(i)); err != nil {
			t.Fatal(err)
		}
	}

	return set
}

func TestDebugPageBounds(t *testing.T) {
	a := NewDebugValueAccess()
	v := runtime.NewSizedArray(1111)

	for _, tc := range []struct {
		page debugger.VariablePage
		want int
	}{
		{debugger.VariablePage{}, 100}, {debugger.VariablePage{Start: 100, Count: 1000}, 1000},
		{debugger.VariablePage{Start: 1109, Count: 20}, 2}, {debugger.VariablePage{Start: math.MaxInt64 - 1, Count: 1}, 0},
		{debugger.VariablePage{Filter: debugger.VariableFilterNamed}, 0},
	} {
		got, err := a.ReadChildren(t.Context(), v, tc.page)
		if err != nil || len(got) != tc.want {
			t.Fatalf("%+v: len=%d err=%v", tc.page, len(got), err)
		}
	}

	for _, p := range []debugger.VariablePage{{Count: -1}, {Start: -1}, {Count: 1001}, {Start: math.MaxInt64}, {Filter: "?"}} {
		if got, err := a.ReadChildren(t.Context(), v, p); got != nil || !errors.Is(err, debugger.ErrInvalidVariablePage) {
			t.Fatalf("invalid %+v: %v %v", p, got, err)
		}
	}

	for _, v := range []runtime.Value{(*runtime.Array)(nil), (*runtime.Object)(nil), (*data.FastObject)(nil), (*data.DataSet)(nil), debugMetadataValue{}} {
		if _, err := a.Lookup(v, runtime.Int(0)); !errors.Is(err, runtime.ErrInvalidOperation) {
			t.Fatalf("lookup %T: %v", v, err)
		}

		if got, err := a.ReadChildren(t.Context(), v, debugger.VariablePage{}); got != nil || !errors.Is(err, errors.ErrUnsupported) {
			t.Fatalf("opaque %T: %v, %v", v, got, err)
		}
	}
}

func TestDebugNamedPage(t *testing.T) {
	for _, tc := range []struct {
		obj     runtime.Map
		name    string
		ordered bool
	}{
		{name: "map", obj: runtime.NewObject()}, {name: "shape", obj: data.NewShapeCache(0, 0).NewFastObject(), ordered: true}, {name: "dict", obj: data.NewShapeCache(0, 1).NewFastObject()},
	} {
		t.Run(tc.name, func(t *testing.T) {
			keys := []string{"\uFFFD", "\xff", "0", "a\x00\n", "é", "e\u0301", ""}

			for i, k := range keys {
				if err := tc.obj.Set(t.Context(), runtime.String(k), runtime.Int(i)); err != nil {
					t.Fatal(err)
				}
			}

			order := []int{6, 2, 3, 5, 4, 0, 1}
			if tc.ordered {
				order = []int{0, 1, 2, 3, 4, 5, 6}
			}

			a := NewDebugValueAccess()

			for i, k := range order {
				got, err := a.ReadChildren(t.Context(), tc.obj, debugger.VariablePage{Start: int64(i), Count: 1})
				if k == 1 {
					if got != nil || err == nil || errors.Is(err, errors.ErrUnsupported) {
						t.Fatalf("malformed: %v, %v", got, err)
					}

					continue
				}

				if err != nil || len(got) != 1 || got[0].Name != keys[k] || got[0].Value != runtime.Int(k) {
					t.Fatalf("position %d: %+v, %v", i, got, err)
				}
			}

			got, err := a.ReadChildren(t.Context(), tc.obj, debugger.VariablePage{Filter: debugger.VariableFilterIndexed})
			if err != nil || len(got) != 0 {
				t.Fatalf("indexed: %v, %v", got, err)
			}
		})
	}
}

func TestDebugObjectChildViews(t *testing.T) {
	access := NewDebugValueAccess()
	for _, tc := range []struct {
		object   runtime.Map
		name     string
		keys     int
		pageable bool
	}{
		{runtime.NewObject(), "small", 1, true},
		{runtime.NewObjectOf(128), "bounded_reservation", 128, true},
		{runtime.NewObjectOf(256), "reserved", 1, false},
		{data.NewShapeCache(64, 64).NewFastObject(), "bounded_dictionary", 128, true},
		{data.NewShapeCache(0, 1).NewFastObject(), "large_dictionary", 129, false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			for i := range tc.keys {
				if err := tc.object.Set(t.Context(), runtime.String(strconv.Itoa(i)), runtime.Int(i)); err != nil {
					t.Fatal(err)
				}
			}

			views := access.Describe(tc.object).Children
			if !views.Named.Known || views.Named.Count != int64(tc.keys) || views.Named.CanExpand() != tc.pageable ||
				!views.Indexed.Known || views.Indexed.Count != 0 || views.Indexed.CanExpand() {
				t.Fatalf("unexpected child views: %+v", views)
			}

			for _, start := range []int64{0, 10} {
				page, err := access.ReadChildren(t.Context(), tc.object, debugger.VariablePage{Filter: debugger.VariableFilterIndexed, Start: start})
				if err != nil || len(page) != 0 {
					t.Errorf("indexed at %d: %v, %v", start, page, err)
				}
			}

			page, err := access.ReadChildren(t.Context(), tc.object, debugger.VariablePage{Filter: debugger.VariableFilterIndexed, Start: -1})
			if page != nil || !errors.Is(err, debugger.ErrInvalidVariablePage) {
				t.Fatalf("invalid indexed range: %v, %v", page, err)
			}

			for _, filter := range []debugger.VariableFilter{debugger.VariableFilterAll, debugger.VariableFilterNamed} {
				page, err := access.ReadChildren(t.Context(), tc.object, debugger.VariablePage{Filter: filter, Count: 1})
				if !tc.pageable {
					if page != nil || !errors.Is(err, errors.ErrUnsupported) {
						t.Errorf("filter %q: %v, %v", filter, page, err)
					}

					continue
				}

				if err != nil || len(page) != 1 || page[0] != (DebugChild{Kind: debugger.VariableFilterNamed, Name: "0", Position: 0, Value: runtime.Int(0)}) {
					t.Errorf("filter %q: %v, %v", filter, page, err)
				}
			}
		})
	}
}

func TestNativeObjectAliasesShareInspectionBounds(t *testing.T) {
	owner := runtime.NewObject()
	alias := *owner

	for i := range 200 {
		_ = alias.Set(t.Context(), runtime.String(strconv.Itoa(i)), runtime.True)
	}

	for i := range 199 {
		_ = alias.RemoveKey(t.Context(), runtime.String(strconv.Itoa(i)))
	}

	if NewDebugValueAccess().Describe(owner).Children.Named.CanExpand() {
		t.Fatal("copy of the Go value lost shared allocation history")
	}

	_ = alias.Clear(t.Context())
	_ = alias.Set(t.Context(), runtime.String("fresh"), runtime.True)
	if !alias.DebugFieldsPageable() || owner.DebugFieldsPageable() {
		t.Fatal("clear must detach the new map and its bounds together")
	}
}

func TestUnorderedMapRemovalBudget(t *testing.T) {
	access := NewDebugValueAccess()
	for _, obj := range []runtime.Map{runtime.NewObject(), data.NewShapeCache(0, 1).NewFastObject()} {
		// Force the fast object into dictionary mode, then start both cases
		// with a fresh unordered backing and no prior deletions.
		_ = obj.Set(t.Context(), runtime.String("a"), runtime.True)
		_ = obj.Set(t.Context(), runtime.String("b"), runtime.True)
		_ = obj.Clear(t.Context())
		for removals := range 129 {
			_ = obj.Set(t.Context(), runtime.String("x"), runtime.True)
			view := access.Describe(obj).Children.Named
			if !view.Known || view.Count != 1 || view.CanExpand() != (removals < 128) {
				t.Fatalf("%T after %d removals: %+v", obj, removals, view)
			}

			_ = obj.RemoveKey(t.Context(), runtime.String("missing"))
			_ = obj.RemoveKey(t.Context(), runtime.String("x"))
		}

		_ = obj.Clear(t.Context())
		_ = obj.Set(t.Context(), runtime.String("fresh"), runtime.True)
		page, err := access.ReadChildren(t.Context(), obj, debugger.VariablePage{})
		if err != nil || len(page) != 1 || page[0].Name != "fresh" {
			t.Fatalf("%T after Clear: %v, %v", obj, page, err)
		}
	}
}

func TestDebugObjectBackingBounds(t *testing.T) {
	a := NewDebugValueAccess()

	for _, o := range []runtime.Map{runtime.NewObject(), data.NewShapeCache(0, 1).NewFastObject(), data.NewShapeCache(0, 0).NewFastObject()} {
		for i := range 129 {
			_ = o.Set(t.Context(), runtime.String(strconv.Itoa(i)), runtime.Int(i))
		}

		_, shape := o.(*data.FastObject)
		shape = shape && o.(*data.FastObject).Shape() != nil
		d := a.Describe(o)
		if d.Children.Named.Count != 129 || d.Children.Named.CanExpand() != shape {
			t.Fatalf("dense %T: %+v", o, d)
		}

		if shape {
			got, err := a.ReadChildren(t.Context(), o, debugger.VariablePage{Start: 128, Count: 1})
			if err != nil || len(got) != 1 || got[0].Name != "128" {
				t.Fatalf("dense shape tail: %v, %v", got, err)
			}
		}

		for i := range 128 {
			_ = o.RemoveKey(t.Context(), runtime.String(strconv.Itoa(i)))
		}

		d = a.Describe(o)
		if d.Children.Named.Count != 1 || d.Children.Named.CanExpand() {
			t.Fatalf("sparse %T: %+v", o, d)
		}

		_ = o.Clear(t.Context())
		_ = o.Set(t.Context(), runtime.String("new"), runtime.True)
		if !a.Describe(o).Children.Named.CanExpand() {
			t.Fatalf("fresh small %T", o)
		}
	}

	o := runtime.NewObjectOf(1 << 16)
	_ = o.Set(t.Context(), runtime.String("k"), runtime.True)
	if a.Describe(o).Children.Named.CanExpand() {
		t.Fatal("oversized allocation")
	}

	o = runtime.NewObject()
	_ = o.Set(t.Context(), runtime.String(strings.Repeat("k", 1<<20)), runtime.True)
	if a.Describe(o).Children.Named.CanExpand() {
		t.Fatal("unbounded sorting work")
	}

	sparse := data.NewShapeCache(0, 0).NewFastObject()

	for _, k := range []string{"b", "a", "c"} {
		_ = sparse.Set(t.Context(), runtime.String(k), runtime.True)
	}

	_ = sparse.RemoveKey(t.Context(), runtime.String("a"))
	got, err := a.ReadChildren(t.Context(), sparse, debugger.VariablePage{Start: 1, Count: 1})
	if err != nil || len(got) != 1 || got[0].Name != "c" {
		t.Fatalf("small sparse shape: %v, %v", got, err)
	}
}

type debugProvider struct {
	read func(context.Context, debugger.VariableFilter, int64) (string, runtime.Value, bool, error)
	debugMetadataValue
	info runtime.DebugChildrenInfo
}

func (p *debugProvider) DebugChildrenInfo() runtime.DebugChildrenInfo { return p.info }
func (p *debugProvider) DebugChildAt(ctx context.Context, kind debugger.VariableFilter, index int64) (string, runtime.Value, bool, error) {
	return p.read(ctx, kind, index)
}

func TestDebugProviderViews(t *testing.T) {
	a := NewDebugValueAccess()
	p := &debugProvider{info: runtime.DebugChildrenInfo{Named: runtime.DebugChildView{Known: true, Count: 2, Readable: true}, Indexed: runtime.DebugChildView{Readable: true}}}
	var reads []string
	p.read = func(_ context.Context, k debugger.VariableFilter, i int64) (string, runtime.Value, bool, error) {
		reads = append(reads, fmt.Sprintf("%s/%d", k, i))

		return strconv.FormatInt(i, 10), runtime.Int(i), i < 3, nil
	}

	if d := a.Describe(p); len(reads) != 0 || d.Children != p.info {
		t.Fatalf("description: %+v, reads %v", d, reads)
	}

	got, err := a.ReadChildren(t.Context(), p, debugger.VariablePage{Start: 1, Count: 5})
	if err != nil || len(got) != 4 || fmt.Sprint(reads) != "[named/1 indexed/0 indexed/1 indexed/2 indexed/3]" {
		t.Fatalf("mixed: %v, reads %v, %v", got, reads, err)
	}

	for i, pos := range []int64{1, 0, 1, 2} {
		kind := debugger.VariableFilterIndexed
		if i == 0 {
			kind = debugger.VariableFilterNamed
		}

		want := DebugChild{Kind: kind, Position: pos, Name: strconv.FormatInt(pos, 10), Value: runtime.Int(pos)}
		if got[i] != want {
			t.Fatalf("combined page entry %d: %+v, want %+v", i, got[i], want)
		}
	}

	p.info.Named.Known = false
	reads = nil
	if got, err := a.ReadChildren(t.Context(), p, debugger.VariablePage{}); got != nil || !errors.Is(err, errors.ErrUnsupported) || len(reads) != 0 {
		t.Fatalf("unknown boundary: %v %v", got, err)
	}

	for _, k := range []debugger.VariableFilter{debugger.VariableFilterNamed, debugger.VariableFilterIndexed} {
		if got, err := a.ReadChildren(t.Context(), p, debugger.VariablePage{Filter: k, Start: 1}); err != nil || len(got) != 2 {
			t.Fatalf("filtered %s: %v, %v", k, got, err)
		}
	}

	p.info.Named = runtime.DebugChildView{Known: true, Count: math.MaxInt64, Readable: true}
	p.info.Indexed = runtime.DebugChildView{Known: true, Count: 10, Readable: true}
	p.read = func(_ context.Context, k debugger.VariableFilter, i int64) (string, runtime.Value, bool, error) {
		return "last", runtime.Int(i), true, nil
	}

	got, err = a.ReadChildren(t.Context(), p, debugger.VariablePage{Start: math.MaxInt64 - 1, Count: 1})
	if err != nil || len(got) != 1 || got[0].Position != math.MaxInt64-1 {
		t.Fatalf("count overflow: %v %v", got, err)
	}

	p.info = runtime.DebugChildrenInfo{Indexed: runtime.DebugChildView{Known: true, Count: 42}}
	p.read = nil
	if d := a.Describe(p); d.Children.Indexed.CanExpand() || d.Children.Indexed.Count != 42 {
		t.Fatalf("summary-only: %+v", d)
	}
}

func TestBufferedPageReadPositions(t *testing.T) {
	const size = 1_000_000

	for _, kind := range []string{"array", "dataset"} {
		t.Run(kind, func(t *testing.T) {
			value := bufferedDebugFixture(t, kind, size)
			var lookup func(context.Context, runtime.Int) (runtime.Value, bool, error)
			switch v := value.(type) {
			case *runtime.Array:
				lookup = v.LookupAt
			case *data.DataSet:
				lookup = v.BufferedAt
			}

			for _, start := range []int64{1000, size - 20} {
				var reads []runtime.Int
				counted := func(ctx context.Context, index runtime.Int) (runtime.Value, bool, error) {
					reads = append(reads, index)

					return lookup(ctx, index)
				}

				children, err := readBufferedPage(t.Context(), debugger.VariablePage{Filter: debugger.VariableFilterIndexed, Start: start, Count: 20}, size, counted)
				if err != nil || len(children) != 20 || len(reads) != 20 {
					t.Fatalf("reads=%v children=%d err=%v", reads, len(children), err)
				}

				for i, idx := range reads {
					if idx != runtime.Int(start+int64(i)) {
						t.Fatalf("read %d at %d", idx, i)
					}
				}
			}
		})
	}
}
