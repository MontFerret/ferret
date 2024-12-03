package internal

import (
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
)

type Stack struct {
	data []core.Value
}

func NewStack(size int) *Stack {
	return &Stack{data: make([]core.Value, 0, size)}
}

func (s *Stack) Length() int {
	return len(s.data)
}

func (s *Stack) Push(value core.Value) {
	s.data = append(s.data, value)
}

func (s *Stack) Pop() core.Value {
	if len(s.data) == 0 {
		return values.None
	}

	last := len(s.data) - 1
	value := s.data[last]
	s.data = s.data[:last]

	return value
}

func (s *Stack) MarshalJSON() ([]byte, error) {
	panic("not supported")
}

func (s *Stack) String() string {
	return "[Stack]"
}

func (s *Stack) Unwrap() interface{} {
	panic("not supported")
}

func (s *Stack) Hash() uint64 {
	panic("not supported")
}

func (s *Stack) Copy() core.Value {
	panic("not supported")
}
