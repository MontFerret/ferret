package collections_test

import (
	"context"
	"errors"
	"io"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

type (
	// This source deliberately has no measurement, collection, or equality methods.
	scanSource struct {
		expected     context.Context
		iterateErr   error
		nextErr      error
		closeErr     error
		onIterate    func()
		onNext       func(int)
		onClose      func()
		values       []runtime.Value
		failAt       int
		iterations   int
		nexts        int
		closes       int
		sourceCloses int
	}

	scanIterator struct {
		source *scanSource
		index  int
	}

	measuredSource struct {
		lengthErr error
		*scanSource
		onLength    func()
		length      runtime.Int
		lengthCalls int
	}

	measurableOnly struct {
		lengthCalls int
	}

	containableSource struct {
		containsErr error
		*scanSource
		onContains    func()
		containsCalls int
		contains      runtime.Boolean
	}

	scanValue struct {
		equalErr error
		onEqual  func()
		closes   int
		copies   int
		clones   int
	}
)

func (s *scanSource) String() string      { return "read-only stream" }
func (s *scanSource) Hash() uint64        { return 1 }
func (s *scanSource) Copy() runtime.Value { return s }

func (s *scanSource) Iterate(ctx context.Context) (runtime.Iterator, error) {
	s.iterations++
	if s.expected != nil && ctx != s.expected {
		return nil, errors.New("iterator context changed")
	}

	if s.onIterate != nil {
		s.onIterate()
	}

	if s.iterateErr != nil {
		return nil, s.iterateErr
	}

	return &scanIterator{source: s}, nil
}

func (s *scanSource) Close() error {
	s.sourceCloses++

	return nil
}

// Next intentionally ignores cancellation to exercise the caller's safepoints.
func (i *scanIterator) Next(ctx context.Context) (runtime.Value, runtime.Value, error) {
	s := i.source
	s.nexts++
	if s.expected != nil && ctx != s.expected {
		return runtime.None, runtime.None, errors.New("next context changed")
	}

	if s.onNext != nil {
		s.onNext(i.index)
	}

	if s.nextErr != nil && i.index == s.failAt {
		return runtime.None, runtime.None, s.nextErr
	}

	if i.index == len(s.values) {
		return runtime.None, runtime.None, io.EOF
	}

	value := s.values[i.index]
	i.index++

	return value, runtime.Int(i.index - 1), nil
}

func (i *scanIterator) Close() error {
	i.source.closes++
	if i.source.onClose != nil {
		i.source.onClose()
	}

	return i.source.closeErr
}

func (s *measuredSource) Length(ctx context.Context) (runtime.Int, error) {
	s.lengthCalls++
	if s.expected != nil && ctx != s.expected {
		return 0, errors.New("length context changed")
	}

	if s.onLength != nil {
		s.onLength()
	}

	return s.length, s.lengthErr
}

func (s *measurableOnly) String() string      { return "measurable only" }
func (s *measurableOnly) Hash() uint64        { return 1 }
func (s *measurableOnly) Copy() runtime.Value { return s }

func (s *measurableOnly) Length(context.Context) (runtime.Int, error) {
	s.lengthCalls++

	return 42, nil
}

func (s *containableSource) Contains(ctx context.Context, _ runtime.Value) (runtime.Boolean, error) {
	s.containsCalls++
	if s.expected != nil && ctx != s.expected {
		return false, errors.New("contains context changed")
	}

	if s.onContains != nil {
		s.onContains()
	}

	return s.contains, s.containsErr
}

func (v *scanValue) String() string { return "scan value" }
func (v *scanValue) Hash() uint64   { return 7 }

func (v *scanValue) Copy() runtime.Value {
	v.copies++

	return v
}

func (v *scanValue) Clone(context.Context) (runtime.Cloneable, error) {
	v.clones++

	return v, nil
}

func (v *scanValue) Close() error {
	v.closes++

	return nil
}

func (v *scanValue) Equal(context.Context, runtime.Value) (bool, error) {
	if v.onEqual != nil {
		v.onEqual()
	}

	return true, v.equalErr
}
