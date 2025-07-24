package core

import (
	"github.com/MontFerret/ferret/pkg/runtime"
	"github.com/MontFerret/ferret/pkg/vm"
)

type AggregateSelector struct {
	name          runtime.String
	args          int
	funcName      runtime.String
	protectedCall bool
	register      vm.Operand
}

func NewAggregateSelector(name runtime.String, args int, funcName runtime.String, protectedCall bool, register vm.Operand) *AggregateSelector {
	return &AggregateSelector{
		name:          name,
		register:      register,
		args:          args,
		funcName:      funcName,
		protectedCall: protectedCall,
	}
}

func (s *AggregateSelector) Name() runtime.String {
	return s.name
}

func (s *AggregateSelector) Args() int {
	return s.args
}

func (s *AggregateSelector) FuncName() runtime.String {
	return s.funcName
}

func (s *AggregateSelector) ProtectedCall() bool {
	return s.protectedCall
}

func (s *AggregateSelector) Register() vm.Operand {
	return s.register
}
