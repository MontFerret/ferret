package objects_test

import (
	"context"
	"io"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

type (
	hostMap struct {
		runtime.Map
		keys        runtime.List
		values      runtime.List
		empty       runtime.Map
		cloned      runtime.Cloneable
		keysErr     error
		valuesErr   error
		lookupErr   error
		containsErr error
		walkErr     error
		emptyErr    error
		cloneErr    error
		setErr      error
		removeErr   error
		equalityErr error
		readingKeys *bool
		walkCalls   *int
		setCalls    *int
		removeCalls *int
		walk        func(context.Context, runtime.KeyReadablePredicate) error
	}
	hostList struct {
		runtime.List
		lengthErr error
		accessErr error
		walkErr   error
		reading   *bool
		failAt    runtime.Int
	}
	entryPair struct {
		items [2]runtime.Value
		runtime.Value
		lengthErr error
		accessErr error
		length    runtime.Int
		failAt    runtime.Int
	}
	entryStream struct {
		runtime.Value
		iterateErr     error
		nextErr        error
		closeErr       error
		entries        []runtime.Value
		iterations     int
		sourceCloses   int
		iteratorCloses int
	}
	entryIterator struct {
		stream *entryStream
		index  int
	}
	copyOnlyValue struct {
		runtime.Value
		copies *int
		number int
	}
	cloneValue struct {
		runtime.Value
		clones   *int
		cloneErr error
		number   int
	}
)

func (m *hostMap) Keys(ctx context.Context) (runtime.List, error) {
	if m.keysErr != nil {
		return nil, m.keysErr
	}

	if m.keys != nil {
		return m.keys, nil
	}

	return m.Map.Keys(ctx)
}

func (m *hostMap) Values(ctx context.Context) (runtime.List, error) {
	if m.valuesErr != nil {
		return nil, m.valuesErr
	}

	if m.values != nil {
		return m.values, nil
	}

	return m.Map.Values(ctx)
}

func (m *hostMap) Lookup(ctx context.Context, key runtime.Value) (runtime.Value, bool, error) {
	if m.lookupErr != nil {
		return runtime.None, false, m.lookupErr
	}

	return m.Map.Lookup(ctx, key)
}

func (m *hostMap) ContainsKey(ctx context.Context, key runtime.Value) (runtime.Boolean, error) {
	if m.containsErr != nil {
		return false, m.containsErr
	}

	return m.Map.ContainsKey(ctx, key)
}

func (m *hostMap) ForEach(ctx context.Context, fn runtime.KeyReadablePredicate) error {
	if m.walkCalls != nil {
		*m.walkCalls++
	}

	if m.walkErr != nil {
		return m.walkErr
	}

	if m.walk != nil {
		return m.walk(ctx, fn)
	}

	return m.Map.ForEach(ctx, fn)
}

func (m *hostMap) Empty(ctx context.Context) (runtime.Map, error) {
	if m.emptyErr != nil {
		return nil, m.emptyErr
	}

	if m.empty != nil {
		return m.empty, nil
	}

	return m.Map.Empty(ctx)
}

func (m *hostMap) Clone(ctx context.Context) (runtime.Cloneable, error) {
	if m.cloneErr != nil {
		return nil, m.cloneErr
	}

	if m.cloned != nil {
		return m.cloned, nil
	}

	return m.Map.Clone(ctx)
}

func (m *hostMap) Set(ctx context.Context, key, value runtime.Value) error {
	if m.setCalls != nil {
		*m.setCalls++
	}

	if m.setErr != nil {
		return m.setErr
	}

	return m.Map.Set(ctx, key, value)
}

func (m *hostMap) RemoveKey(ctx context.Context, key runtime.Value) error {
	if m.removeCalls != nil {
		*m.removeCalls++
	}

	if m.readingKeys != nil && *m.readingKeys {
		return runtime.ErrInvalidOperation
	}

	if m.removeErr != nil {
		return m.removeErr
	}

	return m.Map.RemoveKey(ctx, key)
}

func (l *hostList) Length(ctx context.Context) (runtime.Int, error) {
	if l.lengthErr != nil {
		return 0, l.lengthErr
	}

	return l.List.Length(ctx)
}

func (l *hostList) At(ctx context.Context, index runtime.Int) (runtime.Value, error) {
	if l.accessErr != nil && index == l.failAt {
		return runtime.None, l.accessErr
	}

	return l.List.At(ctx, index)
}

func (l *hostList) ForEach(ctx context.Context, fn runtime.IndexReadablePredicate) error {
	if l.reading != nil {
		*l.reading = true
		defer func() { *l.reading = false }()
	}

	if l.walkErr != nil {
		return l.walkErr
	}

	return l.List.ForEach(ctx, fn)
}

func (p *entryPair) Length(context.Context) (runtime.Int, error) { return p.length, p.lengthErr }

func (p *entryPair) At(_ context.Context, index runtime.Int) (runtime.Value, error) {
	if p.accessErr != nil && index == p.failAt {
		return runtime.None, p.accessErr
	}

	return p.items[index], nil
}

func (s *entryStream) Iterate(context.Context) (runtime.Iterator, error) {
	s.iterations++
	if s.iterateErr != nil {
		return nil, s.iterateErr
	}

	return &entryIterator{stream: s}, nil
}

func (s *entryStream) Close() error {
	s.sourceCloses++

	return nil
}

func (i *entryIterator) Next(ctx context.Context) (runtime.Value, runtime.Value, error) {
	if err := ctx.Err(); err != nil {
		return runtime.None, runtime.None, err
	}

	if i.index == len(i.stream.entries) {
		if i.stream.nextErr != nil {
			return runtime.None, runtime.None, i.stream.nextErr
		}

		return runtime.None, runtime.None, io.EOF
	}

	value := i.stream.entries[i.index]
	i.index++

	return value, runtime.Int(i.index - 1), nil
}

func (i *entryIterator) Close() error {
	i.stream.iteratorCloses++

	return i.stream.closeErr
}

func (v *copyOnlyValue) Copy() runtime.Value {
	*v.copies++
	copied := *v

	return &copied
}

func (v *cloneValue) Clone(context.Context) (runtime.Cloneable, error) {
	*v.clones++
	if v.cloneErr != nil {
		return nil, v.cloneErr
	}

	copied := *v

	return &copied, nil
}

func (m *hostMap) Equal(ctx context.Context, value runtime.Value) (bool, error) {
	if m.equalityErr != nil {
		return false, m.equalityErr
	}

	return m.Map.Equal(ctx, value)
}
