package mem

import (
	"errors"
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

func TestAdoptCloser_SkipsComparableCloser(t *testing.T) {
	// *testCloser is a pointer type, always comparable — no wrapping needed
	c := newTestCloser("x")
	adopted := AdoptCloser(c)

	if _, ok := adopted.(*ManagedResource); ok {
		t.Fatal("expected pointer-comparable closer to NOT be wrapped")
	}

	if adopted != c {
		t.Fatal("expected same value back for comparable closer")
	}
}

func TestAdoptCloser_WrapsNonComparableCloser(t *testing.T) {
	// nonComparableVal is a value type containing a slice — not comparable
	c := newNonComparableVal("x")
	adopted := AdoptCloser(c)

	m, ok := adopted.(*ManagedResource)
	if !ok {
		t.Fatal("expected non-comparable closer to be wrapped in ManagedResource")
	}

	if m.Unwrap().String() != "x" {
		t.Fatalf("expected inner value string 'x', got %q", m.Unwrap().String())
	}
}

func TestAdoptCloser_SkipsNonCloser(t *testing.T) {
	val := runtime.NewString("hello")
	adopted := AdoptCloser(val)

	if adopted != val {
		t.Fatal("expected same value back for non-closer")
	}
}

func TestAdoptCloser_SkipsNil(t *testing.T) {
	adopted := AdoptCloser(nil)

	if adopted != nil {
		t.Fatal("expected nil back for nil input")
	}
}

func TestAdoptCloser_SkipsAlreadyManaged(t *testing.T) {
	c := newNonComparableVal("x")
	first := AdoptCloser(c)
	second := AdoptCloser(first)

	if first != second {
		t.Fatal("expected no double-wrapping")
	}
}

func TestManagedResource_Close(t *testing.T) {
	closed := new(int)
	c := nonComparableVal{data: []byte("x"), closed: closed}
	var val runtime.Value = c
	m := AdoptCloser(val).(*ManagedResource)

	if err := m.Close(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if *closed != 1 {
		t.Fatalf("expected close count 1, got %d", *closed)
	}
}

func TestManagedResource_CloseError(t *testing.T) {
	want := errors.New("close failed")
	c := nonComparableErrVal{
		nonComparableVal: nonComparableVal{data: []byte("x"), closed: new(int)},
		err:              want,
	}
	var val runtime.Value = c
	m := AdoptCloser(val).(*ManagedResource)

	got := m.Close()
	if got != want {
		t.Fatalf("expected error %v, got %v", want, got)
	}
}

func TestManagedResource_String(t *testing.T) {
	c := newNonComparableVal("hello")
	m := AdoptCloser(c).(*ManagedResource)

	if m.String() != "hello" {
		t.Fatalf("expected 'hello', got %q", m.String())
	}
}

func TestManagedResource_Hash(t *testing.T) {
	c := newNonComparableVal("x")
	m := AdoptCloser(c).(*ManagedResource)

	if m.Hash() != 0 {
		t.Fatal("hash should delegate to inner value")
	}
}

func TestManagedResource_Copy(t *testing.T) {
	c := newNonComparableVal("x")
	m := AdoptCloser(c).(*ManagedResource)

	cp := m.Copy()
	if _, ok := cp.(*ManagedResource); ok {
		t.Fatal("copy should not be a ManagedResource")
	}
}

func TestManagedResource_Unwrap(t *testing.T) {
	c := newNonComparableVal("x")
	m := AdoptCloser(c).(*ManagedResource)

	inner := m.Unwrap()
	if inner.String() != "x" {
		t.Fatalf("unwrap should return inner value, got %q", inner.String())
	}
}

func TestManagedResource_PointerComparable(t *testing.T) {
	c1 := newNonComparableVal("x")
	c2 := newNonComparableVal("y")
	m1 := AdoptCloser(c1).(*ManagedResource)
	m2 := AdoptCloser(c2).(*ManagedResource)

	// Different wrappers are distinct map keys
	seen := make(map[*ManagedResource]struct{})
	seen[m1] = struct{}{}

	if _, ok := seen[m1]; !ok {
		t.Fatal("same pointer should be found in map")
	}

	if _, ok := seen[m2]; ok {
		t.Fatal("different pointer should not be found in map")
	}
}

func TestManagedResource_UsableAsCloserMapKey(t *testing.T) {
	c := newNonComparableVal("x")
	m := AdoptCloser(c).(*ManagedResource)

	// ManagedResource can be used as io.Closer map key (pointer-comparable)
	seen := make(map[interface{ Close() error }]struct{})
	seen[m] = struct{}{}

	if _, ok := seen[m]; !ok {
		t.Fatal("ManagedResource should be usable as map key")
	}
}

// nonComparableVal is a value type (not pointer) whose dynamic type contains
// a slice, making it non-comparable when held as an interface value. This
// triggers ManagedResource wrapping in AdoptCloser.
type nonComparableVal struct {
	data   []byte // slice makes this non-comparable
	closed *int
}

func newNonComparableVal(s string) runtime.Value {
	// Return as runtime.Value so the interface holds the value type, not a pointer
	return nonComparableVal{data: []byte(s), closed: new(int)}
}

func (c nonComparableVal) Close() error {
	*c.closed++
	return nil
}

func (c nonComparableVal) String() string {
	return string(c.data)
}

func (c nonComparableVal) Hash() uint64 {
	return 0
}

func (c nonComparableVal) Copy() runtime.Value {
	return nonComparableVal{data: append([]byte(nil), c.data...), closed: new(int)}
}

// nonComparableErrVal returns an error on Close.
type nonComparableErrVal struct {
	nonComparableVal
	err error
}

func (c nonComparableErrVal) Close() error {
	return c.err
}
