package internal

import (
	"context"

	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
)

type (
	Sorter struct {
		data *values.Array
	}

	sorterIterator struct {
		data   *values.Array
		length int
		pos    int
	}
)

func NewSorter(data *values.Array) *Sorter {
	return &Sorter{data}
}

func (iter *sorterIterator) HasNext(_ context.Context) (bool, error) {
	return iter.length > iter.pos, nil
}

func (iter *sorterIterator) Next(_ context.Context) (value core.Value, key core.Value, err error) {
	iter.pos++

	// TODO: Make it less ugly
	return iter.data.Get(iter.pos - 1).(*Tuple).First, values.NewInt(iter.pos - 1), nil
}

func (s *Sorter) Iterate(_ context.Context) (core.Iterator, error) {
	return &sorterIterator{data: s.data, length: s.data.Length(), pos: 0}, nil
}

func (s *Sorter) MarshalJSON() ([]byte, error) {
	//TODO implement me
	panic("implement me")
}

func (s *Sorter) String() string {
	return "[Sorter]"
}

func (s *Sorter) Unwrap() interface{} {
	//TODO implement me
	panic("implement me")
}

func (s *Sorter) Hash() uint64 {
	//TODO implement me
	panic("implement me")
}

func (s *Sorter) Copy() core.Value {
	//TODO implement me
	panic("implement me")
}
