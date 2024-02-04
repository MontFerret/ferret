package runtime

import "github.com/MontFerret/ferret/pkg/runtime/core"

type Stack struct {
	operands  []core.Value
	variables []core.Value
}

func NewStack(operands, variables int) *Stack {
	return &Stack{
		make([]core.Value, 0, operands),
		make([]core.Value, 0, variables),
	}
}

func (s *Stack) Len() int {
	return len(s.operands)
}

func (s *Stack) Peek() core.Value {
	return s.operands[len(s.operands)-1]
}

func (s *Stack) Push(value core.Value) {
	s.operands = append(s.operands, value)
}

func (s *Stack) Pop() core.Value {
	value := s.operands[len(s.operands)-1]
	s.operands = s.operands[:len(s.operands)-1]
	return value
}

func (s *Stack) Get(index int) core.Value {
	return s.operands[index]
}

func (s *Stack) Set(index int, value core.Value) {
	s.operands[index] = value
}

func (s *Stack) GetVariable(index int) core.Value {
	return s.variables[index]
}

func (s *Stack) SetVariable(index int, value core.Value) {
	// TODO: Calculate in advance the number of variables
	if index >= len(s.variables) {
		s.variables = append(s.variables, value)
		return
	}

	s.variables[index] = value
}

func (s *Stack) PushVariable(value core.Value) {
	s.variables = append(s.variables, value)
}
