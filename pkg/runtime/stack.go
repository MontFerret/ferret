package runtime

import (
	"github.com/MontFerret/ferret/pkg/runtime/core"
)

type Stack struct {
	values []core.Value
}

func NewStack(size int) *Stack {
	return &Stack{
		values: make([]core.Value, 0, size),
	}
}

func (s *Stack) Len() int {
	return len(s.values)
}

func (s *Stack) Peek() core.Value {
	return s.values[len(s.values)-1]
}

func (s *Stack) Push(value core.Value) {
	s.values = append(s.values, value)
}

func (s *Stack) Pop() core.Value {
	value := s.values[len(s.values)-1]
	s.values = s.values[:len(s.values)-1]
	return value
}

func (s *Stack) Get(index int) core.Value {
	return s.values[index]
}

func (s *Stack) Set(index int, value core.Value) {
	if index < len(s.values) {
		s.values[index] = value
	} else {
		s.values = append(s.values, value)
	}
}
