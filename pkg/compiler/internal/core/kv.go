package core

import "github.com/MontFerret/ferret/v2/pkg/vm"

type KV struct {
	Key   vm.Operand
	Value vm.Operand
}

func NewKV(key vm.Operand, value vm.Operand) *KV {
	return &KV{key, value}
}
