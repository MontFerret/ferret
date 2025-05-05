package internal

import (
	"context"
	"github.com/MontFerret/ferret/pkg/runtime"
)

type (
	Sequence struct {
		data runtime.List
	}

	sequenceIterator struct {
		data   runtime.List
		length runtime.Int
		pos    runtime.Int
	}
)

func NewSequence(data runtime.List) *Sequence {
	return &Sequence{data}
}

func (iter *sequenceIterator) HasNext(_ context.Context) (bool, error) {
	return iter.length > iter.pos, nil
}

func (iter *sequenceIterator) Next(ctx context.Context) (value runtime.Value, key runtime.Value, err error) {
	iter.pos++

	val, err := iter.data.Get(ctx, iter.pos-1)

	if err != nil {
		return nil, nil, err
	}

	kv := val.(*KeyValuePair)

	return kv.Value, iter.pos - 1, nil
}

func (s *Sequence) Iterate(ctx context.Context) (runtime.Iterator, error) {
	length, err := s.data.Length(ctx)

	if err != nil {
		return nil, err
	}

	return &sequenceIterator{data: s.data, length: length, pos: 0}, nil
}

func (s *Sequence) MarshalJSON() ([]byte, error) {
	//TODO implement me
	panic("implement me")
}

func (s *Sequence) String() string {
	return "[Sequence]"
}

func (s *Sequence) Unwrap() interface{} {
	//TODO implement me
	panic("implement me")
}

func (s *Sequence) Hash() uint64 {
	//TODO implement me
	panic("implement me")
}

func (s *Sequence) Copy() runtime.Value {
	//TODO implement me
	panic("implement me")
}
