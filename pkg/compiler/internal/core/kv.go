package core

import (
	"github.com/MontFerret/ferret/v2/pkg/bytecode"
)

type KV struct {
	Key   bytecode.Operand
	Value bytecode.Operand
}

func NewKV(key bytecode.Operand, value bytecode.Operand) *KV {
	return &KV{key, value}
}
