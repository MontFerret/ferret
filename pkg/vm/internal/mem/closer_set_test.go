package mem

import (
	"errors"
	"io"
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

// errCloser is a test closer that returns a fixed error.
type errCloser struct {
	err error
}

func (c *errCloser) Close() error {
	return c.err
}

type sliceCloser []int

func (sliceCloser) Close() error {
	return nil
}

type udResource struct {
	err    error
	name   string
	id     uint64
	closed bool
}

func newUdResource(id uint64, name string) *udResource {
	return &udResource{
		id:   id,
		name: name,
	}
}

func (u *udResource) Close() error {
	return u.err
}

func (u *udResource) ResourceID() uint64 {
	return u.id
}

func (u *udResource) String() string {
	return u.name
}

func (u *udResource) Hash() uint64 {
	//TODO implement me
	panic("implement me")
}

func (u *udResource) Copy() runtime.Value {
	return &udResource{}
}

func TestCloserSetAddIsIdempotent(t *testing.T) {
	var s CloserSet
	c := newTestCloser("dup")

	if !s.Add(c) {
		t.Fatal("expected first Add to succeed")
	}

	if s.Add(c) {
		t.Fatal("expected duplicate Add to return false")
	}

	if got := s.Len(); got != 1 {
		t.Fatalf("expected 1 closer, got %d", got)
	}

	r := newUdResource(1, "res")

	if !s.Add(r) {
		t.Fatal("expected first Add to succeed")
	}

	if s.Add(r) {
		t.Fatal("expected duplicate Add to return false")
	}

	if got := s.Len(); got != 2 {
		t.Fatalf("expected 2 closers, got %d", got)
	}
}

func TestCloserSetAddRejectsNil(t *testing.T) {
	var s CloserSet

	if s.Add(nil) {
		t.Fatal("expected nil closer to be rejected")
	}

	if got := s.Len(); got != 0 {
		t.Fatalf("expected 0 closers, got %d", got)
	}
}

func TestCloserSetCloseAllJoinsErrors(t *testing.T) {
	var s CloserSet
	e1 := errors.New("err1")
	e2 := errors.New("err2")

	s.Add(&errCloser{err: e1})
	s.Add(&errCloser{err: e2})

	err := s.CloseAll()
	if err == nil {
		t.Fatal("expected error from CloseAll")
	}

	if !errors.Is(err, e1) || !errors.Is(err, e2) {
		t.Fatalf("expected joined errors, got %v", err)
	}

	if got := s.Len(); got != 0 {
		t.Fatalf("expected reset after CloseAll, got %d closers", got)
	}
}

func TestCloserSetMergeDeduplicates(t *testing.T) {
	var a, b CloserSet
	shared := newTestCloser("shared")
	only := newTestCloser("only")

	a.Add(shared)
	b.Add(shared)
	b.Add(only)

	a.Merge(&b)

	if got := a.Len(); got != 2 {
		t.Fatalf("expected 2 closers after merge, got %d", got)
	}

	if got := b.Len(); got != 0 {
		t.Fatalf("expected source reset after merge, got %d", got)
	}
}

func TestCloserSetForEachPreservesOrder(t *testing.T) {
	var s CloserSet
	c1 := newTestCloser("first")
	c2 := newTestCloser("second")
	c3 := newTestCloser("third")

	s.Add(c1)
	s.Add(c2)
	s.Add(c3)

	var order []io.Closer
	s.ForEach(func(c io.Closer) {
		order = append(order, c)
	})

	if len(order) != 3 || order[0] != c1 || order[1] != c2 || order[2] != c3 {
		t.Fatalf("expected insertion order, got %v", order)
	}
}

func TestCloserSetCloseAllCallsEveryCloser(t *testing.T) {
	var s CloserSet
	c1 := newTestCloser("a")
	c2 := newTestCloser("b")

	s.Add(c1)
	s.Add(c2)

	if err := s.CloseAll(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if c1.closed != 1 || c2.closed != 1 {
		t.Fatalf("expected each closer called once, got c1=%d c2=%d", c1.closed, c2.closed)
	}
}

func TestCloserSetAddAcceptsNonComparableValueClosers(t *testing.T) {
	var s CloserSet
	closer := sliceCloser{1, 2, 3}

	if !s.Add(closer) {
		t.Fatal("expected first add to succeed")
	}

	if !s.Add(closer) {
		t.Fatal("expected second add to succeed without dedupe")
	}

	if got, want := s.Len(), 2; got != want {
		t.Fatalf("expected %d closers after adding non-comparable closer twice, got %d", want, got)
	}
}
