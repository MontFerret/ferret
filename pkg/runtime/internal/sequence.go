package internal

import (
	"context"
	"github.com/MontFerret/ferret/pkg/runtime/core"
)

type (
	Sequence struct {
		data *Array
	}

	sequenceIterator struct {
		data   *Array
		length int
		pos    int
	}
)

func NewSequence(data *Array) *Sequence {
	return &Sequence{data}
}

func (iter *sequenceIterator) HasNext(_ context.Context) (bool, error) {
	return iter.length > iter.pos, nil
}

func (iter *sequenceIterator) Next(_ context.Context) (value core.Value, key core.Value, err error) {
	iter.pos++

	// TODO: Make it less ugly
	return iter.data.Get(iter.pos - 1).(*KeyValuePair).Value, core.NewInt(iter.pos - 1), nil
}

func (s *Sequence) Iterate(_ context.Context) (core.Iterator, error) {
	return &sequenceIterator{data: s.data, length: s.data.Length(), pos: 0}, nil
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

func (s *Sequence) Copy() core.Value {
	//TODO implement me
	panic("implement me")
}
