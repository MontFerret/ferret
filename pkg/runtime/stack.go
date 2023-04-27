package runtime

import "github.com/MontFerret/ferret/pkg/runtime/core"

type Stack struct {
	values []core.Value
}

func NewStack(cap int) *Stack {
	return &Stack{make([]core.Value, 0, cap)}
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
